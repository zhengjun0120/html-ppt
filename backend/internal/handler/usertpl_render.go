package handler

// 用户模板 demo 的受控伺服（docs/user-template-history-plan.md §3.5 草稿收口）。
//
// 替代原 utStatic.StaticFS 公开路由。三个入口、一条白名单伺服逻辑：
//   - /user-templates/:id/*filepath —— 公开路由（AuthOptional）：published 免登录
//     （社区画廊 iframe 带不了鉴权头），draft/failed/publishing 仅属主（token）；
//   - /api/user-templates/:id/demo —— 鉴权工作台/编辑弹窗预览（owner）；
//   - /api/user-template-render/:nonce —— 发布门禁的 headless 渲染（一次性授权）。
//
// 白名单只有 index.html 与 style.css：template.json/layouts.md/rules.md 是模板
// "源码"，demo 渲染用不到，永不伺服（修现状问题：StaticFS 全目录暴露）。
// 相对引用 href="style.css" 在各入口下重写到对应的受控资产地址（子请求带不了
// 鉴权头，token 从入口请求的 query 原样带出）。

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/response"
)

// serveUTFile 白名单伺服：只认 index.html 与 style.css，其余一律 404。
func (h *Handler) serveUTFile(c *gin.Context, id, name, styleHref string) {
	switch name {
	case "index.html":
		raw, err := os.ReadFile(filepath.Join(h.usertpl.Dir(id), "index.html"))
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		html := raw
		if styleHref != "" {
			html = []byte(strings.ReplaceAll(string(raw), `href="style.css"`, `href="`+styleHref+`"`))
		}
		deckPageHeaders(c)
		c.Data(http.StatusOK, "text/html; charset=utf-8", html)
	case "style.css":
		raw, err := os.ReadFile(filepath.Join(h.usertpl.Dir(id), "style.css"))
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.Header("Cache-Control", "no-store")
		c.Data(http.StatusOK, "text/css; charset=utf-8", raw)
	default:
		c.Status(http.StatusNotFound)
	}
}

// UserTemplatePublicFile GET /user-templates/:id/*filepath —— 受控公开路由。
// 挂 AuthOptional（不拦截）：published 无条件伺服；其他状态要求属主，
// 不满足按 404 处理（与 GetDeckFile 同一"不泄露存在性"口径）。
func (h *Handler) UserTemplatePublicFile(c *gin.Context) {
	if h.usertpl == nil {
		c.Status(http.StatusNotFound)
		return
	}
	row, err := h.usertpl.Peek(c.Param("id"))
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	if row.Status != "published" {
		uid, ok := authctx.UserID(c.Request.Context())
		if !ok || uid != row.UserID {
			c.Status(http.StatusNotFound)
			return
		}
	}
	// 草稿态的 style.css 子请求带不了 token（iframe 的 query 不遗传）——
	// 把入口请求的 token 拼进重写后的 href，published 下无 token 则不重写
	name := strings.TrimPrefix(c.Param("filepath"), "/")
	styleHref := ""
	if name == "index.html" && row.Status != "published" {
		if token := c.Query("token"); token != "" {
			styleHref = "/user-templates/" + row.ID + "/style.css?token=" + token
		}
	}
	h.serveUTFile(c, row.ID, name, styleHref)
}

// GetUserTemplateDemo GET /api/user-templates/:id/demo —— 鉴权预览（owner）。
// 定制工作台与手动编辑弹窗用：草稿也要能看，走归属校验而非公开路由。
func (h *Handler) GetUserTemplateDemo(c *gin.Context) {
	if h.usertpl == nil {
		response.Err(c, http.StatusServiceUnavailable, "用户模板不可用（需要数据库）")
		return
	}
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	row, err := h.usertpl.GetOwned(uid, c.Param("id"))
	if err != nil {
		response.Err(c, http.StatusNotFound, "模板不存在")
		return
	}
	styleHref := "/api/user-templates/" + row.ID + "/assets/style.css"
	if token := c.Query("token"); token != "" {
		styleHref += "?token=" + token
	}
	h.serveUTFile(c, row.ID, "index.html", styleHref)
}

// GetUserTemplateAsset GET /api/user-templates/:id/assets/:name —— 鉴权资产
// （demo 端点重写出来的子请求）。白名单同上。
func (h *Handler) GetUserTemplateAsset(c *gin.Context) {
	if h.usertpl == nil {
		response.Err(c, http.StatusServiceUnavailable, "用户模板不可用（需要数据库）")
		return
	}
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	row, err := h.usertpl.GetOwned(uid, c.Param("id"))
	if err != nil {
		response.Err(c, http.StatusNotFound, "模板不存在")
		return
	}
	h.serveUTFile(c, row.ID, c.Param("name"), "")
}

// RenderUserTemplateDemo GET /api/user-template-render/:nonce/*filepath ——
// 发布门禁的 headless 渲染端点。nonce 由 Publish 签发（vision.Grants.Peek 语义：
// 渲染要跨多次导航拉 index+css，不能一次即焚），TTL 2 分钟过期自然作废。
func (h *Handler) RenderUserTemplateDemo(c *gin.Context) {
	if h.usertpl == nil {
		c.Status(http.StatusNotFound)
		return
	}
	_, ref, ok := h.renderGrants.Peek(c.Param("nonce"))
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}
	id, found := strings.CutPrefix(ref, "ut:")
	if !found {
		c.Status(http.StatusNotFound)
		return
	}
	if _, err := h.usertpl.Peek(id); err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	name := strings.TrimPrefix(c.Param("filepath"), "/")
	if name == "" {
		name = "index.html" // 无尾路径（gin 通配重定向后的 "/" 形态）兜底
	}
	styleHref := "/api/user-template-render/" + c.Param("nonce") + "/style.css"
	h.serveUTFile(c, id, name, styleHref)
}
