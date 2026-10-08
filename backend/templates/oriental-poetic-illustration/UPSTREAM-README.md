# oriental-poetic-illustration（东方意境插画）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的
`resources/styles/oriental-poetic-illustration`（Apache-2.0 License，其 NOTICE 要求
保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：宣纸白 #F5F5F5 / 竹青淡绿 #E8F5E9 / 月光蓝 #E3F2FD 底，墨竹绿 #2E7D32 与
  苍翠 #4A7C59 主色，墨黑 #1B1B1B 与远黛灰 #5D6D7E 文字，朱红 #D32F2F 只用于印章
  与点缀，淡金 #D4A574 与薄雾青 #B2DFDB 做辅助。
- 意象：修竹、竹叶飘落、月轮、云气、山石、石灯笼；细笔线描、淡彩晕染、留白为主。
- 排版：行楷/瘦金体风标题（本仓库落成系统楷体栈），宋体正文宽行距；引文可竖排；
  数字与英文配无衬线细体。
- 禁忌：高饱和、荧光色、金属质感、复杂渐变、霓虹光、拥挤排版、卡通/赛博混搭。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 月轮（radial-gradient 圆 + 双层光晕）与右侧竹影线描（竹竿、竹节、飘落竹叶一律
  data-URI SVG 椭圆）做成 ambient 层（.slide::before/::after，z-index:-1），
  每页自动衬底、不进骨架，生成侧零成本。
- 色值基本沿用上游：宣纸白/月光蓝/竹青晕染、墨竹绿 #2E7D32 数据主色、朱红 #D32F2F
  印章；仅把印章底色从纯红调至原值并控制墨点尺寸，投屏不刺眼。
- 字体栈：标题 'STKaiti'/'KaiTi'（系统楷体，瘦金气质用字距表达），正文 Georgia/宋体栈
  宽行距，小字与数字 Noto Sans SC；不新增 webfont 文件。
