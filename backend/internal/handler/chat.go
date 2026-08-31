package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Chat POST /api/chat —— agent 循环的入口，阶段1实现。
//
// 计划中的协议（SSE 流式返回，事件类型）：
//   agent_text    LLM 的流式文字回复
//   tool_progress 工具执行进度（如 "正在读取第 3 页…"）
//   ask_user      结构化提问卡片（问题 + 推荐答案 + 选项），循环暂停等用户作答
//   done          本轮结束
func Chat() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{
			"error": "chat 未实现（阶段1：最小 agent 循环）",
			"plan":  "消息 → DeepSeek(带工具定义) → 解析 tool_calls → Go 执行工具 → 结果回填 → 循环；SSE 流式输出",
		})
	}
}
