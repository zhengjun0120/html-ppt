package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/response"
)

// ListDecks GET /api/decks —— 列出所有 deck（id + 标题），前端左侧列表用。
func (h *Handler) ListDecks(c *gin.Context) {
	decks, err := h.decks.List()
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, decks)
}

// GetDeckFile GET /api/decks/:id/file —— 返回 deck.html 给 <iframe> 预览。
// 注意：这是文件流不是 JSON API，走 c.Data 原生写法，不经过 response 包。
// 错误统一映射为 404（骨架阶段从简：非法 id 和不存在的 id 都当"找不到"）。
func (h *Handler) GetDeckFile(c *gin.Context) {
	html, err := h.decks.GetHTML(c.Param("id"))
	if err != nil {
		response.Err(c, http.StatusNotFound, err.Error())
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
