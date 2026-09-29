# oh-my-ppt 84 风格导入为内置模板 —— 施工方案

> 状态：**全部 83 个模板已落地（2026-09-30）**。试点批 12（ef44d45）+ 波 1 24（ebe4cf4）+
> 波 2 26（0dd38f2）+ 波 3 21（f3a341c）。84 风格 - soft-pastel（与既有内置模板同名，拍板跳过）= 83。
> 全部通过：tmp-check 自检 + 注册表 go test（含 TestOMPTemplateCatalog 目录回归）+ Chrome
> 真渲染量测（demo 页填充率/字号/溢出）+ 封面拼图目检。
> 生成质量测试（每模板真实 LLM 生成 12 页统一大纲《二十四节气》）单独报告：
> docs/ohmyppt-import-report.md。
>
> 决策记录（2026-09-29 拍板，全部采纳推荐项）：D1=先试点 12 个；D2=跳过上游 soft-pastel；
> D3=每模板单 default 变体；D4=全中文 demo。2026-09-30 追加：修复预算=修 1 轮重测 1 次；
> 生成主题=统一大纲保证横向可比。
>
> 来源仓库 `D:\go_files\oh-my-ppt\oh-my-ppt`（GitHub arcsin1/oh-my-ppt，
> Apache-2.0，NOTICE 要求保留归属）。

## 实施备注（试点批沉淀的约定，后续批次沿用）

- 统一 9 版式骨架：cover(hero)/contents(table,4行)/keynotes(cards,3卡)/split(split,3步)/
  metrics(chart,3数)/quote(quote)/divider(hero)/moments(chart,4节点)/closing(hero)。
- 类前缀分配：ink-wash-jiangnan=jw、palace-ink-red=pi、chinese-porcelain-rose=pr、
  chinese-cream-blossom=cb、song-rain-poetic=sr、chinese-ink-landscape=cl、
  celadon-bamboo=qc、chinese-fresh-trio=ft、chinese-pastel-spring=ps、indigo-lotus=il、
  oriental-poetic-illustration=oi、gold-ivory=gi。
- ambient 装饰统一走 .slide::before/::after + data-URI SVG（z-index:-1，画在 .deck
  背景之上、内容之下，骨架零成本）；身份色纪律写进各模板 rules.md。
- 字体映射按 §4 执行，未新增 webfont。
- 验收管线：`go test ./internal/service/template/ -run TestLoadRepoTemplates` +
  `go run ./cmd/tplcheck <ids...>`（临时量测工具，全量完成后删除）。

---

## 1. 来源盘点（已核实）

`resources/styles/` 下 **84 个风格目录** + 1 个 manifest.json（v2.0.8, author arcsin1）。
每个风格恰好 3 个文件：

| 文件 | 内容 |
|------|------|
| `style.json` | 元数据：style id、中英文名、description、category（如「学术 · 严谨」）、aliases（搜索别名）、styleCase（适用场景一句）、version、source |
| `SKILL.md` | 风格规范提示词：概述 / 配色（精确色值）/ 排版（字号范围）/ 布局 / 插画与装饰 / 动画 / 适合场景 / 配图 / 不要 |
| `preview.html` | 单页 16:9 预览，内联 CSS，自包含 |

已核查的安全性：84 个 preview.html 全部无外部资源引用（唯一命中是
chinese-porcelain-rose 的 `url(#vaseG)`——SVG 内部引用，不是外链）。

风格谱系粗分：中文学韵（chinese-\* ×7、ink-wash-jiangnan、palace-ink-red、song-rain-poetic 等）、
程序员/终端（terminal-green、tokyo-night、dracula、gruvbox-dark、nord、rose-pine、
catppuccin-\* ×2、solarized-light、sharp-mono、cyberpunk-neon、vaporwave、y2k-chrome 等）、
学术商务（academic-\* ×2、corporate-clean、swiss-\* ×2、blue-orange-analytics 等）、
手绘治愈（hand-drawn-autumn、handdrawn-watercolor、macaron-mist、sakura-soft-healing 等）、
杂志编辑（editorial-serif、e-ink-editorial、magazine-bold、news-broadcast 等）。

## 2. 本质差异：这是「转换 + 再创作」，不是文件搬运

| | oh-my-ppt 的「风格」 | 我们的「模板」 |
|---|---|---|
| 形态 | 一段给 LLM 的风格提示词 + 一张单页预览 | 机器可校验的多版式骨架包 |
| 生成方式 | LLM 拿 SKILL.md 自由发挥，无结构契约 | data-layout 版式契约 + layouts.md + 渲染量测 |
| 页面结构 | **没有**（预览仅 1 页） | 6-14 个版式骨架 + demo 页 |

因此 `SKILL.md + preview.html` 只是「设计 DNA」（配色、字体气质、装饰语言、禁忌清单），
多版式骨架要按我们的体系重建。这一步是本工程的主要工作量，也是质量所在。

## 3. 每风格产出 5 文件的转换规则

| 产物 | 来源 | 规则 |
|------|------|------|
| `template.json` | style.json | id 沿用（soft-pastel 冲突见 §7-D2）；name=中文名；description=原 description+styleCase 合成；scenario=styleCase+aliases 精选 ≤4；tags=category 拆分；canvas=1920×1080；variants 见 §7-D3；`source={derived_from:"oh-my-ppt/resources/styles/<id>", license:"Apache-2.0"}` |
| `style.css` | preview.html + SKILL.md 配色/装饰节 | 提取色板/圆角/描边/阴影/装饰母题 → `.tpl-<id>` 作用域完整设计系统（CSS 变量 + 版式类 + ambient 背景，遵守我们 base.css/html+body 约定） |
| `index.html` | 版式骨架 + SKILL.md 场景 | `SLIDES:START/END` 标记 + §4 版式集的 demo section；demo 文案从 styleCase 场景取材（全中文，见 §7-D4） |
| `layouts.md` | 版式集 + SKILL.md | 每版式 use/constraints；「不要」清单折进各版式约束 |
| `rules.md` | SKILL.md 排版/不要节 | 字号范围、禁用项 → 生成侧质量规则 |
| `UPSTREAM-README.md` | — | 归属声明：源自 arcsin1/oh-my-ppt（Apache-2.0）、SKILL.md 要点摘录、转换说明（满足 NOTICE 义务） |

**字体映射（不新增字体文件）**——fonts.css 已注册字体仅 `Noto Sans SC` / `Inter` /
`JetBrains Mono`：

- 无衬线 → `'Inter','Noto Sans SC',sans-serif`
- 等宽/终端 → `'JetBrains Mono',monospace`
- 衬线风（academic-paper、editorial-serif 等）→ `Georgia,'Times New Roman','Noto Sans SC',serif`（系统回退）
- 手写风 → 无手写 webfont，以圆角/微旋转/笔触边框模拟气质，字体回退 Noto Sans SC

SKILL.md 里的「配图」节忽略（我们模板体系不涉及）；「动画」节与 base.css 动画类对齐，
时长遵守 ≤650ms 截图管线约束。

## 4. 标准版式集（8 版式覆盖 6 角色）

| 版式 | 角色 | 说明 |
|------|------|------|
| cover | cover | 封面：标题 + lede + kicker |
| toc | toc | 目录/议程 |
| cards | content | 三/四卡要点 |
| split | content | 左右分区（图文/论述） |
| metrics | data | 2-4 指标大数字 |
| quote | quote | 金句页（填充率豁免指纹） |
| divider | divider | 章节过渡 |
| thanks | thanks + cta | 收尾致谢 + 行动点 |

- **程序员审美主题**（终端/编辑器配色 10 个左右）附加 **code 版式**（role=code，
  终端窗口母题）——这些风格的天然卖点。
- 每风格允许 ±1 个个性版式（SKILL.md 布局节强烈指向时，如 news-broadcast 的头条横幅），
  其余保持标准集——84 个模板的结构一致性靠这条保底。
- 角色覆盖校验：每模板 roles 并集 ⊇ {cover, content, thanks}，且含 data 或 quote 其一。
- 版式 HTML 结构以现有 19 个内置模板的实现为参照（grid/flex、anim-\*、deck-footer 等约定）。

## 5. 质量验收管线（每模板必过，全部复用现有设施）

1. **注册表校验**：backend 启动即验 template.json/结构契约，不过不注册（现有铁律）。
2. **渲染量测**：Chrome 截图 + evaluateRender 同阈值——填充率 ≥45%（hero/quote 指纹豁免）、
   最小字号 ≥13px、无溢出。批量脚本跑全部 demo 页。
3. **观感抽查**：visual-judge 看每批截图，风格还原度（对照 preview.html）打分。
4. 每批抽 1 个模板跑一次真实生成冒烟（消耗额度，授权范围内）确认契约可用。

## 6. 批量计划（84 个 ≈ 7 批，每批一个 commit 可回滚）

| 批 | 主题 | 约数量 |
|----|------|--------|
| 1（试点） | 中文学韵 | 12 |
| 2 | 程序员/终端/编辑器 | 12 |
| 3 | 暗色霓虹（cyberpunk/vaporwave/y2k/tokyo…） | 12 |
| 4 | 学术商务 | 12 |
| 5 | 手绘治愈 | 12 |
| 6 | 杂志编辑/极简 | 12 |
| 7 | 其余杂项收尾 | 12 |

每批流程：生成 5 文件 → 注册表校验 → 批量量测 → 截图观感 → 修复 → 提交。
试点批（第一批）定调版式选型与皮肤质量，作为后续 6 批的地基。
改动仅新增 `backend/templates/<id>/` 目录 + backend checkout 同步，不碰任何现有代码路径。

## 7. 决策点（请拍板）

- **D1 范围与节奏**：
  - **A. 先试点一批 12 个，你过目后再跑完剩余 72 个（推荐）**——第一批定调，避免方向跑偏后返工 84 个
  - B. 一次性全量 84 个跑完再过目
  - C. 只挑子集（请指定风格）
- **D2 soft-pastel 同名冲突**（两边都有该 id）：
  - **A. 跳过他们的，保留我们已有的（推荐）**——同为淡彩定位，视觉撞车
  - B. 改 id 导入为 `omp-soft-pastel` 共存
- **D3 variants 槽**：
  - **A. 每模板单 default 变体（推荐）**——84 个风格本身就是配色差异，再设变体意义小
  - B. 每模板 2-3 个变体（工作量约 ×1.5，84 个里多数风格撑不起变体）
- **D4 demo 文案语言**：
  - **A. 全中文 demo（推荐）**——画廊与生成侧均为中文场景，SKILL.md 场景本身就是中文
  - B. 中英混排（标题英文 + 正文中文，某些 editor/terminal 风更有味道）

## 8. 风险与不计入项

- **版式同质化**：84 个模板共享标准版式集，观感靠皮肤区分——与上游定位一致
  （他们本就没有结构差异，靠提示词变化），可接受；个性版式出口（§4）兜住强风格。
- **视觉还原偏差**：preview.html 仅单页，版式重建时视觉细节以量测 + 截图对照验收兜底。
- **不计入**：上游的 PPTX 导入、AI 配图、动画编辑系统（超出模板范畴）。
- **版权**：Apache-2.0 允许衍生；以 UPSTREAM-README.md + template.json.source 保留归属，
  满足其 NOTICE 要求；不修改其仓库任何文件。
