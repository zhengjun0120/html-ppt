# editorial-serif（杂志衬线）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/editorial-serif`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：标题 #431407 深褐、正文 #7c2d12 焦棕，强调 #ea580c 深橙与 #f97316 亮橙是
  「编辑的红色铅笔」；背景 #fef7ed → #fff7ed → #fffbeb 奶油三段渐变，卡为
  rgba(255,250,245,.9) 带极淡棕描边 rgba(120,53,15,.18)。
- 排版：衬线标题 + 无衬线正文的经典杂志编排（预览页正文实为衬线，本模板从预览页）；
  行距宽松，段落可缩进；引号与 pull quote 用更大衬线突出。
- 布局：大标题、副标题、正文层次分明；首字母放大做装饰；引文区块用左线与正文区分；
  栏与栏之间用竖色线分隔（预览页 .col-rule）。
- 禁忌：无衬线标题、冷色调、密集排版、装饰打断叙事节奏。

## 转换说明

- 上游预览页的母题全部组件化：mono 标签（.tag）→ es-tag，橙色短线（.line-accent）→ es-rule，
  首字下沉（.dropcap）→ es-drop，左线引文块（.quote-block）→ es-quote-block，
  斜体英文（.english）→ es-en，竖栏线（.col-rule）→ 左缘页边栏线 ambient。
- 页边栏线与刊尾花饰（SVG 曲线加菱点）做成 ambient 层（.slide::before/::after，
  z-index:-1），每页自动衬底、不进骨架。
- letterbox 用纸渐变中段值 #FDF4E3，翻页不跳。
- 字体栈：全衬线 Georgia/宋体（标题与正文），标签与出处 JetBrains Mono + Noto Sans SC；
  不新增 webfont 文件。
