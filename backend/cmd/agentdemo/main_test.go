package main

import (
	"encoding/json"
	"testing"
)

// TestGenerateSchema 离线验证 struct 反射出的工具 schema，
// 不需要 API key：go test ./cmd/agentdemo -v 即可运行。
// 写法上的 tag 手误（比如 required 拼错、enum 少写）都会在这里暴露。
func TestGenerateSchema(t *testing.T) {
	fp := generateSchema[CreatePresentationArgs]()

	raw, err := json.Marshal(fp)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("生成的 schema: %s", raw)

	var s struct {
		Properties map[string]any `json:"properties"`
		Required   []string       `json:"required"`
	}
	if err := json.Unmarshal(raw, &s); err != nil {
		t.Fatal(err)
	}

	// 必填字段齐全
	for _, want := range []string{"title", "theme", "slides"} {
		found := false
		for _, r := range s.Required {
			if r == want {
				found = true
			}
		}
		if !found {
			t.Errorf("required 里缺少 %q（检查 jsonschema:\"required\" 标签）", want)
		}
	}

	// enum 约束生效
	theme, ok := s.Properties["theme"].(map[string]any)
	if !ok {
		t.Fatal("theme 属性缺失或不是对象")
	}
	if enums, _ := theme["enum"].([]any); len(enums) != 3 {
		t.Errorf("theme 的 enum 应有 3 个取值，实际: %v", theme["enum"])
	}

	// 嵌套 struct 内联展开（DoNotReference）：slides.items 里应直接是对象 schema
	slides, ok := s.Properties["slides"].(map[string]any)
	if !ok {
		t.Fatal("slides 属性缺失")
	}
	items, ok := slides["items"].(map[string]any)
	if !ok {
		t.Fatalf("slides.items 缺失，slides=%v", slides)
	}
	if _, hasRef := items["$ref"]; hasRef {
		t.Error("slides.items 不应使用 $ref（检查 DoNotReference 是否开启）")
	}
	if itemsProps, _ := items["properties"].(map[string]any); len(itemsProps) != 2 {
		t.Errorf("Slide 应内联出 2 个属性，实际: %v", itemsProps)
	}
}
