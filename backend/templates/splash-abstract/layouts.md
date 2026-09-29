# 泼彩抽象·多巴胺 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `sp-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-splash-abstract` 作用域生效，骨架里已写全，照抄结构即可。
> 每页的泼彩色块与星芒由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（一页一个爆点）**：五色（橙红/品红/紫罗兰/青绿/亮蓝）只出现在轮换身份件上——
> 编号块（sp-n）、时间线圆点（sp-tl-dot）、大数字（sp-stat-v）由模板按序自动轮换，
> 不用也不许手写内联颜色；明黄只做题签底（sp-kicker）与爆点星芒（sp-burst）。
> 正文与卡内文字一律深墨实色；白字只出现在五色圆点（特大号）与深墨按钮上。
> 星芒 sp-burst 一页至多一枚；粗黑描边 + 硬阴影是卡与钮的统一画法。

---

## cover（开幕主视觉）
指纹：hero

用途：开场页。明黄题签 + 特粗大标题 + 一句策展定位，星芒与药丸同排。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；副题 ≤12 字。

```html
<section class="slide full" data-layout="cover">
  <p class="sp-kicker">{{题签，≤14 字，如 岸线美术馆 · 秋季特展}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:22px">
    <svg class="sp-burst" width="72" height="72" viewBox="0 0 64 64"><path d="M32,2 L38,24 L58,14 L44,32 L62,38 L42,42 L48,62 L32,46 L16,62 L22,42 L2,38 L20,32 L6,14 L26,24 Z" fill="#FFD700"/></svg>
    <span class="sp-pill">{{副题或展期地点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, sp-kicker, h1, mt-m, lede, mt-l, row, sp-burst, sp-pill, deck-footer, slide-number, notes

---

## contents（展签目录）
指纹：table
数量：sp-item=4

用途：议程页。一块粗描边大卡里放 4 行展区：五色编号块 + 展区名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；展区名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="sp-kicker">{{引导语，如 今日展签}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="sp-card mt-l" style="margin-top:44px">
    <div class="sp-item"><span class="sp-n">01</span><span class="sp-t">{{展区名，≤8 字}}</span><span class="sp-d">{{说明，14-26 字}}</span></div>
    <div class="sp-item"><span class="sp-n">02</span><span class="sp-t">{{展区名，≤8 字}}</span><span class="sp-d">{{说明，14-26 字}}</span></div>
    <div class="sp-item"><span class="sp-n">03</span><span class="sp-t">{{展区名，≤8 字}}</span><span class="sp-d">{{说明，14-26 字}}</span></div>
    <div class="sp-item"><span class="sp-n">04</span><span class="sp-t">{{展区名，≤8 字}}</span><span class="sp-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sp-kicker, h2, mt-m, sp-card, mt-l, sp-item, sp-n, sp-t, sp-d, deck-footer, slide-number, notes

---

## keynotes（三件必看）
指纹：cards
数量：sp-card=3

用途：恰好三张粗描边卡。每张：品红题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="sp-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:36px;margin-top:44px">
    <div class="sp-card"><span class="sp-pill sp-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#4A4A66">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="sp-card"><span class="sp-pill sp-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#4A4A66">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="sp-card"><span class="sp-pill sp-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#4A4A66">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sp-kicker, h2, mt-m, grid, g3, mt-l, sp-card, sp-pill, sp-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右线）
指纹：split
数量：sp-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边粗描边卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="sp-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#4A4A66">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="sp-pill">{{要点 1，≤8 字}}</span>
        <span class="sp-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="sp-card">
      <div class="sp-step"><span class="sp-n">01</span><p class="sp-mini-t">{{一步，12-26 字}}</p></div>
      <div class="sp-step"><span class="sp-n">02</span><p class="sp-mini-t">{{一步，12-26 字}}</p></div>
      <div class="sp-step"><span class="sp-n">03</span><p class="sp-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sp-kicker, h2, mt-m, grid, g2, mt-l, lede, row, sp-pill, sp-card, sp-step, sp-n, sp-mini-t, deck-footer, slide-number, notes

---

## metrics（展观数字）
指纹：chart
数量：sp-stat=3

用途：三个关键数据。五色轮换大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="sp-kicker">{{数据语境，如 展前盘点}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="sp-stat"><div class="sp-stat-v">{{数值 ≤6 字符}}<span class="sp-stat-u">{{单位}}</span></div><div class="sp-stat-l">{{指标名，≤8 字}}</div><p class="sp-stat-note">{{口径，14-30 字}}</p></div>
    <div class="sp-stat"><div class="sp-stat-v">{{数值 ≤6 字符}}<span class="sp-stat-u">{{单位}}</span></div><div class="sp-stat-l">{{指标名，≤8 字}}</div><p class="sp-stat-note">{{口径，14-30 字}}</p></div>
    <div class="sp-stat"><div class="sp-stat-v">{{数值 ≤6 字符}}<span class="sp-stat-u">{{单位}}</span></div><div class="sp-stat-l">{{指标名，≤8 字}}</div><p class="sp-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="margin-top:44px;font-size:17px;color:#77748F;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sp-kicker, h2, mt-m, grid, g3, mt-l, sp-stat, sp-stat-v, sp-stat-u, sp-stat-l, sp-stat-note, deck-footer, slide-number, notes

---

## quote（艺术家引文）
指纹：quote

用途：整页一句引文。特粗大字 + 出处 + 两个支撑药丸，留白处一枚星芒。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="sp-kicker">{{语境，如 布展手记}}</p>
  <p class="sp-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="sp-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="sp-pill sp-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="sp-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <svg class="sp-burst" style="position:absolute;right:170px;top:120px" width="88" height="88" viewBox="0 0 64 64"><path d="M32,2 L38,24 L58,14 L44,32 L62,38 L42,42 L48,62 L32,46 L16,62 L22,42 L2,38 L20,32 L6,14 L26,24 Z" fill="#FFD700"/></svg>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sp-kicker, sp-quote, mt-l, sp-src, mt-m, row, sp-pill, sp-pill-accent, sp-burst, deck-footer, slide-number, notes

---

## divider（章节泼墨）
指纹：hero

用途：章节过渡。进度题签 + 特粗章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="sp-kicker">{{进度，如 卷二 · 为什么是泼彩}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="sp-pill sp-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="sp-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, sp-kicker, h1, mt-m, lede, mt-l, row, sp-pill, sp-pill-accent, deck-footer, slide-number, notes

---

## moments（公共项目时间线）
指纹：chart
数量：sp-tl-item=4

用途：3-4 个节点的横向时间线：五色圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="sp-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="sp-tl mt-l" style="margin-top:52px">
    <div class="sp-tl-item"><div class="sp-tl-dot">01</div><div class="sp-tl-t">{{时间点，≤10 字符}}</div><p class="sp-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sp-tl-item"><div class="sp-tl-dot">02</div><div class="sp-tl-t">{{时间点，≤10 字符}}</div><p class="sp-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sp-tl-item"><div class="sp-tl-dot">03</div><div class="sp-tl-t">{{时间点，≤10 字符}}</div><p class="sp-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sp-tl-item"><div class="sp-tl-dot">04</div><div class="sp-tl-t">{{时间点，≤10 字符}}</div><p class="sp-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="margin-top:44px;font-size:17px;color:#77748F;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sp-kicker, h2, mt-m, sp-tl, mt-l, sp-tl-item, sp-tl-dot, sp-tl-t, sp-tl-d, deck-footer, slide-number, notes

---

## closing（闭幕邀请）
指纹：hero

用途：收尾页。特粗大字 + 一句行动提醒 + 深墨按钮 + 描边药丸 + 星芒。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="sp-kicker">{{提醒语境，如 开幕邀请}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:24px">
    <span class="sp-btn">{{按钮文案，≤8 字}}</span>
    <span class="sp-pill">{{次级信息，≤10 字}}</span>
    <svg class="sp-burst" width="72" height="72" viewBox="0 0 64 64"><path d="M32,2 L38,24 L58,14 L44,32 L62,38 L42,42 L48,62 L32,46 L16,62 L22,42 L2,38 L20,32 L6,14 L26,24 Z" fill="#FFD700"/></svg>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, sp-kicker, h1, mt-m, lede, mt-l, row, sp-btn, sp-pill, sp-burst, deck-footer, slide-number, notes
