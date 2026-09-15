package agent

// 这条协议不变式的回归测试：每条带 tool_calls 的 assistant 消息，后面必须紧跟齐
// 全部对应 tool_call_id 的 tool 消息，中间不能夹别的角色。
//
// 为什么值得单独一个文件来守：违反它的表现不是"这一轮报错"，而是**这个会话废了**——
// API 回 400，而坏形状的消息已经被持久化，之后每次请求（包括老老实实回答问题）
// 都是同一个 400。实测就是这么个错：
//
//	400 An assistant message with 'tool_calls' must be followed by tool messages
//	    responding to each 'tool_call_id'
//
// 所以这里不靠"看代码觉得对"，而是把不变式写成断言 —— 每个用例结束时都跑一遍
// assertToolCallsPaired。改坏这条不变式的任何改动都会在这里炸出来。

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
	"testing"

	"html-ppt/backend/internal/store"

	"github.com/openai/openai-go/v3"
)

// assertToolCallsPaired 断言消息数组满足协议不变式。
func assertToolCallsPaired(t *testing.T, messages []openai.ChatCompletionMessageParamUnion) {
	t.Helper()
	for i := 0; i < len(messages); i++ {
		ids := toolCallIDsOf(messages[i])
		if len(ids) == 0 {
			continue
		}
		// 后面必须紧跟一串 tool 消息，且把它们全部答齐
		answered := map[string]bool{}
		for j := i + 1; j < len(messages); j++ {
			id := toolMsgCallID(messages[j])
			if id == "" {
				break // 遇到别的角色就断：tool 应答必须紧挨着
			}
			answered[id] = true
		}
		for _, id := range ids {
			if !answered[id] {
				t.Fatalf("第 %d 条 assistant 声明的 tool_call %q 没有被紧跟的 tool 消息应答"+
					"（API 会回 400，而且这个形状会被持久化、把会话永久弄坏）", i, id)
			}
		}
	}
	// 反向：tool 消息不能是"无主的"（没有对应的 tool_calls 同样会被 API 拒）
	declared := map[string]bool{}
	for _, m := range messages {
		for _, id := range toolCallIDsOf(m) {
			declared[id] = true
		}
	}
	for i, m := range messages {
		if id := toolMsgCallID(m); id != "" && !declared[id] {
			t.Fatalf("第 %d 条 tool 消息回应的 %q 没有任何 tool_calls 声明它（API 同样会拒）", i, id)
		}
	}
}

func assistantWithCalls(content string, calls ...string) openai.ChatCompletionMessageParamUnion {
	// calls 传 "id:name:args" 形式，省得每处都拼一遍结构体
	msg := openai.AssistantMessage(content)
	for _, c := range calls {
		parts := strings.SplitN(c, ":", 3)
		msg.OfAssistant.ToolCalls = append(msg.OfAssistant.ToolCalls, openai.ChatCompletionMessageToolCallUnionParam{
			OfFunction: &openai.ChatCompletionMessageFunctionToolCallParam{
				ID:       parts[0],
				Function: openai.ChatCompletionMessageFunctionToolCallFunctionParam{Name: parts[1], Arguments: parts[2]},
			},
		})
	}
	return msg
}

// TestRepairDoesNotAlarmOnExpectedPause 只有**预期之外**的不成对才该打警告。
//
// 暂停态本来就有一条没应答的 ask_user 调用（正在等用户回答），修复会给它补一条
// 兜底说明——而 AnswerChat 紧接着就会把它换成用户的真实回答。所以这是正常流程，
// 不该报警。不守这条的代价是：每一次"提问 → 回答"都刷一条 [warn]，
// 于是真正需要看的"某条工具调用没应答"被淹在噪声里，没人会再注意到。
//
// 用捕获 log 输出来断言，是因为这里唯一可观察的差别就是"有没有报警"。
// log.SetOutput 是全局的，所以这条测试不能与同包其他测试并行（本包没有用 t.Parallel）。
func TestRepairDoesNotAlarmOnExpectedPause(t *testing.T) {
	var buf bytes.Buffer
	oldOut, oldFlags := log.Writer(), log.Flags()
	log.SetOutput(&buf)
	log.SetFlags(0)
	defer func() { log.SetOutput(oldOut); log.SetFlags(oldFlags) }()

	// 正常暂停态：结尾是一条等回答的 ask_user 调用
	paused := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("做个封面"),
		assistantWithCalls("", "call_ask:ask_user:{...}"),
	}
	out, fixed := repairToolCalls(paused, "call_ask")
	if fixed != 1 {
		t.Errorf("应补 1 条兜底（AnswerChat 随后会把它换成真实回答），实际 %d", fixed)
	}
	assertToolCallsPaired(t, out)
	if buf.Len() != 0 {
		t.Errorf("正常暂停态的兜底不该打警告，实际日志: %s", buf.String())
	}

	// 真异常：一条与提问无关的工具调用没留下结果
	buf.Reset()
	broken := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("改一下第 3 页"),
		assistantWithCalls("", "call_web:web_search:{...}"),
	}
	if _, n := repairToolCalls(broken, "call_ask"); n != 1 {
		t.Errorf("应补 1 条，实际 %d", n)
	}
	if buf.Len() == 0 {
		t.Error("与提问无关的工具调用不成对必须报警——它意味着会话被中断过，不是正常状态")
	}
}

// toolText 取一条 tool 消息的文本内容。Content 是个 union（字符串 或 内容分片数组），
// 测试里只关心文本，在这里摊平，省得每处断言都写一遍判断。
func toolText(m openai.ChatCompletionMessageParamUnion) string {
	if m.OfTool == nil {
		return ""
	}
	if s := m.OfTool.Content.OfString.Value; s != "" {
		return s
	}
	var b strings.Builder
	for _, p := range m.OfTool.Content.OfArrayOfContentParts {
		b.WriteString(p.Text)
	}
	return b.String()
}

func kindsOf(messages []openai.ChatCompletionMessageParamUnion) string {
	var out []string
	for _, m := range messages {
		switch {
		case m.OfSystem != nil:
			out = append(out, "system")
		case m.OfUser != nil:
			out = append(out, "user")
		case m.OfAssistant != nil:
			if n := len(m.OfAssistant.ToolCalls); n > 0 {
				out = append(out, "assistant+tool_calls")
			} else {
				out = append(out, "assistant")
			}
		case m.OfTool != nil:
			out = append(out, "tool("+m.OfTool.ToolCallID+")")
		default:
			out = append(out, "?")
		}
	}
	return strings.Join(out, " → ")
}

// TestRepairClosesDanglingAskUser 暂停态的形状：结尾是一条带 tool_calls 的 assistant
// 消息，而 ask_user 那个调用**故意**没有应答（正在等用户回答）。
//
// 这个形状不能直接发给模型，也不能在它后面追加用户消息——修复要做的就是补一条
// 说明把那个调用闭合。
func TestRepairClosesDanglingAskUser(t *testing.T) {
	msgs := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("你是助手"),
		openai.UserMessage("做个封面"),
		assistantWithCalls("我先确认一下", "call_ask:ask_user:`{\"questions\":[{\"question\":\"深色还是浅色?\"}]}`"),
	}

	out, fixed := repairToolCalls(msgs, "call_ask")
	if fixed != 1 {
		t.Errorf("应修复 1 处，实际 %d", fixed)
	}
	assertToolCallsPaired(t, out)
	if len(out) != 4 {
		t.Fatalf("应补出一条 tool 消息，实际 %s", kindsOf(out))
	}
	// 兜底说明不能像"用户真的答了什么"——否则模型会基于不存在的回答继续编
	got := out[3].OfTool.Content
	if !strings.Contains(toolText(out[3]), "用户没有回答这个提问") {
		t.Errorf("待答的那个调用应得到一句明确的'未作答'说明，实际 %q", got)
	}
	// 提问参数本身不能被改掉：它还要被页面拿回去重建卡片
	if args := askUserArgsOf(out, "call_ask"); !strings.Contains(args, "深色还是浅色") {
		t.Errorf("修复不该动 tool_calls 里的参数原文，实际 %q", args)
	}
}

// TestRepairMovesMisplacedAnswer 被历史 bug 弄坏的形状：用户消息夹在了
// tool_calls 与它的应答之间。这正是线上报 400 的那个形状。
//
// 修复必须把应答**移到**正确位置，而不是再补一条——同一个 tool_call_id 两条应答
// API 也会拒。
func TestRepairMovesMisplacedAnswer(t *testing.T) {
	msgs := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("做个封面"),
		assistantWithCalls("我先确认一下", "call_ask:ask_user:{...}"),
		openai.UserMessage("算了，你先帮我加一页"), // ← bug 就出在这条被追加上来
		openai.ToolMessage(`{"answers":[{"answer":"深色"}]}`, "call_ask"),
	}

	out, fixed := repairToolCalls(msgs, "")
	if fixed != 1 {
		t.Errorf("应记录 1 处修复，实际 %d", fixed)
	}
	assertToolCallsPaired(t, out)
	// 应答要挪到 assistant 后面，用户消息还在、顺序保持
	if kindsOf(out) != "user → assistant+tool_calls → tool(call_ask) → user" {
		t.Errorf("顺序不对: %s", kindsOf(out))
	}
	// 用户真实回答的内容要保留（不能被兜底说明替掉）
	if !strings.Contains(toolText(out[2]), "深色") {
		t.Errorf("已有的真实回答应被保留，实际 %q", out[2].OfTool.Content)
	}
	if len(out) != 4 {
		t.Errorf("不该多出消息，实际 %d 条", len(out))
	}
}

// TestRepairLeavesHealthyConversationAlone 正常的对话不该被动过一个字节。
// 修复逻辑过于激进的话，会把每一轮正常的工具往返都重写一遍——
// 那既浪费又危险（可能打乱模型看到的顺序，而看不出报错）。
func TestRepairLeavesHealthyConversationAlone(t *testing.T) {
	msgs := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("做个封面"),
		assistantWithCalls("", "c1:read_slide:{...}", "c2:list_slides:{}"),
		openai.ToolMessage(`{"html":"<section>…"}`, "c1"),
		openai.ToolMessage(`{"slides":[]}`, "c2"),
		openai.AssistantMessage("做好了"),
	}
	out, fixed := repairToolCalls(msgs, "")
	if fixed != 0 {
		t.Errorf("健康的对话不该被修，实际修了 %d 处", fixed)
	}
	if kindsOf(out) != kindsOf(msgs) {
		t.Errorf("健康的对话不该被重排:\n  前 %s\n  后 %s", kindsOf(msgs), kindsOf(out))
	}
	assertToolCallsPaired(t, out)
}

// TestRepairPartialToolCalls 一条 assistant 发了多个工具调用、只答了一个。
// 缺的那个要补，已答的那个不能被覆盖。
func TestRepairPartialToolCalls(t *testing.T) {
	msgs := []openai.ChatCompletionMessageParamUnion{
		assistantWithCalls("同时查两件事", "c1:read_slide:{...}", "c2:web_search:{...}"),
		openai.ToolMessage(`{"html":"<section>…"}`, "c1"),
	}
	out, fixed := repairToolCalls(msgs, "")
	if fixed != 1 {
		t.Errorf("应补 1 条，实际 %d", fixed)
	}
	assertToolCallsPaired(t, out)
	if kindsOf(out) != "assistant+tool_calls → tool(c1) → tool(c2)" {
		t.Errorf("顺序不对: %s", kindsOf(out))
	}
	if !strings.Contains(toolText(out[1]), "<section>") {
		t.Error("已答的那个不能被兜底说明覆盖")
	}
	if !strings.Contains(toolText(out[2]), "没有留下结果") {
		t.Errorf("缺的那个应得到'未执行'说明，实际 %q", out[2].OfTool.Content)
	}
}

// TestRepairDropsOrphanAndDuplicateTools 两种同样会 400 的脏数据：
// 无主的 tool 消息（没有 tool_calls 声明它）、同一个 id 的两条应答。
func TestRepairDropsOrphanAndDuplicateTools(t *testing.T) {
	msgs := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("hi"),
		openai.ToolMessage(`{"note":"孤儿"}`, "call_ghost"), // 无主
		assistantWithCalls("", "c1:read_slide:{}"),
		openai.ToolMessage(`{"html":"第一次"}`, "c1"),
		openai.ToolMessage(`{"html":"重复的第二次"}`, "c1"), // 同 id 重复
	}
	out, fixed := repairToolCalls(msgs, "")
	if fixed != 2 {
		t.Errorf("应记录 2 处（1 孤儿 + 1 重复），实际 %d", fixed)
	}
	assertToolCallsPaired(t, out)
	if kindsOf(out) != "user → assistant+tool_calls → tool(c1)" {
		t.Errorf("脏数据应被丢掉: %s", kindsOf(out))
	}
	if !strings.Contains(toolText(out[2]), "第一次") {
		t.Error("重复时应保留先出现的那条")
	}
}

// TestSetToolAnswerReplacesInsteadOfAppending 重复回答要"替换"而不是"追加"。
//
// 追加会让同一个 tool_call_id 出现两条应答，API 同样拒——而这条路径很常见：
// 上一次回答因 400 失败过，用户又点了一次提交。
func TestSetToolAnswerReplacesInsteadOfAppending(t *testing.T) {
	msgs := []openai.ChatCompletionMessageParamUnion{
		assistantWithCalls("", "c1:ask_user:{}"),
		openai.ToolMessage(`{"note":"未作答"}`, "c1"), // 修复流程补的兜底
	}
	out := setToolAnswer(msgs, "c1", `{"answers":[{"answer":"深色"}]}`)
	if len(out) != 2 {
		t.Fatalf("应替换而不是追加，实际 %d 条: %s", len(out), kindsOf(out))
	}
	assertToolCallsPaired(t, out)
	if !strings.Contains(toolText(out[1]), "深色") {
		t.Errorf("应答内容应被替换成最新的，实际 %q", out[1].OfTool.Content)
	}

	// 没有现成应答时（例如没经过修复流程）才追加
	appended := setToolAnswer([]openai.ChatCompletionMessageParamUnion{
		assistantWithCalls("", "c1:ask_user:{}"),
	}, "c1", `{"answers":[]}`)
	if len(appended) != 2 {
		t.Errorf("没有现成应答时应追加，实际 %d 条", len(appended))
	}
}

// TestGuardNewMessageBlocksWhilePending 闸门本身。
//
// 硬拦截而不是"当作提问作废、直接继续"：那需要替用户猜他是想答题还是想改需求，
// 猜错时模型会把"算了先加一页"当成对问题的回答——比多一次点击更让人困惑。
func TestGuardNewMessageBlocksWhilePending(t *testing.T) {
	pending := &store.ChatSession{ID: 7, PendingAsk: encodePendingAsk(pendingAsk{ToolCallID: "call_ask", RunID: "1700000000001-a"})}
	if err := guardNewMessage(pending); err != ErrPendingAsk {
		t.Errorf("暂停态下应被拦下（ErrPendingAsk），实际 %v", err)
	}

	// 没有暂停：放行
	if err := guardNewMessage(&store.ChatSession{ID: 7}); err != nil {
		t.Errorf("没有待答提问时不该拦，实际 %v", err)
	}
	// 暂停态字段存在但内容坏了：当作没有暂停放行。
	// 坏掉的暂停态不该让整个会话不可用——真发不出去时 API 会给明确报错
	if err := guardNewMessage(&store.ChatSession{ID: 7, PendingAsk: "{坏掉的 json"}); err != nil {
		t.Errorf("暂停态损坏时应放行（当作没有暂停），实际 %v", err)
	}
	if err := guardNewMessage(&store.ChatSession{ID: 7, PendingAsk: `{"run_id":"x"}`}); err != nil {
		t.Errorf("暂停态缺 tool_call_id 时应放行，实际 %v", err)
	}
}

// TestPendingAskRoundTrip 暂停态的编解码，以及"提问参数能被取回来"
// （页面刷新后重建卡片就靠这一步；取不回来的话那个会话永远答不上）。
func TestPendingAskRoundTrip(t *testing.T) {
	raw := `{"questions":[{"question":"深色还是浅色?","options":["深色","浅色"]}]}`
	sess := &store.ChatSession{
		ID:         3,
		PendingAsk: encodePendingAsk(pendingAsk{ToolCallID: "call_ask", RunID: "1700000000002-b"}),
	}
	got := parsePendingAsk(sess)
	if got.ToolCallID != "call_ask" || got.RunID != "1700000000002-b" {
		t.Errorf("暂停态往返丢了字段: %+v", got)
	}

	msgs := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage("做个封面"),
		assistantWithCalls("", "call_ask:ask_user:"+raw),
	}
	args := askUserArgsOf(msgs, got.ToolCallID)
	if args != raw {
		t.Errorf("应取回提问参数原文，实际 %q", args)
	}
	var parsed struct {
		Questions []struct {
			Question string   `json:"question"`
			Options  []string `json:"options"`
		} `json:"questions"`
	}
	if err := json.Unmarshal([]byte(args), &parsed); err != nil || len(parsed.Questions) != 1 {
		t.Fatalf("取回的参数应能解析成问题列表: err=%v", err)
	}
	if parsed.Questions[0].Question != "深色还是浅色?" {
		t.Errorf("问题内容不对: %+v", parsed.Questions[0])
	}
	if askUserArgsOf(msgs, "不存在的 id") != "" {
		t.Error("查不到的 tool_call_id 应返回空串")
	}
}
