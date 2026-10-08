# 空白模板 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板的 `bl-*` 作用域类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 这是空白起点：九个最简版式，结构都用 base 原语 + 少量 `bl-*` 类。要扩展或改造版式，
> 编辑本文件登记条目，并让 index.html 的 demo 页与 style.css 的规则同步跟上。

---

## blank-cover（空白封面）
指纹：hero

用途：开场页。小引导语 + 大标题 + 一句话定位，垂直居中。
适用 role：cover。
内容约束：标题 ≤16 字；定位一句话 20-40 字。

合法类名：slide, full, kicker, h1, lede, mt-m

```html
<section class="slide full" data-layout="blank-cover">
  <p class="kicker">{{kicker · 这份演示是什么}}</p>
  <h1 class="h1 mt-m">{{主标题}}</h1>
  <p class="lede mt-m">{{一句话定位}}</p>
</section>
```

## blank-toc（空白目录）
指纹：stack

用途：目录/议程页。小节标题 + 编号条目行，每行标题带一句看点。
适用 role：toc。
内容约束：2-5 条；条目标题 4-12 字、看点 12-24 字。

合法类名：slide, kicker, h2, mt-s, bl-toc, mt-l, bl-toc-row, mono, bl-toc-n, bl-toc-t, dim, bl-toc-d

```html
<section class="slide" data-layout="blank-toc">
  <p class="kicker">{{kicker · 目录}}</p>
  <h2 class="h2 mt-s">{{议程标题，≤12 字}}</h2>
  <div class="bl-toc mt-l">
    <div class="bl-toc-row"><span class="mono bl-toc-n">{{01}}</span><span class="bl-toc-t">{{条目标题}}</span><span class="dim bl-toc-d">{{一句看点}}</span></div>
    <div class="bl-toc-row"><span class="mono bl-toc-n">{{02}}</span><span class="bl-toc-t">{{条目标题}}</span><span class="dim bl-toc-d">{{一句看点}}</span></div>
    <div class="bl-toc-row"><span class="mono bl-toc-n">{{03}}</span><span class="bl-toc-t">{{条目标题}}</span><span class="dim bl-toc-d">{{一句看点}}</span></div>
  </div>
</section>
```

## blank-divider（空白章节页）
指纹：hero

用途：章节过渡。大标题 + 一句话引导 + 右下角章节数字，垂直居中。
适用 role：divider。
内容约束：章节标题 ≤12 字；引导语 20-40 字；数字 1-2 字符。

合法类名：slide, full, kicker, h1, mt-m, lede, bl-ghost

```html
<section class="slide full" data-layout="blank-divider">
  <p class="kicker">{{kicker · part 1}}</p>
  <h1 class="h1 mt-m">{{章节标题}}</h1>
  <p class="lede mt-m">{{这一部分回答的问题}}</p>
  <div class="bl-ghost">{{1}}</div>
</section>
```

## blank-content（空白内容页）
指纹：stack

用途：内容页。小节标题 + 2-4 张卡片，自上而下。
适用 role：content。
内容约束：2-4 卡；卡标题 ≤10 字、正文 20-60 字。

合法类名：slide, h2, grid, g2, mt-l, card, h4, dim, mt-s

```html
<section class="slide" data-layout="blank-content">
  <h2 class="h2">{{小节标题}}</h2>
  <div class="grid g2 mt-l">
    <div class="card"><h4 class="h4">{{卡标题}}</h4><p class="dim mt-s">{{卡正文}}</p></div>
    <div class="card"><h4 class="h4">{{卡标题}}</h4><p class="dim mt-s">{{卡正文}}</p></div>
  </div>
</section>
```

## blank-list（空白列表页）
指纹：stack

用途：内容页（编号列表）。小节标题 + 编号要点纵列，每条一句论点加一段展开。
适用 role：content。
内容约束：2-5 条；要点 12-30 字、展开 20-45 字。

合法类名：slide, kicker, h2, mt-s, stack, mt-l, row, bl-list-row, mono, bl-list-n, bl-list-t, dim

```html
<section class="slide" data-layout="blank-list">
  <p class="kicker">{{kicker · 要点}}</p>
  <h2 class="h2 mt-s">{{这组要点回答什么}}</h2>
  <div class="stack mt-l">
    <div class="row bl-list-row"><span class="mono bl-list-n">{{01}}</span><div><p class="bl-list-t">{{要点一句话}}</p><p class="dim mt-s">{{展开说明}}</p></div></div>
    <div class="row bl-list-row"><span class="mono bl-list-n">{{02}}</span><div><p class="bl-list-t">{{要点一句话}}</p><p class="dim mt-s">{{展开说明}}</p></div></div>
    <div class="row bl-list-row"><span class="mono bl-list-n">{{03}}</span><div><p class="bl-list-t">{{要点一句话}}</p><p class="dim mt-s">{{展开说明}}</p></div></div>
  </div>
</section>
```

## blank-split（空白对照页）
指纹：split

用途：内容页（左右对照）。对比主题 + 左右两张对照卡，各带说明与适用前提。
适用 role：content。
内容约束：两侧各一份；方案名 ≤10 字、说明 25-60 字、前提 15-40 字。

合法类名：slide, kicker, h2, mt-s, grid, g2, mt-l, card, bl-panel, h4, dim, bl-panel-note

```html
<section class="slide" data-layout="blank-split">
  <p class="kicker">{{kicker · 对照}}</p>
  <h2 class="h2 mt-s">{{对比主题}}</h2>
  <div class="grid g2 mt-l">
    <div class="card bl-panel"><h4 class="h4">{{左侧方案名}}</h4><p class="dim mt-s">{{左侧说明}}</p><p class="bl-panel-note">{{左侧适用前提}}</p></div>
    <div class="card bl-panel"><h4 class="h4">{{右侧方案名}}</h4><p class="dim mt-s">{{右侧说明}}</p><p class="bl-panel-note">{{右侧适用前提}}</p></div>
  </div>
</section>
```

## blank-data（空白数据页）
指纹：chart

用途：数据/指标页。小节标题 + 2-4 张大数卡 + 来源行。
适用 role：data、content。
内容约束：2-4 卡；数值 ≤6 字符、指标名 ≤12 字、口径 12-30 字；底部注来源。

合法类名：slide, kicker, h2, mt-s, grid, g3, mt-l, bl-stat, bl-stat-v, bl-stat-l, dim, bl-stat-note, bl-source, mt-m

```html
<section class="slide" data-layout="blank-data">
  <p class="kicker">{{kicker · 关键数字}}</p>
  <h2 class="h2 mt-s">{{这些数字回答什么}}</h2>
  <div class="grid g3 mt-l">
    <div class="bl-stat"><div class="bl-stat-v">{{数值}}</div><div class="bl-stat-l">{{指标名}}</div><p class="dim bl-stat-note">{{口径或对比}}</p></div>
    <div class="bl-stat"><div class="bl-stat-v">{{数值}}</div><div class="bl-stat-l">{{指标名}}</div><p class="dim bl-stat-note">{{口径或对比}}</p></div>
    <div class="bl-stat"><div class="bl-stat-v">{{数值}}</div><div class="bl-stat-l">{{指标名}}</div><p class="dim bl-stat-note">{{口径或对比}}</p></div>
  </div>
  <p class="dim bl-source mt-m">{{来源：数据出处}}</p>
</section>
```

## blank-quote（空白金句页）
指纹：quote

用途：引用/金句页。大字引文 + 出处 + 一句注释，垂直居中。
适用 role：quote、content。
内容约束：引文 15-40 字；出处 ≤24 字；注释 20-45 字。

合法类名：slide, full, kicker, bl-quote, mt-l, mono, bl-quote-src, mt-m, dim

```html
<section class="slide full" data-layout="blank-quote">
  <p class="kicker">{{kicker · 引用}}</p>
  <p class="bl-quote mt-l">{{引文}}</p>
  <p class="mono bl-quote-src mt-m">—— {{出处}}</p>
  <p class="dim mt-m">{{为什么记下它}}</p>
</section>
```

## blank-thanks（空白收尾页）
指纹：hero

用途：收尾/行动号召。大标题 + 一句收束 + 联系方式或下一步动作，垂直居中。
适用 role：thanks、cta。
内容约束：标题 ≤14 字；收束 20-45 字；联系行 ≤30 字。

合法类名：slide, full, kicker, h1, mt-m, lede, mono, mt-l, bl-contact

```html
<section class="slide full" data-layout="blank-thanks">
  <p class="kicker">{{kicker · 收尾}}</p>
  <h1 class="h1 mt-m">{{收束标题}}</h1>
  <p class="lede mt-m">{{带走的几件事或一句致谢}}</p>
  <p class="mono bl-contact mt-l">{{联系方式或下一步}}</p>
</section>
```
