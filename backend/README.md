# Deck Agent 后端（骨架阶段）

交互式 HTML PPT 生成 agent 的 Go 后端。当前是阶段 0 骨架：
配置 / 路由 / 中间件 / 分层结构已就位，deck 预览接口可用，chat 留待阶段 1 实现。

## 运行

```bash
# 1. 确保 MySQL 容器在跑（本机已有 mysql8 容器，映射 3306）
docker start mysql8

# 2. 填配置：复制 config.yaml（含本机数据库密码，已被 .gitignore 排除）
#    api_key 阶段 1 才需要，也可用环境变量 LLM_API_KEY

# 3. 启动
go run ./cmd/server
#    或者
go build -o backend-server.exe ./cmd/server && ./backend-server.exe
```

验证：

- `curl http://localhost:8080/api/health` → `{"status":"ok","db":"ok"}`
- `curl http://localhost:8080/api/decks` → 列出 deck
- 浏览器打开 `http://localhost:8080/api/decks/deck-0001/file` → 预览初始 deck

## 接口一览

| 方法 | 路径 | 说明 | 状态 |
|---|---|---|---|
| GET | /api/health | 健康检查（含数据库状态） | ✅ |
| GET | /api/decks | deck 列表（id + 标题） | ✅ |
| GET | /api/decks/:id/file | 返回 deck.html 给 iframe 预览 | ✅ |
| GET | /assets/* | reveal.js / 主题 / 组件库静态资源 | ✅ |
| POST | /api/chat | agent 循环入口（SSE） | 阶段1 |

## 目录结构与分层

```
backend/
├── cmd/server/main.go        入口：装配依赖 → 启动 → 优雅停机
├── internal/                 internal/ 下的包禁止被外部项目 import
│   ├── config/               配置加载：默认值 < config.yaml < 环境变量
│   ├── router/               路由注册：URL → handler 的映射都在这一处
│   ├── middleware/           横切关注点：CORS（后续加日志、鉴权）
│   ├── handler/              HTTP 翻译层：解析请求、调 service、写响应
│   ├── service/deck/         业务逻辑：deck 文件的校验、读取（写入在阶段2）
│   └── store/                数据库层：GORM 连接 + 模型（User/ChatSession 占位）
├── web/assets/               reveal.js + theme.css(变量层) + components.css(结构层)
├── data/decks/<id>/deck.html deck 数据文件（运行时数据，不入库）
└── config.yaml               本地配置（.gitignore 排除）
```

依赖方向单向向下：`handler → service → store / 文件系统`。
handler 不写业务，service 不碰 HTTP，store 不关心谁调用它。

## 已确立的关键约定（和讨论结论一致）

- **deck 文件三层**：框架/脚本（不可改）、组件库（只可用）、`AGENT-EDITABLE` 注释标记之间的 sections（LLM 唯一可写区）
- **id 是唯一外部标识**：白名单正则 `^[a-zA-Z0-9_-]+$` 校验，防路径穿越（有测试用例）
- **组件库分两层**：结构层引用变量、主题层定义变量 → 换主题不换结构
- **数据库可降级**：`db.enabled: true` 但连不上时以无 DB 模式运行并告警

## 路线图

- [x] 阶段0：骨架 + deck 预览
- [ ] 阶段1：最小 agent 循环（go-openai + DeepSeek + SSE 流式）
- [ ] 阶段2：页级工具 + 缩略图 + ask_user 提问卡片
- [ ] 阶段3：Playwright 截图自查环
- [ ] 阶段4：主题系统 + 版本快照 + 导出单文件
- [ ] 阶段5：用户系统 + 免费额度 + BYOK

## 待办备忘

- 上线前：`gin.SetMode(gin.ReleaseMode)`（启动日志里有提示）
- 将来开源时：module 名 `html-ppt/backend` 改成 `github.com/<你>/xxx`，
  `go mod edit -module ...` 后全局替换 import 路径即可
