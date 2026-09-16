package vision

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// 视口跟随 deck 自己的画布：合法值原样用，离谱的值退回默认画布。
//
// 画布值来自 deck 的主题块（模型写的 JSON），是这段代码里唯一不可信的输入。
// 一个手滑写成 999999 的宽会让截图变成几百 MB 的 PNG（截图是逐像素编码的），
// 而这种错从产出的图上**看不出来**——图永远是有内容的，只是内容特别小。
func TestViewportForFallsBackOnInsaneCanvas(t *testing.T) {
	// 三个预设原样通过：它们就是这套系统真正会出现的值
	for _, c := range [][2]int{{1244, 700}, {960, 700}, {933, 700}} {
		if w, h := viewportFor(c[0], c[1]); w != c[0] || h != c[1] {
			t.Errorf("预设画布 %v 应当原样当视口，实际 %d*%d", c, w, h)
		}
	}
	// Reveal 没起来时 canvasJS 返回 0*0，走到这里也不能崩；其余是离谱值
	for _, c := range [][2]int{{0, 0}, {999999, 700}, {1244, 999999}, {-1, -1}, {100, 100}} {
		if w, h := viewportFor(c[0], c[1]); w != DefaultCanvasW || h != DefaultCanvasH {
			t.Errorf("画布 %v 不合理，应当退回默认 %d*%d，实际 %d*%d",
				c, DefaultCanvasW, DefaultCanvasH, w, h)
		}
	}
}

// JS 返回的 key 必须与 Slide 的 json tag 逐字一致，两个方向都查：
//   - 结构体有、JS 没有 → 那个字段永远是零值（正是 "kids 匹配不上 Children" 的形态：
//     "直接子元素超过 5 个"这条判定永远不会触发，而报告会照着 0 说"没问题"）
//   - JS 有、结构体没有 → 量出来的数字被丢掉（原来是 canvasW/canvasH）
//
// 为什么值得单独守：这类错**不报错、不 panic**，只是数字悄悄变成 0。手工对照两张表
// 是不可靠的——第一版就没对上，而且是人肉 review 才发现的（编译器一个字都不会说）。
func TestMeasureKeysMatchSlideStruct(t *testing.T) {
	// 取 measureJS 里**最后一个** JSON.stringify({...})：前面还有两个用于错误路径的
	// {error:...}，抠错了这条测试就白守了。
	idx := strings.LastIndex(measureJS, "JSON.stringify({")
	if idx < 0 {
		t.Fatal("measureJS 里找不到 JSON.stringify({")
	}
	body := measureJS[idx:]
	if end := strings.Index(body, "\n  });"); end > 0 {
		body = body[:end]
	}

	got := map[string]bool{}
	for _, line := range strings.Split(body, "\n") {
		// 只认"行首（允许缩进）就是 key"的行：值里的三元运算符 ? a : b 在行中间，不会误伤
		if m := regexp.MustCompile(`^\s*([A-Za-z_]\w*)\s*:`).FindStringSubmatch(line); m != nil {
			got[m[1]] = true
		}
	}
	if len(got) < 8 {
		t.Fatalf("只从 measureJS 里解析出 %d 个 key，解析逻辑不对：%v", len(got), got)
	}

	want := map[string]bool{}
	typ := reflect.TypeOf(Slide{})
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		want[tag] = true
		if !got[tag] {
			t.Errorf("Slide.%s 的 tag 是 %q，但 measureJS 没返回这个 key——这个字段永远是零值",
				typ.Field(i).Name, tag)
		}
	}
	for k := range got {
		if !want[k] {
			t.Errorf("measureJS 返回了 %q，但 Slide 里没有对应的 json tag——这个数字被丢掉了", k)
		}
	}
	t.Logf("JS 返回 %d 个 key，结构体 %d 个 tag，逐一对上", len(got), len(want))
}

// "这一页要不要截图"是三态：只量测（一页都不拍）、没点名（全拍）、点名（只拍点到的几页）。
//
// 值得单独测，是因为这里的错**不报错也不 panic**：把 1 基页号写进 Shoot 会安静地
// 拍错页——报告里照样写着"第 N 页"，看图看得再仔细也看不出看的是别的一页。
func TestShootPageSelection(t *testing.T) {
	// 三页：没点名（Shoot 为 nil）= 全拍
	for i := 0; i < 3; i++ {
		if !shootPage(Options{}, i) {
			t.Errorf("没点名时第 %d 页也该拍", i+1)
		}
	}
	// 空切片与 nil 同义，调用方少记一条规则（"空 vs nil"这种区别只在写的时候记住、
	// 读的时候永远记不住）
	if !shootPage(Options{Shoot: []int{}}, 1) {
		t.Error("空 Shoot 应当与 nil 同义（全拍），否则调用方要记住空切片是「一页都不拍」")
	}

	// 点名：Shoot 是 0 基，与 Slide.Index 同口径
	only := Options{Shoot: []int{0, 2}}
	if !shootPage(only, 0) {
		t.Error("Shoot=[0,2] 应当拍第 1 页")
	}
	if shootPage(only, 1) {
		t.Error("Shoot=[0,2] 不应当拍第 2 页——会命中这里只有一种原因：有人把 1 基页号塞进了 Shoot")
	}
	if !shootPage(only, 2) {
		t.Error("Shoot=[0,2] 应当拍第 3 页")
	}

	// 只量测：write_deck 之后那次免费体检走这条，一页都不拍
	if shootPage(Options{MeasureOnly: true, Shoot: []int{0}}, 0) {
		t.Error("MeasureOnly 时一页都不该拍：那次体检之所以免费，就是靠不编码 PNG 换来的")
	}
	if shootPage(Options{MeasureOnly: true}, 0) {
		t.Error("MeasureOnly 且没点名时也不该拍")
	}

	// 越界页号不命中也不崩：capture 这里不知道这份 deck 有几页，
	// 范围检查在 Deck.Select 里做（那里才有页数）
	if shootPage(Options{Shoot: []int{99}}, 0) {
		t.Error("越界页号不该命中任何页")
	}
}
