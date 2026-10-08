# sakura-soft-healing（樱花治愈）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/sakura-soft-healing`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：莫兰迪低饱和——浅米粉 `#F7E7DD` 打底，淡粉 `#FBC0CB` + 玫粉 `#CC6F88` 双层樱花，
  草绿 `#72AB60`、湖水蓝 `#5BA0AD` 做清新的对色，深棕 `#3A2A22` 做文字；饱和度刻意压低，
  所有颜色像蒙了一层薄雾。
- 意象：盛开的樱花树冠、飘落花瓣、平静湖水、草地小花、芦苇与天鹅；水彩+彩铅的柔和晕染，
  边缘不锐利；远景/中景/近景的三层风景构图，留白是治愈感的灵魂。
- 排版：圆润手写/软笔标题（500-600 字重），细宋体小字安静待在角落；中英文都避免粗重几何无衬线。
- 禁忌：高饱和明亮色、纯黑大块、动感强烈元素、粗重几何字体、锋利线条、廉价马卡龙粉。
- 动画：节奏极慢（1.2-1.8s、ease-out），花瓣轻飘、文字淡入，绝不快切弹跳。

## 转换说明

- 原上游只有一个 16:9 单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 顶部花影花瓣 + 页底湖岸草色做成 CSS data-URI SVG 的 ambient 层（.slide::before/::after，
  z-index:-1），透明度压到 .24-.48，每页自动衬底、不进骨架，生成侧零成本。
- 与本仓库 chinese-cream-blossom（国风）刻意区分：本模板走现代柔和日系——无印章、无楷体
  宋体、无竖排文字；标题用圆润黑体栈，正文 Noto Sans SC。玫粉 `#CC6F88` 收稳为 `#C06B84`
  作为唯一身份强调色（关键数字/时间点/强调药丸/行动钮），淡樱 `#F4BCC7` 只做圆号与装饰线。
- 翻页动画取上游「慢」的气质：8px 轻移 + .65s ease-out（贴近管线 650ms 上限，不超时）。
  字体栈 'Hiragino Maru Gothic ProN'/'Yuanti SC'/'YouYuan' 打头，回退 Noto Sans SC；
  不新增 webfont 文件。
