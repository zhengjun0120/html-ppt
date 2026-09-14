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
| 8 | `update_theme` | 2 | 改主题配置（风格预设/配色/字体配对/圆角/纹理/动画/画布） | ☑ |
| 9 | `list_templates` | 2 | 列出模板供推荐 | ☐ |
| 10 | `ask_user` | 2 | 向用户提问（human-in-the-loop，暂停循环） | ☑ |
| 11 | `screenshot_slides` | 3 | 截图/程序化检查排版 | ☐ |
| 12 | `move_slide` | 2+ 可选 | 调整页序 | ☐ |
| 13 | `read_history_diff` | 4 | 对比两份状态，看清某一轮改了什么（只读） | ☑ |
| 14 | `list_history` | 4 | 版本历史列表，取版本号（只读，分页） | ☑ |
| 15 | `read_custom_css` | 4 | 读 deck 级自定义样式槽（只读） | ☑ |
| 16 | `update_custom_css` | 4 | 整体替换自定义样式（清洗规则 + 可回滚） | ☑ |
| 17 | `read_component` | 4 | 读共享组件库实现（只读，按类名或全文） | ☑ |
| 18 | `read_theme` | 4 | 读当前主题配置（只读；`update_theme` 的 `vars` 是整体替换制，改前必读） | ☑ |
| — | 样式体检（非工具） | 4 | 三个写工具结果里的 `warning`：写死颜色/px 字号/section 级覆盖/重复内联 | ☑ |

实现顺序：1 → 3/4/5 → 6/7 → 8/9 → 10 →（阶段3）11 → 12 →（阶段4 版本控制）13/14 →（阶段4 样式槽）15/16/17 → 18（含 §19 样式体检）。

> 15/16/17 由 `features.custom_css` 开关统一控制（见 config.go 的 Features）：**关掉时连工具都不挂载**，
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
- **样式体检（与消毒同一条回报通道）**：写死颜色、px 字号、`<section>` 级主题覆盖、大面积重复内联——
  这些不危险，但会让这份 deck 再也换不了风格（`update_theme` 静默失效），同样走 `warning` 提示
  而不阻塞写入。判定在 `service/deck/stylelint.go`，见 §18
- **回声校验（变量合法性，与上面两类不同：这条走拒收）**：把框架 CSS 里**被 `var()` 真正消费过**的
  变量名扫成契约（`service/deck/variable_contract.go`，每次现扫不缓存，新加变量立刻生效）。
  页面内联 style 与自定义 CSS 里引用到的变量逐个校验，**契约外且本处未定义的变量直接拒收**。
  防的是 `var(--acent)`（拼错）、`var(--my-accent)`（凭空造）——页面不报错、结构也没坏，
  只是那处样式不生效，而 AI 会汇报"已改好"，是最难自查的一类假成功。
  豁免：同一份提交里定义过、或该 deck 主题块（`vars` 调色板）里定义过的变量都算已知——
  少了这条，"vars 里定义、页面里引用"这种正常写法会被误判（误拒比漏判更糟）。
  契约读不到时跳过校验（fail-open）：它防的是静默无效，不是安全问题，不能把写入全卡死
- 防工具混乱三件套：description 写清"什么时候不用我"；system prompt 里给标准工作流；按场景分组渐进挂载（阶段2末）
- **工具参数的 `description` 里不能出现半角逗号**：`invopop/jsonschema` 的标签解析按半角逗号切键值对，
  不认引号也不认转义——`description` 里有一个半角逗号，从那里往后整段描述**静默丢失**。
  实测在 `vars` 的描述里放了个 JSON 示例，模型收到的描述就断在 `如 {"--surface":"#1e1836"`，
  恰好把最关键的用法说明吃掉。一律用全角逗号/顿号/分号，且 `description` 必须是标签最后一项；
  `TestSchemaDescriptionsHaveNoASCIIComma` 守着这条（这类"给模型的文档悄悄少一半"联调时发现不了）

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

## 8. update_theme ☑ / 8b. read_theme ☑

- **参数**：`{ deck_id, preset?, accent?, background?, heading_color?, text_color?, surface?,
  border_color?, font?, radius?, texture?, transition?, canvas?, vars? }`
  （只传要改的；`preset` 枚举见 §8a，`font` 枚举 sans/serif/editorial/modern/mono，
  `texture` 枚举 none/grid/dots/rule，`transition` 枚举 slide/fade/zoom/convex/concave/none，
  `canvas` 枚举 standard/wide/classic）
- **要点**：参数是语义化字段而非 CSS；真身是 deck.html 里的
  `<script type="application/json" id="deck-theme">` 配置块，读 JSON → 合并 → 写回 → 同步 CSS 变量；
  校验颜色格式正则 + 枚举；派生变量（border/card-bg/text-muted/accent-soft/on-accent）集中计算。
  `read_theme` 是配对只读工具，返回同一份配置。
- **`surface` / `border_color` 为什么是字段**：原来的派生（`--card-bg` = accent@6%、
  `--border` = accent@25%）只在"暗底 + 单一强调色"这套视觉里成立。换成浅色纸面之后，
  accent@6% 会得到一层粉、accent@25% 会得到一条彩线，而纸感/编辑风要的恰恰是
  "中性面板 + 墨色发丝线"。**派生关系的根基是视觉语言，不是颜色值**，所以它不该被写死成
  一个公式。两个字段留空仍然派生——老 deck 与科技暗色主题的观感一字不变（`fillMissing`
  刻意不补这两项，因为空值在这里是有语义的）。
- **`--on-accent` 修掉一个真 bug**：组件库的 `.badge` 原来写死 `color:#0b132b`，
  只在"强调色一定是亮色"时成立；编辑风预设的强调色是黑，写死就成了黑底黑字。
  现在按强调色的 WCAG 相对亮度算，自动选对比度更高的前景（`onAccentColor`）。

### 8a. 风格预设：`preset` —— 有名字的整套美学

用户说"换个风格/专业一点/好看一点"时，`preset` 是比逐个调色更好的入口：
一个名字背后是一整套设计决定（配色 + 面板色 + 字体配对 + 圆角 + 纹理 + 语义色）。
六个预设：`paper`（纸感，默认）/ `editorial`（编辑风）/ `noir`（暗夜电影）/
`duotone`（双色印刷）/ `terminal`（终端）/ `tech`（科技暗色）。

**为什么需要它——这是整个"去 AI 味"改动的核心判断**：让模型自由发挥配色，它必然回到
训练分布的中位数（深蓝底 + 发光强调色 + 圆角卡片 + 无衬线字体）。这是统计规律，不是态度
问题：它见过的绝大多数"现代网页"就长这样，"求稳"就等于选它。**所以提高自由度治不了
同质化，给出具体的参照物才可以。**

**预设是写路径上的糖，不是活引用**：套用一次就把整套字段落进主题 JSON，渲染路径完全不知道
预设的存在（只认存储的字段）。反过来做（存 `preset` 名、渲染时现查表）意味着改一次预设定义
会静默改掉所有用过它的 deck——用户没碰过的页面自己变了样，正是这个代码库最不想制造的那类
"不可见的变化"。

**`preset` 字段的语义**：它是"最后一次套用的预设名"，**不是**"当前观感等于这个预设"。
所以任何一次不带 preset 的观感字段修改都会清空它（`aestheticFields`）——一个会说谎的字段
比没有这个字段更糟，模型会照着它向用户断言"当前是纸感风格"，而页面早就不是了。
同一调用里既选预设又微调则保留（那种情况下预设是它的底子）；只改 `transition`/`canvas`
不算偏离（那是演示参数，不是视觉语言）。

**默认主题是 `paper` 而不是深底发光色**：默认值决定大多数结果——用户多数时候不会说
"换个配色"，于是默认那一套就是这份工具的观感。深底 + 发光强调色 + 圆角卡片是最容易撞衫的
一类，所以把它降级成可选预设（`tech`），让它在预设表里被明确标注"最容易撞衫，
只在用户点名要科技感时才用"。

**字体从"三档系统栈"升级成配对**（`fontPairs`）：标题与正文用同一套栈是"没有做过字体决策"
的默认结果；而"衬线标题 + 无衬线正文"这一对（`editorial`）光换上去就能让页面从
"AI 生成的网页"变成"有人排过版的读物"。`renderThemeCSS` 分别渲染
`--r-main-font` / `--r-heading-font`。

**枚举的两处真相**：`preset`/`font`/`texture`/`canvas`/`transition` 的取值同时存在于
struct tag 的 `enum`（模型看到的）与 deck 包的表（校验与渲染真正认的）——jsonschema 库不
支持动态枚举，没法合成一处。`TestThemeSchemaEnumsMatchCodeTables` 守两边一致，
漂移的后果是"模型写得出、却被校验拒掉"，报错看起来毫无道理。

**纹理**（`texture`）渲染成一条背景图案规则挂在 `.reveal-viewport` 上：
moon.css 就是往它身上写背景色的（同元素、同优先级，排后面即生效）；不用伪元素，
因为 reveal 会改 section 的 transform，装饰挂在被变换的节点上会跟着翻页抖动。
线色由**正文色**派生而不是固定灰：深色主题上浅、浅色主题上深，同一个枚举两种底色都成立。

### 配色自由度：`vars` —— 让 AI 设计一整套配色，但只有一个权威源

用户要"一套最合适的配色"是真的需求，`vars` 就是它的通道：一组
`变量名 → CSS 值`，渲染进 `#deck-theme-override`，页面里用 `var(--surface)` 引用。
预设已经预置了 `--accent-2` / `--positive` / `--warn` 三个语义色：一份 deck 只有"一个
强调色"时，需要区分正负或画双色对比的页面就只能就地写死颜色（那会让 `update_theme`
之后失效）。这三个刻意**不被组件库消费**（否则没带预设的 deck 上会渲染成不可见），
它们由预设定义、由模型在页面里 `var()` 使用。

**为什么必须走变量而不是允许就地写死颜色**：可再修改（用户说"主色再暖一点"改一处即可，
写死就得重写整个 deck）、不重复（省 token）、一处权威（"现在什么配色"永远答得出来）。
这不是洁癖，是把"自由"和"还能改"同时拿住。

**`vars` 是整体替换制**（不是合并），理由与自定义样式槽一致：一套配色是一个整体，
不能靠增量拼接累积出意外。代价是改一个变量也要提交全文——所以配了 `read_theme` 作为前置读，
和 `read_slide` / `read_custom_css` 是同一套路。

**校验**（`theme_vars.go`）：名字必须是 `--` 开头的合法新变量；值里不许出现
`; { } < > \` `@import` `url(` `expression(`（变量值会拼进 `:root{}`，一个 `}` 就能提前收尾）；
上限 32 个变量、单值 200 字节。

**契约变量不许重定义**：`--accent` `--border` `--card-bg` `--text-muted` `--radius`
`--space-*` `--r-*` 只由结构化字段派生。这不只是洁癖，是堵一个静默失效：

> `#deck-theme-override` 和 `#deck-custom` 都是 `:root{}`（优先级相同），而自定义槽排在主题块
> **之后**——同级靠后者取胜。所以在槽里写 `:root{--accent:#ff8800}` 会静默盖住主题：
> 用户之后说"换个配色"，`update_theme` 会成功返回、页面一动不动。
> 级联顺序本身是对的（槽本来就该能覆盖主题做局部视觉），错的是"让槽重定义契约变量"。

所以 `validateCustomCSS` 现在直接拒绝重定义契约变量，错误信息把模型引到 `update_theme`
（含 `vars`）。判定前先剥注释——否则一句 `/* 别在这里写 --accent: 值 */` 会被误杀，
而那条注释恰恰是对的。反过来，在槽里定义**新**变量（`:root{--brand-ink:…}`）仍然允许。

### 画布自由度：`canvas` —— 只给比例预设，不放开绝对尺寸

`wide`（1244×700，**默认**）/ `standard`（960×700，reveal 自带）/ `classic`（933×700，4:3 老投影仪）。
**高度统一 700，只调宽度**：一页能放多少内容由高度决定，宽度只影响排布宽松度；
高度一动，原本刚好放得下的页面就会被挤爆（触发适配兜底缩小）。

**为什么默认不是 reveal 的 960×700**：字号抬到投影可读下限（正文 29pt）之后，
一页能放多少字由**可用宽度 × 高度**决定，而 standard 在 16:9 屏上左右还要留黑边——
真正用于排版的宽度只有 16:9 画布的 77%。wide 一次解决两件事：横向多约 30% 版面
（双栏每栏从 10 字/行变成 12~13 字/行、三列从"只放数字"变成"放得下 7 字短标签"），
以及投屏/录屏不再有黑边。高度没变，所以"一页能放几行"一个字都没变——
**换 wide 是纯粹的横向增益，不会把任何现有页面推进缩兜底**。

`DefaultCanvas`（`theme.go`）是"新建 deck 用哪个画布"的唯一出处，`defaultTheme()` 与
`init.js` 的兜底表达式都从它来；`TestCanvasPresetsMatchInitJS` 守两边不漂移。

为什么不放开绝对尺寸（deck-0004 写了 `width:1100px` 和 18 处 `font-size:30px`）：
画布尺寸 + em 尺度是**适配兜底、翻页动画、字号缩放**三件事的共同地基，绝对尺寸会让三者同时失灵。

预设表在 Go（`theme.go` 的 `canvasPresets`）和 `init.js`（`CANVAS`）各存一份——
前端必须在 `Reveal.initialize` 之前拿到尺寸，没法从后端要。适配兜底的 `config()`
读的是 `Reveal.getConfig()`，天然与所设尺寸同源。


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
  - 花括号不配对拒绝（漏一个 `}` 会让这之后的**所有规则**一起失效，见表头的 `checkCSSBalance`）
  - 重定义主题契约变量拒绝（见 `theme_vars.go`：会静默盖住 `update_theme`）
  - 引用契约外且本处未定义的变量拒绝（回声校验，见下方通用骨架里的说明）
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

## 17. read_component —— 共享组件库的只读视图 ☑

**为什么需要**：AI 要写 `update_custom_css` 覆盖 `.card`/`.quote`，就得先知道它们现在长什么样。
不给它这个视图，它只能凭想象写选择器和属性，覆盖出来的效果全靠运气。

- **参数**：`{ name? }`——类名（不含点号，如 `card`）。不传 = 返回组件库全文 + 可用类名清单
- **返回（全文模式）**：`{ components_css, classes[], note }`
- **返回（指定类名）**：`{ name, rules, note }`——该类相关的规则原文（含 `.card h3` 这类后代规则）
- **找不到类名时**：报错并附可用类名清单（"错误即反馈"，模型据此改名）
- **`note` 里固定带一句优先级提醒**：组件库的选择器都带 `.reveal` 前缀（`.reveal .card`，两个类），
  自定义槽里写裸的 `.card` 会**因优先级不足被盖掉**（无论级联顺序如何）。这是"写了 CSS 却没生效"
  最常见的原因，必须让模型知道，否则它会一直以为自己写对了。

**要点**：
- 纯只读：工具只读文件，组件库也永远不进任何写工具的目标列表（护栏）
- 每次调用重新读盘（文件 ~2KB）：开发期改了 `components.css`，AI 立刻看到新版本，不做缓存
- 解析用"按 `}` 切块 + 类名 token 边界匹配"（`-`/`_`/字母/数字为边界，所以 `.card` 不会命中 `.card-x`），
  切块前先剥注释（注释里的 `}` 会把块切歪）。组件库是平铺规则、无嵌套 at-rule，够用；
  将来真出现 `@media` 嵌套再换正经解析
- **契约测试**：`TestRealComponentLibraryContract` 会拿真实 `components.css` 校验
  "prompt 承诺的组件类（v1 九个 + v2 页面家具/行清单/行内角色）都真实存在"
  ——文档与实现漂移是最难在联调中发现的一类 bug。
  prompt 里每新增一个 class，这个测试就得跟着加一条

**组件集 v1 → v2 的由来**：v1 只有 9 个 class，表达不了"页码、讲次提头、行清单、术语词"
这些每页都出现的结构，模型只能手写内联补齐。实测三份手写型 deck，可编辑区
**43%~51% 的字符是逐字重复的内联样式**（同一串 177 字符的页码角标在 8 页里一字不差抄了 8 遍）。
v2 把这些高频模式固化成 class：`.page-no` `.kicker` `.rule` `.rows`/`.row`/`.row.line`
`.key` `.num` `.term` `.label` `.sub` `.footnote` `.stat` `.unit` `.grid-3`/`.grid-4`。
每缺一个组件，就会在每一页被重新手写一次——既费 token，又让 8 页之间慢慢长歪。

**组件集 v2 → v3 的由来（版式，不是组件）**：v1+v2 二十多个 class 之后，产出仍然"千页一面"，
因为它们的**形状只有两种**——等分格子和竖排行清单。形状一样，剩下的就只有配色和文字在变，
于是每页轮廓相同、每份 deck 也相同。所以 v3 的挑选标准不是"还缺什么功能"，而是
**画出来的矩形完全不同**：`.split`（不等宽分栏，有主次）`.bleed`（一页一句话）
`.plate`（章节过渡页）`.timeline`（横向序列）`.band`（横向色带）
`.marginal`（正文 + 右侧旁注）。

配套的是 prompt 里的**版式目录**（按"这页要做什么"选版式）与两条硬规矩：同一版式全篇
不超过两次、连续两页不许同版式。**版式多样性不是靠更多 class，而是靠把版式从"默认值"
变成"决策"**——默认值一定趋同，只有显式选择才会分散。

同时 v3 有一条自我约束：**一律不用左侧粗竖条**。彩色侧边条是"一眼认出是 AI 生成"的清单上
排第一的标记，组件库自己不该生产它。`TestRealComponentLibraryContract` 会在
`border-left` 出现在 v3 分段里时直接失败。`.quote` 是唯一保留左侧竖线的组件——它真的是引用块。

### 版式体检：安全边距、字号阶梯、间距节奏

这一节记录三条**量出来**的排版缺陷和它们的修法。起点是把真实渲染的几何量一遍
（`layouts-test.html`，浏览器里读 `offsetLeft/offsetWidth/getComputedStyle`），
结论全部来自数字而不是感觉：

| 缺陷 | 实测 | 规范要求 |
|---|---|---|
| 内容贴边 | `section` 的 padding/margin 都是 `0`，`h2` 的 `x=0`、宽度顶满 960 | 每边留 5~15%（均衡型 10%），投影还会裁掉边缘 5% |
| 字号阶梯塌了 | 一页 8 档，其中 17.7/18.4/19/21.1/21.8px **五档挤在 4px 带宽里**（相邻比例 1.03） | 一个版面 ≤5 档，相邻档必须看得出差别（比例 1.25 附近） |
| 字号低于投影下限 | 卡片正文 21px、次要文字 17.7px ≈ **13.6pt** | 正文 ≥24pt、28~32pt 更好。换算：1 canvas px ≈ 0.77pt |
| 间距只有两档 | `--space-md` 0.6em / `--space-lg` 1.2em | ≥3 档、8px 基准；**组外间距必须大于组内间距** |

后一条的具体症状：卡片间距（40.8px）正好等于卡片自己的内边距（40.8px），
所以三张卡片读起来是一条带子、不是三个独立的块。规范里"间距是比描边更强的分组手段"，
而 `.rows` 原来给每一行画虚线分隔，正是"用描边代替间距"。

**修法**（`theme.css` 定义、`components.css` 只引用）：

- **字号阶梯**：0.7 / 0.92 / 1 / 1.1 / 1.55 / 2.4 em（正文 1em = 34px = 26pt）。
  0.92em 是**投影下限**（31px = 24pt），所以"次要文字"只允许靠色阶区分、不许再缩字号。
  h2 走 reveal 的 2.11em，和上面的台阶合起来共 7 档，但一页通常只用 4 档。
- **间距阶梯**：8/16/24/32/48px 折成 em（`--space-2xs` … `--space-lg`）。
  卡片内边距（24/32）**小于**网格间距（32/48），邻近性才成立。
- **安全边距**：`--slide-pad-x/y` ≈ 6%（默认画布 wide 下 75/41px），内容区约 1094×610。
  必须显式写 `box-sizing: border-box`——reveal.css 里没有全局 `*{box-sizing}`，
  而 section 宽度是画布的 100%，用默认的 content-box 加 padding 会直接撑成 1319px。
  另外 `html.reveal-print`（PDF 导出模式）会把 section 的 padding 强制清零，要单独加回来。
- **`--hairline`**：组内分隔用的发丝线（清单行、旁注顶线），由**正文色**派生而非强调色。
  `.row` 的虚线分隔同时换成实线发丝线——虚线是"表格感"和低质模板的常见特征。
- **`FIT_MARGIN` 32 → 8**（`init.js`）：安全边距现在由 section 自己提供，
  兜底再留 32px 等于白扣 5% 的可用高度。

**连带的连锁反应**（这类改动必须一起想清楚，否则收益会被抵消）：

- 字号抬上去之后，**能放句子的多列版式只剩 `.grid-2`**。默认画布（wide）内容区约 1094px，
  两列每格约 467px 文字宽（12 汉字/行）；三、四列只剩 276/181px，是 7/5 个字一行。
  这是算术结果不是保守：画布宽度固定、字号有下限。所以 `.grid-3`/`.grid-4` 被限定为
  "只放数字与短标签"，`.timeline` 限定 ≤4 格、每格 ≤8 字。
- **默认画布从 standard(960) 换成 wide(1244)**，正是上面那条算术的直接后果：
  standard 在 16:9 屏上还要留黑边，可用于排版的宽度只有 wide 的 77%。
  两者高度都是 700，所以这是纯粹的横向增益，不会把现有页面推进缩兜底。
- **prompt 的容量预算必须跟着重写**：原来是 ≤7 子元素 / ≤12 行 / ≈300 字，
  在新字号下会直接训出溢出的页面。现在是 ≤5 子元素 / ≤200 字，常规演示 8~14 页
  （一页装得少了就多分几页）。**示例页的文字也必须一起压**——示例是模型最强的模仿对象，
  留着旧密度的示例等于把新预算原地推翻。
- **所有已有 deck 的观感会变**：`components.css`/`theme.css` 是共享文件，
  所有 deck 都引用它们（这条正是"改动要落进已被引用的文件"的红利，不需要迁移脚本）。
  但页面的可用面积少了约 20%、字号大了约 40%，原来刚好放满的页会触发收缩兜底。
- **一个没做的**：reveal 的 `@media print`（`html:not(.print-pdf)`）会把 `p/li/td` 强行
  压成 20pt、标题 24pt 的黑白样式，**打印/导出 PDF 时这套字号阶梯会被整个抹平**。
  这是导出的另一个问题，本文档记录在案、尚未处理。

---

## 19. 写入时的样式体检（不是工具，是写工具的 `warning`）☑

**为什么需要**：消毒闸门只拦得住"危险"的东西，拦不住"合法但会让整套 deck 烂掉"的东西。
写死颜色、px 字号、在 `<section>` 上覆盖 `background`/`color`/`font-family`、以及大面积
重复内联——全都不报错、导出也正常，只是这份 deck 从此刻起无法再换风格（`update_theme` 静默失效）。

判定在 `service/deck/stylelint.go`，与消毒走**同一条回报通道**（`warning` 字段），
由 `parseSlideFragment`（update/insert 共用）与 `normalizeSections`（write_deck）各自调用一次。

**五项检测**：

| 检测 | 阈值 | 为什么是问题 |
|---|---|---|
| `<section>` 上写主题级属性 | ≥1 | 整页覆盖主题；`font-size` 还会被适配兜底直接覆盖掉 |
| 写死的颜色值（hex / `rgb()` / `hsl()`） | ≥3 处 | `update_theme` 失效，用户换配色时这里不变 |
| px 字号 | ≥3 处 | 适配兜底靠缩放 section 的 font-size，px 不跟着缩，装不下就被裁 |
| 内联字号（任意单位，含 `font` 简写） | ≥1 处 | 组件库的字号是固定阶梯（0.7/0.92/1.1/1.32/1.6/2.4），手写的字号会插进两档之间，让"哪一层更重要"读不出来 |
| 同一串 style 值重复 | ≥3 次 | 该是个组件，不该是内联 |

（px 字号与内联字号互斥：一条声明只报一次，报更要紧的那条——px 两个毛病都占。
内联字号是唯一"1 处就提示"的项：颜色的偶尔偏离可以是刻意的设计，
字号偏离则是把刚立起来的梯子拆掉，没有"偶尔"这回事。）

**三个关键设计决定**：

1. **跨页聚合**（`styleLinter` 是个累加器，不是纯函数）：页码角标这类"页面家具"在每一页只出现
   一次，逐页检测**永远是"不重复"的**——只有把整份提交放一起数，才看得见"同一串抄了 8 遍"。
   这是这个功能唯一真正的技术含量所在。
2. **归一化后再比**：声明顺序与空白差异不影响判定（模型手写的重复往往就差一个空格）。
3. **只提示、不改写**：不做"自动把重复内联提成 class"——模型下一次 `read_slide` 会读到一份
   它没写过的 HTML，"我写的 = 我读到的"这个自洽性一破，fingerprint 比对和增量修改都会
   开始出现无法解释的差异。代价远高于收益。

**阈值刻意偏保守**：1~2 处写死颜色往往是用户明确要求的刻意偏离（prompt 允许），
只有"成规模"才说明模型抛弃了主题体系。纯组件页面必须完全静默，否则模型很快学会忽略这段提示。
`TestStyleLintSilentOnComponentBasedPage` 和 `TestRealHandWrittenDeckIsFlagged` 一正一反守住这条线。

**手工回归夹具**（浏览器打开即可，不经过 agent）：

- `/assets/components-test.html` —— 组件集 v2 全量渲染验收。最有用的是最后一页
  "对照：组件 vs 原始内联"：同一行内容分别用新组件和 deck-0014 的原始内联写法渲染，
  两边必须完全对齐——这是"组件忠实还原了原效果"的直接判据。
- `/assets/fit-test.html` —— 适配兜底（超载页面等比缩放）的验收页。
- `/assets/layouts-test.html` —— 组件集 v3 的六个版式 + 字号阶梯 + 安全边距的验收页；
  地址栏加 `?preset=noir`（`editorial` / `duotone` / `terminal` / `tech`）可切换预设对比观感。
  页面的文字量是按新容量预算写的（≤150 字/页）——塞旧稿子进去会触发出缩兜底，
  那就分不清"新版式好"还是"被缩小了"。

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
