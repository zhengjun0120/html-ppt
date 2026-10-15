# DeckAgent 后端

AI 幻灯片生成器的 Go 后端：agent 循环、模板体系、deck 存储与账号/账本，单进程 Gin 服务（:8080）。

## 运行

```bash
docker start mysql8                # 本机 MySQL；Redis 可选（验证码注册用）
cd backend && go run ./cmd/server  # :8080
```

- 配置读 `config.yaml`（gitignored，含真实密钥；从任意子目录启动都会向上找到它）。
- 唯一必填 `auth.jwt_secret`（缺失拒绝启动）；LLM key、MySQL 按需，也可走环境变量（如 `LLM_API_KEY`）。
- 健康检查：`curl http://localhost:8080/api/health` → `{"status":"ok","db":"ok"}`

## 分层与目录

依赖单向向下：`handler → service → store / 文件系统`。handler 不写业务，service 不碰 HTTP。

```
cmd/server/          入口：装配依赖 → 启动 → 优雅停机（含模板缩略图预热）
internal/
  agent/             核心：SSE agent 循环、阶段化工具集与提示词、ask_user 暂停协议、
                     视觉审查回环、trace 记录、断线 resume
  handler/           HTTP 翻译层：参数绑定、错误映射（阶段冲突→409 等）、SSE 装配
  router/            URL → handler 映射、静态资源与缓存策略
  middleware/        JWT 鉴权、CORS
  service/
    template/        模板注册表：启动加载+校验、Instantiate 快照实例化、用户模板挂载
    deck/            deck-v2 存储与管线：大纲、写页、节奏校验、AI 味 lint、历史版本
    usertpl/         用户自定义模板：fork / 定制对话 / 发布门禁 / 历史回滚
    auth/            JWT 登录注册、验证码（Redis）、SMTP 邮件
    usage/           用量账本（只增不删，/usage 页聚合）
    usermodel/       用户自选 OpenAI 兼容模型（key 密文不出后端）
    tplthumb/        模板缩略图（chromedp 渲染 + 启动预热）
  vision/            headless Chrome 页面量测与看图审查（一次性 nonce 取页通道）
  export/            PDF / 逐页 PNG / 自包含单文件 HTML
  trace/             run 级观测 JSONL + 内存索引（/trace 观测台数据源）
  store/             GORM：User / UserModel / Deck / ChatSession / UsageEvent / UserTemplate
  cryptox/           AES-GCM 加密盒（BYOK / 自选模型的 key 只落密文）
data/                运行时数据：decks/、user-templates/、template-thumbs/（不入库）
```

## 降级语义（缺什么就少什么，不挡启动）

| 依赖缺失 / 连不上 | 影响 |
|---|---|
| MySQL | 无 DB 模式：只剩 health 可用（其余接口都要求登录） |
| Redis | 验证码注册不可用，登录不受影响 |
| crypto.aes_key | BYOK 与用户自选模型关闭 |
| 模板库加载失败 | 模板相关接口不可用，旧管线照常 |
| Chrome / 回环地址 | 看图审查、导出、缩略图关闭 |

## agent 循环（internal/agent）

- 六阶段状态机 `clarifying → outlining → outline_review → selecting_template → generating → iterating`，
  每阶段独立工具集、提示词与轮数预算（config `deck_v2.max_turns` 可覆盖）。
- 工具失败不终止 run：错误作为 tool result 喂回模型自愈。致命的只有三类：协议破坏、
  事件推送失败（客户端断连）、panic（SSE 入口统一 recover）。
- 防死循环四层闸：阶段轮数预算（耗尽时摘除工具表强制收尾）、单工具配额（如看图审查
  3 次/run）、每页写入预算（write_pages 与 update_slide 合并计数）、ask_user 暂停。
- 质量网：逐页量测（1920×1080 固定画布）→ AI 味 lint → 版式节奏校验 → 看图审查。
- 每条消息即时落库：崩溃最多丢一轮；生成中断 stage 停在 generating，`resume=1`
  对齐已落盘的页续写。

## 接口分组（全表见 router/router.go）

- **公开**：模板清单/详情/preview、`/render/:nonce`（headless 浏览器专用一次性通道）、auth 三件套。
- **JWT 把守**：deck-v2 管线（大纲/模板/生成/导出）、chat SSE 与会话回放、用户模板全生命周期、
  缩略图、trace 读取、用量聚合、模型管理。
- **静态**：`/assets` 回源校验（字体 7 天）、内置模板 1h 强缓存、用户模板受控伺服
  （模板源码与 history 永不外发）。

## 关键约定

- deck / 模板 id 一律白名单正则 `^[a-zA-Z0-9_-]+$`，防路径穿越（有测试用例）。
- 观测 JSONL（滚动删除）与用量账本（只增不删）分离；观测失败只丢记录，不挡生成流。
- 模板是构建产物：改 `template.json` / 模板文件须重启 8080 生效。
