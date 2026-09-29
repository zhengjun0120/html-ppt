# 童趣橙暖 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `co-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-children-warm-orange` 作用域生效，骨架里已写全，照抄结构即可。
> 每页顶部的太阳云朵星星与页底草绿缓坡由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（橙是身份，糖果是点缀）**：橙色只允许出现在四处——题签方块（co-kicker 自带）、
> 关键数字（co-stat-v / co-tl-t）、强调药丸（co-pill-accent）、行动钮（co-btn）。
> 暖黄只落在贴纸（co-bubble）、糖果圆号（co-n）与时间线圆点；粉/蓝紫只出现在背景装饰里，
> 不进正文与卡片。禁荧光色大底、禁满屏卡通贴纸、禁正文使用草绿——热闹靠装饰层，版面保持清爽。

---

## cover（暖阳封面）
指纹：hero

用途：开场页。题签 + 圆润大字标题 + 波浪线 + 一句定位，右下可贴一枚暖黄贴纸。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 26-52 字；题签 ≤14 字；贴纸 ≤12 字。

```html
<section class="slide full" data-layout="cover">
  <p class="co-kicker">{{题签，≤14 字，如 萌芽编程社 · 秋令营}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="co-squiggle mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，26-52 字}}</p>
  <div class="row mt-l" style="gap:22px">
    <span class="co-bubble">{{贴纸副题，≤12 字}}</span>
    <span class="co-pill">{{时间或人群，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, co-kicker, h1, mt-m, co-squiggle, mt-s, lede, mt-m, mt-l, row, co-bubble, co-pill, deck-footer, slide-number, notes

---

## contents（糖果目录）
指纹：table
数量：co-item=4

用途：议程页。一块大圆角卡里放 4 行篇目：糖果圆号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="co-kicker">{{引导语，如 今日旅程}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="co-card mt-l" style="margin-top:44px">
    <div class="co-item"><span class="co-n">01</span><span class="co-t">{{篇名，≤8 字}}</span><span class="co-d">{{说明，14-26 字}}</span></div>
    <div class="co-item"><span class="co-n">02</span><span class="co-t">{{篇名，≤8 字}}</span><span class="co-d">{{说明，14-26 字}}</span></div>
    <div class="co-item"><span class="co-n">03</span><span class="co-t">{{篇名，≤8 字}}</span><span class="co-d">{{说明，14-26 字}}</span></div>
    <div class="co-item"><span class="co-n">04</span><span class="co-t">{{篇名，≤8 字}}</span><span class="co-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, co-kicker, h2, mt-m, co-card, mt-l, co-item, co-n, co-t, co-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：co-card=3

用途：恰好三张大圆角卡。每张：橙描边题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤10 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="co-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="co-card"><span class="co-pill co-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤10 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#7C5B36">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="co-card"><span class="co-pill co-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤10 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#7C5B36">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="co-card"><span class="co-pill co-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤10 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#7C5B36">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, co-kicker, h2, mt-m, grid, g3, mt-l, co-card, co-pill, co-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：co-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边大圆角卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="co-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#7C5B36">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="co-pill">{{要点 1，≤10 字}}</span>
        <span class="co-pill">{{要点 2，≤10 字}}</span>
      </div>
    </div>
    <div class="co-card">
      <div class="co-step"><span class="co-n">01</span><p class="co-mini-t">{{一步，12-26 字}}</p></div>
      <div class="co-step"><span class="co-n">02</span><p class="co-mini-t">{{一步，12-26 字}}</p></div>
      <div class="co-step"><span class="co-n">03</span><p class="co-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, co-kicker, h2, mt-m, grid, g2, mt-l, lede, row, co-pill, co-card, co-step, co-n, co-mini-t, deck-footer, slide-number, notes

---

## metrics（橙子数字）
指纹：chart
数量：co-stat=3

用途：三个关键数据。橙色大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="co-kicker">{{数据语境，如 往期营期 · 复盘}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="co-stat"><div class="co-stat-v">{{数值 ≤6 字符}}<span class="co-stat-u">{{单位}}</span></div><div class="co-stat-l">{{指标名，≤8 字}}</div><p class="co-stat-note">{{口径，14-30 字}}</p></div>
    <div class="co-stat"><div class="co-stat-v">{{数值 ≤6 字符}}<span class="co-stat-u">{{单位}}</span></div><div class="co-stat-l">{{指标名，≤8 字}}</div><p class="co-stat-note">{{口径，14-30 字}}</p></div>
    <div class="co-stat"><div class="co-stat-v">{{数值 ≤6 字符}}<span class="co-stat-u">{{单位}}</span></div><div class="co-stat-l">{{指标名，≤8 字}}</div><p class="co-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#A68757;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, co-kicker, h2, mt-m, grid, g3, mt-l, co-stat, co-stat-v, co-stat-u, co-stat-l, co-stat-note, deck-footer, slide-number, notes

---

## quote（便签引言）
指纹：quote

用途：整页一句引文。圆润大字 + 出处 + 两个支撑药丸，像贴在墙上的一张便签。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤10 字。

```html
<section class="slide" data-layout="quote">
  <p class="co-kicker">{{语境，如 家长反馈}}</p>
  <p class="co-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="co-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="co-pill co-pill-accent">{{支撑点 1，≤10 字}}</span>
    <span class="co-pill">{{支撑点 2，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, co-kicker, co-quote, mt-l, co-src, mt-m, row, co-pill, co-pill-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="co-kicker">{{进度，如 第二章 · 课程三阶}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="co-pill co-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="co-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, co-kicker, h1, mt-m, lede, mt-l, row, co-pill, co-pill-accent, deck-footer, slide-number, notes

---

## moments（流程时间线）
指纹：chart
数量：co-tl-item=4

用途：3-4 个节点的横向时间线：糖果圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="co-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="co-tl mt-l" style="margin-top:52px">
    <div class="co-tl-item"><div class="co-tl-dot">01</div><div class="co-tl-t">{{时间点，≤10 字符}}</div><p class="co-tl-d">{{事件，12-26 字}}</p></div>
    <div class="co-tl-item"><div class="co-tl-dot">02</div><div class="co-tl-t">{{时间点，≤10 字符}}</div><p class="co-tl-d">{{事件，12-26 字}}</p></div>
    <div class="co-tl-item"><div class="co-tl-dot">03</div><div class="co-tl-t">{{时间点，≤10 字符}}</div><p class="co-tl-d">{{事件，12-26 字}}</p></div>
    <div class="co-tl-item"><div class="co-tl-dot">04</div><div class="co-tl-t">{{时间点，≤10 字符}}</div><p class="co-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#A68757;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, co-kicker, h2, mt-m, co-tl, mt-l, co-tl-item, co-tl-dot, co-tl-t, co-tl-d, deck-footer, slide-number, notes

---

## closing（收尾行动）
指纹：hero

用途：收尾页。圆润大字 + 一句行动提醒 + 橙色圆钮 + 贴纸药丸，像海报的报名角。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤12 字。

```html
<section class="slide full" data-layout="closing">
  <p class="co-kicker">{{提醒语境，如 秋令营报名}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:22px">
    <span class="co-btn">{{按钮文案，≤8 字}}</span>
    <span class="co-pill">{{次级信息，≤12 字}}</span>
    <span class="co-bubble">{{落款贴纸，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, co-kicker, h1, mt-m, lede, mt-l, row, co-btn, co-pill, co-bubble, deck-footer, slide-number, notes
