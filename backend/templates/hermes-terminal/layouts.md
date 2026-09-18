# 终端实测 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，
> 只准使用骨架里出现过的类名（base 原语 + 本模板 `hc-*` 专属类）。
> 未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-hermes-cyber-terminal` 作用域前缀生效，骨架里已写全，照抄结构即可。
>
> **结构铁律**：`hc-grid` / `hc-scanlines` 两个空气层 + 顶部 `hc-chrome`（traffic-light 圆点 + 终端标题栏）
> 是本模板的身份，每页都必须保留，只准改标题栏文字；
> 底部 `hc-footer`（左眉标 + 右页码）每页保留。{{占位符}} 只准换文本；
> 骨架的层级与包裹关系不准增删。全局等宽字体，中文也照排，不要另引字体。
>
> 每页讲稿写在页面末尾：`<div class="notes">…</div>`（观众不可见）。

---

## cover（终端封面）
指纹：hero

用途：开场页。命令行 prompt + glow 大标题 + 评测立场一句话 + 分级标签。
适用 role：cover。
内容约束：prompt 6-20 字符；主标题两行、每行 ≤14 字符；lede 30-55 字（评的是什么 + 站位）；3-4 个标签、每个 ≤16 字符。

```html
<section class="slide" data-layout="cover">
  <div class="hc-grid"></div>
  <div class="hc-vignette"></div>
  <div class="hc-scanlines"></div>
  <div class="hc-chrome"><div class="dots"><span></span><span></span><span></span></div><div>{{终端标题栏，如 ~/review · zsh · 03:14}}</div></div>
  <div style="margin:auto 0">
    <p class="hc-prompt">{{命令式开场，如 whoami --honest}}</p>
    <h1 class="hc-h1">{{主标题第一行}}<br>{{主标题第二行，可带版本号}}<span class="hc-cursor"></span></h1>
    <p class="hc-lede">{{评的是什么 + 一句站位，30-55 字}}</p>
    <div style="margin-top:26px">
      <span class="hc-tag">{{标签 1，小写英文}}</span>
      <span class="hc-tag">{{标签 2}}</span>
      <span class="hc-tag amber">{{琥珀标签：卖点点}}</span>
      <span class="hc-tag red">{{红色标签：立场/警示}}</span>
    </div>
  </div>
  <div class="hc-footer"><span>{{deck 名 · 作者 · 年份}}</span><span>{{页码，如 01 / 08}}</span></div>
  <div class="notes">{{讲稿 2-3 句：评什么 + 亮出立场}}</div>
</section>
```

合法类名：slide, hc-grid, hc-vignette, hc-scanlines, hc-chrome, dots, hc-prompt, hc-h1, hc-cursor, hc-lede, hc-tag, amber, red, hc-footer, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。prompt 做读档动作，巨大章节标题 + 一句过门。
适用 role：divider。
内容约束：prompt 6-20 字符；标题 ≤12 字符（`//` 前缀是惯用味）；lede 20-40 字（这节讲什么/多久）。

```html
<section class="slide" data-layout="divider">
  <div class="hc-grid"></div>
  <div class="hc-scanlines"></div>
  <div class="hc-chrome"><div class="dots"><span></span><span></span><span></span></div><div>{{章节进度，如 section · 01/04}}</div></div>
  <div style="margin:auto 0">
    <p class="hc-prompt">{{读档动作，如 cat chapter_01.md}}</p>
    <h1 class="hc-h1" style="font-size:110px">{{// 章节标题，≤12 字符}}</h1>
    <p class="hc-lede">{{一句过门：这一节讲什么、花多久，20-40 字}}</p>
  </div>
  <div class="hc-footer"><span>{{section · 小节名}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, hc-grid, hc-scanlines, hc-chrome, dots, hc-prompt, hc-h1, hc-lede, hc-footer, notes

---

## spec-cards（实测数据卡）
指纹：cards

用途：3 张指标卡（数值 + 口径）+ 2 张正反结论卡。跑分/冷启动/安装数据的主力版式。
适用 role：data / content。
内容约束：3 张指标卡：数值 ≤6 字符 + 英文小标签 + 口径说明 12-24 字；2 张结论卡：结论 8-16 字 + 展开 20-40 字（和谁比、好在哪/坑在哪）；每个数字必须有口径（几次平均、什么条件）。

```html
<section class="slide" data-layout="spec-cards">
  <div class="hc-grid"></div>
  <div class="hc-scanlines"></div>
  <div class="hc-chrome"><div class="dots"><span></span><span></span><span></span></div><div>{{眉标，如 benchmark · cold-start}}</div></div>
  <h2 class="hc-h2">{{数据标题，≤10 字}}</h2>
  <p class="hc-lede">{{口径一句话：测什么、怎么平均，20-40 字}}</p>
  <div class="hc-grid-3">
    <div class="hc-card"><div class="lbl">{{指标名，小写英文}}</div><div class="val">{{数值 ≤6 字符}}</div><div class="desc">{{口径说明，12-24 字}}</div></div>
    <div class="hc-card"><div class="lbl">{{指标名}}</div><div class="val">{{数值}}</div><div class="desc">{{口径说明，12-24 字}}</div></div>
    <div class="hc-card"><div class="lbl">{{指标名}}</div><div class="val">{{数值}}</div><div class="desc">{{口径说明，12-24 字}}</div></div>
  </div>
  <div class="hc-grid-2">
    <div class="hc-card"><div class="lbl">// verdict +</div><div class="val" style="color:var(--hc-green);font-size:18px">{{加分结论，8-16 字}}</div><div class="desc">{{和谁比、好在哪，20-40 字}}</div></div>
    <div class="hc-card"><div class="lbl">// verdict -</div><div class="val" style="color:var(--hc-red);font-size:18px">{{减分结论，8-16 字}}</div><div class="desc">{{坑在哪、多难受，20-40 字}}</div></div>
  </div>
  <div class="hc-footer"><span>{{data · 口径缩写}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, hc-grid, hc-scanlines, hc-chrome, dots, hc-h2, hc-lede, hc-grid-3, hc-card, lbl, val, desc, hc-grid-2, hc-footer, notes

---

## trace（执行 trace）
指纹：code

用途：一段真实执行日志/代码，逐行展示工具或 agent 干了什么。trace 必须真实可读。
适用 role：code。
内容约束：prompt 是触发命令 ≤60 字符；代码 8-16 行、每行 ≤70 字符；`[plan/read/edit/test/commit/push]` 这类阶段词用 `fn`，引号内内容用 `st`，注释用 `cm`，diff 数字用 `hl`。

```html
<section class="slide" data-layout="trace">
  <div class="hc-grid"></div>
  <div class="hc-scanlines"></div>
  <div class="hc-chrome"><div class="dots"><span></span><span></span><span></span></div><div>{{眉标，如 trace · 工具名 run}}</div></div>
  <p class="hc-prompt">{{触发命令，≤60 字符}}</p>
  <h3 class="hc-h3" style="margin-top:12px">{{↓ 一句说明，如 ↓ 真实 trace（节选）}}</h3>
  <pre class="hc-codebox" style="margin-top:10px"><span class="cm">{{# 版本/会话注释}}</span>
[<span class="fn">plan</span>]    <span class="st">{{"计划一句话"}}</span>
[<span class="fn">read</span>]    {{文件路径}}             <span class="cm">// {{行数/备注}}</span>
[<span class="kw">think</span>]   <span class="st">{{"关键决策一句话"}}</span>
[<span class="fn">edit</span>]    {{文件路径}}             <span class="hl">+{{N}} -{{N}}</span>
[<span class="fn">test</span>]    {{测试命令}}        <span class="st">{{结果，如 PASS 18/18}}</span>
[<span class="fn">commit</span>]  <span class="st">{{"commit message"}}</span>

<span class="cm"># {{总耗时 · tokens · 花费}}</span></pre>
  <div class="hc-footer"><span>{{trace · live}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿：这一段里最关键的一步}}</div>
</section>
```

合法类名：slide, hc-grid, hc-scanlines, hc-chrome, dots, hc-prompt, hc-h3, hc-codebox, cm, fn, st, kw, hl, hc-footer, notes

---

## compare-chart（对比柱状图）
指纹：chart

用途：两个方案/版本在同一组指标上的柱状对比。有对比数据就画图，别写成句子。
适用 role：data。
内容约束：标题 ≤12 字；口径 lede 20-40 字；A/B 两组各 2-4 根柱、总数 ≤7；柱顶标数值、柱底标指标名 ≤10 字符；必须带图例，眉标注明样本量（如 `benchmark · n=48`）。柱高 = 百分比 × 2.8，y = 320 − 柱高。

```html
<section class="slide" data-layout="compare-chart">
  <div class="hc-grid"></div>
  <div class="hc-scanlines"></div>
  <div class="hc-chrome"><div class="dots"><span></span><span></span><span></span></div><div>{{眉标，含样本量，如 benchmark · n=48}}</div></div>
  <h2 class="hc-h2">{{对比标题，≤12 字}}</h2>
  <p class="hc-lede">{{口径：同一批任务、各自跑的条件，20-40 字}}</p>
  <svg viewBox="0 0 1000 380" style="width:100%;max-width:1040px;margin-top:24px" xmlns="http://www.w3.org/2000/svg">
    <g font-family="JetBrains Mono, monospace" font-size="16" fill="#8a8892">
      <line x1="80" y1="40" x2="80" y2="320" stroke="rgba(126,211,164,.2)"/>
      <line x1="80" y1="320" x2="960" y2="320" stroke="rgba(126,211,164,.2)"/>
      <!-- A 组柱：柱宽 80、x 间距 110；复制整组加柱 -->
      <g>
        <rect x="130" y="80" width="80" height="240" style="fill:var(--hc-s1-soft);stroke:var(--hc-s1);stroke-width:1.5"/>
        <text x="170" y="76" text-anchor="middle" style="fill:var(--hc-s1)" font-weight="700">{{数值%}}</text>
        <text x="170" y="345" text-anchor="middle">{{指标 ≤10 字符}}</text>
      </g>
      <g>
        <rect x="240" y="150" width="80" height="170" style="fill:var(--hc-s1-soft);stroke:var(--hc-s1);stroke-width:1.5"/>
        <text x="280" y="146" text-anchor="middle" style="fill:var(--hc-s1)" font-weight="700">{{数值%}}</text>
        <text x="280" y="345" text-anchor="middle">{{指标}}</text>
      </g>
      <!-- B 组柱：同样式换 --hc-s2，x 从 570 起、间距 110 -->
      <g>
        <rect x="570" y="150" width="80" height="170" style="fill:var(--hc-s2-soft);stroke:var(--hc-s2);stroke-width:1.5"/>
        <text x="610" y="146" text-anchor="middle" style="fill:var(--hc-s2)" font-weight="700">{{数值%}}</text>
        <text x="610" y="345" text-anchor="middle">{{指标}}</text>
      </g>
      <g transform="translate(820,50)">
        <rect x="0" y="0" width="14" height="14" style="fill:var(--hc-s1-soft);stroke:var(--hc-s1)"/>
        <text x="22" y="12" style="fill:var(--hc-s1)">{{系列 A 名}}</text>
        <rect x="0" y="22" width="14" height="14" style="fill:var(--hc-s2-soft);stroke:var(--hc-s2)"/>
        <text x="22" y="34" style="fill:var(--hc-s2)">{{系列 B 名}}</text>
      </g>
    </g>
  </svg>
  <div class="hc-footer"><span>{{benchmark · n=样本量}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿：差距最大的一项怎么解读}}</div>
</section>
```

合法类名：slide, hc-grid, hc-scanlines, hc-chrome, dots, hc-h2, hc-lede, hc-footer, notes

---

## verdict（评分判决）
指纹：chart

用途：单个大分数 + 一句判决 + 优缺点两卡，收束整个评测。
适用 role：data / content。
内容约束：分数 ≤4 字符（可带分母 span）；判决 lede 12-24 字；两卡各 3 条要点、每条 8-20 字（bullet 用 `•` + `<br>` 分行）。分数必须有依据，没有就用量级表述。

```html
<section class="slide" data-layout="verdict">
  <div class="hc-grid"></div>
  <div class="hc-scanlines"></div>
  <div class="hc-chrome"><div class="dots"><span></span><span></span><span></span></div><div>{{眉标，如 tldr}}</div></div>
  <p class="hc-prompt">{{取值动作，如 echo $VERDICT}}</p>
  <div class="hc-big">{{分数 ≤4 字符}}<span style="font-size:60px;color:var(--hc-ink2)">{{分母，如 / 10}}</span></div>
  <p class="hc-lede" style="margin-top:14px">{{一句判决，12-24 字}}</p>
  <div class="hc-grid-2" style="margin-top:24px">
    <div class="hc-card"><div class="lbl">+ strong points</div><div class="desc">• {{优点 1，8-20 字}}<br>• {{优点 2，8-20 字}}<br>• {{优点 3，8-20 字}}</div></div>
    <div class="hc-card"><div class="lbl">- weak points</div><div class="desc">• {{缺点 1，8-20 字}}<br>• {{缺点 2，8-20 字}}<br>• {{缺点 3，8-20 字}}</div></div>
  </div>
  <div class="hc-footer"><span>{{verdict · honest}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿：这个分怎么来的}}</div>
</section>
```

合法类名：slide, hc-grid, hc-scanlines, hc-chrome, dots, hc-prompt, hc-big, hc-lede, hc-grid-2, hc-card, lbl, desc, hc-footer, notes

---

## install（上手命令）
指纹：code

用途：CTA。让观众自己跑一遍的命令清单，codebox 承载，分步注释。
适用 role：cta / code。
内容约束：标题 ≤12 字；lede 16-32 字（照做能干什么 + 多久）；命令 4-10 行、每行 ≤70 字符、步骤用 `# 注释` 分开；2-3 个环境标签（前置条件用 amber）。

```html
<section class="slide" data-layout="install">
  <div class="hc-grid"></div>
  <div class="hc-scanlines"></div>
  <div class="hc-chrome"><div class="dots"><span></span><span></span><span></span></div><div>{{眉标，如 install}}</div></div>
  <h2 class="hc-h2">{{行动标题，≤12 字}}</h2>
  <p class="hc-lede">{{照做能干什么 + 要多久，16-32 字}}</p>
  <pre class="hc-codebox" style="margin-top:22px"><span class="cm"># 1. {{步骤说明}}</span>
<span class="kw">$</span> {{命令 1}}

<span class="cm"># 2. {{步骤说明}}</span>
<span class="kw">$</span> {{命令 2}}

<span class="cm"># 3. {{步骤说明}}</span>
<span class="kw">$</span> {{命令 3}}</pre>
  <div style="margin-top:26px">
    <span class="hc-tag">{{环境标签 1}}</span>
    <span class="hc-tag">{{环境标签 2}}</span>
    <span class="hc-tag amber">{{前置条件标签}}</span>
  </div>
  <div class="hc-footer"><span>{{try-it-now}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿：最容易卡住的一步}}</div>
</section>
```

合法类名：slide, hc-grid, hc-scanlines, hc-chrome, dots, hc-h2, hc-lede, hc-codebox, cm, kw, hc-tag, amber, hc-footer, notes

---

## thanks（会话收尾）
指纹：hero

用途：收尾页。exit 提示 + 感谢 + 资料仓库链接。
适用 role：thanks / cta。
内容约束：prompt ≤12 字符；大标题 ≤12 字符；lede 20-45 字（资料在哪；链接要真实，没有就写后续计划）。

```html
<section class="slide" data-layout="thanks">
  <div class="hc-grid"></div>
  <div class="hc-scanlines"></div>
  <div class="hc-chrome"><div class="dots"><span></span><span></span><span></span></div><div>{{眉标，如 EOF}}</div></div>
  <div style="margin:auto 0">
    <p class="hc-prompt">{{收尾命令，如 exit 0}}</p>
    <h1 class="hc-h1" style="font-size:120px">{{// 收尾标题，≤12 字符}}<span class="hc-cursor"></span></h1>
    <p class="hc-lede">{{完整资料在哪：仓库/脚本/文章，20-45 字}}</p>
  </div>
  <div class="hc-footer"><span>{{session closed}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿：一句话收尾}}</div>
</section>
```

合法类名：slide, hc-grid, hc-scanlines, hc-chrome, dots, hc-prompt, hc-h1, hc-cursor, hc-lede, hc-footer, notes
