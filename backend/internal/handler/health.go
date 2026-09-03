package handler

import (
	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/response"
)

type healthResponse struct {
	Status string `json:"status"`
	DB     string `json:"db"` // ok / disabled / unreachable
}

// Health GET /api/health —— 部署后给探活和前端做连接检测用。
func (h *Handler) Health(c *gin.Context) {
	resp := healthResponse{Status: "ok", DB: "disabled"}
	if h.st != nil {
		if err := h.st.Ping(c.Request.Context()); err != nil {
			resp.DB = "unreachable"
		} else {
			resp.DB = "ok"
		}
	}
	response.OK(c, resp)
}
