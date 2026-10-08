# e-ink-editorial（墨页叙事）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/e-ink-editorial`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：每份 deck 只选一套纸墨主题——墨水经典 #0A0A0B / #F1EFEA / #E8E5DE，另有靛蓝瓷、
  森林墨、牛皮纸三套；保持低饱和、柔和而清晰的对比，不混搭。
- 排版：标题 Noto Serif SC / Playfair Display，正文 Inter，页码与元数据 JetBrains Mono；
  主标题优雅不夸张，大引文可以成为一页唯一的焦点。
- 布局：大标题、引言、正文、图注与页码建立阅读秩序；留白比卡片更重要，正文页避免连续深色 hero。
- 装饰：只用 1px 细分隔线、低对比章节数字与少量线性图标；流动金属或纸面微光仅限封面/章节页。
- 禁忌：无衬线主标题、密集卡片、重阴影、强装饰打断阅读、同一主题连续三页。

## 转换说明

- 按本仓库定位取「墨水经典」纸墨主题并做墨水屏化收敛：全程无彩色，纯黑 #0A0A0B 是唯一强调色
  （上游其余三套彩色纸墨主题不并入，避免混搭）。
- 上游预览页的母题移植：folio 刊号行 → ei-folio（等宽），1px 通栏 rule → ei-rule，
  图注/出处 → ei-cap（上游 12px 提到 15px，满足本仓库最小字号纪律）。
- 纸面颗粒（feTurbulence 噪声 SVG）与右缘页边细线做成 ambient 层
  （.slide::before/::after，z-index:-1），每页自动衬底、不进骨架。
- 翻页动画设 --page-shift:0：墨水屏只有整页淡入、没有横移，贴 upstream「柔和 cascade」气质。
- 字体栈：标题与正文 Georgia/宋体，元数据 JetBrains Mono + Noto Sans SC；不新增 webfont 文件。
