package deck

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 写入时的样式体检：把"模型悄悄抛弃了主题体系"这类静默失败，变成工具结果里的一句话。
//
// 为什么需要它：消毒闸门只拦得住"危险"的东西，拦不住"合法但会让整套 deck 烂掉"的东西。
// 实测三份手写型 deck，可编辑区里 43%~51% 的字符是逐字重复的内联样式（同一串
// "position:absolute; top:0; right:0; ..." 在 8 页里一字不差抄了 8 遍），另外还有
// 在 <section> 上写死 background/color/font-family 从而让 update_theme 彻底失效的。
// 这些都不报错、导出也正常，只是 deck 变得无法再换风格——是最难自查的一类失败。
//
// 策略与消毒一致：只提示、不改写、不阻塞。
// 不做自动改写（把重复内联提成 class）的原因：模型下一次 read_slide 会读到一份
// 它没写过的 HTML，"我写的 = 我读到的"这个自洽性一破，后面的 fingerprint 比对、
// 增量修改全都会开始出现无法解释的差异。提示它自己改，成本更低也更安全。

const (
	lintHardColorLimit = 3 // 写死的颜色值达到这个数才提示（1~2 处往往是刻意的偏离）
	lintPxFontLimit    = 3 // 同上，px 字号
	// 内联字号是"一次就够"：组件库的字号是六档阶梯，任何手写的字号都会插进两档之间，
	// 而且 prompt 里是明令禁止的（见"字号不要自己定"一节），所以 1 处就提示。
	lintInlineFontLimit = 1
	lintDupLimit        = 3 // 同一串 style 值出现到这个数才提示
	lintMaxSamples      = 2 // 提示里最多举几个例子
	lintSampleRunes     = 46
)

var (
	// 写死的颜色：hex 与 rgb()/hsl() 函数式。刻意不匹配颜色关键字——
	// transparent / currentColor 是合法的、且常被用来做效果。
	hardColorRe = regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b|\b(?:rgba?|hsla?)\(`)

	// emojiRe 抓"表情当图标/点缀"：主表情平面（1F300–1FAFF，🔥🚀📊 这一带）、
	// 杂项符号（2600–26FF，☀⚙⚠ 这一带）和变体选择符 FE0F。
	// 刻意不含 2700–27BF 区（✓✔✗ 这类对勾/叉在对照表里是正当排版符号），
	// 也不含箭头区（→ 同理）——误报会让模型连正当符号都不敢用，宁可少抓。
	emojiRe = regexp.MustCompile(`[\x{1F300}-\x{1FAFF}]|[\x{2600}-\x{26FF}]|[\x{FE0F}]`)

	// emDashRunRe 把连续的 — 归并成"一处"：中文破折号 —— 是两个 U+2014 字符，
	// 按字符数会把这个合法用法误判成两处。"一页最多一个"数的是归并后的处数。
	emDashRunRe = regexp.MustCompile(`—+`)
)

// themeLevelProps 是"写在 <section> 上就等于覆盖整套主题"的属性。
// font-size 也算：适配兜底就是靠改 section 的 font-size 工作的，模型写在这里的值
// 会在第一次 fit 时被直接覆盖掉，属于静默失效。
var themeLevelProps = map[string]bool{
	"background":       true,
	"background-color": true,
	"background-image": true,
	"color":            true,
	"font-family":      true,
	"font-size":        true,
}

type styleReport struct {
	HardColors     int      // 内联样式里写死的颜色值数量
	PxFontSizes    int      // 内联样式里的 px 字号数量
	InlineFontSize int      // 内联样式里自己定的字号（不论单位）——插进阶梯两档之间
	DupStyles      int      // 落在"重复值"上的内联样式总条数
	DupSamples     []string // 重复值样本（形如 "display:flex; …" ×4）
	SectionProps   []string // 直接写在 <section> 上的主题级属性名
	SolePageBlocks int      // "整页只剩一块色块"的页数（想把整页换色却用了页内色块）
	EmojiPages     []string // 文本里出现 emoji 的页（"第 3 页（🔥🚀×2）"）
	DashPages      []string // 破折号超过一处的页（"第 5 页（3 处）"）
}

// styleLinter 是跨页累加器。必须跨页统计：页码角标、讲次提头这类"页面家具"
// 在每一页只出现一次，逐页看永远是"不重复"的——只有把整份提交放一起数，
// 才看得见"同一串 177 字符抄了 8 遍"。这是这个功能存在的核心理由。
type styleLinter struct {
	freq  map[string]int
	orig  map[string]string // 首次出现的原始写法，用于提示时展示
	rep   styleReport
	pages int // 已收的页数；用于把"第 N 页"在单页场景下说成"本页"
}

func newStyleLinter() *styleLinter {
	return &styleLinter{freq: map[string]int{}, orig: map[string]string{}}
}

// pageColorBlockClasses 是"页内整面色块"组件。section 的**唯一元素子节点**是它们
// 之一时，模型想要的其实是"整页换色"，但产物是"贴在页面底色上的一张大卡片"——
// 色块四周露出页面底色，视觉上非常突兀（用户拿真实截图反馈过）。整页底色只归
// section 级 .bg-ink/.bg-accent 管（框架经 reveal 背景层铺满画布，PDF 也正确）。
var pageColorBlockClasses = map[string]bool{
	"accent-block": true, "ink-block": true, "tint-block": true,
}

// add 收一页（含 <section> 自身）。
func (l *styleLinter) add(sec *goquery.Selection) {
	l.pages++
	l.scan(sec.AttrOr("style", ""), true)
	sec.Find("[style]").Each(func(_ int, s *goquery.Selection) {
		l.scan(s.AttrOr("style", ""), false)
	})
	l.checkSolePageColorBlock(sec)
	l.checkText(sec)
}

// checkText 扫页面文本里两类"提示词管不住、只有机制能兜底"的东西：
// emoji（提示词禁用，但模型违规时此前只有截图审查可能偶然看到）和
// 超量的破折号（"一页最多一个"是明文规则，数字数是机器的强项）。
// 都只提示不阻塞——用户明确要求 emoji 时模型保留它是正确行为，措辞里留了这个出口。
func (l *styleLinter) checkText(sec *goquery.Selection) {
	text := sec.Text()
	if hits := emojiRe.FindAllString(text, -1); len(hits) > 0 {
		sample := hits
		if len(sample) > 4 {
			sample = sample[:4]
		}
		l.rep.EmojiPages = append(l.rep.EmojiPages,
			fmt.Sprintf("%s（%s×%d）", l.pageLabel(), strings.Join(sample, ""), len(hits)))
	}
	if n := len(emDashRunRe.FindAllString(text, -1)); n > 1 {
		l.rep.DashPages = append(l.rep.DashPages, fmt.Sprintf("%s（%d 处）", l.pageLabel(), n))
	}
}

// pageLabel 单页场景（update_slide 一次只收一页）说"本页"，
// 说"第 1 页"是错的——它在这份 deck 里未必是第 1 页。
func (l *styleLinter) pageLabel() string {
	if l.pages == 1 {
		return "本页"
	}
	return fmt.Sprintf("第 %d 页", l.pages)
}

// checkSolePageColorBlock 检出"section 里只剩一块色块"的写法。判定是结构性的：
// 元素子节点恰好一个、且它是整面色块、且 section 自身没有挂 bg-ink/bg-accent
//（挂了就是有意的"色面页上放同色块"，不警告）。不阻塞写入——确要"带框海报"
// 的效果可以保留，警告只是把更可能的意图（整页换色）指给模型。
func (l *styleLinter) checkSolePageColorBlock(sec *goquery.Selection) {
	cls, _ := sec.Attr("class")
	for _, c := range strings.Fields(cls) {
		if c == "bg-ink" || c == "bg-accent" {
			return
		}
	}
	kids := sec.Children()
	if kids.Length() != 1 {
		return
	}
	childCls, _ := kids.First().Attr("class")
	for _, c := range strings.Fields(childCls) {
		if pageColorBlockClasses[c] {
			l.rep.SolePageBlocks++
			return
		}
	}
}

func (l *styleLinter) scan(style string, isSection bool) {
	style = strings.TrimSpace(style)
	if style == "" {
		return
	}

	decls := styleDecls(style)
	for _, d := range decls {
		if hardColorRe.MatchString(d.value) {
			l.rep.HardColors++
		}
		// px 字号与"自己定字号"互斥：一条声明只报一次，报的是更要紧的那个问题。
		// （px 字号两个毛病都占：不随适配兜底缩放，也不在阶梯上。）
		switch {
		case isPxFontSize(d.prop, d.value):
			l.rep.PxFontSizes++
		case !isSection && isFontSizeProp(d.prop):
			l.rep.InlineFontSize++
		}
		if isSection && themeLevelProps[d.prop] {
			l.rep.SectionProps = appendUnique(l.rep.SectionProps, d.prop)
		}
	}

	key := canonicalStyle(decls)
	if l.freq[key] == 0 {
		l.orig[key] = style
	}
	l.freq[key]++
}

// report 结算重复值统计。调用时机是"这一批写完之后"，所以必须显式调用一次。
func (l *styleLinter) report() styleReport {
	rep := l.rep
	sort.Strings(rep.SectionProps)

	var keys []string
	for k, n := range l.freq {
		if n >= lintDupLimit {
			keys = append(keys, k)
		}
	}
	// 按"提成组件后能省下的字符数"降序：最该被组件化的排最前面
	sort.Slice(keys, func(i, j int) bool {
		si := (l.freq[keys[i]] - 1) * len(l.orig[keys[i]])
		sj := (l.freq[keys[j]] - 1) * len(l.orig[keys[j]])
		if si != sj {
			return si > sj
		}
		return keys[i] < keys[j]
	})
	for _, k := range keys {
		rep.DupStyles += l.freq[k]
		if len(rep.DupSamples) < lintMaxSamples {
			rep.DupSamples = append(rep.DupSamples,
				fmt.Sprintf("%s ×%d", clipStyle(l.orig[k]), l.freq[k]))
		}
	}
	return rep
}

// Warning 组装提示语。没超阈值时返回空串（omitempty 会省略，模型看不到噪音）。
func (r styleReport) Warning() string {
	var parts []string
	if len(r.SectionProps) > 0 {
		parts = append(parts, fmt.Sprintf(
			"在 <section> 上直接写了 %s——那是整页覆盖主题，用户之后换配色不会生效，整页级视觉请走 update_custom_css",
			strings.Join(r.SectionProps, "、")))
	}
	if r.HardColors >= lintHardColorLimit {
		parts = append(parts, fmt.Sprintf(
			"%d 处写死的颜色值，会让 update_theme 失效，改用 var(--accent)/var(--text-muted) 等主题变量", r.HardColors))
	}
	if r.PxFontSizes >= lintPxFontLimit {
		parts = append(parts, fmt.Sprintf(
			"%d 处 px 字号，不随适配兜底等比缩放（装不下时会被裁掉），改用 em", r.PxFontSizes))
	}
	if r.InlineFontSize >= lintInlineFontLimit {
		parts = append(parts, fmt.Sprintf(
			"%d 处内联 font-size：组件库的字号是固定阶梯（0.7 / 0.92 / 1.1 / 1.32 / 1.6 / 2.4 / 4.6 倍基准），"+
				"手写的字号会插进两档之间，让人读不出哪一层更重要。要突出就用角色类（.stat / .term / .label / .band），"+
				"整页级的锚点用 .cover-title / .big-num（.bg-ink/.bg-accent 页面上它们最出彩）；"+
				"或把它单独放一页用 .bleed；整份 deck 字太小请调 update_theme 的 canvas", r.InlineFontSize))
	}
	if r.SolePageBlocks > 0 {
		parts = append(parts, fmt.Sprintf(
			"%d 页的 section 里只剩一整面色块（.accent-block/.ink-block/.tint-block）——它四周会露出页面底色，"+
				"看起来像一张贴上去的卡片。整页换色请改用 section 级 class（如 <section class=\"bg-accent\">，"+
				"框架会把底色铺满整个画布）；色块只用于页内的一块面，若确要保留这种\"带框\"效果可忽略本提示",
			r.SolePageBlocks))
	}
	if len(r.EmojiPages) > 0 {
		parts = append(parts, fmt.Sprintf(
			"%s 的文本里有 emoji——emoji 是禁用的（用户明确要求除外），图标请用 .ico 放矢量图标"+
				"（path 从 read_icons 的清单里复制），表情当图标是\"内容不够、表情来凑\"",
			strings.Join(r.EmojiPages, "、")))
	}
	if len(r.DashPages) > 0 {
		parts = append(parts, fmt.Sprintf(
			"%s 的破折号——超过一处：一页最多一个，其余的地方用逗号、冒号或句号把句子拆开",
			strings.Join(r.DashPages, "、")))
	}
	if r.DupStyles > 0 {
		parts = append(parts, fmt.Sprintf(
			"有 %d 处内联样式落在重复值上（%s）——同一个 style 值出现第 2 次就该改用组件库 class，先 read_component 看有哪些可用",
			r.DupStyles, strings.Join(r.DupSamples, "、")))
	}
	if len(parts) == 0 {
		return ""
	}
	return "本页样式体检（未阻塞写入，但请修正）：" + strings.Join(parts, "；") + "。"
}

type styleDecl struct{ prop, value string }

// styleDecls 把内联样式切成声明列表。内联样式的值里不会出现分号
// （渐变、rgba() 都在括号内），按分号切足够。
func styleDecls(style string) []styleDecl {
	var out []styleDecl
	for _, part := range strings.Split(style, ";") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		i := strings.Index(part, ":")
		if i <= 0 {
			continue
		}
		out = append(out, styleDecl{
			prop:  strings.ToLower(strings.TrimSpace(part[:i])),
			value: strings.TrimSpace(part[i+1:]),
		})
	}
	return out
}

// canonicalStyle 归一化后再比：声明顺序和空格差异不该影响"这是同一串样式"的判断
// （模型手写的重复往往就是这样差一个空格、差一个顺序）。
func canonicalStyle(decls []styleDecl) string {
	parts := make([]string, 0, len(decls))
	for _, d := range decls {
		parts = append(parts, d.prop+":"+d.value)
	}
	sort.Strings(parts)
	return strings.Join(parts, ";")
}

func isPxFontSize(prop, value string) bool {
	if !isFontSizeProp(prop) {
		return false
	}
	return strings.Contains(strings.ToLower(value), "px")
}

// isFontSizeProp 判断这条声明是不是在设字号（简写 font 也算）。
func isFontSizeProp(prop string) bool {
	return prop == "font-size" || prop == "font"
}

func appendUnique(list []string, v string) []string {
	for _, x := range list {
		if x == v {
			return list
		}
	}
	return append(list, v)
}

// clipStyle 压缩空白并截断，用于提示语里展示样例（按 rune 截，不劈开多字节字符）。
func clipStyle(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	r := []rune(s)
	if len(r) <= lintSampleRunes {
		return s
	}
	return string(r[:lintSampleRunes]) + "…"
}

// joinWarnings 合并两段非空警告，保持"无违规时是空串"这个约定不变。
func joinWarnings(parts ...string) string {
	var out []string
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return strings.Join(out, " ")
}
