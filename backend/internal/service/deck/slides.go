package deck

import (
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