package router

import (
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
