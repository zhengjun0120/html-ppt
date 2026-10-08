# 极简白 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `mw-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-minimal-white` 作用域生效，骨架里已写全，照抄结构即可。
> 每页右侧的巨型「留白」水印与页顶一线蓝由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（字重即层级）**：蓝色 #3B82F6 全模板只有两处——标题短线（mw-rule）与页顶环境细条，
> 任何其他元素不得染色；分隔线一律 1px 发丝线，禁阴影、禁粗边、禁大色块。
> 细字重大标题（h1/h2，字重 200/300）每页最多一组；一页最多两组内容块 + 页脚，留白是构图的一部分。

---

## cover（留白封面）
指纹：hero

用途：开场页。宽字距题签 + 细字重大标题 + 一线蓝短线 + 一句定位，信息签收在标题下。
适用 role：cover。
内容约束：主标题 ≤14 字可两行；lede 30-60 字；题签 ≤14 字；信息签各 ≤12 字。

```html
<section class="slide full" data-layout="cover">
  <p class="mw-kicker">{{题签，≤14 字，如 产品设计复盘 · 2026 夏}}</p>
  <h1 class="h1 mt-m">{{主标题，≤14 字，可 <br> 分两行}}</h1>
  <div class="mw-rule mt-s"></div>
  <p class="lede mt-m" style="max-width:52ch">{{一句话定位，30-60 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="mw-pill">{{信息 1，≤12 字}}</span>
    <span class="mw-pill">{{信息 2，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, mw-kicker, h1, mt-m, mw-rule, mt-s, lede, mt-l, row, mw-pill, deck-footer, slide-number, notes

---

## contents（目录）
指纹：table
数量：mw-item=4

用途：议程页。4 行篇目：细体编号 + 篇名 + 一句说明，发丝线分行。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="mw-kicker">{{引导语，如 本次复盘}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="mw-card mt-l" style="margin-top:44px">
    <div class="mw-item"><span class="mw-n">01</span><span class="mw-t">{{篇名，≤8 字}}</span><span class="mw-d">{{说明，14-26 字}}</span></div>
    <div class="mw-item"><span class="mw-n">02</span><span class="mw-t">{{篇名，≤8 字}}</span><span class="mw-d">{{说明，14-26 字}}</span></div>
    <div class="mw-item"><span class="mw-n">03</span><span class="mw-t">{{篇名，≤8 字}}</span><span class="mw-d">{{说明，14-26 字}}</span></div>
    <div class="mw-item"><span class="mw-n">04</span><span class="mw-t">{{篇名，≤8 字}}</span><span class="mw-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mw-kicker, h2, mt-m, mw-card, mt-l, mw-item, mw-n, mw-t, mw-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：mw-card=3

用途：恰好三张发丝边卡。每张：题签 + 小标题 + 两句说明，层级全靠字重。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="mw-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:36px;margin-top:48px">
    <div class="mw-card"><span class="mw-pill">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:18px;line-height:1.8;color:#475569">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="mw-card"><span class="mw-pill">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:18px;line-height:1.8;color:#475569">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="mw-card"><span class="mw-pill">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:18px;line-height:1.8;color:#475569">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mw-kicker, h2, mt-m, grid, g3, mt-l, mw-card, mw-pill, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：mw-step=3

用途：左边把一件事讲透（lede + 补充 + 签），右边发丝边卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个签；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="mw-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:18px;line-height:1.8;color:#475569">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="mw-pill">{{要点 1，≤8 字}}</span>
        <span class="mw-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="mw-card">
      <div class="mw-step"><span class="mw-n">01</span><p class="mw-mini-t">{{一步，12-26 字}}</p></div>
      <div class="mw-step"><span class="mw-n">02</span><p class="mw-mini-t">{{一步，12-26 字}}</p></div>
      <div class="mw-step"><span class="mw-n">03</span><p class="mw-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mw-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, mw-pill, mw-card, mw-step, mw-n, mw-mini-t, deck-footer, slide-number, notes

---

## metrics（三组数字）
指纹：chart
数量：mw-stat=3

用途：三个关键数据。细字重大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="mw-kicker">{{数据语境，如 数字验证}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="mw-stat"><div class="mw-stat-v">{{数值 ≤6 字符}}<span class="mw-stat-u">{{单位}}</span></div><div class="mw-stat-l">{{指标名，≤8 字}}</div><p class="mw-stat-note">{{口径，14-30 字}}</p></div>
    <div class="mw-stat"><div class="mw-stat-v">{{数值 ≤6 字符}}<span class="mw-stat-u">{{单位}}</span></div><div class="mw-stat-l">{{指标名，≤8 字}}</div><p class="mw-stat-note">{{口径，14-30 字}}</p></div>
    <div class="mw-stat"><div class="mw-stat-v">{{数值 ≤6 字符}}<span class="mw-stat-u">{{单位}}</span></div><div class="mw-stat-l">{{指标名，≤8 字}}</div><p class="mw-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:16px;color:#94A3B8;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mw-kicker, h2, mt-m, grid, g3, mt-l, mw-stat, mw-stat-v, mw-stat-u, mw-stat-l, mw-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（细字引文）
指纹：quote

用途：整页一句引文。细字重大字 + 出处 + 两个细描边签，留白即排版。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；签各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="mw-kicker">{{语境，如 设计手记}}</p>
  <p class="mw-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="mw-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="mw-pill">{{支撑点 1，≤8 字}}</span>
    <span class="mw-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mw-kicker, mw-quote, mt-l, mw-src, mt-m, row, mw-pill, deck-footer, slide-number, notes

---

## divider（章节留白）
指纹：hero

用途：章节过渡。进度题签 + 细字重大字章节名 + 一个过渡问题 + 两个看点签。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-52 字；签各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="mw-kicker">{{进度，如 第二章 · 减法}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="mw-pill">{{看点 1，≤8 字}}</span>
    <span class="mw-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, mw-kicker, h1, mt-m, lede, mt-l, row, mw-pill, deck-footer, slide-number, notes

---

## moments（流程时间线）
指纹：chart
数量：mw-tl-item=4

用途：3-4 个节点的横向时间线：空心圆点 + 时间点 + 一句事件，发丝线串联。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="mw-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="mw-tl mt-l" style="margin-top:56px">
    <div class="mw-tl-item"><div class="mw-tl-dot"></div><div class="mw-tl-t">{{时间点，≤10 字符}}</div><p class="mw-tl-d">{{事件，12-26 字}}</p></div>
    <div class="mw-tl-item"><div class="mw-tl-dot"></div><div class="mw-tl-t">{{时间点，≤10 字符}}</div><p class="mw-tl-d">{{事件，12-26 字}}</p></div>
    <div class="mw-tl-item"><div class="mw-tl-dot"></div><div class="mw-tl-t">{{时间点，≤10 字符}}</div><p class="mw-tl-d">{{事件，12-26 字}}</p></div>
    <div class="mw-tl-item"><div class="mw-tl-dot"></div><div class="mw-tl-t">{{时间点，≤10 字符}}</div><p class="mw-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:16px;color:#94A3B8;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mw-kicker, h2, mt-m, mw-tl, mt-l, mw-tl-item, mw-tl-dot, mw-tl-t, mw-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（收尾）
指纹：hero

用途：收尾页。细字重大字 + 一句行动提醒 + 墨线描边按钮 + 细描边签。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；签 ≤12 字。

```html
<section class="slide full" data-layout="closing">
  <p class="mw-kicker">{{行动语境，如 下半年}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句行动提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:22px">
    <span class="mw-btn">{{按钮文案，≤8 字}}</span>
    <span class="mw-pill">{{次级信息，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, mw-kicker, h1, mt-m, lede, mt-l, row, mw-btn, mw-pill, deck-footer, slide-number, notes
