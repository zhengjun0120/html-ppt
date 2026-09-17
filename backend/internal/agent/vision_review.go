package agent

import (
	"context"
	"errors"
	"fmt"
	"html-ppt/backend/internal/trace"
	"html-ppt/backend/internal/vision"
	"encoding/json"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const reviewChars = 1500 //限制返回给模型的报告字数长度

// 单次最多看几页。图片 token 是整个 agent 最贵的一项（一张 1244×700 实测约 836 token），
// 不设上限的话"让 agent 自己挑页"很容易退化成"每页都挑"——那正好是这次要改掉的行为。
const maxReviewPages = 6

// reviewQuotaPerRun 一次 run 里 review_slides 的硬配额（Tool.MaxPerRun 的值）。
// 预期用法三步：全量审可疑页（1~2 次覆盖完）→ 批量修复 → 确认（1 次，只看改过的页）。
// 它是"审查→修复→再审"死循环的唯一硬闸门：实测一次 run 里模型审了 5 轮、改了 5 轮，
// 直到 maxTurns 被强制收尾——提示词拦不住它，计数器拦得住。
const reviewQuotaPerRun = 3

// 开关没开。与"跑了但失败"必须分开：后者要出声，前者本来就该安静。
var errVisionOff = errors.New("视觉审查没开")

// 视觉审查分两条路，因为代价差两个数量级：
//
//	measureDeck —— 只量测：无头浏览器把每页的 fit / 最小字号 / 溢出 / 字数算出来，几秒钟、零模型开销。
//	reviewPages —— 量测 + 只对点名的几页截图，再让模型看这几张图（图片 token 是整个 agent 最贵的一项）。
//
// 以前只有第二条路、而且是全量跑的：每次 write_deck 都要等所有页看完（9 页实测约 2 分钟）。
// 现在数字免费、先摆出来，要不要花时间看图由 agent 点名——所以它得先知道自己点得动。
func (a *AgentService) captureDeck(ctx context.Context, uid uint, deckID string, pages []int, measureOnly bool) (*vision.Deck, error) {
	if !a.Vision || a.VisionGrants == nil {
		return nil, errVisionOff
	}
	//获取门票
	nonce, err := a.VisionGrants.Issue(uid, deckID)
	if err != nil {
		log.Printf("[warn] 视觉审查：发票据失败,跳过 err：%v", err)
		trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Sub: &trace.SubStep{
			Name: "vision", Stage: "error", Text: "发票据失败: " + err.Error(),
		}})
		return nil, fmt.Errorf("发票据失败: %w", err)
	}
	//预览路由：前缀必须与 router.go 里注册的那条一致（挂在 api 组 → /api/render/）。
	//两处不一致的表现是 404，而这里会把 404 归到"渲染/量测失败"，
	//日志里看不出是路由写错了。
	url := strings.TrimRight(a.VisionBaseURL, "/") + "/api/render/" + nonce

	trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Sub: &trace.SubStep{
		Name: "vision", Stage: "grant",
		Data: map[string]any{"deck_id": deckID, "url": url, "pages": pages, "measure_only": measureOnly},
	}})

	//拿到截图和量测结果
	deck, err := vision.Capture(ctx, vision.Options{
		URL: url, ChromePath: a.ChromePath,
		Shoot: pages, MeasureOnly: measureOnly,
	})
	if err != nil {
		log.Printf("[warn] 视觉审查：渲染/量测失败，跳过 err：%v", err)
		trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Sub: &trace.SubStep{
			Name: "vision", Stage: "error", Text: "渲染/量测失败: " + err.Error(),
		}})
		return nil, fmt.Errorf("渲染/量测失败: %w", err)
	}

	a.traceCapturedDeck(ctx, deck)
	return deck, nil
}

// 只量测，不读图。write_deck 之后的自动体检走这条。
//
// 返回的文本必须自己说清"这只是一半"并且给出继续看的办法：只说"有问题"不说
// "想看画面怎么点"，agent 就没法判断还有没有数字看不出的毛病——而恰恰是那批毛病
// （文字被裁、元素重叠、对比度不足、图与标题讲的不是一回事）只有看图才知道。
func (a *AgentService) measureDeck(ctx context.Context, uid uint, deckID string) string {
	deck, err := a.captureDeck(ctx, uid, deckID, nil, true)
	if err != nil {
		if errors.Is(err, errVisionOff) {
			return "" //开关没开：这一项本来就不该出现
		}
		//跑失败必须出声。返回空串的话，write_deck 的 review 字段会被 omitempty 直接省掉，
		//而"没有这个字段"和"量过了、没问题"在模型眼里长得一模一样——它会当成检查过了。
		//（这条路径不是假想：data/decks 里那份没有 Reveal 的沙箱测试夹具实拍时就走到过这里。）
		return "\n\n版面量测：这次**没跑成**（渲染/量测失败），不代表版面没问题——" +
			"数字缺失和数字正常不是一回事。别把这次当成检查过了；用户问起就说这次没量成。"
	}

	var b strings.Builder
	b.WriteString("\n\n版面量测（自动跑的，只有数字，没有看图）：\n")
	b.WriteString(deck.Text())
	if note := a.deckRhythmNote(uid, deckID, len(deck.Slides)); note != "" {
		b.WriteString(note)
	}
	b.WriteString("\n数字都在范围内**不等于**版面没问题：文字被裁、元素重叠、贴边、对比度不足、" +
		"这页的图和标题讲的不是一回事——这些数字看不出来，只有真看一眼才知道。")
	b.WriteString(fmt.Sprintf("\n要看画面：调 review_slides，pages 传 1 基页号（一次最多 %d 页），"+
		"优先挑 fit 低于 0.90、最小字号低于 31px、溢出大于 1.02 的页，以及自己写的时候就拿不准的页。", maxReviewPages))
	return clipRunes(b.String(), reviewChars)
}

// visionReportMarker 是 review_slides 结果里"真的看到了图"的标记：
// 看图成功的结果一定带它，"只回数字"（pages 留空）和"看图失败"的结果一定不带。
// execTool 用它决定要不要把这次调用从配额里退还——配额只数烧了钱、产出了报告的审查。
const visionReportMarker = "\n看图审查：\n"

// visionFallbackNote 是看图失败时附在量测数字后面的说明。除了告诉模型"这次没有
// 看图的结论"，还要告诉它"这次不占配额"：实测（trace 1789613158323-30a4）模型
// 看到降级结果后想重试，却被配额挡住，只能拿着数字收尾。
const visionFallbackNote = "\n(看图部分失败，以上只有量测数字；这次没有消耗审查配额，可以直接重试)"

// reviewMeasureOnlyArgs 判断 review_slides 的调用参数是不是"数字复查"（pages 缺省或空）。
// 解析失败返回 false：宁把它当真审查计一次配额，也不能让一个坏参数绕过闸门。
// 这个判定只服务配额分流，工具真正的参数解析仍归它自己的实现。
func reviewMeasureOnlyArgs(argsJSON string) bool {
	var a struct {
		Pages []int `json:"pages"`
	}
	if json.Unmarshal([]byte(argsJSON), &a) != nil {
		return false
	}
	return len(a.Pages) == 0
}

// colorPageRe 数"色面页"：挂了整页换色变体类（.bg-ink / .bg-accent）的 section。
// 必须锚定在 <section 标签上：正文文字里提到"bg-ink"不该被算进来。
// \b 保证按 class token 边界匹配（不会命中想象中的 xbg-ink）。
var colorPageRe = regexp.MustCompile(`<section[^>]*class="[^"]*\bbg-(?:ink|accent)\b`)

// deckRhythmNote 是一份**整份级**的免费检查：≥8 页的 deck 一张色面页都没有、
// 也没有背景渐变时，点一句。色面页是"这份被设计过"的第一印象里最便宜的装置
// （提示词的"整页色彩节奏"一节），但逐页的量测数字看不见"整份从头到尾一个底色"——
// 它只在整份维度上成立，所以放在这里而不是页级报告里。
// 语气刻意是"点一下"而不是硬判定：一份素净的纸感 deck 是合法选择，agent 自己判断。
func (a *AgentService) deckRhythmNote(uid uint, deckID string, slides int) string {
	if slides < 8 {
		return ""
	}
	html, err := a.DeckService.GetHTML(uid, deckID)
	if err != nil {
		return ""
	}
	if colorPageRe.MatchString(html) {
		return ""
	}
	// 有背景渐变的主题整份已经有底色氛围，不再要色面页来撑节奏
	if strings.Contains(html, `"bg_gradient":"`) {
		return ""
	}
	return fmt.Sprintf("\n\n整份层面：这 %d 页里没有一张色面页（.bg-ink / .bg-accent），底色从头到尾一个样——"+
		"这是最容易被看成\"没排版过\"的一种单调。章节页、金句页、单数字页是天然的换色页；"+
		"确属有意为之的素净风格可以忽略这条。", slides)
}

// 量测 + 只看点名的这几页。pages 是 **1 基**页号（与 list_slides 的 position、"第 N 页"同一套）。
func (a *AgentService) reviewPages(ctx context.Context, uid uint, deckID string, pages []int) (string, error) {
	//1→0 基只在这一处做：capture 的 Shoot 与 Slide.Index 同口径。换算写错不会报错，
	//只会安静地拍错页，而报告里照样写着"第 N 页"——所以只有一个地方做这件事。
	shoot := make([]int, 0, len(pages))
	for _, p := range pages {
		shoot = append(shoot, p-1)
	}

	deck, err := a.captureDeck(ctx, uid, deckID, shoot, false)
	if err != nil {
		if errors.Is(err, errVisionOff) {
			//开关没开和渲染失败对 agent 的意义不同：前者重试也没用，后者值得再试一次
			return "", errors.New("视觉审查功能没开（features.vision），看不了画面")
		}
		return "", fmt.Errorf("视觉审查不可用：%w（详见服务端日志）", err)
	}

	sel, missing := deck.Select(pages)
	if len(sel.Slides) == 0 {
		return "", fmt.Errorf("要审的页一份也没取到：这份 deck 只有 %d 页，你点的是 %s", len(deck.Slides), joinPages(pages))
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("\n\n视觉审查（无头浏览器实拍 %d 页，共 %d 页）：\n", len(sel.Slides), len(deck.Slides)))
	if len(missing) > 0 {
		//必须说出来：静默丢掉的表现是模型以为那页看过了
		b.WriteString(fmt.Sprintf("注意：这份 deck 只有 %d 页，第 %s 页不存在，这几页没审。\n", len(deck.Slides), joinPages(missing)))
	}
	b.WriteString(vision.Digest(sel))

	//整份层面还有哪些页带着硬判定：不把它们逐页数字铺一遍（模型刚在 write_deck 拿到过
	//全量数字，重复一遍要占上千 token），但要说清"没看的页里还有几页有问题"——
	//否则 agent 会以为没点到的页就是干净的。
	if rest := unselectedFindingsPages(vision.Analyze(deck), pages); len(rest) > 0 {
		b.WriteString(fmt.Sprintf("\n（没看的页里，第 %s 页还有硬判定）", joinPages(rest)))
	}

	onDelta := func(string){}
	if emit := emitFrom(ctx); emit!=nil{
		onDelta = func(text string){
			if emitErr := emit(StreamEvent{Type:EventTypeSubDelta,ToolName: "review_slides",Content: text});emitErr != nil{
				log.Printf("[warn] sub_delta 推送失败 err:%v",emitErr)
			}
		}
	}

	// 获取客户端
	client := a.clientFor(ctx)
	// 把这几页的图 + 整份算出来的相关判定发给模型
	rr, err := vision.Review(ctx, client, a.ModelID, sel, vision.ScopeFindings(vision.Analyze(deck), pages),onDelta)
	if err != nil {
		log.Printf("[warn] 视觉审查：看图失败，只回量测 err:%v", err)
		trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Sub: &trace.SubStep{
			Name: "vision", Stage: "review_error", Text: err.Error(),
		}})
		b.WriteString(visionFallbackNote)
		return clipRunes(b.String(), reviewChars), nil
	}

	// 观测：把实际发出去的提示词记下来。审查的每个结论都建立在这段提示词上，
	// 而它过去只存在于代码里——提示词改过一版之后，旧报告就没法解读了
	trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Sub: &trace.SubStep{
		Name: "vision", Stage: "review_prompt", Text: rr.Prompt,
		Data: map[string]any{"images": rr.Images, "model": a.ModelID, "pages": pages},
	}})

	// 观测：看图是**分开计费**的另一次调用，而且是整个 agent 里最贵的一次
	//（几张 1244×700 的图，图片 token 单独计量）。它的用量过去完全没被计入，
	// 所以"这次对话花了多少 token"一直是偏小的
	trace.Usage(ctx, trace.CompVision, usagePartFrom(rr.Usage))

	trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Sub: &trace.SubStep{
		Name: "vision", Stage: "review_response", Text: rr.Report,
	}})

	b.WriteString(visionReportMarker)
	b.WriteString(rr.Report)
	return clipRunes(b.String(), reviewChars), nil
}

// 没被点名、但带着页级硬判定的页（1 基，升序去重）。只看页级：
// 版式重复那种 deck 级判定会牵出一大半页号，列出来反而让人以为"半份 deck 都有问题"。
func unselectedFindingsPages(fs []vision.Finding, pages []int) []int {
	sel := map[int]bool{}
	for _, p := range pages {
		sel[p] = true
	}
	seen := map[int]bool{}
	var out []int
	for _, f := range fs {
		if f.Level != vision.LevelPage {
			continue
		}
		for _, p := range f.Pages {
			if !sel[p] && !seen[p] {
				seen[p] = true
				out = append(out, p)
			}
		}
	}
	sort.Ints(out)
	return out
}

func joinPages(pages []int) string {
	ps := make([]string, 0, len(pages))
	for _, p := range pages {
		ps = append(ps, strconv.Itoa(p))
	}
	return strings.Join(ps, "、")
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
			continue //只量测时本来就没有图，不是失败
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
				"shot_n":   len(images),
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
