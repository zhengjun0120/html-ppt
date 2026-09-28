# song-rain-poetic（宋人生活·听雨）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/song-rain-poetic`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：宣纸白 #FBFAF5 / 雾灰白 #F2F4F3 底，墨黛 #2E3A3B 主字、青灰 #5A6E70 正文；
  雨雾蓝 #E8F4F8 与嫩芽绿 #C5E8D5 为主色，暖灯黄 #F5E6A3 / 淡杏 #F0E2B6 副色；
  朱砂 #C0392B 与深青 #3A6B6E 只做强调。
- 意象：雨丝、青瓦屋檐、芭蕉叶、茶盏、油纸伞、窗棂、远山；装饰用水墨淡彩、雨纹、细线、印章、回纹。
- 排版：标题宋体/仿宋加粗、墨黛色；副标题楷体，嫩绿或青灰；正文行高 2.0；引文用古风衬线；
  标题旁配竖排诗句。
- 禁忌：高饱和荧光、赛博霓虹、暗黑沉重背景、锋利几何、硬朗金属、拥挤排版。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 屋檐雨丝（右上）与芭蕉叶（左下）做成 CSS data-URI SVG 的 ambient 层
  （.slide::before/::after，z-index:-1），每页自动衬底、不进骨架，生成侧零成本。
- 暖灯黄 #F5E6A3 收敛为唯一身份元素「茶签卡左竖线」（sr-tea）：原色作大面积底色会发闷，
  只做细线点缀后反而成为模板的识别记号。
- 朱砂 #C0392B 沿用上游原值：它本身偏沉，投屏不刺眼，只用于印章与关键数据。
- 字体栈：标题 'Songti SC'/'SimSun'（系统宋体），竖排与印章 'KaiTi'（系统楷体），
  小字 Noto Sans SC；不新增 webfont 文件。
