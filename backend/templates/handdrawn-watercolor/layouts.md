# 治愈手绘水彩 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `hw-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-handdrawn-watercolor` 作用域生效，骨架里已写全，照抄结构即可。
> 每页右上的水彩晕与页底纸纹波浪由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（绘本留白）**：低饱和是底线——湖蓝 #5FA8C9 与天蓝晕做主视觉，珊瑚 #E4693F
> 只做时间点与按钮，一页至多两处；描边是铅笔软线（1.6px、低透明度），禁硬阴影与清晰色块。
> 大标题（h1/h2）每页最多一组；竖排 hw-vert 只放留白处，一页至多一条。

---

## cover（绘本封面）
指纹：hero

用途：开场页。波浪线小签 + 圆体大标题 + 波浪下划线 + 一句定位，像绘本扉页。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；小签 ≤14 字；药丸各 ≤12 字。

```html
<section class="slide full" data-layout="cover">
  <p class="hw-kicker">{{小签，≤14 字，如 课程手册 · 水彩零基础}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="hw-wave mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="hw-pill">{{时间或形式，≤12 字}}</span>
    <span class="hw-pill-coral hw-pill">{{人数或亮点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, hw-kicker, h1, mt-m, hw-wave, mt-s, lede, mt-l, row, hw-pill, hw-pill-coral, deck-footer, slide-number, notes

---

## contents（课程目录）
指纹：table
数量：hw-item=4

用途：议程页。一张水彩纸卡里放 4 行篇目：晕染圆点编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 12-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="hw-kicker">{{引导语，如 今日卷目}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="hw-card mt-l" style="margin-top:44px">
    <div class="hw-item"><span class="hw-n">壹</span><span class="hw-t">{{篇名，≤8 字}}</span><span class="hw-d">{{说明，12-26 字}}</span></div>
    <div class="hw-item"><span class="hw-n">贰</span><span class="hw-t">{{篇名，≤8 字}}</span><span class="hw-d">{{说明，12-26 字}}</span></div>
    <div class="hw-item"><span class="hw-n">叁</span><span class="hw-t">{{篇名，≤8 字}}</span><span class="hw-d">{{说明，12-26 字}}</span></div>
    <div class="hw-item"><span class="hw-n">肆</span><span class="hw-t">{{篇名，≤8 字}}</span><span class="hw-d">{{说明，12-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, hw-kicker, h2, mt-m, hw-card, mt-l, hw-item, hw-n, hw-t, hw-d, deck-footer, slide-number, notes

---

## keynotes（三张画卡）
指纹：cards
数量：hw-card=3

用途：恰好三张纸卡。每张：药丸题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤6 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="hw-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="hw-card"><span class="hw-pill hw-pill-coral">{{题签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#7A7260">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="hw-card"><span class="hw-pill">{{题签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#7A7260">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="hw-card"><span class="hw-pill hw-pill-coral">{{题签，≤6 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#7A7260">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, hw-kicker, h2, mt-m, grid, g3, mt-l, hw-card, hw-pill, hw-pill-coral, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：hw-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边纸卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="hw-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#7A7260">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="hw-pill">{{要点 1，≤8 字}}</span>
        <span class="hw-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="hw-card">
      <div class="hw-step"><span class="hw-n">一</span><p class="hw-mini-t">{{一步，12-26 字}}</p></div>
      <div class="hw-step"><span class="hw-n">二</span><p class="hw-mini-t">{{一步，12-26 字}}</p></div>
      <div class="hw-step"><span class="hw-n">三</span><p class="hw-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, hw-kicker, h2, mt-m, grid, g2, mt-l, lede, row, hw-pill, hw-card, hw-step, hw-n, hw-mini-t, deck-footer, slide-number, notes

---

## metrics（湖蓝数字）
指纹：chart
数量：hw-stat=3

用途：三个关键数据。湖蓝大数字是唯一主视觉，口径写进纸卡，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="hw-kicker">{{数据语境，如 秋季班 · 课程表}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:40px;margin-top:48px">
    <div class="hw-stat"><div class="hw-stat-v">{{数值 ≤6 字符}}<span class="hw-stat-u">{{单位}}</span></div><div class="hw-stat-l">{{指标名，≤8 字}}</div><p class="hw-stat-note">{{口径，14-30 字}}</p></div>
    <div class="hw-stat"><div class="hw-stat-v">{{数值 ≤6 字符}}<span class="hw-stat-u">{{单位}}</span></div><div class="hw-stat-l">{{指标名，≤8 字}}</div><p class="hw-stat-note">{{口径，14-30 字}}</p></div>
    <div class="hw-stat"><div class="hw-stat-v">{{数值 ≤6 字符}}<span class="hw-stat-u">{{单位}}</span></div><div class="hw-stat-l">{{指标名，≤8 字}}</div><p class="hw-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#A79C84;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, hw-kicker, h2, mt-m, grid, g3, mt-l, hw-stat, hw-stat-v, hw-stat-u, hw-stat-l, hw-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（扉页引文）
指纹：quote

用途：整页一句引文。圆体大字 + 出处 + 两个支撑药丸，右侧留白处立一条竖排短语。
适用 role：quote。
内容约束：引文 10-24 字；出处 ≤22 字且真实；药丸各 ≤10 字；竖排短语 ≤6 字。

```html
<section class="slide" data-layout="quote">
  <p class="hw-kicker">{{语境，如 主讲人开场白}}</p>
  <p class="hw-quote mt-l" style="margin-top:40px">「{{引文，10-24 字}}」</p>
  <p class="hw-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="hw-pill hw-pill-coral">{{支撑点 1，≤10 字}}</span>
    <span class="hw-pill">{{支撑点 2，≤10 字}}</span>
  </div>
  <div class="hw-vert" style="position:absolute;right:130px;top:50%;transform:translateY(-50%)">{{竖排短语 ≤6 字}}</div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, hw-kicker, hw-quote, mt-l, hw-src, mt-m, row, hw-pill, hw-pill-coral, hw-vert, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度小签 + 圆体大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤10 字。

```html
<section class="slide full" data-layout="divider">
  <p class="hw-kicker">{{进度，如 第二讲 · 水与色的游戏}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <div class="hw-wave mt-s"></div>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="hw-pill hw-pill-coral">{{看点 1，≤10 字}}</span>
    <span class="hw-pill">{{看点 2，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, hw-kicker, h1, mt-m, hw-wave, mt-s, lede, mt-l, row, hw-pill, hw-pill-coral, deck-footer, slide-number, notes

---

## moments（进度时间线）
指纹：chart
数量：hw-tl-item=4

用途：3-4 个节点的横向时间线：晕染圆点 + 时间点 + 一句事件，铅笔虚线相连。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="hw-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="hw-tl mt-l" style="margin-top:52px">
    <div class="hw-tl-item"><div class="hw-tl-dot">壹</div><div class="hw-tl-t">{{时间点，≤10 字符}}</div><p class="hw-tl-d">{{事件，12-26 字}}</p></div>
    <div class="hw-tl-item"><div class="hw-tl-dot">贰</div><div class="hw-tl-t">{{时间点，≤10 字符}}</div><p class="hw-tl-d">{{事件，12-26 字}}</p></div>
    <div class="hw-tl-item"><div class="hw-tl-dot">叁</div><div class="hw-tl-t">{{时间点，≤10 字符}}</div><p class="hw-tl-d">{{事件，12-26 字}}</p></div>
    <div class="hw-tl-item"><div class="hw-tl-dot">肆</div><div class="hw-tl-t">{{时间点，≤10 字符}}</div><p class="hw-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#A79C84;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, hw-kicker, h2, mt-m, hw-tl, mt-l, hw-tl-item, hw-tl-dot, hw-tl-t, hw-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（收尾页）
指纹：hero

用途：收尾页。圆体大字 + 一句行动提醒 + 珊瑚按钮 + 药丸，如绘本封底。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤12 字。

```html
<section class="slide full" data-layout="closing">
  <p class="hw-kicker">{{提醒语境，如 秋季班报名}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <div class="hw-wave mt-s"></div>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:20px">
    <span class="hw-btn">{{按钮文案，≤8 字}}</span>
    <span class="hw-pill">{{次级信息，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, hw-kicker, h1, mt-m, hw-wave, mt-s, lede, mt-l, row, hw-btn, hw-pill, deck-footer, slide-number, notes
