package usertpl

// 历史模块（history.go）的回归：docs/user-template-history-plan.md §3.1/§6。
// 覆盖：记版与 changed、幂等、回滚（写回+元数据同步+restore 记版）、published 锁、
// 版本号/缺失校验、prune（起点永留 + 旧日抽稀 + 上限截尾）、UpdateMeta 语义、
// execCustomTool 的 dirty 标志、删除/清空。

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"html-ppt/backend/internal/store"
)

// newHistService 每个用例独立的 Service + 模板行 + 五件套目录。
// 共享内存库 → id 必须全局唯一（后缀区分用例）。
func newHistService(t *testing.T, suffix string) (*Service, *store.Store, string, uint) {
	t.Helper()
	st, err := store.OpenMemory(t.Context())
	if err != nil {
		t.Fatalf("内存库: %v", err)
	}
	root := t.TempDir()
	s := New(nil, st, root, "", "", nil)
	id := fmt.Sprintf("ut-hist-%s-%d", suffix, len(root))
	const uid = 9
	if err := st.DB.Create(&store.UserTemplate{
		ID: id, UserID: uid, BaseID: "tech-sharing", Status: "draft",
		Name: "基线名", Description: "基线描述",
	}).Error; err != nil {
		t.Fatalf("建行: %v", err)
	}
	if err := os.MkdirAll(s.Dir(id), 0o755); err != nil {
		t.Fatalf("建目录: %v", err)
	}
	files := map[string]string{
		"template.json": fmt.Sprintf(`{"id":%q,"name":"基线名","description":"基线描述"}`, id),
		"index.html":    "<section>base</section>",
		"style.css":     ".tpl-x{--accent:#000}",
		"layouts.md":    "## cover\n骨架",
		"rules.md":      "# 规则",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(s.Dir(id), name), []byte(content), 0o644); err != nil {
			t.Fatalf("写 %s: %v", name, err)
		}
	}
	return s, st, id, uid
}

func TestRecordAndChanged(t *testing.T) {
	s, _, id, uid := newHistService(t, "rec")
	_ = uid
	// 首版：changed = 全部非空文件（基线）
	if err := s.recordVersionUT(id, OpFork, "克隆自 builtin:tech-sharing"); err != nil {
		t.Fatalf("记首版: %v", err)
	}
	versions, err := s.ListUTVersions(9, id)
	if err != nil || len(versions) != 1 {
		t.Fatalf("应 1 版: %d err=%v", len(versions), err)
	}
	if versions[0].Version != "v000001" || versions[0].Operation != OpFork {
		t.Errorf("首版元信息: %+v", versions[0])
	}
	if len(versions[0].Changed) != len(utFiles) {
		t.Errorf("首版 changed 应为全部 %d 个文件，得到 %v", len(utFiles), versions[0].Changed)
	}
	// 第二版只改 style.css
	if err := os.WriteFile(filepath.Join(s.Dir(id), "style.css"), []byte(".tpl-x{--accent:#fff}"), 0o644); err != nil {
		t.Fatal(err)
	}
	_ = s.recordVersionUT(id, OpChat, "改主色")
	versions, _ = s.ListUTVersions(9, id)
	if len(versions) != 2 {
		t.Fatalf("应 2 版: %d", len(versions))
	}
	if len(versions[0].Changed) != 1 || versions[0].Changed[0] != "style.css" {
		t.Errorf("changed 应为 [style.css]: %v", versions[0].Changed)
	}
	// 无变化不记版
	before, _ := s.ListUTVersions(9, id)
	_ = s.recordVersionUT(id, OpChat, "没改东西")
	after, _ := s.ListUTVersions(9, id)
	if len(before) != len(after) {
		t.Errorf("相同 bundle 不应记版：%d → %d", len(before), len(after))
	}
}

func TestRestoreUTVersion(t *testing.T) {
	s, st, id, uid := newHistService(t, "restore")
	_ = s.recordVersionUT(id, OpFork, "起点")
	// 改 style.css 与名称描述 → 记第二版
	if err := os.WriteFile(filepath.Join(s.Dir(id), "style.css"), []byte(".tpl-x{--accent:#fff}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := s.UpdateMeta(uid, id, "新名字", "新描述"); err != nil {
		t.Fatalf("UpdateMeta: %v", err)
	}
	versions, _ := s.ListUTVersions(uid, id)
	if len(versions) != 2 {
		t.Fatalf("应 2 版（chat/meta 各一）: %d", len(versions))
	}

	// 回滚到起点
	if err := s.RestoreUTVersion(uid, id, "v000001"); err != nil {
		t.Fatalf("回滚: %v", err)
	}
	data, _ := os.ReadFile(filepath.Join(s.Dir(id), "style.css"))
	if string(data) != ".tpl-x{--accent:#000}" {
		t.Errorf("style.css 未写回: %q", data)
	}
	var row store.UserTemplate
	if err := st.DB.First(&row, "id = ?", id).Error; err != nil {
		t.Fatal(err)
	}
	if row.Name != "基线名" || row.Description != "基线描述" {
		t.Errorf("元数据未按快照同步: %q / %q", row.Name, row.Description)
	}
	versions, _ = s.ListUTVersions(uid, id)
	if len(versions) != 3 || versions[0].Operation != OpRestore || versions[0].Detail != "恢复到 v000001" {
		t.Fatalf("回滚应记 restore 版: %+v", versions)
	}
	// 回滚后 status 不变（回滚不改流程状态）
	if row.Status != "draft" {
		t.Errorf("status 被改: %q", row.Status)
	}
}

func TestRestoreGuards(t *testing.T) {
	s, st, id, uid := newHistService(t, "guard")
	_ = s.recordVersionUT(id, OpFork, "起点")

	if err := s.RestoreUTVersion(uid, id, "v999999"); err == nil {
		t.Error("不存在的版本应报错")
	}
	if err := s.RestoreUTVersion(uid, id, "../etc/passwd"); err == nil {
		t.Error("非法版本号应报错")
	}
	if err := s.RestoreUTVersion(uid+1, id, "v000001"); err == nil {
		t.Error("越权应报错")
	}
	// published 锁（D2）
	if err := st.DB.Model(&store.UserTemplate{}).Where("id = ?", id).Update("status", "published").Error; err != nil {
		t.Fatal(err)
	}
	err := s.RestoreUTVersion(uid, id, "v000001")
	if !errors.Is(err, ErrPublished) {
		t.Errorf("published 回滚应报 ErrPublished: %v", err)
	}
}

func TestPruneUTVersions(t *testing.T) {
	s, _, id, _ := newHistService(t, "prune")
	old := time.Now().AddDate(0, 0, -9).Unix()
	// 手工布一个索引：同一天（9 天前）的 3 版 + 起点 v000001
	mk := func(v string, ts int64) UTVersionMeta {
		return UTVersionMeta{Version: v, Time: ts, Operation: OpEdit, Detail: "x"}
	}
	idx := utHistoryIndex{NextSeq: 4, Versions: []UTVersionMeta{
		mk("v000003", old), mk("v000002", old), mk("v000001", old),
	}}
	if err := os.MkdirAll(s.historyDir(id), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, v := range idx.Versions {
		b := snapshotUT{IndexHTML: v.Version}
		data, _ := json.Marshal(&b)
		if err := os.WriteFile(filepath.Join(s.historyDir(id), v.Version+".json"), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.writeUTHistoryIndex(id, idx); err != nil {
		t.Fatal(err)
	}
	// 现在记一版 → 触发 prune：9 天前同日只留 v000003，v000001（起点）强制保留
	if err := s.recordVersionUT(id, OpEdit, "new"); err != nil {
		t.Fatal(err)
	}
	versions, _ := s.ListUTVersions(9, id)
	var got []string
	for _, m := range versions {
		got = append(got, m.Version)
	}
	want := []string{"v000004", "v000003", "v000001"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("裁剪结果 = %v, want %v", got, want)
	}
	if _, err := os.Stat(filepath.Join(s.historyDir(id), "v000002.json")); !os.IsNotExist(err) {
		t.Errorf("v000002 快照应被删除")
	}
}

func TestPruneUTVersionsCap(t *testing.T) {
	s, _, id, _ := newHistService(t, "cap")
	// 105 个"刚刚"的版本（全 recent）+ 记一版 → 上限 100 截尾但起点强制保留
	// （index 约定新→旧：v000105 最新在前、v000001 最老在末）
	now := time.Now().Unix()
	idx := utHistoryIndex{NextSeq: 106}
	for i := 105; i >= 1; i-- {
		idx.Versions = append(idx.Versions, UTVersionMeta{
			Version: fmt.Sprintf("v%06d", i), Time: now - int64(106-i), Operation: OpEdit,
		})
	}
	if err := s.writeUTHistoryIndex(id, idx); err != nil {
		t.Fatal(err)
	}
	_ = s.recordVersionUT(id, OpEdit, "new")
	versions, _ := s.ListUTVersions(9, id)
	if len(versions) != maxVersionsPerUT+1 {
		t.Fatalf("应保留 %d 条（上限+起点豁免），得到 %d", maxVersionsPerUT+1, len(versions))
	}
	if versions[len(versions)-1].Version != "v000001" {
		t.Errorf("起点版本应保留在最末: %s", versions[len(versions)-1].Version)
	}
	if versions[0].Version != "v000106" {
		t.Errorf("最新版应在最前: %s", versions[0].Version)
	}
}

func TestUpdateMeta(t *testing.T) {
	s, st, id, uid := newHistService(t, "meta")
	if err := s.UpdateMeta(uid, id, "新名", "新描述"); err != nil {
		t.Fatalf("UpdateMeta: %v", err)
	}
	// 文件同步（单一真相）
	raw, _ := os.ReadFile(filepath.Join(s.Dir(id), "template.json"))
	var meta map[string]any
	if json.Unmarshal(raw, &meta) != nil || meta["name"] != "新名" {
		t.Errorf("template.json 未同步: %s", raw)
	}
	var row store.UserTemplate
	st.DB.First(&row, "id = ?", id)
	if row.Name != "新名" {
		t.Errorf("DB 未更新: %q", row.Name)
	}
	versions, _ := s.ListUTVersions(uid, id)
	if len(versions) != 1 || versions[0].Operation != OpMeta {
		t.Fatalf("应记一条 meta 版: %+v", versions)
	}
	// 空参数 no-op
	if err := s.UpdateMeta(uid, id, "", ""); err != nil {
		t.Errorf("空参数应放行不做事: %v", err)
	}
	// published：只改 DB（列表管理），不动文件不记版
	if err := st.DB.Model(&store.UserTemplate{}).Where("id = ?", id).Update("status", "published").Error; err != nil {
		t.Fatal(err)
	}
	if err := s.UpdateMeta(uid, id, "发布后名", ""); err != nil {
		t.Fatalf("published 改名应放行: %v", err)
	}
	var row2 store.UserTemplate
	st.DB.First(&row2, "id = ?", id)
	if row2.Name != "发布后名" {
		t.Errorf("published 改名未入库: %q", row2.Name)
	}
	raw2, _ := os.ReadFile(filepath.Join(s.Dir(id), "template.json"))
	var meta2 map[string]any
	json.Unmarshal(raw2, &meta2)
	if meta2["name"] != "新名" {
		t.Errorf("published 改名不应动 template.json: %v", meta2["name"])
	}
	versions, _ = s.ListUTVersions(uid, id)
	if len(versions) != 1 {
		t.Errorf("published 改名不应记版: %d", len(versions))
	}
}

func TestExecCustomToolDirty(t *testing.T) {
	s, _, id, _ := newHistService(t, "dirty")
	row := &store.UserTemplate{ID: id, BaseID: "tech-sharing", Status: "draft"}
	if _, dirty := s.execCustomTool(row, "write_tokens", `{"tokens":{"--accent":"#123456"}}`); !dirty {
		t.Error("write_tokens 有 tokens 应 dirty")
	}
	if _, dirty := s.execCustomTool(row, "write_tokens", `{"tokens":{}}`); dirty {
		t.Error("write_tokens 空 tokens 不应 dirty")
	}
	if _, dirty := s.execCustomTool(row, "set_meta", `{"name":"N"}`); !dirty {
		t.Error("set_meta 有名应 dirty")
	}
	if _, dirty := s.execCustomTool(row, "set_meta", `{"name":""}`); dirty {
		t.Error("set_meta 空参不应 dirty")
	}
	if _, dirty := s.execCustomTool(row, "finish", `{"reply":"done"}`); dirty {
		t.Error("finish 不应 dirty")
	}
}

func TestDeleteAndClearHistory(t *testing.T) {
	s, _, id, uid := newHistService(t, "del")
	_ = s.recordVersionUT(id, OpEdit, "a")
	if err := os.WriteFile(filepath.Join(s.Dir(id), "style.css"), []byte(".x{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	_ = s.recordVersionUT(id, OpEdit, "b")
	if err := s.DeleteUTVersion(uid, id, "v000001"); err != nil {
		t.Fatalf("删单版: %v", err)
	}
	versions, _ := s.ListUTVersions(uid, id)
	if len(versions) != 1 || versions[0].Version != "v000002" {
		t.Fatalf("删后应剩 v000002: %+v", versions)
	}
	if _, err := os.Stat(filepath.Join(s.historyDir(id), "v000001.json")); !os.IsNotExist(err) {
		t.Error("v000001 快照文件应被删除")
	}
	// NextSeq 不清零：清空后下一版从 v000003 继续
	n, err := s.ClearUTHistory(uid, id)
	if err != nil || n != 1 {
		t.Fatalf("清空: n=%d err=%v", n, err)
	}
	_ = s.recordVersionUT(id, OpEdit, "c")
	versions, _ = s.ListUTVersions(uid, id)
	if len(versions) != 1 || versions[0].Version != "v000003" {
		t.Fatalf("清空后编号应延续: %+v", versions)
	}
}
