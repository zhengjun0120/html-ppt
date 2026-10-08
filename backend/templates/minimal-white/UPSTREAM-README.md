# minimal-white（极简白）

**来源归属**：本模板的视觉概念（配色、气质、留白原则）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/minimal-white`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：标题 #0F172A 深墨、正文 #475569 石板灰；强调色 #3B82F6 蓝与 #60A5FA 浅蓝是仅有的彩色，
  背景为 #FFFFFF→#F8FAFC→#F1F5F9 的近纯白渐变；卡片 94% 白、边框 rgba(148,163,184,.14) 细到几乎不存在。
- 排版：Inter 或系统无衬线，标题偏细（字重 200 起步、字号 48px 以上），正文 18-20px；
  行距慷慨、字间距精确——文字层级是唯一的视觉结构。
- 布局：大留白是核心原则，一切对齐到隐形网格；不填满每个角落，「呼吸感」优先于「信息全」。
- 动效：fade-up、stagger-list、rise-in。
- 场景：内部汇报、技术评审、严肃话题、长时间观看的工作坊。禁忌：强渐变、阴影、
  超过两种颜色、装饰性元素。

## 转换说明

- 原上游只有一个 16:9 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 上游预览的三个身份元素全部保留并转为每页自动衬底的 ambient 层（.slide::before/::after，
  z-index:-1）：右侧巨型「留白」水印字、页顶一线蓝渐变细条；骨架零成本。
- 蓝色的使用比上游更收敛：全模板只剩标题短线（mw-rule）与环境细条两处，
  大数字、按钮、卡片一律深墨不染色——把「几乎不可见的点缀」执行到底。
- 字体栈：Inter + Noto Sans SC（fonts.css 已自托管，Inter 可变字重覆盖 100-900 的细字重层级）；
  不新增 webfont。
