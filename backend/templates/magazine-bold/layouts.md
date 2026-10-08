# 杂志大字 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `mb-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-magazine-bold` 作用域生效，骨架里已写全，照抄结构即可。
> 每页右上巨型引号与左下琥珀光晕由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（大字即封面）**：超大衬线标题是灵魂——h1/h2 每页最多一组，让字号本身成为画面；
> 琥珀橙只允许出现在四处——荧光笔短杠（mb-bar / mb-kicker 自带）、关键数据（mb-stat-v / mb-n /
> mb-tl-dot）、题签（mb-pill-accent）、按钮（mb-btn）。禁暗色背景、禁无衬线标题、禁橙色泛滥、
> 禁小字号密集排版。巨型水印字 mb-mark 一页至多一个，只放留白处。

---

## cover（封面大字）
指纹：hero

用途：开场页。期号标签 + 超大衬线标题 + 琥珀短杠 + 一句导语，右下可立巨型水印字。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；标签 ≤14 字；水印 1-2 字。

```html
<section class="slide full" data-layout="cover">
  <p class="mb-kicker">{{期号标签，≤14 字，如 城市月刊 · 第 41 期}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="mb-bar mt-s"></div>
  <p class="lede mt-m" style="max-width:48ch">{{一句导语，24-48 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="mb-pill mb-pill-accent">{{栏目签，≤6 字}}</span>
    <span class="mb-pill">{{季节或日期，≤12 字}}</span>
  </div>
  <div class="mb-mark">{{水印字，1-2 字}}</div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, mb-kicker, h1, mt-m, mb-bar, mt-s, lede, mt-m, row, mt-l, mb-pill, mb-pill-accent, mb-mark, deck-footer, slide-number, notes

---

## contents（目录页）
指纹：table
数量：mb-item=4

用途：议程页。一块奶油白卡里放 4 行篇目：琥珀衬线大序号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="mb-kicker">{{引导语，如 本期目录}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="mb-card mt-l" style="margin-top:44px">
    <div class="mb-item"><span class="mb-n">01</span><span class="mb-t">{{篇名，≤8 字}}</span><span class="mb-d">{{说明，14-26 字}}</span></div>
    <div class="mb-item"><span class="mb-n">02</span><span class="mb-t">{{篇名，≤8 字}}</span><span class="mb-d">{{说明，14-26 字}}</span></div>
    <div class="mb-item"><span class="mb-n">03</span><span class="mb-t">{{篇名，≤8 字}}</span><span class="mb-d">{{说明，14-26 字}}</span></div>
    <div class="mb-item"><span class="mb-n">04</span><span class="mb-t">{{篇名，≤8 字}}</span><span class="mb-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mb-kicker, h2, mt-m, mb-card, mt-l, mb-item, mb-n, mb-t, mb-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：mb-card=3

用途：恰好三张奶油白卡。每张：琥珀题签 + 衬线小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="mb-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="mb-card"><span class="mb-pill mb-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.75;color:#78350F">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="mb-card"><span class="mb-pill mb-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.75;color:#78350F">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="mb-card"><span class="mb-pill mb-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.75;color:#78350F">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mb-kicker, h2, mt-m, grid, g3, mt-l, mb-card, mb-pill, mb-pill-accent, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：mb-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边奶油卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="mb-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.75;color:#78350F">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="mb-pill">{{要点 1，≤8 字}}</span>
        <span class="mb-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="mb-card">
      <div class="mb-step"><span class="mb-n">一</span><p class="mb-mini-t">{{一步，12-26 字}}</p></div>
      <div class="mb-step"><span class="mb-n">二</span><p class="mb-mini-t">{{一步，12-26 字}}</p></div>
      <div class="mb-step"><span class="mb-n">三</span><p class="mb-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mb-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, mt-l, mb-pill, mb-card, mb-step, mb-n, mb-mini-t, deck-footer, slide-number, notes

---

## metrics（琥珀数字）
指纹：chart
数量：mb-stat=3

用途：三个关键数据。琥珀衬线大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="mb-kicker">{{数据语境，如 数据页 · 口径在此}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="mb-stat"><div class="mb-stat-v">{{数值 ≤6 字符}}<span class="mb-stat-u">{{单位}}</span></div><div class="mb-stat-l">{{指标名，≤8 字}}</div><p class="mb-stat-note">{{口径，14-30 字}}</p></div>
    <div class="mb-stat"><div class="mb-stat-v">{{数值 ≤6 字符}}<span class="mb-stat-u">{{单位}}</span></div><div class="mb-stat-l">{{指标名，≤8 字}}</div><p class="mb-stat-note">{{口径，14-30 字}}</p></div>
    <div class="mb-stat"><div class="mb-stat-v">{{数值 ≤6 字符}}<span class="mb-stat-u">{{单位}}</span></div><div class="mb-stat-l">{{指标名，≤8 字}}</div><p class="mb-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:16px;font-weight:500;color:#A16207;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mb-kicker, h2, mt-m, grid, g3, mt-l, mb-stat, mb-stat-v, mb-stat-u, mb-stat-l, mb-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（金句拉页）
指纹：quote

用途：整页一句金句。超大衬线 + 出处 + 两个支撑药丸，右侧留白可立巨型水印字。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="mb-kicker">{{语境，如 记者手记}}</p>
  <p class="mb-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="mb-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="mb-pill mb-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="mb-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="mb-mark" style="bottom:60px">{{水印字，1-2 字}}</div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mb-kicker, mb-quote, mt-l, mb-src, mt-m, row, mt-l, mb-pill, mb-pill-accent, mb-mark, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度标签 + 超大衬线章节名 + 一个过渡句 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡句 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="mb-kicker">{{进度，如 第二章 · 人群画像}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:46ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="mb-pill mb-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="mb-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, mb-kicker, h1, mt-m, lede, mt-l, row, mb-pill, mb-pill-accent, deck-footer, slide-number, notes

---

## moments（流程时间线）
指纹：chart
数量：mb-tl-item=4

用途：3-4 个节点的横向时间线：琥珀圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="mb-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="mb-tl mt-l" style="margin-top:52px">
    <div class="mb-tl-item"><div class="mb-tl-dot">01</div><div class="mb-tl-t">{{时间点，≤10 字符}}</div><p class="mb-tl-d">{{事件，12-26 字}}</p></div>
    <div class="mb-tl-item"><div class="mb-tl-dot">02</div><div class="mb-tl-t">{{时间点，≤10 字符}}</div><p class="mb-tl-d">{{事件，12-26 字}}</p></div>
    <div class="mb-tl-item"><div class="mb-tl-dot">03</div><div class="mb-tl-t">{{时间点，≤10 字符}}</div><p class="mb-tl-d">{{事件，12-26 字}}</p></div>
    <div class="mb-tl-item"><div class="mb-tl-dot">04</div><div class="mb-tl-t">{{时间点，≤10 字符}}</div><p class="mb-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:16px;font-weight:500;color:#A16207;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mb-kicker, h2, mt-m, mb-tl, mt-l, mb-tl-item, mb-tl-dot, mb-tl-t, mb-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（收尾订阅）
指纹：hero

用途：收尾页。超大衬线 + 一句行动提醒 + 琥珀实底按钮 + 描边药丸。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="mb-kicker">{{提醒语境，如 下期预告}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:46ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="mb-btn">{{按钮文案，≤8 字}}</span>
    <span class="mb-pill">{{次级信息，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, mb-kicker, h1, mt-m, lede, mt-l, row, mb-btn, mb-pill, deck-footer, slide-number, notes
