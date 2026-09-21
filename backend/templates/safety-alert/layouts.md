# 安全警示 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，
> 只准使用骨架里出现过的类名（base 原语 + 本模板 `ts-*` 专属类）。
> 未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-testing-safety-alert` 作用域前缀生效，骨架里已写全，照抄结构即可。
>
> **结构铁律**：顶部 `ts-stripe`（红黑斜纹警示带）+ 底部 `ts-stripe-b` + `ts-footer` 每页保留；
> `ts-chrome`（警示标签 + 页码）每页保留，警示标签三色用 `amber` / `green` 修饰（默认红=危险）。
> 红色 `<span class="red">`、删除线 `<span class="strike">`、高亮块 `ts-highlight-red` 一页各最多 1 处。
> {{占位符}} 只准换文本；骨架的层级与包裹关系不准增删。
>
> 每页讲稿写在页面末尾：`<div class="notes">…</div>`（观众不可见）。

---

## cover（警示封面）
指纹：hero

用途：开场页。strike 否定旧问题 + 红色强调真问题 + 风险告警框。
适用 role：cover。
内容约束：kicker ≤16 字；标题 ≤3 行、每行 ≤12 字（strike 包被否定的旧说法、red 包核心判断）；lede 30-55 字（代价量级）；alert-box 标题 ≤10 字 + 正文 25-45 字。

```html
<section class="slide" data-layout="cover">
  <div class="ts-stripe"></div>
  <div class="ts-chrome"><span class="ts-alert-tag">{{眉标，如 ai safety · 高优先级}}</span><span class="ts-page">{{页码，如 01 / 08}}</span></div>
  <div class="ts-kicker">{{一句定调，≤16 字}}</div>
  <h1 class="ts-h1">{{第一行，≤12 字}}<br><span class="strike">{{被否定的旧说法，≤12 字}}</span><br>{{转折词，≤6 字}}：<span class="red">{{真正的问题，≤10 字}}</span></h1>
  <p class="ts-sub">{{为什么严重：代价的量级与场景，30-55 字}}</p>
  <div class="ts-alert-box">
    <h3>{{告警标题，≤10 字}}</h3>
    <p>{{风险在哪 + 为什么是现在，25-45 字，<b> 可强调关键词}}</p>
  </div>
  <div class="ts-stripe-b"></div>
  <div class="ts-footer"><span>{{BRIEF 名 · 作者 · 日期}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, ts-stripe, ts-chrome, ts-alert-tag, ts-page, ts-kicker, ts-h1, strike, red, ts-sub, ts-alert-box, ts-stripe-b, ts-footer, notes

---

## divider（章节幕）
指纹：hero

用途：章节过渡。巨大标题（红色强调一个词）+ 一句过门。
适用 role：divider。
内容约束：kicker ≤14 字；标题 ≤10 字 + 红色强调词 2-4 字；lede 20-40 字、可 `<br>` 分两行。

```html
<section class="slide" data-layout="divider">
  <div class="ts-stripe"></div>
  <div class="ts-chrome"><span class="ts-alert-tag amber">{{章节眉标，如 section · 风险分级}}</span><span class="ts-page">{{页码}}</span></div>
  <div style="margin:auto 0">
    <div class="ts-kicker">{{章节名，≤14 字}}</div>
    <h1 class="ts-h1" style="font-size:130px">{{章节标题，≤10 字}} <span class="red">{{强调词，2-4 字}}</span></h1>
    <p class="ts-sub" style="font-size:28px">{{一句过门：本章先说什么、再说什么，20-40 字，可用 <br> 分行}}</p>
  </div>
  <div class="ts-stripe-b"></div>
  <div class="ts-footer"><span>{{section · 小节名}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ts-stripe, ts-chrome, ts-alert-tag, amber, ts-page, ts-kicker, ts-h1, red, ts-sub, ts-stripe-b, ts-footer, notes

---

## risk-levels（分级色卡）
数量：ts-card=3
指纹：cards

用途：3 档分级并列（绿/琥珀/红顶边）+ 一条琥珀告警框。分级/分类内容的主力版式。
适用 role：content。
内容约束：恰好 3 卡（顶边依次 `--ts-green/--ts-amber/--ts-red`）；卡 lbl ≤10 字符 + 档名 2-6 字 + 说明 2-3 行、每行 8-18 字（末行 `<b>` 加色给策略）；alert-box 标题 ≤12 字 + 正文 20-40 字。

```html
<section class="slide" data-layout="risk-levels">
  <div class="ts-stripe"></div>
  <div class="ts-chrome"><span class="ts-alert-tag">{{眉标，如 风险分级 · 3 levels}}</span><span class="ts-page">{{页码}}</span></div>
  <h2 class="ts-h2">{{分级标题，≤12 字}}</h2>
  <div class="ts-grid-3">
    <div class="ts-card" style="border-top:4px solid var(--ts-green)"><div class="lbl">{{L1 · 绿色}}</div><h4>{{档名，2-6 字}}</h4><p>{{哪些行为，8-18 字}}<br>{{出错的代价，8-18 字}}<br><b style="color:var(--ts-green)">策略：{{处理策略，≤8 字}}</b></p></div>
    <div class="ts-card" style="border-top:4px solid var(--ts-amber)"><div class="lbl">{{L2 · 琥珀}}</div><h4>{{档名，2-6 字}}</h4><p>{{哪些行为，8-18 字}}<br>{{出错的代价，8-18 字}}<br><b style="color:var(--ts-amber)">策略：{{处理策略，≤8 字}}</b></p></div>
    <div class="ts-card" style="border-top:4px solid var(--ts-red)"><div class="lbl">{{L3 · 红色}}</div><h4>{{档名，2-6 字}}</h4><p>{{哪些行为，8-18 字}}<br>{{出错的代价，8-18 字}}<br><b style="color:var(--ts-red)">策略：{{处理策略，≤8 字}}</b></p></div>
  </div>
  <div class="ts-alert-box amber">
    <h3>{{一条警告，≤12 字}}</h3>
    <p>{{最常被做错的地方 + 纠正，20-40 字}}</p>
  </div>
  <div class="ts-stripe-b"></div>
  <div class="ts-footer"><span>{{risk · 分级缩写}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ts-stripe, ts-chrome, ts-alert-tag, ts-page, ts-h2, ts-grid-3, ts-card, lbl, ts-alert-box, amber, ts-stripe-b, ts-footer, notes

---

## policy-code（策略代码）
指纹：code

用途：policy/配置代码整页呈现，红线规则用代码管、不用文档管。
适用 role：code。
内容约束：kicker ≤16 字；标题第一行 ≤10 字 + `ts-highlight-red` 强调块 ≤6 字；代码 10-20 行、每行 ≤70 字符；规则键用 `kw`、值用 `st`、注释用 `cm`、绝对禁止项用 `bad`。

```html
<section class="slide" data-layout="policy-code">
  <div class="ts-stripe"></div>
  <div class="ts-chrome"><span class="ts-alert-tag">{{眉标，如 policy as code}}</span><span class="ts-page">{{页码}}</span></div>
  <div class="ts-kicker">{{方法一句话，≤16 字}}</div>
  <h2 class="ts-h2">{{标题第一行，≤10 字}}<br><span class="ts-highlight-red">{{红线强调，≤6 字}}</span></h2>
  <pre class="ts-codebox"><span class="cm"># {{文件名 · 用途注释}}</span>
<span class="kw">{{允许规则键}}</span>:
  - tools: [<span class="st">{{工具 a}}</span>, <span class="st">{{工具 b}}</span>]

<span class="kw">{{需复核规则键}}</span>:
  - tools: [<span class="st">{{工具 c}}</span>, <span class="st">{{工具 d}}</span>]
    reviewer: <span class="st">{{复核方式}}</span>

<span class="kw">{{硬卡规则键}}</span>:
  - tools: [<span class="st">{{高危工具}}</span>]
    unless: <span class="st">{{豁免条件}}</span>

<span class="bad">{{禁止规则键}}:</span>
  - <span class="bad">{{"绝对禁止项 1"}}</span>
  - <span class="bad">{{"绝对禁止项 2"}}</span></pre>
  <div class="ts-stripe-b"></div>
  <div class="ts-footer"><span>{{policy · 缩写}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿：最关键的一条规则}}</div>
</section>
```

合法类名：slide, ts-stripe, ts-chrome, ts-alert-tag, ts-page, ts-kicker, ts-h2, ts-highlight-red, ts-codebox, cm, kw, st, bad, ts-stripe-b, ts-footer, notes

---

## incident-chart（事故图表）
指纹：chart

用途：按月/按类的事故量柱状图 + 分级图例，事故复盘的数据页。
适用 role：data。
内容约束：标题 ≤14 字 + 红色强调词 2-6 字；lede 20-40 字（口径）；柱 3-6 根（rect 高 = 数量 × 20、y = 320 − 高、x 间距 200）、顶标总量；图例 3 项（红/琥珀/绿）+ 2 行注解各 12-24 字。

```html
<section class="slide" data-layout="incident-chart">
  <div class="ts-stripe"></div>
  <div class="ts-chrome"><span class="ts-alert-tag amber">{{眉标，如 incident report · q1}}</span><span class="ts-page">{{页码}}</span></div>
  <h2 class="ts-h2">{{数据标题，≤14 字}} <span class="red">{{强调词，2-6 字}}</span></h2>
  <p class="ts-sub">{{口径：统计范围与来源，20-40 字}}</p>
  <svg viewBox="0 0 1040 360" style="width:100%;max-width:1040px;margin-top:18px" xmlns="http://www.w3.org/2000/svg">
    <g font-family="Inter,sans-serif" font-size="16" fill="#4a4955">
      <line x1="70" y1="320" x2="1000" y2="320" stroke="#eaecf3" stroke-width="2"/>
      <!-- 每周期一柱：rect 高 = 数量 × 20，y = 320 − 高；复制整组加柱（x 间距 200） -->
      <g transform="translate(120,0)">
        <rect x="0" y="220" width="60" height="100" style="fill:var(--ts-green)"/>
        <text x="30" y="345" text-anchor="middle" font-weight="700">{{周期名}}</text>
        <text x="30" y="210" text-anchor="middle" font-weight="800" fill="#14141a">{{数量}}</text>
      </g>
      <g transform="translate(320,0)">
        <rect x="0" y="240" width="60" height="80" style="fill:var(--ts-green)"/>
        <text x="30" y="345" text-anchor="middle" font-weight="700">{{周期名}}</text>
        <text x="30" y="230" text-anchor="middle" font-weight="800" fill="#14141a">{{数量}}</text>
      </g>
      <g transform="translate(520,0)">
        <rect x="0" y="250" width="60" height="70" style="fill:var(--ts-green)"/>
        <text x="30" y="345" text-anchor="middle" font-weight="700">{{周期名}}</text>
        <text x="30" y="240" text-anchor="middle" font-weight="800" fill="#14141a">{{数量}}</text>
      </g>
      <!-- 图例：分级三色，红在上；rect 颜色可与柱色区分表示另一维度 -->
      <g transform="translate(720,60)">
        <rect x="0" y="0" width="16" height="16" style="fill:var(--ts-red)"/><text x="24" y="13" font-weight="700">{{最高级说明，≤12 字}}</text>
        <rect x="0" y="26" width="16" height="16" style="fill:var(--ts-amber)"/><text x="24" y="39" font-weight="700">{{中等级说明，≤12 字}}</text>
        <rect x="0" y="52" width="16" height="16" style="fill:var(--ts-green)"/><text x="24" y="65" font-weight="700">{{最低级说明，≤12 字}}</text>
        <text x="0" y="100" font-size="15" fill="#8a8892">{{注解：拦截/处置情况，12-24 字}}</text>
        <text x="0" y="118" font-size="15" fill="#8a8892">{{注解：最险的一条，12-24 字}}</text>
      </g>
    </g>
  </svg>
  <div class="ts-stripe-b"></div>
  <div class="ts-footer"><span>{{incident · 缩写}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ts-stripe, ts-chrome, ts-alert-tag, amber, ts-page, ts-h2, red, ts-sub, ts-stripe-b, ts-footer, notes

---

## checklist（红队清单）
指纹：table

用途：上线前必过的问题清单，`✓` 已验证 / `!` 待办 两种状态。
适用 role：content / cta。
内容约束：5-7 条；每条 12-30 字、以问句收尾；ok（✓）与待办（!）混排、待办 ≥2 条；标题含红色数字强调。

```html
<section class="slide" data-layout="checklist">
  <div class="ts-stripe"></div>
  <div class="ts-chrome"><span class="ts-alert-tag green">{{眉标，如 red-team checklist}}</span><span class="ts-page">{{页码}}</span></div>
  <h2 class="ts-h2">{{清单标题前半，≤8 字}} <span class="red">{{数字强调，≤8 字}}</span></h2>
  <div class="ts-checklist">
    <div class="ts-check ok"><div class="box">✓</div><div class="txt">{{已验证的问题，12-30 字，问句收尾}}</div></div>
    <div class="ts-check ok"><div class="box">✓</div><div class="txt">{{已验证的问题，12-30 字，问句收尾}}</div></div>
    <div class="ts-check"><div class="box">!</div><div class="txt">{{待办的问题，12-30 字，问句收尾}}</div></div>
    <div class="ts-check ok"><div class="box">✓</div><div class="txt">{{已验证的问题，12-30 字，问句收尾}}</div></div>
    <div class="ts-check"><div class="box">!</div><div class="txt">{{待办的问题，12-30 字，问句收尾}}</div></div>
  </div>
  <div class="ts-stripe-b"></div>
  <div class="ts-footer"><span>{{checklist · 缩写}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿：! 项怎么办}}</div>
</section>
```

合法类名：slide, ts-stripe, ts-chrome, ts-alert-tag, green, ts-page, ts-h2, red, ts-checklist, ts-check, ok, box, txt, ts-stripe-b, ts-footer, notes

---

## tonight（行动三卡）
指纹：cards

用途：CTA。马上就能做的 2-3 件事（编号卡）+ 绿色告警框收一条原则。
适用 role：cta / content。
内容约束：2-3 卡；卡 lbl（编号 · 动作）≤10 字符 + 卡标题两行、每行 ≤8 字 + 说明 15-30 字；绿色 alert-box 标题 ≤12 字 + 正文 20-40 字。

```html
<section class="slide" data-layout="tonight">
  <div class="ts-stripe"></div>
  <div class="ts-chrome"><span class="ts-alert-tag green">{{眉标，如 今晚就能动}}</span><span class="ts-page">{{页码}}</span></div>
  <h2 class="ts-h2">{{行动标题，≤12 字}} <span class="ts-highlight-red">{{强调块，≤4 字}}</span></h2>
  <div class="ts-grid-3">
    <div class="ts-card"><div class="lbl">{{1 · 动作名}}</div><h4>{{卡标题上行，≤8 字}}<br>{{卡标题下行，≤8 字}}</h4><p>{{怎么做，15-30 字}}</p></div>
    <div class="ts-card"><div class="lbl">{{2 · 动作名}}</div><h4>{{卡标题上行，≤8 字}}<br>{{卡标题下行，≤8 字}}</h4><p>{{怎么做，15-30 字}}</p></div>
    <div class="ts-card"><div class="lbl">{{3 · 动作名}}</div><h4>{{卡标题上行，≤8 字}}<br>{{卡标题下行，≤8 字}}</h4><p>{{怎么做，15-30 字}}</p></div>
  </div>
  <div class="ts-alert-box green">
    <h3>{{收束原则，≤12 字}}</h3>
    <p>{{为什么这条原则成立，20-40 字}}</p>
  </div>
  <div class="ts-stripe-b"></div>
  <div class="ts-footer"><span>{{cta · tonight}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ts-stripe, ts-chrome, ts-alert-tag, green, ts-page, ts-h2, ts-highlight-red, ts-grid-3, ts-card, lbl, ts-alert-box, ts-stripe-b, ts-footer, notes

---

## thanks（收尾页）
指纹：hero

用途：收尾。巨大感谢 + 资料获取方式。
适用 role：thanks / cta。
内容约束：kicker ≤14 字；标题 ≤12 字 + 红点分隔 + 尾词 ≤10 字符；lede 20-40 字（资料怎么拿，要真实可执行）。

```html
<section class="slide" data-layout="thanks">
  <div class="ts-stripe"></div>
  <div class="ts-chrome"><span class="ts-alert-tag amber">{{眉标，如 please stay safe}}</span><span class="ts-page">{{页码}}</span></div>
  <div style="margin:auto 0">
    <div class="ts-kicker">{{收尾标记，如 end of brief，≤14 字}}</div>
    <h1 class="ts-h1" style="font-size:140px">{{感谢语，≤12 字}} <span class="red">·</span> {{尾词，≤10 字符}}</h1>
    <p class="ts-sub" style="font-size:24px">{{资料怎么拿：模板/清单/仓库，20-40 字}}</p>
  </div>
  <div class="ts-stripe-b"></div>
  <div class="ts-footer"><span>{{end of brief}}</span><span>{{页码}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, ts-stripe, ts-chrome, ts-alert-tag, amber, ts-page, ts-kicker, ts-h1, red, ts-sub, ts-stripe-b, ts-footer, notes
