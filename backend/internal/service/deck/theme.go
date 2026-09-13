package deck

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// Theme 是 deck 的语义化主题配置，真身是 deck.html 里的
// <script type="application/json" id="deck-theme"> 块。
// 结构层 components.css 只引用 CSS 变量；本结构决定变量的取值。
// 暴露给 LLM 的只有这些语义字段，不是 CSS 原文——校验做硬、可逆、出错率低。
type Theme struct {
	Accent       string `json:"accent"`        // 强调色，border/card-bg 由它派生
	Background   string `json:"background"`    // 页面背景
	HeadingColor string `json:"heading_color"` // 标题色
	TextColor    string `json:"text_color"`    // 正文色
	Font         string `json:"font"`          // 字体方案枚举：sans/serif/mono
	Radius       string `json:"radius"`        // 卡片圆角
	Transition   string `json:"transition"`    // 翻页动画枚举，init.js 读走
}

// 默认值必须与 web/assets/theme.css 的 :root 一致（两处耦合，改一处必改另一处，
// 骨架 deck_skeleton.html 里的默认 JSON 也存了一份，共三处）
func defaultTheme() Theme {
	return Theme{
		Accent: "#5eead4", Background: "#0b132b",
		HeadingColor: "#5eead4", TextColor: "#e2e8f0",
		Font: "sans", Radius: "10px", Transition: "slide",
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
	sb.WriteString("}")
	return sb.String()
}

// ThemePatch 是工具层的可选字段集合：nil 表示"不改这项"。
// 用指针而不是空串，才能区分"没传"和"传了空值"。
type ThemePatch struct {
	Accent, Background, HeadingColor, TextColor *string
	Font, Radius, Transition                    *string
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
