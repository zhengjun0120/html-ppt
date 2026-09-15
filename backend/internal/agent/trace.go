package agent

import (
	"encoding/json"
	"html-ppt/backend/internal/trace"
	"log"
	"sort"
	"time"

	"github.com/openai/openai-go/v3"
)

// 观测（trace）在这里落地。它是**旁路**：不影响发给模型的东西、不影响工具行为，
// 只是把"已经发生的事"记下来。因此本文里所有函数都遵守同一条纪律：
//
//	观测失败只打 warn，永不返回错误。
//
// 与 main.go 对视觉审查的 fail-open 是同一条原则，但这里的理由更硬：观测是
// 排查手段，为了它把用户正在进行的对话搞挂，等于用一个诊断工具制造故障。

// runTraceInfo 一个 run 的身份信息。
type runTraceInfo struct {
	SessionID   uint
	UserID      uint
	DeckID      string
	UserContent string
	ParentRunID string // ask_user 恢复时指回被暂停的那个 run
}

// openTraceRecorder 为本次 run 开一个记录器。
//
// 三条设计要点：
//   - **永远返回非 nil**（失败时退化成 trace.Discard）。调用方于是可以直接
//     `defer rec.Close()`、直接 Emit，一路上不必判空——观测因此对主流程零侵入。
//   - **实时流里剔除 messages**。它是文件里唯一的大字段（每轮重发整份消息，
//     开了全量上下文就是几 MB），而 serveAgentSSE 的通道只有 16 个槽，
//     推它会把聊天流拖住。观测页要看全文时按需拉 /events/:seq。
//   - **实时回调用 rec.Emit 而不是包级 Emit**：此时 ctx 还没绑定 recorder。
func (as *AgentService) openTraceRecorder(info runTraceInfo, emit func(StreamEvent) error) *trace.Recorder {
	if !as.TraceCfg.Enabled {
		return trace.Discard
	}
	if info.SessionID == 0 {
		// 没有会话号就没有目录归属，也过不了读侧的校验。理论上到不了这里
		//（会话在进循环前已经建好），保守跳过。
		log.Printf("[warn] trace: 会话 id 为空，本轮不记录")
		return trace.Discard
	}

	live := func(e trace.Event) {
		if emit == nil {
			return
		}
		e.Messages = nil
		raw, err := json.Marshal(e)
		if err != nil {
			log.Printf("[warn] trace: 实时事件序列化失败 err: %v", err)
			return
		}
		// 这里**吞掉** emit 的错误：它表示客户端断了（见 serveAgentSSE），
		// 而丢一条观测事件远好过因此中断整轮对话
		_ = emit(StreamEvent{Type: EventTypeTrace, Content: string(raw)})
	}

	runID := trace.NewRunID(time.Now())
	rec, err := trace.New(trace.Options{
		Dir:           as.TraceCfg.Dir,
		SessionID:     info.SessionID,
		RunID:         runID,
		UserID:        info.UserID,
		DeckID:        info.DeckID,
		ParentRunID:   info.ParentRunID,
		CaptureImages: as.TraceCfg.CaptureImages,
		MaxFieldBytes: as.TraceCfg.MaxFieldBytes,
		RetainRuns:    as.TraceCfg.RetainRuns,
		Live:          live,
	})
	if err != nil {
		log.Printf("[warn] trace: 开启观测失败，本轮不记录 err: %v", err)
		return trace.Discard
	}

	// 挂载了哪些工具要记下来：features 开关的效果在这里一眼可见。
	// （排查"模型为什么不用联网搜索"时，第一个要确认的就是这个工具到底有没有挂上）
	names := make([]string, 0, len(as.Exec))
	for name := range as.Exec {
		names = append(names, name)
	}
	sort.Strings(names)

	rec.Emit(trace.Event{
		Kind:        trace.KindRunStart,
		RunID:       runID,
		ParentRunID: info.ParentRunID,
		SessionID:   info.SessionID,
		UserID:      info.UserID,
		DeckID:      info.DeckID,
		UserContent: info.UserContent,
		Model:       as.ModelID,
		Tools:       names,
	})
	return rec
}

// usagePartFrom 把 SDK 的用量换算成观测口径。
//
// Reasoning 与 ImageTokens 单独拎出来，因为它们各自解释一类"看起来很奇怪"的现象：
// 前者是"模型只回了半句话/什么都没回"的常见原因（推理先吃 token，吃完就没有正文，
// 而这不报错）；后者是视觉审查的主要成本（一次审查要塞好几张 1244×700 的图）。
func usagePartFrom(u openai.CompletionUsage) trace.UsagePart {
	return trace.UsagePart{
		Prompt:      u.PromptTokens,
		Completion:  u.CompletionTokens,
		Total:       u.TotalTokens,
		Cached:      u.PromptTokensDetails.CachedTokens,
		Reasoning:   u.CompletionTokensDetails.ReasoningTokens,
		ImageTokens: u.PromptTokensDetails.ImageTokens,
		Calls:       1,
	}
}

// toolCallsOut 记录模型发起过哪些工具调用（含参数原文）。
func toolCallsOut(calls []openai.ChatCompletionMessageToolCallUnion) []trace.ToolCallOut {
	if len(calls) == 0 {
		return nil
	}
	out := make([]trace.ToolCallOut, 0, len(calls))
	for _, c := range calls {
		out = append(out, trace.ToolCallOut{
			ID:        c.ID,
			Name:      c.Function.Name,
			Arguments: c.Function.Arguments, // 原文：模型偶尔给出非法 JSON，那正是要看的东西
		})
	}
	return out
}
