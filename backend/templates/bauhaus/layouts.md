# 包豪斯 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `bh-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-bauhaus` 作用域生效，骨架里已写全，照抄结构即可。
> 暖色画布与角落几何静物由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（几何是主角）**：红黄蓝只用在几何与编号上——题签方块、编号色块、几何静物
> （bh-circ / bh-tri / bh-rect）、时间线圆点、巨号数字与行动钮。卡永远是藏青描边白卡
> （bh-card），禁把卡做成色块、禁圆角阴影。几何静物一页至多一组，正文信息保持低到中密度，
> 不要把每个几何区域都填成内容卡。

---

## cover（几何封面）
指纹：hero

用途：开场页。红方块题签 + 900 大字标题 + 一句定位，页脚上方摆圆、三角、方的几何静物台。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-52 字；标签各 ≤6 字。

```html
<section class="slide full" data-layout="cover">
  <p class="bh-kicker">{{课程标签，≤14 字}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:46ch">{{定位：从哪讲起、讲什么，24-52 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="bh-tag">{{标签 1，≤6 字}}</span>
    <span class="bh-tag">{{标签 2，≤6 字}}</span>
  </div>
  <div class="bh-stage mt-l" style="margin-top:48px">
    <div class="bh-circ" style="width:120px;height:120px"></div>
    <div class="bh-tri" style="width:112px;height:112px"></div>
    <div class="bh-rect" style="width:190px;height:96px"></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 课程}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, bh-kicker, h1, mt-m, lede, row, mt-l, bh-tag, bh-stage, bh-circ, bh-tri, bh-rect, deck-footer, slide-number, notes

---

## contents（课堂目录）
指纹：table
数量：bh-item=4

用途：议程页。藏青描边白卡里放 4 行篇目：编号色块（红蓝金藏青轮换）+ 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="bh-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="bh-card mt-l" style="margin-top:44px">
    <div class="bh-item"><span class="bh-n">壹</span><span class="bh-t">{{篇名，≤8 字}}</span><span class="bh-d">{{说明，14-26 字}}</span></div>
    <div class="bh-item"><span class="bh-n">贰</span><span class="bh-t">{{篇名，≤8 字}}</span><span class="bh-d">{{说明，14-26 字}}</span></div>
    <div class="bh-item"><span class="bh-n">叁</span><span class="bh-t">{{篇名，≤8 字}}</span><span class="bh-d">{{说明，14-26 字}}</span></div>
    <div class="bh-item"><span class="bh-n">肆</span><span class="bh-t">{{篇名，≤8 字}}</span><span class="bh-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 课程}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bh-kicker, h2, mt-m, bh-card, mt-l, bh-item, bh-n, bh-t, bh-d, deck-footer, slide-number, notes

---

## metrics（原色数据）
指纹：chart
数量：bh-stat=3

用途：三个关键数据。藏青顶线统计块 + 900 巨号（红蓝藏青轮换），口径写页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="bh-kicker">{{数据语境}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="bh-stat"><div class="bh-stat-v">{{数值 ≤6 字符}}</div><div class="bh-stat-l">{{指标名，≤8 字}}</div><p class="bh-stat-note">{{口径，14-30 字}}</p></div>
    <div class="bh-stat"><div class="bh-stat-v">{{数值 ≤6 字符}}</div><div class="bh-stat-l">{{指标名，≤8 字}}</div><p class="bh-stat-note">{{口径，14-30 字}}</p></div>
    <div class="bh-stat"><div class="bh-stat-v">{{数值 ≤6 字符}}</div><div class="bh-stat-l">{{指标名，≤8 字}}</div><p class="bh-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="bh-src mt-m">来源：{{出处与统计口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 课程}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bh-kicker, h2, mt-m, grid, g3, mt-l, bh-stat, bh-stat-v, bh-stat-l, bh-stat-note, bh-src, mt-m, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 900 大字章节名 + 过渡问题 + 两个标签 + 一排缩小几何静物。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；标签各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="bh-kicker">{{进度，如 第二讲 · 原色的脾气}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:46ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="bh-tag">{{看点 1，≤8 字}}</span>
    <span class="bh-tag">{{看点 2，≤8 字}}</span>
  </div>
  <div class="bh-stage mt-l" style="margin-top:44px">
    <div class="bh-circ" style="width:64px;height:64px"></div>
    <div class="bh-tri" style="width:60px;height:60px"></div>
    <div class="bh-rect" style="width:110px;height:52px"></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 课程}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, bh-kicker, h1, mt-m, lede, row, mt-l, bh-tag, bh-stage, bh-circ, bh-tri, bh-rect, deck-footer, slide-number, notes

---

## keynotes（三形要点）
指纹：cards
数量：bh-card=3

用途：恰好三张藏青描边白卡。每张：描边标签（红蓝藏青轮换）+ 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；标签 ≤6 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="bh-kicker">{{章节标签}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="bh-card"><span class="bh-tag">{{标签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:18px;line-height:1.8;color:#43536F">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="bh-card"><span class="bh-tag">{{标签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:18px;line-height:1.8;color:#43536F">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="bh-card"><span class="bh-tag">{{标签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:18px;line-height:1.8;color:#43536F">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 课程}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bh-kicker, h2, mt-m, grid, g3, mt-l, bh-card, bh-tag, h4, mt-s, deck-footer, slide-number, notes

---

## split（课堂拆解）
指纹：split
数量：bh-step=3

用途：左边把一段史实讲透（lede + 补充 + 标签），右边白卡装三行练习步骤。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个标签；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="bh-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:64px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：史实或练习的来龙去脉，36-70 字}}</p>
      <p class="mt-m" style="font-size:18px;line-height:1.8;color:#43536F">{{补充：一句判断或提醒，18-40 字}}</p>
      <div class="row mt-l" style="gap:14px">
        <span class="bh-tag">{{要点 1，≤8 字}}</span>
        <span class="bh-tag">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="bh-card">
      <div class="bh-step"><span class="bh-step-i">一</span><p class="bh-mini-t">{{一步，12-26 字}}</p></div>
      <div class="bh-step"><span class="bh-step-i">二</span><p class="bh-mini-t">{{一步，12-26 字}}</p></div>
      <div class="bh-step"><span class="bh-step-i">三</span><p class="bh-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 课程}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bh-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, bh-tag, bh-card, bh-step, bh-step-i, bh-mini-t, deck-footer, slide-number, notes

---

## moments（校史时间线）
指纹：chart
数量：bh-tl-item=4

用途：3-4 个节点的横向时间线：几何圆点（圆、方、菱、环轮换）+ 红色时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="bh-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="bh-tl mt-l" style="margin-top:56px">
    <div class="bh-tl-item"><div class="bh-tl-dot"></div><div class="bh-tl-t">{{时间点，≤10 字符}}</div><p class="bh-tl-d">{{事件，12-26 字}}</p></div>
    <div class="bh-tl-item"><div class="bh-tl-dot"></div><div class="bh-tl-t">{{时间点，≤10 字符}}</div><p class="bh-tl-d">{{事件，12-26 字}}</p></div>
    <div class="bh-tl-item"><div class="bh-tl-dot"></div><div class="bh-tl-t">{{时间点，≤10 字符}}</div><p class="bh-tl-d">{{事件，12-26 字}}</p></div>
    <div class="bh-tl-item"><div class="bh-tl-dot"></div><div class="bh-tl-t">{{时间点，≤10 字符}}</div><p class="bh-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="bh-src mt-m">{{口径或读法，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 课程}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bh-kicker, h2, mt-m, bh-tl, mt-l, bh-tl-item, bh-tl-dot, bh-tl-t, bh-tl-d, bh-src, mt-m, deck-footer, slide-number, notes

---

## quote（大师之言）
指纹：quote

用途：整页一句大师引文。900 大字 + 出处 + 两个描边标签。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；标签各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="bh-kicker">{{语境，如 大师课}}</p>
  <p class="bh-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="bh-src mt-m">—— {{出处：人与身份，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="bh-tag">{{支撑点 1，≤8 字}}</span>
    <span class="bh-tag">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 课程}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bh-kicker, bh-quote, mt-l, bh-src, mt-m, row, bh-tag, deck-footer, slide-number, notes

---

## closing（布置作业）
指纹：hero

用途：收尾页。900 大字 + 一句作业提醒 + 红底藏青描边行动钮 + 标签。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-48 字；按钮 ≤6 字；标签各 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="bh-kicker">{{语境，如 课后作业}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:46ch">{{作业内容与截止提醒，20-48 字}}</p>
  <div class="row mt-l" style="gap:20px">
    <span class="bh-btn">{{按钮文案，≤6 字}}</span>
    <span class="bh-tag">{{次级信息，≤10 字}}</span>
    <span class="bh-tag">{{次级信息，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 课程}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, bh-kicker, h1, mt-m, lede, row, mt-l, bh-btn, bh-tag, deck-footer, slide-number, notes
