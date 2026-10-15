# DeckAgent · AI 生成网页幻灯片

用户描述需求 → agent 澄清并产出**结构化大纲** → 用户挑选**模板** → agent 按模板逐页生成 → 对话式迭代 → 导出 PDF / PNG / 单文件 HTML。

```
frontend/            Vue 3 + Vite + Pinia + Tailwind 工作台（:5173）
backend/             Go (Gin) 单进程服务（:8080）
  internal/agent/    分阶段 agent 循环（SSE）、阶段化工具集、ask_user 暂停协议
  internal/service/  模板注册表 / deck-v2 存储 / 用户模板 / 用量账本 / 自选模型
  internal/vision/   headless Chrome 页面量测与看图审查
  internal/export/   PDF / 逐页 PNG / 单文件 HTML
  internal/trace/    run 级观测（/trace 观测台）
templates/           内置模板库（100+ 套 × 主题变体）
web/assets/deck-v2/  模板共享运行时（vendor 自 html-ppt-skill，MIT）
```

## 核心机制

- **模板承载质量**：LLM 只准使用模板登记的版式与类名（版式锁 + 类名契约，服务端机械校验），
  设计责任从 prompt 转移到模板；systemPrompt 拆为五份阶段提示词。
- **阶段状态机**：`clarifying → outlining → outline_review → selecting_template → generating → iterating`，
  每阶段独立工具集、提示词与轮数预算；断线可 resume 续跑。
- **两道用户闸门**：大纲确认（面板直改 / 对话让 agent 改，双通道同一份 JSON）→ 模板必选
  （画廊 live 预览 + 主题变体 + AI 推荐）。
- **分批生成 + 质量网**：plan_pages 全局规划 → 每批 2-4 页写入 → 每页固定画布量测
  （1920×1080，所见即所得）→ AI 味 lint（中文禁词 / 假精确数字）→ 版式节奏校验 →
  看图审查（硬配额防循环）。
- **用户自定义模板**：fork 内置模板 → 定制对话 / 在线编辑器改稿 → 发布门禁 → 社区画廊。
- **版本与观测**：deck / 用户模板全量历史一键回滚；`/trace` 逐 run 留痕（LLM 请求、
  工具调用、审查截图）；`/usage` 用量账本；支持 BYOK 与用户自选 OpenAI 兼容模型
  （key 密文不出后端）。

## 快速开始

```bash
# 1. 配置 backend/config.yaml（LLM key、MySQL、Redis；文件内有注释）
cd backend && go run ./cmd/server            # :8080
# 2. 前端（另开终端）
cd frontend && npm install && npm run dev    # :5173，代理到 8080
```

浏览器打开 http://localhost:5173 → 注册/登录 → 「新文稿」→ 描述需求，跟着向导走。

生产部署：nginx 伺服 `frontend/dist`（构建产物挂在 `/app/`，避开 deck-v2 运行时的
`/assets` 命名空间）+ 反代 `/api` 到 Go 服务。

## 文档

- `docs/refactor-plan.md` — deck-v2 架构与决策记录（D1-D14）
- `THIRD_PARTY_NOTICES.md` — 第三方组件与致谢（html-ppt-skill MIT / OFL 字体等）
- `backend/templates/tools/` — 模板入库流水线（scaffold + materialize + validate）

## License 与贡献

见 `THIRD_PARTY_NOTICES.md` 的边界声明；向 `backend/templates/` 贡献模板时
**禁止引入 AGPL 代码**（详见该文件 guizang 一节）。
