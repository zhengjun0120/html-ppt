# safety-alert 入库适配记录

来源：html-ppt-skill/templates/full-decks/testing-safety-alert（MIT，见 UPSTREAM-README.md）

适配日期：2026-09-18

## 改动清单

1. **四件套契约**：template.json（8 版式：cover / divider / risk-levels / policy-code / incident-chart / checklist / tonight / thanks，指纹 hero×3 + cards×2 + code + chart + table，5 种指纹）、layouts.md（骨架自 demo 各页提炼，{{占位符}} 带字数区间）、rules.md、本文件。
2. **约束区间化**：分级卡说明 8-18 字/行、清单条目 12-30 字、告警框正文 20-45 字等全部区间化；结构性硬顶（标题行数/字符、代码 10-20 行、柱 3-6 根、清单 5-7 条）保留。
3. **字号 floors**：ts-alert-tag 13→14、ts-page 13→14、ts-card .lbl 12→14、ts-card p 14→18（行高 1.55→1.6）、ts-footer 12→14、ts-codebox 14→17；改后 CSS 无 <14px 字号。
4. **变体槽**：追加 `v-amber`（主警示色 #b45309）、`v-blue`（警蓝 #1d4ed8），只覆盖 `--ts-red/--ts-red-soft`——条纹、警示标签、strike、check 框、codebox 边全部随 token 联动；白字对比度 4.5:1 / 6.3:1 过 AA。
5. **demo 适配**：`<html lang="zh-CN">`（scaffold 已就位）；8 个 section 逐页补 `data-layout`；内容与视觉未动。
6. **layouts.md 图表骨架**：incident-chart 的 SVG 柱/图例用 `style="fill:var(--ts-green/amber/red)"` 写法，变体下图表随主题换色（demo 内 SVG 保留源字面色，未动）。

## 已知取舍

- 目录名/模板 id 是 `safety-alert`，CSS 作用域保持源 body class `tpl-testing-safety-alert`（与 data-dark 先例一致；实例化只往 body 追加变体 class）。
- demo 第 5 页 SVG 轴文字 13-14px 为源视觉，demo 保留；layouts.md 骨架已提到 16px（注解 15px），生成页用骨架。
- demo 第 4 页 bad 项为防复制执行做了实体混淆，rules.md 把这一做法固化成规则。
- 模板无 sidebar、无裸 `.slide` grid，无网格陷阱；封面/章节/收尾的垂直居中由 base.css `.slide` flex（divider/thanks 另有源 `margin:auto 0` 包裹层）完成，骨架原样保留，无需 `full` 类。
- **源色板对比度（保留源设计，未调色）**：white on `--ts-red #e0314a` = 4.46:1（14px 粗体标签，差 0.04）；`--ts-red` 文本 on bg = 4.31:1（kicker 15px 粗体；h1/h2 里的 red span 为 ≥54px 大字，按 AA-large 3:1 达标）；white on `--ts-amber #d97706` = 3.19:1、amber 文本 on bg = 3.07:1（bold ≥18px 场景，按 AA-large 边缘达标）。两变体的白字对比全部 ≥4.5（v-amber 5.02 / v-blue 6.70），v-amber/v-blue 的 soft 底上墨字 16.7:1+；其余源色对（ink 14.9:1、ink2 8.5:1、green 5.5:1）全过 AA。
