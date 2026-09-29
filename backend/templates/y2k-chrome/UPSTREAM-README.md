# y2k-chrome（Y2K 铬）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/y2k-chrome`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：银色四段渐变背景 `linear-gradient(135deg,#e5e5e5,#f5f5f5,#d4d4d4,#a8a8a8)` 模拟铬金属；
  标题 #171717 墨黑、正文 #404040 深灰；#a855f7 紫与 #06b6d4 青是铬面的彩虹反射。
- 意象：银铬反射、彩虹光斑、气泡；大圆角是标志——卡片、按钮、一切可以圆的都圆，
  半透明卡像浮在液态金属上的气泡。
- 排版：现代无衬线、字重轻到中等；全大写标题配宽字距是经典手法。
- 禁忌：不要传统或暗色底；不要小圆角或直角；不要丢掉银色铬金属的核心质感；
  彩虹光斑不做乱，保持光泽感。

## 转换说明

- 原上游只有单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约重建为 9 版式骨架
  （cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 文本级紫/青收深一档（#a855f7 → #5b21b6）：亮紫在银底上对比度不足 4.5:1，
  收暗后用于题签与时间点；亮紫青只保留在渐变按钮、彩虹条与描边等大面积元素上。
- 铬面高光斜扫（.slide::before）与气泡/光斑群（.slide::after，data-URI SVG）做成 ambient 层
  （z-index:-1），每页自动衬底、不进骨架，生成侧零成本。
- 字体栈：Inter + Noto Sans SC（无衬线），数字与标签 JetBrains Mono（等宽，本仓库自托管字体）；
  不新增 webfont 文件。
- letterbox 底色 --bg 取银铬渐变的中段值 #d6d6d8，翻页不跳色。
