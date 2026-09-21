package agent

// deck-v2 的量测与看图审查：与 v1（vision_review.go）同构，但走 CaptureV2/ReviewV2。
//
// 分流的判据是 deck 格式：v2 → 这里的实现；v1 → 原 measureDeck/reviewPages。
// 分流点在 toolReviewSlidesV2 内部（它是 v2 工具集挂的 review_slides），
// v1 的工具仍挂原实现——两套互不干扰。

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"html-ppt/backend/internal/service/deck"
	"html-ppt/backend/internal/trace"
	"html-ppt/backend/internal/vision"
)

// measureDeckV2 只量测不读图（write_pages 之后的自动体检、review_slides 数字复查共用）。
// 返回文本附进工具结果；跑失败显式说明"没跑成 ≠ 没问题"（v1 实测教训的延续）。
func (a *AgentService) measureDeckV2(ctx context.Context, uid uint, deckID string) string {
	d, err := a.captureDeckV2(ctx, uid, deckID, nil, true)
	if err != nil {
		if errors.Is(err, errVisionOff) {
			return ""
		}
		return "\n\n版面量测：这次**没跑成**（渲染/量测失败），不代表版面没问题——" +
			"数字缺失和数字正常不是一回事。别把这次当成检查过了。"
	}

	var b strings.Builder
	b.WriteString("\n\n版面量测（自动跑的，只有数字，没有看图；固定画布 1920×1080）：\n")
	b.WriteString(vision.DigestV2(d))
	b.WriteString("\n溢出（⚠ 标记）的页必须修复：删内容或精简文字，不要缩字号糊弄。")
	b.WriteString("\n填充率偏空（⚠ 标记）的页优先补实质内容（数据、例子、图表行），或换信息密度更高的版式重写；补不动就如实汇报遗留——为凑填充率注水，比留白更伤观感。")
	b.WriteString(fmt.Sprintf("\n要看画面：调 review_slides，pages 传 1 基页号（一次最多 %d 页），"+
		"优先挑溢出页、填充率偏空的页和最小字号偏小的页。", maxReviewPages))
	return clipRunes(b.String(), reviewChars)
}

// captureDeckV2 与 v1 captureDeck 同一条 nonce 门票链路，只是 Capture 换成 CaptureV2。
func (a *AgentService) captureDeckV2(ctx context.Context, uid uint, deckID string, shoot []int, measureOnly bool) (*vision.Deck2, error) {
	if !a.Vision || a.VisionGrants == nil {
		return nil, errVisionOff
	}
	nonce, err := a.VisionGrants.Issue(uid, deckID)
	if err != nil {
		log.Printf("[warn] 视觉审查(v2)：发票据失败,跳过 err：%v", err)
		trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Sub: &trace.SubStep{
			Name: "vision", Stage: "error", Text: "发票据失败: " + err.Error(),
		}})
		return nil, fmt.Errorf("发票据失败: %w", err)
	}
	url := strings.TrimRight(a.VisionBaseURL, "/") + "/api/render/" + nonce

	trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Sub: &trace.SubStep{
		Name: "vision", Stage: "grant",
		Data: map[string]any{"deck_id": deckID, "url": url, "shoot": shoot, "measure_only": measureOnly, "v": 2},
	}})

	d, err := vision.CaptureV2(ctx, vision.OptionsV2{
		URL: url, ChromePath: a.ChromePath,
		Shoot: shoot, MeasureOnly: measureOnly,
	})
	if err != nil {
		log.Printf("[warn] 视觉审查(v2)：渲染/量测失败，跳过 err：%v", err)
		trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Sub: &trace.SubStep{
			Name: "vision", Stage: "error", Text: "渲染/量测失败: " + err.Error(),
		}})
		return nil, fmt.Errorf("渲染/量测失败: %w", err)
	}

	// 版式指纹从模板注册表补进量测结果：hero/quote 的填充率告警豁免要按模板判。
	// 拿不到（deck 还没选模板/读失败）就保持 nil，退化为不豁免——少一个优化，不出错。
	if pats, perr := a.DeckService.LayoutPatterns(uid, deckID); perr == nil {
		d.Patterns = pats
	}

	// 观测：量测数字与截图进 trace（与 v1 同一形态，页字段不同）
	round := trace.From(ctx).NextImageRound()
	pages := make([]map[string]any, 0, len(d.Slides))
	images := make([]trace.ImageEvent, 0, len(d.Slides))
	for _, s := range d.Slides {
		pages = append(pages, map[string]any{
			"page": s.Index + 1, "title": s.Title, "layout": s.Layout,
			"overflow_y": s.OverflowY, "overflow_x": s.OverflowX,
			"min_font_px": s.MinFontPx, "chars": s.Chars,
			"bottom_gap": s.BottomGap, "kids": s.Children,
		})
		if len(s.PNG) == 0 {
			continue
		}
		label := fmt.Sprintf("第 %d 页", s.Index+1)
		if s.Title != "" {
			label += "：" + s.Title
		}
		images = append(images, trace.ImageEvent{
			Name:  fmt.Sprintf("v%d-p%03d.png", round, s.Index+1),
			Label: label,
			Bytes: s.PNG,
		})
	}
	trace.Emit(ctx, trace.Event{
		Kind: trace.KindSubStep,
		Sub: &trace.SubStep{
			Name: "vision", Stage: "capture",
			Text: vision.DigestV2(d),
			Data: map[string]any{
				"pages": pages, "canvas_w": d.CanvasW, "canvas_h": d.CanvasH,
				"round": round, "hard_n": len(vision.HardFindingsV2(d)), "shot_n": len(images),
			},
		},
		Images: images,
	})

	return d, nil
}

// toolReviewSlidesV2 v2 的 review_slides：pages 留空 = 数字复查（不占配额）；
// 点名页 = 量测 + 截图 + 看图子调用。1 基/0 基换算只在这里做一次（v1 的教训）。
func (a *AgentService) toolReviewSlidesV2(ctx context.Context, arguments string) (string, error) {
	if !a.Vision || a.VisionGrants == nil {
		return "", errVisionOff
	}
	uid, err := toolUID(ctx)
	if err != nil {
		return "", err
	}
	var args ReviewSlidesArgs
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("review_slides 参数不是合法json err:%w", err)
	}
	// 归属校验（顺带确认 deck 存在）
	if _, err := a.DeckService.GetHTML(uid, args.DeckID); err != nil {
		return "", err
	}

	pages, err := normalizePages(args.Pages)
	if err != nil {
		return "", err
	}
	if len(pages) > maxReviewPages {
		return "", fmt.Errorf("一次最多看 %d 页，你点了 %d 页（%v）：挑最像有问题的 %d 页，剩下的分几次调用",
			maxReviewPages, len(pages), pages, maxReviewPages)
	}
	if len(pages) == 0 {
		return a.measureDeckV2(ctx, uid, args.DeckID), nil
	}

	shoot := make([]int, 0, len(pages))
	for _, p := range pages {
		shoot = append(shoot, p-1)
	}
	d, err := a.captureDeckV2(ctx, uid, args.DeckID, shoot, false)
	if err != nil {
		if errors.Is(err, errVisionOff) {
			return "", errors.New("视觉审查功能没开（features.vision），看不了画面")
		}
		return "", fmt.Errorf("视觉审查不可用：%w（详见服务端日志）", err)
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("\n\n视觉审查（无头浏览器实拍 %d 页，deck 共 %d 页）：\n", len(pages), len(d.Slides)))
	b.WriteString(vision.DigestV2(d))

	onDelta := func(string) {}
	if emit := emitFrom(ctx); emit != nil {
		onDelta = func(text string) {
			_ = emit(StreamEvent{Type: EventTypeSubDelta, ToolName: "review_slides", Content: text})
		}
	}

	client := a.clientFor(ctx)
	// 主模型关注点 + 逐页大纲意图随图注入：看图是一次隔离调用，主对话的上下文
	// 它拿不到——focus 不传就只做泛泛审查，大纲意图不给就审不了"内容切题"。
	rc := &vision.ReviewContext{Focus: args.Focus, PageBriefs: a.outlineBriefs(uid, args.DeckID)}
	if rc.Focus == "" && len(rc.PageBriefs) == 0 {
		rc = nil
	}
	rr, err := vision.ReviewV2(ctx, client, a.ModelID, d, pages, rc, onDelta)
	if err != nil {
		log.Printf("[warn] 视觉审查(v2)：看图失败，只回量测 err:%v", err)
		trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Sub: &trace.SubStep{
			Name: "vision", Stage: "review_error", Text: err.Error(),
		}})
		b.WriteString(visionFallbackNote)
		return clipRunes(b.String(), reviewChars), nil
	}

	trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Sub: &trace.SubStep{
		Name: "vision", Stage: "review_prompt", Text: rr.Prompt,
		Data: map[string]any{"images": rr.Images, "model": a.ModelID, "pages": pages, "focus": args.Focus},
	}})
	trace.Usage(ctx, trace.CompVision, usagePartFrom(rr.Usage))
	trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Sub: &trace.SubStep{
		Name: "vision", Stage: "review_response", Text: rr.Report,
	}})

	b.WriteString(visionReportMarker)
	b.WriteString(rr.Report)
	return clipRunes(b.String(), reviewChars), nil
}

// outlineBriefs 每页的大纲意图摘要（页号 → 标题+要点+备注）。看图调用拿不到
// 主对话的上下文，这几十个字是它判断"内容切不切题"的唯一依据。
// 读大纲失败返回 nil（裸看图，不阻塞审查——大纲是增强，不是前置条件）。
func (a *AgentService) outlineBriefs(uid uint, deckID string) map[int]string {
	o, err := a.DeckService.ReadOutline(uid, deckID)
	if err != nil {
		return nil
	}
	out := make(map[int]string, len(o.Pages))
	for _, p := range o.Pages {
		out[p.No] = briefOf(p)
	}
	return out
}

// briefOf 单页摘要（纯函数，截断到 ~160 字：意图够了，别把看图提示词撑胖）。
func briefOf(p deck.OutlinePage) string {
	var b strings.Builder
	fmt.Fprintf(&b, "标题「%s」", p.Title)
	if len(p.Points) > 0 {
		b.WriteString("；要点：" + strings.Join(p.Points, "；"))
	}
	if p.Notes != "" {
		b.WriteString("；备注：" + p.Notes)
	}
	s := b.String()
	if r := []rune(s); len(r) > 160 {
		s = string(r[:160]) + "…"
	}
	return s
}
