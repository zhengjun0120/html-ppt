# Agent 观测（trace）

观测台地址：**http://localhost:8080/trace**（与 `/chat-test` 一样是免登录的静态页，
数据接口全部在鉴权组里）。

## 它解决什么

SSE 那条流是**给用户看的对话**（文字增量、思考、提问卡片）。但排查问题时你要问的是
另一类问题，而它们的答案过去**根本不存在**：

| 你想知道的 | 改动前的状况 |
|---|---|
| 模型给 `update_slide` 传了什么参数 | 参数从来没被任何通道外发过（SSE 的 `tool_call` 只有工具名和返回值） |
| 联网搜索到底搜了什么 | `web_search.go` 里 `BetaServerToolUseBlock` 是个空 case，注释写着"后续可以传回给前端显示"——只剩一个 `searches=2` 的计数 |
| 审查报告说"第 3 页文字被裁"，真的是吗 | 截图发给审查模型之后就丢了（`vision.Slide.PNG` 是 `json:"-"`），没有任何办法核对 |
| 这次对话花了多少 token | 只有主循环的数，**视觉审查与联网搜索的用量完全没被计入**（它们是分开计费的另外几次 API 调用） |
| 模型这一轮到底看到了什么 | 不存在 |

## 磁盘布局

```
data/traces/<session_id>/<run_id>.jsonl         一行一个事件（追加写，边跑边可读）
data/traces/<session_id>/<run_id>/img/v1-p003.png   视觉审查截图（v<第几轮>-p<页号>）
```

- `run_id = <13位unix毫秒>-<4位随机hex>`。定宽时间戳在前 → 字典序即时间序，
  列表与裁剪都不需要额外排序。
- **一次对话可能对应多个 run**：`ask_user` 会让循环暂停，用户回答后是新的一次执行。
  恢复时新 run 的 `parent_run_id` 指回被暂停的那个，token 需要跨 run 累加才是"这次对话"。
- `data/` 已在 `.gitignore` 里，无需额外规则。

## 配置

```yaml
features:
  trace: true          # 开关。关掉后连上下文都不会被序列化（真零开销）

trace:
  dir: ""              # 空 = <data.dir>/traces
  capture_images: true # 存截图。关掉后报告仍在，但"报告说的对不对"就无从核对
  retain_runs_per_session: 20   # 每个会话最多留几个 run，0 = 不限
  max_field_bytes: 0   # 0 = 不截断。写满磁盘时打开它瘦身
```

⚠️ `trace.dir` 的相对路径基准是 **config.yaml 所在目录**（与 `data.dir`/`assets.dir` 一致），
不是 cwd。这条在 `config.Load` 里显式处理，漏了的话表现是"换个目录启动就写到别处"，
而 trace 会照常写进一个新目录、页面照常是空的——看起来像"观测没生效"。

## 事件种类

| kind | 何时发 | 关键字段 |
|---|---|---|
| `run_start` | run 开始 | `user_content`、`model`、`tools`（本次挂载了哪些工具）、`parent_run_id` |
| `llm_request` | 每次请求前 | `turn`、`messages`（**全量上下文**）、`message_count`、`bytes` |
| `llm_response` | 每次响应后 | `turn`、`content`、`tool_calls[]`、`finish_reason`、`usage` |
| `tool_call` | 工具执行前 | `tool_name`、`tool_call_id`、`args`（**参数原文**）、`turn` |
| `tool_result` | 工具执行后 | `result`、`duration_ms`、`error`（失败时） |
| `sub_step` | 工具内部 | `sub.name`/`sub.stage`、`sub.text`、`sub.data`、`images[]` |
| `usage` | 每次计费调用后 | `component`（main/vision/web_search）、`usage` |
| `error` | 中途出错 | `error`（流断了、写库失败；原因进 trace，不只留在日志里） |
| `run_end` | 结束/暂停/出错 | `status`（ok/paused/error）、`summary`（含分项用量与合计） |

`turn` 是 0 基的轮次。它是**指针类型**：`第 0 轮` 与 `不适用`（run_start/run_end）必须分得开——
用 `int + omitempty` 的话第一次请求的轮次会被静默丢掉，页面上只有第二轮回合起才显示轮次。

### sub_step 的 stage

| name/stage | 内容 |
|---|---|
| `web_search/request` | 我们让它搜什么（`text`）、请求的 `max_uses` |
| `web_search/queries` | **子模型实际发出的搜索词**（`data.queries`）、计费次数、来源数、stop_reason、工具级错误 |
| `vision/grant` | 为审查签发的一次性取页票据（2 分钟过期、取到即废） |
| `vision/capture` | 逐页量测（`data.pages`：fit / 最小字号 / 溢出 / 子元素 / 字数 / 版式）+ **每页截图** |
| `vision/review_prompt` | 实际发出去的审查提示词全文 + 附了几张图 |
| `vision/review_response` | 审查报告 |
| `vision/error`、`vision/review_error` | 取页/量测失败、看图失败的原因 |

`web_search/request` 与 `web_search/queries` 要成对看：子模型会**自己改写搜索词**
（实测同一个问题两次，它发出的 query 分别是中文和英文两版），两相对照才看得出它有没有跑偏。

## HTTP API

全部挂在鉴权组下（`Authorization: Bearer <jwt>`）。归属复核在 `trace.Store` 内做：
不属于你的 run 一律 **404，且与"不存在"同一个响应**——区分它们就等于泄露"这个 run 存在"。

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/traces` | `{sessions:[...], runs:[...]}`。`sessions` 是按会话**跨 run 累加**的汇总（回答"这次对话花了多少"），`runs` 是扁平列表。可选 `session_id` / `limit` / `offset` |
| GET | `/api/traces/:sid/:runID` | `{events:[...], next_offset, running}`。可选 `from_offset`（增量，**实时跟随用这个**）/ `from_seq` / `include=messages` / `limit` |
| GET | `/api/traces/:sid/:runID/events/:seq` | 单个事件的**全量**版本（含 messages）。页面上"展开完整上下文"走它 |
| GET | `/api/traces/:sid/:runID/img/:name` | 截图（`image/png`）。**用 `?token=` 传令牌**——`<img src>` 带不了 Authorization 头 |

读接口**不受 `features.trace` 影响**：关掉的是"继续记录"，已经落盘的记录应该照样能看。
`trace.Store` 无论开关如何都会装配，否则调一次开关就把之前跑出来的东西变成读不到的孤儿文件。

## 一个刻意的取舍：messages 只落盘、不进实时流

`llm_request.messages` 是文件里唯一的大字段（每轮重发整份消息，开了全量上下文就是几 MB），
而 `serveAgentSSE` 的通道只有 16 个槽——推它会把聊天流拖住。所以：

- **落盘**：全量存（`max_field_bytes: 0`）。
- **实时 SSE 的 `trace` 事件**：剔除 `messages`，只留 `message_count` 与 `bytes`。
- **事件列表接口**：默认剔除，显式 `?include=messages` 才带。
- **单事件接口**：带全文，页面按需拉。

`web/chat-test.html` 里加了一个默认关闭的"观测事件"勾选框：不认识 `trace` 类型的前端
会把它归到 `default` 分支刷"未知事件类型"，所以这个 case 必须显式处理。

## 设计纪律

- **fail-open**：任何写盘/回调错误只 `log.Printf("[warn] ...")` 并继续。观测是排查手段，
  为它把用户正在进行的对话搞挂，等于用一个诊断工具制造故障（与 `main.go` 对视觉审查同一条原则）。
- **实时回调查在锁外**：它要做通道发送，通道满时会阻塞；持着锁阻塞 = 整个 agent 循环卡在观测上。
- **归属走 context**：工具不知道自己在第几轮、call_id 是什么，只调 `trace.Emit(ctx, ...)`；
  拿不到 recorder 时它是空操作（`From` 永不返回 nil），所以测试与其他调用方零改动。
- **`trace.Active(ctx)`** 只用在"要先付出构造代价才能发的事件"上（序列化几 MB 上下文），
  其余不必套——`Discard` 本身就是空操作。
- **图片名过白名单** `^[a-zA-Z0-9_-]{1,32}\.png$`，run_id/会话号同理。
  三个入参都来自 URL，任何一个漏校验都能读到别人的文件。

## 排查

| 现象 | 先看哪里 |
|---|---|
| 观测页一片空白、显示"还没有观测记录" | `features.trace` 是否打开；`trace.dir` 的实际路径（启动日志里有打） |
| 列表有 run，但 token 一栏是空的 | 那个 run 还没结束（没有 `run_end`）——页面会标"运行中" |
| 没有截图，只显示"（未保存图片）" | `trace.capture_images` 是否打开；以及 `features.vision` 是否打开（**关着的话根本没有视觉审查**） |
| 事件列表里 `messages` 是空的 | 这是默认行为。用 `?include=messages`，或在页面上点"展开完整上下文" |
| run 一直显示"运行中" | 进程在那次对话中途被杀（没有 `run_end`）。`running` 由文件尾部判断，不是"这一段里没见到 run_end" |
