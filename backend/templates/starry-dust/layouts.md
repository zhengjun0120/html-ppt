# 星屑柔光 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `sd-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-starry-dust` 作用域生效，骨架里已写全，照抄结构即可。
> 每页右上月光碎星与页底云朵由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（柔光童话，不是国风深色）**：全模板无深色背景，紫罗兰是唯一深色（标题、
> 按钮、徽章 sd-badge / sd-btn）。暖金星光只出现在题签胶囊 ✦、渐变线 sd-line、数据顶线
> （sd-stat 自带）与散点 sd-spark 四处。卡一律棉花糖大圆角 + 柔影，禁直角与硬阴影；
> 圆体大字（h1/h2）每页最多一组；散点 sd-spark 一页至多一处、放留白处。

---

## cover（星夜封面）
指纹：hero

用途：开场页。题签胶囊 + 圆体大字标题 + 渐变线 + 一句定位，左侧可立一枚紫色徽章。
适用 role：cover。
内容约束：主标题 ≤10 字；lede 28-48 字；题签 ≤14 字；徽章 4 字内。

```html
<section class="slide full" data-layout="cover">
  <p class="sd-kicker">{{题签，≤14 字，如 新书首发 · 秋季童书展}}</p>
  <h1 class="h1 mt-m">{{主标题，≤10 字，可 <br> 分行}}</h1>
  <div class="sd-line mt-s"></div>
  <p class="lede mt-m" style="max-width:50ch">{{一句话定位，28-48 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <div class="sd-badge">{{徽章 2×2 字}}</div>
    <span class="sd-pill">{{副题或时间地点，≤12 字}}</span>
    <span class="sd-spark"></span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, sd-kicker, h1, mt-m, sd-line, mt-s, lede, mt-l, row, sd-badge, sd-pill, sd-spark, deck-footer, slide-number, notes

---

## contents（星图目录）
指纹：table
数量：sd-item=4

用途：议程页。一块棉花糖大卡里放 4 行篇目：星圈编号 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 12-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="sd-kicker">{{引导语，如 今晚星图}}</p>
  <h2 class="h2 mt-m">{{标题，≤8 字}}</h2>
  <div class="sd-card mt-l" style="margin-top:44px">
    <div class="sd-item"><span class="sd-n">01</span><span class="sd-t">{{篇名，≤8 字}}</span><span class="sd-d">{{说明，12-26 字}}</span></div>
    <div class="sd-item"><span class="sd-n">02</span><span class="sd-t">{{篇名，≤8 字}}</span><span class="sd-d">{{说明，12-26 字}}</span></div>
    <div class="sd-item"><span class="sd-n">03</span><span class="sd-t">{{篇名，≤8 字}}</span><span class="sd-d">{{说明，12-26 字}}</span></div>
    <div class="sd-item"><span class="sd-n">04</span><span class="sd-t">{{篇名，≤8 字}}</span><span class="sd-d">{{说明，12-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sd-kicker, h2, mt-m, sd-card, mt-l, sd-item, sd-n, sd-t, sd-d, deck-footer, slide-number, notes

---

## keynotes（三星要点）
指纹：cards
数量：sd-card=3

用途：恰好三张棉花糖卡。每张：星签胶囊 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；星签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="sd-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:32px;margin-top:44px">
    <div class="sd-card"><span class="sd-pill sd-pill-accent">{{星签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#6E5F96">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="sd-card"><span class="sd-pill sd-pill-accent">{{星签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#6E5F96">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="sd-card"><span class="sd-pill sd-pill-accent">{{星签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:19px;line-height:1.8;color:#6E5F96">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sd-kicker, h2, mt-m, grid, g3, mt-l, sd-card, sd-pill, sd-pill-accent, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左梦右径）
指纹：split
数量：sd-step=3

用途：左边把一件事讲透（lede + 补充 + 胶囊），右边薰衣草糖卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个胶囊；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="sd-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:19px;line-height:1.8;color:#6E5F96">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="sd-pill">{{要点 1，≤8 字}}</span>
        <span class="sd-pill">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="sd-card-tint">
      <div class="sd-step"><span class="sd-n">1</span><p class="sd-mini-t">{{一步，12-26 字}}</p></div>
      <div class="sd-step"><span class="sd-n">2</span><p class="sd-mini-t">{{一步，12-26 字}}</p></div>
      <div class="sd-step"><span class="sd-n">3</span><p class="sd-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sd-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, sd-pill, sd-card-tint, sd-step, sd-n, sd-mini-t, deck-footer, slide-number, notes

---

## metrics（星光数字）
指纹：chart
数量：sd-stat=3

用途：三个关键数据。紫罗兰大数字是唯一主视觉，口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="sd-kicker">{{数据语境，如 首印筹备 · 数据}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:44px;margin-top:48px">
    <div class="sd-stat"><div class="sd-stat-v">{{数值 ≤6 字符}}<span class="sd-stat-u">{{单位}}</span></div><div class="sd-stat-l">{{指标名，≤8 字}}</div><p class="sd-stat-note">{{口径，14-30 字}}</p></div>
    <div class="sd-stat"><div class="sd-stat-v">{{数值 ≤6 字符}}<span class="sd-stat-u">{{单位}}</span></div><div class="sd-stat-l">{{指标名，≤8 字}}</div><p class="sd-stat-note">{{口径，14-30 字}}</p></div>
    <div class="sd-stat"><div class="sd-stat-v">{{数值 ≤6 字符}}<span class="sd-stat-u">{{单位}}</span></div><div class="sd-stat-l">{{指标名，≤8 字}}</div><p class="sd-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#7E71A4;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sd-kicker, h2, mt-m, grid, g3, mt-l, sd-stat, sd-stat-v, sd-stat-u, sd-stat-l, sd-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（睡前题语）
指纹：quote

用途：整页一句引文。圆体大字 + 出处 + 两个支撑胶囊，右侧留白处缀一颗散点星光。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；胶囊各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="sd-kicker">{{语境，如 作者的话}}</p>
  <p class="sd-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="sd-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="sd-pill sd-pill-accent">{{支撑点 1，≤8 字}}</span>
    <span class="sd-pill">{{支撑点 2，≤8 字}}</span>
  </div>
  <span class="sd-spark" style="position:absolute;right:150px;top:40%"></span>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sd-kicker, sd-quote, mt-l, sd-src, mt-m, row, sd-pill, sd-pill-accent, sd-spark, deck-footer, slide-number, notes

---

## divider（月门章节）
指纹：hero

用途：章节过渡。进度题签 + 圆体大字章节名 + 一个过渡问题 + 两个看点胶囊。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-44 字；胶囊各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="sd-kicker">{{进度，如 第二段 · 三个晚安故事}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="sd-pill sd-pill-accent">{{看点 1，≤8 字}}</span>
    <span class="sd-pill">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, sd-kicker, h1, mt-m, lede, mt-l, row, sd-pill, sd-pill-accent, deck-footer, slide-number, notes

---

## moments（星轨时间线）
指纹：chart
数量：sd-tl-item=4

用途：4 个节点的横向时间线：星圈圆点 + 时间点 + 一句事件。
适用 role：content / data。
内容约束：4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="sd-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="sd-tl mt-l" style="margin-top:52px">
    <div class="sd-tl-item"><div class="sd-tl-dot">01</div><div class="sd-tl-t">{{时间点，≤10 字符}}</div><p class="sd-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sd-tl-item"><div class="sd-tl-dot">02</div><div class="sd-tl-t">{{时间点，≤10 字符}}</div><p class="sd-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sd-tl-item"><div class="sd-tl-dot">03</div><div class="sd-tl-t">{{时间点，≤10 字符}}</div><p class="sd-tl-d">{{事件，12-26 字}}</p></div>
    <div class="sd-tl-item"><div class="sd-tl-dot">04</div><div class="sd-tl-t">{{时间点，≤10 字符}}</div><p class="sd-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:17px;color:#7E71A4;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sd-kicker, h2, mt-m, sd-tl, mt-l, sd-tl-item, sd-tl-dot, sd-tl-t, sd-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（晚安收尾）
指纹：hero

用途：收尾页。圆体大字 + 一句行动提醒 + 星光按钮 + 胶囊，像睡前道一声晚安。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤8 字；胶囊 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="sd-kicker">{{提醒语境，如 首发购买}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句收束提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:26px">
    <span class="sd-btn">{{按钮文案，≤8 字}}</span>
    <span class="sd-pill">{{次级信息，≤10 字}}</span>
    <div class="sd-badge">{{徽章 2 字}}</div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, sd-kicker, h1, mt-m, lede, mt-l, row, sd-btn, sd-pill, sd-badge, deck-footer, slide-number, notes
