# dir-nav-minimal · 质量规则（生成阶段注入）

本文件是这份模板的专属质量规则，与 layouts.md 的版式锁配合使用。

## 一页一事

- 这是 keynote 式极简模板：一页只讲一件事，讲透就走。想堆料请换模板。
- 大标题（dk-h0 / dk-h1 / dk-h2）是主视觉，写结论不写栏目名——「答案不是更大的窗口」优于「第二章：窗口」。
- 正文只有三个去处：dk-lede（一句话）、dk-list（箭头清单）、dk-col（对照栏）。长段落放不下，也别硬塞。

## 底色与 accent

- 每页底色由骨架写死的 t-* 类决定（indigo→cream→crimson→emerald→slate→violet→white→charcoal），照抄，不要换、不要发明新色。
- accent 一律 `<span class="dk-accent">`，每页最多点 1-2 处（一个关键词或主数字）；再多就失去 mono 底色的冲击力。
- 深色页（indigo/crimson/emerald/slate/violet/charcoal）文字保持默认；浅色页（cream/white）同理——不要手动改页面文字颜色，对比度已经调好。

## 代码与数字

- dk-code 放命令/配置/代码：≤12 行、每行 ≤48 字符，`<` `>` `&` 必须转义；放伪代码不如不放。
- 关键数字一律走 stat 版式：dk-big 主数字 + 结论句 + SVG 占比条，不要把数字埋进 lede 的句子里。
- stat 的进度条：accent rect 宽度按占比给（0-900），标签写在条外右端；没有占比关系就不画条，只留大数字。

## 禁令

- 禁 emoji（✅❌✦ 等）：正反对比用「× / ✓」文字前缀，compare 版式已内置。
- 禁在 t-* 与 dk-accent 之外写死颜色：本模板的颜色全部来自骨架已写的底色类与 accent 钩子。
- 禁编造精确数字：stat 页的数字必须有口径（时间窗 / 样本 / 出处），写不出口径就换 list。
- 禁同一版式连排 ≥2 页：八种底色的轮换是节奏本身，同版式连排=同底色连排，观众会以为翻页失败。

## 密度

- 每页可见文字 70-160 字（cover/divider/thanks 这类 hero 页 15 字起；代码行另计）；低于下限会被服务端门禁直接拒收。
- dk-lede 25-50 字：一句完整的话，句号结尾；dk-list 每条 12-24 字，不写半截句。
- **每页要有下半区**：量测会检查填充率——dk-keyhint / dk-page 压住页脚，但正文块自身也要有分量，别只漂一行标题。
- **数字纪律**：关键数字/正反对比一律 stat（主数字）与 compare（×/✓ 两栏）呈现，不准埋进句子里。

## 节奏

- cover → divider → …内容… → cta / thanks 的骨架顺序保持完整；10 页以上的 deck 至少一页 divider。
- stat / code / compare 与 list 穿插排，不要连排；收尾二选一：有行动就 cta，纯收束就 thanks。
