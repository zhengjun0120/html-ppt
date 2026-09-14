package deck

import (
	"encoding/json"
	"fmt"
	"math"
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
//
// Surface/BorderColor 为什么是字段而不是继续派生：原来的派生
// （--card-bg = accent@6%、--border = accent@25%）只在"暗底 + 单一强调色"这套
// 视觉里成立。换成浅色纸面之后，accent@6% 会得到一层粉、accent@25% 会得到一条
// 彩线——而纸感/编辑风要的恰恰是"中性面板 + 墨色发丝线"。派生关系的根基是
// 视觉语言，不是颜色值，所以它不该被写死成一个公式。留空仍然派生，
// 老 deck 与科技暗色主题的观感一字不变。
type Theme struct {
	Preset       string            `json:"preset,omitempty"`        // 最后一次套用的风格预设名，见 theme_presets.go
	Accent       string            `json:"accent"`                  // 强调色，边框/面板底色默认由它派生
	Background   string            `json:"background"`              // 页面背景
	HeadingColor string            `json:"heading_color"`           // 标题色
	TextColor    string            `json:"text_color"`              // 正文色
	Surface      string            `json:"surface,omitempty"`       // 面板/卡片底色；空 = 由 accent 派生
	BorderColor  string            `json:"border_color,omitempty"`  // 边框与分隔线色；空 = 由 accent 派生
	Font         string            `json:"font"`                    // 字体配对枚举，见 fontPairs
	Radius       string            `json:"radius"`                  // 卡片圆角
	Texture      string            `json:"texture"`                 // 页面纹理枚举：none/grid/dots/rule
	Transition   string            `json:"transition"`              // 翻页动画枚举，init.js 读走
	Canvas       string            `json:"canvas"`                  // 画布比例预设，init.js 读走
	Vars         map[string]string `json:"vars,omitempty"`          // 自定义调色板（变量名 → CSS 值）
}

// canvasPresets 是画布比例预设：**高度统一 700，只调宽度**。
//
// 为什么固定高度：一页能放多少内容由高度决定，宽度只影响排布的宽松度。
// 高度一动，原本刚好放得下的页面就会被挤爆（触发适配兜底缩小，牺牲可读性）；
// 固定高度则"换比例"是纯粹的横向增益——所以换 wide 不会让任何现有页面掉进兜底。
//
// 本表必须与 init.js 的 CANVAS 表逐字一致：前端要在 Reveal.initialize 之前拿到尺寸，
// 没法从后端要。`TestCanvasPresetsMatchInitJS` 守两处不漂移。
const (
	CanvasStandard = "standard" // 960×700（1.37:1，在 16:9 屏上左右有黑边）
	CanvasWide     = "wide"     // 1244×700（16:9，比例与投屏/录屏一致，没有黑边）
	CanvasClassic  = "classic"  // 933×700（4:3，老投影仪）
)

// DefaultCanvas 是新建 deck 的画布预设。
//
// 为什么不用 reveal 自己的 960×700（standard）：字号抬到投影可读下限（正文 29pt）之后，
// 一页能放多少字由**可用宽度 × 高度**决定，而 standard 在 16:9 的屏上左右还要留黑边——
// 实打实用于排版的宽度只有 16:9 画布的 77%。换成 wide 一次解决两件事：
//   - 横向多出约 30% 的版面：双栏每栏从 10 字/行变成 12~13 字/行，三列从"只放数字"
//     变成"放得下 7 个字的短标签"
//   - 投屏与录屏时不再有黑边（那两块黑边既是浪费，也一眼看得出没做适配）
//
// 高度仍是 700，所以"一页能放几行"一个字都没变——这是纯粹的横向增益。
// 改这个常量要同时改 init.js 的兜底表达式，有 TestCanvasPresetsMatchInitJS 守着。
const DefaultCanvas = CanvasWide

var canvasPresets = map[string][2]int{
	CanvasStandard: {960, 700},
	CanvasWide:     {1244, 700}, // 700 × 16/9
	CanvasClassic:  {933, 700},  // 700 × 4/3
}

// CanvasSize 返回预设对应的逻辑画布尺寸。未识别的预设退回**默认画布**（不报错：
// 这一步在读路径上，坏配置不该让整个 deck 读不出来；退回默认而不是退回 standard，
// 是为了让"读不出来"和"没配过"表现一致）。
func CanvasSize(preset string) (int, int) {
	if s, ok := canvasPresets[preset]; ok {
		return s[0], s[1]
	}
	s := canvasPresets[DefaultCanvas]
	return s[0], s[1]
}

// 自定义调色板的容量上限：既是块大小的护栏，也是 token 上限。
const maxThemeVars = 32

// ============ 预设 ============
//
// 为什么要有"有名字的整套美学"，而不是只有自由配色：
// AI 在自由发挥时**必然**回到训练分布的中位数（深底 + 发光强调色 + 圆角卡片 +
// 无衬线字体），这是统计规律、不是态度问题——它见过的绝大多数"现代网站"就长这样。
// 所以"提高自由度"治不了同质化，**给出具体的参照物**才能。
//
// 预设是**写路径上的糖**：套用一次就把整套观感字段落进主题 JSON，渲染路径完全不知道
// 预设的存在（只认存储的字段）。刻意不做"存储 preset 名、渲染时现查表"那种活引用——
// 那意味着改一次预设定义会静默改掉所有用它的 deck，正是这个代码库最不想看到的一类
// 假成功（改了别处、这里跟着变，且没有任何痕迹）。
const (
	TextureNone = "none" // 无纹理
	TextureGrid = "grid" // 细网格：稿纸/工程纸
	TextureDots = "dots" // 网点：印刷/漫画
	TextureRule = "rule" // 横线：稿纸/终端扫描线
)

const (
	FontSans      = "sans"      // 无衬线（黑体全家）——最"不出错"因而也最没有性格的一档
	FontSerif     = "serif"     // 全衬线（宋体/Georgia）
	FontEditorial = "editorial" // 衬线标题 + 无衬线正文：中文杂志的经典配对
	FontModern    = "modern"    // 几何无衬线（Helvetica/PingFang）
	FontMono      = "mono"      // 等宽
)

// 默认值必须与 web/assets/theme.css 的 :root 一致（两处耦合，改一处必改另一处）。
// 骨架里的初始主题 JSON 由 renderSkeleton 从这里渲染，不再是第三处拷贝。
//
// 默认是"纸感"而不是深底发光色：默认值决定大多数结果——用户多数时候不会说"换个配色"，
// 于是默认那一套就是这份工具的观感。深底 + 发光强调色 + 圆角卡片是最容易撞衫的一类，
// 把它降级成可选预设（tech），而不是让它当门面。
//
// Vars 从 paper 预设拷过来，而不是留空：这份主题标着 `Preset: PresetPaper`，
// 而纸感预设的 Vars 里就有那三个语义色（--accent-2/--positive/--warn）。
// 留空等于"自称是纸感，却少了纸感的一部分"——具体后果不是观感问题而是**拒收**：
// 回声校验只认"deck 自己的样式块定义过"的变量，缺了这三个，模型照着提示词
// （"预设已经预置了它们，直接 var() 引用，不要重复定义"）写下来就会被拒，
// 而提示词又明确叫它别自己定义，等于把模型送进死胡同。
func defaultTheme() Theme {
	return Theme{
		Preset:       PresetPaper,
		Accent:       "#b23a2e", Background: "#f6f2e9",
		HeadingColor: "#1c1a17", TextColor: "#2b2723",
		Surface:      "#fffdf8", BorderColor: "#ded4c3",
		Font:         FontEditorial, Radius: "2px", Texture: TextureGrid,
		Transition:   "slide", Canvas: DefaultCanvas,
		Vars: copyVars(presets[PresetPaper].Theme.Vars),
	}
}

var (
	hexColorPattern = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
	radiusPattern   = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?(px|em|rem)$`)
	// fontPairs 把"字体"从三档系统栈升级成**配对**：标题和正文用同一套栈，
	// 是"没有做过字体决策"的默认结果；而"衬线标题 + 无衬线正文"这一对，
	// 光是换上去就能让页面从"AI 生成的网页"变成"有人排过版的读物"。
	//
	// 五对里只有 mono 这一对引用了**随附字体**（web/assets/fonts、theme.css 里的
	// @font-face）。为什么只在这一档破例：拉丁等宽字体到处都有，中日文等宽字体却没有
	// ——Consolas 的"等宽"只覆盖拉丁，中文落到系统黑体后步进宽度和拉丁对不上，
	// 于是这一档的整个观感取决于机器上恰好装了什么。中日文等宽是"必须自带才成立"的
	// 那一档；其余四对靠系统字体栈就能得到可预期的结果，不值得各背一份几 MB 的字体。
	fontPairs = map[string][2]string{
		FontSans:  {`"Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif`, `"Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif`},
		FontSerif: {`Georgia, "Times New Roman", "Songti SC", SimSun, serif`, `Georgia, "Times New Roman", "Songti SC", SimSun, serif`},
		FontEditorial: {
			`"Songti SC", "SimSun", Georgia, "Times New Roman", serif`,
			`"PingFang SC", "Microsoft YaHei", "Segoe UI", sans-serif`,
		},
		FontModern: {`"Helvetica Neue", Helvetica, Arial, "PingFang SC", sans-serif`, `"Helvetica Neue", Helvetica, Arial, "PingFang SC", sans-serif`},
		FontMono:   {`"JetBrains Maple Mono", Consolas, "Courier New", monospace`, `"JetBrains Maple Mono", Consolas, "Courier New", monospace`},
	}
	// fontNames 给报错信息与 schema 契约测试用的稳定顺序清单（map 遍历是随机的）
	fontNames = []string{FontSans, FontSerif, FontEditorial, FontModern, FontMono}
	textures  = map[string]bool{
		TextureNone: true, TextureGrid: true, TextureDots: true, TextureRule: true,
	}
	// textureNames 同上：给 schema 契约测试用的稳定顺序
	textureNames = []string{TextureNone, TextureGrid, TextureDots, TextureRule}
	// transitionNames 同上
	transitionNames = []string{"slide", "fade", "zoom", "convex", "concave", "none"}
	transitions     = map[string]bool{
		"slide": true, "fade": true, "zoom": true,
		"convex": true, "concave": true, "none": true,
	}
)

// FontNames / TextureNames / CanvasNames / TransitionNames 返回稳定顺序的枚举取值。
// 存在的理由是让工具 schema 的 enum 与这里的表有一条测试守住：
// enum 只能手写在 struct tag 里（jsonschema 库不支持动态枚举），少了这条测试，
// 加一档取值就会出现"代码认识、schema 不认"的静默拒绝（模型写得出、传不进来）。
func FontNames() []string       { return append([]string(nil), fontNames...) }
func TextureNames() []string    { return append([]string(nil), textureNames...) }
func TransitionNames() []string { return append([]string(nil), transitionNames...) }
func CanvasNames() []string     { return []string{CanvasStandard, CanvasWide, CanvasClassic} }

// validate 对"存储值整体"校验而不是只校验本次 patch——
// 存储块可能被手改或损坏，写入前的全量校验是最后一道闸门。
// 和 update_slide "先指纹后内容"是同一思想：不信任任何进过存储的值。
func (t Theme) validate() error {
	for _, p := range []struct{ name, v string }{
		{"accent", t.Accent}, {"background", t.Background},
		{"heading_color", t.HeadingColor}, {"text_color", t.TextColor},
	} {
		if !hexColorPattern.MatchString(p.v) {
			return fmt.Errorf("%s 必须是 6 位十六进制颜色（如 #b23a2e），收到 %q", p.name, p.v)
		}
	}
	// surface / border_color 是可选覆盖：留空表示"由 accent 派生"，不校验空值
	for _, p := range []struct{ name, v string }{
		{"surface", t.Surface}, {"border_color", t.BorderColor},
	} {
		if p.v != "" && !hexColorPattern.MatchString(p.v) {
			return fmt.Errorf("%s 必须是 6 位十六进制颜色（如 #fffdf8），留空表示由 accent 派生，收到 %q", p.name, p.v)
		}
	}
	if _, ok := fontPairs[t.Font]; !ok {
		return fmt.Errorf("font 只支持 %s，收到 %q", strings.Join(fontNames, "/"), t.Font)
	}
	if !radiusPattern.MatchString(t.Radius) {
		return fmt.Errorf("radius 必须是带单位的长度值（如 2px；0px = 直角），收到 %q", t.Radius)
	}
	if !textures[t.Texture] {
		return fmt.Errorf("texture 只支持 none/grid（细网格）/dots（网点）/rule（横线），收到 %q", t.Texture)
	}
	if t.Preset != "" {
		if _, ok := presets[t.Preset]; !ok {
			return fmt.Errorf("preset 只支持 %s，收到 %q", strings.Join(PresetNames(), "/"), t.Preset)
		}
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

// onAccentColor 选一个"压在强调色上读得清"的前景色（近黑或白）。
//
// 为什么不能靠"暗色主题用深字、浅色主题用白字"来判断：强调色和背景色是两件独立的事，
// 深底主题完全可能配一个亮黄强调（那时要深字），浅底主题也可能配一个深靛强调（那时要白字）。
// 所以要按强调色自身的亮度算。用 WCAG 的相对亮度公式（通道先线性化）再比两种前景的
// 对比度，取高的那个——正好在"两个都不够好"的中间色上也能选到相对能用的那个。
//
// 原来的做法是组件库里写死 color:#0b132b。它只在"强调色都是亮色"时成立：
// 编辑风预设的强调色是黑，写死就成黑底黑字。
func onAccentColor(hex string) string {
	lin := func(c float64) float64 {
		if c <= 0.03928 {
			return c / 12.92
		}
		return math.Pow((c+0.055)/1.055, 2.4)
	}
	n, _ := strconv.ParseUint(hex[1:], 16, 32)
	r := lin(float64(n>>16&0xff) / 255)
	g := lin(float64(n>>8&0xff) / 255)
	b := lin(float64(n&0xff) / 255)
	lum := 0.2126*r + 0.7152*g + 0.0722*b

	contrastWithWhite := 1.05 / (lum + 0.05)
	contrastWithBlack := (lum + 0.05) / 0.05
	if contrastWithBlack >= contrastWithWhite {
		return "#12100e"
	}
	return "#ffffff"
}

// renderThemeCSS 把 Theme 渲染成 `:root{}` 变量文本 +（可选）纹理规则
// （不含 <style> 标签本身——标签由骨架持有，写入时 SetRawText 只填内容）。
// 派生变量在这里集中计算：旋钮越少，LLM 的选择面越小，观感越不容易崩。
//
// 除变量外还渲染两条**非** :root 的规则，它们都是"整个画布的默认值"：
//   - 纹理：背景图案挂在 .reveal-viewport（moon.css 就是往它身上写背景色的，
//     同一元素、同一优先级，排在后面即生效）。不用伪元素是因为 reveal 会改
//     section 的 transform，装饰挂在被变换的节点上会跟着翻页抖动。
//   - 进度条颜色：reveal.css 写死 color:#fff，在浅色主题上等于看不见。
//     它属于"跟主题走"的东西，由主题给出，而不是留在第三方 CSS 里。
func renderThemeCSS(t Theme) string {
	pair := fontPairs[t.Font]
	var sb strings.Builder
	sb.WriteString(":root{")
	fmt.Fprintf(&sb, "--r-background-color:%s;", t.Background)
	fmt.Fprintf(&sb, "--r-main-color:%s;", t.TextColor)
	fmt.Fprintf(&sb, "--r-heading-color:%s;", t.HeadingColor)
	fmt.Fprintf(&sb, "--r-link-color:%s;", t.Accent)
	fmt.Fprintf(&sb, "--r-main-font:%s;--r-heading-font:%s;", pair[1], pair[0])
	fmt.Fprintf(&sb, "--accent:%s;", t.Accent)
	fmt.Fprintf(&sb, "--border:%s;", derivedColor(t.BorderColor, t.Accent, 0.25))
	fmt.Fprintf(&sb, "--card-bg:%s;", derivedColor(t.Surface, t.Accent, 0.06))
	// --accent-soft 给 .band 这类"整条色带"用：色带本来就该是强调色的一层淡染，
	// 和"卡片该不该是中性面板"是两个不同的问题，所以它始终保持派生
	fmt.Fprintf(&sb, "--accent-soft:%s;", hexToRGBA(t.Accent, 0.14))
	// --on-accent 是"压在强调色上的字色"，由强调色亮度算出（见 onAccentColor）
	fmt.Fprintf(&sb, "--on-accent:%s;", onAccentColor(t.Accent))
	fmt.Fprintf(&sb, "--text-muted:%s;", hexToRGBA(t.TextColor, 0.65))
	// --hairline 是组内分隔线（清单行、旁注顶线），由**正文色**派生而不是强调色：
	// 它要在深浅两种主题上都只是"一根几乎看不见的线"，跟着强调色走会在浅色主题上
	// 变成一条彩线（原来 .row 用虚线画分隔，正是这个问题最显眼的地方）
	fmt.Fprintf(&sb, "--hairline:%s;", hexToRGBA(t.TextColor, 0.12))
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
	sb.WriteString(textureCSS(t))
	sb.WriteString(".reveal .progress{color:var(--accent)}")
	return sb.String()
}

// derivedColor 二选一：给了显式值就用它，留空才由强调色派生。
// 派生值只对"暗底 + 单一强调色"成立，所以它是默认而不是规则。
func derivedColor(explicit, accent string, alpha float64) string {
	if explicit != "" {
		return explicit
	}
	return hexToRGBA(accent, alpha)
}

// textureCSS 把纹理枚举展开成一条背景图案规则。
// 线条颜色由正文色派生而不是固定灰：深色主题上浅、浅色主题上深，
// 同一个枚举在两种底色上都成立。
func textureCSS(t Theme) string {
	if t.Texture == "" || t.Texture == TextureNone {
		return ""
	}
	line := hexToRGBA(t.TextColor, 0.07)
	switch t.Texture {
	case TextureGrid:
		return fmt.Sprintf(".reveal-viewport{background-image:"+
			"linear-gradient(%s 1px,transparent 1px),"+
			"linear-gradient(90deg,%s 1px,transparent 1px);background-size:30px 30px}", line, line)
	case TextureDots:
		return fmt.Sprintf(".reveal-viewport{background-image:"+
			"radial-gradient(%s 1.2px,transparent 1.2px);background-size:22px 22px}", line)
	case TextureRule:
		return fmt.Sprintf(".reveal-viewport{background-image:"+
			"repeating-linear-gradient(to bottom,transparent 0 33px,%s 33px 34px)}", line)
	}
	return ""
}

// ThemePatch 是工具层的可选字段集合：nil 表示"不改这项"。
// 用指针而不是空串，才能区分"没传"和"传了空值"。
//
// Preset 的语义是"整套观感赋值"：命中预设就先把它的整套字段铺进去，
// 同一调用里显式传的其它字段再覆盖它（"用纸感，但主色换成靛蓝"是一个调用的事）。
// 它不动 transition/canvas——那两个是演示参数，不是视觉语言。
//
// Vars 的语义与其它字段不同：它是**整体替换**，不是合并。
// 传了就是"这就是新的整套自定义调色板"（传空 map 等于清空），没传就不动。
// 选整体替换而不是合并的理由和自定义样式槽一样：一套配色是一个整体，
// "当前这套配色长什么样"必须只有一个权威来源，不能靠增量拼接累积出意外。
// 代价是改单个变量也要提交全文——所以配套了 read_theme（先读后写），
// 和 read_slide/read_custom_css 是同一套路。
type ThemePatch struct {
	Preset                                      *string
	Accent, Background, HeadingColor, TextColor *string
	Surface, BorderColor                        *string
	Font, Radius, Texture, Transition, Canvas   *string
	Vars                                        *map[string]string
}

// aestheticFields 返回"改过就算偏离预设"的字段名（非空即表示这一项被单独动过）。
// 判据是"会不会改变这份 deck 的视觉语言"：transition/canvas 是演示参数
// （用户说"用纸感，翻页改成淡入"，这份 deck 仍然是纸感），所以不算。
func (p ThemePatch) aestheticFields() []string {
	var out []string
	for _, f := range []struct {
		name string
		set  bool
	}{
		{"accent", p.Accent != nil}, {"background", p.Background != nil},
		{"heading_color", p.HeadingColor != nil}, {"text_color", p.TextColor != nil},
		{"surface", p.Surface != nil}, {"border_color", p.BorderColor != nil},
		{"font", p.Font != nil}, {"radius", p.Radius != nil}, {"texture", p.Texture != nil},
		{"vars", p.Vars != nil},
	} {
		if f.set {
			out = append(out, f.name)
		}
	}
	return out
}

// applyTo 返回 error 而不是静默忽略：预设名写错必须报错。
// 静默忽略会产出一份"报了成功、观感一动不动"的结果——正是这个文件反复在防的那类失败。
func (p ThemePatch) applyTo(t *Theme) error {
	if p.Preset != nil {
		preset, ok := presets[*p.Preset]
		if !ok {
			return fmt.Errorf("没有名为 %q 的预设。可用预设：%s", *p.Preset,
				strings.Join(PresetNames(), " / "))
		}
		preset.applyTo(t)
	}
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
	if p.Surface != nil {
		t.Surface = strings.ToLower(strings.TrimSpace(*p.Surface))
	}
	if p.BorderColor != nil {
		t.BorderColor = strings.ToLower(strings.TrimSpace(*p.BorderColor))
	}
	if p.Font != nil {
		t.Font = strings.ToLower(strings.TrimSpace(*p.Font))
	}
	if p.Radius != nil {
		t.Radius = strings.TrimSpace(*p.Radius)
	}
	if p.Texture != nil {
		t.Texture = strings.ToLower(strings.TrimSpace(*p.Texture))
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
	// 这份 deck 已经不等于某个预设的原样了，就不要再自称是那个预设——
	// 一个会说谎的字段比没有这个字段更糟（模型会据此向用户断言"当前是纸感风格"）。
	// 同一次调用里既选预设又微调的情况保留预设名：那种情况下预设是它的底子。
	if p.Preset == nil {
		if touched := p.aestheticFields(); len(touched) > 0 {
			t.Preset = ""
		}
	}
	return nil
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

// fillMissing 逐项补默认值（老 deck 的配置块里没有新字段）。
// Surface/BorderColor 刻意不补：空值在这里是**有语义的**（= 由 accent 派生），
// 补成默认值会让老 deck 一被读到就换掉观感——而这正是"默认派生"要保住的东西。
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
	if t.Texture == "" {
		t.Texture = TextureNone
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

	if err := patch.applyTo(&theme); err != nil {
		return Theme{}, err
	}
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
