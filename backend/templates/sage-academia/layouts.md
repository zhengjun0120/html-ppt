# 青竹 · 学术风 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `sa-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-sage-academia` 作用域生效，骨架里已写全，照抄结构即可。
> 右上鼠尾草大圆与左下圆环、虚线由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（绿意有分工）**：鼠尾草绿只允许出现在四处——绿环编号（sa-no / sa-tl-dot）、
> 关键数据（sa-stat-v / sa-tl-t）、强调药丸与题签短线（sa-pill-accent / sa-kicker-line）、
> 收尾按钮（sa-btn）。青灰蓝只做题签小字（sa-kicker）。标题一律近黑无衬线，禁绿标题；
> 卡是白面细边轻阴影，禁厚阴影、禁深色底、禁高饱和撞色。

---

## cover（汇报封面）
指纹：hero

用途：开场页。青灰题签 + 鼠尾草短线 + 近黑大标题 + 一句副题，药丸标人与场合。
适用 role：cover。
内容约束：主标题 ≤12 字；lede 24-56 字；题签 ≤20 字；药丸各 ≤14 字。

```html
<section class="slide full" data-layout="cover">
  <div><p class="sa-kicker">{{题签，如 READ &amp; SOCIETY · 期末报告}}</p><div class="sa-kicker-line"></div></div>
  <h1 class="h1 mt-m">{{主标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:52ch">{{一句话定位，24-56 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="sa-pill">{{汇报人 · 院系，≤14 字}}</span>
    <span class="sa-pill sa-pill-accent">{{时间与场合，≤14 字}}</span>
  </div>
  <div class="deck-footer"><span>{{院系 · 课程}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, sa-kicker, sa-kicker-line, h1, mt-m, lede, row, mt-l, sa-pill, sa-pill-accent, deck-footer, slide-number, notes

---

## contents（报告目录）
指纹：table
数量：sa-item=4

用途：目录页。一块白卡里放 4 行条目：绿环编号 + 条目名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；条目名 ≤8 字；说明 12-24 字。

```html
<section class="slide" data-layout="contents">
  <p class="sa-kicker">{{引导语，如 报告目录}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="sa-card mt-l" style="margin-top:44px">
    <div class="sa-item"><span class="sa-no">壹</span><span class="sa-t">{{条目名，≤8 字}}</span><span class="sa-d">{{说明，12-24 字}}</span></div>
    <div class="sa-item"><span class="sa-no">贰</span><span class="sa-t">{{条目名，≤8 字}}</span><span class="sa-d">{{说明，12-24 字}}</span></div>
    <div class="sa-item"><span class="sa-no">叁</span><span class="sa-t">{{条目名，≤8 字}}</span><span class="sa-d">{{说明，12-24 字}}</span></div>
    <div class="sa-item"><span class="sa-no">肆</span><span class="sa-t">{{条目名，≤8 字}}</span><span class="sa-d">{{说明，12-24 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{院系 · 课程}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sa-kicker, h2, mt-m, sa-card, mt-l, sa-item, sa-no, sa-t, sa-d, deck-footer, slide-number, notes

---

## keynotes（三卡观察）
指纹：cards
数量：sa-card=3

用途：恰好三张白卡。每张：绿字药丸题签 + 小标题 + 两句说明，三个并列观察。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="sa-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="sa-card"><span class="sa-pill sa-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#5F5F5F">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="sa-card"><span class="sa-pill sa-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#5F5F5F">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="sa-card"><span class="sa-pill sa-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#5F5F5F">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{院系 · 课程}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sa-kicker, h2, mt-m, grid, g3, mt-l, sa-card, sa-pill, sa-pill-accent, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左叙右列）
指纹：split
数量：sa-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边白卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="sa-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#5F5F5F">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="sa-pill">{{要点 1，≤8 字}}</span>
        <span class="sa-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="sa-card">
      <div class="sa-step"><span class="sa-no">一</span><p class="sa-mini-t">{{一步，12-26 字}}</p></div>
      <div class="sa-step"><span class="sa-no">二</span><p class="sa-mini-t">{{一步，12-26 字}}</p></div>
      <div class="sa-step"><span class="sa-no">三</span><p class="sa-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{院系 · 课程}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sa-kicker, h2, mt-m, grid, g2, mt-l, lede, row, sa-pill, sa-card, sa-step, sa-no, sa-mini-t, deck-footer, slide-number, notes

---

## metrics（关键数据）
指纹：chart
数量：sa-stat=3

用途：三个关键数据。鼠尾草绿大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 12-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="sa-kicker">{{数据语境，如 调查数据}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="sa-stat"><div class="sa-stat-v">{{数值 ≤6 字符}}<span class="sa-stat-u">{{单位}}</span></div><div class="sa-stat-l">{{指标名，≤8 字}}</div><p class="sa-stat-note">{{口径，12-30 字}}</p></div>
    <div class="sa-stat"><div class="sa-stat-v">{{数值 ≤6 字符}}<span class="sa-stat-u">{{单位}}</span></div><div class="sa-stat-l">{{指标名，≤8 字}}</div><p class="sa-stat-note">{{口径，12-30 字}}</p></div>
    <div class="sa-stat"><div class="sa-stat-v">{{数值 ≤6 字符}}<span class="sa-stat-u">{{单位}}</span></div><div class="sa-stat-l">{{指标名，≤8 字}}</div><p class="sa-stat-note">{{口径，12-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#727273;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{院系 · 课程}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sa-kicker, h2, mt-m, grid, g3, mt-l, sa-stat, sa-stat-v, sa-stat-u, sa-stat-l, sa-stat-note, deck-footer, slide-number, notes

---

## quote（题记引文）
指纹：quote

用途：整页一句引文。衬线大字 + 出处 + 两个支撑药丸，人文气质的压轴一笔。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="sa-kicker">{{语境，如 题记}}</p>
  <p class="sa-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="sa-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="sa-pill sa-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="sa-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{院系 · 课程}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sa-kicker, sa-quote, mt-l, sa-src, mt-m, row, sa-pill, sa-pill-accent, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。进度题签 + 近黑大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤12 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="sa-kicker">{{进度，如 第二章 · 数据里的书架}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="sa-pill sa-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="sa-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{院系 · 课程}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, sa-kicker, h1, mt-m, lede, row, mt-l, sa-pill, sa-pill-accent, deck-footer, slide-number, notes

---

## moments（共读时间线）
指纹：chart
数量：sa-tl-item=4

用途：3-4 个节点的横向时间线：绿圆点编号 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="sa-kicker">{{引导语，如 一学期共读}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="sa-tl mt-l" style="margin-top:52px">
    <div class="sa-tl-item"><div class="sa-tl-dot">壹</div><div class="sa-tl-t">{{时间点，≤10 字符}}</div><p class="sa-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sa-tl-item"><div class="sa-tl-dot">贰</div><div class="sa-tl-t">{{时间点，≤10 字符}}</div><p class="sa-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sa-tl-item"><div class="sa-tl-dot">叁</div><div class="sa-tl-t">{{时间点，≤10 字符}}</div><p class="sa-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sa-tl-item"><div class="sa-tl-dot">肆</div><div class="sa-tl-t">{{时间点，≤10 字符}}</div><p class="sa-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#727273;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{院系 · 课程}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sa-kicker, h2, mt-m, sa-tl, mt-l, sa-tl-item, sa-tl-dot, sa-tl-t, sa-tl-d, deck-footer, slide-number, notes

---

## closing（行动收尾）
指纹：hero

用途：收尾页。近黑大字 + 一句行动提醒 + 绿底圆角按钮 + 药丸。
适用 role：thanks / cta / content。
内容约束：标题 ≤12 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤14 字。

```html
<section class="slide full" data-layout="closing">
  <p class="sa-kicker">{{提醒语境，如 行动倡议}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:22px">
    <span class="sa-btn">{{按钮文案，≤8 字}}</span>
    <span class="sa-pill">{{次级信息，≤14 字}}</span>
  </div>
  <div class="deck-footer"><span>{{院系 · 课程}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, sa-kicker, h1, mt-m, lede, row, mt-l, sa-btn, sa-pill, deck-footer, slide-number, notes
