package deck

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func sectionsFrom(t *testing.T, html string) []*goquery.Selection {
	t.Helper()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatal(err)
	}
	var out []*goquery.Selection
	doc.Find("section").Each(func(_ int, s *goquery.Selection) { out = append(out, s) })
	if len(out) == 0 {
		t.Fatal("测试输入里没有 section")
	}
	return out
}

func lintHTML(t *testing.T, html string) styleReport {
	t.Helper()
	l := newStyleLinter()
	for _, s := range sectionsFrom(t, html) {
		l.add(s)
	}
	return l.report()
}

// 这是这个功能存在的核心理由：页面家具在每一页只出现一次，
// 逐页检测永远看不到重复——只有跨页一起数才发现"同一串抄了 N 遍"。
func TestStyleLintCatchesCrossPageDuplication(t *testing.T) {
	const pageFurniture = `position:absolute; top:0; right:0; font-size:1.2em; font-weight:900; letter-spacing:.03em; color:transparent; -webkit-text-stroke:1.6px var(--accent); opacity:.42`
	html := `<section><span style="` + pageFurniture + `">01</span><h1>一</h1></section>
	         <section><span style="` + pageFurniture + `">02</span><h1>二</h1></section>
	         <section><span style="` + pageFurniture + `">03</span><h1>三</h1></section>`

	rep := lintHTML(t, html)
	if rep.DupStyles != 3 {
		t.Fatalf("跨页重复应报 3 条，实际 %d（%+v）", rep.DupStyles, rep)
	}
	if len(rep.DupSamples) == 0 || !strings.Contains(rep.DupSamples[0], "×3") {
		t.Fatalf("样例里应带上重复次数: %v", rep.DupSamples)
	}

	// 同一份输入逐页检测时一条都抓不到——正是要避免的盲区
	for _, s := range sectionsFrom(t, html) {
		l := newStyleLinter()
		l.add(s)
		if got := l.report().DupStyles; got != 0 {
			t.Fatalf("单页视角不该发现重复，实际 %d", got)
		}
	}
}

// 声明顺序与空白差异不该让"同一串样式"漏掉：模型手写的重复往往就差一个顺序。
func TestStyleLintCanonicalizesBeforeComparing(t *testing.T) {
	html := `<section>
	  <p style="display:flex;gap:1em;align-items:center">a</p>
	  <p style="align-items: center; display: flex; gap: 1em">b</p>
	  <p style="display:flex; align-items:center; gap:1em">c</p>
	</section>`
	if got := lintHTML(t, html).DupStyles; got != 3 {
		t.Fatalf("顺序/空白不同但声明相同的三处应被认作重复，实际 %d", got)
	}
}

// 写在 <section> 上 = 覆盖整套主题 + 会被适配兜底直接覆盖掉 font-size，必须点名。
func TestStyleLintFlagsSectionLevelThemeOverride(t *testing.T) {
	html := `<section style="background:#1b1b32; color:#f4f4f8; font-family:'Noto Sans JP',serif; font-size:30px">
	  <h1>标题</h1>
	</section>`
	rep := lintHTML(t, html)
	want := []string{"background", "color", "font-family", "font-size"}
	for _, w := range want {
		if !containsStr(rep.SectionProps, w) {
			t.Errorf("应报出 section 上的 %s，实际 %v", w, rep.SectionProps)
		}
	}
	if !strings.Contains(rep.Warning(), "模板版式接管") {
		t.Errorf("section 级覆盖的提示应指向正确做法: %s", rep.Warning())
	}
}

func TestStyleLintHardcodedColorsAndPxFonts(t *testing.T) {
	lean := `<section><p style="color:#172554">一</p><p style="color:var(--accent)">二</p></section>`
	if rep := lintHTML(t, lean); rep.Warning() != "" {
		t.Fatalf("1 处写死颜色属于刻意偏离，不该提示: %s", rep.Warning())
	}

	heavy := `<section>
	  <p style="color:#172554">一</p>
	  <p style="background:rgba(23,37,84,.07)">二</p>
	  <p style="border:1px solid #e6ebf7; font-size:14.5px">三</p>
	  <p style="background:linear-gradient(180deg,#f8faff,#eef2ff); font-size:18px">四</p>
	  <p style="font-size:19px">五</p>
	</section>`
	rep := lintHTML(t, heavy)
	if rep.HardColors < lintHardColorLimit {
		t.Fatalf("写死颜色应被数出来，实际 %d", rep.HardColors)
	}
	if rep.PxFontSizes < lintPxFontLimit {
		t.Fatalf("px 字号应被数出来，实际 %d", rep.PxFontSizes)
	}
	w := rep.Warning()
	for _, want := range []string{"var(--accent)", "em"} {
		if !strings.Contains(w, want) {
			t.Errorf("提示里应出现 %q: %s", want, w)
		}
	}
}

// transparent / currentColor 是合法且常用的效果写法，不能被当成"写死颜色"。
func TestStyleLintIgnoresLegitColorKeywords(t *testing.T) {
	html := `<section>
	  <p style="color:transparent; -webkit-text-stroke:1.6px var(--accent)">一</p>
	  <p style="background:transparent">二</p>
	  <p style="color:currentColor">三</p>
	</section>`
	if rep := lintHTML(t, html); rep.HardColors != 0 {
		t.Fatalf("transparent/currentColor 不该被计为写死颜色，实际 %d", rep.HardColors)
	}
}

// "整页只剩一块色块"要被检出——那是用户拿真实截图反馈过的"贴在页面底色上的大卡片"。
// section 自身挂了 .bg-ink/.bg-accent 的（有意的色面页）豁免；唯一子节点不是色块的不计。
func TestStyleLintFlagsSoleColorBlockPage(t *testing.T) {
	html := `<section><div class="accent-block"><p class="stat">01</p><h2>认识 Qt</h2></div></section>` +
		`<section class="bg-accent"><div class="accent-block"><p>有意保留的块</p></div></section>` +
		`<section><div class="half-col"><p>不是色块</p></div></section>`
	rep := lintHTML(t, html)
	if rep.SolePageBlocks != 1 {
		t.Fatalf("应检出 1 页整页色块（bg-accent 页豁免、非色块不计），实际 %d", rep.SolePageBlocks)
	}
	if w := rep.Warning(); !strings.Contains(w, "贴上去的卡片") {
		t.Errorf("警告里应说明症状与正确做法，实际：%s", w)
	}
}

// 用组件库搭出来的页面必须是完全静默的：否则每次写入都带一段噪音，
// 模型很快就会学会忽略它。
func TestStyleLintSilentOnComponentBasedPage(t *testing.T) {
	html := `<section>
	  <span class="page-no outline">03</span>
	  <h2><span class="kicker">第三講</span>两个方案对比</h2>
	  <div class="grid-2">
	    <div class="card"><h3>方案 A</h3><p>实现简单</p></div>
	    <div class="card"><h3>方案 B</h3><p>崩溃安全</p></div>
	  </div>
	  <div class="rows">
	    <div class="row line"><div class="key"><p class="term">～おかげで</p><p class="sub">多亏</p></div>
	    <p class="quote">「魔王を倒せた。」</p></div>
	  </div>
	  <p class="footnote">结论：走方案 B</p>
	</section>`
	if rep := lintHTML(t, html); rep.Warning() != "" {
		t.Fatalf("纯组件页面不该有任何提示: %s", rep.Warning())
	}
}

func TestJoinWarningsKeepsEmptyContract(t *testing.T) {
	if got := joinWarnings("", "", "  "); got != "" {
		t.Fatalf("全空时必须返回空串（omitempty 依赖这个约定），实际 %q", got)
	}
	if got := joinWarnings("A", "", "B"); got != "A B" {
		t.Fatalf("非空项应以空格拼接，实际 %q", got)
	}
}

func TestClipStyleIsRuneSafeAndCollapsesWhitespace(t *testing.T) {
	got := clipStyle("  color:red;\n\t margin:0  ")
	if got != "color:red; margin:0" {
		t.Fatalf("应压缩空白，实际 %q", got)
	}
	long := strings.Repeat("样", lintSampleRunes+10)
	got = clipStyle(long)
	if !strings.HasSuffix(got, "…") {
		t.Fatalf("超长应截断，实际 %q", got)
	}
	if len([]rune(got)) != lintSampleRunes+1 {
		t.Fatalf("截断应按 rune 计（不能劈开多字节字符），实际 %d", len([]rune(got)))
	}
}

// 真实回归：仓库里那份手写型 deck 必须被报出来。
// 这个测试守的是"阈值别被调松到形同虚设"——它也顺带证明整套检测在真实产物上有效。
func TestRealHandWrittenDeckIsFlagged(t *testing.T) {
	p := filepath.Join("..", "..", "..", "data", "decks", "deck-0014", "deck.html")
	raw, err := os.ReadFile(p)
	if err != nil {
		t.Skipf("读不到 %s，跳过: %v", p, err)
	}
	rep := lintHTML(t, string(raw))
	if rep.DupStyles == 0 {
		t.Fatal("deck-0014 的内联样式重复量很大，一条都没报说明阈值或解析出了问题")
	}
	if w := rep.Warning(); !strings.Contains(w, "组件库 class") {
		t.Fatalf("提示里应引导改用组件库 class: %s", w)
	}
	// 这份 deck 里同一串页码样式抄了 8 遍，每遍都带一个手写的 font-size——
	// 内联字号检测在真实产物上必须也报得出来
	if rep.InlineFontSize == 0 {
		t.Error("deck-0014 里有多处手写 font-size，一处都没报说明检测没接上")
	}
	t.Logf("deck-0014 体检结果: 重复 %d 条 / 写死颜色 %d 处 / px 字号 %d 处 / 内联字号 %d 处 / section 级覆盖 %v",
		rep.DupStyles, rep.HardColors, rep.PxFontSizes, rep.InlineFontSize, rep.SectionProps)
	t.Logf("样例: %v", rep.DupSamples)
}

func containsStr(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

// 内联字号：组件库的字号是固定阶梯（0.7 / 0.92 / 1.1 / 1.32 / 1.6 / 2.4 倍基准），
// 手写的字号会插进两档之间，让"哪一层更重要"读不出来。
// 与写死颜色不同，这里**1 处就提示**：颜色的偶尔偏离可以是刻意的设计，
// 字号偏离则是把刚立起来的梯子拆掉，没有"偶尔"这回事。
func TestStyleLintFlagsInlineFontSize(t *testing.T) {
	html := `<section>
	  <h2>标题</h2>
	  <p style="font-size:1.3em">手动放大的一句</p>
	  <p>正常的一句</p>
	</section>`
	rep := lintHTML(t, html)
	if rep.InlineFontSize != 1 {
		t.Fatalf("应报 1 处内联字号，实际 %d（%+v）", rep.InlineFontSize, rep)
	}
	// 提示语必须给出正确做法，否则模型只会把 1.3em 改成 1.2em 再试一次
	w := rep.Warning()
	for _, want := range []string{"font-size", "阶梯", ".bleed"} {
		if !strings.Contains(w, want) {
			t.Errorf("提示语里应出现 %q，实际: %s", want, w)
		}
	}

	// font 简写也算；px 字号只算一次，报的是更要紧的那条（不随适配兜底缩放）
	rep2 := lintHTML(t, `<section>
	  <p style="font: 20px/1.5 sans-serif">a</p>
	  <p style="font-size:1.1em">b</p>
	</section>`)
	if rep2.PxFontSizes != 1 || rep2.InlineFontSize != 1 {
		t.Fatalf("px 与内联字号应各算一次且互斥，实际 px=%d inline=%d（%+v）",
			rep2.PxFontSizes, rep2.InlineFontSize, rep2)
	}

	// 没写内联字号时一个字都不该提（warning 占的是模型的注意力，不能白占）
	clean := lintHTML(t, `<section><h2>标题</h2><div class="card"><p>正文</p></div></section>`)
	if strings.Contains(clean.Warning(), "font-size") {
		t.Fatalf("没有内联字号不该提这件事: %s", clean.Warning())
	}
}

// emoji：提示词禁用（用户明确要求除外），机制要兜底指出来——不然模型违规时
// 只有截图审查可能偶然看到。✓ → 这类正当排版符号刻意不抓（抓了会让模型
// 连对照表都不敢写）。范围依据见 emojiRe 的注释。
func TestStyleLintFlagsEmojiButNotTypographicMarks(t *testing.T) {
	rep := lintHTML(t, `<section><h1>营收增长 🔥</h1><p>用户突破 🚀 100 万</p></section>`)
	if len(rep.EmojiPages) != 1 {
		t.Fatalf("应报 1 页 emoji，实际 %+v", rep.EmojiPages)
	}
	if !strings.Contains(rep.EmojiPages[0], "×2") {
		t.Fatalf("应数出 2 个 emoji: %v", rep.EmojiPages)
	}
	w := rep.Warning()
	for _, want := range []string{"emoji", ".ico", "read_icons"} {
		if !strings.Contains(w, want) {
			t.Errorf("提示语里应出现 %q，实际: %s", want, w)
		}
	}

	// ✓ 与 → 是正当排版符号（对照表、流向），不该被抓
	clean := lintHTML(t, `<section><p>A 方案 ✓ 已具备</p><p>下一步 → 部署</p></section>`)
	if len(clean.EmojiPages) != 0 {
		t.Fatalf("✓ 与 → 不该被抓: %+v", clean.EmojiPages)
	}

	// 多页场景提示要带页号；单页场景（update_slide）自动说"本页"
	multi := lintHTML(t, `<section><h1>干净的页</h1></section>
	                      <section><h1>这页有 🔥</h1></section>`)
	if len(multi.EmojiPages) != 1 || !strings.Contains(multi.EmojiPages[0], "第 2 页") {
		t.Fatalf("多页场景应带页号: %+v", multi.EmojiPages)
	}
	single := newStyleLinter()
	secs := sectionsFrom(t, `<section><h1>有 🔥 的页</h1></section>`)
	single.add(secs[0])
	if !strings.Contains(single.report().Warning(), "本页") {
		t.Fatalf("单页场景应说本页: %s", single.report().Warning())
	}
}

// 破折号"一页最多一个"数的是归并后的处数：中文破折号 —— 是两个 U+2014，
// 按字符数会把合法的一个 —— 误判成两处。
func TestStyleLintFlagsExcessEmDash(t *testing.T) {
	rep := lintHTML(t, `<section><h1>背景——现状</h1><p>问题——但是——机会</p></section>`)
	if len(rep.DashPages) != 1 || !strings.Contains(rep.DashPages[0], "3 处") {
		t.Fatalf("应报 1 页 3 处破折号，实际 %+v", rep.DashPages)
	}
	w := rep.Warning()
	if !strings.Contains(w, "一页最多一个") {
		t.Errorf("提示语应带上规则原文: %s", w)
	}

	one := lintHTML(t, `<section><h1>背景——现状</h1><p>其他内容</p></section>`)
	if len(one.DashPages) != 0 {
		t.Fatalf("单个 —— 是合法用法，不该报: %+v", one.DashPages)
	}
}
