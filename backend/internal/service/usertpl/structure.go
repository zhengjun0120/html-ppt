package usertpl

// 结构契约（定制工作台「版式」面板）：
//
// 读（StructureContract）展示"生成侧实际生效的契约"——数据来自 registry 挂载
// 快照（loadTemplate 全量校验后的产物），而不是磁盘文件：面板的意义就是让用户
// 看见生成管线真正会拿到什么；remount 失败的陈旧态因此如实可见（配合写路径的
// warning 提示去历史回滚）。
//
// 写（UpdateLayoutMeta）只放开版式元数据 name/use/roles。权威在 template.json
// （生成提示词的版式索引读它）；layouts.md 的「适用 role：」行是给人看的文档，
// best-effort 同步。骨架/类名/指纹不在这里动——动结构是层 3（加/删版式）的事。

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"html-ppt/backend/internal/service/template"
)

// layoutRoles 版式 role 词表：与 deck/v2.go 的 outlineRoles 同一套 9 种大纲角色。
// （deck 包那份不导出，这里重声明；两边语义耦合——role 决定 plan_pages 的版式
// 选择范围——改词表必须两处同步。）
var layoutRoles = map[string]bool{
	"cover": true, "toc": true, "divider": true, "content": true, "data": true,
	"quote": true, "code": true, "cta": true, "thanks": true,
}

// layoutRoleOrder 词表顺序（角色改写后按它排序，提示词与面板的展示稳定）。
var layoutRoleOrder = []string{"cover", "toc", "divider", "content", "data", "quote", "code", "cta", "thanks"}

const (
	maxLayoutNameLen = 40
	maxLayoutUseLen  = 200
	maxLayoutRoleCnt = 3
)

// StructureLayout 面板里一个版式的完整视图。
type StructureLayout struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Use         string         `json:"use,omitempty"`
	Roles       []string       `json:"roles,omitempty"`
	Constraints string         `json:"constraints,omitempty"`
	Pattern     string         `json:"pattern,omitempty"`  // 视觉指纹（节奏校验按它判"假多样性"）
	Repeats     map[string]int `json:"repeats,omitempty"`  // 数量契约：该类必须恰好 N 个
	Skeleton    string         `json:"skeleton,omitempty"` // layouts.md 的骨架代码
}

// DemoPage demo 区间里一个顶层 section 与它演示的版式。
type DemoPage struct {
	No     int    `json:"no"`
	Layout string `json:"layout"`
}

// StructureContractView GET /structure 的响应。
type StructureContractView struct {
	Layouts   []StructureLayout `json:"layouts"`
	DemoPages []DemoPage        `json:"demo_pages"`
	RulesMD   string            `json:"rules_md"`
}

// StructureContract 归属校验后的结构契约读取（面板数据源）。
func (s *Service) StructureContract(userID uint, id string) (*StructureContractView, error) {
	if _, err := s.GetOwned(userID, id); err != nil {
		return nil, err
	}
	if s.reg == nil {
		return nil, errors.New("注册表不可用")
	}
	tpl, err := s.reg.Get(id)
	if err != nil {
		return nil, fmt.Errorf("模板未挂载（最近一次校验未过）： %w", err)
	}
	return assembleContract(tpl), nil
}

// assembleContract 从挂载态模板组装契约视图（测试也要用，独立成纯函数）。
func assembleContract(tpl *template.Template) *StructureContractView {
	v := &StructureContractView{Layouts: make([]StructureLayout, 0, len(tpl.Layouts))}
	for _, l := range tpl.Layouts {
		skeleton, _ := tpl.Layout(l.ID)
		entry := StructureLayout{
			ID:          l.ID,
			Name:        l.Name,
			Use:         l.Use,
			Roles:       l.Roles,
			Constraints: l.Constraints,
			Pattern:     tpl.Pattern(l.ID),
			Repeats:     tpl.Repeats(l.ID),
			Skeleton:    skeleton,
		}
		v.Layouts = append(v.Layouts, entry)
	}
	v.DemoPages = demoPageMap(tpl.DemoHTML(""))
	v.RulesMD = tpl.Rules()
	return v
}

// demoPageMap 解析 demo HTML 里顶层 section 的 data-layout 序列（顺序即页码）。
// 契约规定只有顶层 section 带 data-layout，区间内平铺正则即可。
var demoLayoutRe = regexp.MustCompile(`data-layout="([^"]+)"`)

func demoPageMap(demoHTML string) []DemoPage {
	start := strings.Index(demoHTML, slidesStartMark)
	if start < 0 {
		return nil
	}
	seg := demoHTML[start:]
	if end := strings.Index(seg, slidesEndMark); end >= 0 {
		seg = seg[:end]
	}
	out := make([]DemoPage, 0, 8)
	for i, m := range demoLayoutRe.FindAllStringSubmatch(seg, -1) {
		out = append(out, DemoPage{No: i + 1, Layout: m[1]})
	}
	return out
}

// SLIDES 挂载标记（与 template 包的字面量一致；那份不导出）。
const (
	slidesStartMark = "<!-- SLIDES:START -->"
	slidesEndMark   = "<!-- SLIDES:END -->"
)

// ---------- 写路径：版式元数据 ----------

// LayoutMetaPatch 版式元数据的可选补丁（nil = 不改；与"清空"语义区分开）。
type LayoutMetaPatch struct {
	Name  *string   `json:"name,omitempty"`
	Use   *string   `json:"use,omitempty"`
	Roles *[]string `json:"roles,omitempty"`
}

// UpdateLayoutMeta 改一个版式的 name/use/roles：预检 → 落盘 template.json →
// best-effort 同步 layouts.md 文档行 → 记结构版本 → 重挂（失败降级为 warning，
// 盘上已是新值、历史可回滚——与 SaveStyleCSS 的"先落地后重挂"同一纪律）。
// 返回的 warning 非空时调用方应把它呈现给用户。
func (s *Service) UpdateLayoutMeta(userID uint, id, layoutID string, patch LayoutMetaPatch) (string, error) {
	row, err := s.GetOwned(userID, id)
	if err != nil {
		return "", err
	}
	if row.Status == "published" || row.Status == "publishing" {
		return "", ErrPublished
	}
	if err := patch.validate(); err != nil {
		return "", err
	}

	// 预检在锁外：layoutID 不存在/角色词非法时不动任何文件。
	// 用 map 保真读写（而非窄结构体）：template.json 里的其他字段原样保留。
	tjPath := filepath.Join(s.Dir(id), "template.json")
	raw, err := os.ReadFile(tjPath)
	if err != nil {
		return "", fmt.Errorf("读 template.json 失败: %w", err)
	}
	var meta map[string]any
	if err := json.Unmarshal(raw, &meta); err != nil {
		return "", fmt.Errorf("template.json 解析失败: %w", err)
	}
	layouts, _ := meta["layouts"].([]any)
	idx := -1
	for i, it := range layouts {
		if m, ok := it.(map[string]any); ok && m["id"] == layoutID {
			idx = i
			break
		}
	}
	if idx < 0 {
		return "", fmt.Errorf("版式 %q 未登记", layoutID)
	}
	entry, _ := layouts[idx].(map[string]any)
	if entry == nil {
		return "", fmt.Errorf("template.json 的 layouts[%d] 形状不对", idx)
	}

	changes := applyPatch(entry, patch)

	unlock := s.lockUT(id)
	defer unlock()

	out, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return "", err
	}
	if err := atomicWriteFile(tjPath, out); err != nil {
		return "", fmt.Errorf("写 template.json 失败: %w", err)
	}
	// layouts.md 的「适用 role：」行同步（文档行，找不到/写不了都不算失败）
	if patch.Roles != nil {
		roles, _ := entry["roles"].([]string)
		s.syncLayoutRolesMDQuiet(id, layoutID, roles)
	}
	_ = s.recordVersionUT(id, OpStruct, changes)
	if err := s.remountUT(id); err != nil {
		return fmt.Sprintf("已保存并记入历史，但注册表重挂未过：%v。生成侧暂用旧契约，可到历史回滚或用质量体检排查。", err), nil
	}
	return "", nil
}

// applyPatch 把补丁落到 template.json 的版式条目 map 上（patch 已过 validate），
// 返回给人看的变更摘要。
func applyPatch(entry map[string]any, patch LayoutMetaPatch) string {
	var parts []string
	if patch.Name != nil {
		name := strings.TrimSpace(*patch.Name)
		entry["name"] = name
		parts = append(parts, fmt.Sprintf("名称改为「%s」", name))
	}
	if patch.Use != nil {
		entry["use"] = strings.TrimSpace(*patch.Use)
		parts = append(parts, "用途已更新")
	}
	if patch.Roles != nil {
		roles := normalizeRoles(*patch.Roles)
		entry["roles"] = roles
		if len(roles) == 0 {
			parts = append(parts, "角色清空（不再限定适用场景）")
		} else {
			parts = append(parts, "角色改为 "+strings.Join(roles, "/"))
		}
	}
	display := layoutDisplay(entry)
	return fmt.Sprintf("版式「%s」：%s", display, strings.Join(parts, "；"))
}

// layoutDisplay 条目的展示名：name 优先，退回 id。
func layoutDisplay(entry map[string]any) string {
	if n, ok := entry["name"].(string); ok && n != "" {
		return n
	}
	if n, ok := entry["id"].(string); ok {
		return n
	}
	return "?"
}

// validate 预检：字段长度与角色词表。layoutID 的存在性在调用方查（要读文件）。
func (p LayoutMetaPatch) validate() error {
	if p.Name == nil && p.Use == nil && p.Roles == nil {
		return errors.New("没有要改的字段")
	}
	if p.Name != nil {
		n := len([]rune(strings.TrimSpace(*p.Name)))
		if n == 0 {
			return errors.New("版式名称不能为空")
		}
		if n > maxLayoutNameLen {
			return fmt.Errorf("版式名称 %d 字，超过 %d 上限", n, maxLayoutNameLen)
		}
	}
	if p.Use != nil && len([]rune(*p.Use)) > maxLayoutUseLen {
		return fmt.Errorf("用途 %d 字，超过 %d 上限", len([]rune(*p.Use)), maxLayoutUseLen)
	}
	if p.Roles != nil {
		if len(*p.Roles) > maxLayoutRoleCnt {
			return fmt.Errorf("角色最多 %d 个（太多会让 plan_pages 的选择失去意义）", maxLayoutRoleCnt)
		}
		for _, r := range *p.Roles {
			if !layoutRoles[r] {
				return fmt.Errorf("角色 %q 不在词表（可用：cover/toc/divider/content/data/quote/code/cta/thanks）", r)
			}
		}
	}
	return nil
}

// normalizeRoles 去重 + 按词表顺序排序（展示与提示词稳定）。
func normalizeRoles(in []string) []string {
	seen := map[string]bool{}
	for _, r := range in {
		seen[r] = true
	}
	out := make([]string, 0, len(seen))
	for _, r := range layoutRoleOrder {
		if seen[r] {
			out = append(out, r)
		}
	}
	return out
}

// syncLayoutRolesMDQuiet 同步 layouts.md 该条目的「适用 role：」行（纯文档：
// 注册表解析不消费它）。找不到条目或行就静默跳过；写失败只记日志——template.json
// 已落盘，这里不该让它把整次编辑变成失败。
func (s *Service) syncLayoutRolesMDQuiet(id, layoutID string, roles []string) {
	p := filepath.Join(s.Dir(id), "layouts.md")
	raw, err := os.ReadFile(p)
	if err != nil {
		return
	}
	out, changed := syncRolesLine(string(raw), layoutID, roles)
	if !changed {
		return
	}
	if err := atomicWriteFile(p, []byte(out)); err != nil {
		log.Printf("[warn] 模板 %s layouts.md 角色行同步失败: %v", id, err)
	}
}

// syncRolesLine 在 layouts.md 的 `## <layoutID>` 条目块内替换「适用 role：」行。
// 返回新全文与是否真的改动。块边界 = 下一个 "\n## " 或文件尾。
func syncRolesLine(md, layoutID string, roles []string) (string, bool) {
	head := "## " + layoutID
	start := strings.Index(md, head)
	if start < 0 {
		return md, false
	}
	end := len(md)
	if i := strings.Index(md[start+1:], "\n## "); i >= 0 {
		end = start + 1 + i
	}
	block := md[start:end]
	lineRe := regexp.MustCompile(`(?m)^适用 role[：:].*$`)
	if !lineRe.MatchString(block) {
		return md, false
	}
	val := strings.Join(roles, "、")
	if val == "" {
		val = "（无，不限场景）"
	}
	newBlock := lineRe.ReplaceAllString(block, "适用 role："+val+"。")
	if newBlock == block {
		return md, false
	}
	return md[:start] + newBlock + md[end:], true
}

// ---------- 层 3：对话加/删版式（骨架 + 类名 + CSS 的结构编辑） ----------
//
// 与层 2（元数据补丁）的本质差别：加/删动的是结构性文件，挂载校验失败的概率
// 不低（模型写的骨架类名可能不存在、双向登记可能不一致），而结构坏的模板会让
// 面板（读挂载快照）与生成侧脱节。所以写路径是"内存快照 → 落盘 → MountUser →
// 失败整体回滚再恢复挂载"，盘上永远不留坏状态；历史版本（五件套 bundle）是
// 第二重兜底。
//
// 数量契约（数量：class=N 行）v1 不开放：骨架类计数自动化会把 h2=1 这类
// 非并列元素误登记成契约，需要模型显式声明才有意义——留待手动编辑或后续迭代。

const (
	maxSkeletonBytes  = 128 << 10
	maxDemoHTMLBytes  = 128 << 10
	maxCSSAppendBytes = 64 << 10
	maxClasses        = 60
	maxConstraintsLen = 120
)

// layoutFingerprints 指纹词表（与 template 包 loadTemplate 的校验同源——
// 8 个视觉模式；add_layout 必填显式声明，不做骨架推断）。
var layoutFingerprints = map[string]bool{
	"hero": true, "stack": true, "cards": true, "split": true,
	"code": true, "table": true, "chart": true, "quote": true,
}

var (
	newLayoutIDRe = regexp.MustCompile(`^[a-z][a-z0-9_-]{2,39}$`)
	newClassRe    = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)
	// 骨架/demo 片段里的 class 属性（片段校验用）
	skeletonClassAttrRe = regexp.MustCompile(`class="([^"]*)"`)
)

// AddLayoutSpec add_layout 工具参数（对话是唯一入口）。CSSAppend 承载新类名的
// 样式规则——合法类名必须在 base∪模板文件里真实存在，而模板类来自 style.css
// 选择器与 index.html class token，所以"新类名"与"新 CSS"必须一次原子提交。
type AddLayoutSpec struct {
	LayoutID    string   `json:"layout_id"`
	Name        string   `json:"name"`
	Use         string   `json:"use"`
	Roles       []string `json:"roles"`
	Fingerprint string   `json:"fingerprint"`
	Classes     []string `json:"classes"`
	Skeleton    string   `json:"skeleton"`
	CSSAppend   string   `json:"css_append,omitempty"`
	DemoHTML    string   `json:"demo_html,omitempty"`
	Constraints string   `json:"constraints,omitempty"`
}

// normalize 去首尾空白、去空项、类名去重（保序）。
func (sp *AddLayoutSpec) normalize() {
	sp.LayoutID = strings.TrimSpace(sp.LayoutID)
	sp.Name = strings.TrimSpace(sp.Name)
	sp.Use = strings.TrimSpace(sp.Use)
	sp.Fingerprint = strings.TrimSpace(sp.Fingerprint)
	sp.Constraints = strings.TrimSpace(sp.Constraints)
	sp.Roles = normalizeRoles(sp.Roles)
	classes := sp.Classes[:0]
	seen := map[string]bool{}
	for _, c := range sp.Classes {
		c = strings.TrimSpace(c)
		if c == "" || seen[c] {
			continue
		}
		seen[c] = true
		classes = append(classes, c)
	}
	sp.Classes = classes
	sp.CSSAppend = strings.TrimSpace(sp.CSSAppend)
	sp.DemoHTML = strings.TrimSpace(sp.DemoHTML)
}

// validate 预检：命名/词表/长度/黑名单。类名是否"真实存在"（base∪模板文件）
// 不在这里查——MountUser 全量校验是权威，失败由回滚兜底。
func (sp AddLayoutSpec) validate() error {
	if !newLayoutIDRe.MatchString(sp.LayoutID) {
		return fmt.Errorf("版式 id %q 不合法：小写字母开头，3-40 位小写字母/数字/短横线", sp.LayoutID)
	}
	if sp.Name == "" {
		return errors.New("缺少版式名称")
	}
	if len([]rune(sp.Name)) > maxLayoutNameLen {
		return fmt.Errorf("版式名称超 %d 字", maxLayoutNameLen)
	}
	if sp.Use == "" {
		return errors.New("缺少用途描述（生成模型按它决定什么时候用这个版式）")
	}
	if len([]rune(sp.Use)) > maxLayoutUseLen {
		return fmt.Errorf("用途超 %d 字", maxLayoutUseLen)
	}
	if len(sp.Constraints) > maxConstraintsLen {
		return fmt.Errorf("内容约束超 %d 字", maxConstraintsLen)
	}
	if len(sp.Roles) > maxLayoutRoleCnt {
		return fmt.Errorf("角色最多 %d 个", maxLayoutRoleCnt)
	}
	if !layoutFingerprints[sp.Fingerprint] {
		return fmt.Errorf("指纹 %q 不在词表（hero/stack/cards/split/code/table/chart/quote）", sp.Fingerprint)
	}
	if len(sp.Classes) == 0 || len(sp.Classes) > maxClasses {
		return fmt.Errorf("合法类名需 1-%d 个", maxClasses)
	}
	for _, c := range sp.Classes {
		if !newClassRe.MatchString(c) {
			return fmt.Errorf("类名 %q 不合法（小写字母开头的字母/数字/短横线）", c)
		}
	}
	if sp.Skeleton == "" {
		return errors.New("缺少骨架代码")
	}
	if len(sp.Skeleton) > maxSkeletonBytes {
		return fmt.Errorf("骨架 %d 字节超上限 %d", len(sp.Skeleton), maxSkeletonBytes)
	}
	if len(sp.DemoHTML) > maxDemoHTMLBytes {
		return fmt.Errorf("demo_html %d 字节超上限 %d", len(sp.DemoHTML), maxDemoHTMLBytes)
	}
	if len(sp.CSSAppend) > maxCSSAppendBytes {
		return fmt.Errorf("css_append %d 字节超上限 %d", len(sp.CSSAppend), maxCSSAppendBytes)
	}
	if err := scanSkeleton(sp.Skeleton, sp.LayoutID, sp.Classes); err != nil {
		return fmt.Errorf("骨架校验未过: %w", err)
	}
	if sp.DemoHTML != "" {
		if err := scanSkeleton(sp.DemoHTML, sp.LayoutID, sp.Classes); err != nil {
			return fmt.Errorf("demo_html 校验未过: %w", err)
		}
	}
	if sp.CSSAppend != "" {
		if err := scanCSS(sp.CSSAppend); err != nil {
			return fmt.Errorf("css_append 校验未过: %w", err)
		}
	}
	return nil
}

// scanSkeleton 版式片段校验：单个顶层 section、data-layout 一致、危险内容黑名单、
// 用到的类必须都在声明清单里（给模型精准报错；类是否真实存在由 MountUser 终审）。
func scanSkeleton(fragment, layoutID string, classes []string) error {
	if !utf8.ValidString(fragment) {
		return errors.New("含非法 UTF-8 字节")
	}
	trim := strings.TrimSpace(fragment)
	if !strings.HasPrefix(trim, "<section") || !strings.HasSuffix(trim, "</section>") {
		return errors.New("必须是单个顶层 <section>…</section> 片段")
	}
	low := strings.ToLower(trim)
	for _, bad := range []string{"<script", "<iframe", "<object", "<embed", "javascript:"} {
		if strings.Contains(low, bad) {
			return fmt.Errorf("含被禁用的 %q", bad)
		}
	}
	for _, ev := range []string{" onload=", " onerror=", " onclick=", " onmouseover=", " onfocus="} {
		if strings.Contains(low, ev) {
			return fmt.Errorf("含内联事件 %q", strings.TrimSpace(ev))
		}
	}
	if strings.Count(trim, "<section") != 1 || strings.Count(trim, "</section>") != 1 {
		return errors.New("必须是单个顶层 <section>（骨架里不允许出现嵌套/额外的 section）")
	}
	m := demoLayoutRe.FindStringSubmatch(trim)
	if m == nil {
		return errors.New(`缺 data-layout 属性`)
	}
	if m[1] != layoutID {
		return fmt.Errorf("data-layout=%q 与版式 id %q 不一致", m[1], layoutID)
	}
	allowed := map[string]bool{}
	for _, c := range classes {
		allowed[c] = true
	}
	for _, attr := range skeletonClassAttrRe.FindAllStringSubmatch(trim, -1) {
		for _, tok := range strings.Fields(attr[1]) {
			if !newClassRe.MatchString(tok) {
				continue // 形状怪的 token 交给 MountUser 报
			}
			if !allowed[tok] {
				return fmt.Errorf("骨架用了类 .%s，但 classes 清单里没有（骨架里用到的每个类都必须声明）", tok)
			}
		}
	}
	return nil
}

// AddLayout 新增版式：预检 → 持锁 → 四文件落盘（css 追加/json 加条目/md 加条目/
// demo 追加页）→ MountUser → 失败整体回滚并恢复挂载。成功返回给模型的摘要。
func (s *Service) AddLayout(userID uint, id string, spec AddLayoutSpec) (string, error) {
	row, err := s.GetOwned(userID, id)
	if err != nil {
		return "", err
	}
	if row.Status == "published" || row.Status == "publishing" {
		return "", ErrPublished
	}
	// 词表校验在 normalize 之前：normalizeRoles 会静默丢弃非法词，
	// 模型写错必须收到明确报错而不是被悄悄清空。
	for _, r := range spec.Roles {
		if !layoutRoles[r] {
			return "", fmt.Errorf("角色 %q 不在词表（cover/toc/divider/content/data/quote/code/cta/thanks）", r)
		}
	}
	spec.normalize()
	if err := spec.validate(); err != nil {
		return "", err
	}

	unlock := s.lockUT(id)
	defer unlock()

	dir := s.Dir(id)
	paths := map[string]string{
		"template.json": filepath.Join(dir, "template.json"),
		"layouts.md":    filepath.Join(dir, "layouts.md"),
		"style.css":     filepath.Join(dir, "style.css"),
		"index.html":    filepath.Join(dir, "index.html"),
	}
	snap := map[string]string{}
	for name, p := range paths {
		raw, err := os.ReadFile(p)
		if err != nil {
			return "", fmt.Errorf("读 %s 失败: %w", name, err)
		}
		snap[name] = string(raw)
	}

	// template.json：重名检查 + 追加条目
	var meta map[string]any
	if err := json.Unmarshal([]byte(snap["template.json"]), &meta); err != nil {
		return "", fmt.Errorf("template.json 解析失败: %w", err)
	}
	layouts, _ := meta["layouts"].([]any)
	for _, it := range layouts {
		if m, ok := it.(map[string]any); ok && m["id"] == spec.LayoutID {
			return "", fmt.Errorf("版式 id %q 已存在", spec.LayoutID)
		}
	}
	entry := map[string]any{
		"id": spec.LayoutID, "name": spec.Name, "use": spec.Use, "roles": spec.Roles,
	}
	if spec.Constraints != "" {
		entry["constraints"] = spec.Constraints
	}
	meta["layouts"] = append(layouts, entry)

	// layouts.md：文件尾追加条目
	newMD := snap["layouts.md"] + "\n" + renderLayoutEntry(spec)
	// style.css：追加新类名规则
	newCSS := snap["style.css"]
	if spec.CSSAppend != "" {
		newCSS += "\n" + spec.CSSAppend + "\n"
	}
	// index.html：demo 示例页
	newIDX := snap["index.html"]
	demoNote := ""
	if spec.DemoHTML != "" {
		var ok bool
		if newIDX, ok = appendDemoSection(newIDX, spec.DemoHTML); !ok {
			return "", errors.New("index.html 缺 SLIDES 挂载标记，无法追加示例页")
		}
	} else {
		demoNote = "（未附 demo 示例页：预览与质量体检看不到它，建议补一页）"
	}

	restore := func() {
		_ = atomicWriteFile(paths["template.json"], []byte(snap["template.json"]))
		_ = atomicWriteFile(paths["layouts.md"], []byte(snap["layouts.md"]))
		_ = atomicWriteFile(paths["style.css"], []byte(snap["style.css"]))
		_ = atomicWriteFile(paths["index.html"], []byte(snap["index.html"]))
		_ = s.remountUT(id) // 回滚后恢复原挂载
	}
	if err := atomicWriteFile(paths["template.json"], marshalIndent(meta)); err != nil {
		restore()
		return "", fmt.Errorf("写 template.json 失败: %w", err)
	}
	if err := atomicWriteFile(paths["layouts.md"], []byte(newMD)); err != nil {
		restore()
		return "", fmt.Errorf("写 layouts.md 失败: %w", err)
	}
	if newCSS != snap["style.css"] {
		if err := atomicWriteFile(paths["style.css"], []byte(newCSS)); err != nil {
			restore()
			return "", fmt.Errorf("写 style.css 失败: %w", err)
		}
	}
	if newIDX != snap["index.html"] {
		if err := atomicWriteFile(paths["index.html"], []byte(newIDX)); err != nil {
			restore()
			return "", fmt.Errorf("写 index.html 失败: %w", err)
		}
	}
	if err := s.remountUT(id); err != nil {
		restore()
		return "", fmt.Errorf("注册表校验未过，已整体回滚：%v。按报错修正后重试", err)
	}
	summary := fmt.Sprintf("新增版式「%s」（%s）：角色 %s，指纹 %s。当前共 %d 个版式。",
		spec.Name, spec.LayoutID, roleListCN(spec.Roles), spec.Fingerprint, len(layouts)+1)
	if demoNote != "" {
		summary += " " + demoNote
	}
	_ = s.recordVersionUT(id, OpStruct, summary)
	return summary, nil
}

// RemoveLayout 删除版式：清理 template.json/layouts.md 条目与 demo 引用页，
// 删前做死锁健康检查（content 候选 <3 阻塞——R101 下 plan_pages 无合法解），
// 指纹降级与 cover/thanks 缺失降级为警告（就近退化可用）。
func (s *Service) RemoveLayout(userID uint, id, layoutID string) (string, error) {
	row, err := s.GetOwned(userID, id)
	if err != nil {
		return "", err
	}
	if row.Status == "published" || row.Status == "publishing" {
		return "", ErrPublished
	}

	unlock := s.lockUT(id)
	defer unlock()

	dir := s.Dir(id)
	paths := map[string]string{
		"template.json": filepath.Join(dir, "template.json"),
		"layouts.md":    filepath.Join(dir, "layouts.md"),
		"index.html":    filepath.Join(dir, "index.html"),
	}
	snap := map[string]string{}
	for name, p := range paths {
		raw, err := os.ReadFile(p)
		if err != nil {
			return "", fmt.Errorf("读 %s 失败: %w", name, err)
		}
		snap[name] = string(raw)
	}

	var meta map[string]any
	if err := json.Unmarshal([]byte(snap["template.json"]), &meta); err != nil {
		return "", fmt.Errorf("template.json 解析失败: %w", err)
	}
	layouts, _ := meta["layouts"].([]any)
	idx := -1
	for i, it := range layouts {
		if m, ok := it.(map[string]any); ok && m["id"] == layoutID {
			idx = i
		}
	}
	if idx < 0 {
		return "", fmt.Errorf("版式 %q 未登记", layoutID)
	}
	victim, _ := layouts[idx].(map[string]any)
	display := layoutDisplay(victim)

	remaining := make([]map[string]any, 0, len(layouts)-1)
	kept := make([]any, 0, len(layouts)-1)
	for i, it := range layouts {
		if i == idx {
			continue
		}
		kept = append(kept, it)
		if m, ok := it.(map[string]any); ok {
			remaining = append(remaining, m)
		}
	}
	if len(remaining) == 0 {
		return "", errors.New("这是最后一个版式，删除后模板无法生成任何页面")
	}
	contentCount, hasCover, hasThanks := 0, false, false
	for _, m := range remaining {
		for _, r := range entryRoles(m) {
			switch r {
			case "content":
				contentCount++
			case "cover":
				hasCover = true
			case "thanks":
				hasThanks = true
			}
		}
	}
	if contentCount < 3 {
		return "", fmt.Errorf("删掉「%s」后正文（content）候选只剩 %d 个：5 页以上的 deck 会因同版式连续超限在 plan_pages 死锁。先加一个新版式再来删", display, contentCount)
	}
	patterns := layoutMDPatterns(snap["layouts.md"])
	warnings := removalWarnings(remaining, patterns, hasCover, hasThanks)

	meta["layouts"] = kept
	newTJ := marshalIndent(meta)
	newMD, mdOK := removeLayoutMD(snap["layouts.md"], layoutID)
	newIDX, removed := removeDemoSections(snap["index.html"], layoutID)

	restore := func() {
		_ = atomicWriteFile(paths["template.json"], []byte(snap["template.json"]))
		_ = atomicWriteFile(paths["layouts.md"], []byte(snap["layouts.md"]))
		_ = atomicWriteFile(paths["index.html"], []byte(snap["index.html"]))
		_ = s.remountUT(id)
	}
	if err := atomicWriteFile(paths["template.json"], newTJ); err != nil {
		restore()
		return "", fmt.Errorf("写 template.json 失败: %w", err)
	}
	if mdOK {
		if err := atomicWriteFile(paths["layouts.md"], []byte(newMD)); err != nil {
			restore()
			return "", fmt.Errorf("写 layouts.md 失败: %w", err)
		}
	}
	if removed > 0 {
		if err := atomicWriteFile(paths["index.html"], []byte(newIDX)); err != nil {
			restore()
			return "", fmt.Errorf("写 index.html 失败: %w", err)
		}
	}
	if err := s.remountUT(id); err != nil {
		restore()
		return "", fmt.Errorf("注册表校验未过，已整体回滚：%v", err)
	}
	summary := fmt.Sprintf("删除版式「%s」（%s），demo 引用页清理 %d 页。当前共 %d 个版式。",
		display, layoutID, removed, len(kept))
	if len(warnings) > 0 {
		summary += " 注意：" + strings.Join(warnings, "；")
	}
	_ = s.recordVersionUT(id, OpStruct, summary)
	return summary, nil
}

// entryRoles 读 template.json 条目的 roles（any 形状安全解包）。
func entryRoles(entry map[string]any) []string {
	raw, _ := entry["roles"].([]any)
	out := make([]string, 0, len(raw))
	for _, r := range raw {
		if s, ok := r.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// layoutMDPatterns 从 layouts.md 提取 版式 id → 指纹（删除警告用；缺指纹行的不计入）。
var mdPatternRe = regexp.MustCompile(`(?m)^指纹[：:]\s*(\S+)\s*$`)

func layoutMDPatterns(md string) map[string]string {
	out := map[string]string{}
	head := regexp.MustCompile(`(?m)^##\s+([A-Za-z][A-Za-z0-9_-]*)`)
	matches := head.FindAllStringSubmatchIndex(md, -1)
	for i, h := range matches {
		id := md[h[2]:h[3]]
		end := len(md)
		if i+1 < len(matches) {
			end = matches[i+1][2]
		}
		if m := mdPatternRe.FindStringSubmatch(md[h[0]:end]); m != nil {
			out[id] = m[1]
		}
	}
	return out
}

// removalWarnings 删除后的非阻塞降级提示。
func removalWarnings(remaining []map[string]any, patterns map[string]string, hasCover, hasThanks bool) []string {
	var out []string
	distinct := map[string]bool{}
	for _, m := range remaining {
		if fp := patterns[fmt.Sprint(m["id"])]; fp != "" && layoutFingerprints[fp] {
			distinct[fp] = true
		}
	}
	if len(distinct) > 0 && len(distinct) < 4 {
		out = append(out, fmt.Sprintf("视觉指纹只剩 %d 种（<4）：假多样性检查会停用，注意别让连续页面长得一样", len(distinct)))
	}
	if !hasCover {
		out = append(out, "没有任何版式声明 cover 角色：封面页会按用途就近退化，建议给某个满版版式勾上封面")
	}
	if !hasThanks {
		out = append(out, "没有任何版式声明 thanks 角色：收尾页会就近退化")
	}
	return out
}

// roleListCN 角色清单的可读形态（报错与摘要用）。
func roleListCN(roles []string) string {
	if len(roles) == 0 {
		return "无（不限场景）"
	}
	return strings.Join(roles, "/")
}

// renderLayoutEntry 把新条目渲染成 layouts.md 的文本块（格式与脚手架条目一致）。
func renderLayoutEntry(spec AddLayoutSpec) string {
	var b strings.Builder
	fmt.Fprintf(&b, "## %s（%s）\n", spec.LayoutID, spec.Name)
	fmt.Fprintf(&b, "指纹：%s\n\n", spec.Fingerprint)
	fmt.Fprintf(&b, "用途：%s\n", spec.Use)
	fmt.Fprintf(&b, "适用 role：%s。\n", roleListMD(spec.Roles))
	if spec.Constraints != "" {
		fmt.Fprintf(&b, "内容约束：%s\n", spec.Constraints)
	}
	fmt.Fprintf(&b, "\n合法类名：%s\n\n", strings.Join(spec.Classes, ", "))
	fmt.Fprintf(&b, "```html\n%s\n```\n", strings.TrimSpace(spec.Skeleton))
	return b.String()
}

// roleListMD layouts.md 文档行的角色形态。
func roleListMD(roles []string) string {
	if len(roles) == 0 {
		return "（无，不限场景）"
	}
	return strings.Join(roles, "、")
}

// appendDemoSection 把示例 section 插到 SLIDES:END 标记前。
func appendDemoSection(idxHTML, demoHTML string) (string, bool) {
	end := strings.Index(idxHTML, slidesEndMark)
	if end < 0 {
		return idxHTML, false
	}
	return idxHTML[:end] + strings.TrimSpace(demoHTML) + "\n" + idxHTML[end:], true
}

// removeDemoSections 删除 SLIDES 区间内所有 data-layout == layoutID 的顶层
// section（契约：不嵌套，首个 </section> 即边界）。返回新全文与删除数。
func removeDemoSections(idxHTML, layoutID string) (string, int) {
	start := strings.Index(idxHTML, slidesStartMark)
	if start < 0 {
		return idxHTML, 0
	}
	relEnd := strings.Index(idxHTML[start:], slidesEndMark)
	if relEnd < 0 {
		return idxHTML, 0
	}
	segStart := start + len(slidesStartMark)
	segEnd := start + relEnd
	seg := idxHTML[segStart:segEnd]

	var b strings.Builder
	b.WriteString(idxHTML[:segStart])
	removed := 0
	rest := seg
	for {
		i := strings.Index(rest, "<section")
		if i < 0 {
			b.WriteString(rest)
			break
		}
		relJ := strings.Index(rest[i:], "</section>")
		if relJ < 0 {
			b.WriteString(rest)
			break
		}
		secEnd := i + relJ + len("</section>")
		sec := rest[i:secEnd]
		if m := demoLayoutRe.FindStringSubmatch(sec); m != nil && m[1] == layoutID {
			removed++
		} else {
			b.WriteString(rest[:secEnd])
		}
		rest = rest[secEnd:]
	}
	b.WriteString(idxHTML[segEnd:])
	return b.String(), removed
}

// removeLayoutMD 从 layouts.md 删除该条目块（多吃的空行按一个换行归还）。
func removeLayoutMD(md, layoutID string) (string, bool) {
	head := "## " + layoutID
	start := strings.Index(md, head)
	if start < 0 {
		return md, false
	}
	end := len(md)
	if i := strings.Index(md[start+1:], "\n## "); i >= 0 {
		end = start + 1 + i
	}
	cutFrom := start
	for cutFrom > 0 && (md[cutFrom-1] == '\n' || md[cutFrom-1] == '\r') {
		cutFrom--
	}
	return md[:cutFrom] + "\n" + md[end:], true
}

// marshalIndent JSON 缩进输出（与模板既有格式一致）。
func marshalIndent(v any) []byte {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return []byte("{}")
	}
	return out
}
