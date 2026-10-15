package usertpl

// 定制对话的观测与流式（2026-09-28）：
//   - 观测：customizeLoop 的每个动作（请求/回复/工具/用量）落 trace 事件，
//     观测台与文稿 run 同页可见（run_kind=customize 区分）。
//   - 流式：LLM 调用改为流式，工具气泡与生成进度以 CustEvent 推给 SSE 客户端，
//     定制页聊天不再"转圈干等"。
//
// 与 deck 侧的关系：SDK 用法（累加器/流分片）照搬 agent.streamOnce 的成熟形状；
// 事件类型与数据模型分开定义——定制聊天没有大纲/页面/闸门管线，照搬
// agent.StreamEvent 只会让前端对着一堆永远用不到的字段做防御。

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"strings"
	"time"

	"github.com/openai/openai-go/v3"

	"html-ppt/backend/internal/chatimg"
	"html-ppt/backend/internal/store"
	"html-ppt/backend/internal/trace"
)

// CustEvent 定制对话推给前端的事件（SSE 帧 data 体）。
type CustEvent struct {
	Type       string `json:"type"`
	Content    string `json:"content,omitempty"` // delta 文本 / tool_done 结果摘要 / error 消息
	ToolName   string `json:"tool_name,omitempty"`
	ToolCallID string `json:"tool_call_id,omitempty"`
	ToolIndex  int64  `json:"tool_index,omitempty"`
	Bytes      int64  `json:"bytes,omitempty"` // tool_progress：该工具已生成的参数字节数
	Dirty      bool   `json:"dirty,omitempty"` // done：本轮是否有文件写入（前端据此刷新预览/历史）
	Reply      string `json:"reply,omitempty"` // done：最终答复全文
}

const (
	CustEvToolStart    = "tool_start"    // 模型开始生成某个工具调用（id+name 就绪）
	CustEvToolProgress = "tool_progress" // 大参数工具（write_style/write_demo）生成中的字节数心跳
	CustEvToolDone     = "tool_done"     // 工具执行完毕（Content = 结果摘要）
	CustEvDelta        = "delta"         // 模型正文增量（不经 finish 直接回话的收尾轮）
	CustEvThink        = "think"         // 模型思考增量（reasoning_content，非推理模型不会有）
	CustEvDone         = "done"
	CustEvError        = "error"
)

// customSessionID 定制 run 的伪会话号。trace 按数字 session_id 分目录，文稿那边
// 是真实自增 id（从 1 起）；定制没有会话表，用模板 id 哈希 + 1e9 起步的专用段：
// 既不与真实会话撞号，也让同一个模板的定制 run 天然聚在一个目录——
// 按 run 数裁剪（RetainRuns）因此等于每模板各自的保留配额。
func customSessionID(templateID string) uint {
	h := fnv.New32a()
	_, _ = h.Write([]byte("ut:" + templateID))
	return 1_000_000_000 + uint(h.Sum32())
}

// custToolNames 进 run_start.tools：本次挂了哪些工具，观测页一眼可见。
var custToolNames = []string{"finish", "set_meta", "write_demo", "write_style", "write_tokens"}

// openCustRecorder 开一个定制 run 的观测文件。fail-open 与 agent 侧同一条原则：
// 观测是增强不是故障源，任何失败只打 warn，定制对话照跑。
// 用户附图构造成 ImageEvent 挂在 run_start 上：Bytes 由 recorder 落到 <run>/img/，
// JSONL 只留相对路径（capture 关闭时只记张数，与视觉截图同一开关语义）。
func (s *Service) openCustRecorder(userID uint, row *store.UserTemplate, message string, imgs []chatimg.Image, model string) *trace.Recorder {
	if !s.traceCfg.Enabled || s.traceCfg.Dir == "" {
		return s.ledgerRecorder(userID, row, model)
	}
	sessID := customSessionID(row.ID)
	rec, err := trace.New(trace.Options{
		Dir:           s.traceCfg.Dir,
		SessionID:     sessID,
		RunID:         trace.NewRunID(time.Now()),
		UserID:        userID,
		DeckID:        row.ID,
		MaxFieldBytes: s.traceCfg.MaxFieldBytes,
		RetainRuns:    s.traceCfg.RetainRuns,
		CaptureImages: s.traceCfg.CaptureImages, // 用户附图与视觉截图同一开关
		Model:         model,
		UsageSink:     s.traceCfg.UsageSink,
	})
	if err != nil {
		log.Printf("[warn] usertpl: 开启定制观测失败，本轮不记录 err: %v", err)
		return s.ledgerRecorder(userID, row, model)
	}
	ev := trace.Event{
		Kind: trace.KindRunStart, RunKind: "customize",
		RunID: rec.RunID(), SessionID: sessID, UserID: userID,
		DeckID: row.ID, UserContent: message, Model: model, Tools: custToolNames,
	}
	for i, im := range imgs {
		ev.Images = append(ev.Images, trace.ImageEvent{
			Name:  fmt.Sprintf("u%03d.png", i+1),
			Label: fmt.Sprintf("用户附图 %d", i+1),
			Bytes: im.Bytes,
		})
	}
	rec.Emit(ev)
	return rec
}

// ledgerRecorder 观测关闭/落盘失败时的账本-only recorder：定制对话的用量照记
// （账本是业务数据），JSONL 不落、无任何序列化开销。
func (s *Service) ledgerRecorder(userID uint, row *store.UserTemplate, model string) *trace.Recorder {
	return trace.LedgerOnly(trace.Options{
		SessionID: customSessionID(row.ID),
		RunID:     trace.NewRunID(time.Now()),
		UserID:    userID,
		DeckID:    row.ID,
		Model:     model,
		UsageSink: s.traceCfg.UsageSink,
	})
}

// custStream 流式跑一轮定制 LLM 调用（请求上下文归属由 ctx 携带）。
// 返回累积完整的消息、finish_reason 与用量；emit 可为 nil（同步端点），
// nil 时不推任何 CustEvent。
//
// **不设 MaxTokens**（2026-09-28 用户拍板）：上限交给模型/服务商的默认值，与
// agent 主循环（setChatOpts）同款。起因是实测坑：推理模型（deepseek-flash）会
// 先输出 reasoning_content，显式小上限被推理吃满后 finish=length、零正文零工具
// （4000 时必现）。空回复兜底重试保留——服务商自己的默认上限仍可能截断推理。
//
// tool_progress 的心跳节奏：参数增量每累计 4KB 推一次累计字节数——write_style
// 实测量级 20-60KB，一轮生成 30-60s，没有心跳的话聊天页在最重要的那段时间里
// 仍然是一片空白；推全量增量则是把 60KB CSS 刷进聊天框，两头都不要。
func (s *Service) custStream(ctx context.Context, sess *customizeSession, llm LLM, emit func(CustEvent) error) (openai.ChatCompletionMessage, string, openai.CompletionUsage, error) {
	// 观测：这一轮实际发给模型的东西（含 style.css 注入的全文）。只记工具输入输出
	// 回答不了"模型为什么这么改"，上下文才是证据；MB 级序列化只在观测开启时做。
	// 附图的 data URL 在这里脱敏——原图已落 <run>/img/，base64 内联进 JSONL 只会
	// 把逐行拉取的观测文件撑爆。
	if trace.Active(ctx) {
		raw, err := json.Marshal(sess.messages)
		if err != nil {
			log.Printf("[warn] usertpl: 序列化定制上下文失败 err: %v", err)
			raw = json.RawMessage(`null`)
		}
		raw = chatimg.RedactDataURLs(raw)
		trace.Emit(ctx, trace.Event{
			Kind: trace.KindLLMRequest, Model: llm.Model,
			Messages: raw, MessageCount: len(sess.messages), Bytes: len(raw),
		})
	}

	stream := llm.Client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:    openai.ChatModel(llm.Model),
		Messages: sess.messages,
		Tools:    customizeTools,
		ToolChoice: openai.ChatCompletionToolChoiceOptionUnionParam{
			OfAuto: openai.String("auto"),
		},
		Temperature:   openai.Float(0.4),
		StreamOptions: openai.ChatCompletionStreamOptionsParam{IncludeUsage: openai.Bool(true)},
	})
	acc := openai.ChatCompletionAccumulator{}
	totalArgs := map[int64]int64{} // tool index → 该工具已生成的参数字节数
	beaten := map[int64]int64{}    // tool index → 上次心跳时累计到的字节数
	for stream.Next() {
		chunk := stream.Current()
		if !acc.AddChunk(chunk) {
			continue
		}
		// 尾分片（只有 usage、无 choices）跳过防 panic
		if len(chunk.Choices) == 0 {
			continue
		}
		delta := chunk.Choices[0].Delta
		for _, tc := range delta.ToolCalls {
			idx := tc.Index
			if idx < 0 {
				idx = 0 // 个别网关对单个工具调用用 -1，SDK 累加器内部归到 0
			}
			if tc.ID != "" && emit != nil {
				// 只有首个分片带 id 和完整 name
				if emitErr := emit(CustEvent{Type: CustEvToolStart, ToolIndex: idx, ToolCallID: tc.ID, ToolName: tc.Function.Name}); emitErr != nil {
					stream.Close()
					return openai.ChatCompletionMessage{}, "", openai.CompletionUsage{}, fmt.Errorf("事件推送失败: %w", emitErr)
				}
			}
			if n := int64(len(tc.Function.Arguments)); n > 0 {
				totalArgs[idx] += n
				if emit != nil && totalArgs[idx]-beaten[idx] >= 4096 {
					beaten[idx] = totalArgs[idx]
					if emitErr := emit(CustEvent{Type: CustEvToolProgress, ToolIndex: idx, Bytes: totalArgs[idx]}); emitErr != nil {
						stream.Close()
						return openai.ChatCompletionMessage{}, "", openai.CompletionUsage{}, fmt.Errorf("事件推送失败: %w", emitErr)
					}
				}
			}
		}
		if delta.Content != "" && emit != nil {
			if emitErr := emit(CustEvent{Type: CustEvDelta, Content: delta.Content}); emitErr != nil {
				stream.Close()
				return openai.ChatCompletionMessage{}, "", openai.CompletionUsage{}, fmt.Errorf("事件推送失败: %w", emitErr)
			}
		}
		// 思考增量（deepseek 系推理模型的 reasoning_content，与文稿侧 EventTypeThink
		// 同一提取方式：SDK 不认识这个字段，走原始 JSON ExtraFields）
		if emit != nil {
			if f, ok := delta.JSON.ExtraFields["reasoning_content"]; ok {
				var reasoning string
				if err := json.Unmarshal([]byte(f.Raw()), &reasoning); err == nil && reasoning != "" {
					if emitErr := emit(CustEvent{Type: CustEvThink, Content: reasoning}); emitErr != nil {
						stream.Close()
						return openai.ChatCompletionMessage{}, "", openai.CompletionUsage{}, fmt.Errorf("事件推送失败: %w", emitErr)
					}
				}
			}
		}
	}
	if err := stream.Err(); err != nil {
		stream.Close()
		trace.Emit(ctx, trace.Event{Kind: trace.KindLLMResponse, Error: "流式输出错误: " + err.Error()})
		return openai.ChatCompletionMessage{}, "", openai.CompletionUsage{}, fmt.Errorf("流式输出错误: %w", err)
	}
	stream.Close()

	if len(acc.Choices) == 0 || acc.Choices[0].FinishReason == "" {
		trace.Emit(ctx, trace.Event{Kind: trace.KindLLMResponse, Error: "流式响应不完整（没有 finish_reason）"})
		return openai.ChatCompletionMessage{}, "", openai.CompletionUsage{}, fmt.Errorf("流式响应不完整，请重试")
	}

	finish := string(acc.Choices[0].FinishReason)
	msg := acc.Choices[0].Message
	// 观测：模型回了什么（文本/工具调用意图/finish_reason），与 llm_request
	// 配对看，才能回答"同样的上下文，模型这次为什么改主意"；用量紧随其后落一条。
	trace.Emit(ctx, trace.Event{
		Kind:         trace.KindLLMResponse,
		FinishReason: finish,
		Content:      msg.Content,
		ToolCalls:    custToolCallsOut(msg.ToolCalls),
	})
	trace.Usage(ctx, trace.CompMain, custUsagePart(acc.Usage))
	return msg, finish, acc.Usage, nil
}

// custUsagePart SDK 用量 → 观测口径。与 agent.usagePartFrom 同一张表，就地复制：
// 为一个换算函数跨包导出 agent 内部件不值得。
func custUsagePart(u openai.CompletionUsage) trace.UsagePart {
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

// custToolCallsOut 模型发起的工具调用意图（参数保留原文——非法 JSON 本身就是要观测的失败）。
func custToolCallsOut(calls []openai.ChatCompletionMessageToolCallUnion) []trace.ToolCallOut {
	if len(calls) == 0 {
		return nil
	}
	out := make([]trace.ToolCallOut, 0, len(calls))
	for _, c := range calls {
		out = append(out, trace.ToolCallOut{ID: c.ID, Name: c.Function.Name, Arguments: c.Function.Arguments})
	}
	return out
}

// toolBrief tool_done 气泡的摘要：聊天页只需要"成了/被拒了 + 一句原因"，
// 完整结果模型已经拿到、观测台看得到原文。
func toolBrief(result string) string {
	r := []rune(strings.TrimSpace(result))
	if len(r) > 100 {
		return string(r[:100]) + "…"
	}
	return string(r)
}
