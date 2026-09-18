# hermes-terminal 入库适配记录

来源：html-ppt-skill/templates/full-decks/hermes-cyber-terminal（MIT，见 UPSTREAM-README.md）

适配日期：2026-09-18

## 改动清单

1. **四件套契约**：template.json（8 版式：cover / divider / spec-cards / trace / compare-chart / verdict / install / thanks，指纹 hero×3 + cards + code×2 + chart×2，4 种指纹）、layouts.md（骨架自 demo 各页提炼，{{占位符}} 带字数区间）、rules.md、本文件。
2. **约束区间化**：卡片说明 12-24 / 20-40 字、lede 16-55 字等，全部写成区间；结构性上限（标题字符数、代码行数、柱数）保留硬顶。
3. **字号 floors**（观众可读性）：hc-chrome 11→14、hc-card .lbl 10→14、hc-card .desc 13→18、hc-codebox 14→17、hc-tag 11→14（padding 3/10→4/12）、hc-footer 10→14；改后除装饰大字外无 <14px 字号。
4. **变体槽**：glow/网格/扫描线/边框/代码盒内发光的磷光色参数化为 `--hc-glow`（RGB 三元组），图表系列色参数化为 `--hc-s1/--hc-s2`（默认取值 = 原字面色，默认观感不变）；追加 `v-amber`（琥珀屏）、`v-ice`（冰蓝）两个变体，只换发光色相，对比度均过 AA（#e9c58a / #79c0ff on #0a0c10 均 >8:1）。
5. **demo 适配**：`<html lang="zh-CN">`（scaffold 已就位）；8 个 section 逐页补 `data-layout`；内容与视觉未动。
6. **layouts.md 图表骨架**：compare-chart 的 SVG 用 `style="fill:var(--hc-s1)…"` 写法，变体下图表随主题换色（demo 内 SVG 保留源字面色，未动）。

## 已知取舍

- 目录名/模板 id 是 `hermes-terminal`，CSS 作用域保持源 body class `tpl-hermes-cyber-terminal`（与 data-dark 的 id≠作用域先例一致；实例化只往 body 追加变体 class，原 class 保留）。
- 模板无 sidebar 包裹层、也无裸 `.slide` grid，故无 `:has(.sidebar)` 守卫需求；封面/章节/收尾的垂直居中由 base.css `.slide` 的 flex column + 源 `margin:auto 0` 包裹层共同完成，骨架原样保留。
- demo 第 5 页 SVG 柱状图的轴文字在源里是 13px（viewBox ≈1:1 渲染）：demo 按"内容不动"保留；layouts.md 骨架已把 SVG 字号提到 16，生成页用骨架。
- amber 变体下，绿色 `hc-tag`/结论字随 `--hc-green` 变琥珀，红/琥珀分级标签的语义对比略减弱（单色终端的固有观感），rules.md 已要求数据页用图表而非颜色下结论。
