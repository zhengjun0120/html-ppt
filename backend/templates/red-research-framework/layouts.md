# 丹朱教研 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `rr-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-red-research-framework` 作用域生效，骨架里已写全，照抄结构即可。
> 右上六边形与左下金黄虚线由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（红是骨架不是墙纸）**：丹朱红只允许出现在骨架元素上——标题带（rr-band 红底白字）、
> 编号（rr-n / rr-tl-dot / rr-num）、徽标（rr-badge）、关键数据（rr-stat-v 与 rr-stat 红左线）、
> 强调药丸（rr-pill-accent）、按钮（rr-btn）、品牌行（rr-brand）。标题带每份 deck 至多两次
> （cover 与 divider）；金黄只做推进虚线与连接，禁给标题描金、禁大面积色块墙。

---

## cover（课题封面）
指纹：hero

用途：开场页。品牌行 + 红色徽标 + 红底白字标题带（主标题 + 副题）+ 一句汇报语境。
适用 role：cover。
内容约束：主标题 ≤12 字；副题 14-26 字；lede 24-48 字；徽标 ≤6 字。

```html
<section class="slide full" data-layout="cover">
  <div class="row" style="justify-content:space-between">
    <p class="rr-brand">{{品牌行，如 SCHOOL RESEARCH · 2026}}</p>
    <span class="rr-badge">{{徽标，≤6 字}}</span>
  </div>
  <div class="rr-band mt-l">
    <h1 class="h1 mt-s">{{主标题，≤12 字}}</h1>
    <p class="rr-band-sub mt-s">{{一句副题，14-26 字}}</p>
  </div>
  <p class="lede mt-m" style="max-width:52ch">{{汇报语境：课题做了什么，24-48 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="rr-pill rr-pill-accent">{{课题编号或属性，≤14 字}}</span>
    <span class="rr-pill">{{汇报人与单位，≤14 字}}</span>
  </div>
  <div class="deck-footer"><span>{{课题组 · 编号}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, row, rr-brand, rr-badge, rr-band, mt-l, h1, mt-s, rr-band-sub, lede, mt-m, rr-pill, rr-pill-accent, deck-footer, slide-number, notes

---

## contents（结题提纲）
指纹：table
数量：rr-item=4

用途：提纲页。一块红描边白卡里放 4 行板块：红圆编号 + 板块名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；板块名 ≤8 字；说明 12-24 字。

```html
<section class="slide" data-layout="contents">
  <p class="rr-brand">{{引导语，如 OUTLINE · 结题提纲}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="rr-card mt-l" style="margin-top:44px">
    <div class="rr-item"><span class="rr-n">01</span><span class="rr-t">{{板块名，≤8 字}}</span><span class="rr-d">{{说明，12-24 字}}</span></div>
    <div class="rr-item"><span class="rr-n">02</span><span class="rr-t">{{板块名，≤8 字}}</span><span class="rr-d">{{说明，12-24 字}}</span></div>
    <div class="rr-item"><span class="rr-n">03</span><span class="rr-t">{{板块名，≤8 字}}</span><span class="rr-d">{{说明，12-24 字}}</span></div>
    <div class="rr-item"><span class="rr-n">04</span><span class="rr-t">{{板块名，≤8 字}}</span><span class="rr-d">{{说明，12-24 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{课题组 · 编号}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rr-brand, h2, mt-m, rr-card, mt-l, rr-item, rr-n, rr-t, rr-d, deck-footer, slide-number, notes

---

## keynotes（三模块要点）
指纹：cards
数量：rr-card=3, rr-num=3

用途：恰好三张红描边白卡。每张：红色大编号 + 小标题 + 两句说明，框架图的三个模块。
适用 role：content。
内容约束：恰好 3 卡；编号 2 字符；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="rr-brand">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="rr-card"><span class="rr-num">01</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#595959">{{说明：一句做法 + 一句机制，22-44 字}}</p></div>
    <div class="rr-card"><span class="rr-num">02</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#595959">{{说明：一句做法 + 一句机制，22-44 字}}</p></div>
    <div class="rr-card"><span class="rr-num">03</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#595959">{{说明：一句做法 + 一句机制，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{课题组 · 编号}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rr-brand, h2, mt-m, grid, g3, mt-l, rr-card, rr-num, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左讲右序）
指纹：split
数量：rr-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边红描边白卡装三行步骤或规程。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="rr-brand">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：机制怎么运转，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#595959">{{补充：判断标准或依据，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="rr-pill">{{要点 1，≤8 字}}</span>
        <span class="rr-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="rr-card">
      <div class="rr-step"><span class="rr-n">01</span><p class="rr-mini-t">{{一步，12-26 字}}</p></div>
      <div class="rr-step"><span class="rr-n">02</span><p class="rr-mini-t">{{一步，12-26 字}}</p></div>
      <div class="rr-step"><span class="rr-n">03</span><p class="rr-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{课题组 · 编号}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rr-brand, h2, mt-m, grid, g2, mt-l, lede, row, rr-pill, rr-card, rr-step, rr-n, rr-mini-t, deck-footer, slide-number, notes

---

## metrics（成效数据）
指纹：chart
数量：rr-stat=3

用途：三个关键数据。淡红底红左线数据块 + 丹朱大数字，口径写进块内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 12-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="rr-brand">{{数据语境，如 OUTCOMES · 成效数据}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:36px;margin-top:48px">
    <div class="rr-stat"><div class="rr-stat-v">{{数值 ≤6 字符}}<span class="rr-stat-u">{{单位}}</span></div><div class="rr-stat-l">{{指标名，≤8 字}}</div><p class="rr-stat-note">{{口径，12-30 字}}</p></div>
    <div class="rr-stat"><div class="rr-stat-v">{{数值 ≤6 字符}}<span class="rr-stat-u">{{单位}}</span></div><div class="rr-stat-l">{{指标名，≤8 字}}</div><p class="rr-stat-note">{{口径，12-30 字}}</p></div>
    <div class="rr-stat"><div class="rr-stat-v">{{数值 ≤6 字符}}<span class="rr-stat-u">{{单位}}</span></div><div class="rr-stat-l">{{指标名，≤8 字}}</div><p class="rr-stat-note">{{口径，12-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#808080;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{课题组 · 编号}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rr-brand, h2, mt-m, grid, g3, mt-l, rr-stat, rr-stat-v, rr-stat-u, rr-stat-l, rr-stat-note, deck-footer, slide-number, notes

---

## quote（教育题记）
指纹：quote

用途：整页一句引文。近黑大字 + 出处 + 两个支撑药丸，教育立场的压轴一笔。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="rr-brand">{{语境，如 EPIGRAPH · 题记}}</p>
  <p class="rr-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="rr-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="rr-pill rr-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="rr-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{课题组 · 编号}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rr-brand, rr-quote, mt-l, rr-src, mt-m, row, rr-pill, rr-pill-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。品牌行 + 徽标 + 红底白字标题带（章节名 + 过渡语）+ 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤12 字；过渡语 14-26 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <div class="row" style="justify-content:space-between">
    <p class="rr-brand">{{进度，如 CHAPTER 03 · 从试点到常态}}</p>
    <span class="rr-badge">{{徽标，≤6 字}}</span>
  </div>
  <div class="rr-band mt-l">
    <h2 class="h2 mt-s">{{章节标题，≤12 字}}</h2>
    <p class="rr-band-sub mt-s">{{一句过渡语，14-26 字}}</p>
  </div>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="rr-pill rr-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="rr-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{课题组 · 编号}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, row, rr-brand, rr-badge, rr-band, mt-l, h2, mt-s, rr-band-sub, lede, mt-m, rr-pill, rr-pill-accent, deck-footer, slide-number, notes

---

## moments（课题历程）
指纹：chart
数量：rr-tl-item=4

用途：3-4 个节点的横向时间线：红圆点编号 + 金黄虚线连接 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="rr-brand">{{引导语，如 TIMELINE · 课题历程}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="rr-tl mt-l" style="margin-top:52px">
    <div class="rr-tl-item"><div class="rr-tl-dot">01</div><div class="rr-tl-t">{{时间点，≤10 字符}}</div><p class="rr-tl-d">{{事件，12-26 字}}</p></div>
    <div class="rr-tl-item"><div class="rr-tl-dot">02</div><div class="rr-tl-t">{{时间点，≤10 字符}}</div><p class="rr-tl-d">{{事件，12-26 字}}</p></div>
    <div class="rr-tl-item"><div class="rr-tl-dot">03</div><div class="rr-tl-t">{{时间点，≤10 字符}}</div><p class="rr-tl-d">{{事件，12-26 字}}</p></div>
    <div class="rr-tl-item"><div class="rr-tl-dot">04</div><div class="rr-tl-t">{{时间点，≤10 字符}}</div><p class="rr-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#808080;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{课题组 · 编号}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rr-brand, h2, mt-m, rr-tl, mt-l, rr-tl-item, rr-tl-dot, rr-tl-t, rr-tl-d, deck-footer, slide-number, notes

---

## closing（结题致谢）
指纹：hero

用途：收尾页。红字大标题 + 一句说明 + 丹朱实底按钮 + 药丸，结题落款。
适用 role：thanks / cta / content。
内容约束：标题 ≤12 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤14 字。

```html
<section class="slide full" data-layout="closing">
  <p class="rr-brand">{{提醒语境，如 CLOSING · 结题致谢}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束说明，20-44 字}}</p>
  <div class="row mt-l" style="gap:22px">
    <span class="rr-btn">{{按钮文案，≤8 字}}</span>
    <span class="rr-pill">{{次级信息，≤14 字}}</span>
  </div>
  <div class="deck-footer"><span>{{课题组 · 编号}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, rr-brand, h1, mt-m, lede, row, mt-l, rr-btn, rr-pill, deck-footer, slide-number, notes
