# PPT 生成流程重构方案（deck-v2）· 执行蓝图

> 状态：全部决策已对齐（D1–D14），方案冻结，待最终确认后按 P0 开工。
> 本文档是**执行蓝图**：细化到真实文件路径、接口签名、数据结构、事件载荷和验收标准，供后续会话直接照此施工，不依赖对话记忆。
> 事实基准：2026-09-17 的仓库结构（backend/internal 各包、frontend/src 各文件已逐一核对）。

---

## 1. 背景与问题

现状是"一句话进来 → agent 提问确认 → `write_deck` 一次性生成整份 deck"，渲染底座是 reveal.js。质量天花板受制于三点：

1. **渲染栈过长**：moon.css 兜底皮肤 + theme.css 变量层 + components.css 组件层 + 自定义 CSS 槽，四层叠加，LLM 在顶端自由发挥 HTML，观感不可预测。历史证据："优化排版反而反向优化"（be4d753）、22 轮不收敛的审查循环、init.js 两级 fit 兜底——全是补丁。
2. **质量靠 prompt 叠规则**：systemPrompt.md 667 行，规则间冲突风险与 token 成本同步增长。
3. **没有大纲环节**：大纲只存在于提示词纪律里，用户对内容的控制在生成之后，返工贵。

外部参照（html-ppt-skill / guizang-ppt-skill / taste-skill，均已调研）证明：**高质量不靠渲染引擎，靠"模板自带完整设计系统，LLM 只准填内容"**。本次重构把质量责任从 prompt 转移到模板。

## 2. 目标与非目标

**目标**

- 五段式生成流程：澄清 → 结构化大纲（双通道编辑）→ 选模板（必选）→ 分批生成 + 自动质检 → 对话迭代。
- 渲染底座换为模板自带运行时（html-ppt-skill runtime.js + base.css），固定画布 1920×1080 transform 缩放，量测/截图所见即所得。
- 模板为目录级自包含单元，契约兼容未来用户自定义模板。
- 三层程序化质量网：AI 味 lint + 版式节奏校验 + 类名契约校验，叠加现有 vision 审查。
- 导出 PDF / 逐页 PNG / 单文件 HTML。

**非目标（v1 不做）**

- .pptx 导出（独立管线，另立项）。
- 用户自定义模板的完整 UI（只保证契约与加载器就绪）。
- 非 16:9 画布（runtime 支持 data-w/h，v2 再放开 xhs-post 等竖版场景）。
- fragments 分步显示、多人协作、生成托管运行（断线续跑做简化版，见 §7.4）。

## 3. 决策记录（全部已对齐）

| # | 决策 | 结论 |
|---|------|------|
| D1 | reveal.js 去留 | **换掉**。模板自带运行时；init.js、moon.css、theme 渲染层、14 preset 主题系统全部退役 |
| D2 | 模板主干 | **html-ppt-skill（MIT）**：runtime.js + base.css + 整 deck 模板目录规范 |
| D3 | guizang（AGPL-3.0） | **只吸收思想、重写表达，不复制任何文件与代码** |
| D4 | 导出形态 | **HTML 主产物 + 服务端 headless Chrome 出 PDF/逐页 PNG** |
| D5 | 大纲对齐 | agent 先提问拿信息 → 结构化大纲 → 用户**面板直改**或**对话让 agent 改**，双通道同一份 JSON，确认过关 |
| D6 | 首批模板 | **精选 8 个**，逐个补齐版式清单与审查标准后上架 |
| D7 | 存量兼容 | **直接废弃**，v1 deck 隐藏不迁移（§6.4） |
| D8 | 版式与主题 | 视觉身份归模板所有，模板内 2-4 个**主题变体槽**；不做跨模板主题组合 |
| D9 | 首批模板名单 | §5.8 的 8 个确认 |
| D10 | AI 味禁词表 v1 | §9.4 词表与"全部提示级"确认 |
| D11 | variants 变体槽 | **v1 就做**，随 P2 画廊上线 |
| D12 | 生成中交互 | **不允许**中途插入对话指令，generating 阶段锁输入 |
| D13 | deck 存储形态（原 R1） | **混合**：runtime/base.css 共享绝对路径引用（版本化逃生舱），模板 style.css 拷贝进 deck 目录（快照隔离），下载时打包单文件 |
| D14 | 生成粒度（原 R2） | **plan_pages 一次规划 + 每批默认 3 页（可配 2-4）分批写入**；批内按页部分成功；同 run 内对话记忆保证跨页一致性 |

## 4. 与现状差异速览

| 维度 | deck-v1（现状） | deck-v2（目标） |
|------|----------------|----------------|
| 渲染底座 | reveal.js 5.1 + init.js fit 兜底 | runtime.js 固定画布 transform 缩放 |
| 视觉来源 | 14 preset 主题 + LLM 自由发挥 | 模板设计系统，LLM 只准用登记版式 |
| 大纲 | 提示词纪律，无数据结构 | 一等公民 outline.json，双通道编辑 |
| 生成方式 | write_deck 一次全量 | plan_pages 规划 + write_pages 每批 3 页 |
| 修改协议 | fingerprint 乐观锁整页替换 | 保留 + 类名契约校验 |
| 质量网 | stylelint + 变量契约 + vision | + AI 味 lint + 节奏校验 + 类名契约 |
| 产物 | 单文件 deck.html（reveal 格式） | deck 目录三件套，导出可打包单文件 |
| systemPrompt | 667 行单文件 | 分阶段五件套，设计规则下沉到模板 |

## 5. 模板体系

### 5.1 目录契约

```
backend/templates/<id>/
  template.json    # 机器可读元数据（附录 A）
  index.html       # 自包含骨架：head/CSS 引用 + demo 数据 deck + <!-- SLIDES:START -->…<!-- SLIDES:END --> 挂载区
  style.css        # .tpl-<id> 作用域隔离的完整设计系统（token 覆盖 + 专属组件 + variants class）
  layouts.md       # 版式登记簿：每个版式一段（用途/适用 role/骨架代码/合法类名清单/内容约束），LLM 的唯一契约
  rules.md         # 本模板专属质量规则（字号分档、图片比例、禁忌），生成阶段注入
  preview/
    cover.png      # 画廊封面（16:9，≥800×450）
    page-*.png     # 关键页截图（可选）
    golden-*.png   # 中文黄金样张（上架标准，见 §5.8）
```

template.json 关键字段：`id / name / description / tags[] / scenario[]（场景关键词，画廊推荐用）/ canvas{w,h} / variants[]（{id,name,class}，class 为 style.css 里预定义的 token 覆盖 class，default 变体 class 为空串）/ layouts[]（{id,name,use,constraints?}）/ fonts[]（自托管字体清单）/ source{derived_from,license}`。

**variants 实现**：每个变体 = style.css 中一个 class（如 `.tpl-tech-sharing.v-forest`），只覆盖 `--accent/--grad/--bg` 等少量 token，其余继承默认。用户选变体 = 实例化时 body 挂 class，纯存储操作，agent 不参与。对比度由 base.css 的 `--accent-ink` 机制兜底（html-ppt-skill 全主题 WCAG 实测过）。

### 5.2 layouts.md 规范（版式锁，guizang 思想重写）

- 每个版式条目固定五段：**用途 / 适用 role / 完整 `<section>` 骨架（占位符文本）/ 合法类名清单 / 内容约束**。格式见附录 C。
- **硬规则：每个 `<section>` 必须带 `data-layout="<id>"`，id 必须在 template.json layouts[] 登记**，服务端回声校验，未登记拒收（§9.2）。
- 版式词汇表来源：以整 deck 模板自带 7-10 页结构为骨架，从 html-ppt-skill 36 个单页版式中挑类兼容的并入，必要时手写。**layouts.md 是每个模板版式的唯一事实来源**，单页版式库不直接暴露给 LLM。
- 首批每模板登记 12-18 个版式，覆盖：cover / toc / divider / bullets / two-column / three-column / stat-hero / kpi-grid / chart-bar / chart-line / chart-pie / big-quote / timeline / process / comparison / image-text / code / cta / thanks 语义族。

### 5.3 运行时与共享资产

```
backend/web/assets/deck-v2/
  runtime.js       # vendor 自 html-ppt-skill（MIT），翻页/演讲者/概览/画布缩放/打印/深链
  base.css         # vendor，设计 token + 布局原语（.slide/.grid/.card/.h1/.lede/.img-frame…）
  fonts.css        # @font-face 指向 fonts/（替代原 Google Fonts @import）
  fonts/           # 自托管 woff2：Noto Sans SC（400/500/700）、JetBrains Mono（400/700）等
  PATCHES.md       # 对 vendor 文件的每一处改动记录（原因 + 原始行为）
```

- deck 的 index.html 以绝对路径引用 `/assets/deck-v2/runtime.js` 等（升级集中生效）；**破坏性升级时发 `/assets/deck-v2.1/` 新路径，老 deck 继续引用老路径，零迁移**（D13 逃生舱）。
- 字体：从仓库现有 `font/` 目录（已 ignore，20MB 副本）取源做 woff2，P0 落进 `deck-v2/fonts/`（预计总量 3-6MB，可接受）；base.css 字体栈同步改本地路径，保留系统字体回退。**解决国内访问 Google Fonts 失败与离线渲染两个问题**。
- runtime 深链能力直接复用：`#/N` 翻页定位、`?preview=N` 单页模式（量测截图用）、`@media print`（PDF 导出用）。**P0 必须做三个 spike 验证**：① headless 环境下 `#/N` 与 `?preview=N` 哪个定位单页更稳；② PrintToPDF 与其 print CSS 的分页配合；③ iframe sandbox 内初始化时序。结论与改动记入 PATCHES.md。

### 5.4 模板入库适配（一次性，入仓库前完成）

html-ppt-skill 原模板用相对路径（`../../../assets/...`）引用共享资产。**入库时预先适配**，deck 实例化时不再做路径改写：

1. 复制模板目录到 `backend/templates/<id>/`；
2. 全部共享资产引用改写为 `/assets/deck-v2/...` 绝对路径；
3. index.html 加 `<!-- SLIDES:START/END -->` 挂载标记、body 加 `data-variant` 占位属性；
4. demo deck 保留（画廊 live 预览直接用它）；
5. 跑入库校验脚本（`backend/templates/tools/validate-template.mjs`，Node）：文件齐全、引用全部可解析、layouts.md 可解析且 id 与 template.json 一致、variants class 在 style.css 中存在。
6. 适配过程记录到 `backend/templates/<id>/ADAPTATION.md`（改了什么、为什么），保证与上游 diff 可追溯。

### 5.5 模板注册与加载（新包 `backend/internal/service/template/`）

- `registry.go`：启动时扫描 `backend/templates/`，逐模板校验（复刻入库校验的 Go 版：文件齐全、template.json schema、index.html 有挂载标记、layouts.md 解析出的 id 集合与 layouts[] 一致、variants class 存在于 style.css）。**校验失败拒绝注册并 log.Error，半成品模板不上线**。
- `instantiate.go`：模板 → deck 实例化（D13）。步骤：
  1. 建 `data/decks/<id>/`；
  2. index.html = 模板骨架，剥离 SLIDES 区间内 demo sections，body 挂所选 variant class，`<title>` 换成大纲标题；
  3. 拷贝 style.css 进 deck 目录；
  4. 写 outline.json（已确认大纲，version=1）、deck.json（元数据）；
  5. **引用核验**：解析 index.html 每个 href/src/url()，归属两类——`/assets/deck-v2/*`（查共享资产存在）与本地相对路径（查 deck 目录内文件存在），断链即整体失败回滚；
  6. 更新 DB Deck 行（template_id/variant/stage）。
- `manifest.go`：**类名清单构建**（§9.2 的数据源）——注册时用正则 `\.([a-zA-Z][a-zA-Z0-9_-]*)` 分别解析 base.css 与模板 style.css 的类选择器，再并入 layouts.md 各版式显式声明的"合法类名清单"，产出三份清单缓存：baseClasses / templateClasses / layoutClasses[layoutID]。layouts.md 的显式清单是**每页校验的权威**，前两者是兜底并集。

### 5.6 画廊预览服务

- `GET /templates/<id>/preview/index.html`：直接静态服务模板的 index.html（demo 数据完整可交互），挂公开路由（模板不含用户数据，免鉴权；`/assets/deck-v2/*` 本就是公开静态）。画廊 iframe 懒加载此 URL。

### 5.7 图标策略

模板自带 inline SVG 图标（html-ppt-skill 模板现状）；lint 禁止 `<script>` 图标加载器与外链图标 CSS。v1 的 Tabler 白名单（`backend/web/assets/icons.md` + `read_icons` 工具）退役。

### 5.8 首批模板清单（D9 确认）

| 模板 id | 来源（html-ppt-skill/templates/full-decks/） | 场景 |
|---------|------|------|
| tech-sharing | tech-sharing | 技术分享（GitHub dark + JetBrains Mono） |
| pitch-deck | pitch-deck | 商业路演（YC 风 10 页结构） |
| product-launch | product-launch | 产品发布 |
| weekly-report | weekly-report | 工作汇报/周报 |
| course-module | course-module | 课程/培训 |
| knowledge-blueprint | knowledge-arch-blueprint | 知识分享/架构讲解 |
| editorial-white | xhs-white-editorial | 杂志图文/轻内容 |
| data-dark | graphify-dark-graph | 数据/研究型暗色 |

暂缓池：hermes-cyber-terminal、obsidian-claude-gradient、testing-safety-alert、xhs-pastel-card、dir-key-nav-minimal、xhs-post（3:4 竖版）。**presenter-mode-reveal 弃用**（reveal 专用，与 D1 冲突）。

**上架标准（缺一不可）**：template.json 完整 / layouts.md ≥12 版式且骨架可渲染 / rules.md 就绪 / 中文黄金样张（3-5 页真实中文内容）通过 vision 审查并存入 preview/ / 画廊 cover.png 就绪 / 入库校验 + 契约单测通过。

### 5.9 模板贡献规范（开源协作）

`backend/templates/CONTRIBUTING.md`：目录契约、template.json 字段说明、layouts.md 条目格式、上架自检清单（= §5.8 标准）、`validate-template.mjs` 用法。社区模板走同一注册校验，不改核心代码。

## 6. 数据模型与产物格式

### 6.1 deck-v2 目录（D13）

```
backend/data/decks/<id>/
  deck.json        # {id, template_id, variant, stage, canvas:{w,h}, title, format:"v2", created_at, updated_at}
  outline.json     # 已确认大纲（附录 B），version 乐观锁
  index.html       # 骨架 + sections；AGENT-EDITABLE 区间仍是 LLM 唯一可写区
  style.css        # 实例化时从模板拷贝（快照隔离）
  history/         # 沿用现有机制；快照对象 = index.html + style.css + outline.json 三件
  exports/         # 导出产物（pdf/png.zip/single.html）
```

- index.html 延续三层结构：框架/脚本不可改；AGENT-EDITABLE 区间内每页一个 `<section class="slide …" data-layout="…" data-id="sN">`；讲稿放每页内 `<div class="notes">`（runtime 原生支持 N 抽屉/演讲者模式读取，display:none）。

### 6.2 DB 变更（`backend/internal/store/model.go` + `store.go`）

Deck 表新增列（GORM AutoMigrate 自动迁移）：

```go
Stage      string `gorm:"index;default:draft"`   // draft/clarifying/outlining/outline_review/selecting_template/generating/iterating
TemplateID string
Variant    string
Format     string `gorm:"index;default:v2"`
```

ChatSession 结构不变（PendingAsk 协议、Messages JSON 沿用）。`store.go` 的 AutoMigrate 会自动加列，无需手工迁移脚本。

### 6.3 outline.json 双通道写路径

- 面板直改：`PUT /api/decks/:id/outline`，整份替换 + version CAS（携带客户端持有的 version，服务端不一致返回 409 + 最新版）；
- agent 修改：`update_outline` 工具（仅 outline_review 阶段可用），内部同样走 version CAS，agent 修改前强制 `read_outline` 拿最新版；
- 两条路径都过 schema 校验（role 枚举、页码唯一递增、title 非空、points ≤6 条）；成功后广播 SSE `outline_updated`，前端面板刷新。

### 6.4 v1 存量处置（D7）

- deck 列表与所有 deck 读取端点过滤 `format = v2`（v1 行自然隐藏，不作废不迁移）；
- v1 的 `data/decks/*` 磁盘文件保留不动（不删用户盘上数据，但产品内不可见）；
- P5 删除 v1 渲染代码后，v1 deck 连查看能力一并消失——已确认接受。

## 7. 会话与状态机

### 7.1 阶段定义与迁移守卫

```
draft → clarifying → outlining → outline_review → selecting_template → generating → iterating
```

| 迁移 | 触发方 | 服务端守卫 |
|------|--------|-----------|
| draft → clarifying | 首条 chat 消息 | 自动 |
| clarifying → outlining | agent 调 submit_outline | 工具层守卫 |
| outlining → outline_review | submit_outline 成功落库 | 自动 |
| outline_review → selecting_template | POST /outline/confirm | outline 存在且 schema 合法 |
| selecting_template → generating | POST /template（选模板+变体）后自动，或显式 POST /generate | template_id 合法且已注册 |
| generating → iterating | 生成 run 正常结束 | run 结束钩子 |
| 任意 → clarifying/outlining 回退 | 不允许（重开新 deck） | — |

- stage 存 `Deck.stage`，所有阶段相关端点先校验当前 stage，非法迁移返回 409。
- clarifying 的暂停/恢复**完全复用** `ask_user` + `ChatSession.PendingAsk` + `POST /api/chat/answer` 协议（`handler/answer.go`），不新造机制。

### 7.2 各阶段 LLM 预算（maxTurns，进 config 可调）

| 阶段 | maxTurns | 说明 |
|------|:---:|------|
| clarifying | 8 | ask_user 一到两轮足够 |
| outlining | 10 | 含 web_search 素材搜集 |
| outline_review | 8 | 大纲修订 |
| generating | 24 | 1(plan) + ceil(pages/3)(写批) + ≤4(修复轮) + 余量 |
| iterating | 20 | 沿用 v1 值 |

`review_slides` 硬配额 3 次/run 与退还规则（`agent/vision_review.go`、`agent/runrecorder.go`）原样保留。预算耗尽注入收尾消息的现有机制（`agent.go` runLoop）沿用。

### 7.3 系统提示词重构（`agent/prompts/`，替代 systemPrompt.md）

```
backend/internal/agent/prompts/
  shared.md       # 角色 + 文案纪律（taste 移植通用规则）——所有阶段拼接在头部
  clarify.md      # 澄清问题清单（受众/时长→页数换算/素材/硬约束）
  outline.md      # 大纲纪律（叙事弧、页数换算、每页信息密度、role 语义）
  generate.md     # 生成纪律：版式锁、先 read_layout 再写、内容约束、节奏；动态注入所选模板 layouts[] 索引（id+用途）与 rules.md 全文
  iterate.md      # 迭代纪律（类名契约、最小改动）
```

- `agent/prompt.go` 的 //go:embed 机制改为嵌入整个 prompts/ 目录，新增 `buildStagePrompt(stage, deck) string`：shared + 阶段文件 + 动态注入（deck_id、当前日期保留；**preset 清单注入删除**；generating 阶段追加模板索引与 rules.md）。LF 归一化逻辑保留。
- 设计规则（字号分档/图片比例/色彩禁忌）全部下沉模板 layouts.md 与 rules.md，提示词只写"怎么用契约"。

### 7.4 生成 run 的发起与断线恢复

- **发起**：新 `agent.StartGenerationRun(session, deck)`。`POST /api/decks/:id/generate` 校验 stage 后调用，**响应即 SSE 流**（与 `handler/chat.go` 的 serveAgentSSE 同一套 channel/事件管线，前端复用 `lib/sse.ts` 消费）。
- run 内消息按 v1 机制逐轮 `persistSession` 落库，run 结束记历史快照 + stage 迁移。
- **断线恢复（简化版）**：每一批 write_pages 落盘即持久化。前端断开重连后读 deck stage（仍 generating）与已写页数，前端展示"生成中断，已写 N/M 页"并提供"继续生成"按钮 → `POST /generate?resume=1` → 新 run 注入"从第 N+1 页继续，先 list_slides 对齐现状"指令。完全托管运行（服务端脱离请求跑完）不做，保持 v1 语义。

### 7.5 SSE 事件扩展（`agent/stream_event.go` + `agent/transcript.go`）

新增事件类型与载荷见附录 E。前端 `types/chat.ts` ViewEvent 联合类型同步扩展，历史回放投影（transcript.go）补新事件的投影规则。

## 8. 生成管线（generating 阶段）

1. **plan_pages**：agent 一次性提交 `{assignments: [{no, layout, reason?}]}`。服务端校验：页数与大纲一致、layout 已登记、role 与版式语义匹配（cover 页必须用 cover 族版式）、**节奏校验（§9.5）**。不合格整体打回（返回违规明细），agent 重排；合格存入 deck.json 的 `page_plan` 字段。
2. **write_pages 分批**（D14）：默认每批 3 页（config `deck_v2.batch_size`，2-4），页数 ≤6 时允许一批提交完。每页提交依次过：sanitize v2 → 类名契约 → AI 味 lint → 落盘（替换 AGENT-EDITABLE 区间内对应 section）。**批内按页部分成功**：3 页中 1 页契约违规，写 2 页、返回 1 页错误明细（含最接近合法类名建议）。
3. **每页自动量测**：chromedp 打开 `/api/render/:nonce`（复用现有 nonce 门票通道，`handler/render.go` + `vision/grant.go`），定位到该页（P0 spike 定：`#/N` 或 `?preview=N`），跑 measureJS v2（附录 F）：溢出 / 最小字号 / 字数 / 版式指纹 / 底部留白。结果拼进 write_pages 返回值（沿用 v1 write_deck 的 review 字段模式；跑失败显式说明"没跑成 ≠ 没问题"）。
4. **修复循环**：量测/lint 不合格页，agent 在后续批或收尾修复轮重写；`review_slides` 看图审查保留，看图发现问题同样只重写问题页。修复轮上限 4，超过则带问题清单收尾（向用户如实报告未解决项）。
5. **收尾**：run 正常结束 → `deck.RecordRunVersion` 快照三件套（`service/deck/history.go` 的快照对象改三件）→ stage → iterating → 前端预览刷新（`deckTouched` 信号机制沿用，`DECK_MODIFYING_TOOLS` 更新为 v2 工具名）。

## 9. 质量体系

### 9.1 规则总表

| 规则 ID | 内容 | 级别 | 执行点 |
|---------|------|------|--------|
| C201 | section 缺 `data-layout` 或未登记 | **阻塞**（该页拒收） | write_pages / update_slide / insert_slide |
| C202 | 使用未声明类名 | **阻塞**（返回最接近合法类名建议，编辑距离 ≤2） | 同上 |
| C203 | 危险标签 / on* / javascript: | **阻塞**（沿用 sanitize.go 现有防线） | 同上 |
| R101 | 同版式连续 ≥3 | **阻塞**（打回重排） | plan_pages |
| R102 | 每 8 页 <4 种版式 | 提示 | plan_pages |
| R103 | 每 8 页无满版 hero 类版式 | 提示 | plan_pages |
| R104 | 左右图文交替连续 >2 | 提示 | plan_pages |
| R105 | 实写页数 ≠ 大纲页数 | **阻塞** | write_pages（页码必须落在 outline 内） |
| T001 | 中文禁词命中 | 提示（附替代写法） | write_pages 返回 lint_report |
| T002 | 英文禁词命中 | 提示 | 同上 |
| T003 | 破折号：标题含"——"提示；正文每页 >1 提示；英文 em-dash 归零 | 提示 | 同上 |
| T004 | 假精确数字（`92%`/`99.9%` 类无出处完美数） | 提示 | 同上 |
| T005 | 页标题 >16 字（中文） | 提示 | 同上 |
| T006 | 要点条 >28 字 | 提示 | 同上 |
| T007 | eyebrow/眉题类元素每 3 页 >1 | 提示 | write_pages（按版式语义类统计） |
| S301 | stylelint v1 提示项（写死颜色/px 字号/内联字号/emoji 等） | 提示 | 沿用 stylelint.go，阈值可被模板 rules.md 覆盖 |

哲学不变：**阻塞级只挡机械可判定的问题（契约），品味级问题提示给 agent 自行修复，绝不改写用户内容**。

### 9.2 类名契约（C201/C202 实现）

- 数据源：§5.5 manifest.go 产出的 baseClasses / templateClasses / layoutClasses[layoutID]。
- 校验：section 及其后代元素的 class 拆分后 ⊆ (base ∪ template ∪ 该页 layoutClasses)。goquery 遍历（依赖已有）。
- 报错信息：`未知类名 .xxx，该版式合法类名：[…]，最接近：.yyy`。

### 9.3 sanitize v2（`service/deck/sanitize.go` 改造）

保留全部现有防线（危险标签整份拒绝、on*/javascript: 剥除）；新增：解析 section 的 data-layout（C201）、notes div 允许（class 合同内）、`<script>/<style>/<link>` 仍禁。

### 9.4 AI 味 lint（`service/deck/lint_taste.go` 新增，taste-skill 规则中文化重写）

- **中文禁词表 v1（D10 确认，config `deck_v2.lint.cjk_banned` 可改）**：赋能、抓手、闭环、沉淀、心智、护城河、组合拳、打法、底层逻辑、顶层设计、颗粒度、拉齐、对齐（非会议义，正则排除"跟…对齐日程"场景难写，先全量提示人工判断）、值得注意的是、综上所述、总而言之、众所周知、毋庸置疑、 impassion 类直译腔（"这不仅…更是…"连用 ≥2 次提示）。
- **英文禁词表（`en_banned`）**：elevate、seamless、unleash、next-gen、revolutionize、empower、cutting-edge、delve、tapestry、robust、leverage（动词义）、holistic。
- 每条命中返回 `{rule, word, occurrence, hint}`，hint 给替代写法（词表内置映射）。

### 9.5 节奏校验（`service/deck/lint_rhythm.go` 新增，guizang 思想重写）

纯函数：输入 page_plan（或实写页的 data-layout 序列），输出 R101-R105 违规明细。单测覆盖边界（恰好 3 连、8 页窗滑动）。

### 9.6 vision 审查适配（`vision/review.go` + `agent/vision_review.go`）

- reviewPrompt 更新：删除 fit 系数相关判断；判定基准改为固定画布语义（溢出、贴边、底部空白、字号过小、密度失衡）；保留【硬/软】分级与"没把握不报"。
- 空报告重试一次、配额 3 次/run、退还规则全部保留。
- capture.go 的 measureJS 重写为附录 F 版本（旧 fit 缩放逻辑删除）。

## 10. 导出（新包 `backend/internal/export/`）

| 格式 | 实现 |
|------|------|
| PDF | chromedp 一次调用 PrintToPDF：`PaperWidth=13.333in, PaperHeight=7.5in, PrintBackground=true, PreferCSSPageSize`（1920×1080 @96dpi 恰为 16:9 整页），依赖 runtime 的 @media print 逐页分页（P0 spike ② 验证） |
| PNG | 逐页 `#/N` 定位 + viewport 1920×1080 + 等字体加载 + 截图，`archive/zip` 打包 |
| 单文件 HTML | 打包器（附录 G）：内联本地 style.css 与 /assets/deck-v2 的 CSS/JS；**CJK 字体不内联**（每字重 MB 级，base64 会把文件撑爆，收益低），保留 /assets 绝对 URL + 系统字体回退；产物可分享可托管 |

- `POST /api/decks/:id/export {format}`，同步实现，超时 60s（config `export.timeout_seconds`），超时 408；产物落 `exports/`，`GET /api/decks/:id/exports/:file` 带 `?token=` 下载。
- 导出复用 chromedp 实例池（与 vision 共用，注意串行化避免资源竞争，与 §18 风险一致）。

## 11. 前端改造（真实文件映射）

### 11.1 改动点位

| 文件 | 改动 |
|------|------|
| `stores/deck.ts`（已存在，扩展） | 增 stage / template_id / variant / outline 状态；stage 迁移 action（confirmOutline / selectTemplate / startGenerate） |
| `stores/chat.ts` | `DECK_MODIFYING_TOOLS` 换 v2 工具名；新事件类型接线（附录 E）；generating 阶段锁 ChatInput |
| `types/chat.ts` | ViewEvent 联合类型 + DeckV2 类型 |
| `lib/previewRefresh.ts` | 触发工具表同步 v2 |
| `lib/sse.ts` | 无协议变化，仅事件种类增多（归一逻辑在 chat store） |
| `api/decks.ts` | 新端点包装：putOutline / confirmOutline / selectTemplate / generate(SSE) / exportDeck / listTemplates |
| `api/templates.ts`（新增） | GET /api/templates |
| `components/wizard/WizardStepper.vue`（新增） | 五段向导条，当前阶段高亮 |
| `components/wizard/OutlinePanel.vue`（新增） | 大纲卡片：拖拽排序/增删页/行内编辑，PUT 乐观更新 + 409 冲突提示；"确认大纲"按钮 |
| `components/wizard/TemplateGallery.vue`（新增） | 模板卡（cover.png + 场景标签 + 懒加载 live iframe 预览）+ variants 色板圆点 + 必选确认；按 outline.meta 场景关键词推荐排序置顶 |
| `components/preview/PreviewPane.vue` | iframe URL 改 `/api/decks/:id/file#/1`；固定画布下**删除一切缩放兜底**；生成中显示 page_generated 进度 chips |
| `views/WorkspaceView.vue` | 按 stage 切换区域可见性（WizardStepper + OutlinePanel/TemplateGallery/ChatMessages 组合），?session=N 深链保留 |
| `components/chat/ChatInput.vue` | generating/selecting_template 阶段禁用（D12/R3：selecting_template 锁输入提示"请先选择模板"） |

### 11.2 交互细节

- 大纲双通道并发：面板编辑挂起时收到 agent 的 `outline_updated` → 提示"agent 已更新大纲，点击加载最新版"，不自动覆盖用户未保存的编辑。
- 模板画廊 iframe 在卡片进入视口才加载，离开销毁（demo deck 带 WebGL/canvas 的模板耗资源）。
- 生成进度：page_generated 事件驱动 chips（页码 + 版式名 + ✓/⚠），完成后整体切换 iterating 布局。

## 12. 后端接口变更

新增（鉴权除注明外均为 JWT，iframe/下载带 `?token=` 回退，沿用 `middleware/auth.go` 机制）：

```
GET  /api/templates                       # 公开：模板列表（meta + cover 路径 + variants）
GET  /templates/:id/preview/*             # 公开静态：画廊 live 预览（router 公开段，参照 /assets）
PUT  /api/decks/:id/outline               # {outline, version} → 200{version} / 409{latest}
POST /api/decks/:id/outline/confirm       # gate 1
POST /api/decks/:id/template              # gate 2：{template_id, variant}
POST /api/decks/:id/generate[?resume=1]   # 响应即 SSE 流（复用 chat 流管线）
POST /api/decks/:id/export                # {format: pdf|png|html}
GET  /api/decks/:id/exports/:file         # ?token= 下载
GET  /api/decks/:id/assets/*              # deck 目录静态资产（style.css 等）
```

变更：`GET /api/decks/:id/file` 返回 index.html（handler/deck.go）。退役：无独立 theme 端点（theme 全在 agent 工具层），custom_css 特性随工具一并摘除。新增路由挂在 `router/router.go`，公开路径沿用 `static_mime.go` 的 MIME 处理。

## 13. 工具层改造（`agent/func_tool.go`）

**新增工具**（schema 详见附录 D）：

| 工具 | 阶段 | 说明 |
|------|------|------|
| submit_outline {outline} | outlining | schema 校验 → 落库 → stage 迁移；校验失败返回逐条错误 |
| read_outline {} | outline_review | 返回当前 outline.json + version |
| update_outline {version, outline} | outline_review | CAS 整份替换 |
| plan_pages {assignments} | generating | §8.1 |
| read_layout {layout} | generating / iterating | 返回骨架代码 + 合法类名 + 内容约束（按需取，控制 token） |
| read_guidelines {} | generating / iterating | 返回模板 rules.md |
| write_pages {pages[]} | generating | §8.2，返回每页 {ok, lint[], measure, retry_hint} |

**改造保留**：list_slides / read_slide（返回值增 data-layout、fingerprint 保留）/ update_slide / insert_slide（schema 增必填 layout，过 C201/C202）/ delete_slide / review_slides / ask_user / web_search / list_history / read_history_diff。

**退役删除**：write_deck、update_theme、read_theme、update_custom_css、read_custom_css、read_component（`agent/component_lib.go` 一并删）、read_icons（icons.md 一并删）、write_deck 的 preset 参数逻辑、`func_tool.go` 中 PresetSummary 注入。

## 14. 文件级改动清单（真实路径）

### 14.1 新增

```
backend/internal/service/template/{registry,instantiate,manifest,validate}.go + *_test.go
backend/internal/service/deck/{lint_taste,lint_rhythm}.go + *_test.go
backend/internal/export/{export,pdf,png,bundle}.go + *_test.go
backend/internal/agent/prompts/{shared,clarify,outline,generate,iterate}.md
backend/web/assets/deck-v2/{runtime.js,base.css,fonts.css,fonts/,PATCHES.md}
backend/templates/<8 个模板目录>/ + CONTRIBUTING.md + tools/validate-template.mjs
frontend/src/api/templates.ts
frontend/src/components/wizard/{WizardStepper,OutlinePanel,TemplateGallery}.vue + 对应 .test.ts
```

### 14.2 重写/大改

```
backend/internal/agent/agent.go          # 阶段化 prompt 装配、StartGenerationRun、preset 注入删除、maxTurns 按阶段
backend/internal/agent/func_tool.go      # 工具表 v2（§13）
backend/internal/agent/prompt.go         # embed 目录 + buildStagePrompt
backend/internal/agent/stream_event.go   # 新事件类型
backend/internal/agent/transcript.go     # 新事件投影
backend/internal/agent/vision_review.go  # 量测结果结构对齐 v2
backend/internal/handler/chat.go         # SSE 管线复用出 generate 流
backend/internal/handler/deck.go         # outline/template/generate/export/assets 端点
backend/internal/router/router.go        # 新路由 + 公开段
backend/internal/service/deck/create.go  # v2 建 deck 流程
backend/internal/service/deck/deck.go    # 三件套读写、AGENT-EDITABLE 区间按新骨架
backend/internal/service/deck/sanitize.go    # +C201
backend/internal/service/deck/stylelint.go   # 阈值可被 rules.md 覆盖
backend/internal/service/deck/history.go     # 快照对象三件套
backend/internal/service/deck/slide_ops.go   # insert 的 layout 必填与契约校验
backend/internal/vision/capture.go       # measureJS v2（附录 F）
backend/internal/vision/review.go        # reviewPrompt v2
backend/internal/store/model.go          # Deck 四列
backend/config.yaml + internal/config/config.go  # §15 新键
backend/web/chat-test.html               # 阶段驱动测试台
frontend/src/stores/{deck,chat}.ts、types/chat.ts、lib/previewRefresh.ts、api/decks.ts
frontend/src/components/preview/PreviewPane.vue、components/chat/ChatInput.vue
frontend/src/views/WorkspaceView.vue
```

### 14.3 退役删除（P5 执行）

```
backend/web/assets/{reveal.js,reveal.css,moon.css,theme.css,components.css,init.js,icons.md,
                    attack.js,components-test.html,fit-test.html,layouts-test.html,theme-test.html}
backend/internal/service/deck/{theme.go,theme_presets.go,theme_vars.go,variable_contract.go,
                               custom_css.go}  # 及对应 *_test.go
backend/internal/agent/component_lib.go      # 及测试
backend/internal/agent/systemPrompt.md       # 被 prompts/ 取代
工具：write_deck/update_theme/read_theme/update_custom_css/read_custom_css/read_component/read_icons
```

### 14.4 明确保留不动

agent 循环与 ask_user 协议（agent.go runLoop、toolcalls.go、handler/answer.go）、persistSession、repairToolCalls、web_search、trace 全家（trace/、handler/trace.go、TraceView）、历史机制骨架（history.go / HistoryDrawer.vue）、鉴权体系（auth/middleware/authctx/cryptox）、前端 UI 组件库与主题（components/ui、theme/tokens.css）、vite 代理（vite.config.ts）。

## 15. config.yaml 新增键

```yaml
templates:
  dir: "./templates"                 # 相对 backend/，默认内嵌仓库目录

deck_v2:
  batch_size: 3                      # write_pages 每批页数，2-4
  max_turns:
    clarify: 8
    outline: 10
    outline_review: 8
    generate: 24
    iterate: 20
  lint:
    cjk_banned: [赋能, 抓手, 闭环, 沉淀, 心智, 护城河, 组合拳, 打法, 底层逻辑, 顶层设计, 颗粒度, 拉齐, 值得注意的是, 综上所述, 总而言之, 众所周知, 毋庸置疑]
    en_banned: [elevate, seamless, unleash, next-gen, revolutionize, empower, cutting-edge, delve, tapestry, robust, leverage, holistic]
    title_max_chars: 16
    bullet_max_chars: 28
    eyebrow_per_pages: 3

export:
  timeout_seconds: 60
```

`internal/config/config.go` 增对应结构体与默认值（零值可用，测试不依赖 config 文件）。

## 16. 实施计划（每阶段独立可合并，P5 前旧管线不破坏）

### P0 底座与技术验证
- vendor runtime.js/base.css + 字体自托管落 `deck-v2/`；PATCHES.md 建立；
- `service/template` 四件（registry/instantiate/manifest/validate）+ 契约单测；
- deck-v2 目录结构与骨架渲染；`GET /templates/:id/preview` 公开路由；
- **三个 spike**：单页定位方式（#/N vs ?preview=N）、PrintToPDF 分页、iframe sandbox 初始化；
- tech-sharing 一个模板完成入库适配全流程（作为流水线试运行）。
- **DoD**：curl 建 v2 deck → 浏览器渲染 demo 正常、深链翻页、单页截图、PrintToPDF 出 16:9 PDF；registry 对坏模板（缺文件/断链/id 不一致）拒绝注册有单测。

### P1 agent 管线（后端闭环）
- 阶段状态机 + 端点（outline/confirm/template/generate）+ prompts 五件套 + 工具 v2 + sanitize v2/类名契约 + measureJS v2 + 修复循环预算 + 断线 resume。
- **DoD**：chat-test.html（改造成阶段测试台）用真实 LLM 走通 澄清→大纲→确认→（curl 选模板）→生成→页级修改 全流程；坏类名/坏版式被拒且报错可指导修复；trace 里 run 带 stage。

### P2 前端向导
- WizardStepper + OutlinePanel + TemplateGallery + stage 接线 + 新 SSE 事件 + 进度 chips + PreviewPane v2；variants 色板（D11）。
- **DoD**：浏览器全流程无 chat-test 走通；双通道大纲编辑含 409 冲突场景；画廊懒加载不卡。

### P3 质量层完整版 + 模板批量上架
- AI 味 lint + 节奏校验 + reviewPrompt v2 + 修复轮；剩余 7 个模板按 §5.8 标准逐个上架（横跨 P1-P3，上架一个用一个）。
- **DoD**：注入坏样张（未登记版式/禁词/3 连版式/超长标题）全部被阻塞或提示且 agent 能自愈；每模板中文黄金样张归档 preview/。

### P4 导出
- PDF / PNG zip / 单文件打包器 + 下载端点。
- **DoD**：8 页 deck 三种导出与预览逐页一致；超时路径有测试。

### P5 摘除旧栈
- §14.3 删除清单执行；v1 deck 隐藏（format 过滤）；THIRD_PARTY_NOTICES.md；README / backend/README 分层说明 / docs 更新。
- **DoD**：全仓 grep 无 reveal/preset 残留；`go build ./... && go test ./...` 与前端 vitest 全绿；附录 H 回归清单过一遍。

## 17. 第三方许可与署名

- **html-ppt-skill（MIT）**：vendor runtime/base.css、模板衍生，保留其 LICENSE，入 `THIRD_PARTY_NOTICES.md`；runtime 改动记 PATCHES.md。
- **taste-skill（MIT）**：规则中文化重写进提示词与 lint，属思想借鉴；NOTICES 致意 + 仓库链接。
- **guizang-ppt-skill（AGPL-3.0）**：**零文件、零代码引入**，仅思想参考（版式锁/节奏/checklist 分级）；本节记录边界，CONTRIBUTING 写明禁止复制 AGPL 代码进模板。
- 三仓库链接进 README 致谢区。

## 18. 风险与缓解

| 风险 | 影响 | 缓解 |
|------|------|------|
| runtime.js 为 file:// 静态场景设计，iframe/sandbox/print 集成未知 | P0 卡壳 | P0 三个 spike 前置；MIT 允许打补丁且 PATCHES.md 记录；最坏 fork 维护（1162 行可控） |
| 模板英文语境设计，中文排版未验证 | 中文样张质量 | 上架标准强制中文黄金样张过 vision；Noto Sans SC 在其字体栈内 |
| DeepSeek 对版式锁遵从不足 | 写入失败率升高 | C202 报错带最接近合法类名 + 修复循环；generate.md 注入版式索引；契约是机械校验不依赖模型自觉 |
| layouts.md 编写量大（8×12-18） | P1-P3 周期 | 写半自动提取脚本从模板 index.html 抽骨架人工修订；先 1 个模板（P0）+3 个（P1）跑通链路再批量 |
| 逐批生成用语/密度漂移 | 观感割裂 | plan_pages 全局前置 + 节奏校验 + shared.md 文案纪律 + 同 run 对话记忆 + 黄金样张比对 |
| 大纲双通道并发冲突 | 编辑丢失 | version CAS + 409 + 前端显式加载提示；agent 修改前强制 read_outline |
| chromedp 与 SSE 同进程资源竞争（既有） | 高峰变慢 | 导出与 vision 串行化共用实例池；固定画布使量测失败率大降；实测后再议拆进程 |
| 生成中断线 | 用户以为丢了 | 每批落盘 + resume 端点 + 前端"继续生成"入口（§7.4） |
| 字体子集工作量 | P0 拖期 | 直接 vendor 完整 woff2（3-6MB 可接受），子集化延后 |

## 19. 开放问题（已全部决议，2026-09-17）

见 §3 D9-D12 与 D13-D14：首批模板认可；禁词表提示级认可；variants v1 就做；生成中锁输入；deck 存储混合形态确认；每批 2-4 页（默认 3）确认。无未决项。

---

## 附录 A：template.json 完整示例

```json
{
  "id": "tech-sharing",
  "name": "极客深色",
  "description": "GitHub 暗底 + JetBrains Mono，面向工程团队的技术分享",
  "tags": ["技术", "深色", "代码"],
  "scenario": ["技术分享", "团队培训", "工程", "开发"],
  "canvas": { "w": 1920, "h": 1080 },
  "variants": [
    { "id": "default", "name": "经典深色", "class": "" },
    { "id": "forest", "name": "森林夜", "class": "v-forest" },
    { "id": "ember", "name": "暖橙", "class": "v-ember" }
  ],
  "layouts": [
    { "id": "cover", "name": "封面", "use": "开场：标题+副标+作者", "roles": ["cover"] },
    { "id": "code-walkthrough", "name": "代码走读", "use": "代码块+要点解说",
      "constraints": "代码 ≤20 行；说明要点 2-4 条" }
  ],
  "fonts": ["Noto Sans SC", "JetBrains Mono"],
  "source": { "derived_from": "html-ppt-skill/templates/full-decks/tech-sharing", "license": "MIT" }
}
```

## 附录 B：outline.json schema

```json
{
  "version": 1,
  "title": "Go 语言入门",
  "meta": { "audience": "后端转语言的工程师", "duration_min": 20, "page_count": 12, "tone": "务实" },
  "narrative": { "hook": "一句话钩子", "arcs": ["为什么是 Go", "核心心智模型", "上手实战", "收束"] },
  "pages": [
    {
      "no": 1,
      "role": "cover",
      "title": "Go 语言入门",
      "points": ["副标题或一句话定位"],
      "layout_hint": "cover",
      "materials": [{ "type": "image", "desc": "可选素材描述" }],
      "notes": "讲稿要点（生成阶段补全进 .notes）"
    }
  ]
}
```

校验规则：`role ∈ cover|toc|divider|content|data|quote|code|cta|thanks`；no 唯一且从 1 连续递增；title 非空 ≤40 字；points ≤6 条；version 单调递增。

## 附录 C：layouts.md 条目格式

```markdown
## L07 · stat-hero（数据大字报）

用途：单个核心指标需要冲击力呈现。
适用 role：data。
内容约束：主数字 ≤6 字符；指标名 ≤12 字；支撑要点 2-3 条、每条 ≤20 字。

<section class="slide stat-hero" data-layout="stat-hero" data-id="{{ID}}">
  <div class="kicker">{{引导语 ≤12 字}}</div>
  <div class="mega">{{主数字}}</div>
  <div class="metric-name">{{指标名}}</div>
  <ul class="support"><li>{{支撑要点}}</li></ul>
  <div class="notes">{{讲稿}}</div>
</section>

合法类名：slide, stat-hero, kicker, mega, metric-name, support, notes
```

## 附录 D：新工具 schema 草案

```json
{
  "name": "write_pages",
  "description": "按 page_plan 分批写入页面。每页必须用已登记版式，先 read_layout 取骨架。",
  "parameters": { "pages": [{ "no": 3, "layout": "stat-hero",
    "html": "<section class=\"slide stat-hero\" data-layout=\"stat-hero\" data-id=\"s3\">…</section>",
    "notes": "讲稿（写入该页 .notes div）" }] },
  "returns": { "results": [{ "no": 3, "ok": true,
    "lint": [{ "rule": "T001", "word": "赋能", "hint": "改为具体动作" }],
    "measure": { "overflow": false, "min_font_px": 18, "words": 86, "bottom_gap_px": 120 },
    "retry_hint": "" }] }
}
```

`plan_pages`：`{assignments: [{no: 1, layout: "cover", reason: "可选"}]}`，返回 `{ok, violations[]}`。
`read_layout`：`{layout: "stat-hero"}`，返回 `{skeleton, classes[], constraints}`。
`submit_outline` / `update_outline`：参数即附录 B 形状（update 带 version）。

## 附录 E：SSE 新事件载荷

```
stage_changed   {deck_id, from, to}
gate_waiting    {deck_id, gate: "outline"|"template"}
outline_updated {deck_id, version, source: "user"|"agent"}
page_generated  {deck_id, no, total, layout, ok, measure?}
lint_report     {deck_id, page, level: "warn"|"block", items: [{rule, word?, msg, hint}]}
```

## 附录 F：measureJS v2 伪代码（固定画布，无缩放换算）

```js
el = 当前页 section（按 spike 选定的定位方式）
({
  overflow: el.scrollHeight > H + 2 || el.scrollWidth > W + 2,   // 设计像素，transform scale 不影响
  min_font_px: 遍历可见文本节点 computedStyle.fontSize 取最小,
  words: textContent 的 CJK 字符 + 英文单词计数,
  layout: el.dataset.layout,
  bottom_gap_px: H - 最后一个可见子元素 bottom,
  density: 直接子元素数
})
```

## 附录 G：单文件打包器步骤

1. 读 index.html；2. 每个 `<link rel=stylesheet>`：本地 style.css 与 `/assets/deck-v2/*.css` → 内联 `<style>`（带来源注释）；3. 每个 `<script src>`（/assets/deck-v2/*）→ 内联；4. 内联 CSS 中 `url()` 指向 fonts/ 的保持绝对 URL（不内联 CJK 字体）；5. 输出 `exports/single.html`；6. 断链引用 → 导出失败并报明。

## 附录 H：P5 回归清单

注册/登录 → 新建 deck → 澄清问答（ask_user）→ 大纲面板直改 + 对话改各一次 → 确认 → 画廊选模板 + 换 variant → 生成（进度 chips 正常）→ 预览翻页/概览/演讲者模式 → 迭代：update_slide / insert_slide / delete_slide 各一次 → 历史快照恢复 → 导出 PDF/PNG/单文件 → /trace 看到 generating run 与 stage → 旧 v1 deck 不出现在列表 → 断线重连出现"继续生成"。
