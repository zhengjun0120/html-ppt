# 渐变美学·樱粉雾蓝 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `dp-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-dreamy-pink-gradient` 作用域生效，骨架里已写全，照抄结构即可。
> 每页的晨昏柔光与页底星子由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（渐变是身份）**：渐变只允许出现在五处——h1 大标题（模板自带渐变填充）、
> 关键数字（dp-stat-v）、编号球（dp-n / dp-tl-dot）、题签短线（dp-kicker 自带）、行动钮（dp-btn）。
> 正文一律深靛实色；禁给 h2/h4/整卡/整页加渐变或彩虹底。月亮（dp-moon）是唯一的发光体，
> 一页至多一枚。磨砂卡半透明白底，禁深色卡、禁硬阴影。

---

## cover（晨昏封面）
指纹：hero

用途：开场页。渐变题签 + 渐变填充大标题 + 一句定位，月亮印与药丸同排。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；副题 ≤12 字。

```html
<section class="slide full" data-layout="cover">
  <p class="dp-kicker">{{题签，≤14 字，如 樱粉雾蓝 · 秋季新品}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="dp-moon"></div>
    <span class="dp-pill">{{副题或时间地点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, dp-kicker, h1, mt-m, lede, mt-l, row, dp-moon, dp-pill, deck-footer, slide-number, notes

---

## contents（星图目录）
指纹：table
数量：dp-item=4

用途：议程页。一块磨砂云卡里放 4 行篇目：渐变编号球 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="dp-kicker">{{引导语，如 今日星图}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="dp-card mt-l" style="margin-top:44px">
    <div class="dp-item"><span class="dp-n">壹</span><span class="dp-t">{{篇名，≤8 字}}</span><span class="dp-d">{{说明，14-26 字}}</span></div>
    <div class="dp-item"><span class="dp-n">贰</span><span class="dp-t">{{篇名，≤8 字}}</span><span class="dp-d">{{说明，14-26 字}}</span></div>
    <div class="dp-item"><span class="dp-n">叁</span><span class="dp-t">{{篇名，≤8 字}}</span><span class="dp-d">{{说明，14-26 字}}</span></div>
    <div class="dp-item"><span class="dp-n">肆</span><span class="dp-t">{{篇名，≤8 字}}</span><span class="dp-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, dp-kicker, h2, mt-m, dp-card, mt-l, dp-item, dp-n, dp-t, dp-d, deck-footer, slide-number, notes

---

## keynotes（三云要点）
指纹：cards
数量：dp-card=3

用途：恰好三张磨砂云卡。每张：玫粉题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="dp-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="dp-card"><span class="dp-pill dp-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#5A669C">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="dp-card"><span class="dp-pill dp-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#5A669C">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="dp-card"><span class="dp-pill dp-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#5A669C">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, dp-kicker, h2, mt-m, grid, g3, mt-l, dp-card, dp-pill, dp-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右云）
指纹：split
数量：dp-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边磨砂卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="dp-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#5A669C">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="dp-pill">{{要点 1，≤8 字}}</span>
        <span class="dp-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="dp-card">
      <div class="dp-step"><span class="dp-n">一</span><p class="dp-mini-t">{{一步，12-26 字}}</p></div>
      <div class="dp-step"><span class="dp-n">二</span><p class="dp-mini-t">{{一步，12-26 字}}</p></div>
      <div class="dp-step"><span class="dp-n">三</span><p class="dp-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, dp-kicker, h2, mt-m, grid, g2, mt-l, lede, row, dp-pill, dp-card, dp-step, dp-n, dp-mini-t, deck-footer, slide-number, notes

---

## metrics（渐变数字）
指纹：chart
数量：dp-stat=3

用途：三个关键数据。渐变大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="dp-kicker">{{数据语境，如 内测一个月 · 复盘}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="dp-stat"><div class="dp-stat-v">{{数值 ≤6 字符}}<span class="dp-stat-u">{{单位}}</span></div><div class="dp-stat-l">{{指标名，≤8 字}}</div><p class="dp-stat-note">{{口径，14-30 字}}</p></div>
    <div class="dp-stat"><div class="dp-stat-v">{{数值 ≤6 字符}}<span class="dp-stat-u">{{单位}}</span></div><div class="dp-stat-l">{{指标名，≤8 字}}</div><p class="dp-stat-note">{{口径，14-30 字}}</p></div>
    <div class="dp-stat"><div class="dp-stat-v">{{数值 ≤6 字符}}<span class="dp-stat-u">{{单位}}</span></div><div class="dp-stat-l">{{指标名，≤8 字}}</div><p class="dp-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="margin-top:44px;font-size:17px;color:#7A82B0;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, dp-kicker, h2, mt-m, grid, g3, mt-l, dp-stat, dp-stat-v, dp-stat-u, dp-stat-l, dp-stat-note, deck-footer, slide-number, notes

---

## quote（星语引文）
指纹：quote

用途：整页一句引文。衬线大字 + 出处 + 两个支撑药丸，留白处立一枚月亮。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="dp-kicker">{{语境，如 主理人手记}}</p>
  <p class="dp-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="dp-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="dp-pill dp-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="dp-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="dp-moon" style="position:absolute;right:140px;top:120px"></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, dp-kicker, dp-quote, mt-l, dp-src, mt-m, row, dp-pill, dp-pill-accent, dp-moon, deck-footer, slide-number, notes

---

## divider（章节月幕）
指纹：hero

用途：章节过渡。进度题签 + 渐变章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="dp-kicker">{{进度，如 卷二 · 新品三件}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="dp-pill dp-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="dp-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, dp-kicker, h1, mt-m, lede, mt-l, row, dp-pill, dp-pill-accent, deck-footer, slide-number, notes

---

## moments（星轨时间线）
指纹：chart
数量：dp-tl-item=4

用途：3-4 个节点的横向时间线：渐变圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="dp-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="dp-tl mt-l" style="margin-top:52px">
    <div class="dp-tl-item"><div class="dp-tl-dot">壹</div><div class="dp-tl-t">{{时间点，≤10 字符}}</div><p class="dp-tl-d">{{事件，12-26 字}}</p></div>
    <div class="dp-tl-item"><div class="dp-tl-dot">贰</div><div class="dp-tl-t">{{时间点，≤10 字符}}</div><p class="dp-tl-d">{{事件，12-26 字}}</p></div>
    <div class="dp-tl-item"><div class="dp-tl-dot">叁</div><div class="dp-tl-t">{{时间点，≤10 字符}}</div><p class="dp-tl-d">{{事件，12-26 字}}</p></div>
    <div class="dp-tl-item"><div class="dp-tl-dot">肆</div><div class="dp-tl-t">{{时间点，≤10 字符}}</div><p class="dp-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="margin-top:44px;font-size:17px;color:#7A82B0;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, dp-kicker, h2, mt-m, dp-tl, mt-l, dp-tl-item, dp-tl-dot, dp-tl-t, dp-tl-d, deck-footer, slide-number, notes

---

## closing（暮色收尾）
指纹：hero

用途：收尾页。渐变大字 + 一句行动提醒 + 渐变按钮 + 描边药丸 + 月亮印。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="dp-kicker">{{提醒语境，如 预售提醒}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="dp-btn">{{按钮文案，≤8 字}}</span>
    <span class="dp-pill">{{次级信息，≤10 字}}</span>
    <div class="dp-moon"></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, dp-kicker, h1, mt-m, lede, mt-l, row, dp-btn, dp-pill, dp-moon, deck-footer, slide-number, notes
