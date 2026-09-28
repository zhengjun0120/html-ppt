package usertpl

// CreateBlank（从空白新建）的回归：docs/user-template-history-plan.md 的姊妹需求。
// 关键断言用**真实 registry**（全量内置模板加载）：空白脚手架必须整体通过
// loadTemplate 的全部校验（meta/layouts 双向登记/骨架自洽/资产引用/挂载标记），
// 否则 MountUser 失败、模板建不出来——这是脚手架内容唯一的权威裁判。

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"html-ppt/backend/internal/service/template"
	"html-ppt/backend/internal/store"
	"html-ppt/backend/internal/trace"
)

func TestCreateBlank(t *testing.T) {
	st, err := store.OpenMemory(t.Context())
	if err != nil {
		t.Fatalf("内存库: %v", err)
	}
	// 测试 CWD = 包目录；仓库内 templates/ 与 web/assets 在 ../../../
	reg, err := template.NewRegistry("../../../templates", "../../../web/assets")
	if err != nil {
		t.Fatalf("加载内置模板: %v", err)
	}
	root := t.TempDir()
	s := New(reg, st, root, "", "", nil, trace.Config{})

	row, err := s.CreateBlank(9, "我的空白")
	if err != nil {
		t.Fatalf("CreateBlank: %v", err)
	}
	if row.Status != "draft" || row.Visibility != "private" || row.BaseID != "" {
		t.Errorf("行状态不对: %+v", row)
	}
	if row.Name != "我的空白" {
		t.Errorf("名称 = %q", row.Name)
	}
	// template.json 的 id/name 与目录/行一致（loadTemplate 的硬规则）
	raw, err := os.ReadFile(filepath.Join(s.Dir(row.ID), "template.json"))
	if err != nil {
		t.Fatalf("读 template.json: %v", err)
	}
	if !strings.Contains(string(raw), fmt.Sprintf(`"id": %q`, row.ID)) {
		t.Errorf("template.json 未替换 id: %s", raw)
	}
	// 五件套 + ADAPTATION.md 都在
	for _, f := range append(scaffoldFiles, "ADAPTATION.md") {
		if _, err := os.Stat(filepath.Join(s.Dir(row.ID), f)); err != nil {
			t.Errorf("缺 %s", f)
		}
	}
	// 注册表里挂上了（生成管线立即可用）
	if _, err := reg.Get(row.ID); err != nil {
		t.Errorf("注册表查不到新建的空白模板: %v", err)
	}
	// 起点版本
	versions, err := s.ListUTVersions(9, row.ID)
	if err != nil || len(versions) != 1 || versions[0].Operation != OpFork || versions[0].Detail != "从空白新建" {
		t.Errorf("起点版本不对: %+v err=%v", versions, err)
	}
	// 空名兜底
	row2, err := s.CreateBlank(9, "")
	if err != nil {
		t.Fatalf("CreateBlank 空名: %v", err)
	}
	if row2.Name != "空白模板" {
		t.Errorf("空名应兜底「空白模板」，得到 %q", row2.Name)
	}
}
