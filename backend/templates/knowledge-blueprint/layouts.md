# 知识蓝图 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 XXX> 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-knowledge-arch-blueprint` 作用域前缀生效，骨架里已写全，照抄结构即可。

---

## pipeline-overview（流程管线）
指纹：cards

用途：整条链路的全景图（kb-step 逐站，hero 标当前站）
适用 role：content / divider。
内容约束：3-5 站；每站 ≤8 字；hero 标 1 个

```html
<section class="slide" data-layout="pipeline-overview">
  <div class="kb-kicker">{{引导语}}</div>
  <h1 class="kb-h1">{{链路名}}</h1>
  <p class="kb-sub">{{链路一句话}}</p>
  <div class="kb-pipeline mt-l">
    <div class="kb-step">{{站点 1}}</div>
    <div class="kb-step">{{站点 2}}</div>
    <div class="kb-step hero">{{当前站点}}</div>
    <div class="kb-step">{{站点 4}}</div>
  </div>
  <div class="kb-footer">{{页脚：模块名 · 页码}}</div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kb-kicker, kb-h1, kb-sub, kb-pipeline, mt-l, kb-step, hero, kb-footer, notes

---

## section-open（章节开篇）
指纹：hero

用途：章节过渡：大标题 + 一句话引导
适用 role：divider。
内容约束：标题 ≤12 字；副句 20-30 字

```html
<section class="slide" data-layout="section-open">
  <div class="kb-kicker">{{章节进度}}</div>
  <h1 class="kb-h1">{{章节标题}}</h1>
  <p class="kb-sub">{{这一章讲什么，20-30 字}}</p>
  <div class="kb-footer">{{页脚}}</div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kb-kicker, kb-h1, kb-sub, kb-footer, notes

---

## dual-cards（双卡对照）
指纹：cards

用途：两个概念/方案的对照卡（各带小标签行）
适用 role：content。
内容约束：恰好 2 卡；卡标题 ≤8 字；卡内说明 22-45 字（一句论断 + 一句展开）

```html
<section class="slide" data-layout="dual-cards">
  <div class="kb-kicker">{{引导语}}</div>
  <h1 class="kb-h1">{{对照标题}}</h1>
  <div class="kb-grid-2 mt-l">
    <div class="kb-card"><div class="kb-kicker">{{卡标签}}</div><p>{{内容，22-45 字}}</p></div>
    <div class="kb-card"><div class="kb-kicker">{{卡标签}}</div><p>{{内容，22-45 字}}</p></div>
  </div>
  <div class="kb-legend mt-l">{{图例/口径说明}}</div>
  <div class="kb-footer">{{页脚}}</div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kb-kicker, kb-h1, kb-grid-2, mt-l, kb-card, kb-legend, kb-footer, notes

---

## codebox（代码盒）
指纹：code

用途：一段核心代码/配置（等宽盒 + 高亮）
适用 role：code。
内容约束：代码 ≤18 行；kw/st 高亮 span 可用

```html
<section class="slide" data-layout="codebox">
  <div class="kb-kicker">{{文件/场景}}</div>
  <h1 class="kb-h1">{{代码在做什么}}</h1>
  <pre class="kb-codebox mt-l"><span class="kw">{{代码（转义 &lt; &gt; &amp;）}}</span></pre>
  <div class="kb-footer">{{页脚}}</div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kb-kicker, kb-h1, kb-codebox, mt-l, kw, st, kb-footer, notes

---

## statement（机制宣言）
指纹：quote

用途：整页一句话机制总结
适用 role：quote / divider。
内容约束：宣言 12-26 字

```html
<section class="slide" data-layout="statement">
  <div class="kb-kicker">{{引导语}}</div>
  <h1 class="kb-h1">{{机制宣言，12-26 字}}</h1>
  <div class="kb-footer">{{页脚}}</div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kb-kicker, kb-h1, kb-footer, notes

---

## big-num（蓝图大数）
指纹：chart

用途：一个关键数字/量级（kb-big-num 承载）
适用 role：data。
内容约束：数字 ≤8 字符；来源写进正文，无出处标“估算”

```html
<section class="slide" data-layout="big-num">
  <div class="kb-kicker">{{指标语境}}</div>
  <div class="kb-big-num mt-l">{{数字}}</div>
  <h1 class="kb-h1">{{指标名}}</h1>
  <div class="kb-footer">{{页脚}}</div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kb-kicker, kb-big-num, mt-l, kb-h1, kb-footer, notes

---

## insight（洞见卡）
指纹：stack

用途：一条反直觉结论/关键洞见（kb-insight 承载）
适用 role：content。
内容约束：洞见 15-32 字；支撑说明 22-45 字

```html
<section class="slide" data-layout="insight">
  <div class="kb-kicker">{{引导语}}</div>
  <div class="kb-insight mt-l">{{洞见一句话，15-32 字}}</div>
  <h1 class="kb-h1">{{洞见标题}}</h1>
  <p class="kb-sub">{{支撑说明，22-45 字}}</p>
  <div class="kb-footer">{{页脚}}</div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kb-kicker, kb-insight, mt-l, kb-h1, kb-sub, kb-footer, notes

---

## legend（结构分解）
指纹：cards

用途：把一个结构拆成组成部件（kb-step 纵列 + 图例）
适用 role：content。
内容约束：3-5 个部件；部件名 ≤12 字 + 部件说明 12-24 字

```html
<section class="slide" data-layout="legend">
  <div class="kb-kicker">{{引导语}}</div>
  <h1 class="kb-h1">{{结构名}}</h1>
  <div class="kb-pipeline mt-l">
    <div class="kb-step">{{部件 1}}</div>
    <div class="kb-step">{{部件 2}}</div>
    <div class="kb-step">{{部件 3}}</div>
  </div>
  <div class="kb-legend mt-l">{{部件间的衔接说明}}</div>
  <div class="kb-footer">{{页脚}}</div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kb-kicker, kb-h1, kb-pipeline, mt-l, kb-step, kb-legend, kb-footer, notes
