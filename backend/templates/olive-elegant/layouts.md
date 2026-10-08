# 橄榄奶白·高级稳重 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `oe-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-olive-elegant` 作用域生效，骨架里已写全，照抄结构即可。
> 每页的金色细框与页角橄榄枝由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（克制即仪式）**：金色只出现在细处——菱形题签、oe-rule 金线、oe-n/oe-tl-dot
> 描边、oe-btn 描边、数字单位与 oe-pill-accent；奶白是唯一的亮色面，奶白衬线大数字是
> 数据页唯一主视觉。禁大面积金色、禁白色卡、禁第三种亮色；页面底部是橄榄渐变，
> 深色文字禁用（正文一律奶白系）。h1/h2 每页最多一组；oe-en 英文题记一页至多一条。

---

## cover（庄园封面）
指纹：hero

用途：开场页。菱形题签 + 衬线大字标题 + 金线 + 英文题记 + 一句定位。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；英文题记 ≤22 字符。

```html
<section class="slide full" data-layout="cover">
  <p class="oe-kicker">{{题签，≤14 字，如 初榨之地 · 新季发布}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="oe-rule mt-s"></div>
  <p class="oe-en mt-m">{{英文题记，≤22 字符，全大写}}</p>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:20px">
    <span class="oe-pill oe-pill-accent">{{时间地点，≤12 字}}</span>
    <span class="oe-pill">{{场合，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{品牌 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, oe-kicker, h1, mt-m, oe-rule, mt-s, oe-en, lede, mt-l, row, oe-pill, oe-pill-accent, deck-footer, slide-number, notes

---

## contents（发布议程）
指纹：table
数量：oe-item=4

用途：议程页。一块信纸大卡里放 4 行篇章：金圈罗马数字 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="oe-kicker">{{引导语，如 发布议程}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="oe-card mt-l" style="margin-top:44px">
    <div class="oe-item"><span class="oe-n">I</span><span class="oe-t">{{篇名，≤8 字}}</span><span class="oe-d">{{说明，14-26 字}}</span></div>
    <div class="oe-item"><span class="oe-n">II</span><span class="oe-t">{{篇名，≤8 字}}</span><span class="oe-d">{{说明，14-26 字}}</span></div>
    <div class="oe-item"><span class="oe-n">III</span><span class="oe-t">{{篇名，≤8 字}}</span><span class="oe-d">{{说明，14-26 字}}</span></div>
    <div class="oe-item"><span class="oe-n">IV</span><span class="oe-t">{{篇名，≤8 字}}</span><span class="oe-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{品牌 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, oe-kicker, h2, mt-m, oe-card, mt-l, oe-item, oe-n, oe-t, oe-d, deck-footer, slide-number, notes

---

## keynotes（三卡承诺）
指纹：cards
数量：oe-card=3

用途：恰好三张信纸卡。每张：金色题签 + 衬线小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="oe-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="oe-card"><span class="oe-pill oe-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#CCC5AE">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="oe-card"><span class="oe-pill oe-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#CCC5AE">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="oe-card"><span class="oe-pill oe-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#CCC5AE">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{品牌 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, oe-kicker, h2, mt-m, grid, g3, mt-l, oe-card, oe-pill, oe-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：oe-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边信纸卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="oe-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#CCC5AE">{{补充：判断标准或注意事项，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="oe-pill">{{要点 1，≤8 字}}</span>
        <span class="oe-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="oe-card">
      <div class="oe-step"><span class="oe-n">一</span><p class="oe-mini-t">{{一步，12-26 字}}</p></div>
      <div class="oe-step"><span class="oe-n">二</span><p class="oe-mini-t">{{一步，12-26 字}}</p></div>
      <div class="oe-step"><span class="oe-n">三</span><p class="oe-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{品牌 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, oe-kicker, h2, mt-m, grid, g2, mt-l, lede, row, oe-pill, oe-card, oe-step, oe-n, oe-mini-t, deck-footer, slide-number, notes

---

## metrics（庄园数字）
指纹：chart
数量：oe-stat=3

用途：三个关键数据。奶白衬线大数字是唯一主视觉，口径写进块内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="oe-kicker">{{数据语境，如 新季数据}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="oe-stat"><div class="oe-stat-v">{{数值 ≤6 字符}}<span class="oe-stat-u">{{单位}}</span></div><div class="oe-stat-l">{{指标名，≤8 字}}</div><p class="oe-stat-note">{{口径，14-30 字}}</p></div>
    <div class="oe-stat"><div class="oe-stat-v">{{数值 ≤6 字符}}<span class="oe-stat-u">{{单位}}</span></div><div class="oe-stat-l">{{指标名，≤8 字}}</div><p class="oe-stat-note">{{口径，14-30 字}}</p></div>
    <div class="oe-stat"><div class="oe-stat-v">{{数值 ≤6 字符}}<span class="oe-stat-u">{{单位}}</span></div><div class="oe-stat-l">{{指标名，≤8 字}}</div><p class="oe-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#B3AA94;letter-spacing:.06em">来源：{{出处与检测口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{品牌 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, oe-kicker, h2, mt-m, grid, g3, mt-l, oe-stat, oe-stat-v, oe-stat-u, oe-stat-l, oe-stat-note, deck-footer, slide-number, notes

---

## quote（产地引文）
指纹：quote

用途：整页一句引文。衬线大字 + 出处 + 两个金描边药丸，右侧留白可立英文题记。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="oe-kicker">{{语境，如 产地谚语}}</p>
  <p class="oe-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="oe-src mt-m">—— {{出处：谚语来源或品牌出处，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="oe-pill oe-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="oe-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <p class="oe-en" style="position:absolute;right:120px;top:50%;transform:translateY(-50%)">{{英文点缀，≤16 字符}}</p>
  <div class="deck-footer"><span>{{品牌 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, oe-kicker, oe-quote, mt-l, oe-src, mt-m, row, oe-pill, oe-pill-accent, oe-en, deck-footer, slide-number, notes

---

## divider（篇章幕）
指纹：hero

用途：篇章过渡。菱形进度题签 + 衬线大字篇章名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：篇章名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="oe-kicker">{{进度，如 第二章 · 酸度与工艺}}</p>
  <h1 class="h1 mt-m">{{篇章标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="oe-pill oe-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="oe-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{品牌 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, oe-kicker, h1, mt-m, lede, mt-l, row, oe-pill, oe-pill-accent, deck-footer, slide-number, notes

---

## moments（流程时间线）
指纹：chart
数量：oe-tl-item=4

用途：3-4 个节点的横向时间线：金圈圆点 + 时间点 + 一句安排。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="oe-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="oe-tl mt-l" style="margin-top:52px">
    <div class="oe-tl-item"><div class="oe-tl-dot">I</div><div class="oe-tl-t">{{时间点，≤10 字符}}</div><p class="oe-tl-d">{{事件，12-26 字}}</p></div>
    <div class="oe-tl-item"><div class="oe-tl-dot">II</div><div class="oe-tl-t">{{时间点，≤10 字符}}</div><p class="oe-tl-d">{{事件，12-26 字}}</p></div>
    <div class="oe-tl-item"><div class="oe-tl-dot">III</div><div class="oe-tl-t">{{时间点，≤10 字符}}</div><p class="oe-tl-d">{{事件，12-26 字}}</p></div>
    <div class="oe-tl-item"><div class="oe-tl-dot">IV</div><div class="oe-tl-t">{{时间点，≤10 字符}}</div><p class="oe-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#B3AA94;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{品牌 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, oe-kicker, h2, mt-m, oe-tl, mt-l, oe-tl-item, oe-tl-dot, oe-tl-t, oe-tl-d, deck-footer, slide-number, notes

---

## closing（首发收尾）
指纹：hero

用途：收尾页。衬线大字 + 一句行动提醒 + 金框按钮 + 药丸，如请柬落款。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="oe-kicker">{{提醒语境，如 首发预订}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:20px">
    <span class="oe-btn">{{按钮文案，≤8 字}}</span>
    <span class="oe-pill">{{次级信息，≤10 字}}</span>
    <span class="oe-pill oe-pill-accent">{{限量信息，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{品牌 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, oe-kicker, h1, mt-m, lede, mt-l, row, oe-btn, oe-pill, oe-pill-accent, deck-footer, slide-number, notes
