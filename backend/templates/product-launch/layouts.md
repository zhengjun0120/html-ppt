# 产品发布 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 XXX> 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-product-launch` 作用域前缀生效，骨架里已写全，照抄结构即可。

---

## cover（发布封面）
指纹：hero

用途：暗色 hero：产品名大标题 + 一句话主张
适用 role：cover。
内容约束：标题 ≤10 字；lede 20-30 字；kicker 是发布主题标签

```html
<section class="slide dark" data-layout="cover">
  <p class="kicker">{{发布标签}}</p>
  <h1 class="h1 anim-fade-up" data-anim="fade-up">{{产品名/主张}}</h1>
  <p class="lede mt-m">{{一句话主张，20-30 字}}</p>
  <div class="deck-footer"><span class="brand">{{品牌}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, dark, kicker, h1, anim-fade-up, lede, mt-m, deck-footer, brand, slide-number, notes

---

## introducing（发布宣言）
指纹：quote

用途：居中宣言页：新版本来袭，一页只说一句话
适用 role：divider。
内容约束：宣言 ≤12 字；副句 15-30 字

```html
<section class="slide center tc" data-layout="introducing">
  <p class="kicker">{{引导语}}</p>
  <h1 class="h1">{{宣言}}</h1>
  <p class="lede">{{副句，15-30 字}}</p>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, center, tc, kicker, h1, lede, notes

---

## feature-trio（特性三连）
指纹：cards
数量：feature-card=3

用途：三个核心特性（feature-card 逐卡）
适用 role：content。
内容约束：恰好 3 卡；卡标题 ≤8 字；说明 22-45 字（写用户收益，一句论断 + 一句展开）

```html
<section class="slide" data-layout="feature-trio">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{特性主题}}</h2>
  <div class="grid g3 mt-l">
    <div class="feature-card"><div class="icon">{{符号，如 ♪}}</div><h4>{{特性}}</h4><p class="dim">{{说明，22-45 字}}</p></div>
    <div class="feature-card"><div class="icon">{{符号}}</div><h4>{{特性}}</h4><p class="dim">{{说明，22-45 字}}</p></div>
    <div class="feature-card"><div class="icon">{{符号}}</div><h4>{{特性}}</h4><p class="dim">{{说明，22-45 字}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, grid, g3, mt-l, feature-card, icon, h4, dim, notes

---

## fit-cards（场景适配）
指纹：cards
数量：card=3

用途：适用场景/人群三卡（暗色页变体）
适用 role：content。
内容约束：恰好 3 卡；说明 22-45 字

```html
<section class="slide dark" data-layout="fit-cards">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{场景标题}}</h2>
  <div class="grid g3 mt-l">
    <div class="card"><h4>{{场景}}</h4><p class="dim">{{说明，22-45 字}}</p></div>
    <div class="card"><h4>{{场景}}</h4><p class="dim">{{说明，22-45 字}}</p></div>
    <div class="card"><h4>{{场景}}</h4><p class="dim">{{说明，22-45 字}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, dark, kicker, h2, grid, g3, mt-l, card, h4, dim, notes

---

## feature-duo（深度两卡）
指纹：cards
数量：feature-card=2

用途：两个重点能力各占一卡，讲深一点
适用 role：content。
内容约束：恰好 2 卡；卡标题 ≤10 字；说明 30-55 字（讲深一点）

```html
<section class="slide" data-layout="feature-duo">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{能力标题}}</h2>
  <div class="grid g2 mt-l">
    <div class="feature-card"><div class="icon">{{符号，如 ✦}}</div><h4>{{能力}}</h4><p class="dim">{{说明，30-55 字}}</p></div>
    <div class="feature-card"><div class="icon">{{符号}}</div><h4>{{能力}}</h4><p class="dim">{{说明，30-55 字}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, grid, g2, mt-l, feature-card, icon, h4, dim, notes

---

## how-it-works（三步上手）
指纹：table
数量：step=3

用途：使用流程三步（step 逐行）
适用 role：content。
内容约束：恰好 3 步；每步 12-24 字、动词开头

```html
<section class="slide" data-layout="how-it-works">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{流程标题}}</h2>
  <div class="stack mt-l">
    <div class="step"><div class="n">1</div><div><h4>{{第 1 步，动词开头 ≤8 字}}</h4><p class="dim">{{展开一句，12-24 字}}</p></div></div>
    <div class="step"><div class="n">2</div><div><h4>{{第 2 步}}</h4><p class="dim">{{展开一句}}</p></div></div>
    <div class="step"><div class="n">3</div><div><h4>{{第 3 步}}</h4><p class="dim">{{展开一句}}</p></div></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, stack, mt-l, step, n, h4, dim, notes

---

## pricing（定价卡）
指纹：cards
数量：price-card=3

用途：定价三档（price-card，可标推荐档）
适用 role：content / cta。
内容约束：恰好 3 档；档名 ≤6 字；每档包含项 3-4 条、每条 12-24 字；价格真实，没有定价就写“联系我们”

```html
<section class="slide" data-layout="pricing">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{定价标题}}</h2>
  <div class="grid g3 mt-l">
    <div class="price-card"><span class="amount">{{价格}}</span><h4>{{档名}}</h4><ul><li>{{包含内容，12-24 字}}</li></ul></div>
    <div class="price-card"><span class="amount">{{价格}}</span><h4>{{档名}}</h4><ul><li>{{包含内容，12-24 字}}</li></ul></div>
    <div class="price-card"><span class="amount">{{价格}}</span><h4>{{档名}}</h4><ul><li>{{包含内容，12-24 字}}</li></ul></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, grid, g3, mt-l, price-card, amount, h4, notes

---

## ship（发布号召）
指纹：hero

用途：收尾号召：现在就用/扫码/链接
适用 role：cta / thanks。
内容约束：号召 ≤10 字；lede 15-30 字

```html
<section class="slide center tc" data-layout="ship">
  <p class="kicker">{{发布日期}}</p>
  <h1 class="h1">{{号召语}}</h1>
  <p class="lede">{{获取方式，15-30 字}}</p>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, center, tc, kicker, h1, lede, notes
