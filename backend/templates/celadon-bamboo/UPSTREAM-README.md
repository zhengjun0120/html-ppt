# celadon-bamboo（青瓷竹影）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/celadon-bamboo`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：浅黄绿纸面 #F3F1D5 / #EEF0D6 底，墨字 #1B1A18（或 #263238）；主色青瓷绿三阶
  #9DCB95 / #6FAF82 / #5DA8A7，远景可过渡至 #2F7D91；强调色仅少量使用朱砂 #B53A2D
  （或 #C64632 / #A87B45）。
- 意象：竹简竖柱、竹叶、刻纹、题跋、朱砂印章、纸伞、薄雾、浅水倒影；边缘柔和，阴影轻长淡。
- 排版：标题用书卷气高对比衬线体；正文与注释清秀细宋体、字重偏轻、行距拉开；
  可少量用竖排短句、小型落款与细线题签，正文保持横向易读。
- 布局：纵向秩序与空间纵深，大片留白承托文字；半透明青绿竖色块、竹简长条、低对比描边分区。
- 禁忌：荧光色、强投影、厚重 3D、金属质感、粗黑边框、强对比商业按钮、塞满数据表格。

## 转换说明

- 原上游只有一个 16:9 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 右侧竹简竖条用 CSS border 实现（粗青瓷深绿左缘 + 细青瓷浅绿右缘 + repeating-gradient
  竹节横纹，整体低透明度），左下竹叶与薄雾做成 data-URI SVG；两层都挂在
  .slide::before/::after（z-index:-1）上，每页自动衬底、不进骨架，生成侧零成本。
- 关键数据色从朱砂改为远景深青 #2F7D91：上游明确「强调色仅少量使用朱砂」，
  把朱砂留给印章与题签点缀（qc-seal / qc-pill-accent / qc-verse-accent），数字用深青更耐看。
- 预览里的竖排诗句「风来竹有声」保留为 quote 版式的 qc-verse 示例；
  题跋小字落款落实为 qc-src（小号、宽字距、无衬线）。
- 字体栈：标题 Georgia / 宋体衬线高对比（沿用上游 Georgia 取向），小字 Noto Sans SC；
  不新增 webfont 文件。
