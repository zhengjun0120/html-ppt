package vision

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
)

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
