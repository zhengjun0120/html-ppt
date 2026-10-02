package tplthumb

// Prewarm 的回归：内置模板全部预渲落盘、ut-* 一张不渲（草稿 demo 只认属主
// token，无 token 渲出来的是无样式页，落盘即污染缓存）、二次调用全磁盘命中
// 零渲染、连续失败熔断。真 Chrome 渲染不在单测范围——走浏览器验证。

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"html-ppt/backend/internal/service/template"
	"html-ppt/backend/internal/vision"
)

// renderLog 记录假渲染出口收到的 URL（并发不敏感，这里串行调用）。
type renderLog struct{ urls []string }

func (l *renderLog) add(u string)        { l.urls = append(l.urls, u) }
func (l *renderLog) count() int          { return len(l.urls) }
func (l *renderLog) has(frag string) bool {
	for _, u := range l.urls {
		if strings.Contains(u, frag) {
			return true
		}
	}
	return false
}

// newTestService 真注册表（内置 + 一个挂载的 ut-*）+ 假健康端点 + 假渲染出口。
func newTestService(t *testing.T) (*Service, *renderLog, string) {
	t.Helper()
	tplDir, _ := filepath.Abs(filepath.Join("..", "..", "..", "templates"))
	assetsDir, _ := filepath.Abs(filepath.Join("..", "..", "..", "web", "assets"))
	reg, err := template.NewRegistry(tplDir, assetsDir)
	if err != nil {
		t.Fatalf("模板注册表: %v", err)
	}
	// 挂一个用户模板（拷贝 tech-sharing 改 id）——预热必须跳过 ut- 前缀
	utDir := filepath.Join(t.TempDir(), "ut-prewarm")
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
	meta["id"] = "ut-prewarm"
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

	// waitSelf 探测 baseURL/api/health——给个真的 200 端点
	health := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(health.Close)

	dir := t.TempDir()
	svc := New(reg, health.URL, "", dir, 60*time.Second)
	lg := &renderLog{}
	svc.renderFn = func(ctx context.Context, url string) (*vision.Deck2, error) {
		lg.add(url)
		return &vision.Deck2{Slides: []vision.Slide2{{PNG: []byte("fake-png")}}}, nil
	}
	return svc, lg, dir
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

func TestPrewarm(t *testing.T) {
	t.Run("内置全渲ut跳过", func(t *testing.T) {
		svc, lg, dir := newTestService(t)
		svc.Prewarm(context.Background())

		if lg.has("/user-templates/") {
			t.Fatalf("预热不应碰 ut-* 的 demo 路由: %v", lg.urls)
		}
		builtins := 0
		for _, m := range svc.reg.List() {
			if strings.HasPrefix(m.ID, "ut-") {
				continue
			}
			builtins++
			v, err := svc.reg.ContentVersion(m.ID)
			if err != nil {
				t.Fatalf("ContentVersion(%s): %v", m.ID, err)
			}
			if _, err := os.Stat(filepath.Join(dir, m.ID+"-"+v[:8]+".png")); err != nil {
				t.Fatalf("内置模板 %s 预热后无落盘文件: %v", m.ID, err)
			}
		}
		if lg.count() != builtins {
			t.Fatalf("渲染次数 = %d, 期望内置模板数 %d", lg.count(), builtins)
		}

		// 二次预热：全磁盘命中，零渲染
		before := lg.count()
		svc.Prewarm(context.Background())
		if lg.count() != before {
			t.Fatalf("二次预热渲染了 %d 张, 期望 0（应全磁盘命中）", lg.count()-before)
		}
	})

	t.Run("连续失败熔断", func(t *testing.T) {
		svc, lg, _ := newTestService(t)
		svc.renderFn = func(ctx context.Context, url string) (*vision.Deck2, error) {
			lg.add(url)
			return nil, errors.New("chrome 挂了")
		}
		svc.Prewarm(context.Background())
		// 每张内部重试一次 = 1 张占 2 次 renderFn 调用；熔断阈值 3 张
		if got := lg.count(); got != 3*2 {
			t.Fatalf("连续失败应 3 张熔断（%d 次 renderFn 调用）, got %d", 3*2, got)
		}
	})
}
