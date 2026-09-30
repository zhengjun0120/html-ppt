package handler

// PreviewTemplate 的协商缓存回归（2026-09-30 预览缓存化）：ETag 存在、
// If-None-Match 命中 304、不同 variant/slide 的 ETag 互不相同。
//
// 走真实仓库模板注册表（tech-sharing 有 demo 文件与多变体）+ gin 测试上下文，
// 不起数据库——预览端点不碰任何存储。

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/service/template"
	"html-ppt/backend/internal/trace"
)

func newPreviewTestHandler(t *testing.T) *Handler {
	t.Helper()
	gin.SetMode(gin.TestMode)
	tplDir := filepath.Clean(filepath.Join("..", "..", "templates"))
	assetsDir := filepath.Clean(filepath.Join("..", "..", "web", "assets"))
	reg, err := template.NewRegistry(tplDir, assetsDir)
	if err != nil {
		t.Fatalf("模板注册表: %v", err)
	}
	return New(nil, nil, nil, nil, nil, nil, reg, nil, nil, trace.Config{})
}

// previewRequest 直调 PreviewTemplate（gin 测试上下文，带 path 参数与请求头）。
func previewRequest(t *testing.T, h *Handler, id string, ifNoneMatch string) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/api/templates/"+id+"/preview", nil)
	if ifNoneMatch != "" {
		req.Header.Set("If-None-Match", ifNoneMatch)
	}
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: id}}
	h.PreviewTemplate(c)
	return w
}

func TestPreviewTemplateCache(t *testing.T) {
	h := newPreviewTestHandler(t)

	t.Run("首响应带ETag与no-cache", func(t *testing.T) {
		w := previewRequest(t, h, "tech-sharing", "")
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		etag := w.Header().Get("ETag")
		if etag == "" {
			t.Fatal("缺 ETag 头")
		}
		if cc := w.Header().Get("Cache-Control"); cc != "no-cache" {
			t.Fatalf("Cache-Control = %q, 期望 no-cache", cc)
		}
	})

	t.Run("IfNoneMatch命中304", func(t *testing.T) {
		first := previewRequest(t, h, "tech-sharing", "")
		etag := first.Header().Get("ETag")
		w := previewRequest(t, h, "tech-sharing", etag)
		if w.Code != http.StatusNotModified {
			t.Fatalf("status = %d, 期望 304", w.Code)
		}
		if w.Body.Len() != 0 {
			t.Fatalf("304 不应带 body, got %d bytes", w.Body.Len())
		}
	})

	t.Run("变体与页数改变ETag", func(t *testing.T) {
		base := previewRequest(t, h, "tech-sharing", "")
		// slide 缩略版：换 URL（slide=1）后内容不同 → ETag 不同。
		// 直接造第二份响应头对比：手动带上 slide 参数调一次
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/api/templates/tech-sharing/preview?slide=1", nil)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "tech-sharing"}}
		h.PreviewTemplate(c)
		if w.Header().Get("ETag") == base.Header().Get("ETag") {
			t.Fatal("slide 版与整本版 ETag 不应相同")
		}
	})

	t.Run("不存在的模板404", func(t *testing.T) {
		w := previewRequest(t, h, "no-such-tpl", "")
		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d", w.Code)
		}
	})
}
