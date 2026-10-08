# neon-haze（深色弥散 · 氛围拉满）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/neon-haze`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：深蓝紫 #3C2466（或更深 #1A1235）做底，弥散光斑用品红 #B23A8C、电光蓝 #4A6FE3、
  霓虹粉 #E94C89、青绿 #3DD9C0，大半径 radial-gradient 做柔光晕染；
  文字近白 #F5F0FF、辅助浅紫 #B0A4E0。
- 装饰：2-4 个弥散光斑叠放、网格线/扫描线、霓虹描边与 glow；光斑斜向分布营造纵深。
- 排版：黑体粗标题字距收紧、字号大；正文细无衬线；数字用等宽字体强化科技调。
- 布局：标题压在主光斑上做视觉锚点，正文用半透明深色卡片浮在画面上；信息密度中等，
  强调氛围而非堆数据。
- 禁忌：不要浅色背景/白底；不要柔光毛玻璃（要霓虹光而非磨砂）；不要衬线、手写体；
  不要写实照片；不要把光斑做得太满太乱。

## 转换说明

- 原上游只有 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 光斑做进 ambient 层：.slide::before 承载四路 radial-gradient 弥散光（多路渐变叠加，
  不用 filter:blur，渲染成本与截图管线更稳），.slide::after 承载渐隐网格（mask 径向淡出）；
  z-index:-1，每页自动衬底、不进骨架。
- 文本级霓虹粉提亮为 #F27BA9（原 #E94C89 在深底上 4.2:1 不达 AA）：亮粉只保留在
  描边、光点与渐变按钮等大面积元素上。
- 深底纪律：正文近白/浅紫（#F5F0FF / #B0A4E0），小字 #9A8FC8（≥4.5:1）；
  letterbox 底色 --bg 取深蓝紫渐变中段值 #281C48，翻页不跳色。
- 字体栈：Noto Sans SC + Inter（黑体粗标题/细正文），标签与数字 JetBrains Mono
  （本仓库自托管字体）；不新增 webfont 文件。
