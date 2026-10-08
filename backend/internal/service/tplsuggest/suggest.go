// Package tplsuggest 选模板阶段的模板推荐：聚大纲与用户诉求喂 LLM，
// 从内置模板里挑 3-5 个，结果缓存进 deck.json（deck 层负责存取）。
//
// 依赖注入照 usertpl 的先例：LLM client/model 由装配层从 agent 服务构造
// （handler 里 h.agent.CustomizeLLMFor），本包不 import agent，避免依赖环。
package tplsuggest

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"html-ppt/backend/internal/service/deck"
	"html-ppt/backend/internal/service/template"
	"html-ppt/backend/internal/trace"

	"github.com/openai/openai-go/v3"
)

// LLM 一次补全的接入参数（usertpl.LLM 同款形状）。
type LLM struct {
	Client *openai.Client
	Model  string
}

// Service 模板推荐服务。reg 只读元数据；decks 负责守卫依据（stage）、
// 大纲读取与推荐缓存落盘；traceCfg 观测落盘（run_kind=tplsugg，2026-09-30 接入观测台）。
type Service struct {
	decks    *deck.Service
	reg      *template.Registry
	traceCfg trace.Config
	// inFlight 同 deck 的在飞去重：gate 1 确认后的后台预热与用户进模板页的
	// 请求撞在一起时，后到者等锁、拿到后读缓存返回——同一 deck 同时只有一次
	// LLM 调用（重复调用除了双倍额度没有任何收益）。
	inFlight sync.Map
}

func New(decks *deck.Service, reg *template.Registry, traceCfg trace.Config) *Service {
	return &Service{decks: decks, reg: reg, traceCfg: traceCfg}
}

// suggestTimeout 一次推荐 LLM 调用的上限。上限交还模型默认后，推理模型首轮
// 可能比原来慢；选模板页是异步加载推荐、网格照常可用，但也不能让请求干等到
// 请求方超时。
const suggestTimeout = 60 * time.Second

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

	// 缓存检查 + 计算全程持锁：后到的请求（用户进页 vs 确认后的预热）在这里
	// 排队，前一个算完写了缓存，后一个直接命中返回——同 deck 永不重复调 LLM。
	unlock := s.lockDeck(deckID)
	defer unlock()

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

	// 观测（fail-open，与 customize 同一条原则）：缓存命中不进这里——只有真调
	// LLM 的那一轮才算一次 run。上下文全量进 llm_request，"模型为什么推这个"
	// 才有证据可查。
	ctx, rec, closeTrace := s.traceRun(ctx, uid, deckID, df.Title, len(o.Pages), llm.Model, sys, user)
	defer closeTrace()
	// run_end 带 recorder 汇总（用量/耗时/轮数）——观测台列表的 tokens 列就从这来
	endRun := func(status string, errMsg string) {
		if rec == nil {
			return
		}
		s := rec.Summary()
		ev := trace.Event{Kind: trace.KindRunEnd, Status: status, Summary: &s}
		if errMsg != "" {
			ev.Error = errMsg
		}
		trace.Emit(ctx, ev)
	}

	raw, finish, err := callLLM(ctx, llm, sys, user)
	if err != nil {
		trace.Emit(ctx, trace.Event{Kind: trace.KindLLMResponse, Error: err.Error()})
		endRun(trace.StatusError, "LLM 调用失败")
		log.Printf("[tplsuggest] %s LLM 调用失败 err:%v", deckID, err)
		return []deck.TplSuggestion{}, nil
	}
	// 空回复重试一次（customize 同款兜底）：推理模型偶尔把整轮预算花在
	// reasoning 上，正文零字——再给一次机会，还空就当失败处理。
	if strings.TrimSpace(raw) == "" {
		trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Content: "首次调用空回复（推理耗尽输出），重试一次"})
		raw, finish, err = callLLM(ctx, llm, sys, user)
		if err != nil {
			trace.Emit(ctx, trace.Event{Kind: trace.KindLLMResponse, Error: err.Error()})
			endRun(trace.StatusError, "重试仍失败")
			log.Printf("[tplsuggest] %s LLM 重试失败 err:%v", deckID, err)
			return []deck.TplSuggestion{}, nil
		}
	}
	sugs := sanitize(raw, candidates)
	trace.Emit(ctx, trace.Event{Kind: trace.KindLLMResponse, FinishReason: finish, Content: raw})
	if len(sugs) == 0 {
		endRun(trace.StatusError, "无合法推荐（解析后 0 条）")
		log.Printf("[tplsuggest] %s 无合法推荐，原始输出：%.80s", deckID, raw)
		return []deck.TplSuggestion{}, nil
	}
	if err := s.decks.SaveTemplateSuggestions(uid, deckID, sugs); err != nil {
		log.Printf("[tplsuggest] %s 写缓存失败 err:%v", deckID, err) // 推荐照常返回，缓存下次重建
	}
	endRun(trace.StatusOK, "")
	return sugs, nil
}

// lockDeck 同 deck 的在飞互斥（deck.Service 的 lockDeck 同款思路，就地小实现）。
func (s *Service) lockDeck(deckID string) func() {
	v, _ := s.inFlight.LoadOrStore(deckID, &sync.Mutex{})
	mu := v.(*sync.Mutex)
	mu.Lock()
	return mu.Unlock
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
// 候选行刻意压缩（简介截 40 字、变体只留 id）——102 个内置模板的清单是输入
// token 的大头，行宽直接换算成首 token 延迟与费用；模型要的风格信号在
// 标签/场景里，长简介是冗余。
func buildPrompt(df *deck.DeckFile, o *deck.Outline, clarify string, candidates []*template.Meta) (string, string) {
	var sys strings.Builder
	sys.WriteString("你是 PPT 模板推荐助手。根据文稿大纲与用户在澄清对话里表达的诉求，从候选清单中挑 3-5 个最合适的模板。\n" +
		"这是一次轻量挑选任务：候选清单很长但判断不复杂，直接挑选并输出，不要展开冗长推理。\n" +
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
		fmt.Fprintf(&b, "- %s · %s · %s", m.ID, m.Name, truncateRunes(m.Description, 40))
		if len(m.Tags) > 0 {
			fmt.Fprintf(&b, " · 标签[%s]", strings.Join(capped(m.Tags, 6), "/"))
		}
		if len(m.Scenario) > 0 {
			fmt.Fprintf(&b, " · 场景[%s]", strings.Join(capped(m.Scenario, 4), "/"))
		}
		if len(m.Variants) > 0 {
			ids := make([]string, 0, len(m.Variants))
			for _, v := range m.Variants {
				ids = append(ids, v.ID)
			}
			fmt.Fprintf(&b, " · 变体[%s]", strings.Join(ids, ","))
		}
		b.WriteByte('\n')
	}
	return sys.String(), b.String()
}

// capped 截断字符串切片（超出丢弃——标签/场景按原序，前面的更核心）。
func capped(list []string, n int) []string {
	if len(list) <= n {
		return list
	}
	return list[:n]
}

// callLLM 非流式一轮补全（llmhello 先例的形状）。**不设 MaxTokens、不设
// ReasoningEffort**（customize 的 custStream 同款决定，2026-09-30 实测重蹈覆辙：
// 4000 上限被 deepseek-flash 的 reasoning 全部吃满、正文零字，观测台 llm_response
// 记到空 content 才定位）——上限交给模型/服务商默认，与 agent 主循环一致。
// 返回正文与 finish_reason（观测要记）。
func callLLM(ctx context.Context, llm LLM, sys, user string) (string, string, error) {
	ctx, cancel := context.WithTimeout(ctx, suggestTimeout)
	defer cancel()
	completion, err := llm.Client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModel(llm.Model),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(sys),
			openai.UserMessage(user),
		},
		Temperature: openai.Float(0.3),
	})
	if err != nil {
		return "", "", err
	}
	if len(completion.Choices) == 0 {
		return "", "", fmt.Errorf("响应没有 choices")
	}
	trace.Usage(ctx, trace.CompMain, suggestUsagePart(completion.Usage))
	return completion.Choices[0].Message.Content, string(completion.Choices[0].FinishReason), nil
}

// suggestSessionID 推荐 run 的伪会话号（customize 的 customSessionID 同款思路）：
// 模板推荐没有会话表，deck id 哈希进 1e9 起步的专用段——同一 deck 的推荐 run
// 天然聚在一个目录，按 run 数裁剪等于每 deck 各自的保留配额。
func suggestSessionID(deckID string) uint {
	h := fnv.New32a()
	_, _ = h.Write([]byte("tplsugg:" + deckID))
	return 1_000_000_000 + uint(h.Sum32())
}

// traceRun 开一个推荐 run 的观测文件并把 recorder 挂进 ctx；发 run_start 与
// llm_request（候选清单+大纲+诉求全量落盘）。返回挂好 recorder 的 ctx、recorder
// 本体（endRun 取汇总用，关闭路径失败时为 nil）与清理闭包（关 recorder）。
// 任何失败只打日志不阻断推荐——观测是增强不是故障源。
func (s *Service) traceRun(ctx context.Context, uid uint, deckID, title string, pages int, model, sys, user string) (context.Context, *trace.Recorder, func()) {
	if !s.traceCfg.Enabled || s.traceCfg.Dir == "" {
		return ctx, nil, func() {}
	}
	sessID := suggestSessionID(deckID)
	rec, err := trace.New(trace.Options{
		Dir:           s.traceCfg.Dir,
		SessionID:     sessID,
		RunID:         trace.NewRunID(time.Now()),
		UserID:        uid,
		DeckID:        deckID,
		MaxFieldBytes: s.traceCfg.MaxFieldBytes,
		RetainRuns:    s.traceCfg.RetainRuns,
	})
	if err != nil {
		log.Printf("[warn] tplsuggest: 开启观测失败，本轮不记录 err: %v", err)
		return ctx, nil, func() {}
	}
	rec.Emit(trace.Event{
		Kind: trace.KindRunStart, RunKind: "tplsugg",
		RunID: rec.RunID(), SessionID: sessID, UserID: uid, DeckID: deckID,
		UserContent: fmt.Sprintf("模板推荐：%s（%d 页）", title, pages),
		Model:       model,
	})
	msgs, err := json.Marshal([]map[string]string{
		{"role": "system", "content": sys},
		{"role": "user", "content": user},
	})
	if err != nil {
		msgs = json.RawMessage(`null`)
	}
	rec.Emit(trace.Event{
		Kind: trace.KindLLMRequest, Model: model,
		Messages: msgs, MessageCount: 2, Bytes: len(msgs),
	})
	return trace.With(ctx, rec), rec, rec.Close
}

// suggestUsagePart SDK 用量 → 观测口径（custUsagePart 同一张表，就地复制：
// 为一个换算函数跨包导出别的服务的内部件不值得）。
func suggestUsagePart(u openai.CompletionUsage) trace.UsagePart {
	return trace.UsagePart{
		Prompt:      u.PromptTokens,
		Completion:  u.CompletionTokens,
		Total:       u.TotalTokens,
		Cached:      u.PromptTokensDetails.CachedTokens,
		Reasoning:   u.CompletionTokensDetails.ReasoningTokens,
		ImageTokens: u.PromptTokensDetails.ImageTokens,
		Calls:       1,
	}
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
