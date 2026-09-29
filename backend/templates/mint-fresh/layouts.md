# 薄荷清新 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `mf-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-mint-fresh` 作用域生效，骨架里已写全，照抄结构即可。
> 每页右上柔光斑与叶脉线稿、左下豆沙绿光斑由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（留白即呼吸）**：薄荷填充只允许三处——题签胶囊（mf-kicker）、编号圆角块
> （mf-n / mf-tl-dot）、数据上边线（mf-stat 顶线）。苔绿只做关键数据（mf-stat-v / mf-tl-t /
> mf-pill-accent），深墨绿只做标题与 mf-btn。禁荧光绿、禁大面积色块铺底、禁硬阴影；
> 卡只有一种（mf-card 半透明大圆角）。衬线轻字 mf-thin 只在 h1 里点一词，一页至多一处。

---

## cover（晨叶封面）
指纹：hero

用途：开场页。薄荷胶囊题签 + 大标题 + 薄荷渐变线 + 一句定位，右侧可立一枚苔绿方印。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；方印 4 字内。

```html
<section class="slide full" data-layout="cover">
  <p class="mf-kicker">{{题签，≤14 字，如 薄荷生活 · 春夏发布}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行，可点一词 <span class="mf-thin">轻字</span>}}</h1>
  <div class="mf-line mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="mf-badge">{{方印 2×2 字}}</div>
    <span class="mf-pill">{{副题或时间地点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, mf-kicker, h1, mt-m, mf-thin, mf-line, mt-s, lede, row, mt-l, mf-badge, mf-pill, deck-footer, slide-number, notes

---

## contents（议程薄卡）
指纹：table
数量：mf-item=4

用途：议程页。一块半透明圆角大卡里放 4 行篇目：薄荷圆角块编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="mf-kicker">{{引导语，如 今日议程}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="mf-card mt-l" style="margin-top:44px">
    <div class="mf-item"><span class="mf-n">01</span><span class="mf-t">{{篇名，≤8 字}}</span><span class="mf-d">{{说明，14-26 字}}</span></div>
    <div class="mf-item"><span class="mf-n">02</span><span class="mf-t">{{篇名，≤8 字}}</span><span class="mf-d">{{说明，14-26 字}}</span></div>
    <div class="mf-item"><span class="mf-n">03</span><span class="mf-t">{{篇名，≤8 字}}</span><span class="mf-d">{{说明，14-26 字}}</span></div>
    <div class="mf-item"><span class="mf-n">04</span><span class="mf-t">{{篇名，≤8 字}}</span><span class="mf-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mf-kicker, h2, mt-m, mf-card, mt-l, mf-item, mf-n, mf-t, mf-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：mf-card=3

用途：恰好三张半透明圆角卡。每张：苔绿题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="mf-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="mf-card"><span class="mf-pill mf-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#5C6B60">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="mf-card"><span class="mf-pill mf-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#5C6B60">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="mf-card"><span class="mf-pill mf-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#5C6B60">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mf-kicker, h2, mt-m, grid, g3, mt-l, mf-card, mf-pill, mf-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：mf-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边圆角卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="mf-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#5C6B60">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="mf-pill">{{要点 1，≤8 字}}</span>
        <span class="mf-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="mf-card">
      <div class="mf-step"><span class="mf-n">一</span><p class="mf-mini-t">{{一步，12-26 字}}</p></div>
      <div class="mf-step"><span class="mf-n">二</span><p class="mf-mini-t">{{一步，12-26 字}}</p></div>
      <div class="mf-step"><span class="mf-n">三</span><p class="mf-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mf-kicker, h2, mt-m, grid, g2, mt-l, lede, row, mf-pill, mf-card, mf-step, mf-n, mf-mini-t, deck-footer, slide-number, notes

---

## metrics（苔绿数字）
指纹：chart
数量：mf-stat=3

用途：三个关键数据。苔绿大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="mf-kicker">{{数据语境，如 内测与实测 · 数据}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="mf-stat"><div class="mf-stat-v">{{数值 ≤6 字符}}<span class="mf-stat-u">{{单位}}</span></div><div class="mf-stat-l">{{指标名，≤8 字}}</div><p class="mf-stat-note">{{口径，14-30 字}}</p></div>
    <div class="mf-stat"><div class="mf-stat-v">{{数值 ≤6 字符}}<span class="mf-stat-u">{{单位}}</span></div><div class="mf-stat-l">{{指标名，≤8 字}}</div><p class="mf-stat-note">{{口径，14-30 字}}</p></div>
    <div class="mf-stat"><div class="mf-stat-v">{{数值 ≤6 字符}}<span class="mf-stat-u">{{单位}}</span></div><div class="mf-stat-l">{{指标名，≤8 字}}</div><p class="mf-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#8CA093;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mf-kicker, h2, mt-m, grid, g3, mt-l, mf-stat, mf-stat-v, mf-stat-u, mf-stat-l, mf-stat-note, deck-footer, slide-number, notes

---

## quote（叶间引文）
指纹：quote

用途：整页一句引文。宋体大字 + 出处 + 两个支撑药丸，右侧留白处立一株蕨叶线稿。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="mf-kicker">{{语境，如 主理人说}}</p>
  <p class="mf-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="mf-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="mf-pill mf-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="mf-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="mf-fern" aria-hidden="true"></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mf-kicker, mf-quote, mt-l, mf-src, mt-m, row, mf-pill, mf-pill-accent, mf-fern, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。薄荷题签 + 大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="mf-kicker">{{进度，如 第二部分 · 配方观}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="mf-pill mf-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="mf-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, mf-kicker, h1, mt-m, lede, mt-l, row, mf-pill, mf-pill-accent, deck-footer, slide-number, notes

---

## moments（节奏时间线）
指纹：chart
数量：mf-tl-item=4

用途：3-4 个节点的横向时间线：薄荷圆角块圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="mf-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="mf-tl mt-l" style="margin-top:52px">
    <div class="mf-tl-item"><div class="mf-tl-dot">壹</div><div class="mf-tl-t">{{时间点，≤10 字符}}</div><p class="mf-tl-d">{{事件，12-26 字}}</p></div>
    <div class="mf-tl-item"><div class="mf-tl-dot">贰</div><div class="mf-tl-t">{{时间点，≤10 字符}}</div><p class="mf-tl-d">{{事件，12-26 字}}</p></div>
    <div class="mf-tl-item"><div class="mf-tl-dot">叁</div><div class="mf-tl-t">{{时间点，≤10 字符}}</div><p class="mf-tl-d">{{事件，12-26 字}}</p></div>
    <div class="mf-tl-item"><div class="mf-tl-dot">肆</div><div class="mf-tl-t">{{时间点，≤10 字符}}</div><p class="mf-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#8CA093;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mf-kicker, h2, mt-m, mf-tl, mt-l, mf-tl-item, mf-tl-dot, mf-tl-t, mf-tl-d, deck-footer, slide-number, notes

---

## closing（收尾圆签）
指纹：hero

用途：收尾页。大字标题 + 一句行动提醒 + 墨绿胶囊按钮 + 描边药丸 + 苔绿方印。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="mf-kicker">{{提醒语境，如 首发预约}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="mf-btn">{{按钮文案，≤8 字}}</span>
    <span class="mf-pill">{{次级信息，≤10 字}}</span>
    <div class="mf-badge">{{方印 2 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, mf-kicker, h1, mt-m, lede, mt-l, row, mf-btn, mf-pill, mf-badge, deck-footer, slide-number, notes
