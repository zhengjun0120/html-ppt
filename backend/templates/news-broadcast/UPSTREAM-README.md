# news-broadcast（新闻播报）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/news-broadcast`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：深红 #7F1D1D 标题、暗红 #991B1B 正文，正红 #DC2626 与浓红 #B91C1C 是新闻的灵魂色；
  背景 145° 白 → #FEF2F2 → #FEE2E2 渐变，冷静中带着紧迫；卡为 94% 白底 + 红系边框。
- 意象：白色画面左侧一道红色竖条、LIVE/BREAKING 标签、信息条（ticker）、硬阴影卡片，
  信息层级像头条、二条、三条一样分明。
- 排版：大写粗黑无衬线、行距紧凑、信息密度高；底部出现新闻标签。
- 适合：突发播报、数据发布、产品公告、季度汇报——「这件事很重要、必须关注」的时刻。
- 禁忌：柔和颜色、圆角设计、暗色背景、轻松娱乐化排版；红竖条不许消失或变细。

## 转换说明

- 原上游只有一个 16:9 单页预览（红竖条 + LIVE 徽标 + 台标 + 数据条 + 底部滚动条）；本模板按
  本仓库 deck-v2 契约重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 左侧红竖条（18px，带深红内缘）与顶部红条做成 .slide::before/::after 的 ambient 层
  （z-index:-1），每页自动衬底、不进骨架；LIVE 徽标未做成常驻元素，由台标框 + 时间戳承担身份。
- 滚动条（nw-ticker）、台标框（nw-station）、时间戳（nw-time）、头条横幅（nw-banner）转成骨架类，
  滚动条 absolute 定位在页脚上方，一页至多一条，只用于 cover/contents/divider/closing。
- 阴影全部为硬阴影（offset 无模糊、正红或深红低透明），卡片与块号一律直角，呼应上游「禁圆角」。
- 字体栈：Inter + Noto Sans SC（800 字重标题），时刻与编号用 JetBrains Mono；不新增 webfont 文件。
- letterbox 底色 --bg 取渐变中段 #FEF2F2，翻页时页外区域不跳色。
