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
