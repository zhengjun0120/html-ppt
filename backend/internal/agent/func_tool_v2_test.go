package agent

// v2 工具 schema 的守门测试。
//
// 为什么值得单独守：generateSchema 依赖 jsonschema 标签，而那个解析器按**半角逗号**
// 切键值对——描述里混进一个半角逗号，后半段描述就被静默吃掉（v1 在 update_theme 上
// 踩过，模型看到的是半截说明）。这批测试把每个阶段每个工具的 schema 完整生成一遍，
// 断言描述完整、必填字段在位——标签写坏的当次提交立刻红。

import (
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
}
