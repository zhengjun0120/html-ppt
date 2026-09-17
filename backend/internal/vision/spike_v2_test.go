package vision

// P0 spike：deck-v2 渲染底座的三个关键假设，一次性验证。
//
// 结论落在 deck-v2/PATCHES.md（若有补丁）与 docs/refactor-plan.md 的 P0 记录里。
// 需要 Chrome，CI/短模式默认跳过；本地跑：SPIKE_V2=1 go test ./internal/vision/ -run Spike -v

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

func spikeChrome(t *testing.T) string {
	t.Helper()
	if testing.Short() {
		t.Skip("short 模式跳过 spike")
	}
	if os.Getenv("SPIKE_V2") == "" {
		t.Skip("设 SPIKE_V2=1 以运行渲染 spike")
	}
	return "" // ChromePath：空 = chromedp 默认查找
}

// spikeServer 以生产同构的路径暴露仓库根：
//
//	/assets/...  → backend/web/assets 需要的文件在仓库根下 web/... 不对——
//	模板引用的是 /assets/deck-v2/*，而磁盘上是 backend/web/assets/deck-v2/*。
//	所以这里做一层路径映射：/assets/ → backend/web/assets/，/templates/ → backend/templates/。
type spikeServer struct {
	srv *httptest.Server
}

func newSpikeServer(t *testing.T) *spikeServer {
	t.Helper()
	backendRoot, _ := filepath.Abs(filepath.Join("..", ".."))
	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(filepath.Join(backendRoot, "web", "assets")))))
	mux.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir(filepath.Join(backendRoot, "templates")))))
	mux.HandleFunc("/spike/host.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<!DOCTYPE html><html><body style="margin:0">
<iframe id="f" src="/templates/tech-sharing/index.html" style="width:1920px;height:1080px;border:0"></iframe>
</body></html>`))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return &spikeServer{srv: srv}
}

func newBrowser(t *testing.T, chromePath string) (context.Context, context.CancelFunc) {
	t.Helper()
	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.Flag("hide-scrollbars", true),
		chromedp.WindowSize(V2DefaultW, V2DefaultH),
	)
	if chromePath != "" {
		allocOpts = append(allocOpts, chromedp.ExecPath(chromePath))
	}
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), allocOpts...)
	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx, chromedp.WithLogf(func(string, ...any) {}))
	ctx, cancelTimeout := context.WithTimeout(browserCtx, 90*time.Second)
	return ctx, func() { cancelTimeout(); cancelBrowser(); cancelAlloc() }
}

// SpikeV2Render 主 spike：
//  1. 假设"非激活页仍参与布局，可一次量完" → 量测数字齐且合理；
//  2. 假设"hash 深链能激活单页用于截图" → 点名页 PNG 非空、未点名页为空；
//  3. 假设"PrintToPDF 出 16:9 整页" → PDF 头合法、体积合理。
func TestSpikeV2Render(t *testing.T) {
	chrome := spikeChrome(t)
	srv := newSpikeServer(t)
	url := srv.srv.URL + "/templates/tech-sharing/index.html"

	// —— 量测全部页 ——
	deck, err := CaptureV2(t.Context(), OptionsV2{URL: url, ChromePath: chrome, MeasureOnly: true})
	if err != nil {
		t.Fatalf("量测失败: %v", err)
	}
	if len(deck.Slides) != 8 {
		t.Errorf("tech-sharing demo 应为 8 页，量到 %d", len(deck.Slides))
	}
	if deck.CanvasW != 1920 || deck.CanvasH != 1080 {
		t.Errorf("画布 %d*%d，应为 1920*1080", deck.CanvasW, deck.CanvasH)
	}
	for _, s := range deck.Slides {
		if s.Layout == "" {
			t.Errorf("第 %d 页缺 data-layout", s.Index+1)
		}
		if s.OverflowY || s.OverflowX {
			t.Errorf("第 %d 页（%s / %s）溢出：demo 页不该溢出", s.Index+1, s.Title, s.Layout)
		}
		if s.MinFontPx <= 0 || s.MinFontPx > 400 {
			t.Errorf("第 %d 页最小字号异常: %v", s.Index+1, s.MinFontPx)
		}
	}

	// —— 单页截图（hash 激活） ——
	withShot, err := CaptureV2(t.Context(), OptionsV2{URL: url, ChromePath: chrome, Shoot: []int{0, 7}})
	if err != nil {
		t.Fatalf("截图失败: %v", err)
	}
	if len(withShot.Slides[0].PNG) < 10_000 {
		t.Errorf("第 1 页截图过小（%d bytes），疑似空白", len(withShot.Slides[0].PNG))
	}
	if len(withShot.Slides[7].PNG) < 10_000 {
		t.Errorf("第 8 页截图过小（%d bytes），疑似空白", len(withShot.Slides[7].PNG))
	}
	if len(withShot.Slides[3].PNG) != 0 {
		t.Error("未点名的页不应有截图")
	}

	// —— PrintToPDF ——
	pdf, err := PrintPDF2(t.Context(), url, chrome, 60*time.Second)
	if err != nil {
		t.Fatalf("PDF 失败: %v", err)
	}
	if !strings.HasPrefix(string(pdf[:5]), "%PDF-") {
		t.Error("PDF 头不合法")
	}
	if len(pdf) < 50_000 {
		t.Errorf("PDF 体积 %d bytes 过小，疑似空白页", len(pdf))
	}

	// 产物落 tmp 供人工核对（视觉核对：PNG 内容必须是真实的幻灯片画面）
	outDir := filepath.Join("..", "..", "..", "tmp", "spike-v2")
	_ = os.MkdirAll(outDir, 0o755)
	_ = os.WriteFile(filepath.Join(outDir, "page-1.png"), withShot.Slides[0].PNG, 0o644)
	_ = os.WriteFile(filepath.Join(outDir, "page-8.png"), withShot.Slides[7].PNG, 0o644)
	_ = os.WriteFile(filepath.Join(outDir, "deck.pdf"), pdf, 0o644)
	t.Logf("spike 产物已写入 %s", outDir)
}


// SpikeV2Iframe 验证 runtime.js 在 iframe 内正常初始化（PreviewPane 的前提）。
func TestSpikeV2Iframe(t *testing.T) {
	chrome := spikeChrome(t)
	srv := newSpikeServer(t)

	ctx, cancel := newBrowser(t, chrome)
	defer cancel()

	checkJS := `(function(){
  var f = document.getElementById('f');
  if (!f) return JSON.stringify({error:'no iframe'});
  try {
    var doc = f.contentDocument || f.contentWindow.document;
    var deck = doc.querySelector('.deck');
    if (!deck) return JSON.stringify({error:'no deck in iframe'});
    var active = doc.querySelectorAll('.slide.is-active').length;
    var scaled = deck.style.getPropertyValue('--deck-scale');
    var hash = f.contentWindow.location.hash;
    var real = doc.querySelectorAll('.deck > .slide');
    var realActive = doc.querySelectorAll('.deck > .slide.is-active').length;
    return JSON.stringify({slides: real.length, active: realActive, scale: scaled, hash: hash});
  } catch(e) { return JSON.stringify({error: String(e)}); }
})()`

	if err := chromedp.Run(ctx,
		emulation.SetDeviceMetricsOverride(V2DefaultW, V2DefaultH, 1, false),
		chromedp.Navigate(srv.srv.URL+"/spike/host.html"),
		chromedp.WaitReady("body", chromedp.ByQuery),
	); err != nil {
		t.Fatalf("打开宿主页失败: %v", err)
	}
	time.Sleep(v2ReadyDelay + 500*time.Millisecond)

	var raw string
	if err := chromedp.Run(ctx, chromedp.Evaluate(checkJS, &raw)); err != nil {
		t.Fatalf("iframe 检查执行失败: %v", err)
	}
	if strings.Contains(raw, "error") {
		t.Fatalf("iframe 内 runtime 未正常初始化: %s", raw)
	}
	if !strings.Contains(raw, `"active":1`) {
		t.Fatalf("iframe 内没有激活页（runtime go() 未执行？）: %s", raw)
	}
	t.Logf("iframe spike 结果: %s", raw)
}
