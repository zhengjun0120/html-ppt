package deck

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// deck 自定义样式槽：<style id="deck-custom">，级联在组件库与主题 override 之后。
//
// 为什么是独立槽而不是让 AI 在页面里写 <style>：
//   1. 页面内容区必须保持声明式（禁 script/style 的闸门不放松），槽是"框架层"的一个洞；
//   2. 槽在 head、不在任何 section 内，所以它不进版本 diff 的页面部分，
//      但整份 deck 的历史快照照样覆盖它——写坏了能一键回滚；
//   3. 导出成单文件时它天然在内，不需要额外处理。
//
// 写入是**整体替换制**（不做增量 diff）：AI 想改就必须先 read 再提交全文，
// 这样"当前样式是什么"永远只有一个权威来源，不会出现"片段式追加"累积出的意外样式。
const customCSSBlockID = "deck-custom"

// maxCustomCSSBytes 槽容量上限。它同时是"防止模型写出一大坨"的护栏和 token 上限。
const maxCustomCSSBytes = 32 * 1024

// customCSSBlockHTML 骨架里的槽（老 deck 没有这个块，见 ensureCustomCSSBlock）。
const customCSSBlockHTML = `<style id="deck-custom"></style>`

// validateCustomCSS 清洗并校验 AI 提交的自定义 CSS。
//
// 这里挡的是"通道""无效写法"和"没主的变量"，不是"风格"：可以改配色、间距、圆角、字号；
// 但不允许通过 CSS 再开一条加载/执行/外发的路（@import/url()/</style>），
// 也不允许写注定不生效的东西——花括号不配对（会让后面所有规则一起失效）、
// 重定义契约变量（会静默盖住 update_theme）、引用没人消费的变量名（拼错/凭空造，静默无效）。
// 风格偏离由用户看着预览决定（不满意就回滚历史），不属于清洗该管的事。
//
// known = 被框架 CSS 消费过的变量；preDefined = 该 deck 主题块里已定义过的变量。
// known 为 nil 表示契约不可得，跳过回声校验（fail-open：这项检查防的是"静默无效"，
// 不是安全问题，不能因为资源目录没配好就把所有写入卡死）。
func validateCustomCSS(css string, known, preDefined map[string]bool) (string, error) {
	clean := strings.TrimSpace(css)
	if clean == "" {
		return "", nil // 空 = 清空自定义样式，合法
	}
	if len(clean) > maxCustomCSSBytes {
		return "", fmt.Errorf("自定义 CSS 超过 %d KB 上限（当前 %d KB）：请只保留真正必要的覆盖，把结构性改动放回组件库 class",
			maxCustomCSSBytes/1024, len(clean)/1024)
	}

	lower := strings.ToLower(clean)
	// </style 会提前终止 raw-text 样式块（写进去之后整个文档结构都会乱），必须拒绝
	if strings.Contains(lower, "</style") {
		return "", errors.New("CSS 里不能出现 </style：它会提前结束样式块、破坏整份文档结构，请移除")
	}
	if strings.Contains(lower, "@import") {
		return "", errors.New("不支持 @import：它会加载外部样式表（既是外部依赖也是数据外发通道）。需要什么样式请直接写在槽里")
	}
	if strings.Contains(lower, "url(") {
		return "", errors.New("不支持 url()：它会发起外部请求。图片请用 <img src=\"/assets/...\"> 或 data: 内联，字体内嵌留到字体功能里统一做")
	}
	if strings.Contains(lower, "expression(") {
		return "", errors.New("不支持 expression()：这是可执行代码的老写法，请移除")
	}
	// 语法底线：花括号不配对会让后面所有规则一起失效，是最值得拦的语法错误
	if err := checkCSSBalance(clean); err != nil {
		return "", err
	}
	// 契约变量不许重定义：自定义槽排在主题块之后、同级选择器靠后者取胜，
	// 在这里定义会静默盖掉 update_theme（详见 theme_vars.go 的说明）。
	// 这不是风格问题而是"谁的权威"问题，所以走拒绝而不是剥除——
	// 悄悄删掉这个变量会让 AI 精心写的整套配色缺一块，页面反而更坏，
	// 不如直接报错让它换 update_theme 的 vars 重写。
	if name := findReservedVarRedefinition(clean); name != "" {
		return "", reservedVarRedefinitionErr(name)
	}
	// 回声校验：引用的变量必须真的有元素消费它（或在同一份 CSS / 主题块里定义过）。
	// 拒绝而不是提示：那处样式注定不生效，放行就是"写了 CSS 却没反应"的假成功。
	if known != nil {
		defined := definedVariables(clean)
		for k := range preDefined {
			defined[k] = true
		}
		if unknown := unknownVariables(clean, known, defined); len(unknown) > 0 {
			return "", unknownVarErr(unknown, known)
		}
	}
	return clean, nil
}

// ensureCustomCSSBlock 保证槽存在并返回它。
// 骨架只在 write_deck 时固化一次，所以老 deck 没有这个块——碰到时补上，
// 这是最轻量的 schema 迁移（升级动作放在读写路径上，不写独立迁移脚本），
// 和 ensureThemeBlocks 是同一套路。
//
// 位置约束：必须排在 #deck-theme-override 之后（自定义样式的优先级最高）。
func ensureCustomCSSBlock(doc *goquery.Document) *goquery.Selection {
	if sel := doc.Find("#" + customCSSBlockID); sel.Length() > 0 {
		return sel.First()
	}
	if override := doc.Find("#deck-theme-override"); override.Length() > 0 {
		override.First().AfterHtml(customCSSBlockHTML)
	} else {
		doc.Find("head").AppendHtml(customCSSBlockHTML)
	}
	return doc.Find("#" + customCSSBlockID).First()
}

// ReadCustomCSS 读当前槽内容（工具 read_custom_css 用）。
// 整体替换制下这是写前的必读步骤——不然 AI 会拿一份想象中的旧内容覆盖真实内容。
func (s *Service) ReadCustomCSS(userID uint, deckID string) (string, error) {
	if err := s.authorize(userID, deckID); err != nil {
		return "", err
	}
	return s.readCustomCSSRaw(deckID)
}

// readCustomCSSRaw 拆出来供白盒测试用（绕开 authorize，st=nil 也能跑）。
func (s *Service) readCustomCSSRaw(deckID string) (string, error) {
	raw, err := s.readRaw(deckID)
	if err != nil {
		return "", err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(raw))
	if err != nil {
		return "", fmt.Errorf("解析 deck 失败 err:%w", err)
	}
	sel := doc.Find("#" + customCSSBlockID)
	if sel.Length() == 0 {
		return "", nil // 老 deck 还没有槽：等价于"没有自定义样式"
	}
	return strings.TrimSpace(sel.First().Text()), nil
}

// UpdateCustomCSS 整体替换自定义样式。返回清洗后的最终内容（空串表示已清空）。
// 不做指纹校验：槽只有一个写入者（agent），且整份 deck 的历史快照覆盖它，
// 冲突代价低——和 update_theme 同一取舍。
func (s *Service) UpdateCustomCSS(userID uint, deckID, css string) (string, error) {
	if err := s.authorize(userID, deckID); err != nil {
		return "", err
	}
	return s.updateCustomCSSLocked(deckID, css)
}

// updateCustomCSSLocked 写入核心（锁内），拆出来同样为白盒测试。
func (s *Service) updateCustomCSSLocked(deckID, css string) (string, error) {
	unlock := s.lockDeck(deckID)
	defer unlock()

	raw, err := s.readRaw(deckID)
	if err != nil {
		return "", err
	}

	// 回声校验的"已知变量"= 框架契约（被组件库/reveal 消费过的）
	//                        ∪ 该 deck 主题块里定义过的（即 update_theme 的 vars 调色板）
	// 第二项不能少：AI 在 vars 里定义了 --surface、再在槽里 var(--surface) 是正常用法，
	// 漏了它就会把正常写法误判成"消费不到的变量"。
	known := s.variableContract()
	preDefined := definedVariables(deckStyleBlocks(raw))

	clean, err := validateCustomCSS(css, known, preDefined)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(raw))
	if err != nil {
		return "", fmt.Errorf("解析 deck 失败 err:%w", err)
	}

	sel := ensureCustomCSSBlock(doc)
	// setRawText 走裸文本节点：绕过 goquery SetText 在 raw-text 元素上的实体转义问题
	// （theme.go 里踩过，注释在那儿）
	setRawText(sel, clean)

	out, err := goquery.OuterHtml(doc.Selection)
	if err != nil {
		return "", fmt.Errorf("序列化 deck 失败 err:%w", err)
	}
	if err := s.atomicWriteDeck(deckID, out); err != nil {
		return "", fmt.Errorf("写入 deck 失败 err:%w", err)
	}
	return clean, nil
}
