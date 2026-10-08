# chinese-pastel-spring（中国传统色·春日嫩柳）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/chinese-pastel-spring`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：嫩柳绿 #A8BF8F、柳芽黄绿 #C5D99E 双主色，桃花粉 #F3A694 与淡樱粉 #F8C9B8
  副色，宣纸米 #F7F3E8 / 象牙白 #FBF8F0 底，墨黛 #3A3A3A 与深棕 #5A4030 文字，
  朱砂 #C0392B 强调、淡金 #D4A574 点缀。
- 意象：桃花枝、柳条、燕剪春风、蝴蝶、纸鸢、团扇；装饰用水墨淡彩晕染、印章、细线分隔，
  一两朵桃花或柳叶做页眉点缀。
- 排版：标题方正清刻本悦宋或思源宋体加粗，副标题楷体，正文思源宋体行高 1.8，
  数字用现代无衬线细体。
- 禁忌：不要高饱和荧光、赛博霓虹；不要硬朗几何、锋利金属质感；不要沉重暗色背景；
  不要过度装饰。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 上游预览顶部的四段渐变色带（band）与右上桃花枝 SVG 都改为 CSS ambient 层
  （.slide::before / ::after，data-URI，z-index:-1），每页自动衬底、不进骨架，
  生成侧零成本；桃花枝整体降透明度，压在内容之下不抢正文。
- 上游预览的四格色卡（palette swatch）收敛为 keynotes 的三张实色小色卡
  （嫩柳绿 / 柳芽黄绿 / 桃花粉），卡上文字用墨黛深色以保证可读性；
  右侧左竖线暖卡（card）保留为 split 右列的 ps-card-line（竖线取桃粉色）。
- 色值全部沿用上游原值；数字从上游的细无衬线落地为 Inter 300 细体大数字。
- 字体栈：标题宋体栈，次级标题与卡题 STKaiti/KaiTi，小字 Noto Sans SC；
  不新增 webfont 文件。
