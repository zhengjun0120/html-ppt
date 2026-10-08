# 空白模板 · 质量规则（生成阶段注入）

本模板从空白开始：视觉与结构由使用者在 style.css / index.html / layouts.md 里定义。

- 版式共 9 个：blank-cover（封面）/ blank-toc（目录）/ blank-divider（章节）/ blank-content（卡片）/ blank-list（编号列表）/ blank-split（左右对照）/ blank-data（大数）/ blank-quote（金句）/ blank-thanks（收尾）。写页前先 `read_layout` 拿骨架，照结构填内容。
- 内容页有五种可选（content/list/split/data/quote 都声明了 content role）：长 deck 交替使用避免连续同版式；关键数字优先用 blank-data，不要把数字埋进句子。
- 字号底线：正文 ≥20px（base.css 的 .slide 已保证），标题用 .h1 / .h2 / .h4。
- 颜色只用 token（--accent / --text-* / --surface…），不要写死色值。
- 填充率门禁：内容页 ≥45%；封面/章节/金句/收尾按 hero/quote 豁免，但不要故意留白凑数。
- 扩展版式先登记 layouts.md（骨架 + 合法类名 + 指纹），并让 demo 页与 style.css 同步。
