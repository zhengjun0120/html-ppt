# blue-orange-analytics（蓝橙经营图谱）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/blue-orange-analytics`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：白底 #FFFFFF 与 #F2F2F2；主色藏蓝系 #1D4273 / #305598 / #4474C5 / #5C7DB9 / #778EAB / #99B0E3；
  强调色赭红 #A13422 与 #C7857A 只用于突出异常、目标或关键节点；正文 #000000，
  表格线 #D1D5DB，关键边界 #0F213A。
- 排版：SF Pro Text、PingFang SC、Helvetica Neue 等系统无衬线；标题 46-60px 以黑为主；
  KPI 仅在需要强调结论时使用蓝色实底和白色粗体。
- 布局：顶部结论、主体图表、底部 KPI 的阅读顺序；卡片与色条严格对齐并保留真实留白。
- 装饰：只使用图表、色带、圆角标签、细描边表格和必要的连接线，装饰必须表达数据关系。
- 禁忌：大面积深色背景、霓虹色、手绘/卡通插画、花哨字体、强阴影、失控大圆角。

## 转换说明

- 上游只有一个单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 上游预览的顶部状态条与柱状图母题转写为页顶藏蓝数据色带与页底柱状剪影
  （.slide::before/::after，z-index:-1，纯装饰、不承载任何数据），每页自动衬底。
- 等宽标签（JetBrains Mono，自托管）承接上游预览的 monospace 读数气质；
  KPI 藏蓝实底 + 白色大数字保留「结论级强调专用」的上游纪律。
- 赭红按上游纪律收进 bo-status（状态牌）单点使用，不作任何大面积色块。
