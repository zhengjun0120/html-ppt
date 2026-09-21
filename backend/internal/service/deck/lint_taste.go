package deck

// AI 味 lint（T001-T007，全部提示级——品味问题提示给模型自修，不阻塞写入）。
// 规则来源：taste-skill（MIT）的 AI-Tells 词表与密度纪律，按"PPT 每页 ≈ 网页一个
// section"翻译成中文语境并程序化——原 skill 全靠模型自律，这里补上机械执行。

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// LintItem 单条 lint 提示（进 write_pages 返回值，模型据此自修）。
type LintItem struct {
	Rule string `json:"rule"`
	Word string `json:"word,omitempty"`
	Msg  string `json:"msg"`
	Hint string `json:"hint,omitempty"`
}

// 词表默认值（config deck_v2.lint 可整体覆盖）。
var (
	defaultCJKBan = []string{
		"赋能", "抓手", "闭环", "沉淀", "心智", "护城河", "组合拳", "打法",
		"底层逻辑", "顶层设计", "颗粒度", "拉齐", "值得注意的是", "综上所述",
		"总而言之", "众所周知", "毋庸置疑",
	}
	defaultENBan = []string{
		"elevate", "seamless", "unleash", "next-gen", "revolutionize", "empower",
		"cutting-edge", "delve", "tapestry", "robust", "holistic",
	}
)

// LintCfg lint 的可配置词表（service 层只读；main 装配时注入）。
type LintCfg struct {
	CJKBanned    []string
	ENBanned     []string
	TitleMax     int
	BulletMax    int
	EyebrowLimit int
}

func defaultLintCfg() LintCfg {
	return LintCfg{
		CJKBanned:    defaultCJKBan,
		ENBanned:     defaultENBan,
		TitleMax:     16, // 页标题字数上限（与 config DefaultTitleMaxChars 同源，deck 包内自持）
		BulletMax:    60, // 曾是 28：把要点条压成一行小句是 deck-0031 过空的推手之一，
		// 22-45 字的"一句论断+一句展开"才是卡片该有的密度（与契约区间一致）
		EyebrowLimit: 3,
	}
}

// lintCfg 当前生效的词表（WithLintConfig 覆盖；并发只读，装载期写一次）。
var lintCfg = defaultLintCfg()

// TitleMax 当前生效的页标题字数上限（T005 的同一来源）。生成提示词的大纲
// 标题预检用它——大纲在模板选定之前写、看不到这条约束，别在别的包抄常量。
func TitleMax() int { return lintCfg.TitleMax }

// WithLintConfig 注入配置词表（main 启动时调一次；nil 值字段保留默认）。
func WithLintConfig(cfg LintCfg) {
	if len(cfg.CJKBanned) > 0 {
		lintCfg.CJKBanned = cfg.CJKBanned
	}
	if len(cfg.ENBanned) > 0 {
		lintCfg.ENBanned = cfg.ENBanned
	}
	if cfg.TitleMax > 0 {
		lintCfg.TitleMax = cfg.TitleMax
	}
	if cfg.BulletMax > 0 {
		lintCfg.BulletMax = cfg.BulletMax
	}
	if cfg.EyebrowLimit > 0 {
		lintCfg.EyebrowLimit = cfg.EyebrowLimit
	}
}

// hintFor 禁词的替代写法建议（表内没有的给通用建议）。
var cjkHints = map[string]string{
	"赋能":   "写具体动作：支持 / 加速 / 让 X 能…",
	"抓手":   "写具体手段或直接删掉",
	"闭环":   "写完整流程的首尾（从 X 到 Y）",
	"沉淀":   "写“留存/积累”的具体结果",
	"心智":   "写“理解/印象”的具体内容",
	"护城河":  "写具体的竞争优势",
	"组合拳":  "列出具体动作",
	"打法":   "写具体做法",
	"底层逻辑": "写真正起作用的机制",
	"顶层设计": "写具体的整体方案",
	"颗粒度":  "写“粒度/精细程度”的具体标准",
	"拉齐":   "写“同步/对齐到一致”的对象",
	"值得注意的是": "直接说值得注意的那件事",
	"综上所述":  "直接给结论",
	"总而言之":  "直接给结论",
	"众所周知":  "删掉，直接陈述",
	"毋庸置疑":  "删掉，用证据说话",
}

var enHints = map[string]string{
	"elevate": "write the concrete effect", "seamless": "say what actually happens",
	"unleash": "use a plain verb", "next-gen": "say what actually improved",
	"revolutionize": "state the real change", "empower": "state what becomes possible",
	"cutting-edge": "name the technique", "delve": "use examine/cover",
	"tapestry": "delete the metaphor", "robust": "state what it withstands",
	"holistic": "say what is included",
}

var (
	// 破折号：先摘掉中文成对破折号再查单个 em-dash；标题内任何破折号都提示
	cjkDashRe    = regexp.MustCompile(`——`)
	enDashRe     = regexp.MustCompile(`—`)
	decimalPctRe = regexp.MustCompile(`\d+\.\d+%`)
	// 疑似英文单词的边界匹配（避免 robust 匹配到 robustness 之外的中文上下文误报过重——保持词边界即可）
	wordInText = func(word, text string) bool {
		return strings.Contains(strings.ToLower(text), word)
	}
)

// LintTaste 对一页做 AI 味检查（提示级，不改写不阻塞）。
// 可见文本范围：整页文本去掉 .notes 讲稿与标题（标题单独按 T005/T003 判）。
func LintTaste(sec *goquery.Selection) []LintItem {
	var titleParts []string
	sec.Find("h1, h2").Each(func(_ int, h *goquery.Selection) {
		titleParts = append(titleParts, h.Text())
	})
	titles := strings.Join(titleParts, " ")
	notes := sec.Find(".notes").Text()
	// 字符串减法：notes 在原文里只出现一次，摘掉即可；标题同理（正文判定不含标题）
	body := strings.Replace(strings.Replace(sec.Text(), notes, "", 1), titles, "", 1)
	text := strings.Join(strings.Fields(body), " ")
	titleText := strings.Join(strings.Fields(titles), " ")
	// 禁词/假数字检查覆盖"正文+标题"：标题是最显眼的位置，更不允许 AI 腔
	allText := text + "\n" + titleText
	var out []LintItem

	// T001 中文禁词
	for _, w := range lintCfg.CJKBanned {
		if strings.Contains(allText, w) {
			out = append(out, LintItem{Rule: "T001", Word: w,
				Msg:  fmt.Sprintf("中文 AI 腔用词「%s」", w),
				Hint: cjkHints[w]})
		}
	}
	// T002 英文禁词（整词命中才算）
	lower := strings.ToLower(allText)
	for _, w := range lintCfg.ENBanned {
		if hasEnglishWord(lower, w) {
			out = append(out, LintItem{Rule: "T002", Word: w,
				Msg:  fmt.Sprintf("英文填充词 %q", w),
				Hint: enHints[w]})
		}
	}

	// T003 破折号：标题里禁；正文每页 ≤1 次；单个 em-dash（英文腔）提示
	if strings.Contains(titleText, "——") || enDashRe.MatchString(titleText) {
		out = append(out, LintItem{Rule: "T003", Msg: "标题里不要用破折号",
			Hint: "标题用冒号或逗号断句"})
	}
	bodyNoPair := cjkDashRe.ReplaceAllString(text, "")
	if n := len(enDashRe.FindAllString(bodyNoPair, -1)); n > 0 {
		out = append(out, LintItem{Rule: "T003", Msg: fmt.Sprintf("正文有 %d 处英文 em-dash（—）", n),
			Hint: "中文用「——」，且每页最多 1 次；英文语境改用冒号"})
	} else if n := len(cjkDashRe.FindAllString(text, -1)); n > 1 {
		out = append(out, LintItem{Rule: "T003", Msg: fmt.Sprintf("中文破折号用了 %d 处（每页最多 1 次）", n),
			Hint: "保留最必要的一处，其余改逗号或冒号"})
	}

	// T004 假精确数字：带小数的百分比（92.5%）没有出处就是编的
	if m := decimalPctRe.FindAllString(allText, -1); len(m) > 0 {
		out = append(out, LintItem{Rule: "T004", Word: strings.Join(m, "、"),
			Msg:  fmt.Sprintf("精确百分比 %s 需要出处", strings.Join(m, "、")),
			Hint: "无来源就写量级（“约九成”）或标注“估算”"})
	}

	// T005 页标题超长（去空格计字）
	for _, t := range titleParts {
		if n := len([]rune(strings.ReplaceAll(t, " ", ""))); n > lintCfg.TitleMax {
			out = append(out, LintItem{Rule: "T005", Msg: fmt.Sprintf("页标题 %d 字（上限 %d）", n, lintCfg.TitleMax),
				Hint: "拆成主标 + lede 副句，或砍修饰成分"})
			break
		}
	}

	// T006 要点条超长（li 元素逐条判）
	sec.Find("li").Each(func(_ int, li *goquery.Selection) {
		if n := len([]rune(strings.Join(strings.Fields(li.Text()), ""))); n > lintCfg.BulletMax {
			out = append(out, LintItem{Rule: "T006", Msg: fmt.Sprintf("要点条 %d 字（上限 %d）：%s", n, lintCfg.BulletMax, truncateRunes(li.Text(), 24)),
				Hint: "拆成两条，或把细节移进讲稿"})
		}
	})

	// T007 眉题限流：一页内 kicker/eyebrow 类元素 >1
	if n := sec.Find(".kicker, .eyebrow").Length(); n > 1 {
		out = append(out, LintItem{Rule: "T007", Msg: fmt.Sprintf("本页有 %d 个眉题类元素（每页最多 1 个）", n),
			Hint: "删掉多余的，或并进正文"})
	}

	return out
}

func hasEnglishWord(lowerText, word string) bool {
	idx := 0
	for {
		i := strings.Index(lowerText[idx:], word)
		if i < 0 {
			return false
		}
		start := idx + i
		end := start + len(word)
		// 边界 = 词的前一个字符与后一个字符（词首/词尾本身当然是字母）
		if isWordBoundaryAt(lowerText, start-1) && isWordBoundaryAt(lowerText, end) {
			return true
		}
		idx = start + 1
	}
}

func isWordBoundaryAt(s string, i int) bool {
	if i <= 0 || i >= len(s) {
		return true
	}
	c := s[i]
	isLetter := (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '_'
	return !isLetter
}

func truncateRunes(s string, n int) string {
	r := []rune(strings.Join(strings.Fields(s), ""))
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}
