# 小红书图文 · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 `<section>` 必须带 `data-layout="<本文件登记的 id>"`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 `read_layout` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 `.tpl-xhs-post` 作用域前缀生效，骨架里已写全，照抄结构即可。
> 画布是竖版 810×1080（小红书 3:4 图文）：字号大、字要少，每页只讲一件事。
>
> **结构铁律**：每页必须以 `<div class="page-dot">{{当前页}} / {{总页数}}</div>` 开头——
> 页码贴纸是九宫格的身份；cover 与 cta 两页必须以 `<div class="bottom-bar">`
> （头像 + 账号名）收尾。sticker 是绝对定位贴纸，每页最多 2 枚，
> 位置/角度照抄骨架给出的 style。{{占位符}} 只准换文本；骨架层级不准增删。

---

## cover（贴纸封面）
指纹：hero

用途：首图。两枚歪贴纸造氛围 + 栏目行 + 大标题（荧光笔高亮关键词）+ 底部账号条。
适用 role：cover。
内容约束：主标题 2-3 行（`<br>` 分行）、每行 ≤9 字，1-2 个关键词包 `<span class="cover-title">`；lede 栏目行 12-20 字；2 枚贴纸各 ≤6 字；贴纸色从 pink/yellow/blue/green 里选。

```html
<section class="slide" data-layout="cover">
  <div class="page-dot">{{当前页}} / {{总页数}}</div>
  <div class="sticker pink" style="top:120px;left:48px;transform:rotate(-6deg)">{{贴纸 1，≤6 字}}</div>
  <div class="sticker yellow" style="top:140px;right:64px;transform:rotate(5deg)">{{贴纸 2，≤6 字}}</div>
  <div style="margin-top:200px">
    <p class="lede" style="font-size:24px;color:var(--text-1);font-weight:600">{{栏目一句话，12-20 字}}</p>
    <h1 class="h1 mt-s">{{主标题，2-3 行 <br> 分行，关键词包 <span class="cover-title">}}</h1>
  </div>
  <div class="bottom-bar"><div><span class="avatar">{{账号首字}}</span> <b style="color:var(--text-1);margin-left:8px">{{账号名}}</b></div><div>← 左滑 查看</div></div>
  <div class="notes">{{配文}}</div>
</section>
```

合法类名：slide, page-dot, sticker, pink, yellow, blue, green, lede, h1, mt-s, cover-title, bottom-bar, avatar, notes

---

## hook（钩子页）
指纹：hero

用途：第 2 图。大 emoji 定格 + 反问钩子 + 两三行共鸣，让读者停下手指。
适用 role：divider / content。
内容约束：1 个 emoji；反问 ≤12 字；共鸣 2-3 行（`<br>` 分行、每行 ≤14 字），关键数字包 `<b style="color:var(--accent)">`；1 枚贴纸 ≤8 字（转折提示）。

```html
<section class="slide" data-layout="hook">
  <div class="page-dot">{{当前页}} / {{总页数}}</div>
  <div class="big-emoji" style="margin-top:80px">{{一个 emoji}}</div>
  <h2 class="h2 tc mt-l">{{反问钩子，≤12 字}}</h2>
  <p class="lede tc mt-m" style="padding:0 20px">{{共鸣 2-3 行 <br> 分行，每行 ≤14 字；关键数字包 <b>}}</p>
  <div class="sticker blue" style="bottom:160px;left:50%;transform:translateX(-50%) rotate(-2deg)">{{转折提示，≤8 字}}</div>
  <div class="notes">{{配文}}</div>
</section>
```

合法类名：slide, page-dot, big-emoji, h2, tc, mt-l, lede, mt-m, sticker, blue, pink, yellow, green, notes

---

## pain（痛点清单）
指纹：stack

用途：痛点共鸣页。❌ 引子 + 短标题 + 3-4 张便签卡，每张一个症状。
适用 role：content。
内容约束：3-4 张 hand-box；每条标题 ≤14 字（emoji 开头）+ 说明 12-22 字（句号结尾）；引子 ≤10 字（❌ 开头）；标题 ≤8 字。

```html
<section class="slide" data-layout="pain">
  <div class="page-dot">{{当前页}} / {{总页数}}</div>
  <p class="lede" style="font-weight:700;color:var(--accent)">{{❌ 引子，≤10 字}}</p>
  <h2 class="h2 mt-s">{{痛点标题，≤8 字}}</h2>
  <div class="stack mt-l">
    <div class="hand-box"><b style="font-size:22px">{{emoji + 症状 1，≤14 字}}</b><p class="dim" style="font-size:18px;margin-top:4px">{{一句扎心说明，12-22 字}}</p></div>
    <div class="hand-box"><b style="font-size:22px">{{emoji + 症状 2，≤14 字}}</b><p class="dim" style="font-size:18px;margin-top:4px">{{说明，12-22 字}}</p></div>
    <div class="hand-box"><b style="font-size:22px">{{emoji + 症状 3，≤14 字}}</b><p class="dim" style="font-size:18px;margin-top:4px">{{说明，12-22 字}}</p></div>
  </div>
  <div class="notes">{{配文}}</div>
</section>
```

合法类名：slide, page-dot, lede, h2, mt-s, stack, mt-l, hand-box, dim, notes

---

## truth（真相断言）
指纹：quote

用途：aha moment 页。一句荧光笔高亮的大字断言 + 两三行原理 + 一句关键方法。整页只交付这一个「原来如此」。
适用 role：quote / content。
内容约束：断言两行（`<br>` 分行）、每行 ≤10 字，高亮 1 处（`<span style="background:var(--accent-3);padding:0 8px">`）；原理解释 2-3 行每行 ≤16 字；关键句 ≤8 字（`<span style="color:var(--accent)">`）；引子 ≤8 字（💡 开头）。

```html
<section class="slide" data-layout="truth">
  <div class="page-dot">{{当前页}} / {{总页数}}</div>
  <div class="sticker green" style="top:100px;right:48px;transform:rotate(4deg)">{{✨ 标签，≤10 字}}</div>
  <p class="lede mt-l" style="color:var(--accent);font-weight:700">{{💡 引子，≤8 字}}</p>
  <h2 class="h2 mt-s">{{断言，两行 <br> 分行，关键词包 <span style="background:var(--accent-3);padding:0 8px">}}</h2>
  <p class="lede mt-l">{{原理解释，2-3 行 <br> 分行，每行 ≤16 字}}</p>
  <p class="lede mt-m" style="color:var(--text-1);font-weight:700">关键是：<span style="color:var(--accent)">{{一句话方法，≤8 字}}</span>。</p>
  <div class="notes">{{配文}}</div>
</section>
```

合法类名：slide, page-dot, sticker, green, pink, yellow, blue, lede, mt-l, h2, mt-s, mt-m, notes

---

## step（干货步骤）
指纹：stack

用途：核心干货页（一篇 9 图里可重复 2-3 页）：num-circle 步骤号 + 步骤名 + 公式/做法便签 + 浅底举例便签。
适用 role：content。
内容约束：步骤号 1-9；步骤名 ≤10 字；主 hand-box 要点 ≤14 字 + 做法 20-45 字；第二个 hand-box（`background:var(--surface-2)`）写举例，例句 ≤20 字；可选 1 枚贴纸。

```html
<section class="slide" data-layout="step">
  <div class="page-dot">{{当前页}} / {{总页数}}</div>
  <div class="num-circle">{{步骤号，1-9}}</div>
  <h2 class="h2 mt-m">{{步骤名，≤10 字}}</h2>
  <div class="hand-box mt-l">
    <p style="font-size:22px;margin:0;color:var(--text-1);font-weight:700">{{👉 要点标题，≤14 字}}</p>
    <p style="font-size:20px;margin:10px 0 0;color:var(--text-2);line-height:1.7">{{公式或做法，20-45 字，关键数字包 <b style="color:var(--accent)">}}</p>
  </div>
  <div class="hand-box mt-m" style="background:var(--surface-2)">
    <p style="font-size:18px;margin:0;color:var(--text-2)">举例：</p>
    <p style="font-size:24px;margin:8px 0 0;color:var(--text-1);font-weight:800">{{具体例子，≤20 字，可两行 <br>}}</p>
  </div>
  <div class="notes">{{配文}}</div>
</section>
```

合法类名：slide, page-dot, num-circle, h2, mt-m, hand-box, mt-l, notes

---

## result（前后对比）
数量：hand-box=4
指纹：cards

用途：效果页。✅ 引子 + 2×2 便签格，每格一个变化 + 一句对比，末格可换底色强调。
适用 role：content。
内容约束：恰好 4 格（grid g2）；每格标题 ≤12 字（emoji 开头）+ 对比 12-24 字；引子 ≤12 字（✅ 开头）；第 4 格可加 `style="background:var(--accent-3);border-color:var(--text-1)"` 强调。

```html
<section class="slide" data-layout="result">
  <div class="page-dot">{{当前页}} / {{总页数}}</div>
  <p class="lede" style="color:var(--good);font-weight:700">{{✅ 引子，≤12 字}}</p>
  <h2 class="h2 mt-s">{{效果标题，≤8 字}}</h2>
  <div class="grid g2 mt-l">
    <div class="hand-box"><b style="font-size:20px">{{emoji + 变化 1，≤12 字}}</b><p class="dim" style="font-size:18px;margin:6px 0 0">{{一句对比，12-24 字}}</p></div>
    <div class="hand-box"><b style="font-size:20px">{{emoji + 变化 2，≤12 字}}</b><p class="dim" style="font-size:18px;margin:6px 0 0">{{对比，12-24 字}}</p></div>
    <div class="hand-box"><b style="font-size:20px">{{emoji + 变化 3，≤12 字}}</b><p class="dim" style="font-size:18px;margin:6px 0 0">{{对比，12-24 字}}</p></div>
    <div class="hand-box" style="background:var(--accent-3);border-color:var(--text-1)"><b style="font-size:20px">{{emoji + 最重要的变化，≤12 字}}</b><p class="dim" style="font-size:18px;margin:6px 0 0">{{对比，12-24 字}}</p></div>
  </div>
  <div class="notes">{{配文}}</div>
</section>
```

合法类名：slide, page-dot, lede, h2, mt-s, grid, g2, mt-l, hand-box, dim, notes

---

## cta（关注收尾）
指纹：hero

用途：末图。大 emoji + 收藏关注号召 + 下期预告 + 话题 tag 行 + 账号条。
适用 role：cta / thanks。
内容约束：1 个 emoji；引导 ≤8 字；号召 ≤8 字（可带 1 个 emoji）；预告两行每行 ≤14 字；3-4 个 `<span class="ht">` 话题（# 开头，各 ≤6 字）。

```html
<section class="slide" data-layout="cta">
  <div class="page-dot">{{当前页}} / {{总页数}}</div>
  <div class="big-emoji" style="margin-top:60px">{{一个 emoji}}</div>
  <h2 class="h2 tc mt-l">{{引导，≤8 字}}</h2>
  <h1 class="h1 tc mt-s" style="color:var(--accent)">{{行动号召，≤8 字}}</h1>
  <p class="lede tc mt-l" style="padding:0 30px">{{下期预告，两行 <br> 分行，每行 ≤14 字}}</p>
  <div class="tag-row" style="justify-content:center;margin-top:36px">
    <span class="ht">#{{话题 1}}</span>
    <span class="ht">#{{话题 2}}</span>
    <span class="ht">#{{话题 3}}</span>
    <span class="ht">#{{话题 4}}</span>
  </div>
  <div class="bottom-bar"><div><span class="avatar">{{账号首字}}</span> <b style="color:var(--text-1);margin-left:8px">{{账号名}}</b></div><div>{{❤️ 点赞数}}</div></div>
  <div class="notes">{{配文}}</div>
</section>
```

合法类名：slide, page-dot, big-emoji, h2, tc, mt-l, h1, mt-s, lede, tag-row, ht, bottom-bar, avatar, notes
