# 情绪疗愈色卡 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `hc-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-healing-color-card` 作用域生效，骨架里已写全，照抄结构即可。
> 右上月牙柔光与页底薄荷雾由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（柔软即秩序）**：三色母题只允许两种形态——色卡渐变（hc-card-lav / -peach / -mint）
> 与小圆色板（hc-chip，成组出现）；文字只用暖墨灰与柔棕，藕荷紫只出现在关键数据、圆点编号与
> 强调签。禁锋利直角、禁高饱和色块、禁 emoji——治愈感靠文字体温与留白，不靠表情符号。

---

## cover（柔光封面）
指纹：hero

用途：开场页。情绪题签 + 圆润大字标题 + 三色色板行 + 一句定位，右上月牙衬底。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 26-52 字；题签 ≤14 字；色板行固定三枚圆点。

```html
<section class="slide full" data-layout="cover">
  <p class="hc-kicker">{{题签，≤14 字，如 HEALING COLOR · 第 03 期}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="row mt-s" style="gap:12px"><span class="hc-chip" style="background:#C8B8D9"></span><span class="hc-chip" style="background:#E8CEC8"></span><span class="hc-chip" style="background:#9CD9C9"></span></div>
  <p class="lede mt-m" style="max-width:48ch">{{一句话定位，26-52 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="hc-pill hc-pill-accent">{{开课或主题信息，≤12 字}}</span>
    <span class="hc-pill">{{次级信息，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, hc-kicker, h1, mt-m, mt-s, row, hc-chip, lede, mt-l, hc-pill, hc-pill-accent, deck-footer, slide-number, notes

---

## contents（课表目录）
指纹：table
数量：hc-item=4

用途：议程页。一块奶油大卡里放 4 行篇目：圆点编号 + 篇名 + 一句说明，虚线分隔。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="hc-kicker">{{引导语，如 本期四周}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="hc-card mt-l" style="padding:18px 44px">
    <div class="hc-item"><span class="hc-n">01</span><span class="hc-t">{{篇名，≤8 字}}</span><span class="hc-d">{{说明，14-26 字}}</span></div>
    <div class="hc-item"><span class="hc-n">02</span><span class="hc-t">{{篇名，≤8 字}}</span><span class="hc-d">{{说明，14-26 字}}</span></div>
    <div class="hc-item"><span class="hc-n">03</span><span class="hc-t">{{篇名，≤8 字}}</span><span class="hc-d">{{说明，14-26 字}}</span></div>
    <div class="hc-item"><span class="hc-n">04</span><span class="hc-t">{{篇名，≤8 字}}</span><span class="hc-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, hc-kicker, h2, mt-m, hc-card, mt-l, hc-item, hc-n, hc-t, hc-d, deck-footer, slide-number, notes

---

## keynotes（三色情绪卡）
指纹：cards
数量：hc-card=3

用途：恰好三张渐变色卡（薰衣草/暖桃/薄荷各一）。每张：题签 + 小标题 + 两句说明 + 底部色板行。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="hc-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:30px;margin-top:44px">
    <div class="hc-card hc-card-lav"><span class="hc-pill hc-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.9;color:#4A4448">{{说明：一句论断 + 一句展开，22-44 字}}</p><div class="row mt-m" style="gap:10px;margin-top:auto"><span class="hc-chip" style="background:#9B7FB5"></span><span class="hc-hex">{{色号 · 情绪，如 #C8B8D9 安抚}}</span></div></div>
    <div class="hc-card hc-card-peach"><span class="hc-pill hc-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.9;color:#4A4448">{{说明：一句论断 + 一句展开，22-44 字}}</p><div class="row mt-m" style="gap:10px;margin-top:auto"><span class="hc-chip" style="background:#E8CEC8"></span><span class="hc-hex">{{色号 · 情绪，如 #E8CEC8 拥抱}}</span></div></div>
    <div class="hc-card hc-card-mint"><span class="hc-pill hc-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.9;color:#4A4448">{{说明：一句论断 + 一句展开，22-44 字}}</p><div class="row mt-m" style="gap:10px;margin-top:auto"><span class="hc-chip" style="background:#9CD9C9"></span><span class="hc-hex">{{色号 · 情绪，如 #9CD9C9 呼吸}}</span></div></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, hc-kicker, h2, mt-m, grid, g3, mt-l, hc-card, hc-card-lav, hc-card-peach, hc-card-mint, hc-pill, hc-pill-accent, h4, mt-s, row, hc-chip, hc-hex, deck-footer, slide-number, notes

---

## split（左文右卡）
指纹：split
数量：hc-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边奶油卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="hc-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.9;color:#6B5D5F">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:14px">
        <span class="hc-pill">{{要点 1，≤8 字}}</span>
        <span class="hc-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="hc-card" style="min-height:0">
      <div class="hc-step"><span class="hc-n">一</span><p class="hc-mini-t">{{一步，12-26 字}}</p></div>
      <div class="hc-step"><span class="hc-n">二</span><p class="hc-mini-t">{{一步，12-26 字}}</p></div>
      <div class="hc-step"><span class="hc-n">三</span><p class="hc-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, hc-kicker, h2, mt-m, grid, g2, mt-l, lede, row, hc-pill, hc-card, hc-step, hc-n, hc-mini-t, deck-footer, slide-number, notes

---

## metrics（柔光数字）
指纹：chart
数量：hc-stat=3

用途：三个关键数据。藕荷紫大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="hc-kicker">{{数据语境，如 春季班 · 结课复盘}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:48px">
    <div class="hc-stat"><div class="hc-stat-v">{{数值 ≤6 字符}}<span class="hc-stat-u">{{单位}}</span></div><div class="hc-stat-l">{{指标名，≤8 字}}</div><p class="hc-stat-note">{{口径，14-30 字}}</p></div>
    <div class="hc-stat"><div class="hc-stat-v">{{数值 ≤6 字符}}<span class="hc-stat-u">{{单位}}</span></div><div class="hc-stat-l">{{指标名，≤8 字}}</div><p class="hc-stat-note">{{口径，14-30 字}}</p></div>
    <div class="hc-stat"><div class="hc-stat-v">{{数值 ≤6 字符}}<span class="hc-stat-u">{{单位}}</span></div><div class="hc-stat-l">{{指标名，≤8 字}}</div><p class="hc-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="margin-top:40px;font-size:17px;color:#A3969A;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, hc-kicker, h2, mt-m, grid, g3, mt-l, hc-stat, hc-stat-v, hc-stat-u, hc-stat-l, hc-stat-note, deck-footer, slide-number, notes

---

## quote（手记引文）
指纹：quote

用途：整页一句引文。楷体大字 + 出处 + 两个支撑药丸，底部三色色板行收尾。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="hc-kicker">{{语境，如 主理人手记}}</p>
  <p class="hc-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="hc-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:14px">
    <span class="hc-pill hc-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="hc-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="row mt-l" style="gap:12px"><span class="hc-chip" style="background:#C8B8D9"></span><span class="hc-chip" style="background:#E8CEC8"></span><span class="hc-chip" style="background:#9CD9C9"></span><span class="hc-hex">THE HEALING PALETTE</span></div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, hc-kicker, hc-quote, mt-l, hc-src, mt-m, row, hc-pill, hc-pill-accent, hc-chip, hc-hex, deck-footer, slide-number, notes

---

## divider（章节柔幕）
指纹：hero

用途：章节过渡。进度题签 + 圆润大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="hc-kicker">{{进度，如 卷贰 · 三色安抚法}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:46ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:14px">
    <span class="hc-pill hc-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="hc-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, hc-kicker, h1, mt-m, lede, mt-l, row, hc-pill, hc-pill-accent, deck-footer, slide-number, notes

---

## moments（流程时间线）
指纹：chart
数量：hc-tl-item=4

用途：3-4 个节点的横向时间线：白底圆点编号 + 时间点 + 一句事件，柔光连接线。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="hc-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="hc-tl mt-l" style="margin-top:52px">
    <div class="hc-tl-item"><div class="hc-tl-dot">1</div><div class="hc-tl-t">{{时间点，≤10 字符}}</div><p class="hc-tl-d">{{事件，12-26 字}}</p></div>
    <div class="hc-tl-item"><div class="hc-tl-dot">2</div><div class="hc-tl-t">{{时间点，≤10 字符}}</div><p class="hc-tl-d">{{事件，12-26 字}}</p></div>
    <div class="hc-tl-item"><div class="hc-tl-dot">3</div><div class="hc-tl-t">{{时间点，≤10 字符}}</div><p class="hc-tl-d">{{事件，12-26 字}}</p></div>
    <div class="hc-tl-item"><div class="hc-tl-dot">4</div><div class="hc-tl-t">{{时间点，≤10 字符}}</div><p class="hc-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="margin-top:44px;font-size:17px;color:#A3969A;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, hc-kicker, h2, mt-m, hc-tl, mt-l, hc-tl-item, hc-tl-dot, hc-tl-t, hc-tl-d, deck-footer, slide-number, notes

---

## closing（温柔收尾）
指纹：hero

用途：收尾页。圆润大字 + 一句行动提醒 + 藕荷紫胶囊按钮 + 药丸 + 三色色板行。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="hc-kicker">{{提醒语境，如 第 03 期 · 报名}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:46ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:20px">
    <span class="hc-btn">{{按钮文案，≤8 字}}</span>
    <span class="hc-pill">{{次级信息，≤10 字}}</span>
  </div>
  <div class="row mt-l" style="gap:12px"><span class="hc-chip" style="background:#C8B8D9"></span><span class="hc-chip" style="background:#E8CEC8"></span><span class="hc-chip" style="background:#9CD9C9"></span></div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, hc-kicker, h1, mt-m, lede, mt-l, row, hc-btn, hc-pill, hc-chip, deck-footer, slide-number, notes
