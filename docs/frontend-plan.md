# 前端方案（Frontend Plan）

> 状态：已评审定稿，待动工（等后端完善后启动 M0）
> 定稿日期：2026-09-16
> 本文是前端开发唯一依据；动工后如有决策变更，直接改本文并在文末变更记录里记一笔。

---

## 0. TL;DR

为 Deck Agent（交互式 HTML PPT 生成 agent）构建正式 Web 前端：**Vue 3 + Vite + TS + Tailwind v4，全手写组件（永不引入组件库）**，明暗双主题，前后端分离部署，桌面优先。首版覆盖：登录注册、工作台（预览为主 + 对话侧栏）、Deck 列表 + 版本历史、观测台；对话历史接口与后端同步设计、排在最后实现。

## 1. 背景与现状

### 1.1 产品

Agent 通过多轮对话生成/修改交互式 HTML PPT（deck.html，reveal.js 渲染）。后端 Go + Gin 已具备完整能力：

| 能力 | 接口 | 备注 |
|---|---|---|
| 邮箱验证码登录 | `POST /api/auth/code` `/register` `/login` | JWT |
| 当前用户 / BYOK | `GET /api/auth/me`、`POST /api/auth/apikey` | |
| Agent 对话 | `POST /api/chat`（SSE 流式响应） | 携带 session_id 续轮 |
| 回答提问 / 恢复提问卡 | `POST /api/chat/answer`、`GET /api/chat/pending` | 刷新页面后重建提问卡片 |
| Deck 列表 / 预览 | `GET /api/decks`、`GET /api/decks/:id/file` | 预览给 iframe 用 |
| 版本历史 | `GET/POST/DELETE /api/decks/:id/history...` | 列表 / 恢复 / 删除 |
| 观测数据 | `GET /api/traces...`（列表 / run / 事件 / 截图 / 导出） | 归属复核在 store 内做 |
| 一次性取页 | `GET /api/render/:nonce` | 视觉审查用，前端不接触 |
| 静态资源 | `/assets/*` | reveal.js + deck 主题/组件库 |

### 1.2 SSE 事件协议（`internal/agent/stream_event.go`，前端必须全覆盖）

| 事件 | 语义 | 工作台 UI |
|---|---|---|
| `delta` | 正文文本增量 | agent 气泡流式追加 |
| `think` | 推理链增量 | 折叠"思考"块（`<details>` 思路：闭合时子树不参与布局） |
| `tool_start` | 工具调用开始（id + name，参数开始生成） | 工具卡片出现 |
| `tool_delta` | 工具参数增量 | 卡片内参数流式填充，默认折叠 |
| `tool_call` | 工具调用成功（含 `tool_index` 区分并行调用） | 卡片定型为成功 |
| `tool_error` | 工具调用失败 | 卡片定型为失败（红） |
| `sub_delta` | 子调用文本增量（视觉审查报告等） | 子模型流式盒，默认展开 |
| `ask_user` | 结构化提问卡片 | 提问卡片，回答走 `/chat/answer` |
| `session` | 下发 session_id | 前端保存，续轮携带 |
| `trace` | 观测事件（Content 是 trace.Event JSON） | **工作台默认不显示**（归观测台），留开发者开关 |
| `done` | 完成，带分项 Usage（prompt/completion/total/cached） | 结束态，用量入对话元信息 |
| `error` | 出错 | Toast + 对话内错误条 |

已知注意点：`tool_index` 用于区分并行工具调用；未知事件类型一律**静默忽略**（协议明确要求）。

### 1.3 现有前端资产

- `backend/web/chat-test.html`、`trace.html`：手工测试台，**验证了交互协议但不迁移代码**。保留作为后端回归工具，不删除。
- `backend/web/assets/`（theme.css / components.css / reveal.js）：**给生成的 deck 用的展示层**，与应用 UI 是两套体系，互不复用。
- 已知妥协：deck 预览与 trace 截图的 `<iframe>/<img>` 带不了 Authorization 头，靠 `?token=` 回退（后端 `middleware/auth.go` 有说明），前端沿用。

## 2. 已定决策（2026-09-16 评审结论）

| 决策项 | 结论 |
|---|---|
| 框架 | Vue 3（Composition API + `<script setup>`）+ Vite + TypeScript |
| 组件方案 | **Tailwind CSS v4 + 全手写组件，永久不引入组件库** |
| 状态管理 | Pinia |
| 路由 | Vue Router，history 模式（部署侧 `try_files` 兜底） |
| 部署形态 | **前后端分离部署**；dev 用 Vite proxy 联调，生产走后端现有 CORS 配置 |
| 视觉 | 明暗双主题可切换，CSS 变量 token 体系，默认跟随系统 |
| 工作台布局 | **预览为主，对话为可收起侧栏** |
| MVP 范围 | 全功能：工作台、登录注册、Deck 列表 + 历史版本、观测台 |
| 对话历史（后端暂缺） | 前后端一并设计，接口契约见 §7，实现在 M5 |
| 移动端 | 桌面优先，窄屏降级为"对话 / 预览"双视图切换 |
| 界面语言 | 中文 |

## 3. 设计系统

### 3.1 双主题机制

- `<html data-theme="light|dark">` 驱动，不用 media query 做切换来源；
- Tailwind v4 里 `@custom-variant dark (&:where([data-theme=dark] *))`，让 `dark:` 前缀跟随 data-theme 而非系统；
- token 分两层：**原始色板**（不直接用）→ **语义 token**（`--color-bg / --color-surface / --color-surface-hover / --color-border / --color-text / --color-text-secondary / --color-text-muted / --color-accent / --color-danger / --color-success / --color-warning` 等），语义 token 在 `:root[data-theme=...]` 下各给一套值；
- `@theme` 把语义 token 映射成 utility（`bg-surface`、`text-muted`…），页面代码**只允许用语义类，不允许直接写色值**；
- 切换偏好存 localStorage，未设置时跟随 `prefers-color-scheme`；
- `index.html` 内联一小段脚本提前设置 `data-theme`，防首帧闪烁（FOUC）。

### 3.2 手写组件政策（硬约束：永不引入组件库）

- 基础件统一放 `src/components/ui/`，**按页面需要渐进补充**，不一次性造全套；起步批：Button、Input、Textarea、Dialog、Toast、Spinner、Empty、Badge、Tooltip、Tabs、Dropdown、Select；
- 涉及的细节按规范做：Dialog 的 focus trap / ESC 关闭 / 背景滚动锁定，Toast 用 portal 挂 body，所有可交互件有 hover/focus-visible/disabled 态；
- **页面里禁止裸写 magic style**：颜色/间距/圆角一律走 token 或基础件，这是没有组件库的项目不失控的关键纪律；
- 图标用 `lucide-vue-next`（SVG 图标库，不属于组件库约定范围）。

## 4. 工程结构

```
frontend/
├── index.html              # 含防 FOUC 的主题内联脚本
├── vite.config.ts          # dev: /api 代理到 localhost:8080
├── src/
│   ├── main.ts
│   ├── App.vue
│   ├── router/             # 路由表 + 鉴权守卫
│   ├── stores/             # auth / deck / chat / trace / theme
│   ├── api/                # fetch 封装 + auth / decks / chat / trace 模块
│   ├── lib/sse.ts          # POST + ReadableStream 的 SSE 解析器
│   ├── theme/              # tokens.css、useTheme、ThemeToggle
│   ├── components/ui/      # 手写基础件（见 §3.2）
│   ├── components/chat/    # 消息流、ask_user 卡片、输入框
│   ├── components/preview/ # iframe 预览器、工具条、历史抽屉
│   └── views/              # Login / Decks / Workspace / Settings / Trace
```

要点：

- **SSE 用 fetch + ReadableStream 手工解析**（`EventSource` 不支持 POST 与 Authorization 头），解析器独立成 `lib/sse.ts`，支持中止（AbortController）与断流错误上报；
- **fetch 轻封装**（不上 axios）：baseURL、token 注入、401 统一跳登录、业务错误规范化为统一 Error 形状；
- 消息流渲染走 **`ViewEvent[]` 视图模型**：SSE 事件先归一成视图模型再渲染，流式态与回放态（M5 历史回放）共用同一套组件，后端历史接口一通即可零改动接入；
- stores 划分：`auth`（token/用户/BYOK）、`deck`（列表/当前 deck/版本历史）、`chat`（按 deck 的会话状态、ViewEvent 流、流式中标志、pending 提问）、`trace`（run 列表/详情）、`theme`（主题偏好）。

## 5. 页面与路由

| 路由 | 页面 | 要点 |
|---|---|---|
| `/login` | 登录/注册 | 验证码 60s 倒计时；登录/注册同页切换；成功后回跳来源页 |
| `/decks` | Deck 列表 | 卡片：标题、更新时间；点击进入工作台；空态引导新建 |
| `/decks/:id` | **工作台** | 核心页，详见 §6 |
| `/settings` | 设置 | API Key 设置（遮显）、账号信息 |
| `/trace` | 观测台 | run 列表 → run 详情：事件时间线、耗时、审查截图、导出；重写 trace.html 为受保护路由 |

全局：路由守卫校验 token（无 token → `/login`）；token 持 localStorage（与现有测试台同源策略一致）。

## 6. 工作台详细设计

### 6.1 布局（桌面）

```
┌──────────────────────────────────────────┬──────────────┐
│ 工具条: 页码 · 刷新 · 新标签打开 · 历史版本   │  对话侧栏      │
├──────────────────────────────────────────┤  ┌──────────┐ │
│                                          │  │ 消息流     │ │
│       deck.html 预览（iframe）             │  │ (可滚动)  │ │
│       占绝大部分宽度                        │  └──────────┘ │
│                                          │  [输入框][↑]  │
└──────────────────────────────────────────┴──────────────┘
```

- 侧栏可收起、宽度可拖拽；
- 版本历史做成**工作台内抽屉**（列表 / 按版本预览 / 恢复 / 删除），不跳页；
- 窄屏（<768px）：预览与对话切换为两个 tab 视图，工具条简化。

### 6.2 预览区

- iframe src = 后端域名 `/api/decks/:id/file?token=...`（token 回退是已确立的已知妥协，沿用）；
- iframe 直连后端域名，不受 CORS 约束；sandbox 属性沿用 chat-test 验证过的安全模型；
- agent 报告"写完某页"或 `done` 后提供刷新（手动 + 自动策略实现时定）；
- 提供新标签页全屏打开、页码/滚动状态提示。

### 6.3 对话侧栏

- 消息流按 §1.2 映射渲染全部事件；`trace` 事件默认不渲染（开发者开关放设置或侧栏菜单）；
- `ask_user` 卡片：结构化表单，提交走 `POST /api/chat/answer`；进入工作台时调 `GET /api/chat/pending` 恢复未答提问；
- 流式中输入框变"中止"按钮（AbortController 断流，会话可续轮）；
- `done` 事件上的分项 Usage 折叠展示在对话元信息里；
- 空态：新 deck 显示引导语与示例指令。

## 7. 对话历史接口契约（草案，后端按此实现）

```
GET /api/decks/:id/messages?after_seq=&limit=
→ 200 { messages: [
    { seq: int64,            // 会话内单调递增，增量拉取游标
      role: "user"|"assistant"|"tool",
      kind: "text"|"think"|"tool_call"|"tool_result"|"ask_user"|"answer",
      content: string,        // kind 决定解读方式
      tool_name?: string,
      created_at: string }
  ] }
```

- **按 deck 维度组织**：进入某 deck 即见与之相关的完整对话；
- `after_seq` 支持增量拉取；进入工作台 = 全量拉取 + 灌入 ViewEvent[]；
- 后端实现建议：**追加式事件日志**（append-only，按 session 存储），本接口从日志投影——观测台与对话历史共用同一份底账，避免两套持久化；
- 前端排期在 M5，但消息流组件从 M2 起就是视图模型驱动（§4），接口不阻塞 UI。

## 8. 分期计划

| 期 | 内容 | 完成标准 |
|---|---|---|
| M0 脚手架 | Vite+Vue+TS+Tailwind v4+Router+Pinia；双主题 token 与切换；dev 代理；fetch 封装；第一批基础件 | 空壳页面跑通登录态路由跳转；明暗切换全站生效无闪烁 |
| M1 鉴权 | 登录/注册页、路由守卫、设置页（API Key） | 验证码注册→登录→进列表全流程可用；401 自动跳登录 |
| M2 工作台 MVP | 双栏布局、iframe 预览与刷新、SSE 全事件渲染、发送/中止、pending 恢复 | §1.2 的 12 种事件全部有对应 UI；刷新后提问卡片可恢复；流式中可中止 |
| M3 管理 | Deck 列表页、版本历史抽屉 | 列表进入工作台；历史可预览/恢复/删除 |
| M4 观测台 | run 列表、事件时间线、耗时、截图查看、导出 | trace.html 的能力全覆盖且挂鉴权 |
| M5 对话历史 | 后端按 §7 实现 + 前端回放 | 进工作台可见旧对话，流式/回放渲染一致 |

约束：M2 起布局就用响应式写法，避免 M5 后返工；每期独立可交付、可演示。

## 9. 已定的小默认（有异议随时推翻）

1. 图标：`lucide-vue-next`；
2. Markdown：agent 正文用 `markdown-it` 渲染、代码高亮用 `shiki`（工具库，不违反无组件库约定）；
3. HTTP：轻封装原生 fetch；
4. 包管理器：pnpm；
5. `chat-test.html` / `trace.html` 保留作后端回归工具，不删；
6. 前端目录为仓库根的 `frontend/`，与单仓库结构一致。

## 10. 待决问题（不阻塞 M0–M2，动工对应期前定）

1. **JWT 过期策略**：目前疑似过期即重新登录。是否引入刷新/静默续期机制？→ 影响 M1 的 401 处理与 token 存储设计；
2. **观测台受众**：纯自用运维工具，还是对普通用户开放？→ 决定 M4 的打磨程度与入口显眼度；
3. 预览 iframe 的 `?token=` 妥协的长期方案（短期沿用现状）。

## 11. 环境要求

- Node ≥ 20.19（Vite 7 / Tailwind v4 要求），pnpm；
- dev 联调走 Vite proxy 无 CORS 问题；**生产分离部署时**，后端 `config.yaml` 的 `server.allow_origins` 需加入前端正式域名；
- 后端接口以 `internal/router/router.go` 为准，本文 §1.1 与其同步维护。

---

## 变更记录

- 2026-09-16：初稿定稿（讨论结论：Vue3 + Tailwind 手写组件 + 双主题 + 分离部署 + 预览为主布局；对话历史接口一并设计）。
