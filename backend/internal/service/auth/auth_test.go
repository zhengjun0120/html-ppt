package auth

import (
	"testing"
	"time"
)

// dayStart 是"每用户每天一次"刷新额度的记账锚点：SQL 闸门用
// token_refreshed_at < dayStart(now) 判定"上次刷新不在今天"。
//
// 必须守住的是"本地日历日"这个语义——这正是不能写 time.Truncate(24h) 的原因：
// Truncate 按距零值时间的绝对时长切分，切点落在 UTC 零点，东八区会变成
// "每天早上 8 点到次日早上 8 点"。刷新额度切错 8 小时，用户会在深夜被拒。
func TestDayStart(t *testing.T) {
	loc := time.FixedZone("UTC+8", 8*3600)

	// 当天任意时刻的 dayStart 都对齐到本地零点
	noon := time.Date(2026, 9, 16, 12, 30, 0, 0, loc)
	got := dayStart(noon)
	want := time.Date(2026, 9, 16, 0, 0, 0, 0, loc)
	if !got.Equal(want) {
		t.Errorf("dayStart(12:30) = %v, want %v", got, want)
	}

	// 刷新记账的判定方向：上次刷新时刻 < 今天零点 = 今天没刷过。
	// 同一天里刷新过的（23:59:59 也算）必须 >= 今天零点，拿不到额度
	lastNight := time.Date(2026, 9, 15, 23, 59, 59, 0, loc)
	todayLate := time.Date(2026, 9, 16, 23, 59, 59, 0, loc)
	if !(lastNight.Before(dayStart(todayLate))) {
		t.Error("昨晚 23:59 的刷新应判为'今天没刷过'")
	}
	if !(dayStart(todayLate).Before(todayLate)) {
		t.Error("今晚 23:59 的刷新应判为'今天已刷过'（>= 今天零点）")
	}

	// 对照组（锁住选型理由）：Truncate 在东八区切在 08:00，不是零点
	if h := noon.Truncate(24 * time.Hour).In(loc).Hour(); h != 8 {
		t.Fatalf("前置假设失效：time.Truncate(24h) 在东八区切在了 %d 点（若 Go 改了此行为，请重估 dayStart 的实现必要性）", h)
	}
}
