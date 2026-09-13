package deck

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// bannedTags 会带来执行、加载或数据渗出通道的标签，属于"结构性违规"：
// 整页拒绝，让模型重写。属性级违规（on* 事件、危险协议 URL）走剥除 + 回报，
// 两条策略的分界见 sanitizeReport。
//
//	script/style          执行代码 / 定义样式（原有规则，集中到这里统一判定）
//	iframe/frame/frameset 内嵌外部文档
//	object/embed/applet   插件通道（CSP 里 object-src 'none' 挡预览，导出后没人挡）
//	form                  表单提交是"导航"不是 fetch，connect-src 管不到；
//	                      导出成单文件后更没有 CSP，是一条现成的外发通道
//	base                  劫持相对 URL 的基准，能把同源引用指到外站
//	meta                  http-equiv refresh 可做跳转
//	link                  外部样式/预加载通道
var bannedTags = []string{
	"script", "style",
	"iframe", "frame", "frameset", "applet", "object", "embed",
	"form", "base", "meta", "link",
}

// dangerousSchemes 前缀（已做过空白剥离与转小写）。
// 只列会在文档里执行或替换文档的协议；data:image/* 是合法的内联图片，不在其中。
var dangerousSchemes = []string{
	"javascript:", "vbscript:", "data:text/html", "data:application/xhtml",
}

type sanitizeReport struct {
	BannedTags []string // 命中的危险标签（去重、字典序）
	Handlers   int      // 剥除的内联事件属性数（on*）
	BadURLs    int      // 剥除的危险协议 URL 数
}

func (r sanitizeReport) Clean() bool {
	return len(r.BannedTags) == 0 && r.Handlers == 0 && r.BadURLs == 0
}

// merge 合并多页的报告（write_deck 一次提交多页，警告要汇总成一条）。
func (r *sanitizeReport) merge(o sanitizeReport) {
	r.Handlers += o.Handlers
	r.BadURLs += o.BadURLs
	seen := map[string]bool{}
	for _, t := range r.BannedTags {
		seen[t] = true
	}
	for _, t := range o.BannedTags {
		if !seen[t] {
			r.BannedTags = append(r.BannedTags, t)
			seen[t] = true
		}
	}
	sort.Strings(r.BannedTags)
}

// RejectErr 结构性违规 → 拒绝写入，让模型重写整页。没有结构性违规时返回 nil。
func (r sanitizeReport) RejectErr() error {
	if len(r.BannedTags) == 0 {
		return nil
	}
	names := make([]string, len(r.BannedTags))
	hasCodeTag := false
	for i, t := range r.BannedTags {
		names[i] = "<" + t + ">"
		if t == "script" || t == "style" {
			hasCodeTag = true
		}
	}
	msg := fmt.Sprintf("幻灯片不支持 %s：这些标签会执行代码、加载外部内容，或形成数据外发通道（导出成单文件后没有任何浏览器策略兜底），请移除后重新提交",
		strings.Join(names, "、"))
	if hasCodeTag {
		// 保留原有的具体指引：模型需要知道"样式该写哪里"
		msg += "。样式请用组件库 class 或内联 style，页面不执行脚本"
	}
	return errors.New(msg)
}

// Warning 属性级违规的提示语，附在成功的工具结果里。
// 刻意不静默：模型以为 onclick 生效、实际已被剥掉，是最典型的"静默失败"——
// 它会向用户汇报"已加上点击交互"，而页面毫无变化。
func (r sanitizeReport) Warning() string {
	if r.Handlers == 0 && r.BadURLs == 0 {
		return ""
	}
	var parts []string
	if r.Handlers > 0 {
		parts = append(parts, fmt.Sprintf("%d 处内联事件属性（on*）", r.Handlers))
	}
	if r.BadURLs > 0 {
		parts = append(parts, fmt.Sprintf("%d 处危险协议的链接", r.BadURLs))
	}
	return "已自动移除 " + strings.Join(parts, "、") + "（为安全起见不写入页面）。" +
		"幻灯片是声明式的：不支持内联事件与脚本，交互请改用组件库 class + data-*；" +
		"若这处改动对方案很关键，请换一种不依赖脚本的实现方式。"
}

// sanitizeSlide 就地清理一页内容并返回报告，是消毒的唯一判定点。
// 调用方按策略处置：结构性违规 → RejectErr 拒绝；属性级 → 放行并把 Warning 带回工具结果。
//
// 剥除与报告同时进行：即使调用方选择拒绝写入，被清理过的树也不会留下脏东西。
// 这是"校验永远先行"之外的第二道保险——将来某个新写入路径忘了判定，也不会漏出去。
func sanitizeSlide(sec *goquery.Selection) sanitizeReport {
	rep := sanitizeReport{}

	// 1) 危险标签：先统计（报告需要标签名），再整棵移除
	hit := map[string]bool{}
	sec.Find(strings.Join(bannedTags, ",")).Each(func(_ int, s *goquery.Selection) {
		if len(s.Nodes) == 0 {
			return
		}
		hit[strings.ToLower(s.Nodes[0].Data)] = true
	})
	for t := range hit {
		rep.BannedTags = append(rep.BannedTags, t)
	}
	sort.Strings(rep.BannedTags)
	if len(rep.BannedTags) > 0 {
		sec.Find(strings.Join(rep.BannedTags, ",")).Remove()
	}

	// 2) on* 事件属性 + 危险协议 URL（根节点自身也要扫，如 <section onclick=...>）
	clean := func(s *goquery.Selection) {
		if len(s.Nodes) == 0 {
			return
		}
		// 复制一份再改：RemoveAttr 会改 Attr 切片，边遍历边改不安全
		for _, a := range append([]html.Attribute(nil), s.Nodes[0].Attr...) {
			name := strings.ToLower(a.Key)
			if strings.HasPrefix(name, "on") && len(name) > 2 {
				s.RemoveAttr(a.Key)
				rep.Handlers++
				continue
			}
			if isDangerousURL(a.Val) {
				s.RemoveAttr(a.Key)
				rep.BadURLs++
			}
		}
	}
	clean(sec)
	sec.Find("*").Each(func(_ int, s *goquery.Selection) { clean(s) })

	return rep
}

// isDangerousURL 判断属性值是否是危险协议。
// 空白与控制字符必须先去干净再比：浏览器解析协议时会忽略它们，
// 所以 `java\tscript:alert(1)`、` javascript:...` 都是可执行的真绕过。
func isDangerousURL(v string) bool {
	var b strings.Builder
	for _, r := range strings.ToLower(v) {
		if r <= ' ' || r == '\u007f' {
			continue
		}
		b.WriteRune(r)
	}
	s := b.String()
	for _, p := range dangerousSchemes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}
