# 藏青学衡 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `an-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-academic-navy` 作用域生效，骨架里已写全，照抄结构即可。
> 右上金环与左下六边形由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（证据即装饰）**：金棕只允许出现在四处——关键数字（an-num / an-stat-v /
> an-tl-t / an-side-v）、强调药丸（an-pill-accent）、题签小字（an-kicker）、金顶线卡
> （an-card-gold，仅限引证场景）。标题一律藏青，禁给标题描金；一页只允许一种强调方式；
> 白卡只用藏青顶线分组，禁硬阴影与大面积色块。

---

## cover（答辩封面）
指纹：hero

用途：开场页。英文题签 + 藏青大标题 + 一句研究定位，右侧立一枚金棕索引数，药丸标答辩人与日期。
适用 role：cover。
内容约束：主标题 ≤14 字；lede 24-56 字；题签 ≤26 字符；索引数 ≤4 字符。

```html
<section class="slide full" data-layout="cover">
  <p class="an-kicker">{{题签，如 RESEARCH DEFENSE · 2026}}</p>
  <h1 class="h1 mt-m">{{论文标题，≤14 字，可 <br> 分行}}</h1>
  <div class="an-rule mt-s"></div>
  <p class="lede mt-m" style="max-width:52ch">{{一句研究定位，24-56 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="an-pill">{{答辩人 · 导师，≤16 字}}</span>
    <span class="an-pill an-pill-accent">{{答辩日期与场合，≤16 字}}</span>
  </div>
  <div class="an-side"><span class="an-side-v">{{索引数，≤4 字符}}</span><span class="an-side-l">{{英文索引名，≤20 字符}}</span></div>
  <div class="deck-footer"><span>{{校名 · 院系}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, an-kicker, h1, mt-m, an-rule, mt-s, lede, row, mt-l, an-pill, an-pill-accent, an-side, an-side-v, an-side-l, deck-footer, slide-number, notes

---

## contents（答辩提纲）
指纹：table
数量：an-item=4

用途：提纲页。一块白卡里放 4 行板块：环形编号 + 板块名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；板块名 ≤8 字；说明 12-24 字。

```html
<section class="slide" data-layout="contents">
  <p class="an-kicker">{{引导语，如 OUTLINE · 答辩提纲}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="an-card mt-l" style="margin-top:44px">
    <div class="an-item"><span class="an-n">01</span><span class="an-t">{{板块名，≤8 字}}</span><span class="an-d">{{说明，12-24 字}}</span></div>
    <div class="an-item"><span class="an-n">02</span><span class="an-t">{{板块名，≤8 字}}</span><span class="an-d">{{说明，12-24 字}}</span></div>
    <div class="an-item"><span class="an-n">03</span><span class="an-t">{{板块名，≤8 字}}</span><span class="an-d">{{说明，12-24 字}}</span></div>
    <div class="an-item"><span class="an-n">04</span><span class="an-t">{{板块名，≤8 字}}</span><span class="an-d">{{说明，12-24 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{校名 · 院系}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, an-kicker, h2, mt-m, an-card, mt-l, an-item, an-n, an-t, an-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：an-card=3, an-num=3

用途：恰好三张藏青顶线白卡。每张：金棕编号 + 小标题 + 两句说明，论文式证据链的一环。
适用 role：content。
内容约束：恰好 3 卡；编号 2 字符；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="an-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="an-card"><span class="an-num">01</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#404040">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="an-card"><span class="an-num">02</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#404040">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="an-card"><span class="an-num">03</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#404040">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{校名 · 院系}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, an-kicker, h2, mt-m, grid, g3, mt-l, an-card, an-num, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左证右列）
指纹：split
数量：an-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边白卡装三行步骤或规程，如观测流程与质控口径。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="an-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#404040">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="an-pill">{{要点 1，≤8 字}}</span>
        <span class="an-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="an-card">
      <div class="an-step"><span class="an-n">01</span><p class="an-mini-t">{{一步，12-26 字}}</p></div>
      <div class="an-step"><span class="an-n">02</span><p class="an-mini-t">{{一步，12-26 字}}</p></div>
      <div class="an-step"><span class="an-n">03</span><p class="an-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{校名 · 院系}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, an-kicker, h2, mt-m, grid, g2, mt-l, lede, row, an-pill, an-card, an-step, an-n, an-mini-t, deck-footer, slide-number, notes

---

## metrics（关键数据）
指纹：chart
数量：an-stat=3

用途：三个关键数据。金棕大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 12-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="an-kicker">{{数据语境，如 DATASET · 数据规模}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="an-stat"><div class="an-stat-v">{{数值 ≤6 字符}}<span class="an-stat-u">{{单位}}</span></div><div class="an-stat-l">{{指标名，≤8 字}}</div><p class="an-stat-note">{{口径，12-30 字}}</p></div>
    <div class="an-stat"><div class="an-stat-v">{{数值 ≤6 字符}}<span class="an-stat-u">{{单位}}</span></div><div class="an-stat-l">{{指标名，≤8 字}}</div><p class="an-stat-note">{{口径，12-30 字}}</p></div>
    <div class="an-stat"><div class="an-stat-v">{{数值 ≤6 字符}}<span class="an-stat-u">{{单位}}</span></div><div class="an-stat-l">{{指标名，≤8 字}}</div><p class="an-stat-note">{{口径，12-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#808080;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{校名 · 院系}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, an-kicker, h2, mt-m, grid, g3, mt-l, an-stat, an-stat-v, an-stat-u, an-stat-l, an-stat-note, deck-footer, slide-number, notes

---

## quote（治学题记）
指纹：quote

用途：整页一句引文。衬线大字 + 出处 + 两个支撑药丸，为答辩立一条治学信条。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="an-kicker">{{语境，如 EPIGRAPH · 题记}}</p>
  <p class="an-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="an-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="an-pill an-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="an-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{校名 · 院系}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, an-kicker, an-quote, mt-l, an-src, mt-m, row, an-pill, an-pill-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。英文进度题签 + 藏青大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤12 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="an-kicker">{{进度，如 CHAPTER 02 · 数据与方法}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="an-pill an-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="an-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{校名 · 院系}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, an-kicker, h1, mt-m, lede, row, mt-l, an-pill, an-pill-accent, deck-footer, slide-number, notes

---

## moments（历程时间线）
指纹：chart
数量：an-tl-item=4

用途：3-4 个节点的横向时间线：藏青圆点编号 + 时间点 + 一句事件，讲研究怎么一步步做下来。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="an-kicker">{{引导语，如 TIMELINE · 研究历程}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="an-tl mt-l" style="margin-top:52px">
    <div class="an-tl-item"><div class="an-tl-dot">01</div><div class="an-tl-t">{{时间点，≤10 字符}}</div><p class="an-tl-d">{{事件，12-26 字}}</p></div>
    <div class="an-tl-item"><div class="an-tl-dot">02</div><div class="an-tl-t">{{时间点，≤10 字符}}</div><p class="an-tl-d">{{事件，12-26 字}}</p></div>
    <div class="an-tl-item"><div class="an-tl-dot">03</div><div class="an-tl-t">{{时间点，≤10 字符}}</div><p class="an-tl-d">{{事件，12-26 字}}</p></div>
    <div class="an-tl-item"><div class="an-tl-dot">04</div><div class="an-tl-t">{{时间点，≤10 字符}}</div><p class="an-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#808080;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{校名 · 院系}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, an-kicker, h2, mt-m, an-tl, mt-l, an-tl-item, an-tl-dot, an-tl-t, an-tl-d, deck-footer, slide-number, notes

---

## closing（致谢收尾）
指纹：hero

用途：收尾页。藏青大字致谢 + 一句说明 + 藏青按钮 + 药丸，如答辩落款。
适用 role：thanks / cta / content。
内容约束：标题 ≤12 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤14 字。

```html
<section class="slide full" data-layout="closing">
  <p class="an-kicker">{{提醒语境，如 Q&A · 敬请指正}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束说明，20-44 字}}</p>
  <div class="row mt-l" style="gap:22px">
    <span class="an-btn">{{按钮文案，≤8 字}}</span>
    <span class="an-pill">{{次级信息，≤14 字}}</span>
  </div>
  <div class="deck-footer"><span>{{校名 · 院系}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, an-kicker, h1, mt-m, lede, row, mt-l, an-btn, an-pill, deck-footer, slide-number, notes
