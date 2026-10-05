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
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/openai/openai-go/v3"

	"html-ppt/backend/internal/service/template"
	"html-ppt/backend/internal/store"
	"html-ppt/backend/internal/trace"
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
		Name:        "set_layout_roles",
		Description: openai.String("改某个版式的适用场景（role，覆盖式不是增量）。roles 是词表的子集、最多 3 个：cover（封面）/toc（目录）/divider（章节）/content（正文）/data（数据）/quote（金句）/code（代码）/cta（行动号召）/thanks（收尾）。用户说「这个版式也能用在数据页」「封面不要用这个」时用它；空数组 = 清空限定（任何场景按内容性质选用）。layout 必须是结构契约里登记的版式 id。"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"layout": map[string]any{"type": "string", "description": "版式 id（结构契约里登记的，如 blank-data）"},
				"roles": map[string]any{
					"type": "array",
					"items": map[string]any{
						"type": "string",
						"enum": []string{"cover", "toc", "divider", "content", "data", "quote", "code", "cta", "thanks"},
					},
					"description": "新的角色清单（整组替换，最多 3 个）；空数组 = 清空限定",
				},
			},
			"required": []string{"layout", "roles"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "set_layout_meta",
		Description: openai.String("改某个版式给生成模型看的名称与用途描述（直接影响生成时的版式选择，写得越具体模型选得越准）。name ≤40 字、use ≤200 字；不改的字段省略。layout 必须是结构契约里登记的版式 id。"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"layout": map[string]any{"type": "string", "description": "版式 id（结构契约里登记的）"},
				"name":   map[string]any{"type": "string", "description": "新版式名称；不改就省略"},
				"use":    map[string]any{"type": "string", "description": "新用途描述（一句话说清什么内容适合用它）；不改就省略"},
			},
			"required": []string{"layout"},
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
//
// emit 为 SSE 客户端的事件出口（nil = 同步调用，只落观测不推流）。
// 整轮同时落 trace 事件（customize_trace.go），观测台可见。
func (s *Service) Customize(ctx context.Context, userID uint, id, message string, llm LLM, emit func(CustEvent) error) (string, error) {
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

	rec := s.openCustRecorder(userID, row, message, llm.Model)
	defer rec.Close()
	ctx = trace.With(ctx, rec)

	reply, dirty, note, loopErr := s.customizeLoop(ctx, sess, row, llm, emit)
	if loopErr != nil {
		// 中途出错的 run 也要收尾：不落 run_end，观测页会永远停在"运行中"并一直轮询
		trace.Emit(ctx, trace.Event{Kind: trace.KindError, Error: loopErr.Error()})
		sum := rec.Summary()
		trace.Emit(ctx, trace.Event{Kind: trace.KindRunEnd, Status: trace.StatusError, Summary: &sum})
		return "", loopErr
	}
	sum := rec.Summary()
	trace.Emit(ctx, trace.Event{Kind: trace.KindRunEnd, Status: trace.StatusOK, Summary: &sum})
	if dirty {
		_ = s.recordVersionUT(row.ID, OpChat, note)
	}
	if emit != nil {
		// 收尾帧：答复全文 + dirty（前端据此把流式草稿落成正式消息、刷新预览与历史）
		if emitErr := emit(CustEvent{Type: CustEvDone, Dirty: dirty, Reply: reply}); emitErr != nil {
			return "", emitErr
		}
	}
	return reply, nil
}

// customizeLoop 工具循环（会话锁内调用）。返回最终答复、是否有文件写入、
// 版本备注（finish 汇报优先，兜底用答复摘要）。
// 每轮的请求/回复/工具调用/用量都在 turnCtx 归属下落 trace；emit 非 nil 时
// 工具气泡（tool_start/tool_done）与生成进度（tool_progress）实时推给客户端。
func (s *Service) customizeLoop(ctx context.Context, sess *customizeSession, row *store.UserTemplate, llm LLM, emit func(CustEvent) error) (string, bool, string, error) {
	const maxTurns = 6
	dirty := false
	note := ""
	retriedEmpty := false
	for turn := 0; turn < maxTurns; turn++ {
		turnCtx := trace.WithTurn(ctx, turn)
		msg, finish, _, err := s.custStream(turnCtx, sess, llm, emit)
		if err != nil {
			return "", false, "", fmt.Errorf("模型调用失败: %w", err)
		}
		// 推理模型把输出预算全花在 reasoning 上（finish=length、零正文零工具）
		// 时，回退这轮重试一次；不重试用户看到的就是"（本轮无回复）"。
		if len(msg.ToolCalls) == 0 && strings.TrimSpace(msg.Content) == "" && finish == "length" && !retriedEmpty {
			retriedEmpty = true
			turn--
			continue
		}
		sess.messages = append(sess.messages, msg.ToParam())

		if len(msg.ToolCalls) == 0 {
			// 没调工具直接回话：视为最终答复
			reply := strings.TrimSpace(msg.Content)
			if reply == "" {
				return "", false, "", fmt.Errorf("模型返回空响应，请重试")
			}
			if note == "" {
				note = truncateNote(reply)
			}
			return reply, dirty, note, nil
		}
		replied := ""
		for _, tc := range msg.ToolCalls {
			// 归属挂进 ctx：工具调用/返回两条事件自动带上轮次与 tool_call_id
			toolCtx := trace.WithTool(turnCtx, tc.Function.Name, tc.ID)
			trace.Emit(toolCtx, trace.Event{Kind: trace.KindToolCall, Args: tc.Function.Arguments})
			started := time.Now()
			result, d := s.execCustomTool(row, tc.Function.Name, tc.Function.Arguments)
			trace.Emit(toolCtx, trace.Event{
				Kind: trace.KindToolResult, Result: result,
				DurationMS: time.Since(started).Milliseconds(),
			})
			if d {
				dirty = true
			}
			if emit != nil {
				if emitErr := emit(CustEvent{Type: CustEvToolDone, ToolCallID: tc.ID, ToolName: tc.Function.Name, Content: toolBrief(result)}); emitErr != nil {
					return "", false, "", fmt.Errorf("事件推送失败: %w", emitErr)
				}
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
	// 轮次上限：模型连续多轮调工具没调 finish（write_demo 这类大输出常见）。
	// 备注兜底不能空着——历史列表里空 detail 不可读。
	reply := "本轮改动已应用（达到对话轮次上限）。可以继续描述，或去预览确认效果。"
	if note == "" {
		note = truncateNote(reply)
	}
	return reply, dirty, note, nil
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
		if err := s.remountUT(row.ID); err != nil {
			return fmt.Sprintf("style.css 已写入并记入历史，但注册表校验未过：%v。预览可见，但生成侧可能仍挂旧版——请检查是否破坏了版式契约或文件结构（可用质量体检核对）。", err), true
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
		if err := s.remountUT(row.ID); err != nil {
			return fmt.Sprintf("index.html 已写入并记入历史，但注册表校验未过：%v。预览可见，但生成侧可能仍挂旧版——最常见原因是发明了 layouts.md 里没有的 data-layout，请改回已登记的版式。", err), true
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
	case "set_layout_roles", "set_layout_meta":
		var in struct {
			Layout string   `json:"layout"`
			Roles  []string `json:"roles"`
			Name   *string  `json:"name"`
			Use    *string  `json:"use"`
		}
		if err := json.Unmarshal([]byte(args), &in); err != nil {
			return "参数不是合法 JSON: " + err.Error(), false
		}
		patch := LayoutMetaPatch{}
		if name == "set_layout_roles" {
			patch.Roles = &in.Roles
		} else {
			if in.Name == nil && in.Use == nil {
				return "name/use 至少传一个（只改适用场景请用 set_layout_roles）", false
			}
			patch.Name, patch.Use = in.Name, in.Use
		}
		warning, err := s.UpdateLayoutMeta(row.UserID, row.ID, in.Layout, patch)
		if err != nil {
			return "修改失败: " + err.Error(), false
		}
		if warning != "" {
			return warning, true
		}
		return "版式元数据已更新（template.json 为权威，注册表已重挂，生成侧立即生效）。", true
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
	// 已挂载的模板重挂一次，让注册表里的 style.css 快照同步；失败不阻断
	// （token 写入本身已成功，校验问题由发布门禁终审）
	if err := s.remountUT(id); err != nil {
		log.Printf("[warn] 模板 %s token 写入后重挂失败: %v", id, err)
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
	if err := s.remountUT(id); err != nil {
		log.Printf("[warn] 模板 %s 元数据写入后重挂失败: %v", id, err)
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
		"- 改版式的适用场景（role）或名称/用途：set_layout_roles / set_layout_meta——" +
		"layout 传下方结构契约里登记的版式 id。\n" +
		"- 改模板名/描述：set_meta。\n\n" +
		"纪律：\n" +
		"- 一次改动要成套（改主色时同步考虑 --accent-2/--accent-3 的协调，以及文字在新底色上的可读性）。\n" +
		"- 结构契约：不要发明 layouts.md 里没有的版式；类名沿用模板既有体系。role 只能从词表选：" +
		"cover/toc/divider/content/data/quote/code/cta/thanks，每个版式最多 3 个。\n" +
		"- 安全预检会拒绝 CSS 网络外链与额外脚本；被拒时按报错修正重试，不要换个写法绕。\n" +
		"- 改动前不需要向用户确认——直接改，然后在 finish 里用中文具体说明改了哪些内容、建议用户看哪一页验证。\n" +
		"- 模板名：" + row.Name + "（base：" + row.BaseID + "）。\n\n" +
		"结构契约（版式清单：id · 名称 · 角色 · 用途）：\n" + s.layoutContractSummary(row.ID) + "\n\n" +
		"当前 style.css 全文：\n" + styleAll
}

// layoutContractSummary 读 template.json 的 layouts 清单拼成模型可读的结构摘要
// （system prompt 注入用；读盘而非注册表——会话里刚发生的修改也能反映）。
// 读不到时返回占位说明，不让定制对话直接失败。
func (s *Service) layoutContractSummary(id string) string {
	raw, err := os.ReadFile(filepath.Join(s.Dir(id), "template.json"))
	if err != nil {
		return "（读不到 template.json，本模板没有版式元数据可改）"
	}
	var meta struct {
		Layouts []template.LayoutMeta `json:"layouts"`
	}
	if err := json.Unmarshal(raw, &meta); err != nil {
		return "（template.json 解析失败，本模板没有版式元数据可改）"
	}
	var b strings.Builder
	for _, l := range meta.Layouts {
		roles := "无（不限场景）"
		if len(l.Roles) > 0 {
			roles = strings.Join(l.Roles, "/")
		}
		use := l.Use
		if use == "" {
			use = "（未填用途）"
		}
		fmt.Fprintf(&b, "- %s（%s）· 角色：%s · %s\n", l.ID, l.Name, roles, use)
	}
	return strings.TrimRight(b.String(), "\n")
}

var bodyClassRe = regexp.MustCompile(`<body class="([^"]*)">`)
