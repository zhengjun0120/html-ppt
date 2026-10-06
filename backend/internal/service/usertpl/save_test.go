package usertpl

// SaveIndexHTML（编辑器手动保存）的回归：docs/user-template-history-plan.md §3.2。
// 历史系统上线后旧滚动备份退役：每次保存记一条 edit 版本。
// 覆盖：保存生效+记版、幂等（相同内容不记）、published 拒绝、越权拒绝。
//
// 注意：store.OpenMemory 是进程级共享内存库，UserTemplate 的 id 用固定值会与
// 其他测试文件撞唯一约束——这里每次运行用唯一 id（纳秒后缀）。

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"html-ppt/backend/internal/store"
	"html-ppt/backend/internal/trace"
)

func TestSaveIndexHTML(t *testing.T) {
	st, err := store.OpenMemory(t.Context())
	if err != nil {
		t.Fatalf("内存库: %v", err)
	}
	root := t.TempDir()
	s := New(nil, st, root, "", "", nil, trace.Config{}) // 该路径不依赖注册表（reg 为 nil，record 不触它）

	id := fmt.Sprintf("ut-edit-%d", len(root))
	const uid = 7
	if err := st.DB.Create(&store.UserTemplate{ID: id, UserID: uid, BaseID: "tech-sharing", Status: "draft"}).Error; err != nil {
		t.Fatalf("建行: %v", err)
	}
	if err := os.MkdirAll(s.Dir(id), 0o755); err != nil {
		t.Fatalf("建目录: %v", err)
	}
	if err := os.WriteFile(filepath.Join(s.Dir(id), "index.html"), []byte(demoHTML("v1")), 0o644); err != nil {
		t.Fatalf("写初版: %v", err)
	}

	t.Run("保存生效并记 edit 版本", func(t *testing.T) {
		if err := s.SaveIndexHTML(uid, id, demoHTML("v2")); err != nil {
			t.Fatalf("保存: %v", err)
		}
		data, _ := os.ReadFile(filepath.Join(s.Dir(id), "index.html"))
		if string(data) != demoHTML("v2") {
			t.Errorf("index.html = %q", data)
		}
		versions, err := s.ListUTVersions(uid, id)
		if err != nil || len(versions) != 1 {
			t.Fatalf("应有 1 版，得到 %d 版 err=%v", len(versions), err)
		}
		if versions[0].Operation != OpEdit || versions[0].Detail != "手动编辑 index.html" {
			t.Errorf("版本元信息不对: %+v", versions[0])
		}
		if versions[0].Version != "v000001" {
			t.Errorf("首版序号 = %s", versions[0].Version)
		}
	})

	t.Run("相同内容不产生噪音版本", func(t *testing.T) {
		if err := s.SaveIndexHTML(uid, id, demoHTML("v2")); err != nil {
			t.Fatalf("重复保存: %v", err)
		}
		versions, _ := s.ListUTVersions(uid, id)
		if len(versions) != 1 {
			t.Errorf("幂等保存后仍应 1 版，得到 %d", len(versions))
		}
	})

	t.Run("内容变化追加新版本", func(t *testing.T) {
		if err := s.SaveIndexHTML(uid, id, demoHTML("v3")); err != nil {
			t.Fatalf("保存: %v", err)
		}
		versions, _ := s.ListUTVersions(uid, id)
		if len(versions) != 2 || versions[0].Version != "v000002" {
			t.Fatalf("应 2 版且新版 v000002，得到 %+v", versions)
		}
		// changed 只含 index.html（其余文件与上一版一致或双方皆空）
		if len(versions[0].Changed) != 1 || versions[0].Changed[0] != "index.html" {
			t.Errorf("changed 应为 [index.html]，得到 %v", versions[0].Changed)
		}
	})

	t.Run("published拒绝", func(t *testing.T) {
		if err := st.DB.Model(&store.UserTemplate{}).Where("id = ?", id).Update("status", "published").Error; err != nil {
			t.Fatalf("置 published: %v", err)
		}
		if err := s.SaveIndexHTML(uid, id, demoHTML("hack")); !errors.Is(err, ErrPublished) {
			t.Errorf("应报 ErrPublished，得到: %v", err)
		}
		data, _ := os.ReadFile(filepath.Join(s.Dir(id), "index.html"))
		if string(data) != demoHTML("v3") {
			t.Errorf("published 状态下内容被改写: %q", data)
		}
		versions, _ := s.ListUTVersions(uid, id)
		if len(versions) != 2 {
			t.Errorf("published 拒绝不应记版，现在 %d 版", len(versions))
		}
	})

	t.Run("越权拒绝", func(t *testing.T) {
		if err := s.SaveIndexHTML(uid+1, id, demoHTML("x")); err == nil || !strings.Contains(err.Error(), "不存在") {
			t.Errorf("越权应报不存在，得到: %v", err)
		}
	})
}

// demoHTML 构造能过 scanHTML 的最小 demo（<section> + 唯一的 runtime.js 脚本 +
// head 的框架样式 link 与 style.css 引用——2026-10-06 起是 scanHTML 的结构前提）。
func demoHTML(tag string) string {
	return `<html><head><link rel="stylesheet" href="/assets/deck-v2/base.css"><link rel="stylesheet" href="style.css"></head>` +
		`<body class="tpl-x"><div class="deck"><section class="slide">` + tag + `</section></div>` +
		`<script src="/assets/deck-v2/runtime.js"></script></body></html>`
}
