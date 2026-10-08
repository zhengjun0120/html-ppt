# neo-brutalism（新野兽派）

**来源归属**：本模板的视觉概念（配色、装饰母题、风格气质）提取自开源项目
[arcsin1/oh-my-ppt](https://github.com/arcsin1/oh-my-ppt) 的 `resources/styles/neo-brutalism`
（Apache-2.0 License，其 NOTICE 要求保留归属，本文件即归属声明）。
页面结构、CSS、demo 与版式契约为本仓库从零手写，未复制上游代码。

## 上游风格要点摘录（原 SKILL.md 大意）

- 配色：标题深墨 #0f172a、正文暗石板 #1e293b，强调色荧光橙 #f97316 与信号红 #ef4444；
  背景 linear-gradient(160deg, #fff7ed, #ffedd5 50%, #fde68a) 从暖白渐变到柔黄；
  卡片 rgba(255,255,255,0.96)，描边 rgba(15,23,42,0.65)——实际按 3px 实线执行；
  硬阴影偏移 4px、零模糊。
- 排版：粗犷无衬线或等宽，字重偏重；标题大且自信，可直接全大写；正文不拖泥带水。
- 布局：卡片排列大胆、不规则，允许少量重叠或偏移，但数量由内容决定；粗边框与硬阴影
  本身已经很重，默认留足空隙，避免积木墙。
- 场景：创业路演、创意比稿、反叛精神演讲、年轻品牌发声。
- 禁忌：微妙设计、柔和颜色、细线边框、模糊阴影、追求优雅精致、像企业汇报。

## 转换说明

- 原上游预览是「斜置描边色块 + 描边标题盒 + 歪斜标签 + 黑色通栏条」的封面构图；
  本模板按本仓库 deck-v2 契约重建为 9 版式骨架
  （cover/contents/keynotes/split/metrics/quote/divider/moments/closing），无 code 版式。
- 招牌动作全部保留：nb-box（描边标题盒 + 橙色硬阴影）、nb-badge（黑底橙字徽章 + 红色
  硬阴影）、nb-strip（黑通栏条）、歪斜标签（.row .nb-tag 交替旋转）；斜置色块做成
  .slide::after 的 SVG ambient 层（低透明度，不与内容抢地盘）。
- 硬阴影纪律统一为「零模糊 + 颜色轮换」：卡片阴影黑/橙/红（nb-card / nb-hot / nb-wild），
  时间线节点奶油/橙/黄/红轮换；荧光黄 #FDE68A 另任「荧光笔划」（nb-mark）。
- 与本仓库已有 brutalist-bold（工业粗野：哑光纸、无阴影、单一警戒红）刻意区分：
  本模板是暖画布 + 厚描边 + 多彩硬阴影 + 歪斜手感的新野兽派。
- 字体栈：Inter / Noto Sans SC（900 黑体），元数据用 JetBrains Mono；
  不新增 webfont 文件。
