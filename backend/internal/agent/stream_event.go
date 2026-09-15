package agent

import "html-ppt/backend/internal/trace"

type StreamEvent struct {
	Type       string `json:"type"`
	Content    string `json:"content"`
	ToolCallID string `json:"tool_call_id,omitempty"`
	ToolName   string `json:"tool_name,omitempty"`
	// 下面四个扁平字段是**主循环口径**（不端子调用的历史语义，保持不动以免打破前端契约）。
	// 视觉审查与联网搜索各自的用量在 Usage 里分项列出——真实总成本看 Usage.Total。
	PromptTokens     int64 `json:"prompt_tokens,omitempty"`     //输入token
	CompletionTokens int64 `json:"completion_tokens,omitempty"` //输出token
	TotalTokens      int64 `json:"total_tokens,omitempty"`      //总token
	CachedTokens     int64 `json:"cached_tokens,omitempty"`     //缓存命中token
	// Usage 分项用量汇总，只在 done 事件上带。为什么必须分项：视觉审查（每次要传好几张图）
	// 和联网搜索（按搜索次数计费）过去完全没被计入，所以"这次对话花了多少 token"
	// 那时的答案是系统性偏小的。
	Usage *trace.Summary `json:"usage,omitempty"`
}

const (
	EventTypeDelta     = "delta"      // 文本增量，前端追加显示
	EventTypeToolCall  = "tool_call"  // 工具调用成功，前端可显示"正在查阅…"
	EventTypeToolError = "tool_error" // 工具调用出错
	EventTypeDone      = "done"       // 全部完成，携带完整回复
	EventTypeError     = "error"      // 出错，前端弹错误提示
	EventTypeThink     = "think"      // 思考内容，前端可显示"正在思考…"

	EventTypeSession = "session"  // 携带 session_id (首轮/续轮/恢复都会发) 由前后端保存
	EventTypeAskUser = "ask_user" // 结构化提问卡片
	// EventTypeTrace 观测事件（工具参数、返回耗时、工具内部过程、分项用量）。
	// Content 是一条 trace.Event 的 JSON，**不含 messages**（那个太大，走观测页按需拉）。
	// 不认识这个类型的前端会把它归到 default 分支——所以客户端要么显式处理，
	// 要么就该默认忽略未知类型，别刷"未知事件类型"。
	EventTypeTrace = "trace"
)
