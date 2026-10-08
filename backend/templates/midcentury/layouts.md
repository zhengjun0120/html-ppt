# 世纪中叶 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `mc-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-midcentury` 作用域生效，骨架里已写全，照抄结构即可。
> 每页右上有机形状簇与左下原子射线由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（温暖即秩序）**：芥末黄/焦橙/青绿只按分工出现——芥末黄管题签圆点、横条与数字上缘，
> 青绿管药丸重音与行动钮，焦橙只给关键数据与时间点（mc-stat-v / mc-tl-t）。禁硬阴影、禁冷色调、
> 禁锐利直角几何。几何无衬线大标题（h1/h2）每页最多一组；奶油渐变是画布，不要再叠色块。

---

## cover（展厅入口）
指纹：hero

用途：开场页。标本签 + 几何大标题 + 三色横条 + 一句定位，左下立一枚年份圆章。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 30-60 字；题签 ≤14 字；圆章 4-8 字符。

```html
<section class="slide full" data-layout="cover">
  <p class="mc-kicker">{{题签，≤14 字，如 中世纪现代家具展 · 导览手册}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="mc-line mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，30-60 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="mc-badge">{{年份或代号，4-8 字符}}</div>
    <span class="mc-pill">{{时间地点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, mc-kicker, h1, mt-m, mc-line, mt-s, lede, mt-l, row, mc-badge, mc-pill, deck-footer, slide-number, notes

---

## contents（观展目录）
指纹：table
数量：mc-item=4

用途：议程页。一张奶白大卡里放 4 行篇目：圆号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="mc-kicker">{{引导语，如 今日动线}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="mc-card mt-l" style="margin-top:44px">
    <div class="mc-item"><span class="mc-n">01</span><span class="mc-t">{{篇名，≤8 字}}</span><span class="mc-d">{{说明，14-26 字}}</span></div>
    <div class="mc-item"><span class="mc-n">02</span><span class="mc-t">{{篇名，≤8 字}}</span><span class="mc-d">{{说明，14-26 字}}</span></div>
    <div class="mc-item"><span class="mc-n">03</span><span class="mc-t">{{篇名，≤8 字}}</span><span class="mc-d">{{说明，14-26 字}}</span></div>
    <div class="mc-item"><span class="mc-n">04</span><span class="mc-t">{{篇名，≤8 字}}</span><span class="mc-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mc-kicker, h2, mt-m, mc-card, mt-l, mc-item, mc-n, mc-t, mc-d, deck-footer, slide-number, notes

---

## keynotes（三件经典）
指纹：cards
数量：mc-card=3

用途：恰好三张奶白卡。每张：青绿题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤10 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="mc-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="mc-card"><span class="mc-pill mc-pill-accent">{{题签，≤10 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#92400E">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="mc-card"><span class="mc-pill mc-pill-accent">{{题签，≤10 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#92400E">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="mc-card"><span class="mc-pill mc-pill-accent">{{题签，≤10 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#92400E">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mc-kicker, h2, mt-m, grid, g3, mt-l, mc-card, mc-pill, mc-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（细读一件）
指纹：split
数量：mc-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边奶白卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="mc-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#92400E">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="mc-pill">{{要点 1，≤8 字}}</span>
        <span class="mc-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="mc-card">
      <div class="mc-step"><span class="mc-n">01</span><p class="mc-mini-t">{{一步，12-26 字}}</p></div>
      <div class="mc-step"><span class="mc-n">02</span><p class="mc-mini-t">{{一步，12-26 字}}</p></div>
      <div class="mc-step"><span class="mc-n">03</span><p class="mc-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mc-kicker, h2, mt-m, grid, g2, mt-l, lede, row, mc-pill, mc-card, mc-step, mc-n, mc-mini-t, deck-footer, slide-number, notes

---

## metrics（展观数字）
指纹：chart
数量：mc-stat=3

用途：三个关键数据。焦橙大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="mc-kicker">{{数据语境，如 展览概览 · 手册数据}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="mc-stat"><div class="mc-stat-v">{{数值 ≤6 字符}}<span class="mc-stat-u">{{单位}}</span></div><div class="mc-stat-l">{{指标名，≤8 字}}</div><p class="mc-stat-note">{{口径，14-30 字}}</p></div>
    <div class="mc-stat"><div class="mc-stat-v">{{数值 ≤6 字符}}<span class="mc-stat-u">{{单位}}</span></div><div class="mc-stat-l">{{指标名，≤8 字}}</div><p class="mc-stat-note">{{口径，14-30 字}}</p></div>
    <div class="mc-stat"><div class="mc-stat-v">{{数值 ≤6 字符}}<span class="mc-stat-u">{{单位}}</span></div><div class="mc-stat-l">{{指标名，≤8 字}}</div><p class="mc-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#B45309;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mc-kicker, h2, mt-m, grid, g3, mt-l, mc-stat, mc-stat-v, mc-stat-u, mc-stat-l, mc-stat-note, deck-footer, slide-number, notes

---

## quote（设计师语录）
指纹：quote

用途：整页一句引文。衬线大字 + 出处 + 两个支撑药丸，右上有机形状自然衬底。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="mc-kicker">{{语境，如 设计师语录}}</p>
  <p class="mc-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="mc-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="mc-pill mc-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="mc-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mc-kicker, mc-quote, mt-l, mc-src, mt-m, row, mc-pill, mc-pill-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 大字章节名 + 一个过渡句 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤12 字；过渡句 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="mc-kicker">{{进度，如 第二厅 · 客厅的民主}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤12 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="mc-pill mc-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="mc-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, mc-kicker, h1, mt-m, lede, mt-l, row, mc-pill, mc-pill-accent, deck-footer, slide-number, notes

---

## moments（年代大事记）
指纹：chart
数量：mc-tl-item=4

用途：3-4 个节点的横向时间线：圆点编号 + 年份 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="mc-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="mc-tl mt-l" style="margin-top:52px">
    <div class="mc-tl-item"><div class="mc-tl-dot">01</div><div class="mc-tl-t">{{时间点，≤10 字符}}</div><p class="mc-tl-d">{{事件，12-26 字}}</p></div>
    <div class="mc-tl-item"><div class="mc-tl-dot">02</div><div class="mc-tl-t">{{时间点，≤10 字符}}</div><p class="mc-tl-d">{{事件，12-26 字}}</p></div>
    <div class="mc-tl-item"><div class="mc-tl-dot">03</div><div class="mc-tl-t">{{时间点，≤10 字符}}</div><p class="mc-tl-d">{{事件，12-26 字}}</p></div>
    <div class="mc-tl-item"><div class="mc-tl-dot">04</div><div class="mc-tl-t">{{时间点，≤10 字符}}</div><p class="mc-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#B45309;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mc-kicker, h2, mt-m, mc-tl, mt-l, mc-tl-item, mc-tl-dot, mc-tl-t, mc-tl-d, deck-footer, slide-number, notes

---

## closing（观展预约）
指纹：hero

用途：收尾页。大字标题 + 一句行动提醒 + 青绿按钮 + 描边药丸 + 年份圆章。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤12 字。

```html
<section class="slide full" data-layout="closing">
  <p class="mc-kicker">{{提醒语境，如 观展预约}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="mc-btn">{{按钮文案，≤8 字}}</span>
    <span class="mc-pill">{{次级信息，≤12 字}}</span>
    <div class="mc-badge">{{年份或代号，4-8 字符}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, mc-kicker, h1, mt-m, lede, mt-l, row, mc-btn, mc-pill, mc-badge, deck-footer, slide-number, notes
