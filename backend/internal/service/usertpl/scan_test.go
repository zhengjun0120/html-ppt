package usertpl

// scanCSS/scanHTML 安全预检（scan.go）与对话整文件重写工具的回归：
// docs/user-template-history-plan.md §3.3/§6。
// 覆盖：CSS 黑名单（url( 只放行 data: 内联）、HTML 黑名单与结构前提、
// write_style/write_demo 的预检拦截与落盘、published 全路径锁（D2）。

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"html-ppt/backend/internal/store"
)

func TestScanCSS(t *testing.T) {
	if err := scanCSS(".tpl-x{--accent:#000}"); err != nil {
		t.Errorf("合法 css 不应报错: %v", err)
	}
	// data: 内联放行（brutalist-bold 的纸纹底即此用法，fork 副本要能过门禁）
	if err := scanCSS(`.a{background:var(--paper) url("data:image/svg+xml,<svg/>") repeat}`); err != nil {
		t.Errorf("data: 内联应放行: %v", err)
	}
	for _, bad := range []string{
		`.a{background:url(https://evil.com/x.png)}`,
		`.a{background:url('//evil.com/x.png')}`,
		`.a{background:url(/local/x.png)}`,
		`@import url("https://evil.com/x.css")`,
		`.a{width:expression(alert(1))}`,
		`.a{-moz-binding:url(https://evil.com/x.xml)}`,
		`body{behavior:url(https://evil.com/x.htc)}`,
	} {
		if err := scanCSS(bad); err == nil {
			t.Errorf("应拒绝: %s", bad)
		}
	}
}

func TestScanHTML(t *testing.T) {
	good := `<html><head><link rel="stylesheet" href="/assets/deck-v2/base.css"><link rel="stylesheet" href="style.css"></head><body class="tpl-x"><div class="deck"><section class="slide">a</section></div><script src="/assets/deck-v2/runtime.js"></script></body></html>`
	if err := scanHTML(good); err != nil {
		t.Errorf("合法 demo 不应报错: %v", err)
	}
	for name, bad := range map[string]string{
		"第二个脚本":        `<section/>` + `<script src="/evil.js"></script>` + `<script src="/assets/deck-v2/runtime.js"></script>`,
		"iframe":       `<iframe src="https://evil.com"></iframe><section/><script src="/assets/deck-v2/runtime.js"></script>`,
		"内联事件":         `<section onload="alert(1)"></section><script src="/assets/deck-v2/runtime.js"></script>`,
		"缺section":     `<div>hello</div><script src="/assets/deck-v2/runtime.js"></script>`,
		"缺runtime":     `<section class="slide">a</section>`,
		"丢style.css链接": `<html><head><link rel="stylesheet" href="/assets/deck-v2/base.css"></head><body><div class="deck"><section class="slide">a</section></div><script src="/assets/deck-v2/runtime.js"></script></body></html>`,
		"丢框架样式link":    `<html><head><link rel="stylesheet" href="style.css"></head><body><div class="deck"><section class="slide">a</section></div><script src="/assets/deck-v2/runtime.js"></script></body></html>`,
	} {
		if err := scanHTML(bad); err == nil {
			t.Errorf("%s 应被拒绝", name)
		}
	}
}

func TestExecWriteStyleAndDemo(t *testing.T) {
	s, st, id, _ := newHistService(t, "ws")
	row := &store.UserTemplate{ID: id, BaseID: "tech-sharing", Status: "draft"}

	// write_style 合法：写盘 + dirty
	newCSS := ".tpl-x{--accent:#123456}\n.card{padding:24px}"
	result, dirty := s.execCustomTool(row, "write_style", fmt.Sprintf(`{"css":%q}`, newCSS))
	if !dirty {
		t.Fatalf("write_style 应 dirty，结果: %s", result)
	}
	data, _ := os.ReadFile(filepath.Join(s.Dir(id), "style.css"))
	if string(data) != newCSS {
		t.Errorf("style.css 未整体替换: %q", data)
	}
	// write_demo 合法：写盘 + dirty
	newHTML := demoHTML("new-cover")
	result, dirty = s.execCustomTool(row, "write_demo", fmt.Sprintf(`{"html":%q}`, newHTML))
	if !dirty {
		t.Fatalf("write_demo 应 dirty，结果: %s", result)
	}
	data, _ = os.ReadFile(filepath.Join(s.Dir(id), "index.html"))
	if string(data) != newHTML {
		t.Errorf("index.html 未整体替换: %q", data)
	}

	// 预检拦截：不落盘、不 dirty
	before, _ := os.ReadFile(filepath.Join(s.Dir(id), "style.css"))
	result, dirty = s.execCustomTool(row, "write_style", fmt.Sprintf(`{"css":%q}`, ".a{background:url(https://evil.com/x.png)}"))
	if dirty {
		t.Error("预检失败不应 dirty")
	}
	if !strings.Contains(result, "安全预检未过") {
		t.Errorf("应返回预检错误: %s", result)
	}
	after, _ := os.ReadFile(filepath.Join(s.Dir(id), "style.css"))
	if string(before) != string(after) {
		t.Error("预检失败不应落盘")
	}
	_, dirty = s.execCustomTool(row, "write_demo", fmt.Sprintf(`{"html":%q}`, "<div>没有骨架</div>"))
	if dirty {
		t.Error("缺 section/runtime 的 demo 应被预检拒绝")
	}

	// published 全路径锁（D2）
	if err := st.DB.Model(&store.UserTemplate{}).Where("id = ?", id).Update("status", "published").Error; err != nil {
		t.Fatal(err)
	}
	publishedRow := &store.UserTemplate{ID: id, BaseID: "tech-sharing", Status: "published"}
	for _, tool := range []string{"write_tokens", "write_style", "write_demo", "set_meta"} {
		args := map[string]string{
			"write_tokens": `{"tokens":{"--accent":"#000"}}`,
			"write_style":  fmt.Sprintf(`{"css":%q}`, newCSS),
			"write_demo":   fmt.Sprintf(`{"html":%q}`, newHTML),
			"set_meta":     `{"name":"x"}`,
		}[tool]
		result, dirty := s.execCustomTool(publishedRow, tool, args)
		if dirty {
			t.Errorf("published 下 %s 不应 dirty", tool)
		}
		if !strings.Contains(result, "下架") {
			t.Errorf("published 下 %s 应提示下架: %s", tool, result)
		}
	}
	// 文件未被改
	data, _ = os.ReadFile(filepath.Join(s.Dir(id), "style.css"))
	if string(data) != newCSS {
		t.Errorf("published 下 style.css 被改写: %q", data)
	}
}

func TestSaveStyleCSS(t *testing.T) {
	s, st, id, uid := newHistService(t, "savestyle")
	// 记一个 fork 基线版，之后 save 的 changed 才有意义
	_ = s.recordVersionUT(id, OpFork, "起点")

	if err := s.SaveStyleCSS(uid, id, ".tpl-x{--accent:#abcdef}"); err != nil {
		t.Fatalf("保存 style: %v", err)
	}
	data, _ := os.ReadFile(filepath.Join(s.Dir(id), "style.css"))
	if string(data) != ".tpl-x{--accent:#abcdef}" {
		t.Errorf("style.css 未写入: %q", data)
	}
	versions, _ := s.ListUTVersions(uid, id)
	if len(versions) != 2 || versions[0].Operation != OpEdit || versions[0].Detail != "手动编辑 style.css" {
		t.Fatalf("应记 edit 版: %+v", versions)
	}
	if len(versions[0].Changed) != 1 || versions[0].Changed[0] != "style.css" {
		t.Errorf("changed 应为 [style.css]: %v", versions[0].Changed)
	}
	// 预检拦截
	if err := s.SaveStyleCSS(uid, id, "@import url('https://evil.com/x.css');"); err == nil {
		t.Error("预检失败应报错")
	}
	// published
	if err := st.DB.Model(&store.UserTemplate{}).Where("id = ?", id).Update("status", "published").Error; err != nil {
		t.Fatal(err)
	}
	if err := s.SaveStyleCSS(uid, id, ".x{}"); !errors.Is(err, ErrPublished) {
		t.Errorf("published 应拒绝: %v", err)
	}
}

// TestSaveIndexHTML_RejectsInvalidDemo 手动保存同样过 scanHTML。
func TestSaveIndexHTML_RejectsInvalidDemo(t *testing.T) {
	s, _, id, uid := newHistService(t, "rejectdemo")
	err := s.SaveIndexHTML(uid, id, "<div>没有骨架也没有脚本</div>")
	if err == nil {
		t.Fatal("缺 section/runtime.js 的保存应被拒绝")
	}
	if !strings.Contains(err.Error(), "index.html") {
		t.Errorf("报错应指向 index.html 规则: %v", err)
	}
}
