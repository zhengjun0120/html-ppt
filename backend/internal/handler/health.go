package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"html-ppt/backend/internal/store"
)

type healthResponse struct {
	Status string `json:"status"`
	DB     string `json:"db"` // ok / disabled / unreachable
}

// Health GET /api/health —— 部署后给探活和前端做连接检测用。
func Health(st *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := healthResponse{Status: "ok", DB: "disabled"}
		if st != nil {
			if err := st.Ping(c.Request.Context()); err != nil {
				resp.DB = "unreachable"
			} else {
				resp.DB = "ok"
			}
		}
		c.JSON(http.StatusOK, resp)
	}
}
