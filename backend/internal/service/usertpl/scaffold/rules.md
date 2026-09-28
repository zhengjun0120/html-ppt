# 空白模板 · 质量规则（生成阶段注入）

本模板从空白开始：视觉与结构由使用者在 style.css / index.html / layouts.md 里定义。

- 版式只有 `blank-cover` / `blank-content` 两个；写页前先 `read_layout` 拿骨架，照结构填内容。
- 字号底线：正文 ≥20px（base.css 的 .slide 已保证），标题用 .h1 / .h2 / .h4。
- 颜色只用 token（--accent / --text-* / --surface…），不要写死色值。
- 填充率门禁：内容页 ≥45%；封面按 hero 豁免，但不要故意留白凑数。
- 扩展版式先登记 layouts.md（骨架 + 合法类名 + 指纹），并让 demo 页与 style.css 同步。
