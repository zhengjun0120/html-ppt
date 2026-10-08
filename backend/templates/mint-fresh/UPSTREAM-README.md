# mint-fresh（薄荷清新 · 高级绿）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/mint-fresh`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：浅薄荷青 #B5D7C3、豆沙绿 #A3C4A9、灰绿 #8FA89A、奶白绿 #E8F0E5 为面，深墨绿
  #2F4A3A 做标题、暖灰 #6B7268 做正文；整体低饱和、高明度，清新但不轻浮。
- 排版：现代无衬线（思源黑体 / HarmonyOS Sans）中等偏粗标题，细无衬线宽行距正文；
  可用衬线做点缀标题增加质感。
- 装饰：植物线稿（叶片、枝条、蕨类）、圆形色块标注、细线分隔、半透明色块卡片。
- 布局：大量留白（≥40%），模块化网格，卡片大圆角 + 半透明浅绿底，柔和呼吸感。
- 禁忌：高饱和霓虹、深暗背景、硬朗几何、塞满画面、商业 KPI 式密集数据。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（FEUILLE 色卡页，无多版式结构）；本模板按本仓库
  deck-v2 契约重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 预览里的模糊光斑（blob）与叶脉线稿改写为 CSS radial-gradient + data-URI SVG 的 ambient
  层（.slide::before/::after，z-index:-1），每页自动衬底、不进骨架；quote 页另立一株蕨叶
  （mf-fern）呼应上游「植物线稿」母题。
- 薄荷青 #B5D7C3 收进题签胶囊与编号圆角块做「身份填充」；关键数据用苔绿 #3E7A5E
  （上游墨绿加深一档）保证浅底上的对比度；未新增 webfont。
- demo 叙事改为「薄荷生活：极简护肤品牌」发布，覆盖上游 styleCase（护肤/疗愈/科普）。
