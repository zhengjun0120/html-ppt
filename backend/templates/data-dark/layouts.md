# 数据深空 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 XXX> 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-graphify-dark-graph` 作用域前缀生效，骨架里已写全，照抄结构即可。

---

## opening（数据开场）
指纹：hero

用途：开场：编号 + 大标题 + 一句话钩子
适用 role：cover。
内容约束：标题 ≤14 字；lede 20-30 字；kicker 是数据集/报告名

```html
<section class="slide" data-layout="opening">
  <div class="gd-ambient"></div>
  <div class="gd-snum">{{页码/编号}}</div>
  <p class="gd-eyebrow">{{报告名}}</p>
  <h1 class="gd-h1">{{标题}}</h1>
  <p class="gd-lede">{{钩子一句话，20-30 字}}</p>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gd-ambient, gd-snum, gd-eyebrow, gd-h1, gd-lede, notes

---

## context（背景铺垫）
指纹：stack

用途：研究背景/问题定义（居中叙事）
适用 role：divider / content。
内容约束：标题 ≤14 字；lede 20-40 字

```html
<section class="slide" data-layout="context">
  <div class="gd-ambient"></div>
  <div class="gd-snum">{{编号}}</div>
  <div class="gd-eyebrow">{{章节标签}}</div>
  <h1 class="gd-h1">{{问题/背景}}</h1>
  <p class="gd-lede">{{展开一句，20-40 字}}</p>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gd-ambient, gd-snum, gd-eyebrow, gd-h1, gd-lede, notes

---

## glass-grid（玻璃数据卡）
指纹：cards

用途：3-4 张玻璃卡并列（每卡一个维度：数字或事实）
适用 role：data / content。
内容约束：3-4 卡；卡标题 ≤8 字；卡内说明 22-40 字（一句论断 + 一句展开）；数字无出处标“估算”

```html
<section class="slide" data-layout="glass-grid">
  <div class="gd-ambient"></div>
  <div class="gd-snum">{{编号}}</div>
  <p class="gd-eyebrow">{{引导语}}</p>
  <h2 class="gd-h2">{{维度标题}}</h2>
  <div class="gd-grid-3 mt-l">
    <div class="gd-glass gd-glass-blue"><span class="gd-tag">{{维度}}</span><p>{{内容，22-40 字}}</p></div>
    <div class="gd-glass gd-glass-green"><span class="gd-tag">{{维度}}</span><p>{{内容，22-40 字}}</p></div>
    <div class="gd-glass gd-glass-warm"><span class="gd-tag">{{维度}}</span><p>{{内容，22-40 字}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gd-ambient, gd-snum, gd-eyebrow, gd-h2, gd-grid-3, mt-l, gd-glass, gd-glass-blue, gd-glass-green, gd-glass-warm, gd-tag, notes

---

## big-num（彩虹大数）
指纹：chart

用途：全页一个核心数字（gd-big 彩虹渐变字）
适用 role：data。
内容约束：数字 ≤8 字符；指标名 ≤12 字；口径说明 20-40 字

```html
<section class="slide" data-layout="big-num">
  <div class="gd-ambient"></div>
  <div class="gd-snum">{{编号}}</div>
  <p class="gd-eyebrow">{{指标语境}}</p>
  <div class="gd-big gd-rainbow mt-l">{{数字}}</div>
  <h2 class="gd-h2">{{指标名}}</h2>
  <p class="gd-lede">{{口径说明，20-40 字（无来源标“估算”）}}</p>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gd-ambient, gd-snum, gd-eyebrow, gd-big, gd-rainbow, mt-l, gd-h2, gd-lede, notes

---

## terminal（命令行实证）
指纹：code

用途：终端命令 + 输出（gd-cmd 行）
适用 role：code。
内容约束：命令 1-2 条 + 输出 ≤6 行；每行 ≤64 字符

```html
<section class="slide" data-layout="terminal">
  <div class="gd-ambient"></div>
  <div class="gd-snum">{{编号}}</div>
  <p class="gd-eyebrow">{{实证语境}}</p>
  <h2 class="gd-h2">{{这条命令证明什么}}</h2>
  <pre class="gd-codebox mt-l"><span class="gd-cmd">$ {{命令}}</span>
{{输出行（转义 &lt; &gt;）}}</pre>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gd-ambient, gd-snum, gd-eyebrow, gd-h2, gd-codebox, mt-l, gd-cmd, kw, st, fn, notes

---

## compare（对照卡）
数量：gd-glass=2
指纹：cards

用途：两列对照（如方案 A/B、前后、国内外）
适用 role：content。
内容约束：恰好 2 列；每列卡内说明 22-45 字（一句论断 + 一句展开）

```html
<section class="slide" data-layout="compare">
  <div class="gd-ambient"></div>
  <div class="gd-snum">{{编号}}</div>
  <p class="gd-eyebrow">{{引导语}}</p>
  <h2 class="gd-h2">{{对照标题}}</h2>
  <div class="gd-grid-3 mt-l" style="grid-template-columns:1fr 1fr">
    <div class="gd-glass"><span class="gd-tag">{{列 A}}</span><p>{{内容，22-45 字}}</p></div>
    <div class="gd-glass"><span class="gd-tag">{{列 B}}</span><p>{{内容，22-45 字}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gd-ambient, gd-snum, gd-eyebrow, gd-h2, gd-grid-3, mt-l, gd-glass, gd-tag, notes

---

## trend（四格趋势）
数量：gd-glass=4
指纹：chart

用途：四个趋势/发现并列（grid-4 玻璃卡）
适用 role：data。
内容约束：恰好 4 卡；每卡说明 22-40 字（一句论断 + 一句展开）

```html
<section class="slide" data-layout="trend">
  <div class="gd-ambient"></div>
  <div class="gd-snum">{{编号}}</div>
  <p class="gd-eyebrow">{{引导语}}</p>
  <h2 class="gd-h2">{{趋势标题}}</h2>
  <div class="gd-grid-4 mt-l">
    <div class="gd-glass"><span class="gd-tag">{{发现}}</span><p>{{内容，22-40 字}}</p></div>
    <div class="gd-glass"><span class="gd-tag">{{发现}}</span><p>{{内容，22-40 字}}</p></div>
    <div class="gd-glass"><span class="gd-tag">{{发现}}</span><p>{{内容，22-40 字}}</p></div>
    <div class="gd-glass"><span class="gd-tag">{{发现}}</span><p>{{内容，22-40 字}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gd-ambient, gd-snum, gd-eyebrow, gd-h2, gd-grid-4, mt-l, gd-glass, gd-tag, notes

---

## closing（收束）
指纹：hero

用途：研究收束：结论一句 + 后续动作
适用 role：thanks / cta。
内容约束：结论 12-20 字；lede 20-30 字

```html
<section class="slide" data-layout="closing">
  <div class="gd-ambient"></div>
  <div class="gd-snum">{{编号}}</div>
  <p class="gd-eyebrow">{{章节}}</p>
  <h1 class="gd-h1">{{结论一句，12-20 字}}</h1>
  <p class="gd-lede">{{后续动作/致谢，20-30 字}}</p>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, gd-ambient, gd-snum, gd-eyebrow, gd-h1, gd-lede, notes
