# pitch-deck-vc（融资路演）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/pitch-deck-vc`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：标题深靛 #1e1b4b、正文靛蓝 #4338ca，强调色紫罗兰 #7c3aed 与蓝紫 #6366f1 是
  渐变的灵魂；背景 linear-gradient(145deg, #ffffff, #f5f3ff 50%, #ede9fe) 从纯白到淡紫；
  卡片 rgba(255,255,255,0.92)，边框 rgba(139,92,246,0.15) 极淡紫描边。
- 排版：现代无衬线、大字号、标题直接有力；数字要大、要粗、要一眼看到；正文精练，
  每页不超过三行文字。
- 布局：大留白是核心，每页只传达一个关键信息；渐变只用在标题或分隔线附近，不喧宾夺主。
- 场景：融资路演、种子轮 Pitch、VC Meeting、创业大赛——让投资人「看到数字就想投」。
- 禁忌：密集内容、小字、超过三种颜色、暗色背景；保持白底的通透。

## 转换说明

- 原上游只有一个 16:9 单页预览（左侧渐变面板 + 右侧 traction 指标与柱阵的封面构图）；
  本模板按本仓库 deck-v2 契约重建为 9 版式骨架
  （cover/contents/keynotes/split/metrics/quote/divider/moments/closing），无 code 版式。
- 封面侧板取自上游预览的左栏构图：靛 #1E1B4B → 靛蓝 #4338CA → 紫罗兰 #7C3AED 渐变，
  右上细环母题同源；等宽小节标签（TRACTION · …）沿用上游预览的 mono 口音。
- 与本仓库已有内置模板 pitch-deck（pd-，纯蓝 #3b5bff、渐变字大数字、圆角大卡）刻意区分：
  本模板身份色是紫罗兰系，封面有渐变侧板，标签与口径用等宽字，画布带淡紫渐变洗底。
- ambient 层：白→淡紫的晨曦渐变 + 右上极淡紫罗兰圆环，画在 .slide::before/::after
  （z-index:-1），每页自动衬底、不进骨架。
- 字体栈：Inter（900 标题与数字）+ Noto Sans SC 正文，元数据用 JetBrains Mono；
  不新增 webfont 文件。
