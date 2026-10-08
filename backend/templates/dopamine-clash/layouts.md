# 多巴胺活力撞色 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `dc-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-dopamine-clash` 作用域生效，骨架里已写全，照抄结构即可。
> 每页的玫粉圆角、圆点矩阵与底边三色条由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（撞色有配比，正文守对比度）**：淡紫与亮绿是大面积打底色（上面只压黑字），
> 玫粉只做小面积焦点——硬阴影、贴纸、时间点（dc-tl-t）与菱形糖（dc-diamond），
> 玫粉上的文字一律近黑。白字只允许出现在黑底（dc-kicker / dc-btn）上。
> 黑描边 + 硬阴影是全部卡与钮的统一画法；菱形糖一页一颗，禁再手画色块。

---

## cover（撞色封面）
指纹：hero

用途：开场页。黑底题签 + 特粗大标题 + 一句定位，菱形糖与旋转贴纸同排。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；贴纸 ≤14 字。

```html
<section class="slide full" data-layout="cover">
  <p class="dc-kicker">{{题签，≤14 字，如 SUGAR HIGH · 2026 秋}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:20px">
    <div class="dc-diamond"></div>
    <div class="dc-stamp">{{贴纸文案，≤14 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, dc-kicker, h1, mt-m, lede, mt-l, row, dc-diamond, dc-stamp, deck-footer, slide-number, notes

---

## contents（色卡目录）
指纹：table
数量：dc-item=4

用途：议程页。一块描边大卡里放 4 行篇目：撞色编号块 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="dc-kicker">{{引导语，如 本月企划}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="dc-card mt-l" style="margin-top:44px">
    <div class="dc-item"><span class="dc-n">01</span><span class="dc-t">{{篇名，≤8 字}}</span><span class="dc-d">{{说明，14-26 字}}</span></div>
    <div class="dc-item"><span class="dc-n">02</span><span class="dc-t">{{篇名，≤8 字}}</span><span class="dc-d">{{说明，14-26 字}}</span></div>
    <div class="dc-item"><span class="dc-n">03</span><span class="dc-t">{{篇名，≤8 字}}</span><span class="dc-d">{{说明，14-26 字}}</span></div>
    <div class="dc-item"><span class="dc-n">04</span><span class="dc-t">{{篇名，≤8 字}}</span><span class="dc-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, dc-kicker, h2, mt-m, dc-card, mt-l, dc-item, dc-n, dc-t, dc-d, deck-footer, slide-number, notes

---

## keynotes（三糖要点）
指纹：cards
数量：dc-card=3

用途：恰好三张描边硬阴影卡。每张：撞色题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="dc-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:36px;margin-top:44px">
    <div class="dc-card"><span class="dc-pill dc-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#4A4763">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="dc-card"><span class="dc-pill dc-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#4A4763">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="dc-card"><span class="dc-pill dc-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#4A4763">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, dc-kicker, h2, mt-m, grid, g3, mt-l, dc-card, dc-pill, dc-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右格）
指纹：split
数量：dc-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边描边卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="dc-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#4A4763">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="dc-pill">{{要点 1，≤8 字}}</span>
        <span class="dc-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="dc-card">
      <div class="dc-step"><span class="dc-n">01</span><p class="dc-mini-t">{{一步，12-26 字}}</p></div>
      <div class="dc-step"><span class="dc-n">02</span><p class="dc-mini-t">{{一步，12-26 字}}</p></div>
      <div class="dc-step"><span class="dc-n">03</span><p class="dc-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, dc-kicker, h2, mt-m, grid, g2, mt-l, lede, row, dc-pill, dc-card, dc-step, dc-n, dc-mini-t, deck-footer, slide-number, notes

---

## metrics（像素数字）
指纹：chart
数量：dc-stat=3

用途：三个关键数据。等宽大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="dc-kicker">{{数据语境，如 会员优先购 · 72 小时}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="dc-stat"><div class="dc-stat-v">{{数值 ≤6 字符}}<span class="dc-stat-u">{{单位}}</span></div><div class="dc-stat-l">{{指标名，≤8 字}}</div><p class="dc-stat-note">{{口径，14-30 字}}</p></div>
    <div class="dc-stat"><div class="dc-stat-v">{{数值 ≤6 字符}}<span class="dc-stat-u">{{单位}}</span></div><div class="dc-stat-l">{{指标名，≤8 字}}</div><p class="dc-stat-note">{{口径，14-30 字}}</p></div>
    <div class="dc-stat"><div class="dc-stat-v">{{数值 ≤6 字符}}<span class="dc-stat-u">{{单位}}</span></div><div class="dc-stat-l">{{指标名，≤8 字}}</div><p class="dc-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="margin-top:44px;font-size:17px;color:#6A6787;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, dc-kicker, h2, mt-m, grid, g3, mt-l, dc-stat, dc-stat-v, dc-stat-u, dc-stat-l, dc-stat-note, deck-footer, slide-number, notes

---

## quote（高亮引文）
指纹：quote

用途：整页一句引文。特粗大字自带荧光笔高亮 + 出处 + 两个支撑药丸。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="dc-kicker">{{语境，如 企划开篇}}</p>
  <p class="dc-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="dc-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="dc-pill dc-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="dc-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="dc-diamond" style="position:absolute;right:170px;top:110px"></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, dc-kicker, dc-quote, mt-l, dc-src, mt-m, row, dc-pill, dc-pill-accent, dc-diamond, deck-footer, slide-number, notes

---

## divider（章节切换）
指纹：hero

用途：章节过渡。进度题签 + 特粗章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="dc-kicker">{{进度，如 卷二 · 本季色盘}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="dc-pill dc-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="dc-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, dc-kicker, h1, mt-m, lede, mt-l, row, dc-pill, dc-pill-accent, deck-footer, slide-number, notes

---

## moments（快闪时间线）
指纹：chart
数量：dc-tl-item=4

用途：3-4 个节点的横向时间线：贴纸圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="dc-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="dc-tl mt-l" style="margin-top:52px">
    <div class="dc-tl-item"><div class="dc-tl-dot">01</div><div class="dc-tl-t">{{时间点，≤10 字符}}</div><p class="dc-tl-d">{{事件，12-26 字}}</p></div>
    <div class="dc-tl-item"><div class="dc-tl-dot">02</div><div class="dc-tl-t">{{时间点，≤10 字符}}</div><p class="dc-tl-d">{{事件，12-26 字}}</p></div>
    <div class="dc-tl-item"><div class="dc-tl-dot">03</div><div class="dc-tl-t">{{时间点，≤10 字符}}</div><p class="dc-tl-d">{{事件，12-26 字}}</p></div>
    <div class="dc-tl-item"><div class="dc-tl-dot">04</div><div class="dc-tl-t">{{时间点，≤10 字符}}</div><p class="dc-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="margin-top:44px;font-size:17px;color:#6A6787;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, dc-kicker, h2, mt-m, dc-tl, mt-l, dc-tl-item, dc-tl-dot, dc-tl-t, dc-tl-d, deck-footer, slide-number, notes

---

## closing（落幕按钮）
指纹：hero

用途：收尾页。特粗大字 + 一句行动提醒 + 黑底硬阴影按钮 + 描边药丸 + 菱形糖。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="dc-kicker">{{提醒语境，如 快闪预约}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:24px">
    <span class="dc-btn">{{按钮文案，≤8 字}}</span>
    <span class="dc-pill">{{次级信息，≤10 字}}</span>
    <div class="dc-diamond"></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, dc-kicker, h1, mt-m, lede, mt-l, row, dc-btn, dc-pill, dc-diamond, deck-footer, slide-number, notes
