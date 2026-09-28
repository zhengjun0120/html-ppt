package usertpl

// 对话定制（plan-v3 B2）：用户用自然语言改模板（色板/字体/圆角/名称描述），
// agent 通过受限工具改 style.css 的 token 覆盖块与 template.json 元数据。
//
// 设计取舍：这是一个 ~6 轮上限的小循环，不是完整的 deck 生成管线——
// 可改面被刻意收窄到 token 级（D1 克隆定制的边界），结构契约（版式/类名）
// 不可动，从源头上杜绝"对话把模板改坏"。

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/openai/openai-go/v3"

	"html-ppt/backend/internal/store"
)

// LLM 定制对话的模型接入（由装配层从 agent 服务构造，避免包依赖环）。
type LLM struct {
	Client *openai.Client
	Model  string
}

// CustomizeRequest 一次定制对话请求。
type CustomizeRequest struct {
	Message string `json:"message"`
}

type customizeSession struct {
	mu       sync.Mutex
	messages []openai.ChatCompletionMessageParamUnion
}

func (s *Service) sessionKey(userID uint, id string) string {
	return fmt.Sprintf("%d|%s", userID, id)
}

// customizeTools 工具 schema（JSON Schema，直接给 function calling）。
var customizeTools = []openai.ChatCompletionToolUnionParam{
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "write_tokens",
		Description: openai.String("把设计 token 写进模板（后缀覆盖块，立即生效于预览）。只改列出的键；值必须是合法 CSS 值。可用的键：--accent（主强调色）--accent-2（次强调色）--accent-3（危险/反差色）--bg（页面底色）--bg-soft（次级底色）--surface（卡片底）--surface-2（次级面）--text-1（主文字）--text-2（次级文字）--text-3（弱文字）--radius（小圆角）--radius-lg（大圆角）。色值写 hex（如 #b45309）。"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"tokens": map[string]any{
					"type":                 "object",
					"description":          "键=token 名（含 -- 前缀），值=CSS 值",
					"additionalProperties": map[string]any{"type": "string"},
				},
			},
			"required": []string{"tokens"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "set_meta",
		Description: openai.String("改模板的名称与描述（展示在画廊里）"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"name":        map[string]any{"type": "string", "description": "新名称；不改就传空串"},
				"description": map[string]any{"type": "string", "description": "新描述；不改就传空串"},
			},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "write_style",
		Description: openai.String("整体重写模板的 style.css（全文替换）。必须基于系统提示里的当前全文改造，不要凭空杜撰既有规则；保留 .tpl- 作用域前缀与 token 声明区。安全预检会拒绝 url( 网络外链（仅允许 data: 内联）、@import、expression 等内容，被拒时根据报错修正后重试。"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"css": map[string]any{"type": "string", "description": "style.css 完整新全文"},
			},
			"required": []string{"css"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "write_demo",
		Description: openai.String("整体重写模板 demo 的 index.html（全文替换）。必须保留 body 上的 tpl- 作用域 class 与 /assets/deck-v2/runtime.js 脚本引用；页面是 .deck 下的 <section class=\"slide\"> 序列，版式与类名沿用模板既有体系（不要发明 layouts.md 里没有的版式）。"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"html": map[string]any{"type": "string", "description": "index.html 完整新全文"},
			},
			"required": []string{"html"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "finish",
		Description: openai.String("本轮定制结束。把做了什么、让用户去看哪里复述给用户。"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"reply": map[string]any{"type": "string", "description": "给用户的中文总结（2-4 句，具体说改了什么）"},
			},
			"required": []string{"reply"},
		},
	}),
}

// Customize 一轮对话：把用户消息追加进会话，跑到 finish 或轮次上限。
// 本轮有实际文件写入时记一条 chat 版本（docs/user-template-history-plan.md §3.2），
// 备注用 finish 的汇报——历史列表因此可读。
func (s *Service) Customize(ctx context.Context, userID uint, id, message string, llm LLM) (string, error) {
	row, err := s.GetOwned(userID, id)
	if err != nil {
		return "", err
	}
	if message == "" {
		return "", fmt.Errorf("消息不能为空")
	}

	key := s.sessionKey(userID, id)
	sess := s.sessionFor(key)
	sess.mu.Lock()
	defer sess.mu.Unlock()

	if len(sess.messages) == 0 {
		sess.messages = append(sess.messages, openai.SystemMessage(s.customizeSystemPrompt(row)))
	}
	sess.messages = append(sess.messages, openai.UserMessage(message))

	reply, dirty, note, err := s.customizeLoop(ctx, sess, row, llm)
	if err != nil {
		return "", err
	}
	if dirty {
		_ = s.recordVersionUT(row.ID, OpChat, note)
	}
	return reply, nil
}

// customizeLoop 工具循环（会话锁内调用）。返回最终答复、是否有文件写入、
// 版本备注（finish 汇报优先，兜底用答复摘要）。
func (s *Service) customizeLoop(ctx context.Context, sess *customizeSession, row *store.UserTemplate, llm LLM) (string, bool, string, error) {
	const maxTurns = 6
	dirty := false
	note := ""
	for turn := 0; turn < maxTurns; turn++ {
		completion, err := llm.Client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
			Model:    openai.ChatModel(llm.Model),
			Messages: sess.messages,
			Tools:    customizeTools,
			ToolChoice: openai.ChatCompletionToolChoiceOptionUnionParam{
				OfAuto: openai.String("auto"),
			},
			MaxTokens:   openai.Int(4000),
			Temperature: openai.Float(0.4),
		})
		if err != nil {
			return "", false, "", fmt.Errorf("模型调用失败: %w", err)
		}
		if len(completion.Choices) == 0 {
			return "", false, "", fmt.Errorf("模型返回空响应")
		}
		msg := completion.Choices[0].Message
		sess.messages = append(sess.messages, msg.ToParam())

		if len(msg.ToolCalls) == 0 {
			// 没调工具直接回话：视为最终答复
			reply := strings.TrimSpace(msg.Content)
			if note == "" {
				note = truncateNote(reply)
			}
			return reply, dirty, note, nil
		}
		replied := ""
		for _, tc := range msg.ToolCalls {
			result, d := s.execCustomTool(row, tc.Function.Name, tc.Function.Arguments)
			if d {
				dirty = true
			}
			replied = result
			sess.messages = append(sess.messages, openai.ToolMessage(result, tc.ID))
		}
		// finish 工具 = 本轮结束
		for _, tc := range msg.ToolCalls {
			if tc.Function.Name == "finish" {
				var out struct {
					Reply string `json:"reply"`
				}
				_ = json.Unmarshal([]byte(tc.Function.Arguments), &out)
				reply := out.Reply
				if reply == "" {
					reply = replied
				}
				if note == "" {
					note = truncateNote(reply)
				}
				return reply, dirty, note, nil
			}
		}
	}
	return "本轮改动已应用（达到对话轮次上限）。可以继续描述，或去预览确认效果。", dirty, note, nil
}

// truncateNote 版本备注上限（index.json 里的人读字段，太长列表没法看）。
func truncateNote(s string) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) > 120 {
		return string(r[:120]) + "…"
	}
	if len(r) == 0 {
		return "对话定制"
	}
	return string(r)
}

// execCustomTool 执行一个定制工具，返回给模型的结果文本与"是否写了文件"
//（dirty = 本轮对话需要记一条历史版本）。
func (s *Service) execCustomTool(row *store.UserTemplate, name, args string) (string, bool) {
	// D2 已发布锁定（修现状问题 C：此前 write_tokens 可绕过门禁改已发布模板）。
	// 门禁运行态（publishing）同理——量测期间内容不能变。finish 不受影响。
	if name != "finish" && (row.Status == "published" || row.Status == "publishing") {
		return "该模板已发布（或正在跑发布门禁），请先下架再修改。", false
	}
	switch name {
	case "write_tokens":
		var in struct {
			Tokens map[string]string `json:"tokens"`
		}
		if err := json.Unmarshal([]byte(args), &in); err != nil {
			return "参数不是合法 JSON: " + err.Error(), false
		}
		if err := s.writeTokenBlock(row.ID, row.BaseID, in.Tokens); err != nil {
			return "写入失败: " + err.Error(), false
		}
		return "已写入 " + fmt.Sprint(len(in.Tokens)) + " 个 token，预览刷新即可看到。", len(in.Tokens) > 0
	case "write_style":
		var in struct {
			CSS string `json:"css"`
		}
		if err := json.Unmarshal([]byte(args), &in); err != nil {
			return "参数不是合法 JSON: " + err.Error(), false
		}
		if len(in.CSS) > maxToolFileBytes {
			return fmt.Sprintf("style.css 内容 %d 字节，超过 %d 上限，请精简。", len(in.CSS), maxToolFileBytes), false
		}
		if err := scanCSS(in.CSS); err != nil {
			return "安全预检未过：" + err.Error() + "。请修正后重试。", false
		}
		if err := s.writeTemplateFile(row.ID, "style.css", in.CSS); err != nil {
			return "写入失败: " + err.Error(), false
		}
		return "style.css 已整体更新，预览刷新即可看到。", true
	case "write_demo":
		var in struct {
			HTML string `json:"html"`
		}
		if err := json.Unmarshal([]byte(args), &in); err != nil {
			return "参数不是合法 JSON: " + err.Error(), false
		}
		if len(in.HTML) > maxToolFileBytes {
			return fmt.Sprintf("index.html 内容 %d 字节，超过 %d 上限，请精简。", len(in.HTML), maxToolFileBytes), false
		}
		if err := scanHTML(in.HTML); err != nil {
			return "安全预检未过：" + err.Error() + "。请修正后重试。", false
		}
		if err := s.writeTemplateFile(row.ID, "index.html", in.HTML); err != nil {
			return "写入失败: " + err.Error(), false
		}
		return "demo index.html 已整体更新，预览刷新即可看到。", true
	case "set_meta":
		var in struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		if err := json.Unmarshal([]byte(args), &in); err != nil {
			return "参数不是合法 JSON: " + err.Error(), false
		}
		if err := s.applyMetaFile(row.ID, in.Name, in.Description); err != nil {
			return "写入失败: " + err.Error(), false
		}
		updates := map[string]any{}
		if in.Name != "" {
			updates["name"] = in.Name
		}
		if in.Description != "" {
			updates["description"] = in.Description
		}
		if len(updates) > 0 {
			if err := s.st.DB.Model(row).Updates(updates).Error; err != nil {
				return "入库失败: " + err.Error(), false
			}
		}
		return "名称/描述已更新。", len(updates) > 0
	case "finish":
		var out struct {
			Reply string `json:"reply"`
		}
		_ = json.Unmarshal([]byte(args), &out)
		return out.Reply, false
	default:
		return "未知工具 " + name, false
	}
}

// customizeOverrideMark 定制覆盖块的标记注释（重写整块用）。
const (
	customizeBegin = "/* [customize] 用户定制 token 覆盖（对话自动生成，重发整块） */"
	customizeEnd   = "/* [/customize] */"
)

var customizeBlockRe = regexp.MustCompile(`(?s)/\* \[customize\][^*]*\*/.*?/\* \[/customize\] \*/\n?`)

// writeTokenBlock 把 token 覆盖块写进 style.css（已有则整块替换，没有则追加）。
// 覆盖块位于文件末尾且带模板作用域前缀，优先级高于默认 token；写完重挂注册表，
// 让"已发布"的模板也立即生效。
func (s *Service) writeTokenBlock(id, baseID string, tokens map[string]string) error {
	if len(tokens) == 0 {
		return nil
	}
	scope := s.scopeOf(id, baseID)
	var b strings.Builder
	b.WriteString(customizeBegin + "\n")
	b.WriteString("." + scope + " {\n")
	for k, v := range tokens {
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		if !strings.HasPrefix(k, "--") || strings.ContainsAny(v, ";{}") {
			continue // 非法键值直接丢弃（防注入）
		}
		fmt.Fprintf(&b, "  %s: %s;\n", k, v)
	}
	b.WriteString("}\n" + customizeEnd + "\n")

	dir := s.Dir(id)
	p := filepath.Join(dir, "style.css")
	raw, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	out := customizeBlockRe.ReplaceAllString(string(raw), "")
	out = out + "\n" + b.String()
	if err := os.WriteFile(p, []byte(out), 0o644); err != nil {
		return err
	}
	// 已挂载的模板（fork 即挂载）重挂一次，让注册表里的 style.css 快照同步
	if s.reg != nil {
		return s.reg.MountUser(dir)
	}
	return nil
}

// applyMetaFile 把对话里的改名/描述同步进 template.json（注册表重挂后生效）。
func (s *Service) applyMetaFile(id, name, description string) error {
	p := filepath.Join(s.Dir(id), "template.json")
	raw, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	var meta map[string]any
	if err := json.Unmarshal(raw, &meta); err != nil {
		return err
	}
	if name != "" {
		meta["name"] = name
	}
	if description != "" {
		meta["description"] = description
	}
	out, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(p, out, 0o644); err != nil {
		return err
	}
	if s.reg != nil {
		return s.reg.MountUser(s.Dir(id))
	}
	return nil
}

// scopeOf 模板的 CSS 作用域类（body class）。fork 沿用 base 的作用域
// （如 data-dark 的 tpl-graphify-dark-graph），读取 index.html 的 body class 判定。
func (s *Service) scopeOf(id, baseID string) string {
	raw, err := os.ReadFile(filepath.Join(s.Dir(id), "index.html"))
	if err == nil {
		if m := bodyClassRe.FindStringSubmatch(string(raw)); m != nil {
			for _, c := range strings.Fields(m[1]) {
				if strings.HasPrefix(c, "tpl-") {
					return c
				}
			}
		}
	}
	return "tpl-" + baseID
}

// maxToolFileBytes 对话整文件重写工具（write_style/write_demo）的内容上限。
// style.css 实测量级 20-60KB、demo 30KB；512KB 给足余量同时防 token 失控。
const maxToolFileBytes = 512 << 10

// customizeSystemPrompt 定制对话的系统提示词（注入当前 style.css 全文，模型有上下文）。
func (s *Service) customizeSystemPrompt(row *store.UserTemplate) string {
	styleAll := ""
	if raw, err := os.ReadFile(filepath.Join(s.Dir(row.ID), "style.css")); err == nil {
		styleAll = string(raw)
	}
	if r := []rune(styleAll); len(r) > 8000 {
		styleAll = string(r[:8000]) + "\n/* …（已截断，重写时以磁盘现状为准并保留未展示部分的结构） */"
	}
	return "你是 PPT 模板定制助手。用户会用自然语言描述想改的视觉风格与页面结构，" +
		"你通过工具落地，改完用 finish 汇报。\n\n" +
		"工具选择：\n" +
		"- 改配色/圆角等设计 token：优先 write_tokens（精准、影响面小）。\n" +
		"- token 表达不了的样式改动（布局、间距、装饰、字体规则）：用 write_style 整体重写 " +
		"style.css——必须基于下方当前全文改造，输出完整新全文。\n" +
		"- 改 demo 页面结构（加删页面/卡片、改布局骨架）：用 write_demo 整体重写 index.html——" +
		"必须保留 body 的 tpl- 作用域类与 /assets/deck-v2/runtime.js 引用。\n" +
		"- 改模板名/描述：set_meta。\n\n" +
		"纪律：\n" +
		"- 一次改动要成套（改主色时同步考虑 --accent-2/--accent-3 的协调，以及文字在新底色上的可读性）。\n" +
		"- 结构契约：不要发明 layouts.md 里没有的版式；类名沿用模板既有体系。\n" +
		"- 安全预检会拒绝 CSS 网络外链与额外脚本；被拒时按报错修正重试，不要换个写法绕。\n" +
		"- 改动前不需要向用户确认——直接改，然后在 finish 里用中文具体说明改了哪些内容、建议用户看哪一页验证。\n" +
		"- 模板名：" + row.Name + "（base：" + row.BaseID + "）。\n\n" +
		"当前 style.css 全文：\n" + styleAll
}

var bodyClassRe = regexp.MustCompile(`<body class="([^"]*)">`)
