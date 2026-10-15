package handler

// deck-v2 的管线端点：大纲（读/改/确认）、模板选择、生成触发、deck 元数据。
//
// 阶段守卫在 service 层（transitionStage / SaveOutline / SelectTemplate），
// handler 只做参数绑定与错误映射：阶段冲突 → 409，版本冲突 → 409 带最新版，
// 归属问题 → 404（与 v1 同一条"不泄露存在性"纪律）。

import (
	"context"
	"errors"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/agent"
	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/export"
	"html-ppt/backend/internal/response"
	"html-ppt/backend/internal/service/deck"
	"html-ppt/backend/internal/service/tplsuggest"
)

func (h *Handler) uid(c *gin.Context) (uint, bool) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return 0, false
	}
	return uid, true
}

// mapDeckErr v2 管线错误的统一映射：版本冲突与阶段守卫都是 409（用户可据此行动），
// 其余按 400。
func mapDeckErr(c *gin.Context, err error) {
	var conflict deck.OutlineConflict
	if errors.As(err, &conflict) {
		c.JSON(http.StatusConflict, gin.H{
			"code": "outline_conflict", "latest_version": conflict.LatestVersion, "message": err.Error(),
		})
		return
	}
	var mismatch deck.StageMismatch
	if errors.As(err, &mismatch) {
		// 重复确认/重复提交落在已前进的阶段上：不是请求错误，是状态冲突
		response.Err(c, http.StatusConflict, err.Error())
		return
	}
	var locked agent.ErrStageLocked
	if errors.As(err, &locked) {
		response.Err(c, http.StatusConflict, err.Error())
		return
	}
	response.Err(c, http.StatusBadRequest, err.Error())
}

// GetDeckV2Meta GET /api/decks/:id/meta —— v2 元数据 + 大纲（前端向导的状态源）。
func (h *Handler) GetDeckV2Meta(c *gin.Context) {
	uid, ok := h.uid(c)
	if !ok {
		return
	}
	df, err := h.decks.GetDeckV2(uid, c.Param("id"))
	if err != nil {
		response.Err(c, http.StatusNotFound, "deck 不存在")
		return
	}
	out := gin.H{"deck": df}
	if o, err := h.decks.ReadOutline(uid, df.ID); err == nil {
		out["outline"] = o
	}
	// 模板的变体清单（画廊选中后展示色板用）
	if h.templates != nil && df.TemplateID != "" {
		if tpl, err := h.templates.Get(df.TemplateID); err == nil {
			out["variants"] = tpl.Variants
		}
	}
	response.OK(c, out)
}

type putOutlineReq struct {
	Version int          `json:"version"`
	Outline deck.Outline `json:"outline"`
}

// PutOutline PUT /api/decks/:id/outline —— 面板直改（双通道之一）。
func (h *Handler) PutOutline(c *gin.Context) {
	uid, ok := h.uid(c)
	if !ok {
		return
	}
	var req putOutlineReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParameterErr(c)
		return
	}
	if err := h.decks.SaveOutline(uid, c.Param("id"), &req.Outline, req.Version); err != nil {
		mapDeckErr(c, err)
		return
	}
	response.OK(c, gin.H{"version": req.Outline.Version})
}

// ConfirmOutline POST /api/decks/:id/outline/confirm —— gate 1。
func (h *Handler) ConfirmOutline(c *gin.Context) {
	uid, ok := h.uid(c)
	if !ok {
		return
	}
	id := c.Param("id")
	to, err := h.decks.ConfirmOutline(uid, id)
	if err != nil {
		mapDeckErr(c, err)
		return
	}
	// gate 1 一过就后台预热模板推荐：等用户读到选模板页，结果多半已写进
	// deck.json 缓存，推荐区秒出——LLM 的十几秒延迟从关键路径上整个挪走。
	// client 必须在请求 ctx 里解析（BYOK 的 key 挂在 authctx），后台换用独立
	// context 活过本次响应；失败静默，预热只是增强。
	if h.suggest != nil && h.agent != nil {
		if cl, ok := h.suggestLLM(c.Request.Context()); ok {
			go h.prewarmSuggestions(cl, uid, id)
		}
	}
	response.OK(c, gin.H{"stage": to})
}

type selectTemplateReq struct {
	TemplateID string `json:"template_id"`
	Variant    string `json:"variant"`
}

// SelectTemplate POST /api/decks/:id/template —— gate 2（实例化）。
func (h *Handler) SelectTemplate(c *gin.Context) {
	uid, ok := h.uid(c)
	if !ok {
		return
	}
	var req selectTemplateReq
	if err := c.ShouldBindJSON(&req); err != nil || req.TemplateID == "" {
		response.ParameterErr(c)
		return
	}
	to, err := h.decks.SelectTemplate(uid, c.Param("id"), req.TemplateID, req.Variant)
	if err != nil {
		mapDeckErr(c, err)
		return
	}
	response.OK(c, gin.H{"stage": to})
}

// TemplateSuggestions POST /api/decks/:id/template-suggestions[?refresh=1] ——
// 选模板阶段的 AI 推荐（tplsuggest）。LLM/解析失败返回空 suggestions（200），
// 前端静默隐藏推荐区；归属走 404 口径，阶段不对走 409（mapDeckErr）。
// 澄清诉求是尽力而为的上下文：取不到（降级模式等）就空着，不阻断推荐。
func (h *Handler) TemplateSuggestions(c *gin.Context) {
	uid, ok := h.uid(c)
	if !ok {
		return
	}
	if h.suggest == nil {
		response.Err(c, http.StatusServiceUnavailable, "推荐服务不可用")
		return
	}
	cl, ok := h.suggestLLM(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusServiceUnavailable, "推荐服务不可用")
		return
	}
	deckID := c.Param("id")
	sugs, err := h.suggest.Suggest(c.Request.Context(), uid, deckID, cl,
		h.clarifyContext(c.Request.Context(), uid, deckID),
		c.Query("refresh") == "1")
	if err != nil {
		mapDeckErr(c, err)
		return
	}
	if sugs == nil {
		sugs = []deck.TplSuggestion{}
	}
	response.OK(c, gin.H{"suggestions": sugs})
}

// suggestLLM 推荐调用的 LLM 接入参数。必须在请求 ctx 里解析——BYOK 的
// 用户 key 挂在 authctx 上，换 ctx 就拿不到了。
func (h *Handler) suggestLLM(ctx context.Context) (tplsuggest.LLM, bool) {
	if h.agent == nil {
		return tplsuggest.LLM{}, false
	}
	cl := h.agent.CustomizeLLMFor(ctx)
	return tplsuggest.LLM{Client: cl.Client, Model: cl.Model}, cl.Client != nil
}

// clarifyContext 澄清诉求（尽力而为）：任何失败都当"没有诉求"，不阻断推荐。
func (h *Handler) clarifyContext(ctx context.Context, uid uint, deckID string) string {
	s, err := h.agent.ClarifyUserMessages(ctx, uid, deckID, 2000)
	if err != nil {
		return ""
	}
	return s
}

// prewarmSuggestions 后台预计算模板推荐（gate 1 确认后启动）。与用户主动请求
// 走同一个 Suggest：缓存命中就空转，同 deck 在飞去重锁保证不会重复调 LLM。
func (h *Handler) prewarmSuggestions(cl tplsuggest.LLM, uid uint, deckID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	_, _ = h.suggest.Suggest(ctx, uid, deckID, cl, h.clarifyContext(ctx, uid, deckID), false)
}

// GenerateDeck POST /api/decks/:id/generate?session_id=N[&resume=1] —— 触发生成 run。
// 响应即 SSE 流（与 /api/chat 同一套事件管线，前端复用消费逻辑）。
func (h *Handler) GenerateDeck(c *gin.Context) {
	uid, ok := h.uid(c)
	if !ok {
		return
	}
	sessionID, err := strconv.ParseUint(c.Query("session_id"), 10, 64)
	if err != nil || sessionID == 0 {
		response.Err(c, http.StatusBadRequest, "缺少或非法的 session_id（生成 run 挂在会话上）")
		return
	}
	if err := h.agent.EnsureSessionOwner(uid, uint(sessionID)); err != nil {
		response.Err(c, http.StatusNotFound, "会话不存在")
		return
	}
	resume := c.Query("resume") == "1"
	serveAgentSSE(c, func(emit func(agent.StreamEvent) error) (uint, error) {
		return h.agent.StartGenerationRun(c.Request.Context(), uid, uint(sessionID), c.Param("id"), resume, emit)
	})
}

// ExportDeck POST /api/decks/:id/export {format: pdf|png|html}
// 同步实现（8 页实测 < 30s）；产物落 exports/，GET 下载。
func (h *Handler) ExportDeck(c *gin.Context) {
	uid, ok := h.uid(c)
	if !ok {
		return
	}
	var req struct {
		Format string `json:"format"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Format == "" {
		response.ParameterErr(c)
		return
	}
	if h.exporter == nil {
		response.Err(c, http.StatusServiceUnavailable, "导出功能不可用")
		return
	}
	res, err := h.exporter.Export(c.Request.Context(), uid, c.Param("id"), req.Format)
	if err != nil {
		mapDeckErr(c, err)
		return
	}
	// 告诉前端下载会用什么名字（随文稿标题），toast 展示用；权威值在下载头里
	if r, ok := res.(export.Result); ok {
		if df, err := h.decks.GetDeckV2(uid, c.Param("id")); err == nil {
			r.DownloadAs = exportDownloadName(df.Title, r.Filename)
		}
		res = r
	}
	response.OK(c, res)
}

// exportDownloadName 把磁盘产物名映射成「文稿标题.ext」的下载名。
// 磁盘文件保持规范名（deck.html 等，白名单/缓存都依赖它），下载名只在
// Content-Disposition 里给：清洗标题里的非法文件名字符与控制符、按 rune 截断、
// 空则回退 deck。中文经 mime.FormatMediaType 走 RFC 2231 编码。
func exportDownloadName(title, diskName string) string {
	clean := strings.Map(func(r rune) rune {
		if r < 0x20 || strings.ContainsRune(`\\/:*?"<>|`, r) {
			return -1
		}
		return r
	}, strings.TrimSpace(title))
	clean = strings.Join(strings.Fields(clean), " ")
	const maxRunes = 60
	if runes := []rune(clean); len(runes) > maxRunes {
		clean = string(runes[:maxRunes])
	}
	clean = strings.Trim(clean, " .")
	if clean == "" {
		clean = "deck"
	}
	return clean + filepath.Ext(diskName)
}

// DownloadExport GET /api/decks/:id/exports/:file?token= —— 导出产物下载。
func (h *Handler) DownloadExport(c *gin.Context) {
	uid, ok := h.uid(c)
	if !ok {
		return
	}
	if err := h.decks.EnsureOwner(uid, c.Param("id")); err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	// 文件名白名单（三种产物），防路径穿越
	name := c.Param("file")
	switch name {
	case "deck.pdf", "deck-png.zip", "deck.html":
	default:
		c.Status(http.StatusNotFound)
		return
	}
	p := filepath.Join(h.decks.ExportDir(c.Param("id")), name)
	data, err := os.ReadFile(p)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	ct := "application/octet-stream"
	switch {
	case strings.HasSuffix(name, ".pdf"):
		ct = "application/pdf"
	case strings.HasSuffix(name, ".zip"):
		ct = "application/zip"
	case strings.HasSuffix(name, ".html"):
		ct = "text/html; charset=utf-8"
	}
	// 下载名跟随文稿标题（重命名后下载即刻生效）；标题读不到就退回磁盘名
	disposition := name
	if df, err := h.decks.GetDeckV2(uid, c.Param("id")); err == nil {
		disposition = exportDownloadName(df.Title, name)
	}
	c.Header("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": disposition}))
	c.Data(http.StatusOK, ct, data)
}
