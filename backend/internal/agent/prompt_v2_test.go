package agent

import (
	"strings"
	"testing"

	"html-ppt/backend/internal/service/deck"
)

// TestPromptTitlePrecheck 大纲标题超过上限时，生成阶段提示词要点名超长页。
// T005 在写页时才报，而大纲是模板选定之前写的、看不到这条约束——deck-0061
// 首批 4 页全部命中标题超长，返工从第一批就开始。预检让 agent 首批就写短。
func TestPromptTitlePrecheck(t *testing.T) {
	max := deck.TitleMax()
	long := "GitHub 入门：从注册到第一次 Pull Request"
	o := &deck.Outline{Pages: []deck.OutlinePage{
		{No: 1, Role: "cover", Title: long},
		{No: 2, Role: "content", Title: "短标题"},
	}}
	got := BuildStagePrompt(deck.StageGenerating, "deck-0001", nil, o)
	if !strings.Contains(got, "标题预检") || !strings.Contains(got, "第 1 页") {
		t.Errorf("超长标题该被点名（%d 字上限）", max)
	}
	if !strings.Contains(got, "逐字母算") {
		t.Errorf("标题计数口径要常驻写死——deck-0062 实测 agent 按「英文单词算 1 个词」自数，首写仍超限")
	}
	if !strings.Contains(got, "可拆成主标「GitHub 入门」") {
		t.Errorf("冒号能拆的要给现成拆法（主标 + lede 承接）")
	}
	if strings.Contains(got, "第 2 页「短标题」") {
		t.Errorf("达标标题不该出现在预检里")
	}

	// 没有超长页时不注入预检，但计数口径的规则行常驻（别让 agent 猜上限）
	o2 := &deck.Outline{Pages: []deck.OutlinePage{{No: 1, Role: "cover", Title: "短标题"}}}
	got2 := BuildStagePrompt(deck.StageGenerating, "deck-0001", nil, o2)
	if strings.Contains(got2, "标题预检") {
		t.Errorf("没有超长标题时不该有预检段")
	}
	if !strings.Contains(got2, "页标题硬上限") {
		t.Errorf("标题规则应常驻（没有超限页也要让 agent 知道上限与口径）")
	}
}
