# 梦幻分界·星河烟火 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `gc-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-gradient-cosmic` 作用域生效，骨架里已写全，照抄结构即可。
> 每页的星河、星点、流星与「分界光带」由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（金色是身份）**：深色底纪律——正文一律浅色（星白/浅紫灰），卡只有星夜卡一种
> （gc-card，金边强调用 gc-card-gold），禁浅色实底卡、禁纯黑死底。金色只允许出现在五处：
> 题签短线（gc-kicker 自带）、大数字（gc-stat-v）、金色题签（gc-pill-accent）、行动钮（gc-btn）、
> 金月（gc-moon）；禁把金色用于正文文字或大面积底色。金月只作装饰、一页至多一枚、只放留白处。
> 引文走衬线斜体（gc-quote 已内建），正文不混用衬线。

---

## cover（星幕封面）
指纹：hero

用途：开场页。题签 + 大字标题 + 金橙光线 + 一句定位，一轮金月衬在留白处。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；金月只作装饰。

```html
<section class="slide full" data-layout="cover">
  <p class="gc-kicker">{{题签，≤14 字，如 拾光社 · 星野摄影入门}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="gc-line mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="gc-moon" style="width:84px;height:84px"></div>
    <span class="gc-pill">{{时间或地点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, gc-kicker, h1, mt-m, gc-line, mt-s, lede, mt-l, row, gc-moon, gc-pill, deck-footer, slide-number, notes

---

## contents（目录）
指纹：table
数量：gc-item=4

用途：议程页。一块星夜大卡里放 4 行篇目：金圈编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="gc-kicker">{{引导语，如 课程四章}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="gc-card mt-l" style="margin-top:44px">
    <div class="gc-item"><span class="gc-n">壹</span><span class="gc-t">{{篇名，≤8 字}}</span><span class="gc-d">{{说明，14-26 字}}</span></div>
    <div class="gc-item"><span class="gc-n">贰</span><span class="gc-t">{{篇名，≤8 字}}</span><span class="gc-d">{{说明，14-26 字}}</span></div>
    <div class="gc-item"><span class="gc-n">叁</span><span class="gc-t">{{篇名，≤8 字}}</span><span class="gc-d">{{说明，14-26 字}}</span></div>
    <div class="gc-item"><span class="gc-n">肆</span><span class="gc-t">{{篇名，≤8 字}}</span><span class="gc-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gc-kicker, h2, mt-m, gc-card, mt-l, gc-item, gc-n, gc-t, gc-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：gc-card=3

用途：恰好三张星夜卡。每张：金色题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="gc-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="gc-card"><span class="gc-pill gc-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#AAB4E8">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="gc-card"><span class="gc-pill gc-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#AAB4E8">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="gc-card"><span class="gc-pill gc-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#AAB4E8">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gc-kicker, h2, mt-m, grid, g3, mt-l, gc-card, gc-pill, gc-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：gc-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边金边卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="gc-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#AAB4E8">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="gc-pill">{{要点 1，≤8 字}}</span>
        <span class="gc-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="gc-card gc-card-gold">
      <div class="gc-step"><span class="gc-n">壹</span><p class="gc-mini-t">{{一步，12-26 字}}</p></div>
      <div class="gc-step"><span class="gc-n">贰</span><p class="gc-mini-t">{{一步，12-26 字}}</p></div>
      <div class="gc-step"><span class="gc-n">叁</span><p class="gc-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gc-kicker, h2, mt-m, grid, g2, mt-l, lede, row, gc-pill, gc-card, gc-card-gold, gc-step, gc-n, gc-mini-t, deck-footer, slide-number, notes

---

## metrics（星光数据）
指纹：chart
数量：gc-stat=3

用途：三个关键数据。金橙渐变大数字是唯一主视觉，口径写进块内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="gc-kicker">{{数据语境，如 参数三件套}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="gc-stat"><div class="row" style="align-items:baseline;gap:10px"><span class="gc-stat-v">{{数值 ≤6 字符}}</span><span class="gc-stat-u">{{单位}}</span></div><div class="gc-stat-l mt-s">{{指标名，≤8 字}}</div><p class="gc-stat-note mt-s">{{口径，14-30 字}}</p></div>
    <div class="gc-stat"><div class="row" style="align-items:baseline;gap:10px"><span class="gc-stat-v">{{数值 ≤6 字符}}</span><span class="gc-stat-u">{{单位}}</span></div><div class="gc-stat-l mt-s">{{指标名，≤8 字}}</div><p class="gc-stat-note mt-s">{{口径，14-30 字}}</p></div>
    <div class="gc-stat"><div class="row" style="align-items:baseline;gap:10px"><span class="gc-stat-v">{{数值 ≤6 字符}}</span><span class="gc-stat-u">{{单位}}</span></div><div class="gc-stat-l mt-s">{{指标名，≤8 字}}</div><p class="gc-stat-note mt-s">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#96A0D6;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gc-kicker, h2, mt-m, grid, g3, mt-l, gc-stat, row, gc-stat-v, gc-stat-u, gc-stat-l, gc-stat-note, mt-s, mt-m, deck-footer, slide-number, notes

---

## quote（题跋引文）
指纹：quote

用途：整页一句引文。衬线斜体金米大字 + 出处 + 两个支撑药丸，一轮金月衬在留白处。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="gc-kicker">{{语境，如 社长开课语}}</p>
  <p class="gc-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="gc-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="gc-pill gc-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="gc-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="gc-moon" style="position:absolute;right:140px;top:120px"></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gc-kicker, gc-quote, mt-l, gc-src, mt-m, row, gc-pill, gc-pill-accent, gc-moon, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="gc-kicker">{{进度，如 第二章 · 核心参数}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="gc-pill gc-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="gc-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, gc-kicker, h1, mt-m, lede, mt-l, row, gc-pill, gc-pill-accent, deck-footer, slide-number, notes

---

## moments（流程时间线）
指纹：chart
数量：gc-tl-item=4

用途：3-4 个节点的横向时间线：星点圆号 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="gc-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="gc-tl mt-l" style="margin-top:52px">
    <div class="gc-tl-item"><div class="gc-tl-dot">壹</div><div class="gc-tl-t">{{时间点，≤10 字符}}</div><p class="gc-tl-d">{{事件，12-26 字}}</p></div>
    <div class="gc-tl-item"><div class="gc-tl-dot">贰</div><div class="gc-tl-t">{{时间点，≤10 字符}}</div><p class="gc-tl-d">{{事件，12-26 字}}</p></div>
    <div class="gc-tl-item"><div class="gc-tl-dot">叁</div><div class="gc-tl-t">{{时间点，≤10 字符}}</div><p class="gc-tl-d">{{事件，12-26 字}}</p></div>
    <div class="gc-tl-item"><div class="gc-tl-dot">肆</div><div class="gc-tl-t">{{时间点，≤10 字符}}</div><p class="gc-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#96A0D6;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gc-kicker, h2, mt-m, gc-tl, mt-l, gc-tl-item, gc-tl-dot, gc-tl-t, gc-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（收尾行动）
指纹：hero

用途：收尾页。大字 + 一句行动提醒 + 金火按钮 + 金月，如夜幕落款。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="gc-kicker">{{行动语境，如 出队报名}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="gc-btn">{{按钮文案，≤8 字}}</span>
    <span class="gc-pill">{{次级信息，≤10 字}}</span>
    <div class="gc-moon" style="width:64px;height:64px"></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, gc-kicker, h1, mt-m, lede, mt-l, row, gc-btn, gc-pill, gc-moon, deck-footer, slide-number, notes
