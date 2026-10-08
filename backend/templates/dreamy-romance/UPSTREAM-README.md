# dreamy-romance（梦幻浪漫·治愈系）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/dreamy-romance`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：雾紫 `#D8C5E8` 大面积底，樱花粉 `#F4C8D8`、星空蓝 `#B8C8E8` 对位柔光，
  奶油白 `#FBF5F0` 做文字与留白，香槟金 `#E8C896` 提亮星光，深紫灰 `#6B5A8A` 作正文。
  整体低饱和、高明度，磨砂柔雾 + 微光质感。
- 意象：柔光晕（大半径 radial-gradient + blur）、星光、花瓣、羽毛、蝴蝶、月亮；圆形与波浪色块。
- 排版：标题手写体或柔和 Display（衬线抒情），正文细无衬线 Light、行距宽、字号偏大；
  卡片大圆角 + 半透明 + 柔光阴影；大量留白，强调氛围。
- 禁忌：高饱和霓虹刺眼、深暗背景、硬朗现代几何、严肃商务字体、白色背景、塞满画面。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 上源预览的三团 glow（blur 80px）改写为 `.slide::before` 上的三层 radial-gradient 柔雾，
  月亮为同层 data-URI SVG 弦月；花瓣与香槟金星光为 `.slide::after` 的 data-URI SVG：
  每页自动衬底（z-index:-1）、不进骨架。
- 上源卡片/题签的 `backdrop-filter: blur` 被移除（截图管线渲染不稳），以
  「半透明白底 + 白描边 + 紫调柔光阴影」等效表达；磨砂质感保留，玻璃拟态不做。
- 对比度纪律：正文深紫灰 `#6B5A8A`（对奶油底 5.6:1）；樱粉大数字用深一档
  `#C0709A`，小号樱粉文字用 `#A85578`；金色只做描环与星光装饰不进正文。
- 字体栈：标题 'Songti SC'/'Noto Serif SC' 宋体抒情，正文 'Noto Sans SC' 细无衬线；
  不新增 webfont 文件。
