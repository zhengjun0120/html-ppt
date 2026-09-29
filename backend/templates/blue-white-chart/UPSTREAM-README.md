# blue-white-chart（蓝白商务图表）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/blue-white-chart`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：白底 #FFFFFF 为主，局部极浅冷蓝灰 #F1F5F9 / #F3F6FC；企业深蓝 #305598 用于标题、
  分隔线、核心模块与数据柱；明亮商务蓝 #4472C4、柔和中蓝 #5B80D1、淡蓝灰 #A0B8E2 /
  #B4C6E7 / #D0DBEF 制造层级；关键数字可用深海军蓝 #00329D；红 #DC2626、绿 #16A34A、
  橙 #F59E0B 只在复杂图表中作数据线辅助。
- 字体：思源黑体或同类粗黑体标题，正文与数字大量用 Roboto；关键数字可放大成数据锚点。
- 排版：左上标题 + 一条深蓝横向分隔线，然后进入内容区；信息网格化、模块化，留白适中。
- 装饰：图表与逻辑图语言（折线、柱状、环形、矩阵、KPI 卡），细长深蓝横线、虚线框、
  圆环、方块、箭头；几何可轻微叠放但保持精确冷静。
- 禁忌：艳丽渐变、霓虹色、手绘水彩卡通、圆润可爱字体、堆叠阴影、破坏网格对齐。

## 转换说明

- 上游预览是一个高密度仪表盘单页（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing），
  把「仪表盘」拆成「一页回答一个业务问题」的报告节奏。
- 预览右上/左下的装饰圆与标题下深蓝分隔线转写为 ambient 层（.slide::before/::after，
  z-index:-1）与 bw-rule（横贯分隔线），每页自动衬底、不进骨架。
- 「左侧深蓝重点区」的仪表盘手法收进 bw-n / bw-tl-dot / bw-badge 的深蓝实底编号块，
  保持版式骨架与黄金样例同构。
- 字体栈：'Roboto' / 'Inter' / 'Noto Sans SC'（Roboto 本机存在则生效，缺省由 Inter 兜底），
  不新增 webfont 文件。
