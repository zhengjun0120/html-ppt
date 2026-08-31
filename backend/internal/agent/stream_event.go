package agent

type StreamEvent struct {
	Type string `json:"type"`
	Content string `json:"content"`
	ToolName string `json:"tool_name"`
}

const (
    EventTypeDelta    = "delta"     // 文本增量，前端追加显示
    EventTypeToolCall = "tool_call" // 开始调用工具，前端可显示"正在查阅…"
    EventTypeDone     = "done"      // 全部完成，携带完整回复
    EventTypeError    = "error"     // 出错，前端弹错误提示
    EventTypeThink    = "think"     // 思考内容，前端可显示"正在思考…"
)