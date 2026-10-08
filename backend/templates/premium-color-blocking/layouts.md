# 高级撞色 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `pc-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-premium-color-blocking` 作用域生效，骨架里已写全，照抄结构即可。
> 每页左缘双色书脊与右上描边大圆由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（两色对峙，一处留白）**：全页主色不许超过三种——珊瑚橙、深墨蓝、奶白纸底；
> 柠黄只出现在 pc-kicker 标签上。撞色靠实色块（pc-card-ink / pc-pill-accent / pc-btn）表达，
> 禁渐变、禁毛玻璃、禁拟物阴影、禁圆角大卡。粗实线（3px 以上）负责分组，不用发丝线。
> h1/h2 超粗大字每页最多一组；大色块一页至多一块，留白是撞色高级感的另一半。

---

## cover（大片封面）
指纹：hero

用途：开场页。柠黄标签 + 超粗大标题 + 斜切橙带 + 一句定位，可配角标药丸。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；标签 ≤14 字；药丸 ≤12 字。

```html
<section class="slide full" data-layout="cover">
  <p class="pc-kicker">{{标签，≤14 字，如 FW26 · 系列企划}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="pc-slash mt-s"></div>
  <p class="lede mt-m" style="max-width:46ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="pc-pill pc-pill-accent">{{属性标签，≤10 字}}</span>
    <span class="pc-pill">{{版本或日期，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, pc-kicker, h1, mt-m, pc-slash, mt-s, lede, mt-m, row, mt-l, pc-pill, pc-pill-accent, deck-footer, slide-number, notes

---

## contents（目录切页）
指纹：table
数量：pc-item=4

用途：议程页。一块锐利白卡里放 4 行章目：橙色大序号 + 章名 + 一句说明，粗实线分组。
适用 role：toc。
内容约束：恰好 4 行；章名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="pc-kicker">{{引导语，如 本案结构}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="pc-card mt-l" style="margin-top:44px">
    <div class="pc-item"><span class="pc-n">01</span><span class="pc-t">{{章名，≤8 字}}</span><span class="pc-d">{{说明，14-26 字}}</span></div>
    <div class="pc-item"><span class="pc-n">02</span><span class="pc-t">{{章名，≤8 字}}</span><span class="pc-d">{{说明，14-26 字}}</span></div>
    <div class="pc-item"><span class="pc-n">03</span><span class="pc-t">{{章名，≤8 字}}</span><span class="pc-d">{{说明，14-26 字}}</span></div>
    <div class="pc-item"><span class="pc-n">04</span><span class="pc-t">{{章名，≤8 字}}</span><span class="pc-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pc-kicker, h2, mt-m, pc-card, mt-l, pc-item, pc-n, pc-t, pc-d, deck-footer, slide-number, notes

---

## keynotes（三块对峙）
指纹：cards
数量：pc-card=3

用途：恰好三张锐利白卡。每张：橙色标签 + 超粗小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；标签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="pc-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="pc-card"><span class="pc-pill pc-pill-accent">{{标签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.75;color:#46557A">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="pc-card"><span class="pc-pill pc-pill-accent">{{标签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.75;color:#46557A">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="pc-card"><span class="pc-pill pc-pill-accent">{{标签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.75;color:#46557A">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pc-kicker, h2, mt-m, grid, g3, mt-l, pc-card, pc-pill, pc-pill-accent, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左文右块）
指纹：split
数量：pc-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边墨蓝实色块装三行步骤，构成撞色对峙。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="pc-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.75;color:#46557A">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="pc-pill">{{要点 1，≤8 字}}</span>
        <span class="pc-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="pc-card-ink">
      <div class="pc-step"><span class="pc-n">一</span><p class="pc-mini-t">{{一步，12-26 字}}</p></div>
      <div class="pc-step"><span class="pc-n">二</span><p class="pc-mini-t">{{一步，12-26 字}}</p></div>
      <div class="pc-step"><span class="pc-n">三</span><p class="pc-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pc-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, mt-l, pc-pill, pc-card-ink, pc-step, pc-n, pc-mini-t, deck-footer, slide-number, notes

---

## metrics（撞色数字）
指纹：chart
数量：pc-stat=3

用途：三个关键数据。橙色超大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="pc-kicker">{{数据语境，如 规模口径}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="pc-stat"><div class="pc-stat-v">{{数值 ≤6 字符}}<span class="pc-stat-u">{{单位}}</span></div><div class="pc-stat-l">{{指标名，≤8 字}}</div><p class="pc-stat-note">{{口径，14-30 字}}</p></div>
    <div class="pc-stat"><div class="pc-stat-v">{{数值 ≤6 字符}}<span class="pc-stat-u">{{单位}}</span></div><div class="pc-stat-l">{{指标名，≤8 字}}</div><p class="pc-stat-note">{{口径，14-30 字}}</p></div>
    <div class="pc-stat"><div class="pc-stat-v">{{数值 ≤6 字符}}<span class="pc-stat-u">{{单位}}</span></div><div class="pc-stat-l">{{指标名，≤8 字}}</div><p class="pc-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:16px;font-weight:600;color:#8B93A8;letter-spacing:.08em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pc-kicker, h2, mt-m, grid, g3, mt-l, pc-stat, pc-stat-v, pc-stat-u, pc-stat-l, pc-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（宣言引文）
指纹：quote

用途：整页一句宣言式引文。超粗大字 + 出处 + 两个支撑药丸，右侧斜切橙带点缀。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="pc-kicker">{{语境，如 视觉手册 · 卷首}}</p>
  <p class="pc-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="pc-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="pc-pill pc-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="pc-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="pc-slash" style="position:absolute;right:140px;top:50%;transform:translateY(-50%) skewX(-24deg);width:180px;height:12px"></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pc-kicker, pc-quote, mt-l, pc-src, mt-m, row, mt-l, pc-pill, pc-pill-accent, pc-slash, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度标签 + 超粗大字章节名 + 一个过渡句 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡句 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="pc-kicker">{{进度，如 第三章 · 造型结构}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:46ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="pc-pill pc-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="pc-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, pc-kicker, h1, mt-m, lede, mt-l, row, pc-pill, pc-pill-accent, deck-footer, slide-number, notes

---

## moments（流程时间线）
指纹：chart
数量：pc-tl-item=4

用途：3-4 个节点的横向时间线：墨蓝方块点 + 时间点 + 一句事件，粗色带连线。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="pc-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="pc-tl mt-l" style="margin-top:52px">
    <div class="pc-tl-item"><div class="pc-tl-dot">01</div><div class="pc-tl-t">{{时间点，≤10 字符}}</div><p class="pc-tl-d">{{事件，12-26 字}}</p></div>
    <div class="pc-tl-item"><div class="pc-tl-dot">02</div><div class="pc-tl-t">{{时间点，≤10 字符}}</div><p class="pc-tl-d">{{事件，12-26 字}}</p></div>
    <div class="pc-tl-item"><div class="pc-tl-dot">03</div><div class="pc-tl-t">{{时间点，≤10 字符}}</div><p class="pc-tl-d">{{事件，12-26 字}}</p></div>
    <div class="pc-tl-item"><div class="pc-tl-dot">04</div><div class="pc-tl-t">{{时间点，≤10 字符}}</div><p class="pc-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:16px;font-weight:600;color:#8B93A8;letter-spacing:.08em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pc-kicker, h2, mt-m, pc-tl, mt-l, pc-tl-item, pc-tl-dot, pc-tl-t, pc-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（收尾订货）
指纹：hero

用途：收尾页。超粗大字 + 一句行动提醒 + 墨蓝实底按钮 + 描边药丸。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="pc-kicker">{{提醒语境，如 订货安排}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:46ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="pc-btn">{{按钮文案，≤8 字}}</span>
    <span class="pc-pill">{{次级信息，≤10 字}}</span>
    <span class="pc-pill pc-pill-accent">{{时限或编号，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, pc-kicker, h1, mt-m, lede, mt-l, row, pc-btn, pc-pill, pc-pill-accent, deck-footer, slide-number, notes
