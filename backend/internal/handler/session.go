package handler

import (
	"errors"
	"html-ppt/backend/internal/agent"
	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListDeckSessions GET /api/decks/:id/chat/sessions —— 某个 deck 名下的会话列表。
//
// 对话历史按 deck 组织：工作台进来先拿这个列表，第一条（最近活跃）当当前对话，
// 其余留给"历史会话"切换。摘要里带 pending，前端能提前知道"这个会话还压着
// 一个没答的提问"，不必等拉完消息才发现发不出新消息。
func (h *Handler) ListDeckSessions(c *gin.Context) {
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
		// 别人的 deck 和不存在的 deck 同一个响应，不泄露存在性（与 GetDeckFile 口径一致）
		response.Err(c, http.StatusNotFound, "deck 不存在")
		return
	}
	sessions, err := h.agent.DeckSessions(c.Request.Context(), uid, deckID)
	if err != nil {
		response.Err(c, sessionErrStatus(err), err.Error())
		return
	}
	response.OK(c, sessions)
}

// GetSessionMessages GET /api/chat/sessions/:id/messages?after_seq=N
//
// 返回一个会话的可回放消息列表（投影形状见 agent.Transcript）。
//
// 给"刷新后重建对话"用：流式事件的临时态活不过一次刷新，有了这个接口，
// 前端把历史灌回消息流组件就能原样恢复。after_seq 支持增量拉取
// （seq 是消息下标，健康会话只追加所以稳定；修复过的坏会话建议全量重拉）。
//
// 暂停中的会话（session.pending = true）：最后那条 ask_user 没有应答
// （修复流程补的占位已在投影时丢掉），提问卡片用 /api/chat/pending 重建，
// 回答仍走 /api/chat/answer。
func (h *Handler) GetSessionMessages(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	sessionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || sessionID == 0 {
		response.Err(c, 400, "会话 id 必须是正整数")
		return
	}
	var afterSeq int64
	if raw := c.Query("after_seq"); raw != "" {
		v, perr := strconv.ParseInt(raw, 10, 64)
		if perr != nil || v < 0 {
			response.Err(c, 400, "after_seq 必须是非负整数")
			return
		}
		afterSeq = v
	}
	if err := h.agent.EnsureSessionOwner(uid, uint(sessionID)); err != nil {
		response.Err(c, http.StatusNotFound, "会话不存在")
		return
	}
	t, err := h.agent.SessionTranscript(c.Request.Context(), uid, uint(sessionID), afterSeq)
	if err != nil {
		response.Err(c, sessionErrStatus(err), err.Error())
		return
	}
	response.OK(c, t)
}

// sessionErrStatus 会话读取错误 → HTTP 状态：
// 降级模式是"功能关了"（503），不该让前端当故障重试；其余按服务故障（500）。
// 归属/存在性问题在上面的 EnsureOwner / EnsureSessionOwner 预检里就变成 404 了。
func sessionErrStatus(err error) int {
	if errors.Is(err, agent.ErrSessionStoreUnavailable) {
		return http.StatusServiceUnavailable
	}
	return http.StatusInternalServerError
}
