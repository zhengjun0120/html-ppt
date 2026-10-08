# 蒸汽波 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `vw-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-vaporwave` 作用域生效，骨架里已写全，照抄结构即可。
> 每页条纹落日、柱影与地平线网格由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（渐变是身份，不是正文）**：粉 #FF71CE 与青 #01CDFE 只出现在三处——
> ambient（日轮/网格/柱影，模板自动衬底）、数据与时间线（vw-stat-v / vw-tl-t / vw-tl-dot /
> vw-n）、行动钮（vw-btn）。正文永远灰白 #E2E8F0 或冷白 #F8FAFC，落在深紫底上；
> 渐变文字禁用于正文段落。卡只有一种——半透明暗紫大圆角（vw-card），禁硬黑框与锐利直角。

---

## cover（落日封面）
指纹：hero

用途：开场页。等宽题签 + 衬线宽距大标题 + 一句定位，粉青描边标签点题。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 28-56 字；题签 ≤14 字；标签各 ≤10 字。

```html
<section class="slide full" data-layout="cover">
  <p class="vw-kicker">{{题签，≤14 字，如 回声电台 · 曲风分享}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句话定位，28-56 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="vw-pill vw-pill-accent">{{标签 1，≤10 字}}</span>
    <span class="vw-pill">{{标签 2，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, vw-kicker, h1, mt-m, lede, mt-l, row, vw-pill, vw-pill-accent, deck-footer, slide-number, notes

---

## contents（目录卷）
指纹：table
数量：vw-item=4

用途：议程页。一块紫卡里放 4 行篇目：青圈编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="vw-kicker">{{引导语，如 今晚卷目}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="vw-card mt-l" style="margin-top:44px">
    <div class="vw-item"><span class="vw-n">01</span><span class="vw-t">{{篇名，≤8 字}}</span><span class="vw-d">{{说明，14-26 字}}</span></div>
    <div class="vw-item"><span class="vw-n">02</span><span class="vw-t">{{篇名，≤8 字}}</span><span class="vw-d">{{说明，14-26 字}}</span></div>
    <div class="vw-item"><span class="vw-n">03</span><span class="vw-t">{{篇名，≤8 字}}</span><span class="vw-d">{{说明，14-26 字}}</span></div>
    <div class="vw-item"><span class="vw-n">04</span><span class="vw-t">{{篇名，≤8 字}}</span><span class="vw-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, vw-kicker, h2, mt-m, vw-card, mt-l, vw-item, vw-n, vw-t, vw-d, deck-footer, slide-number, notes

---

## keynotes（三源流卡）
指纹：cards
数量：vw-card=3

用途：恰好三张紫卡。每张：粉药丸题签 + 衬线小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="vw-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:36px;margin-top:44px">
    <div class="vw-card"><span class="vw-pill vw-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#E2E8F0">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="vw-card"><span class="vw-pill vw-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#E2E8F0">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="vw-card"><span class="vw-pill vw-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#E2E8F0">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, vw-kicker, h2, mt-m, grid, g3, mt-l, vw-card, vw-pill, vw-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：vw-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边紫卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="vw-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#E2E8F0">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="vw-pill">{{要点 1，≤10 字}}</span>
        <span class="vw-pill vw-pill-accent">{{要点 2，≤10 字}}</span>
      </div>
    </div>
    <div class="vw-card">
      <div class="vw-step"><span class="vw-n">一</span><p class="vw-mini-t">{{一步，12-26 字}}</p></div>
      <div class="vw-step"><span class="vw-n">二</span><p class="vw-mini-t">{{一步，12-26 字}}</p></div>
      <div class="vw-step"><span class="vw-n">三</span><p class="vw-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, vw-kicker, h2, mt-m, grid, g2, mt-l, lede, row, vw-pill, vw-pill-accent, vw-card, vw-step, vw-n, vw-mini-t, deck-footer, slide-number, notes

---

## metrics（渐变数字）
指纹：chart
数量：vw-stat=3

用途：三个关键数据。粉青渐变大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="vw-kicker">{{数据语境，如 年度报告}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:48px;margin-top:48px">
    <div class="vw-stat"><div class="vw-stat-v">{{数值 ≤6 字符}}<span class="vw-stat-u">{{单位}}</span></div><div class="vw-stat-l">{{指标名，≤8 字}}</div><p class="vw-stat-note">{{口径，14-30 字}}</p></div>
    <div class="vw-stat"><div class="vw-stat-v">{{数值 ≤6 字符}}<span class="vw-stat-u">{{单位}}</span></div><div class="vw-stat-l">{{指标名，≤8 字}}</div><p class="vw-stat-note">{{口径，14-30 字}}</p></div>
    <div class="vw-stat"><div class="vw-stat-v">{{数值 ≤6 字符}}<span class="vw-stat-u">{{单位}}</span></div><div class="vw-stat-l">{{指标名，≤8 字}}</div><p class="vw-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:16px;color:#A99BC8;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, vw-kicker, h2, mt-m, grid, g3, mt-l, vw-stat, vw-stat-v, vw-stat-u, vw-stat-l, vw-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（题记引文）
指纹：quote

用途：整页一句引文。衬线宽距大字 + 出处 + 两个支撑药丸，右侧留白处可立竖排短语。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="quote">
  <p class="vw-kicker">{{语境，如 电台题记}}</p>
  <p class="vw-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="vw-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="vw-pill vw-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="vw-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="vw-vert" style="position:absolute;right:120px;top:50%;transform:translateY(-50%)">{{竖排短语 ≤7 字}}<span class="vw-vert-accent">{{一字或两字}}</span></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, vw-kicker, vw-quote, mt-l, vw-src, mt-m, row, vw-pill, vw-pill-accent, vw-vert, vw-vert-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 衬线大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="vw-kicker">{{进度，如 第二段 · 视觉语法}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:46ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="vw-pill vw-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="vw-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, vw-kicker, h1, mt-m, lede, mt-l, row, vw-pill, vw-pill-accent, deck-footer, slide-number, notes

---

## moments（巡演时间线）
指纹：chart
数量：vw-tl-item=4

用途：3-4 个节点的横向时间线：渐变圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="vw-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="vw-tl mt-l" style="margin-top:52px">
    <div class="vw-tl-item"><div class="vw-tl-dot"></div><div class="vw-tl-t">{{时间点，≤10 字符}}</div><p class="vw-tl-d">{{事件，12-26 字}}</p></div>
    <div class="vw-tl-item"><div class="vw-tl-dot"></div><div class="vw-tl-t">{{时间点，≤10 字符}}</div><p class="vw-tl-d">{{事件，12-26 字}}</p></div>
    <div class="vw-tl-item"><div class="vw-tl-dot"></div><div class="vw-tl-t">{{时间点，≤10 字符}}</div><p class="vw-tl-d">{{事件，12-26 字}}</p></div>
    <div class="vw-tl-item"><div class="vw-tl-dot"></div><div class="vw-tl-t">{{时间点，≤10 字符}}</div><p class="vw-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:16px;color:#A99BC8;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, vw-kicker, h2, mt-m, vw-tl, mt-l, vw-tl-item, vw-tl-dot, vw-tl-t, vw-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（收尾场刊）
指纹：hero

用途：收尾页。衬线大字 + 一句行动提醒 + 渐变按钮 + 描边药丸。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤12 字。

```html
<section class="slide full" data-layout="closing">
  <p class="vw-kicker">{{提醒语境，如 午夜商场 · 上海首站}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:46ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:20px">
    <span class="vw-btn">{{按钮文案，≤8 字}}</span>
    <span class="vw-pill">{{次级信息，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, vw-kicker, h1, mt-m, lede, mt-l, row, vw-btn, vw-pill, deck-footer, slide-number, notes
