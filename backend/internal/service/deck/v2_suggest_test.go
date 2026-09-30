package deck

// 模板推荐缓存（TplSuggestion / deck.json template_suggestions 字段）的回归。
//
// 注意：包内测试库是进程级共享内存 SQLite，CreateV2Draft 只能被 TestV2FullPipeline
// 调用一次（理由见 v2_edit_test.go 文件头），这里同样手动装配 deck（固定 id）。
//
// deck 层只管缓存读写，不做阶段守卫（那是 tplsuggest 的职责）；roundtrip +
// 落盘形状 + 归属拒绝是本文件的覆盖面。

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"html-ppt/backend/internal/store"
)

// newSuggestDeck 手动装配一份 selecting_template 阶段的 v2 deck（含大纲）。
// id 由调用者给定：共享内存库里 deck 行全局唯一，一个测试二进制里每个 deck
// 只能有一个 id（TestV2FullPipeline 占 deck-0001 的同款约束）。
func newSuggestDeck(t *testing.T, id string) (*Service, uint, string) {
	t.Helper()
	s := newV2TestService(t)
	const uid = 7
	if err := s.st.DB.Create(&store.Deck{
		ID: id, UserID: uid, Format: FormatV2, Stage: StageSelectingTemplate,
	}).Error; err != nil {
		t.Fatalf("登记 deck 归属: %v", err)
	}
	dir := filepath.Join(s.decksDir, id)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("建 deck 目录: %v", err)
	}
	bundle := map[string]any{
		"id": id, "format": FormatV2, "stage": StageSelectingTemplate,
		"title": "推荐缓存测试", "canvas": map[string]int{"w": 1920, "h": 1080},
	}
	raw, _ := json.Marshal(bundle)
	if err := os.WriteFile(filepath.Join(dir, "deck.json"), raw, 0o644); err != nil {
		t.Fatalf("写 deck.json: %v", err)
	}
	outline := demoOutline()
	orb, _ := json.Marshal(outline)
	if err := os.WriteFile(filepath.Join(dir, "outline.json"), orb, 0o644); err != nil {
		t.Fatalf("写 outline.json: %v", err)
	}
	return s, uid, id
}

// TestGetHTMLNotInstantiated 未选模板的 deck（无 index.html）读页面：返回
// ErrNotInstantiated 哨兵（预览端点据此出友好占位页），文案里不再说"不存在"。
func TestGetHTMLNotInstantiated(t *testing.T) {
	s, uid, id := newSuggestDeck(t, "deck-noinst")
	if _, err := s.GetHTML(uid, id); !errors.Is(err, ErrNotInstantiated) {
		t.Fatalf("期望 ErrNotInstantiated, got %v", err)
	}
	// 越权仍按"不存在"口径，不泄露存在性
	if _, err := s.GetHTML(uid+1, id); err == nil || errors.Is(err, ErrNotInstantiated) {
		t.Fatalf("越权读不应得到未实例化哨兵, got %v", err)
	}
}

func TestTemplateSuggestionsCache(t *testing.T) {
	s, uid, id := newSuggestDeck(t, "deck-suggest")

	t.Run("无缓存返回nilnil", func(t *testing.T) {
		sugs, err := s.TemplateSuggestions(uid, id)
		if err != nil {
			t.Fatalf("读缓存: %v", err)
		}
		if sugs != nil {
			t.Fatalf("期望 nil 缓存, 得到 %v", sugs)
		}
	})

	want := []TplSuggestion{
		{TemplateID: "tech-sharing", VariantID: "blue", Reason: "科技感与主题契合"},
		{TemplateID: "weekly-report", Reason: "结构规整适合汇报"},
	}

	t.Run("写入后读回一致", func(t *testing.T) {
		if err := s.SaveTemplateSuggestions(uid, id, want); err != nil {
			t.Fatalf("写缓存: %v", err)
		}
		got, err := s.TemplateSuggestions(uid, id)
		if err != nil {
			t.Fatalf("读缓存: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("缓存 roundtrip 不一致:\n got %#v\nwant %#v", got, want)
		}
	})

	t.Run("缓存随deck.json落盘", func(t *testing.T) {
		df, err := s.GetDeckV2(uid, id)
		if err != nil {
			t.Fatalf("读 deck: %v", err)
		}
		if len(df.TemplateSuggestions) != 2 {
			t.Fatalf("DeckFile 字段未落盘: %#v", df.TemplateSuggestions)
		}
		raw, err := os.ReadFile(filepath.Join(s.decksDir, id, "deck.json"))
		if err != nil {
			t.Fatalf("读文件: %v", err)
		}
		if !strings.Contains(string(raw), "template_suggestions") {
			t.Fatal("deck.json 里没有 template_suggestions 字段")
		}
	})

	t.Run("空切片清空缓存", func(t *testing.T) {
		if err := s.SaveTemplateSuggestions(uid, id, []TplSuggestion{}); err != nil {
			t.Fatalf("清缓存: %v", err)
		}
		sugs, err := s.TemplateSuggestions(uid, id)
		if err != nil || len(sugs) != 0 {
			t.Fatalf("期望空缓存, got %v err %v", sugs, err)
		}
		// 清空后字段应从 deck.json 里消失（omitempty）
		raw, _ := os.ReadFile(filepath.Join(s.decksDir, id, "deck.json"))
		if strings.Contains(string(raw), "template_suggestions") {
			t.Fatal("清空后 deck.json 仍含 template_suggestions")
		}
	})

	t.Run("越权读不到", func(t *testing.T) {
		if err := s.SaveTemplateSuggestions(uid, id, want); err != nil {
			t.Fatalf("写缓存: %v", err)
		}
		if _, err := s.TemplateSuggestions(uid+1, id); err == nil {
			t.Fatal("别人的 deck 读缓存应当报错")
		}
	})
}
