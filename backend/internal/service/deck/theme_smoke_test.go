package deck

// 主题流程回归测试（纯 goquery 逻辑，不依赖数据库）。
// 曾抓到过 SetText 在 raw-text 元素上把引号转义成 &#34; 的真 bug，值得留着：
// 改主题渲染/合并逻辑后跑一下 go test ./internal/service/deck/ 即可。

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

const miniSkeleton = `<!DOCTYPE html>
<html lang="zh-CN"><head>
<title>测试</title>
<link rel="stylesheet" href="/assets/theme.css">
</head><body>
<div class="reveal"><div class="slides">
<section data-id="s1"><h1>第一页</h1></section>
</div></div>
</body></html>`

func TestThemeFlow(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(miniSkeleton))
	if err != nil {
		t.Fatal(err)
	}

	// 1. 老 deck（无主题块）：自动补默认块
	theme, err := ensureThemeBlocks(doc)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Find("#deck-theme").Length() != 1 || doc.Find("#deck-theme-override").Length() != 1 {
		t.Fatalf("主题块未自动补齐")
	}
	// Theme 现在含 map 字段（自定义调色板），不能用 == 比较
	if !reflect.DeepEqual(theme, defaultTheme()) {
		t.Fatalf("老 deck 应得到默认主题: %+v", theme)
	}

	// 2. patch 合并（大小写归一）+ 全量校验 + 写回两个块
	accent := "#E8734A"
	transition := "fade"
	patch := ThemePatch{Accent: &accent, Transition: &transition}
	if err := patch.applyTo(&theme); err != nil {
		t.Fatal(err)
	}
	if err := theme.validate(); err != nil {
		t.Fatal(err)
	}
	cfg, _ := json.Marshal(theme)
	setRawText(doc.Find("#deck-theme").First(), string(cfg))
	setRawText(doc.Find("#deck-theme-override").First(), renderThemeCSS(theme))

	css := doc.Find("#deck-theme-override").First().Text()
	if !strings.Contains(css, "--accent:#e8734a;") {
		t.Fatalf("accent 未生效: %s", css)
	}
	// 面板色与边框色现在是**显式字段**（默认纸感给了明确取值：中性纸面 + 墨色发丝线），
	// 所以只改 accent 不该把它们一起换掉——派生只在字段留空时兜底。
	// 这一条同时守住"老 deck 的观感不被新默认值改写"（老 deck 没这两个字段 → 仍走派生）。
	if !strings.Contains(css, "--border:"+defaultTheme().BorderColor+";") {
		t.Fatalf("border 是显式字段，不该被 accent 改写: %s", css)
	}
	if !strings.Contains(css, "--card-bg:"+defaultTheme().Surface+";") {
		t.Fatalf("card-bg 是显式字段，不该被 accent 改写: %s", css)
	}
	// 强调色改了，跟着它走的派生量必须跟着变
	if !strings.Contains(css, "--accent-soft:rgba(232, 115, 74, 0.14);") {
		t.Fatalf("accent-soft 应随 accent 派生: %s", css)
	}

	// 3. 模拟下一次 update_theme 的读回：改过的值必须原样回来
	theme2, err := ensureThemeBlocks(doc)
	if err != nil {
		t.Fatal(err)
	}
	if theme2.Accent != "#e8734a" || theme2.Transition != "fade" {
		t.Fatalf("读回不一致: %+v", theme2)
	}
	if theme2.Background != defaultTheme().Background {
		t.Fatalf("未 patch 的字段不应变化: %+v", theme2)
	}

	// 4. 损坏的配置块：回退默认值兜底
	brokenDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(
		strings.Replace(miniSkeleton, "</head>",
			`<script type="application/json" id="deck-theme">{broken json</script>`, 1)))
	theme3, err := ensureThemeBlocks(brokenDoc)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(theme3, defaultTheme()) {
		t.Fatalf("损坏配置未回退默认: %+v", theme3)
	}

	// 5. 非法值拦截
	bad := defaultTheme()
	bad.Transition = "bounce"
	if err := bad.validate(); err == nil {
		t.Fatal("非法 transition 未被拦截")
	}
	bad2 := defaultTheme()
	bad2.Accent = "red"
	if err := bad2.validate(); err == nil {
		t.Fatal("非法颜色未被拦截")
	}

	t.Log("css sample:", css)
}
