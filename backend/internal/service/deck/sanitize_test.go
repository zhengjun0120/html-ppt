package deck

// 消毒闸门的回归测试。策略是分级的（见 sanitize.go）：
//   结构性危险标签（script/style/iframe/form/object/embed/base/meta/link）→ 拒绝写入
//   属性级违规（on* 事件、危险协议 URL）→ 剥除 + 回报警告
// 这里覆盖两类，以及几种真实的协议绕过写法（空白/大小写），
// 还有"不能误伤正常内容"的边界——消毒最常见的翻车方式是管得太宽。

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

// sanitizeHTML 走完整闸门：解析出根 section → 消毒 → 返回报告与清理后的 HTML。
func sanitizeHTML(t *testing.T, fragment string) (sanitizeReport, string) {
	t.Helper()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(fragment))
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	sec := doc.Find("body").Children().Filter("section").First()
	if len(sec.Nodes) == 0 {
		t.Fatalf("片段里没有 section: %s", fragment)
	}
	rep := sanitizeSlide(sec)
	out, err := goquery.OuterHtml(sec)
	if err != nil {
		t.Fatal(err)
	}
	return rep, out
}

func TestSanitizeRejectsDangerousTags(t *testing.T) {
	cases := []struct {
		name string
		html string
		tag  string
	}{
		{"script", `<section><script>alert(1)</script></section>`, "script"},
		{"style", `<section><style>.x{color:red}</style></section>`, "style"},
		{"iframe", `<section><iframe src="https://evil.com"></iframe></section>`, "iframe"},
		{"object", `<section><object data="https://evil.com"></object></section>`, "object"},
		{"embed", `<section><embed src="https://evil.com"></section>`, "embed"},
		{"form", `<section><form action="https://evil.com"><input name="x"></form></section>`, "form"},
		{"base", `<section><base href="https://evil.com/"></section>`, "base"},
		{"meta", `<section><meta http-equiv="refresh" content="0;url=https://evil.com"></section>`, "meta"},
		{"link", `<section><link rel="stylesheet" href="https://evil.com/x.css"></section>`, "link"},
		{"嵌在内层", `<section><div><p>x</p><iframe src="//evil.com"></iframe></div></section>`, "iframe"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rep, out := sanitizeHTML(t, c.html)
			if err := rep.RejectErr(); err == nil {
				t.Fatalf("<%s> 应被拒绝，实际放过", c.tag)
			}
			// 报告必须指名到标签：错误信息要对模型可行动，不能只说"有违规"
			named := false
			for _, tag := range rep.BannedTags {
				if tag == c.tag {
					named = true
				}
			}
			if !named {
				t.Fatalf("报告未列出 <%s>: %+v", c.tag, rep.BannedTags)
			}
			if strings.Contains(out, "<"+c.tag) {
				t.Fatalf("被拒标签仍残留在树里: %s", out)
			}
		})
	}
}

func TestSanitizeStripsHandlersAndBadURLs(t *testing.T) {
	cases := []struct {
		name string
		html string
		gone string // 清理后不应再出现的子串（小写比较）
	}{
		{"onclick", `<section><div class="card" onclick="alert(1)">x</div></section>`, "onclick"},
		{"onerror 大写", `<section><img src="x.png" ONERROR="alert(1)"></section>`, "onerror"},
		{"svg onload", `<section><svg onload="alert(1)"><circle r="4"/></svg></section>`, "onload"},
		{"javascript: 链接", `<section><a href="javascript:alert(1)">点</a></section>`, "javascript:"},
		{"javascript: 制表符绕过", "<section><a href=\"java\tscript:alert(1)\">点</a></section>", "javascript:"},
		{"javascript: 换行绕过", "<section><a href=\"java\nscript:alert(1)\">点</a></section>", "javascript:"},
		{"javascript: 前导空白", `<section><a href="  javascript:alert(1)">点</a></section>`, "javascript:"},
		{"vbscript:", `<section><a href="vbscript:msgbox(1)">点</a></section>`, "vbscript:"},
		{"data:text/html", `<section><a href="data:text/html,<b>x</b>">点</a></section>`, "data:text/html"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rep, out := sanitizeHTML(t, c.html)
			if err := rep.RejectErr(); err != nil {
				t.Fatalf("属性级违规不应拒绝整页: %v", err)
			}
			if rep.Clean() {
				t.Fatal("应产出违规报告")
			}
			if rep.Warning() == "" {
				t.Fatal("警告语不能为空——静默剥除就是静默失败")
			}
			if strings.Contains(strings.ToLower(out), c.gone) {
				t.Fatalf("危险值未剥净: %s", out)
			}
		})
	}
}

// TestSanitizeKeepsLegitContent 消毒最容易翻车的地方：管得太宽。
// 主题变量、内联 style、data-*、内联图片、锚点、fragment 类名都必须原样留下。
func TestSanitizeKeepsLegitContent(t *testing.T) {
	html := `<section data-id="s1">` +
		`<h2 style="color:var(--accent);font-size:36px">标题</h2>` +
		`<div class="grid-2">` +
		`<div class="card"><img src="data:image/png;base64,iVBORw0KGgo=" alt="图"></div>` +
		`<div class="card"><a href="#page2">锚点</a><a href="https://example.com/doc">外链</a></div>` +
		`</div>` +
		`<ul class="steps"><li class="fragment">一</li></ul>` +
		`<div class="chart" data-chart-height="320"></div>` +
		`</section>`

	rep, out := sanitizeHTML(t, html)
	if !rep.Clean() {
		t.Fatalf("正常内容被误判违规: %+v", rep)
	}
	for _, want := range []string{
		`data-id="s1"`,
		"color:var(--accent)",
		"data:image/png",
		"#page2",
		"https://example.com/doc", // 外链是用户点击才发生的导航，不属于自动加载/执行通道
		`class="fragment"`,
		`data-chart-height="320"`, // data-* 是组件库的配置通道
	} {
		if !strings.Contains(out, want) {
			t.Errorf("正常内容丢失: %s\n%s", want, out)
		}
	}
}

// TestSanitizeReportMerge write_deck 一次提交多页，警告要汇总成一条。
func TestSanitizeReportMerge(t *testing.T) {
	a := sanitizeReport{Handlers: 1, BadURLs: 1}
	a.merge(sanitizeReport{BannedTags: []string{"iframe"}, Handlers: 2})
	a.merge(sanitizeReport{BannedTags: []string{"iframe", "form"}, BadURLs: 3})

	if a.Handlers != 3 || a.BadURLs != 4 {
		t.Fatalf("计数合并错误: %+v", a)
	}
	if len(a.BannedTags) != 2 || a.BannedTags[0] != "form" || a.BannedTags[1] != "iframe" {
		t.Fatalf("标签应去重并排序: %+v", a.BannedTags)
	}
}

// TestSanitizeRejectMessage 错误信息要能指导模型怎么改（尤其是样式该写哪里）。
func TestSanitizeRejectMessage(t *testing.T) {
	rep, _ := sanitizeHTML(t, `<section><style>.x{}</style><iframe src="//e.com"></iframe></section>`)
	err := rep.RejectErr()
	if err == nil {
		t.Fatal("应被拒绝")
	}
	msg := err.Error()
	for _, want := range []string{"<style>", "<iframe>", "请移除", "组件库 class"} {
		if !strings.Contains(msg, want) {
			t.Errorf("错误信息缺少 %q: %s", want, msg)
		}
	}
}

func TestIsDangerousURL(t *testing.T) {
	bad := []string{
		"javascript:alert(1)", "JavaScript:alert(1)", " java\tscript:alert(1)",
		"\nvbscript:x", "data:text/html,<b>", "DATA:Text/HTML,x",
	}
	good := []string{
		"", "#page2", "https://example.com/doc", "/assets/a.png",
		"data:image/png;base64,iVBOR", "mailto:a@b.com", "javascriptx:y",
	}
	for _, v := range bad {
		if !isDangerousURL(v) {
			t.Errorf("应判定为危险: %q", v)
		}
	}
	for _, v := range good {
		if isDangerousURL(v) {
			t.Errorf("不应判定为危险: %q", v)
		}
	}
}
