# 日光浅 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `sl-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-solarized-light` 作用域生效，骨架里已写全，照抄结构即可。
> 每页暖光与横格纸纹由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（护眼即算术）**：日光蓝 #268BD2 只管强调（题签/大数字/链接/按钮），青 #2AA198
> 只管代码与成功（时间线圆点/语言标签/字符串），黄 #B58900 只管数字与常量点睛。
> 文字只用深青与石板青系（#073642 / #586E75 / #93A1A1），禁纯黑 #000、禁冷灰压暖纸、
> 禁纯白大色块。卡片是方正小圆角纸片，禁粗硬边框；纸面编辑器窗口（sl-term）
> 是 code 版式专属，其他版式不准搬。

---

## cover（日光封面）
指纹：hero

用途：开场页。书签题签 + 大标题 + 蓝青线 + 一句定位，可配时间地点药丸。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；药丸各 ≤14 字。

```html
<section class="slide full" data-layout="cover">
  <p class="sl-kicker">{{题签，≤14 字，如 日光阅读实验室 · 晚课第 3 讲}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="sl-line mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="sl-pill sl-pill-accent">{{关键信息，≤14 字}}</span>
    <span class="sl-pill">{{时间地点，≤14 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, sl-kicker, h1, mt-m, sl-line, mt-s, lede, mt-l, row, sl-pill, sl-pill-accent, deck-footer, slide-number, notes

---

## contents（讲次目录）
指纹：table
数量：sl-item=4

用途：议程页。一块纸黄大卡里放 4 行条目：等宽编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="sl-kicker">{{引导语，如 晚课目录}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="sl-card mt-l" style="margin-top:44px">
    <div class="sl-item"><span class="sl-n">01</span><span class="sl-t">{{篇名，≤8 字}}</span><span class="sl-d">{{说明，14-26 字}}</span></div>
    <div class="sl-item"><span class="sl-n">02</span><span class="sl-t">{{篇名，≤8 字}}</span><span class="sl-d">{{说明，14-26 字}}</span></div>
    <div class="sl-item"><span class="sl-n">03</span><span class="sl-t">{{篇名，≤8 字}}</span><span class="sl-d">{{说明，14-26 字}}</span></div>
    <div class="sl-item"><span class="sl-n">04</span><span class="sl-t">{{篇名，≤8 字}}</span><span class="sl-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sl-kicker, h2, mt-m, sl-card, mt-l, sl-item, sl-n, sl-t, sl-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：sl-card=3

用途：恰好三张纸黄卡。每张：彩色题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="sl-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="sl-card sl-card-top"><span class="sl-pill sl-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#586E75">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="sl-card sl-card-top"><span class="sl-pill sl-pill-teal">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#586E75">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="sl-card sl-card-top"><span class="sl-pill">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#586E75">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sl-kicker, h2, mt-m, grid, g3, mt-l, sl-card, sl-card-top, sl-pill, sl-pill-accent, sl-pill-teal, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：sl-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边纸黄卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="sl-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#586E75">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="sl-pill sl-pill-accent">{{要点 1，≤8 字}}</span>
        <span class="sl-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="sl-card">
      <div class="sl-step"><span class="sl-n">1</span><p class="sl-mini-t">{{一步，12-26 字}}</p></div>
      <div class="sl-step"><span class="sl-n">2</span><p class="sl-mini-t">{{一步，12-26 字}}</p></div>
      <div class="sl-step"><span class="sl-n">3</span><p class="sl-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sl-kicker, h2, mt-m, grid, g2, mt-l, lede, row, sl-pill, sl-pill-accent, sl-card, sl-step, sl-n, sl-mini-t, deck-footer, slide-number, notes

---

## metrics（实证数字）
指纹：chart
数量：sl-stat=3

用途：三个关键数据。日光蓝等宽大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="sl-kicker">{{数据语境，如 配色科学 · 基础三数}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="sl-stat"><div class="sl-stat-v">{{数值 ≤6 字符}}<span class="sl-stat-u">{{单位}}</span></div><div class="sl-stat-l">{{指标名，≤8 字}}</div><p class="sl-stat-note">{{口径，14-30 字}}</p></div>
    <div class="sl-stat"><div class="sl-stat-v">{{数值 ≤6 字符}}<span class="sl-stat-u">{{单位}}</span></div><div class="sl-stat-l">{{指标名，≤8 字}}</div><p class="sl-stat-note">{{口径，14-30 字}}</p></div>
    <div class="sl-stat"><div class="sl-stat-v">{{数值 ≤6 字符}}<span class="sl-stat-u">{{单位}}</span></div><div class="sl-stat-l">{{指标名，≤8 字}}</div><p class="sl-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#93A1A1;letter-spacing:.04em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sl-kicker, h2, mt-m, grid, g3, mt-l, sl-stat, sl-stat-v, sl-stat-u, sl-stat-l, sl-stat-note, deck-footer, slide-number, notes

---

## quote（题记引文）
指纹：quote

用途：整页一句引文。深青大字引文 + 出处 + 两个支撑药丸。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="sl-kicker">{{语境，如 晚课开场白}}</p>
  <p class="sl-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="sl-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="sl-pill sl-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="sl-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sl-kicker, sl-quote, mt-l, sl-src, mt-m, row, sl-pill, sl-pill-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="sl-kicker">{{进度，如 第二段 · 对比度的算术}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="sl-pill sl-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="sl-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, sl-kicker, h1, mt-m, lede, mt-l, row, sl-pill, sl-pill-accent, deck-footer, slide-number, notes

---

## moments（编年时间线）
指纹：chart
数量：sl-tl-item=4

用途：3-4 个节点的横向时间线：青色圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="sl-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="sl-tl mt-l" style="margin-top:52px">
    <div class="sl-tl-item"><div class="sl-tl-dot">1</div><div class="sl-tl-t">{{时间点，≤10 字符}}</div><p class="sl-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sl-tl-item"><div class="sl-tl-dot">2</div><div class="sl-tl-t">{{时间点，≤10 字符}}</div><p class="sl-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sl-tl-item"><div class="sl-tl-dot">3</div><div class="sl-tl-t">{{时间点，≤10 字符}}</div><p class="sl-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sl-tl-item"><div class="sl-tl-dot">4</div><div class="sl-tl-t">{{时间点，≤10 字符}}</div><p class="sl-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#93A1A1;letter-spacing:.04em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sl-kicker, h2, mt-m, sl-tl, mt-l, sl-tl-item, sl-tl-dot, sl-tl-t, sl-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（收尾行动）
指纹：hero

用途：收尾页。大字标题 + 一句行动提醒 + 日光蓝实底按钮 + 药丸。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤14 字。

```html
<section class="slide full" data-layout="closing">
  <p class="sl-kicker">{{提醒语境，如 今晚三件事}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:22px">
    <span class="sl-btn">{{按钮文案，≤8 字}}</span>
    <span class="sl-pill">{{次级信息，≤14 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, sl-kicker, h1, mt-m, lede, mt-l, row, sl-btn, sl-pill, deck-footer, slide-number, notes

---

## code（代码示例）
指纹：code

用途：一段真实代码/配置放进纸面编辑器窗口：标题栏三圆点 + 文件名 + 语言标签 + 等宽代码块 + 一句说明。
适用 role：code。
内容约束：代码 6-10 行、每行 ≤56 字符；文件名带扩展名 ≤20 字符；语言标签 ≤10 字符；说明一句 14-30 字。
代码行可用 `<span class="sl-kw">`（关键字）、`sl-str`（字符串）、`sl-fn`（函数）、`sl-num`（常量/数字）、`sl-cm`（注释）标语法色，不用则整块默认石板青。

```html
<section class="slide" data-layout="code">
  <p class="sl-kicker">{{语境，如 高亮实例}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="sl-term mt-l" style="margin-top:38px">
    <div class="sl-bar"><span class="sl-dot" style="background:#dc322f"></span><span class="sl-dot" style="background:#b58900"></span><span class="sl-dot" style="background:#859900"></span><span class="sl-file">{{文件名，如 palette_demo.py}}</span><span class="sl-lang">{{语言标签，≤10 字符}}</span></div>
    <pre class="sl-codeblock"><code>{{代码 6-10 行、每行 ≤56 字符}}</code></pre>
  </div>
  <p class="mt-m" style="font-size:18px;line-height:1.7;color:#586E75">{{一句说明，14-30 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sl-kicker, h2, mt-m, sl-term, mt-l, sl-bar, sl-dot, sl-file, sl-lang, sl-codeblock, sl-kw, sl-str, sl-fn, sl-cm, sl-num, deck-footer, slide-number, notes
