# dreamy-pink-gradient（渐变美学·樱粉雾蓝）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的
`resources/styles/dreamy-pink-gradient`（Apache-2.0 License，其 NOTICE 要求保留归属，
本文件即归属声明）。页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：淡樱粉 `#FFCDE1`、雾紫蓝 `#7C96FF`、深靛蓝 `#2E4B8F` 三色渐变交融；粉霞 `#FFB3D1`、
  薰衣草 `#A8B8FF`、夜蓝 `#1E2D5C` 辅助；粉白 `#FFF5F9` / 淡紫白 `#F5F0FF` 做底；
  玫粉 `#FF6B9D` 做强调。如清晨黄昏交界的天空，粉与蓝温柔接吻。
- 意象：渐变色块、云朵、月亮、花瓣、星星、液态形状；柔光、光晕、细线、圆点。
- 排版：标题圆润无衬线或现代衬线，可渐变填充；正文无衬线 Light、行高约 1.9；
  信息卡用半透明白底；留白柔软，节奏舒缓。
- 禁忌：高饱和刺眼、暗黑沉重、锋利几何、硬朗工业风、严肃商务、超过四种渐变色。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 上源预览的大渐变球（orb + blur）改写为 `.slide::before` 上的三层 radial-gradient 柔光，
  星子改写为 `.slide::after` 的 data-URI SVG 星尘带：每页自动衬底（z-index:-1）、不进骨架。
- 上源的 `backdrop-filter: blur` 磨砂卡被移除：截图管线与低性能设备下渲染不稳，
  以「半透明白 + 白描边 + 淡紫软阴影」等效表达磨砂感。
- 对比度纪律：正文一律深靛 `#2E4B8F` 实色；渐变填充只给 h1、大数字（dp-stat-v）、
  编号球与按钮，按钮渐变两端收暗为 `#F0568F→#6B82F2` 以保住白字大字号 3:1；
  小号玫粉文字统一用收暗档 `#C22460`。
- 字体栈：'Noto Sans SC'/'PingFang SC' 无衬线为主，引文用 Georgia/宋体衬线；
  不新增 webfont 文件。
