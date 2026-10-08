# 东京夜 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `tn-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-tokyo-night` 作用域生效，骨架里已写全，照抄结构即可。
> 每页顶部霓虹细线、右上辉光与页底天际线由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
> code 版式是 IDE 终端母题，一份 deck 至多 1-2 页。
>
> **气质铁律（夜的纪律）**：蓝 #7AA2F7 只管主强调——题签细线、关键数字（tn-stat-v）、
> 编号（tn-n）、按钮（tn-btn）与时间线圆点；紫 #BB9AF7 只管强调药丸与时间点
> （tn-pill-accent / tn-tl-t）；青 #7DCFFF 只管终端标签与语言标签（tn-tag / tn-lang）。
> 禁暖色、禁正文对比度不足、禁装饰堆叠（ambient 已内建两层，正文不再加光效）。

---

## cover（夜幕封面）
指纹：hero

用途：开场页。题签 + 大字标题 + 蓝紫渐隐线 + 一句定位，终端标签芯片与药丸收底。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 30-52 字；题签 ≤16 字；终端标签 ≤16 字；药丸 ≤12 字。

```html
<section class="slide full" data-layout="cover">
  <p class="tn-kicker">{{题签，≤16 字，如 深夜课堂 · 异步第一课}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="tn-rule mt-s"></div>
  <p class="lede mt-m" style="max-width:52ch">{{一句话定位，30-52 字}}</p>
  <div class="row mt-l" style="gap:24px">
    <span class="tn-tag">{{终端短语，≤16 字，如 $ node intro.js}}</span>
    <span class="tn-pill">{{副题，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, tn-kicker, h1, mt-m, tn-rule, mt-s, lede, mt-l, row, tn-tag, tn-pill, deck-footer, slide-number, notes

---

## contents（夜行目录）
指纹：table
数量：tn-item=4

用途：议程页。一块暗面板大卡里放 4 行篇目：方括号编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 12-24 字。

```html
<section class="slide" data-layout="contents">
  <p class="tn-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="tn-card mt-l" style="margin-top:44px">
    <div class="tn-item"><span class="tn-n">[01]</span><span class="tn-t">{{篇名，≤8 字}}</span><span class="tn-d">{{说明，12-24 字}}</span></div>
    <div class="tn-item"><span class="tn-n">[02]</span><span class="tn-t">{{篇名，≤8 字}}</span><span class="tn-d">{{说明，12-24 字}}</span></div>
    <div class="tn-item"><span class="tn-n">[03]</span><span class="tn-t">{{篇名，≤8 字}}</span><span class="tn-d">{{说明，12-24 字}}</span></div>
    <div class="tn-item"><span class="tn-n">[04]</span><span class="tn-t">{{篇名，≤8 字}}</span><span class="tn-d">{{说明，12-24 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, tn-kicker, h2, mt-m, tn-card, mt-l, tn-item, tn-n, tn-t, tn-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：tn-card=3

用途：恰好三张暗面板卡。每张：蓝调药丸题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="tn-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="tn-card"><span class="tn-pill tn-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="tn-note mt-s">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="tn-card"><span class="tn-pill tn-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="tn-note mt-s">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="tn-card"><span class="tn-pill tn-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="tn-note mt-s">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, tn-kicker, h2, mt-m, grid, g3, mt-l, tn-card, tn-pill, tn-pill-accent, h4, tn-note, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：tn-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边暗面板装三行步骤。
适用 role：content。
内容约束：左 lede 40-80 字 + 补充 16-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="tn-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，40-80 字}}</p>
      <p class="tn-note mt-m">{{补充：判断标准或代价，16-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="tn-pill">{{要点 1，≤8 字}}</span>
        <span class="tn-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="tn-card">
      <div class="tn-step"><span class="tn-n">[01]</span><p class="tn-mini-t">{{一步，12-26 字}}</p></div>
      <div class="tn-step"><span class="tn-n">[02]</span><p class="tn-mini-t">{{一步，12-26 字}}</p></div>
      <div class="tn-step"><span class="tn-n">[03]</span><p class="tn-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, tn-kicker, h2, mt-m, grid, g2, mt-l, lede, tn-note, row, tn-pill, tn-card, tn-step, tn-n, tn-mini-t, deck-footer, slide-number, notes

---

## metrics（霓虹数字）
指纹：chart
数量：tn-stat=3

用途：三个关键数据。蓝色等宽大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-32 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="tn-kicker">{{数据语境}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="tn-stat"><div class="tn-stat-v">{{数值 ≤6 字符}}<span class="tn-stat-u">{{单位}}</span></div><div class="tn-stat-l">{{指标名，≤8 字}}</div><p class="tn-stat-note">{{口径，14-32 字}}</p></div>
    <div class="tn-stat"><div class="tn-stat-v">{{数值 ≤6 字符}}<span class="tn-stat-u">{{单位}}</span></div><div class="tn-stat-l">{{指标名，≤8 字}}</div><p class="tn-stat-note">{{口径，14-32 字}}</p></div>
    <div class="tn-stat"><div class="tn-stat-v">{{数值 ≤6 字符}}<span class="tn-stat-u">{{单位}}</span></div><div class="tn-stat-l">{{指标名，≤8 字}}</div><p class="tn-stat-note">{{口径，14-32 字}}</p></div>
  </div>
  <p class="tn-src mt-m" style="margin-top:40px">来源：{{出处与统计口径，14-44 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, tn-kicker, h2, mt-m, grid, g3, mt-l, tn-stat, tn-stat-v, tn-stat-u, tn-stat-l, tn-stat-note, tn-src, deck-footer, slide-number, notes

---

## quote（夜半引文）
指纹：quote

用途：整页一句引文。大字 + 出处 + 两个支撑药丸，像深夜签入时读到的一句箴言。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="tn-kicker">{{语境}}</p>
  <p class="tn-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="tn-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="tn-pill tn-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="tn-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, tn-kicker, tn-quote, mt-l, tn-src, mt-m, row, tn-pill, tn-pill-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="tn-kicker">{{进度，如 第二章 · 把回拉直}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="tn-pill tn-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="tn-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, tn-kicker, h1, mt-m, lede, mt-l, row, tn-pill, tn-pill-accent, deck-footer, slide-number, notes

---

## moments（演进时间线）
指纹：chart
数量：tn-tl-item=4

用途：4 个节点的横向时间线：蓝框序号 + 紫色时间点 + 一句事件，讲技术演进最稳。
适用 role：content / data。
内容约束：4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-44 字。

```html
<section class="slide" data-layout="moments">
  <p class="tn-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="tn-tl mt-l" style="margin-top:52px">
    <div class="tn-tl-item"><div class="tn-tl-dot">01</div><div class="tn-tl-t">{{时间点，≤10 字符}}</div><p class="tn-tl-d">{{事件，12-26 字}}</p></div>
    <div class="tn-tl-item"><div class="tn-tl-dot">02</div><div class="tn-tl-t">{{时间点，≤10 字符}}</div><p class="tn-tl-d">{{事件，12-26 字}}</p></div>
    <div class="tn-tl-item"><div class="tn-tl-dot">03</div><div class="tn-tl-t">{{时间点，≤10 字符}}</div><p class="tn-tl-d">{{事件，12-26 字}}</p></div>
    <div class="tn-tl-item"><div class="tn-tl-dot">04</div><div class="tn-tl-t">{{时间点，≤10 字符}}</div><p class="tn-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="tn-src mt-m" style="margin-top:44px">{{一句读法或口径，14-44 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, tn-kicker, h2, mt-m, tn-tl, mt-l, tn-tl-item, tn-tl-dot, tn-tl-t, tn-tl-d, tn-src, deck-footer, slide-number, notes

---

## closing（收尾落款）
指纹：hero

用途：收尾页。大字 + 一句行动提醒 + 蓝调描边按钮 + 药丸 + 终端标签芯片，如一次干净的 git push。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-48 字；按钮 ≤8 字；药丸 ≤12 字；终端标签 ≤16 字。

```html
<section class="slide full" data-layout="closing">
  <p class="tn-kicker">{{提醒语境}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-48 字}}</p>
  <div class="row mt-l" style="gap:24px">
    <span class="tn-btn">{{按钮文案，≤8 字}}</span>
    <span class="tn-pill">{{次级信息，≤12 字}}</span>
    <span class="tn-tag">{{终端短语，≤16 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, tn-kicker, h1, mt-m, lede, mt-l, row, tn-btn, tn-pill, tn-tag, deck-footer, slide-number, notes

---

## code（代码示例）
指纹：code

用途：IDE 终端母题。标题栏三圆点 + 文件名 + 语言标签，等宽代码块直排，状态栏一句话收底。
适用 role：code。
内容约束：代码 6-10 行、每行 ≤56 字符（尖括号与 & 要 HTML 转义）；配文件名 + 语言标签 + 状态栏左右各一句 + 页尾一句说明 18-40 字；一份 deck 至多 1-2 页。

```html
<section class="slide" data-layout="code">
  <p class="tn-kicker">{{语境}}</p>
  <h2 class="h2 mt-m">{{标题，≤14 字}}</h2>
  <div class="tn-term mt-l" style="margin-top:44px">
    <div class="tn-term-bar"><span class="tn-dot"></span><span class="tn-dot"></span><span class="tn-dot"></span><span class="tn-term-name">{{文件名，如 demo/pseudo.js}}</span><span class="tn-lang">{{语言标签，≤8 字符}}</span></div>
    <pre class="tn-codeblock"><code>{{代码直排 6-10 行，每行 ≤56 字符，尖括号与 & 转义}}</code></pre>
    <div class="tn-status"><span>{{左：模式或项目一句}}</span><span>{{右：状态一句话}}</span></div>
  </div>
  <p class="tn-note mt-m" style="margin-top:28px">{{一句说明，18-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, tn-kicker, h2, mt-m, tn-term, mt-l, tn-term-bar, tn-dot, tn-term-name, tn-lang, tn-codeblock, tn-status, tn-note, deck-footer, slide-number, notes
