package router

import (
	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/config"
	"html-ppt/backend/internal/handler"
	"html-ppt/backend/internal/middleware"
	"html-ppt/backend/internal/service/deck"
	"html-ppt/backend/internal/store"
)

// New 装配路由。依赖（store、service）从参数注入，handler 只是"HTTP 翻译层"。
func New(cfg *config.Config, st *store.Store, deckSvc *deck.Service) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.CORS(cfg.Server.AllowOrigins))

	api := r.Group("/api")
	{
		api.GET("/health", handler.Health(st))
		api.GET("/decks", handler.ListDecks(deckSvc))
		api.GET("/decks/:id/file", handler.GetDeckFile(deckSvc))
		api.POST("/chat", handler.Chat())
	}

	// reveal.js 等静态资源：deck.html 里的 <link>/<script> 引用 /assets/...
	r.Static("/assets", cfg.Assets.Dir)
	return r
}
