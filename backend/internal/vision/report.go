package vision

import (
	"fmt"
	"sort"
	"strings"
)

const (
	budgetChars     = 200 //整页文字上限
	budgetChildren  = 5   //一个section的直接子元素上限
	floorFontPx     = 31  //24pt 投影下限
	overflowLimit   = 1.02
	fitWarnScale    = 0.90 //被兜底缩到这个比例以下，投屏上已经明显偏小
	layoutRepeatMax = 2    //同一版式最多用两次
)

// 判定分两层，因为两层成立的范围不一样：
//   - 页级（LevelPage）：只看这一页的数字。agent 只点了几页时，可以只把这几页的判定
//     摆到模型面前，其余页的判定与它无关。
//   - 版式级（LevelLayout）：要拿**整份** deck 数才成立。"同一版式用了 3 次"这种判定
//     在 3 页的子集里永远数不出来（子集会把它算成"用了 1 次，没超"），所以子集审查
//     必须沿用整份算出来的结果，不能拿子集现算。
const (
	LevelPage   = "page"
	LevelLayout = "layout"
)

type Finding struct {
	Level string
	Pages []int //涉及的页（1 基）
	Text  string
}

// 逐页量测的硬判定 + 版式层面的硬判定。做成结构化而不是拼好的字符串，
// 是因为"只看第 3、7 页"时要能从里面挑出相关的那几条——按页号挑得先有页号，
// 用正则从"第 5 页：…"里反解页号是能写，但那等于把页号信息丢掉再猜回来。
func Analyze(d *Deck) []Finding {
	var page, layout []Finding

	for _, s := range d.Slides {
		var why []string
		if s.EffScale < fitWarnScale {
			why = append(why, fmt.Sprintf("被适配兜底缩到 %.2f（投屏上明显偏小）", s.EffScale))
		}
		if s.MinFontPx < floorFontPx {
			why = append(why, fmt.Sprintf("最小字号 %.0fpx 低于 24pt 投影下限", s.MinFontPx))
		}
		if s.Chars > budgetChars {
			why = append(why, fmt.Sprintf("整页 %d 字超过 %d 字预算", s.Chars, budgetChars))
		}
		if s.Children > budgetChildren {
			why = append(why, fmt.Sprintf("直接子元素 %d 个超过 %d", s.Children, budgetChildren))
		}
		if s.Overflow > overflowLimit {
			why = append(why, fmt.Sprintf("溢出比 %.2f，内容超出画布会被裁", s.Overflow))
		}
		if len(why) > 0 {
			page = append(page, Finding{
				Level: LevelPage, Pages: []int{s.Index + 1},
				Text: fmt.Sprintf("第 %d 页：%s", s.Index+1, strings.Join(why, "；")),
			})
		}
	}

	pages := map[string][]int{}
	order := []string{}
	for _, s := range d.Slides {
		if _, seen := pages[s.Layout]; !seen {
			order = append(order, s.Layout)
		}
		pages[s.Layout] = append(pages[s.Layout], s.Index+1)
	}
	sort.SliceStable(order, func(i, j int) bool { return len(pages[order[i]]) > len(pages[order[j]]) })
	for _, l := range order {
		if len(pages[l]) > layoutRepeatMax {
			layout = append(layout, Finding{
				Level: LevelLayout, Pages: pages[l],
				Text: fmt.Sprintf("版式重复：%s 用了 %d 次（第 %v 页，规则 ≤%d）", clip(l, 40), len(pages[l]), pages[l], layoutRepeatMax),
			})
		}
	}
	for i := 1; i < len(d.Slides); i++ {
		if d.Slides[i].Layout == d.Slides[i-1].Layout {
			layout = append(layout, Finding{
				Level: LevelLayout, Pages: []int{i, i + 1},
				Text: fmt.Sprintf("连续两页同一版式：第 %d 页与第 %d 页", i, i+1),
			})
		}
	}
	return append(page, layout...)
}

// 由程序处理的量测结果
func HardFindings(d *Deck) []string {
	return Texts(Analyze(d))
}

func Texts(fs []Finding) []string {
	out := make([]string, 0, len(fs))
	for _, f := range fs {
		out = append(out, f.Text)
	}
	return out
}

// 挑出与 pages（1 基）相关的判定：页级看自己那页在不在选里，
// 版式级看它涉及的页和选里的页有没有交集。保持 Analyze 的原有顺序。
func ScopeFindings(fs []Finding, pages []int) []string {
	want := map[int]bool{}
	for _, p := range pages {
		want[p] = true
	}
	var out []string
	for _, f := range fs {
		for _, p := range f.Pages {
			if want[p] {
				out = append(out, f.Text)
				break
			}
		}
	}
	return out
}

// 按 1 基页号取子集。返回的 Slide 保留原始 Index，所以报告里说"第 5 页"仍然指的是
// 整份 deck 的第 5 页，而不是子集里的第 2 张——错位的报告比没有报告更坏。
// 第二个返回值是"你要了但这份 deck 没有"的页号，调用方必须把它告诉模型：
// 静默丢掉的表现是模型以为那页看过了。
func (d *Deck) Select(pages []int) (*Deck, []int) {
	want := map[int]bool{}
	for _, p := range pages {
		want[p] = true
	}
	out := &Deck{CanvasW: d.CanvasW, CanvasH: d.CanvasH}
	hit := map[int]bool{}
	for _, s := range d.Slides { //按 deck 顺序走，出去的报告就是从前往后读的
		if want[s.Index+1] {
			out.Slides = append(out.Slides, s)
			hit[s.Index+1] = true
		}
	}
	var missing []int
	for p := range want {
		if !hit[p] {
			missing = append(missing, p)
		}
	}
	sort.Ints(missing)
	return out, missing
}

// 给模型看的每一页的摘要
func Digest(d *Deck) string {
	var b strings.Builder
	for _, s := range d.Slides {
		fmt.Fprintf(&b, "第 %d 页：fit=%.2f 最小字号=%.0fpx 溢出=%.2f 子元素=%d 字数=%d", s.Index+1, s.EffScale, s.MinFontPx, s.Overflow, s.Children, s.Chars)
		if s.Title != "" {
			fmt.Fprintf(&b, " 标题=%q", s.Title)
		}
		b.WriteString("\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

// 给人看的摘要
func (d *Deck) Text() string {
	var b strings.Builder
	fmt.Fprintf(&b, "共 %d 页，画布 %d*%d\n", len(d.Slides), d.CanvasW, d.CanvasH)
	b.WriteString(Digest(d))
	if f := HardFindings(d); len(f) > 0 {
		b.WriteString("\n硬判定：\n  ")
		b.WriteString(strings.Join(f, "\n  "))
	} else {
		b.WriteString("\n硬判定: 无")
	}
	return b.String()
}

func clip(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "..."
}
