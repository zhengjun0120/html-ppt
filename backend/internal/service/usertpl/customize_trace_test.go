package usertpl

// 定制对话观测与流式（customize_trace.go）的回归：
// 假 LLM = httptest SSE（openai-go 累加器吃的标准 chunk 序列），两轮脚本：
// 先 write_tokens 工具调用，再 finish 收尾。断言：
//   - 工具真的落盘（style.css 出现写入的 token）+ 记一条 chat 版本；
//   - trace 全事件在案：run_start(run_kind=customize) → llm_request（含全量
//     上下文）→ llm_response → tool_call → tool_result → usage → run_end(ok)；
//   - run 落在伪会话目录（customSessionID），归属字段完整；
//   - emit 非 nil 时 SSE 事件序列（tool_start/tool_done/done）齐全。

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"

	"html-ppt/backend/internal/store"
	"html-ppt/backend/internal/trace"
)

// chunkSSE 一条 SSE chunk 帧。choices 传 nil 时输出 usage 尾包。
func chunkSSE(model string, choices string, usage string) string {
	b := `data: {"id":"chatcmpl-x","object":"chat.completion.chunk","created":1700000000,"model":"` + model + `"`
	if choices == "" && usage == "" {
		return b + "}\n\n"
	}
	if choices != "" {
		b += `,"choices":[` + choices + `]`
	} else {
		b += `,"choices":[]`
	}
	if usage != "" {
		b += `,"usage":` + usage
	}
	return b + "}\n\n"
}

// newFakeLLM 造一个按脚本回放的两轮流式 LLM。
// script[0] = 第一轮 chunk 列表（tool 调用），script[1] = 第二轮（finish）。
func newFakeLLM(t *testing.T, rounds [][]string) (LLM, *atomic.Int32) {
	t.Helper()
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/chat/completions") {
			http.NotFound(w, r)
			return
		}
		n := int(calls.Add(1)) - 1
		if n >= len(rounds) {
			http.Error(w, "no script", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		for _, frame := range rounds[n] {
			_, _ = w.Write([]byte(frame))
		}
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
	t.Cleanup(srv.Close)
	client := openai.NewClient(option.WithAPIKey("test-key"), option.WithBaseURL(srv.URL))
	return LLM{Client: &client, Model: "fake-model"}, &calls
}

func toolRound(name, args string) []string {
	return []string{
		chunkSSE("fake-model", `{"index":0,"delta":{"role":"assistant","tool_calls":[{"index":0,"id":"call_1","type":"function","function":{"name":"`+name+`","arguments":""}}]},"finish_reason":null}`, ""),
		chunkSSE("fake-model", `{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":`+jsonStr(args)+`}}]},"finish_reason":null}`, ""),
		chunkSSE("fake-model", `{"index":0,"delta":{},"finish_reason":"tool_calls"}`, ""),
		chunkSSE("fake-model", "", `{"prompt_tokens":12,"completion_tokens":34,"total_tokens":46}`),
	}
}

func jsonStr(s string) string {
	raw, _ := json.Marshal(s)
	return string(raw)
}

func newTraceService(t *testing.T, suffix string) (*Service, *store.Store, string) {
	t.Helper()
	st, err := store.OpenMemory(t.Context())
	if err != nil {
		t.Fatalf("内存库: %v", err)
	}
	traceDir := t.TempDir()
	s := New(nil, st, t.TempDir(), "", "", nil, trace.Config{Enabled: true, Dir: traceDir})
	id := "ut-trace-" + suffix
	if err := st.DB.Create(&store.UserTemplate{ID: id, UserID: 9, BaseID: "tech-sharing", Status: "draft"}).Error; err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(s.Dir(id), 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"template.json": `{"id":"` + id + `","name":"t"}`,
		"index.html":    demoHTML("base"),
		"style.css":     ".tpl-x{--accent:#000}",
		"layouts.md":    "## cover\n骨架",
		"rules.md":      "# 规则",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(s.Dir(id), name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return s, st, id
}

func TestCustomizeTraceAndEmit(t *testing.T) {
	s, st, id := newTraceService(t, "full")
	llm, _ := newFakeLLM(t, [][]string{
		toolRound("write_tokens", `{"tokens":{"--accent":"#123456"}}`),
		toolRound("finish", `{"reply":"已把主色改成蓝色，请看封面页。"}`),
	})
	// 基线版先记下，Customize 的 dirty 记的 chat 版排在其后
	_ = s.recordVersionUT(id, OpFork, "起点")

	var events []CustEvent
	reply, err := s.Customize(t.Context(), 9, id, "把主色换成蓝色", nil, llm, func(ev CustEvent) error {
		events = append(events, ev)
		return nil
	})
	if err != nil {
		t.Fatalf("Customize: %v", err)
	}
	if !strings.Contains(reply, "主色改成蓝色") {
		t.Errorf("答复不对: %q", reply)
	}
	// 工具真的落盘
	css, _ := os.ReadFile(filepath.Join(s.Dir(id), "style.css"))
	if !strings.Contains(string(css), "#123456") {
		t.Errorf("token 未写入 style.css: %q", css)
	}
	versions, _ := s.ListUTVersions(9, id)
	if len(versions) != 2 || versions[0].Operation != OpChat {
		t.Errorf("应有基线+chat 两版: %+v", versions)
	}

	// —— SSE 事件序列 ——//
	var starts, dones int
	for _, ev := range events {
		switch ev.Type {
		case CustEvToolStart:
			starts++
			if ev.ToolName != "write_tokens" && ev.ToolName != "finish" {
				t.Errorf("意外工具: %s", ev.ToolName)
			}
		case CustEvToolDone:
			if ev.Content == "" {
				t.Errorf("tool_done 缺结果摘要")
			}
		case CustEvDone:
			dones++
			if !ev.Dirty || ev.Reply == "" {
				t.Errorf("done 应带 dirty+reply: %+v", ev)
			}
		}
	}
	if starts < 2 || dones != 1 {
		t.Errorf("SSE 事件不齐: tool_start=%d done=%d", starts, dones)
	}

	// —— trace 文件 ——//
	walk := filepath.Join(s.traceCfg.Dir, strconv.FormatUint(uint64(customSessionID(id)), 10))
	entries, err := os.ReadDir(walk)
	if err != nil || len(entries) == 0 {
		t.Fatalf("伪会话目录里没有 run: %s err=%v", walk, err)
	}
	var runFile string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".jsonl") {
			runFile = filepath.Join(walk, e.Name())
		}
	}
	raw, err := os.ReadFile(runFile)
	if err != nil {
		t.Fatal(err)
	}
	var kinds []string
	var start, end trace.Event
	for _, line := range strings.Split(strings.TrimSpace(string(raw)), "\n") {
		var ev trace.Event
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			t.Fatalf("坏事件行: %v", err)
		}
		kinds = append(kinds, ev.Kind)
		switch ev.Kind {
		case trace.KindRunStart:
			start = ev
		case trace.KindRunEnd:
			end = ev
		case trace.KindToolResult:
			if ev.Result == "" || ev.DurationMS < 0 {
				t.Errorf("tool_result 字段缺失: %+v", ev)
			}
		case trace.KindUsage:
			if ev.Usage == nil || ev.Usage.Total != 46 || ev.Usage.Calls != 1 {
				t.Errorf("用量口径不对: %+v", ev.Usage)
			}
		case trace.KindLLMRequest:
			if ev.Messages == nil || ev.MessageCount == 0 {
				t.Errorf("llm_request 应带全量上下文: seq=%d", ev.Seq)
			}
		case trace.KindLLMResponse:
			if len(ev.ToolCalls) == 0 {
				t.Errorf("llm_response 应带工具调用意图: seq=%d", ev.Seq)
			}
		}
	}
	if start.Kind != trace.KindRunStart || start.RunKind != "customize" || start.DeckID != id || start.UserID != 9 {
		t.Errorf("run_start 归属不对: %+v", start)
	}
	if start.UserContent != "把主色换成蓝色" {
		t.Errorf("run_start 缺用户消息: %q", start.UserContent)
	}
	if end.Status != trace.StatusOK || end.Summary == nil || end.Summary.ToolCalls != 2 {
		t.Errorf("run_end 汇总不对: %+v", end)
	}
	// 两轮 = 2×(request+response+usage) + 2×(call+result) + start + end。
	// usage 紧跟本轮回复合拍（与 agent.runLoop 的记录顺序一致），工具在其后。
	want := []string{"run_start",
		"llm_request", "llm_response", "usage", "tool_call", "tool_result",
		"llm_request", "llm_response", "usage", "tool_call", "tool_result",
		"run_end"}
	if strings.Join(kinds, ",") != strings.Join(want, ",") {
		t.Errorf("事件序列不对:\n got  %v\n want %v", kinds, want)
	}
	// 观测台归属校验依赖的第一行 user_id 必须在文件里
	if !strings.Contains(string(raw), `"user_id":9`) {
		t.Errorf("run_start 缺 user_id")
	}
	_ = st
}

// TestCustomizeLLMError 中途失败也要落 error + run_end(error)，
// 否则观测页那次 run 永远"运行中"。
func TestCustomizeLLMError(t *testing.T) {
	s, _, id := newTraceService(t, "err")
	// 第一轮就空流（无任何 chunk 直接结束）→ 累加器没有 finish_reason → 报错
	llm, _ := newFakeLLM(t, [][]string{{}})
	if _, err := s.Customize(t.Context(), 9, id, "改个色", nil, llm, nil); err == nil {
		t.Fatal("断流应报错")
	}
	walk := filepath.Join(s.traceCfg.Dir, strconv.FormatUint(uint64(customSessionID(id)), 10))
	entries, _ := os.ReadDir(walk)
	var sawError, sawEnd bool
	for _, e := range entries {
		raw, _ := os.ReadFile(filepath.Join(walk, e.Name()))
		for _, line := range strings.Split(strings.TrimSpace(string(raw)), "\n") {
			var ev trace.Event
			if json.Unmarshal([]byte(line), &ev) != nil {
				continue
			}
			if ev.Kind == trace.KindError {
				sawError = true
			}
			if ev.Kind == trace.KindRunEnd && ev.Status == trace.StatusError {
				sawEnd = true
			}
		}
	}
	if !sawError || !sawEnd {
		t.Errorf("失败 run 应落 error+run_end(error): error=%v end=%v", sawError, sawEnd)
	}
}

// TestEmitClientGone 客户端断开（emit 返回 ctx 错误）时循环中止，不再白烧轮次。
func TestEmitClientGone(t *testing.T) {
	s, _, id := newTraceService(t, "gone")
	llm, calls := newFakeLLM(t, [][]string{
		toolRound("write_tokens", `{"tokens":{"--accent":"#123456"}}`),
		toolRound("finish", `{"reply":"不该到达"}`),
	})
	gone := errors.New("client gone")
	_, err := s.Customize(t.Context(), 9, id, "改个色", nil, llm, func(ev CustEvent) error {
		if ev.Type == CustEvToolDone {
			return gone // 第一个工具执行完就"断开"
		}
		return nil
	})
	if !errors.Is(err, gone) {
		t.Fatalf("应把断开错误传回: %v", err)
	}
	if n := calls.Load(); n != 1 {
		t.Errorf("断开后不应再请求模型: %d 次", n)
	}
}

// 思考增量透传：推理模型的 reasoning_content 分片 → CustEvThink 事件（正文之外的
// 通道），最终答复不受影响。非推理模型没有该字段，一条 think 事件都不会有。
func TestCustomizeThinkStream(t *testing.T) {
	s, _, id := newTraceService(t, "think")
	llm, _ := newFakeLLM(t, [][]string{
		{
			chunkSSE("fake-model", `{"index":0,"delta":{"role":"assistant","reasoning_content":"先看看当前"},"finish_reason":null}`, ""),
			chunkSSE("fake-model", `{"index":0,"delta":{"reasoning_content":"主色再定对比。"},"finish_reason":null}`, ""),
			chunkSSE("fake-model", `{"index":0,"delta":{"content":"好的"},"finish_reason":null}`, ""),
			chunkSSE("fake-model", `{"index":0,"delta":{},"finish_reason":"stop"}`, ""),
			chunkSSE("fake-model", "", `{"prompt_tokens":5,"completion_tokens":9,"total_tokens":14}`),
		},
	})

	var thinks []string
	reply, err := s.Customize(t.Context(), 9, id, "把主色调一下", nil, llm, func(ev CustEvent) error {
		if ev.Type == CustEvThink {
			thinks = append(thinks, ev.Content)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Customize: %v", err)
	}
	if reply != "好的" {
		t.Errorf("答复 = %q", reply)
	}
	if got := strings.Join(thinks, ""); got != "先看看当前主色再定对比。" {
		t.Errorf("思考增量不对: %q", got)
	}
}
