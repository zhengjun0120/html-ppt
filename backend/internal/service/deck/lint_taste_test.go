package deck

// AI 味 lint 的回归：规则要"报得准"——误报会驱使模型改好内容，漏报等于没有。

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func lintRulesOf(t *testing.T, html string) map[string][]string {
	t.Helper()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatal(err)
	}
	sec := doc.Find("body").Children().Filter("section").First()
	items := LintTaste(sec)
	out := map[string][]string{}
	for _, it := range items {
		out[it.Rule] = append(out[it.Rule], it.Word+":"+it.Msg)
	}
	return out
}

func TestLintTaste(t *testing.T) {
	t.Run("干净页面零提示", func(t *testing.T) {
		rules := lintRulesOf(t, `<section class="slide" data-layout="cards-3">
			<p class="kicker">// context</p>
			<h2 class="h2">启动不要钱，回收要你负责</h2>
			<div class="grid g3"><div class="card card-accent"><h4>标题</h4><p class="dim">说明文字。</p></div></div>
			<div class="notes">讲稿，可以随意口语。</div>
		</section>`)
		if len(rules) != 0 {
			t.Fatalf("干净页不应有提示: %v", rules)
		}
	})

	t.Run("中文禁词", func(t *testing.T) {
		rules := lintRulesOf(t, `<section class="slide"><h2 class="h2">为团队赋能的抓手</h2></section>`)
		if len(rules["T001"]) < 2 {
			t.Fatalf("赋能/抓手 都应命中: %v", rules)
		}
	})

	t.Run("英文禁词整词匹配", func(t *testing.T) {
		rules := lintRulesOf(t, `<section class="slide"><h2 class="h2">A seamless experience</h2></section>`)
		if len(rules["T002"]) != 1 {
			t.Fatalf("seamless 应命中一次: %v", rules)
		}
		// robustness 不该命中 robust（整词边界）
		rules = lintRulesOf(t, `<section class="slide"><h2 class="h2">Robustness matters</h2></section>`)
		if len(rules["T002"]) != 0 {
			t.Fatalf("robustness 不应命中 robust: %v", rules)
		}
	})

	t.Run("破折号规则", func(t *testing.T) {
		// 标题破折号
		rules := lintRulesOf(t, `<section class="slide"><h2 class="h2">并发——那些年踩的坑</h2></section>`)
		if len(rules["T003"]) == 0 {
			t.Fatal("标题破折号应提示")
		}
		// 正文超过 1 处
		rules = lintRulesOf(t, `<section class="slide"><h2 class="h2">标题</h2><p class="lede">这里——一处，那里——又一处。</p></section>`)
		if len(rules["T003"]) == 0 {
			t.Fatal("正文破折号超 1 处应提示")
		}
	})

	t.Run("假精确数字", func(t *testing.T) {
		rules := lintRulesOf(t, `<section class="slide"><h2 class="h2">标题</h2><p class="dim">用户满意度 92.5%。</p></section>`)
		if len(rules["T004"]) != 1 {
			t.Fatalf("带小数百分比应提示: %v", rules)
		}
	})

	t.Run("标题超长", func(t *testing.T) {
		rules := lintRulesOf(t, `<section class="slide"><h2 class="h2">这是一个特别特别长而且信息密度很低只是为了凑字数的页面标题写法示例</h2></section>`)
		if len(rules["T005"]) != 1 {
			t.Fatalf("超长标题应提示: %v", rules)
		}
	})

	t.Run("讲稿不算可见文本", func(t *testing.T) {
		rules := lintRulesOf(t, `<section class="slide"><h2 class="h2">标题</h2><div class="notes">这页要给团队赋能，形成闭环。</div></section>`)
		if len(rules["T001"]) != 0 {
			t.Fatalf("讲稿里的用词不应触发 T001: %v", rules)
		}
	})
}
