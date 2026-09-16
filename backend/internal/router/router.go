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

		// 视觉审查的一次性取页通道：**必须公开**——无头浏览器是"导航"到它的，
		// 导航带不了 Authorization 头。安全性靠一次性 nonce（见 vision/grant.go）：
		// 24 字节随机、只绑一份 deck、取到即废、2 分钟过期；响应头与预览完全一致
		// （同一个 deckPageHeaders），否则"审查看到的页面"和"用户看到的页面"不是一个东西。
		// 路径前缀改了的话，agent/vision_review.go 里拼 URL 那行必须一起改。
		api.GET("/render/:nonce", h.RenderDeck)

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
			// 登录态换发：每用户每天一次（额度记在用户行上，跨设备共享）。
			// 前端"每天上线"时静默续一个 TTL，正常用户从此几乎不再见到登录页
			guarded.POST("/auth/refresh", h.Refresh)
			guarded.POST("/auth/apikey", h.SetAPIKey)
			guarded.GET("/decks", h.ListDecks)
			guarded.GET("/decks/:id/file", h.GetDeckFile)
			guarded.POST("/chat", h.Chat)
			guarded.POST("/chat/answer", h.AskUser)
			// 暂停中的提问（页面刷新后重建提问卡片用；没有则 questions 为空串）
			guarded.GET("/chat/pending", h.PendingAsk)

			// 对话历史的读取端（回放）。会话按 deck 组织：前者列出一个 deck 名下的
			// 全部会话（进工作台先恢复最近对话/切会话用），后者返回单个会话的可回放
			// 消息列表（前端把历史灌回消息流组件，刷新后原样恢复）。
			// 只读，不经过 agent 闸门——暂停中的会话照样能看历史。
			guarded.GET("/decks/:id/chat/sessions", h.ListDeckSessions)
			guarded.GET("/chat/sessions/:id/messages", h.GetSessionMessages)

			guarded.GET("/decks/:id/history", h.ListDeckHistory)
			guarded.POST("/decks/:id/history/:version/restore", h.RestoreDeckVersion)
			guarded.DELETE("/decks/:id/history/:version", h.DeleteDeckVersion)
			guarded.DELETE("/decks/:id/history", h.ClearDeckHistory)

			// 观测记录读取（观测页 /trace 用）。归属复核在 trace.Store 内部做：
			// 不属于你的 run 一律 404，且与"不存在"同一个响应（不泄露存在性）。
			// 取图那条特殊：<img src> 带不了 Authorization 头，靠 ?token= 回退，
			// 与 deck 预览是同一条已知妥协（见 middleware/auth.go 的说明）。
			guarded.GET("/traces", h.ListTraces)
			guarded.GET("/traces/:sessionID/:runID", h.GetTraceRun)
			guarded.GET("/traces/:sessionID/:runID/events/:seq", h.GetTraceEvent)
			guarded.GET("/traces/:sessionID/:runID/img/:name", h.GetTraceImage)
			guarded.GET("/traces/:sessionID/:runID/export", h.ExportTrace)
		}
	}

	// reveal.js 等静态资源：deck.html 里的 <link>/<script> 引用 /assets/...
	// （静态资源不挂 Auth：浏览器加载 <script src> 时不会带 Authorization 头）
	assets := r.Group("", revalidateStatic())
	assets.Static("/assets", cfg.Assets.Dir)

	// SSE 测试台（同源访问，无 CORS 问题）：http://localhost:8080/chat-test
	r.StaticFile("/chat-test", filepath.Join(filepath.Dir(cfg.Assets.Dir), "chat-test.html"))
	// 观测台：http://localhost:8080/trace（免登录的静态页，令牌从 localStorage 取，
	// 与 chat-test 同一个键；数据接口全部在 guarded 组里）
	r.StaticFile("/trace", filepath.Join(filepath.Dir(cfg.Assets.Dir), "trace.html"))
	return r
}

// revalidateStatic 让浏览器每次回源校验静态资源，而不是凭启发式缓存直接复用旧副本。
//
// 为什么必须有：静态文件响应默认只带 Last-Modified、不带 Cache-Control，Chromium 会按
// "距今时长的 10%" 当作新鲜期——实测改完 components.css 刷新页面，浏览器仍然直接用两周前的
// 旧副本（stylesheets 里只剩 12 条规则、新规则一条都没有），因为压根没触发校验。
// 组件库和主题是给用户调视觉用的，改了看不到变化是最误导人的一种失败。
//
// 刻意不用"长缓存 + 内容哈希"：deck.html 引用 /assets/xxx.css 时不带指纹，一旦长缓存就再也换不掉。
// no-cache 不等于不缓存，只是每次都问一句，命中 304 时开销极小。
func revalidateStatic() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache")
		c.Next()
	}
}
