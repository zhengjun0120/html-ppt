package deck

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 回声校验：CSS 变量只有在"某处真的 var() 消费了它"时才有意义。
// 组件库/主题/reveal 里出现的 var(--x) 就是"被消费"的证据——把这份清单当契约，
// AI 写了清单外的变量名（既不在契约里、也没在自己那份 CSS 里定义过）就直接拒收。
//
// 为什么必须拒收而不是只警告：那是一处**静默无效**的样式。页面不报错、结构也没坏，
// 只是"改了个寂寞"——AI 会向用户汇报"已调整标题字号"，而页面毫无变化，
// 用户开始怀疑整个工具。宁可用一次工具调用换掉这类不可见的失败。
//
// 契约从磁盘现扫而不是硬编码：components.css/theme.css 加了新变量，AI 立刻能用上，
// 不存在"文档与实现漂移"。

// contractCSSFiles 参与契约扫描的文件。
// reveal.css/moon.css 是第三方文件，但 --r-* 系列正是从它们那里被消费的，必须一起扫。
var contractCSSFiles = []string{
	"components.css",
	"theme.css",
	"reveal.css",
	"moon.css",
}

var (
	// varUsageRe 抓 var(--name) / var( --name , fallback )，只取变量名
	varUsageRe = regexp.MustCompile(`var\(\s*(--[a-zA-Z0-9_-]+)`)
	// varDefRe 抓自定义属性定义 --name:（var(--x) 里没有冒号，不会被误认成定义）
	varDefRe = regexp.MustCompile(`(--[a-zA-Z0-9_-]+)\s*:`)
	// styleAttrRe 抓 style="..." 的值：内容区禁 <style>，所以页面里的变量只可能在这里
	styleAttrRe = regexp.MustCompile(`(?i)\bstyle\s*=\s*"([^"]*)"`)
	// 注释剥离复用 theme_vars.go 的 cssCommentRe（同包，语义一致：注释里的花括号/变量
	// 不该参与校验）
)

// loadVariableContract 扫框架 CSS，收集所有被 var() 消费过的变量名。
// 每次调用现扫（几个文件、几十 KB），保证开发期改了组件库立刻生效，不做缓存。
func loadVariableContract(assetsDir string) (map[string]bool, error) {
	if strings.TrimSpace(assetsDir) == "" {
		return nil, fmt.Errorf("资源目录未配置")
	}
	known := map[string]bool{}
	scanned := 0
	for _, name := range contractCSSFiles {
		data, err := os.ReadFile(filepath.Join(assetsDir, name))
		if err != nil {
			// 单个文件缺失（比如换了主题没带 moon.css）不该让整项校验失效
			continue
		}
		scanned++
		for _, m := range varUsageRe.FindAllStringSubmatch(string(data), -1) {
			known[m[1]] = true
		}
	}
	if scanned == 0 {
		return nil, fmt.Errorf("在 %s 下没找到任何契约文件（%s）", assetsDir, strings.Join(contractCSSFiles, " "))
	}
	return known, nil
}

// definedVariables 收集一段 CSS/HTML 里**定义**过的变量名。
// 自定自用是合法写法（`.hero{--shadow:...;box-shadow:var(--shadow)}`），必须豁免。
func definedVariables(text string) map[string]bool {
	out := map[string]bool{}
	for _, m := range varDefRe.FindAllStringSubmatch(cssCommentRe.ReplaceAllString(text, ""), -1) {
		out[m[1]] = true
	}
	return out
}

// unknownVariables 返回既不在契约里、也没有定义过的变量名（去重、排序）。
//
// 豁免并入函数内部而不是交给调用方：文本自己定义过的变量（`--glow:…` 与 `var(--glow)`
// 写在同一份内容里）天然自足，调用方很容易忘了合并、把正常写法误拒——
// 误拒比漏判更糟（它会把能用的写法拦住）。所以这里自己兜住，defined 只用于
// "外部已经定义过"的补充（比如 deck 主题块里的 vars 调色板）。
func unknownVariables(text string, known, defined map[string]bool) []string {
	body := cssCommentRe.ReplaceAllString(text, "")
	defs := definedVariables(body)
	for k := range defined {
		defs[k] = true
	}

	seen := map[string]bool{}
	var out []string
	for _, m := range varUsageRe.FindAllStringSubmatch(body, -1) {
		name := m[1]
		if known[name] || defs[name] || seen[name] {
			continue
		}
		seen[name] = true
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// unknownVarErr 生成"给 LLM 看"的错误：指出哪些名字是无效的、可用的有哪些、自定自用该怎么写。
func unknownVarErr(unknown []string, known map[string]bool) error {
	valid := make([]string, 0, len(known))
	for k := range known {
		valid = append(valid, k)
	}
	sort.Strings(valid)
	return fmt.Errorf("这些变量没有任何样式消费它们，写了也不会生效：%s。可用变量：%s。"+
		"如果你需要自己的变量，请在**同一份 CSS 里**先定义再使用（如 .hero{--glow:0 0 20px var(--accent);box-shadow:var(--glow)}）",
		strings.Join(unknown, " "), strings.Join(valid, " "))
}

// collectInlineStyleVars 从页面 HTML 里只取 style="..." 属性的值（精确，不会把正文里
// 顺手提到的 var(--x) 也算进来），再从中提取变量用法。
func collectInlineStyleVars(html string) string {
	var sb strings.Builder
	for _, m := range styleAttrRe.FindAllStringSubmatch(html, -1) {
		sb.WriteString(m[1])
		sb.WriteString("\n")
	}
	return sb.String()
}

// deckStyleBlocks 取 deck 自己的样式块文本（自定义槽 + 主题 override）。
// 这些块里定义过的变量，页面内联样式里也可以用——所以内容侧回声校验必须把它们
// 算作"已定义"，否则"槽里定义变量、页面里 var() 使用"这种正常写法会被误判。
func deckStyleBlocks(deckHTML string) string {
	if strings.TrimSpace(deckHTML) == "" {
		return ""
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(deckHTML))
	if err != nil {
		return ""
	}
	var sb strings.Builder
	for _, id := range []string{customCSSBlockID, "deck-theme-override"} {
		if sel := doc.Find("#" + id); sel.Length() > 0 {
			sb.WriteString(sel.First().Text())
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// checkCSSBalance 粗查花括号配对。不配对的 CSS 会让它后面的所有规则一起失效
// （一个漏掉的 } 就能吃掉整个文件的后续样式），是最值得拦的语法错误。
// 先剥注释与字符串字面量，避免 content:"{" 这类写法误报。
func checkCSSBalance(css string) error {
	text := cssCommentRe.ReplaceAllString(css, "")
	text = regexp.MustCompile(`"(?:\\.|[^"\\])*"`).ReplaceAllString(text, `""`)
	text = regexp.MustCompile(`'(?:\\.|[^'\\])*'`).ReplaceAllString(text, `''`)

	depth := 0
	for _, r := range text {
		switch r {
		case '{':
			depth++
		case '}':
			depth--
			if depth < 0 {
				return fmt.Errorf("CSS 里有多余的 }（花括号不配对）：它会让这处之后的所有规则一起失效，请检查后再提交")
			}
		}
	}
	if depth != 0 {
		return fmt.Errorf("CSS 里有 %d 个 { 没有闭合：漏掉的 } 会让后面所有规则失效，请补齐后再提交", depth)
	}
	return nil
}
