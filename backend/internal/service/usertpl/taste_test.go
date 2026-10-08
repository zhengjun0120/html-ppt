package usertpl

// AI 味 lint 定制侧接入的回归：lintDemoSections 逐页判定与折叠、write_demo /
// add_layout 的工具返回带提示、system prompt 品味纪律注入。
// 判定本身（词表/规则）由 deck 包的 lint_taste_test 覆盖，这里只测接线与形态。

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func demoWithMarkers(sections ...string) string {
	return "<html><body><div class=\"deck\">\n<!-- SLIDES:START -->\n" +
		strings.Join(sections, "\n") +
		"\n<!-- SLIDES:END -->\n</div><script src=\"/assets/deck-v2/runtime.js\"></script></body></html>"
}

func aiFlavoredSection() string {
	return `<section class="slide" data-layout="a"><p class="kicker">k</p><h2 class="h2">标题</h2><p class="dim">赋能所有人，99.99% 可靠</p></section>`
}

func cleanSection() string {
	return `<section class="slide" data-layout="b"><p class="kicker">k</p><h2 class="h2">干净页</h2><p class="dim">这一页说的是具体的事，有真实的细节与数字口径。</p></section>`
}

func TestLintDemoSections(t *testing.T) {
	hints := lintDemoSections(demoWithMarkers(aiFlavoredSection(), cleanSection()))
	joined := strings.Join(hints, "\n")
	if !strings.Contains(joined, "第 1 页 [T001]") || !strings.Contains(joined, "赋能") {
		t.Errorf("第 1 页缺 T001 禁词提示: %q", joined)
	}
	if !strings.Contains(joined, "[T004]") || !strings.Contains(joined, "99.99%") {
		t.Errorf("第 1 页缺 T004 假精确数字提示: %q", joined)
	}
	if strings.Contains(joined, "第 2 页") {
		t.Errorf("干净的第 2 页不应有提示: %q", joined)
	}

	// 超上限折叠：一个 section 塞满禁词（15 个中文禁词全部命中）
	banned := []string{"赋能", "抓手", "闭环", "沉淀", "心智", "护城河", "组合拳", "打法",
		"底层逻辑", "顶层设计", "颗粒度", "拉齐", "值得注意的是", "综上所述", "总而言之"}
	spam := `<section class="slide" data-layout="a"><h2 class="h2">` + strings.Join(banned, "") + `正文</h2></section>`
	hints = lintDemoSections(demoWithMarkers(spam))
	// 15 处 T001 + 1 处 T005（塞满禁词的标题必然超 16 字上限）
	if total := len(banned) + 1; len(hints) != maxTasteHints+1 || !strings.Contains(hints[maxTasteHints], fmt.Sprintf("共 %d 处", total)) {
		t.Fatalf("折叠行应报共 %d 处，得到 %d 行、尾行 %q", total, len(hints), hints[len(hints)-1])
	}
}

func TestLintDemoFragment(t *testing.T) {
	hints := lintDemoFragment(aiFlavoredSection())
	if len(hints) == 0 || !strings.Contains(hints[0], "示例页 [T001]") {
		t.Errorf("单片段提示的 label 不对: %v", hints)
	}
	if hints := lintDemoFragment(cleanSection()); len(hints) != 0 {
		t.Errorf("干净片段不应有提示: %v", hints)
	}
}

func TestWriteDemoTasteHints(t *testing.T) {
	s, _ := structTestService(t)
	const uid = 31
	row, err := s.CreateBlank(uid, "demo 味")
	if err != nil {
		t.Fatalf("CreateBlank: %v", err)
	}
	row, err = s.GetOwned(uid, row.ID)
	if err != nil {
		t.Fatalf("GetOwned: %v", err)
	}

	// 现有 demo 里注入禁词（文本级替换，data-layout 不动 → 重挂照常过）
	idx, err := os.ReadFile(filepath.Join(s.Dir(row.ID), "index.html"))
	if err != nil {
		t.Fatalf("读 index.html: %v", err)
	}
	dirty := strings.Replace(string(idx), "{{一句话定位}}", "赋能所有人的定位", 1)
	if dirty == string(idx) {
		t.Fatal("占位符未找到，注入失败")
	}
	args, _ := json.Marshal(map[string]string{"html": dirty})
	result, dirtyFlag := s.execCustomTool(row, "write_demo", string(args))
	if !dirtyFlag {
		t.Error("write_demo 应记 dirty")
	}
	if !strings.Contains(result, "AI 味提示") || !strings.Contains(result, "[T001]") {
		t.Errorf("工具返回缺味提示: %q", result)
	}
}

func TestAddLayoutTasteHint(t *testing.T) {
	s, _ := structTestService(t)
	const uid = 32
	row, err := s.CreateBlank(uid, "加版式味")
	if err != nil {
		t.Fatalf("CreateBlank: %v", err)
	}
	spec := flowSpec("blank-flow2")
	spec.DemoHTML = strings.Replace(spec.DemoHTML, "三步完成接入", "赋能三步走", 1)
	summary, err := s.AddLayout(uid, row.ID, spec)
	if err != nil {
		t.Fatalf("AddLayout: %v", err)
	}
	if !strings.Contains(summary, "AI 味提示") || !strings.Contains(summary, "[T001]") {
		t.Errorf("摘要缺味提示: %q", summary)
	}
}

func TestCustomizePromptTaste(t *testing.T) {
	s, _ := structTestService(t)
	const uid = 33
	row, err := s.CreateBlank(uid, "味")
	if err != nil {
		t.Fatalf("CreateBlank: %v", err)
	}
	prompt := s.customizeSystemPrompt(row)
	for _, want := range []string{"品味纪律", "em-dash", "99.99%", "纯黑 #000", "三等分"} {
		if !strings.Contains(prompt, want) {
			t.Errorf("system prompt 缺品味纪律项 %q", want)
		}
	}
}
