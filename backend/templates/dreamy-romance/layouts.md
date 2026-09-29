# 梦幻浪漫·治愈系 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `dm-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-dreamy-romance` 作用域生效，骨架里已写全，照抄结构即可。
> 每页的三团柔光、月亮与页底花瓣星尘由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（柔雾是底色，金粉是身份）**：香槟金只允许出现在三处——编号金环（dm-n / dm-tl-dot）、
> 题签星点（dm-kicker 自带）、行动钮渐变（dm-btn）。樱粉只做关键数据（dm-stat-v）、时间点
> （dm-tl-t）与樱粉题签（dm-pill-accent）。正文一律深紫灰实色；禁深色背景、禁高饱和霓虹、
> 禁硬阴影。宋体大字（h1/h2）每页最多一组；爱心（dm-heart）一页至多一枚，只放留白处。

---

## cover（云朵封面）
指纹：hero

用途：开场页。柔雾题签 + 宋体抒情大标题 + 一句定位，爱心印记与药丸同排。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；副题 ≤12 字。

```html
<section class="slide full" data-layout="cover">
  <p class="dm-kicker">{{题签，≤14 字，如 暮云咖啡 · 开业发布}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:24px">
    <div class="dm-heart"></div>
    <span class="dm-pill">{{副题或时间地点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, dm-kicker, h1, mt-m, lede, mt-l, row, dm-heart, dm-pill, deck-footer, slide-number, notes

---

## contents（星愿目录）
指纹：table
数量：dm-item=4

用途：议程页。一块云朵大卡里放 4 行篇目：金环编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="dm-kicker">{{引导语，如 今日愿单}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="dm-card mt-l" style="margin-top:44px">
    <div class="dm-item"><span class="dm-n">壹</span><span class="dm-t">{{篇名，≤8 字}}</span><span class="dm-d">{{说明，14-26 字}}</span></div>
    <div class="dm-item"><span class="dm-n">贰</span><span class="dm-t">{{篇名，≤8 字}}</span><span class="dm-d">{{说明，14-26 字}}</span></div>
    <div class="dm-item"><span class="dm-n">叁</span><span class="dm-t">{{篇名，≤8 字}}</span><span class="dm-d">{{说明，14-26 字}}</span></div>
    <div class="dm-item"><span class="dm-n">肆</span><span class="dm-t">{{篇名，≤8 字}}</span><span class="dm-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, dm-kicker, h2, mt-m, dm-card, mt-l, dm-item, dm-n, dm-t, dm-d, deck-footer, slide-number, notes

---

## keynotes（三朵云要点）
指纹：cards
数量：dm-card=3

用途：恰好三张云朵卡。每张：樱粉题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="dm-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="dm-card"><span class="dm-pill dm-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#746599">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="dm-card"><span class="dm-pill dm-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#746599">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="dm-card"><span class="dm-pill dm-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#746599">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, dm-kicker, h2, mt-m, grid, g3, mt-l, dm-card, dm-pill, dm-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右云）
指纹：split
数量：dm-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边云朵卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="dm-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#746599">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="dm-pill">{{要点 1，≤8 字}}</span>
        <span class="dm-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="dm-card">
      <div class="dm-step"><span class="dm-n">一</span><p class="dm-mini-t">{{一步，12-26 字}}</p></div>
      <div class="dm-step"><span class="dm-n">二</span><p class="dm-mini-t">{{一步，12-26 字}}</p></div>
      <div class="dm-step"><span class="dm-n">三</span><p class="dm-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, dm-kicker, h2, mt-m, grid, g2, mt-l, lede, row, dm-pill, dm-card, dm-step, dm-n, dm-mini-t, deck-footer, slide-number, notes

---

## metrics（柔光数字）
指纹：chart
数量：dm-stat=3

用途：三个关键数据。樱粉大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="dm-kicker">{{数据语境，如 试营业两周 · 复盘}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="dm-stat"><div class="dm-stat-v">{{数值 ≤6 字符}}<span class="dm-stat-u">{{单位}}</span></div><div class="dm-stat-l">{{指标名，≤8 字}}</div><p class="dm-stat-note">{{口径，14-30 字}}</p></div>
    <div class="dm-stat"><div class="dm-stat-v">{{数值 ≤6 字符}}<span class="dm-stat-u">{{单位}}</span></div><div class="dm-stat-l">{{指标名，≤8 字}}</div><p class="dm-stat-note">{{口径，14-30 字}}</p></div>
    <div class="dm-stat"><div class="dm-stat-v">{{数值 ≤6 字符}}<span class="dm-stat-u">{{单位}}</span></div><div class="dm-stat-l">{{指标名，≤8 字}}</div><p class="dm-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="margin-top:44px;font-size:17px;color:#9387B5;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, dm-kicker, h2, mt-m, grid, g3, mt-l, dm-stat, dm-stat-v, dm-stat-u, dm-stat-l, dm-stat-note, deck-footer, slide-number, notes

---

## quote（抒情引文）
指纹：quote

用途：整页一句引文。宋体大字 + 出处 + 两个支撑药丸，留白处一枚爱心。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="dm-kicker">{{语境，如 主理人的话}}</p>
  <p class="dm-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="dm-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="dm-pill dm-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="dm-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="dm-heart" style="position:absolute;right:150px;top:130px"></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, dm-kicker, dm-quote, mt-l, dm-src, mt-m, row, dm-pill, dm-pill-accent, dm-heart, deck-footer, slide-number, notes

---

## divider（章节月幕）
指纹：hero

用途：章节过渡。进度题签 + 抒情章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="dm-kicker">{{进度，如 卷二 · 一杯云朵}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="dm-pill dm-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="dm-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, dm-kicker, h1, mt-m, lede, mt-l, row, dm-pill, dm-pill-accent, deck-footer, slide-number, notes

---

## moments（星轨时间线）
指纹：chart
数量：dm-tl-item=4

用途：3-4 个节点的横向时间线：金环圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="dm-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="dm-tl mt-l" style="margin-top:52px">
    <div class="dm-tl-item"><div class="dm-tl-dot">壹</div><div class="dm-tl-t">{{时间点，≤10 字符}}</div><p class="dm-tl-d">{{事件，12-26 字}}</p></div>
    <div class="dm-tl-item"><div class="dm-tl-dot">贰</div><div class="dm-tl-t">{{时间点，≤10 字符}}</div><p class="dm-tl-d">{{事件，12-26 字}}</p></div>
    <div class="dm-tl-item"><div class="dm-tl-dot">叁</div><div class="dm-tl-t">{{时间点，≤10 字符}}</div><p class="dm-tl-d">{{事件，12-26 字}}</p></div>
    <div class="dm-tl-item"><div class="dm-tl-dot">肆</div><div class="dm-tl-t">{{时间点，≤10 字符}}</div><p class="dm-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="margin-top:44px;font-size:17px;color:#9387B5;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, dm-kicker, h2, mt-m, dm-tl, mt-l, dm-tl-item, dm-tl-dot, dm-tl-t, dm-tl-d, deck-footer, slide-number, notes

---

## closing（晚风收尾）
指纹：hero

用途：收尾页。抒情大字 + 一句行动提醒 + 香槟金按钮 + 描边药丸 + 爱心。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="dm-kicker">{{提醒语境，如 开业邀请}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:24px">
    <span class="dm-btn">{{按钮文案，≤8 字}}</span>
    <span class="dm-pill">{{次级信息，≤10 字}}</span>
    <div class="dm-heart"></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, dm-kicker, h1, mt-m, lede, mt-l, row, dm-btn, dm-pill, dm-heart, deck-footer, slide-number, notes
