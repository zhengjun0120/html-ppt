# 工作周报 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 XXX> 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-weekly-report` 作用域前缀生效，骨架里已写全，照抄结构即可。

---

## cover（周报封面）

用途：周报开头：周期徽章 + 大标题 + 汇报人
适用 role：cover。
内容约束：标题 ≤14 字；kicker 是周期（如“2026 年第 38 周”）

```html
<section class="slide" data-layout="cover">
  <div class="cover-head"><div class="logo">{{团队/项目名}}</div><div class="week-chip">{{周期}}</div></div>
  <p class="kicker">{{汇报主题}}</p>
  <h1 class="h1 mt-s">{{本周标题}}</h1>
  <p class="lede mt-m">{{一句话总览}}</p>
  <div class="deck-footer"><span class="meta">{{汇报人 · 日期}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cover-head, logo, week-chip, kicker, h1, mt-s, lede, mt-m, deck-footer, meta, notes

---

## kpi-grid（KPI 面板）

用途：本周 4-8 个关键指标（good/warn/bad 状态色）
适用 role：data。
内容约束：4-8 卡；每卡 = 指标名 + 数字 + 环比；状态色如实标注

```html
<section class="slide" data-layout="kpi-grid">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{指标面板标题}}</h2>
  <div class="grid g4 mt-l">
    <div class="kpi good"><span class="value">{{数字}}</span><span class="label">{{指标名}} <span class="delta">{{环比}}</span></span></div>
    <div class="kpi warn"><span class="value">{{数字}}</span><span class="label">{{指标名}} <span class="delta">{{环比}}</span></span></div>
    <div class="kpi bad"><span class="value">{{数字}}</span><span class="label">{{指标名}} <span class="delta">{{环比}}</span></span></div>
    <div class="kpi"><span class="value">{{数字}}</span><span class="label">{{指标名}} <span class="delta">{{环比}}</span></span></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, grid, g4, mt-l, kpi, good, warn, bad, value, label, delta, notes

---

## shipped-list（本周交付）

用途：本周完成的事项清单（ship-item 逐条）
适用 role：content。
内容约束：4-8 条；每条一句话 ≤30 字，可带状态 pill

```html
<section class="slide" data-layout="shipped-list">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{交付标题}}</h2>
  <div class="mt-l">
    <div class="ship-item">{{完成事项（可带 <span class="pill">状态</span>）}}</div>
    <div class="ship-item">{{完成事项}}</div>
    <div class="ship-item">{{完成事项}}</div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, mt-l, ship-item, pill, notes

---

## metrics-chart（数据图表）

用途：趋势柱状（chart-bars，高度/宽度按真实比例）
适用 role：data。
内容约束：3-6 根柱；柱高按真实比例给 style；每组 = 名称 + 数值

```html
<section class="slide" data-layout="metrics-chart">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{趋势标题}}</h2>
  <div class="chart mt-l">
    <div class="chart-bars">
      <div class="col"><div class="c" style="height:80%"><span>{{数值}}</span></div><span class="l">{{名称}}</span></div>
      <div class="col"><div class="c" style="height:80%"><span>{{数值}}</span></div><span class="l">{{名称}}</span></div>
      <div class="col"><div class="c" style="height:80%"><span>{{数值}}</span></div><span class="l">{{名称}}</span></div>
    </div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, chart, mt-l, chart-bars, col, b, lbl, notes

---

## blockers（阻塞与风险）

用途：卡住的事、需要的支持（blocker 逐条）
适用 role：content。
内容约束：1-4 条；每条说清“卡在哪 + 需要谁做什么”

```html
<section class="slide" data-layout="blockers">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{阻塞标题}}</h2>
  <div class="mt-l">
    <div class="blocker">{{阻塞描述：卡在哪 + 需要什么支持}}</div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, mt-l, blocker, notes

---

## next-week（下周计划）

用途：下周要做的事（next-row 逐条）
适用 role：content / cta。
内容约束：3-6 条；每条 ≤24 字、动词开头

```html
<section class="slide" data-layout="next-week">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{下周标题}}</h2>
  <div class="mt-l">
    <div class="next-row">{{计划事项}}</div>
    <div class="next-row">{{计划事项}}</div>
    <div class="next-row">{{计划事项}}</div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, mt-l, next-row, notes

---

## thanks（收尾）

用途：周报收尾：一句话 + 联系方式
适用 role：thanks。
内容约束：大字 ≤10 字；补充 ≤16 字

```html
<section class="slide" data-layout="thanks">
  <p class="kicker">{{周期}}</p>
  <h2 class="h2">{{收尾一句话}}</h2>
  <p class="lede mt-m">{{补充}}</p>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, lede, mt-m, notes
