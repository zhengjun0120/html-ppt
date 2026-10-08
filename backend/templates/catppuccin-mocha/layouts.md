# Catppuccin 摩卡 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `mh-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-catppuccin-mocha` 作用域生效，骨架里已写全，照抄结构即可。
> 每页柔光圆点（左上与右下两簇）由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（三颗糖果各管一摊）**：薰衣草只管结构（题签/编号/大数字/按钮），青只管路径与
> 成功（时间线圆点/语言标签），粉只管字符串与提醒（语法色/药丸点睛，一屏 ≤2 处）。
> 禁硬阴影、禁纯黑背景、禁高饱和荧光。h1/h2 大字每页最多一组；
> 终端窗口（mh-term）是 code 版式专属，其他版式不准搬。

---

## cover（摩卡封面）
指纹：hero

用途：开场页。微光题签 + 大标题 + 渐变线 + 一句定位，可配时间地点药丸。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；药丸各 ≤14 字。

```html
<section class="slide full" data-layout="cover">
  <p class="mh-kicker">{{题签，≤14 字，如 终端自习室 · 夜谈第 7 期}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="mh-line mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="mh-pill mh-pill-accent">{{关键信息，≤14 字}}</span>
    <span class="mh-pill">{{时间地点，≤14 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, mh-kicker, h1, mt-m, mh-line, mt-s, lede, mt-l, row, mh-pill, mh-pill-accent, deck-footer, slide-number, notes

---

## contents（议程列表）
指纹：table
数量：mh-item=4

用途：议程页。一块暗色大卡里放 4 行条目：等宽编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="mh-kicker">{{引导语，如 今晚议程}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="mh-card mt-l" style="margin-top:44px">
    <div class="mh-item"><span class="mh-n">01</span><span class="mh-t">{{篇名，≤8 字}}</span><span class="mh-d">{{说明，14-26 字}}</span></div>
    <div class="mh-item"><span class="mh-n">02</span><span class="mh-t">{{篇名，≤8 字}}</span><span class="mh-d">{{说明，14-26 字}}</span></div>
    <div class="mh-item"><span class="mh-n">03</span><span class="mh-t">{{篇名，≤8 字}}</span><span class="mh-d">{{说明，14-26 字}}</span></div>
    <div class="mh-item"><span class="mh-n">04</span><span class="mh-t">{{篇名，≤8 字}}</span><span class="mh-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mh-kicker, h2, mt-m, mh-card, mt-l, mh-item, mh-n, mh-t, mh-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：mh-card=3

用途：恰好三张暗色卡。每张：彩色题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="mh-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="mh-card mh-card-top"><span class="mh-pill mh-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#A6ADC8">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="mh-card mh-card-top"><span class="mh-pill mh-pill-teal">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#A6ADC8">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="mh-card mh-card-top"><span class="mh-pill">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#A6ADC8">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mh-kicker, h2, mt-m, grid, g3, mt-l, mh-card, mh-card-top, mh-pill, mh-pill-accent, mh-pill-teal, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：mh-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边暗色卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="mh-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#A6ADC8">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="mh-pill mh-pill-accent">{{要点 1，≤8 字}}</span>
        <span class="mh-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="mh-card">
      <div class="mh-step"><span class="mh-n">1</span><p class="mh-mini-t">{{一步，12-26 字}}</p></div>
      <div class="mh-step"><span class="mh-n">2</span><p class="mh-mini-t">{{一步，12-26 字}}</p></div>
      <div class="mh-step"><span class="mh-n">3</span><p class="mh-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mh-kicker, h2, mt-m, grid, g2, mt-l, lede, row, mh-pill, mh-pill-accent, mh-card, mh-step, mh-n, mh-mini-t, deck-footer, slide-number, notes

---

## metrics（柔光数字）
指纹：chart
数量：mh-stat=3

用途：三个关键数据。薰衣草等宽大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="mh-kicker">{{数据语境，如 社区问卷 · 2025}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="mh-stat"><div class="mh-stat-v">{{数值 ≤6 字符}}<span class="mh-stat-u">{{单位}}</span></div><div class="mh-stat-l">{{指标名，≤8 字}}</div><p class="mh-stat-note">{{口径，14-30 字}}</p></div>
    <div class="mh-stat"><div class="mh-stat-v">{{数值 ≤6 字符}}<span class="mh-stat-u">{{单位}}</span></div><div class="mh-stat-l">{{指标名，≤8 字}}</div><p class="mh-stat-note">{{口径，14-30 字}}</p></div>
    <div class="mh-stat"><div class="mh-stat-v">{{数值 ≤6 字符}}<span class="mh-stat-u">{{单位}}</span></div><div class="mh-stat-l">{{指标名，≤8 字}}</div><p class="mh-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#7F849C;letter-spacing:.04em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mh-kicker, h2, mt-m, grid, g3, mt-l, mh-stat, mh-stat-v, mh-stat-u, mh-stat-l, mh-stat-note, deck-footer, slide-number, notes

---

## quote（柔光引文）
指纹：quote

用途：整页一句引文。大字引文 + 出处 + 两个支撑药丸。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="mh-kicker">{{语境，如 一位听众的留言}}</p>
  <p class="mh-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="mh-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="mh-pill mh-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="mh-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mh-kicker, mh-quote, mt-l, mh-src, mt-m, row, mh-pill, mh-pill-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="mh-kicker">{{进度，如 第二段 · 摩卡的原则}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="mh-pill mh-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="mh-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, mh-kicker, h1, mt-m, lede, mt-l, row, mh-pill, mh-pill-accent, deck-footer, slide-number, notes

---

## moments（时间线）
指纹：chart
数量：mh-tl-item=4

用途：3-4 个节点的横向时间线：柔光圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="mh-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="mh-tl mt-l" style="margin-top:52px">
    <div class="mh-tl-item"><div class="mh-tl-dot">1</div><div class="mh-tl-t">{{时间点，≤10 字符}}</div><p class="mh-tl-d">{{事件，12-26 字}}</p></div>
    <div class="mh-tl-item"><div class="mh-tl-dot">2</div><div class="mh-tl-t">{{时间点，≤10 字符}}</div><p class="mh-tl-d">{{事件，12-26 字}}</p></div>
    <div class="mh-tl-item"><div class="mh-tl-dot">3</div><div class="mh-tl-t">{{时间点，≤10 字符}}</div><p class="mh-tl-d">{{事件，12-26 字}}</p></div>
    <div class="mh-tl-item"><div class="mh-tl-dot">4</div><div class="mh-tl-t">{{时间点，≤10 字符}}</div><p class="mh-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#7F849C;letter-spacing:.04em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mh-kicker, h2, mt-m, mh-tl, mt-l, mh-tl-item, mh-tl-dot, mh-tl-t, mh-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（收尾行动）
指纹：hero

用途：收尾页。大字标题 + 一句行动提醒 + 薰衣草实底按钮 + 药丸。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤14 字。

```html
<section class="slide full" data-layout="closing">
  <p class="mh-kicker">{{提醒语境，如 今晚三件事}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:22px">
    <span class="mh-btn">{{按钮文案，≤8 字}}</span>
    <span class="mh-pill">{{次级信息，≤14 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, mh-kicker, h1, mt-m, lede, mt-l, row, mh-btn, mh-pill, deck-footer, slide-number, notes

---

## code（代码示例）
指纹：code

用途：一段真实配置/代码放进终端窗口：标题栏三圆点 + 文件名 + 语言标签 + 等宽代码块 + 一句说明。
适用 role：code。
内容约束：代码 6-10 行、每行 ≤56 字符；文件名带扩展名 ≤20 字符；语言标签 ≤10 字符；说明一句 14-30 字。
代码行可用 `<span class="mh-kw">`（结构/关键字）、`mh-str`（字符串）、`mh-fn`（函数/成功）、`mh-num`（数字）、`mh-cm`（注释）标语法色，不用则整块默认灰蓝。

```html
<section class="slide" data-layout="code">
  <p class="mh-kicker">{{语境，如 实操演示}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="mh-term mt-l" style="margin-top:38px">
    <div class="mh-bar"><span class="mh-dot" style="background:#f38ba8"></span><span class="mh-dot" style="background:#f9e2af"></span><span class="mh-dot" style="background:#a6e3a1"></span><span class="mh-file">{{文件名，如 starship.toml}}</span><span class="mh-lang">{{语言标签，≤10 字符}}</span></div>
    <pre class="mh-codeblock"><code>{{代码 6-10 行、每行 ≤56 字符}}</code></pre>
  </div>
  <p class="mt-m" style="font-size:18px;line-height:1.7;color:#A6ADC8">{{一句说明，14-30 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mh-kicker, h2, mt-m, mh-term, mt-l, mh-bar, mh-dot, mh-file, mh-lang, mh-codeblock, mh-kw, mh-str, mh-fn, mh-cm, mh-num, deck-footer, slide-number, notes
