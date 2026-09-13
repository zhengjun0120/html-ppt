# Agent 工具清单

工具是 agent 的"手"：LLM 只输出调用意图，所有真实操作由后端 Go 代码执行。
本文档是工具的权威定义，实现时同步勾选，工具 description 可直接从这里改写。

## 总览

| # | 工具 | 阶段 | 一句话功能 | 状态 |
|---|---|---|---|---|
| 1 | `write_deck` | 1 | 从零生成整份 deck | ☑ |
| 2 | `list_decks` | 2+ | 列出**当前用户**的 deck（依赖会话持久化+用户系统，已推迟） | ☐ |
| 3 | `list_slides` | 2 | 看目录：每页的 id 和标题 | ☑ |
| 4 | `read_slide` | 2 | 读单页完整 HTML + 指纹 | ☑ |
| 5 | `update_slide` | 2 | 整块替换单页 | ☑ |
| 6 | `insert_slide` | 2 | 在某页后插入新页（后端分配 id） | ☑ |
| 7 | `delete_slide` | 2 | 删除单页 | ☑ |
| 8 | `update_theme` | 2 | 改主题配置（颜色/字体/动画） | ☑ |
| 9 | `list_templates` | 2 | 列出模板供推荐 | ☐ |
| 10 | `ask_user` | 2 | 向用户提问（human-in-the-loop，暂停循环） | ☑ |
| 11 | `screenshot_slides` | 3 | 截图/程序化检查排版 | ☐ |
| 12 | `move_slide` | 2+ 可选 | 调整页序 | ☐ |
| 13 | `read_history_diff` | 4 | 对比两份状态，看清某一轮改了什么（只读） | ☑ |
| 14 | `list_history` | 4 | 版本历史列表，取版本号（只读，分页） | ☑ |
| 15 | `read_custom_css` | 4 | 读 deck 级自定义样式槽（只读） | ☑ |
| 16 | `update_custom_css` | 4 | 整体替换自定义样式（清洗规则 + 可回滚） | ☑ |

实现顺序：1 → 3/4/5 → 6/7 → 8/9 → 10 →（阶段3）11 → 12 →（阶段4 版本控制）13/14 →（阶段4 样式槽）15/16。

> 15/16 由 `features.custom_css` 开关控制（见 config.go 的 Features）：**关掉时连工具都不挂载**，
> 而不是"看得到但一律拒绝"——后者只会让模型浪费轮次去试。

## 通用实现骨架

- 每个工具 = 名字 + 参数 JSON Schema（给 LLM 的菜单）+ handler（普通 Go 函数），统一注册表分发
- handler 统一签名：`(ctx, argsJSON) → (resultString, error)`，返回给 LLM 的永远是字符串
- **错误即反馈**：校验失败时返回人类可读的错误说明，LLM 会照着自我修正，不算失败
- 校验永远先行：id 白名单 → 结构合法性 → 业务规则
- 读工具尽量省 token（list/read 分开），写工具尽量严格（校验 + 原子写）
- **输入消毒（分级处置）**：危险标签（`script`/`style`/`iframe`/`frame`/`frameset`/`object`/`embed`/`applet`/`form`/`base`/`meta`/`link`）
  整页拒绝并报错；属性级违规（`on*` 事件属性、`javascript:`/`vbscript:`/`data:text/html` 协议）**剥除**
  并在工具结果里附 `warning` 字段告知模型。判定集中在 `service/deck/sanitize.go`，
  所有写路径（write_deck / update_slide / insert_slide）共用。
  为什么不静默剥除：模型以为 `onclick` 生效、实际被移除却向用户汇报"已加上交互"，是典型的静默失败。
  CSP 保护的是**预览**，消毒保护的是**导出的单文件**（导出后没有任何浏览器策略兜底），两者不可互替
- 防工具混乱三件套：description 写清"什么时候不用我"；system prompt 里给标准工作流；按场景分组渐进挂载（阶段2末）

---

## 1. write_deck —— 从零生成整份 deck

- **参数**：`{ title: string, sections_html: string }`
- **返回**：`{"deck_id": "deck-0002", "slides": 6}`
- **要点**：LLM 只提供可编辑区的 sections；骨架模板由后端拼装（三层权限的强制落实）；
  deck_id 后端生成（扫描现有最大编号 +1）；原子写入；系统提示词须完整描述组件库

## 2. list_decks —— 已推迟到阶段 2

- **参数**：无（用户身份来自登录态/注入的服务上下文，**绝不作为工具参数**）
- **返回**：当前用户的 `[{id, title, updated_at}]`
- **推迟原因**：a) 当前无会话记忆，跨轮列举无意义（列表结果出不了同一回合）；
  b) 选 deck 的职责由前端侧栏 + deck_id 注入承担；c) 正确形态需要归属过滤，
  单机版实现是错误形态，阶段 5 要返工
- **设计原则**：工具签名即权限边界——身份来自认证中间件注入，LLM 全程不知道
  "用户"概念的存在；越权访问一律 404（不泄露存在性）


## 3. list_slides

- **参数**：`{ deck_id }`
- **返回**：`[{position: 1, id: "s1", title: "封面"}, ...]`——只有页序 + id + 标题，
  不含正文（token 预算设计）。**position 是当前页序（随插删变化），id 是唯一标识
  （永不变化、删除后编号不复用）**——两者分离，防止模型对 id 数字做位置推断
- **要点**：goquery 遍历 AGENT-EDITABLE 区顶层 section，取 data-id + 第一个 heading 文本

## 4. read_slide

- **参数**：`{ deck_id, slide_id }`
- **返回**：`{ html: "...", fingerprint: "a3f8c92d" }`（SHA256 前 12 位）
- **要点**：按 `section[data-id]` 定位取 OuterHtml；为乐观锁准备指纹

## 5. update_slide —— 编辑场景最高频，所有写工具的校验模板

- **参数**：`{ deck_id, slide_id, new_html, fingerprint }`（fingerprint **必填**：把
  "先 read 再改"的工作流做进 schema 硬约束，而不是只靠提示词软约束）
- **校验顺序**：id 白名单 → slide 存在 → 指纹比对（过期拒绝；错误信息**不回显当前指纹**，
  否则 LLM 拿到指纹就能跳过 read 绕过闸门）→ new_html 校验（恰好一个根 `<section>`、
  data-id 与 slide_id 一致、禁 doctype/嵌套、消毒闸门：危险标签拒绝 + 属性级剥除）→ 替换写回
- **返回**：`{"slide_id": "s2"}`，若本轮剥除过违规属性则多一个 `warning` 字段（模型据此知道自己的写法没生效）
- **要点**：goquery 定位旧节点 → `ReplaceWithHtml` → 整份文档重序列化 → 原子写回；
  成功返回刻意不含新指纹（同一页再改必须重新 read，防止"旧内容+新指纹"的二次提交
  冲掉上一次修改）；写路径统一持 per-deck 互斥锁（锁覆盖"读-改-写"全程，防并发丢更新），
  版本快照不在这里做——记档统一在 run 结束后进行，见 §13/§14

## 6. insert_slide

- **参数**：`{ deck_id, after_slide_id, new_html }`（`after_slide_id: "end"` 表示追加末尾）
- **返回**：`{"slide_id": "s7"}`——**id 由后端分配（现有最大编号+1，永不复用）**
- **要点**：剥掉 LLM 写的 data-id 重新赋值；`AfterHtml` 插入

## 7. delete_slide

- **参数**：`{ deck_id, slide_id }`
- **要点**：定位 → Remove → 写回；拒绝删空整个 deck（至少留一页）

## 8. update_theme ☑

- **参数**：`{ deck_id, accent?, background?, heading_color?, text_color?, font?, radius?, transition? }`
  （只传要改的；`font` 枚举 sans/serif/mono，`transition` 枚举 slide/fade/zoom/convex/concave/none）
- **要点**：参数是语义化字段而非 CSS；真身是 deck.html 里的
  `<script type="application/json" id="deck-theme">` 配置块，读 JSON → 合并 → 写回 → 同步 CSS 变量；
  校验颜色格式正则 + 枚举；派生变量（border/card-bg/text-muted）集中计算，
  **旋钮越少模型的选择面越小、观感越不容易崩**

## 9. list_templates

- **参数**：无
- **返回**：`[{id, name, description, style_tags}]`
- **要点**：description 写清适用场景，它是 LLM 推荐模板的依据；阶段2读静态 JSON 桩，阶段4动态化

## 10. ask_user —— human-in-the-loop，唯一"结果是一个活人"的工具 ☑

- **参数**：`{ questions: [{question, options}] }`——questions 1~6 个（三处一致：
  schema description / 工具 description / 后端硬校验）；options ≤4 个，
  **第一个 = 推荐答案**（前端默认高亮，跳答时提示采用）
- **返回**：用户作答后由后端合成 `{"answers": [{question, answer}]}`；用户跳答时合成
  `{"note": "用户未作答，请使用推荐答案"}`（防止 tool_call 无结果导致的死锁）
- **要点**（实现于 agent/agent.go + persist.go + handler/answer.go）：
  - 主循环在分发前拦截 ask_user：**不执行、不继续循环**——同消息里的其他工具
    调用合成"未执行"结果（协议要求每个 tool_call 必须有配对结果，否则恢复后
    回放直接 400），SSE 推 `ask_user` 事件，循环返回 ErrPaused 哨兵（控制信号，
    类比 io.EOF，不是失败）
  - 暂停态显式落库：`chat_sessions.pending_ask` 存 `{"tool_call_id"}`，
    messages 全量 JSON blob 存 `chat_sessions.messages`（每个检查点整体重写）
  - 恢复端点 `POST /api/chat/answer`：读回 messages → 回填 ask_user 的 tool 结果 →
    清 pending_ask → 重入循环（可能连环追问再次暂停）
  - 参数不合法（json 坏 / 数量越界）→ 不暂停，作为普通工具错误反馈给模型，
    循环继续（错误即反馈）
  - 会话归属：session_id 越权访问一律 404（与 deck 同款）
- **复用价值**：任何需要用户确认的时刻（如自定义脚本写入前的二次确认）都用它

## 11. screenshot_slides

- **参数**：`{ deck_id, slide_ids: [...] }`
- **返回**：截图路径列表（供视觉模型审图）或程序化排版检查报告
- **要点**：playwright-go 打开 deck URL → `#/2` hash 定位页 → 截图存
  `data/decks/<id>/shots/`；**deepseek-chat 是纯文本模型看不了图**，
  建议程序化检测为主（JS 检查元素溢出/截断，返回结构化报告）+ 可选视觉模型为辅

## 12. move_slide（可选，优先级最低）

- **参数**：`{ deck_id, slide_id, after_slide_id }`
- **要点**：insert + delete 的组合操作

---

## 13. read_history_diff —— 对比两份状态看清改了什么（只读，版本控制）☑

- **参数**：`{ deck_id, from_version?, to_version? }`
  - `from_version` 不传 = 最新记档版本；`to_version` 不传 = 当前使用中的内容
  - **`"current"` 不是可传入的取值**：空值即"当前"，`"current"` 只作为**输出标签**出现
    （输入约定与输出标签解耦，避免两边分支打架）
- **返回**：`{ deck_id, from, to, to_detail, changed, added[], removed[], modified[], unchanged }`
  - `added` / `removed` 每项 `{slide_id, title, change}`；`modified` 每项额外带 `diff` + `truncated`
  - `changed=false` = 两份完全一致（新一轮开始时、以及 run 记档后都会是 false，都是正常状态）
- **语义（同一次比较，区别只在基线选谁）**：
  - 都不传 = 最新记档版本 vs 当前内容 = **本轮已做的修改**。因为记档发生在 run 结束，
    最新版本就是上一轮的终点，两者之差恰好是本轮在途改动 → 模型改完几页后**自查**用它
  - 只传 `from_version` = 从该版本到现在的**累计差异**。快照是全量的，不需要沿版本链累加，
    直接比两份就是累计结果（类比 `git diff <ref>..工作区`）
  - 两个都传 = 两个历史版本之间（"上一轮改了啥" = 上一条 vs 最新一条，版本号先用 §14 查）
- **要点**：
  - 按 `data-id` 对齐两份 deck，两侧都用 goquery 重新解析 + 重新序列化后再比较
    （归一化，避免格式差异误报成 modified）
  - diff 用 `sergi/go-diff` 的行级模式（`DiffLinesToChars → DiffMain → DiffCharsToLines`）
    + 自渲染 unified 格式；选它的原因之一是容错匹配：LLM 生成的 HTML 常有少量错位，
    比纯 LCS 更不容易把"挪了一行"渲染成大片红绿
  - **token 上限（硬约束）**：单页 diff ≤ 80 行（超出截断并在末尾补 `...`）、最多 6 页给完整
    diff，其余修改页只报"改了"；`@@` 头的行数按**截断后**的可见内容统计，避免头与正文数字矛盾
  - 只读工具，不参与 run 记档（`runRecorder` 只记会改 deck 的五个写工具）

## 14. list_history —— 版本历史列表（只读，版本控制）☑

- **参数**：`{ deck_id, limit?, offset? }`——`limit` 默认 15、最大 50，`offset` 用于翻页看更早的
- **返回**：`{ deck_id, total, returned, has_more, oldest_version, oldest_time_str, versions[] }`
  - 每条：`{version, time, operation, detail, slides, time_str}`
- **要点**：
  - **默认截断是硬要求，不是优化**：工具结果常驻对话上下文、之后每一轮请求都要重发，
    200 条全量约 1 万 token 且被反复计费。`total` / `has_more` / `oldest_*` 让模型不翻页
    就能回答"一共有多少版""最早能回到哪"
  - `time_str` 是本地时间格式化（模型判断"多久以前"比读 unix 秒直观）；存储层仍是 unix 秒
  - `version`（如 `v000003`）就是 §13 的 `from_version` / `to_version` 取值——两个工具闭环
  - 只读；**没有 restore 工具**：恢复由用户在界面上操作（防模型误恢复），
    工具 description 里明确告知模型"引导用户去界面操作"

---

## 15. read_custom_css / 16. update_custom_css —— deck 级自定义样式槽 ☑

**槽是什么**：deck.html 的 `<head>` 里一个 `<style id="deck-custom">`，级联在组件库与主题
override 之后。它是"框架层给 AI 开的一个洞"——页面内容区仍然禁 `<style>`，这个洞单独收口。

**为什么放在 deck 里而不是单独文件**：
- 整份 deck 的历史快照覆盖它 → 写坏了能一键回滚（这是放开这份自由度的前提）；
- 导出成单文件时天然在内，不需要额外处理；
- `update_theme` 的真身（JSON 块 + override 块）本来就在这儿，字体槽将来也放这里，
  一个 deck 的视觉配置只在一个地方。

**15. read_custom_css**
- **参数**：`{ deck_id }`
- **返回**：`{ deck_id, css, bytes }`——老 deck 还没槽时返回空串

**16. update_custom_css**
- **参数**：`{ deck_id, css }`——`css` 是**该 deck 的全部自定义样式**（整体替换，不是追加）；
  空字符串 = 清空。所以工作流是：`read_custom_css` → 改 → 提交全文
- **返回**：`{ deck_id, css, bytes, cleared }`
- **校验（清洗规则）**——挡的是"通道"而不是"风格"：
  - `</style` 拒绝（会提前终止 raw-text 样式块，破坏整份文档）
  - `@import` 拒绝（外部样式表 = 外部依赖 + 数据外发通道）
  - `url()` 拒绝（会发起外部请求；图片走 `<img>`，字体内嵌留给字体功能）
  - `expression(` 拒绝（可执行代码的老写法）
  - 32KB 上限（同时是 token 护栏）
  - 风格偏离**不校验**：由用户看着预览决定，不满意就回滚历史
- **刻意不做指纹校验**：槽只有一个写入者（agent），且历史快照兜底，冲突代价低——同 `update_theme`
- **级联不变量**：`#deck-theme-override` 必须在 `#deck-custom` 之前。两条补块路径都要守住
  （先建槽再建主题块 / 反过来），否则自定义样式会被主题派生变量盖掉——这种 bug 在页面上
  表现为"我写了 CSS 但没生效"，最难排查。回归测试在 `custom_css_test.go` 的
  `TestCustomCSSBlockPlacement`
- **老 deck 迁移**：骨架只在 `write_deck` 时固化一次，所以老 deck 没有这个块——
  首次写入时按需补出来（`ensureCustomCSSBlock`），和 `ensureThemeBlocks` 同一套路，
  不写独立迁移脚本

---

## 附：版本控制 REST 接口（非 agent 工具）

工具层刻意不暴露恢复能力，恢复与历史管理走 HTTP，由界面调用：

| 方法 | 路径 | 作用 |
| --- | --- | --- |
| GET | `/api/decks/:id/history` | 版本列表（新→旧） |
| POST | `/api/decks/:id/history/:version/restore` | 恢复到某版本（**恢复本身也记一条版本，可再撤销**） |
| DELETE | `/api/decks/:id/history/:version` | 删除某条历史（快照互相独立，删中间不影响其他） |
| DELETE | `/api/decks/:id/history` | 清空全部历史（`NextSeq` 不清零，编号永不复用） |

**存储**：`data/decks/<id>/history/<version>.html`（整份 deck.html 快照）+ `index.json`
（版本元信息 + `NextSeq`）。快照含主题块，所以**恢复会连主题一起回滚**（设计使然）。

**版本粒度**：一轮 agent run 一条（一轮用户消息 = 一条版本），`Detail` 由 agent 层
`runRecorder` 汇总本轮工具调用生成（"修改了 s2、调整主题"）。工具级的中间态不入历史，
撤销一轮 = 恢复上一条版本。

**记档时机**：run 成功结束或 `ask_user` 暂停时（失败的 run 不记档，历史只留"成功产出过的状态"）；
另加每次 restore 操作即时记档。写路径本身不含记档逻辑。

**保留策略**：最近 7 天全留，更早的每天留最后一条，总数上限 200；由记档时顺手裁剪，无定时器。

**并发**：`Service.deckLocks` 每份 deck 一把 `sync.Mutex`，锁覆盖"读-改-写"全程；
**锁内只调用不带锁的原语（`readRaw` / `atomicWriteDeck`），绝不在已持锁路径里调用会自己加锁的函数**
（`sync.Mutex` 不可重入，重入即死锁）。
