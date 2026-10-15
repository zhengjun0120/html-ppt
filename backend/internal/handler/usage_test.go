package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/service/usage"
	"html-ppt/backend/internal/store"
	"html-ppt/backend/internal/trace"
)

func TestUsageOverview(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 与 service/usage 的测试用户错开段位（OpenMemory 共享库防撞键）。
	seed := func(t *testing.T, uid uint) *usage.Service {
		t.Helper()
		st, err := store.OpenMemory(t.Context())
		if err != nil {
			t.Fatalf("OpenMemory: %v", err)
		}
		svc := usage.New(st.DB)
		row := store.UsageEvent{
			UserID: uid, Day: time.Now().Format("2006-01-02"), Component: "main", Model: "deepseek-flash",
			Prompt: 100, Completion: 20, Total: 120, Cached: 60, Calls: 1,
		}
		if err := st.DB.Create(&row).Error; err != nil {
			t.Fatalf("造行: %v", err)
		}
		return svc
	}

	t.Run("未登录401", func(t *testing.T) {
		h := New(nil, nil, nil, nil, nil, nil, nil, nil, nil, trace.Config{}, nil, nil, nil)
		r := gin.New()
		r.GET("/api/usage/overview", h.UsageOverview)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/usage/overview", nil))
		if w.Code != http.StatusUnauthorized {
			t.Errorf("想要 401，实际 %d", w.Code)
		}
	})

	t.Run("服务未装配503", func(t *testing.T) {
		h := New(nil, nil, nil, nil, nil, nil, nil, nil, nil, trace.Config{}, nil, nil, nil)
		r := gin.New()
		r.GET("/api/usage/overview", h.UsageOverview)
		req := httptest.NewRequest(http.MethodGet, "/api/usage/overview", nil)
		req = req.WithContext(authctx.WithUser(req.Context(), 7))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("想要 503，实际 %d", w.Code)
		}
	})

	t.Run("已登录返回聚合", func(t *testing.T) {
		uid := uint(92_001)
		svc := seed(t, uid)
		h := New(nil, nil, nil, nil, nil, nil, nil, nil, nil, trace.Config{}, nil, svc, nil)
		r := gin.New()
		r.GET("/api/usage/overview", h.UsageOverview)
		req := httptest.NewRequest(http.MethodGet, "/api/usage/overview", nil)
		req = req.WithContext(authctx.WithUser(req.Context(), uid))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("想要 200，实际 %d: %s", w.Code, w.Body.String())
		}
		body := w.Body.String()
		for _, want := range []string{`"today"`, `"month"`, `"daily"`, `"total":120`, `"cached_rate":0.6`} {
			if !strings.Contains(body, want) {
				t.Errorf("响应缺 %s: %s", want, body)
			}
		}
	})
}
