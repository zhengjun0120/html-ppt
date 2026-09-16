package agent

import (
	"encoding/json"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"

	"html-ppt/backend/internal/store"
)

// assistantMsg 构造一条带（或不带）工具调用的 assistant 消息，形状和
// streamOnce 里 acc.Choices[0].Message.ToParam() 落库的东西一致。
func assistantMsg(text, callID, name, args string) openai.ChatCompletionMessageParamUnion {
	m := openai.ChatCompletionAssistantMessageParam{
		Content: openai.ChatCompletionAssistantMessageParamContentUnion{OfString: param.NewOpt(text)},
	}
	if callID != "" {
		m.ToolCalls = []openai.ChatCompletionMessageToolCallUnionParam{{
			OfFunction: &openai.ChatCompletionMessageFunctionToolCallParam{
				ID:       callID,
				Function: openai.ChatCompletionMessageFunctionToolCallFunctionParam{Name: name, Arguments: args},
			},
		}}
	}
	return openai.ChatCompletionMessageParamUnion{OfAssistant: &m}
}

// sampleConversation 一段标准对话：系统提示词 → 用户 → 助手(调工具) → 工具结果 → 助手收尾。
func sampleConversation() []openai.ChatCompletionMessageParamUnion {
	return []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("系统提示词"),
		openai.UserMessage("帮我做一份介绍 Go 的 PPT"),
		assistantMsg("好的，我先写入封面页", "call_1", "write_page", `{"deck_id":"deck-0001"}`),
		openai.ToolMessage("已写入第 1 页", "call_1"),
		assistantMsg("封面页完成了，还需要加几页吗？", "", "", ""),
	}
}

func TestProjectTranscript_FullConversation(t *testing.T) {
	got := projectTranscript(sampleConversation(), "", 0)
	if len(got) != 4 {
		t.Fatalf("应有 4 条（system 跳过），得到 %d 条: %+v", len(got), got)
	}
	wantRoles := []string{"user", "assistant", "tool", "assistant"}
	wantSeqs := []int64{1, 2, 3, 4}
	for i, m := range got {
		if m.Role != wantRoles[i] {
			t.Errorf("第 %d 条 role = %q, want %q", i, m.Role, wantRoles[i])
		}
		if m.Seq != wantSeqs[i] {
			t.Errorf("第 %d 条 seq = %d, want %d（seq 必须按数组下标计，跳过的 system 也要占位）", i, m.Seq, wantSeqs[i])
		}
	}
	if got[0].Content != "帮我做一份介绍 Go 的 PPT" {
		t.Errorf("用户消息内容丢失: %q", got[0].Content)
	}
	// 工具调用挂在 assistant 上，参数原文保留
	if len(got[1].ToolCalls) != 1 || got[1].ToolCalls[0].Name != "write_page" || got[1].ToolCalls[0].ID != "call_1" {
		t.Fatalf("assistant 的 tool_calls 投影不对: %+v", got[1].ToolCalls)
	}
	// 工具结果单独一条，ToolCallID 关联，ToolName 从前面的 assistant 声明解出来
	if got[2].ToolCallID != "call_1" || got[2].ToolName != "write_page" || got[2].Content != "已写入第 1 页" {
		t.Errorf("tool 消息投影不对: %+v", got[2])
	}
}

// 暂停中的会话：修复流程给"还没回答的 ask_user"补了一条占位应答，
// 回放里必须丢掉它——用户还没回答，不能假装他答了。
// 没有暂停时（会话中断过的历史遗留），同样的应答要保留，如实展示"没有结果"。
func TestProjectTranscript_PendingPlaceholder(t *testing.T) {
	msgs := append(sampleConversation(),
		assistantMsg("", "call_p", "ask_user", `{"questions":[{"question":"要几页？"}]}`),
		openai.ToolMessage(`{"note":"用户没有回答这个提问…"}`, "call_p"),
	)

	got := projectTranscript(msgs, "call_p", 0)
	if n := len(got); n != 5 {
		t.Fatalf("暂停态应丢掉占位应答（5 条），得到 %d 条", n)
	}
	// ask_user 的提问本身（assistant 的 tool_calls）必须还在，卡片靠它渲染
	last := got[len(got)-1]
	if last.Role != "assistant" || len(last.ToolCalls) != 1 || last.ToolCalls[0].Name != "ask_user" {
		t.Fatalf("待答提问丢了，前端无法重建卡片: %+v", last)
	}

	kept := projectTranscript(msgs, "", 0)
	if n := len(kept); n != 6 {
		t.Fatalf("非暂停态应保留全部（6 条），得到 %d 条", n)
	}
	if kept[5].Role != "tool" || kept[5].ToolCallID != "call_p" {
		t.Errorf("中断会话的占位说明应如实保留: %+v", kept[5])
	}
}

func TestProjectTranscript_AfterSeq(t *testing.T) {
	got := projectTranscript(sampleConversation(), "", 2)
	if len(got) != 2 {
		t.Fatalf("after_seq=2 应只剩 seq 3、4 两条，得到 %d 条", len(got))
	}
	if got[0].Seq != 3 || got[1].Seq != 4 {
		t.Errorf("增量拉取的 seq 不对: %d, %d", got[0].Seq, got[1].Seq)
	}
}

// 走一遍真实的落库/读回路径：persistSession 是 json.Marshal 整个数组，
// loadOrCreateSession 是 json.Unmarshal 回来——投影必须吃得上反序列化后的形状，
// 这里守住字段名假设（tool_call_id、function.name 之类的 JSON 键）不被悄悄改掉。
func TestProjectTranscript_RoundtripThroughDBJSON(t *testing.T) {
	msgs := sampleConversation()
	raw, err := json.Marshal(msgs)
	if err != nil {
		t.Fatalf("序列化: %v", err)
	}
	var back []openai.ChatCompletionMessageParamUnion
	if err := json.Unmarshal(raw, &back); err != nil {
		t.Fatalf("反序列化: %v", err)
	}

	got := projectTranscript(back, "", 0)
	if len(got) != 4 {
		t.Fatalf("回放应 4 条，得到 %d 条", len(got))
	}
	if got[2].ToolName != "write_page" || got[2].Content != "已写入第 1 页" {
		t.Errorf("roundtrip 后 tool 消息投影不对: %+v", got[2])
	}
	if got[3].Content != "封面页完成了，还需要加几页吗？" {
		t.Errorf("roundtrip 后 assistant 正文不对: %q", got[3].Content)
	}
}

func TestSessionSummary_PendingFlag(t *testing.T) {
	sess := &store.ChatSession{}
	sess.PendingAsk = `{"tool_call_id":"call_p","run_id":"run-1"}`
	if !sessionSummary(sess).Pending {
		t.Error("有暂停态的会话 Pending 应为 true")
	}

	sess.PendingAsk = ""
	if sessionSummary(sess).Pending {
		t.Error("没有暂停态的会话 Pending 应为 false")
	}
	// 坏掉的暂停态按"没有暂停"处理（parsePendingAsk 的约定），列表页不该被单条脏数据打挂
	sess.PendingAsk = "{not-json"
	if sessionSummary(sess).Pending {
		t.Error("暂停态损坏时应按 false 处理")
	}
}
