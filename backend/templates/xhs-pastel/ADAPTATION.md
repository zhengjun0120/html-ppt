# xhs-pastel 入库适配记录

来源：html-ppt-skill/templates/full-decks/xhs-pastel-card（MIT，见 UPSTREAM-README.md）。
共性提取自 20260412-obsidian-skills 软色卡系统 + 20260409 v2-白底版胶囊 chip 顶条。

适配日期：2026-09-18

## 改动清单

1. **资产路径**（scaffold 已就位）：`/assets/deck-v2/fonts.css`、`base.css`、`runtime.js` 绝对引用，挂载标记 `<!-- SLIDES:START/END -->`。
2. **新增文件**：template.json（9 版式登记 / 3 变体 / canvas 1920×1080）、layouts.md、rules.md、本文件。
3. **版式规划**：demo 8 页提炼为 9 个登记版式，指纹 6 种——
   cover(them hero)、chapter(hero)、grid-2x2(cards)、checklist(stack)、quote-card(quote)、
   prompt(code)、donut-stat(chart)、steps-3(cards)、thanks(hero)。
   其中 checklist 从 donut 页图例列提炼、登记但 demo 未单独占页（8 页 demo 各页 data-layout 均已登记）。
4. **style.css 适配**（保持源视觉，只动地板与槽位）：
   - 字号地板提级：`.xp-card p` 15→18px（观众必读说明）；装饰小字 chip/page/kicker/footer 13/13/14/12→15px；donut SVG 中心副标 12→15px；码箱 14→17px（行高 1.85→1.75 保持 16 行容量）。
   - [deck-v2 patch] 新增 `xp-chart-row` / `xp-legend` / `xp-dot` 三个辅助类：上游 demo 里 donut 图与图例是内联 style，登记骨架照抄内联太脆，提成类视觉不变。
   - 主题变体槽：追加 `v-minty` / `v-lilac`（只覆盖 --xp-bg / peach / rose 两族 accent token），template.json variants[] 登记。
5. **demo 适配**：8 个 section 全部加 `data-layout`；donut 页改用 `xp-chart-row/xp-legend/xp-dot` 与登记骨架一致，图例文字 17→20px、SVG 340px。
6. **未改动**：blob / chip / 衬线斜体大字 / 马卡龙整色卡 / 码箱等源设计逐字保留；demo 文案原样。

## 取舍

- 不登记 `.slide.full`：本模板无侧栏、页页居中，无需 full 逃逸。
- emoji 只保留在 steps-3 的 `xp-num` 装饰位（☕🌸🌙），边界写入 rules.md。
- fonts[] 只登记 fonts.css 实际提供的 Inter / Noto Sans SC / JetBrains Mono；
  Playfair Display / Noto Serif SC 是本地衬线回退（与 course-module 同一处理），不登记。
