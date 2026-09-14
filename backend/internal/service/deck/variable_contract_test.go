package deck

// 回声校验的测试。它防的是"静默无效"这一类失败：
// 变量名拼错、凭空造一个没人消费的变量——页面不报错、结构也没坏，只是那处样式不生效，
// 而 AI 会照常汇报"已改好"。所以这里既测"能扫出契约"，也测"该拒的必须拒、该放的必须放"
// （放错了会把正常写法拦住，比不查还糟）。

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const contractFixtureComponents = `/* 组件库 */
.reveal .card { background: var(--card-bg); border: 1px solid var(--border); border-radius: var(--radius); }
.reveal .grid-2 { gap: var(--space-md) var(--space-lg); }
.reveal .card p { color: var(--text-muted); }
.reveal .badge { background: var(--accent); }
`
const contractFixtureTheme = `:root { --accent: #5eead4; --brand-extra: #fff; }
.reveal { font-family: var(--r-main-font); color: var(--r-main-color); }
`
// contractFixtureReveal 模拟 reveal 自带的变量消费
const contractFixtureReveal = `.reveal-viewport{background-color:var(--r-background-color)}h1{color:var(--r-heading-color)}`

func writeContractAssets(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	files := map[string]string{
		"components.css": contractFixtureComponents,
		"theme.css":      contractFixtureTheme,
		"reveal.css":     contractFixtureReveal,
		"moon.css":       `a{color:var(--r-link-color)}code{font-family:var(--r-code-font)}`,
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestLoadVariableContract(t *testing.T) {
	known, err := loadVariableContract(writeContractAssets(t))
	if err != nil {
		t.Fatal(err)
	}
	// 被 var() 消费过的才算契约
	for _, want := range []string{"--card-bg", "--border", "--radius", "--space-md", "--accent", "--text-muted", "--r-main-font", "--r-link-color", "--r-code-font"} {
		if !known[want] {
			t.Errorf("契约应包含 %s", want)
		}
	}
	// --brand-extra 只在 theme.css 里被**定义**、没有任何元素消费它 → 不算契约
	// （这正是回声校验的意义：改它没有任何效果）
	if known["--brand-extra"] {
		t.Error("只被定义、没被消费的变量不该进契约")
	}

	if _, err := loadVariableContract(""); err == nil {
		t.Error("目录未配置应报错")
	}
	if _, err := loadVariableContract(t.TempDir()); err == nil {
		t.Error("目录里一个契约文件都没有时应报错")
	}
}

func TestUnknownVariables(t *testing.T) {
	known := map[string]bool{"--accent": true, "--radius": true}

	cases := []struct {
		name    string
		text    string
		defined map[string]bool
		want    string // 期望的未知变量（逗号分隔，空串=无）
	}{
		{"全部已知", `.x{color:var(--accent)}`, nil, ""},
		{"拼错一个", `.x{color:var(--acent)}`, nil, "--acent"},
		{"凭空造", `.x{color:var(--my-accent);border-radius:var(--radius)}`, nil, "--my-accent"},
		{"带 fallback 也算使用", `.x{color:var(--nope, red)}`, nil, "--nope"},
		{"var 里带空格", `.x{color:var(-- nopeSpaced )}`, nil, ""}, // -- 后跟空格不是合法变量名，正则不认，交给人眼
		{"同份 CSS 里定义过就豁免", `.hero{--glow:0 0 20px var(--accent);box-shadow:var(--glow)}`, nil, ""},
		{"主题块里定义过也豁免", `.x{color:var(--surface)}`, map[string]bool{"--surface": true}, ""},
		{"注释里的用法不算", `/* 别用 var(--ghost) */ .x{color:var(--accent)}`, nil, ""},
		{"多个未知按序去重", `.a{color:var(--z);background:var(--a)} .b{color:var(--z)}`, nil, "--a,--z"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := strings.Join(unknownVariables(c.text, known, c.defined), ",")
			if got != c.want {
				t.Fatalf("未知变量判定不对:\n got=%q\nwant=%q", got, c.want)
			}
		})
	}
}

func TestCheckCSSBalance(t *testing.T) {
	if err := checkCSSBalance(`.a{color:red} .b{color:blue}`); err != nil {
		t.Fatalf("配对的 CSS 不该报错: %v", err)
	}
	if err := checkCSSBalance(`.a{color:red`); err == nil {
		t.Error("漏掉 } 应报错")
	}
	if err := checkCSSBalance(`.a{color:red}}`); err == nil {
		t.Error("多余的 } 应报错")
	}
	// 注释与字符串里的花括号不该参与判定
	if err := checkCSSBalance(`/* 说明：}) 乱写 */ .a{content:"{"}`); err != nil {
		t.Errorf("注释/字符串里的花括号被误判: %v", err)
	}
}

// TestEchoCheckRejectsUnknownVar 端到端走 validateCustomCSS：该拒的拒、该放的放。
func TestEchoCheckRejectsUnknownVar(t *testing.T) {
	preDefined := map[string]bool{"--surface": true} // 模拟 update_theme 的 vars 里已定义

	reject := []struct{ name, css, wantName string }{
		{"拼错变量", `.reveal .card{background:var(--acent)}`, "--acent"},
		{"凭空造变量", `.hero{color:var(--brand-ink)}`, "--brand-ink"},
	}
	for _, c := range reject {
		t.Run("拒_"+c.name, func(t *testing.T) {
			_, err := validateCustomCSS(c.css, testContract, preDefined)
			if err == nil {
				t.Fatalf("应被拒收: %s", c.css)
			}
			// 错误里必须同时有"哪个名字无效"和"可用的有哪些"，模型才能自己改对
			if !strings.Contains(err.Error(), c.wantName) {
				t.Errorf("错误应点明无效变量 %s: %s", c.wantName, err.Error())
			}
			if !strings.Contains(err.Error(), "--accent") {
				t.Errorf("错误应附可用变量清单: %s", err.Error())
			}
		})
	}

	accept := []struct{ name, css string }{
		{"用契约变量", `.reveal .card{border-color:var(--accent)}`},
		{"用主题块定义过的变量", `.x{background:var(--surface)}`},
		{"同份 CSS 自定自用", `.hero{--glow:0 0 20px var(--accent);box-shadow:var(--glow)}`},
	}
	for _, c := range accept {
		t.Run("放_"+c.name, func(t *testing.T) {
			if _, err := validateCustomCSS(c.css, testContract, preDefined); err != nil {
				t.Fatalf("不该拒收（会把正常写法也拦住）: %v", err)
			}
		})
	}

	// 契约不可得（nil）时跳过回声校验：这项检查防的是"静默无效"而不是安全问题，
	// 不能因为资源目录没配好就把所有写入卡死
	if _, err := validateCustomCSS(`.x{color:var(--whatever)}`, nil, nil); err != nil {
		t.Fatalf("契约不可得时应跳过校验，实际: %v", err)
	}
}

// TestCheckInlineStyleVars 页面内容侧：只认 style="..." 里的用法，
// 且 deck 主题块里定义过的变量要算已知（vars 调色板 + 页面里 var() 引用是正常组合）。
func TestCheckInlineStyleVars(t *testing.T) {
	dir := writeContractAssets(t)
	s := &Service{decksDir: t.TempDir(), assetsDir: dir}

	section := `<section data-id="s1"><h2 style="color:var(--accent)">标题</h2></section>`
	if err := s.checkInlineStyleVars(section, ""); err != nil {
		t.Fatalf("已知变量不该报错: %v", err)
	}

	bad := `<section data-id="s1"><h2 style="color:var(--acent)">标题</h2></section>`
	err := s.checkInlineStyleVars(bad, "")
	if err == nil {
		t.Fatal("未知变量应报错")
	}
	if !strings.Contains(err.Error(), "--acent") {
		t.Errorf("错误应点明变量名: %s", err.Error())
	}

	// 正文里提到 var(--ghost) 不算违规（只认 style 属性里的用法）
	prose := `<section data-id="s1"><p>文档里可以写 var(--ghost) 作为说明</p></section>`
	if err := s.checkInlineStyleVars(prose, ""); err != nil {
		t.Errorf("正文里的字样不该被判违规: %v", err)
	}

	// deck 主题块（vars 调色板）里定义过的变量，页面里可以用
	deckHTML := `<html><head><style id="deck-theme-override">:root{--surface:#1e1836;}</style></head><body></body></html>`
	useSurface := `<section data-id="s1"><div style="background:var(--surface)">x</div></section>`
	if err := s.checkInlineStyleVars(useSurface, deckHTML); err != nil {
		t.Fatalf("主题块里定义过的变量应可用: %v", err)
	}
	// 没有那个主题块时，同一个变量就该被拦住
	if err := s.checkInlineStyleVars(useSurface, ""); err == nil {
		t.Fatal("主题块不存在时 --surface 应被判为未知")
	}
}

// TestRealVariableContract 跨产物契约测试：system prompt 的变量白名单里承诺的变量，
// 必须全部真的被框架 CSS 消费（否则模型照着 prompt 写，会被回声校验拒收——
// 文档、CSS、校验器三者互相打架，是最难排查的一类问题）。
func TestRealVariableContract(t *testing.T) {
	known, err := loadVariableContract(filepath.Join("..", "..", "..", "web", "assets"))
	if err != nil {
		t.Skipf("读不到真实资源目录，跳过契约检查: %v", err)
	}

	// 与 systemPrompt.md 的「可用变量白名单」一一对应
	promptWhitelist := []string{
		"--accent", "--border", "--card-bg", "--text-muted", "--radius",
		"--accent-soft", "--on-accent", "--hairline",
		"--space-2xs", "--space-xs", "--space-sm", "--space-md", "--space-lg",
		"--r-background-color", "--r-main-color", "--r-heading-color", "--r-link-color",
		"--r-main-font", "--r-heading-font", "--r-code-font",
	}
	for _, name := range promptWhitelist {
		if !known[name] {
			t.Errorf("prompt 白名单里的 %s 没有被任何框架 CSS 消费：模型按 prompt 写它会被拒收，两边打架了", name)
		}
	}
	t.Logf("真实契约共 %d 个变量", len(known))
}
