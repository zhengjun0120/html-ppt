# 用户模板历史回滚 + 权限放开 + 草稿收口 —— 施工方案

> 状态：**已实施并验收（2026-09-28）**。四步提交：6f4caa1（历史模块）/ 1e8744c（权限放开）/
> 0593ada（草稿收口）/ 1086731（前端）+ 实测修复。本文档是实施者的完整作业依据：
> 所有代码锚点均已对源码核实（frontend 分支 checkout，行号以撰写时为准，实施时若有
> 漂移以语义定位）。配套阅读：docs/deck-editor-plan.md（文稿手动编辑器，本方案的姊妹篇）。
>
> 实测修订（验收发现，已进代码）：
> 1. writeTemplateFile 的"重挂失败上抛"会耦合全部写路径——一个坏 demo 让之后所有
>    style.css 保存 400 且不记版本（盘与历史脱节）。改为写入/记版必落地，重挂降级为
>    警告（对话工具把校验错误转告模型自修，门禁终审兜底）。
> 2. 轮次上限退出路径的版本备注兜底（模型连发工具不调 finish 时 detail 不为空）。
> 3. scanCSS 的 url( 黑名单放行 data: 内联（brutalist fork 的纸纹底是合法用法，
>    原一律禁会让它的 fork 永远过不了门禁）；scanHTML 收紧为每个 script 标签都必须是
>    runtime.js + 要求 section/runtime 结构前提。
> 4. scanCSS/scanHTML 增加非法 UTF-8 拒绝（bundle 经 json.Marshal 会把坏字节换成
>    U+FFFD，磁盘与历史从此对不上——curl GBK 注入实测发现的防御）。

---

## 0. 需求原话与四项决策

用户需求：**给定制模板加上类似文稿的历史修改系统，可回滚；有了兜底，把定制权限放开，
让用户随便改。**

四项决策（AskUserQuestion 拍板）：

| # | 决策点 | 拍板 | 含义 |
|---|--------|------|------|
| D1 | 权限档位 | **style.css 全文可改 + index.html 结构可改** | 对话新增整文件重写工具；layouts.md/rules.md 不放开 |
| D2 | 已发布模板 | **保持锁定** | published 状态下一切写操作（含回滚）409，先下架再改 |
| D3 | 历史粒度 | **每次用户动作一版** | 一轮有改动的定制对话=1版，一次手动保存=1版，一次回滚=1版 |
| D4 | 草稿安全 | **收口：草稿 demo 走鉴权** | 只有发布后的 demo 公开可访问 |

---

## 1. 现状事实锚点（全部已核实）

### 1.1 用户模板链路

- 模型：`internal/store/model.go:56` `UserTemplate{ID(ut-+8hex), UserID, BaseID, Name,
  Description, Visibility(private/public), Status(draft/publishing/published/failed),
  PublishError, PublishReport}`。
- 服务：`internal/service/usertpl/`，共 5 文件：
  - `service.go`(430)：Fork(:67)、List/Community、getRow(:162 归属判断)、UpdateMeta(:183)、
    Delete(:201)、Unpublish(:237)、Publish(:271 两关门禁)、securityScan(:346)、
    layoutPatterns(:381)、newID(:407)。
  - `customize.go`(310)：定制对话。工具 `write_tokens`/`set_meta`/`finish`（:46 schema）、
    6 轮上限循环(:107)、execCustomTool(:155)、`writeTokenBlock`(:214 [customize] 覆盖块)、
    `applyMetaFile`(:248)、`scopeOf`(:276 读 body class 的 tpl- 前缀)、
    `customizeSystemPrompt`(:291)。
  - `save.go`(102)：`SaveIndexHTML`(:30) 手动保存 index.html，published 拒(ErrPublished:22)、
    滚动备份 5 版到 `data/user-template-backups/<ut-id>/`(:41-52)、`atomicWriteFile`(:96)。
  - `sync.go`(91)：同步内置模板目录 → user-templates（fork 外的另一来源，实施时读一遍确认是否也要接历史）。
- Handler：`internal/handler/userTemplates.go`(309)。端点：
  - `POST /api/templates/:id/fork`(:38)、GET 列表(:61)、GET /:id(:105)、
  - `GET /api/user-templates/:id/editor`(:138) —— 编辑模式 demo 读取：读 index.html →
    `href="style.css"` 重写为 `/user-templates/<id>/style.css` → `injectEditorScript` →
    `deckPageHeaders`。**仅 GetOwned**。
  - `PUT /api/user-templates/:id/file`(:169) —— 全量保存 index.html：MaxBytes 2MB
    (maxEditBodyBytes, handler/deck.go:57) → 含 `<section` 粗检 → SaveIndexHTML →
    ErrPublished→409。
  - PUT 元数据(:205)、DELETE(:228)、POST publish(:247 同步 10-20s)、unpublish(:265)、
    POST chat(:285 定制对话)。
- 路由：`internal/router/router.go:91-101`（guarded 组）；**公开静态** `:132-133`
  `utStatic.Static("/user-templates", <data>/user-templates)`（noCacheHeader 组）——
  **现状问题 A**：StaticFS 把目录里**所有文件**都伺服出去（template.json/layouts.md/rules.md
  全公开）；**现状问题 B**：draft 模板的 demo 任何知道 id 的人可访问。
- 数据目录：`data/user-templates/<ut-id>/`（template.json/index.html/style.css/layouts.md/
  rules.md/ADAPTATION.md）+ `data/user-template-backups/<ut-id>/`（滚动备份）。

### 1.2 文稿历史系统（镜像蓝本）

- `internal/service/deck/history.go`(191)：
  - `index.json` = `{next_seq, versions: [VersionMeta]}`，**versions 新→旧**（:64 prepend）。
  - `VersionMeta{Version(v%06d), Time(unix秒), Operation, Detail, Slides}`。
  - 操作常量 OpRun/OpRestore/OpEdit；版本号正则 `^v\d{6}$`。
  - `pruneVersions`(:88)：近 7 天全留；更早的每天只留当天最新一条；总数上限 200 截尾。
  - DeleteVersion/ClearHistory（NextSeq 不清零，编号永不复用）。
- `internal/service/deck/v2_history.go`(155)：
  - bundle = 单 JSON 文件 `v000001.json`（`snapshotV2{index_html, style_css, outline_json}`），
    "记档/恢复都是单文件原子操作，不存在拷了一半"。
  - `recordVersionV2`(:29)：**锁内调用**；读现文件→marshal→原子写→prepend meta→prune→
    写索引；**历史失败只 log 不阻塞主流程**（:33, :60, :73 三处兜底）。
  - `RestoreVersionV2`(:90)：authorize→版本号合法→锁→索引里存在→读 bundle→**整包写回**→
    invalidateThumbs→`recordVersionV2(OpRestore, "恢复到 "+version)`。
- 路由：`router.go:83-84,102-103`：GET history、POST history/:version/restore、
  DELETE history/:version、DELETE history（清空）。
- 前端：`components/preview/HistoryDrawer.vue`（props `{open, deckId}`，内部 deckApi.history/
  restore/deleteVersion/clearHistory，emit close/restored）；`api/decks.ts:21-27`。

### 1.3 相关设施

- vision 一次性授权：`internal/vision/grant.go:25 Issue(uid, deckID) → nonce`、`:47 Take(nonce)
  → (uid, deckID)`；deck 渲染用 `GET /api/render/:nonce`（handler/render.go）。
- registry：`MountUser(dir)`(:328 重挂读 style.css 快照，write_tokens 改完调它立即生效)、
  `UnmountUser(id)`(:372)、`ValidateUserDir(dir)`(:350)、`BuiltinDir(id)`(:382)。
- 鉴权 query token：deck 文件/缩略图 iframe 均用 `?token=`（auth 中间件已支持；受控端点复用
  同一解析入口，实施时以 `internal/auth` 中间件实现为准）。
- deck 生成侧 ut- 归属校验：`internal/service/deck/v2.go:398`（registry 挂载 ≠ 人人可用）。

### 1.4 实施中顺手修掉的现状问题（本方案范围内）

- **问题 C**：`writeTokenBlock`/`applyMetaFile` 对 published 模板**不设防**——已发布模板的
  style.css 会被对话直接改写并 MountUser 立即生效给所有用户（绕过门禁）。本方案 D2 统一封死。
- **问题 D**：Fork 失败路径（MountUser 失败 :124-127）删了 DB 行但**残留磁盘目录**（孤儿目录）。
- **问题 E**：定制对话会话（customizeSession）在内存 map，重启丢失——维持现状（不是重要
  数据），但历史系统上线后"对话改了什么"有版本可查，会话丢失的伤害变小。

---

## 2. 总体架构

```
                         ┌─ 对话定制 (POST /chat) ── write_tokens / write_style / write_demo / set_meta
                         │        ↓ dirty?（有实际文件写入）
 用户动作 ──→ 写入路径 ──┼─ 手动编辑 (PUT /file index.html、PUT /style) ── 安全预检 ──→ 落盘
                         │        ↓                                                    ↓
                         └─ 回滚 (POST /history/:version/restore) ←── 记录版本 OpRestore
                                  ↓
                    data/user-templates/<ut-id>/history/
                    ├── index.json  {next_seq, versions: 新→旧}
                    └── v000001.json  ← 五件套整包 bundle + 元信息

 公开面：/user-templates/<id>/{index.html,style.css} ← 受控端点
         published → 无条件伺服（画廊/社区 iframe 免鉴权）
         draft/failed/publishing → owner（header 或 ?token=）才伺服，其余 404
         其余文件（template.json/layouts.md/rules.md/ADAPTATION.md/history/）永不伺服
 发布门禁渲染：/api/user-template-render/:nonce（一次性授权，Take 后即焚）
```

---

## 3. 详细设计

### 3.1 历史模块 `internal/service/usertpl/history.go`（新建，约 200 行）

镜像 `deck/history.go + v2_history.go`，**不动 deck 包**（不为复用去泛化已验证代码）。

```go
const (
    maxVersionsPerUT = 100   // 模板改动频率低于文稿，上限减半
    keepFullDays     = 7     // 与 deck 同：近 7 天全留，更早每天留最后一条
)
const (
    OpFork    = "fork"    // 克隆初始化
    OpChat    = "chat"    // 定制对话（一轮有实际改动=一版）
    OpEdit    = "edit"    // 手动保存（index.html 或 style.css）
    OpMeta    = "meta"    // 改名/描述（含 UpdateMeta 同步 template.json）
    OpRestore = "restore" // 回滚本身也记版
)

type snapshotUT struct {           // 五件套整包，单 JSON 原子读写
    TemplateJSON string `json:"template_json,omitempty"`
    IndexHTML    string `json:"index_html,omitempty"`
    StyleCSS     string `json:"style_css,omitempty"`
    LayoutsMD    string `json:"layouts_md,omitempty"`
    RulesMD      string `json:"rules_md,omitempty"`
}

type UTVersionMeta struct {
    Version   string   `json:"version"`    // v000001 起单调递增永不复用
    Time      int64    `json:"time"`
    Operation string   `json:"operation"`
    Detail    string   `json:"detail"`     // 人读摘要（对话版=finish 汇报）
    Changed   []string `json:"changed,omitempty"` // 与上一版相比变化的文件名（记录时算好）
}
// 索引 historyIndex / readHistoryIndex / writeHistoryIndex：与 deck 同构（index.json）
```

核心函数（均为 Service 方法）：

- `(s *Service) historyDir(id) string` → `<root>/<ut-id>/history`
- `(s *Service) lockUT(id string) func()` —— **新增 per-template 互斥锁**。
  现状 customize 有 session 锁但手动保存与对话可并发写同一模板 → 需要 `sync.Mutex` map
  （照 deck `lockDeck` 的形态；map + 全局锁保护 map 本身）。
- `(s *Service) recordVersionUT(id, operation, detail string) error` —— **锁内调用**。
  读五件套（允许缺失→空串）→ 若与上一版 bundle 完全一致则**跳过**（幂等保护：同一动作
  重放不产生噪音版本）→ 计算 Changed（与上一版逐文件比对）→ atomicWrite `v%06d.json` →
  prepend meta → prune（deck pruneVersions 逻辑复制，参数化上限）→ 写索引。
  失败 log 不阻塞主流程（deck 同纪律）。
- `(s *Service) ListUTVersions(userID uint, id string) ([]UTVersionMeta, error)`
  —— GetOwned → 读索引。
- `(s *Service) RestoreUTVersion(userID uint, id, version string) error`：
  GetOwned → **published 拒（复用 ErrPublished，handler 映射 409）** → versionPattern →
  lockUT → 索引中存在 → 读 bundle（损坏报可读错误）→ 整包写回五件套（空串文件跳过）→
  从 bundle.TemplateJSON 解出 name/description **同步回 DB**（以快照为准）→
  `reg.MountUser(dir)`（失败则返回错误——文件已写回，用户可再回滚一版；文档化）→
  `recordVersionUT(OpRestore, "恢复到 "+version)`。
- `(s *Service) DeleteUTVersion(userID uint, id, version string) error` /
  `(s *Service) ClearUTHistory(userID uint, id string) (int, error)` —— 照 deck（锁内、
  NextSeq 不清零）。
- `readUTFile(id, name)`：读模板目录内文件的小 helper。

### 3.2 四条写入路径接历史（D3：每次用户动作一版）

1. **Fork**（service.go:67）：现有流程成功收尾处（MountUser 成功后）追加
   `recordVersionUT(id, OpFork, "克隆自 builtin:"+baseID)`。
   顺手修**问题 D**：MountUser 失败分支补 `os.RemoveAll(dir)`。
2. **定制对话**（customize.go Customize:87）：
   - 循环内 execCustomTool 改为返回 `(result string, fileDirty bool)`；
     write_tokens/write_style/write_demo/set_meta(有实际写入) → true。
   - 循环退出（finish 或无工具直答或轮次上限）后：`if dirty { recordVersionUT(id, OpChat, note) }`，
     note = finish 的 reply（截 200 字）；无 finish 则用最后一条工具结果摘要。
   - 记版在会话锁外、lockUT 内（recordUT 自带锁 → 直接调，注意别嵌套：Customize 的
     sess.mu 是会话锁，与 lockUT 不同粒度，无死锁风险；保持先 sess.mu 后 lockUT 的顺序
     全程一致即可）。
3. **手动保存**：SaveIndexHTML 成功后 `recordVersionUT(id, OpEdit, "手动编辑 index.html")`；
   新增 SaveStyleCSS 同理（见 3.3）。**滚动备份退役**：历史系统取代其职责，删除
   SaveIndexHTML 内的备份调用与 keepTemplateBackups/nextBackupName/pruneTemplateBackups
   （保留 `data/user-template-backups/` 目录不清理，历史文件不受影响）。save_test.go 相应改。
4. **元数据**：UpdateMeta（handler PUT）改为同时更新 DB + template.json（修漂移：现状
   set_meta 改两处、UpdateMeta 只改 DB），成功后 `recordVersionUT(id, OpMeta,
   "改名/描述")`；applyMetaFile 保持改文件 + dirty 标志并入当轮 OpChat。

### 3.3 权限放开（D1）

**对话新工具**（customize.go customizeTools 追加两个 schema）：

- `write_style`：`{css: string}` —— style.css 整文件重写。schema description 明确告知模型：
  文件是全文（含 token 声明区与既有规则），重写必须基于当前全文改造（system prompt 注入
  当前 style.css 全文，截断到 8000 字防 prompt 爆炸）；`[customize]` 覆盖块概念随之消失。
- `write_demo`：`{html: string}` —— index.html 整文件重写。description 要求保留
  runtime.js 引用与 body 的 tpl- 作用域类。

**写入前安全预检**（securityScan 从"发布才查"提前到"每次写入"）：

- securityScan(:346) 重构为两个纯函数 `scanCSS(css string) error` /
  `scanHTML(html string) error`（现黑名单规则原样搬入），原 securityScan(dir) 组合调用它们
  （Publish 门禁行为不变）。
- write_style/write_demo 及手动保存（SaveStyleCSS、SaveIndexHTML）落盘前先预检，失败返回
  可读错误（对话工具返回错误文本让模型转告+自行修正；HTTP 端点 400）。
- write_demo 额外校验：必须含 `<section`（demo 是页面集）**且**必须引用
  `/assets/deck-v2/runtime.js`（丢了 demo 不可翻页）。
- 内容上限：工具参数内 css/html 字段 ≤ 512KB；HTTP body 沿用 2MB。
- 预检 ≠ 门禁终审：Publish 的全量扫描（含 demo 渲染量测）原样保留。

**手动编辑 style.css**：新端点 `PUT /api/user-templates/:id/style`（raw CSS body，
MaxBytesReader 512KB）→ `SaveStyleCSS(userID, id, css)`：GetOwned → published 拒 →
scanCSS → atomicWrite → `reg.MountUser`（预览/已挂载立即生效）→ `recordVersionUT(OpEdit,
"手动编辑 style.css")`。

**write_tokens 保留**：轻量路径继续可用；published 时同样拒绝（见 3.4）。

### 3.4 已发布锁定（D2）—— 封死问题 C

统一 guard（usertpl 包内）：

```go
func (s *Service) ensureWritable(row *store.UserTemplate) error {
    if row.Status == "published" { return ErrPublished }  // 409
}
```

接入点：execCustomTool 的全部写工具（write_tokens/write_style/write_demo/set_meta ——
**修问题 C**）、SaveIndexHTML（已有）、SaveStyleCSS（新）、RestoreUTVersion（新）、
writeTokenBlock/applyMetaFile 底层再防御一道。发布门禁运行态（publishing）允许读不允许写，
同 published 处理（写入类也 409，提示"门禁运行中"）。

### 3.5 草稿收口（D4）—— 受控 demo 端点替换公开静态

**删掉** router.go:132-133 的 `utStatic.Static("/user-templates", ...)`，换成受控 handler
（**不挂 Auth 中间件**——published 必须免登录可访问，鉴权在 handler 内部按状态决定）：

```go
// GET /user-templates/:id/*filepath（白名单伺服）
//   filepath ∈ {index.html, style.css}，其余一律 404（修现状问题：源文件全公开）
//   row := db 查；不存在 → 404
//   status == published            → 伺服（无鉴权；社区画廊 iframe 免登录）
//   其他状态（draft/failed/publishing）→ 需要 owner：
//       Authorization header 或 ?token= 解出 uid（复用 auth 包的解析入口，
//       与 deck ?token= 同一机制），GetOwned 通过 → 伺服；否则 404（不泄露存在性）
//   index.html 响应带 deckPageHeaders（CSP 同 deck）；style.css 带 no-cache
```

- 相对引用自洽：demo 里 `href="style.css"` 在 `/user-templates/<id>/index.html` 路径下解析
  为 `/user-templates/<id>/style.css`，同进本 handler；`/assets/...` 绝对引用照走公开静态。
- **发布门禁渲染改 nonce**：Publish(:298) 的渲染 URL 从公开静态改为
  `GET /api/user-template-render/:nonce`（新 handler：grants.Issue(uid, "ut:"+id) →
  Take 后校验 GetOwned → 读 index.html → style.css 引用重写为
  `/api/user-templates/<id>/assets/style.css?token=…`（headless 导航无鉴权头，走 query
  token；该受控资产端点按 owner+token 放行）→ deckPageHeaders）。grants.Take 即焚，渲染
  结束 URL 作废。复用 `vision.Grants`（Issue/Take 的 deckID 参数就是 string，传 "ut:"+id）。
- **前端预览地址切换**：`userTemplatePreviewUrl`（api/userTemplates.ts:77）改为返回鉴权
  demo 端点 `GET /api/user-templates/:id/demo?page=N`（新 handler：GetOwned/GetReadable →
  读 index.html → 相对引用重写为受控资产端点（带 token 由前端 authedUrl 追加）→
  deckPageHeaders）。调用点：TemplateEditView.vue（工作台预览 iframe）、MyTemplatesView.vue
  （编辑弹窗 src，如走 userTemplatePreviewUrl）、HistoryDrawer 类组件的版本预览（新）。
  **注意 token 过期**：iframe URL 的 token 是加载时签的，长会话可能过期——deck 编辑弹窗已有
  同样问题与刷新策略（previewKey 重挂），照抄。

### 3.6 前端

- `api/userTemplates.ts`：+ `history(id)` / `restoreVersion(id, v)` / `deleteVersion(id, v)` /
  `clearHistory(id)` / `saveStyle(id, css)` / 类型 `UTVersionMeta`；`userTemplatePreviewUrl`
  改受控端点。
- `components/usertpl/TemplateHistoryDrawer.vue`（新建，镜像 HistoryDrawer.vue 的骨架）：
  props `{open, templateId}`；列表项 = 序号/时间/来源徽章(fork/chat/edit/meta/restore)/
  备注摘要/Changed 文件 chips；操作 = 回滚（confirm，published 置灰+提示"先下架"）、删除单版、
  清空（双 confirm）。
- `TemplateEditView.vue`：顶栏加「历史」按钮挂抽屉；加「样式」编辑面板（textarea 展示当前
  style.css → 保存调 saveStyle → 成功后 previewKey+1 刷新预览；dirty 离开提示）。预览 iframe
  src 切受控端点。
- `MyTemplatesView.vue`：状态徽章逻辑不动；若卡片预览走公开静态，同步切换。

---

## 4. API 总览（Δ）

| 端点 | 方法 | 变化 |
|---|---|---|
| /api/user-templates/:id/history | GET | 新：版本列表 |
| /api/user-templates/:id/history/:version/restore | POST | 新：回滚（published→409） |
| /api/user-templates/:id/history/:version | DELETE | 新：删单版 |
| /api/user-templates/:id/history | DELETE | 新：清空 |
| /api/user-templates/:id/style | PUT | 新：手动保存 style.css（raw，预检+记版） |
| /api/user-templates/:id/file | PUT | 行为不变 + 记版 + 预检 |
| /api/user-templates/:id/chat | POST | 工具集扩展 + published 锁 + 记版 |
| /api/user-templates/:id/demo?page=N | GET | 新：受控鉴权 demo（工作台/编辑弹窗预览） |
| /api/user-template-render/:nonce | GET | 新：门禁渲染专用（一次性） |
| /user-templates/:id/*filepath | GET | **行为变更**：StaticFS → 受控白名单端点 |

---

## 5. 安全分析

- **威胁模型**：demo 会被非属主渲染（published 公开、draft 原先也公开）。收口后 draft 仅
  owner 可达，published 必须过门禁扫描——渲染面永远只暴露"当前状态通过过扫描的内容 +
  owner 自己的草稿"。
- 每次写入预检（scanCSS/scanHTML）把"发布才发现注入"提前到"写入即拒"；门禁终审保留。
- history/ 目录在模板目录内但**不在白名单伺服清单**——不会被端点泄露；`data/` 本身不进 git。
- token 级防注入（键 -- 前缀、值禁 ;{}）保留；整文件重写的注入面由 scanCSS 黑名单覆盖。
- 回滚不能用于越权：RestoreUTVersion 走 GetOwned；published 锁住"借回滚绕门禁改公开内容"。

## 6. 测试计划

**go 单测**（`usertpl/history_test.go` + 扩展 `customize_test.go`/`save_test.go`，测试基建
照现有 save_test.go 的装配模式）：

1. Fork 记 v1（五件套内容正确、OpFork）；Fork 失败清理目录（问题 D 回归）。
2. write_tokens 有改动记 OpChat（changed=[style.css]，detail=finish reply）；无改动不记。
3. write_style 整文件生效 + 记版；scanCSS 拦截（@import 样例）且**不落盘不记版**。
4. write_demo：缺 runtime.js 拒；含第二 script 拒；合法通过 + 记版。
5. published：四个写工具、SaveIndexHTML、SaveStyleCSS、RestoreUTVersion 全 409。
6. 回滚：写回五件套、DB name/description 同步、记 OpRestore、changed 计算；
   bundle 损坏报错；版本不在索引报错。
7. prune：>100 截尾；7 天前每天留末条；v1 永不被裁（第一条永久保留语义 = deck 一致性，
   实现时确认 prune 对最老一条的处理，必要时显式保留 v1）。
8. 并发：两个 goroutine 同时 SaveStyleCSS 与 write_tokens → 文件无交错（锁生效）。
9. UpdateMeta 同步 template.json + 记 OpMeta。
10. 受控端点：published 无 token 200；draft 无 token 404；draft 带 owner token 200；
    rules.md/template.json 404；history/ 404。

**浏览器验收清单**（5173 真实环境，不保存用户文稿）：

1. fork 一个内置模板 → 工作台历史面板出现 v1「克隆自 builtin:xxx」。
2. 对话「主色换成橙色」→ v2 出现、备注=agent 汇报、预览即时变色。
3. 对话「把封面大标题改成两行并加一个装饰条」→ write_demo 生效、demo 结构变化、记 v3。
4. 手动「样式」面板改一段 CSS 保存 → v4、预览刷新。
5. 回滚到 v2 → 文件与预览回到 v2、历史顶部出现「恢复到 v000002」v5。
6. 发布 → 门禁照常两关；发布后对话改色被拒（提示先下架）；回滚也被拒。
7. 下架 → 再改/再回滚放行；重新发布通过。
8. 无痕窗口访问 draft 模板 demo URL → 404；published 模板 demo URL → 200。
9. /user-templates/<id>/rules.md → 404（不再泄露源文件）。
10. 编辑弹窗（index.html 可视编辑）保存 → 记版、预览一致。

## 7. 实施顺序（每步一个提交，全绿再下一步）

1. **历史模块 + 四路径记版**（3.1/3.2 + 单测 1/2/5/7/9）——不含权限放开，行为纯增量。
2. **权限放开 + 全路径锁定**（3.3/3.4 + 单测 3/4/5）。
3. **草稿收口 + 受控端点 + nonce 渲染**（3.5 + 单测 10）。
4. **前端**（3.6 + 验收清单全量跑）。
每步结束：go build/vet/test 全绿 + 双写 backend 分支 checkout（editor 纪律）+ 8080 重建重启。

## 8. 不做的事（边界）

- layouts.md/rules.md 的对话修改（D1 未选；要改走"未来手动文件编辑"课题）。
- 历史 diff 可视化（只列 changed 文件名，不做行级 diff 视图）。
- 定制对话会话持久化（内存态维持现状；历史已兜底"改了什么"的可追溯）。
- 模板级乐观锁/协作编辑（单 owner 假设不变）。
- deck 历史 prune/索引机制的泛化重构（deck 包零改动）。

## 9. 文件改动总览

| 文件 | 动作 |
|---|---|
| backend/internal/service/usertpl/history.go | 新建（~200 行：bundle/索引/记版/恢复/删除/prune/锁） |
| backend/internal/service/usertpl/history_test.go | 新建 |
| backend/internal/service/usertpl/customize.go | 工具 schema +2、execCustomTool dirty/published guard、system prompt 注入 style.css 全文 |
| backend/internal/service/usertpl/service.go | Fork 收尾记版+失败清理；securityScan 拆纯函数；Publish 渲染 URL→nonce |
| backend/internal/service/usertpl/save.go | 记版接入；滚动备份退役 |
| backend/internal/service/usertpl/scan.go | 新建（scanCSS/scanHTML 纯函数，从 service.go 迁出） |
| backend/internal/handler/userTemplates.go | +4 端点（history×3→复用路由、style、demo）；file 端点接预检 |
| backend/internal/handler/usertpl_render.go | 新建（nonce 渲染 + 受控资产） |
| backend/internal/router/router.go | +5 路由；utStatic.Static 删除 |
| frontend/src/api/userTemplates.ts | +6 方法/类型；preview URL 切换 |
| frontend/src/components/usertpl/TemplateHistoryDrawer.vue | 新建 |
| frontend/src/views/TemplateEditView.vue | 历史按钮+抽屉、样式面板、预览地址 |
| frontend/src/views/MyTemplatesView.vue | 预览地址切换（如涉及） |
| docs/user-template-history-plan.md | 本文档（实施后更新状态行） |
