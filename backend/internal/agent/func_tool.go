package agent

// 工具层共享件：Tool 类型、鉴权、schema 生成、序列化、历史工具、提问参数。
// （v2 阶段化工具集在 func_tool_v2.go；v1 的工具定义已随旧栈摘除。）

import (
	"bytes"
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go/v3"

	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/service/deck"
)

type ToolFunc func(ctx context.Context, arguments string) (string, error)

type Tool struct {
	Definition openai.ChatCompletionToolUnionParam //发给模型的schema
	Execute    ToolFunc                            //执行函数
	// MaxPerRun 一次 run 里允许执行的最大次数；0 = 不限。超额的调用不执行，
	// 模型会收到一条"配额用完"的结果（见 execTool）。它是"审查→修复→再审"这类
	// 无收敛循环的唯一硬闸门——提示词是软约束，实测一次 22 轮的 run 审了 5 轮、
	// 改了 5 轮，直到 maxTurns 兜底才被强制收尾。
	MaxPerRun int
}

// toolUID 从 ctx 取当前登录用户。工具是"以某个用户的身份"执行的：
// 没有身份就拒绝——这是工具层的第一道归属闸门。LLM 的参数里
// 永远没有 user 概念（身份只来自认证注入，工具签名即权限边界）。
func toolUID(ctx context.Context) (uint, error) {
	uid, ok := authctx.UserID(ctx)
	if !ok {
		return 0, errors.New("未登录状态，无法操作 deck")
	}
	return uid, nil
}

type AskQuestion struct {
	Question string   `json:"question" jsonschema:"required,type=string,description=问题本身，一句话说清要确认什么"`
	Options  []string `json:"options,omitempty" jsonschema:"type=array,description=候选项，最多4个，把你最推荐的回答放在第一个"`
}

type AskUserArgs struct {
	Questions []AskQuestion `json:"questions" jsonschema:"required,type=array,description=要问用户的问题，1~6个，尽可能一次问完不要连环调用"`
}

func (a *AgentService) toolAskUser(ctx context.Context, arguments string) (string, error) {
	return "", errors.New("ask_user 应由循环拦截暂停，不应执行到这里")
}

// ---------- 历史（v2 iterating 阶段挂载） ----------

type ReadHistoryDiffArgs struct {
	DeckID      string `json:"deck_id" jsonschema:"required,type=string,description=目标 deck 的ID，例如 deck-0002"`
	FromVersion string `json:"from_version,omitempty" jsonschema:"type=string,description=基线版本号；不传 = 最新记档版本"`
	ToVersion   string `json:"to_version,omitempty" jsonschema:"type=string,description=目标版本号；不传 = 当前使用中的内容。"`
}

func (a *AgentService) toolReadHistorydiff(ctx context.Context, arguments string) (string, error) {
	uid, err := toolUID(ctx)
	if err != nil {
		return "", err
	}
	var args ReadHistoryDiffArgs
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("read_history_diff 参数不是合法json err：%w", err)
	}
	diff, err := a.DeckService.ReadVersionDiff(uid, args.DeckID, args.FromVersion, args.ToVersion)
	if err != nil {
		return "", err
	}
	res, err := marshalNoEscape(diff)
	if err != nil {
		return "", fmt.Errorf("结果序列化json失败 err: %w", err)
	}
	return res, nil
}

type ListHistoryArgs struct {
	DeckID string `json:"deck_id" jsonschema:"required,type=string,description=目标 deck 的ID，例如 deck-0002"`
	Limit  int    `json:"limit,omitempty" jsonschema:"type=integer,description=返回条数上限，默认 15，最大 50"`
	Offset int    `json:"offset,omitempty" jsonschema:"type=integer,description=跳过最近 N 条，用于翻页查看更早的版本，默认 0"`
}

const (
	historyDefaultLimit = 15
	historyMaxLimit     = 50
)

// historyVersion 是给模型看的版本条目：在存储元信息上补一个本地时间字符串。
type historyVersion struct {
	deck.VersionMeta
	TimeStr string `json:"time_str"`
}

// historyResult 带分页元信息：工具结果会常驻对话上下文、之后每轮都要重发，
// 所以列表类工具一律默认截断。
type historyResult struct {
	DeckID        string           `json:"deck_id"`
	Total         int              `json:"total"`
	Returned      int              `json:"returned"`
	HasMore       bool             `json:"has_more"`
	OldestVersion string           `json:"oldest_version,omitempty"`
	OldestTimeStr string           `json:"oldest_time_str,omitempty"`
	Versions      []historyVersion `json:"versions"`
}

func (a *AgentService) toolListHistory(ctx context.Context, arguments string) (string, error) {
	uid, err := toolUID(ctx)
	if err != nil {
		return "", err
	}
	var args ListHistoryArgs
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("list_history 参数不是合法json err：%w", err)
	}

	limit := args.Limit
	if limit <= 0 {
		limit = historyDefaultLimit
	}
	if limit > historyMaxLimit {
		limit = historyMaxLimit
	}
	offset := args.Offset
	if offset < 0 {
		offset = 0
	}

	versions, err := a.DeckService.ListVersions(uid, args.DeckID)
	if err != nil {
		return "", err
	}
	total := len(versions)

	if offset > total {
		offset = total
	}
	end := offset + limit
	if end > total {
		end = total
	}

	out := make([]historyVersion, 0, end-offset)
	for _, v := range versions[offset:end] {
		out = append(out, historyVersion{VersionMeta: v, TimeStr: fmtUnixTime(v.Time)})
	}

	res := historyResult{
		DeckID:   args.DeckID,
		Total:    total,
		Returned: len(out),
		HasMore:  end < total,
		Versions: out,
	}
	if total > 0 {
		oldest := versions[total-1]
		res.OldestVersion = oldest.Version
		res.OldestTimeStr = fmtUnixTime(oldest.Time)
	}

	raw, err := marshalNoEscape(res)
	if err != nil {
		return "", fmt.Errorf("结果序列化json失败 err: %w", err)
	}
	return raw, nil
}

// fmtUnixTime 把 unix 秒格式化成模型好读的本地时间
func fmtUnixTime(unix int64) string {
	return time.Unix(unix, 0).Format("2006-01-02 15:04")
}

type ReviewSlidesArgs struct {
	DeckID string `json:"deck_id" jsonschema:"required,type=string,description=要审查的 deck ID"`
	// 页号留空不是错误：那是"先只给我数字"。描述里绝不能出现半角逗号：
	// jsonschema 标签解析器按半角逗号切键值对，从那里往后整段描述会被静默吃掉。
	Pages []int `json:"pages" jsonschema:"type=array,description=要看画面的页码（1 基，与 list_slides 返回的 position 一致）；例如 [2、5、7]；留空则只返回量测数字、不看图；一次最多 6 页"`
}

// 页号归一：去重、排序，并把 0 基误传挡在门口。
//
// 0 基误传是这里最可能的一种错，而且它**不报错**：传 [0,1] 时 0 只会被当成
// "没有这一页"丢掉，实际拍的是第 1 页，而报告里写着"第 1 页"——模型以为看过了。
func normalizePages(in []int) ([]int, error) {
	if len(in) == 0 {
		return nil, nil
	}
	seen := map[int]bool{}
	var out []int
	for _, p := range in {
		if p <= 0 {
			return nil, fmt.Errorf("页号从 1 开始：pages 你传的是 %v（含 %d）。要看前两页是 pages:[1,2]", in, p)
		}
		if !seen[p] {
			seen[p] = true
			out = append(out, p)
		}
	}
	sort.Ints(out)
	return out, nil
}

// 启动时调用一次并缓存即可，不要放在请求路径上反射。
func generateSchema[T any]() openai.FunctionParameters {
	reflector := jsonschema.Reflector{
		RequiredFromJSONSchemaTags: true,  // 必填由 jsonschema:"required" 标签决定（默认规则是"无 omitempty 即必填"，太隐晦）
		AllowAdditionalProperties:  false, // 生成 additionalProperties:false，配合 strict 拒绝模型幻觉出的字段
		DoNotReference:             true,  // 嵌套 struct 内联展开，而不是拆到 $defs 再 $ref（模型对 $ref 支持参差）
	}
	s := reflector.Reflect(new(T))

	schema := map[string]any{
		"type":                 "object",
		"properties":           s.Properties,
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

func marshalNoEscape(v any) (string, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(v); err != nil {
		return "", err
	}

	return strings.TrimSuffix(buf.String(), "\n"), nil
}

// joinWarnings 多条警告拼一段（分号分隔）。
func joinWarnings(parts ...string) string {
	var out []string
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			out = append(out, p)
		}
	}
	return strings.Join(out, "\n")
}
