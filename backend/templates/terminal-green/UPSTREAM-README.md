# terminal-green（终端绿）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/terminal-green`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：标题 #4ade80 亮绿、正文 #86efac 浅绿，强调 #22c55e 与 #16a34a 是荧光的两种亮度；
  背景 radial-gradient（#0a1a0a → #0d280d → #0f330f）模拟 CRT 暗面；卡为
  rgba(10,26,10,0.85) 暗绿底 + rgba(34,197,94,0.3) 发光描边。
- 排版：等宽字体是唯一选择，终端美学只认等宽；代码块与命令行是核心视觉元素。
- 装饰：扫描线（repeating-linear-gradient）、辉光 text-shadow、命令行提示符（$、>）作装饰。
- 适合：CLI 工具展示、黑客马拉松、安全演讲、复古计算文化。
- 禁忌：非等宽字体、彩色（保持纯绿）、丢掉暗色沉浸感、现代 UI 元素破坏终端纯粹。

## 转换说明

- 上游只有一个 16:9 单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约重建为
  10 版式骨架（9 基础 + code，深色科技类按规范附加代码版式）。
- --bg 取背景渐变的中段值 #0d280d：letterbox 底色随页面翻动不跳色。
- 扫描线与磷光辉光做成 .slide::before/::after 的 ambient 层（z-index:-1，纯 CSS 渐变），
  每页自动衬底、不进骨架，生成侧零成本；辉光只给标题/大数字/光标，正文不加。
- 上游预览的 mac 三色圆点改为纯绿三档透明度：遵守上游「保持纯绿」纪律。
- 字体栈：'JetBrains Mono'（拉丁）+ Maple Mono CJK（中文等宽回落，随 fonts.css 自托管），
  不新增 webfont 文件。
- demo 叙事为原创（行云：一个 Go 命令行工具的诞生），数据口径写在页内来源行。
