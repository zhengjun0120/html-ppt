# sage-academia（青竹 · 学术风）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/sage-academia`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：纯净白底大面积留白；低饱和自然绿为主——橄榄绿 #5b805b、浅豆绿 #a9d3a9、抹茶绿
  #a9c683，辅以青灰蓝 #66a4a4 / #537e85 / #52717b 做冷色平衡；标题近黑 #000000/#3b3b3b，
  正文中性灰 #5f5f5f/#727273/#404040，冷暖各半、绿意主导不刺眼。
- 排版：SF Pro Text + PingFang SC 无衬线，标题醒目加粗，正文 20-23px 行高舒展。
- 装饰：矢量几何（正圆、六边形、菱形、虚线连接 stroke-dasharray 6 4），渐变大圆与粗圆环
  做背景层次；禁厚重阴影、立体特效、霓虹与花哨手写体。
- 布局：模块化、对称或左右分区，元素对齐精准、留白充足；适合毕业答辩、开题、课题汇报。

## 转换说明

- 原上游只有一个 1600×900 单页预览（开题报告封面：渐变大圆、粗圆环、虚线、六边形矩阵）；
  本模板按本仓库 deck-v2 契约重建为 9 版式骨架
  （cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 六边形矩阵不进骨架（clip-path 卡片不利于 22-44 字说明的排版），其「规整几何」气质交给
  ambient 层：右上渐变大圆（.slide::before）+ 左下粗圆环/虚线/青点（.slide::after，
  data-URI SVG，z-index:-1）。
- --bg 取 #FAFBF7 暖白（米白偏绿），呼应任务设定的人文学院气质；卡为白面细边 + 极轻阴影，
  遵守上游「禁厚重阴影」。
- 绿色语义收紧：编号、数据、强调药丸、按钮四处；青灰蓝只做题签。引文页是全模板唯一衬线位，
  用 'Noto Serif SC' 承担人文气质。
- 字体栈：'SF Pro Text'/'PingFang SC' 无衬线主体；不新增 webfont 文件。
