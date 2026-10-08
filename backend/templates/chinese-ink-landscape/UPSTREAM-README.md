# chinese-ink-landscape（中式水墨意境·山水）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/chinese-ink-landscape`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：云雾白 #F4F1EA / 宣纸米 #EDE7DA 底，墨黛 #2C3E50 主字、暖灰 #5D6D7E 正文；
  青黛蓝 #4A90A4 与远山灰蓝 #8B9DC3 为主色，暖沙米 #D4C4A0 / 淡陶土 #C9B891 副色；
  朱砂 #A93226 只用于印章与关键数据。
- 意象：层叠山峦、云雾、孤舟、松柏、亭台、飞鸟；装饰用水墨晕染、留白云气、细线、印章。
- 排版：标题宋体加粗、墨黛色、大字距，气势沉稳；副标题楷体青黛蓝；正文行高 1.9；
  引文用古风衬线；标题旁配竖排诗句。
- 禁忌：高饱和荧光、赛博霓虹、金属质感、复杂渐变、拥挤排版、卡通或工业风混搭。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 页底层叠山峦（灰蓝 / 青黛 / 墨三层，低透明度）做成 CSS data-URI SVG 的 ambient 层，
  云雾用两抹 radial-gradient 的 ::before 层（.slide::before/::after，z-index:-1），
  每页自动衬底、不进骨架，生成侧零成本。
- 上游预览右下的四格色卡收编为 cl-swatches 身份元素：只允许在封面出现一次，作山水四色的
  正面示意，不再承担配色说明职能。
- 朱砂 #A93226 沿用上游原值：本身偏沉的印泥红，只用于印章、题签细线与关键数据。
- 字体栈：标题 'Songti SC'/'SimSun'（系统宋体，大字距），竖排与印章 'KaiTi'（系统楷体），
  小字 Noto Sans SC；不新增 webfont 文件。
