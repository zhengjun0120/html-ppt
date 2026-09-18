# 方向键极简 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-dir-key-nav-minimal` 作用域前缀生效，骨架里已写全，照抄结构即可。
>
> **结构铁律**：每页的底色类（`t-indigo` / `t-cream` / …）已写死在骨架里——一页一整幅色是本模板的身份，
> 照抄，不要换用别的 t-*、不要发明新 t-*、不要给内容再包一层容器（整页天然垂直居中）。
> `dk-snum` / `dk-keyhint` / `dk-page` 是装饰性 chrome 小字（源 demo 键盘导航的残留，纯展示、不承载内容），
> 照抄结构、换成当页文案即可。
> accent 只有一个钩子：`<span class="dk-accent">`，每页点 1 处关键词或主数字。

---

## cover（封面）
指纹：hero

用途：开场页。巨大两行标题 + accent 短线 + 一句话定位。indigo 底。
适用 role：cover。
内容约束：标题两行内、每行 4-10 字；lede 25-45 字；accent 词 ≤1 个；eyebrow 写主题标签。

```html
<section class="slide t-indigo" data-layout="cover">
  <div class="dk-snum">{{页码，如 01 / 08}}</div>
  <p class="dk-eyebrow">{{主题标签，≤20 字符}}</p>
  <h1 class="dk-h0">{{主标题，两行用 <br>；点 1 个关键词加 <span class="dk-accent">accent 色</span>}}</h1>
  <span class="dk-line"></span>
  <p class="dk-lede">{{一句话定位：讲给谁、回答什么，25-45 字}}</p>
  <!-- dk-snum / dk-keyhint / dk-page 是装饰性 chrome 小字（源 demo 键盘导航残留，纯展示、不承载内容），照抄结构、换成当页文案 -->
  <div class="dk-keyhint">{{装饰标注，如 keynote · cover}}</div>
  <div class="dk-page">{{页标，如 cover}}</div>
  <div class="notes">{{讲稿 2-3 句：开场钩子 + 这次演讲的目标}}</div>
</section>
```

合法类名：slide, t-indigo, dk-snum, dk-eyebrow, dk-h0, dk-accent, dk-line, dk-lede, dk-keyhint, dk-page, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。Chapter 号 + 大标题 + 这一章回答的一个问题。cream 底（全册唯一的浅色章节页）。
适用 role：divider。
内容约束：章节号写进 eyebrow（如 Chapter 02）；标题 ≤8 字；lede 25-45 字（一个问题）。

```html
<section class="slide t-cream" data-layout="divider">
  <div class="dk-snum">{{页码，如 02 / 08}}</div>
  <p class="dk-eyebrow">Chapter {{章节号，如 01}} / {{共几章}}</p>
  <h1 class="dk-h0">{{章节标题，≤8 字，可点 <span class="dk-accent">1 个 accent 词</span>}}</h1>
  <span class="dk-line"></span>
  <p class="dk-lede">{{这一章回答的一个问题，25-45 字}}</p>
  <!-- dk-snum / dk-keyhint / dk-page 是装饰性 chrome 小字（源 demo 键盘导航残留，纯展示、不承载内容），照抄结构、换成当页文案 -->
  <div class="dk-keyhint">{{装饰标注，如 chapter · 01 / 04}}</div>
  <div class="dk-page">{{页标，如 section}}</div>
  <div class="notes">{{讲稿：为什么这一章值得听}}</div>
</section>
```

合法类名：slide, t-cream, dk-snum, dk-eyebrow, dk-h0, dk-accent, dk-line, dk-lede, dk-keyhint, dk-page, notes

---

## list（要点清单）
指纹：stack

用途：3-5 条痛点/要点/现象的箭头清单（dk-list，等宽字体 + → 前缀）。crimson 底。
适用 role：content / toc。
内容约束：标题 ≤12 字（结论式，不写栏目名）；3-5 条、每条 12-24 字。

```html
<section class="slide t-crimson" data-layout="list">
  <div class="dk-snum">{{页码，如 03 / 08}}</div>
  <p class="dk-eyebrow">{{这页在讲什么，≤16 字符}}</p>
  <h2 class="dk-h1">{{结论式标题，两行用 <br>；可点 <span class="dk-accent">1 个 accent 词</span>}}</h2>
  <ul class="dk-list">
    <li>{{要点 1，12-24 字}}</li>
    <li>{{要点 2，12-24 字}}</li>
    <li>{{要点 3，12-24 字}}</li>
    <li>{{要点 4，12-24 字}}</li>
  </ul>
  <!-- dk-snum / dk-keyhint / dk-page 是装饰性 chrome 小字（源 demo 键盘导航残留，纯展示、不承载内容），照抄结构、换成当页文案 -->
  <div class="dk-keyhint">{{装饰标注，如 content · list}}</div>
  <div class="dk-page">{{页标，如 03}}</div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, t-crimson, dk-snum, dk-eyebrow, dk-h1, dk-accent, dk-list, dk-keyhint, dk-page, notes

---

## compare（两栏对照）
指纹：cards

用途：× / ✓ 两个做法的正反对照（dk-grid-2 两栏）。emerald 底。
适用 role：content。
内容约束：恰好 2 栏；栏标题 ≤10 字（× 或 ✓ 空格开头）；每栏说明 15-30 字；lede 30-50 字先讲结论。

```html
<section class="slide t-emerald" data-layout="compare">
  <div class="dk-snum">{{页码，如 04 / 08}}</div>
  <p class="dk-eyebrow">{{对比语境，≤16 字符}}</p>
  <h2 class="dk-h1">{{结论式标题，可点 <span class="dk-accent">1 个 accent 词</span>}}</h2>
  <p class="dk-lede" style="margin-top:10px">{{先把结论讲清，30-50 字}}</p>
  <div class="dk-grid-2">
    <div class="dk-col"><h3>× {{被否定的做法，≤10 字}}</h3><p>{{它的代价，15-30 字}}</p></div>
    <div class="dk-col"><h3>✓ {{推荐的做法，≤10 字}}</h3><p>{{它的收益，15-30 字}}</p></div>
  </div>
  <!-- dk-snum / dk-keyhint / dk-page 是装饰性 chrome 小字（源 demo 键盘导航残留，纯展示、不承载内容），照抄结构、换成当页文案 -->
  <div class="dk-keyhint">{{装饰标注，如 content · compare}}</div>
  <div class="dk-page">{{页标，如 04}}</div>
  <div class="notes">{{讲稿：为什么否一个、荐一个}}</div>
</section>
```

合法类名：slide, t-emerald, dk-snum, dk-eyebrow, dk-h1, dk-accent, dk-lede, dk-grid-2, dk-col, dk-keyhint, dk-page, notes

---

## code（代码页）
指纹：code

用途：一段配置/命令/代码（dk-code 块）+ 它带来什么。slate 底。
适用 role：code / content。
内容约束：代码 3-12 行、每行 ≤48 字符（`<` `>` `&` 转义成 &lt; &gt; &amp;）；lede 30-50 字说明它带来什么。

```html
<section class="slide t-slate" data-layout="code">
  <div class="dk-snum">{{页码，如 05 / 08}}</div>
  <p class="dk-eyebrow">{{文件名 / 场景，≤20 字符}}</p>
  <h2 class="dk-h2">{{标题：这一小段在解决什么，可点 <span class="dk-accent">accent 词</span>}}</h2>
  <pre class="dk-code">{{代码/配置 3-12 行（&lt; &gt; &amp; 转义；真实可读，不放伪代码）}}</pre>
  <p class="dk-lede" style="margin-top:16px;font-size:20px">{{这段代码/配置带来什么，30-50 字}}</p>
  <!-- dk-snum / dk-keyhint / dk-page 是装饰性 chrome 小字（源 demo 键盘导航残留，纯展示、不承载内容），照抄结构、换成当页文案 -->
  <div class="dk-keyhint">{{装饰标注，如 content · code}}</div>
  <div class="dk-page">{{页标，如 05}}</div>
  <div class="notes">{{讲稿：逐行/逐项它在做什么}}</div>
</section>
```

合法类名：slide, t-slate, dk-snum, dk-eyebrow, dk-h2, dk-accent, dk-code, dk-lede, dk-keyhint, dk-page, notes

---

## stat（大数字）
指纹：chart

用途：单个关键数字 + SVG 占比条的冲击力呈现。violet 底。
适用 role：data / content。
内容约束：主数字 ≤6 字符；结论 40-60 字（含 1-2 个具体对比数字）；口径行 15-30 字；进度条 rect 宽度按占比给（0-900）。

```html
<section class="slide t-violet" data-layout="stat">
  <div class="dk-snum">{{页码，如 06 / 08}}</div>
  <p class="dk-eyebrow">{{指标语境，如 30-day result，≤16 字符}}</p>
  <div class="dk-big dk-accent">{{主数字，≤6 字符}}</div>
  <p class="dk-lede" style="margin-top:14px;font-size:26px">{{这个数字说明什么，40-60 字，含具体对比}}</p>
  <svg viewBox="0 0 900 80" style="width:100%;max-width:900px;margin-top:30px" aria-hidden="true">
    <rect x="0" y="30" width="900" height="22" rx="11" fill="rgba(127,127,127,.28)"/>
    <rect class="dk-accent" x="0" y="30" width="{{0-900，按占比}}" height="22" rx="11" fill="currentColor"/>
    <text class="dk-accent" x="{{宽度+10，≤870}}" y="47" font-family="JetBrains Mono" font-size="16" font-weight="700" fill="currentColor">{{占比标签}}</text>
  </svg>
  <p class="dk-lede" style="margin-top:18px;font-size:20px">口径：{{样本 / 时间窗 / 出处，15-30 字}}</p>
  <!-- dk-snum / dk-keyhint / dk-page 是装饰性 chrome 小字（源 demo 键盘导航残留，纯展示、不承载内容），照抄结构、换成当页文案 -->
  <div class="dk-keyhint">{{装饰标注，如 chart · big-num}}</div>
  <div class="dk-page">{{页标，如 06}}</div>
  <div class="notes">{{讲稿：数字怎么来的、该怎么读}}</div>
</section>
```

合法类名：slide, t-violet, dk-snum, dk-eyebrow, dk-big, dk-accent, dk-lede, dk-keyhint, dk-page, notes

---

## cta（行动号召）
指纹：code

用途：号召行动 + 1-4 行可复制的起步命令（dk-code 承载命令）。white 底（全册唯一的纯白页）。
适用 role：cta。
内容约束：标题两行内 ≤8 字/行；lede 25-40 字（为什么现在就做）；命令 1-4 行、可直接复制执行。

```html
<section class="slide t-white" data-layout="cta">
  <div class="dk-snum">{{页码，如 07 / 08}}</div>
  <p class="dk-eyebrow">{{行动语境，如 Start tonight}}</p>
  <h2 class="dk-h1">{{号召，两行用 <br>；可点 <span class="dk-accent">1 个 accent 词</span>}}</h2>
  <span class="dk-line"></span>
  <p class="dk-lede">{{为什么现在就做，25-40 字}}</p>
  <pre class="dk-code" style="font-size:18px">{{1-4 行可复制的命令（&lt; &gt; &amp; 转义）}}</pre>
  <!-- dk-snum / dk-keyhint / dk-page 是装饰性 chrome 小字（源 demo 键盘导航残留，纯展示、不承载内容），照抄结构、换成当页文案 -->
  <div class="dk-keyhint">{{装饰标注，如 cta · three-commands}}</div>
  <div class="dk-page">{{页标，如 07}}</div>
  <div class="notes">{{讲稿：第一步做什么、做完 expect 什么}}</div>
</section>
```

合法类名：slide, t-white, dk-snum, dk-eyebrow, dk-h1, dk-accent, dk-line, dk-lede, dk-code, dk-keyhint, dk-page, notes

---

## thanks（收尾页）
指纹：hero

用途：致谢 + 出处/链接。charcoal 底。结构回到与封面同款的大标题居中。
适用 role：thanks。
内容约束：标题 ≤6 字；lede 25-45 字，含真实的链接或出处（没有真实链接就写一句收束，不编造）。

```html
<section class="slide t-charcoal" data-layout="thanks">
  <div class="dk-snum">{{页码，如 08 / 08}}</div>
  <p class="dk-eyebrow">{{收尾语境，如 End · thanks for staying}}</p>
  <h1 class="dk-h0">{{收尾词，≤6 字，可点 <span class="dk-accent">1 个 accent 词</span>}}</h1>
  <span class="dk-line"></span>
  <p class="dk-lede">{{致谢 + 去哪里找到你（真实链接/出处），25-45 字}}</p>
  <!-- dk-snum / dk-keyhint / dk-page 是装饰性 chrome 小字（源 demo 键盘导航残留，纯展示、不承载内容），照抄结构、换成当页文案 -->
  <div class="dk-keyhint">{{装饰标注，如 thanks · fin}}</div>
  <div class="dk-page">{{页标，如 fin}}</div>
  <div class="notes">{{讲稿：收束一句 + 引导提问}}</div>
</section>
```

合法类名：slide, t-charcoal, dk-snum, dk-eyebrow, dk-h0, dk-accent, dk-line, dk-lede, dk-keyhint, dk-page, notes
