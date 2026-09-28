# 空白模板 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板的 `.tpl-blank` 作用域类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 这是空白起点：只有两个最小说明版式。要扩展版式，编辑本文件登记新条目，
> 并让 index.html 的 demo 页与 style.css 的规则同步跟上。

---

## blank-cover（空白封面）
指纹：hero

用途：开场页。小引导语 + 大标题 + 一句话定位，垂直居中。
适用 role：cover。
内容约束：标题 ≤16 字；定位一句话 20-40 字。

合法类名：slide, full, kicker, h1, lede, mt-m

```html
<section class="slide full" data-layout="blank-cover">
  <p class="kicker">{{kicker · 这份演示是什么}}</p>
  <h1 class="h1 mt-m">{{主标题}}</h1>
  <p class="lede mt-m">{{一句话定位}}</p>
</section>
```

## blank-content（空白内容页）
指纹：stack

用途：内容页。小节标题 + 2-4 张卡片，自上而下。
适用 role：content。
内容约束：2-4 卡；卡标题 ≤10 字、正文 20-60 字。

合法类名：slide, h2, grid, g2, mt-l, card, h4, dim, mt-s

```html
<section class="slide" data-layout="blank-content">
  <h2 class="h2">{{小节标题}}</h2>
  <div class="grid g2 mt-l">
    <div class="card"><h4 class="h4">{{卡标题}}</h4><p class="dim mt-s">{{卡正文}}</p></div>
    <div class="card"><h4 class="h4">{{卡标题}}</h4><p class="dim mt-s">{{卡正文}}</p></div>
  </div>
</section>
```
