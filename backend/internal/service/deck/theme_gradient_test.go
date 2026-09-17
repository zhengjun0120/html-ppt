package deck

// 背景渐变（bg_gradient）的回归测试。要守四件事：
//
//  1. 渐变值校验挡得住坏输入——:root 块旁边开的新自由度，结构不校验就是注入面。
//  2. 纹理与渐变合成同一条 background-image——分开写会被 CSS 级联整体替换，
//     "有纹理的渐变主题"只剩其中一个，而且是静默的。
//  3. 纯色主题的渲染产物逐字不变——textureCSS 的重构不许给老 deck 制造历史 diff。
//  4. 预设带出来的渐变在每个颜色断点上都过对比度下限——渐变的中间亮度只有
//     "逐断点都过线"才保证投影可读，只测 Background 字段会漏掉亮端。
import (
	"encoding/json"
	"strings"
	"testing"
)

func TestValidateGradient(t *testing.T) {
	valid := []string{
		"",
		"   ",
		"linear-gradient(168deg, #091830 0%, #10403f 100%)",
		"radial-gradient(circle at 30% 20%, #1d0d08, #5c1f10)",
		"conic-gradient(from 90deg, #0a1626, #10403f)",
	}
	for _, v := range valid {
		if err := validateGradient(v); err != nil {
			t.Errorf("%q 应合法，实际报错: %v", v, err)
		}
	}

	invalid := []string{
		"linear-gradient(#000, #fff); body{background:red}",
		"gradient(#000, #fff)",                       // 不是三种原生渐变函数
		"linear-gradient(#000, #fff",                 // 括号不闭合
		"linear-gradient(#000, #fff))",               // 括号多一个
		"url(https://evil.example/x.png)",            // 外部加载通道
		"linear-gradient(#000, #fff) @import url(x)", // 混入 @import
		"expression(alert(1))",                       // 老式执行通道
		"linear-gradient(#000, #fff) <script>",       // 结构破坏
	}
	for _, v := range invalid {
		if err := validateGradient(v); err == nil {
			t.Errorf("%q 应被拒收，实际通过", v)
		}
	}

	long := "linear-gradient(" + strings.Repeat("#aaa 0%,", 40) + "#bbb 100%)"
	if err := validateGradient(long); err == nil {
		t.Error("超长渐变应被拒收（主题块不是插画载体）")
	}
}

func TestTextureCSSComposesGradientAndTexture(t *testing.T) {
	grad := "linear-gradient(168deg, #091830 0%, #10403f 100%)"

	t.Run("只有渐变", func(t *testing.T) {
		th := defaultTheme()
		th.Texture = TextureNone
		th.BgGradient = grad
		got := textureCSS(th)
		if !strings.Contains(got, grad) {
			t.Errorf("渐变应原样出现在规则里: %s", got)
		}
		if strings.Contains(got, "background-size") {
			t.Errorf("无纹理时不应有 background-size: %s", got)
		}
	})

	t.Run("渐变与网格纹理同一条规则", func(t *testing.T) {
		th := defaultTheme() // paper 自带 grid 纹理
		th.BgGradient = grad
		got := textureCSS(th)
		if !strings.Contains(got, grad) {
			t.Errorf("渐变丢了: %s", got)
		}
		if !strings.Contains(got, "linear-gradient(rgba(43, 39, 35, 0.07) 1px,transparent 1px)") {
			t.Errorf("纹理层丢了: %s", got)
		}
		// 纹理在上、渐变垫底：纹理层必须是 background-image 的第一层
		if strings.Index(got, "rgba(43, 39, 35, 0.07) 1px") > strings.Index(got, grad) {
			t.Errorf("纹理应在渐变之上（多层 background-image 的层序）: %s", got)
		}
	})

	t.Run("清掉渐变回到纯色", func(t *testing.T) {
		th := defaultTheme()
		th.BgGradient = grad
		th.Texture = TextureNone
		clear := ""
		if err := (ThemePatch{BGGradient: &clear}).applyTo(&th); err != nil {
			t.Fatal(err)
		}
		if th.BgGradient != "" {
			t.Fatalf("空串应清掉渐变，实际 %q", th.BgGradient)
		}
		if err := th.validate(); err != nil {
			t.Fatalf("清掉后应通过全量校验: %v", err)
		}
		if got := textureCSS(th); got != "" {
			t.Errorf("无纹理无渐变不应产出规则: %s", got)
		}
	})

	t.Run("校验不过的渐变进不了主题", func(t *testing.T) {
		th := defaultTheme()
		bad := "linear-gradient(#000,#fff); body{x:y}"
		if err := (ThemePatch{BGGradient: &bad}).applyTo(&th); err != nil {
			t.Fatal(err)
		}
		if err := th.validate(); err == nil {
			t.Fatal("带注入片段的渐变应被 validate 拒收")
		}
	})

	t.Run("改渐变算偏离预设", func(t *testing.T) {
		th := defaultTheme()
		g := "radial-gradient(#111, #222)"
		if err := (ThemePatch{BGGradient: &g}).applyTo(&th); err != nil {
			t.Fatal(err)
		}
		if th.Preset != "" {
			t.Errorf("单独改 bg_gradient 后预设名应清空（它不再是原样纸感），实际 %q", th.Preset)
		}
	})
}

// 纯色主题（无渐变）的渲染产物必须逐字等于重构前——老 deck 的历史 diff 里
// 不该出现"升级代码导致每份 deck 的 CSS 都变了"这种假改动。
func TestTextureOnlyCSSUnchanged(t *testing.T) {
	paper := defaultTheme()
	want := ".reveal-viewport{background-image:linear-gradient(rgba(43, 39, 35, 0.07) 1px,transparent 1px)," +
		"linear-gradient(90deg,rgba(43, 39, 35, 0.07) 1px,transparent 1px);background-size:30px 30px}"
	if got := textureCSS(paper); got != want {
		t.Errorf("纸感网格纹理的渲染产物变了：\n got=%s\nwant=%s", got, want)
	}
}

// 带渐变的预设渲染进骨架后，bg_gradient 必须真的出现在主题 JSON 与 override CSS 里——
// write_deck(preset=aurora) 的整条链路靠的就是这一步。
func TestGradientPresetRendersIntoSkeleton(t *testing.T) {
	theme := defaultTheme()
	presets[PresetAurora].applyTo(&theme)
	html, err := renderSkeleton("渐变测试", `<section><h1>hi</h1></section>`, theme)
	if err != nil {
		t.Fatalf("渲染骨架失败: %v", err)
	}
	start := strings.Index(html, `id="deck-theme">`)
	start += len(`id="deck-theme">`)
	end := strings.Index(html[start:], `</script>`)
	var got Theme
	if err := json.Unmarshal([]byte(html[start:start+end]), &got); err != nil {
		t.Fatalf("主题 JSON 不能解析: %v", err)
	}
	if got.BgGradient != presets[PresetAurora].Theme.BgGradient {
		t.Errorf("骨架里的 bg_gradient 与预设不一致: %q", got.BgGradient)
	}
	if !strings.Contains(html, "linear-gradient(168deg") {
		t.Error("override CSS 里没有渐变（浏览器看不到背景渐变）")
	}
}

// 渐变预设的每个断点都要过对比度下限。只测 Background 字段守不住渐变：
// 亮端往往才是对比度最吃紧的那一档。
func TestPresetGradientStopsMeetContrastFloor(t *testing.T) {
	for _, p := range Presets() {
		grad := p.Theme.BgGradient
		if grad == "" {
			continue
		}
		stops := extractHexColors(grad)
		if len(stops) < 2 {
			t.Fatalf("预设 %s 的渐变里抠不出颜色断点: %q", p.Name, grad)
		}
		for _, stop := range stops {
			for _, c := range []struct{ name, fg string }{
				{"正文", p.Theme.TextColor},
				{"标题", p.Theme.HeadingColor},
			} {
				if got := contrastRatio(c.fg, stop); got < 7 {
					t.Errorf("预设 %s 的渐变断点 %s 上，%s 只有 %.2f:1（下限 7:1）",
						p.Name, stop, c.name, got)
				}
			}
		}
	}
}

// extractHexColors 从 CSS 渐变字符串里按出现顺序抠 #rrggbb。
func extractHexColors(s string) []string {
	var out []string
	for i := 0; i < len(s)-6; i++ {
		if s[i] != '#' {
			continue
		}
		cand := s[i+1 : i+7]
		ok := true
		for _, r := range cand {
			if !strings.ContainsRune("0123456789abcdefABCDEF", r) {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, "#"+cand)
			i += 6
		}
	}
	return out
}
