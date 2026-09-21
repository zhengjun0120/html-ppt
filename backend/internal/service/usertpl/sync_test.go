package usertpl

import "testing"

// cssCovers 是 fork 同步 layouts.md 的安全闸：骨架用到的类要么由 deck 外壳
// base.css 提供（公共类），要么出现在 fork 的 css 里——值可以漂移（用户调过
// token/旧字号档），缺类就是渲染裸奔。
func TestCssCovers(t *testing.T) {
	fork := []byte(".slide{}\n.step .n{font-weight:700}\n.dim{color:#666}")
	shell := []byte(".grid{display:grid}\n.g3{}\n.notes{}\n.mt-l{}")
	l := `<section class="slide" data-layout="how-it-works">
  <div class="grid g3 mt-l">
    <div class="step"><div class="n">1</div><p class="dim">说明</p></div>
  </div>
  <div class="notes">讲稿</div>
</section>`
	if !cssCovers(fork, shell, l) {
		t.Error("壳提供公共类、fork 提供专属类时应返回 true")
	}
	if cssCovers(fork, nil, l) {
		t.Error("壳 css 读不到时退化为全查：fork 缺公共类应返回 false")
	}
	everything := []byte(".slide{}\n.grid{}\n.g3{}\n.mt-l{}\n.step{}\n.n{}\n.dim{}\n.notes{}")
	if !cssCovers(everything, nil, l) {
		t.Error("fork css 连公共类都全有时应通过")
	}
	if cssCovers([]byte(".slide{}"), nil, l) {
		t.Error("壳与 fork 都缺 .step/.n/.dim 时应返回 false")
	}
}
