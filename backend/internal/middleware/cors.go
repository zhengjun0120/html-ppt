package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS 手写最小实现：固定白名单 + 放行 OPTIONS 预检请求。
// 手写是为了看清 CORS 到底做了什么（就是几个响应头）；
// 将来规则变复杂了（通配符、携带 Cookie 等）再换 gin-contrib/cors。
func CORS(allowOrigins []string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(allowOrigins))
	for _, o := range allowOrigins {
		allowed[o] = true
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && allowed[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
