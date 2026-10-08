# Gruvbox 暗 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `gd-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-gruvbox-dark` 作用域生效，骨架里已写全，照抄结构即可。
> 每页的琥珀辉光、苔绿余温与扫描线由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
> code 版式是终端窗口母题，一份 deck 至多 1-2 页。
>
> **气质铁律（终端的克制）**：琥珀橙只允许出现在四处——提示符块/状态栏（gd-prompt / gd-status）、
> 关键数字（gd-stat-v）、强调药丸与按钮（gd-pill-accent / gd-btn）。沙黄管编号与时间点
> （gd-n / gd-tl-t），苔绿只管语言标签（gd-lang）。禁蓝紫、禁玻璃拟态、禁现代 IDE 质感；
> 辉光只做环境光，不给正文加彩。

---

## cover（开机封面）
指纹：hero

用途：开场页。题签 + 等宽大字标题 + 琥珀渐隐线 + 一句定位，提示符块与药丸收底。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 30-52 字；题签 ≤16 字；提示符块 ≤16 字；药丸 ≤12 字。

```html
<section class="slide full" data-layout="cover">
  <p class="gd-kicker">{{题签，≤16 字，如 ~/share · 终端效率之夜}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="gd-rule mt-s"></div>
  <p class="lede mt-m" style="max-width:52ch">{{一句话定位，30-52 字}}</p>
  <div class="row mt-l" style="gap:24px">
    <span class="gd-prompt">{{提示符短语，≤16 字，如 :wq}}</span>
    <span class="gd-pill">{{副题或人名，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, gd-kicker, h1, mt-m, gd-rule, mt-s, lede, mt-l, row, gd-prompt, gd-pill, deck-footer, slide-number, notes

---

## contents（目录终端）
指纹：table
数量：gd-item=4

用途：议程页。一块终端大卡里放 4 行篇目：方括号编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 12-24 字。

```html
<section class="slide" data-layout="contents">
  <p class="gd-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="gd-card mt-l" style="margin-top:44px">
    <div class="gd-item"><span class="gd-n">[01]</span><span class="gd-t">{{篇名，≤8 字}}</span><span class="gd-d">{{说明，12-24 字}}</span></div>
    <div class="gd-item"><span class="gd-n">[02]</span><span class="gd-t">{{篇名，≤8 字}}</span><span class="gd-d">{{说明，12-24 字}}</span></div>
    <div class="gd-item"><span class="gd-n">[03]</span><span class="gd-t">{{篇名，≤8 字}}</span><span class="gd-d">{{说明，12-24 字}}</span></div>
    <div class="gd-item"><span class="gd-n">[04]</span><span class="gd-t">{{篇名，≤8 字}}</span><span class="gd-d">{{说明，12-24 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gd-kicker, h2, mt-m, gd-card, mt-l, gd-item, gd-n, gd-t, gd-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：gd-card=3

用途：恰好三张终端卡。每张：琥珀药丸题签 + 等宽小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="gd-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="gd-card"><span class="gd-pill gd-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="gd-note mt-s">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="gd-card"><span class="gd-pill gd-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="gd-note mt-s">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="gd-card"><span class="gd-pill gd-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="gd-note mt-s">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gd-kicker, h2, mt-m, grid, g3, mt-l, gd-card, gd-pill, gd-pill-accent, h4, gd-note, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：gd-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边终端卡装三行步骤。
适用 role：content。
内容约束：左 lede 40-80 字 + 补充 16-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="gd-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，40-80 字}}</p>
      <p class="gd-note mt-m">{{补充：判断标准或代价，16-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="gd-pill">{{要点 1，≤8 字}}</span>
        <span class="gd-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="gd-card">
      <div class="gd-step"><span class="gd-n">[01]</span><p class="gd-mini-t">{{一步，12-26 字}}</p></div>
      <div class="gd-step"><span class="gd-n">[02]</span><p class="gd-mini-t">{{一步，12-26 字}}</p></div>
      <div class="gd-step"><span class="gd-n">[03]</span><p class="gd-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gd-kicker, h2, mt-m, grid, g2, mt-l, lede, gd-note, row, gd-pill, gd-card, gd-step, gd-n, gd-mini-t, deck-footer, slide-number, notes

---

## metrics（琥珀数字）
指纹：chart
数量：gd-stat=3

用途：三个关键数据。琥珀大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-32 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="gd-kicker">{{数据语境}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="gd-stat"><div class="gd-stat-v">{{数值 ≤6 字符}}<span class="gd-stat-u">{{单位}}</span></div><div class="gd-stat-l">{{指标名，≤8 字}}</div><p class="gd-stat-note">{{口径，14-32 字}}</p></div>
    <div class="gd-stat"><div class="gd-stat-v">{{数值 ≤6 字符}}<span class="gd-stat-u">{{单位}}</span></div><div class="gd-stat-l">{{指标名，≤8 字}}</div><p class="gd-stat-note">{{口径，14-32 字}}</p></div>
    <div class="gd-stat"><div class="gd-stat-v">{{数值 ≤6 字符}}<span class="gd-stat-u">{{单位}}</span></div><div class="gd-stat-l">{{指标名，≤8 字}}</div><p class="gd-stat-note">{{口径，14-32 字}}</p></div>
  </div>
  <p class="gd-src mt-m" style="margin-top:40px">来源：{{出处与统计口径，14-44 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gd-kicker, h2, mt-m, grid, g3, mt-l, gd-stat, gd-stat-v, gd-stat-u, gd-stat-l, gd-stat-note, gd-src, deck-footer, slide-number, notes

---

## quote（终端引文）
指纹：quote

用途：整页一句引文。等宽大字 + 出处 + 两个支撑药丸，像一行被反复引用的注释。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="gd-kicker">{{语境}}</p>
  <p class="gd-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="gd-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="gd-pill gd-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="gd-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gd-kicker, gd-quote, mt-l, gd-src, mt-m, row, gd-pill, gd-pill-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 等宽大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="gd-kicker">{{进度，如 卷二 · 键位即效率}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="gd-pill gd-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="gd-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, gd-kicker, h1, mt-m, lede, mt-l, row, gd-pill, gd-pill-accent, deck-footer, slide-number, notes

---

## moments（演进时间线）
指纹：chart
数量：gd-tl-item=4

用途：4 个节点的横向时间线：琥珀方框序号 + 时间点 + 一句事件，讲演进最稳。
适用 role：content / data。
内容约束：4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-44 字。

```html
<section class="slide" data-layout="moments">
  <p class="gd-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="gd-tl mt-l" style="margin-top:52px">
    <div class="gd-tl-item"><div class="gd-tl-dot">01</div><div class="gd-tl-t">{{时间点，≤10 字符}}</div><p class="gd-tl-d">{{事件，12-26 字}}</p></div>
    <div class="gd-tl-item"><div class="gd-tl-dot">02</div><div class="gd-tl-t">{{时间点，≤10 字符}}</div><p class="gd-tl-d">{{事件，12-26 字}}</p></div>
    <div class="gd-tl-item"><div class="gd-tl-dot">03</div><div class="gd-tl-t">{{时间点，≤10 字符}}</div><p class="gd-tl-d">{{事件，12-26 字}}</p></div>
    <div class="gd-tl-item"><div class="gd-tl-dot">04</div><div class="gd-tl-t">{{时间点，≤10 字符}}</div><p class="gd-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="gd-src mt-m" style="margin-top:44px">{{一句读法或口径，14-44 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gd-kicker, h2, mt-m, gd-tl, mt-l, gd-tl-item, gd-tl-dot, gd-tl-t, gd-tl-d, gd-src, deck-footer, slide-number, notes

---

## closing（收尾落款）
指纹：hero

用途：收尾页。等宽大字 + 一句行动提醒 + 琥珀描边按钮 + 药丸 + 提示符块，如一次干净退出。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-48 字；按钮 ≤8 字；药丸 ≤12 字；提示符块 ≤16 字。

```html
<section class="slide full" data-layout="closing">
  <p class="gd-kicker">{{提醒语境}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-48 字}}</p>
  <div class="row mt-l" style="gap:24px">
    <span class="gd-btn">{{按钮文案，≤8 字}}</span>
    <span class="gd-pill">{{次级信息，≤12 字}}</span>
    <span class="gd-prompt">{{提示符短语，≤16 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, gd-kicker, h1, mt-m, lede, mt-l, row, gd-btn, gd-pill, gd-prompt, deck-footer, slide-number, notes

---

## code（代码示例）
指纹：code

用途：终端窗口母题。标题栏三圆点 + 文件名 + 语言标签，等宽代码块直排，状态栏一句话收底。
适用 role：code。
内容约束：代码 6-10 行、每行 ≤56 字符（尖括号与 & 要 HTML 转义）；配文件名 + 语言标签 + 状态栏左右各一句 + 页尾一句说明 18-40 字；一份 deck 至多 1-2 页。

```html
<section class="slide" data-layout="code">
  <p class="gd-kicker">{{语境}}</p>
  <h2 class="h2 mt-m">{{标题，≤14 字}}</h2>
  <div class="gd-term mt-l" style="margin-top:44px">
    <div class="gd-term-bar"><span class="gd-dot"></span><span class="gd-dot"></span><span class="gd-dot"></span><span class="gd-term-name">{{文件名，如 ~/.vimrc}}</span><span class="gd-lang">{{语言标签，≤8 字符}}</span></div>
    <pre class="gd-codeblock"><code>{{代码直排 6-10 行，每行 ≤56 字符，尖括号与 & 转义}}</code></pre>
    <div class="gd-status"><span>{{左：模式或项目一句}}</span><span>{{右：状态一句话}}</span></div>
  </div>
  <p class="gd-note mt-m" style="margin-top:28px">{{一句说明，18-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gd-kicker, h2, mt-m, gd-term, mt-l, gd-term-bar, gd-dot, gd-term-name, gd-lang, gd-codeblock, gd-status, gd-note, deck-footer, slide-number, notes
