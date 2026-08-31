// agentdemo：流式输出 + 工具调用 + agent 循环 的最小可运行教学示例。
//
// 运行前设置环境变量：
//
//	LLM_API_KEY（或 OPENAI_API_KEY）—— 模型服务密钥
//	LLM_BASE_URL（可选）—— 第三方兼容端点，如 https://api.deepseek.com/v1
//	LLM_MODEL（可选）—— 模型名，默认 deepseek-chat
//
// 在 backend 目录执行：go run ./cmd/agentdemo
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"os"
	"strings"
	"time"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// =====================================================================
// 第 1 部分：工具的"定义"与"执行"分离
// 模型只能看到 Definition（JSON Schema），永远看不到 Execute。
// schema 不再手写 map：用参数 struct + invopop/jsonschema 反射生成。
// 同一个 struct 既是 schema 的来源，又是执行函数 json.Unmarshal 的目标——
// 单一事实源，字段改名/加校验两处同步，编译期就能发现拼写错误。
// =====================================================================

// ToolFunc 执行一个工具：入参是模型生成的 arguments（JSON 字符串），
// 返回值会作为 role:"tool" 消息的内容发回给模型。
type ToolFunc func(ctx context.Context, arguments string) (string, error)

type Tool struct {
	Definition openai.ChatCompletionToolUnionParam // 发给模型的 schema
	Execute    ToolFunc                            // 模型点名时真正执行的函数
}

// ---- 参数 struct：tag 即 schema ----
//
// json:"..."            → 属性名（顺带是执行时 Unmarshal 的键名）
// jsonschema_description:"..." → 属性描述，模型填参的唯一依据，每个字段都要写
// jsonschema:"required"        → 必填（需配合 Reflector.RequiredFromJSONSchemaTags）
// jsonschema:"enum=a,enum=b"   → 枚举约束；minLength/maxLength/pattern 等同理

type WeatherArgs struct {
	City string `json:"city" jsonschema:"required" jsonschema_description:"城市名，如：北京"`
}

type NowArgs struct {
	TZ string `json:"tz" jsonschema_description:"IANA 时区名，如 Asia/Shanghai；省略则用服务器本地时间"` // 没有 required → 可选字段
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

// generateSchema 把 Go struct 反射成 OpenAI 工具参数 schema。
// 启动时调用一次并缓存即可，不要放在请求路径上反射。
func generateSchema[T any]() openai.FunctionParameters {
	reflector := jsonschema.Reflector{
		RequiredFromJSONSchemaTags: true,  // 必填由 jsonschema:"required" 标签决定（默认规则是"无 omitempty 即必填"，太隐晦）
		AllowAdditionalProperties:  false, // 生成 additionalProperties:false，配合 strict 拒绝模型幻觉出的字段
		DoNotReference:             true,  // 嵌套 struct 内联展开，而不是拆到 $defs 再 $ref（模型对 $ref 支持参差）
	}
	s := reflector.Reflect(new(T))

	// 只取三件套重建顶层对象，丢掉 $schema/$id 等模型不关心的键
	// （openai-go 官方 README 的结构化输出示例就是同款做法）
	schema := map[string]any{
		"type":       "object",
		"properties": s.Properties, // 有序 map：模型看到的字段顺序和 struct 声明顺序一致
		// 拒绝 schema 之外的字段；开 strict 时这是硬性要求
		"additionalProperties": false,
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

// ---- 执行函数：和参数 struct 在一起放，改哪个都能一眼看到另一个 ----

func toolWeather(_ context.Context, arguments string) (string, error) {
	var args WeatherArgs
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("参数不是合法 JSON: %w", err)
	}
	// 真实场景这里去调天气 API。结果建议用 JSON 字符串返回：
	// 模型对结构化结果的理解比随意格式的文本更稳定。
	return fmt.Sprintf(`{"city":%q,"weather":"晴","temp_c":26}`, args.City), nil
}

func toolNow(_ context.Context, arguments string) (string, error) {
	// 无参/可选参数工具的 arguments 可能是空串，补一个 {} 再解
	if arguments == "" {
		arguments = "{}"
	}
	var args NowArgs
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", err
	}
	return fmt.Sprintf(`{"now":%q}`, time.Now().Format(time.RFC3339)), nil
}

// toolCreatePresentation 是 demo 里唯一有副作用的工具：真的写出 HTML 文件，
// 运行后可在 backend 目录找到 demo_presentation.html。
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
				// Strict: openai.Bool(true), // OpenAI 端可开：服务端保证参数严格符合 schema。
				// 部分第三方兼容厂商不认这个字段，报错就去掉。
			}),
			Execute: toolCreatePresentation,
		},
	}
}

// =====================================================================
// 第 2 部分：Agent = 消息历史 + [流式请求 → 执行工具 → 回填结果] 的循环
// =====================================================================

const maxTurns = 10 // 防止模型无限调工具

type Agent struct {
	client   *openai.Client
	model    string
	messages []openai.ChatCompletionMessageParamUnion
	tools    []openai.ChatCompletionToolUnionParam
	exec     map[string]ToolFunc
}

func NewAgent(client *openai.Client, model string, tools map[string]Tool) *Agent {
	a := &Agent{client: client, model: model, exec: map[string]ToolFunc{}}
	for name, t := range tools {
		a.tools = append(a.tools, t.Definition)
		a.exec[name] = t.Execute
	}
	return a
}

// Run 跑一轮完整对话。onDelta 在每收到一小段文本时被调用——
// 这就是你将来接到 Gin SSE/WebSocket 往前端推的挂载点。
func (a *Agent) Run(ctx context.Context, userMsg string, onDelta func(string)) (string, error) {
	a.messages = append(a.messages, openai.UserMessage(userMsg))

	for turn := 1; turn <= maxTurns; turn++ {
		acc, err := a.step(ctx, onDelta)
		if err != nil {
			return "", fmt.Errorf("第 %d 轮请求失败: %w", turn, err)
		}

		// 累加器内嵌了完整响应对象 ChatCompletion（streamaccumulator.go:20），
		// 流结束后 acc.Choices[0] 就是非流式 completion.Choices[0] 的等价物
		if len(acc.Choices) == 0 {
			return "", errors.New("响应中没有 choices")
		}
		choice := acc.Choices[0]
		if choice.FinishReason == "length" {
			return "", errors.New("输出被 token 上限截断（注意：推理 token 也计入 max_completion_tokens）")
		}
		msg := choice.Message

		// 没有工具调用 → 模型已经给出最终回答，循环结束
		if len(msg.ToolCalls) == 0 {
			a.messages = append(a.messages, msg.ToParam())
			if msg.Content == "" {
				return "", errors.New("模型返回了空内容（常见原因：输出 token 上限太小，全花在思考上）")
			}
			return msg.Content, nil
		}

		// 关键规则 1：assistant 消息（含 tool_calls）必须先原样回填历史
		a.messages = append(a.messages, msg.ToParam())

		// 关键规则 2：逐个执行工具，用 role:"tool" + tool_call_id 一一对应回填。
		// 一条 assistant 消息可以并行带多个 tool_calls（比如同时查两个城市）。
		for _, tc := range msg.ToolCalls {
			result, err := a.callTool(ctx, tc.Function.Name, tc.Function.Arguments)
			if err != nil {
				// 工具失败不要中断流程：把错误告诉模型，让它自己决定下一步
				result = fmt.Sprintf(`{"error":%q}`, err.Error())
			}
			a.messages = append(a.messages, openai.ToolMessage(result, tc.ID))
		}
		// 继续下一轮：模型看到工具结果后，要么继续调工具，要么给出最终回答
	}
	return "", errors.New("超过最大轮数，模型可能陷入工具调用循环")
}

// step 发一次流式请求，把 SSE 分片实时转发给 onDelta，同时累加成完整消息。
func (a *Agent) step(ctx context.Context, onDelta func(string)) (*openai.ChatCompletionAccumulator, error) {
	params := openai.ChatCompletionNewParams{
		Model:    openai.ChatModel(a.model),
		Messages: a.messages,
		Tools:    a.tools, // 一个工具都没注册时，这就是普通流式聊天
		StreamOptions: openai.ChatCompletionStreamOptionsParam{
			IncludeUsage: openai.Bool(true), // 流式下默认拿不到 usage，要显式开
		},
	}

	stream := a.client.Chat.Completions.NewStreaming(ctx, params)
	defer stream.Close()

	acc := openai.ChatCompletionAccumulator{}
	for stream.Next() {
		chunk := stream.Current()
		if !acc.AddChunk(chunk) {
			return nil, errors.New("流式分片累加失败")
		}
		// 文本增量：平时前端"打字机效果"的来源。
		// 工具调用的 arguments 也是分片下发的（按 index 拼接），
		// AddChunk 已经帮你拼好了，不用自己处理。
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			onDelta(chunk.Choices[0].Delta.Content)
		}
		// 事件钩子：某个工具调用的 name+arguments 刚在流里拼完整。
		// 官方文档警告：开启并行工具调用时不可依赖此事件（streamaccumulator.go:358），
		// 所以正式逻辑应以流结束后 acc.Choices[0].Message.ToolCalls 为准，
		// 这个事件只适合打日志/调试。
		if tool, ok := acc.JustFinishedToolCall(); ok {
			fmt.Fprintf(os.Stderr, "[事件] 工具调用接收完毕: %s(%s)\n", tool.Name, tool.Arguments)
		}
	}
	if err := stream.Err(); err != nil {
		return nil, err // 网络层/HTTP 层错误，调用方用 errors.As(*openai.Error) 细分
	}
	if acc.Usage.TotalTokens > 0 {
		fmt.Fprintf(os.Stderr, "[用量] 输入 %d + 输出 %d = %d tokens\n",
			acc.Usage.PromptTokens, acc.Usage.CompletionTokens, acc.Usage.TotalTokens)
	}
	return &acc, nil
}

func (a *Agent) callTool(ctx context.Context, name, arguments string) (string, error) {
	fn, ok := a.exec[name]
	if !ok {
		return "", fmt.Errorf("未知工具 %q", name)
	}
	return fn(ctx, arguments)
}

// =====================================================================
// 第 3 部分：把零件跑起来——两个场景对比
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

	client := openai.NewClient(opts...) // 注意：返回值类型是 openai.Client（值类型）
	agent := NewAgent(&client, model, buildTools())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// 场景 1：不涉及工具的普通流式对话（和场景 2 走同一个 Run，只是模型不调工具）
	fmt.Println("=== 场景 1：普通流式对话 ===")
	answer, err := agent.Run(ctx, "用两句话介绍一下你自己", func(delta string) {
		fmt.Print(delta) // 打字机效果
	})
	if err != nil {
		fmt.Println("出错:", err)
	} else {
		fmt.Printf("\n[最终回答] %s\n\n", answer)
	}

	// 场景 2：需要工具的问题。注意 stderr 上的 [事件] 行：
	// 模型会先"说话"（或直接发起并行工具调用，一次查两个城市），
	// 你的代码执行工具并回填后，模型才基于结果输出最终回答——这就是 agent 循环。
	fmt.Println("=== 场景 2：工具调用 ===")
	answer, err = agent.Run(ctx, "北京和上海今天天气怎么样？顺便现在几点了？", func(delta string) {
		fmt.Print(delta)
	})
	if err != nil {
		fmt.Println("出错:", err)
	} else {
		fmt.Printf("\n[最终回答] %s\n", answer)
	}

	// 场景 3：struct 反射出的嵌套 schema + 有副作用的工具。
	// 运行后在 backend 目录找 demo_presentation.html。
	fmt.Println("=== 场景 3：生成 PPT（嵌套 schema + 副作用工具） ===")
	answer, err = agent.Run(ctx, "帮我做一份关于 Go 并发编程的 3 页 PPT，主题用 dark", func(delta string) {
		fmt.Print(delta)
	})
	if err != nil {
		fmt.Println("出错:", err)
	} else {
		fmt.Printf("\n[最终回答] %s\n", answer)
	}
}
