// deck 初始化脚本：所有 deck 共用这一份（第①层：框架，LLM 不可见不可改）。
// 将来交互组件（翻转卡片、拖拽排序等）的行为绑定也统一加在这里。
Reveal.initialize({
  hash: true,       // URL 带 #/2 页码：刷新不丢位置，截图也能精确定位到某页
  transition: 'slide',
  controls: true,
  progress: true,
});
