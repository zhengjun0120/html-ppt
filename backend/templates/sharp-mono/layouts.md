# 锐利黑白 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `sm-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-sharp-mono` 作用域生效，骨架里已写全，照抄结构即可。
> 每页右上斜切色带与四角裁切线由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（黑白就是身份）**：全页只准无彩色（黑 #000 / 深灰 #333 / 中灰 #666 / 白与浅灰）。
> 硬卡只有一种——白底 3px 黑框 + 10px 硬阴影（sm-card）；反白块（sm-tag-invert / sm-btn）
> 每页至多两处；禁圆角、禁模糊阴影、禁细于 2px 的边框。
> 加粗大字（h1/h2）每页最多一组；正文段落永远 #333333，不要纯黑。

---

## cover（宣言封面）
指纹：hero

用途：开场页。题签 + 超大加粗标题 + 硬横线 + 一句定位，反白标签点题。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 28-56 字；题签 ≤14 字；标签各 ≤8 字。

```html
<section class="slide full" data-layout="cover">
  <p class="sm-kicker">{{题签，≤14 字，如 锐·字社 秋季公开课}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="sm-line mt-s"></div>
  <p class="lede mt-m" style="max-width:48ch">{{一句话定位，28-56 字}}</p>
  <div class="row mt-l" style="gap:24px">
    <span class="sm-tag sm-tag-invert">{{标签 1，≤8 字}}</span>
    <span class="sm-tag">{{标签 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, sm-kicker, h1, mt-m, sm-line, mt-s, lede, mt-l, row, sm-tag, sm-tag-invert, deck-footer, slide-number, notes

---

## contents（目录墙）
指纹：table
数量：sm-item=4

用途：议程页。一块硬卡里放 4 行篇目：方块编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="sm-kicker">{{引导语，如 今晚议程}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="sm-card mt-l" style="margin-top:44px">
    <div class="sm-item"><span class="sm-n">01</span><span class="sm-t">{{篇名，≤8 字}}</span><span class="sm-d">{{说明，14-26 字}}</span></div>
    <div class="sm-item"><span class="sm-n">02</span><span class="sm-t">{{篇名，≤8 字}}</span><span class="sm-d">{{说明，14-26 字}}</span></div>
    <div class="sm-item"><span class="sm-n">03</span><span class="sm-t">{{篇名，≤8 字}}</span><span class="sm-d">{{说明，14-26 字}}</span></div>
    <div class="sm-item"><span class="sm-n">04</span><span class="sm-t">{{篇名，≤8 字}}</span><span class="sm-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sm-kicker, h2, mt-m, sm-card, mt-l, sm-item, sm-n, sm-t, sm-d, deck-footer, slide-number, notes

---

## keynotes（三卡主张）
指纹：cards
数量：sm-card=3

用途：恰好三张硬卡。每张：反白题签 + 加粗小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="sm-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:36px;margin-top:44px">
    <div class="sm-card"><span class="sm-tag sm-tag-invert">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#333333">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="sm-card"><span class="sm-tag sm-tag-invert">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#333333">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="sm-card"><span class="sm-tag sm-tag-invert">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#333333">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sm-kicker, h2, mt-m, grid, g3, mt-l, sm-card, sm-tag, sm-tag-invert, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：sm-step=3

用途：左边把一件事讲透（lede + 补充 + 标签），右边硬卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个标签；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="sm-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#333333">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="sm-tag">{{要点 1，≤10 字}}</span>
        <span class="sm-tag">{{要点 2，≤10 字}}</span>
      </div>
    </div>
    <div class="sm-card">
      <div class="sm-step"><span class="sm-n">1</span><p class="sm-mini-t">{{一步，12-26 字}}</p></div>
      <div class="sm-step"><span class="sm-n">2</span><p class="sm-mini-t">{{一步，12-26 字}}</p></div>
      <div class="sm-step"><span class="sm-n">3</span><p class="sm-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sm-kicker, h2, mt-m, grid, g2, mt-l, lede, row, sm-tag, sm-card, sm-step, sm-n, sm-mini-t, deck-footer, slide-number, notes

---

## metrics（硬数字）
指纹：chart
数量：sm-stat=3

用途：三个关键数据。超大加粗数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="sm-kicker">{{数据语境，如 三年复盘}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:48px;margin-top:48px">
    <div class="sm-stat"><div class="sm-stat-v">{{数值 ≤6 字符}}<span class="sm-stat-u">{{单位}}</span></div><div class="sm-stat-l">{{指标名，≤8 字}}</div><p class="sm-stat-note">{{口径，14-30 字}}</p></div>
    <div class="sm-stat"><div class="sm-stat-v">{{数值 ≤6 字符}}<span class="sm-stat-u">{{单位}}</span></div><div class="sm-stat-l">{{指标名，≤8 字}}</div><p class="sm-stat-note">{{口径，14-30 字}}</p></div>
    <div class="sm-stat"><div class="sm-stat-v">{{数值 ≤6 字符}}<span class="sm-stat-u">{{单位}}</span></div><div class="sm-stat-l">{{指标名，≤8 字}}</div><p class="sm-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:16px;color:#666666;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sm-kicker, h2, mt-m, grid, g3, mt-l, sm-stat, sm-stat-v, sm-stat-u, sm-stat-l, sm-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（宣言引文）
指纹：quote

用途：整页一句引文。加粗大字 + 出处 + 两个支撑标签。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；标签各 ≤8 字。

```html
<section class="slide full" data-layout="quote">
  <p class="sm-kicker">{{语境，如 讲义扉页}}</p>
  <p class="sm-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="sm-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="sm-tag sm-tag-invert">{{支撑点 1，≤8 字}}</span>
    <span class="sm-tag">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, sm-kicker, sm-quote, mt-l, sm-src, mt-m, row, sm-tag, sm-tag-invert, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 加粗大字章节名 + 一个过渡问题 + 两个看点标签。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；标签各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="sm-kicker">{{进度，如 第二讲 · 留白与密度}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:46ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="sm-tag sm-tag-invert">{{看点 1，≤8 字}}</span>
    <span class="sm-tag">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, sm-kicker, h1, mt-m, lede, mt-l, row, sm-tag, sm-tag-invert, deck-footer, slide-number, notes

---

## moments（流程时间线）
指纹：chart
数量：sm-tl-item=4

用途：3-4 个节点的横向时间线：方块节点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="sm-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="sm-tl mt-l" style="margin-top:52px">
    <div class="sm-tl-item"><div class="sm-tl-dot">1</div><div class="sm-tl-t">{{时间点，≤10 字符}}</div><p class="sm-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sm-tl-item"><div class="sm-tl-dot">2</div><div class="sm-tl-t">{{时间点，≤10 字符}}</div><p class="sm-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sm-tl-item"><div class="sm-tl-dot">3</div><div class="sm-tl-t">{{时间点，≤10 字符}}</div><p class="sm-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sm-tl-item"><div class="sm-tl-dot">4</div><div class="sm-tl-t">{{时间点，≤10 字符}}</div><p class="sm-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:16px;color:#666666;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sm-kicker, h2, mt-m, sm-tl, mt-l, sm-tl-item, sm-tl-dot, sm-tl-t, sm-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（收尾号召）
指纹：hero

用途：收尾页。加粗大字 + 一句行动提醒 + 黑底按钮 + 描边标签。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；标签 ≤12 字。

```html
<section class="slide full" data-layout="closing">
  <p class="sm-kicker">{{提醒语境，如 下一场预告}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:46ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:24px">
    <span class="sm-btn">{{按钮文案，≤8 字}}</span>
    <span class="sm-tag">{{次级信息，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, sm-kicker, h1, mt-m, lede, mt-l, row, sm-btn, sm-tag, deck-footer, slide-number, notes

---

## code（代码示例）
指纹：code

用途：一段真实代码。终端窗口母题：黑栏三圆点 + 文件名 + 等宽代码块，配语言标签与一句说明。
适用 role：code。
内容约束：代码 6-10 行、每行 ≤56 字符；文件名 ≤18 字符；说明 18-40 字；代码必须真实可读，禁用伪代码占位。

```html
<section class="slide" data-layout="code">
  <p class="sm-kicker">{{语境，如 实操 · 讲义里的代码}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="sm-term mt-l" style="margin-top:44px">
    <div class="sm-term-bar"><span class="sm-dot" style="background:#FFFFFF"></span><span class="sm-dot" style="background:#8C8C8C"></span><span class="sm-dot" style="background:#4D4D4D"></span><span class="sm-term-file">{{文件名，如 lecture-poster.css}}</span><span class="sm-term-lang">{{语言标签，如 CSS}}</span></div>
    <pre class="sm-codeblock"><code>{{逐行直排代码，6-10 行、每行 ≤56 字符，禁 tab 缩进}}</code></pre>
  </div>
  <p class="mt-m" style="font-size:17px;color:#666666;letter-spacing:.02em">{{一句说明：这段代码解决什么，18-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sm-kicker, h2, mt-m, sm-term, mt-l, sm-term-bar, sm-dot, sm-term-file, sm-term-lang, sm-codeblock, mt-m, deck-footer, slide-number, notes
