# 凝脂杨妃 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `pr-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-chinese-porcelain-rose` 作用域生效，骨架里已写全，照抄结构即可。
> 每页顶部釉色渐变细带与右下缠枝柔线由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（釉色即身份）**：杨妃粉只做"面上的一层釉"——细带、釉线 pr-line、瓷圈 pr-n、
> hex pr-sw-hex 与胶囊描边；胭脂只允许出现在五处——小印 pr-seal、强调胶囊 pr-pill-red、
> 关键数据（pr-stat-v / pr-tl-t）、实底按钮 pr-btn、竖排重点字 pr-vert-accent。
> 禁大面积粉底整卡铺色、禁硬阴影、禁黑体粗标题。楷体大字（色名/展区名）每页至多一组；
> 竖排 pr-vert 只放留白处，一页至多一条。

---

## cover（瓷面封面）
指纹：hero

用途：开场页。题签 + 宋体大字标题 + 釉线 + 一句定位，右下留白处可立一枚胭脂小印。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；小印 4 字内。

```html
<section class="slide full" data-layout="cover">
  <p class="pr-kicker">{{题签，≤14 字，如 青凝堂 · 瓷艺生活展}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="pr-line mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="pr-seal">{{小印 2×2 字}}</div>
    <span class="pr-pill">{{副题或时间地点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, pr-kicker, h1, mt-m, pr-line, mt-s, lede, mt-l, row, pr-seal, pr-pill, deck-footer, slide-number, notes

---

## contents（展目）
指纹：table
数量：pr-item=4

用途：议程页。一块象牙卡里放 4 行展目：瓷圈编号 + 展区名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；展区名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="pr-kicker">{{引导语，如 今日导览}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="pr-card mt-l" style="margin-top:44px">
    <div class="pr-item"><span class="pr-n">壹</span><span class="pr-t">{{展区名，≤8 字}}</span><span class="pr-d">{{说明，14-26 字}}</span></div>
    <div class="pr-item"><span class="pr-n">贰</span><span class="pr-t">{{展区名，≤8 字}}</span><span class="pr-d">{{说明，14-26 字}}</span></div>
    <div class="pr-item"><span class="pr-n">叁</span><span class="pr-t">{{展区名，≤8 字}}</span><span class="pr-d">{{说明，14-26 字}}</span></div>
    <div class="pr-item"><span class="pr-n">肆</span><span class="pr-t">{{展区名，≤8 字}}</span><span class="pr-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pr-kicker, h2, mt-m, pr-card, mt-l, pr-item, pr-n, pr-t, pr-d, deck-footer, slide-number, notes

---

## keynotes（釉色三卡）
指纹：cards
数量：pr-swatch=3

用途：恰好三张釉色色卡。每张：楷体色名大字 + 拼音小字 + 一句释义 + hex。
适用 role：content。
内容约束：恰好 3 卡；色名 2-3 字；拼音 ≤8 字符；释义 22-44 字；hex 为真实色值。

```html
<section class="slide" data-layout="keynotes">
  <p class="pr-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="pr-swatch"><div class="pr-sw-name">{{色名，2-3 字}}</div><div class="pr-sw-py">{{拼音大写，≤8 字符}}</div><p class="pr-sw-mean">{{释义：一句意象 + 一句用途，22-44 字}}</p><div class="pr-sw-hex">{{hex 色值，如 #F5F2E9}}</div></div>
    <div class="pr-swatch"><div class="pr-sw-name">{{色名，2-3 字}}</div><div class="pr-sw-py">{{拼音大写，≤8 字符}}</div><p class="pr-sw-mean">{{释义：一句意象 + 一句用途，22-44 字}}</p><div class="pr-sw-hex">{{hex 色值，如 #F091A0}}</div></div>
    <div class="pr-swatch"><div class="pr-sw-name">{{色名，2-3 字}}</div><div class="pr-sw-py">{{拼音大写，≤8 字符}}</div><p class="pr-sw-mean">{{释义：一句意象 + 一句用途，22-44 字}}</p><div class="pr-sw-hex">{{hex 色值，如 #B23A48}}</div></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pr-kicker, h2, mt-m, grid, g3, mt-l, pr-swatch, pr-sw-name, pr-sw-py, pr-sw-mean, pr-sw-hex, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：pr-step=3

用途：左边把一件事讲透（lede + 补充 + 胶囊），右象牙卡装三行体验步骤。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个胶囊；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="pr-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#666666">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="pr-pill">{{要点 1，≤8 字}}</span>
        <span class="pr-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="pr-card">
      <div class="pr-step"><span class="pr-n">一</span><p class="pr-step-t">{{一步，12-26 字}}</p></div>
      <div class="pr-step"><span class="pr-n">二</span><p class="pr-step-t">{{一步，12-26 字}}</p></div>
      <div class="pr-step"><span class="pr-n">三</span><p class="pr-step-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pr-kicker, h2, mt-m, grid, g2, mt-l, lede, row, pr-pill, pr-card, pr-step, pr-n, pr-step-t, deck-footer, slide-number, notes

---

## metrics（胭脂数字）
指纹：chart
数量：pr-stat=3

用途：三个关键数据。胭脂大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="pr-kicker">{{数据语境，如 本展规模}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="pr-stat"><div class="pr-stat-v">{{数值 ≤6 字符}}<span class="pr-stat-u">{{单位}}</span></div><div class="pr-stat-l">{{指标名，≤8 字}}</div><p class="pr-stat-note">{{口径，14-30 字}}</p></div>
    <div class="pr-stat"><div class="pr-stat-v">{{数值 ≤6 字符}}<span class="pr-stat-u">{{单位}}</span></div><div class="pr-stat-l">{{指标名，≤8 字}}</div><p class="pr-stat-note">{{口径，14-30 字}}</p></div>
    <div class="pr-stat"><div class="pr-stat-v">{{数值 ≤6 字符}}<span class="pr-stat-u">{{单位}}</span></div><div class="pr-stat-l">{{指标名，≤8 字}}</div><p class="pr-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#8F887E;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pr-kicker, h2, mt-m, grid, g3, mt-l, pr-stat, pr-stat-v, pr-stat-u, pr-stat-l, pr-stat-note, deck-footer, slide-number, notes

---

## quote（瓷语引文）
指纹：quote

用途：整页一句引文。宋体大字 + 出处 + 两个支撑胶囊，右侧留白处可立竖排短句。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；胶囊各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="pr-kicker">{{语境，如 布展手记}}</p>
  <p class="pr-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="pr-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="pr-pill pr-pill-red">{{支撑点 1，≤8 字}}</span>
    <span class="pr-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="pr-vert" style="position:absolute;right:120px;top:50%;transform:translateY(-50%)">{{竖排短句 ≤7 字}}<span class="pr-vert-accent">{{一字或两字}}</span></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pr-kicker, pr-quote, mt-l, pr-src, mt-m, row, pr-pill, pr-pill-red, pr-vert, pr-vert-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 大字章节名 + 一个过渡问题 + 两个看点胶囊。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；胶囊各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="pr-kicker">{{进度，如 第二折 · 三系瓷器}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="pr-pill pr-pill-red">{{看点 1，≤8 字}}</span>
    <span class="pr-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, pr-kicker, h1, mt-m, lede, mt-l, row, pr-pill, pr-pill-red, deck-footer, slide-number, notes

---

## moments（导览时间线）
指纹：chart
数量：pr-tl-item=4

用途：3-4 个站点的横向时间线：瓷圈站号 + 时间点 + 一句看点。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="pr-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="pr-tl mt-l" style="margin-top:52px">
    <div class="pr-tl-item"><div class="pr-tl-dot">壹</div><div class="pr-tl-t">{{时间点，≤10 字符}}</div><p class="pr-tl-d">{{事件，12-26 字}}</p></div>
    <div class="pr-tl-item"><div class="pr-tl-dot">贰</div><div class="pr-tl-t">{{时间点，≤10 字符}}</div><p class="pr-tl-d">{{事件，12-26 字}}</p></div>
    <div class="pr-tl-item"><div class="pr-tl-dot">叁</div><div class="pr-tl-t">{{时间点，≤10 字符}}</div><p class="pr-tl-d">{{事件，12-26 字}}</p></div>
    <div class="pr-tl-item"><div class="pr-tl-dot">肆</div><div class="pr-tl-t">{{时间点，≤10 字符}}</div><p class="pr-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#8F887E;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, pr-kicker, h2, mt-m, pr-tl, mt-l, pr-tl-item, pr-tl-dot, pr-tl-t, pr-tl-d, deck-footer, slide-number, notes

---

## closing（收尾相邀）
指纹：hero

用途：收尾页。宋体大字 + 一句行动提醒 + 胭脂实底按钮 + 胶囊与小印，如瓷底落款。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；胶囊 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="pr-kicker">{{提醒语境，如 导览预约}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="pr-btn">{{按钮文案，≤8 字}}</span>
    <span class="pr-pill">{{次级信息，≤10 字}}</span>
    <div class="pr-seal">{{小印 2 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, pr-kicker, h1, mt-m, lede, mt-l, row, pr-btn, pr-pill, pr-seal, deck-footer, slide-number, notes
