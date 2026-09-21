# 融资路演 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 XXX> 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-pitch-deck` 作用域前缀生效，骨架里已写全，照抄结构即可。

---

## cover（路演封面）
指纹：hero

用途：开场：项目一句话定位 + 讲者，渐变光斑背景
适用 role：cover。
内容约束：主标题 ≤12 字；lede 20-40 字；kicker 是品牌短标签

```html
<section class="slide" data-layout="cover">
  <p class="kicker">{{品牌标签}}</p>
  <h1 class="h1 anim-fade-up" data-anim="fade-up">{{主标题}}</h1>
  <p class="lede mt-m">{{一句话定位，20-40 字}}</p>
  <div class="deck-footer"><span class="mono">{{话题标签}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h1, anim-fade-up, lede, mt-m, deck-footer, mono, slide-number, notes

---

## problem-cards（问题三卡）
数量：card=3
指纹：cards

用途：三个并列痛点/问题，开头编号引导
适用 role：content。
内容约束：恰好 3 卡；卡标题 ≤10 字；说明 22-45 字（一句论断 + 一句展开）

```html
<section class="slide" data-layout="problem-cards">
  <span class="section-num">{{章节数字}}</span>
  <p class="num-tag">{{引导语}}</p>
  <h2 class="h2 mt-s">{{问题陈述}}</h2>
  <div class="grid g3 mt-l">
    <div class="card"><h4>{{痛点1}}</h4><p class="dim">{{说明，22-45 字}}</p></div>
    <div class="card"><h4>{{痛点2}}</h4><p class="dim">{{说明，22-45 字}}</p></div>
    <div class="card"><h4>{{痛点3}}</h4><p class="dim">{{说明，22-45 字}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, section-num, num-tag, h2, mt-s, grid, g3, mt-l, card, h4, dim, notes

---

## pill-statement（方案宣言）
指纹：quote

用途：方案核心论点 + 关键词 pill 一排
适用 role：content / divider。
内容约束：宣言 20-40 字（一句完整论断）；lede 补充 20-45 字；pill 3-5 个、每个 ≤6 字

```html
<section class="slide" data-layout="pill-statement">
  <span class="section-num">{{章节数字}}</span>
  <p class="num-tag">{{引导语}}</p>
  <h2 class="h2 mt-s">{{方案宣言，20-40 字}}</h2>
  <p class="lede mt-m">{{补充一句，20-45 字}}</p>
  <div class="row mt-l">
    <span class="pill pill-accent">{{关键词}}</span>
    <span class="pill pill-accent">{{关键词}}</span>
    <span class="pill">{{关键词}}</span>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, section-num, num-tag, h2, mt-s, lede, mt-m, row, mt-l, pill, pill-accent, notes

---

## feature-grid（产品特性四宫格）
数量：card=4
指纹：cards

用途：产品 4 个能力点（卡带悬浮态）
适用 role：content。
内容约束：恰好 4 卡；标题 ≤8 字；说明 22-45 字

```html
<section class="slide" data-layout="feature-grid">
  <span class="section-num">{{章节数字}}</span>
  <p class="num-tag">{{引导语}}</p>
  <h2 class="h2 mt-s">{{产品能力标题}}</h2>
  <div class="grid g2 mt-l">
    <div class="card card-hover"><h4>{{能力}}</h4><p class="dim">{{说明，22-45 字}}</p></div>
    <div class="card card-hover"><h4>{{能力}}</h4><p class="dim">{{说明，22-45 字}}</p></div>
    <div class="card card-hover"><h4>{{能力}}</h4><p class="dim">{{说明，22-45 字}}</p></div>
    <div class="card card-hover"><h4>{{能力}}</h4><p class="dim">{{说明，22-45 字}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, section-num, num-tag, h2, mt-s, grid, g2, mt-l, card, card-hover, h4, dim, notes

---

## market-metrics（市场三指标）
指纹：chart

用途：市场规模 TAM/SAM/SOM 或三个核心数字
适用 role：data。
内容约束：3 个指标；指标名 ≤8 字；数字 ≤6 字符；说明 22-45 字（写口径/对比，无出处标“估算”）

```html
<section class="slide" data-layout="market-metrics">
  <span class="section-num">{{章节数字}}</span>
  <p class="num-tag">{{引导语}}</p>
  <h2 class="h2 mt-s">{{市场标题}}</h2>
  <div class="grid g3 mt-l">
    <div class="metric"><span class="l">{{指标名}}</span><span class="n">{{数字}}</span><p class="dim">{{说明，22-45 字}}</p></div>
    <div class="metric"><span class="l">{{指标名}}</span><span class="n">{{数字}}</span><p class="dim">{{说明，22-45 字}}</p></div>
    <div class="metric"><span class="l">{{指标名}}</span><span class="n">{{数字}}</span><p class="dim">{{说明，22-45 字}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, section-num, num-tag, h2, mt-s, grid, g3, mt-l, metric, l, n, dim, notes

---

## business-model（商业模式）
指纹：split

用途：怎么赚钱：定价/模式三列或两列卡
适用 role：content。
内容约束：2-3 列；列标题 ≤8 字；说明 22-45 字

```html
<section class="slide" data-layout="business-model">
  <span class="section-num">{{章节数字}}</span>
  <p class="num-tag">{{引导语}}</p>
  <h2 class="h2 mt-s">{{模式标题}}</h2>
  <div class="grid g3 mt-l">
    <div class="card card-accent"><h4>{{模式1}}</h4><p class="dim">{{说明，22-45 字}}</p></div>
    <div class="card card-accent"><h4>{{模式2}}</h4><p class="dim">{{说明，22-45 字}}</p></div>
    <div class="card card-accent"><h4>{{模式3}}</h4><p class="dim">{{说明，22-45 字}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, section-num, num-tag, h2, mt-s, grid, g3, mt-l, card, card-accent, h4, dim, notes

---

## traction-bars（增长势头）
指纹：chart

用途：traction 横向条形（用 .bar 宽度表达相对大小）
适用 role：data。
内容约束：4-6 条；条目名 ≤10 字；宽度按真实比例给 style，数字标注在条目名后；口径说明 20-40 字（无来源标“估算”）

```html
<section class="slide" data-layout="traction-bars">
  <span class="section-num">{{章节数字}}</span>
  <p class="num-tag">{{引导语}}</p>
  <h2 class="h2 mt-s">{{增长标题}}</h2>
  <div class="traction-bar mt-l">
    <div class="bar" style="width:90%"><span>{{条目名 数字}}</span></div>
    <div class="bar" style="width:60%"><span>{{条目名 数字}}</span></div>
    <div class="bar" style="width:35%"><span>{{条目名 数字}}</span></div>
  </div>
  <p class="dim mt-l">{{口径说明，20-40 字（无来源标“估算”）}}</p>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, section-num, num-tag, h2, mt-s, traction-bar, mt-l, bar, dim, notes

---

## team-cards（团队三卡）
数量：card=3
指纹：cards

用途：核心成员 3 人（头像位 + 名字 + 背景）
适用 role：content。
内容约束：恰好 3 卡；姓名 ≤6 字；背景 22-40 字

```html
<section class="slide" data-layout="team-cards">
  <span class="section-num">{{章节数字}}</span>
  <p class="num-tag">{{引导语}}</p>
  <h2 class="h2 mt-s">{{团队标题}}</h2>
  <div class="grid g3 mt-l">
    <div class="card team-card"><div class="avatar"></div><h4>{{姓名}}</h4><p class="dim">{{背景，22-40 字}}</p></div>
    <div class="card team-card"><div class="avatar"></div><h4>{{姓名}}</h4><p class="dim">{{背景，22-40 字}}</p></div>
    <div class="card team-card"><div class="avatar"></div><h4>{{姓名}}</h4><p class="dim">{{背景，22-40 字}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, section-num, num-tag, h2, mt-s, grid, g3, mt-l, card, team-card, avatar, h4, dim, notes

---

## ask-box（The Ask）
指纹：stack

用途：要什么：金额/资源/支持，一页只说这一件事
适用 role：cta。
内容约束：Ask 一句话 20-30 字；用途说明 20-45 字；支持点 2-3 条、每条 12-24 字

```html
<section class="slide" data-layout="ask-box">
  <p class="num-tag">{{引导语}}</p>
  <div class="ask-box mt-m">
    <h2 class="h2">{{Ask 一句话，20-30 字}}</h2>
    <p class="lede">{{用途说明，20-45 字}}</p>
  </div>
  <div class="row mt-l">
    <div class="dim">{{支持点，12-24 字}}</div>
    <div class="dim">{{支持点，12-24 字}}</div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, num-tag, ask-box, mt-m, h2, lede, row, mt-l, dim, notes

---

## thanks-mega（致谢大字）
指纹：hero

用途：收尾致谢 + 联系方式 pill
适用 role：thanks / cta。
内容约束：大字 ≤8 字；一句话 15-30 字；pill ≤2 个

```html
<section class="slide center tc" data-layout="thanks-mega">
  <div class="mega">{{大字}}</div>
  <p class="mega-sub">{{一句话，15-30 字}}</p>
  <div class="row mt-l" style="justify-content:center">
    <span class="pill pill-accent">{{联系方式}}</span>
    <span class="pill">{{仓库/链接}}</span>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, center, tc, mega, mega-sub, row, mt-l, pill, pill-accent, notes
