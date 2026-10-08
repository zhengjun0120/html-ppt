# academic-paper（学术论文）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/academic-paper`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：标题墨黑 #171717、正文深灰 #404040；学术链接蓝 #2563EB / #3B82F6；背景
  linear-gradient(145deg, #fafafa, #f5f5f5 55%, #e5e5e5)，如打印纸；卡片 rgba(255,255,255,.95)
  白底 + rgba(0,0,0,.1) 极淡描边。
- 排版：衬线正文 + 无衬线标题（本模板取衬线标题的期刊版式），行距宽松、段落分明；
  脚注与引用用更小字号，图表标题与图注风格化。
- 布局：单栏/双栏论文排版，标题居中、作者在标题下、摘要区用卡片区分；图表与公式是
  视觉核心，编号清晰。
- 禁忌：非正式与彩色、装饰性字体、省略图表标注与引用、把严谨做成枯燥。

## 转换说明

- 原上游只有一个 16:9 单页预览（期刊首页母题：页顶蓝条、期刊头、摘要卡、双栏正文与
  柱状图）；本模板按本仓库 deck-v2 契约重建为 9 版式骨架
  （cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 摘要卡（abs-box）转写为通用纸面卡 ap-card，在 contents/keynotes/split/cover 复用；
  图注（fig-cap）转写为 ap-figcap，承担「注：来源——」的图表注职责。
- 页顶蓝条与纸面渐变做成 .slide::after / .slide::before（z-index:-1）的 ambient 层，
  --bg 取渐变中段值 #EFEFEF，翻页时 letterbox 不跳。
- 学术蓝语义收紧为「可溯源」：只给编号、数字、摘要标签、强调药丸与按钮，不给标题。
- 字体栈：Georgia/'Songti SC' 衬线主体 + Consolas 等宽元数据；不新增 webfont 文件。
