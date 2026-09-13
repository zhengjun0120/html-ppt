package deck

// read_history_diff（版本对比）的回归测试。
// 语义要点：from 不传 = 最新记档版本，to 不传 = 当前使用中的内容；
// 都不传时得到"最新记档点到当前"的差集，也就是"这一轮已做的修改"。

import (
	"fmt"
	"strings"
	"testing"
)

func TestVersionDiffDefaultComparesCurrent(t *testing.T) {
	s := newHistoryTestService(t)
	v1 := deckWithTitles(map[string]string{"s1": "标题一", "s2": "标题二"})
	if err := writeDeckHTML(s, "deck-9006", v1); err != nil {
		t.Fatal(err)
	}
	s.recordVersion("deck-9006", OpRun, "第一轮", v1)

	// 模拟 run 进行中的在途修改：内容已变但还没到记档点
	cur := deckWithTitles(map[string]string{"s1": "标题一（改）", "s3": "新增页"})
	if err := writeDeckHTML(s, "deck-9006", cur); err != nil {
		t.Fatal(err)
	}

	d, err := s.versionDiff("deck-9006", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if d.From != "v000001" || d.To != "current" {
		t.Fatalf("默认应对比 current: %s → %s", d.From, d.To)
	}
	if !d.Changed {
		t.Fatal("应检测到差异")
	}
	if len(d.Modified) != 1 || d.Modified[0].SlideID != "s1" {
		t.Fatalf("s1 应判定为 modified: %+v", d.Modified)
	}
	if !strings.Contains(d.Modified[0].Diff, "标题一（改）") {
		t.Fatalf("diff 应包含新内容: %s", d.Modified[0].Diff)
	}
	if len(d.Added) != 1 || d.Added[0].SlideID != "s3" {
		t.Fatalf("s3 应判定为 added: %+v", d.Added)
	}
	if len(d.Removed) != 1 || d.Removed[0].SlideID != "s2" {
		t.Fatalf("s2 应判定为 removed: %+v", d.Removed)
	}
}

func TestVersionDiffNoChange(t *testing.T) {
	s := newHistoryTestService(t)
	v1 := deckWithTitles(map[string]string{"s1": "A", "s2": "B"})
	if err := writeDeckHTML(s, "deck-9007", v1); err != nil {
		t.Fatal(err)
	}
	s.recordVersion("deck-9007", OpRun, "r1", v1)

	// 记档后当前内容 == 最新版本：新一轮开始时的正常状态，应为"无差异"
	d, err := s.versionDiff("deck-9007", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if d.Changed || d.Unchanged != 2 {
		t.Fatalf("应无差异: changed=%v unchanged=%d", d.Changed, d.Unchanged)
	}
}

func TestVersionDiffCumulativeAndBetweenVersions(t *testing.T) {
	s := newHistoryTestService(t)
	v1 := deckWithTitles(map[string]string{"s1": "一", "s2": "二"})
	if err := writeDeckHTML(s, "deck-9008", v1); err != nil {
		t.Fatal(err)
	}
	s.recordVersion("deck-9008", OpRun, "r1", v1)

	v2 := deckWithTitles(map[string]string{"s1": "一（改）", "s2": "二"})
	if err := writeDeckHTML(s, "deck-9008", v2); err != nil {
		t.Fatal(err)
	}
	s.recordVersion("deck-9008", OpRun, "r2", v2)

	cur := deckWithTitles(map[string]string{"s1": "一（改）", "s2": "二（也改）", "s3": "新"})
	if err := writeDeckHTML(s, "deck-9008", cur); err != nil {
		t.Fatal(err)
	}
	s.recordVersion("deck-9008", OpRun, "r3", cur)

	// 累计差异：v1 → 当前（两处修改 + 一处新增）
	d, err := s.versionDiff("deck-9008", "v000001", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Modified) != 2 || len(d.Added) != 1 {
		t.Fatalf("累计差异不对: modified=%d added=%d", len(d.Modified), len(d.Added))
	}

	// 两个历史版本之间：v1 → v2，不应包含 v2 之后的改动
	d2, err := s.versionDiff("deck-9008", "v000001", "v000002")
	if err != nil {
		t.Fatal(err)
	}
	if d2.To != "v000002" || d2.ToDetail != "r2" {
		t.Fatalf("to 侧应为 v2 及其操作汇总: %+v", d2)
	}
	if len(d2.Modified) != 1 || d2.Modified[0].SlideID != "s1" {
		t.Fatalf("v1→v2 只应有 s1 变化: %+v", d2.Modified)
	}
	if len(d2.Added) != 0 {
		t.Fatalf("s3 是 v2 之后才加的，不该出现: %+v", d2.Added)
	}
}

// TestVersionDiffTruncation 行数上限：AI 提交的 HTML 带换行时，
// 一次改动就是几百行 diff，必须裁到 maxDiffLinesPerSlide。
func TestVersionDiffTruncation(t *testing.T) {
	s := newHistoryTestService(t)
	var b1, b2 strings.Builder
	b1.WriteString("<div class=\"reveal\"><div class=\"slides\"><section data-id=\"s1\">\n")
	b2.WriteString("<div class=\"reveal\"><div class=\"slides\"><section data-id=\"s1\">\n")
	for i := 0; i < 200; i++ {
		fmt.Fprintf(&b1, "<p>旧内容 %d</p>\n", i)
		fmt.Fprintf(&b2, "<p>新内容 %d</p>\n", i)
	}
	b1.WriteString("</section></div></div>")
	b2.WriteString("</section></div></div>")

	if err := writeDeckHTML(s, "deck-9009", b1.String()); err != nil {
		t.Fatal(err)
	}
	s.recordVersion("deck-9009", OpRun, "r1", b1.String())
	if err := writeDeckHTML(s, "deck-9009", b2.String()); err != nil {
		t.Fatal(err)
	}

	d, err := s.versionDiff("deck-9009", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if !d.Modified[0].Truncated {
		t.Fatal("超长 diff 应标记截断")
	}
	// 行数统计：@@ 头 1 行 + 内容行 + 结尾 "..." 1 行
	lines := strings.Count(d.Modified[0].Diff, "\n")
	if content := lines - 2; content > maxDiffLinesPerSlide {
		t.Fatalf("diff 内容行数应不超过 %d，实际 %d", maxDiffLinesPerSlide, content)
	}
	if !strings.HasSuffix(d.Modified[0].Diff, "\n...") {
		t.Fatal("截断的 diff 应以 ... 结尾")
	}
}

// TestVersionDiffSingleLineCap 只按行数截断挡不住"整页一行"：
// 那种情况下 diff 只有一两条超长行，必须靠单行字符上限兜住 token。
func TestVersionDiffSingleLineCap(t *testing.T) {
	s := newHistoryTestService(t)
	// 全程不换行：序列化后整页就是一行
	var b1, b2 strings.Builder
	b1.WriteString(`<div class="reveal"><div class="slides"><section data-id="s1">`)
	b2.WriteString(`<div class="reveal"><div class="slides"><section data-id="s1">`)
	for i := 0; i < 300; i++ {
		fmt.Fprintf(&b1, `<p>旧内容 %d</p>`, i)
		fmt.Fprintf(&b2, `<p>新内容 %d</p>`, i)
	}
	b1.WriteString(`</section></div></div>`)
	b2.WriteString(`</section></div></div>`)

	if err := writeDeckHTML(s, "deck-9011", b1.String()); err != nil {
		t.Fatal(err)
	}
	s.recordVersion("deck-9011", OpRun, "r1", b1.String())
	if err := writeDeckHTML(s, "deck-9011", b2.String()); err != nil {
		t.Fatal(err)
	}

	d, err := s.versionDiff("deck-9011", "", "")
	if err != nil {
		t.Fatal(err)
	}
	diff := d.Modified[0].Diff
	// 每一行（去掉前缀）都不超过上限 + 省略号 + 容忍度
	for _, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "@@") || line == "..." {
			continue
		}
		body := []rune(strings.TrimPrefix(strings.TrimPrefix(line, "-"), "+"))
		if strings.TrimSpace(line) == "" {
			continue
		}
		if len(body) > maxDiffLineChars+2 {
			t.Fatalf("单行超长未被截断（%d 字符）: %s", len(body), string(body[:60]))
		}
	}
	if !strings.Contains(diff, "…") {
		t.Fatal("超长行截断后应留下 … 标记")
	}
	if len([]rune(diff)) > maxDiffLineChars*4 {
		t.Fatalf("整页一行的 diff 仍然过大: %d 字符", len([]rune(diff)))
	}
}

func TestVersionDiffEdgeCases(t *testing.T) {
	s := newHistoryTestService(t)
	if err := writeDeckHTML(s, "deck-9010", deckHTML(1)); err != nil {
		t.Fatal(err)
	}

	// 还没有任何记档版本
	if _, err := s.versionDiff("deck-9010", "", ""); err == nil {
		t.Fatal("无历史时应报错")
	}

	s.recordVersion("deck-9010", OpRun, "r1", deckHTML(1))

	// 只有一个版本也能比较（最新 vs 当前）；记档后应无差异
	d, err := s.versionDiff("deck-9010", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if d.Changed {
		t.Fatal("记档后当前与最新一致，应无差异")
	}

	// 基线与目标同一个版本
	if _, err := s.versionDiff("deck-9010", "v000001", "v000001"); err == nil {
		t.Fatal("from=to 应报错")
	}
	// 不存在的版本
	if _, err := s.versionDiff("deck-9010", "v999999", ""); err == nil {
		t.Fatal("未知版本应报错")
	}
	// 非法版本号格式
	if _, err := s.versionDiff("deck-9010", "../../x", ""); err == nil {
		t.Fatal("非法版本号应报错")
	}
}
