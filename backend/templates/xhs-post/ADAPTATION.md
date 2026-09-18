# xhs-post 入库适配记录

来源：html-ppt-skill/templates/full-decks/xhs-post（MIT，见 UPSTREAM-README.md）。
小红书 3:4 竖版九宫格图文（810×1080）：手帐便签 + 贴纸 + 圆角硬阴影，封面→hook→痛点→aha→步骤→效果→CTA。

适配日期：2026-09-18

## 改动清单

1. **画布**：template.json `canvas {"w":810,"h":1080}`（系统量测/导出按自定义画布走）；
   demo 的 `.deck` 保留 scaffold 的 `data-w="810" data-h="1080"`，style.css 保留
   `--deck-w/--deck-h` 与 810×1080 的 `.slide` 设计，**源竖版比例未动、未改成 1920**。
2. **新增文件**：template.json（7 版式登记 / 3 变体）、layouts.md、rules.md、本文件。
3. **版式规划**：demo 9 页提炼为 7 个登记版式，指纹 4 种——
   cover(hero)、hook(hero)、pain(stack)、truth(quote)、step(stack)、result(cards)、cta(hero)。
   三个步骤页共用 `step` 版式；`result` 由纵向 stack 改为 2×2 便签格（base 原语 grid g2）以补足
   cards 指纹（判重需要 ≥4 种指纹），便签视觉语言不变。
4. **style.css 适配**（保持源视觉）：
   - 字号地板：`.page-dot` 14→15px（其余本就达标：sticker 18、hand-box 说明 demo 内联 16→18、ht 16、bottom-bar 15）。
   - 主题变体槽：追加 `v-mint` / `v-lav`（只覆盖 --bg 族 + accent 三色 + --grad），--accent-ink 深墨不变。
   - 无裸 `.slide` grid 陷阱；不登记 `.slide.full`（无侧栏、无需逃逸）。
5. **demo 适配**：9 个 section 全部加 `data-layout`（step 页 3 连）；痛点卡说明内联 16→18px；
   举例便签底色 `#fff5ef` → `var(--surface-2)`（默认观感不变、变体下不再穿帮）；
   result 页改写为登记的 2×2 格（每格标题 + 一句对比，末格 accent-3 强调）。

## 取舍

- emoji 是本模板的版面语言（big-emoji / 便签词首 / sticker / ❌💡✅ 引子），不做 tech-sharing 式全禁；
  rules.md 给出每页 ≤6 个、只在词首/装饰位的边界。
- `step` 版式骨架的步骤号底色默认 `var(--accent)`，accent-2/accent-3 轮换写进 constraints（demo 6/7 页做法）。
- fonts[] 登记 Inter / Noto Sans SC / JetBrains Mono（页码贴纸与账号条用等宽）。
- 总页数 6-9 页、cover→hook→pain→truth→step×N→result→cta 顺序完整，写入 rules.md 节奏节。
