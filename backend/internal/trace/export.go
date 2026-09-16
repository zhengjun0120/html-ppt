package trace

// 导出（ExportRun）：把一个 run 打包成**主要给 AI 看**的自包含 JSON。
// 与磁盘 jsonl 的差异全部服务于这个受众：
//
//  1. messages 折叠成增量。一次 run 的上下文只增不减，每轮 llm_request 却带着
//     全量——原样导出体积是 O(n²)，而重复的内容对读它的 AI 没有信息量。
//     只有上一条请求的消息序列是本条的**逐字节前缀**时才折叠；任何对不上
//     （会话被修复过、形状异常）都保留全量——宁大勿缺。
//  2. 截图默认不内联（?images=1 才内联）。base64 让文件膨胀一个量级，
//     而文本形态的读者看不见像素；名字/标签/相对路径保留当线索，
//     需要看图排查（视觉审查类问题）时再带图重新导出。
//  3. timeline：每个事件一行的速览。读它的 AI 先扫这里定位可疑环节，
//     再按 seq 号去 events 里取原文。
//
// 保真不变的部分：工具参数、工具返回、模型回复、错误原因一律**原文**——
// 模型偶尔给出非法 JSON，那正是要排查的证据，解析或"修复"它等于销毁证据。

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// ExportFormat / ExportVersion 标识导出文件自身的格式。它与 jsonl 的事件格式是
// 两回事：jsonl 一行一个事件、截图在旁边的目录里；导出是自包含的单文件。
// 将来导出结构演进时，读侧靠 Version 分支处理。
const (
	ExportFormat  = "html-ppt-trace-export"
	ExportVersion = 2 // v2：messages 增量化、截图默认不内联、新增 timeline（v1 从未发布过，但数字照实走）
)

// Export 一个 run 的导出包。
type Export struct {
	Format     string    `json:"format"`
	Version    int       `json:"version"`
	ExportedAt time.Time `json:"exported_at"`
	Run        RunMeta   `json:"run"` // 与列表同口径的摘要（含 run_end 汇总；还在跑则 status=running）

	// Timeline 每个事件一行的速览，与 Events 按 seq 一一对应。它是"先扫这里、
	// 再按 seq 取全文"的索引层：行内的参数/返回/正文都是截断的摘录，完整原文
	// 一律在 events 里，timeline 不承载证据。
	Timeline []string `json:"timeline"`

	// Events 全量事件，按 seq 升序。两处对 AI 无用的东西被处理掉：
	// user_id 清空（归属已在 run 层校验过）、连续 llm_request 的重复上下文折叠成增量
	// （见 messages 字段里的自描述 note）。
	Events []Event `json:"events"`

	// Images 截图的处理方式。默认不内联；hint 只在不内联时给出。
	Images struct {
		Inlined bool   `json:"inlined"`
		Count   int    `json:"count"`
		Hint    string `json:"hint,omitempty"`
	} `json:"images"`
}

// ExportRun 导出一个 run。归属校验与读接口同一条路（runFileFor），
// "不是你的"和"不存在"一样返回 ErrNotFound。
//
// 全量事件一次性进内存——导出是用户点出来的显式动作，这个代价可以接受；
// 但它不该被任何自动流程（轮询、定时备份）走到。
func (s *Store) ExportRun(userID uint, sessionID uint, runID string, inlineImages bool) (*Export, error) {
	path, err := s.runFileFor(userID, sessionID, runID)
	if err != nil {
		return nil, err
	}

	// 摘要：首行 run_start + 尾行 run_end，与列表共用 runMetaFrom，口径不漂移
	headLine, err := readFirstLine(path)
	if err != nil {
		return nil, ErrNotFound
	}
	var head Event
	if err := json.Unmarshal(headLine, &head); err != nil {
		return nil, ErrNotFound
	}
	var end *Event
	if tailLine, terr := readLastLine(path); terr == nil {
		var te Event
		if json.Unmarshal(tailLine, &te) == nil && te.Kind == KindRunEnd {
			end = &te
		}
	}
	var size int64
	if st, serr := os.Stat(path); serr == nil {
		size = st.Size()
	}

	res, err := s.ReadEvents(userID, sessionID, runID, ReadOptions{IncludeMessages: true})
	if err != nil {
		return nil, err
	}

	// 截图。缺图（CaptureImages 关过、文件被清）不失败：URL 保持原样当线索，
	// "导出缺一张图"和"整个导不出来"比起来，前者明显更有用
	imgCount := 0
	for i := range res.Events {
		for j, im := range res.Events[i].Images {
			imgCount++
			if !inlineImages {
				continue
			}
			p, perr := s.ImagePath(userID, sessionID, runID, im.Name)
			if perr != nil {
				log.Printf("[warn] trace: 导出跳过不合法的图片名 %q（run=%s）", im.Name, runID)
				continue
			}
			png, rerr := os.ReadFile(p)
			if rerr != nil {
				continue
			}
			res.Events[i].Images[j].URL = "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
		}
	}

	// 事件瘦身：归属清掉、重复上下文折叠。prevMsgs 永远存上一条请求的**原始全量**
	// messages（折叠只改当前事件，不改游标），链式折叠才不会把增量叠在增量上
	var prevMsgs json.RawMessage
	for i := range res.Events {
		e := &res.Events[i]
		e.UserID = 0
		if e.Kind != KindLLMRequest || len(e.Messages) == 0 {
			continue
		}
		full := e.Messages
		if len(prevMsgs) > 0 {
			e.Messages = foldMessages(prevMsgs, full)
		}
		prevMsgs = full
	}

	exp := &Export{
		Format:     ExportFormat,
		Version:    ExportVersion,
		ExportedAt: s.now(),
		Run:        runMetaFrom(head, end, size, runID, s.now),
		Timeline:   buildTimeline(res.Events),
		Events:     res.Events,
	}
	exp.Images.Inlined = inlineImages
	exp.Images.Count = imgCount
	if !inlineImages && imgCount > 0 {
		exp.Images.Hint = "截图未内联（省体积）。要排查视觉/版面问题时，用 images=1 重新导出即可带上全部截图。"
	}
	return exp, nil
}

// foldMessages 把 cur 里与 prev 逐字节相同的前缀折叠掉，返回一个自描述的 JSON 对象；
// 前缀对不上就原样返回 cur（全量）。delta 对象里带 note，读到它的 AI 不需要
// 任何额外文档就能知道"完整上下文 = 上一条请求 + new_messages"。
func foldMessages(prev, cur json.RawMessage) json.RawMessage {
	var pm, cm []json.RawMessage
	if json.Unmarshal(prev, &pm) != nil || json.Unmarshal(cur, &cm) != nil {
		return cur
	}
	if len(pm) == 0 || len(pm) > len(cm) {
		return cur
	}
	for i := range pm {
		if !bytes.Equal(pm[i], cm[i]) {
			return cur
		}
	}
	delta := struct {
		Note        string            `json:"note"`
		PrevCount   int               `json:"prev_count"`
		NewMessages []json.RawMessage `json:"new_messages"`
	}{
		Note: "增量：本请求的完整上下文 = 上一条 llm_request 的全部消息（见 prev_count，与本次前缀逐字节相同，故不再重复给出）+ new_messages。只有每个 run 的第一条 llm_request 给全量。",
		PrevCount:   len(pm),
		NewMessages: cm[len(pm):],
	}
	b, err := json.Marshal(delta)
	if err != nil {
		return cur
	}
	return b
}

// buildTimeline 每个事件压成一行。行内的原文摘录一律按 rune 截断（按字节截中文
// 会切出半个字），截断只发生在 timeline：完整原文永远在 events 里。
func buildTimeline(events []Event) []string {
	out := make([]string, 0, len(events))
	for _, e := range events {
		out = append(out, timelineLine(e))
	}
	return out
}

func timelineLine(e Event) string {
	turn := ""
	if e.Turn != nil {
		turn = fmt.Sprintf("[第%d轮] ", *e.Turn+1)
	}
	call := ""
	if e.ToolCallID != "" {
		call = " call=" + e.ToolCallID
	}

	switch e.Kind {
	case KindRunStart:
		return fmt.Sprintf("#%d ▶ run开始 user=%q model=%s 工具=%v", e.Seq, tlClip(e.UserContent, 80), e.Model, e.Tools)
	case KindLLMRequest:
		return fmt.Sprintf("#%d %s→请求模型 上下文%d条/%s", e.Seq, turn, e.MessageCount, tlBytes(e.Bytes))
	case KindLLMResponse:
		line := fmt.Sprintf("#%d %s←回复 finish=%s", e.Seq, turn, e.FinishReason)
		if e.Content != "" {
			line += " 正文=" + tlClip(e.Content, 160)
		}
		for _, tc := range e.ToolCalls {
			line += fmt.Sprintf(" 想调用 %s(%s)", tc.Name, tlClip(tc.Arguments, 200))
		}
		if e.Usage != nil {
			line += fmt.Sprintf(" 本轮%d tok", e.Usage.Total)
		}
		return line
	case KindToolCall:
		return fmt.Sprintf("#%d %s🛠 %s%s 参数=%s", e.Seq, turn, e.ToolName, call, tlClip(e.Args, 400))
	case KindToolResult:
		mark := "成功"
		if e.Error != "" {
			mark = "失败:" + tlClip(e.Error, 120)
		}
		return fmt.Sprintf("#%d %s🛠 %s%s 返回(%dms) %s 结果=%s", e.Seq, turn, e.ToolName, call, e.DurationMS, mark, tlClip(e.Result, 400))
	case KindSubStep:
		if e.Sub == nil {
			return fmt.Sprintf("#%d ·子步骤(空)", e.Seq)
		}
		return fmt.Sprintf("#%d ·子步骤 %s/%s %s", e.Seq, e.Sub.Name, e.Sub.Stage, tlClip(e.Sub.Text, 160))
	case KindUsage:
		if e.Usage == nil {
			return fmt.Sprintf("#%d ·用量(空)", e.Seq)
		}
		line := fmt.Sprintf("#%d %s·用量 %s=%d tok", e.Seq, turn, e.Component, e.Usage.Total)
		if e.Usage.ImageTokens > 0 {
			line += fmt.Sprintf("（图片%d）", e.Usage.ImageTokens)
		}
		return line
	case KindError:
		return fmt.Sprintf("#%d ‼错误 %s", e.Seq, tlClip(e.Error, 300))
	case KindRunEnd:
		line := fmt.Sprintf("#%d ■ run结束 status=%s", e.Seq, e.Status)
		if s := e.Summary; s != nil {
			line += fmt.Sprintf(" %d轮 %d次工具 合计%d tok", s.Turns, s.ToolCalls, s.Total.Total)
		}
		return line
	default:
		return fmt.Sprintf("#%d %s", e.Seq, e.Kind)
	}
}

// tlClip 按 rune 截断。timeline 是给人/AI 扫的摘要层，中文按字节切会切出半个字。
func tlClip(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}

func tlBytes(n int) string {
	switch {
	case n >= 1<<20:
		return fmt.Sprintf("%.1fMB", float64(n)/(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%.1fKB", float64(n)/(1<<10))
	default:
		return fmt.Sprintf("%dB", n)
	}
}
