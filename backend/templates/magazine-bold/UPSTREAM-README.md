# magazine-bold（杂志大字）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/magazine-bold`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：标题 #451a03 深棕、正文 #78350f 焦糖，强调 #f59e0b 琥珀与 #d97706 深橙；
  背景是 #fef3c7 → #fef9c3 → #fefce8 的奶油色三段渐变，卡为 rgba(255,255,255,.92)
  带极淡金色描边 rgba(234,179,8,.2)。
- 排版：超大衬线标题是灵魂，字号大到文字成为画面本身；正文用较轻的无衬线，
  行距宽松；橙色点缀只给数字、关键词与 CTA。
- 布局：大留白，标题居中或左对齐；双栏或单栏大图排版。
- 动效：rise-in、stagger-list、shimmer-sweep。
- 禁忌：小字号密集排版、暗色调或深色背景、无衬线标题、橙色强调色泛滥。

## 转换说明

- 上游预览页的三个母题全部保留并组件化：超大衬线标题（.h1 112px/0.95）、
  深橙短杠（mb-bar，对应上游 .accent）、右下巨型水印字（mb-mark，对应上游 .watermark，
  demo 里用「夜」与期号「41」）。
- 巨型引号「“」（430px、10% 透明琥珀）与左下琥珀光晕做成 ambient 层
  （.slide::before/::after，z-index:-1），每页自动衬底、不进骨架。
- letterbox 用纸渐变中段值 #FCF5CF，翻页不跳；正文换 Inter 300 轻无衬线与衬线标题形成对话。
- 字体栈：标题 Georgia/宋体 900，正文 Inter/Noto Sans SC；不新增 webfont 文件。
