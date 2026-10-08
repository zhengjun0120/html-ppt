# memphis-pop（孟菲斯波普）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/memphis-pop`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：标题深色 #0F172A 压阵，正文 #374151；强调色琥珀黄 #F59E0B 与薰衣草紫 #8B5CF6 是主角，
  背景为奶油—粉紫—靛蓝的三色渐变；卡片白底 92% 不透明度 + 粗黑描边稳稳浮在色彩之上。
- 排版：粗犷无衬线、大字重，标题字号可以夸张到让排版本身成为装饰；字间距略宽给呼吸感。
- 布局：大胆不对称，区块间穿插圆点、三角、锯齿线等几何装饰；卡可以是斜的、带投影的；
  热闹来自装饰和色彩，不来自信息密度，正文保持低到中密度。
- 动效：zoom-pop、stagger-list、shimmer-sweep。
- 场景：年轻品牌发布、潮流合作、创意工作坊、毕业展示。禁忌：极简克制、正式商务排版、
  压抑色彩表现力。

## 转换说明

- 原上游只有一个 16:9 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 底色做了降饱和处理：上游 #FEF3C7→#F0ABFC→#818CF8 的强渐变换成同色相的浅三色渐变
  （#FFF7DE→#FBE3F5→#E4E2FB），正文页更耐读；饱和色保留在描边卡、签、按钮与 ambient 装饰上。
- 硬投影（box-shadow 偏移 0 模糊）与锯齿彩带（clip-path 锯齿）是波普海报语言的核心，
  与本仓库其他浅色模板的「禁硬阴影」约定不同——这是本模板的身份特征，rules.md 有豁免说明。
- 几何碎屑与锯齿彩带做成 CSS data-URI SVG / 渐变的 ambient 层（.slide::before/::after，
  z-index:-1），每页自动衬底、不进骨架。
- 字体栈：Inter + Noto Sans SC（fonts.css 已自托管），900 字重大标题为模板身份；不新增 webfont。
