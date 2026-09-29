# corporate-clean（企业洁净）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/corporate-clean`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：纯白到蓝灰的近纯白底（#FFFFFF → #F8FAFC → #F1F5F9），深墨 #0F172A 标题、
  石板灰 #334155 正文，海军蓝 #1E3A8A 与蓝 #3B82F6 组成企业级双色系统；
  卡片 rgba(255,255,255,0.96)，边框是极淡的蓝色描边 rgba(30,58,138,0.15)。
- 排版：Inter 或专业无衬线体，标题偏粗、正文偏中，字号层级分明。
- 布局：规整卡片式布局，边界清晰、对齐严格、间距均匀，默认低到中密度。
- 禁忌：花哨的动画或配色、暗色背景、圆角过大或阴影过重的卡片、任何「随意」的元素。
- 场景：董事会汇报、B2B 销售演示、金融保险行业报告、季度业务 Review。

## 转换说明

- 上游只有一个单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 上游封面左侧的深色渐变侧栏转写为页顶 5px 海军蓝信头细线：保留「信头」身份、
  避免每页出现大面积深色块（上游自己规定「不要使用暗色背景」）。
- 右下三道账目细线与信头线做成 .slide::before/::after 的 CSS 层（z-index:-1），
  每页自动衬底、不进骨架，生成侧零成本。
- 字体栈：'Inter' / 'Noto Sans SC'（自托管 webfont），不新增字体文件。
