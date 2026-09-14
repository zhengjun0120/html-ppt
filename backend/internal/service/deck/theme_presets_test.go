package deck

// 预设表与"观感字段"的回归测试。这里守四件事：
//
//  1. 预设本身是合法的——它是一张会被整套铺进主题的表，任何一个预设里少写一个 #、
//     或误定义了契约变量，都要到"用户真的选了那个预设"时才炸。表驱动一起测掉。
//  2. 预设是**副本**语义——预设表是包级共享变量，套用时没拷贝就会让"改一个 deck 的 vars"
//     穿透到所有用过这个预设的 deck。
//  3. preset 名不会说谎——被单独改过观感的 deck 必须丢掉预设名，否则模型会照着它
//     向用户断言"当前是纸感风格"，而页面早就不是了。
//  4. 老 deck 的观感不被新默认值改写——没写 surface/border_color 的老 deck 仍然走派生。

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestPresetsAreValid(t *testing.T) {
	names := PresetNames()
	if len(names) != len(Presets()) {
		t.Fatalf("PresetNames 与 Presets 数量不一致：%d vs %d", len(names), len(Presets()))
	}
	if len(names) == 0 {
		t.Fatal("预设表是空的")
	}
	for _, p := range Presets() {
		t.Run(p.Name, func(t *testing.T) {
			if p.Name == "" || p.Label == "" || p.About == "" {
				t.Fatal("预设必须有名字、中文名和一句话说明（模型靠 About 选它）")
			}
			if _, ok := presets[p.Name]; !ok {
				t.Fatalf("Preset %q 不在 presets 表里（Presets() 与表不同步）", p.Name)
			}
			// 从默认主题出发套预设再全量校验：预设只给观感字段，
			// transition/canvas 这些必须由默认值兜住——否则套完预设的 deck 是个半成品
			theme := defaultTheme()
			p.applyTo(&theme)
			if err := theme.validate(); err != nil {
				t.Fatalf("预设 %s 套用后校验失败: %v", p.Name, err)
			}
			if theme.Preset != p.Name {
				t.Fatalf("套用预设后 Preset 字段应是 %q，实际 %q", p.Name, theme.Preset)
			}
			if theme.Transition == "" || theme.Canvas == "" {
				t.Fatalf("预设不该清掉 transition/canvas（那是演示参数）: %+v", theme)
			}
		})
	}
}

// 默认主题必须**就是**纸感预设：默认值决定大多数结果（用户多数时候不会说"换个配色"），
// 所以"默认那一套"和"预设里的一套"只能是同一个东西，不能各写一份。
func TestDefaultThemeIsPaperPreset(t *testing.T) {
	def := defaultTheme()
	paper := presets[PresetPaper]

	if def.Preset != PresetPaper {
		t.Errorf("默认主题应标记为 %s 预设，实际 %q", PresetPaper, def.Preset)
	}
	for _, c := range []struct{ name, got, want string }{
		{"accent", def.Accent, paper.Theme.Accent},
		{"background", def.Background, paper.Theme.Background},
		{"heading_color", def.HeadingColor, paper.Theme.HeadingColor},
		{"text_color", def.TextColor, paper.Theme.TextColor},
		{"surface", def.Surface, paper.Theme.Surface},
		{"border_color", def.BorderColor, paper.Theme.BorderColor},
		{"font", def.Font, paper.Theme.Font},
		{"radius", def.Radius, paper.Theme.Radius},
		{"texture", def.Texture, paper.Theme.Texture},
	} {
		if c.got != c.want {
			t.Errorf("默认主题的 %s 应与 %s 预设一致：默认=%q 预设=%q", c.name, PresetPaper, c.got, c.want)
		}
	}
	// 语义色也要算在内。少了这三个，默认主题就是"自称纸感、却没有纸感那三种语义色"
	// 的半份预设（后果见 defaultTheme 的注释：模型引用它们会被回声校验拒收）。
	if !reflect.DeepEqual(def.Vars, paper.Theme.Vars) {
		t.Errorf("默认主题的语义色应与 %s 预设一致：默认=%v 预设=%v", PresetPaper, def.Vars, paper.Theme.Vars)
	}
}

func TestPresetSummaryListsEveryPreset(t *testing.T) {
	summary := PresetSummary()
	for _, p := range Presets() {
		if !strings.Contains(summary, p.Name) || !strings.Contains(summary, p.Label) {
			t.Errorf("预设清单里缺少 %s（%s）:\n%s", p.Name, p.Label, summary)
		}
	}
	if strings.HasSuffix(summary, "\n") {
		t.Error("预设清单不该以换行结尾（它被拼进工具说明里）")
	}
}

// 套用预设必须拷贝 vars。预设表是包级共享的，裸赋值会让之后对某个 deck 的 vars 修改
// 直接改写这张表——表现是"改了 A deck 的配色，B deck 的配色也跟着变"，
// 而且是进程级的串味，重启前一直在。
func TestPresetApplyCopiesVars(t *testing.T) {
	theme := defaultTheme()
	presets[PresetDuotone].applyTo(&theme)
	if len(theme.Vars) == 0 {
		t.Fatal("双色印刷预设应带语义色，实际没有")
	}
	theme.Vars["--accent-2"] = "#pwned"

	if got := presets[PresetDuotone].Theme.Vars["--accent-2"]; got == "#pwned" {
		t.Fatal("套用预设后改 vars 穿透到了预设表：预设表被当成可变对象了")
	}
}

func TestPresetApplyRejectsUnknownName(t *testing.T) {
	theme := defaultTheme()
	bad := "brutalist"
	err := ThemePatch{Preset: &bad}.applyTo(&theme)
	if err == nil {
		t.Fatal("未知预设名必须报错，不能静默忽略（那会产出'报了成功、页面没动'的结果）")
	}
	// 错误里要给出可用清单，否则模型只能反复试名字
	for _, name := range PresetNames() {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("错误信息里应列出可用预设 %q：%v", name, err)
		}
	}
	if theme.Preset == bad {
		t.Error("未知预设不该被写进主题")
	}
}

// 预设是整套赋值，同一调用里显式传的字段覆盖它；只改 transition/canvas 不算偏离。
func TestPresetPatchMergesExplicitFields(t *testing.T) {
	t.Run("显式字段覆盖预设", func(t *testing.T) {
		theme := defaultTheme()
		preset, accent := PresetNoir, "#334155"
		if err := (ThemePatch{Preset: &preset, Accent: &accent}).applyTo(&theme); err != nil {
			t.Fatal(err)
		}
		if theme.Accent != "#334155" {
			t.Errorf("显式 accent 应覆盖预设，实际 %q", theme.Accent)
		}
		if theme.Background != presets[PresetNoir].Theme.Background {
			t.Errorf("未显式传的字段应来自预设，实际 %q", theme.Background)
		}
		if theme.Preset != PresetNoir {
			t.Errorf("同一次调用里既选预设又微调，预设名应保留（那是它的底子），实际 %q", theme.Preset)
		}
	})

	t.Run("只改翻页动画不算偏离", func(t *testing.T) {
		theme := defaultTheme()
		transition := "fade"
		if err := (ThemePatch{Transition: &transition}).applyTo(&theme); err != nil {
			t.Fatal(err)
		}
		if theme.Preset != PresetPaper {
			t.Errorf("transition/canvas 是演示参数，不该清掉预设名，实际 %q", theme.Preset)
		}
	})

	// 单独改观感字段 → 这份 deck 不再是某个预设的原样，预设名必须清掉：
	// 一个会说谎的字段比没有这个字段更糟
	for _, c := range []struct {
		name  string
		patch ThemePatch
	}{
		{"accent", ThemePatch{Accent: strPtr("#334155")}},
		{"font", ThemePatch{Font: strPtr(FontMono)}},
		{"radius", ThemePatch{Radius: strPtr("18px")}},
		{"texture", ThemePatch{Texture: strPtr(TextureNone)}},
		{"surface", ThemePatch{Surface: strPtr("#101010")}},
	} {
		t.Run("改 "+c.name+" 清掉预设名", func(t *testing.T) {
			theme := defaultTheme()
			if err := c.patch.applyTo(&theme); err != nil {
				t.Fatal(err)
			}
			if theme.Preset != "" {
				t.Errorf("观感被单独改过后不该再自称 %q", theme.Preset)
			}
		})
	}
}

// 老 deck（主题块里没有 surface/border_color）必须继续走派生：
// 新加的两个字段是"可选覆盖"，不是"新的必填项"，否则所有历史 deck 一被写到就换脸。
func TestLegacyDeckStillDerivesSurfaceAndBorder(t *testing.T) {
	legacy := defaultTheme()
	legacy.Surface = ""
	legacy.BorderColor = ""
	legacy.Accent = "#5eead4"

	css := renderThemeCSS(legacy)
	if !strings.Contains(css, "--border:rgba(94, 234, 212, 0.25);") {
		t.Errorf("留空时应由 accent 派生 border: %s", css)
	}
	if !strings.Contains(css, "--card-bg:rgba(94, 234, 212, 0.06);") {
		t.Errorf("留空时应由 accent 派生 card-bg: %s", css)
	}
}

// --on-accent 必须按强调色亮度算，不能按"主题是深是浅"猜：
// 深底主题配亮黄强调时要深字，浅底主题配深靛强调时要白字。
func TestOnAccentColorPicksReadableForeground(t *testing.T) {
	cases := []struct{ accent, want string }{
		{"#5eead4", "#12100e"}, // 亮青绿 → 深字（原来的写死值就是这个方向）
		{"#d9a441", "#12100e"}, // 暖金 → 深字
		{"#111111", "#ffffff"}, // 编辑风的黑强调 → 白字（写死 #0b132b 时这里是黑底黑字）
		{"#12233d", "#ffffff"}, // 深靛 → 白字
		{"#ff5a36", "#12100e"}, // 荧光橙：和黑字对比度更高
	}
	for _, c := range cases {
		if got := onAccentColor(c.accent); got != c.want {
			t.Errorf("%s 上的前景色应为 %s，实际 %s", c.accent, c.want, got)
		}
	}
}

func TestTextureCSS(t *testing.T) {
	theme := defaultTheme()

	theme.Texture = TextureNone
	if got := textureCSS(theme); got != "" {
		t.Errorf("none 不该产出任何规则，实际 %q", got)
	}
	theme.Texture = ""
	if got := textureCSS(theme); got != "" {
		t.Errorf("空值等同于 none，实际 %q", got)
	}

	for _, name := range []string{TextureGrid, TextureDots, TextureRule} {
		theme.Texture = name
		got := textureCSS(theme)
		if got == "" {
			t.Errorf("纹理 %s 没有产出规则", name)
			continue
		}
		// 必须挂在 .reveal-viewport 上：moon.css 就是往它身上写背景色的，
		// 挂到 section 上会跟着 reveal 的 transform 一起动（翻页时装饰抖动）
		if !strings.HasPrefix(got, ".reveal-viewport{") {
			t.Errorf("纹理应挂在 .reveal-viewport 上，实际 %q", got)
		}
		if !strings.Contains(got, "background-image:") {
			t.Errorf("纹理应通过 background-image 实现，实际 %q", got)
		}
		if strings.Contains(got, ";") && !strings.HasSuffix(got, "}") {
			t.Errorf("规则没有收尾: %q", got)
		}
	}

	// 线条颜色由正文色派生：深色主题上浅、浅色主题上深，同一个枚举两种底色都成立
	dark := defaultTheme()
	dark.Texture = TextureGrid
	dark.TextColor = "#e2e8f0"
	light := defaultTheme()
	light.Texture = TextureGrid
	if textureCSS(dark) == textureCSS(light) {
		t.Error("纹理线色应随正文色变化，两种底色的纹理应不同")
	}
}

// 字体从"三档系统栈"升级成配对：editorial 的标题与正文必须是两套栈，
// 否则"衬线标题 + 无衬线正文"这个中文杂志的关键组合根本没落地。
func TestFontPairsRenderBothStacks(t *testing.T) {
	theme := defaultTheme()
	theme.Font = FontEditorial
	css := renderThemeCSS(theme)

	pair := fontPairs[FontEditorial]
	if pair[0] == pair[1] {
		t.Fatal("editorial 预设的标题与正文字体栈应当不同")
	}
	if !strings.Contains(css, "--r-main-font:"+pair[1]+";") {
		t.Errorf("正文栈没渲染: %s", css)
	}
	if !strings.Contains(css, "--r-heading-font:"+pair[0]+";") {
		t.Errorf("标题栈没渲染: %s", css)
	}

	// sans 这类单栈方案两边一致，不该被拆出两种
	same := defaultTheme()
	same.Font = FontSans
	if fontPairs[FontSans][0] != fontPairs[FontSans][1] {
		t.Error("sans 应是同一套栈")
	}
}

// 手工验收页 layouts-test.html 里的观感是 renderThemeCSS 输出的**逐字快照**
// （页面注释就是这么写的：它们是数据快照、不是第二份实现）。快照会漂：改了预设表
// 或 fontPairs，Go 的输出变了、验收页还显示旧观感——于是"看图验收"验的是一个
// 已经不存在的页面，而这件事从截图上看不出来。这条测试把两边钉死。
//
// 判定基准是"页面自己那份 deck-theme JSON"而不是 defaultTheme()：页面 head 里的
// 主题块带 paper 预设的语义色（--accent-2/--positive/--warn），而 defaultTheme()
// 目前**不带** Vars（见 defaultTheme 的定义）。拿 defaultTheme() 当基准会把
// "页面多写了三个变量"报成漂移，而真正要守的是"Go 照着这份 JSON 渲染，会得到这份 CSS"。
func TestLayoutsFixtureMatchesRenderer(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "web", "assets", "layouts-test.html"))
	if err != nil {
		t.Fatalf("读取验收页失败（本包到 web/assets 是三层）: %v", err)
	}
	page := string(raw)

	themeJSON := regexp.MustCompile(`(?s)<script type="application/json" id="deck-theme">(.*?)</script>`).
		FindStringSubmatch(page)
	if themeJSON == nil {
		t.Fatal("验收页里找不到 #deck-theme JSON 块")
	}
	var pageTheme Theme
	if err := json.Unmarshal([]byte(themeJSON[1]), &pageTheme); err != nil {
		t.Fatalf("验收页的主题 JSON 不能解析: %v", err)
	}

	t.Run("head 里的 override 块", func(t *testing.T) {
		m := regexp.MustCompile(`(?s)<style id="deck-theme-override">(.*?)</style>`).FindStringSubmatch(page)
		if m == nil {
			t.Fatal("验收页里找不到 #deck-theme-override 块")
		}
		if want := renderThemeCSS(pageTheme); m[1] != want {
			t.Errorf("override 块与 renderThemeCSS(页面主题) 不一致\n got=%s\nwant=%s", m[1], want)
		}
	})

	t.Run("六个预设的 CSS 快照", func(t *testing.T) {
		got := map[string]string{}
		for _, m := range regexp.MustCompile("(\\w+):\\s*`([^`]*)`").FindAllStringSubmatch(page, -1) {
			got[m[1]] = m[2]
		}
		if len(got) != len(PresetNames()) {
			t.Fatalf("快照数量与预设数量不一致：快照 %d 个，预设 %d 个（正则也可能抓到了别的模板字符串）",
				len(got), len(PresetNames()))
		}
		for _, p := range Presets() {
			snap, ok := got[p.Name]
			if !ok {
				t.Errorf("验收页缺少预设 %s 的快照", p.Name)
				continue
			}
			if want := renderThemeCSS(p.Theme); snap != want {
				t.Errorf("验收页里 %s 的快照过期（改了预设没同步这个页面）\n got=%s\nwant=%s", p.Name, snap, want)
			}
		}
	})
}

// 预设清单出现在工具说明里（每轮都要重发），取值必须来自预设表而不是 prompt 里再抄一份。
func TestPresetVarsAreNotReservedVars(t *testing.T) {
	for _, p := range Presets() {
		for name := range p.Theme.Vars {
			if err := validateVarName(name); err != nil {
				t.Errorf("预设 %s 定义了不该由 vars 定义的变量 %s: %v", p.Name, name, err)
			}
		}
	}
}

// 骨架里的初始主题 JSON 由 Go 渲染（原来是手抄的第三份默认值）。
// 模板占位符写错、或渲染出的 JSON 不能反序列化，都只有真建一份 deck 时才会暴露——
// 而那条路径要数据库，单元测试跑不到。所以这里直接验渲染产物。
func TestRenderSkeletonEmbedsDefaultTheme(t *testing.T) {
	html, err := renderSkeleton("标题 &amp; 测试", `<section><h1>hi</h1></section>`)
	if err != nil {
		t.Fatalf("渲染骨架失败: %v", err)
	}
	start := strings.Index(html, `id="deck-theme">`)
	if start < 0 {
		t.Fatal("骨架里没有 deck-theme 块（模板占位符被改坏或漏了 ThemeJSON 参数）")
	}
	start += len(`id="deck-theme">`)
	end := strings.Index(html[start:], `</script>`)
	if end < 0 {
		t.Fatal("deck-theme 块没有闭合")
	}

	var got Theme
	if err := json.Unmarshal([]byte(html[start:start+end]), &got); err != nil {
		t.Fatalf("骨架里的初始主题 JSON 不能解析：%v\n%s", err, html[start:start+end])
	}
	if !reflect.DeepEqual(got, defaultTheme()) {
		t.Errorf("骨架里的初始主题应与 defaultTheme() 一致：\n got=%+v\nwant=%+v", got, defaultTheme())
	}
	// override 块也必须是渲染出来的，而不是空的：它是"这份 deck 定义了哪些变量"的
	// 权威（回声校验认的就是它），空块等于新建的 deck 少了预设预置的那三个语义色
	css := regexp.MustCompile(`(?s)<style id="deck-theme-override">(.*?)</style>`).FindStringSubmatch(html)
	if css == nil {
		t.Fatal("骨架里没有 deck-theme-override 块")
	}
	if want := renderThemeCSS(defaultTheme()); css[1] != want {
		t.Errorf("骨架的 override 块应是 renderThemeCSS(defaultTheme())：\n got=%s\nwant=%s", css[1], want)
	}
	// 占位符没被替换会留下字面的 {{.ThemeJSON}}
	if strings.Contains(html, "{{") {
		t.Error("骨架里还有未替换的模板占位符")
	}
}

// 提示词无条件地告诉模型："预设已经预置了 --accent-2（第二强调色）/--positive/--warn
// 三个语义色，需要区分正负时直接 var() 引用，不要重复定义它们。"
// 新建的 deck 必须真的能用它们——这条测试守的就是这个承诺。
//
// 为什么值得单独守：这三个变量刻意**不被组件库消费**（没带预设的 deck 上会渲染成
// 不可见），所以它们不在变量契约里，只能靠"deck 自己的样式块定义过"这一条通过校验。
// 一旦 override 块没渲染出来（或者 defaultTheme 少了 Vars），模型照着提示词写就会吃到
// "这些变量没有任何样式消费它们"的拒收，而提示词又禁止它自己定义 —— 死胡同。
func TestFreshDeckAcceptsPresetSemanticVars(t *testing.T) {
	assets := filepath.Join("..", "..", "..", "web", "assets")
	if _, err := os.Stat(assets); err != nil {
		t.Fatalf("读不到真实资源目录（本包到 web/assets 是三层）: %v", err)
	}
	html, err := renderSkeleton("测试", `<section><h1>hi</h1></section>`)
	if err != nil {
		t.Fatal(err)
	}

	// 前提：这三个变量必须真的被渲染进 deck 的样式块（只写在主题 JSON 里不算——
	// 回声校验不认 JSON，浏览器也不认）
	for _, v := range []string{"--accent-2", "--positive", "--warn"} {
		if !strings.Contains(html, v+":") {
			t.Errorf("新建 deck 的 override 块里没有定义 %s，模型引用它会被拒收", v)
		}
	}

	s := &Service{assetsDir: assets}
	for _, v := range []string{"--accent-2", "--positive", "--warn", "--accent", "--hairline"} {
		if err := s.checkInlineStyleVars(`<p style="color:var(--`+strings.TrimPrefix(v, "--")+`)">x</p>`, html); err != nil {
			t.Errorf("新建 deck 里 %s 应当可用，却被拒: %v", v, err)
		}
	}

	// 反向断言：确认这道闸门本身还活着。少了它，上面几条就算"校验被整个关掉"
	// 也会全绿——那种绿说明不了任何事。
	if err := s.checkInlineStyleVars(`<p style="color:var(--nope)">x</p>`, html); err == nil {
		t.Error("不存在的变量必须仍被拒收，否则上面的通过不构成证据")
	}
}
