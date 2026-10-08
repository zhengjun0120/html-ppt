# 春日嫩柳 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `ps-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-chinese-pastel-spring` 作用域生效，骨架里已写全，照抄结构即可。
> 每页顶部四段渐变色带与右上桃枝由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（柔色即身份）**：嫩柳绿、柳芽黄绿、桃花粉只出现在实色小色卡（ps-chip-*）、
> 圆点编号（ps-n / ps-tl-dot，配色由 nth-child 自动轮换，不要手工改）与数据顶线（ps-stat）；
> 朱砂只允许四处——印章（ps-seal）、题签短线（ps-kicker 自带）、关键数据
> （ps-stat-v / ps-tl-t）、强调药丸（ps-pill-accent）。淡金只做第三条数据顶线点缀。
> 禁高饱和荧光、禁赛博霓虹、禁硬朗几何、禁沉重暗色。楷宋大字（h1/h2）每页最多一组；
> 实色小色卡只在 keynotes 用，一页至多一组；右上桃枝是模板自动衬底，正文不要画花。

---

## cover（春笺封面）
指纹：hero

用途：开场页。题签 + 宋体大字标题 + 四色短线 + 一句定位，右侧可立一枚朱砂印。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；印章 4 字内。

```html
<section class="slide full" data-layout="cover">
  <p class="ps-kicker">{{题签，≤14 字，如 春系列 · 新品发布}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="ps-line mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="ps-seal">{{印章 2×2 字}}</div>
    <span class="ps-pill">{{时间或场合，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 系列}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, ps-kicker, h1, mt-m, ps-line, mt-s, lede, mt-l, row, ps-seal, ps-pill, deck-footer, slide-number, notes

---

## contents（春日目录）
指纹：table
数量：ps-item=4

用途：议程页。一块宣纸卡里放 4 行篇目：柔色圆点编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="ps-kicker">{{引导语，如 今晚流程}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="ps-card mt-l" style="margin-top:44px">
    <div class="ps-item"><span class="ps-n">01</span><span class="ps-t">{{篇名，≤8 字}}</span><span class="ps-d">{{说明，14-26 字}}</span></div>
    <div class="ps-item"><span class="ps-n">02</span><span class="ps-t">{{篇名，≤8 字}}</span><span class="ps-d">{{说明，14-26 字}}</span></div>
    <div class="ps-item"><span class="ps-n">03</span><span class="ps-t">{{篇名，≤8 字}}</span><span class="ps-d">{{说明，14-26 字}}</span></div>
    <div class="ps-item"><span class="ps-n">04</span><span class="ps-t">{{篇名，≤8 字}}</span><span class="ps-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 系列}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ps-kicker, h2, mt-m, ps-card, mt-l, ps-item, ps-n, ps-t, ps-d, deck-footer, slide-number, notes

---

## keynotes（三色卡要点）
指纹：cards
数量：ps-chip=3

用途：恰好三张实色小色卡。每张：无衬线编号 + 楷体小标题 + 两句说明，嫩柳、柳芽、桃花各一。
适用 role：content。
内容约束：恰好 3 卡；编号由骨架给出；卡名 ≤6 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="ps-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:30px;margin-top:44px">
    <div class="ps-chip ps-chip-green"><span class="ps-chip-num">01</span><h4 class="mt-m">{{卡名，≤6 字}}</h4><p class="mt-s">{{说明：是什么 + 好在哪，22-44 字}}</p></div>
    <div class="ps-chip ps-chip-bud"><span class="ps-chip-num">02</span><h4 class="mt-m">{{卡名，≤6 字}}</h4><p class="mt-s">{{说明：怎么做 + 带走什么，22-44 字}}</p></div>
    <div class="ps-chip ps-chip-pink"><span class="ps-chip-num">03</span><h4 class="mt-m">{{卡名，≤6 字}}</h4><p class="mt-s">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 系列}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ps-kicker, h2, mt-m, grid, g3, mt-l, ps-chip, ps-chip-green, ps-chip-bud, ps-chip-pink, ps-chip-num, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：ps-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边暖卡（左竖线）装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="ps-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:var(--ink-2)">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="ps-pill">{{要点 1，≤8 字}}</span>
        <span class="ps-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="ps-card-line">
      <div class="ps-step"><span class="ps-n">1</span><span class="ps-d">{{一步，12-26 字}}</span></div>
      <div class="ps-step"><span class="ps-n">2</span><span class="ps-d">{{一步，12-26 字}}</span></div>
      <div class="ps-step"><span class="ps-n">3</span><span class="ps-d">{{一步，12-26 字}}</span></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 系列}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ps-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, ps-pill, ps-card-line, ps-step, ps-n, ps-d, deck-footer, slide-number, notes

---

## metrics（朱砂数字）
指纹：chart
数量：ps-stat=3

用途：三个关键数据。朱砂细体大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="ps-kicker">{{数据语境，如 上一季 · 复盘}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="ps-stat"><div class="ps-stat-v">{{数值 ≤6 字符}}<span class="ps-stat-u">{{单位}}</span></div><div class="ps-stat-l">{{指标名，≤8 字}}</div><p class="ps-stat-note">{{口径，14-30 字}}</p></div>
    <div class="ps-stat"><div class="ps-stat-v">{{数值 ≤6 字符}}<span class="ps-stat-u">{{单位}}</span></div><div class="ps-stat-l">{{指标名，≤8 字}}</div><p class="ps-stat-note">{{口径，14-30 字}}</p></div>
    <div class="ps-stat"><div class="ps-stat-v">{{数值 ≤6 字符}}<span class="ps-stat-u">{{单位}}</span></div><div class="ps-stat-l">{{指标名，≤8 字}}</div><p class="ps-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:var(--ink-3);letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 系列}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ps-kicker, h2, mt-m, grid, g3, mt-l, ps-stat, ps-stat-v, ps-stat-u, ps-stat-l, ps-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（题跋引文）
指纹：quote

用途：整页一句引文。楷体大字 + 出处 + 两个支撑药丸。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="ps-kicker">{{语境，如 调色师手记}}</p>
  <p class="ps-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="ps-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="ps-pill ps-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="ps-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 系列}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ps-kicker, ps-quote, mt-l, ps-src, mt-m, row, ps-pill, ps-pill-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="ps-kicker">{{进度，如 第二幕 · 三支单品}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="ps-pill ps-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="ps-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 系列}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, ps-kicker, h1, mt-m, lede, mt-l, row, ps-pill, ps-pill-accent, deck-footer, slide-number, notes

---

## moments（新品节奏）
指纹：chart
数量：ps-tl-item=4

用途：恰好 4 个节点的横向时间线：柔色圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：恰好 4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="ps-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="ps-tl mt-l" style="margin-top:52px">
    <div class="ps-tl-item"><div class="ps-tl-dot">1</div><div class="ps-tl-t">{{时间点，≤10 字符}}</div><p class="ps-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ps-tl-item"><div class="ps-tl-dot">2</div><div class="ps-tl-t">{{时间点，≤10 字符}}</div><p class="ps-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ps-tl-item"><div class="ps-tl-dot">3</div><div class="ps-tl-t">{{时间点，≤10 字符}}</div><p class="ps-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ps-tl-item"><div class="ps-tl-dot">4</div><div class="ps-tl-t">{{时间点，≤10 字符}}</div><p class="ps-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:var(--ink-3);letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 系列}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ps-kicker, h2, mt-m, ps-tl, mt-l, ps-tl-item, ps-tl-dot, ps-tl-t, ps-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（收尾落款）
指纹：hero

用途：收尾页。大字 + 一句行动提醒 + 墨黛描边按钮 + 描边药丸 + 朱砂印。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="ps-kicker">{{提醒语境，如 预约直播}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="ps-btn">{{按钮文案，≤8 字}}</span>
    <span class="ps-pill">{{次级信息，≤10 字}}</span>
    <div class="ps-seal">{{印章 2 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 系列}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, ps-kicker, h1, mt-m, lede, mt-l, row, ps-btn, ps-pill, ps-seal, deck-footer, slide-number, notes
