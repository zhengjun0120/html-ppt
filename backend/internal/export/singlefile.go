// singlefile.go —— 单文件 HTML 导出的自包含化。
//
// 预览 HTML（deck.PreviewHTML）只内联了模板 style.css；翻页与字体仍挂在
// /assets 绝对引用上。文件一旦离开服务器打开（下载后双击、挪去静态托管），
// 这些引用全部 404：slide 失去 absolute 定位与 JS class 切换，退化为整页
// 堆叠滚动——导出的"PPT"不能翻页。这里在落盘前把四类引用替换成内联内容：
//   - runtime.js / base.css / animations.css 原样内联（翻页三件套，~75KB）；
//   - fonts.css 内联时把 7 个常用 @font-face 的 url 换成 data: base64
//     （Inter 拉丁 ×2 + Noto Sans SC ×3 + Maple 拉丁 ×2，约 +5.5MB）；
//     两个 CJK 全量等宽（6.2MB ×2）整块剔除，代码块里的中文回退系统等宽。
//
// 内联后产物不含任何 /assets 引用，任意环境打开都是完整的单页翻页版。
package export

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// 与 export.go 的 styleLinkRe 同款姿势：锚定 href/src 精确匹配，不依赖属性顺序。
var (
	baseCSSLinkRe      = regexp.MustCompile(`<link[^>]*href="/assets/deck-v2/base\.css"[^>]*>`)
	animationCSSLinkRe = regexp.MustCompile(`<link[^>]*href="/assets/deck-v2/animations\.css"[^>]*>`)
	fontsCSSLinkRe     = regexp.MustCompile(`<link[^>]*href="/assets/deck-v2/fonts\.css"[^>]*>`)
	runtimeScriptRe    = regexp.MustCompile(`<script[^>]*src="/assets/deck-v2/runtime\.js"[^>]*>\s*</script>`)
	fontURLRe          = regexp.MustCompile(`url\('([^']+)'\)`)
)

// cjkFullMonoSkip 等宽 CJK 全量两支（单支 6MB+）：base64 后会把导出文件撑到
// 20MB+，不值得。剔除对应 @font-face 后，代码块里的中文注释回退系统等宽字体，
// ASCII 代码仍由拉丁子集精确渲染。
var cjkFullMonoSkip = map[string]bool{
	"JetBrainsMapleMono-Regular.woff2": true,
	"JetBrainsMapleMono-Bold.woff2":    true,
}

// makeSelfContained 把预览 HTML 的四类 /assets 引用全部内联。
// 任何引用标签没匹配上、资产读不到，都返回错误——宁可拒绝导出，不出残次品。
func (s *Service) makeSelfContained(html string) (string, error) {
	asset := func(name string) ([]byte, error) {
		return os.ReadFile(filepath.Join(s.assetsDir, "deck-v2", name))
	}
	// css/js 内容进替换串时用 Literal 变体：内容里的 $ 不是正则占位符
	inlineCSS := func(html string, re *regexp.Regexp, name, body string) (string, error) {
		if !re.MatchString(html) {
			return "", fmt.Errorf("HTML 中找不到 %s 的引用标签（生成器输出已变更？）", name)
		}
		return re.ReplaceAllLiteralString(html, "<style>\n"+body+"\n</style>"), nil
	}

	for _, a := range []struct {
		re   *regexp.Regexp
		name string
	}{
		{baseCSSLinkRe, "base.css"},
		{animationCSSLinkRe, "animations.css"},
	} {
		raw, err := asset(a.name)
		if err != nil {
			return "", fmt.Errorf("读取 %s 失败: %w", a.name, err)
		}
		html, err = inlineCSS(html, a.re, a.name, string(raw))
		if err != nil {
			return "", err
		}
	}

	fonts, err := s.inlineFonts()
	if err != nil {
		return "", err
	}
	html, err = inlineCSS(html, fontsCSSLinkRe, "fonts.css", fonts)
	if err != nil {
		return "", err
	}

	js, err := asset("runtime.js")
	if err != nil {
		return "", fmt.Errorf("读取 runtime.js 失败: %w", err)
	}
	// HTML 解析器不懂 JS 语法：script 块内容里出现 </script 就提前终结块
	//（哪怕是 JS 注释/字符串）。统一转义成 <\/script——JS 里 \/ 恒等 /，
	// 字符串、模板、正则、注释四种语境下语义都不变。
	if strings.Contains(string(js), "</script") {
		js = []byte(strings.ReplaceAll(string(js), "</script", `<\/script`))
	}
	// <!-- 加 <script 的组合会把解析器带进双逃逸态：那时块尾的 </script>
	// 只退回逃逸态而收不掉块，整个文档尾部被吞。现网 runtime 两者皆无，
	// 这里拒绝将来引入该组合的版本（无法用替换安全化解）。
	if strings.Contains(string(js), "<!--") && strings.Contains(string(js), "<script") {
		return "", fmt.Errorf("runtime.js 同时含 <!-- 与 <script，内联会触发 HTML 双逃逸，拒绝导出")
	}
	if !runtimeScriptRe.MatchString(html) {
		return "", fmt.Errorf("HTML 中找不到 runtime.js 的引用标签（生成器输出已变更？）")
	}
	html = runtimeScriptRe.ReplaceAllLiteralString(html, "<script>\n"+string(js)+"\n</script>")

	if strings.Contains(html, "/assets/deck-v2/") {
		return "", fmt.Errorf("单文件化后仍有 /assets 残留，拒绝导出不完整产物")
	}
	return html, nil
}

// inlineFonts 把 fonts.css 变成自包含的 <style> 内容：每个 @font-face 的
// url 换成 data: base64，CJK 全量等宽两支整块剔除。
func (s *Service) inlineFonts() (string, error) {
	raw, err := os.ReadFile(filepath.Join(s.assetsDir, "deck-v2", "fonts.css"))
	if err != nil {
		return "", fmt.Errorf("读取 fonts.css 失败: %w", err)
	}
	blocks := strings.Split(string(raw), "@font-face")
	var b strings.Builder
	for i, block := range blocks {
		if i == 0 {
			b.WriteString(block) // 首段是文件头注释
			continue
		}
		m := fontURLRe.FindStringSubmatch(block)
		if m == nil {
			// 无 url 的块不该存在，原样保留让问题可见（产物守卫会兜底）
			b.WriteString("@font-face" + block)
			continue
		}
		name := filepath.Base(m[1])
		if cjkFullMonoSkip[name] {
			continue
		}
		font, err := os.ReadFile(filepath.Join(s.assetsDir, "deck-v2", "fonts", name))
		if err != nil {
			return "", fmt.Errorf("读取字体 %s 失败: %w", name, err)
		}
		dataURL := "url(data:font/woff2;base64," + base64.StdEncoding.EncodeToString(font) + ")"
		b.WriteString("@font-face" + strings.Replace(block, m[0], dataURL, 1))
	}
	return b.String(), nil
}
