package deck

// deck-v2 存储层的全链路回归：建草稿 → 写大纲 → 确认 → 选模板（实例化）→
// 规划 → 批写入 → 页级修改 → 历史快照/恢复。
//
// 走真实仓库的 tech-sharing 模板 + 内存 SQLite——这是"契约测试"而不只是单测：
// template 包的校验规则、deck 的写入闸门、模板的 layouts.md 三者任何一个漂移，
// 这批测试第一个炸。

import (
	"strings"
	"testing"
)

func newV2TestService(t *testing.T) *Service {
	t.Helper()
	st, err := newTestStore()
	if err != nil {
		t.Fatalf("内存库: %v", err)
	}
	return &Service{decksDir: t.TempDir(), st: st, templates: newTestRegistry(t)}
}

// demoOutline 一份 6 页的合法大纲。
func demoOutline() *Outline {
	return &Outline{
		Title: "Go 并发入门",
		Meta:  OutlineMeta{Audience: "后端工程师", DurationMin: 15, PageCount: 6},
		Pages: []OutlinePage{
			{No: 1, Role: RoleCover, Title: "Go 并发入门"},
			{No: 2, Role: RoleTOC, Title: "路线图"},
			{No: 3, Role: RoleContent, Title: "为什么需要并发", Points: []string{"阻塞", "成本"}},
			{No: 4, Role: RoleCode, Title: "一个 goroutine"},
			{No: 5, Role: RoleData, Title: "百万协程"},
			{No: 6, Role: RoleThanks, Title: "Q&A"},
		},
	}
}

func mustStage(t *testing.T, s *Service, userID uint, id, want string) {
	t.Helper()
	df, err := s.GetDeckV2(userID, id)
	if err != nil {
		t.Fatalf("读 deck: %v", err)
	}
	if df.Stage != want {
		t.Fatalf("阶段 = %s, 期望 %s", df.Stage, want)
	}
}

func TestV2FullPipeline(t *testing.T) {
	s := newV2TestService(t)
	const uid = 7

	// 建草稿（submit_outline 的第一步）
	id, err := s.CreateV2Draft(uid, "Go 并发入门")
	if err != nil {
		t.Fatalf("建草稿: %v", err)
	}
	mustStage(t, s, uid, id, StageOutlining)

	// 写大纲（首次，version=1）
	o := demoOutline()
	if err := s.SaveOutline(uid, id, o, 0); err != nil {
		t.Fatalf("写大纲: %v", err)
	}
	got, err := s.ReadOutline(uid, id)
	if err != nil || got.Version != 1 {
		t.Fatalf("读大纲: %v version=%d", err, got.Version)
	}

	// CAS：拿 v1 改 v2，再用 v1 提交 → 冲突
	o2 := demoOutline()
	o2.Pages[2].Title = "改过的标题"
	if err := s.SaveOutline(uid, id, o2, 1); err != nil {
		t.Fatalf("CAS 写 v2: %v", err)
	}
	o3 := demoOutline()
	err = s.SaveOutline(uid, id, o3, 1)
	if _, ok := err.(OutlineConflict); !ok {
		t.Fatalf("过期版本应报 OutlineConflict，得到: %v", err)
	}

	// 非法大纲被拒：页码不连续
	bad := demoOutline()
	bad.Pages[3].No = 9
	bad.Version = 2
	if err := s.SaveOutline(uid, id, bad, 2); err == nil || !strings.Contains(err.Error(), "连续") {
		t.Fatalf("页码不连续应被拒: %v", err)
	}

	// 确认大纲
	if _, err := s.ConfirmOutline(uid, id); err != nil {
		t.Fatalf("确认大纲: %v", err)
	}
	mustStage(t, s, uid, id, StageSelectingTemplate)

	// 选模板（实例化）
	if _, err := s.SelectTemplate(uid, id, "tech-sharing", "blue"); err != nil {
		t.Fatalf("选模板: %v", err)
	}
	mustStage(t, s, uid, id, StageGenerating)
	html, err := s.GetHTML(uid, id)
	if err != nil {
		t.Fatalf("读实例化产物: %v", err)
	}
	if !strings.Contains(html, `class="tpl-tech-sharing v-blue"`) {
		t.Error("变体 class 没挂上 body")
	}
	if strings.Contains(html, "Rust 异步运行时") {
		t.Error("demo 内容没被剥离")
	}
	// 空骨架时量测/列表为 0 页
	slides, err := s.ListSlidesV2(uid, id)
	if err != nil || len(slides) != 0 {
		t.Fatalf("实例化后应为 0 页: %v %d", err, len(slides))
	}

	// 规划（page_plan）
	plan := []PlanAssignment{
		{1, "cover", ""}, {2, "agenda", ""}, {3, "cards-3", ""},
		{4, "code-terminal", ""}, {5, "stat-hero", ""}, {6, "qa", ""},
	}
	if _, err := s.SavePlanV2(uid, id, plan); err != nil {
		t.Fatalf("保存计划: %v", err)
	}

	// 批写入：2 好 1 坏（未知类名），部分成功
	pages := []PageInput{
		{No: 1, Layout: "cover", HTML: `<section class="slide" data-layout="cover"><p class="kicker">test / 2026</p><h1 class="h1">Go 并发入门</h1><p class="lede mt-m">一次讲清楚。</p><div class="speaker"><div class="av"></div><div><b>@t</b><span>x · 15 min</span></div></div><div class="notes">讲稿</div></section>`},
		{No: 4, Layout: "code-terminal", HTML: `<section class="slide" data-layout="code-terminal"><p class="kicker">main.go · 4 LOC</p><h2 class="h2">起一个 goroutine</h2><div class="terminal mt-m"><div class="bar"><span class="dot"></span><span class="dot"></span><span class="dot"></span><span>main.go</span></div><pre><span class="kw">go</span> <span class="fn">worker</span>()</pre></div><div class="notes">讲稿</div></section>`},
		{No: 3, Layout: "cards-3", HTML: `<section class="slide" data-layout="cards-3"><h2 class="h2">标题</h2><div class="grid g3 mt-l"><div class="card ghost-class">x</div></div></section>`},
	}
	rep, err := s.WritePagesV2(uid, id, pages)
	if err != nil {
		t.Fatalf("批写入: %v", err)
	}
	if rep.Written != 2 {
		t.Fatalf("应写入 2 页，实际 %d（%+v）", rep.Written, rep.Results)
	}
	for _, r := range rep.Results {
		if r.No == 3 && r.OK {
			t.Error("未知类名的页应当失败")
		}
		if r.No == 3 && !strings.Contains(r.Error, "ghost-class") {
			t.Errorf("错误信息应包含未知类名，得到: %s", r.Error)
		}
	}

	// 修复第 3 页后重写 → 成功；页序保持 1,3,4
	pagesFix := []PageInput{
		{No: 3, Layout: "cards-3", HTML: `<section class="slide" data-layout="cards-3"><h2 class="h2">为什么需要并发</h2><div class="grid g3 mt-l"><div class="card card-accent"><h4>阻塞</h4><p class="dim">一个阻塞调用拖住整个线程。</p><span class="tag mt-s">成本高</span></div><div class="card card-accent"><h4>线程贵</h4><p class="dim">栈 2-8MB，万级连接就要几十 GB。</p><span class="tag mt-s">扛不住</span></div><div class="card card-accent"><h4>协程轻</h4><p class="dim">KB 级栈，一台机器百万个。</p><span class="tag mt-s">Go 选它</span></div></div><div class="notes">讲稿</div></section>`},
		{No: 6, Layout: "qa", HTML: `<section class="slide center tc" data-layout="qa"><div><div class="mono" style="font-size:120px;color:var(--accent);font-weight:800">?</div><h2 class="h2">Questions?</h2><p class="lede" style="margin:14px auto">欢迎提问。</p><div class="row mt-l" style="justify-content:center"><span class="tag">repo</span></div></div><div class="notes">讲稿</div></section>`},
	}
	rep2, err := s.WritePagesV2(uid, id, pagesFix)
	if err != nil {
		t.Fatalf("修复批写入: %v", err)
	}
	if rep2.Written != 2 {
		t.Fatalf("修复批应写 2 页: %+v", rep2.Results)
	}

	slides, err = s.ListSlidesV2(uid, id)
	if err != nil {
		t.Fatalf("列页: %v", err)
	}
	if len(slides) != 4 {
		t.Fatalf("应有 4 页，实际 %d", len(slides))
	}
	// 页序：s1, s3, s4, s6
	wantOrder := []string{"s1", "s3", "s4", "s6"}
	for i, w := range wantOrder {
		if slides[i].ID != w {
			t.Errorf("第 %d 页 id=%s，期望 %s（乱序问题）", i+1, slides[i].ID, w)
		}
	}

	// plan 不一致：第 3 页计划是 cards-3，用 cover 提交 → 拒
	pagesWrong := []PageInput{{No: 3, Layout: "cover", HTML: `<section class="slide" data-layout="cover"><h1 class="h1">x</h1></section>`}}
	rep3, err := s.WritePagesV2(uid, id, pagesWrong)
	if err != nil {
		t.Fatalf("plan 拒绝批: %v", err)
	}
	if rep3.Results[0].OK || !strings.Contains(rep3.Results[0].Error, "plan_pages") {
		t.Errorf("plan 不一致应给出可执行错误: %+v", rep3.Results[0])
	}

	// 迭代期页级修改：insert + delete + fingerprint
	newID, _, err := s.InsertSlideV2(uid, id, "end", `<section class="slide" data-layout="big-quote"><div class="mono" style="font-size:120px;color:var(--accent);line-height:.6">"</div><h2 class="h2" style="max-width:62ch">并发不是并行。</h2><p class="dim mt-m">—— Rob Pike</p><div class="notes">讲稿</div></section>`)
	if err != nil {
		t.Fatalf("插入: %v", err)
	}
	h1, fp, err := s.ReadSlideV2(uid, id, "s1")
	if err != nil {
		t.Fatalf("读页: %v", err)
	}
	if !strings.Contains(h1, "Go 并发入门") {
		t.Error("读回的内容不对")
	}
	// 用过期指纹改 → 拒
	if _, err := s.UpdateSlideV2(uid, id, "s1", strings.Replace(h1, "Go 并发入门", "X", 1), "deadbeef0000"); err == nil {
		t.Error("过期指纹应被拒")
	}
	if _, err := s.UpdateSlideV2(uid, id, "s1", h1, fp); err != nil {
		t.Fatalf("正常更新: %v", err)
	}
	// 删页
	if err := s.DeleteSlideV2(uid, id, newID); err != nil {
		t.Fatalf("删页: %v", err)
	}
	if slides, _ := s.ListSlidesV2(uid, id); len(slides) != 4 {
		t.Errorf("删后应剩 4 页，实际 %d", len(slides))
	}

	// 历史：记档 → 恢复
	if err := s.RecordRunVersion(uid, id, "测试快照"); err != nil {
		t.Fatalf("记档: %v", err)
	}
	versions, err := s.ListVersions(uid, id)
	if err != nil || len(versions) == 0 {
		t.Fatalf("历史为空: %v", err)
	}
	before, _ := s.ListSlidesV2(uid, id)
	// 删一页制造差异，再恢复
	last := before[len(before)-1].ID
	_ = s.DeleteSlideV2(uid, id, last)
	if err := s.RestoreVersionV2(uid, id, versions[0].Version); err != nil {
		t.Fatalf("恢复: %v", err)
	}
	after, _ := s.ListSlidesV2(uid, id)
	if len(after) != len(before) {
		t.Errorf("恢复后页数 %d != 快照时 %d", len(after), len(before))
	}
}

// TestV2RhythmBlocking R101 必须阻塞 plan。
func TestV2RhythmBlocking(t *testing.T) {
	vs := ValidateRhythm(
		map[int]string{1: "cards-3", 2: "cards-3", 3: "cards-3", 4: "cover", 5: "cover", 6: "cover", 7: "agenda", 8: "agenda", 9: "agenda"},
		[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
	)
	blocked := false
	for _, v := range vs {
		if v.Rule == "R101" && v.Block {
			blocked = true
		}
	}
	if !blocked {
		t.Fatalf("连续 3 页同版式必须触发阻塞级 R101: %+v", vs)
	}
}
