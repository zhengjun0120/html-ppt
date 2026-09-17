package agent

// 系统消息拼装的测试。这里守两件事：
//
//  1. **日期必须待在系统提示词的最后**。这不是排版偏好，是缓存纪律：前缀缓存要求
//     "完整匹配一个已持久化的前缀单元"，变化的内容越靠后，被打断的越少。
//     实测把末尾日期改一天，cached_tokens 从 2560/2696 只掉到 2432/2696（前面照样命中）；
//     把日期挪到开头就是第一个 token 分叉、后面全部重算。以后有人"顺手整理一下顺序"
//     把日期挪上去，这条测试要当场变红。
//  2. 日期本身要正确、且**不是**写死在 systemPrompt.md 里的（那个文件是 //go:embed
//     编译进二进制的，写死等于发布即过期）。
//
// 末尾几条是**压缩提示词时的守卫**：压缩真正会造成的损伤不是"话变少了"，
// 而是悄悄少掉一条信息——某个组件类不再被提及、某个数字和代码对不上、容量预算被删。
// 它们都从**源头派生**期望值（组件库文件、deck 包的画布表），所以不需要人工维护清单：
// 代码改了而提示词没跟上，测试就会红。

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"html-ppt/backend/internal/service/deck"
)

func TestBuildSystemMessagePutsDateLast(t *testing.T) {
	svc := &AgentService{}
	msg := svc.buildSystemMessage(0, "deck-0001")

	// 静态提示词要原样留在最前面（防拼装时截断）
	if !strings.HasPrefix(msg, systemPrompt) {
		t.Fatal("系统消息应当以静态提示词开头")
	}

	today := time.Now().Format("2006-01-02")
	if !strings.Contains(msg, today) {
		t.Errorf("系统消息里没有今天的日期 %s：模型无从知道「现在」是什么时候", today)
	}

	// 核心断言：日期在 deck_id 那段**之后**，也就是整个消息的最后
	deckLineAt := strings.Index(msg, "deck-0001")
	dateAt := strings.Index(msg, today)
	if deckLineAt < 0 {
		t.Fatal("deck_id 那行不见了")
	}
	if dateAt < deckLineAt {
		t.Error("日期必须排在 deck_id 之后（系统提示词的最末尾）：挪到前面会让缓存前缀提前分叉——" +
			"实测末尾改一天只损失 ~128 token，放开头则整段系统提示词重算")
	}
	if !strings.HasSuffix(strings.TrimSpace(msg), "也不要为了确认日期去联网搜索。") {
		t.Error("日期那段应当收尾在整个系统提示词的最后")
	}

	// 星期要配上：不然「上周末」「这周三」这类说法算不出来
	found := false
	for _, w := range weekdayCN {
		if strings.Contains(msg, w) {
			found = true
			break
		}
	}
	if !found {
		t.Error("日期里没有星期")
	}
}

// 没有 deck_id（或它不是合法 id）时，系统消息只该是"静态提示词 + 日期"，
// 不能冒出一条空的 deck 说明——模型会照着一句没有 id 的话去瞎猜。
func TestBuildSystemMessageWithoutDeck(t *testing.T) {
	svc := &AgentService{}
	for _, bad := range []string{"", "../../etc/passwd", "not a deck id"} {
		msg := svc.buildSystemMessage(0, bad)
		if strings.Contains(msg, "当前用户正在预览的演示文稿") {
			t.Errorf("deckID=%q 时不该有 deck 说明", bad)
		}
		if !strings.Contains(msg, time.Now().Format("2006-01-02")) {
			t.Errorf("deckID=%q 时日期仍应在（它与会话相关，与 deck 无关）", bad)
		}
	}
}

// 星期映射不能靠"看起来对"：Go 的 time.Weekday 是周日=0、周六=6，
// 数组一错位就是整体偏一天，而这种错在页面上只会表现为"星期三的事写成星期四"。
func TestWeekdayCNMatchesGoWeekday(t *testing.T) {
	if weekdayCN[time.Sunday] != "周日" {
		t.Errorf("time.Sunday(%d) 应映射到 周日，实际 %q", time.Sunday, weekdayCN[time.Sunday])
	}
	if weekdayCN[time.Saturday] != "周六" {
		t.Errorf("time.Saturday(%d) 应映射到 周六，实际 %q", time.Saturday, weekdayCN[time.Saturday])
	}
	if len(weekdayCN) != 7 {
		t.Errorf("星期表应有 7 项，实际 %d", len(weekdayCN))
	}
}

// 日期必须由代码现算，不能写进内嵌的 systemPrompt.md：
// 那份文件编译进二进制，写死的日期会静默变成假话，而且没有任何东西会报出来。
func TestDateIsNotHardcodedInPromptFile(t *testing.T) {
	if strings.Contains(systemPrompt, "当前日期") {
		t.Error("systemPrompt.md 里出现了「当前日期」：日期必须由 buildSystemMessage 现算并放在末尾" +
			"（内嵌文件里写死的日期 = 发布即过期）")
	}
}

// 组件库里定义过的类，**提示词的组件库那一节**必须都提到。没提到 = 模型不知道它存在，
// 那处效果就只能手写内联样式去实现——而内联是这套体系里最该避免的写法。
//
// 为什么限定在那一节而不是整个文件：整份文件里搜到一次，可能是别处顺口带过的
// （"或者把它单独放一页用 .bleed"），而模型学"有这个类、什么时候用"靠的是那份目录。
// 期望值从 components.css 现扫，所以清单不用人维护。
func TestPromptMentionsEveryComponentClass(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "web", "assets", componentCSSFile))
	if err != nil {
		t.Fatalf("读不到组件库（本包到 web/assets 是三层）: %v", err)
	}
	names := extractClassNames(string(raw))
	if len(names) < 30 {
		t.Fatalf("只解析出 %d 个组件类，解析逻辑或组件库文件不对", len(names))
	}

	start := strings.Index(systemPrompt, "# 组件库")
	if start < 0 {
		t.Fatal("提示词里没有「# 组件库」一节")
	}
	catalog := systemPrompt[start:]
	if end := strings.Index(catalog[1:], "\n# "); end > 0 {
		catalog = catalog[:end+1]
	}
	for _, n := range names {
		// 用 hasClassToken 而不是 strings.Contains：后者是子串匹配，
		// `.row` 会被 `.rows` 满足、`.card` 会被 `.cards` 满足——删掉了条目的假通过。
		// 这正是 component_lib.go 里那个按类名边界匹配的既有工具。
		if !hasClassToken(catalog, "."+n) {
			t.Errorf("组件库定义了 .%s，但提示词的组件库一节没提它——模型会以为没有这个类，改用内联样式实现", n)
		}
	}
}

// 画布尺寸写在提示词里（模型要按它分配字数），但真相在代码里。
// 这条从代码派生期望值，防的是"改了 DefaultCanvas、提示词还写着旧尺寸"——
// 那种错模型是照着错的预算排版的，页面上表现为"明明按提示词写的却装不下"。
// （这个坑真的踩过：默认画布从 960 换成 1244 之后，提示词里还留着 960×700。）
func TestPromptCanvasNumbersMatchCode(t *testing.T) {
	w, h := deck.CanvasSize(deck.DefaultCanvas)
	want := fmt.Sprintf("%d×%d", w, h)
	if !strings.Contains(systemPrompt, want) {
		t.Errorf("默认画布是 %s，提示词里没有这个尺寸（模型会按错的画布分配字数）", want)
	}
}

// 容量预算的几个硬数字：压缩提示词时最容易顺手删掉的就是它们，
// 而删掉之后模型只会写得更满——溢出的代价比多几十个 token 大得多。
func TestPromptKeepsCapacityBudget(t *testing.T) {
	for _, s := range []string{"200 字", "≤ 5 个", "1094×618", "8~14 页"} {
		if !strings.Contains(systemPrompt, s) {
			t.Errorf("提示词里找不到容量预算的关键数字 %q——它被压缩掉了？", s)
		}
	}
}

// 跨页节奏的四条规则。它们读起来最像"劝告式的话"，压缩时最容易被当成废话删掉，
// 但每条都对应一个在真实 deck 上量到的缺口（量法与完整数字见 docs/tools.md「通用实现骨架」）：
//   - 16 份真实 deck 里 0 份把最后一页做成收尾页（`.center`），多数停在最后一个知识点上
//   - `.bleed`/`.plate`/`.band` 这三个整页低密度版式，16 份里 13 份一次都没用过
//   - deck-0015 里用了提头的 5 页**全部**是标题的复述
//   - 大纲层面只有 deck-0016/0017 看得出钩子与收束，其余是"整体思路 → 各部分 → 小结"的话题清单
//
// 删掉它们，产出就会退回"每页都是标题 + 等分卡片"。
func TestPromptKeepsRhythmRules(t *testing.T) {
	// 分节断言，而不是整份文件里搜一次就算数："叙事弧"在标准工作流那行里也出现过，
	// 只搜全文的话，把决定一里的定义整段删掉，测试仍然是绿的——
	// 和组件库那条守卫同一个坑（那里是限定在"组件库"一节里搜）。
	upfront := promptRegion(t, "# 动手写之前", "\n# 单页容量预算")
	catalog := promptRegion(t, "# 组件库", "\n# 文案纪律")

	// 决定一的大纲骨架 + 决定三的密度落差
	for _, s := range []string{"叙事弧", "钩子", "收束", "每 3~4 页"} {
		if !strings.Contains(upfront, s) {
			t.Errorf("「动手写之前」那几节里找不到 %q——它被压缩掉了？", s)
		}
	}
	// 三个整页低密度版式必须点在"动手写之前"里（而不是只在组件库目录里躺着）：
	// 16 份真实 deck 里 13 份一次都没用过它们，规则不在决策处就没人会用
	for _, c := range []string{"bleed", "plate", "band"} {
		if !hasClassToken(upfront, "."+c) {
			t.Errorf("「动手写之前」里没有 .%s：三个低密度版式少一个，密度节奏就只剩一档", c)
		}
	}
	// 版式目录的结尾页 + 页面家具的提头规则
	for _, s := range []string{"封面 / 结尾页", "不能是标题的缩写或复述"} {
		if !strings.Contains(catalog, s) {
			t.Errorf("组件库那节里找不到 %q——它被压缩掉了？", s)
		}
	}
}

// 视觉审查那几条规则，压缩时同样最像"劝告"、最容易被顺手删掉，但这次的代价不对称：
//
//   - 删掉"量测只有一半"那句，模型会拿数字当"版面没问题"的证明——被裁/重叠/对比度
//     是它唯一能发现这几类问题的通道，通道没了它反而更自信。
//   - 删掉"点名页号"，模型就只剩两种行为：不看画面，或者把每一页都点一遍（那就退回了
//     这次改造要干掉的全量审查：9 页约两分钟、图片 token 是全程最贵的一项）。
//   - "最多 N 页"里的 N 派生自代码常量 maxReviewPages：模型照着提示词挑页，
//     数字对不上就会在工具那里撞一条报错，白烧一轮往返。
func TestPromptKeepsVisionReviewRules(t *testing.T) {
	flow := promptRegion(t, "# 标准工作流", "\n# 交互方式")
	for _, s := range []string{"review_slides", "pages", "1 基", "只回数字", "只有一半"} {
		if !strings.Contains(flow, s) {
			t.Errorf("「标准工作流」里找不到 %q——视觉审查那几条被压缩掉了？", s)
		}
	}
	if want := fmt.Sprintf("最多 %d 页", maxReviewPages); !strings.Contains(flow, want) {
		t.Errorf("提示词里没有 %q（上限来自代码 maxReviewPages，两处必须一致）", want)
	}
	// 挑页顺序里的三个阈值与 vision 包的判定阈值是同一组数字：
	// 提示词说"fit 低于 0.90 优先看"，而程序判定的下限也是 0.90，这样模型挑的页
	// 和报告报的页才是一回事。
	for _, s := range []string{"0.90", "31px", "1.02"} {
		if !strings.Contains(flow, s) {
			t.Errorf("挑页顺序里少了阈值 %q——模型会去点那些本来就没问题的页", s)
		}
	}
	// 防死循环的三条规则 + 配额锚点。配额数字派生自 reviewQuotaPerRun（execTool 的
	// 硬闸门），提示词与代码必须一致。这三条是实测 22 轮不收敛的 run 换来的：
	// "乒乓"禁止一次改一页审一页，"批量修复"+"确认轮"给出唯一收敛路径，
	// 被压缩掉任何一个，模型就会退回"挤牙膏审查 + 被预算强制收尾"的老路。
	if want := fmt.Sprintf("最多调 %d 次", reviewQuotaPerRun); !strings.Contains(flow, want) {
		t.Errorf("提示词里没有 %q（配额来自代码 reviewQuotaPerRun，两处必须一致）", want)
	}
	for _, s := range []string{"改一页→审一页", "批量修复", "确认轮"} {
		if !strings.Contains(flow, s) {
			t.Errorf("「标准工作流」里找不到 %q——防死循环的规则被压缩掉了？", s)
		}
	}
}

// promptRegion 取提示词里从 start 到 end 之间的一段（end 传空串则取到文件末尾）。
// 用来把断言限定在某一节里，避免"别处顺口提过一次"造成的假通过。
func promptRegion(t *testing.T, start, end string) string {
	t.Helper()
	i := strings.Index(systemPrompt, start)
	if i < 0 {
		t.Fatalf("提示词里没有 %q 这一节（改标题了？）", start)
	}
	rest := systemPrompt[i:]
	if end == "" {
		return rest
	}
	j := strings.Index(rest, end)
	if j < 0 {
		t.Fatalf("提示词里 %q 之后找不到 %q（节被删或改名了）", start, end)
	}
	return rest[:j]
}

// 发给模型的提示词不能带 \r。这条守的是一个看不见的浪费：
// 文件在 Windows 上是 CRLF，//go:embed 原样嵌字节，而那 433 个 \r 实测让
// prompt_tokens 从 7003 涨到 7794（同内容只差行尾，各发一次请求实测）——
// 每次请求白交 11%，且没有任何报错会提示你。
func TestEmbeddedPromptHasNoCR(t *testing.T) {
	if i := strings.IndexByte(systemPrompt, '\r'); i >= 0 {
		t.Errorf("发给模型的提示词里有 \\r（位置 %d）：行尾归一化那步被去掉了？"+
			"实测 CRLF 版比 LF 版多 791 个 token（7794 vs 7003）", i)
	}
}
