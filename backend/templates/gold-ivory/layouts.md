# 鎏金象牙 · 高级质感 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `gi-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-gold-ivory` 作用域生效，骨架里已写全，照抄结构即可。
> 细密斜纹与双线金框（外线 40px / 内线 8px）由模板自动衬底（.slide 伪元素，
> pointer-events:none），骨架与正文都不用画，也禁再画第二道框。
>
> **气质铁律（奢华是留白）**：鎏金棕只做大字与色块——h1/h2、gi-t、gi-quote、gi-stat-v、
> gi-mark、gi-btn、gi-pill-accent 底；香槟金只做细线与菱形花饰（gi-orn、gi-n、gi-stat
> 上缘线、gi-tl 连线）。金色渐变字 gi-gold 只许用于标题局部 2-4 字，全份 deck 至多两处。
> 禁硬阴影、禁大圆角、禁无衬线粗体大标题；对称稳重，一页最多两组内容块 + 页脚。

---

## cover（典藏封面）
指纹：hero

用途：开场页。题签 + 衬线大字标题 + 花饰金线 + 一句定位，可立一枚鎏金印记。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；印记 4 字内。

```html
<section class="slide full" data-layout="cover">
  <p class="gi-kicker">{{题签，≤14 字，如 丙午年 · 年度典藏发布}}</p>
  <h1 class="h1 mt-m"><span class="gi-gold">{{局部 2-4 字}}</span>{{其余标题，合计 ≤10 字，可 <br> 分行}}</h1>
  <div class="gi-orn mt-m">◆</div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="gi-mark">{{印记 2×2 字}}</div>
    <span class="gi-pill">{{副题或时间地点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, gi-kicker, h1, mt-m, gi-gold, gi-orn, lede, mt-l, row, gi-mark, gi-pill, deck-footer, slide-number, notes

---

## contents（鉴赏目录）
指纹：table
数量：gi-item=4

用途：议程页。一块象牙大卡里放 4 行篇目：方框衬线编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="gi-kicker">{{引导语，如 今夜卷目}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="gi-card mt-l" style="margin-top:44px">
    <div class="gi-item"><span class="gi-n">壹</span><span class="gi-t">{{篇名，≤8 字}}</span><span class="gi-d">{{说明，14-26 字}}</span></div>
    <div class="gi-item"><span class="gi-n">贰</span><span class="gi-t">{{篇名，≤8 字}}</span><span class="gi-d">{{说明，14-26 字}}</span></div>
    <div class="gi-item"><span class="gi-n">叁</span><span class="gi-t">{{篇名，≤8 字}}</span><span class="gi-d">{{说明，14-26 字}}</span></div>
    <div class="gi-item"><span class="gi-n">肆</span><span class="gi-t">{{篇名，≤8 字}}</span><span class="gi-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gi-kicker, h2, mt-m, gi-card, mt-l, gi-item, gi-n, gi-t, gi-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：gi-card=3

用途：恰好三张象牙卡。每张：棕底编号签 + 衬线小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="gi-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="gi-card"><span class="gi-pill gi-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:18px;line-height:1.8;color:#5A3A1F">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="gi-card"><span class="gi-pill gi-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:18px;line-height:1.8;color:#5A3A1F">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="gi-card"><span class="gi-pill gi-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:18px;line-height:1.8;color:#5A3A1F">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gi-kicker, h2, mt-m, grid, g3, mt-l, gi-card, gi-pill, gi-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：gi-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边象牙卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="gi-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:18px;line-height:1.8;color:#5A3A1F">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="gi-pill">{{要点 1，≤8 字}}</span>
        <span class="gi-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="gi-card">
      <div class="gi-step"><span class="gi-n">一</span><p class="gi-mini-t">{{一步，12-26 字}}</p></div>
      <div class="gi-step"><span class="gi-n">二</span><p class="gi-mini-t">{{一步，12-26 字}}</p></div>
      <div class="gi-step"><span class="gi-n">三</span><p class="gi-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gi-kicker, h2, mt-m, grid, g2, mt-l, lede, row, gi-pill, gi-card, gi-step, gi-n, gi-mini-t, deck-footer, slide-number, notes

---

## metrics（鎏金数字）
指纹：chart
数量：gi-stat=3

用途：三个关键数据。衬线大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="gi-kicker">{{数据语境，如 本季典藏 · 总览}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="gi-stat"><div class="gi-stat-v">{{数值 ≤6 字符}}<span class="gi-stat-u">{{单位}}</span></div><div class="gi-stat-l">{{指标名，≤8 字}}</div><p class="gi-stat-note">{{口径，14-30 字}}</p></div>
    <div class="gi-stat"><div class="gi-stat-v">{{数值 ≤6 字符}}<span class="gi-stat-u">{{单位}}</span></div><div class="gi-stat-l">{{指标名，≤8 字}}</div><p class="gi-stat-note">{{口径，14-30 字}}</p></div>
    <div class="gi-stat"><div class="gi-stat-v">{{数值 ≤6 字符}}<span class="gi-stat-u">{{单位}}</span></div><div class="gi-stat-l">{{指标名，≤8 字}}</div><p class="gi-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#8C7A63;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gi-kicker, h2, mt-m, grid, g3, mt-l, gi-stat, gi-stat-v, gi-stat-u, gi-stat-l, gi-stat-note, deck-footer, slide-number, notes

---

## quote（题铭引言）
指纹：quote

用途：整页一句引文。衬线大字 + 出处 + 两个支撑药丸，右侧留白立一枚花饰。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="gi-kicker">{{语境，如 主理人题铭}}</p>
  <p class="gi-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="gi-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="gi-pill gi-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="gi-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="gi-orn" style="position:absolute;right:120px;top:50%;transform:translateY(-50%)">❦</div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gi-kicker, gi-quote, mt-l, gi-src, mt-m, row, gi-pill, gi-pill-accent, gi-orn, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 衬线大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="gi-kicker">{{进度，如 卷二 · 工艺细览}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="gi-pill gi-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="gi-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, gi-kicker, h1, mt-m, lede, mt-l, row, gi-pill, gi-pill-accent, deck-footer, slide-number, notes

---

## moments（展期时间线）
指纹：chart
数量：gi-tl-item=4

用途：3-4 个节点的横向时间线：菱形节点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="gi-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="gi-tl mt-l" style="margin-top:52px">
    <div class="gi-tl-item"><div class="gi-tl-dot"></div><div class="gi-tl-t">{{时间点，≤10 字符}}</div><p class="gi-tl-d">{{事件，12-26 字}}</p></div>
    <div class="gi-tl-item"><div class="gi-tl-dot"></div><div class="gi-tl-t">{{时间点，≤10 字符}}</div><p class="gi-tl-d">{{事件，12-26 字}}</p></div>
    <div class="gi-tl-item"><div class="gi-tl-dot"></div><div class="gi-tl-t">{{时间点，≤10 字符}}</div><p class="gi-tl-d">{{事件，12-26 字}}</p></div>
    <div class="gi-tl-item"><div class="gi-tl-dot"></div><div class="gi-tl-t">{{时间点，≤10 字符}}</div><p class="gi-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#8C7A63;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gi-kicker, h2, mt-m, gi-tl, mt-l, gi-tl-item, gi-tl-dot, gi-tl-t, gi-tl-d, deck-footer, slide-number, notes

---

## closing（收尾柬帖）
指纹：hero

用途：收尾页。衬线大字 + 一句行动提醒 + 棕底按钮 + 药丸，如一张柬帖落款。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="gi-kicker">{{提醒语境，如 鉴赏预约}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="gi-btn">{{按钮文案，≤8 字}}</span>
    <span class="gi-pill">{{次级信息，≤10 字}}</span>
    <div class="gi-mark">{{印记 2 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, gi-kicker, h1, mt-m, lede, mt-l, row, gi-btn, gi-pill, gi-mark, deck-footer, slide-number, notes
