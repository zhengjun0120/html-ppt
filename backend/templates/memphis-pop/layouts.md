# 孟菲斯波普 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属 `mp-` 类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-memphis-pop` 作用域生效，骨架里已写全，照抄结构即可。
> 每页右上角的几何碎屑与页底锯齿彩带由模板自动衬底（z-index 在内容之下），骨架与正文都不用画。
>
> **气质铁律（热闹来自装饰，不来自密度）**：撞色只许来自琥珀黄 / 薰衣草紫 / 靛蓝三主角，
> 墨蓝黑 #0F172A 压阵；粗黑描边与硬投影只属于卡（mp-card / mp-stat）与按钮（mp-btn）。
> 黑底题签（mp-kicker）每页至多一枚；正文区保持低到中密度，一页最多两组内容块 + 页脚；
> 900 字重大标题（h1/h2）每页最多一组。

---

## cover（派对开幕）
指纹：hero

用途：开场页。黑底题签 + 900 字重大标题 + 波浪线 + 一句定位，信息签收在标题下。
适用 role：cover。
内容约束：主标题 ≤16 字可两行；lede 30-60 字；题签 ≤14 字；信息签各 ≤12 字。

```html
<section class="slide full" data-layout="cover">
  <p class="mp-kicker">{{题签，≤14 字，如 策展方案 · 2026 春}}</p>
  <h1 class="h1 mt-m">{{主标题，≤14 字，可 <br> 分两行}}</h1>
  <div class="mp-wave mt-s"></div>
  <p class="lede mt-m" style="max-width:52ch">{{一句话定位，30-60 字}}</p>
  <div class="row mt-l" style="gap:22px">
    <span class="mp-tag">{{信息 1，≤12 字}}</span>
    <span class="mp-tag mp-tag-solid">{{信息 2，≤10 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿 2-3 句}}</div>
</section>
```

合法类名：slide, full, mp-kicker, h1, mt-m, mp-wave, mt-s, lede, mt-l, row, mp-tag, mp-tag-solid, deck-footer, slide-number, notes

---

## contents（场刊目录）
指纹：table
数量：mp-item=4

用途：议程页。一块粗描边大卡里放 4 行篇目：编号圆球 + 篇名 + 一句说明。
适用 role：toc。
内容约束：恰好 4 行；篇名 ≤8 字；说明 14-26 字。

```html
<section class="slide" data-layout="contents">
  <p class="mp-kicker">{{引导语，如 今日卷目}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="mp-card mt-l" style="margin-top:44px">
    <div class="mp-item"><span class="mp-n">壹</span><span class="mp-t">{{篇名，≤8 字}}</span><span class="mp-d">{{说明，14-26 字}}</span></div>
    <div class="mp-item"><span class="mp-n">贰</span><span class="mp-t">{{篇名，≤8 字}}</span><span class="mp-d">{{说明，14-26 字}}</span></div>
    <div class="mp-item"><span class="mp-n">叁</span><span class="mp-t">{{篇名，≤8 字}}</span><span class="mp-d">{{说明，14-26 字}}</span></div>
    <div class="mp-item"><span class="mp-n">肆</span><span class="mp-t">{{篇名，≤8 字}}</span><span class="mp-d">{{说明，14-26 字}}</span></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mp-kicker, h2, mt-m, mp-card, mt-l, mp-item, mp-n, mp-t, mp-d, deck-footer, slide-number, notes

---

## keynotes（三卡宣言）
指纹：cards
数量：mp-card=3

用途：恰好三张粗描边卡，投影三色轮换。每张：黑底题签 + 小标题 + 两句说明。
适用 role：content。
内容约束：恰好 3 卡；题签 ≤8 字；小标题 ≤8 字；说明 22-44 字。

```html
<section class="slide" data-layout="keynotes">
  <p class="mp-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤10 字}}</h2>
  <div class="grid g3 mt-l" style="gap:36px;margin-top:48px">
    <div class="mp-card"><span class="mp-tag mp-tag-solid">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:18px;line-height:1.8;color:#374151">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="mp-card"><span class="mp-tag mp-tag-solid">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:18px;line-height:1.8;color:#374151">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
    <div class="mp-card"><span class="mp-tag mp-tag-solid">{{题签，≤8 字}}</span><h4 class="mt-m">{{小标题，≤8 字}}</h4><p class="mt-s" style="font-size:18px;line-height:1.8;color:#374151">{{说明：一句论断 + 一句展开，22-44 字}}</p></div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mp-kicker, h2, mt-m, grid, g3, mt-l, mp-card, mp-tag, mp-tag-solid, h4, mt-m, mt-s, deck-footer, slide-number, notes

---

## split（左文右列）
指纹：split
数量：mp-step=3

用途：左边把一件事讲透（lede + 补充 + 签），右边琥珀投影粗描边卡装三行步骤或条目。
适用 role：content。
内容约束：左 lede 36-70 字 + 补充 18-40 字 + 2 个签；右恰好 3 行、每行 12-26 字。

```html
<section class="slide" data-layout="split">
  <p class="mp-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{标题，≤12 字}}</h2>
  <div class="grid g2 mt-l" style="gap:72px;margin-top:44px;align-items:start">
    <div>
      <p class="lede">{{叙述：怎么做、为什么，36-70 字}}</p>
      <p class="mt-m" style="font-size:18px;line-height:1.8;color:#374151">{{补充：判断标准或代价，18-40 字}}</p>
      <div class="row mt-l" style="gap:16px">
        <span class="mp-tag">{{要点 1，≤8 字}}</span>
        <span class="mp-tag">{{要点 2，≤8 字}}</span>
      </div>
    </div>
    <div class="mp-card mp-card-amber">
      <div class="mp-step"><span class="mp-n">1</span><p class="mp-mini-t">{{一步，12-26 字}}</p></div>
      <div class="mp-step"><span class="mp-n">2</span><p class="mp-mini-t">{{一步，12-26 字}}</p></div>
      <div class="mp-step"><span class="mp-n">3</span><p class="mp-mini-t">{{一步，12-26 字}}</p></div>
    </div>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mp-kicker, h2, mt-m, grid, g2, mt-l, lede, mt-m, row, mp-tag, mp-card, mp-card-amber, mp-step, mp-n, mp-mini-t, deck-footer, slide-number, notes

---

## metrics（三个大数字）
指纹：chart
数量：mp-stat=3

用途：三个关键数据。粗描边数字卡是主视觉，投影三色轮换；口径写进卡内，来源写在页尾。
适用 role：data。
内容约束：恰好 3 个；数值 ≤6 字符；指标名 ≤8 字；口径 14-30 字；必须标来源。

```html
<section class="slide" data-layout="metrics">
  <p class="mp-kicker">{{数据语境，如 备展进度 · 数据先行}}</p>
  <h2 class="h2 mt-m">{{这些数字回答什么，≤12 字}}</h2>
  <div class="grid g3 mt-l" style="gap:40px;margin-top:48px">
    <div class="mp-stat"><div class="mp-stat-v">{{数值 ≤6 字符}}<span class="mp-stat-u">{{单位}}</span></div><div class="mp-stat-l">{{指标名，≤8 字}}</div><p class="mp-stat-note">{{口径，14-30 字}}</p></div>
    <div class="mp-stat"><div class="mp-stat-v">{{数值 ≤6 字符}}<span class="mp-stat-u">{{单位}}</span></div><div class="mp-stat-l">{{指标名，≤8 字}}</div><p class="mp-stat-note">{{口径，14-30 字}}</p></div>
    <div class="mp-stat"><div class="mp-stat-v">{{数值 ≤6 字符}}<span class="mp-stat-u">{{单位}}</span></div><div class="mp-stat-l">{{指标名，≤8 字}}</div><p class="mp-stat-note">{{口径，14-30 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:16px;color:#6B7280;letter-spacing:.06em">来源：{{出处与统计区间，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mp-kicker, h2, mt-m, grid, g3, mt-l, mp-stat, mp-stat-v, mp-stat-u, mp-stat-l, mp-stat-note, mt-m, deck-footer, slide-number, notes

---

## quote（巨字引言）
指纹：quote

用途：整页一句引言。900 字重大字 + 出处 + 两个支撑签，像海报上的口号。
适用 role：quote。
内容约束：引文 10-28 字；出处 ≤22 字且真实；签各 ≤8 字。

```html
<section class="slide" data-layout="quote">
  <p class="mp-kicker">{{语境，如 米兰 · 1981}}</p>
  <p class="mp-quote mt-l" style="margin-top:40px">「{{引文，10-28 字}}」</p>
  <p class="mp-src mt-m">—— {{出处：人与来源，≤22 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="mp-tag mp-tag-solid">{{支撑点 1，≤8 字}}</span>
    <span class="mp-tag">{{支撑点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mp-kicker, mp-quote, mt-l, mp-src, mt-m, row, mp-tag, mp-tag-solid, deck-footer, slide-number, notes

---

## divider（章节转场）
指纹：hero

用途：章节过渡。进度题签 + 900 字重大字章节名 + 一个过渡问题 + 两个看点签。
适用 role：divider。
内容约束：章节名 ≤10 字；过渡问题 22-52 字；签各 ≤8 字。

```html
<section class="slide full" data-layout="divider">
  <p class="mp-kicker">{{进度，如 第二章 · 反叛的语言}}</p>
  <h1 class="h1 mt-m">{{章节标题，≤10 字}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{这一部分回答的一个问题，22-44 字}}</p>
  <div class="row mt-l" style="gap:16px">
    <span class="mp-tag mp-tag-solid">{{看点 1，≤8 字}}</span>
    <span class="mp-tag">{{看点 2，≤8 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, mp-kicker, h1, mt-m, lede, mt-l, row, mp-tag, mp-tag-solid, deck-footer, slide-number, notes

---

## moments（节点时间线）
指纹：chart
数量：mp-tl-item=4

用途：3-4 个节点的横向时间线：编号圆球 + 时间点 + 一句事件，虚线串联。
适用 role：content / data。
内容约束：3-4 节点；时间点 ≤10 字符；事件 12-26 字；页尾一句口径 14-40 字。

```html
<section class="slide" data-layout="moments">
  <p class="mp-kicker">{{引导语}}</p>
  <h2 class="h2 mt-m">{{时间线标题，≤12 字}}</h2>
  <div class="mp-tl mt-l" style="margin-top:56px">
    <div class="mp-tl-item"><div class="mp-tl-dot">1</div><div class="mp-tl-t">{{时间点，≤10 字符}}</div><p class="mp-tl-d">{{事件，12-26 字}}</p></div>
    <div class="mp-tl-item"><div class="mp-tl-dot">2</div><div class="mp-tl-t">{{时间点，≤10 字符}}</div><p class="mp-tl-d">{{事件，12-26 字}}</p></div>
    <div class="mp-tl-item"><div class="mp-tl-dot">3</div><div class="mp-tl-t">{{时间点，≤10 字符}}</div><p class="mp-tl-d">{{事件，12-26 字}}</p></div>
    <div class="mp-tl-item"><div class="mp-tl-dot">4</div><div class="mp-tl-t">{{时间点，≤10 字符}}</div><p class="mp-tl-d">{{事件，12-26 字}}</p></div>
  </div>
  <p class="mt-m" style="font-size:16px;color:#6B7280;letter-spacing:.06em">{{一句读法或口径，14-40 字}}</p>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{页码}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, mp-kicker, h2, mt-m, mp-tl, mt-l, mp-tl-item, mp-tl-dot, mp-tl-t, mp-tl-d, mt-m, deck-footer, slide-number, notes

---

## closing（安可收尾）
指纹：hero

用途：收尾页。900 字重大字 + 一句行动提醒 + 琥珀底按钮 + 描边签，如海报角落的售票章。
适用 role：thanks / cta / content。
内容约束：标题 ≤10 字；lede 20-44 字；按钮 ≤6 字；签 ≤12 字。

```html
<section class="slide full" data-layout="closing">
  <p class="mp-kicker">{{行动语境，如 开幕预约}}</p>
  <h1 class="h1 mt-m">{{收束标题，≤10 字，可 <br> 分行}}</h1>
  <p class="lede mt-m" style="max-width:48ch">{{一句行动提醒，20-44 字}}</p>
  <div class="row mt-l" style="gap:24px">
    <span class="mp-btn">{{按钮文案，≤6 字}}</span>
    <span class="mp-tag">{{次级信息，≤12 字}}</span>
  </div>
  <div class="deck-footer"><span>{{署名 · 年份}}</span><span class="slide-number" data-current="{{总页数}}" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>
```

合法类名：slide, full, mp-kicker, h1, mt-m, lede, mt-l, row, mp-btn, mp-tag, deck-footer, slide-number, notes
