package trace

// 读取端的回归测试。读侧比起写侧有一个更尖锐的责任：它同时负责**归属校验**。
// 图片路由直接拿文件系统响应，run_id / 会话 id / 图片名三个都来自 URL——
// 这里漏掉任何一处白名单校验，后果是别人能读到你任意文件（../../ 之类）。
// 所以下面关于路径穿越的几条不是"补充覆盖"，而是这个包的安全闸门本身。

import (
	"context"
	"encoding/json"

	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type seed struct {
	sessID  uint
	userID  uint
	runID   string
	deckID  string
	content string
	end     bool
	usage   map[string]UsagePart
}

// seedRun 造一个"看起来像真跑过"的 run。
func seedRun(t *testing.T, root string, s seed) {
	t.Helper()
	r, err := New(Options{
		Dir: root, SessionID: s.sessID, RunID: s.runID, UserID: s.userID,
		DeckID: s.deckID, CaptureImages: true,
	})
	if err != nil {
		t.Fatalf("建 recorder 失败: %v", err)
	}
	defer r.Close()

	ctx := With(context.Background(), r)
	Emit(ctx, Event{
		Kind: KindRunStart, RunID: s.runID, SessionID: s.sessID, UserID: s.userID,
		DeckID: s.deckID, UserContent: s.content, Model: "test-model",
	})

	ctx = WithTurn(ctx, 0)
	Emit(ctx, Event{
		Kind:         KindLLMRequest,
		Messages:     json.RawMessage(`[{"role":"user","content":"` + s.content + `"}]`),
		MessageCount: 1,
		Bytes:        32,
	})
	toolCtx := WithTool(ctx, "write_deck", "call_1")
	Emit(toolCtx, Event{Kind: KindToolCall, Args: `{"title":"t"}`})
	Emit(toolCtx, Event{Kind: KindToolResult, Result: `{"deck_id":"deck-0001"}`, DurationMS: 88})

	var total UsagePart
	for k, v := range s.usage {
		Usage(ctx, k, v)
		total = total.Add(v)
	}
	if s.end {
		Emit(ctx, Event{Kind: KindRunEnd, Status: StatusOK, Summary: &Summary{
			DurationMS: 1234, Turns: 1, ToolCalls: 1, Usage: s.usage, Total: total,
		}})
	}
}

// TestListRunsNewestFirstAndSummarises run 列表按新→旧，且把 run_end 的汇总摊到元信息上。
//
// 顺序反了或者汇总丢了，列表页看起来仍然"有东西"——只是你会去点开一个错的 run，
// 或者在最需要 token 数字的地方看到空白。所以顺序和汇总一起断言。
func TestListRunsNewestFirstAndSummarises(t *testing.T) {
	root := t.TempDir()
	usage := map[string]UsagePart{
		CompMain:   {Prompt: 100, Completion: 20, Total: 120},
		CompVision: {Prompt: 900, Completion: 60, Total: 960, Calls: 1},
	}
	for _, id := range []string{"1700000000001-a", "1700000000002-b", "1700000000003-c"} {
		seedRun(t, root, seed{sessID: 9, userID: 4, runID: id, deckID: "deck-0007", content: id, end: true, usage: usage})
	}

	st := NewStore(root)
	metas, err := st.ListRuns(4, 0, 0, 0)
	if err != nil {
		t.Fatalf("列 run 失败: %v", err)
	}
	if len(metas) != 3 {
		t.Fatalf("应有 3 个 run，实际 %d", len(metas))
	}
	if metas[0].RunID != "1700000000003-c" || metas[2].RunID != "1700000000001-a" {
		t.Errorf("列表应按新→旧，实际 %s, %s, %s", metas[0].RunID, metas[1].RunID, metas[2].RunID)
	}
	m := metas[0]
	if m.SessionID != 9 || m.DeckID != "deck-0007" || m.Status != StatusOK {
		t.Errorf("元信息不对: %+v", m)
	}
	if m.Usage == nil {
		t.Fatal("已结束的 run 必须带汇总，否则列表上 token 是空的")
	}
	if m.Usage.Total.Total != 1080 {
		t.Errorf("汇总总 token 应为 1080（120+960），实际 %d", m.Usage.Total.Total)
	}
	if m.Turns != 1 || m.ToolCalls != 1 {
		t.Errorf("轮数/工具调用数不对: turns=%d tools=%d", m.Turns, m.ToolCalls)
	}
	if m.Bytes == 0 {
		t.Error("应带上文件大小，页面靠它判断这个 run 是不是很重")
	}

	// 分页
	if got, _ := st.ListRuns(4, 0, 2, 0); len(got) != 2 || got[0].RunID != "1700000000003-c" {
		t.Errorf("limit=2 应返回最新两条，实际 %+v", got)
	}
	if got, _ := st.ListRuns(4, 0, 0, 2); len(got) != 1 || got[0].RunID != "1700000000001-a" {
		t.Errorf("offset=2 应返回最旧那一条，实际 %+v", got)
	}
	if got, _ := st.ListRuns(4, 0, 0, 99); got == nil || len(got) != 0 {
		t.Errorf("offset 越界应返回空切片而不是 nil（JSON 要是 [] 不是 null），实际 %+v", got)
	}
}

// TestListRunsFiltersByOwner 别人的 run 不能出现在你的列表里。
func TestListRunsFiltersByOwner(t *testing.T) {
	root := t.TempDir()
	seedRun(t, root, seed{sessID: 1, userID: 1, runID: "1700000000001-a", content: "我的"})
	seedRun(t, root, seed{sessID: 2, userID: 2, runID: "1700000000002-b", content: "别人的"})

	st := NewStore(root)
	metas, err := st.ListRuns(1, 0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(metas) != 1 || metas[0].UserContent != "我的" {
		t.Errorf("只应看到自己的 run，实际 %+v", metas)
	}
	if got, _ := st.ListRuns(99, 0, 0, 0); len(got) != 0 {
		t.Errorf("没有记录的用户应看到空列表，实际 %+v", got)
	}
}

// TestRejectsPathTraversalAndForeignRuns 三个入参都来自 URL，任何一个漏校验都能读到别人的文件。
//
// 这是本包唯一的安全闸门。特别注意**两种失败必须返回同一个错误**：
// 区分"不存在"和"不是你的"就等于告诉调用方"这个 id 是存在的"，
// 而 run_id 里带着会话号，泄露存在性等于泄露"这个用户聊过几次"。
func TestRejectsPathTraversalAndForeignRuns(t *testing.T) {
	root := t.TempDir()
	seedRun(t, root, seed{sessID: 1, userID: 1, runID: "1700000000001-a", content: "我的"})

	st := NewStore(root)

	// 会话目录外面放一份"别人的"文件，用来验证穿越真的被挡住了
	outside := filepath.Join(root, "..", "secret.jsonl")
	if err := os.WriteFile(outside, []byte(`{"kind":"run_start","user_id":1}`), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(outside) })

	badRuns := []string{
		"../../secret",
		"..",
		"1700000000001-a/../../secret",
		`..\..\secret`,
		"",
		"1700000000001-a.jsonl", // 带后缀：读侧会再拼一次 .jsonl，必须被白名单拦住
		strings.Repeat("x", 65),
	}
	for _, id := range badRuns {
		if _, err := st.ReadEvents(1, 1, id, ReadOptions{}); err != ErrNotFound {
			t.Errorf("run_id=%q 应被拒绝（ErrNotFound），实际 %v", id, err)
		}
		if _, err := st.ImagePath(1, 1, id, "p001.png"); err != ErrNotFound {
			t.Errorf("run_id=%q 在图片路由上应被拒绝，实际 %v", id, err)
		}
	}

	// 图片名比 run_id 宽松（要允许 v1-p003.png 这种"第几轮审查的第几页"），
	// 但仍然必须是"一个普通文件名 + .png"：带路径分隔符、带上级目录、扩展名不对的都要拦
	badImages := []string{"../../secret.png", "a/b.png", "p001.png/../../x.png", `..\x.png`, "p001.jpg", "", "p001.PNG", ".png"}
	for _, n := range badImages {
		if _, err := st.ImagePath(1, 1, "1700000000001-a", n); err != ErrNotFound {
			t.Errorf("图片名 %q 应被拒绝（ErrNotFound），实际 %v", n, err)
		}
	}
	for _, n := range []string{"p001.png", "v1-p003.png"} {
		if _, err := st.ImagePath(1, 1, "1700000000001-a", n); err != nil {
			t.Errorf("合法图片名 %q 不该被拒: %v", n, err)
		}
	}

	// 归属：同一个 run，换个用户问就是"不存在"，错误必须一模一样
	if _, err := st.ReadEvents(2, 1, "1700000000001-a", ReadOptions{}); err != ErrNotFound {
		t.Errorf("非属主读取应返回 ErrNotFound（不泄露存在性），实际 %v", err)
	}
	if _, err := st.ImagePath(2, 1, "1700000000001-a", "p001.png"); err != ErrNotFound {
		t.Errorf("非属主取图应返回 ErrNotFound，实际 %v", err)
	}
	if _, err := st.ReadEvent(2, 1, "1700000000001-a", 1); err != ErrNotFound {
		t.Errorf("非属主读单事件应返回 ErrNotFound，实际 %v", err)
	}

	// 存在但还没写过任何事件的 run（刚建好就来看）→ 同样是"不存在"
	r, err := New(Options{Dir: root, SessionID: 1, RunID: "1700000000009-empty", UserID: 1})
	if err != nil {
		t.Fatal(err)
	}
	r.Close()
	if _, err := st.ReadEvents(1, 1, "1700000000009-empty", ReadOptions{}); err != ErrNotFound {
		t.Errorf("空文件（刚创建还没写）应视为不存在，实际 %v", err)
	}
	if got, _ := st.ListRuns(1, 0, 0, 0); len(got) != 1 {
		t.Errorf("空文件不该出现在列表里，实际 %d 条", len(got))
	}
}

// TestReadEventsIncrementalByOffset 增量拉取必须只返回新增的那一段。
//
// 这是实时跟随的核心：页面每 800ms 拉一次，如果每次都把整个文件（开了全量上下文
// 就是几十 MB）重新解析一遍，浏览器和服务端都会被这个"观测工具"拖垮。
// 同时守住"游标只推进到完整行"——否则会读到半行，表现为偶发的 JSON 解析失败。
func TestReadEventsIncrementalByOffset(t *testing.T) {
	root := t.TempDir()
	runID := "1700000000001-inc"
	seedRun(t, root, seed{sessID: 3, userID: 1, runID: runID, content: "第一段", end: false})

	st := NewStore(root)
	first, err := st.ReadEvents(1, 3, runID, ReadOptions{})
	if err != nil {
		t.Fatalf("首次读取失败: %v", err)
	}
	if len(first.Events) != 4 {
		t.Fatalf("应读到 4 条事件，实际 %d", len(first.Events))
	}
	if !first.Running {
		t.Error("没有 run_end 时应标记为运行中")
	}
	if first.NextOffset <= 0 {
		t.Fatalf("应返回非零游标供下次增量，实际 %d", first.NextOffset)
	}

	// 再追加两条（模拟"下一轮又产生了事件"），游标前的部分不该被重复返回
	path := filepath.Join(sessDir(root, 3), runID+".jsonl")
	appendRaw(t, path,
		Event{Seq: 5, Kind: KindToolCall, ToolName: "read_slide", Args: `{"slide_id":"s1"}`},
		Event{Seq: 6, Kind: KindRunEnd, Status: StatusOK, Summary: &Summary{Turns: 2}},
	)

	second, err := st.ReadEvents(1, 3, runID, ReadOptions{FromOffset: first.NextOffset})
	if err != nil {
		t.Fatalf("增量读取失败: %v", err)
	}
	if len(second.Events) != 2 {
		t.Fatalf("增量应只返回新增的 2 条，实际 %d 条（重复返回旧事件 = 页面会重复渲染）", len(second.Events))
	}
	if second.Events[0].Seq != 5 || second.Events[1].Kind != KindRunEnd {
		t.Errorf("增量内容不对: %+v", second.Events)
	}
	if second.Running {
		t.Error("追加上 run_end 之后应报告已结束")
	}
	if second.NextOffset < first.NextOffset {
		t.Error("游标不该回退")
	}

	// FromSeq 也要能用（"从第 N 条开始看"），但它会从头扫，只用于一次性查询
	fromSeq, err := st.ReadEvents(1, 3, runID, ReadOptions{FromSeq: 3})
	if err != nil {
		t.Fatal(err)
	}
	if len(fromSeq.Events) != 3 || fromSeq.Events[0].Seq != 4 {
		t.Errorf("FromSeq=3 应返回 seq 4/5/6，实际 %+v", seqs(fromSeq.Events))
	}
}

// TestRunningFlagNotDerivedFromDelta 增量轮询时，run_end 不在这一段里，
// 但状态必须仍是"已结束"。
//
// 这条是修过的一个真 bug：Running 一开始是按"这一段里有没有见到 run_end"判断的，
// 增量轮询从游标之后开始读、永远读不到早已过去的 run_end —— 于是页面会一直显示
// "运行中"并一直轮询下去，而实际那次对话早就结束了。
func TestRunningFlagNotDerivedFromDelta(t *testing.T) {
	root := t.TempDir()
	runID := "1700000000001-end"
	seedRun(t, root, seed{sessID: 1, userID: 1, runID: runID, content: "跑完的", end: true})

	st := NewStore(root)
	all, err := st.ReadEvents(1, 1, runID, ReadOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if all.Running {
		t.Fatal("有 run_end 时应报告已结束")
	}

	// 从文件末尾之后读：一条事件都不返回，但状态仍必须是"已结束"
	empty, err := st.ReadEvents(1, 1, runID, ReadOptions{FromOffset: all.NextOffset})
	if err != nil {
		t.Fatal(err)
	}
	if len(empty.Events) != 0 {
		t.Fatalf("没有新增时不该返回事件，实际 %d 条", len(empty.Events))
	}
	if empty.Running {
		t.Error("增量段里没有 run_end 不代表还在运行：状态必须从文件尾部判断")
	}
}

// TestMessagesExcludedUnlessAsked 默认不返回 messages，按需才给。
//
// 全量上下文是文件里唯一的大字段（每轮重发整份消息，含整页 HTML）。
// 默认带上的话，列表页一次拉取就是几十 MB——而 99% 的查看动作并不需要看上下文。
func TestMessagesExcludedUnlessAsked(t *testing.T) {
	root := t.TempDir()
	runID := "1700000000001-msg"
	seedRun(t, root, seed{sessID: 1, userID: 1, runID: runID, content: "封面要深色", end: true})

	st := NewStore(root)
	def, err := st.ReadEvents(1, 1, runID, ReadOptions{})
	if err != nil {
		t.Fatal(err)
	}
	var req *Event
	for i := range def.Events {
		if def.Events[i].Kind == KindLLMRequest {
			req = &def.Events[i]
		}
	}
	if req == nil {
		t.Fatal("没找到 llm_request")
	}
	if req.Messages != nil {
		t.Errorf("默认不应返回 messages，实际 %s", req.Messages)
	}
	if req.MessageCount != 1 || req.Bytes == 0 {
		t.Errorf("剔除 messages 后仍要留下条数与字节数（实时流就靠这个表达上下文多大）: count=%d bytes=%d",
			req.MessageCount, req.Bytes)
	}

	with, err := st.ReadEvents(1, 1, runID, ReadOptions{IncludeMessages: true})
	if err != nil {
		t.Fatal(err)
	}
	for i := range with.Events {
		if with.Events[i].Kind == KindLLMRequest {
			req = &with.Events[i]
		}
	}
	if req == nil || !strings.Contains(string(req.Messages), "封面要深色") {
		t.Errorf("显式要求时应返回 messages 全文，实际 %s", string(req.Messages))
	}

	// 单事件接口是"点开看完整上下文"的唯一入口，它必须带全文
	one, err := st.ReadEvent(1, 1, runID, req.Seq)
	if err != nil {
		t.Fatalf("按 seq 读单事件失败: %v", err)
	}
	if !strings.Contains(string(one.Messages), "封面要深色") {
		t.Errorf("单事件接口应带 messages 全文，实际 %s", string(one.Messages))
	}
	if _, err := st.ReadEvent(1, 1, runID, 9999); err != ErrNotFound {
		t.Errorf("不存在的 seq 应返回 ErrNotFound，实际 %v", err)
	}
}

// TestReadEventsLimit 带 limit 时提前停下，游标落在已读行的末尾（不丢不重）。
func TestReadEventsLimit(t *testing.T) {
	root := t.TempDir()
	runID := "1700000000001-lim"
	seedRun(t, root, seed{sessID: 1, userID: 1, runID: runID, content: "分页", end: true})

	st := NewStore(root)
	page1, err := st.ReadEvents(1, 1, runID, ReadOptions{Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(page1.Events) != 2 || page1.Events[0].Seq != 1 || page1.Events[1].Seq != 2 {
		t.Fatalf("第一页应为 seq 1/2，实际 %+v", seqs(page1.Events))
	}
	page2, err := st.ReadEvents(1, 1, runID, ReadOptions{Limit: 2, FromOffset: page1.NextOffset})
	if err != nil {
		t.Fatal(err)
	}
	if len(page2.Events) != 2 || page2.Events[0].Seq != 3 {
		t.Fatalf("第二页应接着 seq 3 开始，实际 %+v", seqs(page2.Events))
	}
}

// TestSumBySession 跨 run 累加。一次对话被 ask_user 打断后会分成两个 run，
// 而"这次对话花了多少 token"是用户真正会问的那个问题——只看单个 run 只是半场。
func TestSumBySession(t *testing.T) {
	metas := []RunMeta{
		{RunID: "1700000000003-c", SessionID: 7, StartedAt: time.Unix(3, 0), Status: StatusOK,
			Usage: &Summary{Turns: 2, ToolCalls: 3, Total: UsagePart{Total: 100, Prompt: 80},
				Usage: map[string]UsagePart{CompMain: {Total: 100, Prompt: 80, Calls: 2}}}},
		{RunID: "1700000000002-b", SessionID: 7, StartedAt: time.Unix(2, 0), Status: StatusPaused,
			Usage: &Summary{Turns: 1, Total: UsagePart{Total: 40, Prompt: 30},
				Usage: map[string]UsagePart{CompVision: {Total: 40, Prompt: 30, Calls: 1}}}},
		{RunID: "1700000000001-a", SessionID: 8, StartedAt: time.Unix(1, 0), Status: StatusRunning},
	}

	got := SumBySession(metas)
	if len(got) != 2 {
		t.Fatalf("应聚合成 2 个会话，实际 %d", len(got))
	}
	s7 := got[0]
	if s7.SessionID != 7 || s7.Runs != 2 || s7.Latest != "1700000000003-c" {
		t.Errorf("会话 7 的聚合不对: %+v", s7)
	}
	if s7.Usage.Total.Total != 140 {
		t.Errorf("跨 run 的 token 应累加为 140，实际 %d", s7.Usage.Total.Total)
	}
	if len(s7.Usage.Usage) != 2 || s7.Usage.Usage[CompMain].Total != 100 || s7.Usage.Usage[CompVision].Total != 40 {
		t.Errorf("分项应保留各自的来源（这是回答\"token 花在哪\"的关键）: %+v", s7.Usage.Usage)
	}
	if s7.Usage.Turns != 3 || s7.Usage.ToolCalls != 3 {
		t.Errorf("轮数与工具调用数应累加: %+v", s7.Usage)
	}
	if s7.RunningRuns != 0 {
		t.Errorf("会话 7 没有运行中的 run，实际 %d", s7.RunningRuns)
	}
	if got[1].SessionID != 8 || got[1].RunningRuns != 1 {
		t.Errorf("会话 8 应报告 1 个运行中的 run（页面据此继续轮询）: %+v", got[1])
	}
	if got[1].Usage.Total.Total != 0 {
		t.Errorf("还在跑的 run 没有汇总，不该编一个出来: %+v", got[1].Usage)
	}
}

// TestStoreOnMissingDir 一个 trace 都还没产生时，列表要是空数组而不是错误。
// 这会在"刚部署、还没人聊过"时直接触发，报 500 会让人以为观测工具坏了。
func TestStoreOnMissingDir(t *testing.T) {
	st := NewStore(filepath.Join(t.TempDir(), "not-created-yet"))
	metas, err := st.ListRuns(1, 0, 0, 0)
	if err != nil {
		t.Fatalf("目录不存在不该报错: %v", err)
	}
	if metas == nil || len(metas) != 0 {
		t.Errorf("应返回空切片，实际 %+v", metas)
	}
	if _, err := st.ReadEvents(1, 1, "1700000000001-a", ReadOptions{}); err != ErrNotFound {
		t.Errorf("目录不存在时读取应返回 ErrNotFound，实际 %v", err)
	}
}

func seqs(evs []Event) []int {
	out := make([]int, 0, len(evs))
	for _, e := range evs {
		out = append(out, e.Seq)
	}
	return out
}

// appendRaw 直接往 jsonl 追一行，模拟"跑着跑着又产生了事件"。
// 刻意不走 Recorder：一个 run 只该有一个 Recorder（它的 seq 从 1 开始计数），
// 用第二个 Recorder 续写会让 seq 重头开始，那不是真实场景。
func appendRaw(t *testing.T, path string, evs ...Event) {
	t.Helper()
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	for _, e := range evs {
		if e.TS.IsZero() {
			e.TS = time.Now()
		}
		line, err := json.Marshal(e)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := f.Write(append(line, '\n')); err != nil {
			t.Fatal(err)
		}
	}
	if err := f.Sync(); err != nil {
		t.Fatal(err)
	}
}
