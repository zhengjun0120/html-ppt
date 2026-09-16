package agent

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"html-ppt/backend/internal/vision"
)

// 开关没开时**安静**（这一项本来就不该出现在工具结果里），
// 而"跑了但失败"必须出声——两者在模型眼里原本长得一模一样：都是没有 review 字段。
// 静默失败的代价是模型把"没量成"当成"量过了、没问题"。
func TestVisionOffIsSilentAndFailureWouldNotBe(t *testing.T) {
	svc := &AgentService{} // Vision 零值 = 关
	ctx := context.Background()

	if got := svc.measureDeck(ctx, 1, "deck-0001"); got != "" {
		t.Errorf("开关没开时 measureDeck 不该往工具结果里塞话，实际 %q", got)
	}
	_, err := svc.reviewPages(ctx, 1, "deck-0001", []int{1})
	if err == nil || !strings.Contains(err.Error(), "没开") {
		t.Errorf("开关没开时 review_slides 应当说清是功能没开（重试也没用），实际 %v", err)
	}
	if !errors.Is(errVisionOff, errVisionOff) {
		t.Error("errVisionOff 是调用方区分「没开」和「跑失败」的依据，必须是这个哨兵本身")
	}
}

// 0 基误传必须挡在门口。它不是"宽容一点也过得去"的输入：传 [0,1] 时 0 只会被当成
// "这份 deck 没有第 0 页"丢掉，实际拍的是第 1 页，而报告里写着"第 1 页"——
// 模型以为它看过了第 1、2 页，实际看的是第 1 页。这种错在结果里完全看不出来。
func TestNormalizePagesRejectsZeroBased(t *testing.T) {
	for _, bad := range [][]int{{0}, {0, 1}, {-1}, {2, 0}} {
		if _, err := normalizePages(bad); err == nil {
			t.Errorf("%v 应当报错（页号从 1 开始）", bad)
		}
	}
}

// 去重 + 升序：模型很容易把同一页写两遍（"第 2 页好像也有问题"），
// 重复的页号会被拍两张一样的图、占两份图片 token，报告里还会出现两行同一页。
func TestNormalizePagesDedupesAndSorts(t *testing.T) {
	got, err := normalizePages([]int{5, 2, 5, 3})
	if err != nil {
		t.Fatalf("合法页号不该报错: %v", err)
	}
	if !reflect.DeepEqual(got, []int{2, 3, 5}) {
		t.Errorf("应当去重升序得到 [2 3 5]，实际 %v", got)
	}
	if got, err := normalizePages(nil); err != nil || got != nil {
		t.Errorf("没传 pages 是合法的（＝只要数字不看图），应当返回 nil,nil，实际 %v,%v", got, err)
	}
}

// "没看的页里还有几页有问题"这条提示只列**页级**判定。
// 版式重复一次会牵出半份 deck 的页号，列出来的效果是让 agent 以为
// "十页里有八页都有问题"，反而把真正该看的那几页淹掉。
func TestUnselectedFindingsPagesSkipsDeckLevel(t *testing.T) {
	fs := []vision.Finding{
		{Level: vision.LevelPage, Pages: []int{2}, Text: "第 2 页：…"},
		{Level: vision.LevelPage, Pages: []int{7}, Text: "第 7 页：…"},
		{Level: vision.LevelLayout, Pages: []int{1, 3, 5, 7}, Text: "版式重复：…"},
	}

	got := unselectedFindingsPages(fs, []int{2})
	if !reflect.DeepEqual(got, []int{7}) {
		t.Errorf("点第 2 页时应当只提示第 7 页（第 2 页已经在看了，版式级那条不列），实际 %v", got)
	}
	if got := unselectedFindingsPages(fs, []int{2, 7}); got != nil {
		t.Errorf("两页都在看，不该再提示别的页，实际 %v", got)
	}
	// 升序：模型照着这个顺序挑下一批要看的页，乱序会看起来像"优先级"
	if got := unselectedFindingsPages([]vision.Finding{
		{Level: vision.LevelPage, Pages: []int{9}, Text: "第 9 页：…"},
		{Level: vision.LevelPage, Pages: []int{4}, Text: "第 4 页：…"},
	}, nil); !reflect.DeepEqual(got, []int{4, 9}) {
		t.Errorf("应当升序，实际 %v", got)
	}
}
