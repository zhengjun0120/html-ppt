# dir-nav-minimal 入库适配记录

来源：html-ppt-skill/templates/full-decks/dir-key-nav-minimal（MIT，见 UPSTREAM-README.md）
适配日期：2026-09-18（scaffold + materialize 流水线）

## 改动清单

1. **共享资产与挂载标记**（scaffold 完成）：`/assets/deck-v2/*` 绝对路径；`<!-- SLIDES:START/END -->` 就位。
2. **demo 去除 JS 交互残留**：删 `is-active` 类（runtime.js 自管）、`dk-keyhint` 里的 `<kbd>` 键帽（保留纯文本外观）、`margin:auto 0` 居中包裹层（base `.slide` 已垂直居中）。
3. **demo 每页补 `data-layout`**（8 版式：cover / divider / list / compare / code / stat / cta / thanks）；第 6 页 SVG 进度条从写死色改为 `class="dk-accent"` + `fill="currentColor"`（变体可染色）。
4. **style.css 适配**（源视觉逐条保留）：
   - 字号地板：装饰小字 11-12px → 14px（dk-snum / dk-eyebrow / dk-keyhint / kbd / dk-page）；代码 dk-code 16px → 17px。
   - 补 `.slide.full` 居中规则（防御性，与 base 一致；本模板骨架未用 full）。
   - 末尾追加变体槽 `v-mono` / `v-amber`：只覆盖 `.dk-accent` 全局 accent 色板，八张底色轮换不动（v-amber 在浅色页自动换深琥珀保对比度）。
5. **契约文件**：template.json（3 变体、8 版式、5 种指纹：hero / stack / cards / code / chart）、layouts.md、rules.md、本文件。

## 取舍

- 源 demo 的键盘翻页（←→/space/F）由共享 runtime.js 提供，模板契约不登记任何交互类；keyhint 只剩装饰外观。
- 页面底色 t-* 与版式一一绑定（cover=indigo … thanks=charcoal）：保住「一页一色」的轮换节奏，代价是同版式重复页会重复底色——rules.md 的「禁同版式连排」由此而来。
- 变体只动全局 accent，不动页面底色：色板轮换是身份，accent 是变量。
