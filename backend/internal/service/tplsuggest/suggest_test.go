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

	"html-ppt/backend/internal/service/deck"
	"html-ppt/backend/internal/service/template"
	"html-ppt/backend/internal/store"

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
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		body, _ := readAll(r)
		if n == 1 {
			firstReqBody.Store(string(body))
		}
		var content string
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
		default:
			content = "这不是 JSON，模型抽风了"
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
	dataDir := t.TempDir()
	svc := deck.New(dataDir, t.TempDir(), st).WithTemplateRegistry(newTestRegistry(t))
	for _, row := range []struct {
		id    string
		title string
		stage string
	}{{id, "AI 模板推荐", deck.StageSelectingTemplate}, {guardID, "还没确认大纲", deck.StageOutlineReview}} {
		if err := st.DB.Create(&store.Deck{ID: row.id, UserID: uid, Title: row.title, Format: "v2", Stage: row.stage}).Error; err != nil {
			t.Fatalf("登记 deck %s: %v", row.id, err)
		}
		writeDeckFiles(t, dataDir, row.id, row.stage, row.title)
	}

	sg := New(svc, newTestRegistry(t))

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
		// 空结果没覆盖缓存：再读（不 refresh）应拿到上一轮的 minimal-white
		again, err := sg.Suggest(ctx, uid, id, llm, "", false)
		if err != nil || len(again) != 1 || again[0].TemplateID != "minimal-white" {
			t.Fatalf("缓存被空结果污染: %#v err %v", again, err)
		}
		if n := atomic.LoadInt32(&calls); n != 3 {
			t.Fatalf("读缓存不应重调, calls=%d", n)
		}
	})

	t.Run("越权拒绝", func(t *testing.T) {
		if _, err := sg.Suggest(ctx, uid+1, id, llm, "", false); err == nil {
			t.Fatal("别人的 deck 应当报错")
		}
	})

	t.Run("请求prompt组装", func(t *testing.T) {
		body, _ := firstReqBody.Load().(string)
		for _, want := range []string{
			"tech-sharing",       // 候选清单里有真实模板
			"给投资人讲商业模式",         // 澄清诉求进 prompt
			"商业模式",               // 大纲页标题进 prompt
			"template_id",        // system prompt 里的输出格式契约
			"周报",                 // 候选清单带中文名（weekly-report）
		} {
			if !strings.Contains(body, want) {
				t.Errorf("prompt 缺 %q", want)
			}
		}
	})
}
