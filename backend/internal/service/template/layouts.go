package template

import (
	"regexp"
	"strings"
)

// layouts.md 的解析。
//
// 条目格式（见 tech-sharing/layouts.md）：
//
//	## cover（封面）          ← ## 后第一个 token 是版式 id
//	指纹：hero               ← 可选；该版式的视觉模式（节奏守卫按它判重）
//	数量：step=3             ← 可选；class=次数，逗号分隔多个。登记后写入门
//	                           强制该类元素恰好 N 个（"文字堆一页"的闸）
//	...
//	合法类名：slide, kicker, h1, ...
//	```html
//	<section ...>...</section>
//	```
//
// 解析结果是三个视图：骨架代码（read_layout 用）、类名清单（契约校验用）、原始段落。

type layoutEntry struct {
	Skeleton string   // ```html 围栏内的骨架代码
	Classes  []string // 「合法类名：」行声明的清单
	Pattern  string   // 「指纹：」行声明的视觉模式（hero/stack/cards/split/code/table/chart/quote）
	Repeats  []repSpec// 「数量：」行声明的精确数量契约（可选）
	Body     string   // 原始段落（含头部说明）
}

// repSpec 单条数量契约：该类元素在页面里必须恰好出现 Count 次。
type repSpec struct {
	Class string
	Count int
}

var (
	layoutHeadRe  = regexp.MustCompile(`(?m)^##\s+([A-Za-z][A-Za-z0-9_-]*)`)
	layoutClassRe = regexp.MustCompile(`合法类名[：:]\s*(.+)`)
	layoutPtrnRe  = regexp.MustCompile(`指纹[：:]\s*(.+)`)
	layoutRepRe   = regexp.MustCompile(`数量[：:]\s*(.+)`)
	fenceRe       = regexp.MustCompile("(?s)```html\\s*\n(.*?)```")
)

// patternOf 版式的视觉模式指纹。显式声明优先；缺失时从骨架类名兜底推断
//（第三方模板没写指纹也能工作，只是不如显式声明准）。
// 推断顺序有讲究：full/center 最特异（满版居中必是 hero），code 次之，
// grid 再次（卡片阵），sidebar/main 是分栏特征，其余一律 stack（纵向列表
// 是最常见的退化形态，把它当默认值能让缺指纹的模板立即受到节奏保护）。
func patternOf(e layoutEntry) string {
	if e.Pattern != "" {
		return e.Pattern
	}
	s := e.Skeleton
	switch {
	case strings.Contains(s, "full") || strings.Contains(s, `"center`):
		return "hero"
	case strings.Contains(s, `"code`):
		return "code"
	case strings.Contains(s, "grid"):
		return "cards"
	case strings.Contains(s, "sidebar") || strings.Contains(s, `"main`):
		return "split"
	default:
		return "stack"
	}
}

// skeletonClassCounts 统计骨架里每个类名出现的次数（按 class 属性逐 token 计）。
// 「数量：」行的交叉校验用：声明的次数必须与骨架一致，否则契约自相矛盾。
func skeletonClassCounts(skeleton string) map[string]int {
	counts := map[string]int{}
	for _, m := range htmlClassAttrRe.FindAllStringSubmatch(skeleton, -1) {
		for _, c := range strings.Fields(m[1]) {
			if isClassToken(c) {
				counts[c]++
			}
		}
	}
	return counts
}

// parseLayoutsMD 解析 layouts.md，返回 版式 id → 条目。
// id 取自 "## " 标题行的第一个 token（括号里的中文名不是 id）。
func parseLayoutsMD(md string) map[string]layoutEntry {
	out := map[string]layoutEntry{}

	heads := layoutHeadRe.FindAllStringSubmatchIndex(md, -1)
	for i, h := range heads {
		id := md[h[2]:h[3]]
		end := len(md)
		if i+1 < len(heads) {
			end = heads[i+1][2]
		}
		body := md[h[0]:end]

		entry := layoutEntry{Body: body}
		if m := fenceRe.FindStringSubmatch(body); m != nil {
			entry.Skeleton = strings.TrimSpace(m[1])
		}
		if m := layoutClassRe.FindStringSubmatch(body); m != nil {
			for _, c := range strings.FieldsFunc(m[1], func(r rune) bool {
				return r == ',' || r == '，' || r == ' ' || r == '\t'
			}) {
				if c = strings.TrimSpace(c); isClassToken(c) {
					entry.Classes = append(entry.Classes, c)
				}
			}
		}
		if m := layoutPtrnRe.FindStringSubmatch(body); m != nil {
			if p := strings.TrimSpace(m[1]); isClassToken(p) {
				entry.Pattern = p
			}
		}
		if m := layoutRepRe.FindStringSubmatch(body); m != nil {
			for _, pair := range strings.FieldsFunc(m[1], func(r rune) bool {
				return r == ',' || r == '，' || r == ' ' || r == '\t'
			}) {
				kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
				if len(kv) != 2 {
					continue // 格式坏的条目交给 loadTemplate 的交叉校验报错，解析从宽
				}
				n := 0
				for _, r := range strings.TrimSpace(kv[1]) {
					if r < '0' || r > '9' {
						n = 0
						break
					}
					n = n*10 + int(r-'0')
				}
				if c := strings.TrimSpace(kv[0]); isClassToken(c) && n > 0 {
					entry.Repeats = append(entry.Repeats, repSpec{Class: c, Count: n})
				}
			}
		}
		out[id] = entry
	}
	return out
}
