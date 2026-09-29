# geography-classroom（地理课堂·蓝图）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的
`resources/styles/geography-classroom`（Apache-2.0 License，其 NOTICE 要求保留归属，
本文件即归属声明）。页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：白 / 浅冰蓝 / 极浅灰蓝纸张底，深海军蓝 #173B56 标题带与章节图标，青绿 #2E7C83、
  湖蓝 #5E9FB3；重点信息用珊瑚红 #D85B4A、砖红 #B94B3D 或明亮黄 #F0D34A，少量使用、
  不做大面积撞色；正文用墨黑与深灰。
- 意象：地图、等高线、剖面图、流程箭头、指南针、比例尺、图例、不规则白色标题板、
  顶部深蓝/青绿标题带。
- 排版：稳重的中文黑体标题，封面主标题允许衬线/宋体；正文保持投影可读的大字号；
  图例、坐标轴、比例尺与图注分层清楚。
- 动画：ease-out、短距离位移、克制；重点色只用于当前讲解对象。
- 禁忌：霓虹、玻璃拟态、厚重 3D、商务图库风；不编造地理事实、不虚构地点。

## 转换说明

- 上游只有一张 1600×900 的单页预览（封面 + 一页课式布局）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 经纬网淡格网（.slide::before）与页底等高线、山峰三角、指北针（.slide::after）做成
  CSS data-URI SVG 的 ambient 层（z-index:-1），每页自动衬底、不进骨架。
- 珊瑚红收深为砖红 #B94B3D 作正文级重点色（原 #D85B4A 保留在 ambient 山峰标记里）：
  投屏与截图管线下的 16px 药丸文字也能过对比度；青绿 #2E7C83 同理收深为 #28707A。
- 字体栈：标题/正文 Noto Sans SC + Inter（系统黑体回退），引文 Georgia/宋体栈；
  不新增 webfont 文件。
