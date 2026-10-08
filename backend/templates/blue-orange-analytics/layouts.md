# 蓝橙经营图谱 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `bo-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-blue-orange-analytics` 作用域生效，骨架里已写全，照抄结构即可。
> 页顶藏蓝数据色带与页底柱状剪影由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（数据即秩序）**：藏蓝是秩序色（标题、KPI 实底、按钮、题签方块），
> 赭红只允许出现在状态牌（bo-status）与数据口径里的「异常/目标」措辞，禁作大底色、
> 禁给标题或正文上色。等宽字（JetBrains Mono）只用于标签、编号与出处；正文与
> 大数字保持无衬线。禁具象插画、禁强阴影、禁失控大圆角。

---

## cover（看板封面）
指纹：hero

用途：开场页。等宽题签 + 大标题 + 结论细线 + 一句结论，右侧立赭红聚焦状态牌。
适用 role：cover。
内容约束：主标题 ≤12 字；lede 24-52 字；题签 ≤14 字；状态牌 ≤8 字。

```html
<section class="slide full" data-layout="cover">
  <p class="bo-kicker">{{题签，≤14 字，如 华东区 · 月度经营分析}}</p>
  <h1 class="h1 mt-m">{{主标题，≤12 字，可 <br> 分行}}</h1>
  <div class="bo-rule mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句结论，24-52 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="bo-status">{{聚焦状态，≤8 字}}</div>
    <span class="bo-pill">{{时间或场合，≤14 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年月}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, bo-kicker, h1, mt-m, bo-rule, mt-s, lede, mt-l, row, bo-status, bo-pill, deck-footer, slide-number, notes

---

## contents（结论目录）
指纹：table
数量：bo-item=4

用途：议程页。一块看板卡里放 4 行议程：等宽编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="bo-kicker">{{引导语，如 对账顺序}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="bo-card mt-l" style="margin-top:44px">
    <div class="bo-item"><span class="bo-n">01</span><span class="bo-t">{{篇名，≤8 字}}</span><span class="bo-d">{{说明，14-26 字}}</span></div>
    <div class="bo-item"><span class="bo-n">02</span><span class="bo-t">{{篇名，≤8 字}}</span><span class="bo-d">{{说明，14-26 字}}</span></div>
    <div class="bo-item"><span class="bo-n">03</span><span class="bo-t">{{篇名，≤8 字}}</span><span class="bo-d">{{说明，14-26 字}}</span></div>
    <div class="bo-item"><span class="bo-n">04</span><span class="bo-t">{{篇名，≤8 字}}</span><span class="bo-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年月}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bo-kicker, h2, mt-m, bo-card, mt-l, bo-item, bo-n, bo-t, bo-d, deck-footer, slide-number, notes

---

## keynotes（三卡信号）
指纹：cards
数量：bo-card=3

用途：恰好三张看板卡。每张：藏蓝题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤6 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="bo-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="bo-card"><span class="bo-pill bo-pill-accent">{{题签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#404040">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="bo-card"><span class="bo-pill bo-pill-accent">{{题签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#404040">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="bo-card"><span class="bo-pill bo-pill-accent">{{题签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#404040">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年月}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bo-kicker, h2, mt-m, grid, g3, mt-l, bo-card, bo-pill, bo-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：bo-step=3

用途：左边把一个结论讲透（lede + 补充 + 药丸），右边看板卡装三行动作或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="bo-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#404040">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="bo-pill">{{要点 1，≤8 字}}</span>
        <span class="bo-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="bo-card">
      <div class="bo-step"><span class="bo-n">一</span><p class="bo-mini-t">{{一步，12-26 字}}</p></div>
      <div class="bo-step"><span class="bo-n">二</span><p class="bo-mini-t">{{一步，12-26 字}}</p></div>
      <div class="bo-step"><span class="bo-n">三</span><p class="bo-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年月}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bo-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, bo-pill, bo-card, bo-step, bo-n, bo-mini-t, deck-footer, slide-number, notes

---

## metrics（KPI 三数）
指纹：chart
数量：bo-stat=3

用途：三个关键 KPI。藏蓝实底看板块 + 白色大数字是唯一主视觉，口径写进板块，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="bo-kicker">{{数据语境，如 KPI · 8 月实际}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:36px;margin-top:48px">
    <div class="bo-stat"><div class="bo-stat-v">{{数值 ≤6 字符}}<span class="bo-stat-u">{{单位}}</span></div><div class="bo-stat-l">{{指标名，≤8 字}}</div><p class="bo-stat-note">{{口径，14-30 字}}</p></div>
    <div class="bo-stat"><div class="bo-stat-v">{{数值 ≤6 字符}}<span class="bo-stat-u">{{单位}}</span></div><div class="bo-stat-l">{{指标名，≤8 字}}</div><p class="bo-stat-note">{{口径，14-30 字}}</p></div>
    <div class="bo-stat"><div class="bo-stat-v">{{数值 ≤6 字符}}<span class="bo-stat-u">{{单位}}</span></div><div class="bo-stat-l">{{指标名，≤8 字}}</div><p class="bo-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="margin-top:40px;font-size:17px;color:#778EAB;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年月}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bo-kicker, h2, mt-m, grid, g3, mt-l, bo-stat, bo-stat-v, bo-stat-u, bo-stat-l, bo-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（会议金句）
指纹：quote

用途：整页一句会议金句。大字引言 + 出处 + 两个支撑药丸，右侧留白处立虚线批注框。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤24 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="bo-kicker">{{语境，如 月度会定格}}</p>
  <p class="bo-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="bo-src mt-m">—— {{出处：人与场合，≤24 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="bo-pill bo-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="bo-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="bo-note" style="position:absolute;right:120px;top:50%;transform:translateY(-50%)">{{批注：一句读法或纪律，≤26 字}}</div>
  <div class="deck-footer"><span>{{署名 · 年月}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bo-kicker, bo-quote, mt-l, bo-src, mt-m, row, bo-pill, bo-pill-accent, bo-note, deck-footer, slide-number, notes

---

## divider（章节页）
指纹：hero

用途：章节过渡。进度题签 + 大字章节名 + 一句过渡 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤12 字；过渡句 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="bo-kicker">{{进度，如 PART 02 · 九月打法}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="bo-pill bo-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="bo-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年月}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, bo-kicker, h1, mt-m, lede, mt-l, row, bo-pill, bo-pill-accent, deck-footer, slide-number, notes

---

## moments（月度节点）
指纹：chart
数量：bo-tl-item=4

用途：4 个节点的横向时间线：等宽编号方块 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：恰好 4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="bo-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="bo-tl mt-l" style="margin-top:52px">
    <div class="bo-tl-item"><div class="bo-tl-dot">01</div><div class="bo-tl-t">{{时间点，≤10 字符}}</div><p class="bo-tl-d">{{事件，12-26 字}}</p></div>
    <div class="bo-tl-item"><div class="bo-tl-dot">02</div><div class="bo-tl-t">{{时间点，≤10 字符}}</div><p class="bo-tl-d">{{事件，12-26 字}}</p></div>
    <div class="bo-tl-item"><div class="bo-tl-dot">03</div><div class="bo-tl-t">{{时间点，≤10 字符}}</div><p class="bo-tl-d">{{事件，12-26 字}}</p></div>
    <div class="bo-tl-item"><div class="bo-tl-dot">04</div><div class="bo-tl-t">{{时间点，≤10 字符}}</div><p class="bo-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="margin-top:40px;font-size:17px;color:#778EAB;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年月}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bo-kicker, h2, mt-m, bo-tl, mt-l, bo-tl-item, bo-tl-dot, bo-tl-t, bo-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（行动收尾）
指纹：hero

用途：收尾页。大字结论 + 一句行动提醒 + 藏蓝按钮 + 状态牌。
适用 role：thanks / cta / content。
内容约束：标题 ≤12 字；lede 20-48 字；按钮 ≤6 字；药丸 ≤14 字。

```html
<section class="slide full" data-layout="closing">
  <p class="bo-kicker">{{行动语境}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="bo-btn">{{按钮文案，≤6 字}}</span>
    <span class="bo-pill">{{次级信息，≤14 字}}</span>
    <div class="bo-status">{{聚焦状态，≤8 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年月}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, bo-kicker, h1, mt-m, lede, mt-l, row, bo-btn, bo-pill, bo-status, deck-footer, slide-number, notes
