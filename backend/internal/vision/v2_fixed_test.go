package vision

// fixed 装饰排除的回归守卫。
//
// .slide 自带 transform（翻页位移/画布缩放），因此它是 slide 内 position:fixed
// 子孙的包含块——固定装饰按 slide 坐标参与滚动区计算，left:2600px 的角标能把
// 干净页的 scrollWidth 顶到 2800，溢出判定全页误报（观众看到的是被
// overflow:hidden 裁掉的画面，无感）。v2MeasureAllJS 量测期间摘出 fixed 子孙，
// 本测试守住：带 fixed 装饰的页不得被误判溢出。
//
// 需要 Chrome，short 模式跳过；本地：go test ./internal/vision/ -run FixedDecor -v
import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestV2MeasureExcludesFixedDecor(t *testing.T) {
	if testing.Short() {
		t.Skip("short 模式跳过（需要 Chrome）")
	}
	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(
		http.Dir(filepath.Join("..", "..", "web", "assets")))))
	mux.HandleFunc("/fixed-fixture.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<!doctype html><html><head>
<link rel="stylesheet" href="/assets/deck-v2/base.css"></head><body>
<div class="deck" data-w="1920" data-h="1080">
  <section class="slide" data-layout="cover"><h1 class="title">干净页</h1><p>没有 fixed 装饰</p></section>
  <section class="slide" data-layout="split"><h1 class="title">带 fixed 装饰的页</h1>
    <p>内容本身不溢出，角标在画布外（left:2600px）</p>
    <div style="position:fixed;left:2600px;top:40px;width:200px;height:60px">fixed 角标</div>
  </section>
</div></body></html>`))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	d, err := CaptureV2(t.Context(), OptionsV2{URL: srv.URL + "/fixed-fixture.html", MeasureOnly: true})
	if err != nil {
		t.Fatalf("量测失败: %v", err)
	}
	if len(d.Slides) != 2 {
		t.Fatalf("应量到 2 页，实际 %d", len(d.Slides))
	}
	for _, s := range d.Slides {
		if s.OverflowX || s.OverflowY {
			t.Errorf("第 %d 页（%s）被误判溢出——fixed 装饰排除失效（修复前该页 scrollWidth 被 fixed 角标顶到 2800）",
				s.Index+1, s.Layout)
		}
		if s.MinFontPx <= 0 || s.MinFontPx > 400 {
			t.Errorf("第 %d 页最小字号异常: %v", s.Index+1, s.MinFontPx)
		}
	}
}
