package deck

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// 主题契约变量的所有权。这个文件只回答一个问题：
// **哪些 CSS 变量归主题管，谁都不许在别处重新定义它。**
//
// 为什么必须有一条硬边界：`#deck-theme-override` 和 `#deck-custom` 都是 `:root{}`
// 选择器（优先级完全相同），而自定义槽排在主题块之后——同级、靠后者取胜。
// 于是自定义槽里写一句 `:root{--accent:#ff8800}` 就会静默盖住主题，
// 用户之后说"换个配色"，update_theme 返回成功、页面一动不动。
// 级联顺序本身是对的（自定义槽本来就该能覆盖主题做局部视觉），
// 出问题的是"让自定义槽重定义契约变量"——那把它从局部覆盖变成了全局夺权。
//
// 所以规则是：契约变量只能由 Theme 的结构化字段产生，一处定义、一处渲染。

// maxVarValueBytes 单个变量值的长度上限。变量值最终会拼进 :root{} 里，
// 不限长的话一个变量就能把主题块撑成一大坨（也是 token 上限）。
const maxVarValueBytes = 200

// reservedVars 是 Theme 的结构化字段派生出来的变量名，外加**框架层的结构常数**
// （字号阶梯、间距阶梯、安全边距——它们同样只有一处定义，在 theme.css 里）。
// 这些名字只能由 renderThemeCSS / theme.css 产生，不接受外部传入。
var reservedVars = map[string]bool{
	"--accent":      true,
	"--border":      true,
	"--hairline":    true,
	"--card-bg":     true,
	"--accent-soft": true,
	"--on-accent":   true,
	"--text-muted":  true,
	"--radius":      true,
	// 整页翻色的派生 token（renderThemeCSS 每次全量重算，见 theme.go）。
	// 拦掉的理由和其它契约变量一样：vars 拼在同一块 :root 里，同名键会
	// **静默盖住**派生值——墨页/强调页的翻色从此坏掉，且没有任何报错。
	"--ink-bg": true, "--ink-fg": true, "--ink-accent": true,
	"--ink-muted": true, "--ink-hairline": true, "--ink-panel": true,
	"--page-accent": true, "--page-fg": true,
	"--accent-muted": true, "--accent-hairline": true, "--accent-panel": true,
}

// reservedVarPrefixes 按前缀拦一类名字。用前缀而不是把每个名字列出来：
// 阶梯是可扩展的（加一档字号不该要求同步改这张表），而漏改这一处就会让
// "自定义槽静默盖住框架常数"的坑重新打开。
var reservedVarPrefixes = []string{"--r-", "--space-", "--fs-", "--slide-pad-"}

// isReservedVarName 判断一个变量名是否归框架所有（主题字段 + 结构常数）。
func isReservedVarName(name string) bool {
	if reservedVars[name] {
		return true
	}
	for _, p := range reservedVarPrefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
}

// customVarNamePattern 合法的自定义变量名：-- 开头、字母打头、只含字母数字连字符。
var customVarNamePattern = regexp.MustCompile(`^--[a-zA-Z][a-zA-Z0-9-]*$`)

// varValueForbidden 变量值里不许出现的片段。
// 前五个会破坏 CSS 结构（变量值是被拼进 :root{...} 里的，一个 } 就能提前收尾、
// 把后面的声明变成垃圾）；后三个是外部加载/执行通道，和自定义槽同一套标准。
var varValueForbidden = []string{";", "{", "}", "<", ">", "\\", "@import", "url(", "expression("}

// validateThemeVars 校验整套自定义调色板。
// 按键排序后校验：map 遍历顺序随机，不排的话"多个变量都非法"时报哪个错是不确定的，
// 测试也就没法稳定断言。
func validateThemeVars(vars map[string]string) error {
	if len(vars) == 0 {
		return nil
	}
	if len(vars) > maxThemeVars {
		return fmt.Errorf("自定义调色板最多 %d 个变量（收到 %d 个）：看起来在读一整套设计系统的变量，请只保留页面真正会用到的",
			maxThemeVars, len(vars))
	}

	names := make([]string, 0, len(vars))
	for n := range vars {
		names = append(names, n)
	}
	sort.Strings(names)

	for _, name := range names {
		if err := validateVarName(name); err != nil {
			return err
		}
		if err := validateVarValue(name, vars[name]); err != nil {
			return err
		}
	}
	return nil
}

func validateVarName(name string) error {
	if !customVarNamePattern.MatchString(name) {
		return fmt.Errorf("自定义变量名 %q 不合法：要写成 --xxx（两个连字符开头、字母打头，可含数字和连字符）", name)
	}
	if isReservedVarName(name) {
		return fmt.Errorf("%s 是框架契约变量，不能通过 vars 定义——"+
			"颜色类由 accent/background/heading_color/text_color/surface/border_color 这些字段统一派生，"+
			"字号/间距/边距是组件库的结构常数（--fs-* / --space-* / --slide-pad-*），"+
			"你要改观感就直接传对应字段。两处定义同名变量时谁生效取决于渲染顺序，正是'改了没反应'的经典来源", name)
	}
	return nil
}

func validateVarValue(name, val string) error {
	v := strings.TrimSpace(val)
	if v == "" {
		return fmt.Errorf("自定义变量 %s 的值不能为空（要删掉它就提交一份不含它的 vars）", name)
	}
	if len(v) > maxVarValueBytes {
		return fmt.Errorf("自定义变量 %s 的值超过 %d 字节上限：请只写真正必要的值", name, maxVarValueBytes)
	}
	low := strings.ToLower(v)
	for _, f := range varValueForbidden {
		if strings.Contains(low, f) {
			return fmt.Errorf("自定义变量 %s 的值里不能出现 %q：它会破坏 CSS 结构或形成外部加载通道", name, f)
		}
	}
	return nil
}

// cssCommentRe 去掉 CSS 注释。判定"有没有重定义契约变量"之前必须先剥注释——
// 否则一句 `/* 别在这里写 --accent: 值 */` 会被当违规拒绝，而这条注释恰恰是对的。
var cssCommentRe = regexp.MustCompile(`(?s)/\*.*?\*/`)

// reservedVarDeclRe 匹配"重定义契约变量"的声明（后跟冒号才算声明）。
// 与 isReservedVarName 覆盖同一批名字：颜色类逐个列（数量固定），
// 阶梯类按**前缀**匹配（--space-*/--fs-*/--slide-pad-*），加一档字号不用来这里同步。
// --r- 前缀是 reveal.js 自己的变量，同样归主题字段所有，一并纳入。
var reservedVarDeclRe = regexp.MustCompile(
	`--(?:accent|border|hairline|card-bg|accent-soft|accent-muted|accent-hairline|accent-panel` +
		`|ink-bg|ink-fg|ink-accent|ink-muted|ink-hairline|ink-panel` +
		`|page-accent|page-fg|on-accent|text-muted|radius` +
		`|space-[a-z0-9]+|fs-[a-z0-9]+|slide-pad-[a-z]+|r-[a-z0-9-]+)\s*:`)

// findReservedVarRedefinition 在自定义 CSS 里找出第一个被重定义的契约变量，
// 返回变量名（含 --）。没找到返回空串。
func findReservedVarRedefinition(css string) string {
	body := cssCommentRe.ReplaceAllString(css, "")
	m := reservedVarDeclRe.FindString(body)
	if m == "" {
		return ""
	}
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(m), ":"))
}

// reservedVarRedefinitionErr 组装给模型的拒绝理由。
// 按"错误即反馈"写：说清为什么不许（会静默盖住主题）、以及正确做法是什么
// （换成 update_theme 的 vars 参数），模型就能自己改对。
func reservedVarRedefinitionErr(name string) error {
	return fmt.Errorf("自定义样式槽不能重定义主题变量 %s。原因：自定义槽排在主题块之后、且两者都是 :root 选择器（优先级相同），"+
		"写在这里的值会静默盖住主题——用户之后说'换个配色'，update_theme 会成功返回但页面毫无变化，是一类查不出来的假成功。"+
		"要改这些变量请直接调 update_theme 传对应字段（--accent←accent、--r-main-color←text_color 等）；"+
		"要设计一整套自定义配色，请用 update_theme 的 vars 参数——它和主题写在同一个块里，永远只有一处权威。", name)
}
