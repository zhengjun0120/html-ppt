package deck

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//计算某页内容的指纹，read_slide 与 update_slide 都要使用这个
func slideFingerprint(outer string) string {
	sum := sha256.Sum256([]byte(outer))
	return hex.EncodeToString(sum[:])[:12]
}


// prepareReplacement 返回替换用的 HTML 和消毒警告（可能为空串）。
func prepareReplacement (newHTML,slideID string) (string,string,error){
	// 结构闸门与 insert 共用（parseSlideFragment）：唯一根 section、禁 doctype/嵌套、
	// 危险标签拒绝、on* 事件与危险协议 URL 剥除
	sec,warning,err := parseSlideFragment(newHTML)
	if err !=nil{
		return "","",err
	}

	// update 特有闸门：data-id 必须与要替换的页一致，防"把 A 页内容粘到 B 页"
	got := sec.AttrOr("data-id","")
	if got != slideID{
		return "","",fmt.Errorf("new_html 的data-id 是 %q, 与要修改的slide_id %q 不一致: 请基于 read_slide 返回的内容修改,保留原 data-id",got,slideID)
	}

	outer,err := goquery.OuterHtml(sec)
	if err !=nil {
		return "","",fmt.Errorf("序列化 new_html 失败 err:%w",err)
	}

	return outer,warning,nil

}

// UpdateSlide 返回消毒警告（可能为空串），供工具结果原样附给模型。
func (s *Service) UpdateSlide(userID uint, deckID, slideID, newHTML, fingerprint string) (string, error) {
	if err := s.authorize(userID,deckID);err !=nil{
		return "",err
	}
	unlock := s.lockDeck(deckID)
	defer unlock()
	// 鉴权+拿原HTML
	raw, err := s.readRaw(deckID)
	if err != nil {
		return "",err
	}

	doc,err := goquery.NewDocumentFromReader(strings.NewReader(raw))
	if err !=nil{
		return "",fmt.Errorf("解析 deck 失败 err:%w",err)
	}

	sec := doc.Find(".slides > section").FilterFunction(func(_ int, n *goquery.Selection) bool {
		return n.AttrOr("data-id","") == slideID
	})

	if sec.Length() == 0{
		return "",fmt.Errorf("slide %q 不存在, 请先用 list_slides 获取有效的 slide_id",slideID)
	}

	currentOuter,err := goquery.OuterHtml(sec.First())
	if err !=nil{
		return "",fmt.Errorf("提取 slide html 失败 err:%w",err)
	}

	if fingerprint != slideFingerprint(currentOuter){
		return "",errors.New("内容已过期: 请重新调用read_slide 获取最新内容后再提交修改")
	}

	replacement,warning,err := prepareReplacement(newHTML,slideID)
	if err !=nil{
		return "",err
	}

	// 回声校验：内联 style 里引用的变量必须真有元素消费（或已在主题块/本 deck 里定义），
	// 否则这处样式注定不生效——拒绝比"写进去了但没反应"诚实
	if err := s.checkInlineStyleVars(replacement, raw); err != nil {
		return "", err
	}

	// 替换HTML
	sec.First().ReplaceWithHtml(replacement)

	// 整份文档重新序列化
	out,err := goquery.OuterHtml(doc.Selection)
	if err !=nil{
		return "",fmt.Errorf("序列化 deck 失败 err:%w",err)
	}
	if err := s.atomicWriteDeck(deckID,out); err !=nil{
		return "",fmt.Errorf("写入 deck 失败 err:%w",err)
	}
	return warning,nil
}