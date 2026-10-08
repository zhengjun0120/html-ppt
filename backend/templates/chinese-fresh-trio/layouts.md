# 青绿湖蓝 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `ft-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-chinese-fresh-trio` 作用域生效，骨架里已写全，照抄结构即可。
> 每页骨架第一行的三色条（ft-band）照抄不要删；右下几何山峦由模板自动衬底
> （z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（三色即身份）**：翠竹绿、湖蓝、淡米只出现在三色条、圆点编号（ft-n / ft-tl-dot）、
> 实色三联卡（ft-card-*）与数据顶线（ft-stat）；朱红只允许四处——印章（ft-seal）、题签短线
> （ft-kicker 自带）、关键数据（ft-stat-v / ft-tl-t）、强调药丸（ft-pill-accent）。
> 禁硬阴影、禁暗黑背景、禁复杂渐变、禁第四种主色。宋体大字（h1/h2）每页最多一组；
> 实色三联卡只在 keynotes 用，一页至多一组。

---

## cover（三色封面）
指纹：hero

用途：开场页。题签 + 宋体大字标题 + 三色短线 + 一句定位，右侧可立一枚朱砂印。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；印章 4 字内。

```html
<section class="slide full" data-layout="cover">
  <div class="ft-band"><i></i><i></i><i></i></div>
  <p class="ft-kicker">{{题签，≤14 字，如 春季研学 · 第 7 期}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="ft-line mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="ft-seal">{{印章 2×2 字}}</div>
    <span class="ft-pill">{{时间或地点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, ft-band, ft-kicker, h1, mt-m, ft-line, mt-s, lede, mt-l, row, ft-seal, ft-pill, deck-footer, slide-number, notes

---

## contents（课程目录）
指纹：table
数量：ft-item=4

用途：议程页。一块细线卡里放 4 行篇目：三色圆点编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <div class="ft-band"><i></i><i></i><i></i></div>
  <p class="ft-kicker">{{引导语，如 手册目录}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="ft-panel mt-l" style="margin-top:44px">
    <div class="ft-item"><span class="ft-n">01</span><span class="ft-t">{{篇名，≤8 字}}</span><span class="ft-d">{{说明，14-26 字}}</span></div>
    <div class="ft-item"><span class="ft-n">02</span><span class="ft-t">{{篇名，≤8 字}}</span><span class="ft-d">{{说明，14-26 字}}</span></div>
    <div class="ft-item"><span class="ft-n">03</span><span class="ft-t">{{篇名，≤8 字}}</span><span class="ft-d">{{说明，14-26 字}}</span></div>
    <div class="ft-item"><span class="ft-n">04</span><span class="ft-t">{{篇名，≤8 字}}</span><span class="ft-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ft-band, ft-kicker, h2, mt-m, ft-panel, mt-l, ft-item, ft-n, ft-t, ft-d, deck-footer, slide-number, notes

---

## keynotes（三色课程卡）
指纹：cards
数量：ft-card=3

用途：恰好三张实色卡。每张：无衬线编号 + 黑体小标题 + 两句说明，绿、蓝、墨各一。
适用 role：content。
内容约束：恰好 3 卡；编号由骨架给出；线名 ≤6 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <div class="ft-band"><i></i><i></i><i></i></div>
  <p class="ft-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:30px;margin-top:44px">
    <div class="ft-card ft-card-green"><span class="ft-card-num">01</span><h4 class="mt-m">{{线名，≤6 字}}</h4><p class="mt-s">{{说明：面向谁 + 上什么，22-44 字}}</p></div>
    <div class="ft-card ft-card-blue"><span class="ft-card-num">02</span><h4 class="mt-m">{{线名，≤6 字}}</h4><p class="mt-s">{{说明：怎么做 + 带走什么，22-44 字}}</p></div>
    <div class="ft-card ft-card-ink"><span class="ft-card-num">03</span><h4 class="mt-m">{{线名，≤6 字}}</h4><p class="mt-s">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ft-band, ft-kicker, h2, mt-m, grid, g3, mt-l, ft-card, ft-card-green, ft-card-blue, ft-card-ink, ft-card-num, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：ft-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边细线卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <div class="ft-band"><i></i><i></i><i></i></div>
  <p class="ft-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:var(--ink-2)">{{补充：安全或准备事项，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="ft-pill">{{要点 1，≤8 字}}</span>
        <span class="ft-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="ft-panel">
      <div class="ft-step"><span class="ft-n">1</span><span class="ft-d">{{一步，12-26 字}}</span></div>
      <div class="ft-step"><span class="ft-n">2</span><span class="ft-d">{{一步，12-26 字}}</span></div>
      <div class="ft-step"><span class="ft-n">3</span><span class="ft-d">{{一步，12-26 字}}</span></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ft-band, ft-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, ft-pill, ft-panel, ft-step, ft-n, ft-d, deck-footer, slide-number, notes

---

## metrics（三色数据）
指纹：chart
数量：ft-stat=3

用途：三个关键数据。朱红大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <div class="ft-band"><i></i><i></i><i></i></div>
  <p class="ft-kicker">{{数据语境，如 去年春季 · 复盘}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="ft-stat"><div class="ft-stat-v">{{数值 ≤6 字符}}<span class="ft-stat-u">{{单位}}</span></div><div class="ft-stat-l">{{指标名，≤8 字}}</div><p class="ft-stat-note">{{口径，14-30 字}}</p></div>
    <div class="ft-stat"><div class="ft-stat-v">{{数值 ≤6 字符}}<span class="ft-stat-u">{{单位}}</span></div><div class="ft-stat-l">{{指标名，≤8 字}}</div><p class="ft-stat-note">{{口径，14-30 字}}</p></div>
    <div class="ft-stat"><div class="ft-stat-v">{{数值 ≤6 字符}}<span class="ft-stat-u">{{单位}}</span></div><div class="ft-stat-l">{{指标名，≤8 字}}</div><p class="ft-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:var(--ink-3);letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ft-band, ft-kicker, h2, mt-m, grid, g3, mt-l, ft-stat, ft-stat-v, ft-stat-u, ft-stat-l, ft-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（引文页）
指纹：quote

用途：整页一句引文。宋体大字 + 出处 + 两个支撑药丸。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <div class="ft-band"><i></i><i></i><i></i></div>
  <p class="ft-kicker">{{语境，如 主理人说}}</p>
  <p class="ft-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="ft-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="ft-pill ft-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="ft-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ft-band, ft-kicker, ft-quote, mt-l, ft-src, mt-m, row, ft-pill, ft-pill-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <div class="ft-band"><i></i><i></i><i></i></div>
  <p class="ft-kicker">{{进度，如 第一条 · 观鸟线}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="ft-pill ft-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="ft-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, ft-band, ft-kicker, h1, mt-m, lede, mt-l, row, ft-pill, ft-pill-accent, deck-footer, slide-number, notes

---

## moments（流程时间线）
指纹：chart
数量：ft-tl-item=4

用途：恰好 4 个节点的横向时间线：三色圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：恰好 4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <div class="ft-band"><i></i><i></i><i></i></div>
  <p class="ft-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="ft-tl mt-l" style="margin-top:52px">
    <div class="ft-tl-item"><div class="ft-tl-dot">1</div><div class="ft-tl-t">{{时间点，≤10 字符}}</div><p class="ft-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ft-tl-item"><div class="ft-tl-dot">2</div><div class="ft-tl-t">{{时间点，≤10 字符}}</div><p class="ft-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ft-tl-item"><div class="ft-tl-dot">3</div><div class="ft-tl-t">{{时间点，≤10 字符}}</div><p class="ft-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ft-tl-item"><div class="ft-tl-dot">4</div><div class="ft-tl-t">{{时间点，≤10 字符}}</div><p class="ft-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:var(--ink-3);letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ft-band, ft-kicker, h2, mt-m, ft-tl, mt-l, ft-tl-item, ft-tl-dot, ft-tl-t, ft-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（预约收尾）
指纹：hero

用途：收尾页。大字 + 一句行动提醒 + 墨底按钮 + 描边药丸 + 朱砂印。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <div class="ft-band"><i></i><i></i><i></i></div>
  <p class="ft-kicker">{{提醒语境，如 春季预约}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="ft-btn">{{按钮文案，≤8 字}}</span>
    <span class="ft-pill">{{次级信息，≤10 字}}</span>
    <div class="ft-seal">{{印章 2 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 季节}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, ft-band, ft-kicker, h1, mt-m, lede, mt-l, row, ft-btn, ft-pill, ft-seal, deck-footer, slide-number, notes
