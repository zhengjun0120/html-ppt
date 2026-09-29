# sunset-warm（日落暖）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/sunset-warm`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：标题深焦糖 #7c2d12、正文焦橙 #9a3412，强调亮橙 #f97316 与金黄 #fbbf24 是日落的
  两种光泽；背景 linear-gradient(135deg, #fed7aa → #fdba74 → #fb923c) 三段暖渐变；
  卡片 rgba(255,255,255,.88) 白底浮在暖色上，边框 rgba(251,146,60,.3) 橘色柔和描边。
- 排版：友好的无衬线字体，字重中等偏轻；标题温暖但不沉重，字号层级分明但不过于正式
  （上游预览的标题用 Georgia 衬线，本模板保留给引文页做点缀）。
- 布局：卡片在暖色渐变上自由排列，留白充足但不冷清；布局像一张明信片；圆角与柔和阴影。
- 意象：地平线上的太阳、光环、被夕阳染色的云条（预览中的 sun/halo/cloud 母题）。
- 禁忌：冷色调或暗色、过于正式的排版、暖色压迫感、深色背景破坏日落氛围。

## 转换说明

- 原上游只有一个 16:9 单页预览（太阳+光环+云条的渐变封面，无多版式结构）；本模板按
  本仓库 deck-v2 契约重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 天际渐变 + 云条整体做成 ambient 层（.slide::after inset:0），太阳与光环画在 .slide::before，
  z-index:-1 每页自动衬底；quote 页另立一轮光环（su-halo）呼应上游母题。
- 渐变末档由 #fb923c 收浅为 #fa9c55，保证页脚深焦糖小字与正文在暖底上的对比度；
  letterbox 底色取渐变中段 #fdc98f，翻页不跳。
- 关键数字用亮橙收暗一档 #ea580c（白卡上 ≥3:1 的展示字号对比）；等宽大写题签、金黄
  圆点与描边药丸沿用上游预览的元素语言；不新增 webfont。
- demo 叙事改为「海滨落日婚礼方案」，覆盖上游 styleCase（年度庆典/生活方式/正向叙事）。
