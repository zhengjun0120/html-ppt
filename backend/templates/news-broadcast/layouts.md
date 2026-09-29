# 新闻播报 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `nw-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-news-broadcast` 作用域生效，骨架里已写全，照抄结构即可。
> 每页左侧红竖条与顶部红条由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（这是播报，不是故事会）**：正红 #DC2626 只给四类元素——块号/时间线方块、
> 题签方块、提示块重音（nw-pill-accent）、关键数据（nw-stat-v / nw-tl-t）与行动钮；
> 深红 #7F1D1D 是标题色。禁圆角、禁柔和色、禁深色背景；卡片一律白底硬阴影直角。
> 滚动条 nw-ticker 一页至多一条，且只用于 cover / contents / divider / closing。

---

## cover（开播封面）
指纹：hero

用途：开场页。台标框 + 时间戳 + 头条横幅 + 大标题 + 一句定位，页底滚动条压轴。
适用 role：cover。
内容约束：主标题 ≤10 字；头条导语 ≤20 字；lede 30-60 字；台标 ≤6 字；时间戳 ≤16 字符。

```html
<section class="slide full" data-layout="cover">
  <div class="row" style="gap:16px">
    <span class="nw-station">{{台标，≤6 字}}</span>
    <span class="nw-time">{{时间戳，≤16 字符}}</span>
  </div>
  <p class="nw-kicker mt-m">{{栏目语，≤12 字}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="nw-banner mt-m"><span class="nw-banner-tag">头条</span><span class="nw-banner-text">{{头条导语，≤20 字}}</span></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，30-60 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="nw-pill nw-pill-accent">{{要点 1，≤8 字}}</span>
    <span class="nw-pill">{{要点 2，≤8 字}}</span>
  </div>
  <div class="nw-ticker"><span class="nw-ticker-label">快讯</span><span class="nw-ticker-text">{{快讯 1，≤16 字}}</span><span class="nw-ticker-text">{{快讯 2，≤16 字}}</span><span class="nw-ticker-time">{{时刻 · 演播室}}</span></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, row, nw-station, nw-time, nw-kicker, mt-m, h1, nw-banner, nw-banner-tag, nw-banner-text, lede, mt-l, nw-pill, nw-pill-accent, nw-ticker, nw-ticker-label, nw-ticker-text, nw-ticker-time, deck-footer, slide-number, notes

---

## contents（今晚版面）
指纹：table
数量：nw-item=4

用途：议程页。一块白卡里放 4 行板块：红色块号 + 板块名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；板块名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="nw-kicker">{{引导语，如 今晚版面}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="nw-card mt-l" style="margin-top:44px">
    <div class="nw-item"><span class="nw-n">01</span><span class="nw-t">{{板块名，≤8 字}}</span><span class="nw-d">{{说明，14-26 字}}</span></div>
    <div class="nw-item"><span class="nw-n">02</span><span class="nw-t">{{板块名，≤8 字}}</span><span class="nw-d">{{说明，14-26 字}}</span></div>
    <div class="nw-item"><span class="nw-n">03</span><span class="nw-t">{{板块名，≤8 字}}</span><span class="nw-d">{{说明，14-26 字}}</span></div>
    <div class="nw-item"><span class="nw-n">04</span><span class="nw-t">{{板块名，≤8 字}}</span><span class="nw-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="nw-ticker"><span class="nw-ticker-label">快讯</span><span class="nw-ticker-text">{{快讯 1，≤16 字}}</span><span class="nw-ticker-text">{{快讯 2，≤16 字}}</span><span class="nw-ticker-time">{{时刻 · 演播室}}</span></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, nw-kicker, h2, mt-m, nw-card, mt-l, nw-item, nw-n, nw-t, nw-d, nw-ticker, nw-ticker-label, nw-ticker-text, nw-ticker-time, deck-footer, slide-number, notes

---

## keynotes（三条要闻）
指纹：cards
数量：nw-card=3

用途：恰好三张硬阴影白卡。每张：红色题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="nw-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="nw-card"><span class="nw-pill nw-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#991B1B">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="nw-card"><span class="nw-pill nw-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#991B1B">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="nw-card"><span class="nw-pill nw-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#991B1B">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, nw-kicker, h2, mt-m, grid, g3, mt-l, nw-card, nw-pill, nw-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（深度连线）
指纹：split
数量：nw-step=3

用途：左边把一件事讲透（lede + 补充 + 提示块），右边白卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个提示块；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="nw-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#991B1B">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="nw-pill">{{要点 1，≤8 字}}</span>
        <span class="nw-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="nw-card">
      <div class="nw-step"><span class="nw-n">01</span><p class="nw-mini-t">{{一步，12-26 字}}</p></div>
      <div class="nw-step"><span class="nw-n">02</span><p class="nw-mini-t">{{一步，12-26 字}}</p></div>
      <div class="nw-step"><span class="nw-n">03</span><p class="nw-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, nw-kicker, h2, mt-m, grid, g2, mt-l, lede, row, nw-pill, nw-card, nw-step, nw-n, nw-mini-t, deck-footer, slide-number, notes

---

## metrics（数据播报）
指纹：chart
数量：nw-stat=3

用途：三个关键数据。正红大数字是唯一主视觉，口径写进块内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="nw-kicker">{{数据语境，如 防汛指挥部 · 18 时通报}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="nw-stat"><div class="nw-stat-v">{{数值 ≤6 字符}}<span class="nw-stat-u">{{单位}}</span></div><div class="nw-stat-l">{{指标名，≤8 字}}</div><p class="nw-stat-note">{{口径，14-30 字}}</p></div>
    <div class="nw-stat"><div class="nw-stat-v">{{数值 ≤6 字符}}<span class="nw-stat-u">{{单位}}</span></div><div class="nw-stat-l">{{指标名，≤8 字}}</div><p class="nw-stat-note">{{口径，14-30 字}}</p></div>
    <div class="nw-stat"><div class="nw-stat-v">{{数值 ≤6 字符}}<span class="nw-stat-u">{{单位}}</span></div><div class="nw-stat-l">{{指标名，≤8 字}}</div><p class="nw-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#B45454;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, nw-kicker, h2, mt-m, grid, g3, mt-l, nw-stat, nw-stat-v, nw-stat-u, nw-stat-l, nw-stat-note, deck-footer, slide-number, notes

---

## quote（导播金句）
指纹：quote

用途：整页一句引文。粗黑大字 + 出处 + 两个支撑提示块。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；提示块各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="nw-kicker">{{语境，如 通报会答问}}</p>
  <p class="nw-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="nw-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="nw-pill nw-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="nw-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, nw-kicker, nw-quote, mt-l, nw-src, mt-m, row, nw-pill, nw-pill-accent, deck-footer, slide-number, notes

---

## divider（板块转场）
指纹：hero

用途：章节过渡。进度题签 + 大字章节名 + 一个过渡句 + 两个看点提示块，页底滚动条。
适用 role：divider。
内容约束：章节名 ≤12 字；过渡句 22-44 字；提示块各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="nw-kicker">{{进度，如 板块二 · 全市部署}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤12 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="nw-pill nw-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="nw-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="nw-ticker"><span class="nw-ticker-label">快讯</span><span class="nw-ticker-text">{{快讯 1，≤16 字}}</span><span class="nw-ticker-text">{{快讯 2，≤16 字}}</span><span class="nw-ticker-time">{{时刻 · 演播室}}</span></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, nw-kicker, h1, mt-m, lede, mt-l, row, nw-pill, nw-pill-accent, nw-ticker, nw-ticker-label, nw-ticker-text, nw-ticker-time, deck-footer, slide-number, notes

---

## moments（直播时间轴）
指纹：chart
数量：nw-tl-item=4

用途：3-4 个节点的横向时间线：红色方块 + 时刻 + 一句动态。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="nw-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="nw-tl mt-l" style="margin-top:52px">
    <div class="nw-tl-item"><div class="nw-tl-dot">01</div><div class="nw-tl-t">{{时间点，≤10 字符}}</div><p class="nw-tl-d">{{事件，12-26 字}}</p></div>
    <div class="nw-tl-item"><div class="nw-tl-dot">02</div><div class="nw-tl-t">{{时间点，≤10 字符}}</div><p class="nw-tl-d">{{事件，12-26 字}}</p></div>
    <div class="nw-tl-item"><div class="nw-tl-dot">03</div><div class="nw-tl-t">{{时间点，≤10 字符}}</div><p class="nw-tl-d">{{事件，12-26 字}}</p></div>
    <div class="nw-tl-item"><div class="nw-tl-dot">04</div><div class="nw-tl-t">{{时间点，≤10 字符}}</div><p class="nw-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#B45454;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, nw-kicker, h2, mt-m, nw-tl, mt-l, nw-tl-item, nw-tl-dot, nw-tl-t, nw-tl-d, deck-footer, slide-number, notes

---

## closing（结束播报）
指纹：hero

用途：收尾页。大字标题 + 一句行动提醒 + 红底按钮 + 提示块 + 台标框，页底滚动条。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-48 字；按钮 ≤8 字；提示块 ≤14 字。

```html
<section class="slide full" data-layout="closing">
  <p class="nw-kicker">{{提醒语境，≤12 字}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:50ch">{{一句收束提醒，20-48 字}}</p>
  <div class="row mt-l" style="gap:22px">
    <span class="nw-btn">{{按钮文案，≤8 字}}</span>
    <span class="nw-pill">{{次级信息，≤14 字}}</span>
    <span class="nw-station">{{台标，≤6 字}}</span>
  </div>
  <div class="nw-ticker"><span class="nw-ticker-label">快讯</span><span class="nw-ticker-text">{{快讯 1，≤16 字}}</span><span class="nw-ticker-text">{{快讯 2，≤16 字}}</span><span class="nw-ticker-time">{{时刻 · 演播室}}</span></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, nw-kicker, h1, mt-m, lede, mt-l, row, nw-btn, nw-pill, nw-station, nw-ticker, nw-ticker-label, nw-ticker-text, nw-ticker-time, deck-footer, slide-number, notes
