# Agent 工具清单（deck-v2）

工具是 agent 的"手"：LLM 只输出调用意图，所有真实操作由后端 Go 代码执行。
本文档是 **deck-v2 管线** 工具的权威定义，与 `internal/agent/func_tool_v2.go`
的 `buildToolsV2` 逐一对应。v1 工具（write_deck / update_theme / custom_css 等）
已整体下线，历史档案与工程纪律的详细推演见 [tools-v1.md](tools-v1.md)。

## 核心设计：阶段 × 工具矩阵

v2 与 v1 的根本区别：**工具集按阶段裁剪**——outline 阶段看不见写页工具，
generating 阶段看不见大纲工具。阶段边界不靠提示词自觉，靠"工具根本不在场"
（`buildToolsV2(stage)` 是这张矩阵的唯一实现）。

管线阶段：`clarifying → outlining → outline_review →（用户确认）→ selecting_template →（用户选模板）→ generating → iterating`。
两个用户闸门（确认大纲 / 选模板）走 REST 不走对话；`selecting_template` 与
`generating` 阶段对话直接拒绝（`ErrStageLocked` → 409 + 可执行提示）。

| 工具 | clarifying | outlining | outline_review | generating | iterating | 一句话功能 |
|---|:-:|:-:|:-:|:-:|:-:|---|
| `ask_user` | ☑ | ☑ | ☑ | —（D12） | ☑ | 方向性提问，暂停循环等回答 |
| `web_search` | ☑ | ☑ | ☑ | ☑ | ☑ | 联网核实会变的事实（开关 `features.web_search`） |
| `submit_outline` | — | ☑ | — | — | — | 提交结构化大纲，创建 v2 deck |
| `read_outline` | — | — | ☑ | — | — | 读大纲全文与版本号 |
| `update_outline` | — | — | ☑ | — | — | 整份替换大纲（version 乐观锁） |
| `plan_pages` | — | — | — | ☑ | — | 全局版式分配（节奏规则在这里把关） |
| `read_layout` | — | — | — | ☑ | — | 取版式骨架与合法类名 |
| `read_guidelines` | — | — | — | ☑ | — | 取模板专属质量规则全文 |
| `write_pages` | — | — | — | ☑ | — | 分批写页（≤4 页/次），返回逐页结果+量测 |
| `review_slides` | — | — | — | ☑ | ☑ | 渲染量测（免费）+ 点名看图（3 次/run，开关 `features.vision`） |
| `list_slides` | — | — | — | ☑ | ☑ | 页面目录（页序/版式自查） |
| `read_slide` | — | — | — | ☑ | ☑ | 读单页 HTML + 指纹 |
| `update_slide` | — | — | — | ☑ | ☑ | 整页替换（指纹乐观锁） |
| `insert_slide` | — | — | — | — | ☑ | 插入新页（必须用已登记版式） |
| `delete_slide` | — | — | — | — | ☑ | 删页（至少留一页） |
| `list_history` | — | — | — | — | ☑ | 版本历史列表（只读） |
| `read_history_diff` | — | — | — | — | ☑ | 两份状态的页面差异（自查本轮改动） |

每阶段轮数预算（`defaultStageMaxTurns`，`deck_v2.max_turns` 可整体覆盖）：
clarifying 8 / outlining 10 / outline_review 8 / generating 24 / iterating 20。

## 通用纪律（从 v1 档案继承，仍然生效）

- **错误即反馈**：校验失败返回人类可读的报错，模型照着自我修正，不算失败；
- **校验永远先行**：归属（uid 来自认证上下文，绝不作工具参数）→ 存在性 → 结构 → 业务规则；
- **schema 描述禁半角逗号**（`invopop/jsonschema` 标签按半角逗号切分会静默截断描述，
  `TestSchemaDescriptionsHaveNoASCIIComma` 守着）；**嵌入提示词行尾归一化 LF**；
  **系统提示词的动态内容（deck_id/日期）一律追加在最后**（提示词前缀缓存）——
  三条的详细实测记录见 tools-v1.md 通用骨架一节；
- **消毒分级处置**（`service/deck/sanitize.go`，v2 写路径原样复用）：危险标签整页拒绝；
  `on*` 事件属性与危险 URL 协议剥除并在结果里附 `warning`——不让模型把"被剥除"汇报成"已生效"；
- **写路径持 per-deck 互斥锁**，锁覆盖"读-改-写"全程；锁内只调不带锁的原语。

## 管线事件（SSE，`stream_event.go`）

工具执行过程中向后端推送的事件（前端对未知类型一律忽略，保证向后兼容）：

| 事件 | 载荷 | 时机 |
|---|---|---|
| `stage` | `{deck_id, to}` | 阶段迁移（submit_outline、FinishGeneration 等） |
| `gate_waiting` | `{deck_id, gate}` | 到达用户闸门（outline/template） |
| `outline_updated` | `{deck_id, version}` | 大纲被 agent 修改 |
| `page_generated` | `{no, total, ok, error}` | write_pages 的每一页落地（前端进度条） |
| `lint_report` | `{page, items}` | AI 味 lint 报告 |
| `trace` | trace.Event JSON | 观测台（参数/耗时/分项用量，不含 messages） |

另有通用流事件：`session` / `delta` / `think` / `sub_delta` / `tool_start` /
`tool_delta` / `tool_call` / `tool_error` / `ask_user` / `done` / `error`。

---

## 1. ask_user —— human-in-the-loop（唯一"结果是一个活人"的工具）

- **参数**：`{ questions: [{question, options}] }`，1~6 问、每问 ≤4 选项，
  **第一个 = 推荐答案**（前端默认高亮、跳答时采用）
- **返回**：作答后合成 `{"answers":[{question,answer}]}`；跳答合成
  `{"note":"用户未作答，请使用推荐答案"}`（防 tool_call 无结果死锁）
- **机制**（`agent.go` + `persist.go`）：主循环分发前拦截——同消息其他工具合成
  "未执行"结果，SSE 推 `ask_user`，循环返回 `ErrPaused`；暂停态显式落库
  （`pending_ask` + messages 全量 blob），恢复走 `POST /api/chat/answer` 重入循环
- **D12：generating 阶段不挂载**——中途提问会打断生成管线，方向性问题必须在澄清阶段问完
- **触发纪律**（写在工具 description 反面清单）：打招呼/闲聊/能自己定的细节不要问；
  它会暂停整个循环，代价比一句话重得多

## 2. web_search —— 联网核实会变的事实

- **参数**：`{ query }` → **返回** `{answer, sources[], searches, truncated?, tool_error?}`
- 搜索由 DeepSeek 服务端执行（Anthropic 兼容入口）；**必须走 `/anthropic/v1/messages`**，
  Responses 端点会**静默忽略** web_search 工具（不报错，这是接入时踩过的最危险的坑，
  完整实测记录见 tools-v1.md §19）
- **计费口径**：`usage.server_tool_use.web_search_requests`；`max_uses` 在该端点不被强制执行，
  真实计费次数只看返回值 `searches`
- **引用纪律**（v2 修订）：模板没有 footnote 槽位——来源写成页面里的一行小字（用骨架允许的
  类名）或写进讲稿；搜不到就降级成不带数字的说法，**不要编造**
- 开关：`features.web_search`，关闭时连工具都不挂载（每次搜索都计费）

## 3. submit_outline —— 大纲首次落库（outlining）

- **参数**：`{title, audience?, duration_min?, tone?, hook?, arcs?, pages:[{no, role, title,
  points?, layout_hint?, materials?, notes?}]}`——schema 校验（role 枚举：cover/toc/divider/
  content/data/quote/code/cta/thanks；cover 开头、thanks 结尾），被拒就按报错修正重提
- **效果**：建 v2 草稿 deck（`CreateV2Draft`）+ 写 `outline.json`（version 1）；
  发 `stage → outline_review`、`gate_waiting → outline`、`outline_updated` 事件；
  run 结束后 agent 层把新 deck 绑回会话
- **返回**里的 `note` 要求模型向用户逐行展示页面结构——大纲是用户要确认的东西，不是黑箱

## 4. read_outline / update_outline —— 大纲修订（outline_review）

- `read_outline`：`{deck_id}` → 大纲全文 + version。**update 前必须先读**（整份替换制）
- `update_outline`：`{deck_id, version, ...同 submit}`——version 乐观锁（CAS），
  过期报冲突；用户在面板上的直改与对话修改共用同一份 version，天然互不覆盖
- 与用户面板直改（`PUT /api/decks/:id/outline`）走同一个 `SaveOutline`

## 5. plan_pages —— 全局版式分配（generating 第一步）

- **参数**：`{deck_id, assignments:[{no, layout, reason?}]}`——每页恰好一条、覆盖全部大纲页
- **服务端把关**（`SavePlanV2` → `ValidateRhythm`，`v2_rhythm.go`）：

| 规则 | 判定 | 处置 |
|---|---|---|
| R101 同版式连续 ≥3 | 滑动窗口 | **阻塞**，打回重排 |
| R102 每 8 页 <4 种版式 | 计数 | 提示 |
| R103 每 8 页无满版 hero 类 | 计数 | 提示 |
| R104 左右图文交替连续 >2 | 按名匹配 | 提示 |
| R106 同视觉模式连续 ≥3 | 按 `指纹：` 映射判观感 | **阻塞**（版式 id 不同≠观感不同，防假多样性；模板 ≥4 种模式才启用） |
| R107 每 8 页 <3 种视觉模式 | 计数 | 提示 |

- 违规聚合返回，阻塞规则打回重排、提示规则随结果带给模型（后续批次顺手优化）
- 版式的**视觉模式指纹**（hero/stack/cards/split/code/table/chart/quote）来自模板
  layouts.md 的 `指纹：` 行（`template/layouts.go` 解析，缺失时按骨架类名推断）

## 6. read_layout / read_guidelines —— 模板契约的只读视图

- `read_layout`：`{deck_id, layout}` → `{skeleton, classes[], note}`。骨架里的 `{{占位符}}`
  换成真实内容；**只准用列出的类名**；data-id 不要写（后端权威分配）。
  写某页之前必须先取它的骨架
- `read_guidelines`：`{deck_id}` → 模板 rules.md 全文（该模板的质量纪律）。开始写页前读一次

## 7. write_pages —— 分批写页（生成的主力）

- **参数**：`{deck_id, pages:[{no, layout, html}]}`——**每批 2-4 页**（上限 4，D14 护栏：
  不许一把梭，小批次才能每页都过闸门）；页可以乱序
- **每页过六道闸门**（`writeOnePage`，`v2_pages.go`）：

| # | 闸门 | 规则号 | 处置 |
|---|---|---|---|
| 1 | 页码范围 1~大纲页数 | R105 | 拒绝（实写页数必须等于大纲页数） |
| 2 | 版式已登记 + 与 plan_pages 分配一致 | C201 | 拒绝（换版式要全量重排，不许页面级偷换） |
| 3 | 结构+消毒：唯一根 section、禁嵌套、危险标签拒、`on*`/危险 URL 剥 | — | 拒/剥+warning |
| 4 | `data-layout` 属性与参数一致 | — | 拒绝 |
| 5 | 类名契约：只用该版式登记的类（base+template 两族并集） | C202 | 拒绝 |
| 6 | 静态密度下限：hero 15 / code 40 / quote 40 / 其余 70 可见字（CJK 逐字+西文按词，不含 `.notes`） | — | 拒绝（过空的页写进门就拦，不重量测） |

- **data-id 权威在后端**：剥掉 LLM 写的，按大纲页码编 `s<N>`（同页重写 = 替换）
- **AI 味 lint**（`LintTaste`，提示级不阻塞）：中英文 AI 高频词、标题超长、bullet 超长、
  提头超量等，结果随工具返回 + `lint_report` 事件，模型在下一批或修复轮自修
- 缺 `.notes` 讲稿块附 warning（演讲者模式空讲稿是质量缺陷，但有些版式确实没有，不阻塞）
- **成功批次自动跑量测**（免费、只回数字）；跑失败显式说明——"没跑成 ≠ 没问题"，
  不让 omitempty 把两者混成一样
- 每页发 `page_generated` 事件（前端进度条按大纲逐页点亮）

## 8. review_slides —— 渲染量测 + 点名看图

- **参数**：`{deck_id, pages?}`——pages 留空 = 只量测不产图（几秒、免费、随时可调）；
  点名页号才真正截图发给视觉模型（**一次 run 最多 3 次**，慢且贵）
- **量测**（`internal/vision`，chromedp 打开一次性渲染 URL）：每页 fit 比例 / 最小字号 /
  溢出 / 内容填充率（fill%，DOM 内容在画布上的垂直覆盖）/ 字数 / 版式；
  **硬拒绝线**：fill <45%（hero/quote 豁免）、字号 <13px 会在批次结果里直接点名
- **看图**：图片与数字一起给视觉模型，产出可执行的修复清单（点名制：v1 实测每次全量看图
  9 页要 2 分钟，点名的意义是把"看图"从默认动作变成可疑页的精确打击）
- 截图视口跟随 deck 画布（画布值夹取到 400~4000×300~3000，防手滑超大批图）
- 开关：`features.vision`

## 9. list_slides / read_slide / update_slide —— 页级读改（生成自检 + 迭代主力）

- `list_slides`：`{deck_id}` → `[{no, slide_id, layout, title}]`（token 预算设计，不含正文）
- `read_slide`：`{deck_id, slide_id}` → `{html, fingerprint}`（SHA256 前 12 位）
- `update_slide`：`{deck_id, slide_id, new_html, fingerprint}`——指纹**必填**（把"先读再改"
  做进 schema 硬约束）；过期拒绝且**不回显当前指纹**（防跳过 read 绕闸）；提交同样过
  六道闸门（data-id/data-layout/类名契约/消毒/密度）；成功返回刻意不含新指纹——
  同页再改必须重新 read，防"旧内容+新指纹"二次提交冲掉上次修改
- 生成阶段挂这三个工具的意义：写完一批立刻可以 list/read 自查，不等收尾

## 10. insert_slide / delete_slide —— 结构修改（iterating）

- `insert_slide`：`{deck_id, after_slide_id, new_html}`——`after_slide_id` 传某页 id 或
  `end`；new_html 必须带**已登记**的 data-layout，同样过六道闸门；id 后端分配
- `delete_slide`：`{deck_id, slide_id}`；拒绝删空整份 deck（至少留一页）

## 11. list_history / read_history_diff —— 版本历史（iterating，只读）

- 记档粒度：一轮 run 一条（含生成 run）；`list_history` 分页（默认 15/最大 50），
  `total/has_more/oldest_*` 让模型不翻页就能回答"一共多少版"
- `read_history_diff`：`{deck_id, from_version?, to_version?}`——都不传 = 本轮已做的修改
  （自查）；只传 from = 累计差异；都传 = 两版本之间。按 data-id 对齐 + go-diff 行级比较，
  单页 diff ≤80 行、最多 6 页给完整 diff（token 硬上限）
- **没有 restore 工具**：恢复由用户在界面操作（防模型误恢复），工具 description 引导用户去界面

---

## 附 A：用户模板定制工具（`usertpl/customize.go`，独立于 deck 管线）

用户与 agent 对话定制自己的模板（fork 出来的副本），是 ~6 轮上限的小循环
（不是 deck 管线），可改面刻意收窄到 **token 级**——结构契约（版式/类名）不可动，
从源头杜绝"对话把模板改坏"：

| 工具 | 参数 | 效果 |
|---|---|---|
| `write_tokens` | `{tokens: {名: 值}}` | 写 style.css 末尾的 `/* [customize] */ .tpl-<scope>{…} /* [/customize] */` 覆盖块（整块替换制）。键必须 `--` 开头、值禁 `;{}`（防注入）；写完重挂注册表，已发布模板立即生效 |
| `set_meta` | `{name?, description?}` | 改 template.json 与库表的名称/描述（画廊展示用） |
| `finish` | `{reply}` | 本轮结束，向用户复述改了什么、建议看哪页验证 |

可改 token：`--accent/--accent-2/--accent-3/--bg/--bg-soft/--surface/--surface-2/
--text-1/--text-2/--text-3/--radius/--radius-lg`。系统提示词带当前 token 值
（模型有上下文），并要求成套改动（改主色同步考虑次强调色与文字可读性）。

## 附 B：发布门禁（非工具，`usertpl.Publish` 的两道自动闸）

用户模板"公开给社区"前自动执行，不需要模型参与：

1. **结构+安全扫描**：`ValidateUserDir`（与内置模板同一套规则：layouts.md 指纹词表、
   骨架类名 ⊆ 合法类名、grid-trap 检查）+ 安全扫描（style.css 禁 `url(` / `@import` /
   `expression(` / `behavior:`；index.html 禁多余 `<script>`（runtime.js 除外）/ iframe /
   `on*` 事件）
2. **demo 渲染量测**：chromedp 渲染 demo，溢出 / 填充率 ≥45%（hero/quote 豁免）/
   最小字号 ≥13px 达标才可发布；报告落 `publish_report`，失败原因落 `publish_error`
