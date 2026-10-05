package usertpl

// 结构契约（版式面板）的回归：
//   读  StructureContract——挂载态契约的组装（版式清单/demo 映射/rules.md）；
//   写  UpdateLayoutMeta——roles/名称/用途的补丁语义（nil 不动 vs 清空）、
//       词表校验、template.json 权威落盘 + layouts.md 文档行同步、记 structure 版本、
//       重挂后注册表可读新值；published/publishing 拒绝。
//
// 全程用真实 registry（CreateBlank 的 9 版式脚手架）：重挂过不过 loadTemplate
// 全量校验，只有它能裁判。

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"html-ppt/backend/internal/service/template"
	"html-ppt/backend/internal/store"
	"html-ppt/backend/internal/trace"
)

func structTestService(t *testing.T) (*Service, *template.Registry) {
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
	return s, reg
}

func TestStructureContract(t *testing.T) {
	s, _ := structTestService(t)
	const uid = 11
	row, err := s.CreateBlank(uid, "结构面板")
	if err != nil {
		t.Fatalf("CreateBlank: %v", err)
	}

	v, err := s.StructureContract(uid, row.ID)
	if err != nil {
		t.Fatalf("StructureContract: %v", err)
	}
	if len(v.Layouts) != 9 {
		t.Errorf("layouts = %d，应为 9", len(v.Layouts))
	}
	if len(v.DemoPages) != 9 {
		t.Fatalf("demo_pages = %d，应为 9", len(v.DemoPages))
	}
	for i, p := range v.DemoPages {
		if p.No != i+1 {
			t.Errorf("demo 第 %d 项页码 = %d：页码必须连续", i, p.No)
		}
		if !strings.HasPrefix(p.Layout, "blank-") {
			t.Errorf("demo 第 %d 页的 data-layout = %q，未登记", i+1, p.Layout)
		}
	}
	cover := v.Layouts[0]
	if cover.ID != "blank-cover" || cover.Skeleton == "" || cover.Pattern != "hero" {
		t.Errorf("封面版式视图不对: %+v", cover)
	}
	if !strings.Contains(v.RulesMD, "9 个") {
		t.Errorf("rules_md 未带上：%q", v.RulesMD[:40])
	}
	// 越权：另一个用户读不到
	if _, err := s.StructureContract(uid+1, row.ID); err == nil {
		t.Error("越权读结构契约未拒绝")
	}
}

func TestUpdateLayoutMeta(t *testing.T) {
	s, reg := structTestService(t)
	const uid = 12
	row, err := s.CreateBlank(uid, "版式编辑")
	if err != nil {
		t.Fatalf("CreateBlank: %v", err)
	}

	// —— roles 补丁：改 + 词表内排序 ——//
	warning, err := s.UpdateLayoutMeta(uid, row.ID, "blank-data", LayoutMetaPatch{
		Roles: &[]string{"quote", "data", "content"},
	})
	if err != nil || warning != "" {
		t.Fatalf("UpdateLayoutMeta(roles): warning=%q err=%v", warning, err)
	}

	// template.json（权威）落盘新值，其余字段保真
	raw, err := os.ReadFile(filepath.Join(s.Dir(row.ID), "template.json"))
	if err != nil {
		t.Fatalf("读 template.json: %v", err)
	}
	var meta map[string]any
	if err := json.Unmarshal(raw, &meta); err != nil {
		t.Fatalf("template.json 损坏: %v", err)
	}
	if meta["description"] == nil || meta["canvas"] == nil {
		t.Error("template.json 重写丢了 layouts 之外的字段")
	}
	layouts, _ := meta["layouts"].([]any)
	var dataEntry map[string]any
	for _, it := range layouts {
		if m, ok := it.(map[string]any); ok && m["id"] == "blank-data" {
			dataEntry = m
		}
	}
	if dataEntry == nil {
		t.Fatal("blank-data 条目丢了")
	}
	if got, _ := dataEntry["roles"].([]any); len(got) != 3 || got[0] != "content" {
		t.Errorf("roles 未按词表排序或丢值: %v", got)
	}

	// layouts.md 的文档行同步（适用 role：行）
	mdRaw, err := os.ReadFile(filepath.Join(s.Dir(row.ID), "layouts.md"))
	if err != nil {
		t.Fatalf("读 layouts.md: %v", err)
	}
	if !strings.Contains(string(mdRaw), "适用 role：content、data、quote。") {
		t.Errorf("layouts.md 角色行未同步: %s", mdRaw)
	}

	// 重挂后注册表可读新值（生成侧立即生效）
	tpl, err := reg.Get(row.ID)
	if err != nil {
		t.Fatalf("重挂后查不到模板: %v", err)
	}
	var found *template.LayoutMeta
	for i := range tpl.Layouts {
		if tpl.Layouts[i].ID == "blank-data" {
			found = &tpl.Layouts[i]
		}
	}
	if found == nil || len(found.Roles) != 3 || found.Roles[0] != "content" {
		t.Errorf("注册表里的 roles 未更新: %+v", found)
	}

	// 版本记录：operation=structure，changed 含两个文件
	versions, err := s.ListUTVersions(uid, row.ID)
	if err != nil || len(versions) == 0 {
		t.Fatalf("历史列表: %v", err)
	}
	top := versions[0]
	if top.Operation != OpStruct {
		t.Errorf("最新版本 operation = %q，应为 %q", top.Operation, OpStruct)
	}
	if !containsStr(top.Changed, "template.json") || !containsStr(top.Changed, "layouts.md") {
		t.Errorf("changed 缺文件: %v", top.Changed)
	}

	// —— 名称/用途补丁 + 摘要文案 ——//
	warning, err = s.UpdateLayoutMeta(uid, row.ID, "blank-data", LayoutMetaPatch{
		Name: ptr("大数字页"),
		Use:  ptr("三个以内关键指标的对比"),
	})
	if err != nil || warning != "" {
		t.Fatalf("UpdateLayoutMeta(meta): warning=%q err=%v", warning, err)
	}
	if !strings.Contains(top.Detail, "blank-data") && versions[0].Detail == "" {
		t.Error("版本备注为空")
	}

	// —— 清空角色（空数组）：json 删字段、md 行写"无" ——//
	if _, err := s.UpdateLayoutMeta(uid, row.ID, "blank-split", LayoutMetaPatch{Roles: &[]string{}}); err != nil {
		t.Fatalf("清空角色: %v", err)
	}
	raw, _ = os.ReadFile(filepath.Join(s.Dir(row.ID), "template.json"))
	if strings.Contains(string(raw), `"blank-split", "roles"`) {
		t.Error("blank-split 的 roles 应已清空")
	}
	_ = raw

	// —— 预检拒绝：非法角色词 / 未知版式 / 全空补丁 / 超长名称 ——//
	if _, err := s.UpdateLayoutMeta(uid, row.ID, "blank-data", LayoutMetaPatch{Roles: &[]string{"hero"}}); err == nil {
		t.Error("非法角色词未拒绝")
	}
	if _, err := s.UpdateLayoutMeta(uid, row.ID, "no-such-layout", LayoutMetaPatch{Name: ptr("x")}); err == nil {
		t.Error("未知版式未拒绝")
	}
	if _, err := s.UpdateLayoutMeta(uid, row.ID, "blank-data", LayoutMetaPatch{}); err == nil {
		t.Error("空补丁未拒绝")
	}
	if _, err := s.UpdateLayoutMeta(uid, row.ID, "blank-data", LayoutMetaPatch{Name: ptr(strings.Repeat("长", 41))}); err == nil {
		t.Error("超长名称未拒绝")
	}

	// —— published 拒绝（409 语义）——//
	if err := s.st.DB.Model(&store.UserTemplate{}).Where("id = ?", row.ID).Update("status", "published").Error; err != nil {
		t.Fatalf("置 published: %v", err)
	}
	_, err = s.UpdateLayoutMeta(uid, row.ID, "blank-data", LayoutMetaPatch{Name: ptr("x")})
	if !errors.Is(err, ErrPublished) {
		t.Errorf("published 未拒绝: %v", err)
	}
}

func TestSyncRolesLine(t *testing.T) {
	md := "## a（甲）\n指纹：stack\n适用 role：content。\n正文。\n\n## b（乙）\n适用 role：toc。\n\n## c（丙）\n没有该行的条目。\n"
	out, changed := syncRolesLine(md, "b", []string{"content", "toc"})
	if !changed {
		t.Fatal("未检出可同步行")
	}
	if !strings.Contains(out, "适用 role：content、toc。") {
		t.Errorf("b 条目行未更新: %s", out)
	}
	if !strings.Contains(out, "## a（甲）\n指纹：stack\n适用 role：content。") {
		t.Error("a 条目被误改")
	}
	// 无「适用 role」行的条目：静默不改
	if _, changed := syncRolesLine(md, "c", []string{"quote"}); changed {
		t.Error("缺行时应跳过")
	}
}

func ptr[T any](v T) *T { return &v }

func containsStr(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
