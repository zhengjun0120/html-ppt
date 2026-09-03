package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 本包约定（写给所有调用者看的契约）：
// 1. 错误信息靠 HTTP 状态码表达，body 一律是 {"error": "..."}；
//    成功一律 200，body 直接是数据本身（REST 风格，不用 {code, data} 业务码包装）
// 2. 返回给客户端的错误信息是"给人看的"（前端会展示、或转述给用户）；
//    堆栈、SQL、内部路径等细节打日志，不要塞进响应
// 3. 本包只服务 JSON API；文件流（如 deck.html 预览）和 SSE 走原生写法，不经过这里

// OK 统一成功响应：200 + JSON body。
// data 为 nil 时输出 {}，保证前端拿到的永远是对象而不是 null。
func OK(c *gin.Context, data any) {
	if data == nil {
		data = gin.H{}
	}
	c.JSON(http.StatusOK, data)
}

// Err 统一错误响应：HTTP 状态码 + {"error": msg}。
func Err(c *gin.Context, httpStatus int, msg string) {
	c.JSON(httpStatus, gin.H{"error": msg})
}

// Errf 带格式化的错误响应，省去调用方 fmt.Sprintf。
func Errf(c *gin.Context, httpStatus int, format string, args ...any) {
	Err(c, httpStatus, fmt.Sprintf(format, args...))
}

// ParameterErr 参数错误的响应
func ParameterErr(c *gin.Context){
	c.JSON(400,gin.H{"error":"参数错误"})
}
