package router

import (
	"path/filepath"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/config"
	"html-ppt/backend/internal/handler"
	"html-ppt/backend/internal/middleware"
)

// New 装配路由。所有依赖已注入 Handler，这里只做 URL → 方法的映射。
// 路由分两组：公开（登录注册）和受保护（JWT 中间件把守）。
// 新增受保护接口时挂到 guarded 组，新增公开接口时想清楚它为什么不需要登录。
func New(cfg *config.Config, h *handler.Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.CORS(cfg.Server.AllowOrigins))

	api := r.Group("/api")
	{
		api.GET("/health", h.Health)

		// 公开：验证码 / 注册 / 登录
		authg := api.Group("/auth")
		{
			authg.POST("/code", h.RequestCode)
			authg.POST("/register", h.Register)
			authg.POST("/login", h.Login)
		}

		// 受保护：deck 全部读写 + agent 对话 + 用户摘要 / API Key 设置
		guarded := api.Group("", middleware.Auth(cfg.Auth.JWTSecret))
		{
			guarded.GET("/auth/me", h.Me)
			guarded.POST("/auth/apikey", h.SetAPIKey)
			guarded.GET("/decks", h.ListDecks)
			guarded.GET("/decks/:id/file", h.GetDeckFile)
			guarded.POST("/chat", h.Chat)
			guarded.POST("/chat/answer",h.AskUser)

			guarded.GET("/decks/:id/history",h.ListDeckHistory)
			guarded.POST("/decks/:id/history/:version/restore",h.RestoreDeckVersion)
			guarded.DELETE("/decks/:id/history/:version",h.DeleteDeckVersion)
			guarded.DELETE("/decks/:id/history",h.ClearDeckHistory)
		}
	}

	// reveal.js 等静态资源：deck.html 里的 <link>/<script> 引用 /assets/...
	// （静态资源不挂 Auth：浏览器加载 <script src> 时不会带 Authorization 头）
	r.Static("/assets", cfg.Assets.Dir)

	// SSE 测试台（同源访问，无 CORS 问题）：http://localhost:8080/chat-test
	r.StaticFile("/chat-test", filepath.Join(filepath.Dir(cfg.Assets.Dir), "chat-test.html"))
	return r
}
