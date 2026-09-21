package vision

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// v2 版的量测契约守卫：v2MeasureAllJS 返回的 key 必须与 Slide2 的 json tag
// 逐字一致（与 v1 的 TestMeasureKeysMatchSlideStruct 同一套纪律）。
// 这类错不报错不 panic：字段静默为零，报告照着 0 说"没问题"。
// fillPct 就差点翻车——加字段那天两个方向都可能漏。
func TestMeasureKeysMatchSlide2Struct(t *testing.T) {
	// 取 v2MeasureAllJS 里最后一段 out.push({ ... })：它才是每页量测的形状
	idx := strings.LastIndex(v2MeasureAllJS, "out.push({")
	if idx < 0 {
		t.Fatal("v2MeasureAllJS 里找不到 out.push({")
	}
	body := v2MeasureAllJS[idx:]
	if end := strings.Index(body, "});"); end > 0 {
		body = body[:end]
	}

	got := map[string]bool{}
	for _, line := range strings.Split(body, "\n") {
		if m := regexp.MustCompile(`^\s*([A-Za-z_]\w*)\s*:`).FindStringSubmatch(line); m != nil {
			got[m[1]] = true
		}
	}
	if len(got) < 8 {
		t.Fatalf("只从 v2MeasureAllJS 里解析出 %d 个 key，解析逻辑不对：%v", len(got), got)
	}

	want := map[string]bool{}
	typ := reflect.TypeOf(Slide2{})
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		want[tag] = true
		if !got[tag] {
			t.Errorf("Slide2.%s 的 tag 是 %q，但 v2MeasureAllJS 没返回这个 key——这个字段永远是零值",
				typ.Field(i).Name, tag)
		}
	}
	for k := range got {
		if !want[k] {
			t.Errorf("v2MeasureAllJS 返回了 %q，但 Slide2 里没有对应的 json tag——这个数字被丢掉了", k)
		}
	}
}

// fillPct 的判定：45 以下才告警，与上传门禁（usertpl 内容版式 ≥45%）同一把尺；
// hero/quote 版式豁免——它们的留白是设计。曾经摘要档用 55、又不看版式，46%~54%
// 的正常页与天生稀疏的宣言页全被挂 ⚠，agent 为消警把卡片注水到 300+ 字
//（deck-0061 实测 5 轮修复全因此）。
func TestFillPctThresholds(t *testing.T) {
	d := &Deck2{CanvasW: 1920, CanvasH: 1080,
		Patterns: map[string]string{"cover": "hero", "quote-big": "quote"},
		Slides: []Slide2{
			{Index: 0, FillPct: 92},
			{Index: 1, FillPct: 50},
			{Index: 2, FillPct: 8},
			{Index: 3, Layout: "cover", FillPct: 10},
			{Index: 4, Layout: "quote-big", FillPct: 30, BottomGap: 600},
		}}
	hard := HardFindingsV2(d)
	if n := len(hard); n != 1 {
		t.Errorf("只有第 3 页（8%%）该进硬判定（hero/quote 豁免），实际 %d 条：%v", n, hard)
	}
	digest := DigestV2(d)
	flagged := func(prefix string) bool {
		for _, line := range strings.Split(digest, "\n") {
			if strings.HasPrefix(line, prefix) && strings.Contains(line, "⚠") {
				return true
			}
		}
		return false
	}
	if !flagged("第 3 页") {
		t.Errorf("第 3 页（8%%）该在摘要里挂填充率提示：%s", digest)
	}
	for _, no := range []string{"第 1 页", "第 2 页", "第 4 页", "第 5 页"} {
		if flagged(no) {
			t.Errorf("%s 不该被标记（50%% 已过 45 线；hero/quote 豁免）：%s", no, digest)
		}
	}
}
