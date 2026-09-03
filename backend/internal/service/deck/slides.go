package deck

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type SlideMeta struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func (s *Service) ListSlides(id string) ([]SlideMeta, error) {
	//读取html
	html, err := s.GetHTML(id)
	if err != nil {
		return nil, err
	}

	//解析html
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err !=nil{
		return nil,fmt.Errorf("解析 deck 失败 err:%w",err)
	}

	sections := doc.Find(".slides > section")

	out := make([]SlideMeta,0,0)
	sections.Each(func(_ int, sec *goquery.Selection) {
		out = append(out, SlideMeta{
			ID: sec.AttrOr("data-id",""),
			Title: slideTitle(sec),
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

func (s *Service) ReadSlide(deckID,slideID string)(string,string,error){
	raw,err := s.GetHTML(deckID)
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

	sum := sha256.Sum256([]byte(outer))
	return outer,hex.EncodeToString(sum[:])[:12],nil
}
