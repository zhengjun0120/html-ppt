# 杂志衬线 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `es-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-editorial-serif` 作用域生效，骨架里已写全，照抄结构即可。
> 每页左缘页边栏线与右下刊尾花饰由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（深橙是红笔）**：全页只用暖色系——深褐标题、焦棕正文、深橙批注；深橙只允许出现在
> 六处——细杠与短线（es-tag / es-rule 自带）、首字下沉（es-drop）、橙圈序号（es-n / es-tl-dot）、
> 关键数据（es-stat-v）、引言左线（es-quote-block 自带）、药丸圈点（es-pill-accent）。
> 禁冷色调、禁无衬线标题、禁大面积橙色、禁装饰打断叙事。首字下沉 es-drop 全份最多两处，
> 只用于 lede 导语的第一个字；斜体英文 es-en 一页至多一行。

---

## cover（刊首导读）
指纹：hero

用途：开场页。刊名标签 + 衬线标题 + 英文副题 + 首字下沉导语，像刊物的导读页。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；标签 ≤14 字；英文副题 ≤18 字符。

```html
<section class="slide full" data-layout="cover">
  <p class="es-tag">{{刊名标签，≤14 字，如 回声文学季刊 · 总第 28 期}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <p class="es-en mt-s">{{英文副题，≤18 字符}}</p>
  <p class="lede mt-m" style="max-width:48ch"><span class="es-drop">{{首字}}</span>{{导语其余部分，24-48 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="es-pill es-pill-accent">{{栏目签，≤8 字}}</span>
    <span class="es-pill">{{期号或季节，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, es-tag, h1, mt-m, es-en, mt-s, lede, mt-m, es-drop, row, mt-l, es-pill, es-pill-accent, deck-footer, slide-number, notes

---

## contents（要目页）
指纹：table
数量：es-item=4

用途：议程页。一块暖白卡里放 4 行要目：橙圈序号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="es-tag">{{引导语，如 本期要目}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="es-card mt-l" style="margin-top:44px">
    <div class="es-item"><span class="es-n">一</span><span class="es-t">{{篇名，≤8 字}}</span><span class="es-d">{{说明，14-26 字}}</span></div>
    <div class="es-item"><span class="es-n">二</span><span class="es-t">{{篇名，≤8 字}}</span><span class="es-d">{{说明，14-26 字}}</span></div>
    <div class="es-item"><span class="es-n">三</span><span class="es-t">{{篇名，≤8 字}}</span><span class="es-d">{{说明，14-26 字}}</span></div>
    <div class="es-item"><span class="es-n">四</span><span class="es-t">{{篇名，≤8 字}}</span><span class="es-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, es-tag, h2, mt-m, es-card, mt-l, es-item, es-n, es-t, es-d, deck-footer, slide-number, notes

---

## keynotes（三读法）
指纹：cards
数量：es-card=3

用途：恰好三张暖白卡。每张：橙色题签 + 衬线小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="es-tag">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="es-card"><span class="es-pill es-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.75;color:#7C2D12">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="es-card"><span class="es-pill es-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.75;color:#7C2D12">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="es-card"><span class="es-pill es-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.75;color:#7C2D12">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, es-tag, h2, mt-m, grid, g3, mt-l, es-card, es-pill, es-pill-accent, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：es-step=3

用途：左边把一件事讲透（lede + 斜体引言块 + 药丸），右边暖白卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 引言 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="es-tag">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="es-quote-block mt-m">{{引言：编辑批注或一句话总结，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="es-pill">{{要点 1，≤8 字}}</span>
        <span class="es-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="es-card">
      <div class="es-step"><span class="es-n">一</span><p class="es-mini-t">{{一步，12-26 字}}</p></div>
      <div class="es-step"><span class="es-n">二</span><p class="es-mini-t">{{一步，12-26 字}}</p></div>
      <div class="es-step"><span class="es-n">三</span><p class="es-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, es-tag, h2, mt-m, grid, g2, mt-l, lede, es-quote-block, mt-m, row, mt-l, es-pill, es-card, es-step, es-n, es-mini-t, deck-footer, slide-number, notes

---

## metrics（橙批数字）
指纹：chart
数量：es-stat=3

用途：三个关键数据。深橙衬线大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="es-tag">{{数据语境，如 专题体量}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="es-stat"><div class="es-stat-v">{{数值 ≤6 字符}}<span class="es-stat-u">{{单位}}</span></div><div class="es-stat-l">{{指标名，≤8 字}}</div><p class="es-stat-note">{{口径，14-30 字}}</p></div>
    <div class="es-stat"><div class="es-stat-v">{{数值 ≤6 字符}}<span class="es-stat-u">{{单位}}</span></div><div class="es-stat-l">{{指标名，≤8 字}}</div><p class="es-stat-note">{{口径，14-30 字}}</p></div>
    <div class="es-stat"><div class="es-stat-v">{{数值 ≤6 字符}}<span class="es-stat-u">{{单位}}</span></div><div class="es-stat-l">{{指标名，≤8 字}}</div><p class="es-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:16px;color:#9A6B4F;letter-spacing:.08em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, es-tag, h2, mt-m, grid, g3, mt-l, es-stat, es-stat-v, es-stat-u, es-stat-l, es-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（卷首引文）
指纹：quote

用途：整页一句引文。衬线大字 + 等宽出处 + 两个支撑药丸，右侧留白处一行斜体英文。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="es-tag">{{语境，如 卷首语}}</p>
  <p class="es-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="es-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="es-pill es-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="es-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <p class="es-en" style="position:absolute;right:120px;top:50%;transform:translateY(-50%)">{{斜体英文注记，≤18 字符}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, es-tag, es-quote, mt-l, es-src, mt-m, row, mt-l, es-pill, es-pill-accent, es-en, deck-footer, slide-number, notes

---

## divider（辑间页）
指纹：hero

用途：章节过渡。辑次标签 + 衬线大字章节名 + 一个过渡句 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡句 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="es-tag">{{进度，如 辑二 · 腔调里的叙事}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="es-pill es-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="es-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, es-tag, h1, mt-m, lede, mt-l, row, es-pill, es-pill-accent, deck-footer, slide-number, notes

---

## moments（编年时间线）
指纹：chart
数量：es-tl-item=4

用途：3-4 个节点的横向时间线：橙圈圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="es-tag">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="es-tl mt-l" style="margin-top:52px">
    <div class="es-tl-item"><div class="es-tl-dot">一</div><div class="es-tl-t">{{时间点，≤10 字符}}</div><p class="es-tl-d">{{事件，12-26 字}}</p></div>
    <div class="es-tl-item"><div class="es-tl-dot">二</div><div class="es-tl-t">{{时间点，≤10 字符}}</div><p class="es-tl-d">{{事件，12-26 字}}</p></div>
    <div class="es-tl-item"><div class="es-tl-dot">三</div><div class="es-tl-t">{{时间点，≤10 字符}}</div><p class="es-tl-d">{{事件，12-26 字}}</p></div>
    <div class="es-tl-item"><div class="es-tl-dot">四</div><div class="es-tl-t">{{时间点，≤10 字符}}</div><p class="es-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:16px;color:#9A6B4F;letter-spacing:.08em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, es-tag, h2, mt-m, es-tl, mt-l, es-tl-item, es-tl-dot, es-tl-t, es-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（刊尾约稿）
指纹：hero

用途：收尾页。衬线大字 + 一句行动提醒 + 深褐描边按钮 + 描边药丸。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="es-tag">{{提醒语境，如 订阅与投稿}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="es-btn">{{按钮文案，≤8 字}}</span>
    <span class="es-pill">{{次级信息，≤10 字}}</span>
    <span class="es-pill es-pill-accent">{{时限或期号，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, es-tag, h1, mt-m, lede, mt-l, row, es-btn, es-pill, es-pill-accent, deck-footer, slide-number, notes
