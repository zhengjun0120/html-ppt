# 第三方组件与致谢（THIRD_PARTY_NOTICES）

本项目（deck-v2 重构后）引用或借鉴了以下开源项目，在此一并致谢。

## 1. html-ppt-skill — MIT

- 仓库：https://github.com/（本地参照副本：html-ppt-skill-main）
- 使用方式：**vendor 并有改动**
  - `backend/web/assets/deck-v2/runtime.js` / `base.css` / `animations.css` 来自其
    `assets/` 目录，改动逐条记录于 `backend/web/assets/deck-v2/PATCHES.md`；
  - `backend/templates/` 下 8 个模板由其 `templates/full-decks/` 对应模板适配而来，
    每个模板目录内的 `ADAPTATION.md` 记录了完整改动清单；
  - `backend/templates/<id>/style.css` 末尾的主题变体槽（variants）为本项目新增。
- 原 License（MIT）随源码保留；感谢原作者。

## 2. taste-skill — MIT（思想借鉴，未复制代码）

- 仓库：https://github.com/（tasteskill.dev）
- 使用方式：其「AI 味禁令」（§9 AI Tells）与文案密度纪律被**中文化重写**进
  `backend/internal/service/deck/lint_taste.go` 的 T001–T007 规则与
  `backend/internal/agent/prompts/shared.md` 的文案纪律。
  词表为本项目维护的中文版，非逐词翻译。

## 3. guizang-ppt-skill — AGPL-3.0（仅设计思想参考，零代码引入）

- 仓库：https://github.com/op7418/guizang-ppt-skill
- 边界声明：**本项目未复制其任何文件或代码**。参考吸收的思想包括：
  版式锁（每个 section 必须声明 data-layout 且仅用登记版式）、主题节奏规则
  （同版式不连用 3 页、8 页窗口内版式多样性）、checklist 分级（阻塞/提示两级）。
  以上思想以本项目原创代码实现于 `internal/service/deck/v2_rhythm.go`、
  `internal/service/template/` 与各模板 `layouts.md`。
- **贡献者注意**：禁止将 guizang 仓库的任何代码或模板文件复制进本项目（包括
  `backend/templates/`），避免 AGPL-3.0 传染。设计思想不受版权保护，但表达受。

## 4. 字体（均随附 OFL-1.1 授权文本）

- Inter（rsms/inter，拉丁可变字重）
- Noto Sans SC（Google，简体子集）
- JetBrains Maple Mono（SpaceTimee/Fusion-JetBrainsMapleMono，等宽含 CJK）

授权全文见 `backend/web/assets/deck-v2/fonts/` 各 `*OFL.txt`。

## 5. 其他关键依赖

- Go 生态：gin / gorm / chromedp / goquery / openai-go / glebarez(sqlite) 等，
  以 `backend/go.mod` 为准。
- 前端：Vue 3 / Vite / Pinia / Tailwind CSS v4 / Phosphor Icons，以
  `frontend/package.json` 为准。
