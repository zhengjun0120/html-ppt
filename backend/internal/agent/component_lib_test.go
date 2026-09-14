package agent

// read_component 的解析逻辑测试。都是纯函数，不需要起服务、不需要数据库。
// 这里守的是"给模型的组件库视图必须准确"——它拿着这份视图去写覆盖样式，
// 视图错了就会写出一堆匹配不到的选择器。

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// guardRe 抠出 :where(:not(...)) 里的排除名单（解析 :where 守卫用）
var guardRe = regexp.MustCompile(`:where\(:not\(([^)]*)\)\)`)

const fixtureCSS = `/* 组件库 —— 结构层
   注释里故意放一个 } 大括号，验证切块时注释先被剥掉 */
.reveal .slides section { text-align: left; }
.reveal .slides section.center { text-align: center; }

/* 卡片网格 */
.reveal .grid-2 {
  display: grid;
  gap: var(--space-md) var(--space-lg);
}
.reveal .card {
  background: var(--card-bg);
  padding: var(--space-md) var(--space-lg);
}
.reveal .card h3 { font-size: 0.8em; }
.reveal .card-x { color: red; }          /* 干扰项：不应被 .card 命中 */
.reveal .badge { font-size: 0.65em; }
`

func TestSplitRulesHandlesComments(t *testing.T) {
	rules := splitRules(fixtureCSS)
	if len(rules) != 7 {
		t.Fatalf("应切出 7 条规则，实际 %d:\n%s", len(rules), strings.Join(rules, "\n---\n"))
	}
	for _, r := range rules {
		if !strings.HasSuffix(r, "}") {
			t.Fatalf("规则片段应以 } 结尾: %q", r)
		}
	}
}

func TestExtractClassNames(t *testing.T) {
	got := extractClassNames(fixtureCSS)
	// reveal/slides 是框架命名空间，不算可用组件类，应被过滤掉
	want := []string{"badge", "card", "card-x", "center", "grid-2"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("类名清单不对:\n got=%v\nwant=%v", got, want)
	}
	// 数值不能被误认成类名（0.8em / 0.65em）
	for _, n := range got {
		if n == "8em" || n == "65em" {
			t.Fatalf("把数值当成类名了: %v", got)
		}
	}
}

func TestExtractRulesByClass(t *testing.T) {
	out, err := extractRules(fixtureCSS, "card")
	if err != nil {
		t.Fatal(err)
	}
	// .card 及其后代规则都要给到
	for _, want := range []string{".reveal .card {", ".reveal .card h3 {"} {
		if !strings.Contains(out, want) {
			t.Errorf("应有 %q，实际:\n%s", want, out)
		}
	}
	// 边界：不能把 .card-x 一起带出来（子串匹配就会犯这个错）
	if strings.Contains(out, ".card-x") {
		t.Errorf(".card 不该命中 .card-x:\n%s", out)
	}
}

func TestExtractRulesAcceptsLeadingDotAndTrims(t *testing.T) {
	out, err := extractRules(fixtureCSS, "  .grid-2  ")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, ".reveal .grid-2 {") {
		t.Fatalf("带点号/空格的入参应被归一化:\n%s", out)
	}
}

func TestExtractRulesNotFoundListsAvailable(t *testing.T) {
	_, err := extractRules(fixtureCSS, "hero")
	if err == nil {
		t.Fatal("不存在的类名应报错")
	}
	// 错误里必须带上可用类名，模型才能自纠
	for _, want := range []string{"hero", "card", "grid-2", "badge"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("错误信息应包含 %q: %s", want, err.Error())
		}
	}
}

func TestExtractRulesEmptyName(t *testing.T) {
	if _, err := extractRules(fixtureCSS, "   "); err == nil {
		t.Fatal("空 name 应报错并提示可用全文模式")
	}
}

func TestLoadComponentCSS(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, componentCSSFile), []byte(fixtureCSS), 0o644); err != nil {
		t.Fatal(err)
	}
	css, err := loadComponentCSS(dir)
	if err != nil {
		t.Fatal(err)
	}
	if css != fixtureCSS {
		t.Fatal("读回内容不一致")
	}

	if _, err := loadComponentCSS(t.TempDir()); err == nil {
		t.Fatal("文件缺失应报错")
	}
	if _, err := loadComponentCSS(""); err == nil {
		t.Fatal("目录未配置应报错")
	}
}

// TestRealComponentLibraryContract 契约测试：system prompt 里承诺的组件类，
// 真实组件库必须真的定义过——否则模型照着 prompt 写 class，页面上什么都不会发生。
// 这类"文档与实现漂移"是最难在联调中发现的一类 bug。
// 读不到真实文件时跳过（比如只拷贝了包目录跑测试），不让它依赖仓库布局变脆。
func TestRealComponentLibraryContract(t *testing.T) {
	p := filepath.Join("..", "..", "web", "assets", componentCSSFile)
	css, err := os.ReadFile(p)
	if err != nil {
		t.Skipf("读不到真实组件库（%s），跳过契约检查: %v", p, err)
	}
	names := strings.Join(extractClassNames(string(css)), " ")

	// prompt 的组件库章节承诺的组件集 v1 + v2。
	// prompt 里每新增一个 class，这里就得跟着加一条——这个测试存在的唯一目的
	// 就是拦住"prompt 承诺了、组件库里没有"的漂移。
	want := []string{
		// v1
		"center", "grid-2", "card", "quote", "steps", "badge", "tag", "footer", "muted",
		// v2 页面家具
		"page-no", "outline", "kicker", "rule",
		// v2 行清单
		"rows", "row", "line", "top", "key", "num",
		// v2 行内文本角色
		"term", "label", "sub", "footnote", "stat", "unit",
		// v2 网格扩展
		"grid-3", "grid-4",
		// v3 结构性版式（版式目录就是靠这些 class 落地的，缺一个就有一页没得选）
		"split", "flip", "middle", "bleed", "plate", "timeline", "band", "marginal", "note",
	}
	for _, w := range want {
		if !hasClassToken(names, w) {
			t.Errorf("prompt 承诺的 .%s 在真实组件库里找不到（文档与实现漂移了）", w)
		}
	}

	// v3 的每个版式类都必须带 .reveal 前缀：它们是"整块版式"，最可能被自定义槽或
	// 页面内联样式覆盖，前缀不齐就会重演"写了 CSS 却没生效"
	for _, w := range []string{"split", "bleed", "plate", "timeline", "band", "marginal"} {
		rules, err := extractRules(string(css), w)
		if err != nil {
			t.Errorf("抽 .%s 失败: %v", w, err)
			continue
		}
		if !strings.Contains(rules, ".reveal ."+w) {
			t.Errorf(".%s 的规则应带 .reveal 前缀:\n%s", w, rules)
		}
	}

	// v3 的宽容器要把文本角色类排除掉：`.reveal .bleed p`(0,2,1) 比 `.reveal .stat`(0,2,0)
	// 高一级，"一页一个大数字"会被渲染成一坨普通大段落——不报错，只是不对。
	// 这也正是 .card 那条守卫踩过的坑。
	for _, w := range []string{"bleed", "band"} {
		rules, err := extractRules(string(css), w)
		if err != nil {
			t.Fatalf("抽 .%s 失败: %v", w, err)
		}
		if len(guardRe.FindAllStringSubmatch(rules, -1)) == 0 {
			t.Errorf(".%s 的 p 规则缺少 :where(:not(...)) 守卫：里面的角色类会静默失效", w)
		}
	}

	// 组件库自己不许生产"左侧粗竖条"：一个彩色侧边条是模板感最容易被一眼认出来的
	// 标记（也是"AI 生成界面"清单上排第一的那条）。 .quote 是唯一的例外——它真的是引用块。
	if i := strings.Index(string(css), "组件集 v3"); i >= 0 {
		if strings.Contains(string(css)[i:], "border-left") {
			t.Error("v3 版式里出现了 border-left：左侧粗竖条是最典型的模板感标记，不要往组件库里再加")
		}
	} else {
		t.Error("组件库里找不到「组件集 v3」分段标记")
	}

	// 真实文件里 .card 必须能抽到规则，且规则带 .reveal 前缀
	rules, err := extractRules(string(css), "card")
	if err != nil {
		t.Fatalf("真实组件库里抽 .card 失败: %v", err)
	}
	if !strings.Contains(rules, ".reveal .card") {
		t.Fatalf(".card 的规则应带 .reveal 前缀（自定义槽覆盖时必须同前缀）:\n%s", rules)
	}

	// 组件库内部的优先级陷阱：`.reveal .card p`(0,2,1) 比文本角色类 `.reveal .stat`(0,2,0)
	// 高一级，会把卡内角色类整块盖掉——"hero 数字放进卡片"会渲染成普通灰段落，
	// 不报错、只是不对。守卫必须用 :where()（特异性为零），否则会连带改变这条规则本身的
	// 优先级，波及已有 deck 与自定义槽的覆盖能力。
	guards := guardRe.FindAllStringSubmatch(rules, -1)
	if len(guards) == 0 {
		t.Fatal(".card 的规则缺少 :where(:not(...)) 守卫：卡内的 .stat/.term 等会被 .card p 盖掉")
	}
	var excluded strings.Builder
	for _, g := range guards {
		excluded.WriteString(g[1] + ",")
	}
	for _, role := range []string{".stat", ".term", ".label", ".sub", ".footnote", ".kicker"} {
		if !strings.Contains(excluded.String(), role) {
			t.Errorf(".card 规则的排除名单里缺 %s（它会在卡片内部静默失效）", role)
		}
	}

	t.Logf("真实组件库类名: %s", names)
}
