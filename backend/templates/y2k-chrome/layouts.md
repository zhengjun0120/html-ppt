# Y2K 铬 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `yk-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-y2k-chrome` 作用域生效，骨架里已写全，照抄结构即可。
> 每页的铬面高光斜扫与气泡光斑由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（银铬是身份）**：彩虹只允许出现在三处——标题彩虹条（yk-rainbow）、数字块上描边
> （yk-stat 自带）、题签短线（yk-kicker 自带）；文本级紫色只进 yk-pill-accent / yk-tl-t。
> 禁直角与小圆角：卡、徽标、按钮、药丸全部大圆角或正圆（已内建，别改）。
> 大字标题（h1/h2）每页最多一组；正文保持墨黑/深灰，不写浅色字。

---

## cover（铬面封面）
指纹：hero

用途：开场页。题签 + 大字标题 + 彩虹条 + 一句定位，右侧留白处有自动气泡衬底。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；徽标 2-4 字符。

```html
<section class="slide full" data-layout="cover">
  <p class="yk-kicker">{{题签，≤14 字，中英混排}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="yk-rainbow mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="yk-orb">{{徽标 2-4 字符}}</div>
    <span class="yk-pill">{{副题或时间地点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, yk-kicker, h1, mt-m, yk-rainbow, mt-s, lede, mt-l, row, yk-orb, yk-pill, deck-footer, slide-number, notes

---

## contents（气泡目录）
指纹：table
数量：yk-item=4

用途：议程页。一块气泡大卡里放 4 行篇目：铬面编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="yk-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="yk-card mt-l" style="margin-top:44px">
    <div class="yk-item"><span class="yk-n">01</span><span class="yk-t">{{篇名，≤8 字}}</span><span class="yk-d">{{说明，14-26 字}}</span></div>
    <div class="yk-item"><span class="yk-n">02</span><span class="yk-t">{{篇名，≤8 字}}</span><span class="yk-d">{{说明，14-26 字}}</span></div>
    <div class="yk-item"><span class="yk-n">03</span><span class="yk-t">{{篇名，≤8 字}}</span><span class="yk-d">{{说明，14-26 字}}</span></div>
    <div class="yk-item"><span class="yk-n">04</span><span class="yk-t">{{篇名，≤8 字}}</span><span class="yk-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, yk-kicker, h2, mt-m, yk-card, mt-l, yk-item, yk-n, yk-t, yk-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：yk-card=3

用途：恰好三张气泡卡。每张：紫色题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="yk-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="yk-card"><span class="yk-pill yk-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#47474d">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="yk-card"><span class="yk-pill yk-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#47474d">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="yk-card"><span class="yk-pill yk-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#47474d">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, yk-kicker, h2, mt-m, grid, g3, mt-l, yk-card, yk-pill, yk-pill-accent, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：yk-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边气泡卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="yk-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#47474d">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="yk-pill">{{要点 1，≤8 字}}</span>
        <span class="yk-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="yk-card">
      <div class="yk-step"><span class="yk-n">1</span><p class="yk-mini-t">{{一步，12-26 字}}</p></div>
      <div class="yk-step"><span class="yk-n">2</span><p class="yk-mini-t">{{一步，12-26 字}}</p></div>
      <div class="yk-step"><span class="yk-n">3</span><p class="yk-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, yk-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, yk-pill, yk-card, yk-step, yk-n, yk-mini-t, deck-footer, slide-number, notes

---

## metrics（铬面数字）
指纹：chart
数量：yk-stat=3

用途：三个关键数据。等宽大数字是主视觉，彩虹描边做身份，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="yk-kicker">{{数据语境}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="yk-stat"><div class="yk-stat-v">{{数值 ≤6 字符}}<span class="yk-stat-u">{{单位}}</span></div><div class="yk-stat-l">{{指标名，≤8 字}}</div><p class="yk-stat-note">{{口径，14-30 字}}</p></div>
    <div class="yk-stat"><div class="yk-stat-v">{{数值 ≤6 字符}}<span class="yk-stat-u">{{单位}}</span></div><div class="yk-stat-l">{{指标名，≤8 字}}</div><p class="yk-stat-note">{{口径，14-30 字}}</p></div>
    <div class="yk-stat"><div class="yk-stat-v">{{数值 ≤6 字符}}<span class="yk-stat-u">{{单位}}</span></div><div class="yk-stat-l">{{指标名，≤8 字}}</div><p class="yk-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#5d5d64;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, yk-kicker, h2, mt-m, grid, g3, mt-l, yk-stat, yk-stat-v, yk-stat-u, yk-stat-l, yk-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（高光引文）
指纹：quote

用途：整页一句引文。大字 + 出处 + 两个支撑药丸，右侧留白处立一枚铬面气泡徽标。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="yk-kicker">{{语境}}</p>
  <p class="yk-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="yk-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="yk-pill yk-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="yk-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="yk-n" style="position:absolute;right:130px;top:50%;transform:translateY(-50%);width:96px;height:96px;font-size:22px">{{徽标 2-4 字符}}</div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, yk-kicker, yk-quote, mt-l, yk-src, mt-m, row, yk-pill, yk-pill-accent, yk-n, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="yk-kicker">{{进度，如 第一幕 · 声量}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字，可 <br> 分行}}</h1>
  <div class="yk-rainbow mt-s"></div>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="yk-pill yk-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="yk-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, yk-kicker, h1, mt-m, yk-rainbow, mt-s, lede, mt-l, row, yk-pill, yk-pill-accent, deck-footer, slide-number, notes

---

## moments（潮流时间线）
指纹：chart
数量：yk-tl-item=4

用途：3-4 个节点的横向时间线：铬面圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="yk-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="yk-tl mt-l" style="margin-top:52px">
    <div class="yk-tl-item"><div class="yk-tl-dot">1</div><div class="yk-tl-t">{{时间点，≤10 字符}}</div><p class="yk-tl-d">{{事件，12-26 字}}</p></div>
    <div class="yk-tl-item"><div class="yk-tl-dot">2</div><div class="yk-tl-t">{{时间点，≤10 字符}}</div><p class="yk-tl-d">{{事件，12-26 字}}</p></div>
    <div class="yk-tl-item"><div class="yk-tl-dot">3</div><div class="yk-tl-t">{{时间点，≤10 字符}}</div><p class="yk-tl-d">{{事件，12-26 字}}</p></div>
    <div class="yk-tl-item"><div class="yk-tl-dot">4</div><div class="yk-tl-t">{{时间点，≤10 字符}}</div><p class="yk-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#5d5d64;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, yk-kicker, h2, mt-m, yk-tl, mt-l, yk-tl-item, yk-tl-dot, yk-tl-t, yk-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（收尾行动）
指纹：hero

用途：收尾页。大字 + 一句行动提醒 + 渐变按钮 + 药丸 + 铬面徽标。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="yk-kicker">{{提醒语境}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <div class="yk-rainbow mt-s"></div>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="yk-btn">{{按钮文案，≤8 字}}</span>
    <span class="yk-pill">{{次级信息，≤10 字}}</span>
    <div class="yk-orb">{{徽标 2-4 字符}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, yk-kicker, h1, mt-m, yk-rainbow, mt-s, lede, mt-l, row, yk-btn, yk-pill, yk-orb, deck-footer, slide-number, notes
