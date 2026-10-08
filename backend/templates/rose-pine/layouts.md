# 玫瑰松 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `rp-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-rose-pine` 作用域生效，骨架里已写全，照抄结构即可。
> 每页的玫瑰辉光、松青余温与右下松枝由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
> code 版式是柔和终端母题，一份 deck 至多 1-2 页。
>
> **气质铁律（温柔的暗色）**：玫瑰粉是唯一的心跳色——题签细线、关键数字（rp-stat-v）、
> 时间点（rp-tl-t）、强调药丸与按钮（rp-pill-accent / rp-btn）、芯片圆点。金色只管编号
> （rp-n / rp-tl-dot），松青只管语言标签（rp-lang）与松枝衬底。卡必须圆角轻描边；
> 禁霓虹、禁硬直角、禁把暗色压成纯黑。

---

## cover（暮色封面）
指纹：hero

用途：开场页。题签 + 大字标题 + 玫瑰渐隐线 + 一句定位，玫瑰圆点芯片与药丸收底。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 30-52 字；题签 ≤16 字；芯片 ≤16 字；药丸 ≤12 字。

```html
<section class="slide full" data-layout="cover">
  <p class="rp-kicker">{{题签，≤16 字，如 设计系统分享 · 暮松 v2.0}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="rp-rule mt-s"></div>
  <p class="lede mt-m" style="max-width:52ch">{{一句话定位，30-52 字}}</p>
  <div class="row mt-l" style="gap:24px">
    <span class="rp-chip">{{芯片短语，≤16 字，如 npm i duskpine}}</span>
    <span class="rp-pill">{{副题或人名，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, rp-kicker, h1, mt-m, rp-rule, mt-s, lede, mt-l, row, rp-chip, rp-pill, deck-footer, slide-number, notes

---

## contents（林间目录）
指纹：table
数量：rp-item=4

用途：议程页。一块暮色大卡里放 4 行篇目：金圈编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 12-24 字。

```html
<section class="slide" data-layout="contents">
  <p class="rp-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="rp-card mt-l" style="margin-top:44px">
    <div class="rp-item"><span class="rp-n">01</span><span class="rp-t">{{篇名，≤8 字}}</span><span class="rp-d">{{说明，12-24 字}}</span></div>
    <div class="rp-item"><span class="rp-n">02</span><span class="rp-t">{{篇名，≤8 字}}</span><span class="rp-d">{{说明，12-24 字}}</span></div>
    <div class="rp-item"><span class="rp-n">03</span><span class="rp-t">{{篇名，≤8 字}}</span><span class="rp-d">{{说明，12-24 字}}</span></div>
    <div class="rp-item"><span class="rp-n">04</span><span class="rp-t">{{篇名，≤8 字}}</span><span class="rp-d">{{说明，12-24 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rp-kicker, h2, mt-m, rp-card, mt-l, rp-item, rp-n, rp-t, rp-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：rp-card=3

用途：恰好三张暮色卡。每张：玫瑰药丸题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="rp-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="rp-card"><span class="rp-pill rp-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="rp-note mt-s">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="rp-card"><span class="rp-pill rp-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="rp-note mt-s">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="rp-card"><span class="rp-pill rp-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="rp-note mt-s">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rp-kicker, h2, mt-m, grid, g3, mt-l, rp-card, rp-pill, rp-pill-accent, h4, rp-note, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：rp-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边暮色卡装三行步骤。
适用 role：content。
内容约束：左 lede 40-80 字 + 补充 16-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="rp-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，40-80 字}}</p>
      <p class="rp-note mt-m">{{补充：判断标准或代价，16-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="rp-pill">{{要点 1，≤8 字}}</span>
        <span class="rp-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="rp-card">
      <div class="rp-step"><span class="rp-n">01</span><p class="rp-mini-t">{{一步，12-26 字}}</p></div>
      <div class="rp-step"><span class="rp-n">02</span><p class="rp-mini-t">{{一步，12-26 字}}</p></div>
      <div class="rp-step"><span class="rp-n">03</span><p class="rp-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rp-kicker, h2, mt-m, grid, g2, mt-l, lede, rp-note, row, rp-pill, rp-card, rp-step, rp-n, rp-mini-t, deck-footer, slide-number, notes

---

## metrics（玫瑰数字）
指纹：chart
数量：rp-stat=3

用途：三个关键数据。玫瑰大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-32 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="rp-kicker">{{数据语境}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="rp-stat"><div class="rp-stat-v">{{数值 ≤6 字符}}<span class="rp-stat-u">{{单位}}</span></div><div class="rp-stat-l">{{指标名，≤8 字}}</div><p class="rp-stat-note">{{口径，14-32 字}}</p></div>
    <div class="rp-stat"><div class="rp-stat-v">{{数值 ≤6 字符}}<span class="rp-stat-u">{{单位}}</span></div><div class="rp-stat-l">{{指标名，≤8 字}}</div><p class="rp-stat-note">{{口径，14-32 字}}</p></div>
    <div class="rp-stat"><div class="rp-stat-v">{{数值 ≤6 字符}}<span class="rp-stat-u">{{单位}}</span></div><div class="rp-stat-l">{{指标名，≤8 字}}</div><p class="rp-stat-note">{{口径，14-32 字}}</p></div>
  </div>
  <p class="rp-src mt-m" style="margin-top:40px">来源：{{出处与统计口径，14-44 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rp-kicker, h2, mt-m, grid, g3, mt-l, rp-stat, rp-stat-v, rp-stat-u, rp-stat-l, rp-stat-note, rp-src, deck-footer, slide-number, notes

---

## quote（衬线引文）
指纹：quote

用途：整页一句引文。衬线大字（全模板唯一的衬线时刻）+ 出处 + 两个支撑药丸。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="rp-kicker">{{语境}}</p>
  <p class="rp-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="rp-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="rp-pill rp-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="rp-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rp-kicker, rp-quote, mt-l, rp-src, mt-m, row, rp-pill, rp-pill-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="rp-kicker">{{进度，如 第二章 · 语义先行}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="rp-pill rp-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="rp-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, rp-kicker, h1, mt-m, lede, mt-l, row, rp-pill, rp-pill-accent, deck-footer, slide-number, notes

---

## moments（年轮时间线）
指纹：chart
数量：rp-tl-item=4

用途：4 个节点的横向时间线：金圈序号 + 时间点 + 一句事件，讲项目生长最稳。
适用 role：content / data。
内容约束：4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-44 字。

```html
<section class="slide" data-layout="moments">
  <p class="rp-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="rp-tl mt-l" style="margin-top:52px">
    <div class="rp-tl-item"><div class="rp-tl-dot">01</div><div class="rp-tl-t">{{时间点，≤10 字符}}</div><p class="rp-tl-d">{{事件，12-26 字}}</p></div>
    <div class="rp-tl-item"><div class="rp-tl-dot">02</div><div class="rp-tl-t">{{时间点，≤10 字符}}</div><p class="rp-tl-d">{{事件，12-26 字}}</p></div>
    <div class="rp-tl-item"><div class="rp-tl-dot">03</div><div class="rp-tl-t">{{时间点，≤10 字符}}</div><p class="rp-tl-d">{{事件，12-26 字}}</p></div>
    <div class="rp-tl-item"><div class="rp-tl-dot">04</div><div class="rp-tl-t">{{时间点，≤10 字符}}</div><p class="rp-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="rp-src mt-m" style="margin-top:44px">{{一句读法或口径，14-44 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rp-kicker, h2, mt-m, rp-tl, mt-l, rp-tl-item, rp-tl-dot, rp-tl-t, rp-tl-d, rp-src, deck-footer, slide-number, notes

---

## closing（收尾落款）
指纹：hero

用途：收尾页。大字 + 一句行动提醒 + 玫瑰描边按钮 + 药丸 + 玫瑰圆点芯片，如暮色里的一次道别。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-48 字；按钮 ≤8 字；药丸 ≤12 字；芯片 ≤16 字。

```html
<section class="slide full" data-layout="closing">
  <p class="rp-kicker">{{提醒语境}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-48 字}}</p>
  <div class="row mt-l" style="gap:24px">
    <span class="rp-btn">{{按钮文案，≤8 字}}</span>
    <span class="rp-pill">{{次级信息，≤12 字}}</span>
    <span class="rp-chip">{{芯片短语，≤16 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, rp-kicker, h1, mt-m, lede, mt-l, row, rp-btn, rp-pill, rp-chip, deck-footer, slide-number, notes

---

## code（代码示例）
指纹：code

用途：柔和终端母题。标题栏三圆点 + 文件名 + 语言标签，等宽代码块直排，状态栏一句话收底。
适用 role：code。
内容约束：代码 6-10 行、每行 ≤56 字符（尖括号与 & 要 HTML 转义）；配文件名 + 语言标签 + 状态栏左右各一句 + 页尾一句说明 18-40 字；一份 deck 至多 1-2 页。

```html
<section class="slide" data-layout="code">
  <p class="rp-kicker">{{语境}}</p>
  <h2 class="h2 mt-m">{{标题，≤14 字}}</h2>
  <div class="rp-term mt-l" style="margin-top:44px">
    <div class="rp-term-bar"><span class="rp-dot"></span><span class="rp-dot"></span><span class="rp-dot"></span><span class="rp-term-name">{{文件名，如 themes/duskpine.json}}</span><span class="rp-lang">{{语言标签，≤8 字符}}</span></div>
    <pre class="rp-codeblock"><code>{{代码直排 6-10 行，每行 ≤56 字符，尖括号与 & 转义}}</code></pre>
    <div class="rp-status"><span>{{左：项目或版本一句}}</span><span>{{右：状态一句话}}</span></div>
  </div>
  <p class="rp-note mt-m" style="margin-top:28px">{{一句说明，18-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rp-kicker, h2, mt-m, rp-term, mt-l, rp-term-bar, rp-dot, rp-term-name, rp-lang, rp-codeblock, rp-status, rp-note, deck-footer, slide-number, notes
