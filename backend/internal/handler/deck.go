package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/response"
)

// ListDecks GET /api/decks —— 当前用户的 deck 列表（id + 标题）。
func (h *Handler) ListDecks(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	decks, err := h.decks.List(uid)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, decks)
}

// GetDeckFile GET /api/decks/:id/file —— 返回 deck.html 给 <iframe> 预览。
// 注意：这是文件流不是 JSON API，走 c.Data 原生写法，不经过 response 包。
// 错误统一映射为 404（非法 id、不存在的 id、别人的 id 都当"找不到"，
// 不泄露存在性）。deckPageHeaders 挂沙箱配套的安全响应头。
func (h *Handler) GetDeckFile(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	html, err := h.decks.GetHTML(uid, c.Param("id"))
	if err != nil {
		response.Err(c, http.StatusNotFound, "deck 不存在")
		return
	}
	deckPageHeaders(c)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func deckPageHeaders(c *gin.Context) {
	c.Header("Content-Security-Policy",
		"default-src 'self'; "+
			"script-src 'self'; "+
			"style-src 'self' 'unsafe-inline'; "+
			"img-src 'self' data:; "+
			"connect-src 'none'; "+
			"object-src 'none'; "+
			"base-uri 'self'; "+
			"frame-ancestors 'self'")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Cache-Control", "no-store")
}
