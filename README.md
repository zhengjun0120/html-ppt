# DeckAgent · AI 生成网页幻灯片

一句话：用户描述需求 → agent 澄清并产出**结构化大纲** → 用户挑选**模板** → agent 按模板逐页生成 → 对话式迭代 → 导出 PDF/PNG/单文件 HTML。

本项目是一个全栈的开源 AI PPT 生成器（deck-v2 架构）：

```
frontend/   Vue 3 + Vite + Pinia 工作台（向导式生成流程 + 对话迭代 + 预览/历史/观测台）
backend/    Go 单进程服务
  cmd/server/         入口（:8080）
  internal/agent/     agent 循环、分阶段提示词（prompts/）、阶段化工具集、ask_user 暂停协议
  internal/service/template/   模板注册表/实例化/类名清单（8 个内置模板）
  internal/service/deck/       deck-v2 存储（deck.json/outline.json/三件套快照）、
                               版式节奏校验、AI 味 lint、类名契约
  internal/vision/    固定画布量测 + 看图审查 + PrintToPDF
  internal/export/    PDF / 逐页 PNG zip / 单文件 HTML 导出
  internal/trace/     run 级观测（/trace 观测台）
  templates/          模板库（目录契约见 templates/CONTRIBUTING 计划与 docs/refactor-plan.md §5）
  web/assets/deck-v2/ 模板共享运行时（vendor 自 html-ppt-skill，MIT，补丁见 PATCHES.md）
```

## 核心机制

- **模板承载质量**：LLM 只准使用模板登记的版式与类名（版式锁 + 类名契约，服务端机械校验），
  设计责任从 prompt 转移到模板——systemPrompt 从 667 行瘦身为五份阶段提示词。
- **两道用户闸门**：大纲确认（面板直改 / 对话让 agent 改，双通道同一份 JSON）→ 模板必选
  （画廊 live 预览 + 主题变体）。
- **分批生成 + 三层质量网**：plan_pages 全局规划 → 每批 2-4 页写入 → 每页自动量测
  （固定画布 1920×1080，所见即所得）→ AI 味 lint（中文禁词/破折号/假精确数字）→
  版式节奏校验 → 看图审查（硬配额防死循环）。
- **阶段状态机**：`outlining → outline_review → selecting_template → generating → iterating`，
  每阶段独立工具集与提示词；断线可 resume。

## 快速开始

```bash
# 1. 配置 backend/config.yaml（LLM key、MySQL、Redis；均已在文件内有注释）
# 2. 启动后端
cd backend && go run ./cmd/server
# 3. 启动前端
cd frontend && npm install && npm run dev   # http://localhost:5173
```

浏览器打开 http://localhost:5173 → 注册/登录 → 「新文稿」→ 描述需求 → 跟着向导走。

## 文档

- `docs/refactor-plan.md` — deck-v2 重构执行蓝图（架构、契约、决策记录 D1-D14）
- `THIRD_PARTY_NOTICES.md` — 第三方组件与致谢（html-ppt-skill MIT / taste-skill MIT /
  guizang-ppt-skill 仅思想参考 / OFL 字体）
- `backend/web/assets/deck-v2/PATCHES.md` — vendor 运行时的补丁记录
- `backend/templates/tools/` — 模板入库流水线（scaffold + materialize + validate）

## License 与贡献

见 `THIRD_PARTY_NOTICES.md` 的边界声明；向 `backend/templates/` 贡献模板时
**禁止引入 AGPL 代码**（详见该文件 guizang 一节）。
