# 扁豆紫蜜陀僧 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `aa-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-amber-aurora` 作用域生效，骨架里已写全，照抄结构即可。
> 每页暮色天空、落日光晕、萤光与水面倒影由模板自动衬底（z-index 在内容之下），
> 骨架与正文都不用画。
>
> **气质铁律（黄昏治愈，意境重于密度）**：暮色深褐是唯一的深底（印章 aa-seal、
> 方章编号 aa-n / aa-tl-dot、暮色签 aa-pill-accent、按钮 aa-btn），上面只放奶白字；
> 正文一律深琥珀墨落在暮色渐变上（模板已保证对比）。落日、萤光与水面是自动衬底，
> 正文不要再画。楷宋大字（h1/h2）每页最多一组；竖排 aa-vert 只放留白处，一页至多一条。

---

## cover（暮色封面）
指纹：hero

用途：开场页。题签 + 楷宋大字标题 + 飘带线 + 一句定位，右下可立一枚暮色印。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 28-48 字；题签 ≤14 字；印章 4 字内。

```html
<section class="slide full" data-layout="cover">
  <p class="aa-kicker">{{题签，≤14 字，如 拾光香事 · 秋季系列发布}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="aa-line mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，28-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="aa-seal">{{印章 2×2 字}}</div>
    <span class="aa-pill">{{副题或时间地点，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, aa-kicker, h1, mt-m, aa-line, mt-s, lede, mt-l, row, aa-seal, aa-pill, deck-footer, slide-number, notes

---

## contents（香事目录）
指纹：table
数量：aa-item=4

用途：议程页。一块半透米白大卡里放 4 行篇目：暮色方章编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 12-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="aa-kicker">{{引导语，如 今日香目}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="aa-card mt-l" style="margin-top:44px">
    <div class="aa-item"><span class="aa-n">壹</span><span class="aa-t">{{篇名，≤8 字}}</span><span class="aa-d">{{说明，12-26 字}}</span></div>
    <div class="aa-item"><span class="aa-n">贰</span><span class="aa-t">{{篇名，≤8 字}}</span><span class="aa-d">{{说明，12-26 字}}</span></div>
    <div class="aa-item"><span class="aa-n">叁</span><span class="aa-t">{{篇名，≤8 字}}</span><span class="aa-d">{{说明，12-26 字}}</span></div>
    <div class="aa-item"><span class="aa-n">肆</span><span class="aa-t">{{篇名，≤8 字}}</span><span class="aa-d">{{说明，12-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, aa-kicker, h2, mt-m, aa-card, mt-l, aa-item, aa-n, aa-t, aa-d, deck-footer, slide-number, notes

---

## keynotes（三香要点）
指纹：cards
数量：aa-card=3

用途：恰好三张米白卡。每张：暮色签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；暮色签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="aa-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="aa-card"><span class="aa-pill aa-pill-accent">{{暮色签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#54311C">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="aa-card"><span class="aa-pill aa-pill-accent">{{暮色签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#54311C">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="aa-card"><span class="aa-pill aa-pill-accent">{{暮色签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#54311C">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, aa-kicker, h2, mt-m, grid, g3, mt-l, aa-card, aa-pill, aa-pill-accent, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左品右法）
指纹：split
数量：aa-step=3

用途：左边把一件事讲透（lede + 补充 + 胶囊），右边琥珀线米白卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个胶囊；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="aa-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#54311C">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="aa-pill">{{要点 1，≤8 字}}</span>
        <span class="aa-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="aa-card-line">
      <div class="aa-step"><span class="aa-n">一</span><p class="aa-mini-t">{{一步，12-26 字}}</p></div>
      <div class="aa-step"><span class="aa-n">二</span><p class="aa-mini-t">{{一步，12-26 字}}</p></div>
      <div class="aa-step"><span class="aa-n">三</span><p class="aa-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, aa-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, aa-pill, aa-card-line, aa-step, aa-n, aa-mini-t, deck-footer, slide-number, notes

---

## metrics（香事数字）
指纹：chart
数量：aa-stat=3

用途：三个关键数据。深蜜褐大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="aa-kicker">{{数据语境，如 用料与测试 · 口径}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="aa-stat"><div class="aa-stat-v">{{数值 ≤6 字符}}<span class="aa-stat-u">{{单位}}</span></div><div class="aa-stat-l">{{指标名，≤8 字}}</div><p class="aa-stat-note">{{口径，14-30 字}}</p></div>
    <div class="aa-stat"><div class="aa-stat-v">{{数值 ≤6 字符}}<span class="aa-stat-u">{{单位}}</span></div><div class="aa-stat-l">{{指标名，≤8 字}}</div><p class="aa-stat-note">{{口径，14-30 字}}</p></div>
    <div class="aa-stat"><div class="aa-stat-v">{{数值 ≤6 字符}}<span class="aa-stat-u">{{单位}}</span></div><div class="aa-stat-l">{{指标名，≤8 字}}</div><p class="aa-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#6E4526;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, aa-kicker, h2, mt-m, grid, g3, mt-l, aa-stat, aa-stat-v, aa-stat-u, aa-stat-l, aa-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（香语引文）
指纹：quote

用途：整页一句引文。楷宋大字 + 出处 + 两个支撑胶囊，右侧留白处可立竖排短语。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；胶囊各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="aa-kicker">{{语境，如 调香师的话}}</p>
  <p class="aa-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="aa-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="aa-pill aa-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="aa-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="aa-vert" style="position:absolute;right:120px;top:50%;transform:translateY(-50%)">{{竖排短语 ≤7 字}}<span class="aa-vert-accent">{{一字或两字}}</span></div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, aa-kicker, aa-quote, mt-l, aa-src, mt-m, row, aa-pill, aa-pill-accent, aa-vert, aa-vert-accent, deck-footer, slide-number, notes

---

## divider（暮门章节）
指纹：hero

用途：章节过渡。进度题签 + 大字章节名 + 一个过渡问题 + 两个看点胶囊。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；胶囊各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="aa-kicker">{{进度，如 第二折 · 三支古色香}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="aa-pill aa-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="aa-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, aa-kicker, h1, mt-m, lede, mt-l, row, aa-pill, aa-pill-accent, deck-footer, slide-number, notes

---

## moments（当夜时间线）
指纹：chart
数量：aa-tl-item=4

用途：4 个节点的横向时间线：暮色方点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="aa-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="aa-tl mt-l" style="margin-top:52px">
    <div class="aa-tl-item"><div class="aa-tl-dot">壹</div><div class="aa-tl-t">{{时间点，≤10 字符}}</div><p class="aa-tl-d">{{事件，12-26 字}}</p></div>
    <div class="aa-tl-item"><div class="aa-tl-dot">贰</div><div class="aa-tl-t">{{时间点，≤10 字符}}</div><p class="aa-tl-d">{{事件，12-26 字}}</p></div>
    <div class="aa-tl-item"><div class="aa-tl-dot">叁</div><div class="aa-tl-t">{{时间点，≤10 字符}}</div><p class="aa-tl-d">{{事件，12-26 字}}</p></div>
    <div class="aa-tl-item"><div class="aa-tl-dot">肆</div><div class="aa-tl-t">{{时间点，≤10 字符}}</div><p class="aa-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#6E4526;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, aa-kicker, h2, mt-m, aa-tl, mt-l, aa-tl-item, aa-tl-dot, aa-tl-t, aa-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（落款收尾）
指纹：hero

用途：收尾页。大字 + 一句行动提醒 + 暮色按钮 + 胶囊与印章，像黄昏落款。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；胶囊 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="aa-kicker">{{提醒语境，如 发布夜限定}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="aa-btn">{{按钮文案，≤8 字}}</span>
    <span class="aa-pill">{{次级信息，≤10 字}}</span>
    <div class="aa-seal">{{印章 2 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, aa-kicker, h1, mt-m, lede, mt-l, row, aa-btn, aa-pill, aa-seal, deck-footer, slide-number, notes
