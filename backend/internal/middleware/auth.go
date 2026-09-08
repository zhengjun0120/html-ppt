package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"html-ppt/backend/internal/authctx"
)

// Auth 校验 JWT 并把 user id 注入 context。
// token 来源（按优先级）：
//  1. Authorization: Bearer <token> —— 常规 JSON API（fetch 可自定义头）
//  2. ?token=<token> 查询参数 —— 浏览器原生导航带不了头：iframe src、window.open
//     加载 deck 文件只能走这里。
//     已知妥协：URL 里的 token 会进访问日志。生产正解是短时签名 URL 或
//     HttpOnly cookie，TODO 记在阶段5（用户系统完善时）。
func Auth(jwtSecret string) gin.HandlerFunc {
	key := []byte(jwtSecret)
	return func(c *gin.Context) {
		token := ""
		if v, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer "); ok {
			token = strings.TrimSpace(v)
		} else {
			token = c.Query("token")
		}
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			return
		}

		parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
			return key, nil
		}, jwt.WithValidMethods([]string{"HS256"})) // 显式限定算法，防 alg 混淆攻击
		if err != nil || !parsed.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "登录已过期，请重新登录"})
			return
		}

		claims, ok := parsed.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "登录态无效"})
			return
		}
		// sub 里存的是签发时的 user id（JWT 标准字段，我们签的是十进制字符串）
		id, ok := parseUintClaim(claims["sub"])
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "登录态无效"})
			return
		}

		// 同时写进 gin context（handler 直接取）和 request context
		//（agent 工具链取——工具拿到的 ctx 派生自 request）
		c.Set("user_id", id)
		c.Request = c.Request.WithContext(authctx.WithUser(c.Request.Context(), id))
		c.Next()
	}
}

func parseUintClaim(v any) (uint, bool) {
	switch n := v.(type) {
	case string:
		id, err := strconv.ParseUint(n, 10, 64) // 严格解析，拒绝 "12abc" 之类的杂音
		if err != nil {
			return 0, false
		}
		return uint(id), true
	case float64: // JSON 数字反序列化的默认类型
		if n < 0 {
			return 0, false
		}
		return uint(n), true
	default:
		return 0, false
	}
}
