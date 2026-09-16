package router

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/config"
	"html-ppt/backend/internal/handler"
	"html-ppt/backend/internal/trace"
	"html-ppt/backend/internal/vision"
)

// 视觉审查那条一次性取页通道必须真的注册上。
//
// 为什么值得单独守：它漏注册的表现是"404"，而 agent 那边会把 404 归到
// "渲染/量测失败"、只打一行 warn —— 功能整个不工作，日志里看不出是路由的问题。
// 这种事真的发生过：handler.RenderDeck 写好了、没人注册，编译全绿、测试全绿。
// 而 vision_review.go 里拼 URL 用的就是 /api/render/，两处必须逐字一致。
func TestRenderRouteIsRegistered(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// 路由注册阶段不会调用任何 handler 方法，所以依赖可以全传 nil；
	// 但 config 得给个真实存在的 assets 目录，否则 gin 的 Static 会在注册时报错。
	h := handler.New(nil, nil, nil, nil, &vision.Grants{}, trace.NewStore(t.TempDir()))
	engine := New(&config.Config{
		Assets: config.Assets{Dir: t.TempDir()},
		Auth:   config.Auth{JWTSecret: "test"},
	}, h)

	found := false
	for _, r := range engine.Routes() {
		if r.Method == "GET" && r.Path == "/api/render/:nonce" {
			found = true
		}
	}
	if !found {
		t.Error("GET /api/render/:nonce 没注册：视觉审查会静默失效（404 → 只打一行 warn）。" +
			"它必须公开（无头浏览器导航带不了 Authorization 头），安全性靠一次性 nonce")
	}
}

// 观测接口必须真的注册上。守它的理由和上面那条一样：
// 漏注册的表现是观测页一片空白，而页面那边只会显示"没有记录"——
// 你会先去怀疑 features.trace 开关、怀疑没写盘，最后才想到是路由。
//
// 同时守住取图那条的**参数名**：页面里是按 /img/<name> 拼 URL 的，
// 参数名一改（比如 :name 变成 :img）就全是 404，而 PNG 加载失败在页面上
// 只表现为一个空框——和"截图没存下来"长得一模一样。
func TestTraceRoutesAreRegistered(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := handler.New(nil, nil, nil, nil, &vision.Grants{}, trace.NewStore(t.TempDir()))
	engine := New(&config.Config{
		Assets: config.Assets{Dir: t.TempDir()},
		Auth:   config.Auth{JWTSecret: "test"},
	}, h)

	want := map[string]bool{
		"GET /api/traces":                               false,
		"GET /api/traces/:sessionID/:runID":             false,
		"GET /api/traces/:sessionID/:runID/events/:seq": false,
		"GET /api/traces/:sessionID/:runID/img/:name":   false,
		"GET /api/traces/:sessionID/:runID/export":      false,
	}
	for _, r := range engine.Routes() {
		key := r.Method + " " + r.Path
		if _, ok := want[key]; ok {
			want[key] = true
		}
	}
	for route, ok := range want {
		if !ok {
			t.Errorf("%s 没注册：观测页对应功能会静默失效", route)
		}
	}

	// 观测台静态页也要在（它是免登录的，所以不在 /api 组里）
	foundPage := false
	for _, r := range engine.Routes() {
		if r.Method == "GET" && r.Path == "/trace" {
			foundPage = true
		}
	}
	if !foundPage {
		t.Error("GET /trace 没注册：观测台打不开（与 /chat-test 同样挂在引擎上，免登录）")
	}
}

// 对话历史（回放）那两条读取接口必须真的注册上。
//
// 它们漏注册的表现特别有迷惑性：前端"刷新后对话消失了"——你会先去怀疑
// 前端状态管理、怀疑投影逻辑，最后才想到是路由没挂。而 New() 在路由冲突时
// 会直接 panic，所以这个测试顺带守住"新加的 /chat/sessions/:id/* 没有和
// /chat/answer、/chat/pending 撞参数名"这件事。
func TestChatHistoryRoutesAreRegistered(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := handler.New(nil, nil, nil, nil, &vision.Grants{}, trace.NewStore(t.TempDir()))
	engine := New(&config.Config{
		Assets: config.Assets{Dir: t.TempDir()},
		Auth:   config.Auth{JWTSecret: "test"},
	}, h)

	want := map[string]bool{
		"GET /api/decks/:id/chat/sessions":    false,
		"GET /api/chat/sessions/:id/messages": false,
	}
	for _, r := range engine.Routes() {
		key := r.Method + " " + r.Path
		if _, ok := want[key]; ok {
			want[key] = true
		}
	}
	for route, ok := range want {
		if !ok {
			t.Errorf("%s 没注册：对话历史回放会静默失效（前端刷新后看不到之前的对话）", route)
		}
	}
}

// token 换发接口必须存在**且在受保护组里**。
//
// 它挂在 /auth 前缀下，而 /auth 的 code/register/login 都是公开组——手滑把
// refresh 也注册成公开的话，"没有登录态也能换发新 token"，30 天滑续就变成了
// 永久登录。所以这里不只查注册，还实际打一个无凭证请求，验证它真的被 JWT
// 中间件拦下来（401 而不是 200/404）。
func TestAuthRefreshRouteIsRegisteredAndGuarded(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := handler.New(nil, nil, nil, nil, &vision.Grants{}, trace.NewStore(t.TempDir()))
	engine := New(&config.Config{
		Assets: config.Assets{Dir: t.TempDir()},
		Auth:   config.Auth{JWTSecret: "test"},
	}, h)

	registered := false
	for _, r := range engine.Routes() {
		if r.Method == "POST" && r.Path == "/api/auth/refresh" {
			registered = true
		}
	}
	if !registered {
		t.Fatal("POST /api/auth/refresh 没注册：前端每天上线静默续期的通道不存在")
	}

	w := performRequest(engine, "POST", "/api/auth/refresh")
	if w.Code != 401 {
		t.Errorf("无凭证调用 /api/auth/refresh 应被 JWT 中间件拦成 401，得到 %d——它必须在受保护组里", w.Code)
	}
}

// performRequest 走完整个中间件链的最小请求助手（handler 全是 nil 也不会碰到：
// 无凭证请求在 JWT 中间件就被拒了）。
func performRequest(engine *gin.Engine, method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	engine.ServeHTTP(w, req)
	return w
}
