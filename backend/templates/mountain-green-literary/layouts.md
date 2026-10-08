# 山野葱郁 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `mg-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-mountain-green-literary` 作用域生效，骨架里已写全，照抄结构即可。
> 每页顶部的日光飞鸟与页底的山峦草地由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（图重字轻）**：山绿 #3E7228 是身份色，只做描边、编号圆点、大数字与按钮；
> 草黄 #ECCB2D 只做笔触与标签底（黄上深字，禁黄字上浅底）；页底草地之上仍要保证正文可读，
> 禁把内容卡片压到山形顶部以下再叠第二层色块。大标题（h1/h2）每页最多一组；
> 竖排 mg-vert 只放留白处，一页至多一条。

---

## cover（山野封面）
指纹：hero

用途：开场页。草黄标签 + 粗笔大字标题 + 一句定位，页底山峦草地自动铺陈。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；标签 ≤12 字；药丸 ≤12 字。

```html
<section class="slide full" data-layout="cover">
  <p class="mg-kicker">{{题签，≤14 字，如 七月企划 · 山见民宿}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="mg-chip">{{地点或主题，≤12 字}}</span>
    <span class="mg-pill">{{时间或节点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, mg-kicker, h1, mt-m, lede, mt-l, row, mg-chip, mg-pill, deck-footer, slide-number, notes

---

## contents（企划目录）
指纹：table
数量：mg-item=4

用途：议程页。一张纸卡里放 4 行篇目：色卡圆点编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 12-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="mg-kicker">{{引导语，如 企划卷目}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="mg-card mt-l" style="margin-top:44px">
    <div class="mg-item"><span class="mg-n">壹</span><span class="mg-t">{{篇名，≤8 字}}</span><span class="mg-d">{{说明，12-26 字}}</span></div>
    <div class="mg-item"><span class="mg-n">贰</span><span class="mg-t">{{篇名，≤8 字}}</span><span class="mg-d">{{说明，12-26 字}}</span></div>
    <div class="mg-item"><span class="mg-n">叁</span><span class="mg-t">{{篇名，≤8 字}}</span><span class="mg-d">{{说明，12-26 字}}</span></div>
    <div class="mg-item"><span class="mg-n">肆</span><span class="mg-t">{{篇名，≤8 字}}</span><span class="mg-d">{{说明，12-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mg-kicker, h2, mt-m, mg-card, mt-l, mg-item, mg-n, mg-t, mg-d, deck-footer, slide-number, notes

---

## keynotes（三张木牌）
指纹：cards
数量：mg-card=3

用途：恰好三张纸卡。每张：草黄标签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；标签 ≤6 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="mg-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="mg-card"><span class="mg-chip">{{标签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#4E6640">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="mg-card"><span class="mg-chip">{{标签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#4E6640">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="mg-card"><span class="mg-chip">{{标签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#4E6640">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mg-kicker, h2, mt-m, grid, g3, mt-l, mg-card, mg-chip, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：mg-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边纸卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="mg-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#4E6640">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="mg-pill">{{要点 1，≤8 字}}</span>
        <span class="mg-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="mg-card">
      <div class="mg-step"><span class="mg-n">一</span><p class="mg-mini-t">{{一步，12-26 字}}</p></div>
      <div class="mg-step"><span class="mg-n">二</span><p class="mg-mini-t">{{一步，12-26 字}}</p></div>
      <div class="mg-step"><span class="mg-n">三</span><p class="mg-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mg-kicker, h2, mt-m, grid, g2, mt-l, lede, row, mg-pill, mg-card, mg-step, mg-n, mg-mini-t, deck-footer, slide-number, notes

---

## metrics（松白数字）
指纹：chart
数量：mg-stat=3

用途：三个关键数据。深绿大数字是主视觉，指标名带草黄笔触，口径写进条目，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="mg-kicker">{{数据语境，如 经营台账 · 复盘}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="mg-stat"><div class="mg-stat-v">{{数值 ≤6 字符}}<span class="mg-stat-u">{{单位}}</span></div><div class="mg-stat-l">{{指标名，≤8 字}}</div><p class="mg-stat-note">{{口径，14-30 字}}</p></div>
    <div class="mg-stat"><div class="mg-stat-v">{{数值 ≤6 字符}}<span class="mg-stat-u">{{单位}}</span></div><div class="mg-stat-l">{{指标名，≤8 字}}</div><p class="mg-stat-note">{{口径，14-30 字}}</p></div>
    <div class="mg-stat"><div class="mg-stat-v">{{数值 ≤6 字符}}<span class="mg-stat-u">{{单位}}</span></div><div class="mg-stat-l">{{指标名，≤8 字}}</div><p class="mg-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#7E9070;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mg-kicker, h2, mt-m, grid, g3, mt-l, mg-stat, mg-stat-v, mg-stat-u, mg-stat-l, mg-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（手记引文）
指纹：quote

用途：整页一句引文。粗笔大字 + 出处 + 两个支撑标签，右侧留白处立一条竖排诗句。
适用 role：quote。
内容约束：引文 10-24 字；出处 ≤22 字且真实；标签各 ≤10 字；竖排诗句 ≤7 字。

```html
<section class="slide" data-layout="quote">
  <p class="mg-kicker">{{语境，如 主理人手记}}</p>
  <p class="mg-quote mt-l" style="margin-top:40px">「{{引文，10-24 字}}」</p>
  <p class="mg-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="mg-chip">{{支撑点 1，≤10 字}}</span>
    <span class="mg-pill">{{支撑点 2，≤10 字}}</span>
  </div>
  <div class="mg-vert" style="position:absolute;right:120px;top:50%;transform:translateY(-50%)">{{竖排诗句 ≤7 字}}<span class="mg-vert-accent">{{一两字}}</span></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mg-kicker, mg-quote, mt-l, mg-src, mt-m, row, mg-chip, mg-pill, mg-vert, mg-vert-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度标签 + 粗笔大字章节名 + 一个过渡问题 + 两个看点。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；标签各 ≤10 字。

```html
<section class="slide full" data-layout="divider">
  <p class="mg-kicker">{{进度，如 第二幕 · 山野体验单}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="mg-chip">{{看点 1，≤10 字}}</span>
    <span class="mg-pill">{{看点 2，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, mg-kicker, h1, mt-m, lede, mt-l, row, mg-chip, mg-pill, deck-footer, slide-number, notes

---

## moments（档期时间线）
指纹：chart
数量：mg-tl-item=4

用途：3-4 个节点的横向时间线：色卡圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="mg-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="mg-tl mt-l" style="margin-top:52px">
    <div class="mg-tl-item"><div class="mg-tl-dot">壹</div><div class="mg-tl-t">{{时间点，≤10 字符}}</div><p class="mg-tl-d">{{事件，12-26 字}}</p></div>
    <div class="mg-tl-item"><div class="mg-tl-dot">贰</div><div class="mg-tl-t">{{时间点，≤10 字符}}</div><p class="mg-tl-d">{{事件，12-26 字}}</p></div>
    <div class="mg-tl-item"><div class="mg-tl-dot">叁</div><div class="mg-tl-t">{{时间点，≤10 字符}}</div><p class="mg-tl-d">{{事件，12-26 字}}</p></div>
    <div class="mg-tl-item"><div class="mg-tl-dot">肆</div><div class="mg-tl-t">{{时间点，≤10 字符}}</div><p class="mg-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#7E9070;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mg-kicker, h2, mt-m, mg-tl, mt-l, mg-tl-item, mg-tl-dot, mg-tl-t, mg-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（收尾页）
指纹：hero

用途：收尾页。粗笔大字 + 一句行动提醒 + 深绿按钮 + 药丸，像手账落款。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤12 字。

```html
<section class="slide full" data-layout="closing">
  <p class="mg-kicker">{{提醒语境，如 七月看房}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:22px">
    <span class="mg-btn">{{按钮文案，≤8 字}}</span>
    <span class="mg-pill">{{次级信息，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, mg-kicker, h1, mt-m, lede, mt-l, row, mg-btn, mg-pill, deck-footer, slide-number, notes
