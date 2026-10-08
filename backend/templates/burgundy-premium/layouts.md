# 勃艮第红高级感 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `bg-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-burgundy-premium` 作用域生效，骨架里已写全，照抄结构即可。
> 每页的金色细框与页角橡木年轮由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（金线是身份）**：金色只出现在细处——题签细线、bg-rule 金线、bg-mark 徽记、
> bg-n/bg-tl-dot 描边、bg-btn 描边、关键数据（bg-stat-v / bg-tl-t / bg-pill-accent）。
> 奶白是唯一的亮色面；禁大面积金色、禁白色卡、禁第三种亮色；页底 --bg 是酒红渐变中段值，
> 深色文字禁用（正文一律奶白系）。h1/h2 每页最多一组。

---

## cover（晚宴封面）
指纹：hero

用途：开场页。金色题签 + 衬线大字标题 + 金线 + 一句定位，配庄园徽记与时间地点。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；徽记 1 字。

```html
<section class="slide full" data-layout="cover">
  <p class="bg-kicker">{{题签，≤14 字，如 橡木河酒庄 · 品鉴会}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="bg-rule mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="bg-mark">{{徽记 1 字}}</div>
    <span class="bg-pill">{{时间地点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{主办 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, bg-kicker, h1, mt-m, bg-rule, mt-s, lede, mt-l, row, bg-mark, bg-pill, deck-footer, slide-number, notes

---

## contents（晚宴次序）
指纹：table
数量：bg-item=4

用途：议程页。一块酒标大卡里放 4 行次序：金框编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="bg-kicker">{{引导语，如 当晚次序}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="bg-card mt-l" style="margin-top:44px">
    <div class="bg-item"><span class="bg-n">01</span><span class="bg-t">{{篇名，≤8 字}}</span><span class="bg-d">{{说明，14-26 字}}</span></div>
    <div class="bg-item"><span class="bg-n">02</span><span class="bg-t">{{篇名，≤8 字}}</span><span class="bg-d">{{说明，14-26 字}}</span></div>
    <div class="bg-item"><span class="bg-n">03</span><span class="bg-t">{{篇名，≤8 字}}</span><span class="bg-d">{{说明，14-26 字}}</span></div>
    <div class="bg-item"><span class="bg-n">04</span><span class="bg-t">{{篇名，≤8 字}}</span><span class="bg-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{主办 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bg-kicker, h2, mt-m, bg-card, mt-l, bg-item, bg-n, bg-t, bg-d, deck-footer, slide-number, notes

---

## keynotes（三卡看点）
指纹：cards
数量：bg-card=3

用途：恰好三张酒标卡。每张：金色题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="bg-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="bg-card"><span class="bg-pill bg-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#D9C8B6">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="bg-card"><span class="bg-pill bg-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#D9C8B6">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="bg-card"><span class="bg-pill bg-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#D9C8B6">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{主办 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bg-kicker, h2, mt-m, grid, g3, mt-l, bg-card, bg-pill, bg-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左叙右列）
指纹：split
数量：bg-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边酒标卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="bg-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#D9C8B6">{{补充：判断标准或例外安排，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="bg-pill">{{要点 1，≤8 字}}</span>
        <span class="bg-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="bg-card">
      <div class="bg-step"><span class="bg-n">一</span><p class="bg-mini-t">{{一行，12-26 字}}</p></div>
      <div class="bg-step"><span class="bg-n">二</span><p class="bg-mini-t">{{一行，12-26 字}}</p></div>
      <div class="bg-step"><span class="bg-n">三</span><p class="bg-mini-t">{{一行，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{主办 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bg-kicker, h2, mt-m, grid, g2, mt-l, lede, row, bg-pill, bg-card, bg-step, bg-n, bg-mini-t, deck-footer, slide-number, notes

---

## metrics（酒窖数字）
指纹：chart
数量：bg-stat=3

用途：三个关键数据。亮金衬线大数字是唯一主视觉，口径写进块内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="bg-kicker">{{数据语境，如 酒单数字}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="bg-stat"><div class="bg-stat-v">{{数值 ≤6 字符}}<span class="bg-stat-u">{{单位}}</span></div><div class="bg-stat-l">{{指标名，≤8 字}}</div><p class="bg-stat-note">{{口径，14-30 字}}</p></div>
    <div class="bg-stat"><div class="bg-stat-v">{{数值 ≤6 字符}}<span class="bg-stat-u">{{单位}}</span></div><div class="bg-stat-l">{{指标名，≤8 字}}</div><p class="bg-stat-note">{{口径，14-30 字}}</p></div>
    <div class="bg-stat"><div class="bg-stat-v">{{数值 ≤6 字符}}<span class="bg-stat-u">{{单位}}</span></div><div class="bg-stat-l">{{指标名，≤8 字}}</div><p class="bg-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#B79E8A;letter-spacing:.06em">来源：{{出处与统计口径，含「以当日为准」一类限定，14-40 字}}</p>
  <div class="deck-footer"><span>{{主办 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bg-kicker, h2, mt-m, grid, g3, mt-l, bg-stat, bg-stat-v, bg-stat-u, bg-stat-l, bg-stat-note, deck-footer, slide-number, notes

---

## quote（题酒引文）
指纹：quote

用途：整页一句引文。奶白衬线大字 + 出处 + 两个金色描边药丸，页角年轮衬底。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="bg-kicker">{{语境，如 主理人题酒}}</p>
  <p class="bg-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="bg-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="bg-pill bg-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="bg-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{主办 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bg-kicker, bg-quote, mt-l, bg-src, mt-m, row, bg-pill, bg-pill-accent, deck-footer, slide-number, notes

---

## divider（篇章幕）
指纹：hero

用途：篇章过渡。进度题签 + 衬线大字篇章名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：篇章名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="bg-kicker">{{进度，如 第二幕 · 纵评五款}}</p>
  <h1 class="h1 mt-m">{{篇章标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="bg-pill bg-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="bg-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{主办 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, bg-kicker, h1, mt-m, lede, mt-l, row, bg-pill, bg-pill-accent, deck-footer, slide-number, notes

---

## moments（流程时间线）
指纹：chart
数量：bg-tl-item=4

用途：3-4 个节点的横向时间线：金圈圆点 + 时间点 + 一句安排。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="bg-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="bg-tl mt-l" style="margin-top:52px">
    <div class="bg-tl-item"><div class="bg-tl-dot">01</div><div class="bg-tl-t">{{时间点，≤10 字符}}</div><p class="bg-tl-d">{{事件，12-26 字}}</p></div>
    <div class="bg-tl-item"><div class="bg-tl-dot">02</div><div class="bg-tl-t">{{时间点，≤10 字符}}</div><p class="bg-tl-d">{{事件，12-26 字}}</p></div>
    <div class="bg-tl-item"><div class="bg-tl-dot">03</div><div class="bg-tl-t">{{时间点，≤10 字符}}</div><p class="bg-tl-d">{{事件，12-26 字}}</p></div>
    <div class="bg-tl-item"><div class="bg-tl-dot">04</div><div class="bg-tl-t">{{时间点，≤10 字符}}</div><p class="bg-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#B79E8A;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{主办 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, bg-kicker, h2, mt-m, bg-tl, mt-l, bg-tl-item, bg-tl-dot, bg-tl-t, bg-tl-d, deck-footer, slide-number, notes

---

## closing（订席收尾）
指纹：hero

用途：收尾页。衬线大字 + 一句行动提醒 + 金框按钮 + 徽记药丸，如请柬落款。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="bg-kicker">{{提醒语境，如 席位预订}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="bg-btn">{{按钮文案，≤8 字}}</span>
    <span class="bg-pill">{{次级信息，≤10 字}}</span>
    <div class="bg-mark">{{徽记 1 字}}</div>
  </div>
  <div class="deck-footer"><span>{{主办 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, bg-kicker, h1, mt-m, lede, mt-l, row, bg-btn, bg-pill, bg-mark, deck-footer, slide-number, notes
