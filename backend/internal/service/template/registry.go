// Package template 是 deck-v2 模板体系的注册表与加载器。
//
// 模板是目录级自包含单元（契约见 docs/refactor-plan.md §5.1）：
//
//	templates/<id>/
//	  template.json   机器可读元数据
//	  index.html      自包含骨架（demo 数据 + SLIDES:START/END 挂载标记）
//	  style.css       .tpl-<id> 作用域隔离的完整设计系统（含 variants class）
//	  layouts.md      版式登记簿：LLM 的唯一版式契约
//	  rules.md        模板专属质量规则（生成阶段注入）
//
// 铁律：校验失败的模板拒绝注册并打日志——半成品模板上线会让生成管线
// 拿到坏契约，而那类失败（"模型写不出合法页面"）排查成本远高于启动失败。
package template

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
)

// Canvas 画布尺寸（设计像素，deck 按 transform scale 缩放显示）。
type Canvas struct {
	W int `json:"w"`
	H int `json:"h"`
}

// Variant 主题变体槽：class 为 style.css 里预定义的 token 覆盖 class，
// 空串 = 模板默认观感。换变体 = body 挂 class，agent 不参与。
type Variant struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Class string `json:"class"`
}

// LayoutMeta 单个版式的登记信息。
type LayoutMeta struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Use         string   `json:"use"`
	Roles       []string `json:"roles,omitempty"`    // 适用的大纲 role（cover/toc/...）
	Constraints string   `json:"constraints,omitempty"`
}

// Source 模板来源（License 归属必须可追溯）。
type Source struct {
	DerivedFrom string `json:"derived_from"`
	License     string `json:"license"`
}

// Meta template.json 的结构，也是画廊列表接口的返回形状。
type Meta struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Tags        []string     `json:"tags,omitempty"`
	Scenario    []string     `json:"scenario,omitempty"`
	Canvas      Canvas       `json:"canvas"`
	Variants    []Variant    `json:"variants"`
	Layouts     []LayoutMeta `json:"layouts"`
	Fonts       []string     `json:"fonts,omitempty"`
	Source      Source       `json:"source"`
}

// Template 是一个已通过校验的模板：元数据 + 文件内容 + 类名清单。
type Template struct {
	Meta
	dir       string
	assetsDir string // web/assets 绝对路径：资产引用核验与 base.css 类清单的数据源
	indexHTML string
	styleCSS  string
	rulesMD   string

	// layoutSkeleton 版式 id → layouts.md 里该条目的 ```html 骨架代码
	//（read_layout 工具的原材料；LLM 先取骨架再填内容）。
	layoutSkeleton map[string]string
	// layoutClasses 版式 id → 该条目"合法类名："行声明的类名清单。
	// 每页校验的权威清单（C202），base/template 两族是兜底并集。
	layoutClasses map[string][]string
	// layoutPatterns 版式 id → 视觉模式指纹（hero/stack/cards/...）。
	// 显式声明（指纹：行）优先，缺失时按骨架类名推断。节奏守卫按它判视觉重复。
	layoutPatterns map[string]string

	// baseClasses / templateClasses 类名清单（manifest）：
	// base 从 deck-v2/base.css 解析，template 从本模板 style.css + index.html 解析。
	baseClasses     map[string]bool
	templateClasses map[string]bool
}

// HasLayout 版式是否登记。
func (t *Template) HasLayout(id string) bool {
	_, ok := t.layoutSkeleton[id]
	return ok
}

// Layout 返回版式骨架代码（```html 围栏内的内容）。
func (t *Template) Layout(id string) (string, bool) {
	s, ok := t.layoutSkeleton[id]
	return s, ok
}

// LayoutClassList 返回版式的合法类名清单（切片形式，给报错信息用）。
func (t *Template) LayoutClassList(id string) []string {
	return t.layoutClasses[id]
}

// Pattern 返回版式的视觉模式指纹；未登记的版式返回空串。
func (t *Template) Pattern(id string) string {
	return t.layoutPatterns[id]
}

// DistinctPatterns 模板登记版式的不同视觉模式数（节奏守卫判断该模板
// 有没有足够的模式多样性可判——不足 4 种的模板按模式判重会无解）。
func (t *Template) DistinctPatterns() int {
	seen := map[string]bool{}
	for _, p := range t.layoutPatterns {
		seen[p] = true
	}
	return len(seen)
}

// AllowedClasses 版式的合法类名全集：base 原语 ∪ 模板类 ∪ 该版式显式声明。
// 返回 map 供热路径上的逐类查询。
func (t *Template) AllowedClasses(layoutID string) map[string]bool {
	allowed := make(map[string]bool, len(t.baseClasses)+len(t.templateClasses)+8)
	for c := range t.baseClasses {
		allowed[c] = true
	}
	for c := range t.templateClasses {
		allowed[c] = true
	}
	for _, c := range t.layoutClasses[layoutID] {
		allowed[c] = true
	}
	return allowed
}

// Rules 模板专属质量规则全文（generate 阶段注入）。
func (t *Template) Rules() string { return t.rulesMD }

// DemoHTML 返回 demo 页 HTML，供画廊预览 iframe 使用。variant 非空时把
// 变体 class 追加到 body（服务端换肤——前端切变体只换一个 query 参数）。
// demo 里相对引用的 style.css 由调用方改写为绝对路径（见 handler.PreviewTemplate）。
func (t *Template) DemoHTML(variant string) string {
	html := t.indexHTML
	if variant == "" {
		return html
	}
	cls, ok := t.VariantClass(variant)
	if !ok || cls == "" || strings.Contains(html, cls) {
		return html
	}
	return bodyClassRe.ReplaceAllStringFunc(html, func(m string) string {
		sub := bodyClassRe.FindStringSubmatch(m)
		return `<body class="` + sub[1] + " " + cls + `">`
	})
}

var bodyClassRe = regexp.MustCompile(`<body class="([^"]*)">`)

// VariantClass 返回变体 id 对应的 class（未找到时返回 ""，调用方决定是否报错）。
func (t *Template) VariantClass(variantID string) (string, bool) {
	for _, v := range t.Variants {
		if v.ID == variantID {
			return v.Class, true
		}
	}
	return "", false
}

// ---------- 类名解析 ----------

var cssClassRe = regexp.MustCompile(`\.([A-Za-z][A-Za-z0-9_-]*)`)
var htmlClassAttrRe = regexp.MustCompile(`class="([^"]*)"`)

// collectClasses 从 CSS/HTML 文本收集出现过的类名 token。
// CSS 侧取 `.name` 选择器；HTML 侧取 class="..." 拆词。
func collectClasses(texts ...string) map[string]bool {
	set := make(map[string]bool)
	for _, text := range texts {
		for _, m := range cssClassRe.FindAllStringSubmatch(text, -1) {
			set[m[1]] = true
		}
		for _, m := range htmlClassAttrRe.FindAllStringSubmatch(text, -1) {
			for _, c := range strings.Fields(m[1]) {
				if isClassToken(c) {
					set[c] = true
				}
			}
		}
	}
	return set
}

func isClassToken(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_') {
			return false
		}
	}
	return true
}

// ---------- 注册表 ----------

// Registry 已注册模板的只读集合。内置层启动时一次加载（运行期只读——
// 模板是代码库的一部分，进 git，不热载）；用户层（自定义模板）按发布/
// 下架增量 Mount/Unmount。
type Registry struct {
	mu        sync.RWMutex
	items     map[string]*Template
	builtins  map[string]bool // 内置层 id 集合（用户层不可覆盖内置 id）
	assetsDir string          // 用户层加载需要（内置层在 New 时已用过）
	builtinsRoot string       // 内置模板根目录（fork 复制来源）
	// user 用户自定义模板层（MountUser 注册，Get/List 与内置层合并视图）
	user map[string]*Template
	// pendingUser ValidateUserDir → MountUser 之间的校验结果缓存（dir → t）
	pendingUser map[string]*Template
}

// NewRegistry 扫描 dir 下的模板目录并逐个校验。
// assetsDir 用于解析 base.css 的原语类（deck-v2/base.css）。
// 一个模板校验失败不影响其它模板注册，但会返回聚合错误（启动日志要能看到全貌）。
func NewRegistry(dir, assetsDir string) (*Registry, error) {
	baseCSS, err := os.ReadFile(filepath.Join(assetsDir, "deck-v2", "base.css"))
	if err != nil {
		return nil, fmt.Errorf("读取 base.css（类名清单的数据源）失败: %w", err)
	}
	baseClasses := collectClasses(string(baseCSS))

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("扫描模板目录 %s: %w", dir, err)
	}

	r := &Registry{items: make(map[string]*Template), builtins: map[string]bool{}, user: make(map[string]*Template), pendingUser: make(map[string]*Template), assetsDir: assetsDir, builtinsRoot: dir}
	var errs []string
	for _, e := range entries {
		if !e.IsDir() || e.Name() == "tools" {
			continue
		}
		tdir := filepath.Join(dir, e.Name())
		if _, err := os.Stat(filepath.Join(tdir, "template.json")); err != nil {
			continue // 没有 template.json 的目录不是模板，静默跳过
		}
		t, err := loadTemplate(tdir, assetsDir, baseClasses)
		if err != nil {
			errs = append(errs, fmt.Sprintf("模板 %s: %v", e.Name(), err))
			continue
		}
		r.items[t.ID] = t
		r.builtins[t.ID] = true
	}
	if len(errs) > 0 {
		return r, fmt.Errorf("%d 个模板校验失败（已跳过）:\n%s", len(errs), strings.Join(errs, "\n"))
	}
	return r, nil
}

// MountUser 注册一个用户自定义模板（发布时调用；重发布 = 幂等覆盖）。
// 校验用与内置层完全同一套 loadTemplate 规则——不过关的模板进不来。
func (r *Registry) MountUser(dir string) error {
	if err := r.ValidateUserDir(dir); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.pendingUser[dir]
	if !ok {
		return fmt.Errorf("内部错误：校验缓存缺失，请重试")
	}
	if r.builtins[t.ID] {
		return fmt.Errorf("模板 id %q 与内置模板冲突", t.ID)
	}
	delete(r.pendingUser, dir)
	r.user[t.ID] = t
	r.items[t.ID] = t
	return nil
}

// ValidateUserDir 按内置层同一套规则校验一个用户模板目录（不注册）。
// 发布门禁的第一关用它；校验结果缓存一份供紧随其后的 MountUser 复用
//（避免同一目录连跑两遍完整校验）。
func (r *Registry) ValidateUserDir(dir string) error {
	baseCSS, err := os.ReadFile(filepath.Join(r.assetsDir, "deck-v2", "base.css"))
	if err != nil {
		return err
	}
	t, err := loadTemplate(dir, r.assetsDir, collectClasses(string(baseCSS)))
	if err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.builtins[t.ID] {
		return fmt.Errorf("模板 id %q 与内置模板冲突", t.ID)
	}
	if r.pendingUser == nil {
		r.pendingUser = map[string]*Template{}
	}
	r.pendingUser[dir] = t
	return nil
}

// UnmountUser 注销用户模板（下架/删除时调用；未注册时静默）。
func (r *Registry) UnmountUser(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.builtins[id] {
		delete(r.user, id)
		delete(r.items, id)
	}
}

// BuiltinDir 内置模板的源目录（fork 复制来源）。
func (r *Registry) BuiltinDir(id string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if !r.builtins[id] {
		return "", fmt.Errorf("内置模板 %q 不存在", id)
	}
	return filepath.Join(r.builtinsRoot, id), nil
}

// List 返回全部模板元数据（内置 + 已挂载的用户层，按 id 排序，画廊列表稳定）。
func (r *Registry) List() []*Meta {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*Meta, 0, len(r.items))
	for _, t := range r.items {
		m := t.Meta
		out = append(out, &m)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// Get 按 id 取模板。未注册的 id 与"不存在"同一种错误（不区分，避免信息泄露顾虑）。
func (r *Registry) Get(id string) (*Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.items[id]
	if !ok {
		return nil, fmt.Errorf("模板 %q 不存在", id)
	}
	return t, nil
}

// Count 已注册模板数（启动日志用）。
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.items)
}

// ---------- 加载与校验 ----------

func loadTemplate(dir, assetsDir string, baseClasses map[string]bool) (*Template, error) {
	name := filepath.Base(dir)

	read := func(f string) (string, error) {
		b, err := os.ReadFile(filepath.Join(dir, f))
		if err != nil {
			return "", fmt.Errorf("缺少或读不了 %s: %w", f, err)
		}
		return string(b), nil
	}

	metaRaw, err := read("template.json")
	if err != nil {
		return nil, err
	}
	var meta Meta
	if err := json.Unmarshal([]byte(metaRaw), &meta); err != nil {
		return nil, fmt.Errorf("template.json 解析失败: %w", err)
	}
	if meta.ID == "" || meta.Name == "" {
		return nil, fmt.Errorf("template.json 缺 id/name")
	}
	if meta.ID != name {
		return nil, fmt.Errorf("template.json id %q 与目录名 %q 不一致", meta.ID, name)
	}
	if meta.Canvas.W <= 0 || meta.Canvas.H <= 0 {
		return nil, fmt.Errorf("canvas 尺寸非法: %+v", meta.Canvas)
	}
	if len(meta.Layouts) == 0 {
		return nil, fmt.Errorf("layouts 为空：没有版式契约的模板无法生成")
	}
	if len(meta.Variants) == 0 {
		return nil, fmt.Errorf("variants 为空：至少要有 default 变体")
	}

	indexHTML, err := read("index.html")
	if err != nil {
		return nil, err
	}
	styleCSS, err := read("style.css")
	if err != nil {
		return nil, err
	}
	rulesMD, err := read("rules.md")
	if err != nil {
		return nil, err
	}
	layoutMD, err := read("layouts.md")
	if err != nil {
		return nil, err
	}

	t := &Template{
		Meta:            meta,
		dir:             dir,
		assetsDir:       assetsDir,
		indexHTML:       indexHTML,
		styleCSS:        styleCSS,
		rulesMD:         rulesMD,
		layoutSkeleton:  map[string]string{},
		layoutClasses:   map[string][]string{},
		layoutPatterns:  map[string]string{},
		baseClasses:     baseClasses,
		templateClasses: collectClasses(indexHTML, styleCSS),
	}

	// 挂载标记：没有它实例化就无处剥离 demo、无处插入正式页面。
	if !strings.Contains(indexHTML, slidesStartMarker) || !strings.Contains(indexHTML, slidesEndMarker) {
		return nil, fmt.Errorf("index.html 缺 %s / %s 挂载标记", slidesStartMarker, slidesEndMarker)
	}

	// 资产引用核验：index.html 里引用的 /assets/* 必须真实存在。
	// 这是"实例化出的 deck 一定能渲染"的前提，断链在这里拦比在线上拦便宜一个量级。
	if err := t.checkAssetRefs(indexHTML); err != nil {
		return nil, err
	}

	// layouts.md 解析：条目 id 集合必须与 template.json 一致（双向），
	// 每个条目必须有"合法类名"行，且声明的类名在 base ∪ 模板 里真实存在。
	layouts := parseLayoutsMD(layoutMD)
	declared := map[string]bool{}
	for _, l := range meta.Layouts {
		declared[l.ID] = true
	}
	for id := range declared {
		if _, ok := layouts[id]; !ok {
			return nil, fmt.Errorf("template.json 登记的版式 %q 在 layouts.md 无条目", id)
		}
	}
	for id, l := range layouts {
		if !declared[id] {
			return nil, fmt.Errorf("layouts.md 条目 %q 未在 template.json 登记", id)
		}
		if len(l.Classes) == 0 {
			return nil, fmt.Errorf("版式 %q 缺「合法类名：」行", id)
		}
		if l.Skeleton == "" {
			return nil, fmt.Errorf("版式 %q 缺 ```html 骨架代码", id)
		}
		for _, c := range l.Classes {
			if !baseClasses[c] && !t.templateClasses[c] {
				return nil, fmt.Errorf("版式 %q 声明的类名 .%s 不存在（base.css 与模板文件里都没有）", id, c)
			}
		}
		// 骨架自洽：骨架里出现的类必须在「合法类名」清单里，否则 LLM 照抄骨架
		// 写页会被 C202 拒收，而报错指向的是模板自己的 bug（weekly-report
		// metrics-chart 用 .c/.l、CSS 却是 .b/.lbl 的实测教训）。
		for _, cls := range htmlClassAttrRe.FindAllStringSubmatch(l.Skeleton, -1) {
			for _, c := range strings.Fields(cls[1]) {
				if !isClassToken(c) {
					continue
				}
				ok := false
				for _, d := range l.Classes {
					if d == c {
						ok = true
						break
					}
				}
				if !ok {
					return nil, fmt.Errorf("版式 %q 骨架用了类 .%s，但「合法类名：」清单里没有", id, c)
				}
			}
		}
		// 指纹：显式声明必须在词汇表内；缺失时由 patternOf 从骨架推断（兼容旧模板）。
		if l.Pattern != "" {
			switch l.Pattern {
			case "hero", "stack", "cards", "split", "code", "table", "chart", "quote":
			default:
				return nil, fmt.Errorf("版式 %q 的指纹 %q 不在词汇表（hero/stack/cards/split/code/table/chart/quote）", id, l.Pattern)
			}
		}
		t.layoutSkeleton[id] = l.Skeleton
		t.layoutClasses[id] = l.Classes
		t.layoutPatterns[id] = patternOf(l)
	}

	// 网格陷阱：裸 .slide 无条件 grid 会把不含 .sidebar 包裹层的骨架碎片逐格
	// 压进窄列（course-module deck-0031 十二页十坏的根源）。必须 :has(.sidebar) 触发。
	if gridTrapRe.MatchString(t.styleCSS) {
		return nil, fmt.Errorf("style.css 在裸 .slide 上无条件开 grid（且没有 :has(.sidebar) 守卫）")
	}

	// 变体 class 必须在 style.css 里真的定义过（用户选了变体却没生效，是最阴的静默失败）。
	for _, v := range meta.Variants {
		if v.ID == "" {
			return nil, fmt.Errorf("variants 存在空 id 的变体")
		}
		if v.Class != "" && !t.templateClasses[v.Class] {
			return nil, fmt.Errorf("变体 %q 的 class .%s 在 style.css 中无定义", v.ID, v.Class)
		}
	}

	// demo 页面（SLIDES 区间内）的每个 section 都要带已登记的 data-layout：
	// demo 是画廊预览与"照着写"的范本，范本自己违反契约就没有说服力。
	seg := demoSegment(indexHTML)
	for _, m := range sectionTagRe.FindAllStringSubmatch(seg, -1) {
		dm := dataLayoutRe.FindStringSubmatch(m[0])
		if dm == nil {
			return nil, fmt.Errorf("demo section 缺 data-layout: %.60s", m[0])
		}
		if !declared[dm[1]] {
			return nil, fmt.Errorf("demo section 的 data-layout=%q 未登记", dm[1])
		}
	}

	return t, nil
}

var (
	slidesStartMarker = "<!-- SLIDES:START -->"
	slidesEndMarker   = "<!-- SLIDES:END -->"

	sectionTagRe = regexp.MustCompile(`<section[^>]*>`)
	dataLayoutRe = regexp.MustCompile(`data-layout="([^"]*)"`)
	assetRefRe   = regexp.MustCompile(`(?:href|src)="(\/assets\/[^"]+)"`)
	// gridTrapRe 裸 .slide 上的无条件 grid（负向豁免 :has 守卫在更早的 if 里判）
	gridTrapRe = regexp.MustCompile(`\.tpl-[\w-]+ \.slide\{[^}]*display:\s*grid`)
)

// demoSegment 取挂载标记之间的内容（demo 页面区）。
func demoSegment(indexHTML string) string {
	start := strings.Index(indexHTML, slidesStartMarker)
	end := strings.Index(indexHTML, slidesEndMarker)
	if start < 0 || end < 0 || end < start {
		return ""
	}
	return indexHTML[start+len(slidesStartMarker) : end]
}

// checkAssetRefs 核验 index.html 引用的 /assets/* 都存在于磁盘。
func (t *Template) checkAssetRefs(indexHTML string) error {
	for _, m := range assetRefRe.FindAllStringSubmatch(indexHTML, -1) {
		// 引用形如 /assets/deck-v2/xxx；映射到 assetsDir 下的真实文件。
		rel := strings.TrimPrefix(m[1], "/assets/")
		p := filepath.Join(t.assetsDir, filepath.FromSlash(rel))
		if _, err := os.Stat(p); err != nil {
			return fmt.Errorf("index.html 引用不存在的资产 %s", m[1])
		}
	}
	return nil
}
