# chinese-porcelain-rose（中国传统色·凝脂杨妃）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/chinese-porcelain-rose`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：凝脂 #F5F2E9（暖米白）与杨妃 #F091A0（沉静玫瑰粉）双主色；辅助淡杨妃 #F5BFC8、
  暖象牙 #FAF7F0；文字墨黛 #333333、暖灰 #666666；强调胭脂深红 #B23A48、淡金 #D4A574；中性宣纸米 #E8DCC8。
- 排版：标题 Noto Serif SC/宋体加粗，色名用书法体或楷体强调传统命名，正文宋体行高 1.8，hex 用小号无衬线。
- 意象：瓷器（瓶、碗、盏）、缠枝纹、云纹、莲花、水彩晕染；装饰水墨淡彩、细线、印章、回纹、留白；曲线柔和，避免锋利几何。
- 布局：对称或居中，顶部品牌名、中部大标题、下部色名 + hex 纵向排列（上游预览的"双联色卡"）。
- 禁忌：高饱和荧光、赛博霓虹、暗黑背景、锋利金属、拥挤排版、过度装饰、无衬线粗体与传统色名混搭。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（双联色卡 + 瓷瓶插画）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 预览顶部的 8px 渐变色带原样保留为 ambient 层（.slide::before，z-index:-1）；右下瓷瓶插画
  简化为缠枝柔线 data-URI SVG（.slide::after）——具体器物插画交给生成侧配图，模板只留"一层釉"的底纹。
- 预览的"双联色卡"（色名楷体大字 + 拼音 + 释义 + hex）升级为 keynotes 版式的专属组件
  pr-swatch（三卡，g3 网格），hex 行自带 currentColor 色点芯片。
- 杨妃粉定位为"面上的一层釉"（细带/釉线/瓷圈/hex/描边），强调职责交给胭脂 #B23A48：
  上游 100px 大数字若用杨妃粉，在凝脂底上对比度不足 3:1，故 pr-stat-v 用胭脂。
- 字体栈：标题 'Songti SC'/'STSong'/'SimSun' 系统宋体，色名 'STKaiti'/'KaiTi'，小字
  'Noto Sans SC'（自托管）；不新增 webfont 文件。
