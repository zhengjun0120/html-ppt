package deck

// 自定义样式槽的回归测试。两类东西要守住：
//   1. 清洗规则——挡的是"通道"（@import/url()/</style），不是风格；
//   2. 级联不变量——槽必须排在 #deck-theme-override 之后，否则自定义样式会被主题盖掉，
//      而这种 bug 在页面上表现为"我写了 CSS 但没生效"，最难排查。
// 另外补充 round-trip：CSS 文本经 raw-text 元素写回后必须逐字节一致
// （theme.go 当年就在 SetText 的实体转义上栽过，这里同一个坑）。

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

// testContract 模拟"被框架 CSS 消费过"的变量集合（真实契约由 loadVariableContract 现扫）。
var testContract = map[string]bool{
	"--accent": true, "--border": true, "--card-bg": true, "--text-muted": true,
	"--radius": true, "--space-md": true, "--space-lg": true,
	"--r-main-color": true, "--r-heading-color": true, "--r-background-color": true,
}

func TestValidateCustomCSS(t *testing.T) {
	type rejectCase struct {
		name     string
		css      string
		wantWord string
		// 只有"重定义契约变量"这一类的错误信息才需要把模型引到 update_theme 上；
		// 通道类违规（@import/url 等）应该讲它自己的正确做法，不该被塞进无关指引
		wantUpdateTheme bool
	}
	reject := []rejectCase{
		{name: "@import", css: `@import url("/x.css"); .card{}`, wantWord: "@import"},
		{name: "@IMPORT 大写", css: `@IMPORT "/x.css";`, wantWord: "@import"},
		{name: "url()", css: `.card{background:url(https://evil.com/x.png)}`, wantWord: "url()"},
		{name: "URL( 大写", css: `.card{background:URL(/a.png)}`, wantWord: "url()"},
		{name: "结尾标签", css: `/* </style><script>alert(1)</script> */`, wantWord: "</style"},
		{name: "expression", css: `.card{width:expression(alert(1))}`, wantWord: "expression("},
		// 契约变量不许在这里重定义：槽排在主题块之后、同级选择器靠后者取胜，
		// 写在这里会静默盖住 update_theme（用户换配色时页面不动，最难查的一类假成功）
		{name: "重定义 --accent", css: `:root{--accent:#ff8800}`, wantWord: "--accent", wantUpdateTheme: true},
		{name: "重定义 --radius", css: `:root{ --radius : 14px }`, wantWord: "--radius", wantUpdateTheme: true},
		{name: "重定义 reveal 变量", css: `:root{--r-main-color:#fff}`, wantWord: "--r-main-color", wantUpdateTheme: true},
		{name: "重定义 --space-md", css: `.reveal{--space-md:1em}`, wantWord: "--space-md", wantUpdateTheme: true},
	}
	for _, c := range reject {
		t.Run("拒绝_"+c.name, func(t *testing.T) {
			_, err := validateCustomCSS(c.css, testContract, nil)
			if err == nil {
				t.Fatalf("应被拒绝: %s", c.css)
			}
			if !strings.Contains(err.Error(), c.wantWord) {
				t.Errorf("错误信息应点明原因 %q，实际: %s", c.wantWord, err.Error())
			}
			// 重定义类违规必须把模型引到正确做法上，否则它会反复重试同一个写法
			if c.wantUpdateTheme && !strings.Contains(err.Error(), "update_theme") {
				t.Errorf("错误信息应指向 update_theme，实际: %s", err.Error())
			}
		})
	}

	accept := []struct{ name, css string }{
		{"正常覆盖", `.reveal .card{border-color:var(--accent);}`},
		{"引用主题变量", `.reveal .card{border-color:var(--accent);background:var(--card-bg)}`},
		// 自定义调色板走 update_theme 的 vars；在槽里定义**新**变量是合法的
		// （槽自己的样式消费它，不碰主题契约）
		{"定义新变量", `:root{--brand-ink:#ff8800;--brand-soft:rgba(255,136,0,.2)}`},
		{"注释里提到契约变量", `/* 别在这里写 --accent: 值，那是主题的地盘 */ .card{border:0}`},
		{"渐变与阴影", `.hero{background:linear-gradient(135deg,#123,#456);box-shadow:0 8px 30px rgba(0,0,0,.4);}`},
		{"带引号与实体", `.quote::before{content:"> 引用 <";} .a[data-x="1"]{color:#fff}`},
		{"空串=清空", `   `},
	}
	for _, c := range accept {
		t.Run("接受_"+c.name, func(t *testing.T) {
			out, err := validateCustomCSS(c.css, testContract, nil)
			if err != nil {
				t.Fatalf("不该拒绝: %v", err)
			}
			if strings.TrimSpace(c.css) == "" && out != "" {
				t.Fatalf("空白应归一化为空串（表示清空），实际 %q", out)
			}
			if strings.TrimSpace(c.css) != "" && strings.TrimSpace(out) != strings.TrimSpace(c.css) {
				t.Fatalf("内容被改动: in=%q out=%q", c.css, out)
			}
		})
	}
}

func TestValidateCustomCSSSizeCap(t *testing.T) {
	big := ".x{color:#fff}\n" + strings.Repeat("/* pad */\n", maxCustomCSSBytes/10)
	if _, err := validateCustomCSS(big, testContract, nil); err == nil {
		t.Fatal("超过容量上限应被拒绝")
	}
}

// TestCustomCSSBlockPlacement 级联不变量：override 永远在 custom 之前。
// 两条路径都要守住——先建槽再建主题块、先建主题块再建槽。
func TestCustomCSSBlockPlacement(t *testing.T) {
	orderOK := func(t *testing.T, html string) {
		t.Helper()
		iOverride := strings.Index(html, `id="deck-theme-override"`)
		iCustom := strings.Index(html, `id="deck-custom"`)
		if iOverride < 0 || iCustom < 0 {
			t.Fatalf("两个块都应存在: override@%d custom@%d\n%s", iOverride, iCustom, html)
		}
		if iOverride > iCustom {
			t.Fatalf("级联顺序错了：override 必须在 custom 之前\n%s", html)
		}
	}

	t.Run("先建槽再建主题块", func(t *testing.T) {
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(miniSkeleton))
		ensureCustomCSSBlock(doc)
		if _, err := ensureThemeBlocks(doc); err != nil {
			t.Fatal(err)
		}
		out, _ := goquery.OuterHtml(doc.Selection)
		orderOK(t, out)
	})

	t.Run("先建主题块再建槽", func(t *testing.T) {
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(miniSkeleton))
		if _, err := ensureThemeBlocks(doc); err != nil {
			t.Fatal(err)
		}
		ensureCustomCSSBlock(doc)
		out, _ := goquery.OuterHtml(doc.Selection)
		orderOK(t, out)
	})

	t.Run("重复调用幂等", func(t *testing.T) {
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(miniSkeleton))
		ensureCustomCSSBlock(doc)
		ensureCustomCSSBlock(doc)
		ensureCustomCSSBlock(doc)
		if n := doc.Find("#deck-custom").Length(); n != 1 {
			t.Fatalf("槽应只有一个，实际 %d", n)
		}
	})
}

// TestCustomCSSRoundTrip 写入 → 读回必须逐字节一致（含引号、实体、中文注释）。
func TestCustomCSSRoundTrip(t *testing.T) {
	s := newHistoryTestService(t)
	if err := writeDeckHTML(s, "deck-9101", miniSkeleton); err != nil {
		t.Fatal(err)
	}

	css := ".reveal .card{border-color:var(--accent)}\n" +
		".quote::before{content:\"> 引用 <\"}\n" +
		"/* 中文注释与 & 符号 */\n" +
		":root{--brand-ink:#ff8800}"

	got, err := s.updateCustomCSSLocked("deck-9101", css)
	if err != nil {
		t.Fatal(err)
	}
	if got != css {
		t.Fatalf("写入内容被改动:\n in=%q\nout=%q", css, got)
	}

	back, err := s.readCustomCSSRaw("deck-9101")
	if err != nil {
		t.Fatal(err)
	}
	if back != css {
		t.Fatalf("读回不一致（raw-text 转义问题？）:\n in=%q\nout=%q", css, back)
	}

	// 空串 = 清空：块还在，内容为空
	if _, err := s.updateCustomCSSLocked("deck-9101", "   "); err != nil {
		t.Fatal(err)
	}
	back, _ = s.readCustomCSSRaw("deck-9101")
	if back != "" {
		t.Fatalf("清空后应读回空串，实际 %q", back)
	}
	raw, _ := s.readRaw("deck-9101")
	if !strings.Contains(raw, `id="deck-custom"`) {
		t.Fatal("清空不应把槽本身删掉")
	}
}

// TestCustomCSSOldDeckMigration 老 deck（骨架里没有槽）在首次写入时自动补块，
// 且不破坏已有的主题块——这是"升级动作放在读写路径上"的核心保证。
func TestCustomCSSOldDeckMigration(t *testing.T) {
	s := newHistoryTestService(t)
	// 老骨架：只有 reveal 四件套，没有 override、没有 custom
	oldSkeleton := `<!DOCTYPE html><html><head><title>旧</title>
<link rel="stylesheet" href="/assets/components.css"></head>
<body><div class="reveal"><div class="slides">
<section data-id="s1"><h2>老页</h2></section>
</div></div></body></html>`
	if err := writeDeckHTML(s, "deck-9102", oldSkeleton); err != nil {
		t.Fatal(err)
	}

	if _, err := s.updateCustomCSSLocked("deck-9102", ".card{border-radius:2px}"); err != nil {
		t.Fatal(err)
	}
	raw, _ := s.readRaw("deck-9102")
	if !strings.Contains(raw, `id="deck-custom"`) {
		t.Fatal("老 deck 应被补上样式槽")
	}
	// 补块不能碰页面内容
	if !strings.Contains(raw, `data-id="s1"`) {
		t.Fatal("补块不应影响页面内容")
	}
	got, _ := s.readCustomCSSRaw("deck-9102")
	if got != ".card{border-radius:2px}" {
		t.Fatalf("读回不对: %q", got)
	}
}
