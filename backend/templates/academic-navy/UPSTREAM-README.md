# academic-navy（藏青学衡）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/academic-navy`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：藏青 #112D5C / #173C7A / #1F4F97 为主，浅背景 #F3F5F8；金棕 #7C4E2C、#A78D5B、
  #D4BA54 只用于关键数字、图标与少量标签；正文 #404040 / #808080，深色面板用 #FFFFFF。
- 排版：PingFang SC 中文、SF Pro Text 西文数字；标题 38-40px、章节数字 60-96px、正文 18-23px；
  一页最多一种标题强调方式。
- 布局：左右分区、四宫格、编号环形列表、论文式证据链；装饰用规整的圆形、六边形、菱形、盾形
  与低调线性图标，只承担分组和导向，不引入具象场景。
- 禁忌：糖果色、手绘、卡通、大面积装饰、滥用金色、多字体混排。

## 转换说明

- 原上游只有一个 177.78vh 等比的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 预览页右下的对角藏青色域收敛为「浅灰蓝底 + 藏青标题」的整体气质；藏青色域留给收尾按钮
  （an-btn）这一类小面积身份元素，避免整页深色压过投影环境。
- 金棕三档按用途分工：#7C4E2C 题签小字、#A78D5B 关键数字、#D4BA54 细线圆环——对应上游
  「金棕只用于关键数字、图标与少量标签」的约束。
- 右上金环与左下六边形做成 CSS data-URI SVG 的 ambient 层（.slide::before/::after，
  z-index:-1），每页自动衬底、不进骨架，生成侧零成本。
- 字体栈：'SF Pro Text'/'PingFang SC' 无衬线主体，Menlo/Consolas 等宽做题签与时间点；
  不新增 webfont 文件。
