# 课程模块 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 XXX> 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-course-module` 作用域前缀生效，骨架里已写全，照抄结构即可。

---

## cover（课程封面）

用途：课程标题 + 侧栏课程导航出现
适用 role：cover。
内容约束：标题 ≤14 字；副标 ≤20 字

```html
<section class="slide" data-layout="cover">
  <p class="kicker">{{课程模块标签}}</p>
  <h1 class="h1 mt-s">{{课程标题}}</h1>
  <p class="dim mt-m">{{面向谁、学完能做什么}}</p>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h1, mt-s, dim, mt-m, notes

---

## objectives（学习目标）

用途：本模块学完后能做到的 3-4 条目标
适用 role：toc。
内容约束：3-4 条、每条动词开头 ≤20 字

```html
<section class="slide" data-layout="objectives">
  <p class="kicker">objectives</p>
  <h2 class="h2 mt-s">学完这一模块，你会</h2>
  <div class="stack mt-l">
    <div class="callout">{{目标 1（动词开头）}}</div>
    <div class="callout">{{目标 2}}</div>
    <div class="callout">{{目标 3}}</div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, mt-s, stack, mt-l, callout, notes

---

## concept（概念讲解）

用途：核心概念 + 要点解释（concept-box 承载定义）
适用 role：content。
内容约束：定义 ≤40 字；要点 2-4 条、每条 ≤24 字

```html
<section class="slide" data-layout="concept">
  <p class="kicker">concept</p>
  <h2 class="h2 mt-s">{{概念名}}</h2>
  <div class="concept-box mt-l">{{概念定义（一句话）}}</div>
  <div class="stack mt-l">
    <div class="callout">{{要点解释}}</div>
    <div class="callout">{{要点解释}}</div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, mt-s, concept-box, mt-l, stack, callout, notes

---

## example（示例讲解）

用途：代码/操作示例（.code 块承载，语法高亮可选）
适用 role：code / content。
内容约束：代码 ≤16 行；配 1-2 句解说

```html
<section class="slide" data-layout="example">
  <p class="kicker">example</p>
  <h2 class="h2 mt-s">{{示例标题}}</h2>
  <div class="code mt-m">
{{示例代码（可转义 &lt; &gt;）}}
  </div>
  <p class="dim mt-m">{{这段代码在做什么}}</p>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, mt-s, code, mt-m, dim, notes

---

## exercise（随堂练习）

用途：学员动手的任务说明（exercise 容器）
适用 role：content / cta。
内容约束：任务 ≤30 字；要求 2-3 条

```html
<section class="slide" data-layout="exercise">
  <p class="kicker">exercise</p>
  <h2 class="h2 mt-s">动手练一练</h2>
  <div class="exercise mt-l">
    <p>{{任务描述}}</p>
    <p class="dim">{{要求/提示}}</p>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, mt-s, exercise, mt-l, dim, notes

---

## check（选择自测）

用途：一道选择题（mcq 选项，correct 标正确项）
适用 role：content。
内容约束：恰好 1 题；题干 ≤30 字；4 个选项、每项 ≤20 字

```html
<section class="slide" data-layout="check">
  <p class="kicker">quiz</p>
  <h2 class="h2 mt-s">{{题干}}</h2>
  <div class="stack mt-l">
    <div class="mcq correct">{{正确选项}}</div>
    <div class="mcq">{{干扰项}}</div>
    <div class="mcq">{{干扰项}}</div>
    <div class="mcq">{{干扰项}}</div>
  </div>
  <div class="notes">{{讲稿：为什么对、其余错在哪}}</div>
</section>
```

合法类名：slide, kicker, h2, mt-s, stack, mt-l, mcq, correct, notes

---

## summary（模块小结）

用途：小结 + 下一模块预告
适用 role：content / thanks。
内容约束：3 条记忆点、每条 ≤20 字；预告 ≤16 字

```html
<section class="slide" data-layout="summary">
  <p class="kicker">summary</p>
  <h2 class="h2 mt-s">{{小结标题}}</h2>
  <div class="stack mt-l">
    <div class="callout">{{记忆点 1}}</div>
    <div class="callout">{{记忆点 2}}</div>
    <div class="callout">{{记忆点 3}}</div>
  </div>
  <p class="dim mt-l">下一模块：{{预告}}</p>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, mt-s, stack, mt-l, callout, dim, notes
