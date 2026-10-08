# 日式极简 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `jm-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-japanese-minimal` 作用域生效，骨架里已写全，照抄结构即可。
> 每页右侧的円相与右上纸角由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（一页一笔朱红）**：朱红 #DC2626 每页至多出现一处——jm-red-line / jm-stat-accent /
> jm-pill-accent / jm-vert-accent 四选一；其余全部交给墨色与留白。
> 宋体大字（h1/h2）每页最多一组；竖排 jm-vert 只放留白处，一页至多一条。
> 一页最多两组内容块 + 页脚，元素之间的间距要让空气流通。

---

## cover（题字封面）
指纹：hero

用途：开场页。宽字距题签 + 宋体大字标题 + 一笔朱红短线 + 一句定位，右侧留白可立竖排题款。
适用 role：cover。
内容约束：主标题 ≤14 字可两行；lede 30-60 字；题签 ≤14 字；竖排短语 ≤7 字。

```html
<section class="slide full" data-layout="cover">
  <p class="jm-kicker">{{题签，≤14 字，如 侘寂工作坊 · 秋季场}}</p>
  <h1 class="h1 mt-m">{{主标题，≤14 字，可 <br> 分两行}}</h1>
  <div class="jm-red-line mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，30-60 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="jm-pill">{{信息 1，≤12 字}}</span>
    <span class="jm-pill">{{信息 2，≤10 字}}</span>
  </div>
  <div class="jm-vert" style="position:absolute;right:120px;top:50%;transform:translateY(-50%)">{{竖排短语 ≤7 字}}</div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, jm-kicker, h1, mt-m, jm-red-line, mt-s, lede, mt-l, row, jm-pill, jm-vert, deck-footer, slide-number, notes

---

## contents（目次）
指纹：table
数量：jm-item=4

用途：议程页。一块白纸卡里放 4 行条目：汉数字号 + 条目名 + 一句说明，发丝线分行。
适用 role：toc。
内容约束：恰好 4 行；条目名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="jm-kicker">{{引导语，如 本日次第}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="jm-card mt-l" style="margin-top:44px">
    <div class="jm-item"><span class="jm-n">一</span><span class="jm-t">{{条目名，≤8 字}}</span><span class="jm-d">{{说明，14-26 字}}</span></div>
    <div class="jm-item"><span class="jm-n">二</span><span class="jm-t">{{条目名，≤8 字}}</span><span class="jm-d">{{说明，14-26 字}}</span></div>
    <div class="jm-item"><span class="jm-n">三</span><span class="jm-t">{{条目名，≤8 字}}</span><span class="jm-d">{{说明，14-26 字}}</span></div>
    <div class="jm-item"><span class="jm-n">四</span><span class="jm-t">{{条目名，≤8 字}}</span><span class="jm-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, jm-kicker, h2, mt-m, jm-card, mt-l, jm-item, jm-n, jm-t, jm-d, deck-footer, slide-number, notes

---

## keynotes（三帖要点）
指纹：cards
数量：jm-card=3

用途：恰好三张白纸卡。每张：题签 + 小标题 + 两句说明，层级靠字距与留白。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="jm-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:40px;margin-top:48px">
    <div class="jm-card"><span class="jm-pill">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:18px;line-height:1.85;color:#525252">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="jm-card"><span class="jm-pill">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:18px;line-height:1.85;color:#525252">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="jm-card"><span class="jm-pill">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:18px;line-height:1.85;color:#525252">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, jm-kicker, h2, mt-m, grid, g3, mt-l, jm-card, jm-pill, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：jm-step=3

用途：左边把一件事讲透（lede + 补充 + 签），右边白纸卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个签；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="jm-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:18px;line-height:1.85;color:#525252">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="jm-pill">{{要点 1，≤8 字}}</span>
        <span class="jm-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="jm-card">
      <div class="jm-step"><span class="jm-n">一</span><p class="jm-mini-t">{{一步，12-26 字}}</p></div>
      <div class="jm-step"><span class="jm-n">二</span><p class="jm-mini-t">{{一步，12-26 字}}</p></div>
      <div class="jm-step"><span class="jm-n">三</span><p class="jm-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, jm-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, jm-pill, jm-card, jm-step, jm-n, jm-mini-t, deck-footer, slide-number, notes

---

## metrics（三行数字）
指纹：chart
数量：jm-stat=3

用途：三个关键数据。宋体大数字是唯一主视觉，中栏一道朱红封线点睛；口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="jm-kicker">{{数据语境，如 春季场 · 复盘}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="jm-stat"><div class="jm-stat-v">{{数值 ≤6 字符}}<span class="jm-stat-u">{{单位}}</span></div><div class="jm-stat-l">{{指标名，≤8 字}}</div><p class="jm-stat-note">{{口径，14-30 字}}</p></div>
    <div class="jm-stat jm-stat-accent"><div class="jm-stat-v">{{数值 ≤6 字符}}<span class="jm-stat-u">{{单位}}</span></div><div class="jm-stat-l">{{指标名，≤8 字}}</div><p class="jm-stat-note">{{口径，14-30 字}}</p></div>
    <div class="jm-stat"><div class="jm-stat-v">{{数值 ≤6 字符}}<span class="jm-stat-u">{{单位}}</span></div><div class="jm-stat-l">{{指标名，≤8 字}}</div><p class="jm-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#A3A3A3;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, jm-kicker, h2, mt-m, grid, g3, mt-l, jm-stat, jm-stat-accent, jm-stat-v, jm-stat-u, jm-stat-l, jm-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（一笔引文）
指纹：quote

用途：整页一句引文。宋体大字 + 出处 + 两个细描边签，右侧留白处立一条竖排短语（可点朱）。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；签各 ≤8 字；竖排短语 ≤7 字。

```html
<section class="slide" data-layout="quote">
  <p class="jm-kicker">{{语境，如 茶人题记}}</p>
  <p class="jm-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="jm-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="jm-pill jm-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="jm-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="jm-vert" style="position:absolute;right:120px;top:50%;transform:translateY(-50%)">{{竖排短语 ≤7 字}}<span class="jm-vert-accent">·</span></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, jm-kicker, jm-quote, mt-l, jm-src, mt-m, row, jm-pill, jm-pill-accent, jm-vert, jm-vert-accent, deck-footer, slide-number, notes

---

## divider（章节间）
指纹：hero

用途：章节过渡。进度题签 + 宋体大字章节名 + 一笔朱红 + 一个过渡问题 + 两个看点签。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-52 字；签各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="jm-kicker">{{进度，如 第二章 · 观物}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <div class="jm-red-line mt-s"></div>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="jm-pill">{{看点 1，≤8 字}}</span>
    <span class="jm-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, jm-kicker, h1, mt-m, jm-red-line, mt-s, lede, mt-l, row, jm-pill, deck-footer, slide-number, notes

---

## moments（流程时间线）
指纹：chart
数量：jm-tl-item=4

用途：3-4 个节点的横向时间线：空心圆点 + 时间点 + 一句事件，发丝线串联。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="jm-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="jm-tl mt-l" style="margin-top:56px">
    <div class="jm-tl-item"><div class="jm-tl-dot"></div><div class="jm-tl-t">{{时间点，≤10 字符}}</div><p class="jm-tl-d">{{事件，12-26 字}}</p></div>
    <div class="jm-tl-item"><div class="jm-tl-dot"></div><div class="jm-tl-t">{{时间点，≤10 字符}}</div><p class="jm-tl-d">{{事件，12-26 字}}</p></div>
    <div class="jm-tl-item"><div class="jm-tl-dot"></div><div class="jm-tl-t">{{时间点，≤10 字符}}</div><p class="jm-tl-d">{{事件，12-26 字}}</p></div>
    <div class="jm-tl-item"><div class="jm-tl-dot"></div><div class="jm-tl-t">{{时间点，≤10 字符}}</div><p class="jm-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#A3A3A3;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, jm-kicker, h2, mt-m, jm-tl, mt-l, jm-tl-item, jm-tl-dot, jm-tl-t, jm-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（收尾款识）
指纹：hero

用途：收尾页。宋体大字 + 一句行动提醒 + 墨线描边按钮 + 细描边签，如卷末落款。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤6 字；签 ≤12 字。

```html
<section class="slide full" data-layout="closing">
  <p class="jm-kicker">{{行动语境，如 冬季场报名}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句行动提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:22px">
    <span class="jm-btn">{{按钮文案，≤6 字}}</span>
    <span class="jm-pill">{{次级信息，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, jm-kicker, h1, mt-m, lede, mt-l, row, jm-btn, jm-pill, deck-footer, slide-number, notes
