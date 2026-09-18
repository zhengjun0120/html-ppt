# 工业粗野 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-brutalist-bold` 作用域前缀生效，骨架里已写全，照抄结构即可。
>
> **风格铁律**：浅色工业印刷——纸底碳黑 + 唯一强调色，全部直角、无渐变、无投影、无半透明、禁 emoji。
> 红色（或当前变体强调色）只出现在：b-rule 粗线、b-mark 红块、b-strike 删除线、b-val-red 数据高亮、
> b-hot 热板、b-blade-red 红色编号、b-term 红竖线。别处不准上色。
> 顶栏 `.b-top` 与页脚 `.deck-footer` 是每页的机械外框，照抄结构、只换 {{占位符}}。
> 满版页（cover / section / closing）用 `slide full` 垂直居中；其余页沿用骨架即可。

---

## cover（工业封面）
指纹：hero

用途：开场页。文件编号栏 + 巨字主标 + 红色通栏粗线 + 要素标签，像一份贴在车间门口的正式通报。
适用 role：cover。
内容约束：主标题 ≤10 字（1 个关键词外套 `<mark class="b-mark">`）；lede 30-55 字（什么事 + 什么时间 + 谁必须执行）；3 个要素标签各 ≤12 字；顶栏 2 条 mono 编号。

```html
<section class="slide full" data-layout="cover">
  <div class="b-top">
    <span class="b-meta">{{文件编号，如 DOC NO. FD-2026-0924}}</span>
    <span class="b-meta">{{单位/库区代码，如 UNIT D-01 · 3号库区}}</span>
    <span class="b-code">||||| || |||||| |||</span>
  </div>
  <p class="kicker mt-l">{{文种 · 效力，如 内部通报 · 签发即生效}}</p>
  <h1 class="h1 mt-s">{{主标题 ≤10 字，1 个关键词外套 <mark class="b-mark">}}</h1>
  <div class="b-rule"></div>
  <p class="lede mt-m">{{什么事 + 什么时间 + 谁必须执行，30-55 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="b-tag">{{关键要素 1，≤12 字}}</span>
    <span class="b-tag">{{关键要素 2，≤12 字}}</span>
    <span class="b-tag">{{关键要素 3，≤12 字}}</span>
  </div>
  <div class="b-num">{{背景巨号：月份或编号，≤2 位数字}}</div>
  <div class="deck-footer"><span class="b-meta">{{签发方 · 签发日期}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句：为什么现在发、要求谁做什么}}</div>
</section>
```

合法类名：slide, full, b-top, b-meta, b-code, kicker, mt-l, h1, mt-s, b-mark, b-rule, mt-m, lede, row, b-tag, b-num, deck-footer, slide-number, notes

---

## manifest（科目清单）
指纹：table

用途：验收清单表。序号 + 科目 + 验收口径 + 责任人四列，4-5 行逐项过，清单/议程/分工的主力版式。
适用 role：toc / content。
内容约束：4-5 行；标题 ≤12 字；每行科目 ≤8 字 + 验收口径 14-28 字（写时限与合格判据）+ 责任岗位 ≤10 字。

```html
<section class="slide" data-layout="manifest">
  <div class="b-top">
    <span class="b-meta">{{节号 · 主题，如 SEC 01 · 演练科目}}</span>
    <span class="b-meta">{{文件编号 · 版本，如 FD-2026-0924 · REV 2.6}}</span>
    <span class="b-code">|||| ||||| ||</span>
  </div>
  <h2 class="h2 mt-l">{{标题：这张清单在验收什么，≤12 字}}</h2>
  <div class="b-table mt-m">
    <div class="b-tr">
      <span class="b-no">01</span>
      <h4 class="b-tt">{{科目名，≤8 字}}</h4>
      <p class="b-td">{{验收口径：多少时间内完成什么、怎么算过，14-28 字}}</p>
      <span class="b-own">{{责任岗位 + 人名，≤10 字}}</span>
    </div>
    <div class="b-tr">
      <span class="b-no">02</span>
      <h4 class="b-tt">{{科目名，≤8 字}}</h4>
      <p class="b-td">{{验收口径，14-28 字}}</p>
      <span class="b-own">{{责任岗位，≤10 字}}</span>
    </div>
    <div class="b-tr">
      <span class="b-no">03</span>
      <h4 class="b-tt">{{科目名，≤8 字}}</h4>
      <p class="b-td">{{验收口径，14-28 字}}</p>
      <span class="b-own">{{责任岗位，≤10 字}}</span>
    </div>
    <div class="b-tr">
      <span class="b-no">04</span>
      <h4 class="b-tt">{{科目名，≤8 字}}</h4>
      <p class="b-td">{{验收口径，14-28 字}}</p>
      <span class="b-own">{{责任岗位，≤10 字}}</span>
    </div>
  </div>
  <div class="deck-footer"><span class="b-meta">{{落款 · 修订号}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, b-top, b-meta, b-code, h2, mt-l, b-table, mt-m, b-tr, b-no, b-tt, b-td, b-own, deck-footer, slide-number, notes

---

## section（章节幕）
指纹：hero

用途：章节过渡。进度标记 + 巨字章节名 + 红色粗线 + 描边巨号数字压在右上角。
适用 role：divider。
内容约束：kicker 写 `[ SECTION N / M ]` 式进度 ≤16 字符；标题 ≤10 字；lede 22-45 字（这一节要执行的硬要求）；2 个看点标签各 ≤10 字。

```html
<section class="slide full" data-layout="section">
  <p class="kicker">[ SECTION {{N}} / {{M}} ]</p>
  <h1 class="h1 mt-s">{{章节标题，≤10 字}}</h1>
  <div class="b-rule"></div>
  <p class="lede mt-m">{{这一节要执行的硬要求 / 回答的一个问题，22-45 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="b-tag">{{本节看点 1，≤10 字}}</span>
    <span class="b-tag">{{本节看点 2，≤10 字}}</span>
  </div>
  <div class="b-num">{{章节数字，1-2 位}}</div>
  <div class="deck-footer"><span class="b-meta">{{落款 · 修订号}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, kicker, h1, mt-s, b-rule, mt-m, lede, row, mt-l, b-tag, b-num, deck-footer, b-meta, slide-number, notes

---

## split-spec（参数对照）
指纹：split

用途：左栏一句硬判断把要求讲透，右栏 2×2 参数铭牌（数值 + 单位 + 口径），其中 1 块 b-hot 红板给红线指标。
适用 role：content / data。
内容约束：判断句 ≤12 字；lede 30-55 字；2 个动作标签各 ≤8 字（禁做项用 `b-tag b-strike` 划掉）；恰好 4 块参数板：数值 ≤4 字符 + 参数名 ≤8 字 + 口径 12-24 字。

```html
<section class="slide" data-layout="split-spec">
  <div class="b-top">
    <span class="b-meta">{{节号 · 主题}}</span>
    <span class="b-meta">{{文件编号 · 版本}}</span>
    <span class="b-code">||| |||| |||</span>
  </div>
  <div class="b-cols mt-l">
    <div>
      <p class="kicker">{{左栏标签，如 EVACUATION SPEC}}</p>
      <h2 class="h2 mt-s">{{左栏判断句，≤12 字}}</h2>
      <p class="lede mt-m">{{把要求讲透：时限、走法、例外，30-55 字}}</p>
      <div class="row mt-l" style="gap:14px">
        <span class="b-tag">{{动作要领，≤8 字}}</span>
        <span class="b-tag b-strike">{{明令禁止的做法，≤8 字}}</span>
      </div>
    </div>
    <div class="b-blade">
      <div class="b-plate b-hot">
        <span class="b-kv">{{参数名，≤8 字}}</span>
        <div class="b-val b-val-red">{{数值 ≤4 字符}}<span class="b-unit">{{单位}}</span></div>
        <p class="b-note">{{口径：从什么到什么、怎么算，12-24 字}}</p>
      </div>
      <div class="b-plate">
        <span class="b-kv">{{参数名，≤8 字}}</span>
        <div class="b-val">{{数值 ≤4 字符}}<span class="b-unit">{{单位}}</span></div>
        <p class="b-note">{{口径，12-24 字}}</p>
      </div>
      <div class="b-plate">
        <span class="b-kv">{{参数名，≤8 字}}</span>
        <div class="b-val">{{数值 ≤4 字符}}<span class="b-unit">{{单位}}</span></div>
        <p class="b-note">{{口径，12-24 字}}</p>
      </div>
      <div class="b-plate">
        <span class="b-kv">{{参数名，≤8 字}}</span>
        <div class="b-val">{{数值 ≤4 字符}}<span class="b-unit">{{单位}}</span></div>
        <p class="b-note">{{口径，12-24 字}}</p>
      </div>
    </div>
  </div>
  <div class="deck-footer"><span class="b-meta">{{落款 · 修订号}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, b-top, b-meta, b-code, b-cols, mt-l, kicker, h2, mt-s, lede, mt-m, row, b-tag, b-strike, b-blade, b-plate, b-hot, b-kv, b-val, b-val-red, b-unit, b-note, deck-footer, slide-number, notes

---

## stat-plate（数据铭牌）
指纹：chart

用途：三块大数字铭牌回答一个底数问题。第一块数值红色高亮，页尾必须有 mono 来源行。
适用 role：data / content。
内容约束：恰好 3 板；数值 ≤4 字符 + 单位 + 指标名 ≤8 字 + 口径 12-24 字（写清点日期/范围）；来源行写台账/系统名 + 截止日期。

```html
<section class="slide" data-layout="stat-plate">
  <div class="b-top">
    <span class="b-meta">{{节号 · 主题}}</span>
    <span class="b-meta">{{文件编号 · 版本}}</span>
    <span class="b-code">||||| |||| ||</span>
  </div>
  <p class="kicker mt-l">{{数据语境，如 BASELINE · 9 月台账}}</p>
  <h2 class="h2 mt-s">{{这些数字回答什么问题，≤14 字}}</h2>
  <div class="b-blade b-blade-3 mt-m">
    <div class="b-plate">
      <div class="b-val b-val-red">{{数值 ≤4 字符}}<span class="b-unit">{{单位}}</span></div>
      <h4 class="b-tt">{{指标名，≤8 字}}</h4>
      <p class="b-note">{{口径/清点日期，12-24 字}}</p>
    </div>
    <div class="b-plate">
      <div class="b-val">{{数值 ≤4 字符}}<span class="b-unit">{{单位}}</span></div>
      <h4 class="b-tt">{{指标名，≤8 字}}</h4>
      <p class="b-note">{{口径，12-24 字}}</p>
    </div>
    <div class="b-plate">
      <div class="b-val">{{数值 ≤4 字符}}<span class="b-unit">{{单位}}</span></div>
      <h4 class="b-tt">{{指标名，≤8 字}}</h4>
      <p class="b-note">{{口径，12-24 字}}</p>
    </div>
  </div>
  <p class="b-meta mt-m">来源：{{台账/系统名 + 截止日期}}</p>
  <div class="deck-footer"><span class="b-meta">{{落款 · 修订号}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, b-top, b-meta, b-code, kicker, mt-l, h2, mt-s, b-blade, b-blade-3, mt-m, b-plate, b-val, b-val-red, b-unit, b-tt, b-note, deck-footer, slide-number, notes

---

## code-run（指令时序）
指纹：code

用途：时间戳调度序列 / 执行日志（b-term 承载），按秒执行的硬流程一页讲完。
适用 role：code / content。
内容约束：kicker 写日志名与规模（如 `RUNBOOK.LOG · 7 步`）；标题 ≤14 字；代码 8-14 行、每行 ≤64 字符，指令词用 `<span class="kw">` 加粗、注释用 `.cmt`、引值用 `.str`；`<` `>` `&` 要转义；lede 22-45 字写关键一步失败会怎样。

```html
<section class="slide" data-layout="code-run">
  <div class="b-top">
    <span class="b-meta">{{节号 · 主题}}</span>
    <span class="b-meta">{{文件编号 · 版本}}</span>
    <span class="b-code">|| ||||| ||||</span>
  </div>
  <p class="kicker mt-l">{{日志名 · 规模，如 RUNBOOK.LOG · 7 步}}</p>
  <h2 class="h2 mt-s">{{标题：这条时序在保证什么，≤14 字}}</h2>
  <div class="b-term mt-m"><pre style="margin:0"><span class="cmt">{{首行注释：对齐/阅读约定，≤40 字符}}</span>
<span class="kw">{{HH:MM:SS 指令大写}}</span> <span class="str">{{"参数"}}</span>  <span class="cmt">// {{谁做什么、判据是什么}}</span>
{{…… 8-14 行同构，每行 ≤64 字符；时间戳升序，关键行指令用 .kw 加粗 ……}</pre></div>
  <p class="lede mt-m">{{关键一步失败会怎样、谁来兜，22-45 字}}</p>
  <div class="deck-footer"><span class="b-meta">{{落款 · 修订号}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, b-top, b-meta, b-code, kicker, mt-l, h2, mt-s, b-term, mt-m, kw, str, cmt, lede, deck-footer, slide-number, notes

---

## redline-cards（红线三卡）
指纹：cards

用途：刀锋细线三卡并列，红色编号 + 红线名 + 判罚口径。违规/禁区/否决项的主力版式。
适用 role：content。
内容约束：恰好 3 卡；红线名 ≤6 字；判罚口径 22-45 字（写清怎么判罚 + 已经做了什么准备）；lede 20-40 字写兜底裁决归属。

```html
<section class="slide" data-layout="redline-cards">
  <div class="b-top">
    <span class="b-meta">{{节号 · 主题}}</span>
    <span class="b-meta">{{文件编号 · 版本}}</span>
    <span class="b-code">|||| ||| |||||</span>
  </div>
  <p class="kicker mt-l">{{引导语，如 RED LINES · 判罚口径}}</p>
  <h2 class="h2 mt-s">{{标题：什么行为判不合格，≤14 字}}</h2>
  <div class="b-blade b-blade-3 b-blade-red mt-m">
    <div class="b-plate">
      <span class="b-no">01</span>
      <h4 class="b-tt">{{红线名，≤6 字}}</h4>
      <p class="b-td">{{判罚口径 + 已做的准备，22-45 字}}</p>
    </div>
    <div class="b-plate">
      <span class="b-no">02</span>
      <h4 class="b-tt">{{红线名，≤6 字}}</h4>
      <p class="b-td">{{判罚口径 + 已做的准备，22-45 字}}</p>
    </div>
    <div class="b-plate">
      <span class="b-no">03</span>
      <h4 class="b-tt">{{红线名，≤6 字}}</h4>
      <p class="b-td">{{判罚口径 + 已做的准备，22-45 字}}</p>
    </div>
  </div>
  <p class="lede mt-m">{{兜底：拿不准先问谁、谁现场裁决，20-40 字}}</p>
  <div class="deck-footer"><span class="b-meta">{{落款 · 修订号}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, b-top, b-meta, b-code, kicker, mt-l, h2, mt-s, b-blade, b-blade-3, b-blade-red, mt-m, b-plate, b-no, b-tt, b-td, lede, deck-footer, slide-number, notes

---

## closing（收尾铭牌）
指纹：quote

用途：收尾页。巨字断言（红块压关键词）+ 下一步安排 + 签发栏，整页只说一件结论。
适用 role：thanks / quote。
内容约束：断言拆两行、每行 ≤8 字（1 个关键词外套 `<mark class="b-mark">`）；lede 22-45 字（成绩何时出、不合格怎么办）；签发栏 2 条 mono + 条码。

```html
<section class="slide full" data-layout="closing">
  <p class="kicker">{{文末编号，如 FD-2026-0924 · END OF DOC}}</p>
  <h1 class="h1 mt-s">{{断言第一行 ≤8 字}}<br>{{断言第二行，含 <mark class="b-mark">关键词}}</h1>
  <div class="b-rule"></div>
  <p class="lede mt-m">{{下一步：成绩何时出、不合格怎么办，22-45 字}}</p>
  <div class="b-sig mt-l">
    <span class="b-meta">{{签发单位}}</span>
    <span class="b-meta">{{修订号 · 签发日期}}</span>
    <span class="b-code">|||||| ||| ||||</span>
  </div>
  <div class="deck-footer"><span class="b-meta">{{落款}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, kicker, h1, mt-s, b-mark, b-rule, mt-m, lede, b-sig, mt-l, b-meta, b-code, deck-footer, slide-number, notes
