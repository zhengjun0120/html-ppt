# 瑞士网格 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `sg-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-swiss-grid` 作用域生效，骨架里已写全，照抄结构即可。
> 每页的外露 12 栏网格线由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（红色是唯一出口）**：正红 #EF4444 只允许出现在——红色标签 sg-kicker、
> 短线与节点 sg-line / sg-tl-item、红块 sg-square、红顶线卡 sg-card-top、红色数值
> sg-stat-hot、红边标签 sg-tag-red、行动钮与通栏条 sg-btn / sg-bar。一页红色元素
> ≤3 处；其余皆由深墨与近黑承担秩序。全部直角，禁圆角、禁阴影、禁渐变。

---

## cover（宣言封面）
指纹：hero

用途：开场页。红色标签 + 900 大字标题 + 红色短线，红块与近黑块双拼，页底红色通栏条。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 30-56 字；红块内 ≤6 字符；近黑块主张 24-44 字。

```html
<section class="slide full" data-layout="cover">
  <p class="sg-kicker">{{工作室名 · PORTFOLIO 2026，≤16 字符}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="sg-line mt-s" style="width:280px"></div>
  <p class="lede mt-m" style="max-width:52ch">{{定位：几人工作室 + 只做什么，30-56 字}}</p>
  <div class="grid g2 mt-l" style="gap:40px">
    <div class="sg-square" style="min-height:170px"><span>{{工作室缩写，≤6 字符}}</span></div>
    <div class="sg-block"><p class="sg-block-label">STYLE DNA</p><p class="sg-block-d">{{风格主张一句，24-44 字}}</p></div>
  </div>
  <div class="sg-bar mt-l" style="margin-top:44px"><span>{{工作室全名}}</span><span>PORTFOLIO / {{年份}}</span></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, sg-kicker, h1, mt-m, sg-line, mt-s, lede, grid, g2, mt-l, sg-square, sg-block, sg-block-label, sg-block-d, sg-bar, deck-footer, slide-number, notes

---

## contents（目录栏）
指纹：table
数量：sg-toc-row=4

用途：议程页。顶部黑粗线 + 4 行目录：红色 mono 索引 + 900 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="sg-kicker">INDEX · {{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="mt-l" style="margin-top:40px;border-top:2px solid #0F172A">
    <div class="sg-toc-row"><span class="sg-i">01</span><span class="sg-t">{{篇名，≤8 字}}</span><span class="sg-d">{{说明，14-26 字}}</span></div>
    <div class="sg-toc-row"><span class="sg-i">02</span><span class="sg-t">{{篇名，≤8 字}}</span><span class="sg-d">{{说明，14-26 字}}</span></div>
    <div class="sg-toc-row"><span class="sg-i">03</span><span class="sg-t">{{篇名，≤8 字}}</span><span class="sg-d">{{说明，14-26 字}}</span></div>
    <div class="sg-toc-row"><span class="sg-i">04</span><span class="sg-t">{{篇名，≤8 字}}</span><span class="sg-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sg-kicker, h2, mt-m, mt-l, sg-toc-row, sg-i, sg-t, sg-d, deck-footer, slide-number, notes

---

## metrics（数据锚点）
指纹：chart
数量：sg-stat=3

用途：三个关键数据。黑顶线统计块 + 900 巨号，其中一个数值允许红色；口径写页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="sg-kicker">DATA · {{数据语境}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="sg-stat"><div class="sg-stat-v">{{数值 ≤6 字符}}</div><div class="sg-stat-l">{{指标名，≤8 字}}</div><p class="sg-stat-note">{{口径，14-30 字}}</p></div>
    <div class="sg-stat sg-stat-hot"><div class="sg-stat-v">{{数值 ≤6 字符}}</div><div class="sg-stat-l">{{指标名，≤8 字}}</div><p class="sg-stat-note">{{口径，14-30 字}}</p></div>
    <div class="sg-stat"><div class="sg-stat-v">{{数值 ≤6 字符}}</div><div class="sg-stat-l">{{指标名，≤8 字}}</div><p class="sg-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="sg-src mt-m">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sg-kicker, h2, mt-m, grid, g3, mt-l, sg-stat, sg-stat-hot, sg-stat-v, sg-stat-l, sg-stat-note, sg-src, mt-m, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度标签 + 900 大字章节名 + 红色短线 + 过渡问题 + 两个方角标签。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；标签各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="sg-kicker">{{进度，如 02 / METHOD}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字，可 <br> 分行}}</h1>
  <div class="sg-line mt-m" style="width:220px"></div>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:14px">
    <span class="sg-tag">{{看点 1，≤8 字}}</span>
    <span class="sg-tag sg-tag-red">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, sg-kicker, h1, mt-m, sg-line, lede, row, mt-l, sg-tag, sg-tag-red, deck-footer, slide-number, notes

---

## keynotes（三栏作品）
指纹：cards
数量：sg-card=3

用途：恰好三张白卡。首卡带 6px 红顶线做唯一焦点；mono 编号 + 800 标题 + 说明。
适用 role：content。
内容约束：恰好 3 卡；编号 01-03；标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="sg-kicker">{{章节标签}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="sg-card sg-card-top"><p class="sg-card-i">01</p><h4 class="mt-s">{{作品名，≤8 字}}</h4><p class="sg-card-d">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="sg-card"><p class="sg-card-i">02</p><h4 class="mt-s">{{作品名，≤8 字}}</h4><p class="sg-card-d">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="sg-card"><p class="sg-card-i">03</p><h4 class="mt-s">{{作品名，≤8 字}}</h4><p class="sg-card-d">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sg-kicker, h2, mt-m, grid, g3, mt-l, sg-card, sg-card-top, sg-card-i, h4, mt-s, sg-card-d, deck-footer, slide-number, notes

---

## split（工序拆解）
指纹：split
数量：sg-step=3

用途：左边把方法论讲透（lede + 补充 + 方角标签），右边白卡装三行工序。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个标签；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="sg-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:64px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：工序怎么走、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:18px;line-height:1.8;color:#334155">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:14px">
        <span class="sg-tag">{{要点 1，≤8 字}}</span>
        <span class="sg-tag">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="sg-card">
      <div class="sg-step"><span class="sg-i">01</span><p class="sg-mini-t">{{一步，12-26 字}}</p></div>
      <div class="sg-step"><span class="sg-i">02</span><p class="sg-mini-t">{{一步，12-26 字}}</p></div>
      <div class="sg-step"><span class="sg-i">03</span><p class="sg-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sg-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, sg-tag, sg-card, sg-step, sg-i, sg-mini-t, deck-footer, slide-number, notes

---

## moments（排期网格）
指纹：chart
数量：sg-tl-item=4

用途：3-4 个节点的横向时间线：红色方块节点 + mono 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="sg-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="sg-tl mt-l" style="margin-top:52px">
    <div class="sg-tl-item"><div class="sg-tl-t">{{时间点，≤10 字符}}</div><p class="sg-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sg-tl-item"><div class="sg-tl-t">{{时间点，≤10 字符}}</div><p class="sg-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sg-tl-item"><div class="sg-tl-t">{{时间点，≤10 字符}}</div><p class="sg-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sg-tl-item"><div class="sg-tl-t">{{时间点，≤10 字符}}</div><p class="sg-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="sg-src mt-m">{{口径或读法，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sg-kicker, h2, mt-m, sg-tl, mt-l, sg-tl-item, sg-tl-t, sg-tl-d, sg-src, mt-m, deck-footer, slide-number, notes

---

## quote（工作宣言）
指纹：quote

用途：整页一句工作室宣言。900 大字两行（下行套红）+ mono 出处 + 两个方角标签。
适用 role：quote。
内容约束：引文共 10-28 字；出处 ≤22 字且真实；标签各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="sg-kicker">{{语境，如 主理人说}}</p>
  <p class="sg-quote mt-l" style="margin-top:40px">「{{引文上行，≤14 字}}<br><span class="sg-accent">{{引文下行，≤14 字}}</span>」</p>
  <p class="sg-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:14px">
    <span class="sg-tag sg-tag-red">{{支撑点 1，≤8 字}}</span>
    <span class="sg-tag">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sg-kicker, sg-quote, mt-l, sg-accent, sg-src, mt-m, row, sg-tag, sg-tag-red, deck-footer, slide-number, notes

---

## closing（下一格）
指纹：hero

用途：收尾邀约。900 大字 + 一句下一步 + 红底行动钮 + 方角标签 + 页底红色通栏条。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-48 字；按钮 ≤6 字；标签各 ≤10 字；通栏条两段各 ≤14 字符。

```html
<section class="slide full" data-layout="closing">
  <p class="sg-kicker">{{语境，如 2026 Q4 档期}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{下一步与邀约，20-48 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="sg-btn">{{按钮文案，≤6 字}}</span>
    <span class="sg-tag">{{次级信息，≤10 字}}</span>
    <span class="sg-tag sg-tag-red">{{次级信息，≤10 字}}</span>
  </div>
  <div class="sg-bar mt-l" style="margin-top:48px"><span>{{工作室全名}}</span><span>{{结束语，≤14 字符}}</span></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, sg-kicker, h1, mt-m, lede, row, mt-l, sg-btn, sg-tag, sg-tag-red, sg-bar, deck-footer, slide-number, notes
