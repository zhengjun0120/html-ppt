# classic-duo-blue（经典双色·米黄深蓝）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/classic-duo-blue`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：主色深海蓝 #1660AB 与暖米黄 #F9F2E0；辅助浅米 #FCF6E8、藏青 #0D3D6E；
  背景米白 #FBF5E6、象牙 #FFFCF2；文字深蓝 #0D3D6E、暖墨 #2C2416；强调金棕 #B8860B。
- 排版：标题衬线体（Noto Serif SC Bold / Playfair Display）深海蓝、气势沉稳；
  正文思源宋体行高 1.8；数字与英文用大号衬线，体现典籍质感。
- 装饰：书卷、罗盘、几何线条等古典意象；细金线、几何边框、圆点、序号；
  风格简洁高级、强对比双色、大面积色块分区、少量装饰。
- 布局：双色色块明确分区或对比；标题大字 + 序号 + 副标题；信息卡片米色底深蓝字，
  或深蓝底米色字；留白与色块节奏分明。
- 禁忌：不超过三种主色、高饱和荧光、赛博霓虹、卡通可爱、过度装饰。

## 转换说明

- 上游预览是 1600×900 单页封面（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 上游封面的右侧深蓝大色块收进 quote 版式的书脊色块（cd-panel，竖排题铭）与
  cd-badge / cd-btn 的深蓝实底，遵守「双色克制」不再整页铺蓝。
- 四角金色书角与右上暖调纸感光晕做成 .slide::before/::after 的 CSS 层
  （z-index:-1，纯 CSS 渐变实现），每页自动衬底、不进骨架。
- 字体栈：'Playfair Display' / 'Noto Serif SC' / 宋体系统栈（衬线，webfont 缺省时由
  Georgia + Songti 兜底），不新增 webfont 文件。
