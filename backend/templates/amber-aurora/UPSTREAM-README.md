# amber-aurora（扁豆紫蜜陀僧 · 国风治愈）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/amber-aurora`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：扁豆紫 #9E7059（低饱和暖棕紫）作顶部天空，蜜陀僧 #EC954E（高饱和暖蜜橙）作
  底部水面与落日，垂直渐变过渡；点缀浅金 #F4D9B0 萤光，奶白 #FBF3E8 作文字。
- 意象：流动云层、散落光点、水面倒影、柔和曲线；垂直「天—云—水」三段式，紫在上、
  橙在下，标题压在紫橙交界处做视觉锚点，正文卡片半透明米白底浮在画面上。
- 排版：传统书法变体标题，正文现代无衬线平衡；强调意境而非信息密度。
- 禁忌：冷色、硬朗几何与锐角、商业 KPI 卡片、英文 Heavy Grotesk、塞满画面。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 对比度决策：上游奶白文字直接压在 #9E7059→#EC954E 渐变上不足 2.5:1；本模板把渐变
  整体掺一层米白雾（雾化后顶 #B28672、底 #F2B27C），正文改用深琥珀墨 #3A2012 系
  （对全幅 ≥4.5:1），上游原色扁豆紫/蜜陀僧降为飘带线与装饰，奶白只出现在暮色深褐
  印章、方章与按钮上（对 rgba(90,42,16,.88) 底 ≈10:1）。
- 暮色天空、落日光晕、萤光与水面倒影做成 CSS data-URI SVG 的 ambient 层
  （.slide::before/::after，z-index:-1），每页自动衬底、不进骨架，生成侧零成本。
- 字体栈：标题与正文 Songti SC/STSong/KaiTi/SimSun 宋楷系，小字 PingFang SC；
  不新增 webfont 文件。
