// Package usertpl 用户自定义模板（plan-v3 B）：克隆内置模板 → 对话定制 →
// 自动门禁发布 → 社区可用。
//
// 文件即模板：data/user-templates/<ut-id>/ 下是标准四件套（template.json /
// index.html / style.css / layouts.md / rules.md）。发布 = 通过门禁后挂进
// Registry 的用户层（deck.SelectTemplate 直接可用）；下架/删除即注销。
package usertpl

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"html-ppt/backend/internal/service/template"
	"html-ppt/backend/internal/store"
	"html-ppt/backend/internal/trace"
	"html-ppt/backend/internal/vision"
)

type Service struct {
	reg        *template.Registry
	st         *store.Store
	root       string // data/user-templates
	chromePath string
	baseURL    string         // 回环地址（demo 渲染走 nonce 端点）
	grants     *vision.Grants // 体检渲染的一次性授权（Peek 语义，见 vision/grant.go）
	traceCfg   trace.Config   // 定制对话的观测落盘（2026-09-28 定制接入观测台）

	// 定制对话的会话表（内存态；重启即清空——对话历史不是重要数据）
	custMu   sync.Mutex
	sessions map[string]*customizeSession

	// per-template 写锁（history.go lockUT）：对话与手动保存可能并发写同一模板
	locksMu sync.Mutex
	locks   map[string]*sync.Mutex
}

func New(reg *template.Registry, st *store.Store, root, chromePath, baseURL string, grants *vision.Grants, traceCfg trace.Config) *Service {
	return &Service{reg: reg, st: st, root: root, chromePath: chromePath, baseURL: baseURL, grants: grants, traceCfg: traceCfg,
		sessions: make(map[string]*customizeSession)}
}

// Peek 读模板行，不做归属判断（受控公开端点专用：先看状态再决定要不要鉴权，
// 见 handler 的 UserTemplatePublicFile）。
func (s *Service) Peek(id string) (*store.UserTemplate, error) {
	var row store.UserTemplate
	if err := s.st.DB.First(&row, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("模板不存在")
	}
	return &row, nil
}

func (s *Service) sessionFor(key string) *customizeSession {
	s.custMu.Lock()
	defer s.custMu.Unlock()
	if s.sessions == nil {
		s.sessions = make(map[string]*customizeSession)
	}
	sess, ok := s.sessions[key]
	if !ok {
		sess = &customizeSession{}
		s.sessions[key] = sess
	}
	return sess
}

func (s *Service) Dir(id string) string { return filepath.Join(s.root, id) }

// ---------- Fork ----------

// forkCopies 克隆时带走的文件（源目录的 preview/、UPSTREAM/DEMO 文档不跟过来）。
var forkCopies = []string{"template.json", "index.html", "style.css", "layouts.md", "rules.md"}

// Fork 从内置模板克隆一份私有副本。
// baseID == "_blank" 是约定的"空白来源"：走 CreateBlank 脚手架而不是克隆目录。
// 用 :id 的取值而不是静态路由段（/api/templates/blank/fork 之类），避开 gin 的
// 静态段与参数段同位 panic（router.go 的 recent-sessions 注释记过同一个坑）。
func (s *Service) Fork(userID uint, baseID, name string) (*store.UserTemplate, error) {
	if baseID == "_blank" {
		return s.CreateBlank(userID, name)
	}
	src, err := s.reg.BuiltinDir(baseID)
	if err != nil {
		return nil, err
	}
	id, err := newID()
	if err != nil {
		return nil, err
	}
	dir := s.Dir(id)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	for _, f := range forkCopies {
		if err := copyFile(filepath.Join(src, f), filepath.Join(dir, f)); err != nil {
			return nil, fmt.Errorf("复制 %s 失败: %w", f, err)
		}
	}
	// template.json：换 id/名字/来源标注（fork 出去的模板不再指向 html-ppt-skill，
	// 但保留 base 的溯源链）
	var meta map[string]any
	raw, _ := os.ReadFile(filepath.Join(dir, "template.json"))
	if err := json.Unmarshal(raw, &meta); err != nil {
		return nil, fmt.Errorf("源 template.json 损坏: %w", err)
	}
	meta["id"] = id
	if name != "" {
		meta["name"] = name
	} else {
		meta["name"] = fmt.Sprintf("%s（我的版本）", meta["name"])
	}
	meta["source"] = map[string]string{
		"derived_from": fmt.Sprintf("builtin:%s", baseID),
		"license":      "MIT",
	}
	out, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(dir, "template.json"), out, 0o644); err != nil {
		return nil, err
	}
	os.WriteFile(filepath.Join(dir, "ADAPTATION.md"), []byte(fmt.Sprintf(
		"# %s\n\n- fork 自内置模板 `%s`（%s），可自由定制。\n- 定制入口：对话改 token（色板/字体/圆角）与 rules 文案；结构契约（版式/类名）继承 base。\n- 创建时间：%s\n",
		meta["name"], baseID, baseID, time.Now().Format("2006-01-02 15:04"))), 0o644)

	row := &store.UserTemplate{
		ID: id, UserID: userID, BaseID: baseID,
		Name:        fmt.Sprint(meta["name"]),
		Description: fmt.Sprint(meta["description"]),
		Visibility:  "private", Status: "draft",
	}
	if err := s.st.DB.Create(row).Error; err != nil {
		return nil, err
	}
	// 私有模板也挂进注册表：owner 在生成管线里立即可用（见 deck.SelectTemplate
	// 的 ut- 归属校验，非 owner 拿不到私有模板）。
	if err := s.reg.MountUser(dir); err != nil {
		_ = s.st.DB.Delete(row)
		_ = os.RemoveAll(dir) // 连目录一起清：只删行会留孤儿目录（实测遗留问题）
		return nil, fmt.Errorf("克隆出的模板没过校验（源模板损坏？）: %w", err)
	}
	// 起点基线：fork 本身记一版，"回得去起点"从这里开始
	_ = s.recordVersionUT(id, OpFork, "克隆自 builtin:"+baseID)
	return row, nil
}

// ---------- CRUD ----------

func (s *Service) List(userID uint) ([]store.UserTemplate, error) {
	var rows []store.UserTemplate
	err := s.st.DB.Where("user_id = ?", userID).Order("updated_at DESC").Find(&rows).Error
	return rows, err
}

// Community 公开模板清单（画廊「社区模板」；带作者邮箱前缀做署名）。
func (s *Service) Community() ([]map[string]any, error) {
	var rows []store.UserTemplate
	if err := s.st.DB.Where("visibility = ? AND status = ?", "public", "published").
		Order("updated_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		author := ""
		var u store.User
		if err := s.st.DB.First(&u, r.UserID).Error; err == nil {
			author = strings.SplitN(u.Email, "@", 2)[0]
		}
		out = append(out, map[string]any{
			"id": r.ID, "name": r.Name, "description": r.Description,
			"base_id": r.BaseID, "author": author, "updated_at": r.UpdatedAt,
		})
	}
	return out, nil
}

// getRow 取行 + 归属判断：owner 必过；非 owner 仅当公开且已发布（只读场景）。
func (s *Service) getRow(userID uint, id string, publicReadable bool) (*store.UserTemplate, error) {
	var row store.UserTemplate
	if err := s.st.DB.First(&row, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("模板不存在")
	}
	if row.UserID != userID {
		if !(publicReadable && row.Visibility == "public" && row.Status == "published") {
			return nil, fmt.Errorf("模板不存在")
		}
	}
	return &row, nil
}

func (s *Service) GetOwned(userID uint, id string) (*store.UserTemplate, error) {
	return s.getRow(userID, id, false)
}

func (s *Service) GetReadable(userID uint, id string) (*store.UserTemplate, error) {
	return s.getRow(userID, id, true)
}

func (s *Service) UpdateMeta(userID uint, id, name, description string) error {
	row, err := s.GetOwned(userID, id)
	if err != nil {
		return err
	}
	if name == "" && description == "" {
		return nil
	}
	// published：只改 DB（改名/描述是列表管理，不动文件——D2 的文件锁不含它，
	// 但也不记版本：没改设计内容）。社区列表的名称来自 DB，展示即时生效。
	if row.Status == "published" {
		updates := map[string]any{}
		if name != "" {
			updates["name"] = name
		}
		if description != "" {
			updates["description"] = description
		}
		return s.st.DB.Model(row).Updates(updates).Error
	}
	// draft：DB 与 template.json 一起改（单一真相），并记一条 meta 版本。
	// 现状问题修正：此前 UpdateMeta 只改 DB，template.json 会与 DB 漂移。
	if err := s.applyMetaFile(id, name, description); err != nil {
		return err
	}
	updates := map[string]any{}
	if name != "" {
		updates["name"] = name
	}
	if description != "" {
		updates["description"] = description
	}
	if err := s.st.DB.Model(row).Updates(updates).Error; err != nil {
		return err
	}
	return s.recordVersionUT(id, OpMeta, "改名/描述")
}

func (s *Service) Delete(userID uint, id string) error {
	row, err := s.GetOwned(userID, id)
	if err != nil {
		return err
	}
	s.reg.UnmountUser(id)
	if err := s.st.DB.Delete(row).Error; err != nil {
		return err
	}
	return os.RemoveAll(s.Dir(id))
}

// SetPublishing 进入门禁运行态（status=publishing；UI 禁止重复触发）。
func (s *Service) setPublishing(row *store.UserTemplate) error {
	return s.st.DB.Model(row).Updates(map[string]any{
		"status": "publishing", "publish_error": "", "publish_report": "",
	}).Error
}

func (s *Service) markFailed(row *store.UserTemplate, err error) error {
	row.Status = "failed"
	row.PublishError = err.Error()
	return s.st.DB.Model(row).Updates(map[string]any{
		"status": "failed", "publish_error": row.PublishError,
	}).Error
}

func (s *Service) markPublished(row *store.UserTemplate, report string) error {
	row.Status = "published"
	row.Visibility = "public"
	return s.st.DB.Model(row).Updates(map[string]any{
		"status": "published", "visibility": "public", "publish_report": report,
	}).Error
}

// Unpublish 下架：visibility 回 private + 注册表注销（已生成的 deck 快照不受影响）。
func (s *Service) Unpublish(userID uint, id string) error {
	row, err := s.GetOwned(userID, id)
	if err != nil {
		return err
	}
	s.reg.UnmountUser(id)
	return s.st.DB.Model(row).Updates(map[string]any{
		"visibility": "private", "status": "draft",
	}).Error
}

// ---------- 发布与质量体检 ----------

// PublishReport 发布/体检的结果（进 publish_report 字段，前端展示）。
// Render 只在体检（Checkup）里出现；发布自 2026-09-28 起不再跑渲染量测。
// Taste 是 AI 味提示（taste-skill 词表，体检时对 demo 各页跑 lint），只报告。
type PublishReport struct {
	Structure string            `json:"structure"`
	Render    *RenderGateResult `json:"render,omitempty"`
	Taste     []string          `json:"taste,omitempty"`
	Note      string            `json:"note,omitempty"`
}

type RenderGateResult struct {
	Pages   int      `json:"pages"`
	MinFill float64  `json:"min_fill"`
	MaxFont float64  `json:"max_flag_font"`
	Flaws   []string `json:"flaws,omitempty"`
}

// Publish 发布（2026-09-28 门禁降级，用户拍板）：
//  1. 安全扫描（写入时各路径已各自强制，这里兜底防绕过产品的直改）；
//  2. 挂载校验——MountUser 内部就是全量 loadTemplate 结构校验，挂不上 =
//     模板坏了，failed + 可读原因。
//
// 渲染量测（Chrome 实拍：溢出/填充率/最小字号）从发布门禁里拿掉了：它是唯一
// "看法类"的阈值（45%/13px 是审美不是功能），又是唯一贵的（10-30s），还把
// 空白骨架这类稀疏 demo 卡死在线外。完整保留为 Checkup 质量体检（不拦发布）。
// 结构安全不因降级而松动：结构坏 = 挂不上 = 生成侧解析不到，社区拿到手的
// 永远是挂载成功的模板。
func (s *Service) Publish(ctx context.Context, userID uint, id string) (*PublishReport, error) {
	row, err := s.GetOwned(userID, id)
	if err != nil {
		return nil, err
	}
	if row.Status == "publishing" {
		return nil, fmt.Errorf("发布正在运行，请稍候")
	}
	if err := s.setPublishing(row); err != nil {
		return nil, err
	}
	fail := func(err error) (*PublishReport, error) {
		_ = s.markFailed(row, err)
		return nil, err
	}

	dir := s.Dir(id)
	if err := s.securityScan(dir); err != nil {
		return fail(fmt.Errorf("安全扫描未过：%v", err))
	}
	if err := s.reg.MountUser(dir); err != nil {
		return fail(fmt.Errorf("结构校验未过：%v", err))
	}
	report := &PublishReport{Structure: "ok",
		Note: "发布即挂载校验，秒级完成；溢出/填充率等视觉质量可用「质量体检」随时检查"}
	rep, _ := json.Marshal(report)
	if err := s.markPublished(row, string(rep)); err != nil {
		return nil, err
	}
	return report, nil
}

// Checkup 质量体检（2026-09-28，从发布门禁降级而来）：headless 实拍 demo，
// 量测溢出/填充率/最小字号，报告返回给前端并落到 publish_report 字段。
// **不改状态、不拦任何东西**——draft 与 published 都能跑（纯只读诊断）。
// 发布者自己拿它当修复线索；社区质量靠发布者的自觉 + 体检报告可查。
func (s *Service) Checkup(ctx context.Context, userID uint, id string) (*PublishReport, error) {
	row, err := s.GetOwned(userID, id)
	if err != nil {
		return nil, err
	}
	dir := s.Dir(id)
	if err := s.securityScan(dir); err != nil {
		return &PublishReport{Structure: err.Error(), Note: "安全扫描未过：写入时本应拦截，请检查文件是否被绕过产品直接修改"}, nil
	}

	// 渲染走一次性 nonce 端点（草稿收口后 headless 导航带不了鉴权头，
	// nonce 是唯一的桥，docs/user-template-history-plan.md §3.5）
	if s.grants == nil {
		return nil, fmt.Errorf("渲染授权不可用")
	}
	nonce, err := s.grants.Issue(userID, "ut:"+id)
	if err != nil {
		return nil, fmt.Errorf("渲染授权签发失败: %w", err)
	}
	url := strings.TrimRight(s.baseURL, "/") + "/api/user-template-render/" + nonce + "/index.html"
	d, err := vision.CaptureV2(ctx, vision.OptionsV2{URL: url, ChromePath: s.chromePath, Timeout: 2 * time.Minute})
	if err != nil {
		return nil, fmt.Errorf("demo 渲染失败：%v（检查 demo 是否可独立打开）", err)
	}
	// 填充率下限按版式指纹豁免：hero（封面/章节/收尾）与 quote 是刻意的稀疏页，
	// 它们的质量靠"有没有视觉锚点"而不是"塞没塞满"。内容型版式一律 ≥45%。
	report := &PublishReport{Structure: "ok",
		Render: evaluateRender(d.Slides, layoutPatterns(filepath.Join(dir, "layouts.md"))),
		Note:   "体检报告（不影响发布）：flaws 为空即未发现溢出/稀疏/小字号问题"}
	// AI 味提示（taste-skill 词表）：demo 示例文案是生成范本，发布前让作者可见
	if raw, err := os.ReadFile(filepath.Join(dir, "index.html")); err == nil {
		report.Taste = lintDemoSections(string(raw))
	}
	rep, _ := json.Marshal(report)
	// 落到 publish_report 字段，刷新页面后报告还在。发布报告（publish 时写的）
	// 会被最近一次体检覆盖——两份都是"诊断快照"，留最新即可
	_ = s.st.DB.Model(row).Updates(map[string]any{"publish_report": string(rep)}).Error
	return report, nil
}

// evaluateRender 量测评估（纯函数，无 IO）：逐页算 flaw，返回量测结果。
// 从发布门禁时代原样搬来——阈值只在这一处，体检与未来任何调用方共用一套口径。
func evaluateRender(slides []vision.Slide2, patterns map[string]string) *RenderGateResult {
	render := &RenderGateResult{Pages: len(slides), MinFill: 100}
	for _, sl := range slides {
		if sl.OverflowY || sl.OverflowX {
			render.Flaws = append(render.Flaws, fmt.Sprintf("第 %d 页溢出画布", sl.Index+1))
		}
		pat := patterns[sl.Layout]
		exempt := pat == "hero" || pat == "quote"
		if sl.FillPct > 0 && sl.FillPct < render.MinFill && !exempt {
			render.MinFill = sl.FillPct
		}
		if sl.FillPct > 0 && sl.FillPct < 45 && !exempt {
			render.Flaws = append(render.Flaws, fmt.Sprintf("第 %d 页填充率仅 %.0f%%（大面积留白）", sl.Index+1, sl.FillPct))
		}
		if sl.MinFontPx > 0 && sl.MinFontPx < 13 {
			render.Flaws = append(render.Flaws, fmt.Sprintf("第 %d 页最小字号 %.0fpx", sl.Index+1, sl.MinFontPx))
		}
		if sl.MinFontPx > render.MaxFont {
			render.MaxFont = sl.MinFontPx
		}
	}
	return render
}

// securityScan 用户模板的安全黑名单（写入预检与发布门禁共用的组合入口）：
// style.css 与 index.html 的规则见 scan.go（scanCSS/scanHTML 纯函数）。
func (s *Service) securityScan(dir string) error {
	css, err := os.ReadFile(filepath.Join(dir, "style.css"))
	if err != nil {
		return fmt.Errorf("style.css 缺失")
	}
	if err := scanCSS(string(css)); err != nil {
		return err
	}
	html, err := os.ReadFile(filepath.Join(dir, "index.html"))
	if err != nil {
		return fmt.Errorf("index.html 缺失")
	}
	return scanHTML(string(html))
}

// layoutPatterns 从 layouts.md 提取 版式 id → 指纹（发布门禁的稀疏豁免判据）。
// 只认 "## id" 与其后的 "指纹：xxx" 行——与注册表解析器的语义一致，这里故意
// 不复用（注册表在 template 包内，usertpl 只需要这一个字段，不值得引整包解析器）。
func layoutPatterns(layoutsMDPath string) map[string]string {
	out := map[string]string{}
	raw, err := os.ReadFile(layoutsMDPath)
	if err != nil {
		return out
	}
	curID := ""
	for _, line := range strings.Split(string(raw), "\n") {
		t := strings.TrimSpace(line)
		if m := layoutHeadRe.FindStringSubmatch(t); m != nil {
			curID = m[1]
			continue
		}
		if strings.HasPrefix(t, "指纹：") && curID != "" {
			out[curID] = strings.TrimSpace(strings.TrimPrefix(t, "指纹："))
		}
	}
	return out
}

// layoutHeadRe 与注册表解析器同一条规则：## 后第一个 token 是版式 id
//（id 与中文名之间可以没有空格，如 `## qa（问答收尾）`）。
var layoutHeadRe = regexp.MustCompile(`^##\s*([A-Za-z][A-Za-z0-9_-]*)`)

// ---------- helpers ----------

func newID() (string, error) {
	buf := make([]byte, 4)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return "ut-" + hex.EncodeToString(buf), nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return nil
}
