# red-gold-ceremony（初心红典）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/red-gold-ceremony`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：主背景深正红 #CD170A，层次 #BB1718 / #C82C2D；标题用 #FEEBA5 至 #E29B60 的克制金渐变，
  正文白 #FFFFFF，辅助信息 #FFEBA9。红为主、金为点睛，避免大面积反光。
- 排版：封面标题用有骨架感的宋体/高对比衬线，正文简洁无衬线但保持清晰行距；数字年份可用几何字体。
- 布局：居中对称与稳定轴线；细金线与序号负责分组，不把信息单元做成浮卡。
- 装饰：低密度五角星、金色圆点、细分隔线；禁卡通、手绘、多色复杂纹样、厚重金色光晕。
- 场景：主题教育、政府工作总结、企业年会、仪式致辞、文化纪念。

## 转换说明

- 上游只有一个 177.78vh 比例的单页预览（含内框 frame、◆◆◆/✦✦✦ 标记、金渐变标题）；
  本模板按本仓库 deck-v2 契约重建为 9 版式骨架
  （cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 金渐变只保留在 h1（香槟金 #FEEBA5 → 鎏金 #E29B60），h2 用单色金：投屏下渐变过多会显得「喜庆」，
  与「庄重不喜庆」的定位相悖。
- 深红改为径向光场（#D43525 → #CD170A → #8F140F），letterbox 用中段值 #B4160F，翻页不跳。
- 内框与页底金辉做成 CSS ambient 层（.slide::before/::after，z-index:-1），每页自动衬底、不进骨架。
- 字体栈：标题 Georgia/宋体，小字与元数据 JetBrains Mono + Noto Sans SC；不新增 webfont 文件。
