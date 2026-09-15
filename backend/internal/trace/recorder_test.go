package trace

// trace 写入端的回归测试。这个包是观测工具的地基，而它坏掉的方式几乎都是静默的：
// 事件少写一条、图片只存了路径没存内容、token 少累加一项——都不会报错，
// 只会让你在排查别人问题的时候，拿着一份缺了关键几行的"完整记录"下结论。
//
// 所以这里的测试重点不在"能不能跑通"，而在几条具体的、曾经被漏掉过的事实：
// 参数有没有真的落进文件、图片有没有真的写到磁盘、归属有没有串台、
// 分项用量有没有各自独立累加。

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func mustNew(t *testing.T, o Options) *Recorder {
	t.Helper()
	r, err := New(o)
	if err != nil {
		t.Fatalf("建 recorder 失败: %v", err)
	}
	t.Cleanup(r.Close) // 句柄不关会在 Windows 上让后面的删除失败
	return r
}

// readLines 把 jsonl 逐行解析回事件。用 bufio 之外的朴素做法：
// 测试里文件都很小，而"能读到几行"本身就是被测事实之一。
func readLines(t *testing.T, path string) []Event {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("读 trace 文件失败: %v", err)
	}
	var out []Event
	for _, line := range bytes.Split(raw, []byte("\n")) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var e Event
		if err := json.Unmarshal(line, &e); err != nil {
			t.Fatalf("第 %d 行不是合法 JSON: %v\n原文: %s", len(out)+1, err, line)
		}
		out = append(out, e)
	}
	return out
}

func runFilePath(root string, sessID uint, runID string) string {
	return filepath.Join(sessDir(root, sessID), runID+".jsonl")
}

// turnOf 取事件的轮次；没有归属时返回 -1。
// 用 -1（而不是 0）表示"不适用"，因为 0 是合法轮次——这正是 Turn 用指针的原因。
func turnOf(e Event) int {
	if e.Turn == nil {
		return -1
	}
	return *e.Turn
}

func sessDir(root string, sessID uint) string {
	return filepath.Join(root, strconv.FormatUint(uint64(sessID), 10))
}

// TestRecorderWritesPayloadAndImages 交代这个包最基本的三条承诺：
//   - 工具**参数**真的落进了文件（API 上是过去从没外发过的东西，SSE 的 tool_call 只有返回值）
//   - 每行一个事件、seq 单调，读侧才能按行切、按 seq 增量
//   - 图片落成磁盘文件，JSONL 里只留路径——几百 KB 的 PNG 内联进去会让文件不可读
//
// 顺带守住"PNG 字节没有混进 JSONL"：它是这条设计最容易被随手改掉的地方
// （图省事直接 base64 塞进去，测试全绿，文件暴涨）。
func TestRecorderWritesPayloadAndImages(t *testing.T) {
	root := t.TempDir()
	runID := "1700000000000-abcd"
	r := mustNew(t, Options{
		Dir: root, SessionID: 7, RunID: runID, UserID: 42, DeckID: "deck-0001",
		CaptureImages: true,
	})

	ctx := With(context.Background(), r)
	Emit(ctx, Event{Kind: KindRunStart, RunID: runID, SessionID: 7, UserID: 42, UserContent: "做个封面"})

	ctx = WithTurn(ctx, 0)
	Emit(ctx, Event{Kind: KindLLMRequest, Messages: json.RawMessage(`[{"role":"user","content":"做个封面"}]`), MessageCount: 1, Bytes: 38})

	toolCtx := WithTool(ctx, "write_deck", "call_x1")
	Emit(toolCtx, Event{Kind: KindToolCall, Args: `{"title":"封面"}`})

	png := []byte("PNG-BYTES-SHOULD-NOT-BE-INLINED")
	Emit(toolCtx, Event{
		Kind: KindSubStep,
		Sub:  &SubStep{Name: "vision", Stage: "capture"},
		Images: []ImageEvent{
			{Name: "p001.png", Label: "第 1 页", Bytes: png},
			{Name: "p002.png", Label: "第 2 页", Bytes: png},
		},
	})
	Emit(toolCtx, Event{Kind: KindToolResult, Result: `{"deck_id":"deck-0001"}`, DurationMS: 1234})

	Emit(ctx, Event{Kind: KindRunEnd, Status: StatusOK, Summary: &Summary{Turns: 1, ToolCalls: 1}})
	r.Close()

	evs := readLines(t, runFilePath(root, 7, runID))
	if len(evs) != 6 {
		t.Fatalf("应写出 6 条事件，实际 %d 条", len(evs))
	}
	for i, e := range evs {
		if e.Seq != i+1 {
			t.Errorf("第 %d 条的 seq 应为 %d，实际 %d", i+1, i+1, e.Seq)
		}
		if e.TS.IsZero() {
			t.Errorf("第 %d 条没有时间戳", i+1)
		}
	}

	// 参数原文必须原样保留：模型给了非法 JSON 的时候，"原文长什么样"正是要观测的东西
	var call *Event
	for i := range evs {
		if evs[i].Kind == KindToolCall {
			call = &evs[i]
		}
	}
	if call == nil {
		t.Fatal("没找到 tool_call 事件")
	}
	if call.Args != `{"title":"封面"}` {
		t.Errorf("工具参数应原样落盘，实际 %q", call.Args)
	}
	if call.ToolName != "write_deck" || call.ToolCallID != "call_x1" {
		t.Errorf("工具归属不对: name=%q call_id=%q", call.ToolName, call.ToolCallID)
	}

	// 全量上下文要真的存下来（这次没有截断）
	var req *Event
	for i := range evs {
		if evs[i].Kind == KindLLMRequest {
			req = &evs[i]
		}
	}
	if req == nil || !strings.Contains(string(req.Messages), "做个封面") {
		t.Errorf("llm_request 应带全量 messages，实际 %s", string(req.Messages))
	}

	// 图片：磁盘上有文件、JSON 里是路径、字节没有内联
	imgDir := filepath.Join(sessDir(root, 7), runID, "img")
	for _, name := range []string{"p001.png", "p002.png"} {
		got, err := os.ReadFile(filepath.Join(imgDir, name))
		if err != nil {
			t.Fatalf("截图 %s 应落盘: %v", name, err)
		}
		if !bytes.Equal(got, png) {
			t.Errorf("截图 %s 内容与写入不一致", name)
		}
	}
	var imgEv *Event
	for i := range evs {
		if len(evs[i].Images) > 0 {
			imgEv = &evs[i]
		}
	}
	if imgEv == nil {
		t.Fatal("没找到带图片的事件")
	}
	if len(imgEv.Images) != 2 {
		t.Fatalf("应有 2 张图，实际 %d", len(imgEv.Images))
	}
	if imgEv.Images[0].URL != runID+"/img/p001.png" {
		t.Errorf("图片路径应为 run 相对路径，实际 %q", imgEv.Images[0].URL)
	}
	if imgEv.Images[0].Label != "第 1 页" {
		t.Errorf("图片标签丢了: %q", imgEv.Images[0].Label)
	}
	raw, err := os.ReadFile(runFilePath(root, 7, runID))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(raw, png) {
		t.Error("PNG 字节被内联进了 JSONL：图片必须以文件形式引用，否则 trace 文件会膨胀到不可读")
	}
}

// TestCaptureImagesOffStillRecordsCount 关掉截图开关时，仍要留下"审了几页"的痕迹。
// 否则关掉开关的代价是"视觉审查做了什么"重新变成不可见——那正是这个工具要解决的问题。
func TestCaptureImagesOffStillRecordsCount(t *testing.T) {
	root := t.TempDir()
	runID := "1700000000001-off1"
	r := mustNew(t, Options{Dir: root, SessionID: 1, RunID: runID, UserID: 1, CaptureImages: false})

	Emit(With(context.Background(), r), Event{
		Kind:   KindSubStep,
		Images: []ImageEvent{{Name: "p001.png", Label: "第 1 页", Bytes: []byte("x")}},
	})
	r.Close()

	evs := readLines(t, runFilePath(root, 1, runID))
	if len(evs) != 1 || len(evs[0].Images) != 1 {
		t.Fatalf("应记下 1 张图的元信息，实际 %+v", evs)
	}
	if evs[0].Images[0].URL != "" {
		t.Errorf("关掉截图时 URL 应为空（页面据此显示\"未保存\"），实际 %q", evs[0].Images[0].URL)
	}
	if _, err := os.Stat(filepath.Join(sessDir(root, 1), runID, "img", "p001.png")); !os.IsNotExist(err) {
		t.Error("关掉截图时不该有图片文件落盘")
	}
}

// TestEmitFillsScopeFromCtx 归属由 ctx 决定，不由调用方传。
//
// 这条守的是"工具的兄弟调用互相串台"：工具自己不知道自己在第几轮、call_id 是什么，
// 它只调 trace.Emit(ctx, ...)。如果 WithTool 是就地改一个共享的 scope，
// 第二个工具的事件会顶着第一个工具的 call_id 落盘，排查时会把两个工具的
// 子过程拼成一条时间线——而且看不出错了。
func TestEmitFillsScopeFromCtx(t *testing.T) {
	root := t.TempDir()
	runID := "1700000000002-scop"
	r := mustNew(t, Options{Dir: root, SessionID: 1, RunID: runID, UserID: 1})

	ctx := WithTurn(With(context.Background(), r), 3)
	Emit(WithTool(ctx, "read_slide", "call_a"), Event{Kind: KindSubStep, Sub: &SubStep{Name: "x"}})
	Emit(ctx, Event{Kind: KindLLMRequest}) // 兄弟：不该继承上面那个工具
	Emit(WithTool(ctx, "web_search", "call_b"), Event{Kind: KindSubStep, Sub: &SubStep{Name: "y"}})
	// 没有 WithTurn 的 ctx：轮次应当**缺失**而不是 0。
	// 0 是合法轮次（第一次请求），混起来会让页面把 run_start 标成"第 1 轮"
	Emit(With(context.Background(), r), Event{Kind: KindRunStart})
	r.Close()

	evs := readLines(t, runFilePath(root, 1, runID))
	if len(evs) != 4 {
		t.Fatalf("应有 4 条，实际 %d", len(evs))
	}
	if turnOf(evs[0]) != 3 || evs[0].ToolName != "read_slide" || evs[0].ToolCallID != "call_a" {
		t.Errorf("第 1 条归属应为 turn=3/read_slide/call_a，实际 %d/%q/%q",
			turnOf(evs[0]), evs[0].ToolName, evs[0].ToolCallID)
	}
	if turnOf(evs[1]) != 3 {
		t.Errorf("轮次应由 ctx 补齐为 3，实际 %d", turnOf(evs[1]))
	}
	if evs[1].ToolName != "" || evs[1].ToolCallID != "" {
		t.Errorf("工具外的轮级事件不该带工具归属，实际 %q/%q", evs[1].ToolName, evs[1].ToolCallID)
	}
	if evs[2].ToolName != "web_search" || evs[2].ToolCallID != "call_b" {
		t.Errorf("第 3 条应归到 web_search/call_b，实际 %q/%q", evs[2].ToolName, evs[2].ToolCallID)
	}
	if evs[3].Turn != nil {
		t.Errorf("没有 WithTurn 的 ctx 上事件不该带轮次（0 是合法轮次，不能拿它当\"不适用\"），实际 %d", *evs[3].Turn)
	}
}

// TestNoRecorderIsANoOp 没有 recorder 时一切必须安静地什么都不做。
//
// 这是整个设计能"零侵入接入"的前提：工具函数、vision 包、以及各种单元测试
// 都会在**没有** recorder 的 ctx 上跑。这里任何一个 panic 都意味着
// 关掉 features.trace 之后 agent 直接崩——观测工具把主流程搞挂是最糟的失败方式。
func TestNoRecorderIsANoOp(t *testing.T) {
	bg := context.Background()
	Emit(bg, Event{Kind: KindToolCall, Args: "x"})
	Usage(bg, CompMain, UsagePart{Prompt: 1})

	if From(bg) != Discard {
		t.Error("没有 recorder 时 From 应返回 Discard，这样调用方不必判空")
	}
	if got := From(bg).Summary(); got.Turns != 0 || got.Total.Prompt != 0 {
		t.Errorf("Discard 的汇总应为零值，实际 %+v", got)
	}
	Discard.Emit(Event{Kind: KindToolCall})
	Discard.Close()
	Discard.Close() // 重复关闭

	// nil 接收者也要安全：调用方可能持有一个未初始化的 *Recorder
	var nilR *Recorder
	nilR.Emit(Event{Kind: KindToolCall})
	nilR.Close()
	if got := nilR.Summary(); got.Turns != 0 {
		t.Errorf("nil recorder 的汇总应为零值，实际 %+v", got)
	}

	// 关闭后继续写不该 panic（循环收尾阶段的事件晚于 Close 到达是会发生的）
	root := t.TempDir()
	r := mustNew(t, Options{Dir: root, SessionID: 1, RunID: "1700000000003-clos", UserID: 1})
	ctx := With(bg, r)
	Emit(ctx, Event{Kind: KindRunStart})
	r.Close()
	Emit(ctx, Event{Kind: KindRunEnd, Status: StatusOK})
	if evs := readLines(t, runFilePath(root, 1, "1700000000003-clos")); len(evs) != 1 {
		t.Errorf("Close 之后不该再写入，实际共 %d 条", len(evs))
	}
}

// TestUsageAccumulatesPerComponent 分项用量要各算各的，且合计等于各项之和。
//
// 这条直接对着实现前的那个缺口：done 事件里的 prompt/completion 一直只是**主循环**的数，
// 视觉审查（每次要塞好几张图，图片 token 很贵）和联网搜索（按搜索次数计费）
// 的用量拿到了就扔。用户问"这次对话花了多少 token"时，那时的答案是系统性偏小的。
func TestUsageAccumulatesPerComponent(t *testing.T) {
	root := t.TempDir()
	runID := "1700000000004-usag"
	r := mustNew(t, Options{Dir: root, SessionID: 1, RunID: runID, UserID: 1})
	ctx := With(context.Background(), r)

	Usage(ctx, CompMain, UsagePart{Prompt: 100, Completion: 20, Total: 120, Cached: 60})
	Usage(ctx, CompMain, UsagePart{Prompt: 300, Completion: 40, Total: 340, Cached: 200})
	Usage(ctx, CompWebSearch, UsagePart{Prompt: 10, Completion: 5, Total: 15})
	Usage(ctx, CompVision, UsagePart{Prompt: 5000, Completion: 300, Total: 5300, Reasoning: 120})

	got := r.Summary()
	main := got.Usage[CompMain]
	if main.Prompt != 400 || main.Completion != 60 || main.Cached != 260 || main.Calls != 2 {
		t.Errorf("主循环用量累加不对: %+v", main)
	}
	if ws := got.Usage[CompWebSearch]; ws.Total != 15 || ws.Calls != 1 {
		t.Errorf("联网搜索用量不对（Calls 未显式给时应记 1 次）: %+v", ws)
	}
	if v := got.Usage[CompVision]; v.Reasoning != 120 || v.Calls != 1 {
		t.Errorf("视觉审查用量不对: %+v", v)
	}
	wantTotal := int64(120 + 340 + 15 + 5300)
	if got.Total.Total != wantTotal {
		t.Errorf("合计应为 %d，实际 %d", wantTotal, got.Total.Total)
	}
	if got.Total.Prompt != 100+300+10+5000 {
		t.Errorf("合计 prompt 不对: %d", got.Total.Prompt)
	}

	// 每条用量都要当场落一条事件：实时看的时候数字得往上涨，而不是等这轮结束才知道
	r.Close()
	n := 0
	for _, e := range readLines(t, runFilePath(root, 1, runID)) {
		if e.Kind == KindUsage {
			n++
			if e.Component == "" || e.Usage == nil {
				t.Errorf("usage 事件缺 component 或 usage: %+v", e)
			}
		}
	}
	if n != 4 {
		t.Errorf("应写出 4 条 usage 事件，实际 %d 条", n)
	}
}

// TestMaxFieldBytesTruncatesSafely 截断开关打开时：长字段要截、messages 不能截出坏 JSON。
//
// messages 是必须特殊处理的：直接在 JSON 数组里砍一刀会留下语法坏掉的原文，
// 而读的人分不清"上下文本来就这样"和"被我截断了"——于是拿着半截上下文去推断
// 模型为什么那么答，结论必然错。所以这里断言的是"换成一个明确的说明对象"。
func TestMaxFieldBytesTruncatesSafely(t *testing.T) {
	root := t.TempDir()
	long := strings.Repeat("长", 100)

	runID := "1700000000005-trun"
	r := mustNew(t, Options{Dir: root, SessionID: 1, RunID: runID, UserID: 1, MaxFieldBytes: 32})
	ctx := With(context.Background(), r)
	Emit(ctx, Event{Kind: KindToolResult, Result: long})
	Emit(ctx, Event{Kind: KindLLMRequest, Messages: json.RawMessage(`["` + long + `"]`)})
	r.Close()

	evs := readLines(t, runFilePath(root, 1, runID))
	if !strings.Contains(evs[0].Result, "已按 max_field_bytes 截断") {
		t.Errorf("超长 Result 应带截断标记，实际 %q", evs[0].Result)
	}
	if len(evs[0].Result) > 32+64 {
		t.Errorf("截断后不该还是原来那么长: %d", len(evs[0].Result))
	}
	var note struct {
		Truncated bool `json:"truncated"`
		Bytes     int  `json:"bytes"`
	}
	if err := json.Unmarshal(evs[1].Messages, &note); err != nil {
		t.Fatalf("截断后的 messages 必须是合法 JSON，实际 %s（err=%v）", evs[1].Messages, err)
	}
	if !note.Truncated || note.Bytes == 0 {
		t.Errorf("截断后的 messages 应说明截断了以及原大小，实际 %s", evs[1].Messages)
	}
}

// TestNoTruncationByDefault 默认不截断。用户要的是"模型这一轮到底看到了什么"，
// 而这个问题的答案经常就藏在被截断的那一段里（比如整页 HTML 的尾部）。
func TestNoTruncationByDefault(t *testing.T) {
	root := t.TempDir()
	long := strings.Repeat("x", 4096)
	runID := "1700000000006-full"
	r := mustNew(t, Options{Dir: root, SessionID: 1, RunID: runID, UserID: 1})
	Emit(With(context.Background(), r), Event{Kind: KindToolResult, Result: long})
	r.Close()

	evs := readLines(t, runFilePath(root, 1, runID))
	if evs[0].Result != long {
		t.Errorf("默认不该截断：期望 %d 字节，实际 %d 字节", len(long), len(evs[0].Result))
	}
}

// TestPruneRunsKeepsNewest 裁剪留最新、删最旧。
// 写反了的表现是"刚跑的 run 不见了、旧的还在"——而这恰好会让你以为观测工具没生效。
func TestPruneRunsKeepsNewest(t *testing.T) {
	metas := []RunMeta{
		{RunID: "1700000000003-c"},
		{RunID: "1700000000001-a"},
		{RunID: "1700000000002-b"},
	}
	// 输入乱序，输出必须自己排好（列表与裁剪共用同一份排序）
	drop := pruneRuns(metas, 2)
	if len(drop) != 1 || drop[0].RunID != "1700000000001-a" {
		t.Errorf("保留 2 个时应删掉最旧的 1700000000001-a，实际 %+v", drop)
	}
	if got := pruneRuns(metas, 3); len(got) != 0 {
		t.Errorf("保留数等于总数时不该删任何东西，实际 %+v", got)
	}
	if got := pruneRuns(metas, 9); len(got) != 0 {
		t.Errorf("保留数大于总数时不该删任何东西，实际 %+v", got)
	}
	if got := pruneRuns(metas, 0); len(got) != 0 {
		t.Errorf("0 表示不限，不该删任何东西，实际 %+v", got)
	}
	// 不能改到调用方的切片：列表那边还要用
	if metas[0].RunID != "1700000000003-c" {
		t.Error("pruneRuns 不该就地重排传入的切片")
	}
}

// TestPruneOldRunsRemovesFiles 真删文件（连图片目录一起）。
// 只删 jsonl 不删图片目录的表现是磁盘一直在涨，而列表上看不出任何异常。
func TestPruneOldRunsRemovesFiles(t *testing.T) {
	root := t.TempDir()
	old := "1700000000001-old"

	rOld := mustNew(t, Options{Dir: root, SessionID: 5, RunID: old, UserID: 1, CaptureImages: true})
	Emit(With(context.Background(), rOld), Event{Kind: KindRunStart, RunID: old, SessionID: 5, UserID: 1})
	rOld.Close()
	if _, err := os.Stat(filepath.Join(sessDir(root, 5), old+".jsonl")); err != nil {
		t.Fatalf("前置条件失败，old 的 jsonl 应存在: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(sessDir(root, 5), old, "img"), 0o755); err != nil {
		t.Fatal(err)
	}

	// 建两个更新的 run，保留数=2 → 最旧的 old 应被裁掉
	for _, id := range []string{"1700000000002-mid", "1700000000003-new"} {
		rr := mustNew(t, Options{Dir: root, SessionID: 5, RunID: id, UserID: 1, RetainRuns: 2})
		Emit(With(context.Background(), rr), Event{Kind: KindRunStart, RunID: id, SessionID: 5, UserID: 1})
		rr.Close()
	}

	if _, err := os.Stat(filepath.Join(sessDir(root, 5), old+".jsonl")); !os.IsNotExist(err) {
		t.Error("超出保留数的旧 run 应被删除")
	}
	if _, err := os.Stat(filepath.Join(sessDir(root, 5), old)); !os.IsNotExist(err) {
		t.Error("旧 run 的图片目录也应被一并删除，否则磁盘会一直涨")
	}
	for _, id := range []string{"1700000000002-mid", "1700000000003-new"} {
		if _, err := os.Stat(filepath.Join(sessDir(root, 5), id+".jsonl")); err != nil {
			t.Errorf("保留范围内的 %s 不该被删: %v", id, err)
		}
	}
}

// TestNewRunIDSortableAndUnique run_id 定宽时间戳在前，字典序即时间序。
// 列表和裁剪都靠这个性质排序；长度一旦不固定（比如毫秒位数变化），
// 排序会悄悄错位，表现为"列表顺序看起来有点乱"——很难被当成 bug 报上来。
func TestNewRunIDSortableAndUnique(t *testing.T) {
	base := time.UnixMilli(1700000000000)
	a := NewRunID(base)
	b := NewRunID(base.Add(time.Millisecond))
	c := NewRunID(base.Add(time.Millisecond))

	if !(a < b) {
		t.Errorf("晚 1 毫秒的 run_id 字典序应更大: %q vs %q", a, b)
	}
	if b == c {
		t.Error("同一毫秒内的两个 run_id 应靠随机后缀区分开")
	}
	if len(a) != len(b) {
		t.Errorf("run_id 长度应固定（字典序才等于时间序）: %d vs %d", len(a), len(b))
	}
	for _, id := range []string{a, b, c} {
		if !ValidID(id) {
			t.Errorf("生成的 run_id %q 必须过自己的白名单校验", id)
		}
	}
}
