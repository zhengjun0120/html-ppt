# handdrawn-watercolor（治愈手绘水彩）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/handdrawn-watercolor`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：天蓝 #87CEEB、湖蓝与米白 #F5F5DC 打底，浅卡其、暖黄与橘粉过渡，
  珊瑚红 #FF7F50 只做点缀；背景是大面积留白或纸张质感的浅米色底，低饱和是底线。
- 母题：海浪、太阳、纸船、雨伞、云朵、天气图标、丝带横幅、手绘边框与轻微纹理；
  彩铅 + 水彩晕染，轮廓柔和，带儿童绘本式的不规则感。
- 排版：标题用较粗的毛笔/手写风字体；副标题与正文细字，层级清晰；
  上下分区或模块化排布，留白充足、呼吸感明显，不要默认卡片宫格。
- 禁忌：高饱和刺眼色、尖锐几何、复杂渐变、强烈阴影、拥挤排版。

## 转换说明

- 原上游只有一个单页风格预览（色卡 + 晕染色斑，无多版式结构）；本模板按本仓库
  deck-v2 契约重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 水彩感全部用 CSS/SVG 实现：径向渐变软边界做"晕开"的色斑与圆点（无清晰边界），
  1.6px 低透明度边框做铅笔软线，不规则 border-radius 做绘本式圆角；
  纸纹用 data-URI SVG feTurbulence 以 0.05 透明度平铺。不用任何外部图片，不加 webfont。
- 珊瑚 #FF7F50 收暗为 #E4693F、天蓝收深为 #5FA8C9/#3E7C99 作文字级用色：
  原色在浅纸底上对比度不足，投屏与截图管线下观感更稳。
- 右上水彩晕与页底纸纹波浪做成 CSS ambient 层（.slide::before/::after，z-index:-1），
  每页自动衬底、不进骨架，生成侧零成本。
- 字体栈：'Yuanti SC'/'YouYuan' 圆体打头（绘本手感），回退 PingFang SC/Microsoft YaHei；
  不新增 webfont 文件。
