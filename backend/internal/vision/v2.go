package vision

// deck-v2 的渲染捕获层：固定画布（默认 1920×1080）+ runtime.js。
//
// 与 v1（capture.go，reveal.js）的本质区别：
//   - 画布不再来自页面主题块，而是模板常量（.deck 的 data-w/data-h），
//     是可信输入——不再需要 viewportFor 的离谱值防御；
//   - 没有 fit 兜底缩放，"内容放不放得下"只有一个判据：scrollHeight > 画布高；
//   - 页面定位不再依赖 Reveal.slide(i)，用 runtime 的 #/N 深链（hashchange → go()）；
//   - 非激活页只是 opacity:0（仍参与布局，scrollHeight 有效），所以量测可以
//     **一次过全部量完**，只有截图才需要逐页激活——这是比 v1 快一个量级的原因。

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

// V2 画布常量：模板 template.json canvas 的合法区间校验用（模板加载器已校验
// 正数，这里只防"deck.json 被手改坏"这种运行期意外）。
const (
	V2DefaultW = 1920
	V2DefaultH = 1080

	v2ReadyDelay  = 1200 * time.Millisecond // 字体 + runtime 初始化
	v2SlideSwitch = 650 * time.Millisecond  // 激活一页后等过渡（.slide 过渡 .5s）走完
)

// Slide2 固定画布下单页的量测结果。
// **json tag 与 measureJS2 返回的 key 逐字一致**，契约由测试守住（同 v1 的教训：
// tag 对不上不报错、字段静默为零、报告照着 0 说"没问题"）。
type Slide2 struct {
	Index     int     `json:"i"`         // 0 基
	Title     string  `json:"title"`     // h1/h2 文本
	Layout    string  `json:"layout"`    // data-layout（版式指纹的一部分）
	OverflowY bool    `json:"overflowY"` // 内容高于画布（横向滚动条也被算进高度时同样算它）
	OverflowX bool    `json:"overflowX"` // 内容宽于画布
	MinFontPx float64 `json:"minFontPx"` // 页面最小字号（设计像素）
	Chars     int     `json:"chars"`     // 可见文字量（CJK 字符 + 西文单词）
	Children  int     `json:"kids"`      // 直接子元素数
	BottomGap float64 `json:"bottomGap"` // 画布底边距最后可见内容的空隙（px）
	FillPct   float64 `json:"fillPct"`   // 内容纵向填充率（首尾可见内容的高差 / 画布高，%）
	PNG       []byte  `json:"-"`         // 截图不进 JSON
}

// Deck2 一份 deck 的量测集合。
type Deck2 struct {
	Slides  []Slide2
	CanvasW int
	CanvasH int
}

// OptionsV2 捕获参数。
type OptionsV2 struct {
	URL        string
	ChromePath string
	Timeout    time.Duration
	// MeasureOnly 只量测不截图（生成管线的每页自动体检走这条）。
	MeasureOnly bool
	// Shoot 要截图的页，**0 基**（与 Slide2.Index 同口径）。空 = 全拍。
	Shoot []int
}

// v2SlideCountJS 页数与画布：一次取齐，减少往返。
const v2MetaJS = `(function(){
  var deck = document.querySelector('.deck');
  if (!deck) return JSON.stringify({error:'no deck'});
  var slides = deck.querySelectorAll(':scope > .slide');
  var w = parseInt(deck.getAttribute('data-w'),10) || 1920;
  var h = parseInt(deck.getAttribute('data-h'),10) || 1080;
  return JSON.stringify({n: slides.length, w: w, h: h});
})()`

// v2MeasureAllJS 一次量测全部页（利用"非激活页仍参与布局"）。
// 逐页计算溢出/最小字号/字数/版式，一次性返回数组。
const v2MeasureAllJS = `(function(){
  var deck = document.querySelector('.deck');
  if (!deck) return JSON.stringify({error:'no deck'});
  var w = parseInt(deck.getAttribute('data-w'),10) || 1920;
  var h = parseInt(deck.getAttribute('data-h'),10) || 1080;
  var slides = deck.querySelectorAll(':scope > .slide');
  var out = [];
  for (var i = 0; i < slides.length; i++) {
    var sec = slides[i];
    var overflowY = sec.scrollHeight > h + 2;
    var overflowX = sec.scrollWidth > w + 2;

    var minFont = 999, kids = 0;
    var all = sec.querySelectorAll('*');
    for (var a = 0; a < all.length; a++) {
      var el = all[a];
      if (el.closest('.notes, .deck-footer, .deck-header, .kicker, .eyebrow, .tag, .slide-number')) continue; // 讲稿、页脚页眉与 kicker/眉标/tag 是模板定死的小字设计，不参与最小字号判定（只看真实内容）
      var fs = parseFloat(getComputedStyle(el).fontSize);
      if (fs > 0 && fs < minFont) minFont = fs;
    }
    for (var k = 0; k < sec.children.length; k++) {
      if (!sec.children[k].classList.contains('notes')) kids++;
    }

    var text = (sec.innerText || '').replace(/\s+/g, ' ');
    var cjk = (text.match(/[\u4e00-\u9fff\u3400-\u4dbf]/g) || []).length;
    var latin = (text.match(/[A-Za-z0-9][A-Za-z0-9'’\-]*/g) || []).length;
    var hEl = sec.querySelector('h1, h2');
    // getBoundingClientRect 受 transform scale 影响，除回去换算回设计像素
    var scale = parseFloat(deck.style.getPropertyValue('--deck-scale')) || 1;
    var bottom = 0;
    var top = h;
    for (var b = 0; b < sec.children.length; b++) {
      var c = sec.children[b];
      if (c.classList.contains('notes')) continue;
      var r = c.getBoundingClientRect();
      if (r.bottom > bottom) bottom = r.bottom;
      if (r.bottom > r.top && r.top / scale < top) top = r.top / scale;
    }
    var bottomGap = (h - bottom / scale);
    // 填充率 = 首尾可见内容的纵向跨度占画布高的比例。它与 bottomGap 互补：
    // 居中排版的稀疏页 bottomGap 会"看起来还行"（上下各留一半），fillPct 不会说谎。
    var fillPct = Math.max(0, Math.min(100, (bottom / scale - top) / h * 100));

    out.push({
      i: i,
      title: hEl ? (hEl.innerText||'').replace(/\s+/g,' ').trim().slice(0,30) : '',
      layout: sec.getAttribute('data-layout') || '',
      overflowY: overflowY,
      overflowX: overflowX,
      minFontPx: +minFont.toFixed(1),
      chars: cjk + latin,
      kids: kids,
      bottomGap: +bottomGap.toFixed(1),
      fillPct: +fillPct.toFixed(1)
    });
  }
  return JSON.stringify({slides: out, w: w, h: h});
})()`

// v2ActivateJS 把页面切到第 i 页（1 基 hash）。runtime 监听 hashchange → go()。
// 直接写 hash 而不是找内部函数：go() 是闭包私有，hash 是它的公开契约。
const v2ActivateJS = `(function(i){
  var target = '#/' + i;
  if (location.hash === target) return 'same';
  location.hash = target;
  return 'ok';
})`

type v2Meta struct {
	N      int `json:"n"`
	W      int `json:"w"`
	H      int `json:"h"`
	Error  string `json:"error"`
}

type v2MeasureAll struct {
	Slides []Slide2      `json:"slides"`
	W      int           `json:"w"`
	H      int           `json:"h"`
	Error  string        `json:"error"`
}

// CaptureV2 打开 deck（URL 由调用方拼好，含鉴权手段），量测全部页，
// 并按 Shoot 截图。截图前用 hash 深链逐页激活——只影响截图，不影响量测。
func CaptureV2(ctx context.Context, opt OptionsV2) (*Deck2, error) {
	if opt.Timeout == 0 {
		opt.Timeout = 120 * time.Second
	}
	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.Flag("hide-scrollbars", true),
		chromedp.WindowSize(V2DefaultW, V2DefaultH),
	)
	if opt.ChromePath != "" {
		allocOpts = append(allocOpts, chromedp.ExecPath(opt.ChromePath))
	}
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, allocOpts...)
	defer cancelAlloc()
	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx, chromedp.WithLogf(func(string, ...any) {}))
	defer cancelBrowser()
	runCtx, cancelTimeout := context.WithTimeout(browserCtx, opt.Timeout)
	defer cancelTimeout()

	if err := chromedp.Run(runCtx,
		emulation.SetDeviceMetricsOverride(V2DefaultW, V2DefaultH, 1, false),
		chromedp.Navigate(opt.URL),
		chromedp.WaitReady("body", chromedp.ByQuery),
	); err != nil {
		return nil, fmt.Errorf("打开 deck 失败 err:%w", err)
	}
	time.Sleep(v2ReadyDelay)

	var rawMeta string
	if err := chromedp.Run(runCtx, chromedp.Evaluate(v2MetaJS, &rawMeta)); err != nil {
		return nil, fmt.Errorf("读取 deck 元信息失败 err:%w", err)
	}
	var meta v2Meta
	if err := json.Unmarshal([]byte(rawMeta), &meta); err != nil {
		return nil, fmt.Errorf("deck 元信息解析失败(%s) err:%w", rawMeta, err)
	}
	if meta.Error != "" || meta.N <= 0 {
		return nil, fmt.Errorf("页面上没有可渲染的 deck（error=%q n=%d）——runtime 没起来或 SLIDES 区间为空", meta.Error, meta.N)
	}

	// 视口跟随模板画布：截图必须是"观众看到的那一张"。
	if meta.W != V2DefaultW || meta.H != V2DefaultH {
		if err := chromedp.Run(runCtx, emulation.SetDeviceMetricsOverride(int64(meta.W), int64(meta.H), 1, false)); err != nil {
			return nil, fmt.Errorf("按画布 %d*%d 设视口失败 err:%w", meta.W, meta.H, err)
		}
		time.Sleep(v2SlideSwitch)
	}

	var raw string
	if err := chromedp.Run(runCtx, chromedp.Evaluate(v2MeasureAllJS, &raw)); err != nil {
		return nil, fmt.Errorf("量测失败 err:%w", err)
	}
	var res v2MeasureAll
	if err := json.Unmarshal([]byte(raw), &res); err != nil {
		return nil, fmt.Errorf("量测结果解析失败(%s) err:%w", raw, err)
	}
	if res.Error != "" || len(res.Slides) == 0 {
		return nil, fmt.Errorf("量测结果为空（error=%q）", res.Error)
	}

	deck := &Deck2{Slides: res.Slides, CanvasW: res.W, CanvasH: res.H}
	if opt.MeasureOnly {
		return deck, nil
	}

	// 截图：逐页激活。runtime 的 go() 会在 hash 变化时把目标页置为 is-active
	//（opacity 1），过渡 .5s，等 v2SlideSwitch 再拍就是静止画面。
	for i := range deck.Slides {
		if !v2ShootPage(opt.Shoot, i) {
			continue
		}
		var act string
		if err := chromedp.Run(runCtx, chromedp.Evaluate(fmt.Sprintf(v2ActivateJS+"(%d)", i+1), &act)); err != nil {
			return nil, fmt.Errorf("第 %d 页激活失败 err:%w", i+1, err)
		}
		time.Sleep(v2SlideSwitch)

		var buf []byte
		if err := chromedp.Run(runCtx, chromedp.ActionFunc(func(c context.Context) error {
			b, e := page.CaptureScreenshot().WithFormat(page.CaptureScreenshotFormatPng).Do(c)
			buf = b
			return e
		})); err != nil {
			return nil, fmt.Errorf("第 %d 页截图失败 err:%w", i+1, err)
		}
		deck.Slides[i].PNG = buf
	}
	return deck, nil
}

func v2ShootPage(shoot []int, i int) bool {
	if len(shoot) == 0 {
		return true
	}
	for _, p := range shoot {
		if p == i {
			return true
		}
	}
	return false
}

// PrintPDF2 用浏览器打印整份 deck 为 PDF（16:9 整页）。
//
// 纸张尺寸从页面上的 deck 画布（data-w/data-h）按 96dpi 等比换算：
// 1920×1080 → 13.333×7.5 英寸；竖版画布（如 xhs-post 的 810×1080）→ 5.625×7.5，
// 打印出的每页就是设计画布的等比映射。读不到画布时退回 16:9 默认。
// PreferCSSPageSize=false：纸张尺寸由这里显式给定，不依赖页面声明 @page。
func PrintPDF2(ctx context.Context, url, chromePath string, timeout time.Duration) ([]byte, error) {
	if timeout == 0 {
		timeout = 120 * time.Second
	}
	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Headless,
		chromedp.DisableGPU,
	)
	if chromePath != "" {
		allocOpts = append(allocOpts, chromedp.ExecPath(chromePath))
	}
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, allocOpts...)
	defer cancelAlloc()
	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx, chromedp.WithLogf(func(string, ...any) {}))
	defer cancelBrowser()
	runCtx, cancelTimeout := context.WithTimeout(browserCtx, timeout)
	defer cancelTimeout()

	if err := chromedp.Run(runCtx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.ByQuery),
	); err != nil {
		return nil, fmt.Errorf("打开 deck 失败 err:%w", err)
	}
	// 字体加载完再打，否则 PDF 里是回退字体
	if err := chromedp.Run(runCtx, chromedp.Evaluate(`document.fonts ? document.fonts.ready.then(()=>'ok') : 'ok'`, nil)); err != nil {
		return nil, fmt.Errorf("等待字体失败 err:%w", err)
	}
	time.Sleep(v2ReadyDelay)

	// 画布比例决定纸张比例（竖版模板导出横版纸会裁掉内容）
	var canvas string
	pw, ph := 13.333, 7.5
	if err := chromedp.Run(runCtx, chromedp.Evaluate(
		`(function(){var d=document.querySelector('.deck');return d?(d.getAttribute('data-w')||'')+'x'+(d.getAttribute('data-h')||''):''})()`,
		&canvas)); err == nil && canvas != "" {
		var w, h int
		if _, err := fmt.Sscanf(canvas, "%dx%d", &w, &h); err == nil && w > 0 && h > 0 {
			ph = 7.5
			pw = 7.5 * float64(w) / float64(h)
		}
	}

	var pdf []byte
	if err := chromedp.Run(runCtx, chromedp.ActionFunc(func(c context.Context) error {
		b, _, e := page.PrintToPDF().
			WithPaperWidth(pw).
			WithPaperHeight(ph).
			WithPrintBackground(true).
			WithPreferCSSPageSize(false).
			Do(c)
		pdf = b
		return e
	})); err != nil {
		return nil, fmt.Errorf("PrintToPDF 失败 err:%w", err)
	}
	if len(pdf) == 0 {
		return nil, fmt.Errorf("PrintToPDF 返回空内容")
	}
	return pdf, nil
}

// ---------- 看图审查（v2） ----------

// DigestV2 逐页量测的文本摘要（发给审查模型的原料）。
func DigestV2(d *Deck2) string {
	out := make([]string, 0, len(d.Slides))
	for _, s := range d.Slides {
		flags := make([]string, 0, 3)
		if s.OverflowY {
			flags = append(flags, "纵向溢出")
		}
		if s.OverflowX {
			flags = append(flags, "横向溢出")
		}
		// 投影基准：1920 画布等比缩放后，14px 设计字≈9px 物理。全模板字号提级后
		// 内容文字下限 16-18px，这里 15 起提示、13 起硬判（比提级前的 18/14 收紧了实际效果）。
		if s.MinFontPx > 0 && s.MinFontPx < 15 {
			flags = append(flags, fmt.Sprintf("最小字号 %.0fpx 偏小", s.MinFontPx))
		}
		if s.FillPct > 0 && s.FillPct < 55 {
			flags = append(flags, fmt.Sprintf("填充率 %.0f%% 偏空（不到画布一半）", s.FillPct))
		}
		if d.CanvasH > 0 && s.BottomGap > float64(d.CanvasH)*0.35 {
			flags = append(flags, "底部空隙过大")
		}
		line := fmt.Sprintf("第 %d 页 [%s] %s：字数 %d，最小字号 %.0fpx，底部空隙 %.0fpx，填充率 %.0f%%，子元素 %d",
			s.Index+1, s.Layout, s.Title, s.Chars, s.MinFontPx, s.BottomGap, s.FillPct, s.Children)
		if len(flags) > 0 {
			line += " ⚠ " + strings.Join(flags, "、")
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

// HardFindingsV2 程序硬判定（overflow / 过小字号 / 大面积留白）。给 ScopeFindingsV2 与 trace 用。
func HardFindingsV2(d *Deck2) []string {
	var out []string
	for _, s := range d.Slides {
		if s.OverflowY || s.OverflowX {
			out = append(out, fmt.Sprintf("第 %d 页内容溢出画布（程序硬判定）", s.Index+1))
		}
		if s.MinFontPx > 0 && s.MinFontPx < 13 {
			out = append(out, fmt.Sprintf("第 %d 页最小字号 %.0fpx 低于投影下限（程序硬判定）", s.Index+1, s.MinFontPx))
		}
		// 大面积留白与溢出同样会让观众觉得"没做完"——deck-0031 实测教训：
		// 12 页里最空的页只画了 8% 的画布，当时的量测体系对此完全沉默。
		if s.FillPct > 0 && s.FillPct < 45 {
			out = append(out, fmt.Sprintf("第 %d 页内容只覆盖画布的 %.0f%%，大面积留白（程序硬判定）", s.Index+1, s.FillPct))
		}
	}
	return out
}

const reviewPromptV2 = `你是网页幻灯片的版面审查员。这套幻灯片渲染在固定画布上（所见即所得，截图就是观众看到的画面）。下面按页给你：一页渲染截图 + 这一页的量测数字。

量测字段：版式（data-layout）；字数；最小字号（px，设计像素）；底部空隙（内容距画布底边的空隙，接近 0 或负数 = 贴边/溢出）；子元素数。程序已给出的硬判定（溢出等）只是"值得看一眼"的信号，**不要复核数字**。

**只看给你的这几页**：没给你截图的页不要在报告里出现。你的任务：
1. 用眼睛确认硬判定在画面上真的存在，并说清是哪一处（引用页面实际文字）。
2. 找数字看不出来的问题：文字被裁、元素重叠、内容贴边、标题孤字换行、卡片高度不齐、代码块溢出、对比度不足、以及这一页内容与标题讲的不是一回事。

不要提"建议补充内容"这类意见，只说版面。不要重写内容。

**两条纪律，违反的代价是主模型陷入"改→再审→改"的无限返工：**
- **每条问题分级**：【硬】= 投屏会翻车的（被裁/重叠/看不清/明显错位）；【软】= 审美打磨。拿不准算【软】。
- **没把握就不报**：一页没问题就写没问题。报得准远比报得多有价值。

输出格式（严格照做）：
第 N 页｜【硬或软】｜问题类别｜具体是什么（引用页面实际文字）｜期望效果
每页最多一行，只写有问题的页。最后一行：合计：N 页有问题（硬 M 页）
都没问题则只输出一行：未发现问题

"期望效果"只写改成什么样（如"该行换到下一行开头"），不要指定实现手段。`

// ReviewV2 让模型看 v2 deck 的指定页截图。findings 由调用方从整份 deck 算好传入。
func ReviewV2(ctx context.Context, client *openai.Client, model string, d *Deck2, pages []int, onDelta func(string)) (*ReviewResult, error) {
	if onDelta == nil {
		onDelta = func(string) {}
	}
	var sel []Slide2
	picked := map[int]bool{}
	for _, p := range pages {
		if p >= 1 && p <= len(d.Slides) {
			sel = append(sel, d.Slides[p-1])
			picked[p] = true
		}
	}
	if len(sel) == 0 {
		return nil, fmt.Errorf("没有可审查的页")
	}

	var head strings.Builder
	head.WriteString(reviewPromptV2)
	fmt.Fprintf(&head, "\n\n你要看的页：第 %s 页（共 %d 张截图，deck 共 %d 页）\n", joinInts(pages), len(sel), len(d.Slides))
	hard := HardFindingsV2(d)
	head.WriteString("\n程序已判定的硬问题：\n")
	if len(hard) == 0 {
		head.WriteString("无")
	} else {
		head.WriteString(strings.Join(hard, "\n"))
	}
	head.WriteString("\n\n逐页量测（整份）：\n")
	head.WriteString(DigestV2(d))

	parts := []openai.ChatCompletionContentPartUnionParam{openai.TextContentPart(head.String())}
	for _, s := range sel {
		parts = append(parts, openai.TextContentPart(fmt.Sprintf("\n第 %d 页：", s.Index+1)))
		parts = append(parts, openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
			URL:    "data:image/png;base64," + base64.StdEncoding.EncodeToString(s.PNG),
			Detail: "auto",
		}))
	}
	_ = picked

	// 空报告重试一次（与 v1 Review 同一纪律：reasoning 吃光 completion 是概率行为）
	tryReview := func() (string, openai.CompletionUsage, error) {
		stream := client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
			Model:           openai.ChatModel(model),
			ReasoningEffort: shared.ReasoningEffortLow,
			Messages:        []openai.ChatCompletionMessageParamUnion{openai.UserMessage(parts)},
			MaxTokens:       openai.Int(8000),
			StreamOptions:   openai.ChatCompletionStreamOptionsParam{IncludeUsage: openai.Bool(true)},
		})
		acc := openai.ChatCompletionAccumulator{}
		var sb strings.Builder
		for stream.Next() {
			chunk := stream.Current()
			if !acc.AddChunk(chunk) {
				continue
			}
			if len(chunk.Choices) == 0 {
				continue
			}
			if delta := chunk.Choices[0].Delta; delta.Content != "" {
				sb.WriteString(delta.Content)
				onDelta(delta.Content)
			}
		}
		if err := stream.Err(); err != nil {
			stream.Close()
			return "", openai.CompletionUsage{}, err
		}
		stream.Close()
		if len(acc.Choices) == 0 || acc.Choices[0].FinishReason == "" {
			return "", openai.CompletionUsage{}, errors.New("看图审查流式响应不完整（没有 finish_reason）")
		}
		return strings.TrimSpace(sb.String()), acc.Usage, nil
	}

	report, usage, err := tryReview()
	if err == nil && report == "" {
		report, usage, err = tryReview()
	}
	if err != nil {
		return nil, err
	}
	if report == "" {
		return nil, fmt.Errorf("模型返回了空报告（completion=%d）", usage.CompletionTokens)
	}
	return &ReviewResult{Report: report, Prompt: head.String(), Images: len(sel), Usage: usage}, nil
}

func joinInts(xs []int) string {
	strs := make([]string, 0, len(xs))
	for _, x := range xs {
		strs = append(strs, fmt.Sprintf("%d", x))
	}
	return strings.Join(strs, "、")
}
