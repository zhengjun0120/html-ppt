# 一半星河一半烟火 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `sf-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-starlight-fireworks` 作用域生效，骨架里已写全，照抄结构即可。
> 每页左上星点流星与右下烟火粒子由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（双意象二元）**：星光金只允许出现在四处——星月徽章（sf-moon）、题签细线、
> 关键数据（sf-stat-v / sf-tl-t）、金色药丸与按钮（sf-pill-accent / sf-btn）。烟火橙只活在
> ambient 层，不进正文；禁大面积橙红、禁硬阴影、禁荧光色。衬线大字（h1/h2）每页最多一组；
> 星月徽章一页至多一枚。

---

## cover（星幕开场）
指纹：hero

用途：开场页。星辉题签 + 衬线大字标题 + 星光线 + 一句定位，右侧可立一枚星月徽章。
适用 role：cover。
内容约束：主标题 ≤12 字；lede 24-48 字；题签 ≤14 字；徽章 4 字内。

```html
<section class="slide full" data-layout="cover">
  <p class="sf-kicker">{{题签，≤14 字，如 跨年夜活动策划案 · 双城}}</p>
  <h1 class="h1 mt-m">{{主标题，≤12 字，可 <br> 分行}}</h1>
  <div class="sf-starline mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="sf-moon">{{徽章 2×2 字}}</div>
    <span class="sf-pill">{{副题或时间地点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, sf-kicker, h1, mt-m, sf-starline, mt-s, lede, mt-l, row, sf-moon, sf-pill, deck-footer, slide-number, notes

---

## contents（双城篇目）
指纹：table
数量：sf-item=4

用途：议程页。一块星夜大卡里放 4 行篇目：星圈编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="sf-kicker">{{引导语，如 今晚篇目}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="sf-card mt-l" style="margin-top:44px">
    <div class="sf-item"><span class="sf-n">壹</span><span class="sf-t">{{篇名，≤8 字}}</span><span class="sf-d">{{说明，14-26 字}}</span></div>
    <div class="sf-item"><span class="sf-n">贰</span><span class="sf-t">{{篇名，≤8 字}}</span><span class="sf-d">{{说明，14-26 字}}</span></div>
    <div class="sf-item"><span class="sf-n">叁</span><span class="sf-t">{{篇名，≤8 字}}</span><span class="sf-d">{{说明，14-26 字}}</span></div>
    <div class="sf-item"><span class="sf-n">肆</span><span class="sf-t">{{篇名，≤8 字}}</span><span class="sf-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sf-kicker, h2, mt-m, sf-card, mt-l, sf-item, sf-n, sf-t, sf-d, deck-footer, slide-number, notes

---

## keynotes（三幕要点）
指纹：cards
数量：sf-card=3

用途：恰好三张星夜卡。每张：金色题签 + 衬线小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="sf-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="sf-card"><span class="sf-pill sf-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#C4D3EC">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="sf-card"><span class="sf-pill sf-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#C4D3EC">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="sf-card"><span class="sf-pill sf-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#C4D3EC">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sf-kicker, h2, mt-m, grid, g3, mt-l, sf-card, sf-pill, sf-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（双栏分幕）
指纹：split
数量：sf-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边星夜卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="sf-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#C4D3EC">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="sf-pill">{{要点 1，≤8 字}}</span>
        <span class="sf-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="sf-card">
      <div class="sf-step"><span class="sf-n">一</span><p class="sf-mini-t">{{一步，12-26 字}}</p></div>
      <div class="sf-step"><span class="sf-n">二</span><p class="sf-mini-t">{{一步，12-26 字}}</p></div>
      <div class="sf-step"><span class="sf-n">三</span><p class="sf-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sf-kicker, h2, mt-m, grid, g2, mt-l, lede, row, sf-pill, sf-card, sf-step, sf-n, sf-mini-t, deck-footer, slide-number, notes

---

## metrics（星光数据）
指纹：chart
数量：sf-stat=3

用途：三个关键数据。星光金大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="sf-kicker">{{数据语境，如 去年跨年 · 复盘}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="sf-stat"><div class="sf-stat-v">{{数值 ≤6 字符}}<span class="sf-stat-u">{{单位}}</span></div><div class="sf-stat-l">{{指标名，≤8 字}}</div><p class="sf-stat-note">{{口径，14-30 字}}</p></div>
    <div class="sf-stat"><div class="sf-stat-v">{{数值 ≤6 字符}}<span class="sf-stat-u">{{单位}}</span></div><div class="sf-stat-l">{{指标名，≤8 字}}</div><p class="sf-stat-note">{{口径，14-30 字}}</p></div>
    <div class="sf-stat"><div class="sf-stat-v">{{数值 ≤6 字符}}<span class="sf-stat-u">{{单位}}</span></div><div class="sf-stat-l">{{指标名，≤8 字}}</div><p class="sf-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#8FA3C6;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sf-kicker, h2, mt-m, grid, g3, mt-l, sf-stat, sf-stat-v, sf-stat-u, sf-stat-l, sf-stat-note, deck-footer, slide-number, notes

---

## quote（星夜题记）
指纹：quote

用途：整页一句引文。衬线斜体大字 + 出处 + 两个支撑药丸，右侧留白处一枚星月徽章。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="sf-kicker">{{语境，如 主理人题记}}</p>
  <p class="sf-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="sf-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="sf-pill sf-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="sf-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="sf-moon" style="position:absolute;right:120px;top:50%;transform:translateY(-50%)">{{徽章 2×2 字}}</div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sf-kicker, sf-quote, mt-l, sf-src, mt-m, row, sf-pill, sf-pill-accent, sf-moon, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 衬线大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="sf-kicker">{{进度，如 第二幕 · 焰火之章}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="sf-pill sf-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="sf-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, sf-kicker, h1, mt-m, lede, mt-l, row, sf-pill, sf-pill-accent, deck-footer, slide-number, notes

---

## moments（流程时间线）
指纹：chart
数量：sf-tl-item=4

用途：3-4 个节点的横向时间线：星点圆珠 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="sf-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="sf-tl mt-l" style="margin-top:52px">
    <div class="sf-tl-item"><div class="sf-tl-dot">壹</div><div class="sf-tl-t">{{时间点，≤10 字符}}</div><p class="sf-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sf-tl-item"><div class="sf-tl-dot">贰</div><div class="sf-tl-t">{{时间点，≤10 字符}}</div><p class="sf-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sf-tl-item"><div class="sf-tl-dot">叁</div><div class="sf-tl-t">{{时间点，≤10 字符}}</div><p class="sf-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sf-tl-item"><div class="sf-tl-dot">肆</div><div class="sf-tl-t">{{时间点，≤10 字符}}</div><p class="sf-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#8FA3C6;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sf-kicker, h2, mt-m, sf-tl, mt-l, sf-tl-item, sf-tl-dot, sf-tl-t, sf-tl-d, deck-footer, slide-number, notes

---

## closing（收尾邀约）
指纹：hero

用途：收尾页。衬线大字 + 一句行动提醒 + 金边按钮 + 描边药丸 + 星月徽章，如焰火落幕。
适用 role：thanks / cta / content。
内容约束：标题 ≤12 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="sf-kicker">{{提醒语境，如 观演预约}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="sf-btn">{{按钮文案，≤8 字}}</span>
    <span class="sf-pill">{{次级信息，≤10 字}}</span>
    <div class="sf-moon">{{徽章 2 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, sf-kicker, h1, mt-m, lede, mt-l, row, sf-btn, sf-pill, sf-moon, deck-footer, slide-number, notes
