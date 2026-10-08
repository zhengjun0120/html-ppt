# 融资路演 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `pv-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-pitch-deck-vc` 作用域生效，骨架里已写全，照抄结构即可。
> 每页的白→淡紫渐变与右上圆环由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（数字即主角）**：紫罗兰渐变只允许出现在四处——封面侧板（pv-panel）、编号块
> （pv-n）、实心标签与行动钮（pv-tag-solid / pv-btn）。标题与数字永远深靛（--ink），
> 等宽标签（pv-kicker / pv-tag / pv-src）是唯一的「仪表盘口音」。
> 一页只传达一个关键信息：大留白是核心，每页内容块不超过两组 + 页脚。

---

## cover（路演封面）
指纹：hero

用途：开场页。左靛紫渐变面板装品牌与轮次，右侧 traction 三个大数字先声夺人。
适用 role：cover。
内容约束：主标题 ≤12 字；面板副文 30-60 字（含本轮融资额）；3 个数值各 ≤6 字符；口径行 20-44 字。

```html
<section class="slide full" data-layout="cover">
  <div class="grid g2" style="grid-template-columns:46fr 54fr;gap:56px;align-items:stretch;flex:1">
    <div class="pv-panel">
      <div class="pv-brand"><span class="pv-mark"></span><span>{{品牌名，≤8 字}}</span></div>
      <p class="pv-round">SEED ROUND · {{年份}}</p>
      <h1 class="pv-panel-title">{{项目一句话，≤12 字，可 <br> 分行}}</h1>
      <p class="pv-panel-sub">{{一句话定位 + 本轮融资额与出让比例，30-60 字}}</p>
    </div>
    <div class="stack" style="display:flex;flex-direction:column;justify-content:center">
      <p class="pv-kicker">TRACTION · {{阶段标签，≤8 字}}</p>
      <div class="grid g3 mt-l" style="gap:32px">
        <div class="pv-metric"><div class="pv-metric-v">{{数值 ≤6 字符}}</div><div class="pv-metric-l">{{指标名 ≤8 字}}</div></div>
        <div class="pv-metric"><div class="pv-metric-v">{{数值 ≤6 字符}}</div><div class="pv-metric-l">{{指标名 ≤8 字}}</div></div>
        <div class="pv-metric"><div class="pv-metric-v">{{数值 ≤6 字符}}</div><div class="pv-metric-l">{{指标名 ≤8 字}}</div></div>
      </div>
      <p class="pv-src mt-l">口径：{{数据口径与来源，20-44 字}}</p>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 轮次}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, grid, g2, pv-panel, pv-brand, pv-mark, pv-round, pv-panel-title, pv-panel-sub, stack, pv-kicker, g3, mt-l, pv-metric, pv-metric-v, pv-metric-l, pv-src, deck-footer, slide-number, notes

---

## contents（融资议程）
指纹：table
数量：pv-row=4

用途：议程页。一块白卡里放 4 行投资人问题：编号块 + 问题 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；问题 ≤10 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="pv-kicker">AGENDA · {{引导语，≤8 字}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="pv-card mt-l" style="margin-top:44px">
    <div class="pv-row"><span class="pv-n">01</span><span class="pv-t">{{投资人问题，≤10 字}}</span><span class="pv-d">{{说明，14-26 字}}</span></div>
    <div class="pv-row"><span class="pv-n">02</span><span class="pv-t">{{投资人问题，≤10 字}}</span><span class="pv-d">{{说明，14-26 字}}</span></div>
    <div class="pv-row"><span class="pv-n">03</span><span class="pv-t">{{投资人问题，≤10 字}}</span><span class="pv-d">{{说明，14-26 字}}</span></div>
    <div class="pv-row"><span class="pv-n">04</span><span class="pv-t">{{投资人问题，≤10 字}}</span><span class="pv-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 轮次}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pv-kicker, h2, mt-m, pv-card, mt-l, pv-row, pv-n, pv-t, pv-d, deck-footer, slide-number, notes

---

## metrics（增长数字）
指纹：chart
数量：pv-metric=3

用途：三个核心指标 + 月度柱阵。紫罗兰左线大数字是唯一主视觉，口径与来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源并说明柱高含义。

```html
<section class="slide" data-layout="metrics">
  <p class="pv-kicker">{{数据语境，如 TRACTION · 复盘}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:36px;margin-top:44px">
    <div class="pv-metric"><div class="pv-metric-v">{{数值 ≤6 字符}}</div><div class="pv-metric-l">{{指标名 ≤8 字}}</div></div>
    <div class="pv-metric"><div class="pv-metric-v">{{数值 ≤6 字符}}</div><div class="pv-metric-l">{{指标名 ≤8 字}}</div></div>
    <div class="pv-metric"><div class="pv-metric-v">{{数值 ≤6 字符}}</div><div class="pv-metric-l">{{指标名 ≤8 字}}</div></div>
  </div>
  <div class="pv-bars mt-m" style="margin-top:36px">
    <span class="pv-bar pv-bar-dim" style="height:16%"></span>
    <span class="pv-bar pv-bar-dim" style="height:24%"></span>
    <span class="pv-bar pv-bar-dim" style="height:31%"></span>
    <span class="pv-bar pv-bar-dim" style="height:42%"></span>
    <span class="pv-bar pv-bar-dim" style="height:55%"></span>
    <span class="pv-bar" style="height:68%"></span>
    <span class="pv-bar" style="height:84%"></span>
    <span class="pv-bar" style="height:100%"></span>
  </div>
  <p class="pv-src mt-m">来源：{{出处与统计区间，14-40 字，并说明柱高含义}}</p>
  <div class="deck-footer"><span>{{署名 · 轮次}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pv-kicker, h2, mt-m, grid, g3, mt-l, pv-metric, pv-metric-v, pv-metric-l, pv-bars, pv-bar, pv-bar-dim, pv-src, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度标签 + 大字章节名 + 一个过渡问题 + 两个看点标签。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；标签各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="pv-kicker">{{进度，如 02 · 增长与壁垒}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:14px">
    <span class="pv-tag-solid">{{看点 1，≤8 字}}</span>
    <span class="pv-tag">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 轮次}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, pv-kicker, h1, mt-m, lede, row, mt-l, pv-tag-solid, pv-tag, deck-footer, slide-number, notes

---

## keynotes（三卡论点）
指纹：cards
数量：pv-card=3

用途：恰好三张白卡。每张：实心渐变标签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；标签 ≤6 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="pv-kicker">{{章节标签}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="pv-card"><span class="pv-tag-solid">{{标签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#4338CA">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="pv-card"><span class="pv-tag-solid">{{标签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#4338CA">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="pv-card"><span class="pv-tag-solid">{{标签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#4338CA">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 轮次}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pv-kicker, h2, mt-m, grid, g3, mt-l, pv-card, pv-tag-solid, h4, mt-s, deck-footer, slide-number, notes

---

## split（打法拆解）
指纹：split
数量：pv-step=3

用途：左边把打法讲透（lede + 补充 + 标签），右边白卡装三行执行步骤。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个标签；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="pv-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:64px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#4338CA">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:14px">
        <span class="pv-tag">{{要点 1，≤8 字}}</span>
        <span class="pv-tag">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="pv-card">
      <div class="pv-step"><span class="pv-n">01</span><p class="pv-mini-t">{{一步，12-26 字}}</p></div>
      <div class="pv-step"><span class="pv-n">02</span><p class="pv-mini-t">{{一步，12-26 字}}</p></div>
      <div class="pv-step"><span class="pv-n">03</span><p class="pv-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 轮次}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pv-kicker, h2, mt-m, grid, g2, mt-l, lede, row, pv-tag, pv-card, pv-step, pv-n, pv-mini-t, deck-footer, slide-number, notes

---

## moments（里程碑）
指纹：chart
数量：pv-tl-item=4

用途：3-4 个节点的横向时间线：编号圆点 + 等宽时间点 + 一句事件，讲增长故事。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="pv-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="pv-tl mt-l" style="margin-top:48px">
    <div class="pv-tl-item"><div class="pv-tl-dot">01</div><div class="pv-tl-t">{{时间点，≤10 字符}}</div><p class="pv-tl-d">{{事件，12-26 字}}</p></div>
    <div class="pv-tl-item"><div class="pv-tl-dot">02</div><div class="pv-tl-t">{{时间点，≤10 字符}}</div><p class="pv-tl-d">{{事件，12-26 字}}</p></div>
    <div class="pv-tl-item"><div class="pv-tl-dot">03</div><div class="pv-tl-t">{{时间点，≤10 字符}}</div><p class="pv-tl-d">{{事件，12-26 字}}</p></div>
    <div class="pv-tl-item"><div class="pv-tl-dot">04</div><div class="pv-tl-t">{{时间点，≤10 字符}}</div><p class="pv-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="pv-src mt-m">{{口径或读法，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 轮次}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pv-kicker, h2, mt-m, pv-tl, mt-l, pv-tl-item, pv-tl-dot, pv-tl-t, pv-tl-d, pv-src, mt-m, deck-footer, slide-number, notes

---

## quote（资本之言）
指纹：quote

用途：整页一句创始人或资本视角的宣言：大字引文 + 出处 + 两个支撑标签。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；标签各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="pv-kicker">{{语境，如 创始人说}}</p>
  <p class="pv-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="pv-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:14px">
    <span class="pv-tag-solid">{{支撑点 1，≤8 字}}</span>
    <span class="pv-tag">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 轮次}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pv-kicker, pv-quote, mt-l, pv-src, mt-m, row, pv-tag-solid, pv-tag, deck-footer, slide-number, notes

---

## closing（The Ask）
指纹：hero

用途：收尾要钱。大字融资额 + 资金用途 + 渐变行动钮 + 联系标签，一页只说这一件事。
适用 role：thanks / cta / content。
内容约束：标题 ≤12 字；lede 20-48 字；按钮 ≤6 字；标签各 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="pv-kicker">THE ASK · {{语境}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{资金用途与承诺，20-48 字}}</p>
  <div class="row mt-l" style="gap:22px">
    <span class="pv-btn">{{按钮文案，≤6 字}}</span>
    <span class="pv-tag">{{次级信息，≤10 字}}</span>
    <span class="pv-tag-solid">{{次级信息，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 轮次}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, pv-kicker, h1, mt-m, lede, row, mt-l, pv-btn, pv-tag, pv-tag-solid, deck-footer, slide-number, notes
