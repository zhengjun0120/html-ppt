package usertpl

// 发布降级与质量体检（2026-09-28 门禁降级）的回归：
//   - Publish = 安全扫描 + 挂载校验，秒级完成，不再跑 Chrome 渲染量测；
//     稀疏 demo（空白骨架）从此能直接发布。
//   - Checkup 独立成只读诊断：不改状态；渲染评估抽成 evaluateRender 纯函数
//     （阈值唯一出处），Chrome 量测本身由 vision 包自己的测试与真机验收覆盖。

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"html-ppt/backend/internal/service/template"
	"html-ppt/backend/internal/store"
	"html-ppt/backend/internal/trace"
	"html-ppt/backend/internal/vision"
)

// newRegService 真实 registry（全量内置模板加载）+ 独立 Service。
func newRegService(t *testing.T) (*Service, *store.Store) {
	t.Helper()
	st, err := store.OpenMemory(t.Context())
	if err != nil {
		t.Fatalf("内存库: %v", err)
	}
	reg, err := template.NewRegistry("../../../templates", "../../../web/assets")
	if err != nil {
		t.Fatalf("加载内置模板: %v", err)
	}
	s := New(reg, st, t.TempDir(), "", "", nil, trace.Config{})
	return s, st
}

func TestPublishFastPath(t *testing.T) {
	s, st := newRegService(t)
	row, err := s.CreateBlank(9, "体检豁免的骨架")
	if err != nil {
		t.Fatalf("CreateBlank: %v", err)
	}
	// 空白骨架的 demo 是稀疏占位页——旧门禁的 45% 填充率会拦它，降级后应直接发布
	report, err := s.Publish(t.Context(), 9, row.ID)
	if err != nil {
		t.Fatalf("Publish（空白骨架）: %v", err)
	}
	if report.Structure != "ok" || report.Render != nil {
		t.Errorf("报告应为纯结构 ok（无渲染量测）: %+v", report)
	}
	var got store.UserTemplate
	if err := st.DB.First(&got, "id = ?", row.ID).Error; err != nil {
		t.Fatal(err)
	}
	if got.Status != "published" || got.Visibility != "public" {
		t.Errorf("状态未落: %+v", got)
	}
	if _, err := s.reg.Get(row.ID); err != nil {
		t.Errorf("发布后注册表查不到: %v", err)
	}
}

func TestPublishFailsOnBrokenStructure(t *testing.T) {
	s, st := newRegService(t)
	row, err := s.CreateBlank(9, "坏骨架")
	if err != nil {
		t.Fatal(err)
	}
	// 脚手架 demo 的已登记版式改成未登记的 → 挂载（内部 loadTemplate）必须失败
	raw, err := os.ReadFile(filepath.Join(s.Dir(row.ID), "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	bad := strings.Replace(string(raw), `data-layout="blank-cover"`, `data-layout="not-registered"`, 1)
	if bad == string(raw) {
		t.Fatal("脚手架 demo 里找不到 blank-cover，测试前提失效")
	}
	if err := os.WriteFile(filepath.Join(s.Dir(row.ID), "index.html"), []byte(bad), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Publish(t.Context(), 9, row.ID); err == nil {
		t.Fatal("未登记版式应发布失败")
	} else if !strings.Contains(err.Error(), "结构校验未过") {
		t.Errorf("失败原因应指向结构: %v", err)
	}
	var got store.UserTemplate
	if err := st.DB.First(&got, "id = ?", row.ID).Error; err != nil {
		t.Fatal(err)
	}
	if got.Status != "failed" || got.PublishError == "" {
		t.Errorf("应落 failed + 可读原因: %+v", got)
	}
}

func TestEvaluateRender(t *testing.T) {
	patterns := map[string]string{"cover": "hero", "quote-page": "quote", "bullets": "content"}
	mk := func(idx int, layout string, fill float64, font float64, ovX, ovY bool) vision.Slide2 {
		return vision.Slide2{Index: idx, Layout: layout, FillPct: fill, MinFontPx: font, OverflowX: ovX, OverflowY: ovY}
	}
	r := evaluateRender([]vision.Slide2{
		mk(0, "cover", 20, 15, false, false),      // hero 豁免：稀疏不算 flaw
		mk(1, "bullets", 40, 12, false, true),     // 内容页：填充率 + 字号 + 溢出三连
		mk(2, "quote-page", 10, 14, false, false), // quote 豁免：填充率不算，字号仍要达标
	}, patterns)
	if r.Pages != 3 {
		t.Errorf("Pages = %d", r.Pages)
	}
	// MinFill 只在非豁免页里取
	if r.MinFill != 40 {
		t.Errorf("MinFill 应为 40（hero/quote 豁免）: %v", r.MinFill)
	}
	if r.MaxFont != 15 {
		t.Errorf("MaxFont 应为 15: %v", r.MaxFont)
	}
	joined := strings.Join(r.Flaws, "；")
	for _, want := range []string{"第 2 页溢出画布", "第 2 页填充率仅 40%", "第 2 页最小字号 12px"} {
		if !strings.Contains(joined, want) {
			t.Errorf("flaws 缺 %q: %v", want, r.Flaws)
		}
	}
	if strings.Contains(joined, "第 1 页") || strings.Contains(joined, "第 3 页") {
		t.Errorf("hero/quote 页不应有 flaw: %v", r.Flaws)
	}
}

// TestCheckupNoGrants 无授权表时体检报错但不改状态（发布路径不再依赖 grants，
// 这个守卫防止将来有人把渲染重新接回 Publish 时漏配授权）。
func TestCheckupNoGrants(t *testing.T) {
	s, st := newRegService(t)
	row, err := s.CreateBlank(9, "无授权")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Checkup(context.Background(), 9, row.ID); err == nil {
		t.Fatal("无 grants 应报错")
	}
	var got store.UserTemplate
	if err := st.DB.First(&got, "id = ?", row.ID).Error; err != nil {
		t.Fatal(err)
	}
	if got.Status != "draft" {
		t.Errorf("体检不应改状态: %+v", got)
	}
}
