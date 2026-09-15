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

// 由程序处理的量测结果
func HardFindings(d *Deck) []string {
	var out []string
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
			out = append(out, fmt.Sprintf("第 %d 页：%s", s.Index+1, strings.Join(why, "；")))
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
			out = append(out, fmt.Sprintf("版式重复：%s 用了 %d 次（第 %v 页，规则 ≤%d）", clip(l, 40), len(pages[l]), pages[l], layoutRepeatMax))
		}
	}
	for i := 1; i < len(d.Slides); i++ {
		if d.Slides[i].Layout == d.Slides[i-1].Layout {
			out = append(out, fmt.Sprintf("连续两页同一版式：第 %d 页与第 %d 页", i, i+1))
		}
	}
	return out
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
