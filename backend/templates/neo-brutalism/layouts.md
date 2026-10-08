# 新野兽派 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `nb-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-neo-brutalism` 作用域生效，骨架里已写全，照抄结构即可。
> 暖画布与角落斜置色块由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（拳头要省着挥）**：描边永远 3px 起、阴影永远零模糊（Npx Npx 0）；情绪色只有
> 荧光橙、信号红、荧光黄三种，按「黑 / 橙 / 红」轮换阴影与节点。卡片允许少量歪斜与偏移，
> 但一页至多三张卡 + 页脚，卡与卡之间留足空隙——不做积木墙。徽章、行动钮、通栏条是
> 唯一的三处实底色块，别再发明第四处。

---

## cover（拳头封面）
指纹：hero

用途：开场页。粗标签 + 黑底徽章 + 描边标题盒大字，一排歪斜标签，页底黑通栏条压阵。
适用 role：cover。
内容约束：主标题 ≤8 字；lede 30-60 字；标签各 ≤6 字；通栏条两段各 ≤16 字符。

```html
<section class="slide full" data-layout="cover">
  <div class="row" style="justify-content:space-between">
    <p class="nb-kicker">{{音乐节名 · 品牌提案}}</p>
    <span class="nb-badge">PROPOSAL · {{年份}}</span>
  </div>
  <h1 class="h1 mt-l"><span class="nb-box">{{主标题，≤8 字，可 <br> 分行}}</span></h1>
  <p class="lede mt-l" style="max-width:50ch">{{提案背景：几天几夜、多少乐队、要解决什么，30-60 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="nb-tag">{{交付物 1，≤6 字}}</span>
    <span class="nb-tag">{{交付物 2，≤6 字}}</span>
    <span class="nb-tag">{{交付物 3，≤6 字}}</span>
  </div>
  <div class="nb-strip mt-l"><span>{{音乐节英文名，≤16 字符}}</span><span>BRAND PROPOSAL · {{版本}}</span></div>
  <div class="deck-footer"><span>{{署名 · 提案}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, row, nb-kicker, nb-badge, h1, mt-l, nb-box, lede, nb-tag, nb-strip, deck-footer, slide-number, notes

---

## contents（提案四步）
指纹：table
数量：nb-row=4

用途：议程页。白卡里放 4 行条目：黄色编号块 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="nb-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="nb-card mt-l" style="margin-top:44px">
    <div class="nb-row"><span class="nb-n">01</span><span class="nb-t">{{篇名，≤8 字}}</span><span class="nb-d">{{说明，14-26 字}}</span></div>
    <div class="nb-row"><span class="nb-n">02</span><span class="nb-t">{{篇名，≤8 字}}</span><span class="nb-d">{{说明，14-26 字}}</span></div>
    <div class="nb-row"><span class="nb-n">03</span><span class="nb-t">{{篇名，≤8 字}}</span><span class="nb-d">{{说明，14-26 字}}</span></div>
    <div class="nb-row"><span class="nb-n">04</span><span class="nb-t">{{篇名，≤8 字}}</span><span class="nb-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 提案}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, nb-kicker, h2, mt-m, nb-card, mt-l, nb-row, nb-n, nb-t, nb-d, deck-footer, slide-number, notes

---

## metrics（复盘重拳）
指纹：chart
数量：nb-stat=3

用途：三个关键数据。描边数据卡，中间一张橙底重色卡做焦点；口径写页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="nb-kicker">{{数据语境}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:36px;margin-top:48px">
    <div class="nb-stat"><div class="nb-stat-v">{{数值 ≤6 字符}}</div><div class="nb-stat-l">{{指标名，≤8 字}}</div><p class="nb-stat-note">{{口径，14-30 字}}</p></div>
    <div class="nb-stat nb-stat-hot"><div class="nb-stat-v">{{数值 ≤6 字符}}</div><div class="nb-stat-l">{{指标名，≤8 字}}</div><p class="nb-stat-note">{{口径，14-30 字}}</p></div>
    <div class="nb-stat nb-wild"><div class="nb-stat-v">{{数值 ≤6 字符}}</div><div class="nb-stat-l">{{指标名，≤8 字}}</div><p class="nb-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="nb-src mt-m">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 提案}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, nb-kicker, h2, mt-m, grid, g3, mt-l, nb-stat, nb-stat-hot, nb-wild, nb-stat-v, nb-stat-l, nb-stat-note, nb-src, mt-m, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度粗标签 + 描边标题盒大字章节名 + 过渡问题 + 两个歪斜标签。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；标签各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="nb-kicker">{{进度，如 PART 02 · 核心想法}}</p>
  <h1 class="h1 mt-l"><span class="nb-box">{{章节标题，≤10 字，可 <br> 分行}}</span></h1>
  <p class="lede mt-l" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="nb-tag">{{看点 1，≤8 字}}</span>
    <span class="nb-tag">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 提案}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, nb-kicker, h1, mt-l, nb-box, lede, row, mt-l, nb-tag, deck-footer, slide-number, notes

---

## keynotes（三张重拳）
指纹：cards
数量：nb-card=3

用途：恰好三张描边白卡，阴影按黑/橙/红轮换。每张：标签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；标签 ≤6 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="nb-kicker">{{章节标签}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:36px;margin-top:44px">
    <div class="nb-card"><span class="nb-tag">{{标签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:18px;line-height:1.8;color:#1E293B">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="nb-card nb-hot"><span class="nb-tag">{{标签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:18px;line-height:1.8;color:#1E293B">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="nb-card nb-wild"><span class="nb-tag">{{标签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:18px;line-height:1.8;color:#1E293B">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 提案}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, nb-kicker, h2, mt-m, grid, g3, mt-l, nb-card, nb-tag, nb-hot, nb-wild, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（落地拆解）
指纹：split
数量：nb-step=3

用途：左边把方案讲透（lede + 补充 + 标签），右边白卡装三行落地步骤。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个标签；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="nb-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:64px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么落地、凭什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:18px;line-height:1.8;color:#1E293B">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="nb-tag">{{要点 1，≤8 字}}</span>
        <span class="nb-tag">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="nb-card">
      <div class="nb-step"><span class="nb-step-i">01</span><p class="nb-mini-t">{{一步，12-26 字}}</p></div>
      <div class="nb-step"><span class="nb-step-i">02</span><p class="nb-mini-t">{{一步，12-26 字}}</p></div>
      <div class="nb-step"><span class="nb-step-i">03</span><p class="nb-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 提案}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, nb-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, nb-tag, nb-card, nb-step, nb-step-i, nb-mini-t, deck-footer, slide-number, notes

---

## moments（筹备节奏）
指纹：chart
数量：nb-tl-item=4

用途：3-4 个节点的横向时间线：描边方块节点（奶油/橙/黄/红轮换）+ 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="nb-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="nb-tl mt-l" style="margin-top:56px">
    <div class="nb-tl-item"><div class="nb-tl-dot">01</div><div class="nb-tl-t">{{时间点，≤10 字符}}</div><p class="nb-tl-d">{{事件，12-26 字}}</p></div>
    <div class="nb-tl-item"><div class="nb-tl-dot">02</div><div class="nb-tl-t">{{时间点，≤10 字符}}</div><p class="nb-tl-d">{{事件，12-26 字}}</p></div>
    <div class="nb-tl-item"><div class="nb-tl-dot">03</div><div class="nb-tl-t">{{时间点，≤10 字符}}</div><p class="nb-tl-d">{{事件，12-26 字}}</p></div>
    <div class="nb-tl-item"><div class="nb-tl-dot">04</div><div class="nb-tl-t">{{时间点，≤10 字符}}</div><p class="nb-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="nb-src mt-m">{{口径或读法，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 提案}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, nb-kicker, h2, mt-m, nb-tl, mt-l, nb-tl-item, nb-tl-dot, nb-tl-t, nb-tl-d, nb-src, mt-m, deck-footer, slide-number, notes

---

## quote（立场宣言）
指纹：quote

用途：整页一句立场宣言。900 大字两行（下行荧光黄荧光笔划）+ 等宽出处 + 两个标签。
适用 role：quote。
内容约束：引文共 10-28 字；出处 ≤22 字且真实；标签各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="nb-kicker">{{语境，如 发起人说}}</p>
  <p class="nb-quote mt-l" style="margin-top:40px">「{{引文上行，≤14 字}}<br><span class="nb-mark">{{引文下行，≤14 字}}</span>」</p>
  <p class="nb-src mt-m">—— {{出处：人与场合，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="nb-tag">{{支撑点 1，≤8 字}}</span>
    <span class="nb-tag">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 提案}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, nb-kicker, nb-quote, mt-l, nb-mark, nb-src, mt-m, row, nb-tag, deck-footer, slide-number, notes

---

## closing（拍板页）
指纹：hero

用途：收尾页。描边标题盒 + 一句承诺 + 橙底行动钮 + 标签 + 黑通栏条收尾。
适用 role：thanks / cta / content。
内容约束：标题 ≤8 字；lede 20-48 字；按钮 ≤6 字；标签各 ≤10 字；通栏条两段各 ≤16 字符。

```html
<section class="slide full" data-layout="closing">
  <p class="nb-kicker">{{语境，如 拍板与排期}}</p>
  <h1 class="h1 mt-m"><span class="nb-box">{{收束标题，≤8 字，可 <br> 分行}}</span></h1>
  <p class="lede mt-m" style="max-width:48ch">{{承诺与交付节奏，20-48 字}}</p>
  <div class="row mt-l" style="gap:20px">
    <span class="nb-btn">{{按钮文案，≤6 字}}</span>
    <span class="nb-tag">{{次级信息，≤10 字}}</span>
    <span class="nb-tag">{{次级信息，≤10 字}}</span>
  </div>
  <div class="nb-strip mt-l"><span>{{音乐节英文名，≤16 字符}}</span><span>{{结束语，≤16 字符}}</span></div>
  <div class="deck-footer"><span>{{署名 · 提案}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, nb-kicker, h1, mt-m, nb-box, lede, row, mt-l, nb-btn, nb-tag, nb-strip, deck-footer, slide-number, notes
