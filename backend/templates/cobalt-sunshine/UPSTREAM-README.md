# cobalt-sunshine（明黄钴蓝 · 活力专业）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/cobalt-sunshine`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：明黄 #FFE470（高饱和暖黄）大面积底与活力基调，钴蓝 #0B2299（深蓝带紫调）作标题、
  强调与数据图表色；纯白/纯黑文字，浅灰 #F0F0F0 次要背景。
- 排版：无衬线粗体大标题（思源黑体 Bold、Arial Bold），字重大、字号大；「重对比、轻装饰」；
  黄底配蓝字、蓝底配黄字或白字，形成强反差。
- 装饰：矩形、圆形、半圆、斜线条；棋盘格与对角条纹纹理；数据图表钴蓝主色 + 明黄高亮。
- 布局：撞色对半切割或大色块交替，标题压在撞色交界处；数据用大数字 + 钴蓝强调。
- 动画：明快有冲击，cubic-bezier(0.2, 0.8, 0.2, 1)，0.4s-0.6s，色块从边缘快速滑入。
- 禁忌：低饱和灰调、柔光毛玻璃、衬线与手写体、写实照片、弱化对比。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 明黄 #FFE470 收为 #FFD338、钴蓝 #0B2299 收为 #0047AB：保留互补撞色张力的同时，
  把荧光感与紫调压回投屏与截图管线的舒适区间。
- preview 的视觉母题逐项转译：左侧蓝色斜切块 → cs-kicker 题签块与 cs-badge 徽章；
  斜纹条 → cs-stripe；右上网格圆点 → .slide::before 的 data-URI SVG ambient；
  右下旋转贴纸 → cs-tag；底部三格数据 cell → metrics 版式的 cs-cell（蓝底黄字大数字）。
- 翻页动画变量按上游缓动覆盖：--page-ease cubic-bezier(.2,.8,.2,1)、0.5s、横移 38px
  （≤650ms，符合截图管线契约）。
- 字体栈：系统无衬线黑体（Noto Sans SC / PingFang SC / Microsoft YaHei / Arial），
  不新增 webfont 文件。
