# 奶油温柔 · 高级感 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `cp-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-cream-pastel` 作用域生效，骨架里已写全，照抄结构即可。
> 右上尤加利枝与页底抹茶雾由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（奶油即克制）**：焦糖只允许出现在细边——奶油卡描边（cp-card）、焦糖顶边
> （cp-stat）、强调签（cp-pill-accent）；豆沙粉只做题签实底（cp-kicker）与柔光；
> 抹茶绿只做圆号与节点。禁深色大面积底、禁硬直角、禁密集 KPI 卡。
> 信息密度低是这份模板的高级感来源：一页最多两组内容块 + 页脚。

---

## cover（奶油封面）
指纹：hero

用途：开场页。豆沙粉胶囊题签 + 温柔大字标题（可衬线点缀）+ 一句定位 + 两枚药丸。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 26-52 字；题签 ≤14 字；衬线点缀 ≤4 字。

```html
<section class="slide full" data-layout="cover">
  <span class="cp-kicker">{{题签，≤14 字，如 拾光软装 · 方案 A-112}}</span>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行，可嵌 <span class="cp-serif">衬线点缀 ≤4 字</span>}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句话定位，26-52 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="cp-pill">{{时间或期数，≤12 字}}</span>
    <span class="cp-pill cp-pill-accent">{{一句亮点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, cp-kicker, h1, mt-m, cp-serif, lede, mt-l, row, cp-pill, cp-pill-accent, deck-footer, slide-number, notes

---

## contents（方案目录）
指纹：table
数量：cp-item=4

用途：议程页。一块奶油大卡里放 4 行篇目：抹茶圆号 + 篇名 + 一句说明，细虚线分隔。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <span class="cp-kicker">{{引导语，如 本次提案}}</span>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="cp-card mt-l" style="padding:18px 44px">
    <div class="cp-item"><span class="cp-n">01</span><span class="cp-t">{{篇名，≤8 字}}</span><span class="cp-d">{{说明，14-26 字}}</span></div>
    <div class="cp-item"><span class="cp-n">02</span><span class="cp-t">{{篇名，≤8 字}}</span><span class="cp-d">{{说明，14-26 字}}</span></div>
    <div class="cp-item"><span class="cp-n">03</span><span class="cp-t">{{篇名，≤8 字}}</span><span class="cp-d">{{说明，14-26 字}}</span></div>
    <div class="cp-item"><span class="cp-n">04</span><span class="cp-t">{{篇名，≤8 字}}</span><span class="cp-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cp-kicker, h2, mt-m, cp-card, mt-l, cp-item, cp-n, cp-t, cp-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：cp-card=3

用途：恰好三张圆角奶油卡（焦糖细边）。每张：豆沙粉题签 + 小标题 + 两句说明 + 底部色样行。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <span class="cp-kicker">{{引导语}}</span>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:30px;margin-top:44px">
    <div class="cp-card"><span class="cp-pill">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.9;color:#5C5248">{{说明：一句论断 + 一句展开，22-44 字}}</p><div class="row mt-m" style="gap:12px;margin-top:auto"><span class="cp-sw" style="background:#FBF7F2"></span><span class="cp-hex">{{色号 色名，如 #FBF7F2 奶油白}}</span></div></div>
    <div class="cp-card"><span class="cp-pill">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.9;color:#5C5248">{{说明：一句论断 + 一句展开，22-44 字}}</p><div class="row mt-m" style="gap:12px;margin-top:auto"><span class="cp-sw" style="background:#F4E9E8"></span><span class="cp-hex">{{色号 色名，如 #F4E9E8 豆沙粉}}</span></div></div>
    <div class="cp-card"><span class="cp-pill">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.9;color:#5C5248">{{说明：一句论断 + 一句展开，22-44 字}}</p><div class="row mt-m" style="gap:12px;margin-top:auto"><span class="cp-sw" style="background:#A3B18A"></span><span class="cp-hex">{{色号 色名，如 #A3B18A 抹茶绿}}</span></div></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cp-kicker, h2, mt-m, grid, g3, mt-l, cp-card, cp-pill, h4, mt-s, row, cp-sw, cp-hex, deck-footer, slide-number, notes

---

## split（左文右卡）
指纹：split
数量：cp-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边焦糖左边条奶油卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <span class="cp-kicker">{{引导语}}</span>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.9;color:#7A6F65">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:14px">
        <span class="cp-pill">{{要点 1，≤8 字}}</span>
        <span class="cp-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="cp-card cp-card-warm" style="min-height:0">
      <div class="cp-step"><span class="cp-n">一</span><p class="cp-mini-t">{{一步，12-26 字}}</p></div>
      <div class="cp-step"><span class="cp-n">二</span><p class="cp-mini-t">{{一步，12-26 字}}</p></div>
      <div class="cp-step"><span class="cp-n">三</span><p class="cp-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cp-kicker, h2, mt-m, grid, g2, mt-l, lede, row, cp-pill, cp-card, cp-card-warm, cp-step, cp-n, cp-mini-t, deck-footer, slide-number, notes

---

## metrics（奶油数字）
指纹：chart
数量：cp-stat=3

用途：三个关键数据。豆沙色大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <span class="cp-kicker">{{数据语境，如 项目 A-112 · 报价口径}}</span>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:48px">
    <div class="cp-stat"><div class="cp-stat-v">{{数值 ≤6 字符}}<span class="cp-stat-u">{{单位}}</span></div><div class="cp-stat-l">{{指标名，≤8 字}}</div><p class="cp-stat-note">{{口径，14-30 字}}</p></div>
    <div class="cp-stat"><div class="cp-stat-v">{{数值 ≤6 字符}}<span class="cp-stat-u">{{单位}}</span></div><div class="cp-stat-l">{{指标名，≤8 字}}</div><p class="cp-stat-note">{{口径，14-30 字}}</p></div>
    <div class="cp-stat"><div class="cp-stat-v">{{数值 ≤6 字符}}<span class="cp-stat-u">{{单位}}</span></div><div class="cp-stat-l">{{指标名，≤8 字}}</div><p class="cp-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="margin-top:40px;font-size:17px;color:#A89E94;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cp-kicker, h2, mt-m, grid, g3, mt-l, cp-stat, cp-stat-v, cp-stat-u, cp-stat-l, cp-stat-note, deck-footer, slide-number, notes

---

## quote（宋体引文）
指纹：quote

用途：整页一句引文。宋体大字 + 出处 + 两个支撑药丸，留白即奶油质感。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <span class="cp-kicker">{{语境，如 主理人的话}}</span>
  <p class="cp-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="cp-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:14px">
    <span class="cp-pill cp-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="cp-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cp-kicker, cp-quote, mt-l, cp-src, mt-m, row, cp-pill, cp-pill-accent, deck-footer, slide-number, notes

---

## divider（章节柔幕）
指纹：hero

用途：章节过渡。进度题签 + 温柔大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <span class="cp-kicker">{{进度，如 第叁章 · 客厅方案}}</span>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:46ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:14px">
    <span class="cp-pill cp-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="cp-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, cp-kicker, h1, mt-m, lede, mt-l, row, cp-pill, cp-pill-accent, deck-footer, slide-number, notes

---

## moments（工期时间线）
指纹：chart
数量：cp-tl-item=4

用途：3-4 个节点的横向时间线：抹茶圆角节点 + 阶段时间点 + 一句事件，虚线连接。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <span class="cp-kicker">{{引导语}}</span>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="cp-tl mt-l" style="margin-top:52px">
    <div class="cp-tl-item"><div class="cp-tl-dot">01</div><div class="cp-tl-t">{{时间点，≤10 字符}}</div><p class="cp-tl-d">{{事件，12-26 字}}</p></div>
    <div class="cp-tl-item"><div class="cp-tl-dot">02</div><div class="cp-tl-t">{{时间点，≤10 字符}}</div><p class="cp-tl-d">{{事件，12-26 字}}</p></div>
    <div class="cp-tl-item"><div class="cp-tl-dot">03</div><div class="cp-tl-t">{{时间点，≤10 字符}}</div><p class="cp-tl-d">{{事件，12-26 字}}</p></div>
    <div class="cp-tl-item"><div class="cp-tl-dot">04</div><div class="cp-tl-t">{{时间点，≤10 字符}}</div><p class="cp-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="margin-top:44px;font-size:17px;color:#A89E94;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cp-kicker, h2, mt-m, cp-tl, mt-l, cp-tl-item, cp-tl-dot, cp-tl-t, cp-tl-d, deck-footer, slide-number, notes

---

## closing（温柔收尾）
指纹：hero

用途：收尾页。温柔大字 + 一句行动提醒 + 深棕胶囊按钮 + 药丸，如一封奶油信笺。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <span class="cp-kicker">{{提醒语境，如 预约量房}}</span>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:46ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="cp-btn">{{按钮文案，≤8 字}}</span>
    <span class="cp-pill">{{次级信息，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, cp-kicker, h1, mt-m, lede, mt-l, row, cp-btn, cp-pill, deck-footer, slide-number, notes
