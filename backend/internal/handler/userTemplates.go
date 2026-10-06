package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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
// :id 取 "_blank" 时走空白脚手架（usertpl.Fork 内分派到 CreateBlank）——
// 不为它单开静态路由段，避开 gin 静态/参数同位 panic（见 RecentSessions 注释）。
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
		if h.templates != nil {
			if cv, err := h.templates.ContentVersion(r.ID); err == nil && len(cv) >= 8 {
				v.Thumb = cv[:8]
			}
		}
		views = append(views, v)
	}
	etagJSON(c, views)
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
			// 封面缩略图的内容版本前缀：与我的模板清单同口径，前端拼
			// /templates/:id/thumb?v= 拉快照（社区模板已公开，鉴权路由任何
			// 登录用户可读；tplthumb 对 ut-* 走带 token 的 demo 路由渲染）
			if h.templates != nil {
				if v, err := h.templates.ContentVersion(id); err == nil && len(v) >= 8 {
					m["thumb"] = v[:8]
				}
				if t, err := h.templates.Get(id); err == nil && t.DemoPages > 0 {
					m["demo_pages"] = t.DemoPages
				}
			}
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
	// Thumb 缩略图内容版本前缀（registry.ContentVersion，编辑重挂即变）——
	// 卡片 img 地址 /api/templates/<id>/thumb?v=<Thumb> 的 cache-bust 键
	Thumb string `json:"thumb,omitempty"`
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
	// 相对引用 → 受控资产端点（草稿收口后 style.css 不在公开路由上；
	// 子请求带不了鉴权头，token 从本请求的 query 原样带出）
	styleHref := "/api/user-templates/" + row.ID + "/assets/style.css"
	if token := c.Query("token"); token != "" {
		styleHref += "?token=" + token
	}
	html := strings.ReplaceAll(string(raw), `href="style.css"`, `href="`+styleHref+`"`)
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
	if err := h.usertpl.SaveIndexHTML(uid, c.Param("id"), string(body)); err != nil {
		if errors.Is(err, usertpl.ErrPublished) {
			response.Err(c, http.StatusConflict, err.Error())
			return
		}
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"ok": true})
}

// maxTemplateStyleBytes style.css 手动保存的 body 上限（对话 write_style 同限额）。
const maxTemplateStyleBytes = 512 << 10

// SaveUserTemplateStyle PUT /api/user-templates/:id/style —— 手动保存 style.css
// （工作台"样式"面板）。安全预检 + 记 edit 版本；published 拒绝（409）。
func (h *Handler) SaveUserTemplateStyle(c *gin.Context) {
	if h.usertpl == nil {
		response.Err(c, http.StatusServiceUnavailable, "用户模板不可用（需要数据库）")
		return
	}
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	body, err := io.ReadAll(http.MaxBytesReader(c.Writer, c.Request.Body, maxTemplateStyleBytes))
	if err != nil {
		response.Err(c, http.StatusRequestEntityTooLarge, "内容超限或读取失败")
		return
	}
	if strings.TrimSpace(string(body)) == "" {
		response.Err(c, http.StatusBadRequest, "空内容")
		return
	}
	if err := h.usertpl.SaveStyleCSS(uid, c.Param("id"), string(body)); err != nil {
		if errors.Is(err, usertpl.ErrPublished) {
			response.Err(c, http.StatusConflict, err.Error())
			return
		}
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"ok": true})
}

// GetUserTemplateStructure GET /api/user-templates/:id/structure —— 结构契约
// （版式面板数据源）：挂载态的版式清单（含骨架/指纹/数量契约）+ demo 页映射 +
// rules.md。面板展示的是"生成侧实际生效的契约"，所以读注册表快照而非磁盘文件。
func (h *Handler) GetUserTemplateStructure(c *gin.Context) {
	if h.usertpl == nil {
		response.Err(c, http.StatusServiceUnavailable, "用户模板不可用（需要数据库）")
		return
	}
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	v, err := h.usertpl.StructureContract(uid, c.Param("id"))
	if err != nil {
		if errors.Is(err, usertpl.ErrPublished) {
			response.Err(c, http.StatusConflict, err.Error())
			return
		}
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, v)
}

// UpdateUserTemplateLayout PUT /api/user-templates/:id/layouts/:layoutId —— 改
// 版式元数据（name/use/roles，可选补丁）。warning 非空 = 重挂未过（盘上已是新值、
// 已记历史），前端当提示展示。
func (h *Handler) UpdateUserTemplateLayout(c *gin.Context) {
	if h.usertpl == nil {
		response.Err(c, http.StatusServiceUnavailable, "用户模板不可用（需要数据库）")
		return
	}
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	var req usertpl.LayoutMetaPatch
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, http.StatusBadRequest, "请求体不是合法 JSON")
		return
	}
	warning, err := h.usertpl.UpdateLayoutMeta(uid, c.Param("id"), c.Param("layoutId"), req)
	if err != nil {
		if errors.Is(err, usertpl.ErrPublished) {
			response.Err(c, http.StatusConflict, err.Error())
			return
		}
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"ok": true, "warning": warning})
}

// ListUserTemplateHistory GET /api/user-templates/:id/history —— 版本历史（新→旧）。
func (h *Handler) ListUserTemplateHistory(c *gin.Context) {
	if h.usertpl == nil {
		response.Err(c, http.StatusServiceUnavailable, "用户模板不可用（需要数据库）")
		return
	}
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	versions, err := h.usertpl.ListUTVersions(uid, c.Param("id"))
	if err != nil {
		response.Err(c, http.StatusNotFound, err.Error())
		return
	}
	response.OK(c, versions)
}

// RestoreUserTemplateVersion POST /api/user-templates/:id/history/:version/restore —— 回滚。
// published 拒绝（409）；回滚本身记一条 restore 版本。
func (h *Handler) RestoreUserTemplateVersion(c *gin.Context) {
	if h.usertpl == nil {
		response.Err(c, http.StatusServiceUnavailable, "用户模板不可用（需要数据库）")
		return
	}
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	if err := h.usertpl.RestoreUTVersion(uid, c.Param("id"), c.Param("version")); err != nil {
		if errors.Is(err, usertpl.ErrPublished) {
			response.Err(c, http.StatusConflict, err.Error())
			return
		}
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"ok": true})
}

// DeleteUserTemplateVersion DELETE /api/user-templates/:id/history/:version —— 删单版。
func (h *Handler) DeleteUserTemplateVersion(c *gin.Context) {
	if h.usertpl == nil {
		response.Err(c, http.StatusServiceUnavailable, "用户模板不可用（需要数据库）")
		return
	}
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	if err := h.usertpl.DeleteUTVersion(uid, c.Param("id"), c.Param("version")); err != nil {
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"ok": true})
}

// ClearUserTemplateHistory DELETE /api/user-templates/:id/history —— 清空历史。
func (h *Handler) ClearUserTemplateHistory(c *gin.Context) {
	if h.usertpl == nil {
		response.Err(c, http.StatusServiceUnavailable, "用户模板不可用（需要数据库）")
		return
	}
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	n, err := h.usertpl.ClearUTHistory(uid, c.Param("id"))
	if err != nil {
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"deleted": n})
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

// PublishUserTemplate POST /api/user-templates/:id/publish —— 发布（挂载校验+落状态）。
// 2026-09-28 门禁降级后秒级完成；渲染量测独立为 Checkup（POST /checkup）。
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

// CheckupUserTemplate POST /api/user-templates/:id/checkup —— 质量体检。
// headless 实拍 demo 量测溢出/填充率/最小字号，报告返回并落 publish_report；
// 不改状态（draft/published 都能跑），失败只影响本次体检不失败发布。
func (h *Handler) CheckupUserTemplate(c *gin.Context) {
	if h.usertpl == nil {
		response.Err(c, http.StatusServiceUnavailable, "用户模板不可用（需要数据库）")
		return
	}
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	report, err := h.usertpl.Checkup(c.Request.Context(), uid, c.Param("id"))
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

// maxChatBodyBytes 定制对话请求体上限。文本很小，大头是附图：base64 比原始
// 字节膨胀 ~4/3，3 张 × 4MB 原图的极限在 16MB，20MB 留了余量。
const maxChatBodyBytes = 20 << 20

// bindCustomizeChat 两个定制对话端点共用的请求绑定：文本与附图至少一个非空，
// 图片的合法性（形态/大小/真实格式）在服务层 normalizeUserImages 里统一校验。
func bindCustomizeChat(c *gin.Context) (message string, images []string, ok bool) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxChatBodyBytes)
	var body struct {
		Message string   `json:"message"`
		Images  []string `json:"images"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Err(c, http.StatusRequestEntityTooLarge, "请求超限或格式错误（含图片时整体不超过 20MB）")
		return "", nil, false
	}
	if strings.TrimSpace(body.Message) == "" && len(body.Images) == 0 {
		response.Err(c, http.StatusBadRequest, "消息不能为空（文本与图片至少一个）")
		return "", nil, false
	}
	return body.Message, body.Images, true
}

// CustomizeUserTemplate POST /api/user-templates/:id/chat —— 对话定制（同步版）。
// 用 agent 服务的 LLM 接入（BYOK 优先）跑受限工具循环。前端已改走 /chat/stream
// （SSE，工具气泡实时可见）；这个端点保留给 curl/脚本/降级路径，emit 传 nil。
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
	message, images, ok := bindCustomizeChat(c)
	if !ok {
		return
	}
	cl := h.agent.CustomizeLLMFor(c.Request.Context())
	reply, err := h.usertpl.Customize(c.Request.Context(), uid, c.Param("id"), message, images, usertpl.LLM{Client: cl.Client, Model: cl.Model}, nil)
	if err != nil {
		response.Err(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"reply": reply})
}

// CustomizeUserTemplateStream POST /api/user-templates/:id/chat/stream —— 定制对话 SSE。
// 帧格式与 deck 对话一致（event: <type>\ndata: <json>\n\n）；事件类型是
// usertpl.CustEvent 的窄集合（tool_start/tool_progress/tool_done/delta/done/error）。
// 客户端断开时 request context 取消，LLM 流与工具循环随之停止，不再白烧 token。
func (h *Handler) CustomizeUserTemplateStream(c *gin.Context) {
	if h.usertpl == nil {
		response.Err(c, http.StatusServiceUnavailable, "用户模板不可用（需要数据库）")
		return
	}
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	message, images, ok := bindCustomizeChat(c)
	if !ok {
		return
	}

	cl := h.agent.CustomizeLLMFor(c.Request.Context())
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	ch := make(chan usertpl.CustEvent, 16)
	go func() {
		defer close(ch)
		_, err := h.usertpl.Customize(c.Request.Context(), uid, c.Param("id"), message, images,
			usertpl.LLM{Client: cl.Client, Model: cl.Model},
			func(ev usertpl.CustEvent) error {
				select {
				case ch <- ev:
					return nil
				case <-c.Request.Context().Done():
					return c.Request.Context().Err()
				}
			})
		if err != nil {
			log.Printf("定制对话流失败 err:%v", err)
			select {
			case ch <- usertpl.CustEvent{Type: usertpl.CustEvError, Content: err.Error()}:
			case <-c.Request.Context().Done():
			}
		}
	}()

	c.Stream(func(w io.Writer) bool {
		ev, ok := <-ch
		if !ok {
			return false
		}
		data, _ := json.Marshal(ev)
		fmt.Fprintf(w, "event: %s\ndata: %s\n\n", ev.Type, data)
		return true
	})
}
