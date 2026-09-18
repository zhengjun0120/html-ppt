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
	"regexp"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"html-ppt/backend/internal/service/template"
	"html-ppt/backend/internal/store"
	"html-ppt/backend/internal/vision"
)

type Service struct {
	reg        *template.Registry
	st         *store.Store
	root       string // data/user-templates
	chromePath string
	baseURL    string // 回环地址（demo 渲染走公开静态 /user-templates/）

	// 定制对话的会话表（内存态；重启即清空——对话历史不是重要数据）
	custMu   sync.Mutex
	sessions map[string]*customizeSession
}

func New(reg *template.Registry, st *store.Store, root, chromePath, baseURL string) *Service {
	return &Service{reg: reg, st: st, root: root, chromePath: chromePath, baseURL: baseURL,
		sessions: make(map[string]*customizeSession)}
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
func (s *Service) Fork(userID uint, baseID, name string) (*store.UserTemplate, error) {
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
		Name: fmt.Sprint(meta["name"]),
		Description: fmt.Sprint(meta["description"]),
		Visibility: "private", Status: "draft",
	}
	if err := s.st.DB.Create(row).Error; err != nil {
		return nil, err
	}
	// 私有模板也挂进注册表：owner 在生成管线里立即可用（见 deck.SelectTemplate
	// 的 ut- 归属校验，非 owner 拿不到私有模板）。
	if err := s.reg.MountUser(dir); err != nil {
		_ = s.st.DB.Delete(row)
		return nil, fmt.Errorf("克隆出的模板没过校验（源模板损坏？）: %w", err)
	}
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
	updates := map[string]any{}
	if name != "" {
		updates["name"] = name
	}
	if description != "" {
		updates["description"] = description
	}
	if len(updates) == 0 {
		return nil
	}
	return s.st.DB.Model(row).Updates(updates).Error
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

// ---------- 发布门禁（自动，两关） ----------

// PublishReport 两关门禁的结果（进 PublishReport 字段，前端展示）。
type PublishReport struct {
	Structure string           `json:"structure"`
	Render    *RenderGateResult `json:"render,omitempty"`
	Note      string           `json:"note,omitempty"`
}

type RenderGateResult struct {
	Pages   int      `json:"pages"`
	MinFill float64  `json:"min_fill"`
	MaxFont float64  `json:"max_flag_font"`
	Flaws   []string `json:"flaws,omitempty"`
}

// Publish 跑发布门禁：
//  1. 结构校验（与内置模板同一套 loadTemplate 规则 + CSS 黑名单扫描）；
//  2. demo 渲染量测（headless 实拍：无溢出、填充率 ≥45%、无 <13px 内容字号）。
//
// 全过 → public + published + Registry 挂载；任一失败 → failed + 可读原因。
// 注：规划阶段的"LLM 真生成冒烟"在 v1 以渲染量测替代（确定性、零额度消耗）；
// LLM 冒烟留给从零构建一起做。
func (s *Service) Publish(ctx context.Context, userID uint, id string) (*PublishReport, error) {
	row, err := s.GetOwned(userID, id)
	if err != nil {
		return nil, err
	}
	if row.Status == "publishing" {
		return nil, fmt.Errorf("发布门禁正在运行，请稍候")
	}
	if err := s.setPublishing(row); err != nil {
		return nil, err
	}
	fail := func(err error) (*PublishReport, error) {
		_ = s.markFailed(row, err)
		return nil, err
	}

	dir := s.Dir(id)

	// 门禁 1：结构校验（registry 同一套规则）+ 用户模板专属黑名单
	if err := s.reg.ValidateUserDir(dir); err != nil {
		return fail(fmt.Errorf("结构校验未过：%v", err))
	}
	if err := s.securityScan(dir); err != nil {
		return fail(fmt.Errorf("安全扫描未过：%v", err))
	}

	// 门禁 2：demo 渲染量测（公开静态 /user-templates/<id>/index.html）
	url := strings.TrimRight(s.baseURL, "/") + "/user-templates/" + id + "/index.html"
	d, err := vision.CaptureV2(ctx, vision.OptionsV2{URL: url, ChromePath: s.chromePath, Timeout: 2 * time.Minute})
	if err != nil {
		return fail(fmt.Errorf("demo 渲染失败：%v（检查 demo 是否可独立打开）", err))
	}
	// 填充率下限按版式指纹豁免：hero（封面/章节/收尾）与 quote 是刻意的稀疏页，
	// 它们的质量靠"有没有视觉锚点"而不是"塞没塞满"——把 anchor 判断留给量测
	// 的溢出/字号项与人工预览。内容型版式一律 ≥45%。
	patterns := layoutPatterns(filepath.Join(dir, "layouts.md"))
	render := &RenderGateResult{Pages: len(d.Slides), MinFill: 100}
	for _, sl := range d.Slides {
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
	if len(render.Flaws) > 0 {
		return fail(fmt.Errorf("渲染量测未过：%s", strings.Join(render.Flaws, "；")))
	}

	// 全过：挂进注册表 + 落状态
	if err := s.reg.MountUser(dir); err != nil {
		return fail(fmt.Errorf("注册失败：%v", err))
	}
	report := &PublishReport{Structure: "ok", Render: render,
		Note: "v1 冒烟为 demo 渲染量测（确定性）；LLM 真生成冒烟计划于从零构建版本加入"}
	rep, _ := json.Marshal(report)
	if err := s.markPublished(row, string(rep)); err != nil {
		return nil, err
	}
	return report, nil
}

// securityScan 用户模板的安全黑名单：style.css 禁外链与表达式；index.html 禁
// 额外脚本与内联事件（demo 会被其他用户渲染）。
func (s *Service) securityScan(dir string) error {
	css, err := os.ReadFile(filepath.Join(dir, "style.css"))
	if err != nil {
		return fmt.Errorf("style.css 缺失")
	}
	cssLow := strings.ToLower(string(css))
	for _, bad := range []string{"url(", "@import", "expression(", "behavior:", "-moz-binding"} {
		if strings.Contains(cssLow, bad) {
			return fmt.Errorf("style.css 含被禁用的 %q（外链/表达式是数据渗出通道）", bad)
		}
	}
	html, err := os.ReadFile(filepath.Join(dir, "index.html"))
	if err != nil {
		return fmt.Errorf("index.html 缺失")
	}
	htmlLow := strings.ToLower(string(html))
	for _, tag := range []string{"<script", "<iframe", "<object", "<embed"} {
		if idx := strings.Index(htmlLow, tag); idx >= 0 {
			// runtime.js 是唯一的合法脚本（本仓库资产，且 head 引用固定）
			if tag != "<script" || !strings.Contains(htmlLow[idx:], "/assets/deck-v2/runtime.js") {
				return fmt.Errorf("index.html 含被禁用的元素 %q", tag)
			}
		}
	}
	for _, ev := range []string{" onload=", " onerror=", " onclick="} {
		if strings.Contains(htmlLow, ev) {
			return fmt.Errorf("index.html 含内联事件 %q", strings.TrimSpace(ev))
		}
	}
	return nil
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
