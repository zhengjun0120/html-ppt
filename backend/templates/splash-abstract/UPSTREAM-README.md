# splash-abstract（泼彩抽象·多巴胺）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/splash-abstract`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：橙红 `#FF4500`、热品红 `#FF1493`、紫罗兰 `#9370DB` 三主色，青绿 `#20B2AA`、
  亮蓝 `#1E90FF` 副色，纯白/暖白 `#FFF8F0` 底，深墨 `#1A1A2E`、暗紫 `#2D1B69` 文字，
  明黄 `#FFD700` 强调。五色大胆泼洒，如把彩虹泼进画面。
- 意象：泼墨色块、彩虹、爆炸星芒、流动液态、不规则形状、彩色圆点、粗笔触；
  色块可叠加混色溢出边界。
- 排版：粗体无衬线（Noto Sans SC Black / Helvetica Bold）超大标题；信息卡彩色底白字；
  不对称布局有动感；正文信息保持低到中密度，「不要用彩色卡片堆满页面」。
- 禁忌：低饱和沉闷暗淡、严肃商务、古典克制、小字密集排版、单一主色。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 上源的五团 splash 色块与星芒全部收进 ambient 层：`.slide::before` 的左上橙红、右上品红
  两组径向泼点，`.slide::after` 的页底紫罗兰/青绿泼块、明黄八角星芒与飞溅圆点
  （data-URI SVG）。每页自动衬底（z-index:-1）、不进骨架——上游「文字叠加泼彩」改为
  「泼彩衬在边缘、正文区留暖白」，避免色彩抢正文。
- **对比度改造**：上源统计卡的「白字小标签压彩色底」不满足 WCAG，本模板改为
  「粗黑描边白卡 + 深墨文字 + 五色只上大数字（特大号 ≥3:1）」；彩色编号块/时间线圆点
  上的白字限特大号等宽数字；按钮深墨底白字（16:1）。泼彩浓度保留在装饰层。
- 上源的粗描边引文卡（3px 深墨描边）保留为全模板统一的卡片画法（sp-card / sp-stat），
  硬阴影用深墨 14% 或明黄做爆点。
- 字体栈：'Noto Sans SC' 特粗标题，数字 'Helvetica Neue'/JetBrains Mono；
  不新增 webfont 文件。
