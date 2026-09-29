# 北欧 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `no-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-nord` 作用域生效，骨架里已写全，照抄结构即可。
> 每页顶部极光薄雾与四周霜蓝细框由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（克制即身份）**：全模板唯一的强调色是冰蓝 #88C0D0（钢蓝 #81A1C1 只做次级），
> 出现在题签发丝、编号、大数字、时间点、按钮与 accent 标签；正文一律浅灰，禁用霓虹与辉光。
> 轻字重是气质的一部分：h1/h2 已内建 600 字重、大数字 300 字重，不要再加粗。
> h1 / h2 大字每页最多一组；一页最多两组内容块 + 页脚。

---

## cover（极夜封面）
指纹：hero

用途：开场页。冰蓝发丝题签 + 极浅灰白大标题 + 霜蓝细线 + 一句定位 + 描边标签。
适用 role：cover。
内容约束：主标题 ≤12 字；lede 24-48 字；标签各 ≤8 字。

```html
<section class="slide full" data-layout="cover">
  <p class="no-kicker">{{题签，≤14 字，如 极简工作流 · 分享}}</p>
  <h1 class="h1 mt-m">{{主标题，≤12 字，可 <br> 分行}}</h1>
  <div class="no-line mt-s"></div>
  <p class="lede mt-m" style="max-width:52ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="no-pill no-pill-accent">{{标签 1，≤8 字}}</span>
    <span class="no-pill">{{标签 2，≤8 字}}</span>
    <span class="no-pill">{{标签 3，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, no-kicker, h1, mt-m, no-line, mt-s, lede, mt-l, row, no-pill, no-pill-accent, deck-footer, slide-number, notes

---

## contents（目录页）
指纹：table
数量：no-item=4

用途：议程页。一块暗色卡里放 4 行篇目：等宽编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="no-kicker">{{引导语，≤14 字}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="no-card mt-l" style="margin-top:44px">
    <div class="no-item"><span class="no-n">01</span><span class="no-t">{{篇名，≤8 字}}</span><span class="no-d">{{说明，14-26 字}}</span></div>
    <div class="no-item"><span class="no-n">02</span><span class="no-t">{{篇名，≤8 字}}</span><span class="no-d">{{说明，14-26 字}}</span></div>
    <div class="no-item"><span class="no-n">03</span><span class="no-t">{{篇名，≤8 字}}</span><span class="no-d">{{说明，14-26 字}}</span></div>
    <div class="no-item"><span class="no-n">04</span><span class="no-t">{{篇名，≤8 字}}</span><span class="no-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, no-kicker, h2, mt-m, no-card, mt-l, no-item, no-n, no-t, no-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：no-card=3

用途：恰好三张暗色卡。每张：冰蓝标签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；标签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="no-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="no-card"><span class="no-pill no-pill-accent">{{标签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#d8dee9">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="no-card"><span class="no-pill no-pill-accent">{{标签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#d8dee9">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="no-card"><span class="no-pill no-pill-accent">{{标签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#d8dee9">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, no-kicker, h2, mt-m, grid, g3, mt-l, no-card, no-pill, no-pill-accent, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：no-step=3

用途：左边把一件事讲透（lede + 补充 + 键牌），右边暗色卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个键牌；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="no-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#d8dee9">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="no-key">{{快捷键 1，≤10 字符}}</span>
        <span class="no-key">{{快捷键 2，≤10 字符}}</span>
      </div>
    </div>
    <div class="no-card">
      <div class="no-step"><span class="no-n">01</span><p class="no-mini-t">{{一步，12-26 字}}</p></div>
      <div class="no-step"><span class="no-n">02</span><p class="no-mini-t">{{一步，12-26 字}}</p></div>
      <div class="no-step"><span class="no-n">03</span><p class="no-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, no-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, no-key, no-card, no-step, no-n, no-mini-t, deck-footer, slide-number, notes

---

## metrics（冰蓝数字）
指纹：chart
数量：no-stat=3

用途：三个关键数据。冰蓝轻字重大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="no-kicker">{{数据语境，≤14 字}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="no-stat"><div class="no-stat-v">{{数值 ≤6 字符}}<span class="no-stat-u">{{单位}}</span></div><div class="no-stat-l">{{指标名，≤8 字}}</div><p class="no-stat-note">{{口径，14-30 字}}</p></div>
    <div class="no-stat"><div class="no-stat-v">{{数值 ≤6 字符}}<span class="no-stat-u">{{单位}}</span></div><div class="no-stat-l">{{指标名，≤8 字}}</div><p class="no-stat-note">{{口径，14-30 字}}</p></div>
    <div class="no-stat"><div class="no-stat-v">{{数值 ≤6 字符}}<span class="no-stat-u">{{单位}}</span></div><div class="no-stat-l">{{指标名，≤8 字}}</div><p class="no-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#b3bfcd;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, no-kicker, h2, mt-m, grid, g3, mt-l, no-stat, no-stat-v, no-stat-u, no-stat-l, no-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（极光引文）
指纹：quote

用途：整页一句引文。极浅灰白大字 + 出处 + 两个支撑标签，安静得像极夜。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；标签各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="no-kicker">{{语境，≤14 字}}</p>
  <p class="no-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="no-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="no-pill no-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="no-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, no-kicker, no-quote, mt-l, no-src, mt-m, row, no-pill, no-pill-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。冰蓝发丝题签 + 大字章节名 + 一个过渡问题 + 两个看点标签。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；标签各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="no-kicker">{{进度，≤14 字，如 二 · 取舍}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="no-pill no-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="no-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, no-kicker, h1, mt-m, lede, mt-l, row, no-pill, no-pill-accent, deck-footer, slide-number, notes

---

## moments（演进时间线）
指纹：chart
数量：no-tl-item=4

用途：3-4 个节点的横向时间线：冰蓝圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="no-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="no-tl mt-l" style="margin-top:52px">
    <div class="no-tl-item"><div class="no-tl-dot"></div><div class="no-tl-t">{{时间点，≤10 字符}}</div><p class="no-tl-d">{{事件，12-26 字}}</p></div>
    <div class="no-tl-item"><div class="no-tl-dot"></div><div class="no-tl-t">{{时间点，≤10 字符}}</div><p class="no-tl-d">{{事件，12-26 字}}</p></div>
    <div class="no-tl-item"><div class="no-tl-dot"></div><div class="no-tl-t">{{时间点，≤10 字符}}</div><p class="no-tl-d">{{事件，12-26 字}}</p></div>
    <div class="no-tl-item"><div class="no-tl-dot"></div><div class="no-tl-t">{{时间点，≤10 字符}}</div><p class="no-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#b3bfcd;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, no-kicker, h2, mt-m, no-tl, mt-l, no-tl-item, no-tl-dot, no-tl-t, no-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（收尾页）
指纹：hero

用途：收尾页。大字 + 一句行动提醒 + 冰蓝描边按钮 + 标签。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；标签 ≤12 字。

```html
<section class="slide full" data-layout="closing">
  <p class="no-kicker">{{提醒语境，≤14 字}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:22px">
    <span class="no-btn">{{按钮文案，≤8 字}}</span>
    <span class="no-pill">{{次级信息，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, no-kicker, h1, mt-m, lede, mt-l, row, no-btn, no-pill, deck-footer, slide-number, notes

---

## code（代码示例）
指纹：code

用途：一段真实配置或代码。编辑器窗承载：标题栏三圆点 + 文件名 + 语言标签 + 等宽代码块，代码是本页唯一主角。
适用 role：code。
内容约束：代码 6-10 行、每行 ≤56 字符（< > & 需转义为 &lt; &gt; &amp;）；文件名 ≤16 字符；配一句说明 18-40 字。

```html
<section class="slide" data-layout="code">
  <p class="no-kicker">{{语境，≤14 字，如 ~/.config/nvim}}</p>
  <h2 class="h2 mt-m">{{这段代码回答什么，≤12 字}}</h2>
  <div class="no-term mt-l" style="margin-top:40px">
    <div class="no-term-bar"><span class="no-dot"></span><span class="no-dot"></span><span class="no-dot"></span><span class="no-term-file">{{文件名，≤16 字符}}</span><span class="no-lang">{{语言，如 Lua}}</span></div>
    <pre class="no-codeblock"><code>{{6-10 行代码，直排文本行，每行 ≤56 字符，&lt; &gt; &amp; 需转义}}</code></pre>
  </div>
  <p class="no-codenote mt-m">{{一句说明：这段代码在讲什么，18-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, no-kicker, h2, mt-m, no-term, mt-l, no-term-bar, no-dot, no-term-file, no-lang, no-codeblock, no-codenote, deck-footer, slide-number, notes
