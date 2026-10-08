# macaron-mist（柔雾甜梦）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/macaron-mist`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：薄荷绿 `#8FD9D5`、粉珊瑚 `#FFC2B0`、薰衣草紫 `#C0B0E8` 三件套做柔雾低饱和处理；
  奶油白 `#FBF5F0` 打底、浅香槟 `#F0E6D0` 点缀、暖灰 `#8A8090` 正文、深紫 `#6B5A8A` 标题。
- 意象：马卡龙、甜甜圈、糖果线稿；圆形与波浪边卡片；柔和渐变光晕、磨砂柔雾；小花、爱心、丝带。
- 排版：圆润 Display 字体标题、圆润无衬线正文；圆形卡片像散落的马卡龙，大量圆角与柔和阴影。
- 禁忌：高饱和刺眼色、深暗背景、硬朗几何与锐角、严肃商务字体、塞满画面。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 三团柔雾光晕（薄荷/薰衣草/珊瑚的径向渐变）与右下马卡龙线稿、左下糖屑做成 ambient 层
  （.slide::before/::after，z-index:-1），每页自动衬底、不进骨架，生成侧零成本。
- 三色明确分工（写进 rules.md）：薄荷深色 `#3F9A94` 只做关键数字与行动钮，珊瑚深 `#D97B62`
  只做题签圆点/时间点/标签，薰衣草只做编号、圆点与描边；标题墨色为深紫 `#6B5A8A`，
  正文暖灰——低饱和原色直接当文字可读性不够，故各配一个"深一档"的可读变体。
- 字体栈：圆体 'Yuanti SC'/'Hiragino Maru Gothic ProN'/'YouYuan' 打头，回退 Noto Sans SC；
  不新增 webfont 文件。翻页动画按上游建议带轻微弹性（cubic-bezier(.34,1.3,.5,1)，.6s）。
