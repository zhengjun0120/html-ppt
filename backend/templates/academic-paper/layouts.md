# 学术论文 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `ap-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-academic-paper` 作用域生效，骨架里已写全，照抄结构即可。
> 页顶蓝条与纸面渐变由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（蓝色必须可溯源）**：学术蓝 #2563EB 只允许出现在五处——节号与编号
> （ap-kicker / ap-no / ap-tl-dot / ap-tl-t）、关键数字（ap-stat-v）、摘要标签
> （ap-abs-label）、强调药丸（ap-pill-accent）、收尾按钮（ap-btn）。蓝色永远指向
> 「有出处的东西」；标题一律墨黑衬线，禁给标题上蓝。卡只用白纸面 + 极淡描边，方角，无阴影。

---

## cover（论文封面）
指纹：hero

用途：开场页。期刊头（刊名 + 卷期页码）+ 居中标题与作者行 + 摘要卡，如一篇论文的首页。
适用 role：cover。
内容约束：标题 ≤22 字；摘要 44-80 字；作者行 ≤30 字；卷期行 ≤30 字符。

```html
<section class="slide full" data-layout="cover">
  <div class="ap-journal"><span class="ap-jname">{{刊名，如 Survey Research Methods}}</span><span class="ap-jmeta">{{卷期 · 年份 · 页码}}</span></div>
  <div class="tc mt-l">
    <h1 class="h1" style="font-size:64px">{{论文标题，≤22 字}}</h1>
    <p class="ap-authors mt-m">{{作者与机构，≤30 字}}</p>
  </div>
  <div class="ap-card mt-l" style="max-width:1000px;margin:44px auto 0">
    <div class="ap-abs-label">ABSTRACT · 摘要</div>
    <p class="ap-abs-text mt-s">{{摘要：研究什么、怎么综述、得到什么，44-80 字}}</p>
  </div>
  <div class="deck-footer"><span>{{通讯作者 · 单位}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, ap-journal, ap-jname, ap-jmeta, tc, mt-l, h1, ap-authors, mt-m, ap-card, ap-abs-label, ap-abs-text, mt-s, deck-footer, slide-number, notes

---

## contents（章节目录）
指纹：table
数量：ap-item=4

用途：目录页。一块纸面卡里放 4 行章节：蓝色节号 + 章节名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；章节名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="ap-kicker">{{引导语，如 OUTLINE · 本文结构}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="ap-card mt-l" style="margin-top:44px">
    <div class="ap-item"><span class="ap-no">1.</span><span class="ap-t">{{章节名，≤8 字}}</span><span class="ap-d">{{说明，14-26 字}}</span></div>
    <div class="ap-item"><span class="ap-no">2.</span><span class="ap-t">{{章节名，≤8 字}}</span><span class="ap-d">{{说明，14-26 字}}</span></div>
    <div class="ap-item"><span class="ap-no">3.</span><span class="ap-t">{{章节名，≤8 字}}</span><span class="ap-d">{{说明，14-26 字}}</span></div>
    <div class="ap-item"><span class="ap-no">4.</span><span class="ap-t">{{章节名，≤8 字}}</span><span class="ap-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{通讯作者 · 单位}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ap-kicker, h2, mt-m, ap-card, mt-l, ap-item, ap-no, ap-t, ap-d, deck-footer, slide-number, notes

---

## keynotes（三点综述）
指纹：cards
数量：ap-card=3, ap-no=3

用途：恰好三张纸面卡。每张：蓝色节号 + 小标题 + 两句说明，综述的分节陈述。
适用 role：content。
内容约束：恰好 3 卡；节号 ≤4 字符；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="ap-kicker">{{引导语，如 SECTION 2 · 方法谱系}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="ap-card"><span class="ap-no">2.1</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#404040">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="ap-card"><span class="ap-no">2.2</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#404040">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="ap-card"><span class="ap-no">2.3</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#404040">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{通讯作者 · 单位}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ap-kicker, h2, mt-m, grid, g3, mt-l, ap-card, ap-no, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左论右法）
指纹：split
数量：ap-step=3

用途：左边把一个论点讲透（lede + 补充 + 药丸），右边纸面卡装三行步骤或判据。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="ap-kicker">{{引导语，如 SECTION 4 · 选型路径}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{论点：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#404040">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="ap-pill">{{要点 1，≤8 字}}</span>
        <span class="ap-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="ap-card">
      <div class="ap-step"><span class="ap-no">01</span><p class="ap-mini-t">{{一步，12-26 字}}</p></div>
      <div class="ap-step"><span class="ap-no">02</span><p class="ap-mini-t">{{一步，12-26 字}}</p></div>
      <div class="ap-step"><span class="ap-no">03</span><p class="ap-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{通讯作者 · 单位}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ap-kicker, h2, mt-m, grid, g2, mt-l, lede, row, ap-pill, ap-card, ap-step, ap-no, ap-mini-t, deck-footer, slide-number, notes

---

## metrics（数据图表）
指纹：chart
数量：ap-stat=3

用途：三个关键数据。蓝色衬线大数字是唯一主视觉，口径写进卡内，来源以图表注（注：…）写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤10 字；口径 12-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="ap-kicker">{{数据语境，如 SECTION 3 · 效度比较}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="ap-stat"><div class="ap-stat-v">{{数值 ≤6 字符}}<span class="ap-stat-u">{{单位}}</span></div><div class="ap-stat-l">{{指标名，≤10 字}}</div><p class="ap-stat-note">{{口径，12-30 字}}</p></div>
    <div class="ap-stat"><div class="ap-stat-v">{{数值 ≤6 字符}}<span class="ap-stat-u">{{单位}}</span></div><div class="ap-stat-l">{{指标名，≤10 字}}</div><p class="ap-stat-note">{{口径，12-30 字}}</p></div>
    <div class="ap-stat"><div class="ap-stat-v">{{数值 ≤6 字符}}<span class="ap-stat-u">{{单位}}</span></div><div class="ap-stat-l">{{指标名，≤10 字}}</div><p class="ap-stat-note">{{口径，12-30 字}}</p></div>
  </div>
  <p class="ap-figcap mt-m" style="margin-top:44px">注：来源——{{出处与统计口径，14-40 字}}。</p>
  <div class="deck-footer"><span>{{通讯作者 · 单位}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ap-kicker, h2, mt-m, grid, g3, mt-l, ap-stat, ap-stat-v, ap-stat-u, ap-stat-l, ap-stat-note, ap-figcap, deck-footer, slide-number, notes

---

## quote（题记引文）
指纹：quote

用途：整页一句引文。衬线斜体大字 + 出处 + 两个支撑药丸，如论文的题记页。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="ap-kicker">{{语境，如 EPIGRAPH · 题记}}</p>
  <p class="ap-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="ap-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="ap-pill ap-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="ap-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{通讯作者 · 单位}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ap-kicker, ap-quote, mt-l, ap-src, mt-m, row, ap-pill, ap-pill-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。英文节号题签 + 墨黑大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤12 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="ap-kicker">{{进度，如 CHAPTER 3 · 效度之争}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="ap-pill ap-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="ap-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{通讯作者 · 单位}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, ap-kicker, h1, mt-m, lede, row, mt-l, ap-pill, ap-pill-accent, deck-footer, slide-number, notes

---

## moments（研究年表）
指纹：chart
数量：ap-tl-item=4

用途：3-4 个节点的横向时间线：蓝色方框编号 + 时间点 + 一句事件，如文献或研究的演进年表。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="ap-kicker">{{引导语，如 TIMELINE · 方法史}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="ap-tl mt-l" style="margin-top:52px">
    <div class="ap-tl-item"><div class="ap-tl-dot">01</div><div class="ap-tl-t">{{时间点，≤10 字符}}</div><p class="ap-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ap-tl-item"><div class="ap-tl-dot">02</div><div class="ap-tl-t">{{时间点，≤10 字符}}</div><p class="ap-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ap-tl-item"><div class="ap-tl-dot">03</div><div class="ap-tl-t">{{时间点，≤10 字符}}</div><p class="ap-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ap-tl-item"><div class="ap-tl-dot">04</div><div class="ap-tl-t">{{时间点，≤10 字符}}</div><p class="ap-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="ap-figcap mt-m" style="margin-top:44px">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{通讯作者 · 单位}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ap-kicker, h2, mt-m, ap-tl, mt-l, ap-tl-item, ap-tl-dot, ap-tl-t, ap-tl-d, ap-figcap, deck-footer, slide-number, notes

---

## closing（致谢投稿）
指纹：hero

用途：收尾页。墨黑大字 + 一句行动说明 + 蓝色按钮 + 药丸，投稿或预审的落款。
适用 role：thanks / cta / content。
内容约束：标题 ≤12 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤14 字。

```html
<section class="slide full" data-layout="closing">
  <p class="ap-kicker">{{提醒语境，如 PRE-SUBMISSION · 预审}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束说明，20-44 字}}</p>
  <div class="row mt-l" style="gap:22px">
    <span class="ap-btn">{{按钮文案，≤8 字}}</span>
    <span class="ap-pill">{{次级信息，≤14 字}}</span>
  </div>
  <div class="deck-footer"><span>{{通讯作者 · 单位}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, ap-kicker, h1, mt-m, lede, row, mt-l, ap-btn, ap-pill, deck-footer, slide-number, notes
