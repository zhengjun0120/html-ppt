# starlight-fireworks（一半星河一半烟火）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/starlight-fireworks`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：深海军蓝 #0A2E5C（星河半边）与暖焰橙 #FF8C42（烟火半边）沿 45° 对角线相接；
  星辉蓝 #4A90E2、星金 #FFD700、琥珀 #FFA500 为强调，文字用星白 #FFFFFF 与淡金米 #FFF3D6。
- 意象：星空、银河、烟火、流星、柔光过渡；装饰用细密星点、烟火粒子、对角分界线。
- 排版：衬线体大标题（中粗），副标题细无衬线，正文无衬线 Light、行高 1.9、星白带轻投影。
- 适合：庆典、里程碑、愿景、年终回顾与章节开启页；数据、流程、条目页不配图。
- 禁忌：不破坏对角双色分割、不高饱和刺眼、不办公室商务风、不卡通可爱、不密集烟花墙。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 深色·沉稳取向：页面以深海军蓝为主体，右下角对角渐暗琥珀收尾（不做 50/50 硬分割），
  保证全文本对比度与截图管线稳定；星金 #FFD700 收柔为 #F5C469，投屏不刺眼。
- 星点流星（左上）与烟火粒子（右下）做成 CSS data-URI SVG 的 ambient 层
  （.slide::before/::after，z-index:-1），每页自动衬底、不进骨架，生成侧零成本。
- 字体栈：标题 Georgia/'Noto Serif SC' 衬线，小字 Inter/'Noto Sans SC'；
  不新增 webfont 文件。
- demo 叙事为「跨年双城记：星河与烟火」活动策划案，两城场次、预算与观演人数
  均带口径来源行。
