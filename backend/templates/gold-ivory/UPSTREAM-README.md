# gold-ivory（鎏金象牙 · 高级质感）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的
`resources/styles/gold-ivory`（Apache-2.0 License，其 NOTICE 要求保留归属，
本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：鎏金棕 #7A4E1D（高饱和暖棕，鎏金质感）作大色块与大字，象牙白 #E8DCC8
  （低饱和暖白，丝光质感）作对位与底色，深棕 #4A2E15 正文，香槟金 #C9A876
  做高级强调与细线。
- 排版：衬线标题（宋体/Playfair 一路），字重偏粗；正文无衬线 Light 行距宽；
  英文衬线斜体做装饰铭文。
- 装饰：细密斜纹与丝绒噪点、金色细线分隔与边框、菱形/卷草连续纹样、
  半透明色块卡片；古典书页与奢侈品包装感。
- 布局：对称、稳重、有仪式感；直角或小圆角。
- 禁忌：高饱和霓虹荧光、现代极简几何、轻浮弹跳动画、密集 KPI 卡、塞满画面。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 双线金框（外线 inset 40px、内线再退 8px）落成 .slide::before/::after 两个伪元素，
  配合 .slide 上的细密斜纹背景（repeating-linear-gradient，115deg，4% 香槟金），
  每页自动衬底、不进骨架；伪元素 pointer-events:none，不拦点击。
- 上游以鎏金棕做深底大色块；本仓库反转对位：象牙白 #E8DCC8 做整页底（深底页在
  投影与截图管线下的可读性不稳），鎏金棕退到大字、印记与按钮色块，香槟金只做细线。
- 金色渐变字（background-clip:text）保留上游 hero 的点睛手法，收成 gi-gold 类，
  纪律为「只许标题局部 2-4 字，全份 deck 至多 2 处」。
- 字体栈：标题与数字 Georgia/'Songti SC'（衬线），正文 Noto Sans SC Light；
  不新增 webfont 文件。
