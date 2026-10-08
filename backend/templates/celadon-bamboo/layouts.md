# 青瓷竹影 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `qc-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-celadon-bamboo` 作用域生效，骨架里已写全，照抄结构即可。
> 每页右侧竹简竖条与左下竹叶薄雾由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（一页一个核心观点）**：朱砂只允许出现在两处——印章（qc-seal）与强调点缀
> （qc-pill-accent / qc-verse-accent）；关键数据用深青（qc-stat-v / qc-tl-t），不用朱砂。
> 青瓷绿只出现在细线、圈码与竹影 ambient 里，禁大面积色块。禁硬阴影、禁强投影、
> 禁强对比商业按钮。衬线大字（h1/h2）每页最多一组；竖排 qc-verse 只放留白处，一页至多一条。

---

## cover（题跋封面）
指纹：hero

用途：开场页。题签 + 衬线大字标题 + 细线 + 一句定位，留白处可立朱砂小印。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；印章 4 字内。

```html
<section class="slide full" data-layout="cover">
  <p class="qc-kicker">{{题签，≤14 字，如 龙泉青瓷展 · 展陈导览}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="qc-line mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="qc-seal">{{印章 2×2 字}}</div>
    <span class="qc-pill">{{副题或时间地点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, qc-kicker, h1, mt-m, qc-line, mt-s, lede, mt-l, row, qc-seal, qc-pill, deck-footer, slide-number, notes

---

## contents（目录四进）
指纹：table
数量：qc-item=4

用途：议程页。一块纸面大卡里放 4 行展段：圈码编号 + 展段名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；展段名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="qc-kicker">{{引导语，如 观展次第}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="qc-card mt-l" style="margin-top:44px">
    <div class="qc-item"><span class="qc-n">壹</span><span class="qc-t">{{展段名，≤8 字}}</span><span class="qc-d">{{说明，14-26 字}}</span></div>
    <div class="qc-item"><span class="qc-n">贰</span><span class="qc-t">{{展段名，≤8 字}}</span><span class="qc-d">{{说明，14-26 字}}</span></div>
    <div class="qc-item"><span class="qc-n">叁</span><span class="qc-t">{{展段名，≤8 字}}</span><span class="qc-d">{{说明，14-26 字}}</span></div>
    <div class="qc-item"><span class="qc-n">肆</span><span class="qc-t">{{展段名，≤8 字}}</span><span class="qc-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, qc-kicker, h2, mt-m, qc-card, mt-l, qc-item, qc-n, qc-t, qc-d, deck-footer, slide-number, notes

---

## keynotes（三器要点）
指纹：cards
数量：qc-card=3

用途：恰好三张纸面卡。每张：朱砂题签 + 衬线小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="qc-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="qc-card"><span class="qc-pill qc-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#55625B">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="qc-card"><span class="qc-pill qc-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#55625B">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="qc-card"><span class="qc-pill qc-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#55625B">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, qc-kicker, h2, mt-m, grid, g3, mt-l, qc-card, qc-pill, qc-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右简）
指纹：split
数量：qc-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边纸面卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="qc-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#55625B">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="qc-pill">{{要点 1，≤8 字}}</span>
        <span class="qc-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="qc-card">
      <div class="qc-step"><span class="qc-n">一</span><p class="qc-mini-t">{{一步，12-26 字}}</p></div>
      <div class="qc-step"><span class="qc-n">二</span><p class="qc-mini-t">{{一步，12-26 字}}</p></div>
      <div class="qc-step"><span class="qc-n">三</span><p class="qc-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, qc-kicker, h2, mt-m, grid, g2, mt-l, lede, row, qc-pill, qc-card, qc-step, qc-n, qc-mini-t, deck-footer, slide-number, notes

---

## metrics（器物数字）
指纹：chart
数量：qc-stat=3

用途：三个关键数据。深青大数字是唯一主视觉（朱砂留给印），口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="qc-kicker">{{数据语境，如 展览规模}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="qc-stat"><div class="qc-stat-v">{{数值 ≤6 字符}}<span class="qc-stat-u">{{单位}}</span></div><div class="qc-stat-l">{{指标名，≤8 字}}</div><p class="qc-stat-note">{{口径，14-30 字}}</p></div>
    <div class="qc-stat"><div class="qc-stat-v">{{数值 ≤6 字符}}<span class="qc-stat-u">{{单位}}</span></div><div class="qc-stat-l">{{指标名，≤8 字}}</div><p class="qc-stat-note">{{口径，14-30 字}}</p></div>
    <div class="qc-stat"><div class="qc-stat-v">{{数值 ≤6 字符}}<span class="qc-stat-u">{{单位}}</span></div><div class="qc-stat-l">{{指标名，≤8 字}}</div><p class="qc-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#8A958C;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, qc-kicker, h2, mt-m, grid, g3, mt-l, qc-stat, qc-stat-v, qc-stat-u, qc-stat-l, qc-stat-note, deck-footer, slide-number, notes

---

## quote（翠色引文）
指纹：quote

用途：整页一句引文。衬线大字 + 题跋小字出处 + 两个支撑药丸，右侧留白处立竖排短句。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="qc-kicker">{{语境，如 唐人题跋}}</p>
  <p class="qc-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="qc-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="qc-pill qc-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="qc-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="qc-verse" style="position:absolute;right:120px;top:50%;transform:translateY(-50%)">{{竖排短句 ≤7 字}}<span class="qc-verse-accent">{{一字或两字}}</span></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, qc-kicker, qc-quote, mt-l, qc-src, mt-m, row, qc-pill, qc-pill-accent, qc-verse, qc-verse-accent, deck-footer, slide-number, notes

---

## divider（章节竹影）
指纹：hero

用途：章节过渡。进度题签 + 衬线大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="qc-kicker">{{进度，如 第二进 · 问窑}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="qc-pill qc-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="qc-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, qc-kicker, h1, mt-m, lede, mt-l, row, qc-pill, qc-pill-accent, deck-footer, slide-number, notes

---

## moments（展期时间线）
指纹：chart
数量：qc-tl-item=4

用途：3-4 个节点的横向时间线：圈码圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="qc-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="qc-tl mt-l" style="margin-top:52px">
    <div class="qc-tl-item"><div class="qc-tl-dot">壹</div><div class="qc-tl-t">{{时间点，≤10 字符}}</div><p class="qc-tl-d">{{事件，12-26 字}}</p></div>
    <div class="qc-tl-item"><div class="qc-tl-dot">贰</div><div class="qc-tl-t">{{时间点，≤10 字符}}</div><p class="qc-tl-d">{{事件，12-26 字}}</p></div>
    <div class="qc-tl-item"><div class="qc-tl-dot">叁</div><div class="qc-tl-t">{{时间点，≤10 字符}}</div><p class="qc-tl-d">{{事件，12-26 字}}</p></div>
    <div class="qc-tl-item"><div class="qc-tl-dot">肆</div><div class="qc-tl-t">{{时间点，≤10 字符}}</div><p class="qc-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#8A958C;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, qc-kicker, h2, mt-m, qc-tl, mt-l, qc-tl-item, qc-tl-dot, qc-tl-t, qc-tl-d, deck-footer, slide-number, notes

---

## closing（收尾钤印）
指纹：hero

用途：收尾页。衬线大字 + 一句行动提醒 + 细描边按钮 + 描边药丸，如卷尾钤印。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="qc-kicker">{{提醒语境，如 观展预约}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="qc-btn">{{按钮文案，≤8 字}}</span>
    <span class="qc-pill">{{次级信息，≤10 字}}</span>
    <div class="qc-seal">{{印章 2 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, qc-kicker, h1, mt-m, lede, mt-l, row, qc-btn, qc-pill, qc-seal, deck-footer, slide-number, notes
