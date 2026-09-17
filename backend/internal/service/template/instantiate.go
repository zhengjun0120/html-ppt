package template

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

// Instantiate 把模板变成一份 deck 的初始文件集（D13 混合存储的"模板侧"）：
//
//	index.html —— demo 页面被剥离，只留骨架 + 挂载标记（正式页面由生成管线写入）；
//	style.css  —— 原样拷贝进 deck 目录（快照隔离：模板日后改版不影响存量 deck，
//	              "恢复历史版本"恢复的才是当时的样子）。
//
// runtime.js / base.css / 字体不在拷贝范围：它们以 /assets/deck-v2/* 绝对路径
// 共享引用，升级集中生效；真有破坏性升级时发新路径（deck-v2.1），老 deck 不动。
type Instantiation struct {
	IndexHTML string
	StyleCSS  string
	Canvas    Canvas
}

// bodyTagRe 定位 <body ...> 开标签（变体 class 挂在这里）。
var bodyTagRe = regexp.MustCompile(`<body([^>]*)>`)

// Instantiate 产出 deck 初始文件。
//
//	variantID 变体（"" 或 "default" = 模板默认观感）；
//	title     deck 标题（写入 <title>，做 HTML 转义）。
func (t *Template) Instantiate(variantID, title string) (*Instantiation, error) {
	if variantID == "" {
		variantID = "default"
	}
	variantClass, ok := t.VariantClass(variantID)
	if !ok {
		return nil, fmt.Errorf("模板 %s 没有变体 %q", t.ID, variantID)
	}

	out := t.indexHTML

	// 1. 剥离 demo 页面，保留两个标记本身——它们是 deck 写路径定位 AGENT 可写区
	//    的锚点，实例化后的文件里必须还在。
	start := strings.Index(out, slidesStartMarker)
	end := strings.Index(out, slidesEndMarker)
	if start < 0 || end < 0 || end < start {
		return nil, fmt.Errorf("挂载标记损坏（这不应该发生：加载时校验过）")
	}
	inner := out[start+len(slidesStartMarker) : end]
	out = out[:start+len(slidesStartMarker)] + "\n" + out[end:]

	// demo 区间内除了 section 之外的东西（注释等）也一并丢弃了——
	// 挂载区间就是"内容区"，里面不该有任何生成流程关心的东西。
	_ = inner

	// 2. 变体 class 挂到 body。模板的 .tpl-<id> 类必须保留（整个设计系统以它作用域）。
	out = applyBodyClass(out, variantClass)

	// 3. 标题。
	out = replaceTitle(out, title)

	return &Instantiation{
		IndexHTML: out,
		StyleCSS:  t.styleCSS,
		Canvas:    t.Canvas,
	}, nil
}

// applyBodyClass 把 extraClass 追加到 <body> 的 class 属性（无该属性则创建）。
// extraClass 为空时原样返回。
func applyBodyClass(doc, extraClass string) string {
	if extraClass == "" {
		return doc
	}
	m := bodyTagRe.FindStringSubmatchIndex(doc)
	if m == nil {
		return doc // 骨架必须有 body；没有就放弃变体（校验过的模板到不了这里）
	}
	attrs := doc[m[2]:m[3]]
	if cm := classAttrRe.FindStringSubmatch(attrs); cm != nil {
		newAttrs := strings.Replace(attrs, cm[0], `class="`+cm[1]+" "+extraClass+`"`, 1)
		return doc[:m[0]] + "<body" + newAttrs + ">" + doc[m[1]:]
	}
	// body 没有 class 属性：插一个在最前
	return doc[:m[0]] + `<body class="` + extraClass + `"` + attrs + ">" + doc[m[1]:]
}

var classAttrRe = regexp.MustCompile(`class="([^"]*)"`)

// replaceTitle 替换 <title> 文本（HTML 转义，防标题里的 <>& 破坏骨架）。
func replaceTitle(doc, title string) string {
	return titleTagRe.ReplaceAllStringFunc(doc, func(s string) string {
		return "<title>" + html.EscapeString(strings.TrimSpace(title)) + "</title>"
	})
}

var titleTagRe = regexp.MustCompile(`(?s)<title>.*?</title>`)
