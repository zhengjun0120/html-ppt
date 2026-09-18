# presenter-cards · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-presenter-mode-reveal` 作用域前缀生效，骨架里已写全，照抄结构即可。
>
> **结构铁律**：每页末尾的 `<div class="notes">` 是本模板的灵魂——150-300 字口语逐字稿，
> 观众不可见、演讲者视图可见。写页面先想逐字稿：核心点加粗、过渡句成段、数据和名字列清楚。
> 丢了 notes，这份模板就没有存在的意义。

通用约定：
- 类名分两族：**base 原语**（`slide kicker h1 h2 h4 lede dim mono stack row grid g2 g3 card card-accent mt-s mt-m mt-l tc center deck-footer slide-number notes`）与**模板专属**（`speaker av agenda-row feature-row rule-row num accent code-block comment cmd flag blue green orange purple red`）。
- 颜色一律来自主题 token（`var(--accent)` 等）与登记的序号修饰类（`blue green orange purple red`），页面不写死色值。
- 每页讲稿写在页面末尾：`<div class="notes">…</div>`，150-300 字、口语、读一遍不拗口。
- **逐字稿纪律**：notes 不是讲稿是提示信号——少于 150 字提示不够，多于 300 字来不及看。

---

## cover（封面）
指纹：hero

用途：开场页。大标题（渐变点睛词）+ 一句话定位 + 讲者行 + 页脚话题标签。
适用 role：cover。
内容约束：主标题两行内、每行 ≤12 字；lede 20-30 字；讲者行 = 名字 + 身份 · 时长；页脚是话题标签。

```html
<section class="slide" data-layout="cover">
  <p class="kicker">{{主题标签，如 tech-talk / 2026-09}}</p>
  <h1 class="h1">{{主标题，两行用 <br>；点睛词用 <span class="accent">accent 色</span>}}</h1>
  <p class="lede mt-m">{{一句话定位，20-30 字}}</p>
  <div class="speaker"><div class="av"></div><div><b>{{讲者名}}</b><span>{{身份 · 时长}}</span></div></div>
  <div class="deck-footer"><span class="mono">{{话题标签，如 #tech-talk}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{逐字稿 150-300 字：开场钩子 + 自我介绍 + 今天的目标}}</div>
</section>
```

合法类名：slide, kicker, h1, accent, lede, mt-m, speaker, av, deck-footer, mono, slide-number, notes

---

## agenda（议程）
指纹：table

用途：目录/议程。编号行卡片，每行标题 + 预计时长。
适用 role：toc。
内容约束：4-6 行；每行标题 10-20 字 + 时长（`~Nmin`，写进 d 列）。

```html
<section class="slide" data-layout="agenda">
  <p class="kicker">agenda</p>
  <h2 class="h2">{{议程标题，如：今天要讲的几件事}}</h2>
  <div class="stack mt-l">
    <div class="agenda-row"><span class="num">01</span><span class="t">{{条目标题，10-20 字}}</span><span class="d">~{{时长}}min</span></div>
    <div class="agenda-row"><span class="num">02</span><span class="t">{{条目标题，10-20 字}}</span><span class="d">~{{时长}}min</span></div>
    <div class="agenda-row"><span class="num">03</span><span class="t">{{条目标题，10-20 字}}</span><span class="d">~{{时长}}min</span></div>
    <div class="agenda-row"><span class="num">04</span><span class="t">{{条目标题，10-20 字}}</span><span class="d">~{{时长}}min</span></div>
  </div>
  <div class="notes">{{逐字稿 150-300 字：逐条一句预告 + 时间分配逻辑}}</div>
</section>
```

合法类名：slide, kicker, h2, stack, mt-l, agenda-row, num, t, d, notes

---

## cards-3（三卡论点）
指纹：cards

用途：三个并列论点、三类对比、三块拆解。这是本模板的"内容主力"版式。
适用 role：content。
内容约束：恰好 3 卡；卡标题 4-12 字；说明 25-45 字（一句论断 + 一句展开，句号结尾）。

```html
<section class="slide" data-layout="cards-3">
  <p class="kicker">// {{引导语，如 part 01 · problem}}</p>
  <h2 class="h2">{{标题：一句话点破这页在讲什么，关键词用 <span class="accent">accent</span>}}</h2>
  <div class="grid g3 mt-l">
    <div class="card card-accent"><h4>{{卡标题，4-12 字}}</h4><p class="dim">{{说明，25-45 字}}</p></div>
    <div class="card card-accent"><h4>{{卡标题，4-12 字}}</h4><p class="dim">{{说明，25-45 字}}</p></div>
    <div class="card card-accent"><h4>{{卡标题，4-12 字}}</h4><p class="dim">{{说明，25-45 字}}</p></div>
  </div>
  <div class="notes">{{逐字稿 150-300 字：逐卡展开 + 三卡之间的关系}}</div>
</section>
```

合法类名：slide, kicker, h2, accent, grid, g3, mt-l, card, card-accent, h4, dim, notes

---

## feature-grid（特性两列）
指纹：cards

用途：4 个特性/模块/组件的两列清单（feature-row：彩色序号 + 名称 + 一句说明）。
适用 role：content / data。
内容约束：恰好 4 条、g2 两列各 2 条；名称 ≤10 字；说明 18-30 字；序号颜色按 blue/green/orange/purple 轮转。

```html
<section class="slide" data-layout="feature-grid">
  <p class="kicker">// {{引导语，如 part 02 · design}}</p>
  <h2 class="h2">{{标题：这组特性回答什么，关键词用 <span class="accent">accent</span>}}</h2>
  <div class="grid g2 mt-l">
    <div>
      <div class="feature-row"><span class="num blue">①</span><div><b>{{特性名，≤10 字}}</b><p class="dim">{{说明，18-30 字}}</p></div></div>
      <div class="feature-row"><span class="num green">②</span><div><b>{{特性名，≤10 字}}</b><p class="dim">{{说明，18-30 字}}</p></div></div>
    </div>
    <div>
      <div class="feature-row"><span class="num orange">③</span><div><b>{{特性名，≤10 字}}</b><p class="dim">{{说明，18-30 字}}</p></div></div>
      <div class="feature-row"><span class="num purple">④</span><div><b>{{特性名，≤10 字}}</b><p class="dim">{{说明，18-30 字}}</p></div></div>
    </div>
  </div>
  <div class="notes">{{逐字稿 150-300 字：按重要性逐条讲，别按屏幕顺序念}}</div>
</section>
```

合法类名：slide, kicker, h2, accent, grid, g2, mt-l, feature-row, num, blue, green, orange, purple, dim, notes

---

## rule-list（铁律清单）
指纹：stack

用途：3-4 条方法论/铁律/步骤的大号编号行（rule-row：红号 + 加粗标题 + 展开）。
适用 role：content。
内容约束：3-4 条；每条标题 ≤16 字（关键词用 accent）+ 说明 25-40 字。

```html
<section class="slide" data-layout="rule-list">
  <p class="kicker">// part {{序号}} · {{主题}}</p>
  <h2 class="h2">{{标题：几条什么，如：逐字稿的 3 条铁律}}</h2>
  <div class="stack mt-l">
    <div class="rule-row"><span class="num red">01</span><div><b>{{铁律标题，≤16 字，关键词用 <span class="accent">accent</span>}}</b><p class="dim">{{展开，25-40 字}}</p></div></div>
    <div class="rule-row"><span class="num red">02</span><div><b>{{铁律标题，≤16 字}}</b><p class="dim">{{展开，25-40 字}}</p></div></div>
    <div class="rule-row"><span class="num red">03</span><div><b>{{铁律标题，≤16 字}}</b><p class="dim">{{展开，25-40 字}}</p></div></div>
  </div>
  <div class="notes">{{逐字稿 150-300 字：每条为什么成立、最常见的反例}}</div>
</section>
```

合法类名：slide, kicker, h2, stack, mt-l, rule-row, num, red, accent, dim, notes

---

## demo-close（命令演示 + 收尾）
指纹：code

用途：安装/操作命令演示（code-block：comment/cmd/flag 三色）+ 一句收尾 + 页脚 Q&A。兼任整份 deck 的收尾页。
适用 role：code / cta / thanks。
内容约束：代码 6-14 行、每行 ≤60 字符（`<` `>` `&` 转义）；收尾句 20-40 字；关键动作用 `<strong>`。

```html
<section class="slide" data-layout="demo-close">
  <p class="kicker">// part {{序号}} · demo + close</p>
  <h2 class="h2">{{标题：现在你也能做到，关键词用 <span class="accent">accent</span>}}</h2>
  <div class="code-block mt-m"><span class="comment">{{注释：这一步在做什么}}</span>
<span class="cmd">{{命令}}</span> {{参数}}
<span class="flag">{{按键/旗标}}</span> <span class="comment">{{作用说明}}</span></div>
  <p class="lede mt-m tc">{{收尾一句，20-40 字，关键动作用 <strong>}}</p>
  <div class="deck-footer"><span class="mono">#thanks · Q&amp;A</span><span class="slide-number" data-current="{{当前页}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{逐字稿 150-300 字：演示步骤 + 收尾感谢 + 引导提问}}</div>
</section>
```

合法类名：slide, kicker, h2, accent, code-block, mt-m, comment, cmd, flag, lede, tc, deck-footer, mono, slide-number, notes
