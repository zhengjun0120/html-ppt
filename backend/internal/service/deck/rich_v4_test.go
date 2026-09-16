package deck

import (
	"path/filepath"
	"testing"
)

// richDeckV4 用满 v4 的每一个新组件类，且**只用组件 class + 引用主题变量的内联 style**
// （图框比例、bars 柱长这类允许的内联）。它存在的理由：v4 加了三十多个类与一批翻色 token，
// 只要有一个的 var() 消费没进变量契约、或某个写法触雷 stylelint 的重复内联守卫，
// 生产 write_deck 路径就会拒写或报警——而这些只有在真跑一遍 normalize+回声校验时才暴露。
const richDeckV4 = `<section class="cover bg-ink">
  <span class="page-head"><span>留住 1.5℃</span><span>SS · 26.09 · 01</span></span>
  <h1 class="cover-title">一页需要一个落点</h1>
  <div>
    <p class="cover-sub">气候账单，从一杯冷却的咖啡说起</p>
    <span class="tag">校园科普</span>
  </div>
  <span class="page-foot"><span>留住 1.5℃</span><span>1.5℃ 之后呢</span></span>
</section>
<section class="bg-accent center">
  <h1>把结论压在强调色上</h1>
  <p>这是 <span class="mark">高亮的几个字</span> 与 <span class="badge">徽章</span></p>
</section>
<section class="stretch">
  <h2>账单的前三行</h2>
  <div class="grid-3">
    <div class="kpi"><span class="kpi-label">海平面</span><p class="kpi-val">3.4 <span class="unit">mm/年</span></p><p class="kpi-note">近 30 年翻倍</p></div>
    <div class="kpi"><span class="kpi-label">盛夏天数</span><p class="kpi-val">+47 <span class="unit">天</span></p><p class="kpi-note">以 1981 为基线</p></div>
    <div class="kpi"><span class="kpi-label">电力峰值</span><p class="kpi-val">18 <span class="unit">%</span></p><p class="kpi-note">夏季同比抬升</p></div>
  </div>
  <figure class="figure">
    <div class="figure-body" style="aspect-ratio:21/6"></div>
    <p class="caption">1900 以来温度距平</p>
  </figure>
  <div class="dots-deco"></div>
</section>
<section class="half">
  <div class="half-col"><h2>适应</h2><p class="sub">接受新的常态</p></div>
  <div class="half-col"><h2>缓解</h2><p class="sub">改变原因本身</p></div>
</section>
<section>
  <h2>四次关键节点</h2>
  <ul class="timeline-v">
    <li><p class="label">1988</p><p class="sub">IPCC 成立</p></li>
    <li class="tl-hi"><p class="label">2015</p><p class="sub">巴黎协定写入 1.5℃</p></li>
  </ul>
  <div class="duo">
    <div class="duo-col"><p class="term">升温情景</p><p class="sub">+2.7℃ 路径</p></div>
    <div class="duo-col"><p class="term">承诺情景</p><p class="sub">+1.9℃ 路径</p></div>
  </div>
  <div class="split w8-4">
    <div class="ink-block"><p class="label">墨块内的行清单</p><div class="rows"><div class="row"><p class="sub">自动反色的发丝线</p></div></div></div>
    <div><p class="stat">68 <span class="unit">%</span></p><p class="sub">达标城市占比</p></div>
  </div>
  <div class="bars">
    <div class="bar"><p>煤电</p><div class="track"><i class="fill" style="width:100%"></i></div><span class="val">42%</span></div>
    <div class="bar"><p>工业</p><div class="track"><i class="fill" style="width:72%"></i></div><span class="val">30%</span></div>
  </div>
  <svg class="ico ico-lg" viewBox="0 0 24 24"><path d="M12 2v20M2 12h20"/></svg>
  <span class="page-no outline">06</span>
</section>
<section class="plate">
  <p class="stat">02</p>
  <h2>第二部分：怎么办</h2>
  <p class="sub">从个人选择到公共政策</p>
</section>`

// 走的是 write_deck 生产校验的同一条链：normalizeSections（消毒 + stylelint）
// → renderSkeleton → checkInlineStyleVars（变量契约回声）。三段任一不过都说明
// v4 的某个类/变量在真实提交里会被拒或被警告，而提示词恰恰叫模型用它。
func TestRichV4DeckPassesProductionValidation(t *testing.T) {
	s := New(t.TempDir(), filepath.Join("..", "..", "..", "web", "assets"), nil)

	normalized, count, warning, err := s.normalizeSections(richDeckV4)
	if err != nil {
		t.Fatalf("v4 满配 deck 被 normalizeSections 拒写：%v", err)
	}
	if count != 6 {
		t.Fatalf("应有 6 页，实际 %d", count)
	}
	if warning != "" {
		t.Errorf("v4 满配 deck 触发 stylelint 警告（说明有该被组件化的写法被判重复内联）：\n%s", warning)
	}

	// 变量契约回声：用真实主题（含 v4 翻色 token 的 renderThemeCSS 输出）渲染骨架，
	// 再校验内联样式里的 var() 是否都有消费方。
	rendered, err := renderSkeleton("v4 满配", normalized)
	if err != nil {
		t.Fatalf("renderSkeleton 失败：%v", err)
	}
	if err := s.checkInlineStyleVars(normalized, rendered); err != nil {
		t.Errorf("v4 满配 deck 未过变量契约回声校验：%v", err)
	}
}

// 翻色 token 必须被组件库消费（进变量契约）——否则模型在 vars 里改配色时，
// renderThemeCSS 输出的这些名字没有任何 var() 引用它们，回声校验会判"定义了没人用"。
func TestFlipTokensAreConsumed(t *testing.T) {
	contract, err := loadVariableContract(filepath.Join("..", "..", "..", "web", "assets"))
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"--ink-bg", "--ink-fg", "--ink-accent", "--ink-muted", "--ink-hairline", "--ink-panel",
		"--page-accent", "--page-fg", "--accent-muted", "--accent-hairline", "--accent-panel",
		"--hero-weight", "--fs-hero",
	} {
		if !contract[name] {
			t.Errorf("翻色/大字 token %s 没有被组件库 var() 消费（未进变量契约）", name)
		}
	}
}
