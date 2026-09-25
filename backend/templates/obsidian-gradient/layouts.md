# 黑曜渐变 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，
> 只准使用骨架里出现过的类名（base 原语 + 本模板 `oc-*` 专属类）。
> 未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-obsidian-claude-gradient` 作用域前缀生效，骨架里已写全，照抄结构即可。
>
> **结构铁律**：`oc-cbg`（紫蓝 radial 晕染）+ `oc-cgrid`（60px 遮罩网格）两个空气层每页保留；
> 右上角 `oc-snum` 页码每页保留。本模板全部居中排版（slide 自带 flex 居中），
> 渐变字 `<span class="oc-g">` 一页最多 1-2 处、只包 2-6 字的短语，不要整句渐变。
> {{占位符}} 只准换文本；骨架的层级与包裹关系不准增删。
>
> 每页讲稿写在页面末尾：`<div class="notes">…</div>`（观众不可见）。

---

## cover（渐变封面）
指纹：hero

用途：开场页。胶囊眉标 + 渐变字大标题 + 一句话定位 + 能力胶囊。
适用 role：cover。
内容约束：眉标 ≤20 字符（`●` 开头）；标题两行、每行 ≤14 字，渐变 span 只包一个短语；lede 30-55 字（不是什么 + 是什么）；3 个能力胶囊、每个 ≤12 字符。

```html
<section class="slide" data-layout="cover">
  <div class="oc-cbg"></div>
  <div class="oc-cgrid"></div>
  <div class="oc-snum">{{页码，如 01 / 08}}</div>
  <div class="oc-tag">{{● 主题眉标，如 ● OBSIDIAN × CLAUDE · 第二大脑}}</div>
  <h1 class="oc-h1">{{标题第一行，≤14 字}}<br>{{第二行起点，≤8 字}} <span class="oc-g">{{渐变短语，2-6 字}}</span></h1>
  <p class="oc-sub">{{一句话定位：不是什么、而是什么，30-55 字，可用 <br> 分行}}</p>
  <div style="margin-top:32px">
    <span class="oc-pill">{{能力 1，≤12 字符}}</span>
    <span class="oc-pill">{{能力 2，≤12 字符}}</span>
    <span class="oc-pill">{{能力 3，≤12 字符}}</span>
  </div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, oc-cbg, oc-cgrid, oc-snum, oc-tag, oc-h1, oc-g, oc-sub, oc-pill, notes

---

## section（章节幕）
指纹：hero

用途：章节过渡。章节眉标 + 巨大标题（渐变强调一个词）+ 一句过门。
适用 role：divider。
内容约束：眉标 ≤16 字符；标题 ≤16 字符；lede 20-40 字、可 `<br>` 分两行。

```html
<section class="slide" data-layout="section">
  <div class="oc-cbg"></div>
  <div class="oc-cgrid"></div>
  <div class="oc-snum">{{页码}}</div>
  <div class="oc-tag">{{● 章节眉标，如 ● CHAPTER 01}}</div>
  <h1 class="oc-h1" style="font-size:110px">{{章节标题，≤16 字}} <span class="oc-g">{{强调词，2-6 字}}</span></h1>
  <p class="oc-sub">{{一句过门：本章回答什么问题，20-40 字，可用 <br> 分行}}</p>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, oc-cbg, oc-cgrid, oc-snum, oc-tag, oc-h1, oc-g, oc-sub, notes

---

## compare（双卡对照）
数量：oc-card=2
指纹：cards

用途：两个方案/工具的对照（优劣、A vs B）。左卡中性、右卡描边强调立场。
适用 role：content。
内容约束：恰好 2 卡；卡徽标 4-10 字符 + 卡标题 4-10 字 + 说明 3 条、每条 10-22 字（`•` 分行）；底部 hl 洞察 25-45 字（一句点破）。

```html
<section class="slide" data-layout="compare">
  <div class="oc-cbg"></div>
  <div class="oc-cgrid"></div>
  <div class="oc-snum">{{页码}}</div>
  <div class="oc-tag">{{● 眉标，如 ● COMPARE}}</div>
  <h2 class="oc-h2">{{对比标题主干，≤12 字}} <span class="oc-g">{{重点词}}</span></h2>
  <div class="oc-grid-2">
    <div class="oc-card">
      <span class="oc-badge oc-bb">{{方案 A 名，4-10 字符}}</span>
      <h4>{{卡标题，4-10 字}}</h4>
      <p>• {{要点 1，10-22 字}}<br>• {{要点 2，10-22 字}}<br>• {{要点 3，10-22 字}}</p>
    </div>
    <div class="oc-card" style="border-color:rgba(var(--oc-accent2-rgb),.35);background:rgba(var(--oc-accent-rgb),.05)">
      <span class="oc-badge oc-bp">{{方案 B 名，4-10 字符}}</span>
      <h4>{{卡标题，4-10 字}}</h4>
      <p>• {{要点 1，10-22 字}}<br>• {{要点 2，10-22 字}}<br>• {{要点 3，10-22 字}}</p>
    </div>
  </div>
  <div class="oc-hl" style="margin-top:26px">{{关键洞察：一句点破，25-45 字}}</div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, oc-cbg, oc-cgrid, oc-snum, oc-tag, oc-h2, oc-g, oc-grid-2, oc-card, oc-badge, oc-bb, oc-bp, h4, oc-hl, notes

---

## steps（步骤清单）
指纹：table

用途：3-5 步的上手流程/配置步骤，编号圆点 + 标题 + 说明。
适用 role：content / cta。
内容约束：3-5 步；每步标题 8-16 字 + 说明 15-32 字；编号连续不跳。

```html
<section class="slide" data-layout="steps">
  <div class="oc-cbg"></div>
  <div class="oc-cgrid"></div>
  <div class="oc-snum">{{页码}}</div>
  <div class="oc-tag">{{● 眉标，如 ● SETUP · 4 STEPS}}</div>
  <h2 class="oc-h2">{{流程标题，≤14 字}}</h2>
  <div class="oc-steps">
    <div class="oc-step"><div class="oc-sn">1</div><div class="oc-sc"><h4>{{步骤标题，8-16 字}}</h4><p>{{说明，15-32 字}}</p></div></div>
    <div class="oc-step"><div class="oc-sn">2</div><div class="oc-sc"><h4>{{步骤标题，8-16 字}}</h4><p>{{说明，15-32 字}}</p></div></div>
    <div class="oc-step"><div class="oc-sn">3</div><div class="oc-sc"><h4>{{步骤标题，8-16 字}}</h4><p>{{说明，15-32 字}}</p></div></div>
    <div class="oc-step"><div class="oc-sn">4</div><div class="oc-sc"><h4>{{步骤标题，8-16 字}}</h4><p>{{说明，15-32 字}}</p></div></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, oc-cbg, oc-cgrid, oc-snum, oc-tag, oc-h2, oc-steps, oc-step, oc-sn, oc-sc, h4, notes

---

## config（配置代码）
指纹：code

用途：配置文件/命令整页呈现（oc-code 承载），配一句验证结果。
适用 role：code。
内容约束：代码 8-18 行、每行 ≤70 字符；键用 `cc`、字符串值用 `cs`、注释用 `cm`、强调键用 `ca`；代码下方一句验证结果 12-30 字。

```html
<section class="slide" data-layout="config">
  <div class="oc-cbg"></div>
  <div class="oc-cgrid"></div>
  <div class="oc-snum">{{页码}}</div>
  <div class="oc-tag">{{● 眉标，如 ● MCP CONFIG}}</div>
  <h2 class="oc-h2">{{文件名，≤20 字符}}</h2>
  <pre class="oc-code"><span class="cm">// {{文件路径注释}}</span>
{
  <span class="cc">"{{顶层键}}"</span>: {
    <span class="cc">"{{配置键}}"</span>: <span class="cs">"{{值}}"</span>,
    <span class="cc">"{{配置键}}"</span>: [<span class="cs">"{{值 1}}"</span>, <span class="cs">"{{值 2}}"</span>],
    <span class="cc">"env"</span>: {
      <span class="cc">"{{环境变量键}}"</span>: <span class="cs">"{{值，密钥打码}}"</span>
    }
  }
}</pre>
  <p class="oc-sub" style="margin-top:18px">{{验证结果一句：看到什么算配好了，12-30 字}}</p>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, oc-cbg, oc-cgrid, oc-snum, oc-tag, oc-h2, oc-code, cm, cc, cs, ca, cp, oc-sub, notes

---

## metrics（数据战果）
数量：oc-card=3
指纹：chart

用途：3 个大数字并列（用量/收益/降幅）+ 一句体感结论。
适用 role：data / content。
内容约束：恰好 3 卡；数字 ≤6 字符 + 指标说明 6-16 字；hl 结论 25-45 字；数字必须有口径或时间范围，没有就写"估算"。

```html
<section class="slide" data-layout="metrics">
  <div class="oc-cbg"></div>
  <div class="oc-cgrid"></div>
  <div class="oc-snum">{{页码}}</div>
  <div class="oc-tag">{{● 眉标，如 ● 3 MONTHS IN}}</div>
  <h2 class="oc-h2">{{数据标题，≤14 字}}</h2>
  <div class="oc-grid-3" style="margin-top:28px">
    <div class="oc-card" style="text-align:center"><div class="oc-big oc-g" style="font-size:80px">{{数字 ≤6 字符}}</div><p>{{指标说明，6-16 字}}</p></div>
    <div class="oc-card" style="text-align:center"><div class="oc-big oc-g" style="font-size:80px">{{数字 ≤6 字符}}</div><p>{{指标说明，6-16 字}}</p></div>
    <div class="oc-card" style="text-align:center"><div class="oc-big oc-g" style="font-size:80px">{{数字 ≤6 字符}}</div><p>{{指标说明，6-16 字}}</p></div>
  </div>
  <div class="oc-hl" style="margin-top:26px">{{体感结论：最大收益是什么，25-45 字}}</div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, oc-cbg, oc-cgrid, oc-snum, oc-tag, oc-h2, oc-g, oc-grid-3, oc-card, oc-big, oc-hl, notes

---

## quote-cta（金句行动）
指纹：quote

用途：一句立场金句 + 署名 + 下载/行动胶囊行。
适用 role：quote / cta。
内容约束：引文两行、每行 8-20 字（渐变 span 包核心短语）；署名 ≤16 字；2-4 个行动胶囊、每个 ≤14 字符。引文要真实立场，不要编名言。

```html
<section class="slide" data-layout="quote-cta">
  <div class="oc-cbg"></div>
  <div class="oc-cgrid"></div>
  <div class="oc-snum">{{页码}}</div>
  <div class="oc-tag">{{● 眉标，如 ● CTA · 今晚可以做}}</div>
  <div class="oc-quote">
    <blockquote>{{引文第一行，8-20 字}}<br>{{第二行起点}} <span class="oc-g">{{渐变短语，2-6 字}}</span>。</blockquote>
    <div class="attr">— {{署名，≤16 字}}</div>
  </div>
  <div style="margin-top:36px">
    <span class="oc-pill">{{行动 1，≤14 字符}}</span>
    <span class="oc-pill">{{行动 2，≤14 字符}}</span>
    <span class="oc-pill">{{行动 3，≤14 字符}}</span>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, oc-cbg, oc-cgrid, oc-snum, oc-tag, oc-quote, oc-g, attr, oc-pill, notes

---

## thanks（致谢页）
指纹：hero

用途：收尾页。渐变大字 Thanks + 资料仓库链接。
适用 role：thanks / cta。
内容约束：大标题 ≤10 字符；lede 20-45 字（资料链接要真实，没有就写后续计划）。

```html
<section class="slide" data-layout="thanks">
  <div class="oc-cbg"></div>
  <div class="oc-cgrid"></div>
  <div class="oc-snum">{{页码}}</div>
  <div class="oc-big oc-g">{{收尾词，如 Thanks.，≤10 字符}}</div>
  <p class="oc-sub" style="margin-top:26px">{{资料在哪：模板/脚本/仓库，20-45 字}}</p>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, oc-cbg, oc-cgrid, oc-snum, oc-big, oc-g, oc-sub, notes
