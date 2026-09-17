package handler

// deck-v2 的管线端点：大纲（读/改/确认）、模板选择、生成触发、deck 元数据。
//
// 阶段守卫在 service 层（transitionStage / SaveOutline / SelectTemplate），
// handler 只做参数绑定与错误映射：阶段冲突 → 409，版本冲突 → 409 带最新版，
// 归属问题 → 404（与 v1 同一条"不泄露存在性"纪律）。

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/agent"
	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/service/deck"
	"html-ppt/backend/internal/response"
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
	to, err := h.decks.ConfirmOutline(uid, c.Param("id"))
	if err != nil {
		mapDeckErr(c, err)
		return
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
