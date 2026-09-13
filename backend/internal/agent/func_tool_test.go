package agent

// 工具挂载的测试。重点不是"工具能不能跑"（那要靠联调），而是
// "features 开关关掉时，工具必须彻底消失"——挂着一个注定被拒的工具，
// 只会让模型反复尝试、浪费轮次。

import (
	"reflect"
	"strings"
	"testing"
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
