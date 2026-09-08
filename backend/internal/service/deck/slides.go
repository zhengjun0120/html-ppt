package deck

import (
	"fmt"
	"strings"
	"github.com/PuerkitoBio/goquery"
)

type SlideMeta struct {
	Position int    `json:"position"` // 1 起始的当前页序：随插入/删除变化，仅供人/模型定位
	ID       string `json:"id"`       // 每页的唯一标识（身份证）：永不变化，永不复用
	Title    string `json:"title"`
}

func (s *Service) ListSlides(userID uint, id string) ([]SlideMeta, error) {


	//鉴权+读取html（authorize 已过，readRaw 不再做归属校验）
	html, err := s.readOwned(userID,id)
	if err != nil {
		return nil, err
	}

	//解析html
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err !=nil{
		return nil,fmt.Errorf("解析 deck 失败 err:%w",err)
	}

	sections := doc.Find(".slides > section")

	out := make([]SlideMeta,0)
	sections.Each(func(i int, sec *goquery.Selection) {
		out = append(out, SlideMeta{
			Position: i + 1,
			ID:       sec.AttrOr("data-id",""),
			Title:    slideTitle(sec),
		})
	})

	return out,nil
}

func slideTitle(sec *goquery.Selection) string{
	for _,tag := range []string{"h1","h2","h3"}{
		if t := strings.TrimSpace(sec.Find(tag).First().Text()); t!=""{
			return t
		}
	}

	text := strings.Join(strings.Fields(sec.Text())," ")
	if r := []rune(text);len(r)>30{
		text = string(r[:30])+"..."
	}

	return text
}

// ReadSlide 返回指定页的完整 HTML 和内容指纹（公开读入口，先过归属校验）。
func (s *Service) ReadSlide(userID uint, deckID, slideID string) (string, string, error) {
	// 鉴权+读
	raw, err := s.readOwned(userID,deckID)
	if err !=nil{
		return "","",err
	}

	doc,err := goquery.NewDocumentFromReader(strings.NewReader(raw))
	if err !=nil{
		return "","",fmt.Errorf("html解析失败 err:%w",err)
	}

	sec := doc.Find(".slides > section").FilterFunction(func(_ int, n *goquery.Selection) bool {
		return n.AttrOr("data-id","") == slideID
	})

	if sec.Length() == 0{
		return "","",fmt.Errorf("slide %q 不存在，请先用 list_slides 获取有效slide_id",slideID)
	}

	outer,err := goquery.OuterHtml(sec.First())
	if err != nil{
		return "","",fmt.Errorf("提取 slide html失败 err:%w",err)
	}

	
	return outer,slideFingerprint(outer),nil
}
