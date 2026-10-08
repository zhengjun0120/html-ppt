# engineering-whiteprint（工程白图）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/engineering-whiteprint`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：标题海军蓝 #1e3a5f、正文钢灰 #374151，强调深蓝 #1e40af 与亮蓝 #2563eb 两档墨线；
  背景近纯白的三段微灰渐变，卡片白底 + 淡蓝工程描边。
- 排版：等宽字体是唯一选择——系统设计需要等宽的精确感；数据表格、API 路径、架构描述自然融入。
- 布局：网格底纹贯穿每一页，内容沿网格线排列；卡片像工程图纸上的标注框，方正规矩。
- 动画：fade-up、stagger-list、typewriter。
- 禁忌：不用非等宽字体、不用装饰性元素、不用暖色或活泼配色、不让网格底纹消失。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 10 版式骨架（基础 9 版式 + code 代码示例版式，属深色科技/工程终端类的附加契约）。
- 坐标网格（浅蓝细线 + 主线）与双线图框做成 CSS data-URI SVG 的 ambient 层
  （.slide::before/::after，z-index:-1），每页自动衬底、不进骨架，生成侧零成本；
  浅底（#FAFAFA letterbox 同色）配海军蓝/钢灰文字，对比度 ≥5:1。
- code 版式做成「图纸风代码卡」：ew-term 标注框 + 标题栏三圆点（ew-dot 方角圆点，
  遵守无装饰禁令）+ 文件名 + 语言标签（ew-lang），代码块 ew-codeblock 白纸浅蓝底，
  语法着色只用蓝灰两族（cm 灰 / kw 深蓝 / st 亮蓝 / fn 海军蓝）。
- 全站等宽字体栈 'JetBrains Mono'（fonts.css 自托管，CJK 由同族中日韩分段覆盖）；
  不新增 webfont 文件。
- demo 叙事为「智能温室工程方案」，传感器点位、管线长度、施工节点均带口径来源行，
  code 页放一段夜间巡检控制逻辑伪代码。
