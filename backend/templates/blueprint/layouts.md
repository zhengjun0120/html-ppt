# 蓝图 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `bp-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-blueprint` 作用域生效，骨架里已写全，照抄结构即可。
> 每页细网格与双线图框由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（制图纪律）**：全站等宽字体、方角、无阴影、无暖色。亮蓝线色只允许出现在
> 图签章（bp-stamp）、题签方点、编号框（bp-n / bp-tl-dot / bp-btn）与数据顶线（bp-stat）
> 这几处线稿上；正文一律浅蓝白（#DBEAFE/#93C5FD）。标注框（bp-box）是唯一的卡，
> 一页至多一组图签章。

---

## cover（图框封面）
指纹：hero

用途：开场页。图纸题签 + 等宽大标题 + 尺寸标注线 + 一句定位，右下可盖方案图签章。
适用 role：cover。
内容约束：主标题 ≤12 字；lede 24-48 字；题签 ≤14 字；图签章 4 字内。

```html
<section class="slide full" data-layout="cover">
  <p class="bp-kicker">{{题签，≤14 字，如 建筑设计提案 · 社区图书馆}}</p>
  <h1 class="h1 mt-m">{{主标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:52ch">{{一句话定位，24-48 字}}</p>
  <div class="bp-dim mt-m" style="max-width:640px">{{关键尺寸串，如 用地 5,200 ㎡ · 建面 3,860 ㎡}}</div>
  <div class="row mt-l" style="gap:30px">
    <div class="bp-stamp">{{图签 2×2 字}}</div>
    <span class="bp-pill">{{副题或阶段，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 图号}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, bp-kicker, h1, mt-m, lede, bp-dim, row, mt-l, bp-stamp, bp-pill, deck-footer, slide-number, notes

---

## contents（图纸目录）
指纹：table
数量：bp-item=4

用途：议程页。一个图例大框里放 4 行篇目：方框编号 + 篇名 + 一句说明，框头是图例行。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="bp-kicker">{{引导语，如图纸目录}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="bp-box mt-l" style="margin-top:44px">
    <div class="bp-box-head"><span>{{图例标签，如 LEGEND · 图例}}</span><span>{{图号，如 SH-01}}</span></div>
    <div class="bp-item"><span class="bp-n">壹</span><span class="bp-t">{{篇名，≤8 字}}</span><span class="bp-d">{{说明，14-26 字}}</span></div>
    <div class="bp-item"><span class="bp-n">贰</span><span class="bp-t">{{篇名，≤8 字}}</span><span class="bp-d">{{说明，14-26 字}}</span></div>
    <div class="bp-item"><span class="bp-n">叁</span><span class="bp-t">{{篇名，≤8 字}}</span><span class="bp-d">{{说明，14-26 字}}</span></div>
    <div class="bp-item"><span class="bp-n">肆</span><span class="bp-t">{{篇名，≤8 字}}</span><span class="bp-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 图号}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bp-kicker, h2, mt-m, bp-box, mt-l, bp-box-head, bp-item, bp-n, bp-t, bp-d, deck-footer, slide-number, notes

---

## keynotes（三框要点）
指纹：cards
数量：bp-box=3

用途：恰好三个标注框。每框：图例框头（题签 + 图号）+ 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="bp-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="bp-box"><div class="bp-box-head"><span>{{题签，≤8 字}}</span><span>{{图号，如 D-01}}</span></div><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#93C5FD">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="bp-box"><div class="bp-box-head"><span>{{题签，≤8 字}}</span><span>{{图号，如 D-02}}</span></div><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#93C5FD">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="bp-box"><div class="bp-box-head"><span>{{题签，≤8 字}}</span><span>{{图号，如 D-03}}</span></div><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#93C5FD">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 图号}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bp-kicker, h2, mt-m, grid, g3, mt-l, bp-box, bp-box-head, h4, mt-s, deck-footer, slide-number, notes

---

## split（主图侧栏）
指纹：split
数量：bp-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边标注框装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="bp-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#93C5FD">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="bp-pill">{{要点 1，≤8 字}}</span>
        <span class="bp-pill bp-pill-accent">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="bp-box">
      <div class="bp-step"><span class="bp-n">一</span><p class="bp-mini-t">{{一步，12-26 字}}</p></div>
      <div class="bp-step"><span class="bp-n">二</span><p class="bp-mini-t">{{一步，12-26 字}}</p></div>
      <div class="bp-step"><span class="bp-n">三</span><p class="bp-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 图号}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bp-kicker, h2, mt-m, grid, g2, mt-l, lede, row, bp-pill, bp-pill-accent, bp-box, bp-step, bp-n, bp-mini-t, deck-footer, slide-number, notes

---

## metrics（尺寸标注）
指纹：chart
数量：bp-stat=3

用途：三个关键数据，如图纸上引出的一组尺寸。口径写进框内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="bp-kicker">{{数据语境，如 方案总图 · 读数}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="bp-stat"><div class="bp-stat-v">{{数值 ≤6 字符}}<span class="bp-stat-u">{{单位}}</span></div><div class="bp-stat-l">{{指标名，≤8 字}}</div><p class="bp-stat-note">{{口径，14-30 字}}</p></div>
    <div class="bp-stat"><div class="bp-stat-v">{{数值 ≤6 字符}}<span class="bp-stat-u">{{单位}}</span></div><div class="bp-stat-l">{{指标名，≤8 字}}</div><p class="bp-stat-note">{{口径，14-30 字}}</p></div>
    <div class="bp-stat"><div class="bp-stat-v">{{数值 ≤6 字符}}<span class="bp-stat-u">{{单位}}</span></div><div class="bp-stat-l">{{指标名，≤8 字}}</div><p class="bp-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#7FA6DC;letter-spacing:.06em">来源：{{出处与统计口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 图号}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bp-kicker, h2, mt-m, grid, g3, mt-l, bp-stat, bp-stat-v, bp-stat-u, bp-stat-l, bp-stat-note, deck-footer, slide-number, notes

---

## quote（图注引言）
指纹：quote

用途：整页一句引文。等宽大字 + 出处 + 两个支撑药丸，右侧留白处盖一枚图签章。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="bp-kicker">{{语境，如 主持建筑师图注}}</p>
  <p class="bp-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="bp-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="bp-pill bp-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="bp-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="bp-stamp" style="position:absolute;right:130px;top:50%;transform:translateY(-50%)">{{图签 2×2 字}}</div>
  <div class="deck-footer"><span>{{署名 · 图号}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bp-kicker, bp-quote, mt-l, bp-src, mt-m, row, bp-pill, bp-pill-accent, bp-stamp, deck-footer, slide-number, notes

---

## divider（分幅）
指纹：hero

用途：章节过渡。进度题签 + 等宽大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="bp-kicker">{{进度，如 第二幅 · 设计要点}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="bp-pill bp-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="bp-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 图号}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, bp-kicker, h1, mt-m, lede, mt-l, row, bp-pill, bp-pill-accent, deck-footer, slide-number, notes

---

## moments（施工轴）
指纹：chart
数量：bp-tl-item=4

用途：3-4 个节点的横向时间线：方框节点 + 时间点 + 一句事件，如施工横道图。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="bp-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="bp-tl mt-l" style="margin-top:52px">
    <div class="bp-tl-item"><div class="bp-tl-dot">壹</div><div class="bp-tl-t">{{时间点，≤10 字符}}</div><p class="bp-tl-d">{{事件，12-26 字}}</p></div>
    <div class="bp-tl-item"><div class="bp-tl-dot">贰</div><div class="bp-tl-t">{{时间点，≤10 字符}}</div><p class="bp-tl-d">{{事件，12-26 字}}</p></div>
    <div class="bp-tl-item"><div class="bp-tl-dot">叁</div><div class="bp-tl-t">{{时间点，≤10 字符}}</div><p class="bp-tl-d">{{事件，12-26 字}}</p></div>
    <div class="bp-tl-item"><div class="bp-tl-dot">肆</div><div class="bp-tl-t">{{时间点，≤10 字符}}</div><p class="bp-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#7FA6DC;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 图号}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bp-kicker, h2, mt-m, bp-tl, mt-l, bp-tl-item, bp-tl-dot, bp-tl-t, bp-tl-d, deck-footer, slide-number, notes

---

## closing（签章收尾）
指纹：hero

用途：收尾页。等宽大字 + 一句行动提醒 + 双线按钮 + 描边药丸 + 图签章，如图纸报审。
适用 role：thanks / cta / content。
内容约束：标题 ≤12 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="bp-kicker">{{提醒语境，如 评审安排}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:30px">
    <span class="bp-btn">{{按钮文案，≤8 字}}</span>
    <span class="bp-pill">{{次级信息，≤10 字}}</span>
    <div class="bp-stamp">{{图签 2 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 图号}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, bp-kicker, h1, mt-m, lede, mt-l, row, bp-btn, bp-pill, bp-stamp, deck-footer, slide-number, notes
