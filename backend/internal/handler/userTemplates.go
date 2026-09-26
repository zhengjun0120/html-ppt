package handler

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/response"
	"html-ppt/backend/internal/service/template"
	"html-ppt/backend/internal/service/usertpl"
	"html-ppt/backend/internal/store"
)

// 用户自定义模板（plan-v3 B）：fork → 对话定制 → 门禁发布 → 社区使用。
// usertpl 服务可能为 nil（数据库降级模式），所有入口先判空。

// templateVisuals 从注册表取模板画布与变体（画廊预览卡需要；未挂载/注册表不可用时返回零值）。
func (h *Handler) templateVisuals(id string) (*template.Canvas, []template.Variant) {
	if h.templates == nil {
		return nil, nil
	}
	t, err := h.templates.Get(id)
	if err != nil {
		return nil, nil
	}
	canvas := t.Canvas
	return &canvas, t.Variants
}

// ForkTemplate POST /api/templates/:id/fork —— 从内置模板克隆私有副本。
func (h *Handler) ForkTemplate(c *gin.Context) {
	if h.usertpl == nil {
		response.Err(c, http.StatusServiceUnavailable, "用户模板不可用（需要数据库）")
		return
	}
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	var body struct {
		Name string `json:"name"`
	}
	_ = c.ShouldBindJSON(&body)
	row, err := h.usertpl.Fork(uid, c.Param("id"), body.Name)
	if err != nil {
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, row)
}

// ListUserTemplates GET /api/user-templates —— 我的模板。
func (h *Handler) ListUserTemplates(c *gin.Context) {
	if h.usertpl == nil {
		response.Err(c, http.StatusServiceUnavailable, "用户模板不可用（需要数据库）")
		return
	}
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	rows, err := h.usertpl.List(uid)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}
	views := make([]userTemplateView, 0, len(rows))
	for _, r := range rows {
		v := userTemplateView{UserTemplate: r}
		v.Canvas, v.Variants = h.templateVisuals(r.ID)
		views = append(views, v)
	}
	response.OK(c, views)
}

// CommunityTemplates GET /api/templates/community —— 社区模板（公开已发布，带作者署名）。
func (h *Handler) CommunityTemplates(c *gin.Context) {
	if h.usertpl == nil {
		response.OK(c, []any{})
		return
	}
	rows, err := h.usertpl.Community()
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}
	for _, m := range rows {
		if id, ok := m["id"].(string); ok {
			m["canvas"], m["variants"] = h.templateVisuals(id)
		}
	}
	response.OK(c, rows)
}

// GetUserTemplate GET /api/user-templates/:id —— 详情（owner 或公开已发布可读）。
func (h *Handler) GetUserTemplate(c *gin.Context) {
	if h.usertpl == nil {
		response.Err(c, http.StatusServiceUnavailable, "用户模板不可用（需要数据库）")
		return
	}
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	row, err := h.usertpl.GetReadable(uid, c.Param("id"))
	if err != nil {
		response.Err(c, http.StatusNotFound, err.Error())
		return
	}
	canvas, variants := h.templateVisuals(row.ID)
	response.OK(c, userTemplateView{UserTemplate: *row, Canvas: canvas, Variants: variants})
}

// userTemplateView 用户模板行 + 从注册表补的视觉元数据（画廊预览卡需要画布与变体）。
type userTemplateView struct {
	store.UserTemplate
	Canvas   *template.Canvas   `json:"canvas,omitempty"`
	Variants []template.Variant `json:"variants,omitempty"`
}

// GetUserTemplateEditor GET /api/user-templates/:id/editor —— 编辑模式的 demo 读取端点
// （docs/deck-editor-plan.md §4.2）。
//
// 为什么不走公开静态路由：/user-templates/:id/index.html 是纯文件服务，没有
// 服务端钩子，注入编辑器脚本只能走受保护端点。除注入外，相对引用重写与
// PreviewTemplate 同一套；响应头与 deck 预览同一条（deckPageHeaders），
// 保证"编辑看到的页面"和"用户看到的页面"同构。仅属主可读——协作编辑不存在。
func (h *Handler) GetUserTemplateEditor(c *gin.Context) {
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
		response.Err(c, http.StatusNotFound, "模板不存在或无权编辑")
		return
	}
	raw, err := os.ReadFile(filepath.Join(h.usertpl.Dir(row.ID), "index.html"))
	if err != nil {
		response.Err(c, http.StatusNotFound, "模板 demo 缺失")
		return
	}
	// 相对引用 → 静态路由的绝对引用（与 PreviewTemplate 同一套）
	base := "/user-templates/" + row.ID + "/"
	html := strings.ReplaceAll(string(raw), `href="style.css"`, `href="`+base+`style.css"`)
	html = injectEditorScript(html)
	deckPageHeaders(c) // 含 Cache-Control: no-store
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// SaveUserTemplateFile PUT /api/user-templates/:id/file —— 编辑器全量保存模板
// index.html（docs/deck-editor-plan.md §4.2）。覆盖前服务端自动滚动备份最近 5 版；
// published 状态拒绝（409，先下架再编辑）。
func (h *Handler) SaveUserTemplateFile(c *gin.Context) {
	if h.usertpl == nil {
		response.Err(c, http.StatusServiceUnavailable, "用户模板不可用（需要数据库）")
		return
	}
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
	if len(strings.TrimSpace(string(body))) == 0 {
		response.Err(c, http.StatusBadRequest, "空内容")
		return
	}
	if !bytes.Contains(body, []byte("<section")) {
		response.Err(c, http.StatusBadRequest, "内容不含 <section>，疑似非模板 HTML")
		return
	}
	backups, err := h.usertpl.SaveIndexHTML(uid, c.Param("id"), string(body))
	if errors.Is(err, usertpl.ErrPublished) {
		response.Err(c, http.StatusConflict, err.Error())
		return
	}
	if err != nil {
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"ok": true, "backups": backups})
}

// UpdateUserTemplate PUT /api/user-templates/:id —— 改名/描述。
func (h *Handler) UpdateUserTemplate(c *gin.Context) {
	if h.usertpl == nil {
		response.Err(c, http.StatusServiceUnavailable, "用户模板不可用（需要数据库）")
		return
	}
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	_ = c.ShouldBindJSON(&body)
	if err := h.usertpl.UpdateMeta(uid, c.Param("id"), body.Name, body.Description); err != nil {
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"ok": true})
}

// DeleteUserTemplate DELETE /api/user-templates/:id —— 删除（含磁盘目录）。
func (h *Handler) DeleteUserTemplate(c *gin.Context) {
	if h.usertpl == nil {
		response.Err(c, http.StatusServiceUnavailable, "用户模板不可用（需要数据库）")
		return
	}
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	if err := h.usertpl.Delete(uid, c.Param("id")); err != nil {
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"ok": true})
}

// PublishUserTemplate POST /api/user-templates/:id/publish —— 跑门禁并公开。
// 同步执行（渲染关约 10-20s）；失败返回可读原因，状态落 failed。
func (h *Handler) PublishUserTemplate(c *gin.Context) {
	if h.usertpl == nil {
		response.Err(c, http.StatusServiceUnavailable, "用户模板不可用（需要数据库）")
		return
	}
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	report, err := h.usertpl.Publish(c.Request.Context(), uid, c.Param("id"))
	if err != nil {
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, report)
}

// UnpublishUserTemplate POST /api/user-templates/:id/unpublish —— 下架。
func (h *Handler) UnpublishUserTemplate(c *gin.Context) {
	if h.usertpl == nil {
		response.Err(c, http.StatusServiceUnavailable, "用户模板不可用（需要数据库）")
		return
	}
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	if err := h.usertpl.Unpublish(uid, c.Param("id")); err != nil {
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"ok": true})
}

// CustomizeUserTemplate POST /api/user-templates/:id/chat —— 对话定制。
// 用 agent 服务的 LLM 接入（BYOK 优先）跑受限工具循环（token/元数据级）。
func (h *Handler) CustomizeUserTemplate(c *gin.Context) {
	if h.usertpl == nil {
		response.Err(c, http.StatusServiceUnavailable, "用户模板不可用（需要数据库）")
		return
	}
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	var body struct {
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Message == "" {
		response.Err(c, http.StatusBadRequest, "message 不能为空")
		return
	}
	cl := h.agent.CustomizeLLMFor(c.Request.Context())
	reply, err := h.usertpl.Customize(c.Request.Context(), uid, c.Param("id"), body.Message, usertpl.LLM{Client: cl.Client, Model: cl.Model})
	if err != nil {
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"reply": reply})
}
