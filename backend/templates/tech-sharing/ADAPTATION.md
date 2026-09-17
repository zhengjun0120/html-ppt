# tech-sharing 入库适配记录

来源：html-ppt-skill/templates/full-decks/tech-sharing（MIT，见 UPSTREAM-README.md）

适配日期：2026-09-17

## 改动清单

1. **资产路径改绝对路径**（一次性，实例化时不再改写）：
   - `../../../assets/fonts.css` → `/assets/deck-v2/fonts.css`
   - `../../../assets/base.css` → `/assets/deck-v2/base.css`
   - `../../../assets/animations/animations.css` → `/assets/deck-v2/animations.css`
   - `../../../assets/runtime.js` → `/assets/deck-v2/runtime.js`
2. **挂载标记**：demo slides 区域外包 `<!-- SLIDES:START -->` / `<!-- SLIDES:END -->`，实例化时区间内内容被剥离、由生成管线填入正式页面。
3. **主题变体槽**：style.css 末尾追加 `v-blue` / `v-ember` 两个 variant class（只覆盖 `--accent/--accent-2/--grad`），template.json `variants[]` 登记。
4. **新增文件**：template.json、layouts.md（14 版式登记）、rules.md、preview/ 目录、本文件。
5. **未改动**：demo deck 的 8 页内容原样保留（画廊 live 预览用）；style.css 其余规则逐字未动。
