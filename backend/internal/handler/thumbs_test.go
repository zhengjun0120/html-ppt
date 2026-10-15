package handler

// TemplateThumb 的回归：版本命中出图（immutable 头）、版本不匹配 404、
// 未装配 404。渲染路径（Chrome）不在单测覆盖——走浏览器验证。

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/service/template"
	"html-ppt/backend/internal/service/tplthumb"
	"html-ppt/backend/internal/trace"
)

func newThumbTestHandler(t *testing.T) (*Handler, string, string) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	tplDir, _ := filepath.Abs(filepath.Join("..", "..", "templates"))
	assetsDir, _ := filepath.Abs(filepath.Join("..", "..", "web", "assets"))
	reg, err := template.NewRegistry(tplDir, assetsDir)
	if err != nil {
		t.Fatalf("模板注册表: %v", err)
	}
	thumbDir := t.TempDir()
	h := New(nil, nil, nil, nil, nil, trace.NewStore(t.TempDir()), reg, nil, nil, trace.Config{},
		tplthumb.New(reg, "http://127.0.0.1:0", "", thumbDir, 0), nil, nil)
	v, err := reg.ContentVersion("tech-sharing")
	if err != nil {
		t.Fatalf("ContentVersion: %v", err)
	}
	return h, thumbDir, v[:8]
}

// thumbRequest 走真实 gin engine：直调 handler 的话 c.Status() 的延迟
// WriteHeader 不被冲刷，无 body 的 404 在 recorder 上测不准（同 serveUTFile 测试的教训）。
func thumbRequest(h *Handler, id, v string) *httptest.ResponseRecorder {
	r := gin.New()
	r.GET("/api/templates/:id/thumb", func(c *gin.Context) {
		c.Params = gin.Params{{Key: "id", Value: id}}
		h.TemplateThumb(c)
	})
	req := httptest.NewRequest(http.MethodGet, "/api/templates/"+id+"/thumb?v="+v, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestTemplateThumb(t *testing.T) {
	h, dir, version := newThumbTestHandler(t)

	t.Run("未装配404", func(t *testing.T) {
		bare := New(nil, nil, nil, nil, nil, trace.NewStore(t.TempDir()), nil, nil, nil, trace.Config{}, nil, nil, nil)
		if w := thumbRequest(bare, "tech-sharing", version); w.Code != http.StatusNotFound {
			t.Fatalf("status = %d, 期望 404", w.Code)
		}
	})

	t.Run("版本不匹配404", func(t *testing.T) {
		if w := thumbRequest(h, "tech-sharing", "deadbeef"); w.Code != http.StatusNotFound {
			t.Fatalf("status = %d, 期望 404", w.Code)
		}
		if w := thumbRequest(h, "tech-sharing", ""); w.Code != http.StatusNotFound {
			t.Fatalf("缺版本应 404, got %d", w.Code)
		}
	})

	t.Run("磁盘命中出图带immutable", func(t *testing.T) {
		// 预置缓存文件（渲染路径不在单测范围）
		png := filepath.Join(dir, "tech-sharing-"+version+".png")
		if err := os.WriteFile(png, []byte("fake-png-bytes"), 0o644); err != nil {
			t.Fatal(err)
		}
		w := thumbRequest(h, "tech-sharing", version)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		if got := w.Header().Get("Content-Type"); got != "image/png" {
			t.Fatalf("Content-Type = %q", got)
		}
		if cc := w.Header().Get("Cache-Control"); cc != "public, max-age=31536000, immutable" {
			t.Fatalf("Cache-Control = %q", cc)
		}
	})

	t.Run("未知模板404", func(t *testing.T) {
		if w := thumbRequest(h, "no-such", version); w.Code != http.StatusNotFound {
			t.Fatalf("status = %d", w.Code)
		}
	})
}
