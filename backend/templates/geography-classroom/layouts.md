# 地理课堂·蓝图 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `ge-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-geography-classroom` 作用域生效，骨架里已写全，照抄结构即可。
> 每页的经纬网与页底等高线、指北针由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（图册即课堂）**：结构色是海军蓝（标题、编号块、按钮），青绿负责讲解；
> 重点色只有两处——珊瑚红给关键数据（ge-stat-v / ge-tl-t / ge-pill-accent），
> 明黄给 ge-tag 标签（一页至多两枚）。禁大面积撞色、禁硬阴影；
> h1/h2 每页最多一组；一页最多两组内容块 + 页脚，留出图面的呼吸感。

---

## cover（课题封面）
指纹：hero

用途：开场页。课程题签 + 课题大字 + 课线 + 本课要解决的问题，配课时标签与授课信息。
适用 role：cover。
内容约束：课题名 ≤10 字；lede 24-48 字；题签 ≤14 字；课时标签 ≤8 字。

```html
<section class="slide full" data-layout="cover">
  <p class="ge-kicker">{{课程信息，≤14 字，如 七年级地理 · 第三章}}</p>
  <h1 class="h1 mt-m">{{课题名，≤10 字，可 <br> 分行}}</h1>
  <div class="ge-line mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{本课要解决的问题，24-48 字}}</p>
  <div class="row mt-l" style="gap:24px">
    <span class="ge-tag">{{课时标签，≤8 字}}</span>
    <span class="ge-pill">{{授课信息，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{课程 · 授课人}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, ge-kicker, h1, mt-m, ge-line, mt-s, lede, mt-l, row, ge-tag, ge-pill, deck-footer, slide-number, notes

---

## contents（环节目录）
指纹：table
数量：ge-item=4

用途：议程页。一块图册大卡里放 4 行教学环节：编号块 + 环节名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；环节名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="ge-kicker">{{引导语，如 本课环节}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="ge-card mt-l" style="margin-top:44px">
    <div class="ge-item"><span class="ge-n">01</span><span class="ge-t">{{环节名，≤8 字}}</span><span class="ge-d">{{说明，14-26 字}}</span></div>
    <div class="ge-item"><span class="ge-n">02</span><span class="ge-t">{{环节名，≤8 字}}</span><span class="ge-d">{{说明，14-26 字}}</span></div>
    <div class="ge-item"><span class="ge-n">03</span><span class="ge-t">{{环节名，≤8 字}}</span><span class="ge-d">{{说明，14-26 字}}</span></div>
    <div class="ge-item"><span class="ge-n">04</span><span class="ge-t">{{环节名，≤8 字}}</span><span class="ge-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{课程 · 授课人}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ge-kicker, h2, mt-m, ge-card, mt-l, ge-item, ge-n, ge-t, ge-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：ge-card=3

用途：恰好三张图册卡。每张：珊瑚红题签 + 小标题 + 两句讲解。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="ge-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="ge-card"><span class="ge-pill ge-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#43535D">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="ge-card"><span class="ge-pill ge-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#43535D">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="ge-card"><span class="ge-pill ge-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#43535D">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{课程 · 授课人}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ge-kicker, h2, mt-m, grid, g3, mt-l, ge-card, ge-pill, ge-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左讲右证）
指纹：split
数量：ge-step=3

用途：左边把一件事讲透（lede + 判断标准 + 药丸），右边图册卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="ge-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#43535D">{{补充：判断标准或注意事项，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="ge-pill">{{要点 1，≤8 字}}</span>
        <span class="ge-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="ge-card">
      <div class="ge-step"><span class="ge-n">一</span><p class="ge-mini-t">{{一步，12-26 字}}</p></div>
      <div class="ge-step"><span class="ge-n">二</span><p class="ge-mini-t">{{一步，12-26 字}}</p></div>
      <div class="ge-step"><span class="ge-n">三</span><p class="ge-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{课程 · 授课人}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ge-kicker, h2, mt-m, grid, g2, mt-l, lede, row, ge-pill, ge-card, ge-step, ge-n, ge-mini-t, deck-footer, slide-number, notes

---

## metrics（考点数字）
指纹：chart
数量：ge-stat=3

用途：三个关键数据。珊瑚红大数字是唯一主视觉，口径写进块内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="ge-kicker">{{数据语境，如 考点数据}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="ge-stat"><div class="ge-stat-v">{{数值 ≤6 字符}}<span class="ge-stat-u">{{单位}}</span></div><div class="ge-stat-l">{{指标名，≤8 字}}</div><p class="ge-stat-note">{{口径，14-30 字}}</p></div>
    <div class="ge-stat"><div class="ge-stat-v">{{数值 ≤6 字符}}<span class="ge-stat-u">{{单位}}</span></div><div class="ge-stat-l">{{指标名，≤8 字}}</div><p class="ge-stat-note">{{口径，14-30 字}}</p></div>
    <div class="ge-stat"><div class="ge-stat-v">{{数值 ≤6 字符}}<span class="ge-stat-u">{{单位}}</span></div><div class="ge-stat-l">{{指标名，≤8 字}}</div><p class="ge-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#7C8B94;letter-spacing:.06em">来源：{{教材页码或权威机构与统计口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{课程 · 授课人}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ge-kicker, h2, mt-m, grid, g3, mt-l, ge-stat, ge-stat-v, ge-stat-u, ge-stat-l, ge-stat-note, deck-footer, slide-number, notes

---

## quote（谚语引读）
指纹：quote

用途：整页一句引文。宋体大字 + 出处 + 两个支撑药丸，右下留白交给等高线衬底。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="ge-kicker">{{语境，如 课堂引读}}</p>
  <p class="ge-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="ge-src mt-m">—— {{出处：谚语来源或教材出处，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="ge-pill ge-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="ge-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{课程 · 授课人}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ge-kicker, ge-quote, mt-l, ge-src, mt-m, row, ge-pill, ge-pill-accent, deck-footer, slide-number, notes

---

## divider（环节幕）
指纹：hero

用途：环节过渡。进度题签 + 课题大字环节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：环节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="ge-kicker">{{进度，如 第二环节 · 季风成因}}</p>
  <h1 class="h1 mt-m">{{环节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一环节回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="ge-pill ge-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="ge-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{课程 · 授课人}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, ge-kicker, h1, mt-m, lede, mt-l, row, ge-pill, ge-pill-accent, deck-footer, slide-number, notes

---

## moments（课时流程）
指纹：chart
数量：ge-tl-item=4

用途：3-4 个节点的横向时间线：编号块 + 时间点 + 一句安排。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="ge-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="ge-tl mt-l" style="margin-top:52px">
    <div class="ge-tl-item"><div class="ge-tl-dot">01</div><div class="ge-tl-t">{{时间点，≤10 字符}}</div><p class="ge-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ge-tl-item"><div class="ge-tl-dot">02</div><div class="ge-tl-t">{{时间点，≤10 字符}}</div><p class="ge-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ge-tl-item"><div class="ge-tl-dot">03</div><div class="ge-tl-t">{{时间点，≤10 字符}}</div><p class="ge-tl-d">{{事件，12-26 字}}</p></div>
    <div class="ge-tl-item"><div class="ge-tl-dot">04</div><div class="ge-tl-t">{{时间点，≤10 字符}}</div><p class="ge-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#7C8B94;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{课程 · 授课人}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ge-kicker, h2, mt-m, ge-tl, mt-l, ge-tl-item, ge-tl-dot, ge-tl-t, ge-tl-d, deck-footer, slide-number, notes

---

## closing（课堂收束）
指纹：hero

用途：收尾页。课题大字 + 一句课后任务 + 海军蓝按钮 + 标签药丸，如下课前的黑板小结。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="ge-kicker">{{提醒语境，如 课后任务}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束任务，20-44 字}}</p>
  <div class="row mt-l" style="gap:24px">
    <span class="ge-btn">{{按钮文案，≤8 字}}</span>
    <span class="ge-pill">{{次级信息，≤10 字}}</span>
    <span class="ge-tag">{{作业标签，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{课程 · 授课人}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, ge-kicker, h1, mt-m, lede, mt-l, row, ge-btn, ge-pill, ge-tag, deck-footer, slide-number, notes
