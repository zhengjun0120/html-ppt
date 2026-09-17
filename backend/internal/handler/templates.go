package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/response"
)

// ListTemplates GET /api/templates —— 画廊的模板清单（公开）。
//
// 为什么公开：模板元数据不含任何用户数据，画廊在登录页之外的展示场景
// （落地页、分享预览）也需要它。体积小（每模板一份元数据），不做分页。
func (h *Handler) ListTemplates(c *gin.Context) {
	if h.templates == nil {
		response.Err(c, http.StatusServiceUnavailable, "模板库不可用")
		return
	}
	response.OK(c, h.templates.List())
}

// GetTemplate GET /api/templates/:id —— 单个模板详情（含版式索引，前端向导展示用）。
func (h *Handler) GetTemplate(c *gin.Context) {
	if h.templates == nil {
		response.Err(c, http.StatusServiceUnavailable, "模板库不可用")
		return
	}
	t, err := h.templates.Get(c.Param("id"))
	if err != nil {
		response.Err(c, http.StatusNotFound, "模板不存在")
		return
	}
	response.OK(c, t.Meta)
}
