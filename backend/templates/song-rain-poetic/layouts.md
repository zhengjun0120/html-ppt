# 宋人生活·听雨 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `sr-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-song-rain-poetic` 作用域生效，骨架里已写全，照抄结构即可。
> 每页右上屋檐雨丝与左下芭蕉叶由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（留白是诗意的基础）**：朱砂只允许出现在三处——印章（sr-seal）、关键数据
> （sr-stat-v / sr-tl-t）、强调药丸与竖排点缀（sr-pill-accent / sr-vert-accent）。
> 暖灯黄只允许出现在茶签卡左竖线（sr-tea）一处。禁硬阴影、禁大面积色块、禁荧光色。
> 宋体大字（h1/h2）每页最多一组；竖排 sr-vert 只放留白处，一页至多一条。

---

## cover（听雨封面）
指纹：hero

用途：开场页。题签 + 宋体大字标题 + 雨线 + 一句定位，右侧可立一枚朱砂印。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；印章 4 字内。

```html
<section class="slide full" data-layout="cover">
  <p class="sr-kicker">{{题签，≤14 字，如 谷雨茶课 · 丙午年春}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="sr-rain-line mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="sr-seal">{{印章 2×2 字}}</div>
    <span class="sr-pill">{{副题或时间地点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, sr-kicker, h1, mt-m, sr-rain-line, mt-s, lede, mt-l, row, sr-seal, sr-pill, deck-footer, slide-number, notes

---

## contents（目录四盏）
指纹：table
数量：sr-item=4

用途：议程页。一块宣纸大卡里放 4 行篇目：圈码编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="sr-kicker">{{引导语，如 今日课表}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="sr-card mt-l" style="margin-top:44px">
    <div class="sr-item"><span class="sr-n">壹</span><span class="sr-t">{{篇名，≤8 字}}</span><span class="sr-d">{{说明，14-26 字}}</span></div>
    <div class="sr-item"><span class="sr-n">贰</span><span class="sr-t">{{篇名，≤8 字}}</span><span class="sr-d">{{说明，14-26 字}}</span></div>
    <div class="sr-item"><span class="sr-n">叁</span><span class="sr-t">{{篇名，≤8 字}}</span><span class="sr-d">{{说明，14-26 字}}</span></div>
    <div class="sr-item"><span class="sr-n">肆</span><span class="sr-t">{{篇名，≤8 字}}</span><span class="sr-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sr-kicker, h2, mt-m, sr-card, mt-l, sr-item, sr-n, sr-t, sr-d, deck-footer, slide-number, notes

---

## keynotes（三盏要点）
指纹：cards
数量：sr-card=3

用途：恰好三张宣纸卡。每张：朱砂题签 + 宋体小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="sr-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="sr-card"><span class="sr-pill sr-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#5A6E70">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="sr-card"><span class="sr-pill sr-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#5A6E70">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="sr-card"><span class="sr-pill sr-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#5A6E70">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sr-kicker, h2, mt-m, grid, g3, mt-l, sr-card, sr-pill, sr-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右签）
指纹：split
数量：sr-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边茶签卡（暖黄左竖线）装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="sr-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#5A6E70">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="sr-pill">{{要点 1，≤8 字}}</span>
        <span class="sr-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="sr-tea">
      <div class="sr-step"><span class="sr-n">一</span><p class="sr-mini-t">{{一步，12-26 字}}</p></div>
      <div class="sr-step"><span class="sr-n">二</span><p class="sr-mini-t">{{一步，12-26 字}}</p></div>
      <div class="sr-step"><span class="sr-n">三</span><p class="sr-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sr-kicker, h2, mt-m, grid, g2, mt-l, lede, row, sr-pill, sr-tea, sr-step, sr-n, sr-mini-t, deck-footer, slide-number, notes

---

## metrics（听雨数字）
指纹：chart
数量：sr-stat=3

用途：三个关键数据。朱砂大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="sr-kicker">{{数据语境，如 去年谷雨课 · 复盘}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="sr-stat"><div class="sr-stat-v">{{数值 ≤6 字符}}<span class="sr-stat-u">{{单位}}</span></div><div class="sr-stat-l">{{指标名，≤8 字}}</div><p class="sr-stat-note">{{口径，14-30 字}}</p></div>
    <div class="sr-stat"><div class="sr-stat-v">{{数值 ≤6 字符}}<span class="sr-stat-u">{{单位}}</span></div><div class="sr-stat-l">{{指标名，≤8 字}}</div><p class="sr-stat-note">{{口径，14-30 字}}</p></div>
    <div class="sr-stat"><div class="sr-stat-v">{{数值 ≤6 字符}}<span class="sr-stat-u">{{单位}}</span></div><div class="sr-stat-l">{{指标名，≤8 字}}</div><p class="sr-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#8FA3A4;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sr-kicker, h2, mt-m, grid, g3, mt-l, sr-stat, sr-stat-v, sr-stat-u, sr-stat-l, sr-stat-note, deck-footer, slide-number, notes

---

## quote（竖诗引文）
指纹：quote

用途：整页一句引文。宋体大字 + 出处 + 两个支撑药丸，右侧留白处立竖排诗句。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="sr-kicker">{{语境，如 宋人题跋}}</p>
  <p class="sr-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="sr-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="sr-pill sr-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="sr-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="sr-vert" style="position:absolute;right:120px;top:50%;transform:translateY(-50%)">{{竖排诗句 ≤7 字}}<span class="sr-vert-accent">{{一字或两字}}</span></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sr-kicker, sr-quote, mt-l, sr-src, mt-m, row, sr-pill, sr-pill-accent, sr-vert, sr-vert-accent, deck-footer, slide-number, notes

---

## divider（章节雨幕）
指纹：hero

用途：章节过渡。进度题签 + 宋体大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="sr-kicker">{{进度，如 卷二 · 雨声与茶汤}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="sr-pill sr-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="sr-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, sr-kicker, h1, mt-m, lede, mt-l, row, sr-pill, sr-pill-accent, deck-footer, slide-number, notes

---

## moments（流程雨声）
指纹：chart
数量：sr-tl-item=4

用途：3-4 个节点的横向时间线：圈码圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="sr-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="sr-tl mt-l" style="margin-top:52px">
    <div class="sr-tl-item"><div class="sr-tl-dot">壹</div><div class="sr-tl-t">{{时间点，≤10 字符}}</div><p class="sr-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sr-tl-item"><div class="sr-tl-dot">贰</div><div class="sr-tl-t">{{时间点，≤10 字符}}</div><p class="sr-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sr-tl-item"><div class="sr-tl-dot">叁</div><div class="sr-tl-t">{{时间点，≤10 字符}}</div><p class="sr-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sr-tl-item"><div class="sr-tl-dot">肆</div><div class="sr-tl-t">{{时间点，≤10 字符}}</div><p class="sr-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#8FA3A4;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sr-kicker, h2, mt-m, sr-tl, mt-l, sr-tl-item, sr-tl-dot, sr-tl-t, sr-tl-d, deck-footer, slide-number, notes

---

## closing（收尾落款）
指纹：hero

用途：收尾页。宋体大字 + 一句行动提醒 + 墨框按钮 + 描边药丸，如听雨后的落款。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="sr-kicker">{{提醒语境，如 茶课报名}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="sr-btn">{{按钮文案，≤8 字}}</span>
    <span class="sr-pill">{{次级信息，≤10 字}}</span>
    <div class="sr-seal">{{印章 2 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, sr-kicker, h1, mt-m, lede, mt-l, row, sr-btn, sr-pill, sr-seal, deck-footer, slide-number, notes
