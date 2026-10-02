package handler

// ListTraces 的回归：{runs,total} 响应形状（sessions 汇总已下线）、kind/status
// 白名单 400、q 服务端全文搜索（响应里的 user_content 是截断预览，但搜索必须
// 命中全文）、分页窗口。归属隔离由 trace.Store 层测试覆盖，这里只验证接线。

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/trace"
)

func seedTraceRun(t *testing.T, dir string, sessID, uid uint, runID, deckID, content, kind string, end bool) {
	t.Helper()
	r, err := trace.New(trace.Options{Dir: dir, SessionID: sessID, RunID: runID, UserID: uid, DeckID: deckID})
	if err != nil {
		t.Fatalf("建 recorder: %v", err)
	}
	defer r.Close()
	ctx := trace.With(context.Background(), r)
	trace.Emit(ctx, trace.Event{
		Kind: trace.KindRunStart, RunKind: kind, RunID: runID, SessionID: sessID,
		UserID: uid, DeckID: deckID, UserContent: content, Model: "test-model",
	})
	if end {
		trace.Emit(ctx, trace.Event{Kind: trace.KindRunEnd, Status: trace.StatusOK,
			Summary: &trace.Summary{DurationMS: 100, Turns: 1}})
	}
}

func TestListTraces(t *testing.T) {
	gin.SetMode(gin.TestMode)
	dir := t.TempDir()
	long := strings.Repeat("复", 200) // 200 rune，超过 160 的预览上限
	seedTraceRun(t, dir, 1, 7, "1700000000001-a", "deck-0001", "帮我做讲稿", "", true)
	seedTraceRun(t, dir, 2, 7, "1700000000002-b", "ut-abc", "换个配色", "customize", true)
	seedTraceRun(t, dir, 3, 7, "1700000000003-c", "deck-0002", long, "tplsugg", true)
	seedTraceRun(t, dir, 9, 99, "1700000000009-x", "deck-999", "别人的", "", true)

	h := New(nil, nil, nil, nil, nil, trace.NewStore(dir), nil, nil, nil, trace.Config{}, nil)
	r := gin.New()
	r.GET("/api/traces", h.ListTraces)
	do := func(query string, uid uint) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/api/traces"+query, nil)
		req = req.WithContext(authctx.WithUser(req.Context(), uid))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	t.Run("默认列表带total无sessions", func(t *testing.T) {
		w := do("", 7)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var body map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if _, ok := body["sessions"]; ok {
			t.Error("sessions 汇总已下线，不应再出现在响应里")
		}
		if body["total"].(float64) != 3 {
			t.Errorf("total = %v, 期望 3（只看自己的 run）", body["total"])
		}
		if runs := body["runs"].([]any); len(runs) != 3 {
			t.Errorf("runs = %d 条, 期望 3", len(runs))
		}
	})

	t.Run("未登录401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/traces", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, 期望 401", w.Code)
		}
	})

	t.Run("kind筛选与白名单", func(t *testing.T) {
		w := do("?kind=deck", 7)
		var body struct {
			Total int             `json:"total"`
			Runs  []trace.RunMeta `json:"runs"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if body.Total != 1 || len(body.Runs) != 1 || body.Runs[0].RunID != "1700000000001-a" {
			t.Errorf("kind=deck 应只剩空 run_kind 的旧数据那条, got total=%d", body.Total)
		}
		if w := do("?kind=bogus", 7); w.Code != http.StatusBadRequest {
			t.Errorf("未知 kind 应 400, got %d", w.Code)
		}
	})

	t.Run("status筛选", func(t *testing.T) {
		if w := do("?status=running", 7); !strings.Contains(w.Body.String(), `"total":0`) {
			t.Errorf("全部已结束，running 应 0 条: %s", w.Body.String())
		}
		if w := do("?status=ok", 7); !strings.Contains(w.Body.String(), `"total":3`) {
			t.Errorf("status=ok 应 3 条: %s", w.Body.String())
		}
		if w := do("?status=nope", 7); w.Code != http.StatusBadRequest {
			t.Errorf("未知 status 应 400, got %d", w.Code)
		}
	})

	t.Run("q搜索命中全文而响应是预览", func(t *testing.T) {
		// 响应里的 user_content 已截到 160 rune；q 用第 170 个字之后的片段，
		// 命中说明搜索走的是索引里的全文。片段必须按 rune 切——字节切片会把
		// 多字节汉字切成非法 UTF-8，ToLower 会把它换成 U+FFFD 而永远匹配不上
		q := string([]rune(long)[170:175])
		w := do("?q="+q, 7)
		var body struct {
			Total int             `json:"total"`
			Runs  []trace.RunMeta `json:"runs"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if body.Total != 1 || body.Runs[0].RunID != "1700000000003-c" {
			t.Fatalf("q 应命中全文被截断的那条, total=%d", body.Total)
		}
		got := []rune(body.Runs[0].UserContent)
		if len(got) > 161 || !strings.HasSuffix(body.Runs[0].UserContent, "…") {
			t.Errorf("列表 user_content 应截到 160 rune + 省略号, 实际 %d rune", len(got))
		}
	})

	t.Run("limit与offset窗口", func(t *testing.T) {
		w := do("?limit=2", 7)
		var body struct {
			Total int             `json:"total"`
			Runs  []trace.RunMeta `json:"runs"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if body.Total != 3 || len(body.Runs) != 2 || body.Runs[0].RunID != "1700000000003-c" {
			t.Errorf("limit=2 应返回最新 2 条且 total=3, got %+v", body)
		}
		if w := do("?offset=1&limit=1", 7); !strings.Contains(w.Body.String(), "1700000000002-b") {
			t.Errorf("offset=1 应跳过最新一条: %s", w.Body.String())
		}
	})
}
