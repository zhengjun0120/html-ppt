# glassmorphism（磨砂玻璃）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/glassmorphism`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：标题 #0F172A 深墨、正文 #334155 石板灰；强调色 #818CF8 柔紫与 #06B6D4 青蓝是光斑
  的两种颜色；背景为薰衣草 #C7D2FE → #E0E7FF → 浅青 #CFFAFE 的三段渐变；卡片 rgba(255,255,255,0.46)
  高度透明磨砂，边框 rgba(255,255,255,0.55) 半透明白描边。
- 排版：现代无衬线、细字重，标题轻盈地浮在磨砂卡片上；文字与背景的对比靠 backdrop-filter 而非纯色底。
- 布局：磨砂卡片浮在多色渐变上，blur 创造深度，卡片可重叠，大圆角、柔和阴影保持轻盈。
- 动画关键词：blur-in、shimmer-sweep、gradient-flow。
- 禁忌：不用实色背景或硬边框；不丢磨砂模糊的核心效果；不让卡片变得不透明；不过度填满磨砂的呼吸空间。

## 转换说明

- 原上游只有一个单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约重建为 9 版式骨架
  （cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 文字级强调色各加深一档：柔紫 #818CF8 → #4338CA、青蓝 #06B6D4 → #0E7490，
  保证玻璃卡上深墨文字体系之外的彩色文字在投屏与截图管线下仍满足对比度。
- 光斑与顶部光带做成多层 radial/linear 渐变 + data-URI SVG 的 ambient 层
  （.slide::before/::after，z-index:-1），每页自动衬底、不进骨架，生成侧零成本；
  不用 filter:blur 大面积滤镜，避免渲染管线开销。
- --bg 取渐变中段值 #DCE3F8，letterbox 与翻页过渡不跳色。
- 字体栈：Inter / SF Pro / Segoe UI / PingFang SC / Noto Sans SC，无新增 webfont 文件。
- demo 叙事「云雾：一个毛玻璃组件库的发布」：组件数、star 数、框架适配数的口径
  均标注为 v2.0 发布说明与 GitHub 公开数据。
