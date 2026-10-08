# Catppuccin 拿铁 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `lt-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-catppuccin-latte` 作用域生效，骨架里已写全，照抄结构即可。
> 每页奶泡光斑（右上与左下两团）由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（浅色也要有开发者脾气）**：薰衣草只管结构（题签/编号/大数字/按钮），青只管
> 内容与成功（时间线圆点/语言标签/字符串），粉只管提醒（语法色/药丸点睛，一屏 ≤2 处）。
> 文字一律深拿铁灰系，禁纯黑 #000、禁高饱和荧光；卡片是奶白微浮卡，禁死白大色块。
> 浅色编辑器窗口（lt-term）是 code 版式专属，其他版式不准搬。

---

## cover（拿铁封面）
指纹：hero

用途：开场页。色点题签 + 大标题 + 渐变线 + 一句定位，可配时间地点药丸。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；药丸各 ≤14 字。

```html
<section class="slide full" data-layout="cover">
  <p class="lt-kicker">{{题签，≤14 字，如 拿铁主题组 · 设计分享}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="lt-line mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="lt-pill lt-pill-accent">{{关键信息，≤14 字}}</span>
    <span class="lt-pill">{{时间地点，≤14 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, lt-kicker, h1, mt-m, lt-line, mt-s, lede, mt-l, row, lt-pill, lt-pill-accent, deck-footer, slide-number, notes

---

## contents（议程列表）
指纹：table
数量：lt-item=4

用途：议程页。一块奶白大卡里放 4 行条目：等宽编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="lt-kicker">{{引导语，如 今晚议程}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="lt-card mt-l" style="margin-top:44px">
    <div class="lt-item"><span class="lt-n">01</span><span class="lt-t">{{篇名，≤8 字}}</span><span class="lt-d">{{说明，14-26 字}}</span></div>
    <div class="lt-item"><span class="lt-n">02</span><span class="lt-t">{{篇名，≤8 字}}</span><span class="lt-d">{{说明，14-26 字}}</span></div>
    <div class="lt-item"><span class="lt-n">03</span><span class="lt-t">{{篇名，≤8 字}}</span><span class="lt-d">{{说明，14-26 字}}</span></div>
    <div class="lt-item"><span class="lt-n">04</span><span class="lt-t">{{篇名，≤8 字}}</span><span class="lt-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, lt-kicker, h2, mt-m, lt-card, mt-l, lt-item, lt-n, lt-t, lt-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：lt-card=3

用途：恰好三张奶白卡。每张：彩色题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="lt-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="lt-card lt-card-top"><span class="lt-pill lt-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#5C5F77">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="lt-card lt-card-top"><span class="lt-pill lt-pill-teal">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#5C5F77">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="lt-card lt-card-top"><span class="lt-pill">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#5C5F77">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, lt-kicker, h2, mt-m, grid, g3, mt-l, lt-card, lt-card-top, lt-pill, lt-pill-accent, lt-pill-teal, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：lt-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边奶白卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="lt-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#5C5F77">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="lt-pill lt-pill-accent">{{要点 1，≤8 字}}</span>
        <span class="lt-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="lt-card">
      <div class="lt-step"><span class="lt-n">1</span><p class="lt-mini-t">{{一步，12-26 字}}</p></div>
      <div class="lt-step"><span class="lt-n">2</span><p class="lt-mini-t">{{一步，12-26 字}}</p></div>
      <div class="lt-step"><span class="lt-n">3</span><p class="lt-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, lt-kicker, h2, mt-m, grid, g2, mt-l, lede, row, lt-pill, lt-pill-accent, lt-card, lt-step, lt-n, lt-mini-t, deck-footer, slide-number, notes

---

## metrics（数据卡）
指纹：chart
数量：lt-stat=3

用途：三个关键数据。薰衣草等宽大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="lt-kicker">{{数据语境，如 可读性实验 · 2025}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="lt-stat"><div class="lt-stat-v">{{数值 ≤6 字符}}<span class="lt-stat-u">{{单位}}</span></div><div class="lt-stat-l">{{指标名，≤8 字}}</div><p class="lt-stat-note">{{口径，14-30 字}}</p></div>
    <div class="lt-stat"><div class="lt-stat-v">{{数值 ≤6 字符}}<span class="lt-stat-u">{{单位}}</span></div><div class="lt-stat-l">{{指标名，≤8 字}}</div><p class="lt-stat-note">{{口径，14-30 字}}</p></div>
    <div class="lt-stat"><div class="lt-stat-v">{{数值 ≤6 字符}}<span class="lt-stat-u">{{单位}}</span></div><div class="lt-stat-l">{{指标名，≤8 字}}</div><p class="lt-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#8C8FA1;letter-spacing:.04em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, lt-kicker, h2, mt-m, grid, g3, mt-l, lt-stat, lt-stat-v, lt-stat-u, lt-stat-l, lt-stat-note, deck-footer, slide-number, notes

---

## quote（引文页）
指纹：quote

用途：整页一句引文。深灰大字引文 + 出处 + 两个支撑药丸。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="lt-kicker">{{语境，如 发布说明里的一句话}}</p>
  <p class="lt-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="lt-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="lt-pill lt-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="lt-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, lt-kicker, lt-quote, mt-l, lt-src, mt-m, row, lt-pill, lt-pill-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="lt-kicker">{{进度，如 第二段 · 色板设计}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="lt-pill lt-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="lt-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, lt-kicker, h1, mt-m, lede, mt-l, row, lt-pill, lt-pill-accent, deck-footer, slide-number, notes

---

## moments（时间线）
指纹：chart
数量：lt-tl-item=4

用途：3-4 个节点的横向时间线：彩色圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="lt-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="lt-tl mt-l" style="margin-top:52px">
    <div class="lt-tl-item"><div class="lt-tl-dot">1</div><div class="lt-tl-t">{{时间点，≤10 字符}}</div><p class="lt-tl-d">{{事件，12-26 字}}</p></div>
    <div class="lt-tl-item"><div class="lt-tl-dot">2</div><div class="lt-tl-t">{{时间点，≤10 字符}}</div><p class="lt-tl-d">{{事件，12-26 字}}</p></div>
    <div class="lt-tl-item"><div class="lt-tl-dot">3</div><div class="lt-tl-t">{{时间点，≤10 字符}}</div><p class="lt-tl-d">{{事件，12-26 字}}</p></div>
    <div class="lt-tl-item"><div class="lt-tl-dot">4</div><div class="lt-tl-t">{{时间点，≤10 字符}}</div><p class="lt-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#8C8FA1;letter-spacing:.04em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, lt-kicker, h2, mt-m, lt-tl, mt-l, lt-tl-item, lt-tl-dot, lt-tl-t, lt-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（收尾行动）
指纹：hero

用途：收尾页。大字标题 + 一句行动提醒 + 薰衣草实底按钮 + 药丸。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤14 字。

```html
<section class="slide full" data-layout="closing">
  <p class="lt-kicker">{{提醒语境，如 今晚三件事}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:22px">
    <span class="lt-btn">{{按钮文案，≤8 字}}</span>
    <span class="lt-pill">{{次级信息，≤14 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, lt-kicker, h1, mt-m, lede, mt-l, row, lt-btn, lt-pill, deck-footer, slide-number, notes

---

## code（代码示例）
指纹：code

用途：一段真实配置/代码放进浅色编辑器窗口：标题栏三圆点 + 文件名 + 语言标签 + 等宽代码块 + 一句说明。
适用 role：code。
内容约束：代码 6-10 行、每行 ≤56 字符；文件名带扩展名 ≤20 字符；语言标签 ≤10 字符；说明一句 14-30 字。
代码行可用 `<span class="lt-kw">`（结构/键名）、`lt-str`（字符串值）、`lt-fn`（提醒色）、`lt-num`（数字）、`lt-cm`（注释）标语法色，不用则整块默认灰。

```html
<section class="slide" data-layout="code">
  <p class="lt-kicker">{{语境，如 色板速览}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="lt-term mt-l" style="margin-top:38px">
    <div class="lt-bar"><span class="lt-dot" style="background:#d20f39"></span><span class="lt-dot" style="background:#df8e1d"></span><span class="lt-dot" style="background:#40a02b"></span><span class="lt-file">{{文件名，如 latte-colors.jsonc}}</span><span class="lt-lang">{{语言标签，≤10 字符}}</span></div>
    <pre class="lt-codeblock"><code>{{代码 6-10 行、每行 ≤56 字符}}</code></pre>
  </div>
  <p class="mt-m" style="font-size:18px;line-height:1.7;color:#5C5F77">{{一句说明，14-30 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, lt-kicker, h2, mt-m, lt-term, mt-l, lt-bar, lt-dot, lt-file, lt-lang, lt-codeblock, lt-kw, lt-str, lt-fn, lt-cm, lt-num, deck-footer, slide-number, notes
