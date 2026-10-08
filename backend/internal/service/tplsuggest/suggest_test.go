package tplsuggest

// Suggest 全链路回归：阶段守卫 → LLM 调用（httptest 假 OpenAI 服务）→
// 解析容错（围栏/白名单外 id/重复/非法变体/空理由/畸形 JSON）→ 缓存命中与 refresh 强刷。
//
// 走真实仓库模板注册表（候选清单是真实的内置模板元数据）+ 内存 SQLite。
// 所有场景串在一个测试函数里顺序推进：调用计数与缓存内容互相依赖，拆开会失真。
// deck 手动装配、不走 CreateV2Draft（内存库进程级共享、占号按 TempDir 重计数，
// 详见 deck 包 v2_edit_test.go 文件头）。

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"html-ppt/backend/internal/service/deck"
	"html-ppt/backend/internal/service/template"
	"html-ppt/backend/internal/store"
	"html-ppt/backend/internal/trace"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func newTestRegistry(t *testing.T) *template.Registry {
	t.Helper()
	tplDir := filepath.Clean(filepath.Join("..", "..", "..", "templates"))
	assetsDir := filepath.Clean(filepath.Join("..", "..", "..", "web", "assets"))
	reg, err := template.NewRegistry(tplDir, assetsDir)
	if err != nil {
		t.Fatalf("模板注册表: %v", err)
	}
	return reg
}

func completionBody(content string) string {
	quoted, _ := json.Marshal(content)
	return fmt.Sprintf(`{"id":"cmpl-1","object":"chat.completion","created":1700000000,"model":"test",`+
		`"choices":[{"index":0,"message":{"role":"assistant","content":%s},"finish_reason":"stop"}],`+
		`"usage":{"prompt_tokens":10,"completion_tokens":5,"total_tokens":15}}`, quoted)
}

// writeDeckFiles 手写 deck.json + outline.json（Service 的目录约定：dataDir/decks/<id>）。
func writeDeckFiles(t *testing.T, dataDir, id, stage, title string) {
	t.Helper()
	dir := filepath.Join(dataDir, "decks", id)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("建 deck 目录: %v", err)
	}
	deckJSON := fmt.Sprintf(`{"id":%s,"format":"v2","stage":%s,"title":%s,"canvas":{"w":1920,"h":1080}}`,
		mustJSON(id), mustJSON(stage), mustJSON(title))
	if err := os.WriteFile(filepath.Join(dir, "deck.json"), []byte(deckJSON), 0o644); err != nil {
		t.Fatalf("写 deck.json: %v", err)
	}
	outline := fmt.Sprintf(`{"version":1,"title":%s,"meta":{"audience":"投资人","page_count":2},"pages":[`+
		`{"no":1,"role":"cover","title":%s},`+
		`{"no":2,"role":"content","title":"商业模式","points":["订阅制","规模效应"]}]}`,
		mustJSON(title), mustJSON(title))
	if err := os.WriteFile(filepath.Join(dir, "outline.json"), []byte(outline), 0o644); err != nil {
		t.Fatalf("写 outline.json: %v", err)
	}
}

func mustJSON(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func readAll(r *http.Request) ([]byte, error) {
	return io.ReadAll(r.Body)
}

func TestSuggestFlow(t *testing.T) {
	ctx := context.Background()

	// —— 假 OpenAI 服务：按调用次数返回不同形状的输出 ——
	var calls int32
	var firstReqBody atomic.Value // string，第 1 次请求的 body（带澄清诉求的那轮）
	// 去重测试的阻塞开关：置位后新到的请求挂起等 release，返回一条合法推荐
	var blockUntilRelease atomic.Bool
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		body, _ := readAll(r)
		if n == 1 {
			firstReqBody.Store(string(body))
		}
		var content string
		if blockUntilRelease.Load() {
			<-release
			content = `[{"template_id":"tech-sharing","reason":"去重验证"}]`
		} else {
			switch n {
			case 1:
				// 围栏包裹的合法 JSON；混入：白名单外 id、重复 id、非法变体、空理由
				content = "```json\n" + `[
					{"template_id":"tech-sharing","variant_id":"blue","reason":"科技感与工程主题契合"},
					{"template_id":"ut-private-tpl","reason":"不在候选里"},
					{"template_id":"tech-sharing","reason":"重复条目"},
					{"template_id":"weekly-report","variant_id":"no-such-variant","reason":"结构适合周报"},
					{"template_id":"minimal-white","reason":""}
				]` + "\n```"
			case 2:
				content = `[{"template_id":"minimal-white","reason":"极简留白，商务干净"}]`
			case 4:
				content = "" // 空回复（推理耗尽输出的形状）→ 触发重试
			case 5:
				content = `[{"template_id":"minimal-white","reason":"重试后给出的推荐"}]`
			default:
				content = "这不是 JSON，模型抽风了"
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(completionBody(content)))
	}))
	defer srv.Close()

	client := openai.NewClient(option.WithAPIKey("test-key"), option.WithBaseURL(srv.URL))
	llm := LLM{Client: &client, Model: "test-model"}

	// —— deck 装配 ——
	st, err := store.OpenMemory(ctx)
	if err != nil {
		t.Fatalf("内存库: %v", err)
	}
	const uid = uint(7)
	const id = "deck-suggest-flow"
	const guardID = "deck-suggest-guard"
	const dedupID = "deck-suggest-dedup"
	dataDir := t.TempDir()
	reg := newTestRegistry(t)
	svc := deck.New(dataDir, t.TempDir(), st).WithTemplateRegistry(reg)
	// 观测落盘目录：断言推荐 run 会写 trace（run_kind=tplsugg）
	traceDir := filepath.Join(t.TempDir(), "traces")
	traceCfg := trace.Config{Enabled: true, Dir: traceDir, MaxFieldBytes: 1 << 16, RetainRuns: 5}

	// traceRuns 观测断言辅助：该 deck 名下（伪会话段）的推荐 run 数
	traceRuns := func(want int, label string) {
		t.Helper()
		tstore := trace.NewStore(traceDir)
		runs, _, err := tstore.ListRuns(uid, trace.ListFilter{SessionID: suggestSessionID(id)}, 10, 0)
		if err != nil {
			t.Fatalf("%s: ListRuns: %v", label, err)
		}
		if len(runs) != want {
			t.Fatalf("%s: 期望 %d 个推荐 run, got %d", label, want, len(runs))
		}
	}
	for _, row := range []struct {
		id    string
		title string
		stage string
	}{{id, "AI 模板推荐", deck.StageSelectingTemplate}, {guardID, "还没确认大纲", deck.StageOutlineReview}, {dedupID, "去重测试", deck.StageSelectingTemplate}} {
		if err := st.DB.Create(&store.Deck{ID: row.id, UserID: uid, Title: row.title, Format: "v2", Stage: row.stage}).Error; err != nil {
			t.Fatalf("登记 deck %s: %v", row.id, err)
		}
		writeDeckFiles(t, dataDir, row.id, row.stage, row.title)
	}

	sg := New(svc, reg, traceCfg)

	t.Run("阶段守卫", func(t *testing.T) {
		_, err := sg.Suggest(ctx, uid, guardID, llm, "", false)
		var mismatch deck.StageMismatch
		if !errors.As(err, &mismatch) {
			t.Fatalf("期望 StageMismatch, got %v", err)
		}
		if mismatch.Current != deck.StageOutlineReview {
			t.Fatalf("Current = %s", mismatch.Current)
		}
		if n := atomic.LoadInt32(&calls); n != 0 {
			t.Fatalf("守卫拒绝不应触达 LLM, calls=%d", n)
		}
	})

	// 第一轮：LLM 返回 5 条，清洗后应剩 2 条合法推荐
	var first []deck.TplSuggestion
	t.Run("首调清洗", func(t *testing.T) {
		got, err := sg.Suggest(ctx, uid, id, llm, "给投资人讲商业模式，最好商务一点", false)
		if err != nil {
			t.Fatalf("Suggest: %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("期望 2 条合法推荐, got %d: %#v", len(got), got)
		}
		if got[0].TemplateID != "tech-sharing" || got[0].VariantID != "blue" {
			t.Fatalf("第 1 条不符: %#v", got[0])
		}
		if got[1].TemplateID != "weekly-report" || got[1].VariantID != "" {
			t.Fatalf("第 2 条应保留模板、丢掉非法变体: %#v", got[1])
		}
		first = got
		// 观测：真调 LLM 的一轮落一个 run，run_kind=tplsugg、user_content 带标题页数
		traceRuns(1, "首调后")
		tstore := trace.NewStore(traceDir)
		runs, _, _ := tstore.ListRuns(uid, trace.ListFilter{SessionID: suggestSessionID(id)}, 10, 0)
		if runs[0].RunKind != "tplsugg" {
			t.Errorf("run_kind = %s, 期望 tplsugg", runs[0].RunKind)
		}
		if !strings.Contains(runs[0].UserContent, "AI 模板推荐") || !strings.Contains(runs[0].UserContent, "2 页") {
			t.Errorf("user_content = %q", runs[0].UserContent)
		}
		if runs[0].Status != "ok" {
			t.Errorf("status = %s, 期望 ok", runs[0].Status)
		}
		// run_end 带 recorder 汇总：观测台列表 tokens 列的数据源
		if runs[0].Usage == nil || runs[0].Usage.Total.Total == 0 {
			t.Errorf("run 汇总缺用量: %+v", runs[0].Usage)
		}
	})

	t.Run("缓存命中不重调", func(t *testing.T) {
		got, err := sg.Suggest(ctx, uid, id, llm, "", false)
		if err != nil {
			t.Fatalf("Suggest: %v", err)
		}
		if len(got) != len(first) {
			t.Fatalf("缓存内容不符: %#v", got)
		}
		if n := atomic.LoadInt32(&calls); n != 1 {
			t.Fatalf("缓存命中不应重调 LLM, calls=%d", n)
		}
		traceRuns(1, "缓存命中后") // 不新增 run
	})

	t.Run("refresh强刷覆盖缓存", func(t *testing.T) {
		got, err := sg.Suggest(ctx, uid, id, llm, "", true)
		if err != nil {
			t.Fatalf("Suggest: %v", err)
		}
		if len(got) != 1 || got[0].TemplateID != "minimal-white" {
			t.Fatalf("强刷结果不符: %#v", got)
		}
		if n := atomic.LoadInt32(&calls); n != 2 {
			t.Fatalf("refresh 应重调 LLM, calls=%d", n)
		}
		traceRuns(2, "refresh 后")
	})

	t.Run("畸形输出返回空且不写缓存", func(t *testing.T) {
		got, err := sg.Suggest(ctx, uid, id, llm, "", true)
		if err != nil {
			t.Fatalf("畸形输出不应报错, got %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("期望空结果, got %#v", got)
		}
		if n := atomic.LoadInt32(&calls); n != 3 {
			t.Fatalf("calls=%d", n)
		}
		traceRuns(3, "畸形输出后") // 失败的 LLM 轮也要留痕（run_end 带 error）
		{
			tstore := trace.NewStore(traceDir)
			runs, _, _ := tstore.ListRuns(uid, trace.ListFilter{SessionID: suggestSessionID(id)}, 10, 0)
			if runs[0].Status != "error" || runs[0].Usage == nil || runs[0].Usage.Total.Total == 0 {
				t.Errorf("失败 run 应带 error 状态与用量: status=%s usage=%+v", runs[0].Status, runs[0].Usage)
			}
		}
		// 空结果没覆盖缓存：再读（不 refresh）应拿到上一轮的 minimal-white
		again, err := sg.Suggest(ctx, uid, id, llm, "", false)
		if err != nil || len(again) != 1 || again[0].TemplateID != "minimal-white" {
			t.Fatalf("缓存被空结果污染: %#v err %v", again, err)
		}
		if n := atomic.LoadInt32(&calls); n != 3 {
			t.Fatalf("读缓存不应重调, calls=%d", n)
		}
		traceRuns(3, "再读缓存后")
	})

	t.Run("空回复重试一次后成功", func(t *testing.T) {
		got, err := sg.Suggest(ctx, uid, id, llm, "", true)
		if err != nil {
			t.Fatalf("Suggest: %v", err)
		}
		if len(got) != 1 || got[0].TemplateID != "minimal-white" || got[0].Reason != "重试后给出的推荐" {
			t.Fatalf("重试结果不符: %#v", got)
		}
		if n := atomic.LoadInt32(&calls); n != 5 {
			t.Fatalf("空回复应恰好重试一次, calls=%d", n)
		}
		traceRuns(4, "重试成功后") // 空回复轮+重试轮同属一个 run
	})

	t.Run("同deck在飞去重", func(t *testing.T) {
		// 模拟确认后的预热（A）与用户进页请求（B）同时到达：A 先拿到锁在 LLM
		// 上阻塞，B 排队；放行后 A 写缓存，B 拿到锁直接命中——server 只被调一次
		blockUntilRelease.Store(true)
		type res struct {
			sugs []deck.TplSuggestion
			err  error
		}
		doneA := make(chan res, 1)
		go func() {
			sugs, err := sg.Suggest(ctx, uid, dedupID, llm, "", false)
			doneA <- res{sugs, err}
		}()
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) && atomic.LoadInt32(&calls) != 6 {
			time.Sleep(30 * time.Millisecond)
		}
		if atomic.LoadInt32(&calls) != 6 {
			t.Fatal("A 未进入 LLM 调用（阻塞点没到）")
		}
		doneB := make(chan res, 1)
		go func() {
			sugs, err := sg.Suggest(ctx, uid, dedupID, llm, "", false)
			doneB <- res{sugs, err}
		}()
		time.Sleep(300 * time.Millisecond) // 给 B 时间排上锁队
		close(release)
		var rA, rB res
		select {
		case rA = <-doneA:
		case <-time.After(5 * time.Second):
			t.Fatal("A 未返回")
		}
		select {
		case rB = <-doneB:
		case <-time.After(5 * time.Second):
			t.Fatal("B 未返回")
		}
		if rA.err != nil || rB.err != nil {
			t.Fatalf("去重双调报错: A=%v B=%v", rA.err, rB.err)
		}
		if n := atomic.LoadInt32(&calls); n != 6 {
			t.Fatalf("去重失败：B 又调了一次 LLM, calls=%d", n)
		}
		if len(rA.sugs) != 1 || rA.sugs[0].Reason != "去重验证" {
			t.Fatalf("A 结果不符: %#v", rA.sugs)
		}
		if len(rB.sugs) != 1 || rB.sugs[0].Reason != "去重验证" {
			t.Fatalf("B 应命中 A 写的缓存: %#v", rB.sugs)
		}
		// A 的 run 落在 dedupID 自己的伪会话目录（每 deck 一段），带 ok 状态与用量
		dstore := trace.NewStore(traceDir)
		druns, _, err := dstore.ListRuns(uid, trace.ListFilter{SessionID: suggestSessionID(dedupID)}, 10, 0)
		if err != nil || len(druns) != 1 {
			t.Fatalf("去重 deck 应有 1 个 run: %v %d", err, len(druns))
		}
		if druns[0].Status != "ok" || druns[0].Usage == nil || druns[0].Usage.Total.Total == 0 {
			t.Errorf("去重 run 汇总不符: status=%s usage=%+v", druns[0].Status, druns[0].Usage)
		}
		// B 命中缓存不算 run：主 deck 的 run 数不变
		traceRuns(4, "去重后")
	})

	t.Run("越权拒绝", func(t *testing.T) {
		if _, err := sg.Suggest(ctx, uid+1, id, llm, "", false); err == nil {
			t.Fatal("别人的 deck 应当报错")
		}
	})

	t.Run("请求prompt组装", func(t *testing.T) {
		body, _ := firstReqBody.Load().(string)
		for _, want := range []string{
			"tech-sharing", // 候选清单里有真实模板
			"给投资人讲商业模式",    // 澄清诉求进 prompt
			"商业模式",         // 大纲页标题进 prompt
			"template_id",  // system prompt 里的输出格式契约
			"周报",           // 候选清单带中文名（weekly-report）
		} {
			if !strings.Contains(body, want) {
				t.Errorf("prompt 缺 %q", want)
			}
		}
	})
}
