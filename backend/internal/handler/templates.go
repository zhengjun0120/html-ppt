package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/response"
	"html-ppt/backend/internal/service/template"
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

// PreviewTemplate GET /api/templates/:id/preview?variant=<vid>&slide=<n> —— demo 页 HTML。
//
// 选模板页的实时换肤预览：variant 非空时服务端把变体 class 挂到 body，
// 前端切变体只改 iframe src 的 query。demo 里的相对引用（style.css）改写为
// 公开静态路由 /templates/<id>/ 的绝对路径——本端点自身带鉴权，资产走静态。
//
// slide>=1 是缩略模式：只返回第 n 页。画廊一屏十几张卡，每张都为整本 demo
// 付解析+布局的 CPU 账，缩略卡只看得到第一页，没必要把其余十几页也渲染了
// （TrimToSlide 的 fail-open 语义见 service 层注释）。
func (h *Handler) PreviewTemplate(c *gin.Context) {
	if h.templates == nil {
		response.Err(c, http.StatusServiceUnavailable, "模板库不可用")
		return
	}
	t, err := h.templates.Get(c.Param("id"))
	if err != nil {
		response.Err(c, http.StatusNotFound, "模板不存在")
		return
	}
	html := t.DemoHTML(c.Query("variant"))
	if v := c.Query("slide"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 1 {
			html = template.TrimToSlide(html, n)
		}
	}
	// 相对引用 → 对应静态路由的绝对引用（内置 = /templates/<id>/，用户 = /user-templates/<id>/）
	base := "/templates/" + t.ID + "/"
	if strings.HasPrefix(t.ID, "ut-") {
		base = "/user-templates/" + t.ID + "/"
	}
	html = strings.ReplaceAll(html, `href="style.css"`, `href="`+base+`style.css"`)
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
