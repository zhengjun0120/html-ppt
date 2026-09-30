// Package tplsuggest 选模板阶段的模板推荐：聚大纲与用户诉求喂 LLM，
// 从内置模板里挑 3-5 个，结果缓存进 deck.json（deck 层负责存取）。
//
// 依赖注入照 usertpl 的先例：LLM client/model 由装配层从 agent 服务构造
//（handler 里 h.agent.CustomizeLLMFor），本包不 import agent，避免依赖环。
package tplsuggest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	"html-ppt/backend/internal/service/deck"
	"html-ppt/backend/internal/service/template"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

// LLM 一次补全的接入参数（usertpl.LLM 同款形状）。
type LLM struct {
	Client *openai.Client
	Model  string
}

// Service 模板推荐服务。reg 只读元数据；decks 负责守卫依据（stage）、
// 大纲读取与推荐缓存落盘。
type Service struct {
	decks *deck.Service
	reg   *template.Registry
}

func New(decks *deck.Service, reg *template.Registry) *Service {
	return &Service{decks: decks, reg: reg}
}

// suggestTimeout 一次推荐 LLM 调用的上限。选模板页是异步加载推荐、
// 网格照常可用，但也不能让请求干等到请求方超时。
const suggestTimeout = 45 * time.Second

// maxSuggestions 最多推荐条数；reason 上限（rune）。
const (
	maxSuggestions = 5
	maxReasonRunes = 60
)

// Suggest 返回 deck 的模板推荐。refresh=false 时优先读缓存；
// LLM 调用或解析失败一律返回空切片 + nil error（前端静默隐藏推荐区，
// 筛选栏兜底），只有阶段/归属这类请求本身的问题才返回错误。
func (s *Service) Suggest(ctx context.Context, uid uint, deckID string, llm LLM, clarify string, refresh bool) ([]deck.TplSuggestion, error) {
	if s.decks == nil {
		return nil, fmt.Errorf("文稿服务不可用")
	}
	if s.reg == nil {
		return nil, fmt.Errorf("模板库不可用")
	}
	if llm.Client == nil || llm.Model == "" {
		return nil, fmt.Errorf("LLM 不可用")
	}

	df, err := s.decks.GetDeckV2(uid, deckID)
	if err != nil {
		return nil, err
	}
	// 与 gate 2 同阶段窗：推荐是选模板的辅助，生成开始后再调没有意义。
	// StageMismatch 会被 handler 映射成 409（"冲突但用户可行动"）。
	if df.Stage != deck.StageSelectingTemplate {
		return nil, deck.StageMismatch{Current: df.Stage, Require: deck.StageSelectingTemplate}
	}

	if !refresh {
		if cached, err := s.decks.TemplateSuggestions(uid, deckID); err == nil && len(cached) > 0 {
			return cached, nil
		}
	}

	candidates := builtinMetas(s.reg)
	if len(candidates) == 0 {
		log.Printf("[tplsuggest] %s 候选清单为空", deckID)
		return []deck.TplSuggestion{}, nil
	}
	o, err := s.decks.ReadOutline(uid, deckID)
	if err != nil {
		// selecting_template 阶段大纲必然已确认（gate 1 前置）；读不到属于异常，
		// 但对用户仍是"给不出推荐"而不是 5xx
		log.Printf("[tplsuggest] %s 读大纲失败 err:%v", deckID, err)
		return []deck.TplSuggestion{}, nil
	}

	sys, user := buildPrompt(df, o, clarify, candidates)
	raw, err := callLLM(ctx, llm, sys, user)
	if err != nil {
		log.Printf("[tplsuggest] %s LLM 调用失败 err:%v", deckID, err)
		return []deck.TplSuggestion{}, nil
	}
	sugs := sanitize(raw, candidates)
	if len(sugs) == 0 {
		log.Printf("[tplsuggest] %s 无合法推荐，原始输出：%.80s", deckID, raw)
		return []deck.TplSuggestion{}, nil
	}
	if err := s.decks.SaveTemplateSuggestions(uid, deckID, sugs); err != nil {
		log.Printf("[tplsuggest] %s 写缓存失败 err:%v", deckID, err) // 推荐照常返回，缓存下次重建
	}
	return sugs, nil
}

// builtinMetas 内置模板元数据（候选白名单的来源）。用户模板（ut- 前缀）
// 本轮不进候选——可见性因人而异，推荐了别人看不到就是废条目。
func builtinMetas(reg *template.Registry) []*template.Meta {
	all := reg.List()
	out := make([]*template.Meta, 0, len(all))
	for _, m := range all {
		if strings.HasPrefix(m.ID, "ut-") {
			continue
		}
		out = append(out, m)
	}
	return out
}

// buildPrompt 组推荐请求：system 定格式契约，user 段按 标题→大纲→诉求→候选 排。
// 候选行带变体清单，模型才有给 variant_id 的依据。
func buildPrompt(df *deck.DeckFile, o *deck.Outline, clarify string, candidates []*template.Meta) (string, string) {
	var sys strings.Builder
	sys.WriteString("你是 PPT 模板推荐助手。根据文稿大纲与用户在澄清对话里表达的诉求，从候选清单中挑 3-5 个最合适的模板。\n" +
		"只输出 JSON 数组，不要输出任何其他文字：\n" +
		`[{"template_id":"…","variant_id":"…（可选）","reason":"不超过40字的推荐理由"}]` + "\n" +
		"template_id 必须逐字来自候选清单；reason 要结合文稿主题与模板风格说具体的话，不要套话；" +
		"variant_id 仅在该模板有多个变体且你确有倾向时给。")

	var b strings.Builder
	fmt.Fprintf(&b, "## 文稿\n标题：%s\n\n", df.Title)
	b.WriteString("## 已确认大纲（逐页）\n")
	for _, pg := range o.Pages {
		fmt.Fprintf(&b, "- 第 %d 页 [%s] %s\n", pg.No, pg.Role, pg.Title)
	}
	if s := strings.TrimSpace(clarify); s != "" {
		fmt.Fprintf(&b, "\n## 用户诉求（澄清对话里的原话）\n%s\n", s)
	}
	b.WriteString("\n## 候选模板\n")
	for _, m := range candidates {
		var variants string
		if len(m.Variants) > 0 {
			ids := make([]string, 0, len(m.Variants))
			for _, v := range m.Variants {
				ids = append(ids, v.ID+":"+v.Name)
			}
			variants = " · 变体 " + strings.Join(ids, ",")
		}
		fmt.Fprintf(&b, "- %s · %s · %s", m.ID, m.Name, m.Description)
		if len(m.Tags) > 0 {
			fmt.Fprintf(&b, " · 标签[%s]", strings.Join(m.Tags, "/"))
		}
		if len(m.Scenario) > 0 {
			fmt.Fprintf(&b, " · 场景[%s]", strings.Join(m.Scenario, "/"))
		}
		b.WriteString(variants)
		b.WriteByte('\n')
	}
	return sys.String(), b.String()
}

// callLLM 非流式一轮补全（llmhello 先例的形状）。MaxTokens 给足——推理模型
// 可能把预算花在推理上，正文被截成空串（usertpl 侧实测过的坑），宁多勿少。
func callLLM(ctx context.Context, llm LLM, sys, user string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, suggestTimeout)
	defer cancel()
	completion, err := llm.Client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:          openai.ChatModel(llm.Model),
		ReasoningEffort: shared.ReasoningEffortLow,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(sys),
			openai.UserMessage(user),
		},
		MaxTokens: openai.Int(4000),
	})
	if err != nil {
		return "", err
	}
	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("响应没有 choices")
	}
	return completion.Choices[0].Message.Content, nil
}

// ---------- 模型输出的解析与白名单过滤 ----------

type rawSuggestion struct {
	TemplateID string `json:"template_id"`
	VariantID  string `json:"variant_id"`
	Reason     string `json:"reason"`
}

// sanitize 清洗模型输出：剥 markdown 围栏与隐形字符 → 反序列化 → 白名单校验
// → 理由截断 → 去重 → 截条数。template_id 非法的条目整条丢弃；
// variant_id 非法只丢变体（模板本身仍是有效推荐）。
func sanitize(text string, candidates []*template.Meta) []deck.TplSuggestion {
	byID := make(map[string]*template.Meta, len(candidates))
	for _, m := range candidates {
		byID[m.ID] = m
	}

	var raws []rawSuggestion
	clean := stripFences(text)
	if err := json.Unmarshal([]byte(clean), &raws); err != nil {
		// 兜底：模型把数组包进对象时取单键的值（decodeToolArgs 的同款思路）
		var wrapped map[string]json.RawMessage
		if err2 := json.Unmarshal([]byte(clean), &wrapped); err2 == nil && len(wrapped) == 1 {
			for _, v := range wrapped {
				if err3 := json.Unmarshal(v, &raws); err3 != nil {
					raws = nil
				}
			}
		}
		if len(raws) == 0 {
			log.Printf("[tplsuggest] 解析模型输出失败 err:%v", err)
			return nil
		}
	}

	seen := make(map[string]bool, len(raws))
	out := make([]deck.TplSuggestion, 0, len(raws))
	for _, r := range raws {
		meta, ok := byID[strings.TrimSpace(r.TemplateID)]
		if !ok {
			continue
		}
		if seen[meta.ID] {
			continue
		}
		reason := truncateRunes(strings.TrimSpace(r.Reason), maxReasonRunes)
		if reason == "" {
			continue // 没有理由的推荐前端没法展示，当废条目
		}
		variant := strings.TrimSpace(r.VariantID)
		if variant != "" {
			okVariant := variant == "default"
			for _, v := range meta.Variants {
				if v.ID == variant {
					okVariant = true
					break
				}
			}
			if !okVariant {
				variant = ""
			}
		}
		seen[meta.ID] = true
		out = append(out, deck.TplSuggestion{TemplateID: meta.ID, VariantID: variant, Reason: reason})
		if len(out) >= maxSuggestions {
			break
		}
	}
	return out
}

// stripFences 剥掉模型偶尔裹上的 markdown 代码围栏，并清掉 JSON 解析必炸的
// 隐形字符（BOM/零宽空格；agent 包 sanitizeModelJSON 的同款思路，本地小实现）。
func stripFences(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	s = strings.TrimSpace(s)
	replacer := strings.NewReplacer("\uFEFF", "", "\u200B", "", "\u200C", "", "\u200D", "", "\u2060", "")
	return replacer.Replace(s)
}

func truncateRunes(s string, max int) string {
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	runes := []rune(s)
	return string(runes[:max])
}
