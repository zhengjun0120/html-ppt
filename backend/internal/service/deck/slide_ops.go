package deck

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// parseSlideFragment 校验 LLM 提交的单页片段（insert 与 update 共用的闸门）：
// 唯一根 <section>、禁 script/style/doctype、禁嵌套。返回该页的 selection，
// 供调用方继续加工（update 核对 data-id，insert 剥掉重编号）。
func parseSlideFragment(newHTML string) (*goquery.Selection, error) {
	raw := strings.TrimSpace(newHTML)
	if raw == "" {
		return nil, errors.New("new_html 不能为空")
	}

	lower := strings.ToLower(raw)
	if strings.Contains(lower, "<!doctype") || strings.Contains(lower, "<html") {
		return nil, errors.New("new_html 只能是一个 <section> 元素，不要输出完整的 HTML 文档")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("解析 new_html 失败 err:%w", err)
	}

	if doc.Find("script").Length() > 0 {
		return nil, errors.New("幻灯片不支持 <script>，请移除后重新提交")
	}
	if doc.Find("style").Length() > 0 {
		return nil, errors.New("幻灯片不支持 <style>，请使用组件库 class 或内联样式，请移除后重新提交")
	}

	// Children() 只算元素节点，纯注释不算
	roots := doc.Find("body").Children().Filter("section")
	if roots.Length() == 0 {
		return nil, errors.New("new_html 必须以 <section> 为根元素")
	}
	// 多个 section 或夹带其他元素都拒绝——只取根节点会静默丢数据
	if roots.Length() > 1 || doc.Find("body").Children().Length() != roots.Length() {
		return nil, errors.New("new_html 只能包含一个 <section> 元素，多页修改请分多次调用")
	}
	if roots.First().Find("section").Length() > 0 {
		return nil, errors.New("不支持嵌套 <section>（垂直子页），页内不要嵌套 section")
	}

	return roots.First(), nil
}

var slideIDPattern = regexp.MustCompile(`^s(\d+)$`)

// nextSlideNumber 扫描现有页的 data-id（s1、s2…），返回最大编号。
// 删除不回填编号，所以"最大值+1"永远不会撞上历史 id。
func nextSlideNumber(doc *goquery.Document) int {
	max := 0
	doc.Find(".slides > section").Each(func(_ int, sec *goquery.Selection) {
		if m := slideIDPattern.FindStringSubmatch(sec.AttrOr("data-id", "")); m != nil {
			if n, err := strconv.Atoi(m[1]); err == nil && n > max {
				max = n
			}
		}
	})
	return max
}

// InsertSlide 在 afterSlideID 页之后插入新页（"end" = 追加到末尾），
// 返回后端分配的新 slide_id。LLM 写的 data-id 一律剥掉重编号。
func (s *Service) InsertSlide(userID uint, deckID, afterSlideID, newHTML string) (string, error) {
	if err := s.authorize(userID, deckID); err != nil {
		return "", err
	}
	raw, err := s.readRaw(deckID)
	if err != nil {
		return "", err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(raw))
	if err != nil {
		return "", fmt.Errorf("解析 deck 失败 err:%w", err)
	}

	sections := doc.Find(".slides > section")
	if sections.Length() == 0 {
		return "", errors.New("deck 没有任何页面，属于异常状态，请用 write_deck 重建")
	}

	sec, err := parseSlideFragment(newHTML)
	if err != nil {
		return "", err
	}
	// 编号是后端的权威：LLM 给的 data-id 不可信，剥掉重发
	sec.RemoveAttr("data-id")
	newID := fmt.Sprintf("s%d", nextSlideNumber(doc)+1)
	sec.SetAttr("data-id", newID)
	canonical, err := goquery.OuterHtml(sec)
	if err != nil {
		return "", fmt.Errorf("序列化新页失败 err:%w", err)
	}

	if afterSlideID == "end" {
		sections.Last().AfterHtml(canonical)
	} else {
		target := sections.FilterFunction(func(_ int, n *goquery.Selection) bool {
			return n.AttrOr("data-id", "") == afterSlideID
		})
		if target.Length() == 0 {
			return "", fmt.Errorf("after_slide_id %q 不存在，请先用 list_slides 获取有效的 slide_id，或传 \"end\" 追加到末尾", afterSlideID)
		}
		target.First().AfterHtml(canonical)
	}

	out, err := goquery.OuterHtml(doc.Selection)
	if err != nil {
		return "", fmt.Errorf("序列化 deck 失败 err:%w", err)
	}
	if err := s.atomicWriteDeck(deckID, out); err != nil {
		return "", fmt.Errorf("写入 deck 失败 err:%w", err)
	}
	return newID, nil
}

// DeleteSlide 删除一页。拒绝删空（至少保留一页，全删的合法路径是 write_deck 重做）。
// 被删页的编号不回收：后续新页继续从"最大值+1"分配，旧 id 永远不会指向新页。
func (s *Service) DeleteSlide(userID uint, deckID, slideID string) error {
	if err := s.authorize(userID, deckID); err != nil {
		return err
	}
	raw, err := s.readRaw(deckID)
	if err != nil {
		return err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(raw))
	if err != nil {
		return fmt.Errorf("解析 deck 失败 err:%w", err)
	}

	sections := doc.Find(".slides > section")
	if sections.Length() <= 1 {
		return errors.New("deck 至少要保留一页，不允许删空；若要全部重来请用 write_deck")
	}

	target := sections.FilterFunction(func(_ int, n *goquery.Selection) bool {
		return n.AttrOr("data-id", "") == slideID
	})
	if target.Length() == 0 {
		return fmt.Errorf("slide %q 不存在，请先用 list_slides 获取有效的 slide_id", slideID)
	}

	target.Remove()

	out, err := goquery.OuterHtml(doc.Selection)
	if err != nil {
		return fmt.Errorf("序列化 deck 失败 err:%w", err)
	}
	if err := s.atomicWriteDeck(deckID, out); err != nil {
		return fmt.Errorf("写入 deck 失败 err:%w", err)
	}
	return nil
}
