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

// fillPct 的分段判定：45 以下是硬伤，55 以下是提示，之上正常。
// 阈值分两档是有意的——硬判定进 review 的"程序已判定"区（模型必须回应），
// 提示档只在摘要里挂 ⚠。两档写反了要么_model 被硬判定轰炸，要么空页漏网。
func TestFillPctThresholds(t *testing.T) {
	d := &Deck2{CanvasW: 1920, CanvasH: 1080, Slides: []Slide2{
		{Index: 0, FillPct: 92},
		{Index: 1, FillPct: 50},
		{Index: 2, FillPct: 8},
	}}
	hard := HardFindingsV2(d)
	if n := len(hard); n != 1 {
		t.Errorf("只有第 3 页（8%%）该进硬判定，实际 %d 条：%v", n, hard)
	}
	digest := DigestV2(d)
	if !strings.Contains(digest, "填充率 50") || !strings.Contains(digest, "填充率 8") {
		t.Errorf("第 2、3 页该在摘要里挂填充率提示：%s", digest)
	}
	if strings.Contains(digest, "第 1 页") && strings.Contains(digest, "92") && strings.Contains(digest, "⚠") {
		// 第 1 页 92% 不该有任何 ⚠（该行不含 ⚠ 即可，这里做兜底检查）
		for _, line := range strings.Split(digest, "\n") {
			if strings.HasPrefix(line, "第 1 页") && strings.Contains(line, "⚠") {
				t.Errorf("填充率 92%% 的页不该被标记：%s", line)
			}
		}
	}
}
