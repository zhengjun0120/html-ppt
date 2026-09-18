# presenter-cards 入库适配记录

来源：html-ppt-skill/templates/full-decks/presenter-mode-reveal（MIT，见 UPSTREAM-README.md）
适配日期：2026-09-18（scaffold + materialize 流水线）

## 改动清单

1. **共享资产与挂载标记**（scaffold 完成）：`/assets/deck-v2/*` 绝对路径；`<!-- SLIDES:START/END -->` 就位——并修正了 scaffold 里 SLIDES:END 的位置（原先落在 `.deck` 闭合之后的固定提示条内，实例化锚点会错位）。
2. **主题依赖落地**：源 demo 依赖 `assets/themes/tokyo-night.css`（deck-v2 世界无 themes 目录，链接 404）。改为在 style.css 顶部内置同族 GitHub 深色 token（--bg #0d1117 / --accent #58a6ff / --grad 蓝紫橙，与源文件自带的 fallback 值一致），demo 头部删掉死链 `<link id="theme-link">` 与 `data-themes` 属性。
3. **demo 去除 JS 交互残留**：删掉 `.deck` 外的 fixed 键盘提示条（11px、不随画布缩放，S/T/O/R 是 runtime 演讲者视图功能，不属模板契约）。
4. **demo 每页补 `data-layout`**（6 版式：cover / agenda / cards-3 / feature-grid / rule-list / demo-close）；`<aside class="notes">` 统一为 `<div class="notes">`（与契约骨架一致，runtime 同样识别）。
5. **style.css 适配**（源视觉逐条保留）：
   - 字号地板：kicker 13→15px；agenda-row t 17→18px、num 14→16px、d 12→14px；card dim 14→18px；feature-row b 17→18px、dim 14→18px；rule-row dim 15→18px；speaker span 13→14px；code-block 15→17px。
   - 补 `.slide.full` 居中规则（防御性，与 base 一致；本模板骨架未用 full）。
   - 末尾追加变体槽 `v-ember` / `v-lilac`：只覆盖 `--accent/--accent-2/--grad`，结构/组件全部继承。
6. **契约文件**：template.json（3 变体、6 版式、5 种指纹：hero / table / cards / stack / code）、layouts.md、rules.md、本文件。

## 取舍

- 源模板的磁吸卡片式演讲者视图（S 键弹窗、拖拽卡片）集成在上游 runtime 里，deck-v2 共享 runtime 未包含；本模板保留其内容侧契约——notes 逐字稿（runtime 演讲者/备注视图照常读取）。
- 5 主题切换（T 键）改为 3 个 template.json 变体（夜蓝 / 暖橙 / 丁香），换变体 = body 挂 class，agent 不参与。
- demo 页 6 的"命令演示 + 收尾"合并页保留原样，登记为单一 demo-close 版式（roles 覆盖 code/cta/thanks）。
