package deck

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// Theme 是 deck 的语义化主题配置，真身是 deck.html 里的
// <script type="application/json" id="deck-theme"> 块。
// 结构层 components.css 只引用 CSS 变量；本结构决定变量的取值。
// 暴露给 LLM 的只有这些语义字段，不是 CSS 原文——校验做硬、可逆、出错率低。
//
// 两处"自由"的边界，刻意分得很清：
//   - 配色可以完全自由：Vars 让 AI 设计一整套自定义调色板，
//     但必须写成"一处变量定义"，而不是散落在元素上的字面量。
//   - 画布不许自由：Canvas 只给比例预设。画布尺寸 + em 尺度是适配兜底、
//     翻页动画、字号缩放三件事的共同地基，放开绝对尺寸会让它们同时失灵
//     （deck-0004 那 18 处 font-size:30px 就是这个后果）。
type Theme struct {
	Accent       string            `json:"accent"`        // 强调色，border/card-bg 由它派生
	Background   string            `json:"background"`    // 页面背景
	HeadingColor string            `json:"heading_color"` // 标题色
	TextColor    string            `json:"text_color"`    // 正文色
	Font         string            `json:"font"`          // 字体方案枚举：sans/serif/mono
	Radius       string            `json:"radius"`        // 卡片圆角
	Transition   string            `json:"transition"`    // 翻页动画枚举，init.js 读走
	Canvas       string            `json:"canvas"`        // 画布比例预设，init.js 读走
	Vars         map[string]string `json:"vars,omitempty"` // 自定义调色板（变量名 → CSS 值）
}

// canvasPresets 是画布比例预设：**高度统一 700，只调宽度**。
//
// 为什么固定高度：一页能放多少内容由高度决定，宽度只影响排布的宽松度。
// 高度一动，原本刚好放得下的页面就会被挤爆（触发适配兜底缩小，牺牲可读性）；
// 固定高度则"换比例"是纯粹的横向增益。
//
// 960×700 是 reveal.js 的默认逻辑画布，也是 init.js 里 config() 的兜底值——
// 三处必须一致：本表、init.js 的 CANVAS 表、reveal 默认值。
// 有 TestCanvasPresetsMatchInitJS 守着前两处不漂移。
const (
	CanvasStandard = "standard" // 通用（1.37:1，在 16:9 屏上左右有黑边）
	CanvasWide     = "wide"     // 16:9（比例与投屏/录屏一致，没有黑边）
	CanvasClassic  = "classic"  // 4:3（老投影仪）
)

var canvasPresets = map[string][2]int{
	CanvasStandard: {960, 700},
	CanvasWide:     {1244, 700}, // 700 × 16/9
	CanvasClassic:  {933, 700},  // 700 × 4/3
}

// CanvasSize 返回预设对应的逻辑画布尺寸。未识别的预设退回默认（不报错：
// 这一步在读路径上，坏配置不该让整个 deck 读不出来）。
func CanvasSize(preset string) (int, int) {
	if s, ok := canvasPresets[preset]; ok {
		return s[0], s[1]
	}
	s := canvasPresets[CanvasStandard]
	return s[0], s[1]
}

// 自定义调色板的容量上限：既是块大小的护栏，也是 token 上限。
const maxThemeVars = 32

// 默认值必须与 web/assets/theme.css 的 :root 一致（两处耦合，改一处必改另一处，
// 骨架 deck_skeleton.html 里的默认 JSON 也存了一份，共三处）
func defaultTheme() Theme {
	return Theme{
		Accent: "#5eead4", Background: "#0b132b",
		HeadingColor: "#5eead4", TextColor: "#e2e8f0",
		Font: "sans", Radius: "10px", Transition: "slide",
		Canvas: CanvasStandard,
	}
}

var (
	hexColorPattern = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
	radiusPattern   = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?(px|em|rem)$`)
	fontStacks      = map[string]string{
		"sans":  `"Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif`,
		"serif": `Georgia, "Times New Roman", "Songti SC", SimSun, serif`,
		"mono":  `Consolas, "JetBrains Mono", "Courier New", monospace`,
	}
	transitions = map[string]bool{
		"slide": true, "fade": true, "zoom": true,
		"convex": true, "concave": true, "none": true,
	}
)

// validate 对"存储值整体"校验而不是只校验本次 patch——
// 存储块可能被手改或损坏，写入前的全量校验是最后一道闸门。
// 和 update_slide "先指纹后内容"是同一思想：不信任任何进过存储的值。
func (t Theme) validate() error {
	for _, p := range []struct{ name, v string }{
		{"accent", t.Accent}, {"background", t.Background},
		{"heading_color", t.HeadingColor}, {"text_color", t.TextColor},
	} {
		if !hexColorPattern.MatchString(p.v) {
			return fmt.Errorf("%s 必须是 6 位十六进制颜色（如 #5eead4），收到 %q", p.name, p.v)
		}
	}
	if _, ok := fontStacks[t.Font]; !ok {
		return fmt.Errorf("font 只支持 sans/serif/mono，收到 %q", t.Font)
	}
	if !radiusPattern.MatchString(t.Radius) {
		return fmt.Errorf("radius 必须是带单位的长度值（如 10px），收到 %q", t.Radius)
	}
	if !transitions[t.Transition] {
		return fmt.Errorf("transition 只支持 slide/fade/zoom/convex/concave/none，收到 %q", t.Transition)
	}
	if _, ok := canvasPresets[t.Canvas]; !ok {
		return fmt.Errorf("canvas 只支持 standard（960×700，通用）/ wide（16:9，投屏录屏无黑边）/ classic（4:3，老投影仪），收到 %q", t.Canvas)
	}
	if err := validateThemeVars(t.Vars); err != nil {
		return err
	}
	return nil
}

// hexToRGBA 把 #rrggbb 转成 rgba()。调用前必须先过 validate，
// 所以这里忽略解析错误是安全的（validate-then-parse 约定）
func hexToRGBA(hex string, alpha float64) string {
	n, _ := strconv.ParseUint(hex[1:], 16, 32)
	return fmt.Sprintf("rgba(%d, %d, %d, %g)", n>>16&0xff, n>>8&0xff, n&0xff, alpha)
}

// renderThemeCSS 把 Theme 渲染成 :root 变量文本（不含 <style> 标签本身——
// 标签由骨架持有，写入时 SetText 只填内容）。派生变量在这里集中计算：
// 旋钮越少，LLM 的选择面越小，观感越不容易崩。
func renderThemeCSS(t Theme) string {
	stack := fontStacks[t.Font]
	var sb strings.Builder
	sb.WriteString(":root{")
	fmt.Fprintf(&sb, "--r-background-color:%s;", t.Background)
	fmt.Fprintf(&sb, "--r-main-color:%s;", t.TextColor)
	fmt.Fprintf(&sb, "--r-heading-color:%s;", t.HeadingColor)
	fmt.Fprintf(&sb, "--r-link-color:%s;", t.Accent)
	fmt.Fprintf(&sb, "--r-main-font:%s;--r-heading-font:%s;", stack, stack)
	fmt.Fprintf(&sb, "--accent:%s;", t.Accent)
	fmt.Fprintf(&sb, "--border:%s;", hexToRGBA(t.Accent, 0.25))
	fmt.Fprintf(&sb, "--card-bg:%s;", hexToRGBA(t.Accent, 0.06))
	fmt.Fprintf(&sb, "--text-muted:%s;", hexToRGBA(t.TextColor, 0.65))
	fmt.Fprintf(&sb, "--radius:%s;", t.Radius)

	if len(t.Vars) > 0 {
		names := make([]string, 0, len(t.Vars))
		for n := range t.Vars {
			names = append(names, n)
		}
		// 必须排序：map 遍历顺序是随机的，不排的话每次写入产出的 :root 文本都在变，
		// 历史 diff 会把整块报成"改过"——一堆假改动比没有 diff 更糟
		sort.Strings(names)
		for _, n := range names {
			fmt.Fprintf(&sb, "%s:%s;", n, t.Vars[n])
		}
	}

	sb.WriteString("}")
	return sb.String()
}

// ThemePatch 是工具层的可选字段集合：nil 表示"不改这项"。
// 用指针而不是空串，才能区分"没传"和"传了空值"。
//
// Vars 的语义与其它字段不同：它是**整体替换**，不是合并。
// 传了就是"这就是新的整套自定义调色板"（传空 map 等于清空），没传就不动。
// 选整体替换而不是合并的理由和自定义样式槽一样：一套配色是一个整体，
// "当前这套配色长什么样"必须只有一个权威来源，不能靠增量拼接累积出意外。
// 代价是改单个变量也要提交全文——所以配套了 read_theme（先读后写），
// 和 read_slide/read_custom_css 是同一套路。
type ThemePatch struct {
	Accent, Background, HeadingColor, TextColor *string
	Font, Radius, Transition, Canvas            *string
	Vars                                        *map[string]string
}

func (p ThemePatch) applyTo(t *Theme) {
	if p.Accent != nil {
		t.Accent = strings.ToLower(strings.TrimSpace(*p.Accent))
	}
	if p.Background != nil {
		t.Background = strings.ToLower(strings.TrimSpace(*p.Background))
	}
	if p.HeadingColor != nil {
		t.HeadingColor = strings.ToLower(strings.TrimSpace(*p.HeadingColor))
	}
	if p.TextColor != nil {
		t.TextColor = strings.ToLower(strings.TrimSpace(*p.TextColor))
	}
	if p.Font != nil {
		t.Font = strings.ToLower(strings.TrimSpace(*p.Font))
	}
	if p.Radius != nil {
		t.Radius = strings.TrimSpace(*p.Radius)
	}
	if p.Transition != nil {
		t.Transition = strings.ToLower(strings.TrimSpace(*p.Transition))
	}
	if p.Canvas != nil {
		t.Canvas = strings.ToLower(strings.TrimSpace(*p.Canvas))
	}
	if p.Vars != nil {
		// 整体替换。拷贝一份再存：调用方的 map 不该被存储层持有
		// （裸赋值会让之后对入参的修改穿透到主题里）
		if len(*p.Vars) == 0 {
			t.Vars = nil // 空 map = 清空自定义调色板，别留下一个空壳
		} else {
			next := make(map[string]string, len(*p.Vars))
			for k, v := range *p.Vars {
				next[strings.TrimSpace(k)] = strings.TrimSpace(v)
			}
			t.Vars = next
		}
	}
}

// setRawText 用裸文本替换元素内容。必须绕过 goquery 的 SetText——它内部走
// "转义成 HTML 片段再解析插入"，而 script/style 是 raw-text 元素，片段解析
// 在 raw-text 模式下不还原实体：JSON 的引号会被写成 &#34; 字面量（实测确认），
// 读回即废。直接清空子节点、追加 TextNode 没有这个问题。
func setRawText(sel *goquery.Selection, text string) {
	n := sel.Nodes[0]
	for c := n.FirstChild; c != nil; {
		next := c.NextSibling
		n.RemoveChild(c)
		c = next
	}
	if text != "" {
		n.AppendChild(&html.Node{Type: html.TextNode, Data: text})
	}
}

// ensureThemeBlocks 保证主题块存在并解析出当前配置。
// 骨架改造前的老 deck（deck-0001/0002）没有这两个块——碰到时补上默认值，
// 这是最轻量的 schema 迁移：升级动作放在读写路径上，不写独立迁移脚本。
// 配置损坏时退回默认值兜底，缺字段时逐项补默认，不让坏块卡死主题功能。
func ensureThemeBlocks(doc *goquery.Document) (Theme, error) {
	js := doc.Find(`#deck-theme`)
	if js.Length() == 0 {
		def, _ := json.Marshal(defaultTheme())
		doc.Find("head").AppendHtml(
			`<script type="application/json" id="deck-theme">` + string(def) + `</script>`)
		js = doc.Find(`#deck-theme`)
	}
	if doc.Find(`#deck-theme-override`).Length() == 0 {
		// 级联顺序不变量：override 必须排在自定义样式槽之前（否则自定义样式会被
		// 主题派生变量盖掉）。槽可能已经被 ensureCustomCSSBlock 补出来，所以
		// 有槽时插在它前面，而不是无脑 append 到 head 末尾。
		if custom := doc.Find("#" + customCSSBlockID); custom.Length() > 0 {
			custom.First().BeforeHtml(`<style id="deck-theme-override"></style>`)
		} else {
			doc.Find("head").AppendHtml(`<style id="deck-theme-override"></style>`)
		}
	}

	theme := defaultTheme()
	var parsed Theme
	if err := json.Unmarshal([]byte(strings.TrimSpace(js.First().Text())), &parsed); err == nil {
		theme = parsed
	}
	fillMissing(&theme)
	return theme, nil
}

func fillMissing(t *Theme) {
	def := defaultTheme()
	if t.Accent == "" {
		t.Accent = def.Accent
	}
	if t.Background == "" {
		t.Background = def.Background
	}
	if t.HeadingColor == "" {
		t.HeadingColor = def.HeadingColor
	}
	if t.TextColor == "" {
		t.TextColor = def.TextColor
	}
	if t.Font == "" {
		t.Font = def.Font
	}
	if t.Radius == "" {
		t.Radius = def.Radius
	}
	if t.Transition == "" {
		t.Transition = def.Transition
	}
	if t.Canvas == "" {
		t.Canvas = def.Canvas
	}
}

// UpdateTheme 读当前配置 → 合并 patch → 全量校验 → 重渲染两块 → 原子写回。
// 两个主题块由这一条路径一起写：JSON 是真身，CSS 是派生品，永不分家。
// 主题编辑不做指纹校验：字段少、冲突代价低，是相对 update_slide 的刻意简化。
// 主题块在 head、不在任何 section 内，所以已读页的指纹不失效——
// 内容编辑和主题编辑两个维度互不干扰。
func (s *Service) UpdateTheme(userID uint, deckID string, patch ThemePatch) (Theme, error) {
	if err := s.authorize(userID, deckID); err != nil {
		return Theme{}, err
	}
	return s.updateThemeLocked(deckID, patch)
}

// updateThemeLocked 写入核心（锁内），拆出来为白盒测试——与 updateCustomCSSLocked 同一惯例：
// "读→合并→全量校验→重渲染→写回"这条链路本身才是最该被测的东西，
// 而它藏在 authorize 后面，不拆开就只能靠一个需要数据库的集成测试覆盖。
func (s *Service) updateThemeLocked(deckID string, patch ThemePatch) (Theme, error) {
	unlock := s.lockDeck(deckID)
	defer unlock()
	raw, err := s.readRaw(deckID)
	if err != nil {
		return Theme{}, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(raw))
	if err != nil {
		return Theme{}, fmt.Errorf("解析 deck 失败 err:%w", err)
	}

	theme, err := ensureThemeBlocks(doc)
	if err != nil {
		return Theme{}, err
	}

	patch.applyTo(&theme)
	if err := theme.validate(); err != nil {
		return Theme{}, err
	}

	cfg, err := json.Marshal(theme)
	if err != nil {
		return Theme{}, err
	}
	setRawText(doc.Find(`#deck-theme`).First(), string(cfg))
	setRawText(doc.Find(`#deck-theme-override`).First(), renderThemeCSS(theme))

	out, err := goquery.OuterHtml(doc.Selection)
	if err != nil {
		return Theme{}, fmt.Errorf("序列化 deck 失败 err:%w", err)
	}
	if err := s.atomicWriteDeck(deckID, out); err != nil {
		return Theme{}, fmt.Errorf("写入 deck 失败 err:%w", err)
	}
	return theme, nil
}

// ReadTheme 读当前主题配置（工具 read_theme 用）。
//
// 为什么必须有这个读工具：Vars 是整体替换制，AI 改配色前必须先看到当前整套，
// 否则会拿一份想象中的旧配色覆盖真实配色——和 read_custom_css 存在的理由完全一样。
// 读路径只解析不落盘：ensureThemeBlocks 对老 deck 会往内存里的 doc 补块，
// 但这里不写回文件，所以"读"永远是安全的。
func (s *Service) ReadTheme(userID uint, deckID string) (Theme, error) {
	if err := s.authorize(userID, deckID); err != nil {
		return Theme{}, err
	}
	return s.readThemeRaw(deckID)
}

// readThemeRaw 拆出来供白盒测试用（绕开 authorize，st=nil 也能跑）。
func (s *Service) readThemeRaw(deckID string) (Theme, error) {
	raw, err := s.readRaw(deckID)
	if err != nil {
		return Theme{}, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(raw))
	if err != nil {
		return Theme{}, fmt.Errorf("解析 deck 失败 err:%w", err)
	}
	return ensureThemeBlocks(doc)
}
