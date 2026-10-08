# 彩虹渐变 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `rg-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-rainbow-gradient` 作用域生效，骨架里已写全，照抄结构即可。
> 每页顶部七彩丝带与页底彩虹浪由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（庆祝但不混乱）**：正文一律海军蓝/钢灰墨色，长文字优先放进白色卡
> （rg-card / rg-band），不许把正文直接铺在彩底上。七彩只做条、环、点：
> rg-stripe / rg-kicker 丝带 / rg-n 与 rg-tl-dot 彩虹环 / rg-chip 圆点；大数字只许
> 玫红→紫双色渐变（rg-stat-v 已内建），不许给正文上七彩渐变。禁深色背景、禁
> 七彩大色块直衬文字。

---

## cover（彩虹开场）
指纹：hero

用途：开场页。题签胶囊 + 大字标题 + 彩虹条 + 一句定位，三颗彩点芯片与胶囊并列。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 28-48 字；题签 ≤14 字；胶囊 ≤12 字。

```html
<section class="slide full" data-layout="cover">
  <p class="rg-kicker">{{题签，≤14 字，如 活动策划案 · 2026 秋}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="rg-stripe mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，28-48 字}}</p>
  <div class="row mt-l" style="gap:18px">
    <span class="rg-chip" style="--c:#F43F5E"></span>
    <span class="rg-chip" style="--c:#EAB308"></span>
    <span class="rg-chip" style="--c:#8B5CF6"></span>
    <span class="rg-pill">{{时间地点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, rg-kicker, h1, mt-m, rg-stripe, mt-s, lede, mt-l, row, rg-chip, rg-pill, deck-footer, slide-number, notes

---

## contents（庆典目录）
指纹：table
数量：rg-item=4

用途：议程页。一块白色大卡里放 4 行篇目：彩虹环编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 12-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="rg-kicker">{{引导语，如 方案目录}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="rg-band mt-l" style="margin-top:44px">
    <div class="rg-item"><span class="rg-n">01</span><span class="rg-t">{{篇名，≤8 字}}</span><span class="rg-d">{{说明，12-26 字}}</span></div>
    <div class="rg-item"><span class="rg-n">02</span><span class="rg-t">{{篇名，≤8 字}}</span><span class="rg-d">{{说明，12-26 字}}</span></div>
    <div class="rg-item"><span class="rg-n">03</span><span class="rg-t">{{篇名，≤8 字}}</span><span class="rg-d">{{说明，12-26 字}}</span></div>
    <div class="rg-item"><span class="rg-n">04</span><span class="rg-t">{{篇名，≤8 字}}</span><span class="rg-d">{{说明，12-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rg-kicker, h2, mt-m, rg-band, mt-l, rg-item, rg-n, rg-t, rg-d, deck-footer, slide-number, notes

---

## keynotes（三彩要点）
指纹：cards
数量：rg-card=3

用途：恰好三张白卡。每张：彩签胶囊 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；彩签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="rg-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="rg-card"><span class="rg-pill rg-pill-accent">{{彩签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#374151">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="rg-card"><span class="rg-pill rg-pill-accent">{{彩签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#374151">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="rg-card"><span class="rg-pill rg-pill-accent">{{彩签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#374151">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rg-kicker, h2, mt-m, grid, g3, mt-l, rg-card, rg-pill, rg-pill-accent, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左讲右办）
指纹：split
数量：rg-step=3

用途：左边把一件事讲透（lede + 补充 + 胶囊），右边白卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个胶囊；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="rg-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#374151">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="rg-pill">{{要点 1，≤8 字}}</span>
        <span class="rg-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="rg-card">
      <div class="rg-step"><span class="rg-n">1</span><p class="rg-mini-t">{{一步，12-26 字}}</p></div>
      <div class="rg-step"><span class="rg-n">2</span><p class="rg-mini-t">{{一步，12-26 字}}</p></div>
      <div class="rg-step"><span class="rg-n">3</span><p class="rg-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rg-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, rg-pill, rg-card, rg-step, rg-n, rg-mini-t, deck-footer, slide-number, notes

---

## metrics（庆祝数字）
指纹：chart
数量：rg-stat=3

用途：三个关键数据。玫红到紫的渐变大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="rg-kicker">{{数据语境，如 关键数字 · 策划口径}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="rg-stat"><div class="rg-stat-v">{{数值 ≤6 字符}}<span class="rg-stat-u">{{单位}}</span></div><div class="rg-stat-l">{{指标名，≤8 字}}</div><p class="rg-stat-note">{{口径，14-30 字}}</p></div>
    <div class="rg-stat"><div class="rg-stat-v">{{数值 ≤6 字符}}<span class="rg-stat-u">{{单位}}</span></div><div class="rg-stat-l">{{指标名，≤8 字}}</div><p class="rg-stat-note">{{口径，14-30 字}}</p></div>
    <div class="rg-stat"><div class="rg-stat-v">{{数值 ≤6 字符}}<span class="rg-stat-u">{{单位}}</span></div><div class="rg-stat-l">{{指标名，≤8 字}}</div><p class="rg-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#5F6B7A;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rg-kicker, h2, mt-m, grid, g3, mt-l, rg-stat, rg-stat-v, rg-stat-u, rg-stat-l, rg-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（欢呼引言）
指纹：quote

用途：整页一句引文。大字 + 出处 + 两个支撑胶囊，右侧留白处缀一段彩虹条。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；胶囊各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="rg-kicker">{{语境，如 组委会的话}}</p>
  <p class="rg-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="rg-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="rg-pill rg-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="rg-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="rg-stripe" style="position:absolute;right:140px;top:44%"></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rg-kicker, rg-quote, mt-l, rg-src, mt-m, row, rg-pill, rg-pill-accent, rg-stripe, deck-footer, slide-number, notes

---

## divider（彩门章节）
指纹：hero

用途：章节过渡。进度题签 + 大字章节名 + 一个过渡问题 + 两个看点胶囊。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；胶囊各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="rg-kicker">{{进度，如 第二段 · 彩虹关卡}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="rg-pill rg-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="rg-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, rg-kicker, h1, mt-m, lede, mt-l, row, rg-pill, rg-pill-accent, deck-footer, slide-number, notes

---

## moments（赛程时间线）
指纹：chart
数量：rg-tl-item=4

用途：4 个节点的横向时间线：彩虹环圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="rg-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="rg-tl mt-l" style="margin-top:52px">
    <div class="rg-tl-item"><div class="rg-tl-dot">01</div><div class="rg-tl-t">{{时间点，≤10 字符}}</div><p class="rg-tl-d">{{事件，12-26 字}}</p></div>
    <div class="rg-tl-item"><div class="rg-tl-dot">02</div><div class="rg-tl-t">{{时间点，≤10 字符}}</div><p class="rg-tl-d">{{事件，12-26 字}}</p></div>
    <div class="rg-tl-item"><div class="rg-tl-dot">03</div><div class="rg-tl-t">{{时间点，≤10 字符}}</div><p class="rg-tl-d">{{事件，12-26 字}}</p></div>
    <div class="rg-tl-item"><div class="rg-tl-dot">04</div><div class="rg-tl-t">{{时间点，≤10 字符}}</div><p class="rg-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#5F6B7A;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, rg-kicker, h2, mt-m, rg-tl, mt-l, rg-tl-item, rg-tl-dot, rg-tl-t, rg-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（冲线收尾）
指纹：hero

用途：收尾页。大字 + 一句行动提醒 + 海军蓝按钮 + 胶囊与彩点，像冲过终点线。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；胶囊 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="rg-kicker">{{提醒语境，如 报名通道}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="rg-btn">{{按钮文案，≤8 字}}</span>
    <span class="rg-pill">{{次级信息，≤10 字}}</span>
    <span class="rg-chip" style="--c:#F43F5E"></span>
    <span class="rg-chip" style="--c:#EAB308"></span>
    <span class="rg-chip" style="--c:#8B5CF6"></span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, rg-kicker, h1, mt-m, lede, mt-l, row, rg-btn, rg-pill, rg-chip, deck-footer, slide-number, notes
