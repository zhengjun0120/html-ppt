package agent

import (
	"context"
	"fmt"
	"html-ppt/backend/internal/trace"
	"html-ppt/backend/internal/vision"
	"log"
	"strings"
)

const reviewChars = 1500 //限制返回给模型的报告字数长度

// 运行视觉审查
//
// 观测在这里最能体现价值：capture 阶段每页都实拍了一张 1244×700 的 PNG，而它过去
// 发给审查模型之后就丢了（vision.Slide.PNG 是 json:"-"），于是报告里说"第 3 页文字被裁"
// 的时候，人没有任何办法核对——只能选择相信或不信。现在截图落盘、逐页量测进 trace。
func (a *AgentService) runVisionReview(ctx context.Context, uid uint, deckID string) string {
	if !a.Vision || a.VisionGrants == nil {
		return ""
	}
	//获取门票
	nonce, err := a.VisionGrants.Issue(uid, deckID)
	if err != nil {
		log.Printf("[warn] 视觉审查：发票据失败,跳过 err：%v", err)
		trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Sub: &trace.SubStep{
			Name: "vision", Stage: "error", Text: "发票据失败: " + err.Error(),
		}})
		return ""
	}
	//预览路由：前缀必须与 router.go 里注册的那条一致（挂在 api 组 → /api/render/）。
	//两处不一致的表现是 404，而 runVisionReview 会把 404 归到"渲染/量测失败"，
	//日志里看不出是路由写错了。
	url := strings.TrimRight(a.VisionBaseURL, "/") + "/api/render/" + nonce

	trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Sub: &trace.SubStep{
		Name: "vision", Stage: "grant",
		Data: map[string]any{"deck_id": deckID, "url": url},
	}})

	//拿到截图和量测结果
	deck, err := vision.Capture(ctx, vision.Options{URL: url, ChromePath: a.ChromePath})
	if err != nil {
		log.Printf("[warn] 视觉审查：渲染/量测失败，跳过 err：%v", err)
		trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Sub: &trace.SubStep{
			Name: "vision", Stage: "error", Text: "渲染/量测失败: " + err.Error(),
		}})
		return ""
	}

	a.traceCapturedDeck(ctx, deck)

	var b strings.Builder
	b.WriteString("\n\n视觉审查（无头浏览器实拍 + 量测）：\n")
	b.WriteString(deck.Text())

	// 获取客户端
	client := a.clientFor(ctx)
	// 把所有信息整合发给模型 并获取结果
	rr, err := vision.Review(ctx, client, a.ModelID, deck)
	if err != nil {
		log.Printf("[warn] 视觉审查：看图失败，只回量测 err:%v", err)
		trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Sub: &trace.SubStep{
			Name: "vision", Stage: "review_error", Text: err.Error(),
		}})
		b.WriteString("\n (看图部分失败，以上是程序量测结果)")
		return clipRunes(b.String(), reviewChars)
	}

	// 观测：把实际发出去的提示词记下来。审查的每个结论都建立在这段提示词上，
	// 而它过去只存在于代码里——提示词改过一版之后，旧报告就没法解读了
	trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Sub: &trace.SubStep{
		Name: "vision", Stage: "review_prompt", Text: rr.Prompt,
		Data: map[string]any{"images": rr.Images, "model": a.ModelID},
	}})

	// 观测：看图是**分开计费**的另一次调用，而且是整个 agent 里最贵的一次
	//（几张 1244×700 的图，图片 token 单独计量）。它的用量过去完全没被计入，
	// 所以"这次对话花了多少 token"一直是偏小的
	trace.Usage(ctx, trace.CompVision, usagePartFrom(rr.Usage))

	trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Sub: &trace.SubStep{
		Name: "vision", Stage: "review_response", Text: rr.Report,
	}})

	b.WriteString("\n看图审查：\n")
	b.WriteString(rr.Report)
	return clipRunes(b.String(), reviewChars)
}

// traceCapturedDeck 把每一页的量测数字与截图交给观测。
//
// 量测数字单独成一张结构化卡片（而不是只塞进那个给模型看的文本摘要）：
// 页面上要按页对比 fit / 最小字号 / 溢出比，从一整段文本里抠数字既费劲又容易看错行。
func (a *AgentService) traceCapturedDeck(ctx context.Context, deck *vision.Deck) {
	round := trace.From(ctx).NextImageRound()

	pages := make([]map[string]any, 0, len(deck.Slides))
	images := make([]trace.ImageEvent, 0, len(deck.Slides))
	for _, s := range deck.Slides {
		pages = append(pages, map[string]any{
			"page":        s.Index + 1,
			"title":       s.Title,
			"fit":         s.EffScale,
			"min_font_px": s.MinFontPx,
			"overflow":    s.Overflow,
			"children":    s.Children,
			"chars":       s.Chars,
			"layout":      s.Layout,
		})
		if len(s.PNG) == 0 {
			continue
		}
		label := fmt.Sprintf("第 %d 页", s.Index+1)
		if s.Title != "" {
			label += "：" + s.Title
		}
		images = append(images, trace.ImageEvent{
			// 轮次前缀见 Recorder.NextImageRound 的说明：不加的话，
			// 同一 run 里第二次审查的图会盖掉第一次的
			Name:  fmt.Sprintf("v%d-p%03d.png", round, s.Index+1),
			Label: label,
			Bytes: s.PNG,
		})
	}

	trace.Emit(ctx, trace.Event{
		Kind: trace.KindSubStep,
		Sub: &trace.SubStep{
			Name: "vision", Stage: "capture",
			Text: deck.Text(),
			Data: map[string]any{
				"pages":    pages,
				"canvas_w": deck.CanvasW,
				"canvas_h": deck.CanvasH,
				"round":    round,
				"hard_n":   len(vision.HardFindings(deck)),
			},
		},
		Images: images,
	})
}

func clipRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "\n(报告过长已截断)"
}
