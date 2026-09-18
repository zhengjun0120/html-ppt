# obsidian-gradient 入库适配记录

来源：html-ppt-skill/templates/full-decks/obsidian-claude-gradient（MIT，见 UPSTREAM-README.md）

适配日期：2026-09-18

## 改动清单

1. **四件套契约**：template.json（8 版式：cover / section / compare / steps / config / metrics / quote-cta / thanks，指纹 hero×3 + cards + table + code + chart + quote，6 种指纹）、layouts.md（骨架自 demo 各页提炼，{{占位符}} 带字数区间）、rules.md、本文件。
2. **约束区间化**：卡内要点 10-22 字/条、hl 洞察 25-45 字、lede 30-55 字等全部区间化；结构性硬顶（步数 ≤5、代码 ≤18 行、胶囊字符数）保留。
3. **字号 floors**：oc-snum 12→14、oc-tag 11→14、oc-badge 11→14、oc-code 14→17、oc-sc p 14→18、oc-sc h4 17→20、oc-quote .attr 13→15、oc-hl 16→18、oc-pill 14→16；新增 `.oc-card h4`（20px）/`.oc-card p`（18px）类级规则，骨架不再依赖行内字号。改后 CSS 无 <14px 字号。
4. **变体槽**：accent 系底色参数化为 `--oc-accent-rgb` / `--oc-accent2-rgb` 三元组，渐变字三停参数化为 `--oc-g1/g2/g3`（默认取值 = 原字面值，默认观感不变；`.oc-cbg` 紫晕随之联动）；追加 `v-emerald`（翠青）、`v-sunset`（落日）两个变体，accent3 文字在 `#0d1117` 上均 >8:1。
5. **demo 适配**：`<html lang="zh-CN">`（scaffold 已就位）；8 个 section 逐页补 `data-layout`；内容与结构未动。

## 已知取舍

- 目录名/模板 id 是 `obsidian-gradient`，CSS 作用域保持源 body class `tpl-obsidian-claude-gradient`（与 data-dark 先例一致；实例化只往 body 追加变体 class）。
- demo 第 3 页对比卡的第二张用了行内紫描边（字面色）与 14px 行内正文字号：按"demo 内容不动"保留；生成页走骨架（var() 描边 + 18px 类级字号），变体下描边/底色随主题换色。
- demo 第 6/7 页的 🧠⚡🔗⬇ emoji 为源视觉的一部分，demo 保留；rules.md 已对新生成页面禁 emoji。
- 本模板 `.slide` 自带 flex 居中（非 grid），无网格陷阱；封面/章节/收尾天然垂直居中，无需 `full` 类。
