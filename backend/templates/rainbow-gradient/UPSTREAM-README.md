# rainbow-gradient（彩虹渐变）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/rainbow-gradient`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：标题 #1E3A5F 海军蓝、正文 #374151 钢灰，强调 #F43F5E 玫红与 #8B5CF6 紫；
  七彩流动渐变如节日彩带飘过白底，卡片 rgba(255,255,255,.9) 白底浮在色彩之上。
- 排版：友好无衬线、字重中等，标题温暖不厚重，层级分明但调性轻松。
- 布局：彩虹渐变作背景或分隔元素贯穿全页，白色卡片浮在色彩之上，圆角柔和；
  欢快但不混乱，保持节奏。
- 动画：zoom-pop、stagger-list、confetti-burst。
- 禁忌：暗色或严肃元素、彩虹变混乱、过多装饰、密集布局压缩庆祝空间。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 对比度决策：整幅背景用上游的柔彩七段渐变（opacity 降半作 wash），正文一律海军蓝/
  钢灰墨色、长文字进白卡；大数字改用玫红 #E11D48→紫 #7C3AED 双色渐变字——上游
  七彩直接上文字会让黄色段对白底不足 2:1，改后两端对白底均 ≥4.5:1。
- 柔彩 wash、顶部七彩丝带与页底彩虹浪做成 CSS data-URI SVG 的 ambient 层
  （.slide::before/::after，z-index:-1），每页自动衬底、不进骨架，生成侧零成本。
- 字体栈：PingFang SC / Microsoft YaHei / Noto Sans SC；不新增 webfont 文件。
