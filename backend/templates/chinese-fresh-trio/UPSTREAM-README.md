# chinese-fresh-trio（中国传统色·青绿湖蓝）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/chinese-fresh-trio`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：翠竹绿 #7BCFA6 与湖蓝 #66A9C9 双主色，米白 #E6F5F2 / 淡青白 #F0F8F6 底，
  墨黑 #333333 主字、深炭灰 #4A4A4A 正文，朱红 #C0392B 只用于关键数据与印章，
  淡米 #D9D2C5 做三色条的第三段中性色。
- 意象：山峦、溪流、青竹、云气、几何色块；装饰用三色色块拼接、细线分割、圆点、印章。
- 排版：标题思源宋体加粗，副标题楷体，正文思源黑体行高 1.7，数字与英文用现代无衬线。
- 布局：三色色块明确分区或穿插，标题大字 + 三色色卡横排，留白与色块节奏交替。
- 禁忌：不要超过四种主色；不要高饱和荧光、暗黑背景；不要复杂渐变、金属质感；不要拥挤排版。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 上游预览顶部的 flex 三色条（trio-bar）保留为骨架内真实元素（ft-band，flex 三段色带），
  每页第一行照抄；右下几何山峦折线改为 CSS data-URI SVG 的 ambient 层
  （.slide::after，z-index:-1），每页自动衬底、不进骨架，生成侧零成本。
- 实色三联卡沿用上游预览的绿/蓝/墨三色与白字（读感按上游 palette 展示卡处理，
  卡内正文固定 18px 白字，不作正文长段落容器——长说明交给 ft-panel 细线卡）。
- 色值全部沿用上游原值，未做收暗调整；朱红 #C0392B 在投屏下实测观感沉稳，维持原样。
- 字体栈：标题 Georgia/宋体栈，正文黑体回退 Noto Sans SC，数字 Inter；
  不新增 webfont 文件。
