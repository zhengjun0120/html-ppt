# aurora（极光）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/aurora`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：标题 #1E1B4B 深靛、正文 #312E81；强调色 #4F46E5 靛紫与 #06B6D4 青蓝；预览页采用
  深蓝夜空 #0A0A1A 打底，翠绿 #A7F3D0/#6EE7B7、靛紫 #6366F1/#A855F7 的多色光晕融化在 blur 里，
  构成极光带；卡片 rgba(255,255,255,0.55) 高透明、边框 rgba(255,255,255,0.4) 半透明白描边。
- 排版：现代无衬线，标题大而醒目，正文极简；文字浮在渐变之上，可带模糊背景增强可读性。
- 意象：极地夜空、翠绿-蓝紫光带、远山剪影；气质是「第一眼与最后一眼都让人屏住呼吸」。
- 动画关键词：blur-in、gradient-flow、shimmer-sweep。
- 禁忌：不用硬边框或锐利布局；不让文字量超过极光的美。

## 转换说明

- 原上游定位偏「封面/CTA 页风格」，只有一个单页预览；本模板按本仓库 deck-v2 契约重建为
  9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing），
  依据任务定位采用预览页的深蓝天幕方向而非浅色四段渐变。
- 深色模板纪律：正文一律浅色（星白 #F1F5FF / 浅靛 #C3CDF2 / 小字 #8E99CE），
  卡为暗夜玻璃 rgba(16,20,52,.55) + blur + 细浅描边；--bg 取天幕渐变中段值 #131735，
  letterbox 与翻页过渡不跳色。
- 极光带做成多层径向渐变、星野与挪威山脊剪影做成 data-URI SVG 的 ambient 层
  （.slide::before/::after，z-index:-1），每页自动衬底、不进骨架，生成侧零成本。
- 极光渐变文字（au-stat-v / au-quote）两端色均为浅色（#6EE7B7→#67E8F9、#A7F3D0→#67E8F9→#C4B5FD），
  暗底上全程满足对比度；行动钮为浅色渐变实底配深墨绿字。
- 字体栈：Inter / Segoe UI / PingFang SC / Noto Sans SC，无新增 webfont 文件。
- demo 叙事「追光者：极光观测旅行计划」：观测点位、KP 指数（NOAA 口径）、
  成团场次与成功率的口径均标注为领队日志与空间天气预报。
