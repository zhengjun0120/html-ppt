# industrial-kaizen（表格 · 现代周报）

**来源归属**：本模板的视觉概念（配色、表格母题、报表气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/industrial-kaizen`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：暖白底 #FBFAF6 画布，行间发丝线 #EDE9DC，深墨 #1F1E1A 强线只出现在标题下沿、
  表头底部与合计行；文字三级灰阶 #1F1E1A / #2A2924 / #8B8478。
  **禁止任何彩色**——状态用「实底反白 / 浅底 / 描边」三种灰阶徽章表达。
- 排版：粗壮无衬线（Semibold+），字距收紧 -.015em；**所有数字必须 tabular-nums 等宽数字并右对齐**，
  多列数字要能连成一条竖线；表头全大写、宽字距、灰色。
- 装饰母题：几乎没有装饰。只有细分割线、深墨强线、状态徽章、▲▼ 趋势符四种元素。
  禁渐变、阴影、圆角超 2px、位图、图标、斑马纹。
- 布局：顶部标题区 + KPI 横栏（列间细竖线）+ 主表格 + 页脚的报表结构；
  「对齐是这套风格的信仰」。
- 场景：经营周报、季度复盘、KPI 看板、销售明细、项目进度汇总。

## 转换说明

- 原上游只有一个 16:9 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 上游单页的「头部 + KPI 横栏 + 表格 + 页脚」结构被拆解映射到各版式：
  metrics 版式即 KPI 横栏（上下细线封边 + 竖细线分列，▲▼ 趋势符保留在 note 行）；
  contents/keynotes/split 用发丝线行与 2px 圆角细边卡承接报表行语言；
  页顶深墨封边与右侧栏线做成 ambient 层（.slide::before/::after，z-index:-1）每页自动衬底。
- 等宽数字纪律升级为全局：模板根节点设 font-variant-numeric:tabular-nums，
  大数字、时间点、编号全部继承；关键数字必须带口径与来源（rules.md 硬性条款）。
- 字体栈：Inter + PingFang SC + Noto Sans SC（fonts.css 已自托管 Inter/Noto Sans SC）；
  不新增 webfont。
