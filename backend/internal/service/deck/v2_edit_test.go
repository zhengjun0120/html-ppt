package deck

// 手动编辑保存链路（SaveHTML）的回归：docs/deck-editor-plan.md §4.1。
//
// 注意：包内测试库是进程级共享内存 SQLite（store.OpenMemory 用 cache=shared），
// 而 CreateV2Draft 的占号按各自 TempDir 从 deck-0001 重新计数——同一测试二进制里
// 再走一次 CreateV2Draft 必撞唯一约束（TestV2FullPipeline 已占用 deck-0001）。
// 所以这里完全手动建 deck：固定 id "deck-edit" + 手写 deck.json/index.html，
// 不经过 CreateV2Draft。
//
// 覆盖：正常保存落盘 + 记 OpEdit 版本（detail 缺省文案）、快照内容 = 保存后的
// index.html、NextSeq 单调递增、越权与非 v2 拒绝。

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"html-ppt/backend/internal/store"
)

const editTestHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head><meta charset="utf-8"><title>编辑器验收</title>
<link rel="stylesheet" href="style.css"></head>
<body class="tpl-tech-sharing">
<div class="deck">
<section class="slide" data-layout="cover" data-id="s1"><h1 class="h1">第一页</h1><div class="notes">notes 1</div></section>
<section class="slide" data-layout="two-column" data-id="s2"><h2 class="h2">第二页</h2><div class="notes">notes 2</div></section>
</div>
</body></html>
`

// newEditableDeck 手动装配一份可编辑的 v2 deck（不走 CreateV2Draft，理由见文件头）。
func newEditableDeck(t *testing.T) (*Service, uint, string) {
	t.Helper()
	s := newV2TestService(t)
	const uid = 7
	const id = "deck-edit"
	if err := s.st.DB.Create(&store.Deck{
		ID: id, UserID: uid, Format: FormatV2, Stage: StageIterating, TemplateID: "tech-sharing",
	}).Error; err != nil {
		t.Fatalf("登记 deck 归属: %v", err)
	}
	dir := filepath.Join(s.decksDir, id)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("建 deck 目录: %v", err)
	}
	bundle := map[string]any{
		"id": id, "format": FormatV2, "stage": StageIterating,
		"title": "编辑器验收", "template_id": "tech-sharing",
		"canvas": map[string]int{"w": 1920, "h": 1080},
	}
	raw, _ := json.Marshal(bundle)
	if err := os.WriteFile(filepath.Join(dir, "deck.json"), raw, 0o644); err != nil {
		t.Fatalf("写 deck.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte(editTestHTML), 0o644); err != nil {
		t.Fatalf("写 index.html: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "style.css"), []byte("/* stub */"), 0o644); err != nil {
		t.Fatalf("写 style.css: %v", err)
	}
	return s, uid, id
}

func TestSaveHTML(t *testing.T) {
	s, uid, id := newEditableDeck(t)

	orig, err := s.GetHTML(uid, id)
	if err != nil {
		t.Fatalf("读 deck: %v", err)
	}
	edited := strings.Replace(orig, "</body>", "<!-- hand-edited --></body>", 1)

	t.Run("保存落盘并记版本", func(t *testing.T) {
		if err := s.SaveHTML(uid, id, edited, ""); err != nil {
			t.Fatalf("保存: %v", err)
		}
		got, err := s.GetHTML(uid, id)
		if err != nil {
			t.Fatalf("读回: %v", err)
		}
		if !strings.Contains(got, "<!-- hand-edited -->") {
			t.Error("编辑内容没有落盘")
		}
		versions, err := s.ListVersions(uid, id)
		if err != nil || len(versions) == 0 {
			t.Fatalf("版本列表: %v len=%d", err, len(versions))
		}
		top := versions[0]
		if top.Operation != OpEdit {
			t.Errorf("operation = %q, 期望 %q", top.Operation, OpEdit)
		}
		if top.Detail != "手动编辑" {
			t.Errorf("detail = %q, 期望缺省文案", top.Detail)
		}
		if top.Slides != 2 {
			t.Errorf("slides = %d, 期望 2", top.Slides)
		}
		// 快照整包 = 编辑后的 index.html（bundle 是 json.Marshal 产物，< 被转义成
		// \u003c，必须 unmarshal 后比对原文，不能对原始字节做 Contains）
		raw, err := os.ReadFile(filepath.Join(s.historyDir(id), top.Version+".json"))
		if err != nil {
			t.Fatalf("读快照: %v", err)
		}
		var bundle snapshotV2
		if err := json.Unmarshal(raw, &bundle); err != nil {
			t.Fatalf("快照损坏: %v", err)
		}
		if !strings.Contains(bundle.IndexHTML, "<!-- hand-edited -->") {
			t.Error("快照里不是编辑后的内容")
		}
	})

	t.Run("二次保存版本递增", func(t *testing.T) {
		if err := s.SaveHTML(uid, id, edited+"<!-- second -->", "第二轮手改"); err != nil {
			t.Fatalf("二次保存: %v", err)
		}
		versions, _ := s.ListVersions(uid, id)
		if len(versions) < 2 {
			t.Fatalf("二次保存后版本数 = %d", len(versions))
		}
		if versions[0].Version == versions[1].Version || versions[0].Detail != "第二轮手改" {
			t.Errorf("版本号未递增或最新 detail 不对: %+v", versions[:2])
		}
	})

	t.Run("越权拒绝且与不存在同文案", func(t *testing.T) {
		err := s.SaveHTML(uid+1, id, "<html></html>", "")
		if err == nil || !strings.Contains(err.Error(), "不存在") {
			t.Errorf("越权保存应报不存在，得到: %v", err)
		}
	})

	t.Run("非 v2 拒绝", func(t *testing.T) {
		// DB 有行但 deck.json 缺失（v1/半初始化的形态）
		if err := s.st.DB.Create(&store.Deck{ID: "deck-nonv2", UserID: uid, Format: "v1"}).Error; err != nil {
			t.Fatalf("造非 v2 行: %v", err)
		}
		err := s.SaveHTML(uid, "deck-nonv2", "<html></html>", "")
		if err == nil || !strings.Contains(err.Error(), "不支持编辑") {
			t.Errorf("非 v2 应报不支持编辑，得到: %v", err)
		}
	})
}
