package router

// 传输压缩的端到端回归：清单端点在 Accept-Encoding: gzip 下返回压缩体，
// 解压后是合法 JSON 且已瘦身（无 layouts）；不带 gzip 头照常 identity。
// 走真实 engine + 真实模板注册表——中间件挂错路由、Vary 丢失这类装配问题
// 只有整条管线能暴露。

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/config"
	"html-ppt/backend/internal/handler"
	"html-ppt/backend/internal/service/template"
	"html-ppt/backend/internal/trace"
	"html-ppt/backend/internal/vision"
)

func newGzipTestEngine(t *testing.T) *gin.Engine {
	t.Helper()
	tplDir, err := filepath.Abs(filepath.Join("..", "..", "templates"))
	if err != nil {
		t.Fatal(err)
	}
	assetsDir, err := filepath.Abs(filepath.Join("..", "..", "web", "assets"))
	if err != nil {
		t.Fatal(err)
	}
	reg, err := template.NewRegistry(tplDir, assetsDir)
	if err != nil {
		t.Fatalf("模板注册表: %v", err)
	}
	h := handler.New(nil, nil, nil, nil, &vision.Grants{}, trace.NewStore(t.TempDir()), reg, nil, nil, trace.Config{}, nil, nil)
	return New(&config.Config{
		Assets:    config.Assets{Dir: t.TempDir()},
		Templates: config.Templates{Dir: tplDir},
		Auth:      config.Auth{JWTSecret: "test"},
		Features:  config.Features{Trace: false},
	}, h)
}

func TestListTemplatesGzipped(t *testing.T) {
	engine := newGzipTestEngine(t)

	req := httptest.NewRequest(http.MethodGet, "/api/templates", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	if enc := w.Header().Get("Content-Encoding"); enc != "gzip" {
		t.Fatalf("Content-Encoding = %q, 期望 gzip", enc)
	}
	zr, err := gzip.NewReader(bytes.NewReader(w.Body.Bytes()))
	if err != nil {
		t.Fatalf("响应不是合法 gzip 流: %v", err)
	}
	raw, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("解压失败: %v", err)
	}
	var list []map[string]any
	if err := json.Unmarshal(raw, &list); err != nil {
		t.Fatalf("解压后不是合法 JSON: %v", err)
	}
	if len(list) == 0 {
		t.Fatal("清单为空")
	}
	if _, ok := list[0]["layouts"]; ok {
		t.Error("压缩后的清单仍含 layouts（瘦身没生效或挂错了 handler）")
	}

	// 无 gzip 头：identity 照常可用（内容协商没把不支持压缩的客户端拒之门外）
	w2 := httptest.NewRecorder()
	engine.ServeHTTP(w2, httptest.NewRequest(http.MethodGet, "/api/templates", nil))
	if w2.Code != http.StatusOK {
		t.Fatalf("identity status = %d", w2.Code)
	}
	if w2.Header().Get("Content-Encoding") == "gzip" {
		t.Error("客户端没带 Accept-Encoding: gzip 却返回了压缩体")
	}
}
