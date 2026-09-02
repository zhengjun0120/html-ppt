// agentdemo（完整版）：流式输出 + 累加器 + 多轮上下文 + 工具调用 + 思考/回答分离。
//
// 演示的核心数据流（对照 Run 和 step 的注释阅读）：
//
//	流式分片 ──AddChunk──▶ 累加器(完整消息) ──ToParam()──▶ messages 历史 ──▶ 下一轮请求
//
// 运行前环境变量（在 backend 目录 go run ./cmd/agentdemo）：
//
//	LLM_API_KEY（或 OPENAI_API_KEY）—— 模型服务密钥
//	LLM_BASE_URL（可选）—— 第三方兼容端点，如 https://api.deepseek.com/v1
//	LLM_MODEL（可选，默认 deepseek-chat；换成 deepseek-reasoner 可看到"思考"流）
//	LLM_BUDGET（可选，跨轮 token 总预算，默认 200000）
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// =====================================================================
// 第 1 部分：工具——schema 由参数 struct 反射生成，执行函数与 struct 同源
// =====================================================================

// ToolFunc 执行一个工具：入参是模型生成的 arguments（JSON 字符串），
// 返回值会作为 role:"tool" 消息的内容发回给模型。
type ToolFunc func(ctx context.Context, arguments string) (string, error)

type Tool struct {
	Definition openai.ChatCompletionToolUnionParam // 发给模型的 schema
	Execute    ToolFunc                            // 模型点名时真正执行的函数
}

// ---- 参数 struct：tag 即 schema ----
// json:"..."                   → 属性名
// jsonschema_description:"..." → 属性描述（模型填参的依据，每个字段都要写）
// jsonschema:"required"        → 必填（配合 RequiredFromJSONSchemaTags）
// jsonschema:"enum=a,enum=b"   → 枚举等约束

type WeatherArgs struct {
	City string `json:"city" jsonschema:"required" jsonschema_description:"城市名，如：北京"`
}

type NowArgs struct {
	TZ string `json:"tz" jsonschema_description:"IANA 时区名，如 Asia/Shanghai；省略则用服务器本地时间"` // 无 required → 可选
}

type Slide struct {
	Title   string   `json:"title" jsonschema:"required" jsonschema_description:"本页标题"`
	Bullets []string `json:"bullets" jsonschema:"required" jsonschema_description:"本页要点，每条一句话"`
}

type CreatePresentationArgs struct {
	Title  string  `json:"title" jsonschema:"required,minLength=1" jsonschema_description:"演示文稿主标题"`
	Theme  string  `json:"theme" jsonschema:"required,enum=dark,enum=light,enum=gradient" jsonschema_description:"视觉主题，用户未指定时选 dark"`
	Slides []Slide `json:"slides" jsonschema:"required" jsonschema_description:"幻灯片列表，按播放顺序"`
}

// generateSchema 把 Go struct 反射成 OpenAI 工具参数 schema（启动时调用一次并缓存）。
func generateSchema[T any]() openai.FunctionParameters {
	reflector := jsonschema.Reflector{
		RequiredFromJSONSchemaTags: true,  // 必填由 jsonschema:"required" 标签决定
		AllowAdditionalProperties:  false, // 生成 additionalProperties:false，拒绝幻觉字段
		DoNotReference:             true,  // 嵌套 struct 内联展开，不用 $defs/$ref
	}
	s := reflector.Reflect(new(T))

	// 只取三件套重建顶层，丢掉 $schema/$id 等模型不关心的键（官方 README 同款做法）
	schema := map[string]any{
		"type":                 "object",
		"properties":           s.Properties, // 有序 map：字段顺序 = struct 声明顺序
		"additionalProperties": false,        // 开 strict 时的硬性要求
	}
	if len(s.Required) > 0 {
		schema["required"] = s.Required
	}
	raw, err := json.Marshal(schema)
	if err != nil {
		panic(fmt.Sprintf("生成工具 schema 失败: %v", err))
	}
	var fp openai.FunctionParameters
	if err := json.Unmarshal(raw, &fp); err != nil {
		panic(err)
	}
	return fp
}

func toolWeather(_ context.Context, arguments string) (string, error) {
	var args WeatherArgs
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("参数不是合法 JSON: %w", err)
	}
	return fmt.Sprintf(`{"city":%q,"weather":"晴","temp_c":26}`, args.City), nil
}

func toolNow(_ context.Context, arguments string) (string, error) {
	if arguments == "" { // 无参/可选参数工具的 arguments 可能是空串
		arguments = "{}"
	}
	var args NowArgs
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", err
	}
	return fmt.Sprintf(`{"now":%q}`, time.Now().Format(time.RFC3339)), nil
}

// 有副作用的工具：真的写出 HTML 文件，运行后在 backend 目录可见。
func toolCreatePresentation(_ context.Context, arguments string) (string, error) {
	var args CreatePresentationArgs
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("参数不是合法 JSON: %w", err)
	}
	var b strings.Builder
	b.WriteString("<!doctype html><html><head><meta charset=\"utf-8\"><title>")
	b.WriteString(html.EscapeString(args.Title))
	b.WriteString("</title></head><body>")
	for _, s := range args.Slides {
		b.WriteString("<section><h2>" + html.EscapeString(s.Title) + "</h2><ul>")
		for _, item := range s.Bullets {
			b.WriteString("<li>" + html.EscapeString(item) + "</li>")
		}
		b.WriteString("</ul></section>")
	}
	b.WriteString("</body></html>")

	const path = "demo_presentation.html"
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		return "", err
	}
	return fmt.Sprintf(`{"path":%q,"slides":%d}`, path, len(args.Slides)), nil
}

func buildTools() map[string]Tool {
	return map[string]Tool{
		"get_weather": {
			Definition: openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
				Name:        "get_weather",
				Description: openai.String("查询中国某个城市的实时天气"),
				Parameters:  generateSchema[WeatherArgs](),
			}),
			Execute: toolWeather,
		},
		"get_now": {
			Definition: openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
				Name:        "get_now",
				Description: openai.String("获取服务器当前时间"),
				Parameters:  generateSchema[NowArgs](),
			}),
			Execute: toolNow,
		},
		"create_presentation": {
			Definition: openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
				Name:        "create_presentation",
				Description: openai.String("根据主题和逐页大纲生成一份 HTML 演示文稿并保存为文件。当用户要求制作 PPT/幻灯片/演示文稿时使用"),
				Parameters:  generateSchema[CreatePresentationArgs](),
				// Strict: openai.Bool(true), // OpenAI 端可开；部分第三方厂商不认，报错就去掉
			}),
			Execute: toolCreatePresentation,
		},
	}
}

// =====================================================================
// 第 2 部分：流式分类——每个分片必属于以下几类之一
//
//	思考   delta.reasoning_content（厂商私有字段，从 JSON.ExtraFields 挖）
//	回答   delta.Content
//	工具   delta.ToolCalls（碎片，累加器按 index 归并）
//	结束   choice.FinishReason（收尾分片才有值）
//	用量   choices 为空的收尾分片（需 IncludeUsage）
// =====================================================================

// StreamCallbacks 把增量实时分发给界面。
// OnDelta / OnReasoning 的分离就是前端"思考折叠区 + 打字机正文"的后端来源。
type StreamCallbacks struct {
	OnDelta     func(string) // 回答增量
	OnReasoning func(string) // 思考增量（OpenAI 官方接口没有思考流，仅部分厂商提供）
}

// step 发一次流式请求：增量走回调，完整结果走累加器。
func (a *Agent) step(ctx context.Context, tools []openai.ChatCompletionToolUnionParam, cb StreamCallbacks) (*openai.ChatCompletionAccumulator, error) {
	params := openai.ChatCompletionNewParams{
		Model:    openai.ChatModel(a.model),
		Messages: a.messages,
		StreamOptions: openai.ChatCompletionStreamOptionsParam{
			IncludeUsage: openai.Bool(true), // 流式下默认拿不到 usage
		},
	}
	if len(tools) > 0 { // 不带工具 = 普通流式聊天；强制收尾轮也靠传 nil 禁用工具
		params.Tools = tools
	}

	stream := a.client.Chat.Completions.NewStreaming(ctx, params)
	defer stream.Close() // 中途 return 时必须释放连接（正常走完会被自动关闭，重复调无害）

	acc := openai.ChatCompletionAccumulator{}
	for stream.Next() {
		chunk := stream.Current()
		if !acc.AddChunk(chunk) {
			return nil, errors.New("流式分片累加失败（分片 ID 不一致或超出协议上限）")
		}

		// 纯 usage 收尾分片：choices 为空（IncludeUsage 开启时的最后一个分片）
		if len(chunk.Choices) == 0 {
			continue
		}
		d := chunk.Choices[0].Delta

		// 思考增量：reasoning_content 是 DeepSeek/Qwen 等厂商的私有扩展，
		// 官方 SDK 没这个字段，落到 JSON.ExtraFields；Raw() 是原始 JSON 片段（带引号），要再解一层。
		// 注意：累加器不收集 ExtraFields，所以思考文本必须在这里逐片取走。
		if f, ok := d.JSON.ExtraFields["reasoning_content"]; ok {
			var reasoning string
			if json.Unmarshal([]byte(f.Raw()), &reasoning) == nil && reasoning != "" && cb.OnReasoning != nil {
				cb.OnReasoning(reasoning)
			}
		}
		// 回答增量
		if d.Content != "" && cb.OnDelta != nil {
			cb.OnDelta(d.Content)
		}

		// 工具调用碎片已由 AddChunk 按 index 归并，无需手工拼接。
		// [事件] 行仅调试用——官方警告：并行工具调用时 JustFinishedToolCall 不可依赖，
		// 正式逻辑以流结束后的 acc.Choices[0].Message.ToolCalls 为准。
		if tool, ok := acc.JustFinishedToolCall(); ok {
			fmt.Fprintf(os.Stderr, "[事件] 工具调用接收完毕: %s(%s)\n", tool.Name, tool.Arguments)
		}
	}
	if err := stream.Err(); err != nil {
		return nil, err // 网络/HTTP/模型端中途错误（*ssestream.StreamError 也在其中）
	}
	return &acc, nil
}

// =====================================================================
// 第 3 部分：多轮循环——把"流式分片"组装成"下一轮的上下文"
//
//	历史只追加不改写；每轮新增：
//	  user(本轮输入) → assistant(tool_calls) → tool(结果)×N → … → assistant(最终回答)
// =====================================================================

const (
	maxTurns  = 25 // 轮数护栏
	repeatCut = 3  // 连续相同调用的中止阈值
)

type Agent struct {
	client      *openai.Client
	model       string
	tools       []openai.ChatCompletionToolUnionParam
	exec        map[string]ToolFunc
	messages    []openai.ChatCompletionMessageParamUnion // 会话上下文：跨 Run 持续追加
	usedTokens  int64                                    // 跨轮累计用量（token 预算护栏）
	tokenBudget int64
	repeats     map[string]int // 重复调用检测（会话级）
}

func NewAgent(client *openai.Client, model string, tokenBudget int64, tools map[string]Tool) *Agent {
	a := &Agent{
		client: client, model: model, tokenBudget: tokenBudget,
		exec: map[string]ToolFunc{}, repeats: map[string]int{},
	}
	for name, t := range tools {
		a.tools = append(a.tools, t.Definition)
		a.exec[name] = t.Execute
	}
	return a
}

// Run 跑一轮完整对话（多次调用 Run 共享同一份 messages，即多轮会话）。
func (a *Agent) Run(ctx context.Context, userMsg string, cb StreamCallbacks) (string, error) {
	a.messages = append(a.messages, openai.UserMessage(userMsg))

	for turn := 1; ; turn++ {
		tools := a.tools
		// 双护栏：轮数或累计 token 任一耗尽 → 注入收尾指令并禁用工具，强制总结（优雅退出）
		if turn > maxTurns || a.usedTokens >= a.tokenBudget {
			fmt.Fprintf(os.Stderr, "[护栏] 触发收尾（轮数=%d, 累计tokens=%d）\n", turn-1, a.usedTokens)
			a.messages = append(a.messages, openai.UserMessage(
				"调用预算已用完：不要再调用任何工具，直接基于以上已获取的信息给出最终回答。"))
			tools = nil
			cb.OnDelta("\n[预算/轮数已达上限，强制收尾]\n")

			acc, err := a.step(ctx, tools, cb)
			if err != nil {
				return "", err
			}
			if len(acc.Choices) == 0 {
				return "", errors.New("收尾轮响应不完整")
			}
			return acc.Choices[0].Message.Content, nil
		}

		acc, err := a.step(ctx, tools, cb)
		if err != nil {
			return "", fmt.Errorf("第 %d 轮失败: %w", turn, err)
		}

		// ★ 健壮性防线：不完整的流绝不能进上下文。
		// FinishReason 为空 = 没收到收尾分片 = assistant 消息残缺，
		// 其 tool_calls 参数可能拼了一半，进历史后下一轮必然配对失败。
		if len(acc.Choices) == 0 || acc.Choices[0].FinishReason == "" {
			return "", errors.New("流式响应不完整（未收到 finish_reason），本轮不进入上下文")
		}
		choice := acc.Choices[0]
		a.usedTokens += acc.Usage.TotalTokens // 跨轮记账（token 预算的依据）
		fmt.Fprintf(os.Stderr, "[第 %d 轮] finish=%s tokens=%d(累计%d)\n",
			turn, choice.FinishReason, acc.Usage.TotalTokens, a.usedTokens)

		msg := choice.Message
		// 没有工具调用 → 最终回答，进历史后结束循环
		if len(msg.ToolCalls) == 0 {
			a.messages = append(a.messages, msg.ToParam())
			if msg.Content == "" {
				return "", errors.New("模型返回空内容（常见原因：输出上限被推理 token 吃掉）")
			}
			return msg.Content, nil
		}

		// 关键规则 1：assistant 消息（含拼好的 tool_calls）先原样进历史
		a.messages = append(a.messages, msg.ToParam())

		// 关键规则 2：逐个执行，role:"tool" + tool_call_id 一一对应回填
		for _, tc := range msg.ToolCalls {
			// 死循环检测：连续相同 name+arguments 达到阈值 → 中止
			key := tc.Function.Name + ":" + tc.Function.Arguments
			a.repeats[key]++
			if a.repeats[key] >= repeatCut {
				return "", fmt.Errorf("模型连续 %d 次重复调用 %s，疑似死循环，已中止", repeatCut, tc.Function.Name)
			}

			result, err := a.callTool(ctx, tc.Function.Name, tc.Function.Arguments)
			if err != nil {
				// 工具失败不中断流程：把错误作为结果回填，让模型自己调整
				result = fmt.Sprintf(`{"error":%q}`, err.Error())
			}
			a.messages = append(a.messages, openai.ToolMessage(result, tc.ID))
		}
		// 继续下一轮：新上下文 = 全部历史（累加器拼好的消息已在其中）
	}
}

func (a *Agent) callTool(ctx context.Context, name, arguments string) (string, error) {
	fn, ok := a.exec[name]
	if !ok {
		return "", fmt.Errorf("未知工具 %q", name) // 模型幻觉出的工具名：作为错误结果回填
	}
	tctx, cancel := context.WithTimeout(ctx, 30*time.Second) // 单工具超时
	defer cancel()
	out, err := fn(tctx, arguments)
	if err != nil {
		return "", err
	}
	if len(out) > 8192 { // 工具结果截断，防止撑爆上下文
		out = out[:8192] + "...[结果已截断]"
	}
	return out, nil
}

// =====================================================================
// 第 4 部分：main——三个场景
// =====================================================================

func main() {
	apiKey := os.Getenv("LLM_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKey == "" {
		fmt.Println("请先设置环境变量 LLM_API_KEY（或 OPENAI_API_KEY）")
		os.Exit(1)
	}

	opts := []option.RequestOption{option.WithAPIKey(apiKey)}
	if base := os.Getenv("LLM_BASE_URL"); base != "" {
		opts = append(opts, option.WithBaseURL(base))
	}
	model := os.Getenv("LLM_MODEL")
	if model == "" {
		model = "deepseek-chat"
	}
	budget := int64(200000)
	if v := os.Getenv("LLM_BUDGET"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 {
			budget = n
		}
	}

	client := openai.NewClient(opts...) // 值类型
	agent := NewAgent(&client, model, budget, buildTools())

	// 回调：思考流灰色输出（需终端支持 ANSI，不支持就打到日志），
	// 回答流原样输出——前端把这两个回调换成两条 SSE 事件即是"思考折叠区"效果。
	cb := StreamCallbacks{
		OnDelta:     func(s string) { fmt.Print(s) },
		OnReasoning: func(s string) { fmt.Print("\033[90m" + s + "\033[0m") },
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	fmt.Println("=== 场景 1：普通流式对话（LLM_MODEL=deepseek-reasoner 时可见灰色思考流） ===")
	answer, err := agent.Run(ctx, "用两句话介绍一下你自己", cb)
	if err != nil {
		fmt.Println("出错:", err)
	} else {
		fmt.Printf("\n[最终回答] %s\n\n", answer)
	}

	fmt.Println("=== 场景 2：并行工具调用（注意 stderr 的 [事件] 与 [第 N 轮] 行） ===")
	answer, err = agent.Run(ctx, "北京和上海今天天气怎么样？顺便现在几点了？", cb)
	if err != nil {
		fmt.Println("出错:", err)
	} else {
		fmt.Printf("\n[最终回答] %s\n\n", answer)
	}

	fmt.Println("=== 场景 3：生成 PPT（嵌套 schema + 副作用工具） ===")
	answer, err = agent.Run(ctx, "帮我做一份关于 Go 并发编程的 3 页 PPT，主题用 dark", cb)
	if err != nil {
		fmt.Println("出错:", err)
	} else {
		fmt.Printf("\n[最终回答] %s\n", answer)
	}
	// 场景 3 跑完看 backend/demo_presentation.html
}
