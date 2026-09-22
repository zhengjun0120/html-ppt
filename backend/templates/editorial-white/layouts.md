# 杂志白 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 XXX> 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-xhs-white-editorial` 作用域前缀生效，骨架里已写全，照抄结构即可。

---

## hero-quote（杂志封面）
指纹：hero

用途：刊头 + 大标题 + hero 区（可放一张图或纯排版）
适用 role：cover。
内容约束：标题 ≤12 字；副标 20-30 字；kicker 是栏目名

```html
<section class="slide" data-layout="hero-quote">
  <div class="xw-topline"></div>
  <div class="xw-topbar"><div class="xw-tag"><span class="dot"></span>{{页标签，2-8 字}}</div><div class="xw-page">{{页码，如 04 / 08}}</div></div>
  <p class="xw-kicker">{{栏目名}}</p>
  <h1 class="xw-title">{{标题主干}}<br><span class="xw-grad">{{重点短语}}</span></h1>
  <p class="xw-sub">{{副标，20-30 字，关键处可用 <span class="xw-focus">重点标注</span>}}</p>
  <div class="xw-hero"><div class="xw-quote">{{金句，18-36 字，可带 <span class="xw-focus-orange">橙色重点</span>}}</div></div>
  <div class="xw-footer"><span>{{刊名/栏目}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, xw-topline, xw-topbar, xw-tag, dot, xw-page, xw-kicker, xw-title, xw-grad, xw-sub, xw-focus, xw-hero, xw-quote, xw-focus-orange, xw-footer, notes

---

## statement（单句页）
指纹：quote

用途：整页只放一句转场/观点，杂志翻页感
适用 role：divider / quote。
内容约束：单句 12-24 字

```html
<section class="slide" data-layout="statement">
  <div class="xw-topline"></div>
  <div class="xw-topbar"><div class="xw-tag"><span class="dot"></span>{{页标签，2-8 字}}</div><div class="xw-page">{{页码，如 04 / 08}}</div></div>
  <p class="xw-kicker">{{栏目}}</p>
  <h1 class="xw-title">{{单句，12-24 字}}</h1>
  <div class="xw-footer"><span>{{刊名/栏目}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, xw-topline, xw-topbar, xw-tag, dot, xw-page, xw-kicker, xw-title, xw-footer, notes

---

## quad-cards（马卡龙四卡）
数量：xw-card=4
指纹：cards

用途：四张浅色卡（soft-pink/blue/green/orange 各一，颜色即分类）
适用 role：content。
内容约束：恰好 4 卡；卡标题 ≤8 字；卡内说明 22-40 字（一句论断 + 一句展开）；颜色与内容情感匹配

```html
<section class="slide" data-layout="quad-cards">
  <div class="xw-topline"></div>
  <div class="xw-topbar"><div class="xw-tag"><span class="dot"></span>{{页标签，2-8 字}}</div><div class="xw-page">{{页码，如 04 / 08}}</div></div>
  <h2 class="xw-title-md">{{主题}}</h2>
  <div class="xw-grid-2 mt-l">
    <div class="xw-card soft-pink"><div class="xw-label">{{栏目标签，2-4 字}}</div><div class="main">{{主体，8-14 字}}</div><div class="desc">{{一句注脚，10-20 字}}</div></div>
    <div class="xw-card soft-blue"><div class="xw-label">{{栏目标签}}</div><div class="main">{{主体}}</div><div class="desc">{{注脚}}</div></div>
    <div class="xw-card soft-green"><div class="xw-label">{{栏目标签}}</div><div class="main">{{主体}}</div><div class="desc">{{注脚}}</div></div>
    <div class="xw-card soft-orange"><div class="xw-label">{{栏目标签}}</div><div class="main">{{主体}}</div><div class="desc">{{注脚}}</div></div>
  </div>
  <div class="xw-footer"><span>{{栏目名}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, xw-topline, xw-topbar, xw-tag, dot, xw-page, xw-title-md, xw-grid-2, mt-l, xw-card, soft-pink, soft-blue, soft-green, soft-orange, xw-label, main, desc, xw-footer, notes

---

## steps（步骤指南）
指纹：table

用途：操作/方法步骤（xw-step 编号步进）
适用 role：content。
内容约束：3-5 步；每步 12-24 字、动词开头

```html
<section class="slide" data-layout="steps">
  <div class="xw-topline"></div>
  <div class="xw-topbar"><div class="xw-tag"><span class="dot"></span>{{页标签，2-8 字}}</div><div class="xw-page">{{页码，如 04 / 08}}</div></div>
  <h2 class="xw-title-md">{{标题主干}}<span class="xw-grad">{{重点词，2-6 字}}</span>{{尾缀（可空）}}</h2>
  <div class="xw-steps mt-l">
    <div class="xw-step"><div class="xw-num">1</div><div class="xw-txt">{{第 1 步，12-24 字}}</div></div>
    <div class="xw-step"><div class="xw-num">2</div><div class="xw-txt">{{第 2 步，12-24 字}}</div></div>
    <div class="xw-step"><div class="xw-num">3</div><div class="xw-txt">{{第 3 步，12-24 字}}</div></div>
    <div class="xw-step"><div class="xw-num">4</div><div class="xw-txt">{{第 4 步，12-24 字}}</div></div>
  </div>
  <div class="xw-hero"><div class="xw-quote">{{收束金句，18-36 字}}</div></div>
  <div class="xw-footer"><span>{{栏目名}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, xw-topline, xw-topbar, xw-tag, dot, xw-page, xw-title-md, xw-grad, xw-steps, mt-l, xw-step, xw-num, xw-txt, xw-hero, xw-quote, xw-footer, notes

---

## big-stat（大数字）
指纹：chart

用途：一个关键数字/结论的杂志式呈现（渐变大字）
适用 role：data。
内容约束：数字 ≤8 字符；来源/口径说明 20-40 字，无出处标“估算”

```html
<section class="slide" data-layout="big-stat">
  <div class="xw-topline"></div>
  <div class="xw-topbar"><div class="xw-tag"><span class="dot"></span>{{页标签，2-8 字}}</div><div class="xw-page">{{页码，如 05 / 08}}</div></div>
  <p class="xw-kicker">{{指标语境}}</p>
  <div class="xw-big-stat xw-grad mt-l">{{数字}}</div>
  <h2 class="xw-title-md">{{指标名}}</h2>
  <p class="xw-sub">{{来源/口径说明，20-40 字}}</p>
  <div class="xw-footer"><span>{{栏目名}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, xw-topline, xw-topbar, xw-tag, dot, xw-page, xw-kicker, xw-big-stat, xw-grad, mt-l, xw-title-md, xw-sub, xw-footer, notes

---

## two-column（双栏对照）
数量：xw-card=2
指纹：cards

用途：两个角度/方案的双栏卡片
适用 role：content。
内容约束：恰好 2 卡；卡内说明 22-45 字（一句论断 + 一句展开）

```html
<section class="slide" data-layout="two-column">
  <div class="xw-topline"></div>
  <div class="xw-topbar"><div class="xw-tag"><span class="dot"></span>{{页标签，2-8 字}}</div><div class="xw-page">{{页码，如 04 / 08}}</div></div>
  <h2 class="xw-title-md">{{对照标题}}</h2>
  <div class="xw-grid-2 mt-l">
    <div class="xw-card soft-blue"><div class="xw-label">{{栏 A 标签}}</div><div class="main">{{主体，8-14 字}}</div><div class="desc">{{一句注脚}}</div></div>
    <div class="xw-card soft-purple"><div class="xw-label">{{栏 B 标签}}</div><div class="main">{{主体}}</div><div class="desc">{{注脚}}</div></div>
  </div>
  <div class="xw-footer"><span>{{栏目名}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, xw-topline, xw-topbar, xw-tag, dot, xw-page, xw-title-md, xw-grid-2, mt-l, xw-card, soft-blue, soft-purple, xw-label, main, desc, xw-footer, notes

---

## quote（金句页）
指纹：quote

用途：一句金句/书中原话（xw-quote 衬线排版）
适用 role：quote。
内容约束：引文 15-34 字；出处真实，没有就不写出处行

```html
<section class="slide" data-layout="quote">
  <div class="xw-topline"></div>
  <div class="xw-topbar"><div class="xw-tag"><span class="dot"></span>{{页标签，2-8 字}}</div><div class="xw-page">{{页码，如 04 / 08}}</div></div>
  <h2 class="xw-title-md">{{标题主干}}<br>开始练 <span class="xw-grad">{{重点词}}</span></h2>
  <div class="xw-grid-3 mt-l">
    <div class="xw-card soft-purple"><div class="xw-label">{{时间标签 1，如 Tonight}}</div><div class="main">{{主体，8-14 字}}</div><div class="desc">{{一句注脚，10-20 字}}</div></div>
    <div class="xw-card soft-blue"><div class="xw-label">{{时间标签 2}}</div><div class="main">{{主体}}</div><div class="desc">{{注脚}}</div></div>
    <div class="xw-card soft-green"><div class="xw-label">{{时间标签 3}}</div><div class="main">{{主体}}</div><div class="desc">{{注脚}}</div></div>
  </div>
  <div class="xw-hero"><div class="xw-quote">{{收束金句，18-36 字，可带 <span class="xw-focus">重点标注</span>}}</div></div>
  <div class="xw-footer"><span>{{栏目名}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, xw-topline, xw-topbar, xw-tag, dot, xw-page, xw-title-md, xw-grad, xw-grid-3, mt-l, xw-card, soft-purple, soft-blue, soft-green, xw-label, main, desc, xw-hero, xw-quote, xw-focus, xw-footer, notes

---

## grid-3（三栏要点）
数量：xw-card=3
指纹：cards

用途：三个并列要点的三栏卡（可混用马卡龙色）
适用 role：content。
内容约束：恰好 3 卡；卡内说明 22-40 字（一句论断 + 一句展开）

```html
<section class="slide" data-layout="grid-3">
  <div class="xw-topline"></div>
  <div class="xw-topbar"><div class="xw-tag"><span class="dot"></span>{{页标签，2-8 字}}</div><div class="xw-page">{{页码，如 04 / 08}}</div></div>
  <h2 class="xw-title-md">{{主题}}</h2>
  <div class="xw-grid-3 mt-l">
    <div class="xw-card soft-blue"><div class="xw-label">{{要点标签}}</div><div class="main">{{主体，8-14 字}}</div><div class="desc">{{一句注脚}}</div></div>
    <div class="xw-card soft-green"><div class="xw-label">{{要点标签}}</div><div class="main">{{主体}}</div><div class="desc">{{注脚}}</div></div>
    <div class="xw-card soft-orange"><div class="xw-label">{{要点标签}}</div><div class="main">{{主体}}</div><div class="desc">{{注脚}}</div></div>
  </div>
  <div class="xw-footer"><span>{{栏目名}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, xw-topline, xw-topbar, xw-tag, dot, xw-page, xw-title-md, xw-grid-3, mt-l, xw-card, soft-blue, soft-green, soft-orange, xw-label, main, desc, xw-footer, notes
