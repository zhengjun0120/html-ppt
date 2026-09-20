package deck

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

// versionDiff 必须从 v2 的 .json 快照 bundle 里取 index_html。
// 回归背景：这里曾硬读 <version>.html（v1 形状），v2 deck 的历史目录里
// 根本没有 .html，read_history_diff 工具每次必失败。
func TestVersionDiffReadsV2JsonSnapshot(t *testing.T) {
	dir := t.TempDir()
	s := New(dir, filepath.Join(dir, "assets"), nil)

	deckID := "deck-0037"
	// 当前 index：两页（v2 页面挂在 .deck 下）
	current := `<!DOCTYPE html><html><body><div class="deck">` +
		`<section data-id="s1"><h1>封面</h1></section>` +
		`<section data-id="s2"><h1>目录</h1></section>` +
		`</div></body></html>`
	// v000001 快照：只有一页
	oldHTML := `<!DOCTYPE html><html><body><div class="deck">` +
		`<section data-id="s1"><h1>封面</h1></section>` +
		`</div></body></html>`

	if err := os.MkdirAll(filepath.Join(dir, "decks", deckID, "history"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "decks", deckID, "index.html"), []byte(current), 0o644); err != nil {
		t.Fatal(err)
	}
	bundle := `{"index_html":` + strconv.Quote(oldHTML) + `}`
	if err := os.WriteFile(filepath.Join(dir, "decks", deckID, "history", "v000001.json"), []byte(bundle), 0o644); err != nil {
		t.Fatal(err)
	}
	idx := `{"next_seq":2,"versions":[{"version":"v000001","time":1,"operation":"run","detail":"初版","slides":1}]}`
	if err := os.WriteFile(filepath.Join(dir, "decks", deckID, "history", "index.json"), []byte(idx), 0o644); err != nil {
		t.Fatal(err)
	}

	d, err := s.versionDiff(deckID, "v000001", "")
	if err != nil {
		t.Fatalf("versionDiff 失败: %v", err)
	}
	if !d.Changed || len(d.Added) != 1 || d.Added[0].SlideID != "s2" {
		t.Fatalf("diff 结果不符合预期: %+v", d)
	}
	if d.Unchanged != 1 {
		t.Fatalf("unchanged 应为 1: %+v", d)
	}
}

// 缺 data-id 的 section 用序号兜底参与 diff（整页重写没带 id 时不再静默消失）
func TestMapSlidesSynthesizesMissingIDs(t *testing.T) {
	html := `<!DOCTYPE html><html><body><div class="deck">` +
		`<section data-id="s1"><h1>封面</h1></section>` +
		`<section><h1>无 id 页</h1></section>` +
		`</div></body></html>`

	m, order := mapSlides(html)
	if len(order) != 2 {
		t.Fatalf("应识别出 2 页，实际 %d（%v）", len(order), order)
	}
	if snap, ok := m["#2"]; !ok || snap.title != "无 id 页" {
		t.Fatalf("缺 id 的页应有兜底 id #2: %+v", m)
	}
}
