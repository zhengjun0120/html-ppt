package agent

// agent 侧插桩的回归测试。这一层守的是"事件到底有没有在正确的地方发出来"——
// 观测工具最坏的失败方式不是崩溃，而是**少记了一条关键事件**：
// 你照样能看到一份看起来很完整的 trace，然后在缺了工具参数的那一栏上得出错误结论。
//
// 所以这里逐条断言：参数有没有进事件、失败的调用有没有留痕、
// 关掉开关时是不是真的一点文件都不产生、以及"实时流不带 messages"这条设计
// 有没有被后续改动悄悄破掉。

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"html-ppt/backend/internal/trace"

	"github.com/openai/openai-go/v3"
)

// readRunEvents 读出某个 run 落盘的全部事件。
func readRunEvents(t *testing.T, root string, sessID uint, runID string) []trace.Event {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(root, itoaU(sessID), runID+".jsonl"))
	if err != nil {
		t.Fatalf("读 trace 失败: %v", err)
	}
	var out []trace.Event
	for _, line := range strings.Split(string(raw), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var e trace.Event
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			t.Fatalf("trace 行不是合法 JSON: %v", err)
		}
		out = append(out, e)
	}
	return out
}

func itoaU(u uint) string {
	b, _ := json.Marshal(u)
	return strings.Trim(string(b), `"`)
}

func toolCall(id, name, args string) openai.ChatCompletionMessageToolCallUnion {
	return openai.ChatCompletionMessageToolCallUnion{
		ID:       id,
		Type:     "function",
		Function: openai.ChatCompletionMessageFunctionToolCallFunction{Name: name, Arguments: args},
	}
}

type collector struct{ evs []StreamEvent }

func (c *collector) emit(e StreamEvent) error {
	c.evs = append(c.evs, e)
	return nil
}

// TestExecToolRecordsArgsAndResult 工具调用的**参数**必须进 trace。
//
// 这是本次改动最核心的一条：在此之前参数从来没有被任何通道外发过——
// SSE 的 tool_call 只带工具名和返回值，runRecorder 虽然拿到了参数却只用来给
// deck 历史拼一句"修改了 s3"。结果是"模型给 update_slide 传了什么"这个问题
// 在记录里完全无法回答，而这恰恰是排查改坏页面时的第一个问题。
func TestExecToolRecordsArgsAndResult(t *testing.T) {
	root := t.TempDir()
	rec, err := trace.New(trace.Options{Dir: root, SessionID: 1, RunID: "1700000000001-exec", UserID: 1})
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()

	c := &collector{}
	as := &AgentService{Exec: map[string]ToolFunc{
		"update_slide": func(ctx context.Context, args string) (string, error) {
			// 工具内部也要能发子过程事件，并且归属自动到位
			trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Sub: &trace.SubStep{Name: "inner", Stage: "x"}})
			return `{"slide_id":"s3","warning":"写死了颜色"}`, nil
		},
	}}

	ctx := trace.WithTurn(trace.With(context.Background(), rec), 2)
	result, execErr := as.execTool(ctx, toolCall("call_1", "update_slide", `{"deck_id":"deck-0001","slide_id":"s3"}`), c.emit, newRunRecorder(), nil)
	if execErr != nil {
		t.Fatalf("工具执行不该返回错误: %v", execErr)
	}
	if !strings.Contains(result, "写死了颜色") {
		t.Errorf("工具结果应原样返回给循环，实际 %q", result)
	}
	rec.Close()

	evs := readRunEvents(t, root, 1, "1700000000001-exec")
	var call, res, inner *trace.Event
	for i := range evs {
		switch {
		case evs[i].Kind == trace.KindToolCall:
			call = &evs[i]
		case evs[i].Kind == trace.KindToolResult:
			res = &evs[i]
		case evs[i].Kind == trace.KindSubStep:
			inner = &evs[i]
		}
	}
	if call == nil || res == nil {
		t.Fatalf("应各有一条 tool_call 与 tool_result，实际 %+v", evs)
	}
	if !strings.Contains(call.Args, `"slide_id":"s3"`) {
		t.Errorf("tool_call 必须带参数原文，实际 %q", call.Args)
	}
	if call.ToolName != "update_slide" || call.ToolCallID != "call_1" {
		t.Errorf("tool_call 归属不对: name=%q id=%q", call.ToolName, call.ToolCallID)
	}
	if call.Turn == nil || *call.Turn != 2 {
		t.Errorf("tool_call 应带上轮次 2，实际 %v", call.Turn)
	}
	if call.Error != "" {
		t.Errorf("成功的调用不该带 error，实际 %q", call.Error)
	}
	if !strings.Contains(res.Result, "写死了颜色") {
		t.Errorf("tool_result 应带返回值，实际 %q", res.Result)
	}
	// 工具内部的子过程必须归到这次工具调用下，页面才能把它折叠进对应卡片
	if inner == nil || inner.ToolName != "update_slide" || inner.ToolCallID != "call_1" {
		t.Errorf("工具内部的事件应继承工具归属，实际 %+v", inner)
	}

	// 既有的 SSE 契约不能变：成功仍然发 tool_call
	if len(c.evs) != 1 || c.evs[0].Type != EventTypeToolCall {
		t.Errorf("成功的工具调用仍应推一条 tool_call 事件，实际 %+v", c.evs)
	}
}

// TestExecToolRecordsFailures 失败的调用同样要留痕。
//
// 原来只记成功调用（runRecorder 那句注释"只记录成功的"），于是"模型连试三次
// 都被参数校验拒掉"这种最需要解释的现象，在记录里完全看不到——
// 你只看得到它最后放弃了，看不到它试过什么。
func TestExecToolRecordsFailures(t *testing.T) {
	root := t.TempDir()
	cases := []struct {
		name     string
		toolName string
		exec     map[string]ToolFunc
		wantErr  string
	}{
		{
			name:     "工具自己返回错误",
			toolName: "read_slide",
			exec: map[string]ToolFunc{
				"read_slide": func(context.Context, string) (string, error) {
					return "", errors.New("slide_id 不合法")
				},
			},
			wantErr: "slide_id 不合法",
		},
		{
			name:     "模型调了不存在的工具",
			toolName: "list_decks",
			exec:     map[string]ToolFunc{},
			wantErr:  "未知工具",
		},
	}
	for i, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			runID := "170000000000" + string(rune('1'+i)) + "-fail"
			rec, err := trace.New(trace.Options{Dir: root, SessionID: 1, RunID: runID, UserID: 1})
			if err != nil {
				t.Fatal(err)
			}
			coll := &collector{}
			as := &AgentService{Exec: c.exec}
			ctx := trace.With(context.Background(), rec)

			result, execErr := as.execTool(ctx, toolCall("c1", c.toolName, `{"x":1}`), coll.emit, newRunRecorder(), nil)
			rec.Close()
			if execErr != nil {
				t.Fatalf("execTool 不该把工具错误往上抛（循环要继续跑）: %v", execErr)
			}
			if !strings.Contains(result, c.wantErr) {
				t.Errorf("返回给模型的结果应包含失败原因 %q，实际 %q", c.wantErr, result)
			}

			evs := readRunEvents(t, root, 1, runID)
			var res *trace.Event
			for i := range evs {
				if evs[i].Kind == trace.KindToolResult {
					res = &evs[i]
				}
			}
			if res == nil {
				t.Fatalf("失败也要有 tool_result（否则记录里只剩成功的那几次）: %+v", evs)
			}
			if !strings.Contains(res.Error, c.wantErr) {
				t.Errorf("tool_result.Error 应写明失败原因，实际 %q", res.Error)
			}
			// 模型想调但没执行成的工具，参数尤其要看：它是"模型想干什么"的唯一线索
			var call *trace.Event
			for i := range evs {
				if evs[i].Kind == trace.KindToolCall {
					call = &evs[i]
				}
			}
			if call == nil || call.Args != `{"x":1}` {
				t.Errorf("失败调用的参数也必须记下来，实际 %+v", call)
			}
			// SSE 契约：失败仍然发 tool_error
			if len(coll.evs) != 1 || coll.evs[0].Type != EventTypeToolError {
				t.Errorf("失败的调用应推 tool_error，实际 %+v", coll.evs)
			}
		})
	}
}

// TestOpenTraceRecorderGatingAndLivePolicy 两条设计承诺：
//
//  1. **关掉开关时一个文件都不产生**。开关失效是最容易悄悄发生的事——
//     "以为关掉了，其实一直在写"，磁盘涨起来之前没有任何征兆。
//  2. **实时流里剔除 messages**。它是唯一的大字段（几 MB），而 SSE 通道只有 16 槽；
//     推它会把聊天流拖住。这条一旦被后续改动"顺手补全"就会发现不出来，
//     因为功能看起来更完整了。
func TestOpenTraceRecorderGatingAndLivePolicy(t *testing.T) {
	as := &AgentService{
		ModelID: "test-model",
		Exec:    map[string]ToolFunc{"write_deck": nil, "web_search": nil, "ask_user": nil},
	}

	// —— 关掉：什么都不该发生 ——
	rootOff := t.TempDir()
	as.TraceCfg = trace.Config{Enabled: false, Dir: rootOff, CaptureImages: true}
	c := &collector{}
	rec := as.openTraceRecorder(runTraceInfo{SessionID: 1, UserID: 1}, c.emit)
	if rec != trace.Discard {
		t.Error("关掉观测时应返回 Discard")
	}
	trace.Emit(trace.With(context.Background(), rec), trace.Event{Kind: trace.KindToolCall, Args: "x"})
	rec.Close()
	if entries, _ := os.ReadDir(rootOff); len(entries) != 0 {
		t.Errorf("关掉观测时不该产生任何文件，实际 %v", entries)
	}
	if len(c.evs) != 0 {
		t.Errorf("关掉观测时不该推实时事件，实际 %+v", c.evs)
	}
	if trace.Active(trace.With(context.Background(), rec)) {
		t.Error("Discard 不该报告为活跃（否则每轮都要白序列化一遍上下文）")
	}

	// —— 打开：run_start 要带工具清单，实时事件要剔除 messages ——
	rootOn := t.TempDir()
	as.TraceCfg = trace.Config{Enabled: true, Dir: rootOn, CaptureImages: true, RetainRuns: 20}
	c = &collector{}
	rec = as.openTraceRecorder(runTraceInfo{SessionID: 7, UserID: 42, DeckID: "deck-0001", UserContent: "做个封面"}, c.emit)
	if rec == trace.Discard {
		t.Fatal("打开观测时不该返回 Discard")
	}
	ctx := trace.With(context.Background(), rec)

	// 模拟循环发一条带全量上下文的事件
	trace.Emit(ctx, trace.Event{
		Kind:     trace.KindLLMRequest,
		Messages: json.RawMessage(`[{"role":"user","content":"做个封面"}]`),
		Bytes:    38,
	})
	rec.Close()

	// run_start 落盘了，且带上了本次挂载的工具清单
	//（排查"模型为什么不用联网搜索"时，第一个要确认的就是它到底有没有挂上）
	runID := rec.RunID()
	evs := readRunEvents(t, rootOn, 7, runID)
	if len(evs) < 2 {
		t.Fatalf("应有 run_start 与 llm_request 两条，实际 %d 条", len(evs))
	}
	start := evs[0]
	if start.Kind != trace.KindRunStart || start.SessionID != 7 || start.UserID != 42 || start.DeckID != "deck-0001" {
		t.Errorf("run_start 的元信息不对: %+v", start)
	}
	if len(start.Tools) != 3 {
		t.Errorf("run_start 应带上本次挂载的工具清单，实际 %v", start.Tools)
	}

	// 落盘的那条要有 messages 全文
	if !strings.Contains(string(evs[1].Messages), "做个封面") {
		t.Errorf("落盘的事件应保留 messages 全文，实际 %s", evs[1].Messages)
	}

	// 实时推的那条必须没有 messages，但要有"多大"的元信息
	var live *StreamEvent
	for i := range c.evs {
		if c.evs[i].Type == EventTypeTrace && strings.Contains(c.evs[i].Content, trace.KindLLMRequest) {
			live = &c.evs[i]
		}
	}
	if live == nil {
		t.Fatal("实时流里应有对应的事件")
	}
	var got trace.Event
	if err := json.Unmarshal([]byte(live.Content), &got); err != nil {
		t.Fatalf("实时事件的 content 应是一条可解析的 trace 事件: %v", err)
	}
	if got.Messages != nil {
		t.Error("实时流里绝不能带 messages：几 MB 一条会把 16 槽的 SSE 通道拖住")
	}
	if got.Bytes != 38 {
		t.Errorf("剔除 messages 后要留下字节数，页面靠它表达\"上下文多大\"，实际 %d", got.Bytes)
	}
}

// 每页写入/修改次数闸：同一页第 4 次 update_slide 不执行、给自愈指令；
// 成功调用按 pageKey(slide_id)="p<页码>" 计数（失败不计）——与 write_pages
// 的整页重写共用同一张预算表。
func TestExecToolPerSlideUpdateCap(t *testing.T) {
	c := &collector{}
	calls := 0
	as := &AgentService{Exec: map[string]ToolFunc{
		"update_slide": func(ctx context.Context, args string) (string, error) {
			calls++
			return `{"ok":true}`, nil
		},
	}}

	// 已达上限：不执行、不算错误、给替代路径
	scope := &runScope{exec: as.Exec, updates: map[string]int{"p3": maxUpdatesPerSlide}}
	result, execErr := as.execTool(context.Background(),
		toolCall("c1", "update_slide", `{"deck_id":"d","slide_id":"s3"}`), c.emit, newRunRecorder(), scope)
	if execErr != nil {
		t.Fatalf("次数闸是策略结果，不该返回错误: %v", execErr)
	}
	if calls != 0 {
		t.Fatalf("已达上限的第 4 次调用不应执行")
	}
	for _, want := range []string{"s3", "3 次", "delete_slide"} {
		if !strings.Contains(result, want) {
			t.Errorf("拦截信息缺 %q: %q", want, result)
		}
	}

	// 未达上限：执行并按页计数；失败不计
	scope2 := &runScope{exec: as.Exec}
	ctx := context.Background()
	as.execTool(ctx, toolCall("c2", "update_slide", `{"slide_id":"s5"}`), c.emit, newRunRecorder(), scope2)
	as.execTool(ctx, toolCall("c3", "update_slide", `{"slide_id":"s5"}`), c.emit, newRunRecorder(), scope2)
	as.execTool(ctx, toolCall("c4", "update_slide", `{"slide_id":"s6"}`), c.emit, newRunRecorder(), scope2)
	if scope2.updates["p5"] != 2 || scope2.updates["p6"] != 1 {
		t.Errorf("计数错误: %+v", scope2.updates)
	}
	if calls != 3 {
		t.Errorf("三次调用都应执行，实际 %d", calls)
	}
}

// TestExecToolWritePagesRewriteCap write_pages 的整页重写与 update_slide 共用每页预算：
// 超预算的页从批里剔除（其余页照常落盘、以单页错误并进返回）；被质量闸拒收的页
// 不占预算——deck-0061 的修复循环全走 write_pages，p3/p4 各重写 3 次绕开了旧闸门。
func TestExecToolWritePagesRewriteCap(t *testing.T) {
	c := &collector{}
	var gotArgs string
	resp := `{"results":[{"no":4,"ok":true}],"written":1,"slides":4}`
	as := &AgentService{Exec: map[string]ToolFunc{
		"write_pages": func(ctx context.Context, args string) (string, error) {
			gotArgs = args
			return resp, nil
		},
	}}
	ctx := context.Background()

	// p3 已达上限：批里 3、4 两页，只有 4 该执行；拒因带页号与替代路径
	scope := &runScope{exec: as.Exec, updates: map[string]int{"p3": maxUpdatesPerSlide}}
	result, execErr := as.execTool(ctx, toolCall("c1", "write_pages",
		`{"deck_id":"d","pages":[{"no":3,"layout":"a","html":"<section></section>"},{"no":4,"layout":"b","html":"<section></section>"}]}`),
		c.emit, newRunRecorder(), scope)
	if execErr != nil {
		t.Fatalf("次数闸是策略结果，不该返回错误: %v", execErr)
	}
	if !strings.Contains(gotArgs, `"no":4`) || strings.Contains(gotArgs, `"no":3`) {
		t.Errorf("只有第 4 页该交给工具，实际参数 %s", gotArgs)
	}
	if !strings.Contains(result, `"no":3`) || !strings.Contains(result, "delete_slide") {
		t.Errorf("被拒页要以单页错误并进返回: %s", result)
	}
	if !strings.Contains(result, `"no":4`) {
		t.Errorf("好页的结果不能被合并弄丢: %s", result)
	}
	if scope.updates["p4"] != 1 || scope.updates["p3"] != maxUpdatesPerSlide {
		t.Errorf("只有落盘成功的页计数: %+v", scope.updates)
	}

	// 整批都超预算：不执行、给替代路径
	scope2 := &runScope{exec: as.Exec, updates: map[string]int{"p7": maxUpdatesPerSlide}}
	gotArgs = ""
	result, execErr = as.execTool(ctx, toolCall("c2", "write_pages",
		`{"deck_id":"d","pages":[{"no":7,"layout":"a","html":"<section></section>"}]}`),
		c.emit, newRunRecorder(), scope2)
	if execErr != nil || gotArgs != "" {
		t.Fatalf("整批超预算不应执行: err=%v args=%s", execErr, gotArgs)
	}
	if !strings.Contains(result, "delete_slide") || !strings.Contains(result, "7") {
		t.Errorf("拦截信息缺页号与替代路径: %s", result)
	}

	// 被质量闸拒收的页（ok:false）不占预算
	resp = `{"results":[{"no":5,"ok":false,"error":"密度下限"}],"written":0}`
	scope3 := &runScope{exec: as.Exec}
	as.execTool(ctx, toolCall("c3", "write_pages",
		`{"deck_id":"d","pages":[{"no":5,"layout":"a","html":"<section></section>"}]}`),
		c.emit, newRunRecorder(), scope3)
	if n := scope3.updates["p5"]; n != 0 {
		t.Errorf("落盘失败的页不该计数: %+v", scope3.updates)
	}
}
