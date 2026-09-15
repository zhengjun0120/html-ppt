// Package trace 是 agent 的观测通道：把"agent 到底做了什么"落成一条条事件。
//
// 它刻意与 SSE 那条通道分开。SSE 是**给用户看的对话流**（增量文字、思考、提问卡片），
// 面向的是"正在生成的这份 deck"；trace 是**给人排查用的执行流**（每轮请求的完整上下文、
// 工具参数、工具内部过程、分项 token），面向的是"这次为什么这么做"。
// 两个受众要的东西不一样，硬塞进一条流会让聊天流被几 MB 的上下文拖住。
//
// 三个设计约束，都是踩过才会显形的地方：
//
//  1. **fail-open**：观测失败绝不能中断对话。所有写盘/回调错误只打 warn——与
//     cmd/server/main.go 对视觉审查的处理同一条原则（它防的是"增强功能反而变成故障源"）。
//  2. **事件只追加**（JSONL，一行一个）：实时查看靠的就是"写到哪看到哪"，
//     整文件覆盖写会让运行中的 run 在观测页上一直是空的。
//  3. **归属信息走 context**：工具不需要知道"我是第几轮、我的 tool_call_id 是什么"，
//     那些由循环在调用前注入（见 recorder.go 的 WithTurn / WithTool）。工具只管
//     trace.Emit(ctx, ...)，拿不到 recorder 时自动变成空操作——测试和其他调用方零改动。
package trace

import (
	"encoding/json"
	"time"
)

// 事件种类。一次 run 的生命周期：
//
//	run_start → [llm_request → llm_response → (tool_call → tool_result|sub_step*)+ → usage*] × N → run_end
const (
	KindRunStart    = "run_start"    // 一轮对话开始，带用户消息与本次挂载的工具
	KindLLMRequest  = "llm_request"  // 发给模型的东西（含全量上下文，只落盘）
	KindLLMResponse = "llm_response" // 模型回的东西（文本 / 工具调用意图 / 本轮用量）
	KindToolCall    = "tool_call"    // 工具即将执行：名字 + **参数原文**
	KindToolResult  = "tool_result"  // 工具执行完毕：返回值 + 耗时；失败时带 error
	KindSubStep     = "sub_step"     // 工具内部过程：联网搜了什么、视觉审查拍了哪些图
	KindUsage       = "usage"        // 某个计费来源的用量（主循环/视觉/联网各算各的）
	KindError       = "error"        // 中途出错：流断了、写库失败。原因进 trace，别只留在服务端日志里
	KindRunEnd      = "run_end"      // 本轮结束（正常结束 / ask_user 暂停 / 出错）
)

// 用量归属。这三项是**分开计费的三次 API 调用**，合并成一个 total 会看不出成本花在哪：
// 联网搜索按搜索次数计费、视觉审查每次要传好几张图（图片 token 很贵），
// 而它们过去都被漏掉了——done 事件里的 prompt/completion 一直只是主循环的数。
const (
	CompMain      = "main"       // 主 agent 循环
	CompVision    = "vision"     // 视觉审查子调用
	CompWebSearch = "web_search" // 联网搜索子调用
)

// run 的收尾状态。
const (
	StatusOK     = "ok"
	StatusPaused = "paused" // ask_user 让循环停在等回答处
	StatusError  = "error"
	// StatusRunning 不是 run_end 里写的值，是**读侧推断**出来的：
	// 文件里没有 run_end（还在跑，或者进程被杀了）。实时查看时这个状态最重要。
	StatusRunning = "running"
)

// UsagePart 一次或多次调用的用量。Cached 单独列是因为它决定实际花费
// （前缀命中部分便宜一个量级），而它跟 Prompt 不是一回事。
type UsagePart struct {
	Prompt      int64 `json:"prompt"`
	Completion  int64 `json:"completion"`
	Total       int64 `json:"total"`
	Cached      int64 `json:"cached"`
	Reasoning   int64 `json:"reasoning,omitempty"`    // 推理 token：它先于 content 吃预算，用完就没有正文
	ImageTokens int64 `json:"image_tokens,omitempty"` // 图片输入 token：视觉审查的成本主要就在这一项
	Calls       int   `json:"calls"`                  // 产生这份用量的 API 调用次数
}

func (u UsagePart) Add(o UsagePart) UsagePart {
	return UsagePart{
		Prompt:      u.Prompt + o.Prompt,
		Completion:  u.Completion + o.Completion,
		Total:       u.Total + o.Total,
		Cached:      u.Cached + o.Cached,
		Reasoning:   u.Reasoning + o.Reasoning,
		ImageTokens: u.ImageTokens + o.ImageTokens,
		Calls:       u.Calls + o.Calls,
	}
}

// Event 一条观测记录。字段按种类分组织，用 omitempty 让 JSONL 保持可读
// （一行几 MB 已经够难受了，不想再夹一堆空字段）。
type Event struct {
	Seq  int       `json:"seq"` // run 内单调递增，从 1 开始；增量拉取与去重都靠它
	TS   time.Time `json:"ts"`
	Kind string    `json:"kind"`

	// —— run_start ——
	RunID       string   `json:"run_id,omitempty"`
	ParentRunID string   `json:"parent_run_id,omitempty"` // ask_user 恢复时指回被暂停的那个 run
	SessionID   uint     `json:"session_id,omitempty"`
	UserID      uint     `json:"user_id,omitempty"` // 读侧归属校验的唯一依据
	DeckID      string   `json:"deck_id,omitempty"`
	UserContent string   `json:"user_content,omitempty"`
	Model       string   `json:"model,omitempty"`
	Tools       []string `json:"tools,omitempty"` // 本次挂载了哪些工具（features 开关的效果一眼可见）

	// —— agent 循环归属（由 ctx 自动补齐，见 recorder.go）——
	// Turn 是指针是因为"第 0 轮"（第一次请求）与"不适用"（run_start / run_end）
	// 必须区分得开：用 int + omitempty 的话第一次请求的轮次会被静默丢掉，
	// 页面上只有第二轮回合起才显示"第 N 轮"——看起来像第一轮没记录。
	Turn       *int   `json:"turn,omitempty"`
	ToolName   string `json:"tool_name,omitempty"`
	ToolCallID string `json:"tool_call_id,omitempty"`

	// —— 工具调用 / 返回 ——
	Args       string `json:"args,omitempty"`   // 模型给的参数原文（不解析：解析失败的样子本身就是信息）
	Result     string `json:"result,omitempty"` // 工具返回的原文
	DurationMS int64  `json:"duration_ms,omitempty"`
	Error      string `json:"error,omitempty"` // 失败调用的原因；成功为空

	// —— LLM ——
	FinishReason string          `json:"finish_reason,omitempty"`
	Messages     json.RawMessage `json:"messages,omitempty"` // 全量上下文；默认不推实时、按需再拉
	MessageCount int             `json:"message_count,omitempty"`
	Bytes        int             `json:"bytes,omitempty"` // Messages 的字节数：实时流里就靠它表达"这次上下文多大"
	Content      string          `json:"content,omitempty"`
	ToolCalls    []ToolCallOut   `json:"tool_calls,omitempty"`

	// —— 用量 ——
	Component string     `json:"component,omitempty"`
	Usage     *UsagePart `json:"usage,omitempty"`

	// —— 子过程 ——
	Sub    *SubStep     `json:"sub,omitempty"`
	Images []ImageEvent `json:"images,omitempty"`

	// —— run_end ——
	Status  string   `json:"status,omitempty"`
	Summary *Summary `json:"summary,omitempty"`
}

// ToolCallOut 模型发起的一次工具调用意图。Arguments 保留原文而不是解析成对象：
// 模型偶尔会给出非法 JSON（那正是要观测的失败），解析失败也得能原样看见。
type ToolCallOut struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// SubStep 工具内部的一步。这是"工具返回"之外真正能看到内部过程的地方：
// 联网搜索的子模型实际发了什么 query、视觉审查到底拍了哪几张图。
type SubStep struct {
	Name  string `json:"name"`           // web_search / vision
	Stage string `json:"stage"`          // request / queries / response / capture / review_prompt / review_response
	Text  string `json:"text,omitempty"` // 原始文本（提示词、报告正文）
	Data  any    `json:"data,omitempty"` // 结构化内容（来源列表、逐页量测数字）
}

// ImageEvent 一张要落盘的图。
//
// Bytes 刻意是 json:"-"：图片有几百 KB，内联进 JSONL 会让文件膨胀到不可读、
// 也让"逐行拉取"这件事失去意义。Emit 时写到 <run>/img/<name>，再把相对路径填进 URL。
type ImageEvent struct {
	Name  string `json:"name"` // 文件名，如 p003.png（读侧按白名单 ^p\d{3}\.png$ 校验）
	Label string `json:"label,omitempty"`
	URL   string `json:"url,omitempty"` // run 目录相对路径，如 <run_id>/img/p003.png
	Bytes []byte `json:"-"`
}

// Summary 一次 run 的收尾汇总。分项用量在这里落地，所以观测页能回答
// "这次对话的 token 到底花在哪儿"——而不只是"一共花了多少"。
type Summary struct {
	DurationMS int64                `json:"duration_ms"`
	Turns      int                  `json:"turns"`      // 主循环请求 LLM 的轮数
	ToolCalls  int                  `json:"tool_calls"` // 工具调用次数（成功+失败）
	Usage      map[string]UsagePart `json:"usage"`      // 按 CompMain/CompVision/CompWebSearch 分项
	Total      UsagePart            `json:"total"`      // 三项之和
}

// RunMeta 列表页需要的一条 run 摘要。只从 run_start（文件首行）和一个可选的
// run_end（文件尾行）推出来——列表要能扫几十上百个 run，不能把每个文件都全解析一遍
// （一个开了全量上下文的 run 可能有几十 MB）。
type RunMeta struct {
	RunID       string    `json:"run_id"`
	ParentRunID string    `json:"parent_run_id,omitempty"`
	SessionID   uint      `json:"session_id"`
	UserID      uint      `json:"-"` // 只用于按属主过滤，不下发（列表里全是自己的）
	DeckID      string    `json:"deck_id,omitempty"`
	UserContent string    `json:"user_content,omitempty"`
	Model       string    `json:"model,omitempty"`
	StartedAt   time.Time `json:"started_at"`
	DurationMS  int64     `json:"duration_ms"`
	Status      string    `json:"status"`
	Turns       int       `json:"turns,omitempty"`
	ToolCalls   int       `json:"tool_calls,omitempty"`
	Usage       *Summary  `json:"usage,omitempty"` // 只有已结束的 run 才有
	Bytes       int64     `json:"bytes,omitempty"` // jsonl 文件大小，用来判断"这个 run 很重"
}

// SessionMeta 一个会话（= 用户眼里的"一次对话"）的汇总。
// 一次对话可能因 ask_user 暂停而分成多个 run，所以"这次对话花了多少 token"
// 必须跨 run 累加——单个 run 的 total 只是半场。
type SessionMeta struct {
	SessionID   uint      `json:"session_id"`
	DeckID      string    `json:"deck_id,omitempty"`
	Runs        int       `json:"runs"`
	LastAt      time.Time `json:"last_at"`
	Latest      string    `json:"latest_run_id"` // 最近一个 run，页面默认展开它
	Usage       Summary   `json:"usage"`         // 跨 run 累加
	RunningRuns int       `json:"running_runs"`  // 还在跑的 run 数（页面据此决定要不要继续轮询）
}

// Truncate 按字节上限截断字符串。0 表示不截断（默认：你选了"全量存"）。
// 抽成函数是为了让"截断了"这件事在返回值里可见——静默截断会让人把半截 JSON
// 当成完整的去读，和 web_search 里那个"截断了要标 Truncated"是同一个理由。
func Truncate(s string, max int) (string, bool) {
	if max <= 0 || len(s) <= max {
		return s, false
	}
	return s[:max] + "\n…(已按 max_field_bytes 截断)", true
}
