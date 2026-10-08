# cyberpunk-neon（赛博霓虹）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/cyberpunk-neon`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：标题 #f8fafc 冷白、正文 #cbd5e1 灰白，强调是霓虹粉与霓虹青两道撕裂夜幕的光；
  背景 radial-gradient(circle at 30% 10%, #1a1a2e, #0f0f1a 50%, #000000) 深紫灰到纯黑，
  如无星的夜空；卡片 rgba(15,15,26,.85)，边框是霓虹粉的发光描边。
- 排版：等宽或赛博感字体第一，标题可带 text-shadow 发光，全大写标题与代码风格片段常见。
- 布局：暗底上发光的文字与线条是唯一视觉结构；故障效果做装饰；布局不对称。
- 禁忌：暖色调、柔和元素；不要丢掉霓虹发光核心；不要让暗色沉闷；不要规整对称。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 霓虹双色调取 #FF2E88（粉，结构/描边）与 #00F0FF（青，数据/高亮）——比上游
  #ec4899/#06b6d4 更高饱和，但在浅色投屏与截图管线下仍只承担描边与点缀；
  正文保持上游冷白 #F8FAFC / 灰白 #CBD5E1，量测对比度 ≥12:1。
- 故障网格（青色细网格）与霓虹横线 + 扫描带做成 CSS ambient 层（.slide::before/::after，
  z-index:-1），每页自动衬底、不进骨架；发光用多层 background 渐变模拟光晕，不依赖
  box-shadow 大模糊。
- 字体栈：正文 Inter/Noto Sans SC，数字、参数与小字 JetBrains Mono；不新增 webfont 文件。

## demo 叙事口径

demo 为「霓虹夜行：城市摄影集」分享（虚构摄影师陈夜航）：42 个点位、6 场巡展、
18 万次快门、首印 1200 册均标注了拍摄日志/付印清单等内部口径与来源行，
属演示用虚构数据。
