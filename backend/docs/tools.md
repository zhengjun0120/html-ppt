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
| 8 | `update_theme` | 2 | 改主题配置（颜色/字体/动画） | ☐ |
| 9 | `list_templates` | 2 | 列出模板供推荐 | ☐ |
| 10 | `ask_user` | 2 | 向用户提问（human-in-the-loop，暂停循环） | ☑ |
| 11 | `screenshot_slides` | 3 | 截图/程序化检查排版 | ☐ |
| 12 | `move_slide` | 2+ 可选 | 调整页序 | ☐ |

实现顺序：1 → 3/4/5 → 6/7 → 8/9 → 10 →（阶段3）11 → 12。

## 通用实现骨架

- 每个工具 = 名字 + 参数 JSON Schema（给 LLM 的菜单）+ handler（普通 Go 函数），统一注册表分发
- handler 统一签名：`(ctx, argsJSON) → (resultString, error)`，返回给 LLM 的永远是字符串
- **错误即反馈**：校验失败时返回人类可读的错误说明，LLM 会照着自我修正，不算失败
- 校验永远先行：id 白名单 → 结构合法性 → 业务规则
- 读工具尽量省 token（list/read 分开），写工具尽量严格（校验 + 原子写）
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
  data-id 与 slide_id 一致、禁 script/style、禁嵌套 section）→ 替换写回
- **要点**：goquery 定位旧节点 → `ReplaceWithHtml` → 整份文档重序列化 → 原子写回；
  成功返回刻意不含新指纹（同一页再改必须重新 read，防止"旧内容+新指纹"的二次提交
  冲掉上一次修改）；写前快照（阶段4完善）

## 6. insert_slide

- **参数**：`{ deck_id, after_slide_id, new_html }`（`after_slide_id: "end"` 表示追加末尾）
- **返回**：`{"slide_id": "s7"}`——**id 由后端分配（现有最大编号+1，永不复用）**
- **要点**：剥掉 LLM 写的 data-id 重新赋值；`AfterHtml` 插入

## 7. delete_slide

- **参数**：`{ deck_id, slide_id }`
- **要点**：定位 → Remove → 写回；拒绝删空整个 deck（至少留一页）

## 8. update_theme

- **参数**：`{ deck_id, primary_color?, background_color?, transition?, font? }`（只传要改的）
- **要点**：参数是语义化字段而非 CSS；真身是 deck.html 里的
  `<script type="application/json" id="deck-theme">` 配置块，读 JSON → 合并 → 写回 → 同步 CSS 变量；
  校验颜色格式正则 + transition 枚举（slide/fade/zoom/convex/concave）

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
