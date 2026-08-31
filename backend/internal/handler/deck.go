package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/service/deck"
)

// ListDecks GET /api/decks —— 列出所有 deck（id + 标题），前端左侧列表用。
func ListDecks(svc *deck.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		decks, err := svc.List()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, decks)
	}
}

// GetDeckFile GET /api/decks/:id/file —— 返回 deck.html 给 <iframe> 预览。
// 错误统一映射为 404（骨架阶段从简：非法 id 和不存在的 id 都当"找不到"）。
func GetDeckFile(svc *deck.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		html, err := svc.GetHTML(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	}
}
