# cream-pastel（奶油温柔 · 高级感）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/cream-pastel`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：奶油白 #FBF7F2 大面积打底，豆沙粉 #F4E9E8、抹茶绿 #A3B18A、暖灰绿 #9A9488 作
  辅助色块；深棕灰 #5C5248 正文，暖驼 #C9A876 偶尔作高级强调；整体低饱和、高明度、
  奶油磨砂质感。
- 排版：无衬线中等字重标题、字距宽松；正文细体行距宽；可混排衬线做点缀标题。
- 装饰：柔和色块卡片（圆角 + 柔和阴影）、植物线稿（尤加利叶、橄榄枝、棉花）、
  圆形/椭圆形色块标注 + HEX 值、细虚线分隔。
- 布局：大量留白、模块化网格、大圆角半透明奶油卡；克制、温柔、信息密度低。
- 适合：家居生活、护肤美妆、咖啡茶饮、烘焙、母婴、慢生活。
- 禁忌：高饱和、霓虹、荧光、深暗背景、硬朗现代几何、密集 KPI 卡、塞满画面。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 「焦糖边」定为身份语言：奶油卡细描边（cp-card）、数字卡顶边（cp-stat）、强调签
  （cp-pill-accent）三处，正文与大面积底色不进焦糖。
- 上游的色块标注 + HEX 值落地为色样行（cp-sw 圆角色块 + cp-hex 色号），用于 keynotes
  卡底；植物线稿（尤加利枝）做成 CSS ambient 层（.slide::before，data-URI SVG，z-index:-1），
  每页自动衬底、不进骨架。
- 字体栈：标题与正文 'Noto Sans SC'/'PingFang SC'，衬线点缀 'Songti SC'/'STSong'。
  不新增 webfont 文件。
