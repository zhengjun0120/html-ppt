# 夏日暖色 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `sw-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-summer-warm-color` 作用域生效，骨架里已写全，照抄结构即可。
> 每页的天空→草地渐变、太阳、云朵与纸飞机由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（蓝是主调，黄是温度）**：阳光黄只允许四处——小太阳印（sw-sun）、编号圆
> （sw-n）、题签药丸（sw-pill-accent / sw-kicker 的粉点）、数据上边线（sw-stat 顶线）。
> 暖橙只做时间点（sw-tl-t）。禁霓虹荧光、禁尖锐直角与硬阴影、禁大面积冷蓝灰无暖色、
> 禁把画面塞满。楷书大字（h1/h2）每页最多一组；sw-cloud 只在 quote 留白处一朵。

---

## cover（晴空封面）
指纹：hero

用途：开场页。云朵题签 + 手写气质大标题 + 五圆色卡阵列 + 一句定位，右侧立一枚小太阳印。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；太阳印 4 字内。

```html
<section class="slide full" data-layout="cover">
  <p class="sw-kicker">{{题签，≤14 字，如 夏日冰饮产品线企划}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="sw-swatches mt-m" aria-hidden="true"><span class="sw-swatch" style="background:#20A1EE"></span><span class="sw-swatch" style="background:#9DDF67"></span><span class="sw-swatch" style="background:#FFF061"></span><span class="sw-swatch" style="background:#FE7F08"></span><span class="sw-swatch" style="background:#FFADDD"></span></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="sw-sun">{{太阳印 2×2 字}}</div>
    <span class="sw-pill">{{副题或时间地点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, sw-kicker, h1, mt-m, sw-swatches, sw-swatch, lede, row, mt-l, sw-sun, sw-pill, deck-footer, slide-number, notes

---

## contents（目录草坡）
指纹：table
数量：sw-item=4

用途：议程页。一块白色圆角大卡里放 4 行篇目：阳光圆号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="sw-kicker">{{引导语，如 议程四拍}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="sw-card mt-l" style="margin-top:44px">
    <div class="sw-item"><span class="sw-n">01</span><span class="sw-t">{{篇名，≤8 字}}</span><span class="sw-d">{{说明，14-26 字}}</span></div>
    <div class="sw-item"><span class="sw-n">02</span><span class="sw-t">{{篇名，≤8 字}}</span><span class="sw-d">{{说明，14-26 字}}</span></div>
    <div class="sw-item"><span class="sw-n">03</span><span class="sw-t">{{篇名，≤8 字}}</span><span class="sw-d">{{说明，14-26 字}}</span></div>
    <div class="sw-item"><span class="sw-n">04</span><span class="sw-t">{{篇名，≤8 字}}</span><span class="sw-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sw-kicker, h2, mt-m, sw-card, mt-l, sw-item, sw-n, sw-t, sw-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：sw-card=3

用途：恰好三张白色圆角卡。每张：阳光黄题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="sw-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="sw-card"><span class="sw-pill sw-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#3F6478">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="sw-card"><span class="sw-pill sw-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#3F6478">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="sw-card"><span class="sw-pill sw-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#3F6478">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sw-kicker, h2, mt-m, grid, g3, mt-l, sw-card, sw-pill, sw-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：sw-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边圆角卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="sw-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#3F6478">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="sw-pill">{{要点 1，≤8 字}}</span>
        <span class="sw-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="sw-card">
      <div class="sw-step"><span class="sw-n">一</span><p class="sw-mini-t">{{一步，12-26 字}}</p></div>
      <div class="sw-step"><span class="sw-n">二</span><p class="sw-mini-t">{{一步，12-26 字}}</p></div>
      <div class="sw-step"><span class="sw-n">三</span><p class="sw-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sw-kicker, h2, mt-m, grid, g2, mt-l, lede, row, sw-pill, sw-card, sw-step, sw-n, sw-mini-t, deck-footer, slide-number, notes

---

## metrics（阳光数字）
指纹：chart
数量：sw-stat=3

用途：三个关键数据。深蓝大数字是主视觉，阳光黄只做上边线，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="sw-kicker">{{数据语境，如 试饮与市场 · 数据}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="sw-stat"><div class="sw-stat-v">{{数值 ≤6 字符}}<span class="sw-stat-u">{{单位}}</span></div><div class="sw-stat-l">{{指标名，≤8 字}}</div><p class="sw-stat-note">{{口径，14-30 字}}</p></div>
    <div class="sw-stat"><div class="sw-stat-v">{{数值 ≤6 字符}}<span class="sw-stat-u">{{单位}}</span></div><div class="sw-stat-l">{{指标名，≤8 字}}</div><p class="sw-stat-note">{{口径，14-30 字}}</p></div>
    <div class="sw-stat"><div class="sw-stat-v">{{数值 ≤6 字符}}<span class="sw-stat-u">{{单位}}</span></div><div class="sw-stat-l">{{指标名，≤8 字}}</div><p class="sw-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#4E7E96;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sw-kicker, h2, mt-m, grid, g3, mt-l, sw-stat, sw-stat-v, sw-stat-u, sw-stat-l, sw-stat-note, deck-footer, slide-number, notes

---

## quote（云朵引文）
指纹：quote

用途：整页一句引文。楷书大字 + 出处 + 两个支撑药丸，右侧留白处浮一朵会写字的云。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="sw-kicker">{{语境，如 试饮笔记}}</p>
  <p class="sw-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="sw-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="sw-pill sw-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="sw-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="sw-cloud">{{云上短语，≤6 字}}</div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sw-kicker, sw-quote, mt-l, sw-src, mt-m, row, sw-pill, sw-pill-accent, sw-cloud, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。云朵题签 + 大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="sw-kicker">{{进度，如 第二拍 · 配方与实测}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="sw-pill sw-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="sw-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, sw-kicker, h1, mt-m, lede, mt-l, row, sw-pill, sw-pill-accent, deck-footer, slide-number, notes

---

## moments（档期时间线）
指纹：chart
数量：sw-tl-item=4

用途：3-4 个节点的横向时间线：白圈圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="sw-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="sw-tl mt-l" style="margin-top:52px">
    <div class="sw-tl-item"><div class="sw-tl-dot">壹</div><div class="sw-tl-t">{{时间点，≤10 字符}}</div><p class="sw-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sw-tl-item"><div class="sw-tl-dot">贰</div><div class="sw-tl-t">{{时间点，≤10 字符}}</div><p class="sw-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sw-tl-item"><div class="sw-tl-dot">叁</div><div class="sw-tl-t">{{时间点，≤10 字符}}</div><p class="sw-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sw-tl-item"><div class="sw-tl-dot">肆</div><div class="sw-tl-t">{{时间点，≤10 字符}}</div><p class="sw-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#4E7E96;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sw-kicker, h2, mt-m, sw-tl, mt-l, sw-tl-item, sw-tl-dot, sw-tl-t, sw-tl-d, deck-footer, slide-number, notes

---

## closing（收尾太阳）
指纹：hero

用途：收尾页。大字标题 + 一句行动提醒 + 深蓝胶囊按钮 + 药丸 + 小太阳印。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="sw-kicker">{{提醒语境，如 企划评审}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="sw-btn">{{按钮文案，≤8 字}}</span>
    <span class="sw-pill">{{次级信息，≤10 字}}</span>
    <div class="sw-sun">{{太阳印 2 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, sw-kicker, h1, mt-m, lede, mt-l, row, sw-btn, sw-pill, sw-sun, deck-footer, slide-number, notes
