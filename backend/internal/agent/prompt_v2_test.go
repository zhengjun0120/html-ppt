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
	if !strings.Contains(got, "「"+long+"」") {
		t.Errorf("预检要列出原标题让 agent 对照改写")
	}
	if strings.Contains(got, "第 2 页「短标题」") {
		t.Errorf("达标标题不该出现在预检里")
	}

	// 全部达标时不注入（别白白打断前缀缓存）
	o2 := &deck.Outline{Pages: []deck.OutlinePage{{No: 1, Role: "cover", Title: "短标题"}}}
	if got2 := BuildStagePrompt(deck.StageGenerating, "deck-0001", nil, o2); strings.Contains(got2, "标题预检") {
		t.Errorf("没有超长标题时不该有预检段")
	}
}
