# acid-blue-business（酸蓝张力）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/acid-blue-business`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：科技蓝 #435BE1 是唯一主强调色；背景 #FFFFFF、#F2F2F2、#D9D9D9；
  标题 #000000 / #0D0D0D，正文 #404040。只用一组蓝色和中性色，不加粉彩或渐变。
- 排版：SF Pro Text、PingFang SC、Helvetica Neue 等现代无衬线；主标题 90-96px 重字重
  左对齐；章节编号是独立的蓝色锚点。
- 布局：不对称构图——标题与论点居中左，右侧留白给视觉或几何；垂直装饰条、横向虚线、
  网格线与少量透视几何建立方向，不与内容争夺注意力。
- 装饰：轮廓圆、线条、小三角、直角几何块；轻量棱角感图标；禁圆润卡通与厚重玻璃卡。
- 禁忌：柔和粉嫩色、手写字体、复杂纹理、霓虹赛博效果、对称且拥挤的布局。

## 转换说明

- 上游预览是单页海报式封面（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing），
  把「海报」拆成「一页一个判断」的评审节奏。
- 预览左缘垂直蓝条转写为每页左侧电光蓝色轨（.slide::before，z-index:-1），
  右下的轮廓圆、虚线与小三角转写为 .slide::after 的 data-URI SVG 几何层；
  两者每页自动衬底、不进骨架。
- 应分派的「荧光点缀」收敛为 #CDFF3D 的三个微观几何刻度（题签刻度、数字顶线刻度、
  圆环菱形），永不承载文字——标题正文保持黑/深灰对白底与浅灰底，对比度合规。
- 字体栈：'Inter' / 'Noto Sans SC'（自托管 webfont）+ 'JetBrains Mono'（等宽标签），
  不新增字体文件。
