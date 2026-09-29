# midcentury（世纪中叶）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/midcentury`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：深棕 #78350F 标题、焦糖 #92400E 正文，芥末黄 #CA8A04 与青绿 #0D9488 是经典强调对，
  背景 160° 奶油三段渐变 #FEF3C7 → #FEF9C3 → #FEF3C7；卡为 90% 白底 + 淡焦橙 rgba(180,83,9,.2) 轻描。
- 意象：Eames 椅、Nelson 钟、Saarinen 餐桌的精神；装饰用有机几何——杏仁形、回旋镖、
  原子图案，内容区块用几何线条连接而非硬边框。
- 排版：几何无衬线或带曲线美的衬线，字号层级优雅、行距宽松，温暖但不拥挤。
- 适合：设计史叙事、家居美学展示、复古品牌故事、生活方式分享——「经典永不过时」的时刻。
- 禁忌：现代科技感、冷色调、过于锐利的几何形状、装饰喧宾夺主、先锋未来感。

## 转换说明

- 原上游只有一个 16:9 单页预览（有机形状簇 + 标本签 + 色板圆点）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 右上有机形状簇（芥末黄肾形 + 青绿回旋镖 + 焦橙圆点）与左下原子射线做成 CSS data-URI SVG 的
  ambient 层（.slide::before/::after，z-index:-1），每页自动衬底、不进骨架，生成侧零成本。
- 三色分工：芥末黄管题签/横条/数字上缘，青绿管药丸重音/行动钮，焦橙只给关键数据与时间点，
  避免芥末黄大字在奶油底上对比不足。
- 字体栈：标题 Century Gothic/Futura 几何无衬线（系统栈，失败回落 Trebuchet MS），正文
  Georgia/宋体，小标签与数字 JetBrains Mono；不新增 webfont 文件。
- letterbox 底色 --bg 取渐变中段 #FEF9C3，翻页时页外区域不跳色。
