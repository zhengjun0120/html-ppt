# deck 手动编辑器（PPT 式手改）施工方案

> 状态：**已实现并验收**（2026-09-26，frontend 分支提交 fa710c7 → 8e519af；E2E 十三条验收全过，见 §7）。
> 本文档是施工蓝图：所有现状事实都已在仓库里核实过并标注出处；实现时若发现与文内事实不符，先回来改文档再动手。
> 仓库约定：本 checkout（D:\go_files\html-ppt-frontend，frontend 分支）只跑前端 5173 + 测试后端 8081；backend 分支 checkout 在 D:\go_files\html-ppt（后端 8080，用户正式环境）。**前端改动在 frontend 分支提交，后端改动既改本 checkout 的 backend/ 也同步到 backend 分支 checkout**（历史上后端文件两边同内容，合并线一致）。

---

## 0. 一句话

在预览弹窗里对 deck-v2 文稿和用户模板做 PPT 式手改：双击改文字、拖动/缩放卡片（绝对定位 + 智能参考线 + 原位占位块），显式保存——文稿写进现成 history 版本体系，用户模板覆盖保存 + 服务端留最近 5 版备份。

## 1. 已确认的产品决策（与用户两轮问答定稿，勿再反复）

| # | 决策点 | 结论 |
|---|--------|------|
| 1 | 首版能力 | **文本自由编辑 + 拖动/缩放卡片** 两项。样式面板、元素增删复制、页级操作全部二期。（2026-09-26 晚追加：**字号缩放**提前做了——工具栏 A−/A+ 步进器 + Ctrl+=/Ctrl+-，选中块及其内部文字等比缩放，px 行高随动） |
| 2 | 编辑器形态 | **预览弹窗内编辑**（不做独立全屏路由） |
| 3 | 保存机制 | **显式保存**（Ctrl+S / 按钮）+ 文稿写 history 版本可回滚 |
| 4 | 可编辑对象 | **用户文稿（deck）+ 我的模板（ut-*）**；内置模板画廊 demo 只读 |
| 5 | 拖动语义 | 拖走后**原位置留等大占位块**，其余元素不动（PPT 心理模型）；占位块是真实内容，随保存持久化，用户可后续手动删 |
| 6 | 可选中元素 | **两层**：`.slide` 直接子元素 + grid/flex 容器内的卡片；点击选中最内层，Esc 逐级上选 |
| 7 | AI 与手改 | **覆盖 + 历史可回滚**：AI 重新生成/对话改稿会盖掉手改，不锁页；每次 AI 写入和手动保存都在 history 里 |
| 8 | 模板备份 | 用户模板保存时服务端**自动留最近 5 版**（滚动） |
| — | 临时默认（用户可验收前改） | 占位块带 `data-ed-placeholder` 标记；编辑入口按钮在弹窗右上角 / 预览工具栏 |

## 2. 现状事实（全部已在仓库核实，附锚点）

> 锚点行号是 2026-09-26 的值，会漂移；以函数名为准。

### 2.1 deck（文稿）侧

- 一份 deck = `backend/data/decks/<deck-id>/` 目录：`index.html`（全部页面内联，实测 ~32KB/18 页）、`style.css`（模板 CSS 的 deck 私有拷贝）、`deck.json`（元数据）、`outline.json`、`history/`、`thumbs/`。
- `deck.json` 形如：`{id, format:"v2", stage:"iterating", title, template_id, variant, canvas:{w:1920,h:1080}, page_plan:[{no,layout,reason}...], created_at, updated_at}`。**画布固定 1920×1080**（`deck.json` canvas 字段；base.css DESIGN CANVAS 段确认 `.deck{width:var(--deck-w,1920px);height:var(--deck-h,1080px);transform:scale(var(--deck-scale,1))}`，runtime.js 按窗口写 `--deck-scale`）。
- index.html 每页一个 `<section class="slide" data-layout="..." data-id="sN">`，页内普通流式布局（flex/grid）；`.slide` 是 `position:absolute;inset:0; display:flex;flex-direction:column;justify-content:center; padding:72px 96px`（base.css `.slide` 规则）。**绝对定位子元素以 .slide 的 padding box 为包含块，边框为 0，故 getBoundingClientRect 差值坐标可直接用**。
- 读链路：`deck.Service.GetHTML(userID,id)`（service/deck/deck.go:93，authorize 后读 index.html）→ `PreviewHTML`（service/deck/v2.go:498，v2 时把 style.css **内联**进 `<head>`，即预览 iframe 拿到的 HTML 自带全部样式）→ `handler.GetDeckFile`（handler/deck.go:31，`GET /api/decks/:id/file`，受 JWT 保护，`?token=` 回退兼容 iframe/img）。
- **安全响应头（deckPageHeaders，handler/deck.go:48）**：`CSP: script-src 'self'; style-src 'self' 'unsafe-inline'; connect-src 'none'; form-action 'none'; frame-ancestors 'self'` 等。**两条铁律由此而来：①编辑器脚本必须同源自托管（/assets/...）；②iframe 内绝不能发 fetch——保存数据只能 postMessage 出来由父页调 API。**
- **history 设施（service/deck/history.go + v2_history.go）**：
  - 快照是 bundle JSON：`snapshotV2{index_html, style_css, outline_json}`，一版一个 `history/v%06d.json`；索引 `history/index.json`：`{next_seq, versions:[VersionMeta{version,time,operation,detail,slides}]}`。
  - 操作常量：`OpRun="run"`、`OpRestore="restore"`（history.go 顶部）→ **手改新增 `OpEdit="edit"`**。
  - 写入入口：`recordVersionV2(deckID, operation, detail)`（锁内调用，先读当前 index.html/style.css/outline.json 打包，写 `<version>.json`，NextSeq++ 单调不复用，`pruneVersions` 分层保留：近 7 天全留 + 更早每天一条 + 总量 200）。
  - 已有恢复链路：`RestoreVersionV2`（bundle 整包写回 + 记一条 restore + `invalidateThumbs`）。**前端 api/decks.ts 的 deckApi.history/restore 已存在**（`VersionMeta.operation` 类型目前是 `'run'|'restore'`，要加 `'edit'`）。
  - deck 级互斥锁 `lockDeck(deckID)`（service/deck/deck.go）包住「写文件 + 记版本」。
  - `atomicWriteFile(path,data)` 在 service/deck/create.go:60。
- 缩略图：`/api/decks/:id/thumbs(/:no)` 按内容版本缓存，写 index.html 后**必须 `s.invalidateThumbs(deckID)`**（RestoreVersionV2 里就是这么做的）。
- v1 老格式（data/decks/deck-0001 那种 reveal.js deck.html）：`IsV2(id)` 以 deck.json 存在为准（deck.go:115）。**v1 不支持编辑**。

### 2.2 用户模板（ut）侧

- 存储：`backend/data/user-templates/<ut-id>/`，标准四件套 `template.json / index.html / style.css / layouts.md / rules.md`（usertpl/service.go forkCopies）。fork = 从内置模板目录拷四件套 + 改 template.json 的 id/name/source + DB 建行（`store.UserTemplate{ID,UserID,BaseID,Name,Visibility,Status:...}`）。
- 服务方法：`usertpl.Service`（usertpl/service.go）——`Dir(id)`、`GetOwned(userID,id)`（:175，**保存用这个，仅属主**）、`GetReadable`（属主或已发布，读取用）、`Fork/List/UpdateMeta/Delete/Publish/Unpublish/Community`。
- 发布状态机：`draft → publishing → published / failed`。**发布门禁（Publish）含无头渲染量测 + securityScan(dir) 目录扫描**。
- demo 预览走公开静态路由 `/user-templates/<id>/index.html`（router.go utStatic 组，noCacheHeader），**没有服务端注入钩子** → 编辑模式必须新开受保护读取端点（见 4.2）。
- 定制对话（customize.go）会改 style.css（不改 index.html 结构）——与手动编辑 index.html 并行，last-write-wins。

### 2.3 runtime.js 预览协议与沙箱

- URL 带 `?preview=N`（1-based）→ runtime 进预览模式：`documentElement/body` 打 `data-preview="1"`，只显示第 N 页（`.slide.is-active` / `.is-prev` class 切换，带 .5s 过渡），监听 `{type:'preview-goto', idx}`（0-based）无刷新翻页，就绪后回报 `{type:'preview-ready'}`（runtime.js:40,162-173,188-209）。
- runtime 会改的 DOM 状态（**serialize 清单的依据**）：`.slide` 的 `is-active/is-prev` class、html/body 的 `data-preview` 属性、`.deck` 的 inline style `--deck-scale/--deck-w/--deck-h`（:58-63）、`syncAmbient` 写的 `--ambient-*`、`.slide-number` 的 `data-current/data-total`（:337，生成时 HTML 里本就有这两个属性，值被 runtime 更新，无害）、预览分支给页内 `.notes`/`.speaker-notes` 打的 inline `display:none`（:168-170，hideSel；base.css:201 本就有 `.notes{display:none!important}`，inline 是冗余的，serialize 剥掉即还原原始形状）。
- **预览分支以 `return;` 结束（runtime.js:210）**：键盘监听（:999 presenter、:1099 正常模式）、进度条、notes overlay、overview 全部只在正常模式路径注册——预览模式零键盘监听，与编辑器键盘体系无冲突（已核实）。
- runtime.js 里 830-873 行的 drag/resize 是**演示者视图浮窗（.pcard）**的，与内容编辑无关；可参考其 mousedown/mousemove/mouseup 模式。
- 模板预览 iframe 沙箱（frontend）：内置模板 `allow-scripts allow-same-origin`，用户模板 `allow-scripts`（MyTemplatesView.vue:394-400、TemplateGallery）。deck 预览 iframe（DeckPreview.vue）`sandbox="allow-scripts allow-popups"`。**编辑弹窗统一用 `sandbox="allow-scripts"`（opaque origin）**，postMessage 跨沙箱可用（preview-goto 已验证），不放开同源、不碰安全模型。
- 编辑模式下 runtime 跑**预览模式**（src 同时带 `preview=1&edit=1`）：预览模式没有键盘监听（翻页全靠 postMessage），与编辑器的键盘体系零冲突——这是选预览模式做宿主的决定性理由。

### 2.4 前端挂载点

- `frontend/src/components/templates/TemplatePreviewModal.vue`：只读预览弹窗。单 iframe `:key="src"`、`pointer-events-none`（键盘留在父页）、preview-ready 握手 + preview-goto 翻页 + 就绪前翻页补发。**编辑弹窗不复用它（交互需求相反：iframe 必须可点），但结构参照它抄**。
- `frontend/src/views/DeckView.vue`（工作台）挂 `DeckPreview.vue`（components/preview/DeckPreview.vue）：工具条（导出/历史/刷新/新标签）+ 缩略图栏 + iframe（src=`authedUrl('/api/decks/:id/file')` + `#/N` 整帧重载式翻页）。**文稿编辑入口 = 这个工具条加「编辑」按钮，弹窗挂 DeckView**。
- `frontend/src/views/MyTemplatesView.vue`：卡片预览 chip 已存在（bottom-2 right-2）；**用户模板卡加「编辑」chip**，直接开编辑弹窗。
- `frontend/src/api/client.ts`：`request()` 自动带 `Authorization: Bearer`；`authedUrl(path, query)` 生成 `?token=` URL（iframe/img 场景）。
- 历史回滚 UI 已存在（DeckView 的 historyOpen 弹窗），operation=edit 的版本会自动出现在列表里，无需改 UI，只改类型。

## 3. 总体架构

```
┌─ 父页（Vue，5173 同源代理）────────────────────────────┐
│ DeckEditModal.vue（新）                                  │
│  工具栏：标题·脏标记·撤销/重做·保存(Ctrl+S)·翻页·关闭    │
│  职责：只有「外壳 + 保存流」。不发键盘命令给画布。        │
└────────────┬────────────────────────────▲──────────────┘
   postMessage│（扩展 preview 协议，见 §5）  │editor-* 消息
┌─────────────▼────────────────────────────┴────────────┐
│ iframe（sandbox=allow-scripts，opaque origin）          │
│  src: deck → /api/decks/:id/file?token=..&preview=1&edit=1
│        ut   → /api/user-templates/:id/editor?token=..&preview=1&edit=1
│  runtime.js：预览模式（单页显示、preview-goto 翻页）     │
│  editor.js（新，服务端 ?edit=1 时注入）：                │
│    选中/拖动/缩放/智能参考线/占位块/文本编辑/撤销栈      │
│    serialize（清理编辑器痕迹）→ 全量 HTML 发回父页       │
└────────────┬───────────────────────────────────────────┘
             │ 父页 PUT 全量 HTML
┌────────────▼───────────────────────────────────────────┐
│ 后端（Go）：                                             │
│  PUT /api/decks/:id/file          → index.html + history |
│  GET /api/user-templates/:id/editor（注入版读端点）       │
│  PUT /api/user-templates/:id/file → 覆盖 + 滚动备份 5 版  │
└─────────────────────────────────────────────────────────┘
```

要点：
- **编辑器全部跑在 iframe 内**。父页不解析、不修改 HTML，只搬运字节——沙箱模型（用户模板 scripts-only）零改动。
- **保存 = iframe 内把当前 DOM 清理成"干净 HTML"后整包发出**，父页原样 PUT。单文件 30KB 量级，全量无压力；不做 op/diff 协议（二期也大概率不需要）。
- 翻页复用 `preview-goto`（runtime 处理），编辑器不参与翻页。

## 4. 详细设计

### 4.1 后端 A：deck 全量保存（挂 history）

**路由**（backend/internal/router/router.go，guarded 组，放在 `GetDeckFile` 旁）：

```go
guarded.PUT("/decks/:id/file", h.SaveDeckFile)
```

**Handler**（backend/internal/handler/deck.go，紧跟 GetDeckFile）：

```go
// SaveDeckFile PUT /api/decks/:id/file —— 编辑器全量保存 index.html。
// body 是原始 HTML（不是 JSON）：编辑器序列化产物原样落盘。
// 校验从宽（结构粗检 + 大小上限），因为内容本来就是用户自己的稿子。
func (h *Handler) SaveDeckFile(c *gin.Context) {
    uid, ok := authctx.UserID(c.Request.Context())
    if !ok { response.Err(c, http.StatusUnauthorized, "未登录"); return }
    body, err := io.ReadAll(c.Request.Body)  // 限流：http.MaxBytesReader 包一层，2MB
    if err != nil || len(body) == 0 { response.Err(c, http.StatusBadRequest, "空内容"); return }
    if len(body) > 2<<20 { response.Err(c, http.StatusRequestEntityTooLarge, "超过 2MB"); return }
    if !bytes.Contains(body, []byte("<section")) {
        response.Err(c, http.StatusBadRequest, "内容不含 <section>，疑似非 deck HTML"); return
    }
    if err := h.decks.SaveHTML(uid, c.Param("id"), string(body), c.Query("detail")); err != nil {
        response.Err(c, http.StatusNotFound, err.Error())  // 与 GetDeckFile 同口径：不存在/非 v2 都 404，不泄露存在性
        return
    }
    response.OK(c, gin.H{"ok": true})
}
```

**Service**（backend/internal/service/deck/，新文件 v2_edit.go）：

```go
// SaveHTML 手动编辑保存：全量覆盖 index.html + 记一条 edit 版本。
// 与 agent 写入同一条「锁内写文件 → 记版本 → 缩略图失效」链路；stage 不动。
func (s *Service) SaveHTML(userID uint, id, html, detail string) error {
    if err := s.authorize(userID, id); err != nil { return err }
    if !s.IsV2(id) { return fmt.Errorf("deck %s 不支持编辑（仅 v2 格式）", id) }
    unlock := s.lockDeck(id)
    defer unlock()
    ip, err := s.IndexPathV2(id)   // 内含 idPattern 白名单校验
    if err != nil { return err }
    if err := atomicWriteFile(ip, []byte(html)); err != nil { return err }
    s.invalidateThumbs(id)
    if detail == "" { detail = "手动编辑" }
    return s.recordVersionV2(id, OpEdit, detail)  // recordVersionV2 读刚写完的 index.html 打包
}
```

- `OpEdit = "edit"` 常量加在 history.go 的 `OpRun/OpRestore` 旁。
- **锁 + 版本的次序**：先写 index.html 再 `recordVersionV2`（快照读的就是新内容，与 RecordRunVersion 的语义一致）。
- 并发语义：last-write-wins，不做乐观锁（两标签页/与 AI run 并发时后写覆盖，双方都留了版本可回滚）——已与用户确认接受。

### 4.2 后端 B：用户模板的编辑读取 + 覆盖保存

**读端点**（静态路由没有服务端钩子，编辑态必须走受保护端点做注入）：

```go
// router.go guarded 组：
guarded.GET("/user-templates/:id/editor", h.GetUserTemplateEditor)
guarded.PUT("/user-templates/:id/file", h.SaveUserTemplateFile)
```

**GetUserTemplateEditor**（handler/userTemplates.go）：
1. `usertpl.GetOwned(uid, id)` 归属校验（不是 GetReadable——协作编辑不该存在，仅属主）。
2. 读 `Dir(id)/index.html`。
3. 相对引用重写：`href="style.css"` → `href="/user-templates/<id>/style.css"`（照抄 templates.go:49 `PreviewTemplate` 的写法）。
4. `injectEditorScript(html)`（见 4.3）。
5. `deckPageHeaders(c)` + `Cache-Control: no-store` 输出。CSP 与 deck 预览同一条（script-src 'self' 允许 /assets/deck-v2/editor.js）。

**SaveUserTemplateFile**：
1. `GetOwned` 校验；**`row.Status == "published"` 时拒绝**（409，"先下架再编辑"）——已发布模板是社区内容，绕过发布门禁改它不可接受；草稿/failed 随便改。
2. body 校验同 4.1（非空、2MB、含 `<section`；模板 demo 一定含 section）。
3. **滚动备份**：备份目录放在模板目录**外**（避免污染 Publish 的 securityScan 目录扫描和静态服务）：`data/user-template-backups/<ut-id>/index-<unixms>.html`。写新文件前先把当前 index.html 拷过去，然后按文件名排序删最旧的只留 5 个。
4. `atomicWriteFile`（usertpl 包内自己写一个 10 行的同名 helper，或把 deck 的导出——选前者，别跨包引私有）覆盖 `index.html`。
5. 响应 `{ok:true, backups:n}`。
- **不动 template.json / style.css / layouts.md**（编辑器只改 index.html）。定制对话改 style.css 与本链路天然无冲突。
- 已知妥协：保存后**不自动触发**发布门禁重测——反正 published 状态根本进不来；draft 发布时门禁会全量重跑。

### 4.3 后端 C：编辑模式注入（deck 复用同一 helper）

handler 包内新增共享 helper：

```go
// injectEditorScript 在 </body> 前注入编辑器脚本。找不到 </body> 时整个追加
// 在尾部（fail-open：编辑器脚本 defer 执行，位置不严谨也能活）。
// 仅在 ?edit=1 时调用——正常预览、导出、量测、无头审查拿到的一律是纯净 HTML。
func injectEditorScript(html string) string {
    tag := `<script src="/assets/deck-v2/editor.js" defer></script>`
    if i := strings.LastIndex(html, "</body>"); i >= 0 {
        return html[:i] + tag + html[i:]
    }
    return html + tag
}
```

- `GetDeckFile`（handler/deck.go:31）改造：`PreviewHTML` 之后——
  ```go
  if c.Query("edit") == "1" {
      if !h.decks.IsV2(c.Param("id")) { /* 仍按 404 口径返回"deck 不存在" */ }
      html = injectEditorScript(html)
  }
  ```
  （v1 deck 不注入也不报错分支——编辑弹窗靠 editor-ready 超时兜底提示，见 4.6。）
- `editor.js` 是第一方新文件，**不适用 deck-v2 PATCHES.md 补丁纪律**（那是 vendor 文件的规矩）；但在 PATCHES.md 末尾加一行说明「editor.js 为第一方文件，改动无需登记」。
- CSP `script-src 'self'` 天然放行同源 `/assets/deck-v2/editor.js`；`connect-src 'none'` 强制编辑器内不能 fetch——协议上编辑器只有 postMessage 一条出路，正好。

### 4.4 editor.js（核心，最大块）

文件：`backend/web/assets/deck-v2/editor.js`。纯第一方、无构建步骤、无依赖，风格对齐 runtime.js（IIFE + ES2017 子集：const/let/箭头函数可以，不用 import）。

**激活**：脚本只在编辑态被注入（注入即激活，无需再查 URL），`defer` 保证 DOM 就绪。流程：

```
init():
  1. 注入 <style id="ed-styles">（全部编辑器样式，见下）
  2. 收集 slides = [...document.querySelectorAll('.deck > section.slide')]
  3. 给每个候选元素打标记（见「选中模型」）
  4. 绑定 document 级 pointerdown / dblclick / keydown / pointermove / pointerup
  5. postMessage({type:'editor-ready', pages: slides.length})
```

**编辑器样式（#ed-styles 内容要点）**：
- `.ed-selected{outline:2px solid #6366f1; outline-offset:2px; cursor:move}`（outline 不占布局，不引发回流）
- `.ed-hoverable{outline:1px dashed rgba(99,102,241,.45)}`（hover 提示）
- `.ed-editing{outline:2px dashed #6366f1; cursor:text}`（文本编辑态）
- 手柄 `.ed-handle{position:absolute;width:12px;height:12px;background:#fff;border:2px solid #6366f1;border-radius:2px;z-index:2147483647}`，八个方向各一个，定位在选中元素四角/四边（-7px 偏移），`data-ed-h="nw|n|ne|e|se|s|sw|w"`。
- 参考线 `[data-ed-guide]{position:absolute;background:#f59e0b;pointer-events:none;z-index:2147483646}`（竖线 width:1px;height:100%；横线反之）。
- `[data-ed-placeholder]{}` 不需要样式——它就是一个普通空块，占位即全部使命。
- **不用浮动工具条**（iframe 内 UI 只有 outline/手柄/参考线三类），所有命令在父页工具栏，iframe 内只留键盘。

**选中模型（两层 + Esc 上选）**：

```js
function collectCandidates(slide) {
  const sel = [];
  for (const el of slide.children) {          // 第一层：直接子元素
    if (el.classList.contains('notes')) continue;      // 讲稿永不选
    if (el.tagName === 'SCRIPT' || el.tagName === 'STYLE') continue;
    sel.push({ el, level: 1 });
    const d = getComputedStyle(el).display;
    if (d.includes('grid') || d.includes('flex')) {    // 第二层：容器内的卡片
      for (const c of el.children) {
        if (c.classList.contains('notes')) continue;
        if (isDecorative(c)) continue;                 // .dot、.slide-number 等装饰排除表
        sel.push({ el: c, level: 2 });
      }
    }
  }
  return sel;
}
```

- 排除表 `isDecorative`：`.dot`、`.slide-number`、空元素（无文本无子元素）。
- **命中**：click → `document.elementFromPoint(e.clientX, e.clientY)`（只有 is-active 的 slide `pointer-events:auto`，天然只会命中当前页）→ 向上走到最近的候选元素 → 选中；点到空白 → 取消选中。
- **Esc 链**：文本编辑态 → 退出编辑；选中 level2 → 上选其父（level1）；选中 level1 → 取消选中。
- 选中即渲染 outline + 八手柄（手柄 append 进选中元素，deselect 时移除）。

**拖动（流式 → 绝对定位 + 占位块）**：

```
pointerdown 在已选中元素上（非手柄）→ 开始拖动：
  scale = slide.getBoundingClientRect().width / slide.clientWidth   // .deck 被 --deck-scale 缩放，全部坐标换算除以它
  第一次拖动该元素时（transform 化）：
    rect = el.getBoundingClientRect()
    // 占位块：插在原 DOM 位次，顶住流式空间——其余元素纹丝不动的关键
    ph = <div data-ed-placeholder>，style: width=(rect.width/scale)px; height=(rect.height/scale)px;
         margin = 拷贝 el 的 computed margin; flex:'0 0 auto'; display:'block'
    el.before(ph)
    // 转绝对定位（包含块 = .slide 的 padding box，rect 差值坐标直接可用）
    el.style.position='absolute'
    el.style.left  = (rect.left - slide.getBoundingClientRect().left)/scale + 'px'
    el.style.top   = (rect.top  - slide.getBoundingClientRect().top )/scale + 'px'
    el.style.width = rect.width/scale + 'px'    // 锁死拖动起点尺寸
    el.style.margin='0'
    slide.appendChild(el)                        // 挪到最后一个子节点 → 后画的在上层，不用写 z-index
  pointermove：left/top = 起点 + delta/scale，过参考线吸附（见下），钳位在 [0, slide尺寸-元素尺寸]
  pointerup：清参考线；pushUndo()；reportDirty()
```

- 已绝对定位的元素再拖 = 直接改 left/top，不再生成占位块。
- pushUndo 在 pointerup 一次（拖动过程不进栈）。
- 占位块**随 serialize 保留**（是用户内容），用户后续可手动删（选中删除是二期，v1 可以拖别的卡盖上去）。

**缩放（八手柄）**：

```
pointerdown 在 .ed-handle 上 → beginResize 必须先 transformToAbsolute（与拖动同一条纪律）。
  ——实测教训：流式元素上直接写 width/height 会引发居中布局（justify-content:center）
  整页回流，元素带着选中框"跳走"，且 left/top 对流式元素不生效。
转换后按方向改 width/height（角手柄同时改 left/top），
最小 48×24（slide-local），实时写 style，pointerup 时 pushUndo + reportDirty。
**角手柄（nw/ne/se/sw）= 内容等比缩放**：拖动中以「起始字号 ×（当前宽/起始宽）」
实时换算选中块内所有文字的 font-size 与 px 行高（基准是拖动开始时的现值，避免连乘漂移），
与盒子同一条 undo；边手柄（n/s/e/w）不动字号。注意模板 CSS 的 max-width 会钳制
实际渲染宽度——拖到上限框就停，属模板护栏不是 bug。
文本自然回流（流式内容的固有行为，接受）。
```

**智能参考线（吸附）**：

```
参照集 = 当前 slide 内其他已 transform 化/候选元素的 getBoundingClientRect（slide-local 换算后）
       + 内容区四边（padding 72/96）+ slide 水平/垂直中线
阈值 6px（slide-local 单位）
拖动中：我的 left/center/right vs 参照的 left/center/right 取最近差 ≤ 阈值 → 吸附 + 画竖参考线；
        top/middle/bottom 同理横线。参考线挂 slide 上，pointerup 全部移除。
不做等距分布、不做同尺寸吸附（二期）。
```

**文本编辑**：

```
dblclick 选中元素 → el.contentEditable='plaintext-only'（Chrome 支持；不支持则 'true'）+ class ed-editing
  paste 拦截：preventDefault + execCommand('insertText', false, 纯文本)
  Enter 拦截：preventDefault + execCommand('insertLineBreak')   // 防止 h1 里长出嵌套 <div>
  blur / Esc → 提交：移除 contenteditable、class；pushUndo + reportDirty
span 保留：contenteditable 天然保留既有内联 span（xw-grad/xw-focus/mono…）——只要不全选删光
全选删光会丢 span：v1 接受，写进已知妥协
```

**撤销/重做**：

- 栈存**全文档序列化字符串**（30KB × 100 = 3MB 内存，无所谓），简单可靠、跨翻页不丢。
- push 时机：拖动/缩放结束、文本编辑提交、（未来）一切程序化变更前。
- `Ctrl+Z` / `Ctrl+Shift+Z`（或 Ctrl+Y）在 iframe 的 keydown 里拦（文本编辑进行中放行给浏览器原生 undo，先 blur 再栈 undo）。
- undo/redo 后 reportDirty()。栈只进不出（不与保存点对齐，v1 就这样，保存不清栈）。

**serialize（保存前清理，逐项照单执行）**：

```js
function serialize() {
  deselect(); exitTextEdit(); clearGuides();
  // 1. 编辑器自身痕迹
  document.getElementById('ed-styles')?.remove();
  document.querySelectorAll('[data-ed-guide]').forEach(n => n.remove());
  // 手柄/outline 都挂在选中态上，deselect 已清；双保险再扫一遍：
  document.querySelectorAll('.ed-handle, .ed-selected, .ed-hoverable, .ed-editing')
    .forEach(n => { n.classList.remove('ed-handle','ed-selected','ed-hoverable','ed-editing'); });
  document.querySelectorAll('[contenteditable]').forEach(n => n.removeAttribute('contenteditable'));
  // 2. runtime 运行态（还原到"生成产物"的原始形状）
  document.querySelectorAll('.slide').forEach(s => s.classList.remove('is-active','is-prev'));
  document.documentElement.removeAttribute('data-preview');
  document.body.removeAttribute('data-preview');
  const deck = document.querySelector('.deck'); if (deck) deck.removeAttribute('style'); // 去掉 --deck-scale/--deck-w/--deck-h/--ambient-*
  document.querySelectorAll('.notes, .speaker-notes').forEach(n => n.removeAttribute('style')); // 剥预览分支打的 inline display:none
  document.querySelectorAll('.progress-bar, .notes-overlay, .overview').forEach(n => n.remove()); // 正常模式才有的 chrome，保险起见
  // data-ed-placeholder 保留！data-current/data-total 保留（runtime 会重算）
  // 3. 取整份 HTML
  return '<!DOCTYPE html>\n' + document.documentElement.outerHTML;
}
```

- serialize 前必须切到什么页面无所谓——所有 section 都在 DOM 里（只是 opacity 隐藏），整包天然含全部页面。
- `data-id="sN"`（生成时就有）保持原样，是 AI 链路和缩略图的锚点，编辑器绝不增删改它。

**键盘（iframe 内）**：Esc 链；Ctrl+Z/Y；方向键微移选中元素（1px，Shift=10px，未选中不拦）；Ctrl+S → `preventDefault()` + 直接 serialize + post editor-serialize（与工具栏保存同一条路）。

**editor.js 不做的事**：不发任何请求（CSP 也不允许）；不碰 runtime 的翻页/主题消息；不实现多选、对齐分布、元素增删复制、页级操作（全部二期）。

### 4.5 消息协议总表（父页 ⇄ iframe）

| 方向 | 消息 | 说明 |
|---|---|---|
| →iframe | `{type:'preview-goto', idx}` | 翻页（0-based），runtime 处理，编辑器不关心 |
| →iframe | `{type:'editor-save'}` | 请求序列化；editor 回 `editor-serialize` |
| →iframe | `{type:'editor-saved'}` | 保存成功回执；editor 清脏标记 |
| →iframe | `{type:'editor-undo'}` / `{type:'editor-redo'}` | 工具栏按钮触发（iframe 失焦时用） |
| ←iframe | `{type:'preview-ready'}` | runtime 既有 |
| ←iframe | `{type:'editor-ready', pages:N}` | 就绪 + 页数（顺带解决页数来源） |
| ←iframe | `{type:'editor-dirty', dirty:boolean}` | 脏标记变化（首次变更 true、保存后 false） |
| ←iframe | `{type:'editor-selection', selected:boolean, fontSize:number\|null}` | 选中态变化；父页据此启用字号步进器并显示当前字号 |
| ←iframe | `{type:'editor-serialize', html:string}` | 全量 HTML（对 editor-save / Ctrl+S 的响应） |
| →iframe | `{type:'editor-font', factor:number}` | 字号等比缩放选中块及其内部文字（0.9/1.1） |

安全：`onMessage` 一律先 `e.source === frameEl.contentWindow` 再处理（TemplatePreviewModal 既有做法）；目标侧 `postMessage(...,'*')`（opaque origin 无源可指定，既有约定）。

### 4.6 前端外壳：DeckEditModal.vue（新）

文件：`frontend/src/components/editor/DeckEditModal.vue`。**不复用/不改动 TemplatePreviewModal**（交互前提相反：那边 iframe pointer-events-none 键盘在父页，这边 iframe 必须可点；抄结构、不复用）。

**Props/Emits**：

```ts
props: {
  open: boolean
  kind: 'deck' | 'usertpl'
  id: string            // deck-xxxx | ut-xxxx
  title: string
  canvas: { w: number; h: number }
}
emits: { close: []; saved: [] }
```

**src 计算**：

```ts
deck:    authedUrl(`/api/decks/${id}/file`,    { preview: 1, edit: 1 })
usertpl: authedUrl(`/api/user-templates/${id}/editor`, { preview: 1, edit: 1 })
```

**iframe**：`sandbox="allow-scripts"`（opaque origin，与用户模板预览同级）；`pointer-events` 默认（可交互）；`:key="src"` 打开后不重建。

**布局**（弹窗全屏黑底，中间画布同 TemplatePreviewModal 的等比容器算法 `min(94vw, (100vh-130px)*w/h)`，容器头部换成本编辑工具栏）：

```
┌ 标题 · [未保存●] · 撤销 重做 │ 保存(Ctrl+S) │ ◀ 第X/N页 ▶ │ 关闭 ┐
│                    画布（iframe）                                   │
└────────────────────────────────────────────────────────────────────┘
```

- N 来源：优先 `editor-ready.pages`，兜底 `props.canvas` 外调用方传入的 pages（文稿可用 `thumbPages()`，模板用基模板页数）。
- **父页不装全局键盘翻页**（与预览弹窗最大的差异）：iframe 聚焦时按键全归编辑器；父页工具栏按钮 + iframe 内快捷键已覆盖。Esc 关闭走「关闭按钮 + 脏确认」而非全局键盘（iframe 聚焦时父页收不到 Esc）。
- **脏确认**：`dirty && 点关闭` → 工具栏原地变成「放弃修改并关闭？[放弃] [继续编辑]」内联条（不弹原生 confirm）。

**保存流**：

```
点保存 / editor-serialize 到达：
  saving=true → PUT（deck: deckApi.saveFile(id, html, detail?)
                       usertpl: userTemplateApi.saveFile(id, html)）
  成功 → toast「已保存」+ postMessage editor-saved + dirty=false + emit('saved')
  失败 → toast 错误原文，保持 dirty（用户可重试）
```

- `editor-ready` 超时兜底：8s 未收到（v1 deck、后端未部署、脚本 404）→ 弹窗内提示「该文稿不支持在线编辑」+ 只有关闭按钮。
- 保存期间禁止重复点击（saving 状态）；保存与翻页互不阻塞（serialize 抓全文档，与当前页无关）。

**API 层**（`frontend/src/api/decks.ts` / `userTemplates.ts` 各加一个）：

```ts
saveFile: (id, html) => request(`/api/decks/${id}/file`, { method:'PUT', body: html, raw: true })
```

——`request()` 目前只发 JSON（client.ts：`headers.Authorization`，body 直接给 fetch）。**需给 client.ts 的 request 加可选 `raw` 支持：body 为字符串时 `Content-Type: text/html; charset=utf-8`、不 JSON.stringify**（改动最小，向后兼容）。
- `VersionMeta.operation` 类型补 `'edit'`。

**入口（两处）**：
1. `DeckPreview.vue` 工具条加「编辑」按钮（PhPencilSimple 图标，放「历史」旁）→ emit('edit') → `DeckView.vue` 挂 `<DeckEditModal kind="deck">`（v1 deck 无判断，靠超时兜底）。
2. `MyTemplatesView.vue` 用户模板卡（`ut-` 前缀/来自 userTemplateApi 的行）预览 chip 旁加「编辑」chip → 挂 `<DeckEditModal kind="usertpl">`（canvas 用该行 templateVisuals 补的 canvas）。内置模板卡不显示。

## 5. 边界与已知妥协（每条都是决策，不是遗漏）

| 边界 | 处理 |
|---|---|
| grid 容器里拖走一张卡，占位块没有 grid-area，可能引起自动布局的格子位移 | v1 接受；占位块尽量拷贝 `grid-column/grid-row` inline（有就拷），其余不管 |
| 全选删光文本会丢高亮 span | v1 接受 |
| 无元素 z 序控制 | 转换时把元素挪到 slide 末位（后画在上），不写 z-index；精确层级二期 |
| 两标签页/与 AI run 并发编辑 | last-write-wins，双方都有 history 版本可回滚 |
| published 用户模板 | 409 拒绝编辑，先下架 |
| v1（reveal.js）deck | 不注入；弹窗 8s editor-ready 超时提示不支持 |
| 移动端/触屏 | 不支持，仅桌面鼠标 |
| 保存的 HTML 含用户输入 | 与 AI 生成 HTML 同一信任级别：只在 sandbox iframe 里渲染（opaque origin + CSP），不额外消毒 |
| 编辑弹窗打开时 deck 缩略图/预览可能过期 | 保存后 `invalidateThumbs` 已处理；DeckPreview 的 iframe 不自动刷新（用户手点「刷新」，与现状一致） |

## 6. 实施顺序（每步独立可验证，按此 commit）

1. **后端 deck 保存链路**：`OpEdit` 常量 + `SaveHTML`（v2_edit.go）+ `SaveDeckFile` handler + 路由 + client.ts raw 支持 + `deckApi.saveFile`。验证：curl PUT 8081 → history 出现 edit 版本 → 前端历史弹窗能看到 → thumbs 失效重建。
2. **后端注入 + usertpl 两端点**：`injectEditorScript` + GetDeckFile 改造 + GetUserTemplateEditor + SaveUserTemplateFile（含备份轮转 + published 拒绝）。验证：curl 编辑态 src 看到注入 tag；备份目录出现文件。
3. **editor.js 文本编辑最小闭环**：注入生效 → editor-ready → 双击改字 → editor-save → serialize 干净（diff 检查只有文字与预期不同）→ 父页保存 → 重开还原。此时前端可先用 curl/临时页面验证。
4. **editor.js 拖动/缩放/吸附/占位块/撤销**：完整编辑能力。
5. **DeckEditModal.vue + 两处入口**：完整 UI + 保存流 + 脏确认 + 超时兜底。
6. **联调验收**（§7）+ 回归：只读预览弹窗、模板画廊、生成流程不受影响。

## 7. 验收清单（手测脚本，全部过才算完）

1. 文稿改字：DeckView → 编辑 → 改一段含 `mono` span 的文字（span 保留）→ 保存 → 重开预览内容在；history 列表顶部出现「手动编辑」版本；缩略图刷新。
2. 拖卡：把 grid 里一张卡拖到空白处 → 原位占位块、兄弟卡不动、参考线出现并吸附 → 保存 → 重开位置保持、其他页无变化。
3. 缩放：se 手柄放大卡 → 内部文字回流正常 → 保存重开一致。
4. 跨页：改第 2 页 → 翻第 5 页再改 → 保存一次 → 重开两页改动都在。
5. 撤销：拖错 → Ctrl+Z 回滚 → 重做正常；保存后栈仍可用。
6. Esc 链：文本编辑 → 退出编辑；level2 → 上选 level1 → 取消选中。
7. 模板：fork 的模板 → 编辑 index.html 改字改布局 → 保存 → 画廊/预览立即可见（no-cache 生效）→ `user-template-backups/<ut-id>/` 出现备份；连续保存 7 次只留 5 版。
8. published 模板点保存 → 409 + 错误提示。
9. AI 覆盖与回滚：手改保存 → 对话让 agent 改稿 → 手改被覆盖 → 历史恢复「手动编辑」版本 → 手改内容回来。
10. 纯净性：保存后的 index.html 里 grep 无 `ed-`、`contenteditable`、`data-preview`、`--deck-scale`、`is-active`；正常预览渲染无样式残留。
11. 沙箱/CSP 回归：编辑态 iframe console 无 CSP 报错；用户模板编辑弹窗全程 `allow-scripts` 沙箱。
12. v1 deck（deck-0001）点编辑 → 超时提示，不白屏。
13. 回归：模板预览弹窗（只读）、TemplateGallery、生成向导、导出全部照旧。

## 8. 二期 backlog（本期明确不做）

元素样式面板（字号/颜色/加粗/对齐）、元素复制/删除/层级调整、多选与对齐分布、页级操作（加/删/复制/排序页）、讲稿（.notes）编辑、乐观锁防并发覆盖、模板改动下推文稿、移动端。

## 9. 附录：改动文件总览

| 层 | 文件 | 动作 |
|---|---|---|
| 后端 | `backend/internal/router/router.go` | +3 路由：`PUT /decks/:id/file`、`GET /user-templates/:id/editor`、`PUT /user-templates/:id/file`（均 guarded 组） |
| 后端 | `backend/internal/handler/deck.go` | +SaveDeckFile；GetDeckFile 加 ?edit=1 注入分支 |
| 后端 | `backend/internal/handler/userTemplates.go` | +GetUserTemplateEditor、+SaveUserTemplateFile |
| 后端 | `backend/internal/handler/inject.go`（新） | injectEditorScript helper |
| 后端 | `backend/internal/service/deck/v2_edit.go`（新） | SaveHTML |
| 后端 | `backend/internal/service/deck/history.go` | +OpEdit 常量 |
| 后端 | `backend/internal/service/usertpl/`（service.go 或新文件） | SaveFile + 滚动备份 |
| 资产 | `backend/web/assets/deck-v2/editor.js`（新） | 编辑器本体 |
| 资产 | `backend/web/assets/deck-v2/PATCHES.md` | +一行说明 editor.js 是第一方文件 |
| 前端 | `frontend/src/components/editor/DeckEditModal.vue`（新） | 编辑弹窗外壳 |
| 前端 | `frontend/src/components/preview/DeckPreview.vue` | +编辑按钮 |
| 前端 | `frontend/src/views/DeckView.vue` | 挂弹窗 |
| 前端 | `frontend/src/views/MyTemplatesView.vue` | 用户模板卡 +编辑 chip + 挂弹窗 |
| 前端 | `frontend/src/api/client.ts` | request 支持 raw body |
| 前端 | `frontend/src/api/decks.ts` | +saveFile；VersionMeta 加 'edit' |
| 前端 | `frontend/src/api/userTemplates.ts` | +saveFile |
