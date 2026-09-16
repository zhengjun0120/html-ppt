package vision

import (
	"strings"
	"testing"
)

// 造一页"什么都不越界"的页面，单个字段要越界时在测试里改。
func okSlide(i int, layout string) Slide {
	return Slide{
		Index: i, Layout: layout, Title: "标题",
		FitScale: 1, InnerScale: 1, EffScale: 1, Overflow: 1,
		Children: 3, Chars: 50, MinFontPx: 40,
	}
}

// 五个阈值各守一次边界：这些数字是提示词里对模型的承诺
// （"fit 低于 0.90 / 最小字号低于 31px / 溢出大于 1.02 就要重排"），
// 代码里改成一边开区间，模型照承诺挑页时就会点空——而报告只会说"无问题"。
func TestPageFindingThresholds(t *testing.T) {
	at := func(mut func(*Slide)) *Deck {
		s := okSlide(0, "s.hero")
		mut(&s)
		return &Deck{Slides: []Slide{s}}
	}

	over := []struct {
		name string
		mut  func(*Slide)
		want string
	}{
		{"fit 低于下限", func(s *Slide) { s.EffScale = 0.89 }, "适配兜底"},
		{"最小字号低于下限", func(s *Slide) { s.MinFontPx = 30 }, "最小字号"},
		{"字数超预算", func(s *Slide) { s.Chars = 201 }, "字预算"},
		{"子元素超上限", func(s *Slide) { s.Children = 6 }, "直接子元素"},
		{"溢出超过容差", func(s *Slide) { s.Overflow = 1.03 }, "溢出"},
	}
	for _, c := range over {
		got := HardFindings(at(c.mut))
		if len(got) != 1 || !strings.Contains(got[0], c.want) {
			t.Errorf("%s：应当判出一条含 %q 的判定，实际 %v", c.name, c.want, got)
		}
	}

	// 正好压线不报：越界判定用 < / >，不是 <= / >=
	for _, c := range []struct {
		name string
		mut  func(*Slide)
	}{
		{"fit 正好 0.90", func(s *Slide) { s.EffScale = 0.90 }},
		{"字号正好 31px", func(s *Slide) { s.MinFontPx = 31 }},
		{"字数正好 200", func(s *Slide) { s.Chars = 200 }},
		{"子元素正好 5 个", func(s *Slide) { s.Children = 5 }},
		{"溢出正好 1.02", func(s *Slide) { s.Overflow = 1.02 }},
	} {
		if got := HardFindings(at(c.mut)); len(got) != 0 {
			t.Errorf("%s：压线不该报，实际 %v", c.name, got)
		}
	}
}

// 版式重复是 **deck 级**规则，这一条是 ScopeFindings 存在的全部理由：
// 只点两页看的时候，拿那两页的子集去现算，重复次数永远数不够（子集里就是"用了 2 次，没超"），
// 于是模型会照着"没有版式重复"去画面上一页页确认——它确认的是一个假结论。
func TestLayoutRepeatIsDeckScopedNotSubsetScoped(t *testing.T) {
	d := &Deck{Slides: []Slide{
		okSlide(0, "s.hero"), okSlide(1, "s.hero"), okSlide(2, "s.hero"), okSlide(3, "s.grid"),
	}}
	all := Analyze(d)
	if !strings.Contains(strings.Join(Texts(all), "\n"), "版式重复") {
		t.Fatal("整份 deck 上应当判出版式重复（3 次 > 上限 2）")
	}

	sel, _ := d.Select([]int{1, 2})
	// 子集上现算："用了 3 次"这条根本数不出来（子集里只有 2 页），
	// 而且它还会派生出别的东西（连续两页同版式也变成"符合子集事实"的判定）——
	// 拿子集的判定去喂模型，模型确认的是另一份 deck 的结论。
	if f := strings.Join(HardFindings(sel), "\n"); strings.Contains(f, "用了 3 次") {
		t.Errorf("只选 2 页时不可能数出 3 次，实际 %v", f)
	}
	if scoped := strings.Join(ScopeFindings(all, []int{1, 2}), "\n"); !strings.Contains(scoped, "用了 3 次") {
		t.Errorf("子集审查必须沿用整份算出来的判定（否则模型照着假结论确认画面），实际 %v", scoped)
	}
	// 没牵涉到的页不带上：留着会让模型去确认它看不到的页
	if f := ScopeFindings(all, []int{4}); len(f) != 0 {
		t.Errorf("第 4 页与那条版式重复无关，不该带上：%v", f)
	}
}

// 页级判定只带选中页的那几条（带多了模型会去评论它没看到的页）。
func TestScopeFindingsKeepsOnlySelectedPages(t *testing.T) {
	small := okSlide(1, "s.grid")
	small.EffScale = 0.73
	tiny := okSlide(2, "s.cards")
	tiny.MinFontPx = 22
	d := &Deck{Slides: []Slide{okSlide(0, "s.hero"), small, tiny}}

	fs := Analyze(d)
	if len(fs) != 2 {
		t.Fatalf("应当有两条页级判定，实际 %d 条：%v", len(fs), Texts(fs))
	}
	got := ScopeFindings(fs, []int{3})
	if len(got) != 1 || !strings.Contains(got[0], "第 3 页") || !strings.Contains(got[0], "最小字号") {
		t.Errorf("只点第 3 页时应当只带第 3 页那条，实际 %v", got)
	}
	if got := ScopeFindings(fs, []int{1}); len(got) != 0 {
		t.Errorf("第 1 页没有判定，不该带任何判定：%v", got)
	}
}

// 取子集时页号口径：1 基、按 deck 顺序、保留原始 Index。
// 保留 Index 是关键——报告里"第 5 页"必须还是整份 deck 的第 5 页，
// 而不是子集里的第 2 张（错位的报告比没有报告更坏：改的时候会改错页）。
func TestSelectKeepsOriginalIndexAndReportsMissing(t *testing.T) {
	d := &Deck{CanvasW: 1244, CanvasH: 700, Slides: []Slide{
		okSlide(0, "a"), okSlide(1, "b"), okSlide(2, "c"),
	}}
	sel, missing := d.Select([]int{3, 1, 9})
	if len(sel.Slides) != 2 {
		t.Fatalf("点第 1、3、9 页应当取到 2 页，实际 %d 页", len(sel.Slides))
	}
	if sel.Slides[0].Index != 0 || sel.Slides[1].Index != 2 {
		t.Errorf("子集要按 deck 顺序排列且保留原始 Index，实际 %d、%d",
			sel.Slides[0].Index, sel.Slides[1].Index)
	}
	if len(missing) != 1 || missing[0] != 9 {
		t.Errorf("要了但没有的页号必须报出来（静默丢掉＝模型以为那页看过了），实际 %v", missing)
	}
	if sel.CanvasW != 1244 || sel.CanvasH != 700 {
		t.Errorf("画布尺寸要跟着子集走，实际 %d*%d", sel.CanvasW, sel.CanvasH)
	}
}
