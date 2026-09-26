package usertpl

// SaveIndexHTML（编辑器手动保存）的回归：docs/deck-editor-plan.md §4.2。
// 覆盖：覆盖前滚动备份、备份上限 5 版按时间戳淘汰、published 拒绝、越权拒绝。
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
)

func TestSaveIndexHTML(t *testing.T) {
	st, err := store.OpenMemory(t.Context())
	if err != nil {
		t.Fatalf("内存库: %v", err)
	}
	root := t.TempDir()
	s := New(nil, st, root, "", "") // SaveIndexHTML 路径不依赖注册表

	id := fmt.Sprintf("ut-edit-%d", len(root)) // 唯一且合法（ut- 前缀）
	const uid = 7
	if err := st.DB.Create(&store.UserTemplate{ID: id, UserID: uid, BaseID: "tech-sharing", Status: "draft"}).Error; err != nil {
		t.Fatalf("建行: %v", err)
	}
	if err := os.MkdirAll(s.Dir(id), 0o755); err != nil {
		t.Fatalf("建目录: %v", err)
	}
	writeIndex := func(content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(s.Dir(id), "index.html"), []byte(content), 0o644); err != nil {
			t.Fatalf("写 index.html: %v", err)
		}
	}
	writeIndex("v1")

	t.Run("覆盖并备份上一版", func(t *testing.T) {
		kept, err := s.SaveIndexHTML(uid, id, "v2")
		if err != nil || kept != 1 {
			t.Fatalf("首次保存: kept=%d err=%v", kept, err)
		}
		data, _ := os.ReadFile(filepath.Join(s.Dir(id), "index.html"))
		if string(data) != "v2" {
			t.Errorf("index.html = %q", data)
		}
		bak, _ := os.ReadFile(filepath.Join(s.backupDir(id), firstBackupName(t, s.backupDir(id))))
		if string(bak) != "v1" {
			t.Errorf("备份内容 = %q, 期望上一版 v1", bak)
		}
	})

	t.Run("滚动保留最近5版", func(t *testing.T) {
		for i := 3; i <= 8; i++ { // 再保存 6 次，共 7 个历史版本
			if _, err := s.SaveIndexHTML(uid, id, fmt.Sprintf("v%d", i)); err != nil {
				t.Fatalf("第 %d 次保存: %v", i, err)
			}
		}
		names := backupNames(t, s.backupDir(id))
		if len(names) != keepTemplateBackups {
			t.Fatalf("备份数 = %d, 期望 %d", len(names), keepTemplateBackups)
		}
		// 8 次保存（写 v2..v8）产生 7 个备份（v1..v7），裁到 5 个后
		// 保留 v3..v7：最旧的是 v3 的内容，v1/v2 已被淘汰
		first := readFile(t, filepath.Join(s.backupDir(id), names[0]))
		if first != "v3" {
			t.Errorf("最旧备份 = %q, 期望 v3", first)
		}
	})

	t.Run("published拒绝", func(t *testing.T) {
		if err := st.DB.Model(&store.UserTemplate{}).Where("id = ?", id).Update("status", "published").Error; err != nil {
			t.Fatalf("置 published: %v", err)
		}
		_, err := s.SaveIndexHTML(uid, id, "hack")
		if !errors.Is(err, ErrPublished) {
			t.Errorf("应报 ErrPublished，得到: %v", err)
		}
		data, _ := os.ReadFile(filepath.Join(s.Dir(id), "index.html"))
		if string(data) != "v8" {
			t.Errorf("published 状态下内容被改写: %q", data)
		}
	})

	t.Run("越权拒绝", func(t *testing.T) {
		if _, err := s.SaveIndexHTML(uid+1, id, "x"); err == nil || !strings.Contains(err.Error(), "不存在") {
			t.Errorf("越权应报不存在，得到: %v", err)
		}
	})
}

func backupNames(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("读备份目录: %v", err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".html") {
			names = append(names, e.Name())
		}
	}
	if len(names) == 0 {
		t.Fatal("备份目录为空")
	}
	return names
}

func firstBackupName(t *testing.T, dir string) string {
	t.Helper()
	names := backupNames(t, dir)
	return names[0]
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("读 %s: %v", path, err)
	}
	return string(data)
}
