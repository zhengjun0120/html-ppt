package deck

// 版本控制（后悔药）的回归测试。白盒：直接打内部方法，绕开 authorize，
// 所以 st=nil（数据库降级态）也能跑——鉴权本身由 handler 层的 EnsureOwner 覆盖。
//
// 为什么这批测试值得存在：自由度放开的全部安全感都建立在"随时能退回去"上，
// 恢复要是坏的，后悔药就是假的，而这个 bug 平时不会被发现（没人天天回滚）。

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func newHistoryTestService(t *testing.T) *Service {
	t.Helper()
	return &Service{decksDir: t.TempDir()}
}

// writeDeckHTML 模拟一次落盘（等价于 atomicWriteDeck 的效果），返回错误供并发测试用。
func writeDeckHTML(s *Service, id, html string) error {
	p := filepath.Join(s.decksDir, id, "deck.html")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return atomicWriteFile(p, []byte(html))
}

// deckHTML 造一份含 n 页的最小 deck。
func deckHTML(n int) string {
	var sb strings.Builder
	sb.WriteString(`<div class="reveal"><div class="slides">`)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, `<section data-id="s%d"><h2>页 %d</h2></section>`, i+1, i+1)
	}
	sb.WriteString(`</div></div>`)
	return sb.String()
}

// deckWithTitles 造带指定标题的 deck（用于 diff 分类测试，页面取 s1/s2/s3）。
func deckWithTitles(titles map[string]string) string {
	var sb strings.Builder
	sb.WriteString(`<div class="reveal"><div class="slides">`)
	for _, id := range []string{"s1", "s2", "s3"} {
		title, ok := titles[id]
		if !ok {
			continue
		}
		fmt.Fprintf(&sb, `<section data-id="%s"><h2>%s</h2></section>`, id, title)
	}
	sb.WriteString(`</div></div>`)
	return sb.String()
}

func TestRecordVersionListsNewestFirst(t *testing.T) {
	s := newHistoryTestService(t)
	if err := writeDeckHTML(s, "deck-9001", deckHTML(3)); err != nil {
		t.Fatal(err)
	}
	s.recordVersion("deck-9001", OpRun, "第一轮", deckHTML(3))
	if err := writeDeckHTML(s, "deck-9001", deckHTML(4)); err != nil {
		t.Fatal(err)
	}
	s.recordVersion("deck-9001", OpRun, "第二轮", deckHTML(4))

	idx := s.readHistoryIndex("deck-9001")
	if idx.NextSeq != 3 || len(idx.Versions) != 2 {
		t.Fatalf("序号或条数不对: next=%d len=%d", idx.NextSeq, len(idx.Versions))
	}
	// 新→旧：第一条是 v2 且页数为 4，Detail 是当轮汇总
	if idx.Versions[0].Version != "v000002" || idx.Versions[0].Slides != 4 {
		t.Fatalf("排序或页数不对: %+v", idx.Versions[0])
	}
	if idx.Versions[0].Detail != "第二轮" || idx.Versions[1].Operation != OpRun {
		t.Fatalf("元信息不对: %+v", idx.Versions)
	}
	// 快照文件必须真实存在——恢复依赖它，索引写了文件没写就是假数据
	for _, m := range idx.Versions {
		if _, err := os.Stat(filepath.Join(s.historyDir("deck-9001"), m.Version+".html")); err != nil {
			t.Fatalf("快照文件缺失 %s: %v", m.Version, err)
		}
	}
}

func TestPruneVersionsTiered(t *testing.T) {
	now := time.Now()
	mk := func(daysAgo, seq int) VersionMeta {
		return VersionMeta{Version: fmt.Sprintf("v%06d", seq), Time: now.AddDate(0, 0, -daysAgo).Unix()}
	}
	// 新→旧：今天两条全留、昨天一条留、8 天前两条只留当天最后记的那条
	versions := []VersionMeta{mk(0, 5), mk(0, 4), mk(1, 3), mk(8, 2), mk(8, 1)}

	kept, dropped := pruneVersions(versions)
	if len(kept) != 4 || len(dropped) != 1 || dropped[0] != "v000001" {
		t.Fatalf("分层裁剪结果不对: kept=%d dropped=%v", len(kept), dropped)
	}
	if kept[0].Version != "v000005" || kept[3].Version != "v000002" {
		t.Fatalf("保留列表顺序不对: %+v", kept)
	}
}

func TestPruneCap(t *testing.T) {
	now := time.Now()
	total := maxVersionsPerDeck + 5
	// 按真实顺序造数据：新→旧，编号递减（最新的一条编号最大）
	versions := make([]VersionMeta, total)
	for i := range versions {
		versions[i] = VersionMeta{Version: fmt.Sprintf("v%06d", total-i), Time: now.Unix()}
	}
	kept, dropped := pruneVersions(versions)
	if len(kept) != maxVersionsPerDeck || len(dropped) != 5 {
		t.Fatalf("总数兜底不对: kept=%d dropped=%d", len(kept), len(dropped))
	}
	// 截尾：保留最新的 200 条，丢掉最旧的 5 条
	if kept[0].Version != fmt.Sprintf("v%06d", total) {
		t.Fatalf("应保留最新一条，实际 %s", kept[0].Version)
	}
	if kept[len(kept)-1].Version != "v000006" {
		t.Fatalf("保留区间的尾部不对: %s", kept[len(kept)-1].Version)
	}
	for _, v := range dropped {
		if v != "v000005" && v != "v000004" && v != "v000003" && v != "v000002" && v != "v000001" {
			t.Fatalf("丢掉的不该是最新的: %s", v)
		}
	}
}

func TestDeleteVersion(t *testing.T) {
	s := newHistoryTestService(t)
	if err := writeDeckHTML(s, "deck-9002", deckHTML(1)); err != nil {
		t.Fatal(err)
	}
	s.recordVersion("deck-9002", OpRun, "r1", deckHTML(1))
	s.recordVersion("deck-9002", OpRun, "r2", deckHTML(2))
	s.recordVersion("deck-9002", OpRun, "r3", deckHTML(3))

	// 非法版本号要拦住（白名单，防路径穿越）
	if err := s.deleteVersionLocked("deck-9002", "../../etc/passwd"); err == nil {
		t.Fatal("非法版本号未被拦截")
	}
	if err := s.deleteVersionLocked("deck-9002", "v000002"); err != nil {
		t.Fatal(err)
	}

	idx := s.readHistoryIndex("deck-9002")
	if len(idx.Versions) != 2 {
		t.Fatalf("删除后应剩 2 条: %d", len(idx.Versions))
	}
	if _, err := os.Stat(filepath.Join(s.historyDir("deck-9002"), "v000002.html")); !os.IsNotExist(err) {
		t.Fatal("v000002 快照文件应已删除")
	}
	// 删中间不影响其他版本恢复——快照互相独立
	if err := s.restoreSnapshot("deck-9002", "v000001"); err != nil {
		t.Fatalf("删中间版本后其余版本应仍可恢复: %v", err)
	}
}

func TestClearHistoryKeepsSeq(t *testing.T) {
	s := newHistoryTestService(t)
	if err := writeDeckHTML(s, "deck-9003", deckHTML(1)); err != nil {
		t.Fatal(err)
	}
	s.recordVersion("deck-9003", OpRun, "r1", deckHTML(1))
	s.recordVersion("deck-9003", OpRun, "r2", deckHTML(2))

	n, err := s.clearHistoryLocked("deck-9003")
	if err != nil || n != 2 {
		t.Fatalf("清空失败: n=%d err=%v", n, err)
	}
	// 清空后编号继续：新版本是 v000003 而不是 v000001，编号永不复用
	s.recordVersion("deck-9003", OpRun, "r3", deckHTML(3))
	idx := s.readHistoryIndex("deck-9003")
	if len(idx.Versions) != 1 || idx.Versions[0].Version != "v000003" {
		t.Fatalf("清空后编号应延续: %+v", idx.Versions)
	}
}

func TestRestoreSnapshotReversible(t *testing.T) {
	s := newHistoryTestService(t)
	v1 := deckHTML(3)
	if err := writeDeckHTML(s, "deck-9004", v1); err != nil {
		t.Fatal(err)
	}
	s.recordVersion("deck-9004", OpRun, "r1", v1)

	v2 := deckHTML(5)
	if err := writeDeckHTML(s, "deck-9004", v2); err != nil {
		t.Fatal(err)
	}
	s.recordVersion("deck-9004", OpRun, "r2", v2)

	if err := s.restoreSnapshot("deck-9004", "v000001"); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(s.decksDir, "deck-9004", "deck.html"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != v1 {
		t.Fatal("恢复后内容应逐字节等于快照")
	}
	if s.readHistoryIndex("deck-9004").Versions[0].Operation != OpRestore {
		t.Fatal("恢复动作本身也要记一条历史，否则这一步不可撤销")
	}

	// 可逆：再回到 v2
	if err := s.restoreSnapshot("deck-9004", "v000002"); err != nil {
		t.Fatal(err)
	}
	got, _ = os.ReadFile(filepath.Join(s.decksDir, "deck-9004", "deck.html"))
	if string(got) != v2 {
		t.Fatal("二次恢复后内容应为 v2")
	}

	// 未知版本要报错，不能静默成功
	if err := s.restoreSnapshot("deck-9004", "v999999"); err == nil {
		t.Fatal("恢复不存在的版本应报错")
	}
}

// TestConcurrentRecordVersions 并发记档：锁必须覆盖"写文件 + 读改索引"整段，
// 否则会丢版本（索引被后写者覆盖）。这里用 n 个 goroutine 压一遍。
func TestConcurrentRecordVersions(t *testing.T) {
	s := newHistoryTestService(t)
	if err := writeDeckHTML(s, "deck-9005", deckHTML(1)); err != nil {
		t.Fatal(err)
	}
	s.recordVersion("deck-9005", OpRun, "r0", deckHTML(1))

	const n = 20
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			html := deckHTML(i + 2)
			unlock := s.lockDeck("deck-9005")
			defer unlock()
			if err := writeDeckHTML(s, "deck-9005", html); err != nil {
				t.Error(err) // goroutine 里不能用 t.Fatal
				return
			}
			s.recordVersion("deck-9005", OpRun, fmt.Sprintf("g%d", i), html)
		}(i)
	}
	wg.Wait()

	idx := s.readHistoryIndex("deck-9005")
	if len(idx.Versions) != n+1 || idx.NextSeq != n+2 {
		t.Fatalf("并发记档有丢失: len=%d next=%d（期望 %d/%d）",
			len(idx.Versions), idx.NextSeq, n+1, n+2)
	}
}
