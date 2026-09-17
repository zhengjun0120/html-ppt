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
	Body     string   // 原始段落（含头部说明）
}

var (
	layoutHeadRe  = regexp.MustCompile(`(?m)^##\s+([A-Za-z][A-Za-z0-9_-]*)`)
	layoutClassRe = regexp.MustCompile(`合法类名[：:]\s*(.+)`)
	fenceRe       = regexp.MustCompile("(?s)```html\\s*\n(.*?)```")
)

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
		out[id] = entry
	}
	return out
}
