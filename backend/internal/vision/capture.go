package vision

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// 启动窗口取最宽的预设画布（wide 1244×700）；真正的视口在页面加载后按 deck 自己的
// 画布再覆盖一次，见 Capture 里 SetDeviceMetricsOverride 那段。
const (
	DefaultCanvasW = 1244
	DefaultCanvasH = 700

	readyDelay  = 1500 * time.Millisecond //等待字体异步加载完成
	settleDelay = 300 * time.Millisecond  //每次调用 Reveal.slide(i) 后等待 300ms

	// 画布区间。画布值来自 deck 自己的主题块（模型写的 JSON），是个不可信输入：
	// 一个手滑写成 999999 的宽会让截图变成几百 MB 的 PNG（截图是逐像素编码的），
	// 而"离谱的画布"没有任何合法的用途。超出区间就退回默认画布。
	minCanvasW, maxCanvasW = 400, 4000
	minCanvasH, maxCanvasH = 300, 3000
)

// viewportFor 把 deck 声明的画布换算成截图视口，不可信的值退回默认画布。
func viewportFor(w, h int) (int, int) {
	if w < minCanvasW || w > maxCanvasW || h < minCanvasH || h > maxCanvasH {
		return DefaultCanvasW, DefaultCanvasH
	}
	return w, h
}

//图+数字
//
//**每个字段都要写 json tag，而且必须与 measureJS 返回的 key 逐字一致**：
//无 tag 时 encoding/json 按字段名忽略大小写匹配，FitScale 匹配不上 fit、
//Children 匹配不上 kids——不报错，字段静默为零（"子元素超 5 个"这条判定永远不会触发，
//而报告会照着 0 说"没问题"）。TestMeasureKeysMatchSlideStruct 守着这条。
type Slide struct {
	Index      int     `json:"i"`         //0 基
	Title      string  `json:"title"`     //该页标题 h1或h2的文本
	FitScale   float64 `json:"fit"`       //适配兜底一级缩放 1.0表示没有缩放
	InnerScale float64 `json:"innerK"`    //二级缩放
	EffScale   float64 `json:"effScale"`  //两者相乘
	Overflow   float64 `json:"overflow"`  // >1.02 表示内容超出画布
	Children   int     `json:"kids"`      //section 直接子元素个数 框架要求<=5
	Chars      int     `json:"chars"`     //整页可见字数 预算<=200
	MinFontPx  float64 `json:"minFontPx"` //页面上最小的字号：低于31px 就是低于 24pt 投影下限
	Layout     string  `json:"layout"`    //版式指纹，用来查版式重复
	CanvasW    int     `json:"canvasW"`   //画布宽：主题可以换 canvas，所以由页面回报
	CanvasH    int     `json:"canvasH"`   //画布高
	PNG        []byte  `json:"-"`         //图片不进 JSON（否则一序列化就拖着几百 KB）
}

type Deck struct {
	Slides  []Slide
	CanvasW int
	CanvasH int
}

type Options struct {
	URL        string
	ChromePath string
	Timeout    time.Duration

	// MeasureOnly：只量测、不截图。
	//
	// 量测是浏览器自己算的（几秒、零模型开销），截图才是贵的那一半（PNG 编码 +
	// 图片 token）。write_deck 之后的自动体检走这条路：先把数字摆出来，"要不要花
	// 时间看画面"由 agent 点名页号决定。
	MeasureOnly bool

	// Shoot：只对这几页截图，**0 基下标**（与 Slide.Index 同一套口径，调用方不用
	// 在这里再做一次 1→0 换算——换算写错既不报错也不改变图片本身的样子，
	// 只会安静地拍错页，而拍错页的报告和拍对了看着一样可信）。
	// len(Shoot)==0（nil 和空切片一样）表示全拍。
	Shoot []int
}

// 这一页要不要截图。单独拆出来是为了能脱离浏览器单测：
// 0 基/1 基写错不会报错，只会拍错页。
func shootPage(opt Options, i int) bool {
	if opt.MeasureOnly {
		return false
	}
	if len(opt.Shoot) == 0 {
		return true
	}
	for _, p := range opt.Shoot {
		if p == i {
			return true
		}
	}
	return false
}

//把页面调到 某一页静止呈现的状态
//
//两条都必须做，各踩过一次，而且两次量测数字都是对的、只有看图才发现图是坏的：
//  - 关翻页动画：reveal 的 slide 过渡要 ~500ms，紧接着截图会拍到动画中间态，
//    画面里的内容被切在边缘，看着像"版式崩了"，其实只是拍早了。
//  - 展开 .fragment：分步渐显的条目初始是隐藏的，不展开会把一页拍成"缺了一半内容"。
const prepareJS = `(function(i){
  if (typeof Reveal === 'undefined') return 'no Reveal';
  Reveal.configure({transition:'none', backgroundTransition:'none'});
  Reveal.slide(i);
  var fr = document.querySelectorAll('.slides > section .fragment');
  for (var n = 0; n < fr.length; n++) fr[n].classList.add('visible');
  return 'ok';
})(%d)`

// measureJS 只量、不改变页面状态（页面定位已经在 prepareJS 里做完）。
// 判据刻意与 init.js 的兜底逻辑同口径（纵向比画布高、横向比自身盒宽），
// 否则会出现"兜底认为放得下、审查认为溢出了"这种自相矛盾。
//
// 返回对象刻意写成"一行一个 key"：契约测试要从这段 JS 里把 key 抠出来跟结构体的 tag
// 对比，一行一个才能用 ^\s*(\w+): 精确匹配，不会被值里的三元运算符 ? a : b 误伤。
const measureJS = `(function(i){
  if (typeof Reveal === 'undefined') return JSON.stringify({error:'no Reveal'});
  var sec = Reveal.getSlides()[i];
  if (!sec) return JSON.stringify({error:'no slide'});
  var cfg = Reveal.getConfig() || {};
  var canvasH = cfg.height || 700, canvasW = cfg.width || 1244;

  var fit = sec.style.fontSize ? parseFloat(sec.style.fontSize) / 100 : 1;
  if (!isFinite(fit) || fit <= 0) fit = 1;
  var k = 1;
  var inner = sec.querySelector(':scope > .ppt-fit-inner');
  if (inner) {
    var m = /scale\(([0-9.]+)\)/.exec(inner.style.transform || '');
    if (m) k = parseFloat(m[1]) || 1;
  }

  var availH = Math.max(1, canvasH - 8);
  var overflow = Math.max(sec.scrollHeight / availH, sec.scrollWidth / Math.max(1, sec.clientWidth));

  var kids = 0;
  for (var n = 0; n < sec.children.length; n++) {
    var c = sec.children[n];
    kids += c.classList.contains('ppt-fit-inner') ? c.children.length : 1;
  }

  // 版式指纹：直接子元素的 tag+class，去掉页面家具（它们每页都可能出现，不参与轮廓判断）
  var FURN = {'page-no':1,'outline':1,'kicker':1,'rule':1,'footer':1,'footnote':1,'muted':1,
              'sub':1,'label':1,'term':1,'unit':1,'num':1,'fragment':1,'ppt-fit-inner':1,'val':1};
  var parts = [], scan = inner ? inner.children : sec.children;
  for (var j = 0; j < scan.length; j++) {
    var e = scan[j], cls = [];
    for (var q = 0; q < e.classList.length; q++) if (!FURN[e.classList[q]]) cls.push(e.classList[q]);
    cls.sort();
    parts.push(e.tagName.toLowerCase() + (cls.length ? '.' + cls.join('.') : ''));
  }

  var text = (sec.innerText || '').replace(/\s+/g, '');
  var minFont = 999, all = sec.querySelectorAll('*');
  for (var a = 0; a < all.length; a++) {
    var fs = parseFloat(getComputedStyle(all[a]).fontSize);
    if (fs > 0 && fs < minFont) minFont = fs;
  }

  var h = sec.querySelector('h1, h2');
  return JSON.stringify({
    i: i,
    title: h ? (h.innerText||'').replace(/\s+/g,' ').trim().slice(0,30) : '',
    fit: +fit.toFixed(3),
    innerK: +k.toFixed(3),
    effScale: +(fit*k).toFixed(3),
    overflow: +overflow.toFixed(3),
    kids: kids,
    chars: text.length,
    minFontPx: +minFont.toFixed(1),
    layout: parts.join(' + '),
    canvasW: canvasW,
    canvasH: canvasH
  });
})(%d)`

const slideCountJS = `(function(){
  if (typeof Reveal === 'undefined') return -1;
  return Reveal.getSlides().filter(function(s){return s.tagName==='SECTION';}).length;
})()`

// deck 自己声明的逻辑画布（主题块里的 canvas 决定它）。读 Reveal.getConfig() 而不是
// 去解析页面里的主题 JSON：那才是真正生效的那个值（主题块缺失时 init.js 有兜底）。
const canvasJS = `(function(){
  if (typeof Reveal === 'undefined') return JSON.stringify({w:0,h:0});
  var cfg = Reveal.getConfig() || {};
  return JSON.stringify({w: cfg.width||0, h: cfg.height||0});
})()`

// 让 reveal 重算缩放。视口变了页面会收到 resize、reveal 自己也会重排，但那一步在什么时候
// 完成取决于它内部的调度，而紧接着我们就要量测和截图——不留这个竞态。这里显式调一次，
// 把重排变成同步的，代价只是一次重排（同一个视口不会重复触发，见下面只在换视口时才调）。
const relayoutJS = `(function(){
  if (typeof Reveal !== 'undefined' && Reveal.layout) Reveal.layout();
  return 'ok';
})()`

type canvasSize struct {
	W int `json:"w"`
	H int `json:"h"`
}

func Capture(ctx context.Context, opt Options) (*Deck, error) {
	if opt.Timeout == 0 {
		opt.Timeout = 120 * time.Second
	}
	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.Flag("hide-scrollbars", true),
		chromedp.WindowSize(DefaultCanvasW, DefaultCanvasH),
	)
	if opt.ChromePath != "" {
		allocOpts = append(allocOpts, chromedp.ExecPath(opt.ChromePath))
	}
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, allocOpts...)
	defer cancelAlloc()

	//关闭 chromedp日志
	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx, chromedp.WithLogf(func(string, ...any) {}))
	defer cancelBrowser()

	runCtx, cancelTimeout := context.WithTimeout(browserCtx, opt.Timeout)
	defer cancelTimeout()

	if err := chromedp.Run(runCtx,
		emulation.SetDeviceMetricsOverride(DefaultCanvasW, DefaultCanvasH, 1, false), //缩放1 非移动
		chromedp.Navigate(opt.URL),
		chromedp.WaitReady("body", chromedp.ByQuery),
	); err != nil {
		return nil, fmt.Errorf("打开 deck 失败 err:%w", err)
	}
	time.Sleep(readyDelay)

	var n int
	if err := chromedp.Run(runCtx, chromedp.Evaluate(slideCountJS, &n)); err != nil {
		return nil, fmt.Errorf("读取页数失败 err:%w", err)
	}
	if n <= 0 {
		// 这里原来漏了参数：Errorf 的 %d 没给值，日志会写成 "%!d(MISSING)"，
		// 而"量到几页"正是报错时最想知道的那个数——和 web_search.go 里那个
		// %w 漏参数是同一类错（那次是丢了失败原因，这次是丢了实际页数）
		return nil, fmt.Errorf("量到 %d 页：Reveal 没起来（js报错？）", n)
	}

	// 视口跟随 deck 自己的画布。
	//
	// 为什么不能让视口写死 1244：**拍出来的画面得是"观众看到的那一张"**。画布比例和
	// 视口比例不一致时，reveal 会按短边缩放并把内容居中，多出来的部分就是空白——
	// 实测（classic 4:3 画布、模拟 16:9 屏）视口用 1244 拍，两侧各留 149px 空白；
	// 视口换成画布自己的 933，空白只剩 reveal 那 2% 边距。那 149px 是**只在截图里存在**的
	// 版面问题，而视觉模型会老老实实把它报成"内容没铺满、两侧大片空白"。
	// canvasJS 返回的是**字符串**（JSON.stringify 的结果），所以要先取字符串再解——
	// 直接往 struct 里解会报 "cannot unmarshal string into Go value of type canvasSize"。
	// 与 measureJS 同一套：Evaluate 拿到的是 JS 表达式的返回值，不是它代表的那个对象。
	var rawCanvas string
	if err := chromedp.Run(runCtx, chromedp.Evaluate(canvasJS, &rawCanvas)); err != nil {
		return nil, fmt.Errorf("读画布尺寸失败 err:%w", err)
	}
	var cs canvasSize
	if err := json.Unmarshal([]byte(rawCanvas), &cs); err != nil {
		return nil, fmt.Errorf("画布尺寸解析失败(%s) err:%w", rawCanvas, err)
	}
	vw, vh := viewportFor(cs.W, cs.H)
	if vw != DefaultCanvasW || vh != DefaultCanvasH {
		// 只有真的要换才覆盖：默认画布那一次覆盖是白走一趟 resize + 重排
		if err := chromedp.Run(runCtx,
			emulation.SetDeviceMetricsOverride(int64(vw), int64(vh), 1, false),
			chromedp.Evaluate(relayoutJS, nil),
		); err != nil {
			return nil, fmt.Errorf("按画布 %d*%d 设视口失败 err:%w", vw, vh, err)
		}
		time.Sleep(settleDelay) //重排之后再量：量的是重排后的结果
	}

	deck := &Deck{CanvasW: vw, CanvasH: vh}
	for i := 0; i < n; i++ {
		var prep string
		if err := chromedp.Run(runCtx, chromedp.Evaluate(fmt.Sprintf(prepareJS, i), &prep)); err != nil {
			return nil, fmt.Errorf("第 %d 页定位失败 err:%w", i+1, err)
		}
		time.Sleep(settleDelay)

		var raw string
		if err := chromedp.Run(runCtx, chromedp.Evaluate(fmt.Sprintf(measureJS, i), &raw)); err != nil {
			return nil, fmt.Errorf("第 %d 页量测失败 err:%w", i+1, err)
		}
		var s Slide
		//这里拿到结果
		if err := json.Unmarshal([]byte(raw), &s); err != nil {
			return nil, fmt.Errorf("第 %d 页量测结果解析失败(%s) err:%w", i+1, raw, err)
		}
		s.Index = i //以循环下标为准：页面回报的值一旦错位，报告就会指错页
		if s.CanvasW > 0 && s.CanvasH > 0 {
			// 报告里写的画布是 deck 自己声明的那个（离谱的值也要如实报出来——
			// 那是 deck 的毛病，藏起来反而没人会去改）
			deck.CanvasW, deck.CanvasH = s.CanvasW, s.CanvasH
		}

		//整视口截图。**不要用元素截图**：reveal 会给 .slides 加 transform，
		//元素截图走的是未变换的盒子模型，裁出来的区域会整体偏移（实测左侧内容被切掉）。
		//视口已经固定成画布尺寸，整视口截图就是"观众实际看到的那一张"。
		//没被点名的页连拍都不拍：PNG 编码是这段里最费时间的一步，而没人看的图
		//只会让 write_deck 之后的自动体检从几秒变成几十秒。
		if shootPage(opt, i) {
			var buf []byte
			if err := chromedp.Run(runCtx, chromedp.ActionFunc(func(c context.Context) error {
				b, e := page.CaptureScreenshot().WithFormat(page.CaptureScreenshotFormatPng).Do(c)
				buf = b
				return e
			})); err != nil {
				return nil, fmt.Errorf("第 %d 页截图失败 err:%w", i+1, err)
			}
			s.PNG = buf
		}
		deck.Slides = append(deck.Slides, s)
	}
	return deck, nil
}