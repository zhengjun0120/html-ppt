// deck 初始化脚本：所有 deck 共用这一份（第①层：框架，LLM 不可见不可改）。
// 将来交互组件（翻转卡片、拖拽排序等）的行为绑定也统一加在这里。

// 翻页动画不在 CSS 变量体系里，从 deck 内的主题配置块读取；
// 配置缺失/损坏时退回默认，框架层永不因配置挂掉
let transition = 'slide';
try {
  transition = JSON.parse(document.getElementById('deck-theme').textContent).transition || transition;
} catch (e) {}

Reveal.initialize({
  hash: true,       // URL 带 #/2 页码：刷新不丢位置，截图也能精确定位到某页
  transition: transition,
  controls: true,
  progress: true,
});
