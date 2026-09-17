package deck

// deck-v2 的存储层：deck.json（元数据/阶段/页计划）、outline.json（结构化大纲）、
// 阶段状态机、模板实例化。
//
// 与 v1 的关系：v1 的代码（deck.html / theme / preset）原样保留到 P5 摘除，
// 两套格式靠 Deck.Format 列区分；所有 v2 方法都走 index.html + SLIDES 挂载标记。

import (
	"encoding/json"
	"regexp"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"html-ppt/backend/internal/service/template"
	"html-ppt/backend/internal/store"
)

// 生成流程阶段（权威在 DB Deck.Stage，deck.json 是冗余副本，写路径以 DB 为准）。
// 迁移守卫见 transitionStage。
const (
	StageDraft             = "draft"
	StageOutlining         = "outlining"
	StageOutlineReview     = "outline_review"
	StageSelectingTemplate = "selecting_template"
	StageGenerating        = "generating"
	StageIterating         = "iterating"
)

// FormatV2 deck 格式标记（Deck.Format 列 + deck.json.format）。
const FormatV2 = "v2"

// OutlineRole 大纲页的语义角色（决定 plan_pages 的版式选择范围）。
const (
	RoleCover   = "cover"
	RoleTOC     = "toc"
	RoleDivider = "divider"
	RoleContent = "content"
	RoleData    = "data"
	RoleQuote   = "quote"
	RoleCode    = "code"
	RoleCTA     = "cta"
	RoleThanks  = "thanks"
)

// ---------- outline.json ----------

// Outline 结构化大纲（D5：一等公民，双通道编辑都写这份 JSON）。
type Outline struct {
	Version   int          `json:"version"`
	Title     string       `json:"title"`
	Meta      OutlineMeta  `json:"meta,omitempty"`
	Narrative Narrative    `json:"narrative,omitempty"`
	Pages     []OutlinePage `json:"pages"`
}

type OutlineMeta struct {
	Audience    string `json:"audience,omitempty"`
	DurationMin int    `json:"duration_min,omitempty"`
	PageCount   int    `json:"page_count,omitempty"`
	Tone        string `json:"tone,omitempty"`
}

type Narrative struct {
	Hook string   `json:"hook,omitempty"`
	Arcs []string `json:"arcs,omitempty"`
}

type OutlineMaterial struct {
	Type string `json:"type,omitempty"` // image / data / screenshot / ...
	Desc string `json:"desc"`
}

type OutlinePage struct {
	No         int               `json:"no"`
	Role       string            `json:"role"`
	Title      string            `json:"title"`
	Points     []string          `json:"points,omitempty"`
	LayoutHint string            `json:"layout_hint,omitempty"`
	Materials  []OutlineMaterial `json:"materials,omitempty"`
	Notes      string            `json:"notes,omitempty"`
}

var outlineRoles = map[string]bool{
	RoleCover: true, RoleTOC: true, RoleDivider: true, RoleContent: true,
	RoleData: true, RoleQuote: true, RoleCode: true, RoleCTA: true, RoleThanks: true,
}

// Validate 大纲的 schema 校验（双通道写路径共用：面板直改与 agent 修改同一套规则）。
func (o *Outline) Validate() error {
	if strings.TrimSpace(o.Title) == "" {
		return fmt.Errorf("大纲缺少 title")
	}
	if len(o.Pages) == 0 {
		return fmt.Errorf("大纲没有任何页面")
	}
	if len(o.Pages) > 60 {
		return fmt.Errorf("页数 %d 超出上限 60：拆成多份或精简", len(o.Pages))
	}
	for i, p := range o.Pages {
		if p.No != i+1 {
			return fmt.Errorf("第 %d 页的 no=%d：页码必须从 1 连续递增", i+1, p.No)
		}
		if !outlineRoles[p.Role] {
			return fmt.Errorf("第 %d 页 role=%q 不合法（cover/toc/divider/content/data/quote/code/cta/thanks）", p.No, p.Role)
		}
		if strings.TrimSpace(p.Title) == "" {
			return fmt.Errorf("第 %d 页缺少标题", p.No)
		}
		if len(p.Points) > 6 {
			return fmt.Errorf("第 %d 页要点 %d 条，上限 6：合并或下沉到讲稿", p.No, len(p.Points))
		}
	}
	return nil
}

// ---------- deck.json ----------

// PlanAssignment plan_pages 的产物：页 → 版式的全局分配（生成前一次性规划）。
type PlanAssignment struct {
	No     int    `json:"no"`
	Layout string `json:"layout"`
	Reason string `json:"reason,omitempty"`
}

// DeckFile deck.json 的结构：deck 的自描述元数据。
type DeckFile struct {
	ID         string           `json:"id"`
	Format     string           `json:"format"` // 恒为 v2
	Stage      string           `json:"stage"`
	Title      string           `json:"title"`
	TemplateID string           `json:"template_id,omitempty"`
	Variant    string           `json:"variant,omitempty"`
	Canvas     template.Canvas  `json:"canvas"`
	PagePlan   []PlanAssignment `json:"page_plan,omitempty"`
	CreatedAt  int64            `json:"created_at"`
	UpdatedAt  int64            `json:"updated_at"`
}

// ---------- Service 上的 v2 方法 ----------

// WithTemplateRegistry 挂载模板注册表（v2 页写入的类名契约数据源）。
// 返回自身以便链式装配；不挂载时 v2 写路径返回明确错误（模板契约不可用）。
func (s *Service) WithTemplateRegistry(reg *template.Registry) *Service {
	s.templates = reg
	return s
}

func (s *Service) templateFor(df *DeckFile) (*template.Template, error) {
	if s.templates == nil {
		return nil, fmt.Errorf("模板库不可用")
	}
	return s.templates.Get(df.TemplateID)
}

// deckFilePath / outlineFilePath：v2 的两个数据文件。
func (s *Service) deckFilePath(id string) string {
	return filepath.Join(s.decksDir, id, "deck.json")
}
func (s *Service) outlineFilePath(id string) string {
	return filepath.Join(s.decksDir, id, "outline.json")
}
func (s *Service) stylePath(id string) string {
	return filepath.Join(s.decksDir, id, "style.css")
}

// IndexPathV2 返回 v2 deck 的 index.html 路径（id 白名单校验后）。
func (s *Service) IndexPathV2(id string) (string, error) {
	if !idPattern.MatchString(id) {
		return "", fmt.Errorf("invalid deck id: %q", id)
	}
	return filepath.Join(s.decksDir, id, "index.html"), nil
}

// readDeckFile 读 deck.json（不校验归属——只允许 authorize 之后的内部调用）。
func (s *Service) readDeckFile(id string) (*DeckFile, error) {
	data, err := os.ReadFile(s.deckFilePath(id))
	if err != nil {
		return nil, fmt.Errorf("deck %s 不是 v2 格式（deck.json 缺失）", id)
	}
	var df DeckFile
	if err := json.Unmarshal(data, &df); err != nil {
		return nil, fmt.Errorf("deck.json 损坏: %w", err)
	}
	if df.Format != FormatV2 {
		return nil, fmt.Errorf("deck %s 不是 v2 格式", id)
	}
	return &df, nil
}

func (s *Service) writeDeckFile(id string, df *DeckFile) error {
	df.UpdatedAt = time.Now().Unix()
	data, err := json.MarshalIndent(df, "", "  ")
	if err != nil {
		return err
	}
	return atomicWriteFile(s.deckFilePath(id), data)
}

// GetDeckV2 归属校验后的元数据读取（handler 与 agent 工具都用它）。
func (s *Service) GetDeckV2(userID uint, id string) (*DeckFile, error) {
	if err := s.authorize(userID, id); err != nil {
		return nil, err
	}
	return s.readDeckFile(id)
}

// CreateV2Draft 新建 v2 deck（submit_outline 时调用）：占号 + 写 deck.json + DB 登记。
// 此时还没有 index.html（模板未选），预览端点对这种 deck 返回"尚未生成"。
func (s *Service) CreateV2Draft(userID uint, title string) (string, error) {
	if s.st == nil {
		return "", errStorage
	}
	if err := os.MkdirAll(s.decksDir, 0o755); err != nil {
		return "", fmt.Errorf("创建文件夹失败 %w", err)
	}
	id, err := s.claimDeckDir()
	if err != nil {
		return "", err
	}
	now := time.Now().Unix()
	df := &DeckFile{
		ID: id, Format: FormatV2, Stage: StageOutlining,
		Title: strings.TrimSpace(title), CreatedAt: now, UpdatedAt: now,
	}
	if err := s.writeDeckFile(id, df); err != nil {
		os.RemoveAll(filepath.Join(s.decksDir, id))
		return "", fmt.Errorf("写入 deck.json 失败: %w", err)
	}
	row := store.Deck{
		ID: id, UserID: userID, Title: strings.TrimSpace(title),
		Format: FormatV2, Stage: StageOutlining,
	}
	if err := s.st.DB.Create(&row).Error; err != nil {
		os.RemoveAll(filepath.Join(s.decksDir, id))
		return "", fmt.Errorf("登记 deck 归属: %w", err)
	}
	return id, nil
}

// ReadOutline 读大纲（归属校验）。
func (s *Service) ReadOutline(userID uint, id string) (*Outline, error) {
	if err := s.authorize(userID, id); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(s.outlineFilePath(id))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("deck %s 还没有大纲", id)
		}
		return nil, err
	}
	var o Outline
	if err := json.Unmarshal(data, &o); err != nil {
		return nil, fmt.Errorf("outline.json 损坏: %w", err)
	}
	return &o, nil
}

// OutlineConflict 版本冲突：调用方（REST 409 / 工具报错）需要最新版本号。
type OutlineConflict struct {
	LatestVersion int
}

func (e OutlineConflict) Error() string {
	return fmt.Sprintf("大纲版本过期（当前 v%d）：请先读取最新版再修改", e.LatestVersion)
}

// SaveOutline 写大纲（双通道共用）：schema 校验 + version CAS + 标题同步。
// expectVersion <= 0 表示首次写入（submit_outline），此时要求 outline.json 不存在。
func (s *Service) SaveOutline(userID uint, id string, o *Outline, expectVersion int) error {
	if err := s.authorize(userID, id); err != nil {
		return err
	}
	unlock := s.lockDeck(id)
	defer unlock()

	df, err := s.readDeckFile(id)
	if err != nil {
		return err
	}

	if expectVersion <= 0 {
		if _, err := os.Stat(s.outlineFilePath(id)); err == nil {
			return fmt.Errorf("大纲已存在：修改请用 update 语义（携带 version）")
		}
		o.Version = 1
	} else {
		cur, err := s.ReadOutline(userID, id)
		if err != nil {
			return fmt.Errorf("读取现有大纲失败: %w", err)
		}
		if cur.Version != expectVersion {
			return OutlineConflict{LatestVersion: cur.Version}
		}
		o.Version = cur.Version + 1
	}

	if err := o.Validate(); err != nil {
		return err
	}
	data, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return err
	}
	if err := atomicWriteFile(s.outlineFilePath(id), data); err != nil {
		return fmt.Errorf("写入大纲失败: %w", err)
	}

	// 标题随大纲走（deck.json + DB 行同步；列表页显示的是 DB 标题）
	if o.Title != df.Title {
		df.Title = o.Title
		_ = s.st.DB.Model(&store.Deck{}).Where("id = ?", id).Update("title", df.Title).Error
	}

	// 首次落大纲 = 大纲产出完成：outlining → outline_review（等用户确认）。
	// 后续修改（outline_review 阶段的 update_outline）不改变阶段。
	if df.Stage == StageOutlining {
		df.Stage = StageOutlineReview
		_ = s.writeDeckFile(id, df)
		if err := s.st.DB.Model(&store.Deck{}).Where("id = ?", id).Update("stage", StageOutlineReview).Error; err != nil {
			return fmt.Errorf("推进阶段失败: %w", err)
		}
	}
	return nil
}

// transitionStage 阶段迁移的唯一入口：校验 from（期望的当前阶段）后写 DB 与 deck.json。
// from 传空 = 不检查当前阶段（仅限初始化路径使用）。
// 返回冲突时的实际阶段，调用方据此给用户可执行的提示。
func (s *Service) transitionStage(userID uint, id, from, to string) (string, error) {
	if err := s.authorize(userID, id); err != nil {
		return "", err
	}
	unlock := s.lockDeck(id)
	defer unlock()

	var row store.Deck
	if err := s.st.DB.Where("id = ?", id).First(&row).Error; err != nil {
		return "", fmt.Errorf("deck %q 不存在", id)
	}
	if from != "" && row.Stage != from {
		return row.Stage, fmt.Errorf("阶段不对：当前 %s，此操作要求 %s", row.Stage, from)
	}
	if err := s.st.DB.Model(&store.Deck{}).Where("id = ?", id).Update("stage", to).Error; err != nil {
		return row.Stage, fmt.Errorf("更新阶段失败: %w", err)
	}
	if df, err := s.readDeckFile(id); err == nil {
		df.Stage = to
		_ = s.writeDeckFile(id, df)
	}
	return to, nil
}

// ConfirmOutline gate 1：大纲确认（outline_review → selecting_template）。
func (s *Service) ConfirmOutline(userID uint, id string) (string, error) {
	if _, err := s.ReadOutline(userID, id); err != nil {
		return "", fmt.Errorf("还没有可确认的大纲: %w", err)
	}
	return s.transitionStage(userID, id, StageOutlineReview, StageSelectingTemplate)
}

// UpdateOutlinePage 大纲确认后的面板直改入口已封死（确认即冻结）。
// 迭代阶段的页增删走 insert_slide/delete_slide，大纲保持"生成时的契约"不被回写——
// 否则"用户看到的大纲"和"实际的 deck"两份真相会互相漂移。

// SelectTemplate gate 2：选择模板 + 变体（selecting_template → generating）。
// 实例化产物（index.html + style.css）落 deck 目录；任一步失败整体回滚，
// 不留下"选了一半"的 deck。
func (s *Service) SelectTemplate(userID uint, id, templateID, variantID string) (string, error) {
	if err := s.authorize(userID, id); err != nil {
		return "", err
	}
	if s.templates == nil {
		return "", fmt.Errorf("模板库不可用")
	}
	tpl, err := s.templates.Get(templateID)
	if err != nil {
		return "", err
	}
	if variantID == "" {
		variantID = "default"
	}
	if _, ok := tpl.VariantClass(variantID); !ok {
		return "", fmt.Errorf("模板 %s 没有变体 %q（可选：%s）", templateID, variantID, variantNames(tpl))
	}

	// 大纲必须存在：实例化出的骨架马上要被生成管线填内容，没有大纲的 v2 deck 不该走到这
	o, err := s.ReadOutline(userID, id)
	if err != nil {
		return "", fmt.Errorf("选择模板前需要已确认的大纲: %w", err)
	}

	ins, err := tpl.Instantiate(variantID, o.Title)
	if err != nil {
		return "", fmt.Errorf("模板实例化失败: %w", err)
	}

	unlock := s.lockDeck(id)
	defer unlock()

	// 阶段守卫放锁内（与实例化原子），防 confirm 与 select 并发竞态
	var row store.Deck
	if err := s.st.DB.Where("id = ?", id).First(&row).Error; err != nil {
		return "", fmt.Errorf("deck %q 不存在", id)
	}
	if row.Stage != StageSelectingTemplate {
		return row.Stage, fmt.Errorf("阶段不对：当前 %s，此操作要求 %s", row.Stage, StageSelectingTemplate)
	}

	ip, err := s.IndexPathV2(id)
	if err != nil {
		return "", err
	}
	if err := atomicWriteFile(ip, []byte(ins.IndexHTML)); err != nil {
		return "", fmt.Errorf("写入 index.html 失败: %w", err)
	}
	if err := atomicWriteFile(s.stylePath(id), []byte(ins.StyleCSS)); err != nil {
		return "", fmt.Errorf("写入 style.css 失败: %w", err)
	}

	df, err := s.readDeckFile(id)
	if err != nil {
		return "", err
	}
	df.TemplateID = templateID
	df.Variant = variantID
	df.Canvas = tpl.Canvas
	df.Stage = StageGenerating
	if err := s.writeDeckFile(id, df); err != nil {
		return "", err
	}
	if err := s.st.DB.Model(&store.Deck{}).Where("id = ?", id).Updates(map[string]any{
		"stage": StageGenerating, "template_id": templateID, "variant": variantID,
	}).Error; err != nil {
		return "", fmt.Errorf("更新 deck 行失败: %w", err)
	}
	return StageGenerating, nil
}

func variantNames(tpl *template.Template) string {
	names := make([]string, 0, len(tpl.Variants))
	for _, v := range tpl.Variants {
		names = append(names, v.ID)
	}
	return strings.Join(names, " / ")
}

// FinishGeneration 生成 run 正常结束的落点：generating → iterating。
// run 失败/中断时不调用——stage 停在 generating，恢复端点靠它判断续跑。
func (s *Service) FinishGeneration(userID uint, id string) (string, error) {
	return s.transitionStage(userID, id, StageGenerating, StageIterating)
}

// ---------- v2 的渲染读取（预览/量测/导出共用） ----------

var styleLinkRe = regexp.MustCompile(`<link[^>]*href="style\.css"[^>]*>`)

// PreviewHTML 预览用 HTML：v2 把 style.css 内联进 <head>。
//
// 为什么内联：deck 页面位于 /api/decks/:id/file（或一次性 render 通道），
// 相对引用 style.css 会解析到不存在的路径。内联让单条 HTML 响应自包含，
// 预览 iframe、无头浏览器、导出三条通道共用同一份产出——"量测看到的"和
// "用户看到的"才是同一个东西。runtime/base.css/字体是 /assets 绝对引用，不受影响。
func (s *Service) PreviewHTML(userID uint, id string) (string, error) {
	html, err := s.GetHTML(userID, id)
	if err != nil {
		return "", err
	}
	if !s.IsV2(id) {
		return html, nil
	}
	css, err := os.ReadFile(s.stylePath(id))
	if err != nil {
		return html, nil // style.css 缺失时退回原文（v2 骨架在实例化后就有了，这不该发生）
	}
	inline := "<style>\n" + string(css) + "\n</style>"
	return styleLinkRe.ReplaceAllLiteralString(html, inline), nil
}
