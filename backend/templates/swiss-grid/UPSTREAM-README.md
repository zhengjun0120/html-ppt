# swiss-grid（瑞士网格）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/swiss-grid`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：标题深墨 #0f172a、正文石板灰 #334155，强调色近黑 #111827 与正红 #ef4444；
  背景 linear-gradient(160deg, #ffffff, #f8fafc 65%, #f1f5f9) 几乎纯白带一丝蓝灰；
  卡片 rgba(255,255,255,0.96)，边框 rgba(15,23,42,0.14) 细如发丝。
- 排版：Helvetica 或无衬线是唯一选择，一切左对齐；字号层级严格（标题 48px /
  副标题 24px / 正文 18px 量级，本模板在 1920 画布上按比例放大）。
- 布局：12 栏网格神圣不可侵犯，内容严格沿网格线排列；偶尔一根红色横线贯穿全页
  作为视觉锚点——红色是唯一的情绪出口。
- 场景：设计行业分享、品牌规范发布、严肃的排版展示、建筑设计叙事。
- 禁忌：自由布局、装饰性元素（渐变、纹理、图标）、超过两种字体、红色超过一次（本模板
  放宽为「一页红色元素 ≤3 处」以适配多版式结构，气质不变）。

## 转换说明

- 原上游预览是一个「外露网格线」的 12 栏 × 8 行构图（红块 + 近黑块 + 红色通栏条）；
  本模板按本仓库 deck-v2 契约重建为 9 版式骨架
  （cover/contents/keynotes/split/metrics/quote/divider/moments/closing），无 code 版式。
- 外露网格线做成 .slide::before ambient 层（repeating-linear-gradient 画 12 栏 × 8 行
  浅灰线 + 蓝灰洗底，z-index:-1），每页自动衬底、不进骨架——这是与
  swiss-international（极细环境网格、细体、克制）最直接的可感知区分。
- 红色通栏条 sg-bar 与红块 sg-square 取自上游预览的首尾母题，在 cover/closing 呼应；
  红色纪律收紧为「一页 ≤3 处」，写进 rules.md。
- 字体栈：Inter / Helvetica Neue / Noto Sans SC（900 大字），元数据用 JetBrains Mono；
  不新增 webfont 文件。
