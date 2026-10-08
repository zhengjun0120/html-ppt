# 故宫墨红 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `pi-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-palace-ink-red` 作用域生效，骨架里已写全，照抄结构即可。
> 每页左侧墨红竖脊（正文页）与顶部墨红宽条（封面/章节/收尾页）、右上角窗棂格，
> 由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（大色块即分区）**：墨红只允许出现在六处——自动衬底色条、短尺 pi-rule、
> 方印 pi-seal、编号块 pi-num、关键数据与强调（pi-stat-v / pi-tl-t / pi-pill-red）、
> 实底按钮 pi-btn。禁把墨红用于整页背景或正文长段落。黑体大字（h1/h2）每页最多一组；
> 玄石灰深卡 pi-card-dark 一页至多一块；卡是大圆角厚卡，禁改直角细边。

---

## cover（展板封面）
指纹：hero

用途：开场页。题签 + 黑体大字标题 + 墨红短尺 + 一句定位，右下可立一枚方印。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；方印 4 字内。

```html
<section class="slide full" data-layout="cover">
  <p class="pi-kicker">{{题签，≤14 字，如 品牌提案 · 甲辰年}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="pi-rule mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="pi-seal">{{方印 2×2 字}}</div>
    <span class="pi-pill">{{副题或时间地点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, pi-kicker, h1, mt-m, pi-rule, mt-s, lede, mt-l, row, pi-seal, pi-pill, deck-footer, slide-number, notes

---

## contents（目录板）
指纹：table
数量：pi-item=4

用途：议程页。一块浅宣卡里放 4 行条目：墨红编号块 + 条目名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；条目名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="pi-kicker">{{引导语，如 今日议程}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="pi-card mt-l" style="margin-top:44px">
    <div class="pi-item"><span class="pi-num">壹</span><span class="pi-item-t">{{条目名，≤8 字}}</span><span class="pi-item-d">{{说明，14-26 字}}</span></div>
    <div class="pi-item"><span class="pi-num">贰</span><span class="pi-item-t">{{条目名，≤8 字}}</span><span class="pi-item-d">{{说明，14-26 字}}</span></div>
    <div class="pi-item"><span class="pi-num">叁</span><span class="pi-item-t">{{条目名，≤8 字}}</span><span class="pi-item-d">{{说明，14-26 字}}</span></div>
    <div class="pi-item"><span class="pi-num">肆</span><span class="pi-item-t">{{条目名，≤8 字}}</span><span class="pi-item-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pi-kicker, h2, mt-m, pi-card, mt-l, pi-item, pi-num, pi-item-t, pi-item-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：pi-card=3

用途：恰好三张浅宣卡。每张：强调胶囊 + 黑体小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="pi-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="pi-card"><span class="pi-pill pi-pill-red">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#555B54">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="pi-card"><span class="pi-pill pi-pill-red">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#555B54">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="pi-card"><span class="pi-pill pi-pill-red">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#555B54">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pi-kicker, h2, mt-m, grid, g3, mt-l, pi-card, pi-pill, pi-pill-red, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左文右卡）
指纹：split
数量：pi-step=3

用途：左边把一件事讲透（lede + 补充 + 胶囊），右玄石灰深卡装三行工序或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个胶囊；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="pi-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#555B54">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="pi-pill">{{要点 1，≤8 字}}</span>
        <span class="pi-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="pi-card-dark">
      <div class="pi-step"><span class="pi-num">一</span><p class="pi-step-t">{{一步，12-26 字}}</p></div>
      <div class="pi-step"><span class="pi-num">二</span><p class="pi-step-t">{{一步，12-26 字}}</p></div>
      <div class="pi-step"><span class="pi-num">三</span><p class="pi-step-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pi-kicker, h2, mt-m, grid, g2, mt-l, lede, row, pi-pill, pi-card-dark, pi-step, pi-num, pi-step-t, deck-footer, slide-number, notes

---

## metrics（朱红数字）
指纹：chart
数量：pi-stat=3

用途：三个关键数据。墨红大数字是唯一主视觉，口径写进数据板，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="pi-kicker">{{数据语境，如 去年台账 · 复盘}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="pi-stat"><div class="pi-stat-v">{{数值 ≤6 字符}}<span class="pi-stat-u">{{单位}}</span></div><div class="pi-stat-l">{{指标名，≤8 字}}</div><p class="pi-stat-note">{{口径，14-30 字}}</p></div>
    <div class="pi-stat"><div class="pi-stat-v">{{数值 ≤6 字符}}<span class="pi-stat-u">{{单位}}</span></div><div class="pi-stat-l">{{指标名，≤8 字}}</div><p class="pi-stat-note">{{口径，14-30 字}}</p></div>
    <div class="pi-stat"><div class="pi-stat-v">{{数值 ≤6 字符}}<span class="pi-stat-u">{{单位}}</span></div><div class="pi-stat-l">{{指标名，≤8 字}}</div><p class="pi-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#525850;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pi-kicker, h2, mt-m, grid, g3, mt-l, pi-stat, pi-stat-v, pi-stat-u, pi-stat-l, pi-stat-note, deck-footer, slide-number, notes

---

## quote（题铭引文）
指纹：quote

用途：整页一句引文。黑体大字 + 出处 + 两个支撑胶囊，像展板上的一句话题铭。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；胶囊各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="pi-kicker">{{语境，如 主理人说}}</p>
  <p class="pi-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="pi-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="pi-pill pi-pill-red">{{支撑点 1，≤8 字}}</span>
    <span class="pi-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pi-kicker, pi-quote, mt-l, pi-src, mt-m, row, pi-pill, pi-pill-red, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 大字章节名 + 一个过渡问题 + 两个看点胶囊。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；胶囊各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="pi-kicker">{{进度，如 卷二 · 制器之道}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="pi-pill pi-pill-red">{{看点 1，≤8 字}}</span>
    <span class="pi-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, pi-kicker, h1, mt-m, lede, mt-l, row, pi-pill, pi-pill-red, deck-footer, slide-number, notes

---

## moments（流程时间线）
指纹：chart
数量：pi-tl-item=4

用途：3-4 个节点的横向时间线：墨红编号块 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="pi-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="pi-tl mt-l" style="margin-top:52px">
    <div class="pi-tl-item"><div class="pi-tl-dot">壹</div><div class="pi-tl-t">{{时间点，≤10 字符}}</div><p class="pi-tl-d">{{事件，12-26 字}}</p></div>
    <div class="pi-tl-item"><div class="pi-tl-dot">贰</div><div class="pi-tl-t">{{时间点，≤10 字符}}</div><p class="pi-tl-d">{{事件，12-26 字}}</p></div>
    <div class="pi-tl-item"><div class="pi-tl-dot">叁</div><div class="pi-tl-t">{{时间点，≤10 字符}}</div><p class="pi-tl-d">{{事件，12-26 字}}</p></div>
    <div class="pi-tl-item"><div class="pi-tl-dot">肆</div><div class="pi-tl-t">{{时间点，≤10 字符}}</div><p class="pi-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#525850;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pi-kicker, h2, mt-m, pi-tl, mt-l, pi-tl-item, pi-tl-dot, pi-tl-t, pi-tl-d, deck-footer, slide-number, notes

---

## closing（收尾落款）
指纹：hero

用途：收尾页。黑体大字 + 一句行动提醒 + 墨红实底按钮 + 胶囊与方印，如展板落款。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；胶囊 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="pi-kicker">{{提醒语境，如 订购方式}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="pi-btn">{{按钮文案，≤8 字}}</span>
    <span class="pi-pill">{{次级信息，≤10 字}}</span>
    <div class="pi-seal">{{方印 2 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, pi-kicker, h1, mt-m, lede, mt-l, row, pi-btn, pi-pill, pi-seal, deck-footer, slide-number, notes
