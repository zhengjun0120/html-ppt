# 赛博霓虹 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `cn-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-cyberpunk-neon` 作用域生效，骨架里已写全，照抄结构即可。
> 每页故障网格、霓虹横线与扫描带由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（霓虹是身份，不是正文）**：双霓虹分工——粉 #FF2E88 管结构与描边
> （cn-card / cn-pill / cn-stat 顶线 / cn-btn），青 #00F0FF 管数据与高亮
> （cn-stat-v / cn-tl-t / cn-n / cn-kicker）。正文永远冷白 #F8FAFC 或灰白 #CBD5E1，
> 霓虹色禁用于正文段落；发光（text-shadow）只给标题、数字与线条，正文段落禁加。

---

## cover（霓虹封面）
指纹：hero

用途：开场页。等宽题签 + 辉光大标题 + 一句定位，粉青描边标签点题。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 28-56 字；题签 ≤16 字；标签各 ≤10 字。

```html
<section class="slide full" data-layout="cover">
  <p class="cn-kicker">{{题签，≤16 字，如 NIGHT DRIVE · 摄影分享}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句话定位，28-56 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="cn-pill">{{标签 1，≤10 字}}</span>
    <span class="cn-pill cn-pill-cyan">{{标签 2，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, cn-kicker, h1, mt-m, lede, mt-l, row, cn-pill, cn-pill-cyan, deck-footer, slide-number, notes

---

## contents（目录屏）
指纹：table
数量：cn-item=4

用途：议程页。一块暗卡里放 4 行篇目：青色描边编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="cn-kicker">{{引导语，如 分享目录}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="cn-card mt-l" style="margin-top:44px">
    <div class="cn-item"><span class="cn-n">01</span><span class="cn-t">{{篇名，≤8 字}}</span><span class="cn-d">{{说明，14-26 字}}</span></div>
    <div class="cn-item"><span class="cn-n">02</span><span class="cn-t">{{篇名，≤8 字}}</span><span class="cn-d">{{说明，14-26 字}}</span></div>
    <div class="cn-item"><span class="cn-n">03</span><span class="cn-t">{{篇名，≤8 字}}</span><span class="cn-d">{{说明，14-26 字}}</span></div>
    <div class="cn-item"><span class="cn-n">04</span><span class="cn-t">{{篇名，≤8 字}}</span><span class="cn-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cn-kicker, h2, mt-m, cn-card, mt-l, cn-item, cn-n, cn-t, cn-d, deck-footer, slide-number, notes

---

## keynotes（三卡点位）
指纹：cards
数量：cn-card=3

用途：恰好三张暗卡。每张：粉描边题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="cn-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:36px;margin-top:44px">
    <div class="cn-card"><span class="cn-pill">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#CBD5E1">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="cn-card"><span class="cn-pill">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#CBD5E1">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="cn-card"><span class="cn-pill">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#CBD5E1">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cn-kicker, h2, mt-m, grid, g3, mt-l, cn-card, cn-pill, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：cn-step=3

用途：左边把一件事讲透（lede + 补充 + 标签），右边暗卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个标签；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="cn-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#CBD5E1">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="cn-pill">{{要点 1，≤10 字}}</span>
        <span class="cn-pill cn-pill-cyan">{{要点 2，≤10 字}}</span>
      </div>
    </div>
    <div class="cn-card">
      <div class="cn-step"><span class="cn-n">1</span><p class="cn-mini-t">{{一步，12-26 字}}</p></div>
      <div class="cn-step"><span class="cn-n">2</span><p class="cn-mini-t">{{一步，12-26 字}}</p></div>
      <div class="cn-step"><span class="cn-n">3</span><p class="cn-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cn-kicker, h2, mt-m, grid, g2, mt-l, lede, row, cn-pill, cn-pill-cyan, cn-card, cn-step, cn-n, cn-mini-t, deck-footer, slide-number, notes

---

## metrics（霓虹数字）
指纹：chart
数量：cn-stat=3

用途：三个关键数据。青色发光大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="cn-kicker">{{数据语境，如 三年夜行 · 复盘}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:48px;margin-top:48px">
    <div class="cn-stat"><div class="cn-stat-v">{{数值 ≤6 字符}}<span class="cn-stat-u">{{单位}}</span></div><div class="cn-stat-l">{{指标名，≤8 字}}</div><p class="cn-stat-note">{{口径，14-30 字}}</p></div>
    <div class="cn-stat"><div class="cn-stat-v">{{数值 ≤6 字符}}<span class="cn-stat-u">{{单位}}</span></div><div class="cn-stat-l">{{指标名，≤8 字}}</div><p class="cn-stat-note">{{口径，14-30 字}}</p></div>
    <div class="cn-stat"><div class="cn-stat-v">{{数值 ≤6 字符}}<span class="cn-stat-u">{{单位}}</span></div><div class="cn-stat-l">{{指标名，≤8 字}}</div><p class="cn-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:16px;color:#7C86A0;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cn-kicker, h2, mt-m, grid, g3, mt-l, cn-stat, cn-stat-v, cn-stat-u, cn-stat-l, cn-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（夜行引文）
指纹：quote

用途：整页一句引文。辉光大字 + 出处 + 两个支撑标签。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；标签各 ≤8 字。

```html
<section class="slide full" data-layout="quote">
  <p class="cn-kicker">{{语境，如 画册扉页}}</p>
  <p class="cn-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="cn-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="cn-pill">{{支撑点 1，≤8 字}}</span>
    <span class="cn-pill cn-pill-cyan">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, cn-kicker, cn-quote, mt-l, cn-src, mt-m, row, cn-pill, cn-pill-cyan, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 辉光大字章节名 + 一个过渡问题 + 两个看点标签。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；标签各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="cn-kicker">{{进度，如 第二段 · 器材与参数}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:46ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="cn-pill">{{看点 1，≤8 字}}</span>
    <span class="cn-pill cn-pill-cyan">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, cn-kicker, h1, mt-m, lede, mt-l, row, cn-pill, cn-pill-cyan, deck-footer, slide-number, notes

---

## moments（巡展时间线）
指纹：chart
数量：cn-tl-item=4

用途：3-4 个节点的横向时间线：发光方块节点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="cn-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="cn-tl mt-l" style="margin-top:52px">
    <div class="cn-tl-item"><div class="cn-tl-dot"></div><div class="cn-tl-t">{{时间点，≤10 字符}}</div><p class="cn-tl-d">{{事件，12-26 字}}</p></div>
    <div class="cn-tl-item"><div class="cn-tl-dot"></div><div class="cn-tl-t">{{时间点，≤10 字符}}</div><p class="cn-tl-d">{{事件，12-26 字}}</p></div>
    <div class="cn-tl-item"><div class="cn-tl-dot"></div><div class="cn-tl-t">{{时间点，≤10 字符}}</div><p class="cn-tl-d">{{事件，12-26 字}}</p></div>
    <div class="cn-tl-item"><div class="cn-tl-dot"></div><div class="cn-tl-t">{{时间点，≤10 字符}}</div><p class="cn-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:16px;color:#7C86A0;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, cn-kicker, h2, mt-m, cn-tl, mt-l, cn-tl-item, cn-tl-dot, cn-tl-t, cn-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（收尾入场）
指纹：hero

用途：收尾页。辉光大字 + 一句行动提醒 + 粉底按钮 + 描边标签。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；标签 ≤12 字。

```html
<section class="slide full" data-layout="closing">
  <p class="cn-kicker">{{提醒语境，如 上海首展}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:46ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:20px">
    <span class="cn-btn">{{按钮文案，≤8 字}}</span>
    <span class="cn-pill cn-pill-cyan">{{次级信息，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, cn-kicker, h1, mt-m, lede, mt-l, row, cn-btn, cn-pill, cn-pill-cyan, deck-footer, slide-number, notes
