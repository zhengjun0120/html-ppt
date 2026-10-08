# 终端绿 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `tg-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-terminal-green` 作用域生效，骨架里已写全，照抄结构即可。
> 每页扫描线与磷光辉光由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（纯绿磷光）**：全模板只允许一种彩色——绿。标题、强调、数据、按钮一律是
> #4ade80 系磷光绿，禁出现第二种色相。辉光只给标题、大数字与光标，正文不加发光；
> h1 / h2 大字每页最多一组；等宽是唯一字体，不要在骨架外引入别的字体栈。

---

## cover（终端封面）
指纹：hero

用途：开场页。$ 提示符题签 + 磷光大标题 + 一行命令与光标 + 一句定位 + 描边标签。
适用 role：cover。
内容约束：主标题 ≤12 字；lede 24-48 字；命令 ≤40 字符；标签各 ≤8 字。

```html
<section class="slide full" data-layout="cover">
  <p class="tg-kicker">{{题签，≤14 字，如 v1.0.0 · 发布分享}}</p>
  <h1 class="h1 mt-m">{{主标题，≤12 字，可 <br> 分行}}</h1>
  <p class="tg-cmdline mt-s">{{一行命令，≤40 字符}}<span class="tg-cursor"></span></p>
  <p class="lede mt-m" style="max-width:52ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="tg-pill tg-pill-accent">{{标签 1，≤8 字}}</span>
    <span class="tg-pill">{{标签 2，≤8 字}}</span>
    <span class="tg-pill">{{标签 3，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, tg-kicker, h1, mt-m, tg-cmdline, mt-s, tg-cursor, lede, mt-l, row, tg-pill, tg-pill-accent, deck-footer, slide-number, notes

---

## contents（目录页）
指纹：table
数量：tg-item=4

用途：议程页。一块终端窗卡里放 4 行篇目：等宽编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="tg-kicker">{{引导语，≤14 字，如 cat AGENDA.md}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="tg-card mt-l" style="margin-top:44px">
    <div class="tg-item"><span class="tg-n">01</span><span class="tg-t">{{篇名，≤8 字}}</span><span class="tg-d">{{说明，14-26 字}}</span></div>
    <div class="tg-item"><span class="tg-n">02</span><span class="tg-t">{{篇名，≤8 字}}</span><span class="tg-d">{{说明，14-26 字}}</span></div>
    <div class="tg-item"><span class="tg-n">03</span><span class="tg-t">{{篇名，≤8 字}}</span><span class="tg-d">{{说明，14-26 字}}</span></div>
    <div class="tg-item"><span class="tg-n">04</span><span class="tg-t">{{篇名，≤8 字}}</span><span class="tg-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, tg-kicker, h2, mt-m, tg-card, mt-l, tg-item, tg-n, tg-t, tg-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：tg-card=3

用途：恰好三张终端窗卡。每张：磷光标签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；标签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="tg-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="tg-card"><span class="tg-pill tg-pill-accent">{{标签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#7fc98e">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="tg-card"><span class="tg-pill tg-pill-accent">{{标签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#7fc98e">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="tg-card"><span class="tg-pill tg-pill-accent">{{标签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#7fc98e">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, tg-kicker, h2, mt-m, grid, g3, mt-l, tg-card, tg-pill, tg-pill-accent, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：tg-step=3

用途：左边把一件事讲透（lede + 补充 + 标签），右边终端卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个标签；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="tg-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#7fc98e">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="tg-pill">{{要点 1，≤8 字}}</span>
        <span class="tg-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="tg-card">
      <div class="tg-step"><span class="tg-n">01</span><p class="tg-mini-t">{{一步，12-26 字}}</p></div>
      <div class="tg-step"><span class="tg-n">02</span><p class="tg-mini-t">{{一步，12-26 字}}</p></div>
      <div class="tg-step"><span class="tg-n">03</span><p class="tg-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, tg-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, tg-pill, tg-card, tg-step, tg-n, tg-mini-t, deck-footer, slide-number, notes

---

## metrics（磷光数字）
指纹：chart
数量：tg-stat=3

用途：三个关键数据。磷光绿大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="tg-kicker">{{数据语境，≤14 字}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="tg-stat"><div class="tg-stat-v">{{数值 ≤6 字符}}<span class="tg-stat-u">{{单位}}</span></div><div class="tg-stat-l">{{指标名，≤8 字}}</div><p class="tg-stat-note">{{口径，14-30 字}}</p></div>
    <div class="tg-stat"><div class="tg-stat-v">{{数值 ≤6 字符}}<span class="tg-stat-u">{{单位}}</span></div><div class="tg-stat-l">{{指标名，≤8 字}}</div><p class="tg-stat-note">{{口径，14-30 字}}</p></div>
    <div class="tg-stat"><div class="tg-stat-v">{{数值 ≤6 字符}}<span class="tg-stat-u">{{单位}}</span></div><div class="tg-stat-l">{{指标名，≤8 字}}</div><p class="tg-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#63a86c;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, tg-kicker, h2, mt-m, grid, g3, mt-l, tg-stat, tg-stat-v, tg-stat-u, tg-stat-l, tg-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（终端引言）
指纹：quote

用途：整页一句引文。等宽大字 + 出处 + 两个支撑标签，像一条被钉住的 commit message。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；标签各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="tg-kicker">{{语境，≤14 字，如 git log -1 v1.0.0}}</p>
  <p class="tg-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="tg-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="tg-pill tg-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="tg-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, tg-kicker, tg-quote, mt-l, tg-src, mt-m, row, tg-pill, tg-pill-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。$ 进度题签 + 磷光大字章节名 + 一个过渡问题 + 两个看点标签。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；标签各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="tg-kicker">{{进度，≤14 字，如 第二章 · 解析设计}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="tg-pill tg-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="tg-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, tg-kicker, h1, mt-m, lede, mt-l, row, tg-pill, tg-pill-accent, deck-footer, slide-number, notes

---

## moments（版本时间线）
指纹：chart
数量：tg-tl-item=4

用途：3-4 个节点的横向时间线：发光圆点 + 版本/时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="tg-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="tg-tl mt-l" style="margin-top:52px">
    <div class="tg-tl-item"><div class="tg-tl-dot"></div><div class="tg-tl-t">{{时间点，≤10 字符}}</div><p class="tg-tl-d">{{事件，12-26 字}}</p></div>
    <div class="tg-tl-item"><div class="tg-tl-dot"></div><div class="tg-tl-t">{{时间点，≤10 字符}}</div><p class="tg-tl-d">{{事件，12-26 字}}</p></div>
    <div class="tg-tl-item"><div class="tg-tl-dot"></div><div class="tg-tl-t">{{时间点，≤10 字符}}</div><p class="tg-tl-d">{{事件，12-26 字}}</p></div>
    <div class="tg-tl-item"><div class="tg-tl-dot"></div><div class="tg-tl-t">{{时间点，≤10 字符}}</div><p class="tg-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#63a86c;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, tg-kicker, h2, mt-m, tg-tl, mt-l, tg-tl-item, tg-tl-dot, tg-tl-t, tg-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（收尾页）
指纹：hero

用途：收尾页。磷光大字 + 一句行动提醒 + 描边按钮 + 标签，像一条欢迎加入的 README。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；标签 ≤12 字。

```html
<section class="slide full" data-layout="closing">
  <p class="tg-kicker">{{提醒语境，≤14 字，如 brew install xingyun}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:22px">
    <span class="tg-btn">{{按钮文案，≤8 字}}</span>
    <span class="tg-pill">{{次级信息，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, tg-kicker, h1, mt-m, lede, mt-l, row, tg-btn, tg-pill, deck-footer, slide-number, notes

---

## code（代码示例）
指纹：code

用途：一段真实代码或命令。终端窗承载：标题栏三圆点 + 文件名 + 语言标签 + 等宽代码块，代码是本页唯一主角。
适用 role：code。
内容约束：代码 6-10 行、每行 ≤56 字符（< > & 需转义为 &lt; &gt; &amp;）；文件名 ≤16 字符；配一句说明 18-40 字。

```html
<section class="slide" data-layout="code">
  <p class="tg-kicker">{{语境，≤14 字，如 internal/options.go}}</p>
  <h2 class="h2 mt-m">{{这段代码回答什么，≤12 字}}</h2>
  <div class="tg-term mt-l" style="margin-top:40px">
    <div class="tg-term-bar"><span class="tg-dot"></span><span class="tg-dot"></span><span class="tg-dot"></span><span class="tg-term-file">{{文件名，≤16 字符}}</span><span class="tg-lang">{{语言，如 Go}}</span></div>
    <pre class="tg-codeblock"><code>{{6-10 行代码，直排文本行，每行 ≤56 字符，&lt; &gt; &amp; 需转义}}</code></pre>
  </div>
  <p class="tg-codenote mt-m">{{一句说明：这段代码在讲什么，18-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, tg-kicker, h2, mt-m, tg-term, mt-l, tg-term-bar, tg-dot, tg-term-file, tg-lang, tg-codeblock, tg-codenote, deck-footer, slide-number, notes
