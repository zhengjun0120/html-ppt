# arctic-cool（北极冷）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/arctic-cool`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 定位：数据和分析的理想栖所——冰蓝渐变冷静而清澈，「每一组数据都在理性之光中获得尊严」。
- 配色：上游 SKILL.md 给出浅冰蓝变体（#f0f9ff→#bae6fd 三段渐变底，#0c4a6e/#0369a1 深海蓝文字，
  #0284c7 明蓝与 #06b6d4 青蓝强调）；其 style.json 分类为「深色 · 沉稳」，
  preview.html 实际使用深海渐变（#082f49→#0284c7）+ 青蓝描边与冰面高光。
- 排版：专业无衬线、字重适中不张扬；数字字号略大、等宽对齐；图表与数字是主角。
- 布局：卡片式突出关键数据模块，低到中密度，大面积留白让数据呼吸。
- 禁忌：不要暖色调或活泼元素；不要手写体或装饰性字体；不要让装饰抢走数据的注意力。

## 转换说明

- 原上游只有单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约重建为 9 版式骨架
  （cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- **深色取向**：按 style.json 的「深色·沉稳」分类与 preview.html 的深海气质落地为深色版——
  冰蓝灰极地色板 #1B2B3A 底（渐变 #16242F→#22394A）+ 冰川青 #7FB3D5 身份色；
  上游 SKILL.md 的浅冰蓝变体未采用（浅色数据底可由仓库内 minimal-white 等模板承载）。
- 对比度：正文冰白 #EAF4FB（约 13:1）、雾蓝灰 #9DB8CC（约 7:1）、小字 #8FA9BF（约 4.8:1），
  冰川青 #7FB3D5（约 6.4:1）在深底上全部通过 AA。
- 装饰做进 ambient 层：.slide::before 为极光余晖（两路极低透明度 radial-gradient），
  .slide::after 为冰面等高线（data-URI SVG 细波浪线），z-index:-1，每页自动衬底。
- 字体栈：Inter + Noto Sans SC（克制字重），数字/日期/标签 JetBrains Mono（等宽，
  本仓库自托管字体）；不新增 webfont 文件。
- letterbox 底色 --bg 取冰蓝灰渐变中段值 #1B2F3D，翻页不跳色。
