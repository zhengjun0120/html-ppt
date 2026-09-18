# 软糖结构 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `sp-*` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-soft-pastel` 作用域前缀生效，骨架里已写全，照抄结构即可。
>
> **结构铁律（双层表圈）**：凡是「卡」，一律 `.sp-shell`（灰壳：4% 墨底 + 发丝边 + 8px 内衬 + 32px 圆角）
> 包 `.sp-core`（白芯：24px 圆角）。不准把白卡或 1px 灰边框直接裸在银灰底上——浮在底上的
> 只允许超大标题、药丸和时间线圆点。

---

## cover（发布封面）
指纹：hero

用途：开场页。眉标药丸 + 超大标题 + 一句话定位 + 粉彩信息药丸。
适用 role：cover。
内容约束：标题 ≤12 字；lede 20-45 字；3 个信息药丸各 ≤10 字。

```html
<section class="slide full" data-layout="cover">
  <p class="sp-eyebrow">{{品牌 · 事件，如 山雾 SÁNWÙ · 新品发布}}</p>
  <h1 class="h1 mt-m">{{主标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:54ch">{{一句话定位，20-45 字}}</p>
  <div class="row mt-l" style="gap:14px">
    <span class="sp-pill sp-pill-pink">{{信息 1，≤10 字}}</span>
    <span class="sp-pill sp-pill-blue">{{信息 2，≤10 字}}</span>
    <span class="sp-pill sp-pill-yellow">{{信息 3，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{品牌署名 · 季节}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, sp-eyebrow, h1, mt-m, lede, mt-l, row, sp-pill, sp-pill-pink, sp-pill-blue, sp-pill-yellow, deck-footer, slide-number, notes

---

## launch-steps（发布节奏）
指纹：table

用途：一块双壳面板里放 3-4 行节奏：药丸圆号 + 节点（事项 · 日期）+ 一句说明。
适用 role：toc / content。
内容约束：3-4 行；节点 ≤14 字 + 说明 12-24 字。

```html
<section class="slide" data-layout="launch-steps">
  <p class="sp-eyebrow">{{引导语，如 发布节奏}}</p>
  <h2 class="h2 mt-m">{{标题，≤14 字}}</h2>
  <div class="sp-shell mt-l">
    <div class="sp-core">
      <div class="sp-step"><span class="sp-step-n">{{01}}</span><span class="sp-step-t">{{节点 · 日期，≤14 字}}</span><span class="sp-step-d">{{说明，12-24 字}}</span></div>
      <div class="sp-step"><span class="sp-step-n">{{02}}</span><span class="sp-step-t">{{节点 · 日期，≤14 字}}</span><span class="sp-step-d">{{说明，12-24 字}}</span></div>
      <div class="sp-step"><span class="sp-step-n">{{03}}</span><span class="sp-step-t">{{节点 · 日期，≤14 字}}</span><span class="sp-step-d">{{说明，12-24 字}}</span></div>
      <div class="sp-step"><span class="sp-step-n">{{04}}</span><span class="sp-step-t">{{节点 · 日期，≤14 字}}</span><span class="sp-step-d">{{说明，12-24 字}}</span></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{品牌署名}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sp-eyebrow, h2, mt-m, sp-shell, mt-l, sp-core, sp-step, sp-step-n, sp-step-t, sp-step-d, deck-footer, slide-number, notes

---

## warmup-stats（预热数据）
指纹：chart

用途：三张双壳数字卡。数字用超大加粗字重说话，口径写进卡内，来源写在页尾。
适用 role：data / content。
内容约束：恰好 3 卡；数值 ≤6 字符 + 指标名 ≤10 字 + 口径 12-30 字；必须标来源。

```html
<section class="slide" data-layout="warmup-stats">
  <p class="sp-eyebrow">{{数据语境，如 预热数据 · 截至 9.15}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤14 字}}</h2>
  <div class="grid g3 mt-l" style="gap:28px">
    <div class="sp-shell"><div class="sp-core"><div class="sp-stat-v">{{数值 ≤6 字符}}<span class="sp-stat-u">{{单位}}</span></div><div class="sp-stat-l">{{指标名，≤10 字}}</div><p class="sp-stat-note">{{口径，12-30 字}}</p></div></div>
    <div class="sp-shell"><div class="sp-core"><div class="sp-stat-v">{{数值 ≤6 字符}}<span class="sp-stat-u">{{单位}}</span></div><div class="sp-stat-l">{{指标名，≤10 字}}</div><p class="sp-stat-note">{{口径，12-30 字}}</p></div></div>
    <div class="sp-shell"><div class="sp-core"><div class="sp-stat-v">{{数值 ≤6 字符}}<span class="sp-stat-u">{{单位}}</span></div><div class="sp-stat-l">{{指标名，≤10 字}}</div><p class="sp-stat-note">{{口径，12-30 字}}</p></div></div>
  </div>
  <p class="dim mt-m" style="font-size:17px">来源：{{出处与统计区间，15-40 字}}</p>
  <div class="deck-footer"><span>{{品牌署名}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sp-eyebrow, h2, mt-m, grid, g3, mt-l, sp-shell, sp-core, sp-stat-v, sp-stat-u, sp-stat-l, sp-stat-note, dim, mt-m, deck-footer, slide-number, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。超大标题 + 这一部分回答的一个问题 + 两个看点药丸。
适用 role：divider。
内容约束：标题 ≤12 字；lede 25-45 字；2 个看点药丸各 ≤10 字。

```html
<section class="slide full" data-layout="divider">
  <p class="sp-eyebrow">{{进度，如 part 2 · 产品本身}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤12 字}}</h1>
  <p class="lede mt-m" style="max-width:52ch">{{这一部分回答的一个问题，25-45 字}}</p>
  <div class="row mt-l" style="gap:14px">
    <span class="sp-pill sp-pill-pink">{{看点 1，≤10 字}}</span>
    <span class="sp-pill sp-pill-blue">{{看点 2，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{品牌署名}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, sp-eyebrow, h1, mt-m, lede, mt-l, row, sp-pill, sp-pill-pink, sp-pill-blue, deck-footer, slide-number, notes

---

## feature-cards（三壳卖点）
指纹：cards

用途：恰好三张双壳卡。每张：粉彩药丸段位（带时间窗）+ 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；药丸 ≤10 字 + 小标题 ≤10 字 + 说明 22-45 字。

```html
<section class="slide" data-layout="feature-cards">
  <p class="sp-eyebrow">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:28px">
    <div class="sp-shell"><div class="sp-core"><span class="sp-pill sp-pill-pink">{{段位 · 时间窗，≤10 字}}</span><h4 class="mt-m">{{小标题，≤10 字}}</h4><p class="dim mt-s" style="font-size:18px">{{说明：一句论断 + 一句展开，22-45 字}}</p></div></div>
    <div class="sp-shell"><div class="sp-core"><span class="sp-pill sp-pill-blue">{{段位 · 时间窗，≤10 字}}</span><h4 class="mt-m">{{小标题，≤10 字}}</h4><p class="dim mt-s" style="font-size:18px">{{说明：一句论断 + 一句展开，22-45 字}}</p></div></div>
    <div class="sp-shell"><div class="sp-core"><span class="sp-pill sp-pill-yellow">{{段位 · 时间窗，≤10 字}}</span><h4 class="mt-m">{{小标题，≤10 字}}</h4><p class="dim mt-s" style="font-size:18px">{{说明：一句论断 + 一句展开，22-45 字}}</p></div></div>
  </div>
  <div class="deck-footer"><span>{{品牌署名}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sp-eyebrow, h2, mt-m, grid, g3, mt-l, sp-shell, sp-core, sp-pill, sp-pill-pink, sp-pill-blue, sp-pill-yellow, h4, mt-m, dim, mt-s, deck-footer, slide-number, notes

---

## ritual-split（场景左文右步骤）
指纹：split

用途：左边把使用场景讲透（lede + 补充 + 参数药丸），右边三张迷你双壳步骤卡。
适用 role：content。
内容约束：左 lede 40-70 字 + dim 补充 20-40 字 + 2 个参数药丸；右 3 步、每步 12-24 字。

```html
<section class="slide" data-layout="ritual-split">
  <p class="sp-eyebrow">{{引导语，如 使用仪式}}</p>
  <h2 class="h2 mt-m">{{标题，≤14 字}}</h2>
  <div class="grid g2 mt-l" style="gap:56px;align-items:start">
    <div>
      <p class="lede">{{叙述：场景里怎么做、为什么，40-70 字}}</p>
      <p class="dim mt-m" style="font-size:18px">{{补充：判断标准或代价，20-40 字}}</p>
      <div class="row mt-l" style="gap:12px">
        <span class="sp-pill sp-pill-blue">{{关键参数 1，≤8 字}}</span>
        <span class="sp-pill sp-pill-yellow">{{关键参数 2，≤8 字}}</span>
      </div>
    </div>
    <div class="stack">
      <div class="sp-shell"><div class="sp-core sp-core-row"><span class="sp-pill sp-pill-pink">{{第 1 步}}</span><p class="sp-mini-t">{{一句做法，12-24 字}}</p></div></div>
      <div class="sp-shell"><div class="sp-core sp-core-row"><span class="sp-pill sp-pill-pink">{{第 2 步}}</span><p class="sp-mini-t">{{一句做法，12-24 字}}</p></div></div>
      <div class="sp-shell"><div class="sp-core sp-core-row"><span class="sp-pill sp-pill-pink">{{第 3 步}}</span><p class="sp-mini-t">{{一句做法，12-24 字}}</p></div></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{品牌署名}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sp-eyebrow, h2, mt-m, grid, g2, mt-l, lede, dim, mt-m, row, mt-l, sp-pill, sp-pill-blue, sp-pill-yellow, stack, sp-shell, sp-core, sp-core-row, sp-pill-pink, sp-mini-t, deck-footer, slide-number, notes

---

## launch-timeline（里程碑时间线）
指纹：chart

用途：3-5 个节点的横向时间线：药丸圆点 + 时间点 + 事件。
适用 role：content / data。
内容约束：3-5 节点；时间点 ≤10 字符 + 事件 12-24 字；页尾一句口径 15-40 字。

```html
<section class="slide" data-layout="launch-timeline">
  <p class="sp-eyebrow">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤14 字}}</h2>
  <div class="sp-tl mt-l">
    <div class="sp-tl-item"><div class="sp-tl-dot"></div><div class="sp-tl-t">{{时间点，≤10 字符}}</div><p class="sp-tl-d">{{事件，12-24 字}}</p></div>
    <div class="sp-tl-item"><div class="sp-tl-dot"></div><div class="sp-tl-t">{{时间点，≤10 字符}}</div><p class="sp-tl-d">{{事件，12-24 字}}</p></div>
    <div class="sp-tl-item"><div class="sp-tl-dot"></div><div class="sp-tl-t">{{时间点，≤10 字符}}</div><p class="sp-tl-d">{{事件，12-24 字}}</p></div>
    <div class="sp-tl-item"><div class="sp-tl-dot"></div><div class="sp-tl-t">{{时间点，≤10 字符}}</div><p class="sp-tl-d">{{事件，12-24 字}}</p></div>
  </div>
  <p class="dim mt-m" style="font-size:17px">{{一句读法或口径，15-40 字}}</p>
  <div class="deck-footer"><span>{{品牌署名}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sp-eyebrow, h2, mt-m, sp-tl, mt-l, sp-tl-item, sp-tl-dot, sp-tl-t, sp-tl-d, dim, mt-m, deck-footer, slide-number, notes

---

## voice-quote（用户原声）
指纹：quote

用途：整页一句用户原话。超大加粗引文 + 出处 + 支撑药丸与口径。没有真实引文就换版式。
适用 role：quote。
内容约束：引文 12-30 字；出处 ≤24 字且真实；支撑口径 20-45 字。

```html
<section class="slide" data-layout="voice-quote">
  <p class="sp-eyebrow">{{语境，如 试香团盲测 · 9.14}}</p>
  <p class="sp-quote mt-l">「{{引文，12-30 字}}」</p>
  <p class="sp-src mt-m">—— {{出处：人与来源，≤24 字}}</p>
  <div class="row mt-l" style="gap:14px">
    <span class="sp-pill sp-pill-pink">{{支撑数字 1，≤8 字}}</span>
    <span class="sp-pill sp-pill-blue">{{支撑数字 2，≤8 字}}</span>
  </div>
  <p class="dim mt-m" style="font-size:18px">{{支撑口径一句，20-45 字}}</p>
  <div class="deck-footer"><span>{{品牌署名}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, sp-eyebrow, sp-quote, mt-l, sp-src, mt-m, row, mt-l, sp-pill, sp-pill-pink, sp-pill-blue, dim, deck-footer, slide-number, notes

---

## closing（发布收尾）
指纹：hero

用途：收尾页。超大标题 + 一句提醒 + 药丸按钮（右端嵌套圆位图标）+ 次级信息药丸。垂直居中。
适用 role：content / cta / thanks。
内容约束：标题 ≤12 字；lede 20-45 字；按钮 ≤8 字 + 次级药丸 ≤10 字。

```html
<section class="slide full" data-layout="closing">
  <p class="sp-eyebrow">{{提醒语境，如 首发提醒}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤12 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:52ch">{{一句收束提醒，20-45 字}}</p>
  <div class="row mt-l" style="gap:20px">
    <span class="sp-btn">{{按钮文案，≤8 字}}<span class="sp-orb">→</span></span>
    <span class="sp-pill sp-pill-yellow">{{次级信息，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{品牌署名}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, sp-eyebrow, h1, mt-m, lede, mt-l, row, sp-btn, sp-orb, sp-pill, sp-pill-yellow, deck-footer, slide-number, notes
