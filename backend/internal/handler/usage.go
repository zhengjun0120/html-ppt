package handler

// 用量统计的读取端（/usage 页）。写端在 service/usage（经 trace 装配的账本 sink）。

import (
	"net/http"

	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/response"

	"github.com/gin-gonic/gin"
)

// UsageOverview 当前登录用户的用量汇总（今日/本月 + 近 30 天逐日）。
func (h *Handler) UsageOverview(c *gin.Context) {
	uid, ok := authctx.UserID(c.Request.Context())
	if !ok {
		response.Err(c, http.StatusUnauthorized, "未登录")
		return
	}
	if h.usage == nil {
		// 数据库降级模式没有账本服务
		response.Err(c, http.StatusServiceUnavailable, "用量统计不可用")
		return
	}
	ov, err := h.usage.Overview(uid)
	if err != nil {
		response.Err(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, ov)
}
