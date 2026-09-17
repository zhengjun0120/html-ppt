package deck

// 版式节奏校验（guizang 思想重写，服务端机械执行——提示词拦不住的东西这里拦）。
//
// 约束级别（与 docs/refactor-plan.md §9.1 一致）：
//   R101 同版式连续 ≥3           —— 阻塞（plan_pages 打回重排）
//   R102 每 8 页 <4 种版式        —— 提示
//   R103 每 8 页无满版 hero 类    —— 提示（本模板的 hero 版式由 rules 语义约定）
//   R104 左右图文交替连续 >2      —— 提示（左右由模板在版式名里约定，此处按名匹配）
//   R105 实写页数 ≠ 大纲页数      —— 阻塞（写在 v2_pages 的页码范围校验里）

import (
	"fmt"
	"strings"

	"html-ppt/backend/internal/service/template"
)

// RhythmViolation 单条节奏违规。
type RhythmViolation struct {
	Rule  string `json:"rule"`   // R101 / R102 / ...
	Block bool   `json:"block"`  // true = 阻塞级
	Msg   string `json:"msg"`
}

// heroLayoutNames 满版 hero 类版式的识别词（版式 id 里含这些词即视为 hero）。
var heroLayoutNames = []string{"cover", "hero", "divider", "stat-hero", "big-quote", "qa"}

// alternatingLayoutNames 左右图文类版式的识别词。
var alternatingLayoutNames = []string{"split", "image-text", "two-column", "left-right"}

func nameHas(id string, names []string) bool {
	for _, n := range names {
		if strings.Contains(id, n) {
			return true
		}
	}
	return false
}

// ValidateRhythm 对「页码 → 版式 id」序列做节奏校验。
// layouts 是模板登记的全部版式 id（用于 R103 的存在性判断）。
func ValidateRhythm(assignments map[int]string, order []int) []RhythmViolation {
	var out []RhythmViolation
	layouts := make([]string, 0, len(order))
	for _, no := range order {
		layouts = append(layouts, assignments[no])
	}

	// R101 同版式连续 ≥3（阻塞）：滑动窗口线性扫
	runStart := 0
	for i := 1; i <= len(layouts); i++ {
		if i == len(layouts) || layouts[i] != layouts[runStart] {
			if i-runStart >= 3 {
				out = append(out, RhythmViolation{
					Rule: "R101", Block: true,
					Msg: fmt.Sprintf("第 %d~%d 页连续 %d 页都是 %q：同一版式最多连用 2 页，穿插 divider/stat/quote 类版式换节奏",
						order[runStart], order[i-1], i-runStart, layouts[runStart]),
				})
			}
			runStart = i
		}
	}

	// R102 / R103：每 8 页一个窗口
	for start := 0; start < len(layouts); start += 8 {
		end := start + 8
		if end > len(layouts) {
			end = len(layouts)
		}
		win := layouts[start:end]
		if len(win) < 8 {
			break // 不满一个完整窗口不判：短 deck 尾部窗口样本太少
		}
		distinct := map[string]bool{}
		for _, l := range win {
			distinct[l] = true
		}
		if len(distinct) < 4 {
			out = append(out, RhythmViolation{Rule: "R102", Block: false,
				Msg: fmt.Sprintf("第 %d~%d 页只用了 %d 种版式（至少 4 种）：整份会显得单调",
					start+1, end, len(distinct))})
		}
		hasHero := false
		for _, l := range win {
			if nameHas(l, heroLayoutNames) {
				hasHero = true
				break
			}
		}
		if !hasHero {
			out = append(out, RhythmViolation{Rule: "R103", Block: false,
				Msg: fmt.Sprintf("第 %d~%d 页没有满版 hero 类版式（cover/divider/stat-hero/big-quote）：缺一页\"呼吸感\"的锚点",
					start+1, end)})
		}
	}

	// R104 左右图文交替连续 >2
	runStart = 0
	isAlt := func(l string) bool { return nameHas(l, alternatingLayoutNames) }
	for i := 1; i <= len(layouts); i++ {
		if i == len(layouts) || isAlt(layouts[i]) != isAlt(layouts[runStart]) {
			if isAlt(layouts[runStart]) && i-runStart > 2 {
				out = append(out, RhythmViolation{Rule: "R104", Block: false,
					Msg: fmt.Sprintf("第 %d~%d 页连续 %d 页都是左右图文类：连续超过 2 页会 mechanical（机械感）",
						order[runStart], order[i-1], i-runStart)})
			}
			runStart = i
		}
	}

	return out
}

// SavePlanV2 plan_pages 的落盘：校验覆盖完整性、版式登记、节奏（R101 阻塞），
// 全过才写入 deck.json.PagePlan。
func (s *Service) SavePlanV2(userID uint, deckID string, assignments []PlanAssignment) ([]RhythmViolation, error) {
	if err := s.authorize(userID, deckID); err != nil {
		return nil, err
	}
	df, err := s.readDeckFile(deckID)
	if err != nil {
		return nil, err
	}
	if df.Stage != StageGenerating {
		return nil, fmt.Errorf("当前阶段 %s 不允许规划版式", df.Stage)
	}
	tpl, err := s.templateFor(df)
	if err != nil {
		return nil, err
	}
	outline, err := s.ReadOutline(userID, deckID)
	if err != nil {
		return nil, err
	}

	// 覆盖完整性：每页一条、不多不少
	if len(assignments) != len(outline.Pages) {
		return nil, fmt.Errorf("计划覆盖 %d 页，大纲有 %d 页：必须每页恰好一条", len(assignments), len(outline.Pages))
	}
	byNo := map[int]string{}
	order := make([]int, 0, len(assignments))
	for _, a := range assignments {
		if a.No < 1 || a.No > len(outline.Pages) {
			return nil, fmt.Errorf("计划里的页码 %d 超出大纲范围", a.No)
		}
		if _, dup := byNo[a.No]; dup {
			return nil, fmt.Errorf("页码 %d 出现了两次", a.No)
		}
		if !tpl.HasLayout(a.Layout) {
			return nil, fmt.Errorf("第 %d 页的版式 %q 未登记（可用：%s）", a.No, a.Layout, strings.Join(layoutIDs(tpl), "、"))
		}
		byNo[a.No] = a.Layout
		order = append(order, a.No)
	}
	// order 按页码排序，违规信息里的页号区间才有意义
	for i := 0; i < len(order); i++ {
		for j := i + 1; j < len(order); j++ {
			if order[j] < order[i] {
				order[i], order[j] = order[j], order[i]
			}
		}
	}

	violations := ValidateRhythm(byNo, order)
	for _, v := range violations {
		if v.Block {
			return violations, fmt.Errorf("节奏校验未通过：%s", v.Msg)
		}
	}

	unlock := s.lockDeck(deckID)
	defer unlock()
	df, err = s.readDeckFile(deckID) // 锁内重读，避免与并发写互相覆盖
	if err != nil {
		return nil, err
	}
	df.PagePlan = assignments
	if err := s.writeDeckFile(deckID, df); err != nil {
		return nil, err
	}
	return violations, nil
}

func layoutIDs(tpl *template.Template) []string {
	ids := make([]string, 0, len(tpl.Layouts))
	for _, l := range tpl.Layouts {
		ids = append(ids, l.ID)
	}
	return ids
}
