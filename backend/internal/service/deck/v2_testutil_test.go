package deck

// 测试辅助：内存库 + 真实模板注册表（跨测试文件共用）。

import (
	"context"
	"path/filepath"
	"testing"

	"html-ppt/backend/internal/service/template"
	"html-ppt/backend/internal/store"
)

func newTestStore() (*store.Store, error) {
	return store.OpenMemory(context.Background())
}

func newTestRegistry(t *testing.T) *template.Registry {
	t.Helper()
	tplDir := filepath.Clean(filepath.Join("..", "..", "..", "templates"))
	assetsDir := filepath.Clean(filepath.Join("..", "..", "..", "web", "assets"))
	reg, err := template.NewRegistry(tplDir, assetsDir)
	if err != nil {
		t.Fatalf("模板注册表: %v", err)
	}
	return reg
}
