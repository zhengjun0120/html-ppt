# 手绘秋日旅行手账 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `ha-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-hand-drawn-autumn` 作用域生效，骨架里已写全，照抄结构即可。
> 每页左上枫叶与页底小屋虚线路由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（贴纸手账感）**：暖棕 #8B5E3C 是唯一勾线笔触（边框 2-2.5px 实线或虚线），
> 卡片一律半透明白底 + 不规则圆角，元素自带微旋转不再额外加 transform；南瓜橙 #E87D3E
> 只做焦点（大数字、标签、按钮），砖红 #C04851 只做强调标签与按钮，蓝 #6BA3BE 每页至多一处。
> 大标题（h1/h2）每页最多一组；胶带 ha-tape 一页至多一条。

---

## cover（手账封面）
指纹：hero

用途：开场页。虚线胶囊题签 + 圆体大标题 + 胶带 + 一句定位，下面贴一枚异形贴纸。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；贴纸 2-4 字。

```html
<section class="slide full" data-layout="cover">
  <p class="ha-kicker">{{题签，≤14 字，如 京郊一日 · 骑行手账}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="ha-tape mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="ha-sticker">{{贴纸 2-4 字}}</div>
    <span class="ha-pill">{{时间或天气，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, ha-kicker, h1, mt-m, ha-tape, mt-s, lede, mt-l, row, ha-sticker, ha-pill, deck-footer, slide-number, notes

---

## contents（行程目录）
指纹：table
数量：ha-item=4

用途：议程页。一张便签大卡装 4 行行程：手绘圈编号 + 站名 + 一句说明，卡顶压一条胶带。
适用 role：toc。
内容约束：恰好 4 行；站名 ≤8 字；说明 12-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="ha-kicker">{{引导语，如 今日路线}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="mt-l" style="margin-top:52px">
    <div class="ha-tape" style="margin-bottom:-15px;position:relative;z-index:1"></div>
    <div class="ha-card">
      <div class="ha-item"><span class="ha-n">壹</span><span class="ha-t">{{站名，≤8 字}}</span><span class="ha-d">{{说明，12-26 字}}</span></div>
      <div class="ha-item"><span class="ha-n">贰</span><span class="ha-t">{{站名，≤8 字}}</span><span class="ha-d">{{说明，12-26 字}}</span></div>
      <div class="ha-item"><span class="ha-n">叁</span><span class="ha-t">{{站名，≤8 字}}</span><span class="ha-d">{{说明，12-26 字}}</span></div>
      <div class="ha-item"><span class="ha-n">肆</span><span class="ha-t">{{站名，≤8 字}}</span><span class="ha-d">{{说明，12-26 字}}</span></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ha-kicker, h2, mt-m, mt-l, ha-tape, ha-card, ha-item, ha-n, ha-t, ha-d, deck-footer, slide-number, notes

---

## keynotes（三张贴纸）
指纹：cards
数量：ha-card=3

用途：恰好三张贴纸卡。每张：彩色标签 + 小标题 + 两句说明，标签按橙/砖红/蓝轮换。
适用 role：content。
内容约束：恰好 3 卡；标签 ≤6 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="ha-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="ha-card"><span class="ha-tag">{{标签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#6B543E">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="ha-card"><span class="ha-tag-r">{{标签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#6B543E">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="ha-card"><span class="ha-tag-b">{{标签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#6B543E">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ha-kicker, h2, mt-m, grid, g3, mt-l, ha-card, ha-tag, ha-tag-r, ha-tag-b, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：ha-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边便签卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="ha-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#6B543E">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="ha-pill">{{要点 1，≤8 字}}</span>
        <span class="ha-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="ha-card">
      <div class="ha-step"><span class="ha-n">一</span><p class="ha-mini-t">{{一步，12-26 字}}</p></div>
      <div class="ha-step"><span class="ha-n">二</span><p class="ha-mini-t">{{一步，12-26 字}}</p></div>
      <div class="ha-step"><span class="ha-n">三</span><p class="ha-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ha-kicker, h2, mt-m, grid, g2, mt-l, lede, row, ha-pill, ha-card, ha-step, ha-n, ha-mini-t, deck-footer, slide-number, notes

---

## metrics（南瓜大数字）
指纹：chart
数量：ha-stat=3

用途：三个关键数据。橙色大数字是主视觉，口径写进贴纸卡，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="ha-kicker">{{数据语境，如 里程表 · 当日实骑}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:40px;margin-top:48px">
    <div class="ha-stat"><div class="ha-stat-v">{{数值 ≤6 字符}}<span class="ha-stat-u">{{单位}}</span></div><div class="ha-stat-l">{{指标名，≤8 字}}</div><p class="ha-stat-note">{{口径，14-30 字}}</p></div>
    <div class="ha-stat"><div class="ha-stat-v">{{数值 ≤6 字符}}<span class="ha-stat-u">{{单位}}</span></div><div class="ha-stat-l">{{指标名，≤8 字}}</div><p class="ha-stat-note">{{口径，14-30 字}}</p></div>
    <div class="ha-stat"><div class="ha-stat-v">{{数值 ≤6 字符}}<span class="ha-stat-u">{{单位}}</span></div><div class="ha-stat-l">{{指标名，≤8 字}}</div><p class="ha-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#9A7E5F;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ha-kicker, h2, mt-m, grid, g3, mt-l, ha-stat, ha-stat-v, ha-stat-u, ha-stat-l, ha-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（扉页引文）
指纹：quote

用途：整页一句引文。圆体大字 + 出处 + 两个支撑药丸，右侧贴一张竖排便利贴。
适用 role：quote。
内容约束：引文 10-24 字；出处 ≤22 字且真实；药丸各 ≤10 字；便利贴 5-7 字。

```html
<section class="slide" data-layout="quote">
  <p class="ha-kicker">{{语境，如 手账扉页}}</p>
  <p class="ha-quote mt-l" style="margin-top:40px">「{{引文，10-24 字}}」</p>
  <p class="ha-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="ha-pill ha-pill-accent">{{支撑点 1，≤10 字}}</span>
    <span class="ha-pill">{{支撑点 2，≤10 字}}</span>
  </div>
  <div class="ha-note">{{竖排短句 5-7 字}}</div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ha-kicker, ha-quote, mt-l, ha-src, mt-m, row, ha-pill, ha-pill-accent, ha-note, deck-footer, slide-number, notes

---

## divider（章节站牌）
指纹：hero

用途：章节过渡。进度题签 + 圆体大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤10 字。

```html
<section class="slide full" data-layout="divider">
  <p class="ha-kicker">{{进度，如 第二程 · 枫林与湖}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="ha-pill ha-pill-accent">{{看点 1，≤10 字}}</span>
    <span class="ha-pill">{{看点 2，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, ha-kicker, h1, mt-m, lede, mt-l, row, ha-pill, ha-pill-accent, deck-footer, slide-number, notes

---

## moments（时刻路线）
指纹：chart
数量：ha-tl-item=4

用途：3-4 个节点的横向时间线：手绘圈圆点 + 时间点 + 一句事件，虚线连成当日路线。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="ha-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="ha-tl mt-l" style="margin-top:52px">
    <div class="ha-tl-item"><div class="ha-tl-dot">壹</div><div class="ha-tl-t">{{时间点，≤10 字符}}</div><p class="ha-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ha-tl-item"><div class="ha-tl-dot">贰</div><div class="ha-tl-t">{{时间点，≤10 字符}}</div><p class="ha-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ha-tl-item"><div class="ha-tl-dot">叁</div><div class="ha-tl-t">{{时间点，≤10 字符}}</div><p class="ha-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ha-tl-item"><div class="ha-tl-dot">肆</div><div class="ha-tl-t">{{时间点，≤10 字符}}</div><p class="ha-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#9A7E5F;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ha-kicker, h2, mt-m, ha-tl, mt-l, ha-tl-item, ha-tl-dot, ha-tl-t, ha-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（收尾待续）
指纹：hero

用途：收尾页。圆体大字 + 一句行动提醒 + 砖红按钮 + 药丸 + 一枚「待续」贴纸。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤12 字。

```html
<section class="slide full" data-layout="closing">
  <p class="ha-kicker">{{提醒语境，如 下一次骑行}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="ha-btn">{{按钮文案，≤8 字}}</span>
    <span class="ha-pill">{{次级信息，≤12 字}}</span>
    <div class="ha-sticker">{{贴纸 2 字，如 待续}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, ha-kicker, h1, mt-m, lede, mt-l, row, ha-btn, ha-pill, ha-sticker, deck-footer, slide-number, notes
