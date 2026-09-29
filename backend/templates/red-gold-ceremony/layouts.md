# 初心红典 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `rc-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-red-gold-ceremony` 作用域生效，骨架里已写全，照抄结构即可。
> 每页仪典内框与页底金辉由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（红为底，金为光）**：鎏金只允许出现在五处——金线（rc-kicker / rc-line 自带）、
> 徽记（rc-emblem）、标题（h1 金渐变 / h2 金色）、关键数据（rc-stat-v / rc-tl-t / rc-pill-accent）、
> 星点（rc-stars）。禁荧光橙、禁大面积金色色块、禁圆角大卡、禁硬阴影。深红是唯一的大面积色。
> h1/h2 衬线大字每页最多一组；rc-stars 一页至多一条，只放留白处。

---

## cover（典礼封面）
指纹：hero

用途：开场页。金签题行 + 金渐变大标题 + 对称金线 + 一句定位，可立一枚徽记。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；徽记 4 字内。

```html
<section class="slide full" data-layout="cover">
  <p class="rc-kicker">{{题签，≤14 字，如 星辰集团 · 年度盛典}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="rc-line mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="rc-emblem">{{徽记 2×2 字}}</div>
    <span class="rc-pill">{{时间或地点，≤12 字}}</span>
    <span class="rc-stars">✦ ✦ ✦</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, rc-kicker, h1, mt-m, rc-line, mt-s, lede, mt-m, row, mt-l, rc-emblem, rc-pill, rc-stars, deck-footer, slide-number, notes

---

## contents（议程册）
指纹：table
数量：rc-item=4

用途：议程页。一块仪典卡里放 4 行议程：金框序号 + 议程名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；议程名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="rc-kicker">{{引导语，如 今晚议程}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="rc-card mt-l" style="margin-top:44px">
    <div class="rc-item"><span class="rc-n">壹</span><span class="rc-t">{{议程名，≤8 字}}</span><span class="rc-d">{{说明，14-26 字}}</span></div>
    <div class="rc-item"><span class="rc-n">贰</span><span class="rc-t">{{议程名，≤8 字}}</span><span class="rc-d">{{说明，14-26 字}}</span></div>
    <div class="rc-item"><span class="rc-n">叁</span><span class="rc-t">{{议程名，≤8 字}}</span><span class="rc-d">{{说明，14-26 字}}</span></div>
    <div class="rc-item"><span class="rc-n">肆</span><span class="rc-t">{{议程名，≤8 字}}</span><span class="rc-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rc-kicker, h2, mt-m, rc-card, mt-l, rc-item, rc-n, rc-t, rc-d, deck-footer, slide-number, notes

---

## keynotes（三项并列）
指纹：cards
数量：rc-card=3

用途：恰好三张仪典卡。每张：金色题签 + 衬线小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="rc-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="rc-card"><span class="rc-pill rc-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#FFE2AC">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="rc-card"><span class="rc-pill rc-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#FFE2AC">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="rc-card"><span class="rc-pill rc-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#FFE2AC">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rc-kicker, h2, mt-m, grid, g3, mt-l, rc-card, rc-pill, rc-pill-accent, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：rc-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边仪典卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="rc-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#FFE2AC">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="rc-pill">{{要点 1，≤8 字}}</span>
        <span class="rc-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="rc-card">
      <div class="rc-step"><span class="rc-n">一</span><p class="rc-mini-t">{{一步，12-26 字}}</p></div>
      <div class="rc-step"><span class="rc-n">二</span><p class="rc-mini-t">{{一步，12-26 字}}</p></div>
      <div class="rc-step"><span class="rc-n">三</span><p class="rc-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rc-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, mt-l, rc-pill, rc-card, rc-step, rc-n, rc-mini-t, deck-footer, slide-number, notes

---

## metrics（鎏金数字）
指纹：chart
数量：rc-stat=3

用途：三个关键数据。金色大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="rc-kicker">{{数据语境，如 评选规模 · 备案}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="rc-stat"><div class="rc-stat-v">{{数值 ≤6 字符}}<span class="rc-stat-u">{{单位}}</span></div><div class="rc-stat-l">{{指标名，≤8 字}}</div><p class="rc-stat-note">{{口径，14-30 字}}</p></div>
    <div class="rc-stat"><div class="rc-stat-v">{{数值 ≤6 字符}}<span class="rc-stat-u">{{单位}}</span></div><div class="rc-stat-l">{{指标名，≤8 字}}</div><p class="rc-stat-note">{{口径，14-30 字}}</p></div>
    <div class="rc-stat"><div class="rc-stat-v">{{数值 ≤6 字符}}<span class="rc-stat-u">{{单位}}</span></div><div class="rc-stat-l">{{指标名，≤8 字}}</div><p class="rc-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#D9B98C;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rc-kicker, h2, mt-m, grid, g3, mt-l, rc-stat, rc-stat-v, rc-stat-u, rc-stat-l, rc-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（致辞引文）
指纹：quote

用途：整页一句引文。衬线大字 + 出处 + 两个支撑药丸，右侧留白处可立星点。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="rc-kicker">{{语境，如 致辞手记}}</p>
  <p class="rc-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="rc-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="rc-pill rc-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="rc-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <span class="rc-stars" style="position:absolute;right:120px;top:50%;transform:translateY(-50%)">✦ ✦ ✦</span>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rc-kicker, rc-quote, mt-l, rc-src, mt-m, row, mt-l, rc-pill, rc-pill-accent, rc-stars, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 金渐变大字章节名 + 一个过渡句 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡句 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="rc-kicker">{{进度，如 第二项 · 颁奖礼}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="rc-pill rc-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="rc-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, rc-kicker, h1, mt-m, lede, mt-l, row, rc-pill, rc-pill-accent, deck-footer, slide-number, notes

---

## moments（流程时间线）
指纹：chart
数量：rc-tl-item=4

用途：3-4 个节点的横向时间线：金框方点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="rc-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="rc-tl mt-l" style="margin-top:52px">
    <div class="rc-tl-item"><div class="rc-tl-dot">壹</div><div class="rc-tl-t">{{时间点，≤10 字符}}</div><p class="rc-tl-d">{{事件，12-26 字}}</p></div>
    <div class="rc-tl-item"><div class="rc-tl-dot">贰</div><div class="rc-tl-t">{{时间点，≤10 字符}}</div><p class="rc-tl-d">{{事件，12-26 字}}</p></div>
    <div class="rc-tl-item"><div class="rc-tl-dot">叁</div><div class="rc-tl-t">{{时间点，≤10 字符}}</div><p class="rc-tl-d">{{事件，12-26 字}}</p></div>
    <div class="rc-tl-item"><div class="rc-tl-dot">肆</div><div class="rc-tl-t">{{时间点，≤10 字符}}</div><p class="rc-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#D9B98C;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rc-kicker, h2, mt-m, rc-tl, mt-l, rc-tl-item, rc-tl-dot, rc-tl-t, rc-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（礼成收尾）
指纹：hero

用途：收尾页。金渐变大字 + 一句行动提醒 + 金框按钮 + 描边药丸，如典礼落款。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="rc-kicker">{{提醒语境，如 礼成之后}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="rc-btn">{{按钮文案，≤8 字}}</span>
    <span class="rc-pill">{{次级信息，≤10 字}}</span>
    <div class="rc-emblem">{{徽记 2 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, rc-kicker, h1, mt-m, lede, mt-l, row, rc-btn, rc-pill, rc-emblem, deck-footer, slide-number, notes
