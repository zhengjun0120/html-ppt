# 樱花治愈 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `ss-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-sakura-soft-healing` 作用域生效，骨架里已写全，照抄结构即可。
> 每页顶部的花影花瓣与页底的湖岸草色由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（玫粉是唯一强调）**：玫粉只允许出现在四处——题签圆点之后的强调描边
> （ss-pill-rose）、关键数字（ss-stat-v）、时间点（ss-tl-t）、行动钮（ss-btn）。
> 淡樱只做圆号（ss-n）、圆点（ss-tl-dot）与线（ss-stat 顶线）；湖水蓝与草绿只活在背景
> ambient 里，不进正文与组件。现代柔和日系：无印章、无楷体、无竖排、无高饱和少女粉。

---

## cover（薄樱封面）
指纹：hero

用途：开场页。题签 + 圆润大字标题 + 樱线 + 一句定位。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 26-52 字；题签 ≤14 字；药丸 ≤14 字。

```html
<section class="slide full" data-layout="cover">
  <p class="ss-kicker">{{题签，≤14 字，如 沐樱研肤 · 限定发布}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="ss-line mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，26-52 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="ss-pill ss-pill-rose">{{强调副题，≤14 字}}</span>
    <span class="ss-pill">{{时间或渠道，≤14 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, ss-kicker, h1, mt-m, ss-line, mt-s, lede, mt-m, mt-l, row, ss-pill, ss-pill-rose, deck-footer, slide-number, notes

---

## contents（花瓣目录）
指纹：table
数量：ss-item=4

用途：议程页。一块柔白大卡里放 4 行篇目：樱粉圆号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="ss-kicker">{{引导语，如 今日议程}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="ss-card mt-l" style="margin-top:44px">
    <div class="ss-item"><span class="ss-n">01</span><span class="ss-t">{{篇名，≤8 字}}</span><span class="ss-d">{{说明，14-26 字}}</span></div>
    <div class="ss-item"><span class="ss-n">02</span><span class="ss-t">{{篇名，≤8 字}}</span><span class="ss-d">{{说明，14-26 字}}</span></div>
    <div class="ss-item"><span class="ss-n">03</span><span class="ss-t">{{篇名，≤8 字}}</span><span class="ss-d">{{说明，14-26 字}}</span></div>
    <div class="ss-item"><span class="ss-n">04</span><span class="ss-t">{{篇名，≤8 字}}</span><span class="ss-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ss-kicker, h2, mt-m, ss-card, mt-l, ss-item, ss-n, ss-t, ss-d, deck-footer, slide-number, notes

---

## keynotes（三支单品）
指纹：cards
数量：ss-card=3

用途：恰好三张柔白卡。每张：玫粉题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="ss-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="ss-card"><span class="ss-pill ss-pill-rose">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.85;color:#6E5A4E">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="ss-card"><span class="ss-pill ss-pill-rose">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.85;color:#6E5A4E">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="ss-card"><span class="ss-pill ss-pill-rose">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.85;color:#6E5A4E">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ss-kicker, h2, mt-m, grid, g3, mt-l, ss-card, ss-pill, ss-pill-rose, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：ss-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边柔白卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="ss-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.85;color:#6E5A4E">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="ss-pill">{{要点 1，≤10 字}}</span>
        <span class="ss-pill">{{要点 2，≤10 字}}</span>
      </div>
    </div>
    <div class="ss-card">
      <div class="ss-step"><span class="ss-n">01</span><p class="ss-mini-t">{{一步，12-26 字}}</p></div>
      <div class="ss-step"><span class="ss-n">02</span><p class="ss-mini-t">{{一步，12-26 字}}</p></div>
      <div class="ss-step"><span class="ss-n">03</span><p class="ss-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ss-kicker, h2, mt-m, grid, g2, mt-l, lede, row, ss-pill, ss-card, ss-step, ss-n, ss-mini-t, deck-footer, slide-number, notes

---

## metrics（玫粉数字）
指纹：chart
数量：ss-stat=3

用途：三个关键数据。玫粉大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="ss-kicker">{{数据语境，如 去年樱花季 · 复盘}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="ss-stat"><div class="ss-stat-v">{{数值 ≤6 字符}}<span class="ss-stat-u">{{单位}}</span></div><div class="ss-stat-l">{{指标名，≤8 字}}</div><p class="ss-stat-note">{{口径，14-30 字}}</p></div>
    <div class="ss-stat"><div class="ss-stat-v">{{数值 ≤6 字符}}<span class="ss-stat-u">{{单位}}</span></div><div class="ss-stat-l">{{指标名，≤8 字}}</div><p class="ss-stat-note">{{口径，14-30 字}}</p></div>
    <div class="ss-stat"><div class="ss-stat-v">{{数值 ≤6 字符}}<span class="ss-stat-u">{{单位}}</span></div><div class="ss-stat-l">{{指标名，≤8 字}}</div><p class="ss-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#9C8578;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ss-kicker, h2, mt-m, grid, g3, mt-l, ss-stat, ss-stat-v, ss-stat-u, ss-stat-l, ss-stat-note, deck-footer, slide-number, notes

---

## quote（手记引文）
指纹：quote

用途：整页一句引文。圆润大字 + 出处 + 两个支撑药丸，像研发手记里的一页。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤10 字。

```html
<section class="slide" data-layout="quote">
  <p class="ss-kicker">{{语境，如 研发手记}}</p>
  <p class="ss-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="ss-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="ss-pill ss-pill-rose">{{支撑点 1，≤10 字}}</span>
    <span class="ss-pill">{{支撑点 2，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ss-kicker, ss-quote, mt-l, ss-src, mt-m, row, ss-pill, ss-pill-rose, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="ss-kicker">{{进度，如 第一章 · 为什么是樱前}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="ss-pill ss-pill-rose">{{看点 1，≤8 字}}</span>
    <span class="ss-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, ss-kicker, h1, mt-m, lede, mt-l, row, ss-pill, ss-pill-rose, deck-footer, slide-number, notes

---

## moments（发售时间线）
指纹：chart
数量：ss-tl-item=4

用途：3-4 个节点的横向时间线：樱粉圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="ss-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="ss-tl mt-l" style="margin-top:52px">
    <div class="ss-tl-item"><div class="ss-tl-dot">01</div><div class="ss-tl-t">{{时间点，≤10 字符}}</div><p class="ss-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ss-tl-item"><div class="ss-tl-dot">02</div><div class="ss-tl-t">{{时间点，≤10 字符}}</div><p class="ss-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ss-tl-item"><div class="ss-tl-dot">03</div><div class="ss-tl-t">{{时间点，≤10 字符}}</div><p class="ss-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ss-tl-item"><div class="ss-tl-dot">04</div><div class="ss-tl-t">{{时间点，≤10 字符}}</div><p class="ss-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#9C8578;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ss-kicker, h2, mt-m, ss-tl, mt-l, ss-tl-item, ss-tl-dot, ss-tl-t, ss-tl-d, deck-footer, slide-number, notes

---

## closing（收尾预约）
指纹：hero

用途：收尾页。圆润大字 + 一句行动提醒 + 玫粉圆钮 + 描边药丸，像樱花季的一张预约卡。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤14 字。

```html
<section class="slide full" data-layout="closing">
  <p class="ss-kicker">{{提醒语境，如 会员预约}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="ss-btn">{{按钮文案，≤8 字}}</span>
    <span class="ss-pill">{{次级信息，≤14 字}}</span>
    <span class="ss-pill ss-pill-rose">{{限定提醒，≤14 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, ss-kicker, h1, mt-m, lede, mt-l, row, ss-btn, ss-pill, ss-pill-rose, deck-footer, slide-number, notes
