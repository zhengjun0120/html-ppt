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
	"encoding/json"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
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
      if (el.closest('.notes')) continue; // 讲稿观众不可见，不参与字号/密度判定
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
    var bottom = 0;
    for (var b = 0; b < sec.children.length; b++) {
      var c = sec.children[b];
      if (c.classList.contains('notes')) continue;
      var r = c.getBoundingClientRect();
      if (r.bottom > bottom) bottom = r.bottom;
    }
    // getBoundingClientRect 受 transform scale 影响，除回去换算回设计像素
    var scale = parseFloat(deck.style.getPropertyValue('--deck-scale')) || 1;
    var bottomGap = (h - bottom / scale);

    out.push({
      i: i,
      title: hEl ? (hEl.innerText||'').replace(/\s+/g,' ').trim().slice(0,30) : '',
      layout: sec.getAttribute('data-layout') || '',
      overflowY: overflowY,
      overflowX: overflowX,
      minFontPx: +minFont.toFixed(1),
      chars: cjk + latin,
      kids: kids,
      bottomGap: +bottomGap.toFixed(1)
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
// 依赖 base.css 的 @media print 规则：每页 .slide 一张纸（break-after:page）。
// 纸张 13.333×7.5 英寸 = 1920×1080 @96dpi，打印出的每页就是设计画布的等比映射。
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

	var pdf []byte
	if err := chromedp.Run(runCtx, chromedp.ActionFunc(func(c context.Context) error {
		b, _, e := page.PrintToPDF().
			WithPaperWidth(13.333).
			WithPaperHeight(7.5).
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
