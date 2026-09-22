package deck

import (
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestProbeFingerprintStability(t *testing.T) {
	html, err := os.ReadFile("D:/go_files/html-ppt/backend/data/decks/deck-0039/index.html")
	if err != nil {
		t.Skip("no deck file")
	}
	prefix, segment, suffix, err := readSlideSegment(string(html))
	if err != nil {
		t.Fatal(err)
	}
	// 读侧：parse → 取 s7 → OuterHtml → 指纹
	doc1, err := parseSegmentDoc(segment)
	if err != nil {
		t.Fatal(err)
	}
	s7a := doc1.Find("body").Children().Filter("section").FilterFunction(func(_ int, n *goquery.Selection) bool {
		return n.AttrOr("data-id", "") == "s7"
	}).First()
	outerA, _ := goquery.OuterHtml(s7a)
	fpA := slideFingerprint(outerA)

	// 写侧：替换 s3 → segmentHTML
	doc2, _ := parseSegmentDoc(segment)
	s3 := doc2.Find("body").Children().Filter("section").FilterFunction(func(_ int, n *goquery.Selection) bool {
		return n.AttrOr("data-id", "") == "s3"
	}).First()
	s3.ReplaceWithHtml(`<section class="slide" data-layout="cover" data-id="s3"><h1 class="h1">改动</h1></section>`)
	segOut, err := segmentHTML(doc2)
	if err != nil {
		t.Fatal(err)
	}
	_ = prefix
	_ = suffix

	// 读侧再来一轮：parse 新段 → s7 指纹
	doc3, _ := parseSegmentDoc(segOut)
	s7b := doc3.Find("body").Children().Filter("section").FilterFunction(func(_ int, n *goquery.Selection) bool {
		return n.AttrOr("data-id", "") == "s7"
	}).First()
	outerB, _ := goquery.OuterHtml(s7b)
	fpB := slideFingerprint(outerB)

	if fpA != fpB {
		t.Logf("指纹漂移! A=%s B=%s", fpA, fpB)
		// 找出第一个差异点
		for i := 0; i < len(outerA) && i < len(outerB); i++ {
			if outerA[i] != outerB[i] {
				t.Logf("首个差异 @%d:\nA: ...%s\nB: ...%s", i, outerA[max(0,i-60):i+60], outerB[max(0,i-60):i+60])
				break
			}
		}
		if len(outerA) != len(outerB) {
			t.Logf("长度不同: %d vs %d", len(outerA), len(outerB))
		}
	} else {
		t.Logf("指纹稳定 %s", fpA)
	}
	if !strings.Contains(segOut, "改动") {
		t.Fatal("替换未生效")
	}
}
