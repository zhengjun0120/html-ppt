# 北极冷 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `ac-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-arctic-cool` 作用域生效，骨架里已写全，照抄结构即可。
> 每页的极光余晖与冰面等高线由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（冰川青是身份）**：冰川青只允许出现在六处——题签细线（ac-kicker 自带）、
> 关键词（ac-hl）、编号与徽牌（ac-n）、指标名（ac-stat-l）、时间点（ac-tl-t）、
> 行动钮描边（ac-btn）。数字与正文一律冰白/雾蓝灰；无发光、无玻璃拟态、无暖色。
> 大字标题（h1/h2）每页最多一组；等宽字体只管数字与标签，不上中文长句。

---

## cover（冰面封面）
指纹：hero

用途：开场页。等宽题签 + 克制大字（关键词冰川青）+ 细光线 + 一句定位，配编号徽牌。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；徽标 2-4 字符。

```html
<section class="slide full" data-layout="cover">
  <p class="ac-kicker">{{题签，≤14 字，中英混排}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行，含 <span class="ac-hl">关键词</span>}}</h1>
  <div class="ac-line mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="ac-n" style="width:64px;height:64px">{{徽标 2-4 字符}}</div>
    <span class="ac-pill">{{副题或时间地点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, ac-kicker, h1, mt-m, ac-hl, ac-line, mt-s, lede, mt-l, row, ac-n, ac-pill, deck-footer, slide-number, notes

---

## contents（冰架目录）
指纹：table
数量：ac-item=4

用途：议程页。一块冰面大卡里放 4 行篇目：等宽编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="ac-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="ac-card mt-l" style="margin-top:44px">
    <div class="ac-item"><span class="ac-n">01</span><span class="ac-t">{{篇名，≤8 字}}</span><span class="ac-d">{{说明，14-26 字}}</span></div>
    <div class="ac-item"><span class="ac-n">02</span><span class="ac-t">{{篇名，≤8 字}}</span><span class="ac-d">{{说明，14-26 字}}</span></div>
    <div class="ac-item"><span class="ac-n">03</span><span class="ac-t">{{篇名，≤8 字}}</span><span class="ac-d">{{说明，14-26 字}}</span></div>
    <div class="ac-item"><span class="ac-n">04</span><span class="ac-t">{{篇名，≤8 字}}</span><span class="ac-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ac-kicker, h2, mt-m, ac-card, mt-l, ac-item, ac-n, ac-t, ac-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：ac-card=3

用途：恰好三张冰面卡。每张：冰川青题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="ac-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="ac-card"><span class="ac-pill ac-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#9db8cc">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="ac-card"><span class="ac-pill ac-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#9db8cc">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="ac-card"><span class="ac-pill ac-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#9db8cc">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ac-kicker, h2, mt-m, grid, g3, mt-l, ac-card, ac-pill, ac-pill-accent, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：ac-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边冰面卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="ac-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#9db8cc">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="ac-pill">{{要点 1，≤8 字}}</span>
        <span class="ac-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="ac-card">
      <div class="ac-step"><span class="ac-n">1</span><p class="ac-mini-t">{{一步，12-26 字}}</p></div>
      <div class="ac-step"><span class="ac-n">2</span><p class="ac-mini-t">{{一步，12-26 字}}</p></div>
      <div class="ac-step"><span class="ac-n">3</span><p class="ac-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ac-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, ac-pill, ac-card, ac-step, ac-n, ac-mini-t, deck-footer, slide-number, notes

---

## metrics（等宽数字）
指纹：chart
数量：ac-stat=3

用途：三个关键数据。等宽大数字是主视觉，冰川青指标名做身份，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="ac-kicker">{{数据语境}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="ac-stat"><div class="ac-stat-v">{{数值 ≤6 字符}}<span class="ac-stat-u">{{单位}}</span></div><div class="ac-stat-l">{{指标名，≤8 字}}</div><p class="ac-stat-note">{{口径，14-30 字}}</p></div>
    <div class="ac-stat"><div class="ac-stat-v">{{数值 ≤6 字符}}<span class="ac-stat-u">{{单位}}</span></div><div class="ac-stat-l">{{指标名，≤8 字}}</div><p class="ac-stat-note">{{口径，14-30 字}}</p></div>
    <div class="ac-stat"><div class="ac-stat-v">{{数值 ≤6 字符}}<span class="ac-stat-u">{{单位}}</span></div><div class="ac-stat-l">{{指标名，≤8 字}}</div><p class="ac-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#8fa9bf;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ac-kicker, h2, mt-m, grid, g3, mt-l, ac-stat, ac-stat-v, ac-stat-u, ac-stat-l, ac-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（领队引文）
指纹：quote

用途：整页一句引文。克制大字 + 等宽出处 + 两个支撑药丸，右侧留白处立一枚坐标标签。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="ac-kicker">{{语境}}</p>
  <p class="ac-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="ac-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="ac-pill ac-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="ac-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="ac-tag" style="position:absolute;right:120px;top:50%;transform:translateY(-50%)">{{坐标或编号 ≤12 字符}}</div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ac-kicker, ac-quote, mt-l, ac-src, mt-m, row, ac-pill, ac-pill-accent, ac-tag, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="ac-kicker">{{进度，如 第一幕 · 航程}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字，可 <br> 分行}}</h1>
  <div class="ac-line mt-s"></div>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="ac-pill ac-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="ac-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, ac-kicker, h1, mt-m, ac-line, mt-s, lede, mt-l, row, ac-pill, ac-pill-accent, deck-footer, slide-number, notes

---

## moments（航程时间线）
指纹：chart
数量：ac-tl-item=4

用途：3-4 个节点的横向时间线：冰川青圆点 + 日期 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="ac-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="ac-tl mt-l" style="margin-top:52px">
    <div class="ac-tl-item"><div class="ac-tl-dot">1</div><div class="ac-tl-t">{{时间点，≤10 字符}}</div><p class="ac-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ac-tl-item"><div class="ac-tl-dot">2</div><div class="ac-tl-t">{{时间点，≤10 字符}}</div><p class="ac-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ac-tl-item"><div class="ac-tl-dot">3</div><div class="ac-tl-t">{{时间点，≤10 字符}}</div><p class="ac-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ac-tl-item"><div class="ac-tl-dot">4</div><div class="ac-tl-t">{{时间点，≤10 字符}}</div><p class="ac-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#8fa9bf;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ac-kicker, h2, mt-m, ac-tl, mt-l, ac-tl-item, ac-tl-dot, ac-tl-t, ac-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（收尾归档）
指纹：hero

用途：收尾页。大字 + 一句行动提醒 + 冰川青描边按钮 + 药丸 + 编号徽牌。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="ac-kicker">{{提醒语境}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <div class="ac-line mt-s"></div>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="ac-btn">{{按钮文案，≤8 字}}</span>
    <span class="ac-pill">{{次级信息，≤10 字}}</span>
    <div class="ac-n" style="width:64px;height:64px">{{徽标 2-4 字符}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, ac-kicker, h1, mt-m, ac-line, mt-s, lede, mt-l, row, ac-btn, ac-pill, ac-n, deck-footer, slide-number, notes
