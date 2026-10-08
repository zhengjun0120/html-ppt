# 瑞士国际主义 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `si-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-swiss-international` 作用域生效，骨架里已写全，照抄结构即可。
> 每页的极细灰网格由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（结构即视觉）**：全模板只用一种强调色——IKB 蓝 #002FA7（si-square / si-accent /
> si-metric-accent / si-cell-accent / si-tag-accent / si-btn）。禁阴影、禁圆角、禁渐变、
> 禁圆点；标题细体大字、小字加粗，一切左对齐；每页页首必须有一条 si-meta 元信息行。

---

## cover（索引封面）
指纹：hero

用途：开场页。mono 元信息行 + 细体巨字标题（下行套蓝）+ 通栏细线 + 三格 KPI 条。
适用 role：cover。
内容约束：主标题上行 ≤6 字、下行强调词 ≤4 字；lede 30-60 字；3 个数值各 ≤6 字符；指标行 ≤10 字。

```html
<section class="slide full" data-layout="cover">
  <div class="si-meta"><span>{{左标签，如 BERICHT / GRID}}</span><span class="si-square"></span><span>{{右编号，如 01 / 09}}</span></div>
  <h1 class="h1 mt-l">{{主标题上行，≤6 字}}<br><span class="si-accent">{{下行强调词，≤4 字}}</span></h1>
  <div class="si-line mt-m"></div>
  <p class="lede mt-m" style="max-width:52ch">{{报告定位：做了什么、覆盖什么，30-60 字}}</p>
  <div class="grid g3 si-mesh mt-l" style="margin-top:48px">
    <div class="si-metric si-metric-accent"><div class="si-metric-v">{{数值 ≤6 字符}}</div><div class="si-metric-l">{{指标行 ≤10 字}}</div></div>
    <div class="si-metric"><div class="si-metric-v">{{数值 ≤6 字符}}</div><div class="si-metric-l">{{指标行 ≤10 字}}</div></div>
    <div class="si-metric"><div class="si-metric-v">{{数值 ≤6 字符}}</div><div class="si-metric-l">{{指标行 ≤10 字}}</div></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, si-meta, si-square, h1, mt-l, si-accent, si-line, mt-m, lede, grid, g3, si-mesh, si-metric, si-metric-accent, si-metric-v, si-metric-l, deck-footer, slide-number, notes

---

## contents（目录表）
指纹：table
数量：si-toc-row=4

用途：议程页。hairline 分隔的 4 行目录：mono 蓝色索引 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <div class="si-meta"><span>INDEX</span><span class="si-square"></span><span>{{右编号，如 02 / 09}}</span></div>
  <h2 class="h2 mt-m">{{标题，≤6 字}}</h2>
  <div class="mt-l" style="margin-top:40px">
    <div class="si-toc-row"><span class="si-i">01</span><span class="si-t">{{篇名，≤8 字}}</span><span class="si-d">{{说明，14-26 字}}</span></div>
    <div class="si-toc-row"><span class="si-i">02</span><span class="si-t">{{篇名，≤8 字}}</span><span class="si-d">{{说明，14-26 字}}</span></div>
    <div class="si-toc-row"><span class="si-i">03</span><span class="si-t">{{篇名，≤8 字}}</span><span class="si-d">{{说明，14-26 字}}</span></div>
    <div class="si-toc-row"><span class="si-i">04</span><span class="si-t">{{篇名，≤8 字}}</span><span class="si-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, si-meta, si-square, h2, mt-m, mt-l, si-toc-row, si-i, si-t, si-d, deck-footer, slide-number, notes

---

## metrics（KPI 塔）
指纹：chart
数量：si-metric=3

用途：三个关键数据。黑缝拼合三格，第一格是蓝色焦点；细体大数字 + mono 指标行，口径写页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标行 ≤10 字；来源 14-40 字。

```html
<section class="slide" data-layout="metrics">
  <div class="si-meta"><span>KPI / {{数据语境}}</span><span class="si-square"></span><span>{{右编号}}</span></div>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 si-mesh mt-l" style="margin-top:44px">
    <div class="si-metric si-metric-accent"><div class="si-metric-v">{{数值 ≤6 字符}}</div><div class="si-metric-l">{{指标行 ≤10 字}}</div></div>
    <div class="si-metric"><div class="si-metric-v">{{数值 ≤6 字符}}</div><div class="si-metric-l">{{指标行 ≤10 字}}</div></div>
    <div class="si-metric"><div class="si-metric-v">{{数值 ≤6 字符}}</div><div class="si-metric-l">{{指标行 ≤10 字}}</div></div>
  </div>
  <p class="si-src mt-m">来源：{{出处与统计口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, si-meta, si-square, h2, mt-m, grid, g3, si-mesh, mt-l, si-metric, si-metric-accent, si-metric-v, si-metric-l, si-src, mt-m, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度元信息 + 细体大字章节名 + 细线 + 过渡问题，右下浅灰巨号。
适用 role：divider。
内容约束：章节名上行 ≤6 字、下行 ≤4 字；过渡问题 22-44 字；看点各 ≤8 字；巨号 2 字符。

```html
<section class="slide full" data-layout="divider">
  <div class="si-meta"><span>{{进度，如 KAPITEL 02 / FARBE}}</span><span class="si-square"></span><span>{{右编号}}</span></div>
  <h1 class="h1 mt-l">{{章节上行，≤6 字}}<br><span class="si-accent">{{章节下行，≤4 字}}</span></h1>
  <div class="si-line mt-m" style="width:220px"></div>
  <p class="lede mt-m" style="max-width:44ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:14px">
    <span class="si-tag si-tag-accent">{{看点 1，≤8 字}}</span>
    <span class="si-tag">{{看点 2，≤8 字}}</span>
  </div>
  <div class="si-num" style="position:absolute;right:96px;bottom:130px">{{章节号，如 02}}</div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, si-meta, si-square, h1, mt-l, si-accent, si-line, mt-m, lede, row, mt-l, si-tag, si-tag-accent, si-num, deck-footer, slide-number, notes

---

## keynotes（定义三格）
指纹：cards
数量：si-cell=3

用途：恰好三个定义格（黑缝拼合）。中间格做蓝色焦点；mono 索引 + 粗小标题 + 细说明。
适用 role：content。
内容约束：恰好 3 格；索引 01-03；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <div class="si-meta"><span>SYSTEM / {{章节标签}}</span><span class="si-square"></span><span>{{右编号}}</span></div>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 si-mesh mt-l" style="margin-top:44px">
    <div class="si-cell"><p class="si-cell-i">01</p><h4 class="mt-s">{{小标题，≤8 字}}</h4><p class="si-cell-d">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="si-cell si-cell-accent"><p class="si-cell-i">02</p><h4 class="mt-s">{{小标题，≤8 字}}</h4><p class="si-cell-d">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="si-cell"><p class="si-cell-i">03</p><h4 class="mt-s">{{小标题，≤8 字}}</h4><p class="si-cell-d">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, si-meta, si-square, h2, mt-m, grid, g3, si-mesh, mt-l, si-cell, si-cell-accent, si-cell-i, h4, mt-s, si-cell-d, deck-footer, slide-number, notes

---

## split（分屏论点）
指纹：split
数量：si-step=3

用途：左边把论点讲透（lede + 补充 + 方角标签），右边 hairline 三行步骤。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <div class="si-meta"><span>METHOD / {{引导语}}</span><span class="si-square"></span><span>{{右编号}}</span></div>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:18px;line-height:1.8;color:#737373">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:14px">
        <span class="si-tag">{{要点 1，≤8 字}}</span>
        <span class="si-tag">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div>
      <div class="si-step"><span class="si-i">01</span><p class="si-mini-t">{{一步，12-26 字}}</p></div>
      <div class="si-step"><span class="si-i">02</span><p class="si-mini-t">{{一步，12-26 字}}</p></div>
      <div class="si-step"><span class="si-i">03</span><p class="si-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, si-meta, si-square, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, si-tag, si-step, si-i, si-mini-t, deck-footer, slide-number, notes

---

## moments（时间线）
指纹：chart
数量：si-tl-item=4

用途：3-4 个节点的横向时间线：黑色通栏线 + 直角蓝方块节点 + mono 时间点。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <div class="si-meta"><span>TIMELINE / {{引导语}}</span><span class="si-square"></span><span>{{右编号}}</span></div>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="si-tl mt-l" style="margin-top:56px">
    <div class="si-tl-item"><div class="si-tl-t">{{时间点，≤10 字符}}</div><p class="si-tl-d">{{事件，12-26 字}}</p></div>
    <div class="si-tl-item"><div class="si-tl-t">{{时间点，≤10 字符}}</div><p class="si-tl-d">{{事件，12-26 字}}</p></div>
    <div class="si-tl-item"><div class="si-tl-t">{{时间点，≤10 字符}}</div><p class="si-tl-d">{{事件，12-26 字}}</p></div>
    <div class="si-tl-item"><div class="si-tl-t">{{时间点，≤10 字符}}</div><p class="si-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="si-src mt-m">{{口径或读法，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, si-meta, si-square, h2, mt-m, si-tl, mt-l, si-tl-item, si-tl-t, si-tl-d, si-src, mt-m, deck-footer, slide-number, notes

---

## quote（宣言页）
指纹：quote

用途：整页一句宣言。细体大字两行（下行套蓝）+ mono 出处 + 两个方角标签。
适用 role：quote。
内容约束：引文共 10-28 字；出处 ≤22 字且真实；标签各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <div class="si-meta"><span>MANIFEST / {{语境}}</span><span class="si-square"></span><span>{{右编号}}</span></div>
  <p class="si-quote mt-l" style="margin-top:48px">「{{引文上行，10-16 字}}<br><span class="si-accent">{{引文下行，≤12 字}}</span>」</p>
  <p class="si-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:14px">
    <span class="si-tag si-tag-accent">{{支撑点 1，≤8 字}}</span>
    <span class="si-tag">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, si-meta, si-square, si-quote, mt-l, si-accent, si-src, mt-m, row, si-tag, si-tag-accent, deck-footer, slide-number, notes

---

## closing（收尾）
指纹：hero

用途：收尾页。细体大字 + 通栏细线 + 一句下一步 + 蓝底方角行动钮。
适用 role：thanks / cta / content。
内容约束：收束上行 ≤6 字、下行 ≤4 字；lede 20-48 字；按钮 ≤6 字；标签各 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <div class="si-meta"><span>NÄCHSTE / {{语境}}</span><span class="si-square"></span><span>{{右编号，如 09 / 09}}</span></div>
  <h1 class="h1 mt-l">{{收束上行，≤6 字}}<br><span class="si-accent">{{收束下行，≤4 字}}</span></h1>
  <div class="si-line mt-m"></div>
  <p class="lede mt-m" style="max-width:44ch">{{下一步与承诺，20-48 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="si-btn">{{按钮文案，≤6 字}}</span>
    <span class="si-tag">{{次级信息，≤10 字}}</span>
    <span class="si-tag">{{次级信息，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, si-meta, si-square, h1, mt-l, si-accent, si-line, mt-m, lede, row, si-btn, si-tag, deck-footer, slide-number, notes
