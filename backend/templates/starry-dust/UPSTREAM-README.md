# starry-dust（星屑柔光 · 童话感）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/starry-dust`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：雾紫 #D8C5E8、香槟粉 #F4D9D0、星空蓝 #B8C8E8 作大面积柔和底色，暖金 #E8C896
  作星光点缀，奶白 #FBF5F0 作文字；整体低饱和、高明度，带磨砂柔雾感。
- 意象：四角星与六角星、细碎光点、月亮、云朵、独角兽、彩虹；圆形与棉花糖形状的色块、
  柔和渐变光晕。
- 排版：圆润手写体或柔和 Display 标题，圆润无衬线正文；字号偏大、字距宽松。
- 布局：元素散落分布，中心大标题 + 周边散落的星星；卡片大圆角 + 柔和阴影，像棉花糖质感。
- 禁忌：高饱和刺眼色、硬朗几何与锐角、深暗色背景、严肃商务字体、塞满画面。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 定位差异：与 indigo-lotus（深色紫蓝国风）刻意区分——本模板全浅色柔光，紫罗兰
  #6B4FA8 只做墨色、按钮与徽章等小面积深色；暖金 #E8C896 收暗为 #C89A54 做星光，
  投屏不刺眼。
- 柔光色晕、月亮碎星与页底云朵做成 CSS data-URI SVG 的 ambient 层（.slide::before/::after，
  z-index:-1），每页自动衬底、不进骨架，生成侧零成本。
- 字体栈：标题与正文用圆体系（'Yuanti SC'/'STYuanti'/'YouYuan'/'幼圆'，回退
  PingFang SC/Microsoft YaHei）；不新增 webfont 文件。
