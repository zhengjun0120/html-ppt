# plan-v3 · 模板库扩充 + 用户自定义模板 + 前端向导分步

> 状态：**待审查**。审查通过后按此文执行；执行期间不再提问，除非遇到文中"开放问题"或事实性阻塞。
> 前置已锁定决策：D1 克隆定制（自定义模板从内置模板派生）/ D2 自动门禁（公开无人工审核）/
> D3 向导式独立路由页 / D4 搬运剩余 7 个 + taste 风格原创 2-3 个。

---

## 0. 总览与执行顺序

三条工作流：

| 工作流 | 内容 | 风险 | 依赖 |
|---|---|---|---|
| **A 模板库扩充** | 移植 7 个 + taste 风格原创 2-3 个 → 内置模板 8→17~18 | 低（管线成熟） | 无 |
| **C 前端向导分步 + 预览强化** | 五步拆成独立路由页；deck 预览（缩略图/翻页）；主题预览（变体实时切换） | 中（前端重构） | 无 |
| **B 用户自定义模板系统** | fork → 对话定制 → 自动门禁发布 → 社区使用/管理 | 高（全栈新系统） | A 的管线标准；C 的页面框架 |

**执行顺序：A → C → B。** 理由：A 纯搬运零风险先落袋；C 是用户可见收益最大的重构，且 B 的
管理页（/templates/mine、/templates/:id/edit）要长在 C 的页面框架上；B 最后做，直接复用
A 固化下来的模板质量标准与 C 的画廊/预览组件。

每条工作流独立成 commit，可单独回滚。全程本地提交、不推送。

---

## 1. Phase A · 模板库扩充（8 → 17~18）

### A1 移植 html-ppt-skill 剩余 7 个（MIT，来源 `D:\go_files\html-ppt-skills\html-ppt-skill-main\...\templates\full-decks\`）

| 源模板 | 入库 id | 备注 |
|---|---|---|
| hermes-cyber-terminal | hermes-terminal | 赛博终端风，与 tech-sharing 差异化（更浓的霓虹） |
| obsidian-claude-gradient | obsidian-gradient | 暗底渐变 |
| xhs-pastel-card | xhs-pastel | 小红书粉彩卡片 |
| xhs-post | xhs-post | 小红书图文 |
| dir-key-nav-minimal | dir-nav-minimal | 极简目录导航风 |
| presenter-mode-reveal | presenter-cards | **更名**，避免与已废弃的 reveal.js 混淆 |
| testing-safety-alert | safety-alert | 告警/安全通报风 |

每个模板走既有管线，**硬性达标线（本仓库已固化的质量标准，缺一不入库）**：

1. `tools/scaffold-template.mjs` + spec → materialize 生成骨架；
2. 手工适配：`.tpl-<id>` 作用域前缀、中文 demo（lang="zh-CN"）、1-3 个 variants、
   demo 完整覆盖全部版式；
3. layouts.md：每个版式带**指纹行**、内容约束用**区间**（说明 22-45 字等）、
   骨架类名 ⊆ 合法类名清单、cover/收尾页垂直居中（`slide full` 或等价）；
4. 字号投影 floors：正文 ≥18px、半内容 ≥16、装饰小字 ≥14、代码 ≥17（无 <14px 残留）；
5. `validate-template.mjs` 全过 + Go registry 装载通过；
6. demo 渲染量测：无溢出、填充率达标、截图人工目视（标准同 deck-0032 验收）；
7. `THIRD_PARTY_NOTICES.md` 更新搬运清单（MIT 归属）。

### A2 taste 风格原创 2-3 个（输入：taste-skill MIT 的 skills/{brutalist-skill, minimalist-skill, soft-skill}）

| 新模板 id | 风格来源 | 说明 |
|---|---|---|
| brutalist-bold | brutalist-skill | 粗野主义大字报：高对比、硬边框、超大号字 |
| minimal-quiet | minimalist-skill | 极简留白：大标题+极少文字，靠字号层级撑画面 |
| soft-pastel | soft-skill / brandkit | 柔和粉彩圆角卡（与 xhs-pastel 差异化：更通用的工作坊风） |

做法：从 taste skill 的 SKILL.md 提炼设计规范（色板/字体气质/布局原则）写进 ADAPTATION.md，
**结构契约复用内置模板的成熟骨架族**（hero/stack/cards/split/code/chart/table 指纹齐备，
pattern 数 ≥4），同样走 A1 的 7 条达标线。license 记 MIT（借鉴自 taste-skill）。

### A3 全库回归

17~18 个模板全部重跑 demo 渲染量测 + 截图；judge 抽查 3-4 个新模板。 commit 切分：
A1 分 2-3 个 commit（每批 2-3 个模板）、A2 一个 commit。

---

## 2. Phase C · 前端向导分步 + 预览强化

### C1 路由重构（五步 = 五个 URL）

```
/decks            文稿列表（保留）
/new              第 1 步 · 澄清：全屏对话（顶部 WizardStepper）
/new/outline      第 2 步 · 大纲：全屏大纲面板（右侧保留对话抽屉，可收起）
/new/template     第 3 步 · 模板：全屏画廊 + 主题预览（无对话，输入锁定语义不变）
/new/generating   第 4 步 · 生成进度：全屏进度页
/decks/:id        第 5 步 · 成品预览 + 迭代（双栏：预览 + 对话，即现在的工作台瘦身版）
```

- 旧 `/decks/new` 301 到 `/new`；`?session=N` 深链语义保留（/new 与 /decks/:id 都支持）。
- **步骤守卫**：进阶路由时按 deck.stage 校验——没到对应步骤就重定向到当前步骤页；
  stage 被 agent 推进时自动向前跳路由（沿用 wizard store 的 stage watch，把"切主区域"
  改成"推路由"）。
- WizardStepper 变成可点击的进度导航：已完成步骤可回退查看（只读或可编辑按步骤语义）。
- WorkspaceView 拆解：对话侧栏抽成 `components/chat/ChatPanel.vue` 复用；
  主区域五个形态各自成 view；wizard store 增加 `routeForStage(stage)` 映射。
- 回退语义：outline 步可回 /new 继续聊需求（stage 回 outlining 由 agent 管，前端只跳页）；
  template 步是后端 gate，回退仅展示。

### C2 主题/模板预览强化（/new/template）

- **变体实时预览**：后端新增 `GET /api/templates/:id/preview?variant=<vid>` ——
  服务端读模板 demo，把 body 的 variant class 换成查询参数指定的变体后返回
  （内置模板 demo body 挂 `tpl-<id> <variant-class>`，替换是纯字符串操作）；
  前端画廊卡片内的 iframe 换变体时只换 src 的 query，即时换肤。
- **整本 demo 翻页预览**：demo 是带 runtime 的完整 index.html，支持 `#/N` 深链——
  卡片展开态给「上一页/下一页/页码」控制条，直接改 iframe hash，零后端成本。
- 卡片信息补全：版式数、变体色点、适用场景 tags；「预览大图」按钮 = 展开为 60% 宽的
  大 iframe + 翻页条。
- 预留「我的模板」「社区模板」区块位（B 落地前隐藏）。

### C3 成品预览强化（/decks/:id 与列表页）

- **缩略图服务（后端）**：`GET /api/decks/:id/thumbs/:no` —— headless 一次过渲染整本
  deck 输出每页 PNG，落盘 `data/decks/<id>/thumbs/<no>.png`；按 deck 内容版本号做缓存 key
  （版本不变直接回文件）；**失效钩子**挂在 writeSlideSegment / restore / UpdateSlide /
  InsertSlide / DeleteSlide 后清空 thumbs 目录；`GET /api/decks/:id/thumbs` 返回清单
  （页数、是否已生成）。第一页在生成完成时异步预生成（generate run 收尾时），
  其余页懒加载（首次请求触发整本渲染）。
- **预览翻页**：PreviewPane 重做——左侧缩略图栏（点击跳页、当前页高亮）+ 主体 iframe
  （`#/N` hash 驱动）+ 底部翻页条（上一页/下一页/页码输入/键盘 ←→）；缩略图与 iframe
  双向联动。
- **文稿列表封面**：DecksView 卡片封面改为缩略图第一页（`/decks/:id/thumbs/1`，
  懒加载 + 骨架屏占位），点击直达 /decks/:id。
- 预览工具条加「导出」入口（PDF/PNG/单文件，后端已有，补 UI + 下载态）。

### C4 质量线

vue-tsc / vitest 全过；wizard 路由同步逻辑加单测（stage→route 映射、守卫重定向）。
commit 切分：C1 路由拆分一个、C2 预览强化一个、C3 缩略图+翻页一个。

---

## 3. Phase B · 用户自定义模板系统（克隆定制 D1 + 自动门禁 D2）

### B0 数据模型与存储

```
user_templates 表（GORM AutoMigrate；store.OpenMemory 同步建表供测试）：
  ID            string  "ut-0001"（复用 deck-xxxx 的编号风格）
  UserID        uint
  BaseID        string  fork 来源内置模板 id（如 "tech-sharing"）
  Name/Desc     string
  Dir           string  data/user-templates/ut-0001/
  Visibility    string  private | public
  Status        string  draft | publishing | published | failed
  PublishError  string  门禁失败原因（给 UI 展示）
  PublishReport string  三关门禁的结果 JSON（渲染量测数字/截图引用）
  CreatedAt/UpdatedAt
```

目录即模板：`data/user-templates/<id>/{template.json, index.html, style.css, layouts.md,
rules.md, ADAPTATION.md（记录 fork 来源与定制历史）}`。删除 = 删目录 + 删行。

### B1 后端服务 `internal/service/usertpl`

- **Fork**：`POST /api/templates/:id/fork` → 复制内置模板四件套 → 改 id/name/owner/
  来源标注（ADAPTATION.md 记 base）→ 返回 ut id。复制后立刻跑一次 demo 渲染存预览图。
- **CRUD + 可见性**：
  - `GET /api/user-templates`（我的列表）
  - `GET /api/user-templates/:id`（meta；owner 或 public 可读）
  - `PUT /api/user-templates/:id/meta`（改名/描述）
  - `DELETE /api/user-templates/:id`（公开中的先隐式下架）
  - `POST /api/user-templates/:id/publish` / `POST /:id/unpublish`
- **定制 = agent 改文件（白名单式）**。新 agent 工具集（复用 mountTool/decodeToolArgs/
  emitV2 骨架，新阶段 `template_customizing`）：
  - `read_template` / `write_template_file`：可写路径白名单 = style.css 的 **token 块**
    （要求所有用户模板 style.css 的 `:root`/`.tpl-x` token 定义集中在文件头部
    `/* == TOKENS == */` 标记之间，程序化改色/字体只动这一段）、template.json 的
    name/description/variants（色值）、rules.md 文案、layouts.md 的约束行（不可加删版式——v1 不开放结构级增删，防契约烂掉）；
  - `set_palette` / `set_fonts`：高层工具，直接改 token 块（主色/辅色/背景/文字/圆角/
    字体族；字体只允许自托管清单内的 Inter/Noto Sans SC/JetBrains Maple Mono/serif 栈切换）；
  - `preview_template`：实例化 demo → CaptureV2 渲染 → 量测 + 截图回给模型（复用 vision
    nonce 管线），模型自己看效果再迭代；
  - `submit_template`：触发发布门禁（见下）。
  - 对话入口：`POST /api/user-templates/:id/chat`（SSE，复用 serveAgentSSE）；系统提示词
    = 定制阶段提示词（新 prompts/customize.md）+ 模板结构摘要 + token 现值。
- **发布门禁（自动三关，串行，任一失败 → status=failed + PublishError）**：
  1. **结构校验**：Go registry loadTemplate 同一套规则（指纹/类名自洽/网格陷阱/资产引用），
     外加用户模板专属：token 块完整、无禁用 CSS（见 B2 安全）、demo 页 lang=zh-CN；
  2. **渲染量测**：CaptureV2(demo) 全页无溢出、fillPct 达标（与生成验收同一阈值）、
     截图落盘 PublishReport 引用；
  3. **一页冒烟**：用该模板走最小生成环（read_layout → write_pages 1 页，一次 LLM 调用，
     额度记在发布者账号），页必须过全部写入门禁（类契约/密度/溢出）。
  全过 → visibility=public；注册表热加载（见下）。
- **Registry 扩展**：`Registry` 支持**动态用户模板层**——启动时扫描
  `data/user-templates/` 下 published 的目录注册；publish/unpublish/delete 时增量
  注册/注销（`sync.Mutex`；`Get` 先查内置再查用户层；`List` 合并并带 `origin: builtin|user`
  与 owner 信息）。**内置模板永远只读**：fork 是唯一入口。

### B2 安全（用户文件会被其他用户渲染/使用）

- **style.css 黑名单扫描**：禁 `url(`/`@import`（防外链追踪/加载恶意资源）、禁
  `expression(`/`behavior`/`-moz-binding`、禁 `@charset` 之外的 at-rule 里夹带 import；
  白名单扫描（token 块外不允许改结构选择器？——不做到这么严，v1 依赖 demo 渲染在
  iframe sandbox + 成品只用骨架，CSS 注入的实际面很小，黑名单够用并在计划里注明）。
- **index.html demo 清洗**：渲染给非 owner 前过整文档 sanitize（复用 sanitizeSlide 思路
  扩展：剥 script/iframe/object/on* 属性/javascript: URL）。demo 只用于预览，不进成品。
- **成品生成面**：生成管线只消费 layouts.md 骨架 + 类名契约回显校验 + 写入门禁，用户
  模板的 demo/HTML 不进成品 deck——攻击面天然收窄。
- 生成走用户模板时 `deck.SelectTemplate` 接受 `ut-xxxx`：校验（owner 或 public）→
  从用户目录实例化（与内置同一条 Instantiate 路径）。

### B3 前端（长在 C 的框架上）

- `api/userTemplates.ts` + `stores/userTemplates.ts`；
- **/templates/mine**（我的模板）：卡片列表（预览图/状态徽标/公开开关/删除/复制 fork id），
  空态引导"从任一内置模板派生"；
- **/templates/:id/edit**（定制工作台）：左 = 实时预览 iframe（demo + 变体切换 + 翻页），
  右 = agent 定制对话，顶栏 = 名称编辑 + 「发布」按钮（三关门禁的进度与失败原因展示）；
- **画廊两个新区块**（/new/template）：「我的模板」（private+public 都可选用）、
  「社区模板」（published，标作者名）；
- 入口：顶栏导航加「模板」（mine 页）；TemplateGallery 卡片上加「派生我的模板」。

### B4 验收（真实 LLM E2E）

1. 账号 A：fork tech-sharing → 对话"改成暖橙色系、圆角更大、标题用衬线"→ 预览确认 →
   发布 → 三关门禁全过；
2. 账号 B：画廊社区区块看到该模板 → 直接用它生成一份 deck → 全流程过；
3. 负例：把 style.css 改出 `url(` 外链 → 发布被门禁 1 拦，PublishError 可读；
4. 账号 A：unpublish → 账号 B 的已有 deck 不受影响（快照隔离），但新 deck 选不到该模板；
5. 单测：fork/crud/publish 门禁各状态/registry 热载/权限（非 owner 读 private 404）。

commit 切分：B 后端服务+门禁（1-2 个）、agent 定制工具+提示词（1 个）、前端管理 UI（1-2 个）。

---

## 4. 工期估算（自主执行）

| 阶段 | 预计 |
|---|---|
| A1 七个移植 | 最重（每个都要适配+验证+目视），约总工时 1/3 |
| A2 三个原创 | 1/6 |
| C 前端分步+预览 | 1/3（路由重构 + 缩略图服务是两块硬骨头） |
| B 自定义模板 | 1/3~1/2（全栈 + agent 工具 + 门禁 + E2E） |

每阶段收尾做一次完整回归（单测 + validator + 受影响 demo/真机 E2E）。

## 5. 开放问题（有默认值，审查时可推翻）

1. **taste 原创选型**：默认做 brutalist-bold / minimal-quiet / soft-pastel 三个（取自
   taste-skill 里规范最成体系的三个风格）；只要两个就去掉 soft-pastel。
2. **presenter-mode-reveal 更名 presenter-cards**：避免与已废弃的 reveal.js 混淆。
3. **缩略图生成时机**：默认"第一页生成完成时异步预生成，其余懒加载"；若你希望全部
   页在生成完成时同步出图（生成结束多等 10-20 秒），可改。
4. **社区模板的署名**：默认只展示昵称（现有 users 表有邮箱，无昵称字段——需要加
   nickname 列，默认取邮箱前缀；若不想加列就展示脱敏邮箱）。
5. **发布冒烟的额度**：默认记发布者账号（发布 = 消耗一次单页生成）；若希望免费引流
   可改为系统承担（防滥用弱化）。

## 6. 不做（本期明确排除）

- 从零构建模板（D1 已锁定 v1 只做克隆定制；从零构建留给后续版本，需模板质量验收
  管线更成熟）；
- 模板市场/评分/下载计数（本期只做"公开可用"）；
- 用户模板的结构级增删版式（layouts 只允许改约束文案，不允许加删版式——防契约烂掉；
  结构定制留给 v2 从零构建一起做）；
- guizang-ppt-skill 的任何文件搬运（AGPL，只吸收思想且不复制贡献者代码——继续遵守）。
