package handler

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// 阶段0 的沙箱防线一半在响应头上：预览页永远带着这套策略，平时没人会去看它，
// 所以改动 header 时最容易被顺手删掉一条。这里把它钉死。
//
// 注意分工：CSP 保护的是"预览"；导出的单文件没有 CSP，
// 那条路径靠 service 层的输入消毒（sanitize.go）保护，两者不可互相替代。
func TestDeckPageHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	deckPageHeaders(c)

	csp := w.Header().Get("Content-Security-Policy")
	if csp == "" {
		t.Fatal("deck 响应必须带 CSP")
	}
	// 每一条都对应一类真实的执行/渗漏通道，漏一条就是开一个口子
	for _, want := range []string{
		"default-src 'self'",
		"script-src 'self'",                // 挡内联 <script> 与内联事件处理器
		"style-src 'self' 'unsafe-inline'", // 内联样式是允许的：AI 靠它排版
		"font-src 'self'",                  // 字体只许同源：随附字体在 /assets/fonts/，不许走 CDN
		"img-src 'self' data:",             // 只允许同源与内联图片
		"connect-src 'none'",               // 挡 fetch/XHR/WebSocket 外发
		"form-action 'none'",               // 表单提交是导航不是 fetch，不受 connect-src 管辖
		"object-src 'none'",
		"base-uri 'self'",
		"frame-ancestors 'self'",
	} {
		if !strings.Contains(csp, want) {
			t.Errorf("CSP 缺少 %q\n当前值: %s", want, csp)
		}
	}
	// script-src 一旦放开 unsafe-inline，deck 里的内联脚本与 on* 处理器就能执行，
	// 这正是消毒闸门要防的东西——两边不能同时松
	if strings.Contains(csp, "script-src 'self' 'unsafe-inline'") {
		t.Error("script-src 不应放开 unsafe-inline")
	}
	// ?token= 走 URL 的妥协靠它兜底：任何外发导航都不把带 token 的 URL 泄给第三方
	if got := w.Header().Get("Referrer-Policy"); got != "no-referrer" {
		t.Errorf("Referrer-Policy 应为 no-referrer，实际 %q", got)
	}
	if got := w.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Errorf("缺 nosniff: %q", got)
	}
	if got := w.Header().Get("Cache-Control"); got != "no-store" {
		t.Errorf("Cache-Control 应为 no-store: %q", got)
	}
}
