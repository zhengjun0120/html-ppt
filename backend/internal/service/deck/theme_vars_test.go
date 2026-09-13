package deck

// 主题自由度的边界测试。这里守两件事：
//
//  1. **配色可以完全自由**：Vars 让 AI 设计一整套自有调色板（这是用户明确要的能力），
//     且这些变量必须渲染成"一处定义"，不能散落成字面量。
//  2. **配色只有一个权威源**：契约变量（--accent/--border/--r-* 等）只许由结构化字段派生，
//     谁都不许在别处重定义。自定义槽排在主题块之后、同级选择器靠后者取胜，
//     所以"在槽里重定义"会让 update_theme 静默失效——这是本次改动要堵掉的核心坑。
//
// 另外守一条工程细节：Vars 是 map，渲染必须排序。不排的话每次写入产出的 :root
// 文本都在变，历史 diff 会把整块报成"改过"，一堆假改动比没有 diff 更糟。

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// TestCanvasPresetsMatchInitJS 契约测试：画布预设表在 Go 和 init.js 各存一份
// （前端要在 Reveal.initialize 之前拿到尺寸，没法从后端要），两份必须逐字一致。
// 漂移的后果很难查：后端换了尺寸、前端还用旧的，适配兜底按新高度算、reveal 按旧高度摆，
// 页面会莫名其妙地被判成溢出或留出空白。
func TestCanvasPresetsMatchInitJS(t *testing.T) {
	p := filepath.Join("..", "..", "web", "assets", "init.js")
	js, err := os.ReadFile(p)
	if err != nil {
		t.Skipf("读不到 %s，跳过契约检查: %v", p, err)
	}
	src := string(js)

	for name, size := range canvasPresets {
		// 形如：wide: { width: 1244, height: 700 },
		re := regexp.MustCompile(regexp.QuoteMeta(name) + `\s*:\s*\{\s*width:\s*(\d+)\s*,\s*height:\s*(\d+)\s*\}`)
		m := re.FindStringSubmatch(src)
		if m == nil {
			t.Errorf("init.js 的 CANVAS 表里找不到预设 %q", name)
			continue
		}
		w, _ := strconv.Atoi(m[1])
		h, _ := strconv.Atoi(m[2])
		if w != size[0] || h != size[1] {
			t.Errorf("预设 %s 尺寸不一致：Go=%dx%d init.js=%dx%d", name, size[0], size[1], w, h)
		}
	}

	// 兜底必须与 standard 一致：主题块缺失/损坏时两边各自退回默认，退到不同值就会错位
	def := canvasPresets[CanvasStandard]
	fallbackRe := regexp.MustCompile(`CANVAS\[deckTheme\.canvas\]\s*\|\|\s*CANVAS\.(\w+)`)
	if m := fallbackRe.FindStringSubmatch(src); m == nil {
		t.Error("init.js 里找不到画布兜底表达式（CANVAS[deckTheme.canvas] || CANVAS.xxx）")
	} else if m[1] != CanvasStandard {
		t.Errorf("init.js 的兜底预设是 %q，应与 CanvasStandard=%q 一致", m[1], CanvasStandard)
	}
	if !strings.Contains(src, "width: "+strconv.Itoa(def[0])) && !strings.Contains(src, "width: "+strconv.Itoa(def[0])+",") {
		t.Errorf("init.js 里应出现 standard 的宽度 %d", def[0])
	}
}

func TestCanvasSizeFallsBackToDefault(t *testing.T) {
	w, h := CanvasSize("nonsense")
	if w != 960 || h != 700 {
		t.Fatalf("未知预设应退回 960×700，实际 %dx%d", w, h)
	}
	// 高度必须恒为 700：容量由高度决定，换比例不该改变一页能放多少内容
	for name, size := range canvasPresets {
		if size[1] != 700 {
			t.Errorf("预设 %s 的高度是 %d，应统一为 700", name, size[1])
		}
	}
}

func TestValidateThemeVarsAcceptsBespokePalette(t *testing.T) {
	palette := map[string]string{
		"--surface":    "#1e1836",
		"--surface-2":  "rgba(30, 24, 54, 0.6)",
		"--positive":   "#4ade80",
		"--danger":     "#f87171",
		"--hero-glow":  "0 0 32px rgba(255, 107, 53, .45)",
		"--brand-font": `"Noto Serif JP", serif`,
	}
	if err := validateThemeVars(palette); err != nil {
		t.Fatalf("一整套自有配色应被接受，实际: %v", err)
	}
}

func TestValidateThemeVarsRejects(t *testing.T) {
	type varCase struct {
		name     string
		vars     map[string]string
		wantWord string
	}
	cases := []varCase{
		{"占用 accent", map[string]string{"--accent": "#fff"}, "--accent"},
		{"占用 border", map[string]string{"--border": "#fff"}, "--border"},
		{"占用 reveal 变量", map[string]string{"--r-main-color": "#fff"}, "--r-main-color"},
		{"占用间距变量", map[string]string{"--space-md": "1em"}, "--space-md"},
		{"名字没有双连字符", map[string]string{"surface": "#fff"}, "不合法"},
		{"名字以数字开头", map[string]string{"--2x": "#fff"}, "不合法"},
		{"空值", map[string]string{"--surface": "  "}, "不能为空"},
		{"值里带分号", map[string]string{"--surface": "#fff;--accent:#000"}, `";"`},
		{"值里带花括号", map[string]string{"--surface": "#fff}body{display:none"}, `"{"`},
		{"值里带标签", map[string]string{"--surface": "</style><script>x"}, "<"},
		{"值里带 url", map[string]string{"--surface": "url(https://evil.com/a.png)"}, "url("},
		{"值里带 import", map[string]string{"--surface": "@import '/x.css'"}, "@import"},
		{"值太长", map[string]string{"--surface": strings.Repeat("a", maxVarValueBytes+1)}, "上限"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := validateThemeVars(c.vars)
			if err == nil {
				t.Fatalf("应被拒绝: %v", c.vars)
			}
			if !strings.Contains(err.Error(), c.wantWord) {
				t.Errorf("错误信息应点明 %q，实际: %s", c.wantWord, err.Error())
			}
		})
	}

	// 占用契约变量时，必须告诉模型"该改成传哪个字段"，否则它会反复重试同一写法
	err := validateThemeVars(map[string]string{"--accent": "#fff"})
	if err == nil || !strings.Contains(err.Error(), "accent") {
		t.Errorf("应指引模型直接传 accent 字段，实际: %v", err)
	}

	// 数量上限
	big := map[string]string{}
	for i := 0; i < maxThemeVars+1; i++ {
		big["--v"+strconv.Itoa(i)] = "#fff"
	}
	if err := validateThemeVars(big); err == nil {
		t.Error("超过变量数量上限应被拒绝")
	}
}

// 渲染必须排序：map 遍历顺序随机，不排的话每次写入产出的 :root 都在变，
// 历史 diff 会把整块报成"改过了"。
func TestRenderThemeCSSIsDeterministic(t *testing.T) {
	theme := defaultTheme()
	theme.Vars = map[string]string{"--zeta": "#111", "--alpha": "#222", "--mid": "#333"}

	first := renderThemeCSS(theme)
	for i := 0; i < 20; i++ {
		if got := renderThemeCSS(theme); got != first {
			t.Fatalf("同一主题渲染结果不稳定（map 没排序）:\n%s\n%s", first, got)
		}
	}
	// 自定义变量要出现在 :root 里，且按名字排序
	iAlpha := strings.Index(first, "--alpha:#222")
	iMid := strings.Index(first, "--mid:#333")
	iZeta := strings.Index(first, "--zeta:#111")
	if iAlpha < 0 || iMid < 0 || iZeta < 0 {
		t.Fatalf("自定义变量没渲染进 :root: %s", first)
	}
	if !(iAlpha < iMid && iMid < iZeta) {
		t.Errorf("自定义变量应按名字排序渲染: alpha@%d mid@%d zeta@%d", iAlpha, iMid, iZeta)
	}
	// 派生变量仍然在位（自定义调色板不该挤掉它们）
	for _, want := range []string{"--accent:#5eead4", "--border:", "--card-bg:", "--radius:10px"} {
		if !strings.Contains(first, want) {
			t.Errorf("派生变量 %q 丢了: %s", want, first)
		}
	}
}

func TestThemePatchVarsSemantics(t *testing.T) {
	t.Run("传了就是整体替换", func(t *testing.T) {
		theme := defaultTheme()
		theme.Vars = map[string]string{"--old": "#111"}
		vars := map[string]string{"--new": "#222"}
		ThemePatch{Vars: &vars}.applyTo(&theme)
		if !reflect.DeepEqual(theme.Vars, map[string]string{"--new": "#222"}) {
			t.Fatalf("vars 应为整体替换（--old 不该留下），实际 %v", theme.Vars)
		}
	})

	t.Run("空 map 表示清空", func(t *testing.T) {
		theme := defaultTheme()
		theme.Vars = map[string]string{"--old": "#111"}
		empty := map[string]string{}
		ThemePatch{Vars: &empty}.applyTo(&theme)
		if theme.Vars != nil {
			t.Fatalf("空 map 应清空自定义调色板（不留空壳），实际 %v", theme.Vars)
		}
	})

	t.Run("nil 不动", func(t *testing.T) {
		theme := defaultTheme()
		theme.Vars = map[string]string{"--keep": "#111"}
		ThemePatch{}.applyTo(&theme)
		if !reflect.DeepEqual(theme.Vars, map[string]string{"--keep": "#111"}) {
			t.Fatalf("没传 vars 时不该改动，实际 %v", theme.Vars)
		}
	})

	t.Run("存的是副本", func(t *testing.T) {
		theme := defaultTheme()
		vars := map[string]string{"--x": "#111"}
		ThemePatch{Vars: &vars}.applyTo(&theme)
		vars["--x"] = "#pwned" // 调用方之后改了入参
		if theme.Vars["--x"] != "#111" {
			t.Fatal("主题里存的应是副本：调用方改入参不该穿透到存储层")
		}
	})
}

// TestThemeVarsRoundTrip 完整走一遍真实写入路径（读→合并→校验→渲染→落盘→读回）。
// 之前这条链路藏在 authorize 后面测不到，而它恰恰是最该被测的：Vars 是新加的
// map 字段，序列化/反序列化/渲染/清空每一段都可能断。
func TestThemeVarsRoundTrip(t *testing.T) {
	s := newHistoryTestService(t)
	if err := writeDeckHTML(s, "deck-9201", miniSkeleton); err != nil {
		t.Fatal(err)
	}

	palette := map[string]string{"--surface": "#1e1836", "--positive": "#4ade80"}
	got, err := s.updateThemeLocked("deck-9201", ThemePatch{
		Accent: strPtr("#ff6b35"),
		Canvas: strPtr(CanvasWide),
		Vars:   &palette,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Accent != "#ff6b35" || got.Canvas != CanvasWide {
		t.Fatalf("patch 未生效: %+v", got)
	}
	if !reflect.DeepEqual(got.Vars, palette) {
		t.Fatalf("vars 未写回: %+v", got.Vars)
	}

	// 落盘内容：override 块里必须同时有派生变量和自定义变量
	raw, _ := s.readRaw("deck-9201")
	if !strings.Contains(raw, "--surface:#1e1836") || !strings.Contains(raw, "--positive:#4ade80") {
		t.Fatalf("自定义调色板没进 #deck-theme-override:\n%s", raw)
	}
	if !strings.Contains(raw, "--accent:#ff6b35") {
		t.Fatalf("强调色没进 override: \n%s", raw)
	}
	// JSON 真身里也要有
	if !strings.Contains(raw, `"canvas":"wide"`) || !strings.Contains(raw, `"vars"`) {
		t.Fatalf("主题 JSON 应带上 canvas 与 vars:\n%s", raw)
	}

	// 读回一致（raw-text 转义坑：主题 JSON 里有引号）
	back, err := s.readThemeRaw("deck-9201")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(back, got) {
		t.Fatalf("读回不一致:\n got=%+v\nback=%+v", got, back)
	}

	// 不带 vars 的 patch 不该动它
	after, err := s.updateThemeLocked("deck-9201", ThemePatch{Background: strPtr("#101010")})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(after.Vars, palette) {
		t.Fatalf("没传 vars 时不该改动: %+v", after.Vars)
	}

	// 空 map = 清空，且不能留下空壳
	empty := map[string]string{}
	cleared, err := s.updateThemeLocked("deck-9201", ThemePatch{Vars: &empty})
	if err != nil {
		t.Fatal(err)
	}
	if cleared.Vars != nil {
		t.Fatalf("空 map 应清空，实际 %+v", cleared.Vars)
	}
	raw, _ = s.readRaw("deck-9201")
	if strings.Contains(raw, "--surface:") {
		t.Fatal("清空后 override 里不该还有自定义变量")
	}
}

// TestThemeUpdateRejectsBadVarsAndKeepsFile 校验失败必须不落盘——
// 半套非法配色写进去比直接报错坏得多。
func TestThemeUpdateRejectsBadVarsAndKeepsFile(t *testing.T) {
	s := newHistoryTestService(t)
	if err := writeDeckHTML(s, "deck-9202", miniSkeleton); err != nil {
		t.Fatal(err)
	}
	before, _ := s.readRaw("deck-9202")

	bad := map[string]string{"--accent": "#fff"} // 占用契约变量
	if _, err := s.updateThemeLocked("deck-9202", ThemePatch{Vars: &bad}); err == nil {
		t.Fatal("占用契约变量的 vars 应被拒绝")
	}
	after, _ := s.readRaw("deck-9202")
	if before != after {
		t.Fatal("校验失败不应改动文件")
	}

	// 画布预设也是白名单
	if _, err := s.updateThemeLocked("deck-9202", ThemePatch{Canvas: strPtr("superwide")}); err == nil {
		t.Fatal("未知画布预设应被拒绝")
	}
	if after2, _ := s.readRaw("deck-9202"); after2 != before {
		t.Fatal("画布预设校验失败不应改动文件")
	}
}

func strPtr(s string) *string { return &s }
