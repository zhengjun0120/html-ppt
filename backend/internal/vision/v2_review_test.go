package vision

import (
	"strings"
	"testing"
)

// 看图调用是隔离的，主对话的上下文只能靠 ReviewContext 显式注入——
// 这两个测试守的就是"注入真的到了提示词里"。注入静默丢失的后果：
// 看图退化成泛泛审查，focus 参数形同虚设，而且没有任何报错。
func TestReviewHeadInjectsFocusAndFindings(t *testing.T) {
	d := &Deck2{CanvasW: 1920, CanvasH: 1080, Slides: []Slide2{
		{Index: 0, Layout: "cards", Title: "封面", Chars: 80, MinFontPx: 18, FillPct: 60},
		{Index: 1, Layout: "cards", Title: "去向", Chars: 80, MinFontPx: 18, FillPct: 30, OverflowY: true},
	}}
	rc := &ReviewContext{
		Focus:      "第 2 页数据卡与来源小字是否对齐",
		PageBriefs: map[int]string{2: "标题「毕业去向」；要点：九成留粤"},
	}
	head := reviewHead(d, []int{2}, rc)

	for _, want := range []string{
		"特别想验证",
		"第 2 页数据卡与来源小字是否对齐",
		"优先回答这些疑问",
		"程序已判定的硬问题",
		"第 2 页内容溢出画布", // HardFindingsV2 的硬判定必须仍在
		"逐页量测",
	} {
		if !strings.Contains(head, want) {
			t.Errorf("reviewHead 缺少 %q\nhead=%s", want, head)
		}
	}

	// rc 为 nil（裸看图）时不得出现 focus 段
	if got := reviewHead(d, []int{2}, nil); strings.Contains(got, "特别想验证") {
		t.Errorf("rc=nil 不应出现 focus 段:\n%s", got)
	}
}

func TestPageLabelCarriesOutlineBrief(t *testing.T) {
	rc := &ReviewContext{PageBriefs: map[int]string{2: "标题「毕业去向」"}}

	if got := pageLabel(2, rc); !strings.Contains(got, "大纲意图：标题「毕业去向」") {
		t.Errorf("pageLabel 丢了大纲意图: %q", got)
	}
	// 没有摘要的页：只保留页号标头，不出现空意图段
	if got := pageLabel(3, rc); strings.Contains(got, "大纲意图") {
		t.Errorf("没有摘要的页不应出现意图段: %q", got)
	}
	// rc=nil 不 panic、不出现意图段
	if got := pageLabel(2, nil); strings.Contains(got, "大纲意图") {
		t.Errorf("rc=nil 不应出现意图段: %q", got)
	}
}
