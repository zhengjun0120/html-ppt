# 酸蓝张力 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `ab-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-acid-blue-business` 作用域生效，骨架里已写全，照抄结构即可。
> 左侧电光蓝色轨与右下轮廓几何由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（单色即纪律）**：电光蓝 #435BE1 是唯一主强调色——只用于题签（ab-meta）、
> 横线（ab-hline）、标识块（ab-badge）、大数字（ab-stat-v）、时间线方块（ab-tl-dot）与
> 按钮（ab-btn）。荧光黄绿只作几何点缀（ab-meta 刻度、ab-stat 刻度、ab-circle 菱形），
> 永不承载文字；标题正文保持黑灰，对比度不过渡到蓝上。禁粉彩、禁渐变、禁霓虹赛博、
> 禁对称拥挤布局。

---

## cover（张力封面）
指纹：hero

用途：开场页。等宽题签 + 超重字重大标题 + 酸蓝横线 + 一句论点，左下立酸蓝标识块。
适用 role：cover。
内容约束：主标题 ≤12 字；lede 24-52 字；题签 ≤14 字；标识 2-4 字。

```html
<section class="slide full" data-layout="cover">
  <p class="ab-meta">{{题签，≤14 字，如 云迁移项目 · 方案评审}}</p>
  <h1 class="h1 mt-m">{{主标题，≤12 字，可 <br> 分行}}</h1>
  <div class="ab-hline mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句论点，24-52 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="ab-badge">{{标识 2-4 字}}</div>
    <span class="ab-pill">{{时间或场合，≤14 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年月}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, ab-meta, h1, mt-m, ab-hline, mt-s, lede, mt-l, row, ab-badge, ab-pill, deck-footer, slide-number, notes

---

## contents（议程页）
指纹：table
数量：ab-item=4

用途：议程页。一块白卡黑描边里放 4 行议程：等宽编号方块 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="ab-meta">{{引导语，如 AGENDA / 评审议程}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="ab-card mt-l" style="margin-top:44px">
    <div class="ab-item"><span class="ab-n">01</span><span class="ab-t">{{篇名，≤8 字}}</span><span class="ab-d">{{说明，14-26 字}}</span></div>
    <div class="ab-item"><span class="ab-n">02</span><span class="ab-t">{{篇名，≤8 字}}</span><span class="ab-d">{{说明，14-26 字}}</span></div>
    <div class="ab-item"><span class="ab-n">03</span><span class="ab-t">{{篇名，≤8 字}}</span><span class="ab-d">{{说明，14-26 字}}</span></div>
    <div class="ab-item"><span class="ab-n">04</span><span class="ab-t">{{篇名，≤8 字}}</span><span class="ab-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年月}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ab-meta, h2, mt-m, ab-card, mt-l, ab-item, ab-n, ab-t, ab-d, deck-footer, slide-number, notes

---

## keynotes（三卡策略）
指纹：cards
数量：ab-card=3

用途：恰好三张白卡黑描边。每张：酸蓝题签 + 重字重小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤6 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="ab-meta">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:36px;margin-top:44px">
    <div class="ab-card"><span class="ab-pill ab-pill-accent">{{题签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#404040">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="ab-card"><span class="ab-pill ab-pill-accent">{{题签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#404040">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="ab-card"><span class="ab-pill ab-pill-accent">{{题签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#404040">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年月}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ab-meta, h2, mt-m, grid, g3, mt-l, ab-card, ab-pill, ab-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：ab-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边白卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="ab-meta">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#404040">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="ab-pill">{{要点 1，≤8 字}}</span>
        <span class="ab-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="ab-card">
      <div class="ab-step"><span class="ab-n">一</span><p class="ab-mini-t">{{一步，12-26 字}}</p></div>
      <div class="ab-step"><span class="ab-n">二</span><p class="ab-mini-t">{{一步，12-26 字}}</p></div>
      <div class="ab-step"><span class="ab-n">三</span><p class="ab-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年月}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ab-meta, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, ab-pill, ab-card, ab-step, ab-n, ab-mini-t, deck-footer, slide-number, notes

---

## metrics（关键数字）
指纹：chart
数量：ab-stat=3

用途：三个关键数字。酸蓝大数字是唯一主视觉，黑顶线压荧光刻度，口径写进块内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="ab-meta">{{数据语境，如 BASELINE / 盘点口径}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:40px;margin-top:48px">
    <div class="ab-stat"><div class="ab-stat-v">{{数值 ≤6 字符}}<span class="ab-stat-u">{{单位}}</span></div><div class="ab-stat-l">{{指标名，≤8 字}}</div><p class="ab-stat-note">{{口径，14-30 字}}</p></div>
    <div class="ab-stat"><div class="ab-stat-v">{{数值 ≤6 字符}}<span class="ab-stat-u">{{单位}}</span></div><div class="ab-stat-l">{{指标名，≤8 字}}</div><p class="ab-stat-note">{{口径，14-30 字}}</p></div>
    <div class="ab-stat"><div class="ab-stat-v">{{数值 ≤6 字符}}<span class="ab-stat-u">{{单位}}</span></div><div class="ab-stat-l">{{指标名，≤8 字}}</div><p class="ab-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="margin-top:40px;font-size:17px;color:#8A8F99;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年月}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ab-meta, h2, mt-m, grid, g3, mt-l, ab-stat, ab-stat-v, ab-stat-u, ab-stat-l, ab-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（原则引言）
指纹：quote

用途：整页一句原则。超重字重引言 + 等宽出处 + 两个支撑药丸，右侧留白处衬酸蓝圆环。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤24 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="ab-meta">{{语境，如 评审原则}}</p>
  <p class="ab-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="ab-src mt-m">—— {{出处：人与场合，≤24 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="ab-pill ab-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="ab-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="ab-circle" style="position:absolute;right:150px;top:50%;transform:translateY(-50%)"></div>
  <div class="deck-footer"><span>{{署名 · 年月}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ab-meta, ab-quote, mt-l, ab-src, mt-m, row, ab-pill, ab-pill-accent, ab-circle, deck-footer, slide-number, notes

---

## divider（章节页）
指纹：hero

用途：章节过渡。进度题签 + 超大章节名 + 一句过渡 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤12 字；过渡句 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="ab-meta">{{进度，如 PART 02 / 迁移策略}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="ab-pill ab-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="ab-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年月}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, ab-meta, h1, mt-m, lede, mt-l, row, ab-pill, ab-pill-accent, deck-footer, slide-number, notes

---

## moments（路线时间线）
指纹：chart
数量：ab-tl-item=4

用途：4 个节点的横向时间线：酸蓝编号方块 + 等宽时间点 + 一句事件。
适用 role：content / data。
内容约束：恰好 4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="ab-meta">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="ab-tl mt-l" style="margin-top:52px">
    <div class="ab-tl-item"><div class="ab-tl-dot">01</div><div class="ab-tl-t">{{时间点，≤10 字符}}</div><p class="ab-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ab-tl-item"><div class="ab-tl-dot">02</div><div class="ab-tl-t">{{时间点，≤10 字符}}</div><p class="ab-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ab-tl-item"><div class="ab-tl-dot">03</div><div class="ab-tl-t">{{时间点，≤10 字符}}</div><p class="ab-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ab-tl-item"><div class="ab-tl-dot">04</div><div class="ab-tl-t">{{时间点，≤10 字符}}</div><p class="ab-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="margin-top:40px;font-size:17px;color:#8A8F99;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年月}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ab-meta, h2, mt-m, ab-tl, mt-l, ab-tl-item, ab-tl-dot, ab-tl-t, ab-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（行动请求）
指纹：hero

用途：收尾页。超大结论 + 一句行动提醒 + 酸蓝按钮 + 描边药丸。
适用 role：thanks / cta / content。
内容约束：标题 ≤12 字；lede 20-48 字；按钮 ≤6 字；药丸 ≤14 字。

```html
<section class="slide full" data-layout="closing">
  <p class="ab-meta">{{行动语境}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="ab-btn">{{按钮文案，≤6 字}}</span>
    <span class="ab-pill">{{次级信息，≤14 字}}</span>
    <div class="ab-badge">{{标识 2-4 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年月}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, ab-meta, h1, mt-m, lede, mt-l, row, ab-btn, ab-pill, ab-badge, deck-footer, slide-number, notes
