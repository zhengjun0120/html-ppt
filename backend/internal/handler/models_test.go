package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/cryptox"
	"html-ppt/backend/internal/service/usermodel"
	"html-ppt/backend/internal/store"
	"html-ppt/backend/internal/trace"
)

// 模型管理端到端（handler 层）：CRUD + 设为当前使用 + 归属隔离。
// OpenMemory 共享库：用户 id 独占 96_000 段位。
func TestModelsCRUD(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uid := uint(96_001)

	st, err := store.OpenMemory(t.Context())
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	u := store.User{Email: "models-crud@example.com"}
	u.ID = uid
	if err := st.DB.Create(&u).Error; err != nil {
		t.Fatalf("造用户: %v", err)
	}
	box, err := cryptox.NewBox("MDEyMzQ1Njc4OWFiY2RlZjAxMjM0NTY3ODlhYmNkZWY=") // 32 字节固定测试密钥
	if err != nil {
		t.Fatalf("NewBox: %v", err)
	}
	umSvc := usermodel.New(st.DB, box)

	h := New(nil, nil, nil, nil, nil, nil, nil, nil, nil, trace.Config{}, nil, nil, umSvc)
	r := gin.New()
	r.POST("/api/models", func(c *gin.Context) { h.CreateModel(c) })
	r.GET("/api/models", func(c *gin.Context) { h.ListModels(c) })
	r.PUT("/api/models/:id", func(c *gin.Context) { h.UpdateModel(c) })
	r.DELETE("/api/models/:id", func(c *gin.Context) { h.DeleteModel(c) })
	r.POST("/api/models/active", func(c *gin.Context) { h.SetActiveModel(c) })

	do := func(method, path, body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(authctx.WithUser(req.Context(), uid))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	// 新增：name 缺省取 model_id
	w := do(http.MethodPost, "/api/models", `{"model_id":"gpt-4o","base_url":"https://api.openai.com/v1","api_key":"sk-x"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("新增失败 %d: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"name":"gpt-4o"`) || strings.Contains(w.Body.String(), "sk-x") {
		t.Errorf("默认名/密钥泄漏检查: %s", w.Body.String())
	}

	// 列表：has_key 布尔而非密文；带平台模型名
	w = do(http.MethodGet, "/api/models", "")
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"has_key":true`) {
		t.Errorf("列表不对: %d %s", w.Code, w.Body.String())
	}

	// 设为当前使用 → 列表 active_id=1 段位内
	w = do(http.MethodPost, "/api/models/active", `{"model_id":1}`)
	if w.Code != http.StatusOK {
		t.Fatalf("设当前失败 %d: %s", w.Code, w.Body.String())
	}
	w = do(http.MethodGet, "/api/models", "")
	if !strings.Contains(w.Body.String(), `"active_id":1`) {
		t.Errorf("active_id 不对: %s", w.Body.String())
	}

	// 改名（api_key 空 = 保留）
	w = do(http.MethodPut, "/api/models/1", `{"name":"我的GPT","model_id":"gpt-4o","base_url":"https://api.openai.com/v1"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("更新失败 %d: %s", w.Code, w.Body.String())
	}
	w = do(http.MethodGet, "/api/models", "")
	if !strings.Contains(w.Body.String(), `"name":"我的GPT"`) || !strings.Contains(w.Body.String(), `"has_key":true`) {
		t.Errorf("更新后列表不对: %s", w.Body.String())
	}

	// 删除当前使用的模型 → active_id 归零
	w = do(http.MethodDelete, "/api/models/1", "")
	if w.Code != http.StatusOK {
		t.Fatalf("删除失败 %d: %s", w.Code, w.Body.String())
	}
	w = do(http.MethodGet, "/api/models", "")
	if !strings.Contains(w.Body.String(), `"active_id":0`) {
		t.Errorf("删除后 active_id 应归零: %s", w.Body.String())
	}
}
