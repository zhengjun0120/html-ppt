# vaporwave（蒸汽波）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/vaporwave`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：标题 #f8fafc 冷白、正文 #e2e8f0 灰白，强调粉红与青蓝经典双色；背景是深紫到
  青蓝的四段渐变，如日落的数字模拟；卡片 rgba(46,16,101,.6) 半透明暗紫，
  边框是粉红的晕染描边。
- 排版：复古或未来感字体，标题可全大写带字间距；正文轻柔浮在渐变上；希腊字母与
  日文假名可做装饰。
- 布局：渐变本身就是最大的视觉元素，内容浮在色彩之上，卡片半透明、边框柔和、自由流动。
- 禁忌：朴素正式、硬边框锐利布局；渐变不要变成混乱；不需要理由。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 粉青双色取 #FF71CE / #01CDFE，渐变紫 #7C3AED 入画布渐变（按任务方向收成
  「紫粉青」三色系）；画布渐变刻意收敛在深紫-暗蓝区间，亮粉亮青只出现在 ambient
  与点缀层，正文区量测对比度 ≥7:1（文字 #E2E8F0 对最深亮带 #4C1D95）。
- 落日（条纹渐变日轮）、希腊柱影（左右两根细线柱）与地平线透视网格做成 data-URI SVG
  的 ambient 层（.slide::before/::after，z-index:-1），每页自动衬底、不进骨架。
- 字体栈：标题 Georgia/Noto Serif SC（大理石碑文气质），正文 Inter/Noto Sans SC，
  数字与标签 JetBrains Mono；不新增 webfont 文件。
- 渐变数字 vw-stat-v 用 background-clip:text 实现，外包 @supports 回退为纯青色，
  防止不支持环境里文字透明不可见。

## demo 叙事口径

demo 为「蒸汽波：一场听觉与视觉的复古」分享（虚构厂牌「回声电台」）：2.6 亿次年度播放、
1.2 万首在架曲目、15 场线下演出均标注了《2025 年度曲风报告》与统计区间等口径来源，
属演示用虚构数据。
