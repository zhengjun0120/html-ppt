# 青莲碧蓝 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `il-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-indigo-lotus` 作用域生效，骨架里已写全，照抄结构即可。
> 深色渐变底、云雾星光与月光圆由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（白是身份，留白是呼吸）**：文字一律白色与浅灰，强调只准用白
> （il-pill-accent / il-btn 的白底、il-stat-v / il-vert-accent 的月光白），
> 禁朱砂等高饱和暖色、禁霓虹。卡只有玻璃卡一种（半透明白 + 细白描边），
> 一页至多一组三张；禁衬线粗体大标题、禁英文 Grotesk 粗黑（h1/h2 已内建细无衬线，
> 不要改字重）。竖排 il-vert 只放留白处，一页至多一条。一页最多两组内容块 + 页脚，
> 把夜空留出来。

---

## cover（月下封面）
指纹：hero

用途：开场页。细字距题签 + 白色细体大标题 + 月光线 + 一句定位，星月云雾自动衬底。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字。

```html
<section class="slide full" data-layout="cover">
  <p class="il-kicker">{{题签，≤14 字，如 夜航船 · 秋季共读}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="il-line mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="il-pill">{{时间，≤12 字}}</span>
    <span class="il-pill">{{地点或人数，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, il-kicker, h1, mt-m, il-line, mt-s, lede, mt-l, row, il-pill, deck-footer, slide-number, notes

---

## contents（航程目录）
指纹：table
数量：il-item=4

用途：议程页。一块玻璃大卡里放 4 行篇目：描边圆点编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="il-kicker">{{引导语，如 四期主题}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="il-card mt-l" style="margin-top:44px">
    <div class="il-item"><span class="il-n">01</span><span class="il-t">{{篇名，≤8 字}}</span><span class="il-d">{{说明，14-26 字}}</span></div>
    <div class="il-item"><span class="il-n">02</span><span class="il-t">{{篇名，≤8 字}}</span><span class="il-d">{{说明，14-26 字}}</span></div>
    <div class="il-item"><span class="il-n">03</span><span class="il-t">{{篇名，≤8 字}}</span><span class="il-d">{{说明，14-26 字}}</span></div>
    <div class="il-item"><span class="il-n">04</span><span class="il-t">{{篇名，≤8 字}}</span><span class="il-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, il-kicker, h2, mt-m, il-card, mt-l, il-item, il-n, il-t, il-d, deck-footer, slide-number, notes

---

## keynotes（星卡要点）
指纹：cards
数量：il-card=3

用途：恰好三张玻璃卡。每张：白色药丸 + 细黑体小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；药丸 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="il-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:30px;margin-top:44px">
    <div class="il-card"><span class="il-pill il-pill-accent">{{药丸，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s">{{说明：是什么 + 怎么做，22-44 字}}</p></div>
    <div class="il-card"><span class="il-pill il-pill-accent">{{药丸，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s">{{说明：怎么做 + 带走什么，22-44 字}}</p></div>
    <div class="il-card"><span class="il-pill il-pill-accent">{{药丸，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, il-kicker, h2, mt-m, grid, g3, mt-l, il-card, il-pill, il-pill-accent, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：il-item=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边玻璃卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="il-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:var(--ink-2)">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="il-pill">{{要点 1，≤8 字}}</span>
        <span class="il-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="il-card">
      <div class="il-item"><span class="il-n">一</span><span class="il-d">{{一步，12-26 字}}</span></div>
      <div class="il-item"><span class="il-n">二</span><span class="il-d">{{一步，12-26 字}}</span></div>
      <div class="il-item"><span class="il-n">三</span><span class="il-d">{{一步，12-26 字}}</span></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, il-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, il-pill, il-card, il-item, il-n, il-d, deck-footer, slide-number, notes

---

## metrics（月光数据）
指纹：chart
数量：il-stat=3

用途：三个关键数据。白色细体大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="il-kicker">{{数据语境，如 前两季 · 复盘}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="il-stat"><div class="il-stat-v">{{数值 ≤6 字符}}<span class="il-stat-u">{{单位}}</span></div><div class="il-stat-l">{{指标名，≤8 字}}</div><p class="il-stat-note">{{口径，14-30 字}}</p></div>
    <div class="il-stat"><div class="il-stat-v">{{数值 ≤6 字符}}<span class="il-stat-u">{{单位}}</span></div><div class="il-stat-l">{{指标名，≤8 字}}</div><p class="il-stat-note">{{口径，14-30 字}}</p></div>
    <div class="il-stat"><div class="il-stat-v">{{数值 ≤6 字符}}<span class="il-stat-u">{{单位}}</span></div><div class="il-stat-l">{{指标名，≤8 字}}</div><p class="il-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:var(--ink-3);letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, il-kicker, h2, mt-m, grid, g3, mt-l, il-stat, il-stat-v, il-stat-u, il-stat-l, il-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（诗句引文）
指纹：quote

用途：整页一句引文。白色细体大字 + 出处 + 两个支撑药丸，右侧留白可立竖排诗句。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字；竖排 ≤7 字。

```html
<section class="slide" data-layout="quote">
  <p class="il-kicker">{{语境，如 开卷第一问}}</p>
  <p class="il-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="il-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="il-pill il-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="il-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="il-vert" style="position:absolute;right:120px;top:50%;transform:translateY(-50%)">{{竖排诗句 ≤7 字}}<span class="il-vert-accent">·</span></div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, il-kicker, il-quote, mt-l, il-src, mt-m, row, il-pill, il-pill-accent, il-vert, il-vert-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="il-kicker">{{进度，如 第二期 · 夜航之志}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="il-pill il-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="il-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, il-kicker, h1, mt-m, lede, mt-l, row, il-pill, il-pill-accent, deck-footer, slide-number, notes

---

## moments（航期表）
指纹：chart
数量：il-tl-item=4

用途：恰好 4 个节点的横向时间线：描边圆点 + 日期 + 一句事件。
适用 role：content / data。
内容约束：恰好 4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="il-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="il-tl mt-l" style="margin-top:52px">
    <div class="il-tl-item"><div class="il-tl-dot">1</div><div class="il-tl-t">{{时间点，≤10 字符}}</div><p class="il-tl-d">{{事件，12-26 字}}</p></div>
    <div class="il-tl-item"><div class="il-tl-dot">2</div><div class="il-tl-t">{{时间点，≤10 字符}}</div><p class="il-tl-d">{{事件，12-26 字}}</p></div>
    <div class="il-tl-item"><div class="il-tl-dot">3</div><div class="il-tl-t">{{时间点，≤10 字符}}</div><p class="il-tl-d">{{事件，12-26 字}}</p></div>
    <div class="il-tl-item"><div class="il-tl-dot">4</div><div class="il-tl-t">{{时间点，≤10 字符}}</div><p class="il-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:var(--ink-3);letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, il-kicker, h2, mt-m, il-tl, mt-l, il-tl-item, il-tl-dot, il-tl-t, il-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（登船收尾）
指纹：hero

用途：收尾页。大字 + 一句行动提醒 + 白底按钮 + 描边药丸。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="il-kicker">{{提醒语境，如 预约船票}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="il-btn">{{按钮文案，≤8 字}}</span>
    <span class="il-pill">{{次级信息，≤10 字}}</span>
    <span class="il-pill">{{地点或价格，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, il-kicker, h1, mt-m, lede, mt-l, row, il-btn, il-pill, deck-footer, slide-number, notes
