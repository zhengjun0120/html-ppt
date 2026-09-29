# blueprint（蓝图）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/blueprint`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：标题浅蓝白 #dbeafe、正文中蓝 #93c5fd，强调亮蓝 #60a5fa 与制图蓝 #3b82f6；
  深蓝三段渐变如传统蓝图纸张，卡片暗蓝底 + 浅蓝工程描边。
- 排版：等宽或技术感字体，深蓝底上白/浅蓝文字；连接线与箭头是排版的一部分。
- 布局：网格底纹贯穿每一页，卡片像蓝图上的标注框，方正严谨；架构图、流程图是天然元素。
- 动画：path-draw、fade-up、stagger-list。
- 禁忌：不用暖色调或圆角、不丢网格底纹、不用非等宽字体、不让深底文字对比度不足。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 纸底收深为 #0D2A52（取代上游亮蓝渐变的右端），保证正文 #93C5FD 对比度 ≥7:1，
  投屏与截图管线下更稳；网格与双线图框做成 CSS data-URI SVG 的 ambient 层
  （.slide::before/::after，z-index:-1），每页自动衬底、不进骨架，生成侧零成本。
- 全站等宽字体栈 'JetBrains Mono'（fonts.css 自托管，CJK 由同族中日韩分段覆盖）；
  不新增 webfont 文件。
- 身份元素：图签章 bp-stamp（双线方章）、图例框头 bp-box-head、尺寸标注线 bp-dim；
  全部方角、无阴影、无暖色，遵守上游禁忌清单。
- demo 叙事为「河湾社区图书馆建筑设计提案」，面积、层数、造价、工期均带口径来源行。
