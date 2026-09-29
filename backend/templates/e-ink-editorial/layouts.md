# 墨页叙事 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `ei-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-e-ink-editorial` 作用域生效，骨架里已写全，照抄结构即可。
> 每页纸面颗粒与右缘页边细线由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（无彩色，靠秩序）**：全模板没有任何彩色——纯黑 #0A0A0B 就是唯一的强调色；
> 阅读秩序靠三级灰（纯黑/炭灰/元数据灰）、1px 细线（ei-rule / ei-card / ei-item 自带）与
> 等宽元数据（ei-folio / ei-cap / ei-pill / deck-footer）建立。整页最多一个反转面
> （ei-card-ink 或 ei-pill-accent 二选一）；禁彩色、禁阴影、禁圆角卡、禁粗黑大色带。

---

## cover（卷首页）
指纹：hero

用途：开场页。等宽刊号行 + 衬线大标题 + 通栏付印线 + 一句导语，像长文的开卷。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；刊号行 ≤16 字符；药丸 ≤10 字。

```html
<section class="slide full" data-layout="cover">
  <p class="ei-folio">{{刊号行，≤16 字符，如 纸上观察 · FIELD NOTES 07}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="ei-rule mt-s"></div>
  <p class="lede mt-m" style="max-width:52ch">{{一句导语，24-48 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="ei-pill ei-pill-accent">{{栏目签，≤6 字}}</span>
    <span class="ei-pill">{{季节或期号，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, ei-folio, h1, mt-m, ei-rule, mt-s, lede, mt-m, row, mt-l, ei-pill, ei-pill-accent, deck-footer, slide-number, notes

---

## contents（目次页）
指纹：table
数量：ei-item=4

用途：议程页。一块细线框卡里放 4 行目次：等宽编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="ei-folio">{{引导语，如 本期结构}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="ei-card mt-l" style="margin-top:44px">
    <div class="ei-item"><span class="ei-n">01</span><span class="ei-t">{{篇名，≤8 字}}</span><span class="ei-d">{{说明，14-26 字}}</span></div>
    <div class="ei-item"><span class="ei-n">02</span><span class="ei-t">{{篇名，≤8 字}}</span><span class="ei-d">{{说明，14-26 字}}</span></div>
    <div class="ei-item"><span class="ei-n">03</span><span class="ei-t">{{篇名，≤8 字}}</span><span class="ei-d">{{说明，14-26 字}}</span></div>
    <div class="ei-item"><span class="ei-n">04</span><span class="ei-t">{{篇名，≤8 字}}</span><span class="ei-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ei-folio, h2, mt-m, ei-card, mt-l, ei-item, ei-n, ei-t, ei-d, deck-footer, slide-number, notes

---

## keynotes（三段论）
指纹：cards
数量：ei-card=3

用途：恰好三张细线框卡。每张：黑底反白题签 + 衬线小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="ei-folio">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="ei-card"><span class="ei-pill ei-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.75;color:#45423C">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="ei-card"><span class="ei-pill ei-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.75;color:#45423C">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="ei-card"><span class="ei-pill ei-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.75;color:#45423C">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ei-folio, h2, mt-m, grid, g3, mt-l, ei-card, ei-pill, ei-pill-accent, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：ei-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边墨黑实块装三行步骤，纸墨反转。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="ei-folio">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.75;color:#45423C">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="ei-pill">{{要点 1，≤8 字}}</span>
        <span class="ei-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="ei-card-ink">
      <div class="ei-step"><span class="ei-n">一</span><p class="ei-mini-t">{{一步，12-26 字}}</p></div>
      <div class="ei-step"><span class="ei-n">二</span><p class="ei-mini-t">{{一步，12-26 字}}</p></div>
      <div class="ei-step"><span class="ei-n">三</span><p class="ei-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ei-folio, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, mt-l, ei-pill, ei-card-ink, ei-step, ei-n, ei-mini-t, deck-footer, slide-number, notes

---

## metrics（墨黑数字）
指纹：chart
数量：ei-stat=3

用途：三个关键数据。纯黑衬线大数字是唯一主视觉，口径写进卡内，来源用等宽字写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="ei-folio">{{数据语境，如 样本口径}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="ei-stat"><div class="ei-stat-v">{{数值 ≤6 字符}}<span class="ei-stat-u">{{单位}}</span></div><div class="ei-stat-l">{{指标名，≤8 字}}</div><p class="ei-stat-note">{{口径，14-30 字}}</p></div>
    <div class="ei-stat"><div class="ei-stat-v">{{数值 ≤6 字符}}<span class="ei-stat-u">{{单位}}</span></div><div class="ei-stat-l">{{指标名，≤8 字}}</div><p class="ei-stat-note">{{口径，14-30 字}}</p></div>
    <div class="ei-stat"><div class="ei-stat-v">{{数值 ≤6 字符}}<span class="ei-stat-u">{{单位}}</span></div><div class="ei-stat-l">{{指标名，≤8 字}}</div><p class="ei-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="ei-cap mt-m">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ei-folio, h2, mt-m, grid, g3, mt-l, ei-stat, ei-stat-v, ei-stat-u, ei-stat-l, ei-stat-note, ei-cap, mt-m, deck-footer, slide-number, notes

---

## quote（访谈引文）
指纹：quote

用途：整页一句引文。衬线大字 + 等宽出处 + 两个支撑药丸，右侧留白处一行页码注记。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="ei-folio">{{语境，如 店主访谈}}</p>
  <p class="ei-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="ei-cap mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="ei-pill ei-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="ei-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <p class="ei-cap" style="position:absolute;right:120px;top:50%;transform:translateY(-50%)">{{页码注记，≤14 字符}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ei-folio, ei-quote, mt-l, ei-cap, mt-m, row, mt-l, ei-pill, ei-pill-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。等宽进度行 + 衬线大字章节名 + 一个过渡句 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡句 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="ei-folio">{{进度，如 第二章 · 活动即生意}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="ei-pill ei-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="ei-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, ei-folio, h1, mt-m, lede, mt-l, row, ei-pill, ei-pill-accent, deck-footer, slide-number, notes

---

## moments（流程时间线）
指纹：chart
数量：ei-tl-item=4

用途：3-4 个节点的横向时间线：细线方点 + 时间点 + 一句事件，极细连接线。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="ei-folio">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="ei-tl mt-l" style="margin-top:52px">
    <div class="ei-tl-item"><div class="ei-tl-dot">01</div><div class="ei-tl-t">{{时间点，≤10 字符}}</div><p class="ei-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ei-tl-item"><div class="ei-tl-dot">02</div><div class="ei-tl-t">{{时间点，≤10 字符}}</div><p class="ei-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ei-tl-item"><div class="ei-tl-dot">03</div><div class="ei-tl-t">{{时间点，≤10 字符}}</div><p class="ei-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ei-tl-item"><div class="ei-tl-dot">04</div><div class="ei-tl-t">{{时间点，≤10 字符}}</div><p class="ei-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="ei-cap mt-m">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ei-folio, h2, mt-m, ei-tl, mt-l, ei-tl-item, ei-tl-dot, ei-tl-t, ei-tl-d, ei-cap, mt-m, deck-footer, slide-number, notes

---

## closing（卷尾互动）
指纹：hero

用途：收尾页。衬线大字 + 一句行动提醒 + 细线描边按钮 + 描边药丸。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="ei-folio">{{提醒语境，如 互动}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="ei-btn">{{按钮文案，≤8 字}}</span>
    <span class="ei-pill">{{次级信息，≤10 字}}</span>
    <span class="ei-pill ei-pill-accent">{{时间或期号，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, ei-folio, h1, mt-m, lede, mt-l, row, ei-btn, ei-pill, ei-pill-accent, deck-footer, slide-number, notes
