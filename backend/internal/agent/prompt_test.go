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
// 末尾三条是**压缩提示词时的守卫**：压缩真正会造成的损伤不是"话变少了"，
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
	msg := svc.buildSystemMessage("deck-0001")

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
		msg := svc.buildSystemMessage(bad)
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
	for _, s := range []string{"200 字", "≤ 5 个", "1094×610", "8~14 页"} {
		if !strings.Contains(systemPrompt, s) {
			t.Errorf("提示词里找不到容量预算的关键数字 %q——它被压缩掉了？", s)
		}
	}
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
