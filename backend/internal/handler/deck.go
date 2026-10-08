package handler

import (
	"bytes"
	"errors"
	"html-ppt/backend/internal/service/deck"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/response"
)

// ListDecks GET /api/decks —— 当前用户的 deck 列表（id + 标题）。
func (h *Handler) ListDecks(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	decks, err := h.decks.List(uid)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, decks)
}

// GetDeckFile GET /api/decks/:id/file —— 返回 deck.html 给 <iframe> 预览。
// 注意：这是文件流不是 JSON API，走 c.Data 原生写法，不经过 response 包。
// 错误统一映射为 404（非法 id、不存在的 id、别人的 id 都当"找不到"，
// 不泄露存在性）。deckPageHeaders 挂沙箱配套的安全响应头。
func (h *Handler) GetDeckFile(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	html, err := h.decks.PreviewHTML(uid, c.Param("id"))
	if err != nil {
		// 未选模板的 v2 deck：出友好占位页而不是 404 JSON——这个端点的消费者
		// 是 <iframe>，裸 JSON 会原样糊进页面（2026-09-30 用户实测反馈）。
		// 200 而非 404：语义是"还没到能预览的状态"，不是"找不到"。
		if errors.Is(err, deck.ErrNotInstantiated) {
			deckPageHeaders(c)
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(notInstantiatedHTML))
			return
		}
		response.Err(c, http.StatusNotFound, "deck 不存在")
		return
	}
	if c.Query("edit") == "1" {
		// 编辑模式：注入编辑器脚本（deck-editor-plan §4.3）。v1 deck 不注入——
		// 编辑弹窗靠 editor-ready 超时兜底提示"不支持编辑"。
		if h.decks.IsV2(c.Param("id")) {
			html = injectEditorScript(html)
		}
	}
	deckPageHeaders(c)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// notInstantiatedHTML 未实例化 deck 的预览占位页（配色贴近暗色主题的中性灰，
// 亮色主题下也不刺眼）。
const notInstantiatedHTML = `<!DOCTYPE html>
<html lang="zh-CN"><head><meta charset="utf-8"><title>尚未生成</title></head>
<body style="margin:0;display:flex;align-items:center;justify-content:center;height:100vh;background:#15181d;color:#8a929e;font:13px/1.7 system-ui,-apple-system,'Segoe UI',sans-serif;text-align:center;">
<div>这份文稿还没有生成页面<br><span style="font-size:12px;opacity:.7">在向导里选好模板就会开始生成</span></div>
</body></html>`

// maxEditBodyBytes 手动编辑保存的 body 上限。index.html 实测量级是几十 KB，
// 2MB 已经放宽了一个数量级还多；再大基本可以断定不是编辑器序列化产物。
const maxEditBodyBytes = 2 << 20

// SaveDeckFile PUT /api/decks/:id/file —— 编辑器全量保存（docs/deck-editor-plan.md §4.1）。
// body 是原始 HTML（编辑器在 iframe 内序列化 DOM 的产物，父页原样搬运），
// 不是 JSON，所以直接读 body 不走 binding。校验从宽：结构粗检 + 大小上限——
// 内容本来就是用户自己的稿子，语义层面的把关交给前端的 serialize 清理清单。
func (h *Handler) SaveDeckFile(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	body, err := io.ReadAll(http.MaxBytesReader(c.Writer, c.Request.Body, maxEditBodyBytes))
	if err != nil {
		response.Err(c, http.StatusRequestEntityTooLarge, "内容超限或读取失败")
		return
	}
	if len(bytes.TrimSpace(body)) == 0 {
		response.Err(c, http.StatusBadRequest, "空内容")
		return
	}
	if !bytes.Contains(body, []byte("<section")) {
		response.Err(c, http.StatusBadRequest, "内容不含 <section>，疑似非 deck HTML")
		return
	}
	if err := h.decks.SaveHTML(uid, c.Param("id"), string(body), c.Query("detail")); err != nil {
		// 与 GetDeckFile 同口径：不存在/别人的/非 v2 一律 404，不泄露存在性
		response.Err(c, http.StatusNotFound, err.Error())
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func deckPageHeaders(c *gin.Context) {
	c.Header("Content-Security-Policy",
		"default-src 'self'; "+
			"script-src 'self'; "+
			"style-src 'self' 'unsafe-inline'; "+
			// font-src 只允许同源：字体必须自托管（web/assets/fonts/）。
			// 严格说它是 default-src 'self' 的复述，写出来是为了让"字体不走 CDN"
			// 这条决定出现在强制执行它的地方——将来谁想往 <head> 里加一行
			// Google Fonts，撞上的第一堵墙就是这里（而且它会顺便把访问者的 IP
			// 送给第三方，对一个本地工具没有理由）。
			"font-src 'self'; "+
			"img-src 'self' data:; "+
			"connect-src 'none'; "+
			// form-action 不受 connect-src 管辖（表单提交是导航不是 fetch），
			// 不限的话隐藏表单自动提交外站是一条数据渗出通道
			"form-action 'none'; "+
			"object-src 'none'; "+
			"base-uri 'self'; "+
			"frame-ancestors 'self'")
	// ?token= 走 URL 的妥协（见 auth.go）：no-referrer 保证任何外发导航
	// 都不会把带 token 的完整 URL 泄给第三方
	c.Header("Referrer-Policy", "no-referrer")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Cache-Control", "no-store")
}

func (h *Handler) DeleteDeckVersion(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}

	deckID := c.Param("id")
	version := c.Param("version")
	if deckID == "" || version == "" {
		response.ParameterErr(c)
		return
	}

	if err := h.decks.EnsureOwner(uid, deckID); err != nil {
		response.Err(c, http.StatusNotFound, "deck 不存在")
		return
	}
	if err := h.decks.DeleteVersion(uid, deckID, version); err != nil {
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *Handler) ClearDeckHistory(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	deckID := c.Param("id")
	if deckID == "" {
		response.ParameterErr(c)
		return
	}

	if err := h.decks.EnsureOwner(uid, deckID); err != nil {
		response.Err(c, http.StatusNotFound, "deck 不存在 err:"+err.Error())
		return
	}
	n, err := h.decks.ClearHistory(uid, deckID)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, gin.H{"deleted": n})
}

func (h *Handler) ListDeckHistory(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}

	deckID := c.Param("id")
	if deckID == "" {
		response.ParameterErr(c)
		return
	}

	if err := h.decks.EnsureOwner(uid, deckID); err != nil {
		response.Err(c, http.StatusNotFound, "deck 不存在 err:"+err.Error())
		return
	}

	versions, err := h.decks.ListVersions(uid, deckID)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, versions)
}

func (h *Handler) RestoreDeckVersion(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	deckID := c.Param("id")
	version := c.Param("version")
	if deckID == "" || version == "" {
		response.ParameterErr(c)
		return
	}

	if err := h.decks.EnsureOwner(uid, deckID); err != nil {
		response.Err(c, http.StatusNotFound, "deck 不存在 err:"+err.Error())
		return
	}
	if err := h.decks.RestoreVersionV2(uid, deckID, version); err != nil {
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"restored": version})
}
