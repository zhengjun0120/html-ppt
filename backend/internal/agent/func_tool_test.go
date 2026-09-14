package agent

// 工具挂载的测试。重点不是"工具能不能跑"（那要靠联调），而是
// "features 开关关掉时，工具必须彻底消失"——挂着一个注定被拒的工具，
// 只会让模型反复尝试、浪费轮次。

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"html-ppt/backend/internal/service/deck"
)

func TestBuildToolsFeatureGating(t *testing.T) {
	svc := &AgentService{}

	// features.custom_css 关闭（零值）
	off := svc.buildTools()
	for _, name := range []string{"read_custom_css", "update_custom_css", "read_component"} {
		if _, ok := off[name]; ok {
			t.Errorf("custom_css 关闭时不该挂载 %s", name)
		}
	}
	// 核心工具不受 features 影响
	for _, name := range []string{"write_deck", "list_slides", "read_slide", "update_slide", "insert_slide", "delete_slide", "update_theme", "read_theme", "ask_user", "list_history", "read_history_diff"} {
		if _, ok := off[name]; !ok {
			t.Errorf("核心工具 %s 应始终挂载", name)
		}
	}
	// read_theme 必须常驻：update_theme 的 vars 是整体替换制，
	// 没有读工具就只能靠猜，会拿想象的旧配色覆盖真实配色
	if _, ok := off["read_theme"]; !ok {
		t.Fatal("read_theme 是 vars 整体替换制的前置读工具，必须始终挂载")
	}

	// 打开开关：三个样式相关工具都要出现
	svc.CustomCSS = true
	on := svc.buildTools()
	for _, name := range []string{"read_custom_css", "update_custom_css", "read_component"} {
		if _, ok := on[name]; !ok {
			t.Errorf("custom_css 打开后应挂载 %s", name)
		}
	}
	if len(on) != len(off)+3 {
		t.Errorf("打开开关应只多 3 个工具：off=%d on=%d", len(off), len(on))
	}

	// 每个工具都必须有定义和执行函数（漏一个就是运行时才暴露的坑）
	for name, tool := range on {
		if tool.Definition.OfFunction == nil {
			t.Errorf("%s 缺少 Definition", name)
		}
		if tool.Execute == nil {
			t.Errorf("%s 缺少 Execute", name)
		}
	}
}

// TestSchemaDescriptionsHaveNoASCIIComma 守一个很容易踩、而且完全不报错的坑：
// jsonschema 标签解析器按半角逗号切分键值对，**不认引号也不认转义**——
// description 里出现一个半角逗号，从那里往后整段描述会被静默丢弃。
// 实测：vars 的描述里放了个 JSON 示例（{"--surface":"#1e1836","--positive":...}），
// 模型拿到的描述就断在 `如 {"--surface":"#1e1836"` —— 恰好把最关键的用法说明吃掉。
// 这类"给模型的文档悄悄少了一半"的问题联调时几乎发现不了，只能靠断言守住。
//
// 所以：工具参数的 description 一律用全角逗号（，）、顿号（、）或分号（；）。
// 半角逗号只用于分隔标签里的键值对，且 description 必须是最后一项。
func TestSchemaDescriptionsHaveNoASCIIComma(t *testing.T) {
	argTypes := []any{
		WriteDeckArgs{}, ListSlidesArgs{}, ReadSlideArgs{}, UpdateSlideArgs{},
		InsertSlideArgs{}, DeleteSlideArgs{}, UpdateThemeArgs{}, ReadThemeArgs{},
		AskUserArgs{}, ReadHistoryDiffArgs{}, ListHistoryArgs{},
		ReadCustomCSSArgs{}, UpdateCustomCSSArgs{}, ReadComponentArgs{},
	}
	for _, v := range argTypes {
		typ := reflect.TypeOf(v)
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			tag, ok := f.Tag.Lookup("jsonschema")
			if !ok {
				continue
			}
			idx := strings.Index(tag, "description=")
			if idx < 0 {
				continue
			}
			desc := tag[idx+len("description="):]
			if c := strings.Index(desc, ","); c >= 0 {
				t.Errorf("%s.%s 的 jsonschema description 里有半角逗号（位置 %d），"+
					"解析器会在这里截断、后面的描述全部丢失：\n  %q",
					typ.Name(), f.Name, c, desc)
			}
			if strings.TrimSpace(desc) == "" {
				t.Errorf("%s.%s 的 description 是空的", typ.Name(), f.Name)
			}
		}
	}
}

// schemaEnum 从生成的工具 schema 里取某个字段的 enum 取值。
// 直接拿真实 schema 而不是手写的期望值：这个测试要守的正是"模型看到的那份 schema"。
func schemaEnum(t *testing.T, field string) []string {
	t.Helper()
	raw, err := json.Marshal(generateSchema[UpdateThemeArgs]())
	if err != nil {
		t.Fatalf("序列化 schema 失败: %v", err)
	}
	var schema map[string]any
	if err := json.Unmarshal(raw, &schema); err != nil {
		t.Fatalf("反序列化 schema 失败: %v", err)
	}
	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema 里没有 properties: %v", schema)
	}
	prop, ok := props[field].(map[string]any)
	if !ok {
		t.Fatalf("schema 里没有字段 %s", field)
	}
	list, ok := prop["enum"].([]any)
	if !ok {
		t.Fatalf("字段 %s 没有 enum（模型会以为可以随便传值）: %v", field, prop)
	}
	out := make([]string, 0, len(list))
	for _, v := range list {
		s, _ := v.(string)
		out = append(out, s)
	}
	return out
}

// TestThemeSchemaEnumsMatchCodeTables 守一条只能靠测试维持的一致性：
// 枚举取值同时存在于两处——struct tag 里的 enum（模型看到的）与 deck 包里的表
// （校验与渲染真正认的）。jsonschema 库不支持动态枚举，所以两处无法合成一处，
// 只能测。漂移的后果是"模型写得出、却被校验拒掉"，而且报错看起来毫无道理。
func TestThemeSchemaEnumsMatchCodeTables(t *testing.T) {
	for _, c := range []struct {
		field string
		want  []string
	}{
		{"preset", deck.PresetNames()},
		{"font", deck.FontNames()},
		{"texture", deck.TextureNames()},
		{"canvas", deck.CanvasNames()},
		{"transition", deck.TransitionNames()},
	} {
		got := schemaEnum(t, c.field)
		if strings.Join(got, ",") != strings.Join(c.want, ",") {
			t.Errorf("update_theme 的 %s 枚举与代码里的表不一致：\n  模型看到的=%v\n  代码实际认的=%v",
				c.field, got, c.want)
		}
	}
}

// 预设清单必须真的出现在 update_theme 的工具说明里：模型选预设靠的就是这段文字，
// 而它由 deck.PresetSummary() 生成（不在 prompt 里另抄一份，否则迟早只说一半）。
func TestPresetSummaryReachesToolDescription(t *testing.T) {
	svc := &AgentService{}
	tool, ok := svc.buildTools()["update_theme"]
	if !ok {
		t.Fatal("没有 update_theme 工具")
	}
	desc := tool.Definition.OfFunction.Function.Description.Value
	if desc == "" {
		t.Fatal("update_theme 没有 description")
	}
	for _, p := range deck.Presets() {
		if !strings.Contains(desc, p.Name) {
			t.Errorf("工具说明里缺少预设 %s：模型没有机会选它", p.Name)
		}
		if !strings.Contains(desc, p.About) {
			t.Errorf("工具说明里缺少预设 %s 的适用场合说明", p.Name)
		}
	}
}
