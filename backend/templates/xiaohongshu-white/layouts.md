# 小红书白 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `xh-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-xiaohongshu-white` 作用域生效，骨架里已写全，照抄结构即可。
> 右上暖红书签飘带与晨光渐变由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（暖红是点缀，衬线是灵魂）**：暖红 #F56565 只允许出现在书签（ambient 自带）、
> 贴纸签（xh-tag）、关键数据（xh-stat-v / xh-tl-t / xh-tl-dot）、强调签（xh-pill-accent）与
> 话题标签（xh-hash）；禁大面积红底、禁给整卡上红。标题一律衬线（h1/h2/h4 与 xh-t 已内建
> 宋体栈），禁无衬线粗体做标题。每页最多一枚实底贴纸签（xh-tag），日期签用虚线款
> （xh-tag-soft），活泼感靠贴纸与话题标签，不靠 emoji。

---

## cover（笔记封面）
指纹：hero

用途：开场页。贴纸标签行 + 衬线大标题 + 带荧光划重点的定位句 + 话题标签行，右上书签衬底。
适用 role：cover。
内容约束：主标题 ≤12 字；lede 30-60 字；贴纸签 ≤6 字；日期签 ≤12 字；话题 3 个。

```html
<section class="slide full" data-layout="cover">
  <div class="row" style="gap:12px"><span class="xh-tag">{{贴纸签，≤6 字}}</span><span class="xh-tag-soft">{{日期或期数，≤12 字}}</span></div>
  <h1 class="h1 mt-m">{{主标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:46ch">{{一句话定位，30-60 字，可嵌 <span class="xh-mark">划重点，≤8 字</span>}}</p>
  <div class="row mt-l" style="gap:18px"><span class="xh-hash">#{{话题 1，≤8 字}}</span><span class="xh-hash">#{{话题 2，≤8 字}}</span><span class="xh-hash">#{{话题 3，≤8 字}}</span></div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, row, xh-tag, xh-tag-soft, h1, mt-m, lede, xh-mark, mt-l, xh-hash, deck-footer, slide-number, notes

---

## contents（笔记目录）
指纹：table
数量：xh-entry=4

用途：议程页。一块白卡里放 4 行篇目：圆角红号 + 衬线篇名 + 一句说明，晨粉细线分隔。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="xh-kicker">{{引导语，如 复盘四讲}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="xh-card mt-l" style="padding:18px 44px">
    <div class="xh-entry"><span class="xh-n">01</span><span class="xh-t">{{篇名，≤8 字}}</span><span class="xh-d">{{说明，14-26 字}}</span></div>
    <div class="xh-entry"><span class="xh-n">02</span><span class="xh-t">{{篇名，≤8 字}}</span><span class="xh-d">{{说明，14-26 字}}</span></div>
    <div class="xh-entry"><span class="xh-n">03</span><span class="xh-t">{{篇名，≤8 字}}</span><span class="xh-d">{{说明，14-26 字}}</span></div>
    <div class="xh-entry"><span class="xh-n">04</span><span class="xh-t">{{篇名，≤8 字}}</span><span class="xh-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, xh-kicker, h2, mt-m, xh-card, mt-l, xh-entry, xh-n, xh-t, xh-d, deck-footer, slide-number, notes

---

## keynotes（三卡笔记）
指纹：cards
数量：xh-card=3

用途：恰好三张白卡（珊瑚左条，中间可换琥珀暖卡）。每张：贴纸签 + 衬线小标题 + 说明 + 收藏行。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字；收藏行数字需与口径一致。

```html
<section class="slide" data-layout="keynotes">
  <p class="xh-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:30px;margin-top:44px">
    <div class="xh-card"><span class="xh-tag">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#78716C">{{说明：一句论断 + 一句展开，22-44 字}}</p><div class="row mt-m" style="gap:9px;margin-top:auto"><span class="xh-heart"></span><span class="xh-like">{{收藏数，如 2.1w 收藏}}</span></div></div>
    <div class="xh-card xh-card-warm"><span class="xh-tag">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#78716C">{{说明：一句论断 + 一句展开，22-44 字}}</p><div class="row mt-m" style="gap:9px;margin-top:auto"><span class="xh-heart"></span><span class="xh-like">{{收藏数，如 9.6k 收藏}}</span></div></div>
    <div class="xh-card"><span class="xh-tag">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#78716C">{{说明：一句论断 + 一句展开，22-44 字}}</p><div class="row mt-m" style="gap:9px;margin-top:auto"><span class="xh-heart"></span><span class="xh-like">{{收藏数，如 1.3w 收藏}}</span></div></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, xh-kicker, h2, mt-m, grid, g3, mt-l, xh-card, xh-card-warm, xh-tag, h4, mt-s, row, xh-heart, xh-like, deck-footer, slide-number, notes

---

## split（左文右卡）
指纹：split
数量：xh-item=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边白卡装三行步骤或条目（珊瑚圆点）。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="xh-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#78716C">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:14px">
        <span class="xh-pill">{{要点 1，≤8 字}}</span>
        <span class="xh-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="xh-card" style="min-height:0">
      <div class="xh-item"><span class="xh-dot"></span><p class="xh-mini-t">{{一步，12-26 字}}</p></div>
      <div class="xh-item"><span class="xh-dot"></span><p class="xh-mini-t">{{一步，12-26 字}}</p></div>
      <div class="xh-item"><span class="xh-dot"></span><p class="xh-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, xh-kicker, h2, mt-m, grid, g2, mt-l, lede, row, xh-pill, xh-card, xh-item, xh-dot, xh-mini-t, deck-footer, slide-number, notes

---

## metrics（后台数字）
指纹：chart
数量：xh-stat=3

用途：三个关键数据。珊瑚衬线大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="xh-kicker">{{数据语境，如 后台数据 · 2026 年 3-6 月}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:48px">
    <div class="xh-stat"><div class="xh-stat-v">{{数值 ≤6 字符}}<span class="xh-stat-u">{{单位}}</span></div><div class="xh-stat-l">{{指标名，≤8 字}}</div><p class="xh-stat-note">{{口径，14-30 字}}</p></div>
    <div class="xh-stat"><div class="xh-stat-v">{{数值 ≤6 字符}}<span class="xh-stat-u">{{单位}}</span></div><div class="xh-stat-l">{{指标名，≤8 字}}</div><p class="xh-stat-note">{{口径，14-30 字}}</p></div>
    <div class="xh-stat"><div class="xh-stat-v">{{数值 ≤6 字符}}<span class="xh-stat-u">{{单位}}</span></div><div class="xh-stat-l">{{指标名，≤8 字}}</div><p class="xh-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="margin-top:40px;font-size:17px;color:#A8A29E;letter-spacing:.04em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, xh-kicker, h2, mt-m, grid, g3, mt-l, xh-stat, xh-stat-v, xh-stat-u, xh-stat-l, xh-stat-note, deck-footer, slide-number, notes

---

## quote（扉页引文）
指纹：quote

用途：整页一句引文。衬线大字 + 出处 + 两个支撑药丸 + 话题标签行收底。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="xh-kicker">{{语境，如 复盘笔记的扉页}}</p>
  <p class="xh-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="xh-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:14px">
    <span class="xh-pill xh-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="xh-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="row mt-l" style="gap:18px"><span class="xh-hash">#{{话题 1，≤8 字}}</span><span class="xh-hash">#{{话题 2，≤8 字}}</span></div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, xh-kicker, xh-quote, mt-l, xh-src, mt-m, row, xh-pill, xh-pill-accent, xh-hash, deck-footer, slide-number, notes

---

## divider（章节翻页）
指纹：hero

用途：章节过渡。贴纸标签行（章号 + 章名）+ 衬线大字 + 一个过渡问题 + 两个药丸。
适用 role：divider。
内容约束：章号 ≤6 字；章名 ≤8 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <div class="row" style="gap:12px"><span class="xh-tag">{{章号，≤6 字}}</span><span class="xh-tag-soft">{{章名，≤8 字}}</span></div>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:46ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:14px">
    <span class="xh-pill xh-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="xh-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, row, xh-tag, xh-tag-soft, h1, mt-m, lede, mt-l, xh-pill, xh-pill-accent, deck-footer, slide-number, notes

---

## moments（阶段时间线）
指纹：chart
数量：xh-tl-item=4

用途：3-4 个节点的横向时间线：珊瑚圆角方块节点 + 阶段区间 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="xh-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="xh-tl mt-l" style="margin-top:52px">
    <div class="xh-tl-item"><div class="xh-tl-dot">01</div><div class="xh-tl-t">{{时间点，≤10 字符}}</div><p class="xh-tl-d">{{事件，12-26 字}}</p></div>
    <div class="xh-tl-item"><div class="xh-tl-dot">02</div><div class="xh-tl-t">{{时间点，≤10 字符}}</div><p class="xh-tl-d">{{事件，12-26 字}}</p></div>
    <div class="xh-tl-item"><div class="xh-tl-dot">03</div><div class="xh-tl-t">{{时间点，≤10 字符}}</div><p class="xh-tl-d">{{事件，12-26 字}}</p></div>
    <div class="xh-tl-item"><div class="xh-tl-dot">04</div><div class="xh-tl-t">{{时间点，≤10 字符}}</div><p class="xh-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="margin-top:44px;font-size:17px;color:#A8A29E;letter-spacing:.04em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, xh-kicker, h2, mt-m, xh-tl, mt-l, xh-tl-item, xh-tl-dot, xh-tl-t, xh-tl-d, deck-footer, slide-number, notes

---

## closing（互动收尾）
指纹：hero

用途：收尾页。衬线大字 + 一句行动提醒 + 珊瑚胶囊按钮 + 话题标签，像笔记的评论区置顶。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；话题各 ≤8 字。

```html
<section class="slide full" data-layout="closing">
  <div class="row" style="gap:12px"><span class="xh-tag">{{贴纸签，≤6 字}}</span><span class="xh-tag-soft">{{时间或方式，≤12 字}}</span></div>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:46ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="xh-btn">{{按钮文案，≤8 字}}</span>
    <span class="xh-hash">#{{话题 1，≤8 字}}</span>
    <span class="xh-hash">#{{话题 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, row, xh-tag, xh-tag-soft, h1, mt-m, lede, mt-l, xh-btn, xh-hash, deck-footer, slide-number, notes
