package deck

// deck-v2 的页操作层：生成批写入（write_pages）与迭代期页级修改。
//
// 所有写入都走同一道闸门：parseSlideFragmentV2（结构 + 消毒）→ 类名契约（C202）→
// 版式登记校验（C201）→ 段内重编号 → 原子写回。
//
// index.html 的编辑采用**段拼接**而不是整文档 goquery 序列化：
// 模板的 head（字体/CSS 引用/runtime）必须保持字节原样——整文档重序列化
// 会把模板作者手写的属性顺序、自闭合写法全部改掉，diff 噪声与回归风险都大。
// SLIDES:START/END 标记之间的段落才是 goquery 的管辖范围。

import (
	"errors"
	"sort"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"html-ppt/backend/internal/service/template"
)

// ---------- 段读写 ----------

// readSlideSegment 读出挂载标记之间的内容与前后缀。
func readSlideSegment(indexHTML string) (prefix, segment, suffix string, err error) {
	start := strings.Index(indexHTML, slidesStartMarker)
	end := strings.Index(indexHTML, slidesEndMarker)
	if start < 0 || end < 0 || end < start {
		return "", "", "", errors.New("index.html 缺 SLIDES:START/END 挂载标记（骨架损坏）")
	}
	prefix = indexHTML[:start+len(slidesStartMarker)]
	suffix = indexHTML[end:]
	return prefix, indexHTML[start+len(slidesStartMarker) : end], suffix, nil
}

// writeSlideSegment 把段拼回整份 index.html 并原子落盘（锁内调用）。
func (s *Service) writeSlideSegment(deckID, prefix, segment, suffix string) error {
	out := prefix + segment + suffix
	p, err := s.IndexPathV2(deckID)
	if err != nil {
		return err
	}
	return atomicWriteFile(p, []byte(out))
}

// parseSegmentDoc 把段解析成 goquery 文档（段的顶层子元素就是各页 section）。
func parseSegmentDoc(segment string) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(segment))
	if err != nil {
		return nil, fmt.Errorf("解析页面段失败 err:%w", err)
	}
	return doc, nil
}

// ---------- 页清单与读取 ----------

// SlideMetaV2 v2 的页摘要（比 v1 多 layout 字段：节奏/审查/前端进度条都要用）。
type SlideMetaV2 struct {
	Position int    `json:"position"`
	ID       string `json:"id"`
	Title    string `json:"title"`
	Layout   string `json:"layout"`
}

// ListSlidesV2 归属校验后的页清单。
func (s *Service) ListSlidesV2(userID uint, id string) ([]SlideMetaV2, error) {
	if err := s.authorize(userID, id); err != nil {
		return nil, err
	}
	html, err := s.GetHTML(userID, id)
	if err != nil {
		return nil, err
	}
	_, segment, _, err := readSlideSegment(html)
	if err != nil {
		return nil, err
	}
	doc, err := parseSegmentDoc(segment)
	if err != nil {
		return nil, err
	}
	out := make([]SlideMetaV2, 0)
	doc.Find("body").Children().Filter("section").Each(func(i int, sec *goquery.Selection) {
		out = append(out, SlideMetaV2{
			Position: i + 1,
			ID:       sec.AttrOr("data-id", ""),
			Title:    slideTitle(sec),
			Layout:   sec.AttrOr("data-layout", ""),
		})
	})
	return out, nil
}

// ReadSlideV2 返回页 HTML 与指纹（update_slide 的乐观锁用）。
func (s *Service) ReadSlideV2(userID uint, deckID, slideID string) (string, string, error) {
	html, err := s.GetHTML(userID, deckID)
	if err != nil {
		return "", "", err
	}
	sec, err := findSlideV2(html, slideID)
	if err != nil {
		return "", "", err
	}
	outer, err := goquery.OuterHtml(sec)
	if err != nil {
		return "", "", fmt.Errorf("提取 slide html 失败 err:%w", err)
	}
	return outer, slideFingerprint(outer), nil
}

// findSlideV2 在整份 index.html 的段里找指定 data-id 的页。
func findSlideV2(indexHTML, slideID string) (*goquery.Selection, error) {
	_, segment, _, err := readSlideSegment(indexHTML)
	if err != nil {
		return nil, err
	}
	doc, err := parseSegmentDoc(segment)
	if err != nil {
		return nil, err
	}
	sec := doc.Find("body").Children().Filter("section").FilterFunction(func(_ int, n *goquery.Selection) bool {
		return n.AttrOr("data-id", "") == slideID
	})
	if sec.Length() == 0 {
		return nil, fmt.Errorf("slide %q 不存在，请先用 list_slides 获取有效 slide_id", slideID)
	}
	return sec.First(), nil
}

// ---------- 类名契约（C202）----------

// extractUsedClasses 收集 section 及全部后代的 class token（含 section 自身）。
func extractUsedClasses(sec *goquery.Selection) []string {
	seen := map[string]bool{}
	var out []string
	add := func(cls string) {
		for _, c := range strings.Fields(cls) {
			if c != "" && !seen[c] {
				seen[c] = true
				out = append(out, c)
			}
		}
	}
	if cls, exists := sec.Attr("class"); exists {
		add(cls)
	}
	sec.Find("*").Each(func(_ int, el *goquery.Selection) {
		if cls, exists := el.Attr("class"); exists {
			add(cls)
		}
	})
	return out
}

// levenshtein 编辑距离（类名建议用；输入都是短 token，O(mn) 无所谓）。
func levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	prev := make([]int, len(rb)+1)
	cur := make([]int, len(rb)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(ra); i++ {
		cur[0] = i
		for j := 1; j <= len(rb); j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			cur[j] = min3(cur[j-1]+1, prev[j]+1, prev[j-1]+cost)
		}
		prev, cur = cur, prev
	}
	return prev[len(rb)]
}

func min3(a, b, c int) int {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	return m
}

// suggestClass 在合法类名里找编辑距离 ≤2 的最近 suggestion；找不到返回空串。
func suggestClass(unknown string, allowed map[string]bool) string {
	best, bestDist := "", 3
	for c := range allowed {
		if d := levenshtein(unknown, c); d < bestDist {
			best, bestDist = c, d
		}
	}
	return best
}

// checkClassContract C202：section 用到的每个类必须 ∈ 允许集合。
// 报错带"最接近的合法类名"——修复循环的自愈能力靠这个，不靠模型猜。
func checkClassContract(sec *goquery.Selection, layoutID string, allowed map[string]bool) error {
	var unknown []string
	for _, c := range extractUsedClasses(sec) {
		if !allowed[c] {
			unknown = append(unknown, c)
		}
	}
	if len(unknown) == 0 {
		return nil
	}
	var parts []string
	for _, c := range unknown {
		if s := suggestClass(c, allowed); s != "" {
			parts = append(parts, fmt.Sprintf(".%s（你是不是想写 .%s ？）", c, s))
		} else {
			parts = append(parts, "."+c)
		}
	}
	return fmt.Errorf("页面使用了版式 %q 契约之外的类名：%s。只准使用 base 原语类、本模板 style.css 声明的类、与 layouts.md 该版式登记的类", layoutID, strings.Join(parts, "、"))
}

// ---------- 生成批写入（write_pages）----------

// PageInput write_pages 的单页输入。
type PageInput struct {
	No     int    `json:"no"`     // 大纲页码（1 基），data-id 由此确定（s{no}）
	Layout string `json:"layout"` // 必须与 page_plan 一致
	HTML   string `json:"html"`   // 完整 <section>（含 data-layout 与 .notes 讲稿）
}

// PageWriteResult 单页写入结果（部分成功语义：好页落盘、坏页带原因返回）。
type PageWriteResult struct {
	No      int    `json:"no"`
	OK      bool   `json:"ok"`
	Error   string `json:"error,omitempty"`
	Warning string `json:"warning,omitempty"`
}

// WritePagesReport write_pages 的整体返回。
type WritePagesReport struct {
	Results   []PageWriteResult `json:"results"`
	Written   int               `json:"written"`
	Slides    int               `json:"slides"`    // 写入后的总页数
	Warning   string            `json:"warning,omitempty"`
}

// WritePagesV2 生成阶段的批写入（D14）。
//
// 闸门顺序（每页独立，坏页不影响好页）：
//  1. 页码必须落在大纲范围内（R105：实写页数 = 大纲页数）；
//  2. layout 已登记（C201）且与 page_plan 一致（要换版式改 plan，别偷改页面）；
//  3. 结构 + 消毒（复用 v1 的 sanitizeSlide 判定）；
//  4. 类名契约（C202）；
//  5. data-id 强制为 s{no}（LLM 写的 data-id 一律剥掉）。
func (s *Service) WritePagesV2(userID uint, deckID string, pages []PageInput) (*WritePagesReport, error) {
	if err := s.authorize(userID, deckID); err != nil {
		return nil, err
	}
	df, err := s.readDeckFile(deckID)
	if err != nil {
		return nil, err
	}
	if df.Stage != StageGenerating && df.Stage != StageIterating {
		return nil, fmt.Errorf("当前阶段 %s 不允许写页（生成管线未启动）", df.Stage)
	}
	if df.TemplateID == "" {
		return nil, fmt.Errorf("deck 还没选模板")
	}
	tpl, err := s.templateFor(df)
	if err != nil {
		return nil, err
	}
	outline, err := s.ReadOutline(userID, deckID)
	if err != nil {
		return nil, fmt.Errorf("读取大纲失败: %w", err)
	}

	planByNo := map[int]string{}
	for _, pa := range df.PagePlan {
		planByNo[pa.No] = pa.Layout
	}

	unlock := s.lockDeck(deckID)
	defer unlock()

	html, err := s.GetHTML(userID, deckID)
	if err != nil {
		return nil, err
	}
	prefix, segment, suffix, err := readSlideSegment(html)
	if err != nil {
		return nil, err
	}
	doc, err := parseSegmentDoc(segment)
	if err != nil {
		return nil, err
	}

	rep := &WritePagesReport{Results: make([]PageWriteResult, 0, len(pages))}
	var warnings []string
	writtenSomething := false

	for _, pg := range pages {
		res := PageWriteResult{No: pg.No}
		res.OK, res.Error, res.Warning = s.writeOnePage(doc, tpl, outline, planByNo, pg)
		if res.Warning != "" {
			warnings = append(warnings, res.Warning)
		}
		if res.OK {
			writtenSomething = true
			rep.Written++
		}
		rep.Results = append(rep.Results, res)
	}

	// 至少一页成功才落盘：全坏的批次不碰文件（模型拿到的错误信息已足够重写）
	if writtenSomething {
		// 段内按页码重排（批可能乱序到达；大纲顺序 = 页面顺序）
		if err := sortSegmentSections(doc); err != nil {
			return nil, err
		}
		segOut, err := segmentHTML(doc)
		if err != nil {
			return nil, err
		}
		if err := s.writeSlideSegment(deckID, prefix, segOut, suffix); err != nil {
			return nil, fmt.Errorf("写入 index.html 失败: %w", err)
		}
		rep.Slides = doc.Find("body").Children().Filter("section").Length()
	}
	rep.Warning = strings.Join(warnings, "\n")
	return rep, nil
}

// writeOnePage 单页的全闸门链。doc/sections 是当前段（含本批已写入的前序页）。
func (s *Service) writeOnePage(doc *goquery.Document, tpl *template.Template, outline *Outline, planByNo map[int]string, pg PageInput) (bool, string, string) {
	// 1. 页码范围（R105）
	if pg.No < 1 || pg.No > len(outline.Pages) {
		return false, fmt.Sprintf("页码 %d 超出大纲范围（1~%d）：实写页数必须等于大纲页数，改页数请改大纲", pg.No, len(outline.Pages)), ""
	}
	// 2. 版式登记（C201）+ 与计划一致
	if !tpl.HasLayout(pg.Layout) {
		return false, fmt.Sprintf("版式 %q 未在本模板登记；可用版式见 generate 提示词的版式索引，或 read_guidelines", pg.Layout), ""
	}
	if planned, ok := planByNo[pg.No]; ok && planned != pg.Layout {
		return false, fmt.Sprintf("第 %d 页的版式计划是 %q，你提交的是 %q。要换版式请重新走 plan_pages（全量重排），不要页面级偷换", pg.No, planned, pg.Layout), ""
	}
	// 3. 结构 + 消毒（v1 判定复用：唯一根 section、禁嵌套、危险标签拒、on*/危险URL 剥）
	sec, warning, err := parseSlideFragment(pg.HTML)
	if err != nil {
		return false, err.Error(), ""
	}
	// 4. data-layout 必须显式存在且与提交的 layout 字段一致
	if got := sec.AttrOr("data-layout", ""); got != pg.Layout {
		return false, fmt.Sprintf("<section> 的 data-layout=%q 与参数 layout=%q 不一致", got, pg.Layout), ""
	}
	// 5. 类名契约（C202）
	allowed := tpl.AllowedClasses(pg.Layout)
	if err := checkClassContract(sec, pg.Layout, allowed); err != nil {
		return false, err.Error(), ""
	}

	// data-id 权威在后端：剥掉 LLM 写的，按大纲页码编号
	sec.RemoveAttr("data-id")
	slideID := fmt.Sprintf("s%d", pg.No)
	sec.SetAttr("data-id", slideID)

	// 讲稿约定：骨架里应含 .notes div；漏了不阻塞（有些版式确实没有讲稿可写），
	// 但提示一句——演讲者模式空讲稿是质量缺陷。
	if sec.Find(".notes").Length() == 0 {
		warning = joinWarnings(warning, "该页没有 .notes 讲稿块：有讲稿的页请保留 layouts.md 骨架里的 <div class=\"notes\">")
	}

	canonical, err := goquery.OuterHtml(sec)
	if err != nil {
		return false, fmt.Sprintf("序列化失败: %v", err), ""
	}

	// 替换或追加：每页都从 doc 现查（前序页可能刚追加过，选择器不做跨页缓存）
	sections := doc.Find("body").Children().Filter("section")
	target := sections.FilterFunction(func(_ int, n *goquery.Selection) bool {
		return n.AttrOr("data-id", "") == slideID
	})
	if target.Length() > 0 {
		target.First().ReplaceWithHtml(canonical)
	} else {
		doc.Find("body").AppendHtml(canonical)
	}
	return true, "", warning
}

// sortSegmentSections 段内 section 按 data-id 页码升序（s2 必须排在 s10 前，字符串序会错）。
// goquery 没有原生的重排 API：收集 OuterHtml + 页码后在 Go 侧排序再整体重建。
func sortSegmentSections(doc *goquery.Document) error {
	type item struct {
		no  int
		html string
	}
	var items []item
	var firstErr error
	doc.Find("body").Children().Filter("section").Each(func(_ int, sec *goquery.Selection) {
		o, err := goquery.OuterHtml(sec)
		if err != nil && firstErr == nil {
			firstErr = err
			return
		}
		items = append(items, item{no: slideNoOf(sec), html: o})
	})
	if firstErr != nil {
		return firstErr
	}
	sort.Slice(items, func(i, j int) bool { return items[i].no < items[j].no })
	doc.Find("body").Children().Remove()
	for _, it := range items {
		doc.Find("body").AppendHtml(it.html)
	}
	return nil
}

// segmentHTML 序列化段（body 的 innerHTML，不带 body 标签）。
func segmentHTML(doc *goquery.Document) (string, error) {
	var sb strings.Builder
	doc.Find("body").Contents().Each(func(_ int, n *goquery.Selection) {
		h, err := goquery.OuterHtml(n)
		if err != nil {
			return
		}
		sb.WriteString("\n")
		sb.WriteString(h)
	})
	return sb.String(), nil
}

func slideNoOf(sec *goquery.Selection) int {
	var n int
	fmt.Sscanf(sec.AttrOr("data-id", ""), "s%d", &n)
	return n
}

// v2 的挂载标记（与模板包 slidesStartMarker/EndMarker 逐字一致；
// 两边各自声明避免 deck→template 之外的耦合，由契约测试守一致性）。
const (
	slidesStartMarker = "<!-- SLIDES:START -->"
	slidesEndMarker   = "<!-- SLIDES:END -->"
)

// ---------- 迭代期页级修改 ----------

// UpdateSlideV2 整页替换（fingerprint 乐观锁 + 全闸门）。
func (s *Service) UpdateSlideV2(userID uint, deckID, slideID, newHTML, fingerprint string) (string, error) {
	if err := s.authorize(userID, deckID); err != nil {
		return "", err
	}
	df, err := s.readDeckFile(deckID)
	if err != nil {
		return "", err
	}
	tpl, err := s.templateFor(df)
	if err != nil {
		return "", err
	}

	unlock := s.lockDeck(deckID)
	defer unlock()

	html, err := s.GetHTML(userID, deckID)
	if err != nil {
		return "", err
	}
	old, err := findSlideV2(html, slideID)
	if err != nil {
		return "", err
	}
	currentOuter, err := goquery.OuterHtml(old)
	if err != nil {
		return "", fmt.Errorf("提取 slide html 失败 err:%w", err)
	}
	if fingerprint != slideFingerprint(currentOuter) {
		return "", errors.New("内容已过期: 请重新调用 read_slide 获取最新内容后再提交修改")
	}

	sec, warning, err := parseSlideFragment(newHTML)
	if err != nil {
		return "", err
	}
	layoutID := sec.AttrOr("data-layout", "")
	if layoutID == "" {
		return "", fmt.Errorf("data-layout 缺失：本模板要求每页带已登记的 data-layout")
	}
	if !tpl.HasLayout(layoutID) {
		return "", fmt.Errorf("版式 %q 未登记", layoutID)
	}
	if err := checkClassContract(sec, layoutID, tpl.AllowedClasses(layoutID)); err != nil {
		return "", err
	}
	// data-id 不可变（页的身份），与被替换页一致
	got := sec.AttrOr("data-id", "")
	if got != slideID {
		return "", fmt.Errorf("new_html 的 data-id 是 %q，与要修改的 slide_id %q 不一致：请基于 read_slide 的内容修改，保留原 data-id", got, slideID)
	}

	canonical, err := goquery.OuterHtml(sec)
	if err != nil {
		return "", fmt.Errorf("序列化失败 err:%w", err)
	}

	prefix, segment, suffix, err := readSlideSegment(html)
	if err != nil {
		return "", err
	}
	doc, err := parseSegmentDoc(segment)
	if err != nil {
		return "", err
	}
	doc.Find("body").Children().Filter("section").FilterFunction(func(_ int, n *goquery.Selection) bool {
		return n.AttrOr("data-id", "") == slideID
	}).First().ReplaceWithHtml(canonical)
	segOut, err := segmentHTML(doc)
	if err != nil {
		return "", err
	}
	if err := s.writeSlideSegment(deckID, prefix, segOut, suffix); err != nil {
		return "", fmt.Errorf("写入 index.html 失败: %w", err)
	}
	return warning, nil
}

// InsertSlideV2 迭代期插入新页（layout 必填 + 契约校验；编号 = 最大值+1，不复用）。
func (s *Service) InsertSlideV2(userID uint, deckID, afterSlideID, newHTML string) (string, string, error) {
	if err := s.authorize(userID, deckID); err != nil {
		return "", "", err
	}
	df, err := s.readDeckFile(deckID)
	if err != nil {
		return "", "", err
	}
	tpl, err := s.templateFor(df)
	if err != nil {
		return "", "", err
	}

	unlock := s.lockDeck(deckID)
	defer unlock()

	html, err := s.GetHTML(userID, deckID)
	if err != nil {
		return "", "", err
	}
	sec, warning, err := parseSlideFragment(newHTML)
	if err != nil {
		return "", "", err
	}
	layoutID := sec.AttrOr("data-layout", "")
	if layoutID == "" {
		return "", "", fmt.Errorf("data-layout 缺失：本模板要求每页带已登记的 data-layout")
	}
	if !tpl.HasLayout(layoutID) {
		return "", "", fmt.Errorf("版式 %q 未登记", layoutID)
	}
	if err := checkClassContract(sec, layoutID, tpl.AllowedClasses(layoutID)); err != nil {
		return "", "", err
	}

	prefix, segment, suffix, err := readSlideSegment(html)
	if err != nil {
		return "", "", err
	}
	doc, err := parseSegmentDoc(segment)
	if err != nil {
		return "", "", err
	}
	sections := doc.Find("body").Children().Filter("section")
	if sections.Length() == 0 {
		return "", "", errors.New("deck 没有任何页面（异常状态），请走生成管线重建")
	}

	sec.RemoveAttr("data-id")
	newID := fmt.Sprintf("s%d", nextSlideNumberV2(sections)+1)
	sec.SetAttr("data-id", newID)
	canonical, err := goquery.OuterHtml(sec)
	if err != nil {
		return "", "", fmt.Errorf("序列化失败 err:%w", err)
	}

	if afterSlideID == "end" {
		sections.Last().AfterHtml(canonical)
	} else {
		target := sections.FilterFunction(func(_ int, n *goquery.Selection) bool {
			return n.AttrOr("data-id", "") == afterSlideID
		})
		if target.Length() == 0 {
			return "", "", fmt.Errorf("after_slide_id %q 不存在，请先用 list_slides 获取有效的 slide_id，或传 \"end\" 追加到末尾", afterSlideID)
		}
		target.First().AfterHtml(canonical)
	}

	segOut, err := segmentHTML(doc)
	if err != nil {
		return "", "", err
	}
	if err := s.writeSlideSegment(deckID, prefix, segOut, suffix); err != nil {
		return "", "", fmt.Errorf("写入 index.html 失败: %w", err)
	}
	return newID, warning, nil
}

func nextSlideNumberV2(sections *goquery.Selection) int {
	max := 0
	sections.Each(func(_ int, sec *goquery.Selection) {
		if m := slideIDPattern.FindStringSubmatch(sec.AttrOr("data-id", "")); m != nil {
			var n int
			fmt.Sscanf(m[1], "%d", &n)
			if n > max {
				max = n
			}
		}
	})
	return max
}

// DeleteSlideV2 迭代期删页（拒绝删空，与 v1 同一规则）。
func (s *Service) DeleteSlideV2(userID uint, deckID, slideID string) error {
	if err := s.authorize(userID, deckID); err != nil {
		return err
	}
	unlock := s.lockDeck(deckID)
	defer unlock()

	html, err := s.GetHTML(userID, deckID)
	if err != nil {
		return err
	}
	prefix, segment, suffix, err := readSlideSegment(html)
	if err != nil {
		return err
	}
	doc, err := parseSegmentDoc(segment)
	if err != nil {
		return err
	}
	sections := doc.Find("body").Children().Filter("section")
	if sections.Length() <= 1 {
		return errors.New("deck 至少要保留一页，不允许删空")
	}
	target := sections.FilterFunction(func(_ int, n *goquery.Selection) bool {
		return n.AttrOr("data-id", "") == slideID
	})
	if target.Length() == 0 {
		return fmt.Errorf("slide %q 不存在，请先用 list_slides 获取有效的 slide_id", slideID)
	}
	target.Remove()
	segOut, err := segmentHTML(doc)
	if err != nil {
		return err
	}
	return s.writeSlideSegment(deckID, prefix, segOut, suffix)
}
