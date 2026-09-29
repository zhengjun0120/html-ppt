# red-research-framework（丹朱教研）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/red-research-framework`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：主色丹朱红 #D02300，辅 #9C1A00 / #BF3C2E / #C00000 / #8F2D23；强调 #FD905F、
  #E17467、金黄 #FFC000 / #FFDC6D；白底为主，卡片可用 #F6D3CC、#F9ECEA、#FCF1F0；
  正文 #808080 / #595959。
- 排版：PingFang SC 无衬线黑体，标题 40-47px 红字或红底白字；副标题红字主标 + 灰字栏目；
  正文 20-33px 左对齐分点，关键词只用红色突出一次。
- 布局：标题栏固定左上，主体用三段流程、模块化宫格、闭环、鱼骨对比或层级框架；
  箭头、虚线、连接线只表达明确的顺序或依赖；保持白底留白，不把每个模块做成大卡片。
- 装饰：圆角矩形、圆形/胶囊标签、六边形、红描边；金黄箭头强调推进；禁具象人物、
  手绘插画、3D、霓虹、卡通与强烈光效；不削弱红底白字与红字白底的对比。

## 转换说明

- 原上游只有一个 16:9 单页预览（品牌行 + 红徽标 + 红底白字标题带 + 四节点流程图，
  节点间金黄箭头）；本模板按本仓库 deck-v2 契约重建为 9 版式骨架
  （cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 预览的四节点流程图转写为 moments 时间线（红圆点 + 金黄虚线连接）；标题带（band）
  转写为 rr-band，在 cover 与 divider 使用并约定每份 deck 至多两次；徽标（badge）与
  品牌行（brand）原样转写为 rr-badge / rr-brand。
- 数据块用上游淡红卡底 #FCF1F0 + 红左线（rr-stat），避开白卡与数据块的混淆；
  右上六边形与左下金黄虚线做成 data-URI SVG 的 ambient 层（.slide::before/::after，
  z-index:-1）。
- --bg 取 #F4E6E1 淡丹朱暖白，letterbox 衬托白页；--page-shift 24px 强化「由左到右」
  的推进感（时长 .5s，遵守 ≤650ms 截图契约）。
- 字体栈：'PingFang SC'/'Noto Sans SC' 无衬线主体 + Menlo/Consolas 等宽编号；
  不新增 webfont 文件。
