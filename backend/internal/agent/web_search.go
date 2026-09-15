package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html-ppt/backend/internal/trace"
	"log"
	"strings"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

const (
	webSearchDefaultMaxUses = 3 // 模型没传 max_uses 时我们请求里填的次数
	// webSearchMaxUses 只是"我们请求里写的上限"，**不是成本上限**——实测这个兼容端点
	// 不强制执行它：传 max_uses=1 时它照样起了 2 次搜索（两条 server_tool_use、
	// web_search_requests=2），传 2 时起了 4 次。这家的兼容层对不支持的参数一贯静默忽略。
	// 所以真正可控的是：主模型调了几次 web_search（我们这边一次调用 = 一个请求）、
	// 提示词里的纪律、以及事后看返回值里的 Searches。要盯成本，只有 Searches 可信。
	webSearchMaxUses        = 8
	webSearchMaxSources     = 16               //返回给主模型的来源条数
	webSearchMaxAnswerRunes = 3000             // 答案文本上限，是字符不是字节
	webSearchMaxTokens      = 4096             // 子调用的max_token
	webSearchTimeout        = 90 * time.Second // 整个搜索调用的context截止时间
)

const webSearchSystemPrompt = "你是一个只做事实核查的检索助手。用 web_search 找到能回答问题的可靠来源，" +
	"然后给出简洁、具体、带数字与时间的答案，每条关键事实后面附上它的来源链接。" +
	"搜不到就明确说'没有找到'，不要用常识补全，不要编造数字或链接。" +
	"只输出答案本身，不要写'根据搜索结果'这类开场白。"

type webSearchArgs struct {
	Query string `json:"query" jsonschema:"required,type=string,description=要联网搜索的问题，写成一句完整的话最准（带上年份、地区、主题，如 2026年9月 中国新能源汽车 月度销量）。堆关键词反而容易搜到泛泛的结果"`
	// 刻意不写成"允许你搜几次"：它实测不被执行（见 const 块说明），
	// 写成配额反而会让模型把它当预算去花——成本纪律在工具说明的正文里。
	MaxUses int `json:"max_uses" jsonschema:"type=integer,description=搜索轮数上限，默认 3，最大 8（服务端可能不严格遵守）"`
}

type webSearchSource struct {
	Title string `json:"title"` //来源的标题
	URL   string `json:"url"`   //来源链接
}

type webSearchResult struct {
	Query     string            `json:"query"`                //问题
	Answer    string            `json:"answer"`               //返回的回答
	Sources   []webSearchSource `json:"sources"`              //来源
	Searches  int               `json:"searches"`             //搜索次数
	Truncated bool              `json:"truncated,omitempty"`  //是否发生截断
	ToolError string            `json:"tool_error,omitempty"` //搜索工具自己的失败码
	// Queries 子模型**实际发出的**搜索词。
	//
	// 刻意不给模型看（json:"-"）：子搜索词回灌进上下文只是白烧 token，还会诱导它
	// 去复述"我搜了什么"而不是给结论。但这是观测里最有价值的一项——"模型为什么写出
	// 这个数字"往往要回到它到底搜了什么，而"搜错方向"正是这类工具最典型的失败。
	Queries []string `json:"-"`
}

// webSearchQueryFrom 从 server_tool_use 的 input 里取搜索词。
//
// input 在 SDK 里是 any（它不做具体化），所以必须容忍"根本不是 map"的情况：
// 协议一变，硬断言就 panic，而解析不该成为把服务搞崩的原因。
func webSearchQueryFrom(input any) string {
	m, ok := input.(map[string]any)
	if !ok {
		return ""
	}
	q, _ := m["query"].(string)
	return strings.TrimSpace(q)
}

// anthropicBase 返回 Anthropic 兼容入口基址。
func (a *AgentService) anthropicBase() string {
	if base := strings.TrimSpace(a.AnthropicBaseURL); base != "" {
		return strings.TrimRight(base, "/")
	}
	base := strings.TrimSpace(a.BaseURL)
	// base_url 写成 .../v1 时别拼成 .../v1/anthropic
	base = strings.TrimSuffix(strings.TrimRight(base, "/"), "/v1")
	return base + "/anthropic"
}

func parseWebSearchMsg(query string, msg *anthropic.BetaMessage) *webSearchResult {
	res := &webSearchResult{
		Query:    query,
		Searches: int(msg.Usage.ServerToolUse.WebSearchRequests),
	}
	if msg.StopReason == anthropic.BetaStopReasonMaxTokens {
		res.Truncated = true
	}

	seen := map[string]bool{}
	var answer strings.Builder
	for _, b := range msg.Content {
		switch v := b.AsAny().(type) {
		case anthropic.BetaTextBlock:
			// 模型回答的文本
			if t := strings.TrimSpace(v.Text); t != "" {
				if answer.Len() > 0 {
					answer.WriteString("\n\n")
				}
				answer.WriteString(t)
			}

		case anthropic.BetaServerToolUseBlock:
			// 子模型实际发出的搜索词。这一段过去是个空 case，注释还写着"后续可以传回给
			// 前端显示"——于是联网搜索在观测里只剩一个 searches=2 的计数：知道它搜了两次，
			// 但不知道搜的是什么。
			if q := webSearchQueryFrom(v.Input); q != "" {
				res.Queries = append(res.Queries, q)
			}
		case anthropic.BetaWebSearchToolResultBlock:
			// 模型工具结果
			if code := v.Content.ErrorCode; code != "" {
				res.ToolError = string(code)
				continue
			}
			for _, it := range v.Content.OfBetaWebSearchResultBlockArray {
				if !it.JSON.URL.Valid() || it.URL == "" || seen[it.URL] {
					continue
				}
				if !it.JSON.Title.Valid() {
					// 有链接没标题 可能是协议变了
					log.Printf("[warn] 搜索结果缺 title 字段（协议可能变了）: %s", it.URL)
				}

				seen[it.URL] = true
				res.Sources = append(res.Sources, webSearchSource{Title: it.Title, URL: it.URL})
			}
		}
	}

	if len(res.Sources) > webSearchMaxSources {
		res.Sources = res.Sources[:webSearchMaxSources]
		res.Truncated = true
	}

	text := []rune(strings.TrimSpace(answer.String()))
	if len(text) > webSearchMaxAnswerRunes {
		text = text[:webSearchMaxAnswerRunes]
		res.Truncated = true
	}
	res.Answer = string(text)
	return res
}

// clampMaxUses 把模型要的搜索次数收敛到 [默认值, 请求上限]。
// 抽出来是为了能离线测：这段边界不该靠"真发一次搜索"来验。
// 注意它收敛的只是**我们请求里写的那个值**——实测该端点不执行 max_uses，
// 所以这是个"表达意图"的参数，不是成本闸门（详见上面 const 块的说明）。
func clampMaxUses(maxUses int) int {
	if maxUses <= 0 {
		return webSearchDefaultMaxUses
	}
	if maxUses > webSearchMaxUses {
		return webSearchMaxUses
	}
	return maxUses
}

// runWebSearch 发一次真正的服务端搜索 返回的 Searches 是计费口径的次数
func (a *AgentService) runWebSearch(ctx context.Context, apiKey, query string, maxUses int) (*webSearchResult, error) {
	maxUses = clampMaxUses(maxUses)

	// 观测：请求前先记下"我们让它搜什么"。这一步是必要的，因为子模型会自己改写搜索词
	//（实测同一个问题两次，它发出的 query 分别是中文和英文两版），
	// 所以这条与下面从响应里解析出的 queries 不是一回事——两相对照才看得出它有没有跑偏。
	trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Sub: &trace.SubStep{
		Name: "web_search", Stage: "request", Text: query,
		Data: map[string]any{
			"max_uses": maxUses, "model": a.ModelID,
			"timeout_s": int(webSearchTimeout / time.Second),
		},
	}})

	callCtx, cancel := context.WithTimeout(ctx, webSearchTimeout)
	defer cancel()

	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(a.anthropicBase()),
		//重试
		option.WithMaxRetries(2),
	)

	msg, err := client.Beta.Messages.New(callCtx, anthropic.BetaMessageNewParams{
		Model:     a.ModelID,
		MaxTokens: webSearchMaxTokens,
		System:    []anthropic.BetaTextBlockParam{{Text: webSearchSystemPrompt}},
		Messages: []anthropic.BetaMessageParam{
			anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock(query)),
		},

		Tools: []anthropic.BetaToolUnionParam{
			{
				OfWebSearchTool20250305: &anthropic.BetaWebSearchTool20250305Param{
					MaxUses: anthropic.Int(int64(maxUses)),
				},
			},
		},
	})

	if err != nil {
		// 观测：失败要留痕。HTTP 码和原始响应体是判断"这次钱花没花、下次该不该重试"
		// 的唯一依据，而它们过去只出现在服务端日志里，翻不到就没了
		trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Sub: &trace.SubStep{
			Name: "web_search", Stage: "error", Text: err.Error(),
		}})
		var apiErr *anthropic.Error
		if errors.As(err, &apiErr) {
			return nil, fmt.Errorf("联网搜索被拒绝(HTTP %d) : %s", apiErr.StatusCode, apiErr.RawJSON())
		}
		return nil, fmt.Errorf("联网搜索请求失败(网络不通或超时) ： %w", err)
	}

	res := parseWebSearchMsg(query, msg)

	// 观测：子调用的用量。这是**分开计费**的另一次 API 调用，过去完全没进过统计——
	// 于是"这次对话花了多少 token"这个数系统性偏小。图片 token 这里是 0，
	// 但 input/output 必须算上。
	trace.Usage(ctx, trace.CompWebSearch, trace.UsagePart{
		Prompt:     msg.Usage.InputTokens,
		Completion: msg.Usage.OutputTokens,
		Total:      msg.Usage.InputTokens + msg.Usage.OutputTokens,
		Cached:     msg.Usage.CacheReadInputTokens,
	})

	// 观测：子模型到底搜了几轮、搜的是什么、为什么停。这张卡片就是联网搜索的"内部过程"
	trace.Emit(ctx, trace.Event{Kind: trace.KindSubStep, Sub: &trace.SubStep{
		Name: "web_search", Stage: "queries",
		Data: map[string]any{
			"queries":     res.Queries,
			"searches":    res.Searches,
			"stop_reason": string(msg.StopReason),
			"tool_error":  res.ToolError,
			"truncated":   res.Truncated,
			"source_n":    len(res.Sources),
		},
	}})

	if res.Searches > 0 && len(res.Sources) == 0 && res.ToolError == "" {
		log.Printf("[warn] 搜索执行了 %d 次但没解析出任何来源，协议字段可能变了", res.Searches)
	}

	return res, nil
}

func (a *AgentService) toolWebSearch(ctx context.Context, arguments string) (string, error) {
	var args webSearchArgs
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("web_search 参数json不合法 err:%w", err)
	}

	query := strings.TrimSpace(args.Query)
	if query == "" {
		return "", errors.New("query 不能为空")
	}

	key, _ := a.resolveKey(ctx)
	if strings.TrimSpace(key) == "" {
		return "", errors.New("没有可用的 API key，无法联网搜索")
	}

	res, err := a.runWebSearch(ctx, key, query, args.MaxUses)
	if err != nil {
		// 这条错误是**给模型看的**（见 agent.go 的 execTool）：失败原因决定它下一步该换
		// 问法还是干脆放弃联网。所以要带上 %w——漏掉 err 的那版只剩一句
		// %!w(MISSING)，模型无从判断，go vet 也会直接报出来。
		return "", fmt.Errorf("搜索失败 err：%w；不要重试搜索，改用你确定的知识完成这一页，不要写具体数字，并在回复里说明未能联网核实", err)
	}

	switch {
	case res.ToolError != "" && res.Answer == "" && len(res.Sources) == 0:
		return "", fmt.Errorf("联网搜索工具报错(%s)且没有返回任何内容；改用你确定的知识，不要写具体数字", res.ToolError)
	case res.Searches == 0 && res.Answer == "" && len(res.Sources) == 0:
		return "", fmt.Errorf("模型没有发起搜索(可能它认为无需联网)。把问题写得更具体，或换一种问法")
	case res.Answer == "" && len(res.Sources) == 0:
		return "", fmt.Errorf("联网搜索没有返回可用结果，请改用其他方式获取这条信息")
	}

	out, err := json.Marshal(res)
	if err != nil {
		return "", fmt.Errorf("序列化搜索结果失败 err:%w", err)
	}
	return string(out), nil
}
