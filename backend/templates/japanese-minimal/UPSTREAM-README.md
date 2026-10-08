# japanese-minimal（日式极简）

**来源归属**：本模板的视觉概念（配色、「间」的美学、装饰母题）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/japanese-minimal`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：标题 #1C1C1C 墨黑、正文 #525252 灰；强调色 #DC2626 朱红与 #B91C1C 深朱红是整页唯一的
  彩色——一笔即可；背景 #FAF8F5→#F5F0E8→#EDE5D8 从象牙白到暖灰，如和纸质感；
  卡片 90% 纯白、边框 rgba(189,28,28,.12) 是极淡的朱红描边。
- 排版：衬线字体带书法气韵，标题不过大但气场十足，正文行距极为宽松；留白本身是排版的一部分。
- 布局：极致留白，每页只承载一个核心信息；朱红只出现一次，像书法的最后一笔点睛。
- 母题：円相（一笔圆）、竖排文字、纸角折痕、枯山水的少即是多。
- 场景：品牌升级、匠人故事、禅意叙事、高端产品发布。禁忌：密集布局、多彩、粗重阴影边框。

## 转换说明

- 原上游只有一个 16:9 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 上游预览的円相、朱红竖线、纸角三个母题全部保留：円相（含一段淡朱笔意）与纸角做成
  CSS data-URI SVG 的 ambient 层（.slide::before/::after，z-index:-1），每页自动衬底。
- 「一页一笔朱红」执行得比上游更严：每页只允许 jm-red-line / jm-stat-accent /
  jm-pill-accent / jm-vert-accent 之一出现，其余全部交给墨色与留白（rules.md 有硬性条款）。
- 竖排题款（jm-vert，writing-mode:vertical-rl）升格为模板身份元素，封面与引文页各立一条。
- 与本仓库水墨江南模板刻意区隔：本模板用宋体细字（非楷书）、和纸暖底（非宣纸冷米）、
  无印章无远山，朱红更少更轻，气质是「间」而非「诗意」。
- 字体栈：Georgia + 宋体系统栈（'Songti SC'/'STSong'/'SimSun'），不新增 webfont 文件。
