package agent

// decodeToolArgs 的回归：三种病态形状全部来自真实 trace（1789662713586-33b8）。
// 那次 run 烧掉 17 万 token 才发现这层需要兜底——这批测试是那次学费的固化。

import (
	"strings"
	"testing"
)

func TestDecodeToolArgs(t *testing.T) {
	type target struct {
		Title string   `json:"title"`
		Pages []struct {
			No   int    `json:"no"`
			Role string `json:"role"`
		} `json:"pages"`
	}

	t.Run("正常", func(t *testing.T) {
		var out target
		if err := decodeToolArgs(`{"title":"x","pages":[{"no":1,"role":"cover"}]}`, &out); err != nil {
			t.Fatal(err)
		}
		if out.Title != "x" || len(out.Pages) != 1 {
			t.Fatalf("%+v", out)
		}
	})

	t.Run("BOM与零宽字符（trace 实测的 'ï' 报错）", func(t *testing.T) {
		var out target
		// BOM 在开头、零宽空格混在数组元素之后——两种位置都要能活
		dirty := "\uFEFF{\"title\":\"x\",\"pages\":[{\"no\":1,\"role\":\"cover\"}\u200B, {\"no\":2,\"role\":\"toc\"}\u2060]}"
		if err := decodeToolArgs(dirty, &out); err != nil {
			t.Fatalf("隐形字符应被清洗: %v", err)
		}
		if out.Title != "x" || len(out.Pages) != 2 {
			t.Fatalf("%+v", out)
		}
	})

	t.Run("arguments 字符串包裹（失败后的退化形态）", func(t *testing.T) {
		var out target
		wrapped := `{"arguments": "{\"title\":\"x\",\"pages\":[{\"no\":1,\"role\":\"cover\"}]}"}`
		if err := decodeToolArgs(wrapped, &out); err != nil {
			t.Fatalf("字符串化包裹应被拆包: %v", err)
		}
		if out.Title != "x" {
			t.Fatalf("%+v", out)
		}
	})

	t.Run("arguments 对象包裹", func(t *testing.T) {
		var out target
		wrapped := `{"arguments": {"title": "测试标题", "pages": [{"no": 1, "role": "cover"}]}}`
		if err := decodeToolArgs(wrapped, &out); err != nil {
			t.Fatalf("对象包裹应被拆包: %v", err)
		}
		if out.Title != "测试标题" {
			t.Fatalf("%+v", out)
		}
	})

	t.Run("彻底的坏 JSON 报错带修复指令", func(t *testing.T) {
		var out target
		err := decodeToolArgs(`{"title": "x",,,,}`, &out)
		if err == nil {
			t.Fatal("应报错")
		}
		if !strings.Contains(err.Error(), "顶层 JSON") {
			t.Fatalf("报错应包含可执行的修复指令: %v", err)
		}
	})
}
