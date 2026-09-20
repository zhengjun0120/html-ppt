package deck

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

// C203 数量契约：repeats 登记的类必须恰好 N 个，偏多偏少都拒，报错给自愈路径。
func TestCheckRepeatCounts(t *testing.T) {
	mk := func(t *testing.T, html string) *goquery.Selection {
		t.Helper()
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			t.Fatal(err)
		}
		return doc.Find("section")
	}

	threeSteps := mk(t, `<section data-layout="how-it-works"><div class="stack">`+
		`<div class="step">a</div><div class="step">b</div><div class="step">c</div></div></section>`)

	if err := checkRepeatCounts(threeSteps, "how-it-works", map[string]int{"step": 3}); err != nil {
		t.Fatalf("3/3 应当通过: %v", err)
	}
	if err := checkRepeatCounts(threeSteps, "how-it-works", nil); err != nil {
		t.Fatalf("无契约不应拦截: %v", err)
	}

	elevenSteps := mk(t, `<section>`+strings.Repeat(`<div class="step">x</div>`, 11)+`</section>`)
	err := checkRepeatCounts(elevenSteps, "how-it-works", map[string]int{"step": 3})
	if err == nil {
		t.Fatal("11/3 应当被拒收")
	}
	for _, want := range []string{"step", "恰好 3 个", "11 个", "plan_pages"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("报错缺关键字 %q: %v", want, err)
		}
	}

	// 偏少同样拒：登记了数量就意味着骨架结构必须给满
	err = checkRepeatCounts(threeSteps, "feature-duo", map[string]int{"feature-card": 2})
	if err == nil || !strings.Contains(err.Error(), "0 个") {
		t.Fatalf("缺元素也应当报: %v", err)
	}
}
