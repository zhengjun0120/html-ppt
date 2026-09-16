package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// 组件库对 AI 只读：read_component 只读文件、不写文件，组件库也永远不进任何写工具的
// 目标列表（护栏）。这个文件里全是纯函数，方便直接测。
//
// 为什么需要这个工具：AI 要写自定义样式（update_custom_css）覆盖 .card / .quote，
// 就得先知道它们现在长什么样——不然只能凭想象写选择器和属性，覆盖出来的效果全靠运气。

const componentCSSFile = "components.css"

// iconCatalogFile 是预置图标清单（Tabler 描边 path，.ico 槽的唯一合法 path 来源）。
// 单独成文件而不是写进组件库 CSS：组件库会被浏览器加载，图标清单只有模型消费，
// 不该让每个 deck 页面白白多下载几 KB。
const iconCatalogFile = "icons.md"

// frameworkClasses 是框架命名空间（reveal 的根类），不算"可用组件类"：
// 它们出现在每条选择器的前缀里，但 prompt 明确规定 AI 不许动它们。
var frameworkClasses = map[string]bool{"reveal": true, "slides": true}

// specificityNote 是必须转达给模型的一句话：组件规则都带 .reveal 前缀，
// 自定义槽里写裸的 .card 会因优先级不足被组件库盖掉（无论级联顺序如何）——
// 这正是"我写了 CSS 但页面没变化"最常见的原因。
//
// 还要带上"层级一样深"这半句：组件库里不只有 .reveal .card 这种两段式，
// 还有 .reveal .rows .quote 这种三段式（0,3,0）。只补前缀、深度不够，照样盖不动——
// 而 read_component 会把目标规则的原文返回，照抄它的选择器永远是最稳的做法。
const specificityNote = "组件库的选择器都带 .reveal 前缀（如 .reveal .card），有些还多一层（如 .reveal .rows .quote）。覆盖它们时选择器要带上同样的前缀、且层级一样深，否则优先级不足、会被组件库盖掉——表现就是'写了 CSS 却没生效'。最稳的做法是先用 read_component 看目标规则的选择器原文，照抄结构再改属性。"

// classRe 从选择器里抠类名。要求点号后跟字母，所以 0.6em 这类数值不会被误抓。
var classRe = regexp.MustCompile(`\.([a-zA-Z][a-zA-Z0-9_-]*)`)

// loadComponentCSS 读共享组件库。每次调用都重新读盘：文件很小，
// 而且开发期改了 components.css 应该立刻被 AI 看到，不做缓存。
func loadComponentCSS(assetsDir string) (string, error) {
	if strings.TrimSpace(assetsDir) == "" {
		return "", fmt.Errorf("组件库目录未配置")
	}
	p := filepath.Join(assetsDir, componentCSSFile)
	data, err := os.ReadFile(p)
	if err != nil {
		return "", fmt.Errorf("读取组件库失败: %w", err)
	}
	return string(data), nil
}

// loadIconCatalog 读图标清单。为什么不给模型"自己画 path"的自由：LLM 手绘的
// SVG 路径大多是不成形的曲线团，而图标画错比没有图标更扎眼；审核过的清单
// （Tabler 描边，与 .ico 的 currentColor/fill:none 口径逐字匹配）复制即用、零适配。
func loadIconCatalog(assetsDir string) (string, error) {
	if strings.TrimSpace(assetsDir) == "" {
		return "", fmt.Errorf("组件库目录未配置")
	}
	p := filepath.Join(assetsDir, iconCatalogFile)
	data, err := os.ReadFile(p)
	if err != nil {
		return "", fmt.Errorf("读取图标清单失败: %w", err)
	}
	return string(data), nil
}

// splitRules 把 CSS 粗切成一条条规则（整块返回，含外层 @media 时整块是一条）。
// 必须按花括号配对切，不能按 } 朴素切：v4 之后的组件库有真正的嵌套块
// （@media print 里两条打印规则），按 } 切会把 @media 的**第二条**内层规则
// 泄漏成一条独立"选择器"——extractClassNames 会从里面抠出 .reveal-print
// 这种框架类并要求 prompt 收录它（真实发生过的测试失败）。
// 嵌套块的"选择器"取首 个 { 之前的文本（@media print），类名提取自然为空，
// 不会误伤 extractClassNames；extractRules 也只会命中平铺规则，可接受。
func splitRules(css string) []string {
	// 先去掉注释，免得注释里的 } 把配对切歪
	noComment := regexp.MustCompile(`(?s)/\*.*?\*/`).ReplaceAllString(css, "")
	var rules []string
	depth, start := 0, 0
	for i, r := range noComment {
		switch r {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				if rule := strings.TrimSpace(noComment[start : i+1]); rule != "" {
					rules = append(rules, rule)
				}
				start = i + 1
			}
		}
	}
	return rules
}

// extractClassNames 列出组件库定义过的类名（去重、排序，跳过框架命名空间）。
func extractClassNames(css string) []string {
	seen := map[string]bool{}
	var out []string
	for _, rule := range splitRules(css) {
		selector := rule
		if i := strings.Index(rule, "{"); i >= 0 {
			selector = rule[:i]
		}
		for _, m := range classRe.FindAllStringSubmatch(selector, -1) {
			name := m[1]
			if frameworkClasses[name] || seen[name] {
				continue
			}
			seen[name] = true
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}

// extractRules 取出所有选择器里包含该类名的规则（原样返回 CSS 片段）。
// 找不到时返回错误并附可用类名清单——"错误即反馈"，让模型自己改对名字，
// 而不是拿到一段空文本在那儿猜。
func extractRules(css, name string) (string, error) {
	name = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(name), "."))
	if name == "" {
		return "", fmt.Errorf("name 不能为空；不传 name 可以拿到组件库全文")
	}
	want := "." + name
	var hits []string
	for _, rule := range splitRules(css) {
		selector := rule
		if i := strings.Index(rule, "{"); i >= 0 {
			selector = rule[:i]
		}
		// 按"类名边界"匹配：.card 命中 .card 与 .card h3，但不会命中 .card-x
		if hasClassToken(selector, want) {
			hits = append(hits, rule)
		}
	}
	if len(hits) == 0 {
		return "", fmt.Errorf("组件库里没有 .%s 这个类。可用类名：%s",
			name, strings.Join(extractClassNames(css), " "))
	}
	return strings.Join(hits, "\n"), nil
}

// hasClassToken 判断选择器里是否真的出现了这个类（按 token 边界，不是子串匹配）。
func hasClassToken(selector, want string) bool {
	idx := 0
	for {
		i := strings.Index(selector[idx:], want)
		if i < 0 {
			return false
		}
		end := idx + i + len(want)
		// 后面不能再跟类名字符，否则是别的类（.card vs .card-x）
		if end >= len(selector) || !isClassChar(selector[end]) {
			return true
		}
		idx = end
	}
}

func isClassChar(b byte) bool {
	return b == '-' || b == '_' || (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}
