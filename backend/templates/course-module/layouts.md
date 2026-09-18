# 课程模块 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-course-module` 作用域前缀生效，骨架里已写全，照抄结构即可。
>
> **结构铁律**：非满版页（不带 `slide full` 的页）必须保持骨架的
> `<aside class="sidebar">…</aside><div class="main">…</div>` 包裹结构——
> 侧栏是本模板的身份（品牌 + 目标进度 + 关键词），丢了它页面就散架。
> {{占位符}} 只准换文本；骨架的层级与包裹关系不准增删。

---

## cover（课程封面）
指纹：hero

用途：开场页。大标题 + 一句话定位 + 课程信息胶囊。
适用 role：cover。
内容约束：标题 ≤14 字；lede 30-50 字（面向谁 + 学完能做什么）；3 个胶囊。

```html
<section class="slide full" data-layout="cover">
  <p class="kicker">{{课程标签，如 校园科普宣讲 · 10 分钟}}</p>
  <h1 class="h1 mt-s">{{课程标题}}</h1>
  <p class="lede mt-m" style="max-width:62ch">{{面向谁、学完能做什么，30-50 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="pill-academic">{{时长 / 讲次}}</span>
    <span class="pill-academic">{{适合谁 / 先修}}</span>
    <span class="pill-academic">{{形式，如 讲解 + 练习}}</span>
  </div>
  <div class="ghost-num">{{主视觉数字或符号：本课程最核心的一个数、编号或符号，≤6 字符}}</div>
  <div class="deck-footer"><span>{{讲师 / 课程名}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句：开场钩子 + 本节目标}}</div>
</section>
```

合法类名：slide, full, kicker, h1, mt-s, lede, mt-m, row, mt-l, pill-academic, ghost-num, deck-footer, slide-number, notes

---

## objectives（学习目标）
指纹：stack

用途：本模块学完后能做到的 3-4 条目标，每条带一句解释。
适用 role：toc。
内容约束：3-4 条；每条标题动词开头 ≤16 字 + 解释 18-36 字；侧栏目标列表与主区一一对应。

```html
<section class="slide" data-layout="objectives">
  <aside class="sidebar">
    <div class="brand">{{课程简称}}</div>
    <h5>Learning objectives</h5>
    <ul class="obj-list">
      <li class="current">{{目标 1（与主区第 1 条一致）}}</li>
      <li>{{目标 2}}</li>
      <li>{{目标 3}}</li>
      <li>{{目标 4}}</li>
    </ul>
    <h5>进度</h5>
    <p class="dim" style="font-size:16px">第 {{页码}} / {{总页数}} 页 · 约 {{分钟数}} 分钟</p>
  </aside>
  <div class="main">
    <p class="kicker">objectives</p>
    <h2 class="h2 mt-s">学完这一模块，你会</h2>
    <div class="stack mt-l">
      <div class="concept-box"><h4>① {{目标 1，动词开头 ≤16 字}}</h4><p class="dim">{{一句话解释，18-36 字}}</p></div>
      <div class="concept-box"><h4>② {{目标 2，动词开头 ≤16 字}}</h4><p class="dim">{{一句话解释，18-36 字}}</p></div>
      <div class="concept-box"><h4>③ {{目标 3，动词开头 ≤16 字}}</h4><p class="dim">{{一句话解释，18-36 字}}</p></div>
    </div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sidebar, brand, obj-list, current, main, kicker, h2, mt-s, stack, mt-l, concept-box, dim, notes

---

## concept（概念讲解）
指纹：cards

用途：核心概念 + 两个并列要点 + 一个易错点纠正。
适用 role：content。
内容约束：概念一句话 30-50 字；两张卡各 25-50 字解释 + 具体例子；callout 30-50 字。

```html
<section class="slide" data-layout="concept">
  <aside class="sidebar">
    <div class="brand">{{课程简称}}</div>
    <h5>Key terms</h5>
    <p class="dim" style="font-size:16px">{{术语 a · 术语 b · 术语 c}}</p>
  </aside>
  <div class="main">
    <p class="kicker">core concept</p>
    <h2 class="h2 mt-s">{{概念名，≤10 字}}</h2>
    <p class="lede mt-m">{{这个概念是什么，30-50 字}}</p>
    <div class="grid g2 mt-l">
      <div class="concept-box"><h4>{{要点 1 标题，≤10 字}}</h4><p class="dim">{{解释，25-50 字}}</p><p class="pill-academic">例如 <b>{{具体例子}}</b></p></div>
      <div class="concept-box"><h4>{{要点 2 标题，≤10 字}}</h4><p class="dim">{{解释，25-50 字}}</p><p class="pill-academic">例如 <b>{{具体例子}}</b></p></div>
    </div>
    <div class="callout"><b>易错点。</b>{{最常被理解错的地方 + 纠正，30-50 字}}</div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sidebar, brand, main, kicker, h2, mt-s, lede, mt-m, grid, g2, mt-l, concept-box, dim, pill-academic, callout, notes

---

## example（示例讲解）
指纹：code

用途：代码/操作示例（.code 块承载）+ 逐步拆解。
适用 role：code / content。
内容约束：代码 8-16 行；callout 拆解 40-80 字（逐步讲清发生了什么）。

```html
<section class="slide" data-layout="example">
  <aside class="sidebar">
    <div class="brand">{{课程简称}}</div>
    <h5>Try it yourself</h5>
    <p class="dim" style="font-size:16px">{{给学员的一句动手提示，15-30 字}}</p>
  </aside>
  <div class="main">
    <p class="kicker">worked example</p>
    <h2 class="h2 mt-s">{{示例标题，≤12 字}}</h2>
    <div class="code mt-m"><pre style="margin:0">{{示例代码 8-16 行（&lt; &gt; 要转义；.kw/.str/.cmt 标语法）}}</pre></div>
    <div class="callout">
      <b>逐行看。</b>{{这段代码在做什么，从输入到输出，40-80 字}}
    </div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sidebar, brand, main, kicker, h2, mt-s, code, mt-m, callout, kw, str, cmt, dim, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。一页只说一件事：进入第几节、这节讲什么。右下角巨大章节数字是视觉锚点。
适用 role：divider。
内容约束：标题 ≤12 字；lede 25-45 字（这一节回答的一个问题）；2 个看点胶囊。

```html
<section class="slide full" data-layout="divider">
  <p class="kicker">{{第 X 节 · 共 N 节}}</p>
  <h1 class="h1 mt-s">{{章节标题，≤12 字}}</h1>
  <p class="lede mt-m" style="max-width:56ch">{{这一节回答的一个问题，25-45 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="pill-academic">{{本节看点 1}}</span>
    <span class="pill-academic">{{本节看点 2}}</span>
  </div>
  <div class="ghost-num">{{X}}</div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, kicker, h1, mt-s, lede, mt-m, row, mt-l, pill-academic, ghost-num, notes

---

## stat（关键数字）
指纹：chart

用途：用 2-4 张大数字卡回答一个量化问题。有数字就该用它，别把数字埋进句子里。
适用 role：data / content。
内容约束：2-4 张卡；每卡数值 ≤6 字符 + 单位 + 指标名 ≤12 字 + 口径注释 12-30 字；页尾写数据来源。

```html
<section class="slide" data-layout="stat">
  <aside class="sidebar">
    <div class="brand">{{课程简称}}</div>
    <h5>怎么读</h5>
    <p class="dim" style="font-size:16px">{{口径说明或容易误读的点，20-45 字}}</p>
  </aside>
  <div class="main">
    <p class="kicker">key numbers</p>
    <h2 class="h2 mt-s">{{这些数字回答什么问题，≤14 字}}</h2>
    <div class="stat-row mt-l">
      <div class="stat"><div class="stat-v">{{数值}}<span class="stat-u">{{单位}}</span></div><div class="stat-l">{{指标名，≤12 字}}</div><p class="stat-note">{{口径/年份，12-30 字}}</p></div>
      <div class="stat"><div class="stat-v">{{数值}}<span class="stat-u">{{单位}}</span></div><div class="stat-l">{{指标名，≤12 字}}</div><p class="stat-note">{{口径/年份，12-30 字}}</p></div>
      <div class="stat"><div class="stat-v">{{数值}}<span class="stat-u">{{单位}}</span></div><div class="stat-l">{{指标名，≤12 字}}</div><p class="stat-note">{{口径/年份，12-30 字}}</p></div>
    </div>
    <p class="dim mt-m" style="font-size:17px">来源：{{数据来源}}</p>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sidebar, brand, main, kicker, h2, mt-s, stat-row, mt-l, stat, stat-v, stat-u, stat-l, stat-note, dim, mt-m, notes

---

## timeline（时间线 / 里程碑）
指纹：chart

用途：3-4 个节点的发展脉络（年份、阶段、步骤）。
适用 role：data / content。
内容约束：3-4 个节点；节点标题 ≤8 字 + 说明 20-40 字；callout 给一句关键解读 30-50 字。

```html
<section class="slide" data-layout="timeline">
  <aside class="sidebar">
    <div class="brand">{{课程简称}}</div>
    <h5>怎么看</h5>
    <p class="dim" style="font-size:16px">{{读这条时间线的方法，15-35 字}}</p>
  </aside>
  <div class="main">
    <p class="kicker">milestones</p>
    <h2 class="h2 mt-s">{{时间线标题，≤14 字}}</h2>
    <div class="timeline mt-l">
      <div class="tl-item"><div class="tl-t">{{节点 1，≤8 字}}</div><p class="tl-d">{{发生了什么 / 要发生什么，20-40 字}}</p></div>
      <div class="tl-item"><div class="tl-t">{{节点 2，≤8 字}}</div><p class="tl-d">{{说明，20-40 字}}</p></div>
      <div class="tl-item"><div class="tl-t">{{节点 3，≤8 字}}</div><p class="tl-d">{{说明，20-40 字}}</p></div>
      <div class="tl-item"><div class="tl-t">{{节点 4，≤8 字}}</div><p class="tl-d">{{说明，20-40 字}}</p></div>
    </div>
    <div class="callout"><b>关键解读。</b>{{这条线最重要的含义，30-50 字}}</div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sidebar, brand, main, kicker, h2, mt-s, timeline, mt-l, tl-item, tl-t, tl-d, callout, dim, notes

---

## exercise（随堂练习）
指纹：stack

用途：学员动手的任务说明（exercise 容器）+ 编号要求。
适用 role：content / cta。
内容约束：任务 ≤30 字；要求 3-4 条、每条 15-35 字；侧栏写时间与形式。

```html
<section class="slide" data-layout="exercise">
  <aside class="sidebar">
    <div class="brand">{{课程简称}}</div>
    <h5>Time</h5>
    <p class="dim" style="font-size:16px">{{时长 · 独立/组队}}</p>
  </aside>
  <div class="main">
    <p class="kicker">exercise</p>
    <h2 class="h2 mt-s">{{动手任务，≤14 字}}</h2>
    <p class="lede mt-m">{{任务描述，≤30 字}}</p>
    <div class="exercise mt-l">
      <p style="margin:0;font-size:20px;color:var(--text-1)"><b>要求</b></p>
      <ol style="color:var(--text-2);line-height:1.8;margin:10px 0 0">
        <li>{{要求 1，15-35 字}}</li>
        <li>{{要求 2，15-35 字}}</li>
        <li>{{要求 3，15-35 字}}</li>
      </ol>
    </div>
    <p class="dim mt-m" style="font-size:17px">卡住了？{{一句提示}}</p>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sidebar, brand, main, kicker, h2, mt-s, lede, mt-m, exercise, mt-l, dim, notes

---

## check（选择自测）
指纹：stack

用途：一道选择题（mcq 选项，correct 标正确项）。
适用 role：content。
内容约束：恰好 1 题；题干 ≤30 字；4 个选项、每项 ≤20 字 + 一句 12-30 字的对错解释。

```html
<section class="slide" data-layout="check">
  <aside class="sidebar">
    <div class="brand">{{课程简称}}</div>
    <h5>Self-assess</h5>
    <p class="dim" style="font-size:16px">应该能拿到 1/1。</p>
  </aside>
  <div class="main">
    <p class="kicker">check your understanding</p>
    <h2 class="h2 mt-s">{{题干，≤30 字}}</h2>
    <div class="stack mt-l">
      <div class="mcq"><div class="letter">A</div><div><b>{{选项 A，≤20 字}}</b><p class="dim" style="font-size:16px;margin:4px 0 0">{{对/错的一句话解释}}</p></div></div>
      <div class="mcq correct"><div class="letter">B</div><div><b>{{正确选项，≤20 字}}</b><p class="dim" style="font-size:16px;margin:4px 0 0"><b style="color:var(--accent)">✓ 正确。</b>{{为什么对}}</p></div></div>
      <div class="mcq"><div class="letter">C</div><div><b>{{选项 C，≤20 字}}</b><p class="dim" style="font-size:16px;margin:4px 0 0">{{错在哪}}</p></div></div>
      <div class="mcq"><div class="letter">D</div><div><b>{{选项 D，≤20 字}}</b><p class="dim" style="font-size:16px;margin:4px 0 0">{{错在哪}}</p></div></div>
    </div>
  </div>
  <div class="notes">{{讲稿：为什么对、其余错在哪}}</div>
</section>
```

合法类名：slide, sidebar, brand, main, kicker, h2, mt-s, stack, mt-l, mcq, correct, letter, dim, notes

---

## summary（模块小结）
指纹：hero

用途：小结（✓ 记忆点阵列）+ 下一模块预告。
适用 role：content / thanks。
内容约束：3-4 个记忆点、每卡 15-35 字；预告 16-30 字。

```html
<section class="slide full" data-layout="summary">
  <p class="kicker">summary · {{课程简称}}</p>
  <h1 class="h1 mt-s">{{小结标题，≤12 字}}</h1>
  <div class="grid g2 mt-l">
    <div class="concept-box"><h4>✓ {{记忆点 1，≤14 字}}</h4><p class="dim">{{15-35 字}}</p></div>
    <div class="concept-box"><h4>✓ {{记忆点 2，≤14 字}}</h4><p class="dim">{{15-35 字}}</p></div>
    <div class="concept-box"><h4>✓ {{记忆点 3，≤14 字}}</h4><p class="dim">{{15-35 字}}</p></div>
    <div class="concept-box"><h4>✓ {{记忆点 4，≤14 字}}</h4><p class="dim">{{15-35 字}}</p></div>
  </div>
  <div class="callout mt-l"><b>下一节。</b>{{预告，16-30 字}}</div>
  <div class="deck-footer"><span>{{课程名}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, kicker, h1, mt-s, grid, g2, mt-l, concept-box, dim, callout, deck-footer, slide-number, notes
