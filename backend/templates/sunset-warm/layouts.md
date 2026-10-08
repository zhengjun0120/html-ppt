# 日落暖 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `su-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-sunset-warm` 作用域生效，骨架里已写全，照抄结构即可。
> 每页的天际渐变、云条与右上太阳光环由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（暖而不压迫）**：亮橙与金黄只做光泽——题签圆点（su-kicker）、太阳印（su-badge）、
> 金环圆号（su-n / su-tl-dot）、数字块上边线（su-stat 顶线）与 accent 药丸软底。深焦糖 #7C2D12
> 才是标题与正文的骨架色。禁冷色与暗色背景、禁大面积实底橙、禁正式衬线排正文；
> 白卡（su-card / su-stat）是暖色上唯一的「实底」，内容尽量进卡。

---

## cover（落日封面）
指纹：hero

用途：开场页。等宽字题签 + 大字标题 + 一句定位，右侧可立一枚太阳印。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；太阳印 4 字内。

```html
<section class="slide full" data-layout="cover">
  <p class="su-kicker">{{题签，≤14 字，如 WEDDING PLAN · 落日仪式}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="su-badge">{{太阳印 2×2 字}}</div>
    <span class="su-pill">{{副题或时间地点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, su-kicker, h1, mt-m, lede, row, mt-l, su-badge, su-pill, deck-footer, slide-number, notes

---

## contents（目录白卡）
指纹：table
数量：su-item=4

用途：议程页。一块白卡里放 4 行篇目：金环圆号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="su-kicker">{{引导语，如 AGENDA · 流程目录}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="su-card mt-l" style="margin-top:44px">
    <div class="su-item"><span class="su-n">01</span><span class="su-t">{{篇名，≤8 字}}</span><span class="su-d">{{说明，14-26 字}}</span></div>
    <div class="su-item"><span class="su-n">02</span><span class="su-t">{{篇名，≤8 字}}</span><span class="su-d">{{说明，14-26 字}}</span></div>
    <div class="su-item"><span class="su-n">03</span><span class="su-t">{{篇名，≤8 字}}</span><span class="su-d">{{说明，14-26 字}}</span></div>
    <div class="su-item"><span class="su-n">04</span><span class="su-t">{{篇名，≤8 字}}</span><span class="su-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, su-kicker, h2, mt-m, su-card, mt-l, su-item, su-n, su-t, su-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：su-card=3

用途：恰好三张白卡。每张：金黄题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="su-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="su-card"><span class="su-pill su-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#9A3412">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="su-card"><span class="su-pill su-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#9A3412">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="su-card"><span class="su-pill su-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#9A3412">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, su-kicker, h2, mt-m, grid, g3, mt-l, su-card, su-pill, su-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：su-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边白卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="su-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#9A3412">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="su-pill">{{要点 1，≤8 字}}</span>
        <span class="su-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="su-card">
      <div class="su-step"><span class="su-n">一</span><p class="su-mini-t">{{一步，12-26 字}}</p></div>
      <div class="su-step"><span class="su-n">二</span><p class="su-mini-t">{{一步，12-26 字}}</p></div>
      <div class="su-step"><span class="su-n">三</span><p class="su-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, su-kicker, h2, mt-m, grid, g2, mt-l, lede, row, su-pill, su-card, su-step, su-n, su-mini-t, deck-footer, slide-number, notes

---

## metrics（鎏金数字）
指纹：chart
数量：su-stat=3

用途：三个关键数据。白卡数字块 + 亮橙大数字是主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="su-kicker">{{数据语境，如 KEY NUMBERS · 预算}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="su-stat"><div class="su-stat-v">{{数值 ≤6 字符}}<span class="su-stat-u">{{单位}}</span></div><div class="su-stat-l">{{指标名，≤8 字}}</div><p class="su-stat-note">{{口径，14-30 字}}</p></div>
    <div class="su-stat"><div class="su-stat-v">{{数值 ≤6 字符}}<span class="su-stat-u">{{单位}}</span></div><div class="su-stat-l">{{指标名，≤8 字}}</div><p class="su-stat-note">{{口径，14-30 字}}</p></div>
    <div class="su-stat"><div class="su-stat-v">{{数值 ≤6 字符}}<span class="su-stat-u">{{单位}}</span></div><div class="su-stat-l">{{指标名，≤8 字}}</div><p class="su-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#7C2D12;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, su-kicker, h2, mt-m, grid, g3, mt-l, su-stat, su-stat-v, su-stat-u, su-stat-l, su-stat-note, deck-footer, slide-number, notes

---

## quote（光环引文）
指纹：quote

用途：整页一句引文。衬线大字 + 出处 + 两个支撑药丸，右侧留白处浮一轮太阳光环。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="su-kicker">{{语境，如 策划手记}}</p>
  <p class="su-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="su-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="su-pill su-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="su-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="su-halo" aria-hidden="true"></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, su-kicker, su-quote, mt-l, su-src, mt-m, row, su-pill, su-pill-accent, su-halo, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。等宽字题签 + 大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="su-kicker">{{进度，如 第三幕 · 光线与誓言}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="su-pill su-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="su-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, su-kicker, h1, mt-m, lede, mt-l, row, su-pill, su-pill-accent, deck-footer, slide-number, notes

---

## moments（流程时间线）
指纹：chart
数量：su-tl-item=4

用途：3-4 个节点的横向时间线：白圈金环圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="su-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="su-tl mt-l" style="margin-top:52px">
    <div class="su-tl-item"><div class="su-tl-dot">壹</div><div class="su-tl-t">{{时间点，≤10 字符}}</div><p class="su-tl-d">{{事件，12-26 字}}</p></div>
    <div class="su-tl-item"><div class="su-tl-dot">贰</div><div class="su-tl-t">{{时间点，≤10 字符}}</div><p class="su-tl-d">{{事件，12-26 字}}</p></div>
    <div class="su-tl-item"><div class="su-tl-dot">叁</div><div class="su-tl-t">{{时间点，≤10 字符}}</div><p class="su-tl-d">{{事件，12-26 字}}</p></div>
    <div class="su-tl-item"><div class="su-tl-dot">肆</div><div class="su-tl-t">{{时间点，≤10 字符}}</div><p class="su-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#7C2D12;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, su-kicker, h2, mt-m, su-tl, mt-l, su-tl-item, su-tl-dot, su-tl-t, su-tl-d, deck-footer, slide-number, notes

---

## closing（收尾暖印）
指纹：hero

用途：收尾页。大字标题 + 一句行动提醒 + 焦橙渐变按钮 + 药丸 + 太阳印。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="su-kicker">{{提醒语境，如 档期确认}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="su-btn">{{按钮文案，≤8 字}}</span>
    <span class="su-pill">{{次级信息，≤10 字}}</span>
    <div class="su-badge">{{太阳印 2 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, su-kicker, h1, mt-m, lede, mt-l, row, su-btn, su-pill, su-badge, deck-footer, slide-number, notes
