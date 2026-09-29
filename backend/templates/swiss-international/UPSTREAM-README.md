# swiss-international（瑞士国际主义）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/swiss-international`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：背景 #FAFAF8 暖白，灰阶 #F0F0EE / #D4D4D2 / #737373，文字 #0A0A0A；强调色可选
  IKB 蓝 #002FA7、柠檬黄、柠檬绿或安全橙——一个 deck 只能使用一种强调色（本模板取 IKB 蓝）。
- 排版：Inter / Helvetica Neue / Noto Sans SC，禁衬线；主标题与正文比例至少 8:1，
  越大越细（200-300 字重）、越小越粗；正文 ≥18px、标签 ≥14px。
- 布局：元素吸附 12/16 栏网格，左对齐 + 大幅留白的不对称美学；色块全部直角、无阴影、
  无圆角；并列卡片只允许一个强调焦点。
- 装饰：只允许 1px hairline 分割线、8x8 直角方块与极细网格/点阵；禁圆点、渐变、阴影。
- 场景：科技产品发布、数据汇报、设计工程、年度总结、AI 技术分享。

## 转换说明

- 原上游只有一个 16:9 单页预览（细网格纸底 + 细体巨字 + 黑缝 KPI 拼格的封面构图）；
  本模板按本仓库 deck-v2 契约重建为 9 版式骨架
  （cover/contents/keynotes/split/metrics/quote/divider/moments/closing），无 code 版式。
- 封面的黑缝拼格（1px 缝隙当分隔线）、mono 元信息行与 8x8 直角方块取自上游预览母题；
  每页的 5% 极细灰网格做成 .slide::before ambient 层（z-index:-1，linear-gradient 画线）。
- h1/h2 锁定 200 字重细体：Inter 可变字重覆盖拉丁，中文回退到 Noto Sans SC 400，
  观感仍是「细」——上游「越大越细」的气质在 CJK 下同样成立。
- 与上游同名的「swiss-grid」风格（外露网格线、红色大字、近黑色块）在本仓库另建
  swiss-grid 模板承接；本模板保持克制理性路线，两者刻意区分。
- 字体栈：Inter / Helvetica Neue / Noto Sans SC，元数据用 JetBrains Mono；
  不新增 webfont 文件。
