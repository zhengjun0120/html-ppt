# 经典双色·米黄深蓝 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `cd-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-classic-duo-blue` 作用域生效，骨架里已写全，照抄结构即可。
> 四角金色书角与右上暖调纸感由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（双色即身份）**：全模板只有两主色——暖米黄与深海蓝，金棕只作细线与序号点缀。
> 深海蓝大面积实底只允许三处：标识块（cd-badge）、书脊色块（cd-panel）、按钮（cd-btn）；
> 禁把米黄或深蓝外的第三种主色引入版面。衬线大字（h1/h2）每页最多一组；
> 竖排 cd-vert 只出现在书脊色块内，一页至多一条。

---

## cover（典籍封面）
指纹：hero

用途：开场页。金线题签 + 衬线大字标题 + 金棕细线 + 一句定位，左下立深蓝标识块。
适用 role：cover。
内容约束：主标题 ≤12 字；lede 24-52 字；题签 ≤14 字；标识 2-4 字。

```html
<section class="slide full" data-layout="cover">
  <p class="cd-kicker">{{题签，≤14 字，如 致瑞制造 · 供应链中心}}</p>
  <h1 class="h1 mt-m">{{主标题，≤12 字，可 <br> 分行}}</h1>
  <div class="cd-rule mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-52 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="cd-badge">{{标识 2-4 字}}</div>
    <span class="cd-pill">{{时间或场合，≤14 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, cd-kicker, h1, mt-m, cd-rule, mt-s, lede, mt-l, row, cd-badge, cd-pill, deck-footer, slide-number, notes

---

## contents（目录卷页）
指纹：table
数量：cd-item=4

用途：议程页。一块米色卡里放 4 行篇目：金棕衬线序号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="cd-kicker">{{引导语，如 方案目录}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="cd-card mt-l" style="margin-top:44px">
    <div class="cd-item"><span class="cd-n">01</span><span class="cd-t">{{篇名，≤8 字}}</span><span class="cd-d">{{说明，14-26 字}}</span></div>
    <div class="cd-item"><span class="cd-n">02</span><span class="cd-t">{{篇名，≤8 字}}</span><span class="cd-d">{{说明，14-26 字}}</span></div>
    <div class="cd-item"><span class="cd-n">03</span><span class="cd-t">{{篇名，≤8 字}}</span><span class="cd-d">{{说明，14-26 字}}</span></div>
    <div class="cd-item"><span class="cd-n">04</span><span class="cd-t">{{篇名，≤8 字}}</span><span class="cd-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cd-kicker, h2, mt-m, cd-card, mt-l, cd-item, cd-n, cd-t, cd-d, deck-footer, slide-number, notes

---

## keynotes（三卡举措）
指纹：cards
数量：cd-card=3

用途：恰好三张米色卡。每张：金棕题签 + 衬线小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤6 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="cd-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="cd-card"><span class="cd-pill cd-pill-accent">{{题签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#5A4E38">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="cd-card"><span class="cd-pill cd-pill-accent">{{题签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#5A4E38">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="cd-card"><span class="cd-pill cd-pill-accent">{{题签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#5A4E38">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cd-kicker, h2, mt-m, grid, g3, mt-l, cd-card, cd-pill, cd-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：cd-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边米色卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="cd-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#5A4E38">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="cd-pill">{{要点 1，≤8 字}}</span>
        <span class="cd-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="cd-card">
      <div class="cd-step"><span class="cd-n">一</span><p class="cd-mini-t">{{一步，12-26 字}}</p></div>
      <div class="cd-step"><span class="cd-n">二</span><p class="cd-mini-t">{{一步，12-26 字}}</p></div>
      <div class="cd-step"><span class="cd-n">三</span><p class="cd-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cd-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, cd-pill, cd-card, cd-step, cd-n, cd-mini-t, deck-footer, slide-number, notes

---

## metrics（基线数字）
指纹：chart
数量：cd-stat=3

用途：三个关键基线数据。深蓝衬线大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="cd-kicker">{{数据语境，如 现状基线 · 2025}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:40px;margin-top:48px">
    <div class="cd-stat"><div class="cd-stat-v">{{数值 ≤6 字符}}<span class="cd-stat-u">{{单位}}</span></div><div class="cd-stat-l">{{指标名，≤8 字}}</div><p class="cd-stat-note">{{口径，14-30 字}}</p></div>
    <div class="cd-stat"><div class="cd-stat-v">{{数值 ≤6 字符}}<span class="cd-stat-u">{{单位}}</span></div><div class="cd-stat-l">{{指标名，≤8 字}}</div><p class="cd-stat-note">{{口径，14-30 字}}</p></div>
    <div class="cd-stat"><div class="cd-stat-v">{{数值 ≤6 字符}}<span class="cd-stat-u">{{单位}}</span></div><div class="cd-stat-l">{{指标名，≤8 字}}</div><p class="cd-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="margin-top:40px;font-size:17px;color:#97876B;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cd-kicker, h2, mt-m, grid, g3, mt-l, cd-stat, cd-stat-v, cd-stat-u, cd-stat-l, cd-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（题铭引言）
指纹：quote

用途：整页一句引言。衬线大字 + 出处 + 两个支撑药丸，右侧深蓝书脊色块衬竖排题铭。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤24 字且真实；药丸各 ≤8 字；竖排题铭 ≤7 字。

```html
<section class="slide" data-layout="quote">
  <p class="cd-kicker">{{语境，如 评审寄语}}</p>
  <p class="cd-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="cd-src mt-m">—— {{出处：人与场合，≤24 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="cd-pill cd-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="cd-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="cd-panel" style="position:absolute;right:0;top:0;bottom:0"><span class="cd-vert">{{竖排题铭 ≤7 字}}</span></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cd-kicker, cd-quote, mt-l, cd-src, mt-m, row, cd-pill, cd-pill-accent, cd-panel, cd-vert, deck-footer, slide-number, notes

---

## divider（章节扉页）
指纹：hero

用途：章节过渡。进度题签 + 衬线大字章节名 + 一句过渡 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤12 字；过渡句 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="cd-kicker">{{进度，如 PART 02 · 三大举措}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="cd-pill cd-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="cd-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, cd-kicker, h1, mt-m, lede, mt-l, row, cd-pill, cd-pill-accent, deck-footer, slide-number, notes

---

## moments（季度路线）
指纹：chart
数量：cd-tl-item=4

用途：4 个节点的横向时间线：衬线序号圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：恰好 4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="cd-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="cd-tl mt-l" style="margin-top:52px">
    <div class="cd-tl-item"><div class="cd-tl-dot">01</div><div class="cd-tl-t">{{时间点，≤10 字符}}</div><p class="cd-tl-d">{{事件，12-26 字}}</p></div>
    <div class="cd-tl-item"><div class="cd-tl-dot">02</div><div class="cd-tl-t">{{时间点，≤10 字符}}</div><p class="cd-tl-d">{{事件，12-26 字}}</p></div>
    <div class="cd-tl-item"><div class="cd-tl-dot">03</div><div class="cd-tl-t">{{时间点，≤10 字符}}</div><p class="cd-tl-d">{{事件，12-26 字}}</p></div>
    <div class="cd-tl-item"><div class="cd-tl-dot">04</div><div class="cd-tl-t">{{时间点，≤10 字符}}</div><p class="cd-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="margin-top:40px;font-size:17px;color:#97876B;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cd-kicker, h2, mt-m, cd-tl, mt-l, cd-tl-item, cd-tl-dot, cd-tl-t, cd-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（落款收尾）
指纹：hero

用途：收尾页。衬线大字 + 一句行动提醒 + 深蓝按钮 + 描边药丸，如典籍落款。
适用 role：thanks / cta / content。
内容约束：标题 ≤12 字；lede 20-48 字；按钮 ≤6 字；药丸 ≤14 字。

```html
<section class="slide full" data-layout="closing">
  <p class="cd-kicker">{{行动语境}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="cd-btn">{{按钮文案，≤6 字}}</span>
    <span class="cd-pill">{{次级信息，≤14 字}}</span>
    <div class="cd-badge">{{标识 2-4 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, cd-kicker, h1, mt-m, lede, mt-l, row, cd-btn, cd-pill, cd-badge, deck-footer, slide-number, notes
