# 中式水墨意境·山水 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `cl-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-chinese-ink-landscape` 作用域生效，骨架里已写全，照抄结构即可。
> 每页上半区的两抹云雾与页底三层山峦由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（留白即云雾）**：朱砂只允许出现在三处——印章（cl-seal）、题签细线
> （cl-kicker 自带）、关键数据（cl-stat-v / cl-tl-t / cl-pill-accent / cl-vert-accent）。
> 青黛蓝只出现在山棱线（cl-ridge）与数字块顶线（cl-stat），禁大面积色块。
> 禁硬阴影、禁复杂渐变、禁无衬线粗体大标题。宋体大字（h1/h2）每页最多一组；
> 竖排 cl-vert 只放留白处，一页至多一条；色卡 cl-swatches 只在 cover / closing 出现。

---

## cover（云山封面）
指纹：hero

用途：开场页。题签 + 宋体大字标题 + 山棱线 + 一句定位，右下可立朱砂印与四色色卡。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；印章 4 字内。

```html
<section class="slide full" data-layout="cover">
  <p class="cl-kicker">{{题签，≤14 字，如 山地文旅提案 · 丙午年秋}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="cl-ridge mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="cl-seal">{{印章 2×2 字}}</div>
    <span class="cl-pill">{{副题或时间地点，≤12 字}}</span>
    <div class="cl-swatches"><span class="cl-sw" style="background:#4A90A4"></span><span class="cl-sw" style="background:#8B9DC3"></span><span class="cl-sw" style="background:#D4C4A0"></span><span class="cl-sw" style="background:#2C3E50"></span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, cl-kicker, h1, mt-m, cl-ridge, mt-s, lede, mt-l, row, cl-seal, cl-pill, cl-swatches, cl-sw, deck-footer, slide-number, notes

---

## contents（目录四卷）
指纹：table
数量：cl-item=4

用途：议程页。一块宣纸大卡里放 4 行卷目：墨圈编号 + 卷名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；卷名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="cl-kicker">{{引导语，如 今日议程}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="cl-card mt-l" style="margin-top:44px">
    <div class="cl-item"><span class="cl-n">壹</span><span class="cl-t">{{卷名，≤8 字}}</span><span class="cl-d">{{说明，14-26 字}}</span></div>
    <div class="cl-item"><span class="cl-n">贰</span><span class="cl-t">{{卷名，≤8 字}}</span><span class="cl-d">{{说明，14-26 字}}</span></div>
    <div class="cl-item"><span class="cl-n">叁</span><span class="cl-t">{{卷名，≤8 字}}</span><span class="cl-d">{{说明，14-26 字}}</span></div>
    <div class="cl-item"><span class="cl-n">肆</span><span class="cl-t">{{卷名，≤8 字}}</span><span class="cl-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cl-kicker, h2, mt-m, cl-card, mt-l, cl-item, cl-n, cl-t, cl-d, deck-footer, slide-number, notes

---

## keynotes（三境要点）
指纹：cards
数量：cl-card=3

用途：恰好三张宣纸卡。每张：朱砂题签 + 宋体小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="cl-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="cl-card"><span class="cl-pill cl-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#5D6D7E">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="cl-card"><span class="cl-pill cl-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#5D6D7E">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="cl-card"><span class="cl-pill cl-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#5D6D7E">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cl-kicker, h2, mt-m, grid, g3, mt-l, cl-card, cl-pill, cl-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右径）
指纹：split
数量：cl-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边宣纸卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="cl-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#5D6D7E">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="cl-pill">{{要点 1，≤8 字}}</span>
        <span class="cl-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="cl-card">
      <div class="cl-step"><span class="cl-n">一</span><p class="cl-mini-t">{{一步，12-26 字}}</p></div>
      <div class="cl-step"><span class="cl-n">二</span><p class="cl-mini-t">{{一步，12-26 字}}</p></div>
      <div class="cl-step"><span class="cl-n">三</span><p class="cl-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cl-kicker, h2, mt-m, grid, g2, mt-l, lede, row, cl-pill, cl-card, cl-step, cl-n, cl-mini-t, deck-footer, slide-number, notes

---

## metrics（山河数字）
指纹：chart
数量：cl-stat=3

用途：三个关键数据。朱砂大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="cl-kicker">{{数据语境，如 一期 · 核心数字}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="cl-stat"><div class="cl-stat-v">{{数值 ≤6 字符}}<span class="cl-stat-u">{{单位}}</span></div><div class="cl-stat-l">{{指标名，≤8 字}}</div><p class="cl-stat-note">{{口径，14-30 字}}</p></div>
    <div class="cl-stat"><div class="cl-stat-v">{{数值 ≤6 字符}}<span class="cl-stat-u">{{单位}}</span></div><div class="cl-stat-l">{{指标名，≤8 字}}</div><p class="cl-stat-note">{{口径，14-30 字}}</p></div>
    <div class="cl-stat"><div class="cl-stat-v">{{数值 ≤6 字符}}<span class="cl-stat-u">{{单位}}</span></div><div class="cl-stat-l">{{指标名，≤8 字}}</div><p class="cl-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#8C9BAA;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cl-kicker, h2, mt-m, grid, g3, mt-l, cl-stat, cl-stat-v, cl-stat-u, cl-stat-l, cl-stat-note, deck-footer, slide-number, notes

---

## quote（绝顶引文）
指纹：quote

用途：整页一句引文。宋体大字 + 出处 + 两个支撑药丸，右侧留白处立竖排诗句。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="cl-kicker">{{语境，如 山客题跋}}</p>
  <p class="cl-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="cl-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="cl-pill cl-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="cl-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="cl-vert" style="position:absolute;right:120px;top:50%;transform:translateY(-50%)">{{竖排诗句 ≤7 字}}<span class="cl-vert-accent">{{一字或两字}}</span></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cl-kicker, cl-quote, mt-l, cl-src, mt-m, row, cl-pill, cl-pill-accent, cl-vert, cl-vert-accent, deck-footer, slide-number, notes

---

## divider（章节云雾）
指纹：hero

用途：章节过渡。进度题签 + 宋体大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="cl-kicker">{{进度，如 卷二 · 景区三境}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="cl-pill cl-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="cl-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, cl-kicker, h1, mt-m, lede, mt-l, row, cl-pill, cl-pill-accent, deck-footer, slide-number, notes

---

## moments（季节时间线）
指纹：chart
数量：cl-tl-item=4

用途：3-4 个节点的横向时间线：墨圈圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="cl-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="cl-tl mt-l" style="margin-top:52px">
    <div class="cl-tl-item"><div class="cl-tl-dot">壹</div><div class="cl-tl-t">{{时间点，≤10 字符}}</div><p class="cl-tl-d">{{事件，12-26 字}}</p></div>
    <div class="cl-tl-item"><div class="cl-tl-dot">贰</div><div class="cl-tl-t">{{时间点，≤10 字符}}</div><p class="cl-tl-d">{{事件，12-26 字}}</p></div>
    <div class="cl-tl-item"><div class="cl-tl-dot">叁</div><div class="cl-tl-t">{{时间点，≤10 字符}}</div><p class="cl-tl-d">{{事件，12-26 字}}</p></div>
    <div class="cl-tl-item"><div class="cl-tl-dot">肆</div><div class="cl-tl-t">{{时间点，≤10 字符}}</div><p class="cl-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#8C9BAA;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cl-kicker, h2, mt-m, cl-tl, mt-l, cl-tl-item, cl-tl-dot, cl-tl-t, cl-tl-d, deck-footer, slide-number, notes

---

## closing（收尾远望）
指纹：hero

用途：收尾页。宋体大字 + 一句行动提醒 + 墨框按钮 + 描边药丸，如登高后的落款。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="cl-kicker">{{提醒语境，如 提案合作}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="cl-btn">{{按钮文案，≤8 字}}</span>
    <span class="cl-pill">{{次级信息，≤10 字}}</span>
    <div class="cl-seal">{{印章 2 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, cl-kicker, h1, mt-m, lede, mt-l, row, cl-btn, cl-pill, cl-seal, deck-footer, slide-number, notes
