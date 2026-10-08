# 晴橙落日海 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `os-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-orange-sea` 作用域生效，骨架里已写全，照抄结构即可。
> 每页的天空→暮色渐变、太阳、海鸥与波浪纹由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（橙蓝撞色）**：落日橙只允许四处——题签胶囊（os-kicker）、数据与时间点
> （os-stat 顶线 + os-stat-v / os-tl-t 收暗档）、accent 药丸软底、太阳印（os-badge）。
> 深海蓝 #1F6E96 是标题与 os-btn 的骨架色，橙蓝互为撞色、不许单边缺席。
> 禁霓虹高饱和、禁深暗背景、禁硬朗锐角与密集数据表；卡只有一种（os-card 白色大圆角）。

---

## cover（海湾封面）
指纹：hero

用途：开场页。落日橙题签 + 大字标题 + 一句定位，右侧可立一枚太阳印。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 24-48 字；题签 ≤14 字；太阳印 4 字内。

```html
<section class="slide full" data-layout="cover">
  <p class="os-kicker">{{题签，≤14 字，如 澜港文旅 · 城市推广案}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，24-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="os-badge">{{太阳印 2×2 字}}</div>
    <span class="os-pill">{{副题或时间地点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, os-kicker, h1, mt-m, lede, row, mt-l, os-badge, os-pill, deck-footer, slide-number, notes

---

## contents（目录白卡）
指纹：table
数量：os-item=4

用途：议程页。一块白色大圆角卡里放 4 行篇目：杏白圆角块编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="os-kicker">{{引导语，如 目录四章}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="os-card mt-l" style="margin-top:44px">
    <div class="os-item"><span class="os-n">01</span><span class="os-t">{{篇名，≤8 字}}</span><span class="os-d">{{说明，14-26 字}}</span></div>
    <div class="os-item"><span class="os-n">02</span><span class="os-t">{{篇名，≤8 字}}</span><span class="os-d">{{说明，14-26 字}}</span></div>
    <div class="os-item"><span class="os-n">03</span><span class="os-t">{{篇名，≤8 字}}</span><span class="os-d">{{说明，14-26 字}}</span></div>
    <div class="os-item"><span class="os-n">04</span><span class="os-t">{{篇名，≤8 字}}</span><span class="os-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, os-kicker, h2, mt-m, os-card, mt-l, os-item, os-n, os-t, os-d, deck-footer, slide-number, notes

---

## keynotes（三卡要点）
指纹：cards
数量：os-card=3

用途：恰好三张白色圆角卡。每张：落日橙题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="os-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="os-card"><span class="os-pill os-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#3E5F6E">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="os-card"><span class="os-pill os-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#3E5F6E">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="os-card"><span class="os-pill os-pill-accent">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#3E5F6E">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, os-kicker, h2, mt-m, grid, g3, mt-l, os-card, os-pill, os-pill-accent, h4, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：os-step=3

用途：左边把一件事讲透（lede + 补充 + 药丸），右边白卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个药丸；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="os-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#3E5F6E">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="os-pill">{{要点 1，≤8 字}}</span>
        <span class="os-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="os-card">
      <div class="os-step"><span class="os-n">一</span><p class="os-mini-t">{{一步，12-26 字}}</p></div>
      <div class="os-step"><span class="os-n">二</span><p class="os-mini-t">{{一步，12-26 字}}</p></div>
      <div class="os-step"><span class="os-n">三</span><p class="os-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, os-kicker, h2, mt-m, grid, g2, mt-l, lede, row, os-pill, os-card, os-step, os-n, os-mini-t, deck-footer, slide-number, notes

---

## metrics（落日数字）
指纹：chart
数量：os-stat=3

用途：三个关键数据。深橙大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="os-kicker">{{数据语境，如 澜港的数字}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="os-stat"><div class="os-stat-v">{{数值 ≤6 字符}}<span class="os-stat-u">{{单位}}</span></div><div class="os-stat-l">{{指标名，≤8 字}}</div><p class="os-stat-note">{{口径，14-30 字}}</p></div>
    <div class="os-stat"><div class="os-stat-v">{{数值 ≤6 字符}}<span class="os-stat-u">{{单位}}</span></div><div class="os-stat-l">{{指标名，≤8 字}}</div><p class="os-stat-note">{{口径，14-30 字}}</p></div>
    <div class="os-stat"><div class="os-stat-v">{{数值 ≤6 字符}}<span class="os-stat-u">{{单位}}</span></div><div class="os-stat-l">{{指标名，≤8 字}}</div><p class="os-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#7A4515;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, os-kicker, h2, mt-m, grid, g3, mt-l, os-stat, os-stat-v, os-stat-u, os-stat-l, os-stat-note, deck-footer, slide-number, notes

---

## quote（海风引文）
指纹：quote

用途：整页一句引文。大字标题气质 + 出处 + 两个支撑药丸，右侧留白处飞一对海鸥。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；药丸各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="os-kicker">{{语境，如 宣传片旁白}}</p>
  <p class="os-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="os-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="os-pill os-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="os-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="os-gulls" aria-hidden="true"></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, os-kicker, os-quote, mt-l, os-src, mt-m, row, os-pill, os-pill-accent, os-gulls, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。落日橙题签 + 大字章节名 + 一个过渡问题 + 两个看点药丸。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；药丸各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="os-kicker">{{进度，如 第二部分 · 路线与玩法}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="os-pill os-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="os-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, os-kicker, h1, mt-m, lede, mt-l, row, os-pill, os-pill-accent, deck-footer, slide-number, notes

---

## moments（日历时间线）
指纹：chart
数量：os-tl-item=4

用途：3-4 个节点的横向时间线：杏白圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="os-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="os-tl mt-l" style="margin-top:52px">
    <div class="os-tl-item"><div class="os-tl-dot">壹</div><div class="os-tl-t">{{时间点，≤10 字符}}</div><p class="os-tl-d">{{事件，12-26 字}}</p></div>
    <div class="os-tl-item"><div class="os-tl-dot">贰</div><div class="os-tl-t">{{时间点，≤10 字符}}</div><p class="os-tl-d">{{事件，12-26 字}}</p></div>
    <div class="os-tl-item"><div class="os-tl-dot">叁</div><div class="os-tl-t">{{时间点，≤10 字符}}</div><p class="os-tl-d">{{事件，12-26 字}}</p></div>
    <div class="os-tl-item"><div class="os-tl-dot">肆</div><div class="os-tl-t">{{时间点，≤10 字符}}</div><p class="os-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#7A4515;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, os-kicker, h2, mt-m, os-tl, mt-l, os-tl-item, os-tl-dot, os-tl-t, os-tl-d, deck-footer, slide-number, notes

---

## closing（收尾潮印）
指纹：hero

用途：收尾页。大字标题 + 一句行动提醒 + 深海蓝按钮 + 药丸 + 太阳印。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="os-kicker">{{提醒语境，如 投放启动}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="os-btn">{{按钮文案，≤8 字}}</span>
    <span class="os-pill">{{次级信息，≤10 字}}</span>
    <div class="os-badge">{{太阳印 2 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, os-kicker, h1, mt-m, lede, mt-l, row, os-btn, os-pill, os-badge, deck-footer, slide-number, notes
