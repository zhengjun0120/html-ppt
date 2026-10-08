# 米白樱粉 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `cb-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-chinese-cream-blossom` 作用域生效，骨架里已写全，照抄结构即可。
> 每页右上圆瓣花影与左下飘落花瓣由模板自动衬底（z-index 在内容之下），骨架与正文都不用画，
> 也不要再叠加任何花瓣装饰。
>
> **气质铁律（樱粉是温度）**：樱粉只用于柔性元素——花影衬底、樱线 cb-line、卡左线 cb-card、
> 花圈 cb-n、按钮 cb-btn；胭脂只允许出现在五处——小印 cb-seal、强调胶囊 cb-pill-red、
> 关键数据（cb-stat-v / cb-tl-t）、竖排重点字 cb-vert-accent。淡金只做花蕊级小点缀。
> 禁硬阴影、禁锋利几何、禁黑体粗标题。宋楷大字（h1/h2）每页最多一组；
> 竖排 cb-vert 只放留白处，一页至多一条。

---

## cover（花笺封面）
指纹：hero

用途：开场页。题签 + 宋楷大字标题 + 樱线 + 一句定位，花影衬底，可立一枚胭脂小印。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；小印 4 字内。

```html
<section class="slide full" data-layout="cover">
  <p class="cb-kicker">{{题签，≤14 字，如 囍时婚礼 · 二〇二七春}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="cb-line mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="cb-seal">{{小印 2×2 字}}</div>
    <span class="cb-pill">{{副题或时间地点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, cb-kicker, h1, mt-m, cb-line, mt-s, lede, mt-l, row, cb-seal, cb-pill, deck-footer, slide-number, notes

---

## contents（花目）
指纹：table
数量：cb-item=4

用途：议程页。一块左线暖卡里放 4 行条目：花圈编号 + 条目名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；条目名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="cb-kicker">{{引导语，如 方案四幕}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="cb-card mt-l" style="margin-top:44px">
    <div class="cb-item"><span class="cb-n">壹</span><span class="cb-t">{{条目名，≤8 字}}</span><span class="cb-d">{{说明，14-26 字}}</span></div>
    <div class="cb-item"><span class="cb-n">贰</span><span class="cb-t">{{条目名，≤8 字}}</span><span class="cb-d">{{说明，14-26 字}}</span></div>
    <div class="cb-item"><span class="cb-n">叁</span><span class="cb-t">{{条目名，≤8 字}}</span><span class="cb-d">{{说明，14-26 字}}</span></div>
    <div class="cb-item"><span class="cb-n">肆</span><span class="cb-t">{{条目名，≤8 字}}</span><span class="cb-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cb-kicker, h2, mt-m, cb-card, mt-l, cb-item, cb-n, cb-t, cb-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：cb-card=3

用途：恰好三张左线暖卡。每张：胶囊题签 + 楷体小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="cb-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="cb-card"><span class="cb-pill cb-pill-red">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#7A6B5D">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="cb-card"><span class="cb-pill cb-pill-red">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#7A6B5D">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="cb-card"><span class="cb-pill cb-pill-red">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#7A6B5D">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cb-kicker, h2, mt-m, grid, g3, mt-l, cb-card, cb-pill, cb-pill-red, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：cb-step=3

用途：左边把一件事讲透（lede + 补充 + 胶囊），右暖卡装三行流程。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个胶囊；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="cb-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#7A6B5D">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="cb-pill">{{要点 1，≤8 字}}</span>
        <span class="cb-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="cb-card">
      <div class="cb-step"><span class="cb-n">一</span><p class="cb-step-t">{{一步，12-26 字}}</p></div>
      <div class="cb-step"><span class="cb-n">二</span><p class="cb-step-t">{{一步，12-26 字}}</p></div>
      <div class="cb-step"><span class="cb-n">三</span><p class="cb-step-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cb-kicker, h2, mt-m, grid, g2, mt-l, lede, row, cb-pill, cb-card, cb-step, cb-n, cb-step-t, deck-footer, slide-number, notes

---

## metrics（胭脂数字）
指纹：chart
数量：cb-stat=3

用途：三个关键数据。胭脂大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="cb-kicker">{{数据语境，如 方案基准 · 口径}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="cb-stat"><div class="cb-stat-v">{{数值 ≤6 字符}}<span class="cb-stat-u">{{单位}}</span></div><div class="cb-stat-l">{{指标名，≤8 字}}</div><p class="cb-stat-note">{{口径，14-30 字}}</p></div>
    <div class="cb-stat"><div class="cb-stat-v">{{数值 ≤6 字符}}<span class="cb-stat-u">{{单位}}</span></div><div class="cb-stat-l">{{指标名，≤8 字}}</div><p class="cb-stat-note">{{口径，14-30 字}}</p></div>
    <div class="cb-stat"><div class="cb-stat-v">{{数值 ≤6 字符}}<span class="cb-stat-u">{{单位}}</span></div><div class="cb-stat-l">{{指标名，≤8 字}}</div><p class="cb-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#97897B;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cb-kicker, h2, mt-m, grid, g3, mt-l, cb-stat, cb-stat-v, cb-stat-u, cb-stat-l, cb-stat-note, deck-footer, slide-number, notes

---

## quote（花笺引文）
指纹：quote

用途：整页一句引文。宋楷大字 + 出处 + 两个支撑胶囊，右侧留白处可立竖排短句。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；胶囊各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="cb-kicker">{{语境，如 写给二位}}</p>
  <p class="cb-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="cb-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="cb-pill cb-pill-red">{{支撑点 1，≤8 字}}</span>
    <span class="cb-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="cb-vert" style="position:absolute;right:120px;top:50%;transform:translateY(-50%)">{{竖排短句 ≤7 字}}<span class="cb-vert-accent">{{一字或两字}}</span></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cb-kicker, cb-quote, mt-l, cb-src, mt-m, row, cb-pill, cb-pill-red, cb-vert, cb-vert-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 大字章节名 + 一个过渡问题 + 两个看点胶囊。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；胶囊各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="cb-kicker">{{进度，如 第二幕 · 色与花}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="cb-pill cb-pill-red">{{看点 1，≤8 字}}</span>
    <span class="cb-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, cb-kicker, h1, mt-m, lede, mt-l, row, cb-pill, cb-pill-red, deck-footer, slide-number, notes

---

## moments（筹备时间线）
指纹：chart
数量：cb-tl-item=4

用途：3-4 个节点的横向时间线：花圈节点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="cb-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="cb-tl mt-l" style="margin-top:52px">
    <div class="cb-tl-item"><div class="cb-tl-dot">壹</div><div class="cb-tl-t">{{时间点，≤10 字符}}</div><p class="cb-tl-d">{{事件，12-26 字}}</p></div>
    <div class="cb-tl-item"><div class="cb-tl-dot">贰</div><div class="cb-tl-t">{{时间点，≤10 字符}}</div><p class="cb-tl-d">{{事件，12-26 字}}</p></div>
    <div class="cb-tl-item"><div class="cb-tl-dot">叁</div><div class="cb-tl-t">{{时间点，≤10 字符}}</div><p class="cb-tl-d">{{事件，12-26 字}}</p></div>
    <div class="cb-tl-item"><div class="cb-tl-dot">肆</div><div class="cb-tl-t">{{时间点，≤10 字符}}</div><p class="cb-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#97897B;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cb-kicker, h2, mt-m, cb-tl, mt-l, cb-tl-item, cb-tl-dot, cb-tl-t, cb-tl-d, deck-footer, slide-number, notes

---

## closing（收尾相邀）
指纹：hero

用途：收尾页。宋楷大字 + 一句行动提醒 + 樱粉实底按钮 + 胶囊与小印，如花笺落款。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；胶囊 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="cb-kicker">{{提醒语境，如 订档提醒}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="cb-btn">{{按钮文案，≤8 字}}</span>
    <span class="cb-pill">{{次级信息，≤10 字}}</span>
    <div class="cb-seal">{{小印 2 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, cb-kicker, h1, mt-m, lede, mt-l, row, cb-btn, cb-pill, cb-seal, deck-footer, slide-number, notes
