# 磨砂玻璃 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `gl-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-glassmorphism` 作用域生效，骨架里已写全，照抄结构即可。
> 每页的彩色光斑与顶部光带由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（玻璃是身份）**：卡只有一种——rgba 白玻璃 + backdrop-blur + 细白描边（gl-glass，
> 长文本用加实档 gl-glass-deep），禁实色卡、禁硬边框、禁把卡片做成不透明。柔紫与青蓝只允许
> 出现在题签短线、大数字、时间点、强调药丸与按钮；正文永远是深墨与石板灰。h1/h2 每页最多
> 一组；渐变文字只用于 gl-quote 与 gl-stat-v（两端色已够深，不要再自行改浅）。

---

## cover（启幕封面）
指纹：hero

用途：开场页。题签 + 大字标题 + 品牌渐变线 + 一句定位，右侧可立一枚玻璃徽章。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；徽章 6 字内。

```html
<section class="slide full" data-layout="cover">
  <p class="gl-kicker">{{题签，≤14 字，如 云雾 UI · v2.0 发布}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="gl-line mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="gl-badge">{{徽章 2×2 字}}</div>
    <span class="gl-chip">{{时间或场合，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, gl-kicker, h1, mt-m, gl-line, mt-s, lede, mt-l, row, gl-badge, gl-chip, deck-footer, slide-number, notes

---

## contents（目录）
指纹：table
数量：gl-item=4

用途：议程页。一块磨砂大卡里放 4 行篇目：玻璃圆号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="gl-kicker">{{引导语，如 发布说明目录}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="gl-glass gl-glass-deep mt-l" style="margin-top:44px">
    <div class="gl-item"><span class="gl-n">01</span><span class="gl-t">{{篇名，≤8 字}}</span><span class="gl-d">{{说明，14-26 字}}</span></div>
    <div class="gl-item"><span class="gl-n">02</span><span class="gl-t">{{篇名，≤8 字}}</span><span class="gl-d">{{说明，14-26 字}}</span></div>
    <div class="gl-item"><span class="gl-n">03</span><span class="gl-t">{{篇名，≤8 字}}</span><span class="gl-d">{{说明，14-26 字}}</span></div>
    <div class="gl-item"><span class="gl-n">04</span><span class="gl-t">{{篇名，≤8 字}}</span><span class="gl-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gl-kicker, h2, mt-m, gl-glass, gl-glass-deep, mt-l, gl-item, gl-n, gl-t, gl-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：gl-glass=3

用途：恰好三张玻璃卡。每张：强调题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="gl-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="gl-glass"><span class="gl-chip gl-chip-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#475569">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="gl-glass"><span class="gl-chip gl-chip-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#475569">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="gl-glass"><span class="gl-chip gl-chip-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#475569">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gl-kicker, h2, mt-m, grid, g3, mt-l, gl-glass, gl-chip, gl-chip-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：gl-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边玻璃卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="gl-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#475569">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="gl-chip">{{要点 1，≤8 字}}</span>
        <span class="gl-chip">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="gl-glass gl-glass-deep">
      <div class="gl-step"><span class="gl-n">01</span><p class="gl-mini-t">{{一步，12-26 字}}</p></div>
      <div class="gl-step"><span class="gl-n">02</span><p class="gl-mini-t">{{一步，12-26 字}}</p></div>
      <div class="gl-step"><span class="gl-n">03</span><p class="gl-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gl-kicker, h2, mt-m, grid, g2, mt-l, lede, row, gl-chip, gl-glass, gl-glass-deep, gl-step, gl-n, gl-mini-t, deck-footer, slide-number, notes

---

## metrics（光白数字）
指纹：chart
数量：gl-stat=3

用途：三个关键数据。玻璃卡内的渐变大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="gl-kicker">{{数据语境，如 发布说明 · 关键数字}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="gl-stat"><div class="row" style="align-items:baseline;gap:10px"><span class="gl-stat-v">{{数值 ≤6 字符}}</span><span class="gl-stat-u">{{单位}}</span></div><div class="gl-stat-l mt-s">{{指标名，≤8 字}}</div><p class="gl-stat-note mt-s">{{口径，14-30 字}}</p></div>
    <div class="gl-stat"><div class="row" style="align-items:baseline;gap:10px"><span class="gl-stat-v">{{数值 ≤6 字符}}</span><span class="gl-stat-u">{{单位}}</span></div><div class="gl-stat-l mt-s">{{指标名，≤8 字}}</div><p class="gl-stat-note mt-s">{{口径，14-30 字}}</p></div>
    <div class="gl-stat"><div class="row" style="align-items:baseline;gap:10px"><span class="gl-stat-v">{{数值 ≤6 字符}}</span><span class="gl-stat-u">{{单位}}</span></div><div class="gl-stat-l mt-s">{{指标名，≤8 字}}</div><p class="gl-stat-note mt-s">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#64748B;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gl-kicker, h2, mt-m, grid, g3, mt-l, gl-stat, row, gl-stat-v, gl-stat-u, gl-stat-l, gl-stat-note, mt-s, mt-m, deck-footer, slide-number, notes

---

## quote（引文）
指纹：quote

用途：整页一句引文。渐变大字 + 出处 + 两个支撑药丸，光斑在字后流转。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="gl-kicker">{{语境，如 设计团队卷首语}}</p>
  <p class="gl-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="gl-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="gl-chip gl-chip-accent">{{支撑点 1，≤8 字}}</span>
    <span class="gl-chip">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gl-kicker, gl-quote, mt-l, gl-src, mt-m, row, gl-chip, gl-chip-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="gl-kicker">{{进度，如 第二章 · 玻璃体系}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="gl-chip gl-chip-accent">{{看点 1，≤8 字}}</span>
    <span class="gl-chip">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, gl-kicker, h1, mt-m, lede, mt-l, row, gl-chip, gl-chip-accent, deck-footer, slide-number, notes

---

## moments（流程时间线）
指纹：chart
数量：gl-tl-item=4

用途：3-4 个节点的横向时间线：玻璃圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="gl-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="gl-tl mt-l" style="margin-top:52px">
    <div class="gl-tl-item"><div class="gl-tl-dot">01</div><div class="gl-tl-t">{{时间点，≤10 字符}}</div><p class="gl-tl-d">{{事件，12-26 字}}</p></div>
    <div class="gl-tl-item"><div class="gl-tl-dot">02</div><div class="gl-tl-t">{{时间点，≤10 字符}}</div><p class="gl-tl-d">{{事件，12-26 字}}</p></div>
    <div class="gl-tl-item"><div class="gl-tl-dot">03</div><div class="gl-tl-t">{{时间点，≤10 字符}}</div><p class="gl-tl-d">{{事件，12-26 字}}</p></div>
    <div class="gl-tl-item"><div class="gl-tl-dot">04</div><div class="gl-tl-t">{{时间点，≤10 字符}}</div><p class="gl-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#64748B;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gl-kicker, h2, mt-m, gl-tl, mt-l, gl-tl-item, gl-tl-dot, gl-tl-t, gl-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（收尾行动）
指纹：hero

用途：收尾页。大字 + 一句行动提醒 + 渐变按钮 + 玻璃徽章，如发布会谢幕。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="gl-kicker">{{行动语境，如 获取方式}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="gl-btn">{{按钮文案，≤8 字}}</span>
    <span class="gl-chip">{{次级信息，≤10 字}}</span>
    <div class="gl-badge">{{徽章 2×2 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, gl-kicker, h1, mt-m, lede, mt-l, row, gl-btn, gl-chip, gl-badge, deck-footer, slide-number, notes
