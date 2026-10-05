package usertpl

// AI 味 lint 的定制侧接入。词表与判定复用 deck.LintTaste（taste-skill 的
// AI-Tells 翻译版，生成侧 write_pages 闸门同源）——定制侧三处消费：
//
//  1. write_demo：整份 demo 逐页过 lint，提示随工具返回给模型自修；
//  2. add_layout：新增示例页单页 lint，提示进返回摘要；
//  3. Checkup：体检报告加 taste 段，发布前可见。
//
// 全部提示级——品味问题不阻塞写入，与生成侧同一哲学。demo 示例文案是生成时
// 模型模仿的范本，这里的 AI 腔会传染给每一次生成，所以定制侧拦截比生成侧更
// 靠前、性价比更高。

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"html-ppt/backend/internal/service/deck"
)

// maxTasteHints 工具返回/报告里的提示行上限（9 页 demo 全脏时防刷屏）。
const maxTasteHints = 12

// formatLint 单条提示的人读形态（label = "第 N 页" / "示例页"）。
func formatLint(label string, item deck.LintItem) string {
	line := fmt.Sprintf("%s [%s] %s", label, item.Rule, item.Msg)
	if item.Hint != "" {
		line += "：" + item.Hint
	}
	return line
}

// lintDemoSections 对 demo HTML 的 SLIDES 区间逐页跑 AI 味 lint。
// 返回人读提示行；超过上限折叠成计数行。
func lintDemoSections(html string) []string {
	start := strings.Index(html, slidesStartMark)
	if start < 0 {
		return nil
	}
	relEnd := strings.Index(html[start:], slidesEndMark)
	if relEnd < 0 {
		return nil
	}
	seg := html[start+len(slidesStartMark) : start+relEnd]

	var out []string
	total := 0
	no, rest := 0, seg
	for {
		i := strings.Index(rest, "<section")
		if i < 0 {
			break
		}
		j := strings.Index(rest[i:], "</section>")
		if j < 0 {
			break
		}
		secEnd := i + j + len("</section>")
		frag := rest[i:secEnd]
		no++
		rest = rest[secEnd:]

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(frag))
		if err != nil {
			continue
		}
		for _, item := range deck.LintTaste(doc.Find("section").First()) {
			total++
			if len(out) < maxTasteHints {
				out = append(out, formatLint(fmt.Sprintf("第 %d 页", no), item))
			}
		}
	}
	if total > len(out) {
		out = append(out, fmt.Sprintf("…共 %d 处提示，其余同类", total))
	}
	return out
}

// lintDemoFragment 单个示例 section 的 lint（add_layout 的 demo_html 用）。
func lintDemoFragment(frag string) []string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(frag))
	if err != nil {
		return nil
	}
	var out []string
	for _, item := range deck.LintTaste(doc.Find("section").First()) {
		out = append(out, formatLint("示例页", item))
	}
	return out
}
