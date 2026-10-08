# children-warm-orange（童趣橙暖）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/children-warm-orange`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：暖橙渐变天空 `#FF7F00 → #FFA500` 打底，明亮草绿 `#70D422` 强对比；亮橙 `#FE7F08`、
  暖黄 `#FEBB3E` 做强调；淡粉 `#FFADDD`、浅蓝紫 `#A6B0FF`、纯白做童趣点缀；深棕 `#5A3A1A` 做小字。
- 意象：木马、纸飞机、彩色气球、蓬松云朵、太阳与草地；所有轮廓「圆、胖、软」，禁细线与锐角。
- 排版：圆润粗手写体大标题（700-800 字重），圆润无衬线正文（400-500 字重）；
  卡片圆角 20-28px + 柔和投影；热闹来自装饰与色彩节奏，正文保持低到中密度。
- 禁忌：冷色调主导、锋利几何、低饱和灰调、成人化精致排版、写实插画、深沉暗背景。

## 转换说明

- 原上游只有一个 16:9 单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 满屏橙渐变天收进 ambient：页面底色改暖橙奶油 `#FFF6E9`，太阳、简笔云朵、星星与草绿缓坡
  做成 CSS data-URI SVG 的 ambient 层（.slide::before/::after，z-index:-1），每页自动衬底、
  不进骨架，生成侧零成本；正文区保持清爽，照顾 9 页连排的可读性与排版成熟度。
- 亮橙 `#FE7F08` 收暗为 `#EE7009`（投屏不晃眼），且只落在题签、关键数字、强调药丸与行动钮；
  暖黄 `#FEBB3E` 做贴纸与糖果圆号；草绿只画草地；粉/蓝紫只在背景装饰。
- 字体栈：圆体 'Yuanti SC'/'YouYuan'/'Hiragino Maru Gothic ProN' 打头，回退 Noto Sans SC；
  不新增 webfont 文件。翻页动画按上游建议带轻微回弹（cubic-bezier(.34,1.56,.64,1)，.55s）。
