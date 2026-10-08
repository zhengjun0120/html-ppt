# healing-color-card（情绪疗愈色卡）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/healing-color-card`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：薰衣草紫 #C8B8D9、暖桃米 #E8CEC8、薄荷绿 #9CD9C9 三色协奏；淡丁香 #DCC8E8、
  奶杏 #F0DDD2、淡雾绿 #BFE8DB 作辅助；奶油白 #FBF6F2 / 暖象牙 #FAF3EC 底；
  文字暖墨灰 #4A4448 与柔棕 #6B5D5F，强调藕荷紫 #9B7FB5。
- 排版：标题圆润无衬线，副标题细体，正文 Light 行高 2.0；引文可用手写体增加治愈感。
- 意象：云朵、月亮、心形、花朵、波浪、气泡；柔和渐变、圆角形状、细线、大面积柔和色块。
- 布局：大色块分区柔和过渡，信息卡大圆角淡色底，充足留白。
- 禁忌：高饱和刺眼色、暗黑背景、锋利几何、硬朗金属、密集排版、严肃商务风。

## 转换说明

- 原上游只有一个 1600×900 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 色卡母题落地为两种身份元素：三张渐变色卡（hc-card-lav / -peach / -mint）与小圆色板
  （hc-chip，封面与引文页成组出现）；藕荷紫 #9B7FB5 收敛为唯一强调色（数据、编号、强调签）。
- 上游预览中的 emoji（🌙🍑🌿）按本仓库规则弃用，治愈感改由文字体温与留白承担。
- 月牙与柔光色晕做成 CSS ambient 层（.slide::before/::after，radial-gradient + data-URI SVG，
  z-index:-1），每页自动衬底、不进骨架，生成侧零成本。
- 字体栈：标题 'Yuanti SC'/'YouYuan'（系统圆体），正文同栈细体宽行距；引文 'STKaiti'/'KaiTi'。
  不新增 webfont 文件。
