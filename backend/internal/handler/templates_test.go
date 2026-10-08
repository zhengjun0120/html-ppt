package handler

// PreviewTemplate 的协商缓存回归（2026-09-30 预览缓存化）：ETag 存在、
// If-None-Match 命中 304、不同 variant/slide 的 ETag 互不相同。
//
// 走真实仓库模板注册表（tech-sharing 有 demo 文件与多变体）+ gin 测试上下文，
// 不起数据库——预览端点不碰任何存储。

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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
	// 挂一个用户模板（拷贝 tech-sharing 改 id）——ut- 前缀走 no-cache 协商分支
	utDir := filepath.Join(t.TempDir(), "ut-test")
	if err := copyDir(filepath.Join(tplDir, "tech-sharing"), utDir); err != nil {
		t.Fatalf("拷贝模板: %v", err)
	}
	metaPath := filepath.Join(utDir, "template.json")
	raw, err := os.ReadFile(metaPath)
	if err != nil {
		t.Fatalf("读 template.json: %v", err)
	}
	var meta map[string]any
	if err := json.Unmarshal(raw, &meta); err != nil {
		t.Fatalf("解析 template.json: %v", err)
	}
	meta["id"] = "ut-test"
	patched, _ := json.Marshal(meta)
	if err := os.WriteFile(metaPath, patched, 0o644); err != nil {
		t.Fatalf("写 template.json: %v", err)
	}
	if err := reg.ValidateUserDir(utDir); err != nil {
		t.Fatalf("校验用户模板: %v", err)
	}
	if err := reg.MountUser(utDir); err != nil {
		t.Fatalf("挂载用户模板: %v", err)
	}
	return New(nil, nil, nil, nil, nil, nil, reg, nil, nil, trace.Config{}, nil)
}

func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			if err := copyDir(filepath.Join(src, e.Name()), filepath.Join(dst, e.Name())); err != nil {
				return err
			}
			continue
		}
		data, err := os.ReadFile(filepath.Join(src, e.Name()))
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dst, e.Name()), data, 0o644); err != nil {
			return err
		}
	}
	return nil
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

	t.Run("首响应带ETag与内置强缓存", func(t *testing.T) {
		w := previewRequest(t, h, "tech-sharing", "")
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		etag := w.Header().Get("ETag")
		if etag == "" {
			t.Fatal("缺 ETag 头")
		}
		if cc := w.Header().Get("Cache-Control"); cc != "private, max-age=3600" {
			t.Fatalf("内置模板 Cache-Control = %q, 期望 private, max-age=3600", cc)
		}
	})

	t.Run("用户模板走no-cache协商304", func(t *testing.T) {
		first := previewRequest(t, h, "ut-test", "")
		if cc := first.Header().Get("Cache-Control"); cc != "no-cache" {
			t.Fatalf("用户模板 Cache-Control = %q, 期望 no-cache", cc)
		}
		etag := first.Header().Get("ETag")
		if etag == "" {
			t.Fatal("缺 ETag 头")
		}
		w := previewRequest(t, h, "ut-test", etag)
		if w.Code != http.StatusNotModified {
			t.Fatalf("status = %d, 期望 304", w.Code)
		}
		if w.Body.Len() != 0 {
			t.Fatalf("304 不应带 body, got %d bytes", w.Body.Len())
		}
	})

	t.Run("内置模板新鲜期内不做协商", func(t *testing.T) {
		first := previewRequest(t, h, "tech-sharing", "")
		w := previewRequest(t, h, "tech-sharing", first.Header().Get("ETag"))
		if w.Code != http.StatusOK {
			t.Fatalf("max-age 窗口内应直接 200, got %d", w.Code)
		}
		if w.Body.Len() == 0 {
			t.Fatal("200 应带完整 body")
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

// ListTemplates 的协商缓存：清单 ~296KB，重访没变应 304 零传输。
func TestListTemplatesCache(t *testing.T) {
	h := newPreviewTestHandler(t)

	listRequest := func(t *testing.T, ifNoneMatch string) *httptest.ResponseRecorder {
		t.Helper()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/api/templates", nil)
		if ifNoneMatch != "" {
			req.Header.Set("If-None-Match", ifNoneMatch)
		}
		c.Request = req
		h.ListTemplates(c)
		return w
	}

	first := listRequest(t, "")
	if first.Code != http.StatusOK {
		t.Fatalf("status = %d", first.Code)
	}
	etag := first.Header().Get("ETag")
	if etag == "" || first.Header().Get("Cache-Control") != "no-cache" {
		t.Fatalf("头不符: etag=%q cc=%q", etag, first.Header().Get("Cache-Control"))
	}
	again := listRequest(t, etag)
	if again.Code != http.StatusNotModified {
		t.Fatalf("第二次应 304, got %d", again.Code)
	}
}

// TestListTemplatesSlim 清单瘦身契约：列表不下发 layouts/fonts/source
// （占体积 ~80%，前端零消费）；单模板详情仍回全量 Meta。
func TestListTemplatesSlim(t *testing.T) {
	h := newPreviewTestHandler(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/templates", nil)
	h.ListTemplates(c)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var list []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Fatalf("解析清单: %v", err)
	}
	if len(list) == 0 {
		t.Fatal("清单为空")
	}
	for _, m := range list {
		for _, banned := range []string{"layouts", "fonts", "source"} {
			if _, ok := m[banned]; ok {
				t.Errorf("清单条目 %s 不应含 %q 字段", m["id"], banned)
			}
		}
		if m["id"] == "" || m["name"] == "" || m["variants"] == nil {
			t.Errorf("清单条目缺必备字段: %v", m)
		}
	}

	// 详情仍全量：layouts 在
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = httptest.NewRequest(http.MethodGet, "/api/templates/tech-sharing", nil)
	c2.Params = gin.Params{{Key: "id", Value: "tech-sharing"}}
	h.GetTemplate(c2)
	if w2.Code != http.StatusOK {
		t.Fatalf("详情 status = %d", w2.Code)
	}
	var detail map[string]any
	if err := json.Unmarshal(w2.Body.Bytes(), &detail); err != nil {
		t.Fatalf("解析详情: %v", err)
	}
	if _, ok := detail["layouts"]; !ok {
		t.Error("详情应仍含 layouts 字段")
	}
}

// serveUTFile 的协商缓存：ut demo index/style 从 no-store 换成 ETag+no-cache，
// 重访没变 304 零传输，改文件 ETag 变立刻拿新版。走假 usertpl 服务（只覆写 Dir）。
type fakeUTService struct {
	UserTemplateService
	dir string
}

func (f *fakeUTService) Dir(string) string { return f.dir }

func TestServeUTFileCache(t *testing.T) {
	dir := t.TempDir()
	indexHTML := `<!DOCTYPE html><html><head><link rel="stylesheet" href="style.css"></head><body>demo</body></html>`
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte(indexHTML), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "style.css"), []byte("body{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	h := &Handler{usertpl: &fakeUTService{dir: dir}}

	// 走真实 gin engine：直调 handler 的话 c.Status() 的延迟 WriteHeader 不会被
	// 冲刷（engine 的兜底 WriteHeaderNow 不在），404 这类无 body 响应测不准
	utRequest := func(t *testing.T, name, styleHref, ifNoneMatch string) *httptest.ResponseRecorder {
		t.Helper()
		r := gin.New()
		r.GET("/api/user-templates/:id/*any", func(c *gin.Context) {
			h.serveUTFile(c, "ut-x", name, styleHref)
		})
		req := httptest.NewRequest(http.MethodGet, "/api/user-templates/ut-x/"+name, nil)
		if ifNoneMatch != "" {
			req.Header.Set("If-None-Match", ifNoneMatch)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	t.Run("index协商304", func(t *testing.T) {
		first := utRequest(t, "index.html", "/api/user-templates/ut-x/assets/style.css?token=t", "")
		if first.Code != http.StatusOK {
			t.Fatalf("status = %d", first.Code)
		}
		if first.Header().Get("Cache-Control") != "no-cache" {
			t.Fatalf("Cache-Control = %q", first.Header().Get("Cache-Control"))
		}
		// styleHref 重写发生在 ETag 之前：body 里应是重写后的引用
		if !strings.Contains(first.Body.String(), "token=t") {
			t.Fatal("style.css 引用未重写")
		}
		again := utRequest(t, "index.html", "/api/user-templates/ut-x/assets/style.css?token=t", first.Header().Get("ETag"))
		if again.Code != http.StatusNotModified {
			t.Fatalf("第二次应 304, got %d", again.Code)
		}
	})

	t.Run("style协商304", func(t *testing.T) {
		first := utRequest(t, "style.css", "", "")
		if first.Header().Get("Cache-Control") != "no-cache" {
			t.Fatalf("Cache-Control = %q", first.Header().Get("Cache-Control"))
		}
		again := utRequest(t, "style.css", "", first.Header().Get("ETag"))
		if again.Code != http.StatusNotModified {
			t.Fatalf("第二次应 304, got %d", again.Code)
		}
	})

	t.Run("改文件后ETag变", func(t *testing.T) {
		first := utRequest(t, "style.css", "", "")
		if err := os.WriteFile(filepath.Join(dir, "style.css"), []byte("body{color:red}"), 0o644); err != nil {
			t.Fatal(err)
		}
		again := utRequest(t, "style.css", "", first.Header().Get("ETag"))
		if again.Code != http.StatusOK {
			t.Fatalf("内容变了应 200 新版, got %d", again.Code)
		}
	})

	t.Run("缺失文件404", func(t *testing.T) {
		w := utRequest(t, "no-such.txt", "", "")
		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d", w.Code)
		}
	})
}
