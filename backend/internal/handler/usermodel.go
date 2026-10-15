package handler

// 模型管理端点（/api/models）：用户自选 OpenAI 兼容模型的增删改查、
// 设为当前使用、连通性测试。key 永不出后端——列表只回 has_key，测试由
// 后端拿解密后的配置代跑。

import (
	"net/http"
	"strconv"

	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/response"
	"html-ppt/backend/internal/service/usermodel"

	"github.com/gin-gonic/gin"
)

type modelRequest struct {
	Name    string `json:"name"`
	ModelID string `json:"model_id"`
	BaseURL string `json:"base_url"`
	APIKey  string `json:"api_key"`
}

type modelItem struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	ModelID   string `json:"model_id"`
	BaseURL   string `json:"base_url"`
	HasKey    bool   `json:"has_key"`
	CreatedAt string `json:"created_at"`
}

type modelListResp struct {
	Models        []modelItem `json:"models"`
	ActiveID      uint        `json:"active_id"`
	PlatformModel string      `json:"platform_model"`
}

func (h *Handler) ListModels(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	if !h.um.Available() {
		response.Err(c, http.StatusServiceUnavailable, "模型管理不可用")
		return
	}
	rows, err := h.um.List(c.Request.Context(), uid)
	if err != nil {
		response.Err(c, usermodel.HTTPStatus(err), err.Error())
		return
	}
	items := make([]modelItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, modelItem{
			ID: r.ID, Name: r.Name, ModelID: r.ModelID, BaseURL: r.BaseURL,
			HasKey: r.APIKeyEnc != "", CreatedAt: r.CreatedAt.Format("2006-01-02 15:04"),
		})
	}
	activeID, err := h.um.ActiveID(c.Request.Context(), uid)
	if err != nil {
		response.Err(c, usermodel.HTTPStatus(err), err.Error())
		return
	}
	platformModel := ""
	if h.agent != nil { // 测试装配可为 nil；平台模型名只作展示
		platformModel = h.agent.PlatformModelID()
	}
	response.OK(c, modelListResp{
		Models:        items,
		ActiveID:      activeID,
		PlatformModel: platformModel,
	})
}

func (h *Handler) CreateModel(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	if !h.um.Available() {
		response.Err(c, http.StatusServiceUnavailable, "模型管理不可用")
		return
	}
	var req modelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, http.StatusBadRequest, "参数不合法")
		return
	}
	row, err := h.um.Create(c.Request.Context(), uid, req.Name, req.ModelID, req.BaseURL, req.APIKey)
	if err != nil {
		response.Err(c, usermodel.HTTPStatus(err), err.Error())
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateModel(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	if !h.um.Available() {
		response.Err(c, http.StatusServiceUnavailable, "模型管理不可用")
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Err(c, http.StatusBadRequest, "模型 id 不合法")
		return
	}
	var req modelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, http.StatusBadRequest, "参数不合法")
		return
	}
	if err := h.um.Update(c.Request.Context(), uid, uint(id), req.Name, req.ModelID, req.BaseURL, req.APIKey); err != nil {
		response.Err(c, usermodel.HTTPStatus(err), err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) DeleteModel(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	if !h.um.Available() {
		response.Err(c, http.StatusServiceUnavailable, "模型管理不可用")
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Err(c, http.StatusBadRequest, "模型 id 不合法")
		return
	}
	if err := h.um.Delete(c.Request.Context(), uid, uint(id)); err != nil {
		response.Err(c, usermodel.HTTPStatus(err), err.Error())
		return
	}
	response.OK(c, nil)
}

type setActiveRequest struct {
	ModelID uint `json:"model_id"` // 0 = 切回平台模型
}

func (h *Handler) SetActiveModel(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	if !h.um.Available() {
		response.Err(c, http.StatusServiceUnavailable, "模型管理不可用")
		return
	}
	var req setActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, http.StatusBadRequest, "参数不合法")
		return
	}
	if err := h.um.SetActive(c.Request.Context(), uid, req.ModelID); err != nil {
		response.Err(c, usermodel.HTTPStatus(err), err.Error())
		return
	}
	response.OK(c, nil)
}

type testModelRequest struct {
	ID      uint   `json:"id"` // 给了 id：按已存配置测（api_key 可选覆盖）
	BaseURL string `json:"base_url"`
	ModelID string `json:"model_id"`
	APIKey  string `json:"api_key"`
}

func (h *Handler) TestModel(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	if !h.um.Available() {
		response.Err(c, http.StatusServiceUnavailable, "模型管理不可用")
		return
	}
	var req testModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, http.StatusBadRequest, "参数不合法")
		return
	}
	baseURL, modelID, apiKey := req.BaseURL, req.ModelID, req.APIKey
	if req.ID != 0 {
		storedBase, storedModel, storedKey, err := h.um.ConfigFor(c.Request.Context(), uid, req.ID)
		if err != nil {
			response.Err(c, usermodel.HTTPStatus(err), err.Error())
			return
		}
		baseURL, modelID = storedBase, storedModel
		if apiKey == "" {
			apiKey = storedKey // 编辑场景不回填 key：测试时用库存的
		}
	}
	if err := h.um.Test(c.Request.Context(), baseURL, modelID, apiKey); err != nil {
		response.Err(c, http.StatusBadGateway, "连通失败："+err.Error())
		return
	}
	response.OK(c, gin.H{"ok": true})
}
