package export

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func base64Of(s string) string { return base64.StdEncoding.EncodeToString([]byte(s)) }

// fakeDeck 只服务单文件导出路径：PDF/PNG 依赖 Chrome，不在单测范围。
type fakeDeck struct {
	html string
	dir  string
}

func (f *fakeDeck) PreviewHTML(uint, string) (string, error) { return f.html, nil }
func (f *fakeDeck) ExportDir(string) string                  { return filepath.Join(f.dir, "exports") }
func (f *fakeDeck) PageCount(string) int                     { return 8 }
func (f *fakeDeck) EnsureOwner(uint, string) error           { return nil }
func (f *fakeDeck) ThumbsDir(string) string                  { return filepath.Join(f.dir, "thumbs") }
func (f *fakeDeck) IndexPath(string) (string, error)         { return filepath.Join(f.dir, "index.html"), nil }

type fakeGrants struct{}

func (fakeGrants) Issue(uint, string) (string, error) { return "nonce", nil }

// 预览 HTML 的形态（v2.go 生成器输出）：三外链 css + 已内联 style + 尾部 runtime。
const sampleHTML = `<!DOCTYPE html><html lang="zh-CN"><head>
<meta charset="utf-8">
<link rel="stylesheet" href="/assets/deck-v2/fonts.css">
<link rel="stylesheet" href="/assets/deck-v2/base.css">
<link rel="stylesheet" href="/assets/deck-v2/animations.css">
<style>
.t{color:red}
</style>
<title>t</title>
</head><body><div class="deck"></div>
<script src="/assets/deck-v2/runtime.js"></script>
</body></html>`

const sampleFontsCSS = `/* head comment */
@font-face {
  font-family: 'Inter';
  src: url('/assets/deck-v2/fonts/inter-latin.woff2') format('woff2');
}
@font-face {
  font-family: 'JetBrains Mono';
  src: url('/assets/deck-v2/fonts/JetBrainsMapleMono-Regular.woff2') format('woff2');
}
`

func newSingleFileService(t *testing.T, html string) (*Service, string) {
	t.Helper()
	root := t.TempDir()
	av := filepath.Join(root, "assets", "deck-v2")
	if err := os.MkdirAll(filepath.Join(av, "fonts"), 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(name, content string) {
		if err := os.WriteFile(filepath.Join(av, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("base.css", ".slide{position:absolute}")
	write("animations.css", "@keyframes fade{}")
	write("runtime.js", "(function(){})()")
	write("fonts.css", sampleFontsCSS)
	if err := os.WriteFile(filepath.Join(av, "fonts", "inter-latin.woff2"), []byte("INTER-FONT"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(av, "fonts", "JetBrainsMapleMono-Regular.woff2"), []byte("CJK-FONT"), 0o644); err != nil {
		t.Fatal(err)
	}

	deckDir := filepath.Join(root, "deck-xxx")
	if err := os.MkdirAll(deckDir, 0o755); err != nil {
		t.Fatal(err)
	}
	svc := New(&fakeDeck{html: html, dir: deckDir}, "http://127.0.0.1:1", "", time.Second, fakeGrants{}, filepath.Join(root, "assets"))
	return svc, deckDir
}

func TestExportSingleFileSelfContained(t *testing.T) {
	svc, deckDir := newSingleFileService(t, sampleHTML)
	res, err := svc.Export(context.Background(), 1, "deck-xxx", "html")
	if err != nil {
		t.Fatalf("Export: %v", err)
	}
	r, ok := res.(Result)
	if !ok {
		t.Fatalf("Export 返回类型 %T，期望 Result", res)
	}
	raw, err := os.ReadFile(filepath.Join(deckDir, "exports", "deck.html"))
	if err != nil {
		t.Fatal(err)
	}
	out := string(raw)

	if strings.Contains(out, "/assets/") {
		t.Errorf("产物仍含 /assets 引用:\n%s", out)
	}
	for _, want := range []string{
		"<style>\n.slide{position:absolute}\n</style>",
		"<style>\n@keyframes fade{}\n</style>",
		"<script>\n(function(){})()\n</script>",
		/* 文件头注释保留 + Inter 转 data: */
		"/* head comment */",
		"url(data:font/woff2;base64," + base64Of("INTER-FONT") + ")",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("产物缺内联内容 %q", want)
		}
	}
	// CJK 全量等宽整块剔除：家族的 CJK 段不出现，字体文件内容不出现
	if strings.Contains(out, "JetBrainsMapleMono") || strings.Contains(out, "CJK-FONT") {
		t.Errorf("CJK 全量等宽应整块剔除")
	}
	if strings.Count(out, "data:font/woff2;base64,") != 1 {
		t.Errorf("data:font 应恰 1 处（跳过 CJK 全量）")
	}
	if r.Size != int64(len(out)) {
		t.Errorf("Size %d != 实际 %d", r.Size, len(out))
	}
}

func TestExportSingleFileRefusesWhenAssetMissing(t *testing.T) {
	svc, deckDir := newSingleFileService(t, sampleHTML)
	if err := os.Remove(filepath.Join(deckDir, "..", "assets", "deck-v2", "runtime.js")); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Export(context.Background(), 1, "deck-xxx", "html"); err == nil {
		t.Fatal("runtime.js 缺失时应拒绝导出")
	}
}

func TestExportSingleFileEscapesScriptCloseInRuntime(t *testing.T) {
	svc, deckDir := newSingleFileService(t, sampleHTML)
	av := filepath.Join(deckDir, "..", "assets", "deck-v2")
	// 注释/字符串里的 </script 是真实 runtime 的形态（现网 1 处，在注释里）
	if err := os.WriteFile(filepath.Join(av, "runtime.js"), []byte("/* x </script> y */ var s = '<\\/script>';"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Export(context.Background(), 1, "deck-xxx", "html"); err != nil {
		t.Fatalf("内联时应转义而非拒绝: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(deckDir, "exports", "deck.html"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "</script> y */") {
		t.Errorf("内联 JS 里的 </script 未被转义")
	}
	if !strings.Contains(string(raw), `<\/script> y */`) {
		t.Errorf("应转义为 <\\/script")
	}
}

func TestExportSingleFileRefusesDoubleEscapeCombo(t *testing.T) {
	svc, deckDir := newSingleFileService(t, sampleHTML)
	av := filepath.Join(deckDir, "..", "assets", "deck-v2")
	if err := os.WriteFile(filepath.Join(av, "runtime.js"), []byte("var doc = '<!--<script>';"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Export(context.Background(), 1, "deck-xxx", "html"); err == nil || !strings.Contains(err.Error(), "双逃逸") {
		t.Fatalf("<!-- + <script 组合应拒绝导出，got %v", err)
	}
}

func TestExportSingleFileRefusesStyleNotInlined(t *testing.T) {
	svc, _ := newSingleFileService(t, strings.Replace(sampleHTML, "<style>\n.t{color:red}\n</style>", `<link rel="stylesheet" href="style.css">`, 1))
	if _, err := svc.Export(context.Background(), 1, "deck-xxx", "html"); err == nil || !strings.Contains(err.Error(), "style.css 未内联") {
		t.Fatalf("style.css 未内联应拒绝导出，got %v", err)
	}
}
