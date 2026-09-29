# retro-tv（复古电视）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/retro-tv`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：深棕 #78350F 标题、焦糖 #92400E 正文，琥珀 #F59E0B 与深琥珀 #D97706 是画面灵魂；
  背景 145° 奶油 → #FDE68A → #FCD34D 三段渐变，模拟 CRT 的暖光；卡为暖黄 85% 底 +
  琥珀 rgba(217,119,6,.3) 描边。
- 意象：穿越时光的 CRT——暖黄屏幕光晕、琥珀扫描线、奶油色外壳、机身旋钮；
  内容像电视画面一样框在圆角屏框里，保留模拟颗粒的复古质感。
- 排版：复古等宽或衬线字体，字号层级分明但不过分现代，手工排版的温度。
- 适合：怀旧叙事、八零九零年代主题、复古品牌故事、回忆录式演讲——唤起「那时真好」的时刻。
- 禁忌：现代极简、冷色调、过强黑色对比、过于「数字化」的界面感。

## 转换说明

- 原上游只有一个 16:9 单页预览（棕色外壳 bezel + 圆角屏幕 + 扫描线 + 辉光 + 底部三旋钮）；
  本模板按本仓库 deck-v2 契约重建为 9 版式骨架
  （cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 整页做成「屏幕」：.slide::before 画 18px 内缩的圆角屏框（2px 琥珀描边）+ 扫描线
  （repeating-linear-gradient 2px/4px）；.slide::after 画中心辉光、底部两枚旋钮圆点与
  feTurbulence 雪花噪点（opacity 0.05 平铺），全部 z-index:-1，每页自动衬底、骨架零成本。
- 频道角标（rt-badge，切角方牌）、频道签（rt-kicker，等宽字 + LED 点）、屏幕圆角卡（rt-screen）
  转成骨架专属类；关键数字用 #B45309 深琥珀棕，避免 #D97706 在暖黄底上对比不足。
- 字体栈：标题/正文 Georgia/宋体衬线，标签与编号 JetBrains Mono；不新增 webfont 文件。
- letterbox 底色 --bg 取渐变中段 #FDE68A，翻页时页外区域不跳色。
