package usertpl

import "testing"

// 复现 v8 的 changed 多了 style.css 的疑点：上一版 style=stub，磁盘 style=stub，
// 只改 index.html → changed 应只含 index.html。
func TestChangedOnlyIndex(t *testing.T) {
	s, _, id, uid := newHistService(t, "chgidx")
	_ = s.recordVersionUT(id, OpFork, "起点") // v1 基线（五件套全非空）
	// 模拟 v7：style 改成 stub
	if err := s.SaveStyleCSS(uid, id, "/* 下架后可改 */ .x{color:red}"); err != nil {
		t.Fatalf("SaveStyleCSS: %v", err)
	}
	// 模拟 v8：编辑弹窗只改 index.html（style 保持 stub）
	newIndex := "<body class=\"tpl-x\"><div class=\"deck\"><section class=\"slide\">v8</section></div><script src=\"/assets/deck-v2/runtime.js\"></script></body>"
	if err := s.SaveIndexHTML(uid, id, newIndex); err != nil {
		t.Fatalf("SaveIndexHTML: %v", err)
	}
	versions, _ := s.ListUTVersions(uid, id)
	last := versions[0]
	if last.Detail != "手动编辑 index.html" {
		t.Fatalf("detail=%q", last.Detail)
	}
	if len(last.Changed) != 1 || last.Changed[0] != "index.html" {
		t.Errorf("changed 应为 [index.html]，得到 %v", last.Changed)
	}
}
