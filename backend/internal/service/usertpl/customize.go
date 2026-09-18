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

	const maxTurns = 6
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
			return "", fmt.Errorf("模型调用失败: %w", err)
		}
		if len(completion.Choices) == 0 {
			return "", fmt.Errorf("模型返回空响应")
		}
		msg := completion.Choices[0].Message
		sess.messages = append(sess.messages, msg.ToParam())

		if len(msg.ToolCalls) == 0 {
			// 没调工具直接回话：视为最终答复
			return strings.TrimSpace(msg.Content), nil
		}
		replied := ""
		for _, tc := range msg.ToolCalls {
			result := s.execCustomTool(row, tc.Function.Name, tc.Function.Arguments)
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
				if out.Reply != "" {
					return out.Reply, nil
				}
				return replied, nil
			}
		}
	}
	return "本轮改动已应用（达到对话轮次上限）。可以继续描述，或去预览确认效果。", nil
}

// execCustomTool 执行一个定制工具，返回给模型的结果文本。
func (s *Service) execCustomTool(row *store.UserTemplate, name, args string) string {
	switch name {
	case "write_tokens":
		var in struct {
			Tokens map[string]string `json:"tokens"`
		}
		if err := json.Unmarshal([]byte(args), &in); err != nil {
			return "参数不是合法 JSON: " + err.Error()
		}
		if err := s.writeTokenBlock(row.ID, row.BaseID, in.Tokens); err != nil {
			return "写入失败: " + err.Error()
		}
		return "已写入 " + fmt.Sprint(len(in.Tokens)) + " 个 token，预览刷新即可看到。"
	case "set_meta":
		var in struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		if err := json.Unmarshal([]byte(args), &in); err != nil {
			return "参数不是合法 JSON: " + err.Error()
		}
		if err := s.applyMetaFile(row.ID, in.Name, in.Description); err != nil {
			return "写入失败: " + err.Error()
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
				return "入库失败: " + err.Error()
			}
		}
		return "名称/描述已更新。"
	case "finish":
		var out struct {
			Reply string `json:"reply"`
		}
		_ = json.Unmarshal([]byte(args), &out)
		return out.Reply
	default:
		return "未知工具 " + name
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
	return s.reg.MountUser(dir)
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
	return s.reg.MountUser(s.Dir(id))
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

// customizeSystemPrompt 定制对话的系统提示词（含当前 token 值，模型有上下文）。
func (s *Service) customizeSystemPrompt(row *store.UserTemplate) string {
	tokens := ""
	if raw, err := os.ReadFile(filepath.Join(s.Dir(row.ID), "style.css")); err == nil {
		if m := tokenBlockRe.FindStringSubmatch(string(raw)); m != nil {
			tokens = m[1]
		}
	}
	return "你是 PPT 模板定制助手。用户会用自然语言描述想改的视觉风格（配色、圆角、观感），" +
		"你通过 write_tokens 工具落地，改完用 finish 汇报。\n\n" +
		"纪律：\n" +
		"- 只改列出的 token；一次改动要成套（比如改主色时同步考虑 --accent-2/--accent-3 的协调，以及文字在新底色上的可读性）。\n" +
		"- 用户说「更圆/直角」时调 --radius/--radius-lg；说「暗色/亮色」时成套改 --bg/--surface/--text-*。\n" +
		"- 改动前不需要向用户确认——直接改，然后在 finish 里用中文具体说明改了哪些 token、建议用户看哪一页验证。\n" +
		"- 模板名：" + row.Name + "（base：" + row.BaseID + "）。\n\n" +
		"当前 token 值（style.css 的 token 声明区）：\n" + tokens
}

var bodyClassRe = regexp.MustCompile(`<body class="([^"]*)">`)

var tokenBlockRe = regexp.MustCompile(`(?s)\.tpl-[a-z0-9-]+\{([^}]*)\}`)
