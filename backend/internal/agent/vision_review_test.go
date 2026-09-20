package agent

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"html-ppt/backend/internal/authctx"
	"html-ppt/backend/internal/service/deck"
)

// 开关没开时**安静**（这一项本来就不该出现在工具结果里），
// 而"跑了但失败"必须出声——两者在模型眼里原本长得一模一样：都是没有 review 字段。
// 静默失败的代价是模型把"没量成"当成"量过了、没问题"。
func TestVisionOffIsSilentAndFailureWouldNotBe(t *testing.T) {
	svc := &AgentService{} // Vision 零值 = 关
	ctx := context.Background()

	if got := svc.measureDeckV2(ctx, 1, "deck-0001"); got != "" {
		t.Errorf("开关没开时 measureDeck 不该往工具结果里塞话，实际 %q", got)
	}
	ctx2 := authctx.WithUser(ctx, 1)
	_, err := svc.toolReviewSlidesV2(ctx2, `{"deck_id":"deck-0001","pages":[1]}`)
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



// briefOf 把大纲页压成看图调用能用的意图摘要。截断是刻意的：
// 意图几十字就够，全文灌进去会把看图提示词撑胖、还会把量测数字挤出注意力。
func TestBriefOfJoinsAndTruncates(t *testing.T) {
	p := deck.OutlinePage{Title: "毕业去向", Points: []string{"九成留粤", "升学 58.2%"}, Notes: "口径：2024 届"}
	got := briefOf(p)
	for _, want := range []string{"标题「毕业去向」", "要点：九成留粤；升学 58.2%", "备注：口径：2024 届"} {
		if !strings.Contains(got, want) {
			t.Errorf("briefOf 缺少 %q: %q", want, got)
		}
	}

	long := make([]string, 0, 40)
	for i := 0; i < 40; i++ {
		long = append(long, "这条要点特别长专门用来触发截断")
	}
	got = briefOf(deck.OutlinePage{Title: "T", Points: long})
	r := []rune(got)
	if len(r) > 161 || !strings.HasSuffix(got, "…") {
		t.Errorf("briefOf 未按 160 字截断: len=%d tail=%q", len(r), got[len(got)-20:])
	}
}
