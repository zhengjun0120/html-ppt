# 小红书柔色 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-xhs-pastel-card` 作用域前缀生效，骨架里已写全，照抄结构即可。
>
> **结构铁律**：每页必须保持「`<div class="xp-topbar">`（chip + 页码）开头、
> `<div class="xp-footer">` 收尾」的包裹结构，并至少保留一颗 `xp-blob` 柔光——
> chip 顶条 + blob 是本模板的身份，丢了页面就散架。chip 用小写英文短语，
> 页码格式照抄骨架（如 `01 · 08`）。{{占位符}} 只准换文本；
> 骨架的层级与包裹关系不准增删。

---

## cover（柔色封面）
指纹：hero

用途：开场页。chip 顶条 + 衬线大标题（斜体强调词）+ 渐变小条 + 一句话定位。
适用 role：cover。
内容约束：主标题 2-3 行（`<br>` 分行）、每行 ≤7 字，强调词包 `<em>` 或 `<span class="rose">`（各 ≤4 字，共 ≤2 处）；xp-sub 30-55 字（这不是 X，是把 Y 还给你）；chip 为小写英文短语 ≤3 词。

```html
<section class="slide" data-layout="cover">
  <div class="xp-blob b1"></div>
  <div class="xp-blob b2"></div>
  <div class="xp-blob b3"></div>
  <div class="xp-topbar"><div class="xp-chip">{{英文栏目名，≤3 词}}</div><div class="xp-page">{{页码，如 01 · 08}}</div></div>
  <div class="xp-kicker">{{主题行，如 Living With AI · 2026}}</div>
  <h1 class="xp-h1">{{主标题，2-3 行 <br> 分行，强调词包 <em> 或 <span class="rose">}}</h1>
  <div class="xp-divider"></div>
  <p class="xp-sub">{{一句话定位：这不是 X，是把 Y 还给你，30-55 字}}</p>
  <div class="xp-footer"><span>{{作者 · 栏目，如 by lewis · pastel edition}}</span><span>cover</span></div>
  <div class="notes">{{发布文案或讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, xp-blob, b1, b2, b3, xp-topbar, xp-chip, xp-page, xp-kicker, xp-h1, rose, mint, xp-divider, xp-sub, xp-footer, notes

---

## chapter（提问式章节页）
指纹：hero

用途：章节过渡。一页只问一个问题：大号衬线问句 + 一句戳心的展开，不带任何卡片。
适用 role：divider。
内容约束：问句两行（`<br>` 分行）、每行 ≤8 字，强调词包 `<span class="mint">`（≤6 字）；xp-sub 18-35 字；xp-kicker 引导 4-8 字。

```html
<section class="slide" data-layout="chapter">
  <div class="xp-blob b2"></div>
  <div class="xp-blob b3"></div>
  <div class="xp-topbar"><div class="xp-chip mint">{{英文栏目名，如 Chapter one}}</div><div class="xp-page">{{页码，如 02 · 08}}</div></div>
  <div style="margin:auto 0">
    <div class="xp-kicker">{{小引导，4-8 字}}</div>
    <h1 class="xp-h1" style="font-size:120px">{{一个问句，两行 <br> 分行，强调词包 <span class="mint">}}</h1>
    <p class="xp-sub">{{一句戳心的展开，18-35 字}}</p>
  </div>
  <div class="xp-footer"><span>section · chapter {{章节号}}</span><span>{{页码，如 02 · 08}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, xp-blob, b1, b2, b3, xp-topbar, xp-chip, mint, sky, lilac, rose, xp-page, xp-kicker, xp-h1, xp-sub, xp-footer, notes

---

## grid-2x2（四事网格）
指纹：cards

用途：四件并列的小事/小物/小方法：2×2 马卡龙整色卡阵列（斜体序号 01-04）。这是本模板的内容主力版式。
适用 role：content。
内容约束：恰好 4 卡；卡标题 ≤14 字；卡内说明 22-45 字（一句为什么 + 一句省了多少，句号结尾）；序号固定 01-04；四卡颜色从 peach/mint/sky/lilac 里取（可换 lemon/rose）。

```html
<section class="slide" data-layout="grid-2x2">
  <div class="xp-blob b1"></div>
  <div class="xp-topbar"><div class="xp-chip rose">{{英文栏目名，如 Four little escapes}}</div><div class="xp-page">{{页码，如 03 · 08}}</div></div>
  <h2 class="xp-h2">{{标题：N 件可以交出去的小事，两行 <br>、每行 ≤8 字，强调词包 <em>}}</h2>
  <div class="xp-grid-2">
    <div class="xp-card peach"><div class="xp-num">01</div><h4>{{小事 1，≤14 字}}</h4><p>{{为什么交给它 + 省了多少，22-45 字}}</p></div>
    <div class="xp-card mint"><div class="xp-num">02</div><h4>{{小事 2，≤14 字}}</h4><p>{{说明，22-45 字}}</p></div>
    <div class="xp-card sky"><div class="xp-num">03</div><h4>{{小事 3，≤14 字}}</h4><p>{{说明，22-45 字}}</p></div>
    <div class="xp-card lilac"><div class="xp-num">04</div><h4>{{小事 4，≤14 字}}</h4><p>{{说明，22-45 字}}</p></div>
  </div>
  <div class="xp-footer"><span>content · 2x2</span><span>{{页码，如 03 · 08}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, xp-blob, b1, b2, b3, xp-topbar, xp-chip, mint, sky, lilac, rose, xp-page, xp-h2, xp-grid-2, xp-card, peach, mint, sky, lilac, lemon, rose, xp-num, xp-footer, notes

---

## checklist（单列清单）
指纹：stack

用途：4-6 条要点清单/打卡列表：一行一条的小卡纵列，色点呼应卡片色。条目多、每条短时用它，别硬塞进 2×2。
适用 role：content / toc。
内容约束：4-6 条；每条 12-24 字（动词或名词开头的一整句）；页尾可加一句收束 16-32 字。

```html
<section class="slide" data-layout="checklist">
  <div class="xp-blob b3"></div>
  <div class="xp-topbar"><div class="xp-chip sky">{{英文栏目名，如 Tiny checklist}}</div><div class="xp-page">{{页码，如 04 · 08}}</div></div>
  <h2 class="xp-h2">{{清单标题，≤14 字}}</h2>
  <div class="xp-legend">
    <div class="xp-card peach" style="padding:18px 26px;display:flex;align-items:center;gap:16px"><div class="xp-dot" style="background:var(--xp-peach-d)"></div><h4 style="margin:0;font-size:21px">{{条目 1，12-24 字}}</h4></div>
    <div class="xp-card mint" style="padding:18px 26px;display:flex;align-items:center;gap:16px"><div class="xp-dot" style="background:var(--xp-mint-d)"></div><h4 style="margin:0;font-size:21px">{{条目 2，12-24 字}}</h4></div>
    <div class="xp-card sky" style="padding:18px 26px;display:flex;align-items:center;gap:16px"><div class="xp-dot" style="background:var(--xp-sky-d)"></div><h4 style="margin:0;font-size:21px">{{条目 3，12-24 字}}</h4></div>
    <div class="xp-card lilac" style="padding:18px 26px;display:flex;align-items:center;gap:16px"><div class="xp-dot" style="background:var(--xp-lilac-d)"></div><h4 style="margin:0;font-size:21px">{{条目 4，12-24 字}}</h4></div>
  </div>
  <p class="xp-sub" style="margin-top:22px">{{一句收束或提醒，16-32 字}}</p>
  <div class="xp-footer"><span>checklist</span><span>{{页码，如 04 · 08}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, xp-blob, b1, b2, b3, xp-topbar, xp-chip, mint, sky, lilac, rose, xp-page, xp-h2, xp-legend, xp-card, peach, mint, sky, lilac, lemon, rose, xp-dot, xp-sub, xp-footer, notes

---

## quote-card（金句卡）
指纹：quote

用途：金句/态度宣言。白色大圆角卡承载衬线斜体金句，下面一段落地展开。整页只说这一句。
适用 role：quote / content。
内容约束：金句 15-35 字（强调词包 `<em>`，分行用 `<br>`）；展开 25-50 字；不用编造的引用，没有真金句就换版式。

```html
<section class="slide" data-layout="quote-card">
  <div class="xp-blob b3"></div>
  <div class="xp-blob b2"></div>
  <div class="xp-topbar"><div class="xp-chip lilac">{{英文栏目名，如 A small pause}}</div><div class="xp-page">{{页码，如 05 · 08}}</div></div>
  <div class="xp-hero-card">
    <p class="xp-quote">{{金句 15-35 字，分行用 <br>，强调词包 <em>}}</p>
    <div class="xp-divider"></div>
    <p class="xp-sub">{{把金句落到生活的展开，25-50 字}}</p>
  </div>
  <div class="xp-footer"><span>quote</span><span>{{页码，如 05 · 08}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, xp-blob, b1, b2, b3, xp-topbar, xp-chip, mint, sky, lilac, rose, xp-page, xp-hero-card, xp-quote, xp-divider, xp-sub, xp-footer, notes

---

## prompt（Prompt 卡）
指纹：code

用途：一段可以直接抄走用的 prompt / 配置 / 短代码。深色圆角码箱是页面的绝对主角。
适用 role：code / content。
内容约束：代码 8-16 行、每行 ≤60 字符；着色只用 kw（指令词）/ st（字符串值）/ hl（关键数字）/ cm（注释）四类；首行与末行各一句 cm 注释（做什么 / 省了多少）；`<` `>` `&` 要转义。

```html
<section class="slide" data-layout="prompt">
  <div class="xp-blob b1"></div>
  <div class="xp-topbar"><div class="xp-chip">{{英文栏目名，如 My auto-reply prompt}}</div><div class="xp-page">{{页码，如 06 · 08}}</div></div>
  <h2 class="xp-h2">{{标题：把 X 交给 AI 的一段 prompt，两行 <br>，强调词包 <em> 或 <span class="rose">}}</h2>
  <pre class="xp-codebox"><span class="cm"># {{这段配置做什么，一句}}</span>
<span class="kw">when</span> {{触发条件，如 email matches}}:
  reply:
    tone: <span class="st">"{{语气要求}}"</span>
    max_lines: <span class="hl">{{关键数字}}</span>

<span class="kw">always_skip</span>:
  - from: [<span class="st">"{{例外 1}}"</span>, <span class="st">"{{例外 2}}"</span>]
  - contains: [<span class="st">"{{排除关键词}}"</span>]

<span class="cm"># {{效果：一周省 N 分钟，测过}}</span></pre>
  <div class="xp-footer"><span>content · prompt</span><span>{{页码，如 06 · 08}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, xp-blob, b1, b2, b3, xp-topbar, xp-chip, mint, sky, lilac, rose, xp-page, xp-h2, xp-codebox, kw, st, hl, cm, rose, xp-footer, notes

---

## donut-stat（环形统计）
指纹：chart

用途：一组「省出来 / 花出去」的时间或占比：SVG 环形图（中心写总数）+ 右侧色点图例列。有占比就该用它，别把数字埋进句子。
适用 role：data / content。
内容约束：3-5 段；每段 stroke-dasharray 第一值 = 占比% × 6.28（整圈周长 628），dashoffset = 前面各段长度的负数累加；中心总数 ≤4 字符；每条图例 ≤14 字（数值 · 名称）；数字要有出处或用约数。

```html
<section class="slide" data-layout="donut-stat">
  <div class="xp-blob b2"></div>
  <div class="xp-topbar"><div class="xp-chip mint">{{英文栏目名，如 Your week rebuilt}}</div><div class="xp-page">{{页码，如 07 · 08}}</div></div>
  <h2 class="xp-h2">{{标题：这些数字还给你什么，≤14 字，强调词包 <span class="mint">}}</h2>
  <div class="xp-chart-row">
    <svg viewBox="0 0 260 260" style="width:340px;flex-shrink:0">
      <circle cx="130" cy="130" r="100" fill="none" stroke="#fef0e4" stroke-width="40"/>
      <circle cx="130" cy="130" r="100" fill="none" stroke="#f48b5c" stroke-width="40" stroke-dasharray="{{占比1×6.28}} 628" stroke-dashoffset="0" transform="rotate(-90 130 130)"/>
      <circle cx="130" cy="130" r="100" fill="none" stroke="#2e9d70" stroke-width="40" stroke-dasharray="{{占比2×6.28}} 628" stroke-dashoffset="-{{占比1×6.28}}" transform="rotate(-90 130 130)"/>
      <circle cx="130" cy="130" r="100" fill="none" stroke="#4e7ed6" stroke-width="40" stroke-dasharray="{{占比3×6.28}} 628" stroke-dashoffset="-{{(占比1+占比2)×6.28}}" transform="rotate(-90 130 130)"/>
      <circle cx="130" cy="130" r="100" fill="none" stroke="#7b5dc4" stroke-width="40" stroke-dasharray="{{占比4×6.28}} 628" stroke-dashoffset="-{{(占比1+占比2+占比3)×6.28}}" transform="rotate(-90 130 130)"/>
      <text x="130" y="132" text-anchor="middle" font-family="Playfair Display" font-size="52" font-weight="900" fill="#2a2340">{{总数，≤4 字符}}</text>
      <text x="130" y="160" text-anchor="middle" font-family="Inter" font-size="15" fill="#9089a8">{{单位说明，≤12 字符}}</text>
    </svg>
    <div class="xp-legend">
      <div class="xp-card peach" style="padding:16px 24px;display:flex;align-items:center;gap:16px"><div class="xp-dot" style="background:var(--xp-peach-d)"></div><h4 style="margin:0;font-size:20px">{{数值 · 名称，≤14 字}}</h4></div>
      <div class="xp-card mint" style="padding:16px 24px;display:flex;align-items:center;gap:16px"><div class="xp-dot" style="background:var(--xp-mint-d)"></div><h4 style="margin:0;font-size:20px">{{数值 · 名称，≤14 字}}</h4></div>
      <div class="xp-card sky" style="padding:16px 24px;display:flex;align-items:center;gap:16px"><div class="xp-dot" style="background:var(--xp-sky-d)"></div><h4 style="margin:0;font-size:20px">{{数值 · 名称，≤14 字}}</h4></div>
      <div class="xp-card lilac" style="padding:16px 24px;display:flex;align-items:center;gap:16px"><div class="xp-dot" style="background:var(--xp-lilac-d)"></div><h4 style="margin:0;font-size:20px">{{数值 · 名称，≤14 字}}</h4></div>
    </div>
  </div>
  <div class="xp-footer"><span>chart · donut</span><span>{{页码，如 07 · 08}}</span></div>
  <div class="notes">{{讲稿：数字口径要真实，没有来源就用约数}}</div>
</section>
```

合法类名：slide, xp-blob, b1, b2, b3, xp-topbar, xp-chip, mint, sky, lilac, rose, xp-page, xp-h2, xp-chart-row, xp-legend, xp-card, peach, mint, sky, lilac, lemon, rose, xp-dot, xp-footer, notes

---

## steps-3（周末三步）
指纹：cards

用途：行动号召页：三个时段/三步安排的横排马卡龙卡，emoji 作装饰序号。放在 decks 的倒数第二页给读者一个「这周末就做」的抓手。
适用 role：content / cta。
内容约束：恰好 3 卡；卡标题（时段）≤12 字；卡内说明 22-40 字（做什么 + 怎么算赢）；emoji 只出现在 xp-num 位（每卡一个）。

```html
<section class="slide" data-layout="steps-3">
  <div class="xp-blob b1"></div>
  <div class="xp-blob b3"></div>
  <div class="xp-topbar"><div class="xp-chip rose">{{英文栏目名，如 This weekend}}</div><div class="xp-page">{{页码，如 08 · 08}}</div></div>
  <h2 class="xp-h2">{{行动号召标题，两行 <br>、每行 ≤8 字，强调词包 <em>}}</h2>
  <div class="xp-grid-3">
    <div class="xp-card lemon"><div class="xp-num">{{一个 emoji}}</div><h4>{{时段 1，≤12 字}}</h4><p>{{做什么，22-40 字}}</p></div>
    <div class="xp-card peach"><div class="xp-num">{{一个 emoji}}</div><h4>{{时段 2，≤12 字}}</h4><p>{{做什么，22-40 字}}</p></div>
    <div class="xp-card sky"><div class="xp-num">{{一个 emoji}}</div><h4>{{时段 3，≤12 字}}</h4><p>{{做什么 + 怎么算赢，22-40 字}}</p></div>
  </div>
  <div class="xp-footer"><span>cta</span><span>{{页码，如 08 · 08}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, xp-blob, b1, b2, b3, xp-topbar, xp-chip, mint, sky, lilac, rose, xp-page, xp-h2, xp-grid-3, xp-card, peach, mint, sky, lilac, lemon, rose, xp-num, xp-footer, notes

---

## thanks（居中收尾）
指纹：hero

用途：谢谢观看/互动引导。居中大字 + 渐变小条 + 一句评论区互动。
适用 role：thanks。
内容约束：主文案 ≤8 字（可夹一个 `<em>`）；互动一句 20-40 字（问读者一件具体的事）；至多一个 ♡。

```html
<section class="slide" data-layout="thanks">
  <div class="xp-blob b2"></div>
  <div style="margin:auto 0;text-align:center">
    <div class="xp-kicker" style="text-align:center">thanks for reading</div>
    <h1 class="xp-h1" style="font-size:150px;text-align:center">{{收尾主文案，≤8 字，可夹 <em>}}</h1>
    <div class="xp-divider" style="margin:24px auto"></div>
    <p class="xp-sub" style="margin:0 auto">{{一句互动引导：问读者一件具体的事，20-40 字}}</p>
  </div>
  <div class="xp-footer"><span>end</span><span>{{页码，如 09 · 09}}</span></div>
  <div class="notes">{{发布文案}}</div>
</section>
```

合法类名：slide, xp-blob, b1, b2, b3, xp-kicker, xp-h1, xp-divider, xp-sub, xp-footer, notes
