package usage

import (
	"testing"
	"time"

	"html-ppt/backend/internal/store"
	"html-ppt/backend/internal/trace"
)

// OpenMemory 是进程级共享库（cache=shared），用户 id 必须各测试独占一段，防撞键。
var usageTestUser uint = 91_000

func newTestService(t *testing.T) (*Service, *time.Location) {
	t.Helper()
	st, err := store.OpenMemory(t.Context())
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	return New(st.DB), time.Local
}

// dayOf 与服务同款取日串，测试里据此算期望值，月初月末跑都不翻车。
func dayOf(offsetDays int) string {
	return time.Now().AddDate(0, 0, offsetDays).Format(dayLayout)
}

func TestOverviewAggregates(t *testing.T) {
	s, _ := newTestService(t)
	uid := usageTestUser + 1
	other := usageTestUser + 2

	mk := func(day string, uid uint, comp string, prompt, completion, cached, calls int64) store.UsageEvent {
		return store.UsageEvent{
			UserID: uid, Day: day, Component: comp, Model: "deepseek-flash",
			Prompt: prompt, Completion: completion, Total: prompt + completion,
			Cached: cached, Calls: int(calls),
		}
	}
	rows := []store.UsageEvent{
		// 今天：两条，含一条 vision 子调用
		mk(dayOf(0), uid, "main", 1000, 200, 800, 1),
		mk(dayOf(0), uid, "vision", 500, 300, 0, 1),
		// 5 天前：在 30 天窗口内；是否在"本月"取决于今天几号，断言按窗口算
		mk(dayOf(-5), uid, "main", 400, 100, 0, 2),
		// 40 天前：窗口外，哪都不该出现
		mk(dayOf(-40), uid, "main", 99999, 99999, 0, 9),
		// 别人的今天：隔离
		mk(dayOf(0), other, "main", 7777, 777, 0, 3),
	}
	for i := range rows {
		if err := s.db.Create(&rows[i]).Error; err != nil {
			t.Fatalf("造行失败: %v", err)
		}
	}

	ov, err := s.Overview(uid, "")
	if err != nil {
		t.Fatalf("Overview: %v", err)
	}

	// 今日：main 1000+200 + vision 500+300
	if ov.Today.Tokens != 2000 || ov.Today.Calls != 2 {
		t.Errorf("Today = %+v，想要 tokens 2000 calls 2", ov.Today)
	}
	// rate 是全窗口 Σcached/Σprompt：今日输入 1000+500，命中 800
	if r := ov.Today.CachedRate; r < 0.5332 || r > 0.5334 {
		t.Errorf("Today.CachedRate = %v，想要 0.5333（缓存 800 / 输入 1500）", r)
	}

	// 本月：今天两条必在本月；5 天前那条约不约在月内按日历算
	monthStart := time.Now().Format("2006-01") + "-01"
	wantMonth := int64(2000)
	if dayOf(-5) >= monthStart {
		wantMonth += 500
	}
	if ov.Month.Tokens != wantMonth || ov.Month.Calls != int64(2+b2i(dayOf(-5) >= monthStart)*2) {
		t.Errorf("Month = %+v，想要 tokens %d", ov.Month, wantMonth)
	}

	// 近 30 天：30 槽、首尾日期对、零填充、数字归位
	if len(ov.Daily) != windowDays {
		t.Fatalf("Daily 长度 = %d，想要 %d", len(ov.Daily), windowDays)
	}
	first, last := ov.Daily[0], ov.Daily[windowDays-1]
	if first.Date != dayOf(-(windowDays-1)) || last.Date != dayOf(0) {
		t.Errorf("Daily 日期范围 %s..%s 不对", first.Date, last.Date)
	}
	mid := ov.Daily[10]
	if mid.Date != dayOf(-(windowDays-1-10)) || mid.Total != 0 || mid.Calls != 0 {
		t.Errorf("中间空日没补零：%+v（date=%s）", mid, dayOf(-(windowDays - 1 - 10)))
	}
	d5 := ov.Daily[windowDays-1-5]
	if d5.Total != 500 || d5.Calls != 2 || d5.UncachedIn != 400 || d5.Output != 100 {
		t.Errorf("5 天前的柱子 = %+v，想要 total 500 / uncached 400 / output 100", d5)
	}
	if d5.Cached != 0 {
		t.Errorf("5 天前没缓存命中，Cached = %d", d5.Cached)
	}
	// 今日柱：uncached_in = 1000+500-800
	if last.UncachedIn != 700 || last.Cached != 800 || last.Output != 500 || last.Total != 2000 {
		t.Errorf("今日柱 = %+v，想要 cached 800 / uncached 700 / output 500 / total 2000", last)
	}
}

func TestSinkLandsRows(t *testing.T) {
	s, _ := newTestService(t)
	uid := usageTestUser + 3

	sink := s.Sink()
	sink(trace.UsageRecord{
		UserID: uid, SessionID: 7, RunID: "1790000000000-ab12",
		Component: "main", Model: "deepseek-flash", At: time.Now(),
		UsagePart: trace.UsagePart{Prompt: 120, Completion: 30, Total: 150, Cached: 60, Calls: 1},
	})
	sink(trace.UsageRecord{
		UserID: uid, SessionID: 7, RunID: "1790000000000-ab12",
		Component: "vision", Model: "deepseek-flash", At: time.Now(),
		UsagePart: trace.UsagePart{Prompt: 5000, Completion: 400, Total: 5400, Calls: 1},
	})
	s.Flush()

	var rows []store.UsageEvent
	if err := s.db.Where("user_id = ?", uid).Order("component ASC").Find(&rows).Error; err != nil {
		t.Fatalf("查账本: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("账本应有 2 行，实际 %d", len(rows))
	}
	main := rows[0]
	if main.Component != "main" || main.Day != dayOf(0) || main.SessionID != 7 ||
		main.RunID != "1790000000000-ab12" || main.Model != "deepseek-flash" ||
		main.Prompt != 120 || main.Total != 150 || main.Cached != 60 || main.Calls != 1 {
		t.Errorf("main 行字段不对: %+v", main)
	}
	if rows[1].Component != "vision" || rows[1].Reasoning != 0 {
		t.Errorf("vision 行不对: %+v", rows[1])
	}

	// 落库后 Overview 能看到（走通 sink→DB→聚合 整条链）
	ov, err := s.Overview(uid, "")
	if err != nil {
		t.Fatalf("Overview: %v", err)
	}
	if ov.Today.Tokens != 150+5400 || ov.Today.Calls != 2 {
		t.Errorf("链路今日汇总 = %+v，想要 5550/2", ov.Today)
	}
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

func TestOverviewFiltersByModel(t *testing.T) {
	s, ctx := newTestService(t)
	uid := uint(95_010)
	day := dayOf(0)
	mk := func(model string, total int64) store.UsageEvent {
		return store.UsageEvent{UserID: uid, Day: day, Component: "main", Model: model, Total: total, Calls: 1}
	}
	for _, r := range []store.UsageEvent{mk("deepseek-flash", 100), mk("my-gpt", 200), mk("deepseek-flash", 50)} {
		if err := s.db.Create(&r).Error; err != nil {
			t.Fatalf("造行: %v", err)
		}
	}

	all, err := s.Overview(uid, "")
	if err != nil {
		t.Fatalf("Overview 全部: %v", err)
	}
	if all.Today.Tokens != 350 {
		t.Errorf("全部模型 tokens = %d，想要 350", all.Today.Tokens)
	}
	// 选项取全部账本（不受过滤影响），已排序去重由 DISTINCT 保证
	if len(all.Models) != 2 {
		t.Fatalf("模型选项 = %v，想要 2 个", all.Models)
	}

	one, err := s.Overview(uid, "my-gpt")
	if err != nil {
		t.Fatalf("Overview 筛选: %v", err)
	}
	if one.Today.Tokens != 200 || one.Today.Calls != 1 {
		t.Errorf("my-gpt 筛选 = %+v，想要 200/1", one.Today)
	}
	// 筛选状态下选项仍然齐全
	if len(one.Models) != 2 {
		t.Errorf("筛选状态下模型选项 = %v，仍应有 2 个", one.Models)
	}
	_ = ctx
}
