package router

import (
	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/config"
	"html-ppt/backend/internal/handler"
	"html-ppt/backend/internal/middleware"
)

// New 装配路由。所有依赖已注入 Handler，这里只做 URL → 方法的映射。
func New(cfg *config.Config, h *handler.Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.CORS(cfg.Server.AllowOrigins))

	api := r.Group("/api")
	{
		api.GET("/health", h.Health)
		api.GET("/decks", h.ListDecks)
		api.GET("/decks/:id/file", h.GetDeckFile)
		api.POST("/chat", h.Chat)
	}

	// reveal.js 等静态资源：deck.html 里的 <link>/<script> 引用 /assets/...
	r.Static("/assets", cfg.Assets.Dir)

	// SSE 测试台（同源访问，无 CORS 问题）：http://localhost:8080/chat-test
	r.StaticFile("/chat-test", "../../web/chat-test.html")
	return r
}
