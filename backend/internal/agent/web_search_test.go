package agent

// 联网搜索的测试。这里守两件事：
//
//  1. **协议字段 → 我们的结构体** 这个映射是对的。它是会静默坏掉的地方：
//     字段一改名，手写的结构体只会读到零值，来源悄悄变空，而页面上看不出任何异常。
//  2. **features 开关关掉时工具必须彻底消失**（在 func_tool_test.go 里，与 custom_css 同一条测试）。
//
// 为什么不用"真发一次搜索"来测第 1 条：真实搜索要网络、要计费（每次调用都真的花搜索次数），
// 而且子模型自己决定发什么 query、结果每次都不同——实测同一个问题两次，它先发的 query 分别是
// "科技新闻 2026年9月14日" 和 "tech news September 14 2026"（结果里那条新闻倒是一致的）。
// 那种测试只能验"没报错"，验不了映射本身。所以映射用合成响应离线锁住，
// 真实的连通性另有一条用 LLM_API_KEY 门控的集成测试（见文件末尾）。

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	anthropic "github.com/anthropics/anthropic-sdk-go"
)

// 合成响应里刻意放进了四种"手写解析容易漏"的情况：
// 缺 title 的结果项、重复 URL、工具级错误码（不是 HTTP 错误，夹在结果块里回来）、
// stop_reason=max_tokens（服务端截断）。这几样在真实响应里都出现过或被文档确认存在。
func TestParseWebSearchMsg(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("testdata", "web_search_resp.json"))
	if err != nil {
		t.Fatal(err)
	}
	var msg anthropic.BetaMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		t.Fatalf("合成响应解析失败（SDK 的协议类型变了？）: %v", err)
	}

	got := parseWebSearchMsg("提问", &msg)

	// 计费口径来自 usage.server_tool_use.web_search_requests：
	// 一次工具调用可能对应多次计费搜索，这个数字是唯一能对账的地方
	if got.Searches != 2 {
		t.Errorf("计费搜索次数应为 2，实际 %d", got.Searches)
	}

	// 三条结果里只有两个不同 URL（第三条重复第一条）→ 去重后 2 条
	if len(got.Sources) != 2 {
		t.Fatalf("来源应按 URL 去重后剩 2 条，实际 %d：%+v", len(got.Sources), got.Sources)
	}
	if got.Sources[0].URL != "https://a.example" || got.Sources[0].Title != "标题A" {
		t.Errorf("第 1 条来源不对: %+v", got.Sources[0])
	}
	// 缺 title 的那条要**保留**（URL 本身仍然可用）且标题为空。
	// 这正是手写 struct 判不出来的地方：它无法区分"字段是空串"和"字段根本没来"，
	// 而 respjson 的 Valid() 可以——所以这条断言同时守住了"我们没有偷偷把它丢掉"
	if got.Sources[1].URL != "https://b.example" || got.Sources[1].Title != "" {
		t.Errorf("缺 title 的来源应保留且标题为空: %+v", got.Sources[1])
	}

	// 工具级错误码不是 HTTP 错误。漏掉它，模型就分不清"搜了没搜到"和"搜失败了"，
	// 而这两件事的下一步完全不同：前者换问法，后者不该在页面上写具体数字
	if got.ToolError != "max_uses_exceeded" {
		t.Errorf("工具级错误码应为 max_uses_exceeded，实际 %q", got.ToolError)
	}

	// 全部 text 块都要拼进来（第一段是过渡语、第二段是答案）
	if !strings.Contains(got.Answer, "I'll search for that.") || !strings.Contains(got.Answer, "答案正文") {
		t.Errorf("答案应包含全部 text 块，实际 %q", got.Answer)
	}
	// thinking 块绝不能进答案：把模型的思考当事实写进 deck 是很常见的一类错，
	// 而 content 数组里 thinking 与 text 是混在一起的，只按"哪个字段非空"取就会中招
	if strings.Contains(got.Answer, "我应该先搜索") {
		t.Error("thinking 块不该进入答案")
	}

	// stop_reason=max_tokens 是**服务端**截断。我们自己截答案/来源时会置 Truncated，
	// 但服务端截断同样必须让主模型知道——否则半句话看起来是完整的
	if !got.Truncated {
		t.Error("stop_reason=max_tokens 应标记为截断")
	}
}

// anthropicBase 决定联网搜索打到哪里。三种输入都要成立：
// 显式配置优先、从 base_url 推导、base_url 带 /v1 时不能拼成 .../v1/anthropic/v1/messages。
func TestAnthropicBase(t *testing.T) {
	cases := []struct {
		name     string
		explicit string
		baseURL  string
		want     string
	}{
		{"显式配置优先", "https://gw.example/anthropic", "https://api.deepseek.com", "https://gw.example/anthropic"},
		{"显式配置去掉尾部斜杠", "https://gw.example/anthropic/", "", "https://gw.example/anthropic"},
		{"从 base_url 推导", "", "https://api.deepseek.com", "https://api.deepseek.com/anthropic"},
		{"base_url 带 /v1 时不重复", "", "https://api.deepseek.com/v1", "https://api.deepseek.com/anthropic"},
		{"base_url 带尾部斜杠", "", "https://api.deepseek.com/", "https://api.deepseek.com/anthropic"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a := &AgentService{AnthropicBaseURL: c.explicit, BaseURL: c.baseURL}
			if got := a.anthropicBase(); got != c.want {
				t.Errorf("anthropicBase() = %q，期望 %q", got, c.want)
			}
		})
	}
}

// max_uses 的夹取：模型对"多搜几次更保险"没有成本意识，所以它要多少不能给多少。
// 每种越界都要收敛到明确的边界值，不能原样传下去（每一次都是计费搜索）。
func TestClampMaxUses(t *testing.T) {
	cases := []struct{ in, want int }{
		{0, webSearchDefaultMaxUses},  // 没传
		{-5, webSearchDefaultMaxUses}, // 传了负数（模型偶尔会）
		{2, 2},                        // 正常值原样保留
		{webSearchMaxUses, webSearchMaxUses},
		{999, webSearchMaxUses}, // 越界夹到硬上限
	}
	for _, c := range cases {
		if got := clampMaxUses(c.in); got != c.want {
			t.Errorf("max_uses=%d 应夹成 %d，实际 %d", c.in, c.want, got)
		}
	}
}

// 联网集成测试：验的是"这个外部服务今天还认这套协议"，不是仓库里的不变量——
// 所以它用环境变量门控、默认跳过是合理的。但要保证它至少被人手工跑通过一次
// （首次接入时实测：stop_reason=end_turn、web_search_requests=2、每次搜索 10 条结果）。
func TestWebSearchLiveIntegration(t *testing.T) {
	key := os.Getenv("LLM_API_KEY")
	if key == "" {
		t.Skip("未设置 LLM_API_KEY，跳过联网集成测试")
	}
	svc := &AgentService{ModelID: "deepseek-flash", BaseURL: "https://api.deepseek.com"}
	// 刻意把"请求的 max_uses"和"实际计费的 searches"一起打出来：实测这个兼容端点
	// **不执行 max_uses**（请求 1 得到 2、请求 2 得到 4），所以两者对不上是**预期现象**，
	// 不是这条测试失败。这里是唯一能看见成本真实口径的地方。
	const maxUsesWanted = 2
	res, err := svc.runWebSearch(t.Context(), key, "今天的日期是几号", maxUsesWanted)
	if err != nil {
		t.Fatalf("联网搜索失败: %v", err)
	}
	// 只断言结构，不断言内容：子模型每次发的 query 与拿到的结果都不同
	if res.Searches == 0 || len(res.Sources) == 0 {
		t.Fatalf("没有产生搜索或来源，协议可能变了: %+v", res)
	}
	t.Logf("请求 max_uses=%d，实际 searches=%d，sources=%d，answer=%d 字",
		maxUsesWanted, res.Searches, len(res.Sources), len([]rune(res.Answer)))
}
