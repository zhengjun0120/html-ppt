package usage

// 用量账本服务：trace 观测的每条 usage 事件经 Sink() 异步落一张 store.UsageEvent，
// Overview() 按 天聚合出用量页需要的全部数字。
//
// 为什么异步：Sink 在 trace.Usage 的调用路径上（agent 主循环/看图审查每轮一次），
// 落库慢了会拖住生成流；账本是统计用途，丢几条远好过卡住对话。所以带缓冲队列 +
// 后台批量写，队满丢弃并打日志。
//
// 为什么 Go 侧聚合不用 SQL GROUP BY：查询要同时兼容生产 MySQL 和测试的内存
// sqlite（日期函数方言不同），而且个人用量 30 天撑死几千行，全捞回来内存聚合
// 最稳也最好测。

import (
	"context"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"

	"html-ppt/backend/internal/store"
	"html-ppt/backend/internal/trace"
)

const (
	queueSize  = 256 // 约等于几十次完整 run 的用量事件数，正常远远到不了
	flushAt    = 64  // 攒够这么多条立刻落一批
	flushEvery = 500 * time.Millisecond
	dayLayout  = "2006-01-02"
	windowDays = 30
)

// Service 用量账本。Sink 由 trace 装配注入，Overview 由 handler 调。
type Service struct {
	db *gorm.DB
	ch chan store.UsageEvent

	mu      sync.Mutex // 串行化 Flush 与后台写之间的等待判断
	drained *sync.Cond
	pending int
}

// New 建服务并启动后台落库 goroutine。db 复用全局 *gorm.DB（AutoMigrate 已含 UsageEvent）。
func New(db *gorm.DB) *Service {
	s := &Service{
		db: db,
		ch: make(chan store.UsageEvent, queueSize),
	}
	s.drained = sync.NewCond(&s.mu)
	go s.loop()
	return s
}

// Sink 交给 trace.Config.UsageSink。非阻塞：队满丢弃并打日志（统计缺一条无所谓，
// 生成流不能等数据库）。
func (s *Service) Sink() func(trace.UsageRecord) {
	return func(rec trace.UsageRecord) {
		ev := store.UsageEvent{
			UserID:     rec.UserID,
			Day:        rec.At.Local().Format(dayLayout),
			Component:  rec.Component,
			Model:      rec.Model,
			SessionID:  rec.SessionID,
			RunID:      rec.RunID,
			Prompt:     rec.Prompt,
			Completion: rec.Completion,
			Total:      rec.Total,
			Cached:     rec.Cached,
			Reasoning:  rec.Reasoning,
			Calls:      rec.Calls,
		}
		s.mu.Lock()
		s.pending++
		s.mu.Unlock()
		select {
		case s.ch <- ev:
		default:
			log.Printf("[warn] usage: 账本队列已满，丢弃一条（component=%s user=%d）", ev.Component, ev.UserID)
			s.mu.Lock()
			s.pending--
			s.mu.Unlock()
		}
	}
}

// Flush 等队列里的存量全部落库（测试用：发完事件后调它再查库）。
func (s *Service) Flush() {
	s.mu.Lock()
	for s.pending > 0 {
		s.drained.Wait()
	}
	s.mu.Unlock()
}

func (s *Service) loop() {
	batch := make([]store.UsageEvent, 0, flushAt)
	tick := time.NewTicker(flushEvery)
	defer tick.Stop()
	for {
		select {
		case ev := <-s.ch:
			batch = append(batch, ev)
			if len(batch) < flushAt {
				continue
			}
		case <-tick.C:
		}
		if len(batch) == 0 {
			continue
		}
		s.write(batch)
		batch = batch[:0]
	}
}

func (s *Service) write(batch []store.UsageEvent) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.db.WithContext(ctx).CreateInBatches(batch, len(batch)).Error; err != nil {
		// 落库失败只打日志：账本是观测性质，绝不能反过来影响生成主流程
		log.Printf("[warn] usage: 账本落库失败，丢 %d 条 err: %v", len(batch), err)
	}
	s.mu.Lock()
	s.pending -= len(batch)
	if s.pending <= 0 {
		s.pending = 0
		s.drained.Broadcast()
	}
	s.mu.Unlock()
}

// TokenStat 一个时间窗的汇总。CachedRate = Cached/Prompt（0 ≤ rate ≤ 1），
// prompt 为 0 时记 0。
type TokenStat struct {
	Tokens     int64   `json:"tokens"`
	Calls      int64   `json:"calls"`
	CachedRate float64 `json:"cached_rate"`
}

// DayUsage 柱状图一天的三段：Cached（输入里命中缓存的部分）+ UncachedIn（输入里
// 未命中的部分）+ Output。三者相加等于当天 Total——缓存是输入的子集，从输入里
// 拆出来单列，柱子语义不重叠。
type DayUsage struct {
	Date       string `json:"date"`
	Cached     int64  `json:"cached"`
	UncachedIn int64  `json:"uncached_in"`
	Output     int64  `json:"output"`
	Total      int64  `json:"total"`
	Calls      int64  `json:"calls"`
}

// Overview 用量页一次拉全：今日/本月两块汇总 + 近 30 天逐日。
type Overview struct {
	Today TokenStat  `json:"today"`
	Month TokenStat  `json:"month"`
	Daily []DayUsage `json:"daily"`
}

// Overview 聚合某用户的用量。时间基准是服务器本地时区；月窗口从当月 1 号起，
// 日窗口固定近 30 天（含今天），空的日期补零。
func (s *Service) Overview(uid uint) (*Overview, error) {
	now := time.Now()
	today := now.Format(dayLayout)
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format(dayLayout)
	windowStart := now.AddDate(0, 0, -(windowDays - 1)).Format(dayLayout)
	if monthStart < windowStart {
		windowStart = monthStart
	}

	var rows []store.UsageEvent
	err := s.db.WithContext(context.Background()).
		Where("user_id = ? AND day >= ?", uid, windowStart).
		Order("day ASC").Find(&rows).Error
	if err != nil {
		return nil, err
	}

	ov := &Overview{Daily: make([]DayUsage, 0, windowDays)}
	byDay := make(map[string]*DayUsage, windowDays)
	for i := windowDays - 1; i >= 0; i-- {
		d := now.AddDate(0, 0, -i).Format(dayLayout)
		ov.Daily = append(ov.Daily, DayUsage{Date: d})
		byDay[d] = &ov.Daily[len(ov.Daily)-1]
	}

	var monthPrompt, monthCached int64
	var todayPrompt, todayCached int64
	for _, r := range rows {
		slot := byDay[r.Day]
		if slot != nil {
			slot.Cached += r.Cached
			slot.UncachedIn += r.Prompt - r.Cached
			slot.Output += r.Completion
			slot.Total += r.Total
			slot.Calls += int64(r.Calls)
		}
		if r.Day >= monthStart {
			st := &ov.Month
			st.Tokens += r.Total
			st.Calls += int64(r.Calls)
			monthPrompt += r.Prompt
			monthCached += r.Cached
		}
		if r.Day == today {
			st := &ov.Today
			st.Tokens += r.Total
			st.Calls += int64(r.Calls)
			todayPrompt += r.Prompt
			todayCached += r.Cached
		}
	}
	ov.Month.CachedRate = cachedRate(monthPrompt, monthCached)
	ov.Today.CachedRate = cachedRate(todayPrompt, todayCached)
	return ov, nil
}

func cachedRate(prompt, cached int64) float64 {
	if prompt <= 0 {
		return 0
	}
	return float64(cached) / float64(prompt)
}
