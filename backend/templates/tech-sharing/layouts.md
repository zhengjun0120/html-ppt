# tech-sharing · 版式登记簿（LLM 唯一契约）

> **版式锁**：你写的每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，
> 只准使用本文件骨架里出现过的类名（base.css 原语 + 本模板专属类）。
> 未登记的版式和未知类名会被服务端**直接拒收**。
> 每页开头先 `read_layout` 拿骨架，再替换 `{{占位符}}` 为真实内容。
> 不要写 `data-id`（系统自动编号）；不要写 `data-title`（系统取 h1/h2 文本）。

通用约定：
- 类名分两族：**base 原语**（`slide kicker h1 h2 h3 h4 lede dim dim2 mono stack row grid g2 g3 g4 center card card-accent card-soft card-outline pill pill-accent divider divider-accent mt-s mt-m mb-m mt-l mb-l fill nowrap`）与**模板专属**（`terminal bar dot kw fn str cmt num tag agenda-row speaker av t d ghost-num`）。
- 中文正文用默认字体，代码/命令/文件名/数字指标一律 `<span class="mono">` 或 `.tag`。
- 每页讲稿写在页面末尾：`<div class="notes">…</div>`（观众不可见，演讲者模式可见）。
- **数字纪律**：关键数字/对比/规模一律用 stat-hero、kpi-grid 或 two-column 呈现，不准埋进句子里。

---

## cover（封面）
指纹：hero

用途：开场页。大标题 + 一句话定位 + 讲者行。
适用 role：cover。
内容约束：主标题 8-14 字（可一行或两行 `<br>` 分行）；lede 20-30 字；kicker 是主题标签（如 `tech-sharing / 2026-09`）。

```html
<section class="slide" data-layout="cover">
  <p class="kicker">{{主题标签}}</p>
  <h1 class="h1">{{主标题，两行用 <br> 分隔}}</h1>
  <p class="lede mt-m">{{一句话定位，20-30 字}}</p>
  <div class="speaker"><div class="av"></div><div><b>{{讲者名}}</b><span>{{身份 · 时长}}</span></div></div>
  <div class="deck-footer"><span class="mono">{{话题标签 如 #golang #microservice}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h1, lede, mt-m, speaker, av, deck-footer, mono, slide-number, notes

---

## agenda（议程路线图）
指纹：table

用途：目录/议程。编号行列表，每行标题 + 一句看点 + 预计时长。
适用 role：toc。
内容约束：3-6 行；每行标题 8-16 字 + 看点 12-24 字；d 列写预计时长（如 `~5min`）。

```html
<section class="slide" data-layout="agenda">
  <p class="kicker">agenda.toml</p>
  <h2 class="h2">{{议程标题，如：今天的路线图}}</h2>
  <div class="stack mt-l">
    <div class="agenda-row"><span class="num">01</span><span class="t">{{条目标题}}</span><span class="d">{{一句看点，12-24 字}} · ~{{时长}}min</span></div>
    <div class="agenda-row"><span class="num">02</span><span class="t">{{条目标题}}</span><span class="d">{{一句看点，12-24 字}} · ~{{时长}}min</span></div>
    <div class="agenda-row"><span class="num">03</span><span class="t">{{条目标题}}</span><span class="d">{{一句看点，12-24 字}} · ~{{时长}}min</span></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, stack, mt-l, agenda-row, num, t, d, notes

---

## section-divider（章节幕）
指纹：hero

用途：章节过渡。巨大幽灵章节数字压住右下角，左上是本章进度与一个问题。
适用 role：divider。
内容约束：kicker 写章节进度（如 `part · 1 / 3`）；标题 6-16 字；lede 20-40 字；2-3 个本章看点 tag；ghost-num 写章节数字。

```html
<section class="slide" data-layout="section-divider">
  <p class="kicker">{{章节进度，如 part · 1 / 3}}</p>
  <h2 class="h2">{{章节标题，6-16 字}}</h2>
  <p class="lede mt-m">{{本章要解决的一句问题，20-40 字}}</p>
  <div class="row mt-l">
    <span class="tag">{{本章看点 1}}</span>
    <span class="tag">{{本章看点 2}}</span>
  </div>
  <div class="ghost-num">{{章节数字，如 1}}</div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, lede, mt-m, row, mt-l, tag, ghost-num, notes

---

## cards-3（三卡要点）
指纹：cards

用途：三个并列论点、三种方案、三类问题。这是本模板的"内容主力"版式。
适用 role：content。
内容约束：恰好 3 卡；卡标题 4-12 字；说明 25-45 字（一句论断 + 一句展开）；每卡一个结论 tag（4-8 字）。

```html
<section class="slide" data-layout="cards-3">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{标题：一句话点破这页在对比什么}}</h2>
  <div class="grid g3 mt-l">
    <div class="card card-accent"><h4>{{卡标题}}</h4><p class="dim">{{说明，25-45 字}}</p><span class="tag mt-s">{{结论标签}}</span></div>
    <div class="card card-accent"><h4>{{卡标题}}</h4><p class="dim">{{说明，25-45 字}}</p><span class="tag mt-s">{{结论标签}}</span></div>
    <div class="card card-accent"><h4>{{卡标题}}</h4><p class="dim">{{说明，25-45 字}}</p><span class="tag mt-s">{{结论标签}}</span></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, grid, g3, mt-l, card, card-accent, h4, dim, tag, mt-s, notes

---

## split-terminal（左文右终端）
指纹：split

用途：左边把概念讲透（两段），右边放代码/命令/输出佐证。
适用 role：content / code。
内容约束：左侧两段：lede 40-70 字 + dim 补充 20-40 字 + 2-4 个 tag；terminal 内代码 6-16 行、每行 ≤60 字符。

```html
<section class="slide" data-layout="split-terminal">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{标题}}</h2>
  <div class="grid g2 mt-l" style="align-items:start">
    <div>
      <p class="lede">{{概念讲解，40-70 字，代码名词用 <span class="mono">}}</p>
      <p class="dim mt-m" style="font-size:18px">{{补充：边界、代价或反例，20-40 字}}</p>
      <div class="mt-l">
        <span class="tag">{{关键术语1}}</span> <span class="tag">{{关键术语2}}</span>
      </div>
    </div>
    <div class="terminal">
      <div class="bar"><span class="dot"></span><span class="dot"></span><span class="dot"></span><span>{{文件名}}</span></div>
<pre><span class="kw">{{代码 6-16 行，语法着色 kw/fn/str/cmt/num}}</span></pre>
    </div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, grid, g2, mt-l, lede, dim, mt-m, tag, terminal, bar, dot, kw, fn, str, cmt, num, mono, notes

---

## code-terminal（满版代码）
指纹：code

用途：一段需要逐行讲解的核心代码，独占一页。
适用 role：code。
内容约束：代码 10-22 行、每行 ≤70 字符；kicker 写文件名与行数（如 `runtime.rs · ~40 LOC`）；标题 ≤16 字；配 1-2 句 dim 解说。

```html
<section class="slide" data-layout="code-terminal">
  <p class="kicker">{{文件名 · 行数}}</p>
  <h2 class="h2">{{标题：这段代码在干什么}}</h2>
  <div class="terminal mt-m">
    <div class="bar"><span class="dot"></span><span class="dot"></span><span class="dot"></span><span>{{文件名}}</span></div>
<pre><span class="kw">{{代码 10-22 行，语法着色用 kw/fn/str/cmt/num}}</span></pre>
  </div>
  <p class="dim mt-m" style="font-size:17px">{{这段代码的关键一行在做什么，20-40 字}}</p>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, terminal, mt-m, bar, dot, kw, fn, str, cmt, num, dim, notes

---

## two-column（两列对比）
指纹：cards

用途：两个方案的对照（优劣/前后/我们 vs 常规）。
适用 role：content。
内容约束：恰好 2 列；每列卡内标题 4-10 字、说明 30-55 字；可用 tag 行给结论（4-8 字）。

```html
<section class="slide" data-layout="two-column">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{对比标题}}</h2>
  <div class="grid g2 mt-l">
    <div class="card"><h4>{{列标题A}}</h4><p class="dim">{{说明，30-55 字}}</p><span class="tag mt-s">{{结论}}</span></div>
    <div class="card card-accent"><h4>{{列标题B}}</h4><p class="dim">{{说明，30-55 字}}</p><span class="tag mt-s">{{结论}}</span></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, grid, g2, mt-l, card, card-accent, h4, dim, tag, mt-s, notes

---

## stat-hero（数据大字报）
指纹：chart

用途：单个核心数字的冲击力呈现（性能倍数、规模、耗时）。
适用 role：data。
内容约束：主数字 ≤8 字符（用 `.mono`，可带单位）；指标名 6-14 字；支撑要点 2-3 条、每条 12-24 字（写口径/对比/出处）。数字必须有出处，没有就用区间或量级。

```html
<section class="slide" data-layout="stat-hero">
  <p class="kicker">{{指标语境}}</p>
  <div class="mono" style="font-size:200px;font-weight:800;color:var(--accent);line-height:1;letter-spacing:-.04em">{{主数字}}</div>
  <h2 class="h2 mt-m">{{指标名，6-14 字}}</h2>
  <div class="stack mt-m">
    <span class="tag">{{支撑要点：口径/对比，12-24 字}}</span>
    <span class="tag">{{支撑要点：出处/年份，12-24 字}}</span>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, mono, h2, mt-m, stack, tag, notes

---

## kpi-grid（KPI 网格）
指纹：chart

用途：4 个并列指标或特性小卡（性能四项、四个能力）。
适用 role：data / content。
内容约束：恰好 4 卡；卡标题（数字或特性名）≤10 字符；说明 15-30 字（写清口径或含义）。

```html
<section class="slide" data-layout="kpi-grid">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{标题}}</h2>
  <div class="grid g4 mt-l">
    <div class="card card-accent"><h4 class="mono" style="color:var(--accent)">{{数字/特性}}</h4><p class="dim">{{说明，15-30 字}}</p></div>
    <div class="card card-accent"><h4 class="mono" style="color:var(--accent)">{{数字/特性}}</h4><p class="dim">{{说明，15-30 字}}</p></div>
    <div class="card card-accent"><h4 class="mono" style="color:var(--accent)">{{数字/特性}}</h4><p class="dim">{{说明，15-30 字}}</p></div>
    <div class="card card-accent"><h4 class="mono" style="color:var(--accent)">{{数字/特性}}</h4><p class="dim">{{说明，15-30 字}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, grid, g4, mt-l, card, card-accent, h4, mono, dim, notes

---

## big-quote（大引用）
指纹：quote

用途：金句、权威结论、用户原话。整页只放一句话。
适用 role：quote。
内容约束：引文 15-40 字；出处 ≤20 字；不用编造的引用，没有真实引文就换版式。

```html
<section class="slide" data-layout="big-quote">
  <div class="mono" style="font-size:120px;color:var(--accent);line-height:.6">"</div>
  <h2 class="h2" style="max-width:62ch">{{引文}}</h2>
  <p class="dim mt-m">—— {{出处}}</p>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mono, h2, dim, mt-m, notes

---

## timeline（时间线）
指纹：table

用途：版本演进、里程碑、历史脉络。
适用 role：content / data。
内容约束：4-6 个节点；每节点 = 时间点（mono）+ 事件 12-24 字。

```html
<section class="slide" data-layout="timeline">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{标题}}</h2>
  <div class="stack mt-l">
    <div class="agenda-row"><span class="num">{{时间点}}</span><span class="t">{{事件，12-24 字}}</span></div>
    <div class="agenda-row"><span class="num">{{时间点}}</span><span class="t">{{事件，12-24 字}}</span></div>
    <div class="agenda-row"><span class="num">{{时间点}}</span><span class="t">{{事件，12-24 字}}</span></div>
    <div class="agenda-row"><span class="num">{{时间点}}</span><span class="t">{{事件，12-24 字}}</span></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, stack, mt-l, agenda-row, num, t, notes

---

## process-steps（流程步骤）
指纹：table

用途：3-5 步的执行流程、请求链路、安装步骤。
适用 role：content。
内容约束：3-5 步；每步 = 步骤号 + 动词短语 6-14 字 + 说明 15-30 字。

```html
<section class="slide" data-layout="process-steps">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{流程名}}</h2>
  <div class="stack mt-l">
    <div class="agenda-row"><span class="num">1.</span><span class="t">{{动词短语}}</span><span class="d">{{说明，15-30 字}}</span></div>
    <div class="agenda-row"><span class="num">2.</span><span class="t">{{动词短语}}</span><span class="d">{{说明，15-30 字}}</span></div>
    <div class="agenda-row"><span class="num">3.</span><span class="t">{{动词短语}}</span><span class="d">{{说明，15-30 字}}</span></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, stack, mt-l, agenda-row, num, t, d, notes

---

## takeaways-3（三件事带回去）
指纹：cards

用途：收束页。3 个记忆点 + 一行延伸阅读。
适用 role：content / cta。
内容约束：恰好 3 卡；卡标题 4-12 字；说明 22-40 字；延伸阅读行 ≤40 字（真实的链接，没有就省略整行）。

```html
<section class="slide" data-layout="takeaways-3">
  <p class="kicker">// takeaways</p>
  <h2 class="h2">{{收束标题，如：三件事带回去}}</h2>
  <div class="grid g3 mt-l">
    <div class="card card-accent"><h4>{{记忆点1}}</h4><p class="dim">{{说明，22-40 字}}</p></div>
    <div class="card card-accent"><h4>{{记忆点2}}</h4><p class="dim">{{说明，22-40 字}}</p></div>
    <div class="card card-accent"><h4>{{记忆点3}}</h4><p class="dim">{{说明，22-40 字}}</p></div>
  </div>
  <p class="lede mt-l">{{延伸阅读（可省略）}}</p>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, grid, g3, mt-l, card, card-accent, h4, dim, lede, notes

---

## qa（问答收尾）
指纹：hero

用途：Q&A / 谢谢观看。居中收尾 + 联系方式标签。
适用 role：thanks / cta。
内容约束：主文案 ≤10 字；lede 15-30 字；标签 ≤2 个。

```html
<section class="slide center tc" data-layout="qa">
  <div>
    <div class="mono" style="font-size:120px;color:var(--accent);font-weight:800;letter-spacing:-.04em">?</div>
    <h2 class="h2">{{主文案}}</h2>
    <p class="lede" style="margin:14px auto">{{一句话收尾，15-30 字}}</p>
    <div class="row mt-l" style="justify-content:center">
      <span class="tag">{{联系方式/仓库}}</span>
    </div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, center, tc, mono, h2, lede, row, mt-l, tag, notes
