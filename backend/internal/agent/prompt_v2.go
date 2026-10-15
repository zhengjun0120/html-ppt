package agent

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
	"time"

	"html-ppt/backend/internal/service/deck"
	"html-ppt/backend/internal/service/template"
)

// 阶段化提示词（deck-v2）：shared.md 全阶段共享，其余按阶段拼接。
//
// 保留 v1 的 LF 归一化纪律：embed 原样嵌入文件字节，Windows 下是 CRLF，
// 逐字发给模型等于白交 11% 输入费用（实测见 git 历史）。所有阶段文件在
// init 里统一归一化进缓存。

//go:embed prompts/*.md
var promptsFS embed.FS

var promptCache = map[string]string{}

func init() {
	entries, err := fs.ReadDir(promptsFS, "prompts")
	if err != nil {
		panic("prompts 目录缺失: " + err.Error())
	}
	for _, e := range entries {
		raw, err := promptsFS.ReadFile("prompts/" + e.Name())
		if err != nil {
			panic("读 prompts 失败: " + err.Error())
		}
		promptCache[e.Name()] = strings.ReplaceAll(string(raw), "\r\n", "\n")
	}
}

// prompt 加载阶段提示词（已归一化）。
func prompt(name string) string {
	s, ok := promptCache[name]
	if !ok {
		panic("提示词文件不存在: " + name)
	}
	return s
}

// stagePromptName 阶段 → 提示词文件名。返回空串表示该阶段没有独立提示词文件
//（没有 deck 的会话用 clarify.md，v1 deck 走旧 systemPrompt，都不在这里）。
func stagePromptName(stage string) string {
	switch stage {
	case "clarifying":
		return "clarify.md"
	case deck.StageOutlining, deck.StageOutlineReview:
		return "outline.md"
	case deck.StageGenerating:
		return "generate.md"
	case deck.StageIterating:
		return "iterate.md"
	default:
		return ""
	}
}

// BuildStagePrompt 组装 v2 阶段系统提示词：shared + 阶段文件 + 动态注入。
//
// 动态注入保持 v1 的前缀缓存纪律：变化的部分（deck_id、日期）放最后，
// 静态规则放前面——打断前缀缓存的代价见 v1 buildSystemMessage 的实测注释。
//
//	tpl 非 nil 时（generating/iterating）追加「模板契约」段：版式索引 + rules.md 全文。
//
//	density 是会话级内容密度档位（concise/standard/rich，空/未知 = standard）：
//	在末尾注入对应档位的字数指令，并显式声明覆盖 shared.md 的「每条 ≤28 字」
//	上限——不覆盖的话通用规则会把加料压回去。澄清阶段（stage==""）不注入。
func BuildStagePrompt(stage, deckID string, tpl *template.Template, o *deck.Outline, density string) string {
	var b strings.Builder
	b.WriteString(prompt("shared.md"))
	b.WriteString("\n\n")
	if stage == "" {
		// 无 deck 的冷启动：阶段就是澄清
		b.WriteString(prompt("clarify.md"))
	} else {
		name := stagePromptName(stage)
		if name == "" {
			name = "iterate.md" // 未识别的阶段兜底：迭代纪律最通用
		}
		b.WriteString(prompt(name))
	}

	if deckID != "" {
		fmt.Fprintf(&b, "\n\n当前 deck：{\"deck_id\":\"%s\"}。涉及它的读写直接用这个 deck_id，不要向用户询问。", deckID)
	}
	if stage == deck.StageOutlineReview {
		b.WriteString("\n\n当前处于「大纲待确认」状态：用户可能在大纲面板里直接改，也可能在对话里让你改。让你改就走修订模式；用户说\"确认了/没问题了\"时提示他点界面上的「确认大纲」按钮——你没有确认工具，确认动作在用户手里。")
	}
	// 大纲注入：生成阶段是全文（LLM 写页的唯一内容来源）；迭代阶段是紧凑索引
	//（页码+role+标题，改哪页对哪页）。outline 为 nil 时跳过（v1 路径不会到这里）。
	if o != nil {
		if stage == deck.StageGenerating {
			// 标题规则常驻：T005 只在写页后才报，而大纲是模板选定之前写的、看不到
			// 这条约束。deck-0062 实测：只说「≤16 字」agent 会按“英文单词算 1 个词”
			// 的口径自数出 12 字，T005 实际逐字符计数出 17 字——首写照样超限，白烧
			// 一轮修复。口径必须写死：去空格后逐字符计，英文逐字母算。
			fmt.Fprintf(&b, "\n\n## 大纲（已确认，页数与内容以此为准）\n"+
				"页标题硬上限 %d 字：按去掉空格后的字符数计，英文单词也逐字母算（GitHub 算 6 字，不是 1 个词）。"+
				"主标写短，修饰成分挪进 lede 副句；页内其余 h2 小标题同限。\n", deck.TitleMax())
			b.WriteString(o.ToPromptText())
			// 超限预检：点名具体页，冒号能拆的附上现成拆法（主标 + lede 承接），
			// 让 agent 首批一次写对。
			if over := overlongOutlineTitles(o.Pages, deck.TitleMax()); len(over) > 0 {
				fmt.Fprintf(&b, "\n**标题预检**：以下大纲标题超限，写页时直接改短、不要照抄：%s。",
					strings.Join(over, "；"))
			}
		} else {
			b.WriteString("\n\n## 大纲索引（页码 · role · 标题）\n")
			for _, pg := range o.Pages {
				fmt.Fprintf(&b, "\n- 第 %d 页 [%s] %s", pg.No, pg.Role, pg.Title)
			}
		}
	}
	if tpl != nil && (stage == deck.StageGenerating || stage == deck.StageIterating) {
		b.WriteString("\n\n## 模板契约：")
		fmt.Fprintf(&b, "\n模板：%s（%s）。画布 %d×%d 固定设计像素。", tpl.Name, tpl.ID, tpl.Canvas.W, tpl.Canvas.H)
		b.WriteString("\n可用版式（id · 名称 · 用途）：")
		for _, l := range tpl.Layouts {
			roles := ""
			if len(l.Roles) > 0 {
				roles = "〔role: " + strings.Join(l.Roles, "/") + "〕"
			}
			fmt.Fprintf(&b, "\n- %s · %s · %s%s", l.ID, l.Name, l.Use, roles)
		}
		b.WriteString("\n每个版式的骨架代码与合法类名用 read_layout 获取；模板专属规则如下：\n")
		b.WriteString(tpl.Rules())
	}
	// 密度指令放末尾（前缀缓存纪律：会话内不变、跨会话才变化的部分归尾）。
	// 澄清阶段也要注入：冷启动流程的大纲正是在 clarifying 这个 run 里写出来的
	//（实测教训：只对 stage!="" 注入时，rich 会话的首份大纲每页只有 2-3 条，
	// 密度指令没机会生效）。
	b.WriteString("\n\n" + densityDirectives[normalizeDensity(density)])
	return b.String()
}

// 内容密度三档（会话级，ChatSession.Density）：决定大纲要点数量、每条长度与
// 页面正文目标。措辞必须显式声明覆盖 shared.md 的「每条 ≤28 字」上限——那段
// 通用规则在前、密度段在后，不点名覆盖的话模型会按更严的那条执行。
// 封面/章节/收尾/大数字这类 hero 页一律豁免：留白是设计，不为凑字数破坏。
const (
	DensityConcise  = "concise"
	DensityStandard = "standard"
	DensityRich     = "rich"
)

var densityDirectives = map[string]string{
	DensityConcise: "## 内容密度：精炼\n" +
		"沿用前文的要点规矩：每页 2-5 条、每条 ≤28 字、卡片说明 22-45 字。" +
		"封面/章节/收尾/大数字等 hero 页保持留白设计。",
	DensityStandard: "## 内容密度：标准（本段覆盖前文「每条 ≤28 字」的上限）\n" +
		"- 大纲要点：每页 3-6 条，每条 ≤45 字；要点要带具体信息（数字、例子、出处），不写光杆名词短语。\n" +
		"- 写页：卡片说明 40-70 字，把要点扩写成完整句；内容页（概念/数据/对比/时间线/步骤等）正文总量目标 180-280 字，填充率目标 ≥60%。\n" +
		"- 封面/章节/收尾/大数字等 hero 页豁免：留白是设计，不要为凑字数硬写。",
	DensityRich: "## 内容密度：详实（本段覆盖前文「每条 ≤28 字」的上限）\n" +
		"- 大纲要点：每页 4-7 条，每条 ≤55 字；必须带数字、例子或出处。\n" +
		"- 写页：卡片说明 50-90 字；内容页（概念/数据/对比/时间线/步骤等）正文总量目标 ≥300 字，填充率目标 ≥70%，偏空时优先换信息密度更高的版式。\n" +
		"- 封面/章节/收尾/大数字等 hero 页豁免：留白是设计，不要为凑字数硬写。",
}

// normalizeDensity 收敛档位值：空与未知一律 standard（老会话无此字段的默认取向）。
func normalizeDensity(d string) string {
	if d == DensityConcise || d == DensityRich {
		return d
	}
	return DensityStandard
}

// appendDate 尾部追加当前日期（前缀缓存纪律：变化内容永远放最后）。
func appendDate(msg string) string {
	now := time.Now()
	return msg + fmt.Sprintf("\n\n当前日期：%s（%s）。涉及「今天」「本月」「最近」这类时间说法时以它为准——"+
		"你的训练数据有截止时间，不要按它推断当前时间，也不要为了确认日期去联网搜索。",
		now.Format("2006-01-02"), weekdayCN[int(now.Weekday())])
}

// overlongOutlineTitles 按 T005 同口径（去空格计 rune 数）找出超长的大纲标题；
// 冒号能拆出「达标主标 + 副句」的附上现成拆法——deck-0062 的教训：只报字数不给
// 拆法，agent 首写还是超。
func overlongOutlineTitles(pages []deck.OutlinePage, max int) []string {
	var out []string
	for _, p := range pages {
		n := len([]rune(strings.ReplaceAll(p.Title, " ", "")))
		if n <= max {
			continue
		}
		entry := fmt.Sprintf("第 %d 页「%s」（%d 字）", p.No, p.Title, n)
		if head, tail, ok := splitTitleCandidate(p.Title, max); ok {
			entry += fmt.Sprintf("，可拆成主标「%s」+ lede 承接「%s」", head, tail)
		}
		out = append(out, entry)
	}
	return out
}

// splitTitleCandidate 在第一个冒号处把标题拆成 主标候选 + 副句；主标候选仍超限
// 或冒号两侧有空就放弃——宁可不给建议，不给一个过不了的候选。
func splitTitleCandidate(title string, max int) (string, string, bool) {
	idx := strings.IndexAny(title, "：:")
	if idx < 0 {
		return "", "", false
	}
	head := strings.TrimSpace(title[:idx])
	tail := strings.TrimSpace(strings.TrimLeft(title[idx:], "：:"))
	if head == "" || tail == "" {
		return "", "", false
	}
	if len([]rune(strings.ReplaceAll(head, " ", ""))) > max {
		return "", "", false
	}
	return head, tail, true
}
