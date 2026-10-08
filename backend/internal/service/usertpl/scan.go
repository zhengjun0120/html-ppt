package usertpl

// 模板内容安全预检（docs/user-template-history-plan.md §3.3）。
//
// 从 service.go 的 securityScan 拆出的纯函数：同样的黑名单语义，但以内容为参数、
// 在"每次写入"（对话 write_style/write_demo、手动保存）时就拦，而不是等发布门禁
// 才发现改坏了。发布门禁仍组合调用它们做终审。
//
// 对 deck-v2 沙箱模型的关系：demo 会被非属主渲染（published 公开），脚本是唯一
// 执行通道，因此 index.html 只允许 runtime.js 一个脚本来源；CSS 的外链（url(/
// @import）是数据渗出通道，一律禁——唯一例外是 data: 内联（无网络请求，
// brutalist-bold 的纸纹底即此用法，fork 副本因此可以过门禁）。

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// scanCSS style.css 安全预检。禁 @import / expression / behavior / -moz-binding；
// url( 仅放行 data: 内联，网络引用一律拒。非法 UTF-8 一并拒——bundle 快照经
// json.Marshal 会把坏字节替换成 U+FFFD，磁盘与历史从此对不上（实测教训）。
func scanCSS(css string) error {
	if !utf8.ValidString(css) {
		return fmt.Errorf("style.css 含非法 UTF-8 字节")
	}
	low := strings.ToLower(css)
	for i := 0; ; {
		j := strings.Index(low[i:], "url(")
		if j < 0 {
			break
		}
		k := i + j + 4 // 指向 "(" 之后
		for k < len(low) && (low[k] == '"' || low[k] == '\'' || low[k] == ' ') {
			k++
		}
		if !strings.HasPrefix(low[k:], "data:") {
			return fmt.Errorf("style.css 含被禁用的 url() 外链（仅允许 data: 内联资源，外链是数据渗出通道）")
		}
		i = k
	}
	for _, bad := range []string{"@import", "expression(", "behavior:", "-moz-binding"} {
		if strings.Contains(low, bad) {
			return fmt.Errorf("style.css 含被禁用的 %q（外链/表达式是数据渗出通道）", bad)
		}
	}
	return nil
}

// scanHTML index.html 安全预检。禁 runtime.js 之外的一切脚本与 iframe/object/embed、
// 禁内联事件；并要求 demo 的结构前提：有 <section>（页面集）、引用
// /assets/deck-v2/runtime.js（没有它 demo 无法翻页/单页化）、head 里保留
// style.css 链接与 /assets/deck-v2/ 样式引用——demo 端点靠前者重写到受控资产
// 地址，后者承载翻页堆叠与字号底线（write_demo 丢 head link 的实测事故
// 2026-10-06：base.css 没了，所有 slide 掉进文档流，翻页失效且注册表校验不查）。
func scanHTML(html string) error {
	if !utf8.ValidString(html) {
		return fmt.Errorf("index.html 含非法 UTF-8 字节")
	}
	low := strings.ToLower(html)
	// 每个 <script> 开标签都必须指向 runtime.js（旧实现只查第一个 script，
	// 夹在 runtime.js 之前的额外脚本能溜过去——write_demo 全文重写后必须收紧）
	for rest := low; ; {
		idx := strings.Index(rest, "<script")
		if idx < 0 {
			break
		}
		seg := rest[idx:]
		if tagEnd := strings.Index(seg, ">"); tagEnd >= 0 {
			seg = seg[:tagEnd]
		}
		if !strings.Contains(seg, "/assets/deck-v2/runtime.js") {
			return fmt.Errorf("index.html 含被禁用的脚本（只允许 /assets/deck-v2/runtime.js）")
		}
		rest = rest[idx+len("<script"):]
	}
	for _, tag := range []string{"<iframe", "<object", "<embed"} {
		if strings.Contains(low, tag) {
			return fmt.Errorf("index.html 含被禁用的元素 %q", tag)
		}
	}
	for _, ev := range []string{" onload=", " onerror=", " onclick="} {
		if strings.Contains(low, ev) {
			return fmt.Errorf("index.html 含内联事件 %q", strings.TrimSpace(ev))
		}
	}
	if !strings.Contains(low, "<section") {
		return fmt.Errorf("index.html 不含 <section>（demo 必须是页面集）")
	}
	if !strings.Contains(low, "/assets/deck-v2/runtime.js") {
		return fmt.Errorf("index.html 缺少 /assets/deck-v2/runtime.js 引用（没有它 demo 无法翻页）")
	}
	if !strings.Contains(low, `href="style.css"`) {
		return fmt.Errorf(`index.html 缺 href="style.css" 引用（demo 端点靠它重写到受控资产地址，丢了整页没样式）`)
	}
	// 必须有 link 标签引用 /assets/deck-v2/ 样式——runtime.js 的 script 也含这个
	// 子串，所以按 <link + href 联合匹配，不能用裸子串判
	if !linkAssetRe.MatchString(low) {
		return fmt.Errorf("index.html 缺 /assets/deck-v2/ 样式引用（head 里的 base.css 等 link 必须保留，丢了翻页堆叠与字号底线全失效）")
	}
	return nil
}

// linkAssetRe head 里指向框架样式资产的 link 标签。
var linkAssetRe = regexp.MustCompile(`<link[^>]+href="[^"]*/assets/deck-v2/`)
