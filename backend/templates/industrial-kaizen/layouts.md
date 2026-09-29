# 表格 · 现代周报 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `ik-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-industrial-kaizen` 作用域生效，骨架里已写全，照抄结构即可。
> 每页顶部的账本封边与右侧栏线由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（对齐是信仰，零彩色）**：全模板禁彩色——状态一律用灰阶徽章
> （ik-badge-solid 实底反白 / ik-badge-soft 浅底 / ik-badge-line 描边）表达；
> 深墨强线只许出现在标题下沿（ik-rule）、KPI 横栏封边与按钮实底；
> 所有数字用等宽数字（模板已内建 tabular-nums），多列数字要能连成一条竖线。
> 粗壮大标题（h1/h2）每页最多一组；一页最多两组内容块 + 页脚。

---

## cover（周报封面）
指纹：hero

用途：开场页。大写眉题 + 粗壮标题 + 封边强线 + 一句导语，周期与编制用徽章收在标题下。
适用 role：cover。
内容约束：主标题 ≤14 字可两行；lede 30-60 字；眉题 ≤24 字符；徽章各 ≤12 字。

```html
<section class="slide full" data-layout="cover">
  <p class="ik-kicker">{{眉题，≤24 字符，如 WEEKLY KAIZEN REPORT · W37}}</p>
  <h1 class="h1 mt-m">{{主标题，≤14 字，可 <br> 分两行}}</h1>
  <div class="ik-rule mt-s"></div>
  <p class="lede mt-m" style="max-width:52ch">{{一句导语，30-60 字}}</p>
  <div class="row mt-l" style="gap:14px">
    <span class="ik-badge-solid">{{周期或编号，≤12 字}}</span>
    <span class="ik-badge-soft">{{编制或署名，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, ik-kicker, h1, mt-m, ik-rule, mt-s, lede, mt-l, row, ik-badge-solid, ik-badge-soft, deck-footer, slide-number, notes

---

## contents（本期目录）
指纹：table
数量：ik-item=4

用途：议程页。4 行条目：两位数编号 + 条目名 + 一句说明，发丝线分行如报表行。
适用 role：toc。
内容约束：恰好 4 行；条目名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="ik-kicker">{{引导语，如 本期目录}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="ik-card mt-l" style="margin-top:44px">
    <div class="ik-item"><span class="ik-n">01</span><span class="ik-t">{{条目名，≤8 字}}</span><span class="ik-d">{{说明，14-26 字}}</span></div>
    <div class="ik-item"><span class="ik-n">02</span><span class="ik-t">{{条目名，≤8 字}}</span><span class="ik-d">{{说明，14-26 字}}</span></div>
    <div class="ik-item"><span class="ik-n">03</span><span class="ik-t">{{条目名，≤8 字}}</span><span class="ik-d">{{说明，14-26 字}}</span></div>
    <div class="ik-item"><span class="ik-n">04</span><span class="ik-t">{{条目名，≤8 字}}</span><span class="ik-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ik-kicker, h2, mt-m, ik-card, mt-l, ik-item, ik-n, ik-t, ik-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：ik-card=3

用途：恰好三张细边卡。每张：状态徽章题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤6 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="ik-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:48px">
    <div class="ik-card"><span class="ik-badge ik-badge-solid">{{题签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:17px;line-height:1.75;color:#2A2924">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="ik-card"><span class="ik-badge ik-badge-solid">{{题签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:17px;line-height:1.75;color:#2A2924">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="ik-card"><span class="ik-badge ik-badge-soft">{{题签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:17px;line-height:1.75;color:#2A2924">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ik-kicker, h2, mt-m, grid, g3, mt-l, ik-card, ik-badge, ik-badge-solid, ik-badge-soft, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：ik-step=3

用途：左边把一件事讲透（lede + 补充 + 徽章），右边细边卡装三行日程或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个徽章；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="ik-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:18px;line-height:1.75;color:#2A2924">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:14px">
        <span class="ik-badge ik-badge-line">{{要点 1，≤8 字}}</span>
        <span class="ik-badge ik-badge-soft">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="ik-card">
      <div class="ik-step"><span class="ik-n">01</span><p class="ik-mini-t">{{一步，12-26 字}}</p></div>
      <div class="ik-step"><span class="ik-n">02</span><p class="ik-mini-t">{{一步，12-26 字}}</p></div>
      <div class="ik-step"><span class="ik-n">03</span><p class="ik-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ik-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, ik-badge, ik-badge-line, ik-badge-soft, ik-card, ik-step, ik-n, ik-mini-t, deck-footer, slide-number, notes

---

## metrics（KPI 横栏）
指纹：chart
数量：ik-stat=3

用途：三个关键数据。KPI 横栏（上下细线封边、竖细线分列）是主视觉，▲▼ 趋势符只在这里用一次；口径写进栏内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="ik-kicker">{{数据语境，如 第一节 · 指标总览}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:0;margin-top:48px;border-top:1px solid #EDE9DC;border-bottom:1px solid #EDE9DC;padding:28px 0">
    <div class="ik-stat"><div class="ik-stat-v">{{数值 ≤6 字符}}<span class="ik-stat-u">{{单位}}</span></div><div class="ik-stat-l">{{指标名，≤8 字}}</div><p class="ik-stat-note">{{▲▼ 趋势 + 口径，14-30 字}}</p></div>
    <div class="ik-stat"><div class="ik-stat-v">{{数值 ≤6 字符}}<span class="ik-stat-u">{{单位}}</span></div><div class="ik-stat-l">{{指标名，≤8 字}}</div><p class="ik-stat-note">{{▲▼ 趋势 + 口径，14-30 字}}</p></div>
    <div class="ik-stat"><div class="ik-stat-v">{{数值 ≤6 字符}}<span class="ik-stat-u">{{单位}}</span></div><div class="ik-stat-l">{{指标名，≤8 字}}</div><p class="ik-stat-note">{{▲▼ 趋势 + 口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:16px;color:#8B8478;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ik-kicker, h2, mt-m, grid, g3, mt-l, ik-stat, ik-stat-v, ik-stat-u, ik-stat-l, ik-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（摘录引文）
指纹：quote

用途：整页一句引文。粗壮标语体 + 出处 + 两个徽章，像报表卷首的扉页语。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；徽章各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="ik-kicker">{{语境，如 本周一句}}</p>
  <p class="ik-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="ik-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:14px">
    <span class="ik-badge ik-badge-line">{{支撑点 1，≤8 字}}</span>
    <span class="ik-badge ik-badge-soft">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ik-kicker, ik-quote, mt-l, ik-src, mt-m, row, ik-badge, ik-badge-line, ik-badge-soft, deck-footer, slide-number, notes

---

## divider（章节页）
指纹：hero

用途：章节过渡。大写眉题 + 粗壮章节名 + 封边强线 + 一个过渡问题 + 两个徽章。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-52 字；徽章各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="ik-kicker">{{进度，如 第二节 · 改善落地}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <div class="ik-rule mt-s"></div>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:14px">
    <span class="ik-badge ik-badge-line">{{看点 1，≤8 字}}</span>
    <span class="ik-badge ik-badge-soft">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, ik-kicker, h1, mt-m, ik-rule, mt-s, lede, mt-l, row, ik-badge, ik-badge-line, ik-badge-soft, deck-footer, slide-number, notes

---

## moments（周历时间线）
指纹：chart
数量：ik-tl-item=4

用途：3-4 个节点的横向时间线：方点 + 时间点 + 一句事件，细线串联如排程表。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="ik-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="ik-tl mt-l" style="margin-top:56px">
    <div class="ik-tl-item"><div class="ik-tl-dot"></div><div class="ik-tl-t">{{时间点，≤10 字符}}</div><p class="ik-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ik-tl-item"><div class="ik-tl-dot"></div><div class="ik-tl-t">{{时间点，≤10 字符}}</div><p class="ik-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ik-tl-item"><div class="ik-tl-dot"></div><div class="ik-tl-t">{{时间点，≤10 字符}}</div><p class="ik-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ik-tl-item"><div class="ik-tl-dot"></div><div class="ik-tl-t">{{时间点，≤10 字符}}</div><p class="ik-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:16px;color:#8B8478;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ik-kicker, h2, mt-m, ik-tl, mt-l, ik-tl-item, ik-tl-dot, ik-tl-t, ik-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（收尾行）
指纹：hero

用途：收尾页。粗壮大字 + 一句行动提醒 + 深墨实底按钮 + 徽章，如报表末行的呈报栏。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤6 字；徽章 ≤12 字。

```html
<section class="slide full" data-layout="closing">
  <p class="ik-kicker">{{行动语境，如 第四节 · 风险求助}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句行动提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="ik-btn">{{按钮文案，≤6 字}}</span>
    <span class="ik-badge ik-badge-soft">{{次级信息，≤12 字}}</span>
    <span class="ik-badge ik-badge-line">{{次级信息，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, ik-kicker, h1, mt-m, lede, mt-l, row, ik-btn, ik-badge, ik-badge-soft, ik-badge-line, deck-footer, slide-number, notes
