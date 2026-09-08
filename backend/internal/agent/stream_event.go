package agent

type StreamEvent struct {
	Type             string `json:"type"`
	Content          string `json:"content"`
	ToolCallID       string `json:"tool_call_id,omitempty"`
	ToolName         string `json:"tool_name,omitempty"`
	PromptTokens     int64  `json:"prompt_tokens,omitempty"`     //输入token
	CompletionTokens int64  `json:"completion_tokens,omitempty"` //输出token
	TotalTokens      int64  `json:"total_tokens,omitempty"`      //总token
	CachedTokens     int64  `json:"cached_tokens,omitempty"`     //缓存命中token
}

const (
	EventTypeDelta     = "delta"      // 文本增量，前端追加显示
	EventTypeToolCall  = "tool_call"  // 开始调用工具，前端可显示"正在查阅…"
	EventTypeToolError = "tool_error" // 工具调用出错
	EventTypeDone      = "done"       // 全部完成，携带完整回复
	EventTypeError     = "error"      // 出错，前端弹错误提示
	EventTypeThink     = "think"      // 思考内容，前端可显示"正在思考…"

	EventTypeSession = "session"  // 携带 session_id (首轮/续轮/恢复都会发) 由前后端保存
	EventTypeAskUser = "ask_user" // 结构化提问卡片
)
