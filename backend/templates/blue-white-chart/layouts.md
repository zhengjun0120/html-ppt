# 蓝白商务图表 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `bw-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-blue-white-chart` 作用域生效，骨架里已写全，照抄结构即可。
> 右上浅蓝实心圆与左下粗描边圆环由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（深蓝是骨架）**：企业深蓝 #305598 只用于标题（h1/h2/h4 已内建）、分隔线
> （bw-rule）、编号块（bw-n / bw-tl-dot）、KPI 顶线（bw-stat）与按钮（bw-btn）；
> 深海军蓝 #00329D 只给大数字（bw-stat-v）。红绿橙只可在口径文字里表达数据升降语义，
> 禁作装饰。禁大面积渐变、禁手绘元素、禁破坏网格对齐——所有模块都要像咨询报告一样精确。

---

## cover（图表封面）
指纹：hero

用途：开场页。胶囊题签 + 深蓝大标题 + 横贯分隔线 + 一句结论，左下立标识块。
适用 role：cover。
内容约束：主标题 ≤12 字；lede 24-52 字；题签 ≤14 字；标识 2-4 字。

```html
<section class="slide full" data-layout="cover">
  <p class="bw-kicker">{{题签，≤14 字，如 Q3 REVIEW · 产品线}}</p>
  <h1 class="h1 mt-m">{{主标题，≤12 字，可 <br> 分行}}</h1>
  <div class="bw-rule mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句结论，24-52 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="bw-badge">{{标识 2-4 字}}</div>
    <span class="bw-pill">{{时间或场合，≤14 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 季度}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, bw-kicker, h1, mt-m, bw-rule, mt-s, lede, mt-l, row, bw-badge, bw-pill, deck-footer, slide-number, notes

---

## contents（汇报目录）
指纹：table
数量：bw-item=4

用途：议程页。一块浅蓝信息卡里放 4 行议程：深蓝编号块 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="bw-kicker">{{引导语，如 AGENDA}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="bw-card mt-l" style="margin-top:44px">
    <div class="bw-item"><span class="bw-n">01</span><span class="bw-t">{{篇名，≤8 字}}</span><span class="bw-d">{{说明，14-26 字}}</span></div>
    <div class="bw-item"><span class="bw-n">02</span><span class="bw-t">{{篇名，≤8 字}}</span><span class="bw-d">{{说明，14-26 字}}</span></div>
    <div class="bw-item"><span class="bw-n">03</span><span class="bw-t">{{篇名，≤8 字}}</span><span class="bw-d">{{说明，14-26 字}}</span></div>
    <div class="bw-item"><span class="bw-n">04</span><span class="bw-t">{{篇名，≤8 字}}</span><span class="bw-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 季度}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bw-kicker, h2, mt-m, bw-card, mt-l, bw-item, bw-n, bw-t, bw-d, deck-footer, slide-number, notes

---

## keynotes（三卡引擎）
指纹：cards
数量：bw-card=3

用途：恰好三张浅蓝卡。每张：深蓝题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤6 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="bw-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="bw-card"><span class="bw-pill bw-pill-accent">{{题签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#404040">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="bw-card"><span class="bw-pill bw-pill-accent">{{题签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#404040">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="bw-card"><span class="bw-pill bw-pill-accent">{{题签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#404040">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 季度}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bw-kicker, h2, mt-m, grid, g3, mt-l, bw-card, bw-pill, bw-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：bw-step=3

用途：左边把一个问题讲透（lede + 补充 + 药丸），右边浅蓝卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="bw-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#404040">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="bw-pill">{{要点 1，≤8 字}}</span>
        <span class="bw-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="bw-card">
      <div class="bw-step"><span class="bw-n">一</span><p class="bw-mini-t">{{一步，12-26 字}}</p></div>
      <div class="bw-step"><span class="bw-n">二</span><p class="bw-mini-t">{{一步，12-26 字}}</p></div>
      <div class="bw-step"><span class="bw-n">三</span><p class="bw-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 季度}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bw-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, bw-pill, bw-card, bw-step, bw-n, bw-mini-t, deck-footer, slide-number, notes

---

## metrics（核心 KPI）
指纹：chart
数量：bw-stat=3

用途：三个关键 KPI。深蓝顶线白卡 + 深海军蓝大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="bw-kicker">{{数据语境，如 KPI · 2026 Q3}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:36px;margin-top:48px">
    <div class="bw-stat"><div class="bw-stat-v">{{数值 ≤6 字符}}<span class="bw-stat-u">{{单位}}</span></div><div class="bw-stat-l">{{指标名，≤8 字}}</div><p class="bw-stat-note">{{口径，14-30 字}}</p></div>
    <div class="bw-stat"><div class="bw-stat-v">{{数值 ≤6 字符}}<span class="bw-stat-u">{{单位}}</span></div><div class="bw-stat-l">{{指标名，≤8 字}}</div><p class="bw-stat-note">{{口径，14-30 字}}</p></div>
    <div class="bw-stat"><div class="bw-stat-v">{{数值 ≤6 字符}}<span class="bw-stat-u">{{单位}}</span></div><div class="bw-stat-l">{{指标名，≤8 字}}</div><p class="bw-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="margin-top:40px;font-size:17px;color:#8C9BB5;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 季度}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bw-kicker, h2, mt-m, grid, g3, mt-l, bw-stat, bw-stat-v, bw-stat-u, bw-stat-l, bw-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（季度金句）
指纹：quote

用途：整页一句季度金句。深蓝大字引言 + 出处 + 两个支撑药丸，右侧留白处衬圆环。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤24 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="bw-kicker">{{语境，如 季度会金句}}</p>
  <p class="bw-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="bw-src mt-m">—— {{出处：人与场合，≤24 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="bw-pill bw-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="bw-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="bw-ring" style="position:absolute;right:130px;top:50%;transform:translateY(-50%)"></div>
  <div class="deck-footer"><span>{{署名 · 季度}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bw-kicker, bw-quote, mt-l, bw-src, mt-m, row, bw-pill, bw-pill-accent, bw-ring, deck-footer, slide-number, notes

---

## divider（章节页）
指纹：hero

用途：章节过渡。进度题签 + 深蓝大字章节名 + 一句过渡 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤12 字；过渡句 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="bw-kicker">{{进度，如 PART 02 · Q4 打法}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="bw-pill bw-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="bw-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 季度}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, bw-kicker, h1, mt-m, lede, mt-l, row, bw-pill, bw-pill-accent, deck-footer, slide-number, notes

---

## moments（季度大事记）
指纹：chart
数量：bw-tl-item=4

用途：4 个节点的横向时间线：深蓝编号方块 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：恰好 4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="bw-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="bw-tl mt-l" style="margin-top:52px">
    <div class="bw-tl-item"><div class="bw-tl-dot">01</div><div class="bw-tl-t">{{时间点，≤10 字符}}</div><p class="bw-tl-d">{{事件，12-26 字}}</p></div>
    <div class="bw-tl-item"><div class="bw-tl-dot">02</div><div class="bw-tl-t">{{时间点，≤10 字符}}</div><p class="bw-tl-d">{{事件，12-26 字}}</p></div>
    <div class="bw-tl-item"><div class="bw-tl-dot">03</div><div class="bw-tl-t">{{时间点，≤10 字符}}</div><p class="bw-tl-d">{{事件，12-26 字}}</p></div>
    <div class="bw-tl-item"><div class="bw-tl-dot">04</div><div class="bw-tl-t">{{时间点，≤10 字符}}</div><p class="bw-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="margin-top:40px;font-size:17px;color:#8C9BB5;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 季度}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bw-kicker, h2, mt-m, bw-tl, mt-l, bw-tl-item, bw-tl-dot, bw-tl-t, bw-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（收尾行动）
指纹：hero

用途：收尾页。深蓝大字结论 + 一句行动提醒 + 深蓝按钮 + 描边药丸。
适用 role：thanks / cta / content。
内容约束：标题 ≤12 字；lede 20-48 字；按钮 ≤6 字；药丸 ≤14 字。

```html
<section class="slide full" data-layout="closing">
  <p class="bw-kicker">{{行动语境}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="bw-btn">{{按钮文案，≤6 字}}</span>
    <span class="bw-pill">{{次级信息，≤14 字}}</span>
    <div class="bw-badge">{{标识 2-4 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 季度}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, bw-kicker, h1, mt-m, lede, mt-l, row, bw-btn, bw-pill, bw-badge, deck-footer, slide-number, notes
