# 静默极简 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `mq-*` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-minimal-quiet` 作用域前缀生效，骨架里已写全，照抄结构即可。
>
> **安静纪律**：颜色是稀缺资源——粉彩只以「色卡 + 深字」成对出现在胶囊与状态标记上；
> 正文永远炭黑 #2F3437 系（token），不准写纯黑，不准加渐变、大阴影和药丸形容器。

---

## cover（复盘封面）
指纹：hero

用途：开场页。衬线大标题 + 一句话定位 + 周期信息胶囊。
适用 role：cover。
内容约束：标题 ≤12 字；lede 25-45 字（这份复盘是什么、给谁看）；3 个 mono 胶囊。

```html
<section class="slide full" data-layout="cover">
  <p class="kicker">{{性质 · 周期，如 个人季度复盘 · 2026 Q2}}</p>
  <h1 class="h1 mt-s">{{主标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:56ch">{{这份复盘是什么、只写给谁看，25-45 字}}</p>
  <div class="row mt-l" style="gap:14px">
    <span class="mq-chip mono">{{周期区间，如 04.01 – 06.30}}</span>
    <span class="mq-chip mono">{{身份 · 姓名}}</span>
    <span class="mq-chip mono">{{天数或期数，如 共 91 天}}</span>
  </div>
  <div class="deck-footer"><span class="mono">{{署名 · 复盘编号}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句：这份复盘的规矩（只摆事实）}}</div>
</section>
```

合法类名：slide, full, kicker, h1, mt-s, lede, mt-m, row, mt-l, mq-chip, mono, deck-footer, slide-number, notes

---

## agenda（复盘路线）
指纹：table

用途：目录。编号行列表，每行条目标题 + 一句看点 + 预计时长。
适用 role：toc。
内容约束：3-5 行；每行标题 4-10 字 + 看点 12-24 字；时长列写 `~Nmin`。

```html
<section class="slide" data-layout="agenda">
  <p class="kicker">{{引导语，如 index · 目录}}</p>
  <h2 class="h2 mt-s">{{议程标题，≤12 字}}</h2>
  <div class="mq-toc mt-l">
    <div class="mq-toc-row"><span class="mq-toc-n mono">{{01}}</span><span class="mq-toc-t">{{条目标题，4-10 字}}</span><span class="mq-toc-d">{{一句看点，12-24 字}}</span><span class="mq-toc-n mono">~{{时长}}min</span></div>
    <div class="mq-toc-row"><span class="mq-toc-n mono">{{02}}</span><span class="mq-toc-t">{{条目标题，4-10 字}}</span><span class="mq-toc-d">{{一句看点，12-24 字}}</span><span class="mq-toc-n mono">~{{时长}}min</span></div>
    <div class="mq-toc-row"><span class="mq-toc-n mono">{{03}}</span><span class="mq-toc-t">{{条目标题，4-10 字}}</span><span class="mq-toc-d">{{一句看点，12-24 字}}</span><span class="mq-toc-n mono">~{{时长}}min</span></div>
    <div class="mq-toc-row"><span class="mq-toc-n mono">{{04}}</span><span class="mq-toc-t">{{条目标题，4-10 字}}</span><span class="mq-toc-d">{{一句看点，12-24 字}}</span><span class="mq-toc-n mono">~{{时长}}min</span></div>
  </div>
  <div class="deck-footer"><span class="mono">{{署名}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, mt-s, mq-toc, mt-l, mq-toc-row, mq-toc-n, mq-toc-t, mq-toc-d, mono, deck-footer, slide-number, notes

---

## okr-table（目标进度表）
指纹：table

用途：3-4 行目标进度。每行 = 目标 + 关键结果口径 + mono 百分比 + 状态胶囊。
适用 role：data / content。
内容约束：3-4 行；目标 ≤12 字 + 关键结果 18-36 字 + 状态 ≤6 字；页尾一句口径说明 20-45 字。

```html
<section class="slide" data-layout="okr-table">
  <p class="kicker">{{引导语，如 objectives · 目标盘点}}</p>
  <h2 class="h2 mt-s">{{总览标题，≤14 字}}</h2>
  <div class="mq-table mt-l">
    <div class="mq-tr"><span class="mq-tr-t">{{目标，≤12 字}}</span><span class="mq-tr-k">{{关键结果与口径，18-36 字}}</span><span class="mq-tr-num mono">{{进度，如 100%}}</span><span class="mq-chip mq-chip-green">{{状态，≤6 字}}</span></div>
    <div class="mq-tr"><span class="mq-tr-t">{{目标，≤12 字}}</span><span class="mq-tr-k">{{关键结果与口径，18-36 字}}</span><span class="mq-tr-num mono">{{进度百分比}}</span><span class="mq-chip mq-chip-yellow">{{状态，≤6 字}}</span></div>
    <div class="mq-tr"><span class="mq-tr-t">{{目标，≤12 字}}</span><span class="mq-tr-k">{{关键结果与口径，18-36 字}}</span><span class="mq-tr-num mono">{{进度百分比}}</span><span class="mq-chip mq-chip-green">{{状态，≤6 字}}</span></div>
  </div>
  <p class="dim mt-m" style="font-size:17px">口径：{{进度怎么折算、数据来自哪，20-45 字}}</p>
  <div class="deck-footer"><span class="mono">{{署名}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿：延期/超额的那一行发生了什么}}</div>
</section>
```

合法类名：slide, kicker, h2, mt-s, mq-table, mt-l, mq-tr, mq-tr-t, mq-tr-k, mq-tr-num, mono, mq-chip, mq-chip-green, mq-chip-yellow, dim, mt-m, deck-footer, slide-number, notes

---

## metrics（关键数字）
指纹：chart

用途：2-4 张大数字卡回答一个量化问题。数字是这页唯一的主角，别埋进句子。
适用 role：data / content。
内容约束：2-4 卡；数值 ≤6 字符 + 指标名 ≤12 字 + 口径 12-30 字；页尾标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="kicker">{{引导语，如 numbers · 关键数字}}</p>
  <h2 class="h2 mt-s">{{这些数字回答什么，≤14 字}}</h2>
  <div class="mq-stat-row mt-l">
    <div class="mq-stat"><div class="mq-stat-v">{{数值 ≤6 字符}}<span class="mq-stat-u">{{单位}}</span></div><div class="mq-stat-l">{{指标名，≤12 字}}</div><p class="mq-stat-note">{{口径 / 区间，12-30 字}}</p></div>
    <div class="mq-stat"><div class="mq-stat-v">{{数值 ≤6 字符}}<span class="mq-stat-u">{{单位}}</span></div><div class="mq-stat-l">{{指标名，≤12 字}}</div><p class="mq-stat-note">{{口径 / 区间，12-30 字}}</p></div>
    <div class="mq-stat"><div class="mq-stat-v">{{数值 ≤6 字符}}<span class="mq-stat-u">{{单位}}</span></div><div class="mq-stat-l">{{指标名，≤12 字}}</div><p class="mq-stat-note">{{口径 / 区间，12-30 字}}</p></div>
    <div class="mq-stat"><div class="mq-stat-v">{{数值 ≤6 字符}}<span class="mq-stat-u">{{单位}}</span></div><div class="mq-stat-l">{{指标名，≤12 字}}</div><p class="mq-stat-note">{{口径 / 区间，12-30 字}}</p></div>
  </div>
  <p class="dim mt-m" style="font-size:17px">来源：{{数据出处与统计区间，15-40 字}}</p>
  <div class="deck-footer"><span class="mono">{{署名}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, mt-s, mq-stat-row, mt-l, mq-stat, mq-stat-v, mq-stat-u, mq-stat-l, mq-stat-note, dim, mt-m, mono, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。一页只说一件事：进入第几部分、这部分回答什么。右下描边数字是唯一装饰。
适用 role：divider。
内容约束：标题 ≤12 字；lede 25-45 字（这一部分回答的一个问题）；2 个看点胶囊。

```html
<section class="slide full" data-layout="divider">
  <p class="kicker">{{进度，如 part 2 · 一个案例}}</p>
  <h1 class="h1 mt-s">{{章节标题，≤12 字}}</h1>
  <p class="lede mt-m" style="max-width:52ch">{{这一部分回答的一个问题，25-45 字}}</p>
  <div class="row mt-l" style="gap:14px">
    <span class="mq-chip">{{看点 1，≤10 字}}</span>
    <span class="mq-chip">{{看点 2，≤10 字}}</span>
  </div>
  <div class="mq-ghost">{{章节数字，1 字符}}</div>
  <div class="deck-footer"><span class="mono">{{署名}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, kicker, h1, mt-s, lede, mt-m, row, mt-l, mq-chip, mq-ghost, mono, deck-footer, slide-number, notes

---

## case-split（案例左文右摘录）
指纹：split

用途：左边把一个案例讲透（叙述 + 补充 + 数据胶囊），右边 mono 日志/记录摘录面板佐证。
适用 role：content。
内容约束：左 lede 40-70 字 + dim 补充 20-40 字 + 2-3 个 mono 胶囊；右摘录 6-10 行、每行 ≤24 字。

```html
<section class="slide" data-layout="case-split">
  <p class="kicker">{{引导语，如 case · 案例代号}}</p>
  <h2 class="h2 mt-s">{{案例标题，≤14 字}}</h2>
  <div class="grid g2 mt-l" style="gap:56px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做的、为什么这样做，40-70 字}}</p>
      <p class="dim mt-m" style="font-size:18px">{{补充：代价、转折或反例，20-40 字}}</p>
      <div class="row mt-l" style="gap:12px">
        <span class="mq-chip mono">{{关键数字 1，≤8 字}}</span>
        <span class="mq-chip mono">{{关键数字 2，≤8 字}}</span>
        <span class="mq-chip mono">{{关键数字 3，≤8 字}}</span>
      </div>
    </div>
    <div class="mq-panel">
      <p class="mq-panel-h mono">{{摘录标题，如 # week 3-6 · 工作日志}}</p>
<pre class="mono" style="margin:0">{{摘录 6-10 行，每行 ≤24 字（&lt; &gt; &amp; 要转义）}}</pre>
    </div>
  </div>
  <div class="deck-footer"><span class="mono">{{署名}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, mt-s, grid, g2, mt-l, lede, dim, mt-m, row, mq-chip, mono, mq-panel, mq-panel-h, deck-footer, slide-number, notes

---

## reflections（三卡经验）
数量：mq-card=3
指纹：cards

用途：恰好三张白卡。每张：定性粉彩胶囊（有效 / 待改进 / 保持……）+ 经验标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；胶囊 ≤4 字 + 标题 ≤10 字 + 说明 22-45 字（一句论断 + 一句展开）。

```html
<section class="slide" data-layout="reflections">
  <p class="kicker">{{引导语，如 retro · 经验}}</p>
  <h2 class="h2 mt-s">{{标题，≤14 字}}</h2>
  <div class="grid g3 mt-l" style="gap:28px">
    <div class="mq-card"><span class="mq-chip mq-chip-green">{{定性词，≤4 字}}</span><h4 class="mt-m">{{经验标题，≤10 字}}</h4><p class="dim mt-s" style="font-size:18px">{{说明：一句论断 + 一句展开，22-45 字}}</p></div>
    <div class="mq-card"><span class="mq-chip mq-chip-yellow">{{定性词，≤4 字}}</span><h4 class="mt-m">{{经验标题，≤10 字}}</h4><p class="dim mt-s" style="font-size:18px">{{说明：一句论断 + 一句展开，22-45 字}}</p></div>
    <div class="mq-card"><span class="mq-chip mq-chip-blue">{{定性词，≤4 字}}</span><h4 class="mt-m">{{经验标题，≤10 字}}</h4><p class="dim mt-s" style="font-size:18px">{{说明：一句论断 + 一句展开，22-45 字}}</p></div>
  </div>
  <div class="deck-footer"><span class="mono">{{署名}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, h2, mt-s, grid, g3, mt-l, mq-card, mq-chip, mq-chip-green, mq-chip-yellow, mq-chip-blue, h4, mt-m, dim, mt-s, mono, deck-footer, slide-number, notes

---

## quote-note（笔记摘句）
指纹：quote

用途：整页只放一句自己写下的话。衬线大字 + 出处 + 一句「后来怎么用上」。没有真实摘录就换版式。
适用 role：quote。
内容约束：引文 15-40 字；出处 ≤24 字（日期与来源）；注释 20-45 字。

```html
<section class="slide" data-layout="quote-note">
  <p class="kicker">{{引导语，如 notebook · 5 月 12 日}}</p>
  <p class="mq-quote mt-l">{{引文，15-40 字}}</p>
  <div class="divider mt-l"></div>
  <p class="mq-quote-src mono mt-m">—— {{出处：日期与笔记来源，≤24 字}}</p>
  <p class="dim mt-m" style="font-size:18px;max-width:56ch">{{为什么记下它 / 后来怎么用上，20-45 字}}</p>
  <div class="deck-footer"><span class="mono">{{署名}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, kicker, mq-quote, mt-l, divider, mq-quote-src, mono, mt-m, dim, deck-footer, slide-number, notes

---

## closing（下期三件事）
指纹：hero

用途：收尾页。编号列出下周期只做的那几件事 + 一句为什么。垂直居中，安静收束。
适用 role：content / cta / thanks。
内容约束：3-4 条、每条 12-24 字动词开头；lede 20-45 字。

```html
<section class="slide full" data-layout="closing">
  <p class="kicker">{{下周期，如 next · Q3}}</p>
  <h1 class="h1 mt-s">{{收束标题，≤12 字}}</h1>
  <p class="lede mt-m" style="max-width:52ch">{{为什么是这几件，20-45 字}}</p>
  <div class="mq-toc mt-l">
    <div class="mq-toc-row"><span class="mq-toc-n mono">{{01}}</span><span class="mq-toc-t">{{事项，12-24 字，动词开头}}</span></div>
    <div class="mq-toc-row"><span class="mq-toc-n mono">{{02}}</span><span class="mq-toc-t">{{事项，12-24 字，动词开头}}</span></div>
    <div class="mq-toc-row"><span class="mq-toc-n mono">{{03}}</span><span class="mq-toc-t">{{事项，12-24 字，动词开头}}</span></div>
  </div>
  <div class="deck-footer"><span class="mono">{{署名}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿：收束一句，不加空话}}</div>
</section>
```

合法类名：slide, full, kicker, h1, mt-s, lede, mt-m, mq-toc, mt-l, mq-toc-row, mq-toc-n, mq-toc-t, mono, deck-footer, slide-number, notes
