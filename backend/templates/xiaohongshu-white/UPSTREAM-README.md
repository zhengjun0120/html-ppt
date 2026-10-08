# xiaohongshu-white（小红书白）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/xiaohongshu-white`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：标题深焦糖 #7c2d12、正文暖棕，强调暖红 #f56565 与柔红 #fc8181 是小红书的
  标志性红；背景 linear-gradient(160deg, #ffffff → #fff5f5 → #fff1f0) 如晨光；
  卡片近纯白 #fffaf9，边框是极淡红描边 #fde2e1。
- 排版：衬线标题是风格的核心（宋体/衬线英文优雅落笔），无衬线正文，行距宽松，
  像一本精美生活方式杂志的内页。
- 母题：暖红「随手标注的重点」、收藏感、图文笔记的节奏；暖红必须克制，白底必须干净。
- 适合：小红书图文、生活方式分享、美妆美食、旅行日记——「美好生活值得记录」。
- 禁忌：冷色调、科技感元素、无衬线标题、暖红泛滥、暗色背景。

## 转换说明

- 原上游只有一个 16:9 的单页预览（无多版式结构）；本模板按本仓库 deck-v2 契约
  重建为 9 版式骨架（cover/contents/keynotes/split/metrics/quote/divider/moments/closing）。
- 社媒身份元素落地为：右上暖红书签飘带（.slide::before，clip-path，每页 ambient）、
  贴纸标签（xh-tag 实底微倾斜 / xh-tag-soft 虚线款）、白卡珊瑚左条（xh-card::before）、
  CSS 小红心收藏行（xh-heart + xh-like）、话题标签（xh-hash）与荧光划重点（xh-mark）。
- 上游预览的琥珀暖卡（card-warm）保留为 keynotes 的中间变体，一页至多一张。
- 暖红 #F56565 保持克制：只进书签、贴纸签、关键数据、强调签与话题标签，不进正文段落。
- 字体栈：标题 Georgia + 'Songti SC'/'STSong' 衬线，正文 'PingFang SC'/'Noto Sans SC'。
  不新增 webfont 文件。
