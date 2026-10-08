package agent

// v2 工具 schema 的守门测试。
//
// 为什么值得单独守：generateSchema 依赖 jsonschema 标签，而那个解析器按**半角逗号**
// 切键值对——描述里混进一个半角逗号，后半段描述就被静默吃掉（v1 在 update_theme 上
// 踩过，模型看到的是半截说明）。这批测试把每个阶段每个工具的 schema 完整生成一遍，
// 断言描述完整、必填字段在位——标签写坏的当次提交立刻红。

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"html-ppt/backend/internal/service/deck"
)

func TestV2ToolsetsGenerateValidSchemas(t *testing.T) {
	as := &AgentService{Exec: map[string]ToolFunc{}}
	stages := []string{
		"clarifying",
		deck.StageOutlining,
		deck.StageOutlineReview,
		deck.StageGenerating,
		deck.StageIterating,
	}
	for _, stage := range stages {
		tools := as.buildToolsV2(stage)
		if tools == nil {
			t.Fatalf("阶段 %s 的工具集为 nil", stage)
		}
		for name, tool := range tools {
			raw, err := json.Marshal(tool.Definition)
			if err != nil {
				t.Fatalf("阶段 %s 工具 %s schema 序列化失败: %v", stage, name, err)
			}
			var def struct {
				Function struct {
					Name        string          `json:"name"`
					Description string          `json:"description"`
					Parameters  json.RawMessage `json:"parameters"`
				} `json:"function"`
			}
			if err := json.Unmarshal(raw, &def); err != nil {
				t.Fatalf("工具 %s 定义解析失败: %v", name, err)
			}
			if def.Function.Name != name {
				t.Errorf("工具名不一致: %q vs %q", def.Function.Name, name)
			}
			if def.Function.Description == "" {
				t.Errorf("工具 %s 缺描述", name)
			}
			// 半角逗号吞描述的回归位：描述必须以句号/问号/引号等收尾字符之外的
			// 正常文本结束，且不含明显的截断特征（以逗号或冒号结尾）。
			desc := def.Function.Description
			if strings.HasSuffix(desc, "，") || strings.HasSuffix(desc, "：") || strings.HasSuffix(desc, ",") || strings.HasSuffix(desc, ":") {
				t.Errorf("工具 %s 的描述疑似被截断（以 %q 结尾）: %q", name, desc[len(desc)-1:], desc)
			}
			// parameters 必须是合法 JSON object
			var params map[string]any
			if err := json.Unmarshal(def.Function.Parameters, &params); err != nil {
				t.Fatalf("工具 %s parameters 不是合法 JSON: %v", name, err)
			}
			if params["type"] != "object" {
				t.Errorf("工具 %s parameters.type != object", name)
			}
		}
	}
}

// 关键工具必须在正确的阶段出现/缺席（阶段边界的抽样守卫）。
func TestV2StageToolBoundaries(t *testing.T) {
	as := &AgentService{Exec: map[string]ToolFunc{}}

	has := func(stage, tool string) bool {
		_, ok := as.buildToolsV2(stage)[tool]
		return ok
	}
	if !has(deck.StageOutlining, "submit_outline") {
		t.Error("outlining 缺 submit_outline")
	}
	if has(deck.StageOutlining, "write_pages") {
		t.Error("outlining 不该看见 write_pages")
	}
	if !has(deck.StageGenerating, "plan_pages") || !has(deck.StageGenerating, "write_pages") || !has(deck.StageGenerating, "read_layout") {
		t.Error("generating 缺生成三件套")
	}
	if has(deck.StageGenerating, "submit_outline") || has(deck.StageGenerating, "update_outline") {
		t.Error("generating 不该看见大纲工具")
	}
	if has(deck.StageGenerating, "insert_slide") || has(deck.StageGenerating, "delete_slide") {
		t.Error("generating 不该看见结构变更工具（增删页是大纲级决定）")
	}
	if !has(deck.StageIterating, "insert_slide") || !has(deck.StageIterating, "update_slide") {
		t.Error("iterating 缺页级修改工具")
	}
	if !has("clarifying", "ask_user") {
		t.Error("clarifying 缺 ask_user")
	}
	// 冷启动首轮也给 submit_outline：信息给全时一步出稿，API 直连不用两步。
	// 写侧工具仍然不给——澄清阶段没有 deck，update_outline/write_pages 无处落笔。
	if !has("clarifying", "submit_outline") {
		t.Error("clarifying 缺 submit_outline（冷启动首轮直接提交大纲）")
	}
	if has("clarifying", "update_outline") || has("clarifying", "write_pages") {
		t.Error("clarifying 不该看见写侧工具")
	}
}

// emitV2 的载荷必须同时进 Content（遗留 JSON 字符串）与 Data（对象）——
// API 直连的消费者读 data 就不必对 content 做二次 JSON.parse；漏掉 Data
// 等于把双重包装又变回唯一的形态。
func TestEmitV2SetsDataField(t *testing.T) {
	var got StreamEvent
	ctx := withEmit(context.Background(), func(ev StreamEvent) error {
		got = ev
		return nil
	})
	as := &AgentService{}
	as.emitV2(ctx, EventTypeStage, map[string]any{"deck_id": "deck-1", "to": "generating"})
	if got.Type != EventTypeStage {
		t.Fatalf("事件类型 = %s", got.Type)
	}
	if len(got.Data) == 0 {
		t.Fatal("Data 未填充：API 直连消费者仍需二次解析 content")
	}
	var obj map[string]any
	if err := json.Unmarshal(got.Data, &obj); err != nil {
		t.Fatalf("Data 不是合法 JSON 对象: %v", err)
	}
	if obj["to"] != "generating" {
		t.Fatalf("Data 内容不对: %v", obj)
	}
	// Content 保持遗留形态（JSON 字符串），旧前端契约不破
	var legacy map[string]any
	if err := json.Unmarshal([]byte(got.Content), &legacy); err != nil || legacy["to"] != "generating" {
		t.Fatalf("Content 遗留形态被破坏: %q err=%v", got.Content, err)
	}
	// 整个事件序列化后 data 必须是内联对象而非字符串（SSE 线上形态）
	wire, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(wire), `"data":{"deck_id":"deck-1","to":"generating"}`) {
		t.Fatalf("线上形态 data 不是内联对象: %s", wire)
	}
}
