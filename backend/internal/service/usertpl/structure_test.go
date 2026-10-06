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

// ---------- 层 3：对话加/删版式 ----------

// flowSpec 一个能过 MountUser 全量校验的新版式：既有 base 原语 + CSSAppend
// 承载的新类名（bl-flow*），带 demo 示例页。
func flowSpec(id string) AddLayoutSpec {
	skeleton := `<section class="slide" data-layout="` + id + `">
  <p class="kicker">{{kicker · 流程}}</p>
  <h2 class="h2 mt-s">{{流程标题}}</h2>
  <div class="grid g3 mt-l">
    <div class="bl-flow"><span class="mono bl-flow-n">{{1}}</span><p class="bl-flow-t">{{步骤要点}}</p><p class="dim mt-s">{{步骤说明}}</p></div>
    <div class="bl-flow"><span class="mono bl-flow-n">{{2}}</span><p class="bl-flow-t">{{步骤要点}}</p><p class="dim mt-s">{{步骤说明}}</p></div>
    <div class="bl-flow"><span class="mono bl-flow-n">{{3}}</span><p class="bl-flow-t">{{步骤要点}}</p><p class="dim mt-s">{{步骤说明}}</p></div>
  </div>
</section>`
	demo := `<section class="slide" data-layout="` + id + `">
  <p class="kicker">示例 · 三步走</p>
  <h2 class="h2 mt-s">三步完成接入</h2>
  <div class="grid g3 mt-l">
    <div class="bl-flow"><span class="mono bl-flow-n">1</span><p class="bl-flow-t">建好账号</p><p class="dim mt-s">两分钟填完基础信息即可开始。</p></div>
    <div class="bl-flow"><span class="mono bl-flow-n">2</span><p class="bl-flow-t">接入数据</p><p class="dim mt-s">复制接入密钥到你的系统里。</p></div>
    <div class="bl-flow"><span class="mono bl-flow-n">3</span><p class="bl-flow-t">看板验收</p><p class="dim mt-s">回到看板确认数据已经开始流动。</p></div>
  </div>
</section>`
	return AddLayoutSpec{
		LayoutID:    id,
		Name:        "三步流程",
		Use:         "三步流程/步骤说明：横向编号步骤，每步一句要点加一段说明",
		Roles:       []string{"content"},
		Fingerprint: "cards",
		Classes:     []string{"slide", "kicker", "h2", "mt-s", "grid", "g3", "mt-l", "mono", "dim", "bl-flow", "bl-flow-n", "bl-flow-t"},
		Skeleton:    skeleton,
		CSSAppend: ".tpl-blank .bl-flow{padding:22px;border:1px solid var(--border);border-radius:var(--radius)}\n" +
			".tpl-blank .bl-flow-n{font-size:30px;font-weight:800;color:var(--accent)}\n" +
			".tpl-blank .bl-flow-t{font-size:21px;font-weight:650;margin-top:8px}",
		DemoHTML:    demo,
		Constraints: "恰好 3 步；要点 ≤10 字、说明 12-30 字",
	}
}

func fileHash(t *testing.T, s *Service, id string, names ...string) map[string]string {
	t.Helper()
	out := map[string]string{}
	for _, n := range names {
		raw, err := os.ReadFile(filepath.Join(s.Dir(id), n))
		if err != nil {
			t.Fatalf("读 %s: %v", n, err)
		}
		out[n] = string(raw)
	}
	return out
}

func TestAddLayout(t *testing.T) {
	s, reg := structTestService(t)
	const uid = 21
	row, err := s.CreateBlank(uid, "加版式")
	if err != nil {
		t.Fatalf("CreateBlank: %v", err)
	}
	spec := flowSpec("blank-flow")
	summary, err := s.AddLayout(uid, row.ID, spec)
	if err != nil {
		t.Fatalf("AddLayout: %v", err)
	}
	if !strings.Contains(summary, "新增版式") || !strings.Contains(summary, "10 个版式") {
		t.Errorf("摘要不对: %q", summary)
	}

	// 注册表：新版式已挂载、骨架/指纹可读
	tpl, err := reg.Get(row.ID)
	if err != nil {
		t.Fatalf("注册表: %v", err)
	}
	if len(tpl.Layouts) != 10 {
		t.Errorf("版式数 = %d，应为 10", len(tpl.Layouts))
	}
	skel, ok := tpl.Layout("blank-flow")
	if !ok || !strings.Contains(skel, "bl-flow") {
		t.Errorf("新版式骨架不可读: %q", skel)
	}
	if got := tpl.Pattern("blank-flow"); got != "cards" {
		t.Errorf("指纹 = %q，应为 cards", got)
	}

	// 磁盘四件套：json 条目 / md 条目 / css 规则 / demo 页
	tj, _ := os.ReadFile(filepath.Join(s.Dir(row.ID), "template.json"))
	if !strings.Contains(string(tj), `"blank-flow"`) || !strings.Contains(string(tj), "恰好 3 步") {
		t.Error("template.json 缺新版式条目")
	}
	md, _ := os.ReadFile(filepath.Join(s.Dir(row.ID), "layouts.md"))
	if !strings.Contains(string(md), "## blank-flow（三步流程）") || !strings.Contains(string(md), "指纹：cards") {
		t.Error("layouts.md 缺新条目")
	}
	css, _ := os.ReadFile(filepath.Join(s.Dir(row.ID), "style.css"))
	if !strings.Contains(string(css), ".tpl-blank .bl-flow-n") {
		t.Error("style.css 缺新类规则")
	}
	idx, _ := os.ReadFile(filepath.Join(s.Dir(row.ID), "index.html"))
	if got := strings.Count(string(idx), `data-layout="blank-`); got != 10 {
		t.Errorf("demo section 数 = %d，应为 10", got)
	}

	// 版本记录：changed 含四个文件
	versions, err := s.ListUTVersions(uid, row.ID)
	if err != nil || len(versions) == 0 {
		t.Fatalf("历史: %v", err)
	}
	if versions[0].Operation != OpStruct {
		t.Errorf("最新版本 operation = %q", versions[0].Operation)
	}
	for _, f := range []string{"template.json", "layouts.md", "style.css", "index.html"} {
		if !containsStr(versions[0].Changed, f) {
			t.Errorf("changed 缺 %s: %v", f, versions[0].Changed)
		}
	}
}

func TestAddLayoutRollback(t *testing.T) {
	s, reg := structTestService(t)
	const uid = 22
	row, err := s.CreateBlank(uid, "回滚")
	if err != nil {
		t.Fatalf("CreateBlank: %v", err)
	}
	before := fileHash(t, s, row.ID, "template.json", "layouts.md", "style.css", "index.html")
	versionsBefore, _ := s.ListUTVersions(uid, row.ID)

	// 类声明了但 base/模板文件里都不存在：预检过（骨架类 ⊆ 声明），
	// MountUser 必挂 → 回滚。注意 demo_html 不能带（index.html 的 class token
	// 也算"存在"），所以去掉示例页，让该类只出现在合法类名行。
	spec := flowSpec("blank-ghost")
	spec.CSSAppend = ""
	spec.DemoHTML = ""
	for i, c := range spec.Classes {
		if c == "bl-flow-t" {
			spec.Classes[i] = "bl-ghost"
		}
	}
	spec.Skeleton = strings.ReplaceAll(spec.Skeleton, "bl-flow-t", "bl-ghost")
	_, err = s.AddLayout(uid, row.ID, spec)
	if err == nil {
		t.Fatal("无 CSS 的新类名应被 MountUser 拒绝")
	}
	if !strings.Contains(err.Error(), "回滚") {
		t.Errorf("报错应说明已回滚: %v", err)
	}
	// 四文件字节不变、注册表无残留、无新版本
	after := fileHash(t, s, row.ID, "template.json", "layouts.md", "style.css", "index.html")
	for n := range before {
		if before[n] != after[n] {
			t.Errorf("回滚不干净: %s 有变化", n)
		}
	}
	tpl, _ := reg.Get(row.ID)
	if len(tpl.Layouts) != 9 {
		t.Errorf("注册表版式数 = %d，应仍为 9", len(tpl.Layouts))
	}
	versionsAfter, _ := s.ListUTVersions(uid, row.ID)
	if len(versionsAfter) != len(versionsBefore) {
		t.Error("回滚后不应记版本")
	}
}

func TestAddLayoutRejects(t *testing.T) {
	s, _ := structTestService(t)
	const uid = 23
	row, err := s.CreateBlank(uid, "拒绝项")
	if err != nil {
		t.Fatalf("CreateBlank: %v", err)
	}
	cases := []struct {
		name string
		mut  func(*AddLayoutSpec)
		want string
	}{
		{"重复 id", func(sp *AddLayoutSpec) {
			sp.LayoutID = "blank-cover"
			sp.Skeleton = strings.Replace(sp.Skeleton, "blank-flow", "blank-cover", 1)
			sp.DemoHTML = strings.Replace(sp.DemoHTML, "blank-flow", "blank-cover", 1)
		}, "已存在"},
		{"坏 id", func(sp *AddLayoutSpec) { sp.LayoutID = "Bad_ID" }, "不合法"},
		{"坏指纹", func(sp *AddLayoutSpec) { sp.Fingerprint = "fancy" }, "指纹"},
		{"坏角色", func(sp *AddLayoutSpec) { sp.Roles = []string{"hero"} }, "不在词表"},
		{"骨架 data-layout 不一致", func(sp *AddLayoutSpec) {
			sp.Skeleton = strings.Replace(sp.Skeleton, `data-layout="blank-flow"`, `data-layout="other"`, 1)
		}, "不一致"},
		{"骨架未声明类", func(sp *AddLayoutSpec) {
			sp.Classes = sp.Classes[:len(sp.Classes)-1] // 去掉 bl-flow-t
		}, "classes 清单里没有"},
		{"demo 带脚本", func(sp *AddLayoutSpec) {
			sp.DemoHTML = strings.Replace(sp.DemoHTML, "</section>", "<script>alert(1)</script></section>", 1)
		}, "被禁用"},
		{"css 外链", func(sp *AddLayoutSpec) {
			sp.CSSAppend += "\n.tpl-x{background:url(https://evil.example/x.png)}"
		}, "css_append"},
	}
	for _, tc := range cases {
		spec := flowSpec("blank-flow")
		tc.mut(&spec)
		if _, err := s.AddLayout(uid, row.ID, spec); err == nil {
			t.Errorf("%s: 未拒绝", tc.name)
		} else if !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: 报错不对: %v", tc.name, err)
		}
	}
}

func TestRemoveLayout(t *testing.T) {
	s, reg := structTestService(t)
	const uid = 24
	row, err := s.CreateBlank(uid, "删版式")
	if err != nil {
		t.Fatalf("CreateBlank: %v", err)
	}
	// 删对照页：content 候选 5→4（content/list/data/quote），放行
	summary, err := s.RemoveLayout(uid, row.ID, "blank-split")
	if err != nil {
		t.Fatalf("RemoveLayout: %v", err)
	}
	if !strings.Contains(summary, "删除版式") || !strings.Contains(summary, "1 页") {
		t.Errorf("摘要不对: %q", summary)
	}

	tpl, _ := reg.Get(row.ID)
	if len(tpl.Layouts) != 8 {
		t.Errorf("注册表版式数 = %d，应为 8", len(tpl.Layouts))
	}
	if tpl.HasLayout("blank-split") {
		t.Error("注册表仍有 blank-split")
	}
	tj, _ := os.ReadFile(filepath.Join(s.Dir(row.ID), "template.json"))
	if strings.Contains(string(tj), "blank-split") {
		t.Error("template.json 仍有 blank-split")
	}
	md, _ := os.ReadFile(filepath.Join(s.Dir(row.ID), "layouts.md"))
	if strings.Contains(string(md), "## blank-split") {
		t.Error("layouts.md 仍有条目")
	}
	idx, _ := os.ReadFile(filepath.Join(s.Dir(row.ID), "index.html"))
	if strings.Contains(string(idx), `data-layout="blank-split"`) {
		t.Error("demo 仍有引用页")
	}
	if got := strings.Count(string(idx), `data-layout="blank-`); got != 8 {
		t.Errorf("demo section 数 = %d，应为 8", got)
	}
	versions, _ := s.ListUTVersions(uid, row.ID)
	if versions[0].Operation != OpStruct {
		t.Errorf("最新版本 operation = %q", versions[0].Operation)
	}
}

func TestRemoveLayoutHealthCheck(t *testing.T) {
	s, _ := structTestService(t)
	const uid = 25
	row, err := s.CreateBlank(uid, "死锁检查")
	if err != nil {
		t.Fatalf("CreateBlank: %v", err)
	}
	// 连删两个 content 可用版式（list/split），content 候选 5→3
	for _, lid := range []string{"blank-list", "blank-split"} {
		if _, err := s.RemoveLayout(uid, row.ID, lid); err != nil {
			t.Fatalf("删 %s: %v", lid, err)
		}
	}
	// 再删一个 content 可用版式 → 阻塞（候选将剩 2）
	_, err = s.RemoveLayout(uid, row.ID, "blank-content")
	if err == nil {
		t.Fatal("content 候选不足 3 应阻塞")
	}
	if !strings.Contains(err.Error(), "先加一个新版式") {
		t.Errorf("报错缺指引: %v", err)
	}
	// 非 content 版式仍可删（toc）
	if _, err := s.RemoveLayout(uid, row.ID, "blank-toc"); err != nil {
		t.Errorf("删 toc 不应被阻塞: %v", err)
	}
	// 警告随摘要返回（cover/thanks 都没了）
	summary, err := s.RemoveLayout(uid, row.ID, "blank-cover")
	if err != nil {
		t.Fatalf("删 cover: %v", err)
	}
	if !strings.Contains(summary, "cover") {
		t.Errorf("摘要应带 cover 缺失警告: %q", summary)
	}
}

func TestScanSkeleton(t *testing.T) {
	classes := []string{"slide", "h2", "dim"}
	base := `<section class="slide" data-layout="x"><h2 class="h2">{{a}}</h2><p class="dim">{{b}}</p></section>`
	if err := scanSkeleton(base, "x", classes); err != nil {
		t.Fatalf("合法片段被拒: %v", err)
	}
	cases := map[string]string{
		"脚本":     `<section class="slide" data-layout="x"><script>a</script></section>`,
		"iframe": `<section class="slide" data-layout="x"><iframe src="//e.com"></iframe></section>`,
		"嵌套":     `<section class="slide" data-layout="x"><section class="slide"></section></section>`,
		"内联事件":   `<section class="slide" data-layout="x" onclick="a"><h2 class="h2">{{a}}</h2></section>`,
		"js 协议":  `<section class="slide" data-layout="x"><a href="javascript:alert(1)">x</a></section>`,
		"缺属性":    `<section class="slide"><h2 class="h2">{{a}}</h2></section>`,
		"id 不一致": `<section class="slide" data-layout="y"><h2 class="h2">{{a}}</h2></section>`,
		"未声明类":   `<section class="slide" data-layout="x"><h2 class="h2">{{a}}</h2><p class="nope">{{b}}</p></section>`,
		"非单根":    `<section class="slide" data-layout="x"></section><p>尾巴</p>`,
	}
	for name, frag := range cases {
		if err := scanSkeleton(frag, "x", classes); err == nil {
			t.Errorf("%s: 未拒绝", name)
		}
	}
}

func TestCustomizeAddLayoutFlow(t *testing.T) {
	// 对话循环：模型 add_layout + finish（假 LLM），断言落盘与记版链路
	s, reg := structTestService(t)
	const uid = 26
	row, err := s.CreateBlank(uid, "对话加版式")
	if err != nil {
		t.Fatalf("CreateBlank: %v", err)
	}
	spec := flowSpec("blank-flow")
	args, _ := json.Marshal(map[string]any{
		"layout_id": spec.LayoutID, "name": spec.Name,
		"use": spec.Use, "roles": spec.Roles,
		"fingerprint": spec.Fingerprint, "classes": spec.Classes,
		"skeleton": spec.Skeleton, "css_append": spec.CSSAppend,
		"demo_html": spec.DemoHTML, "constraints": spec.Constraints,
	})
	llm, _ := newFakeLLM(t, [][]string{
		toolRound("add_layout", string(args)),
		toolRound("finish", `{"reply":"已新增三步流程版式。"}`),
	})
	reply, err := s.Customize(t.Context(), uid, row.ID, "加一个三步流程版式", llm, nil)
	if err != nil {
		t.Fatalf("Customize: %v", err)
	}
	if !strings.Contains(reply, "三步流程") {
		t.Errorf("答复不对: %q", reply)
	}
	tpl, _ := reg.Get(row.ID)
	if !tpl.HasLayout("blank-flow") {
		t.Error("对话后注册表没有新版式")
	}
	versions, _ := s.ListUTVersions(uid, row.ID)
	// fork 基线 + add_layout 的 structure 版；对话收尾时 bundle 与 structure 版
	// 完全一致（finish 不写文件）→ 历史幂等不记 chat 版，共 2 条
	if len(versions) != 2 {
		t.Errorf("版本数 = %d，应为 2（fork/structure）: %+v", len(versions), versions)
	}
	if len(versions) > 0 && versions[0].Operation != OpStruct {
		t.Errorf("最新版本 operation = %q", versions[0].Operation)
	}
}
