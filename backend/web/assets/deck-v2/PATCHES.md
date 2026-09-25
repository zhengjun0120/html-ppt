# deck-v2 vendor 文件补丁记录

本目录的 `runtime.js` / `base.css` / `animations.css` vendor 自
[html-ppt-skill](https://github.com/)（MIT）。原则：**能不改就不改**——上游越干净，
将来同步修复越容易。每一处改动在这里记一行（文件、原因、原始行为）。

改动纪律：
- 补丁必须最小化，不做顺手重构；
- 每个补丁在源文件里用 `/* [deck-v2 patch] ... */` 注释标记锚点，方便与上游 diff；
- 行为性改动必须配 spike 验证结论。

## 补丁清单

1. runtime.js go()：history.replaceState 包 try/catch。原因：预览 iframe 带 sandbox
   （无 allow-same-origin，安全姿态沿用 v1），opaque origin 下 replaceState 抛
   SecurityError 并打断 go() 的后续状态更新（进度条/激活页切换）。深链语义由
   location.hash 赋值承担（沙箱内允许自导航），replaceState 只服务地址栏展示。

1b. runtime.js 预览模式 showSlide（2026-09-25）：去掉逐 slide 的 display:none 硬切，
    改为 is-active/is-prev class 切换（与 go() 同款）。原因：模板预览弹窗经
    preview-goto postMessage 无刷新翻页，display 切换会杀死 .slide 自带的
    .5s opacity/transform 过渡，观感是"闪一下再突兀替换"。slide 本就 absolute
    叠放 + opacity:0 隐藏，class 切换即得淡入+方向位移过渡，正常模式同款观感。

2. fonts.css unicode-range 分段 + 等宽拉丁子集（2026-09-25，性能）。原因：画廊/
   预览的字体成本实测撑爆首屏——sandboxed iframe 是独立不透明源，Chrome 按
   源分区 HTTP 缓存，同一批字体被每张卡全量重复下载（单次画廊打开 91 次字体
   请求 / ~180MB）；且 @font-face 无 range 时页面里任何"中文栈元素中的西文"
   都会触发 6MB+ 的 CJK 等宽全量文件。改法：
   - 各家族按"拉丁 U+0000-2E7F / CJK 及以上 U+2E80-10FFFF"拆两条 @font-face，
     浏览器按需下载所在分段；
   - 新增 JetBrainsMapleMono-{Regular,Bold}-latin.woff2（58/60KB，由全量文件
     `python -m fontTools.subset <full>.woff2 --unicodes=U+0000-2E7F
     --flavor=woff2 --output-file=<out>` 生成）。**上游字体更新后必须重跑**
     这两条命令再生成子集，否则拉丁字形停在上旧版本。
   效果实测：单次画廊打开字体请求 91 → 9（首次，之后 max-age 内 0 次），
   传输量 ~180MB → ~17MB；纯拉丁等宽页不再触碰 CJK 全量文件。
   配套（非本目录）：router.revalidateStatic 对字体给 7 天 max-age；
   画廊内置模板 iframe 加 allow-same-origin 恢复缓存共享（内容是仓库静态
   文件，与宿主同源可信；用户模板 ut-* 保持 scripts-only 沙箱不变）；
   /api/templates/:id/preview?slide=N 只返回第 N 页（缩略卡不再为整本
   demo 付解析/布局账，TrimToSlide fail-open）。

## 上游版本

- 导入日期：2026-09-17
- 来源仓库：html-ppt-skill-main（本地副本 D:\go_files\html-ppt-skills\html-ppt-skill-main）
- runtime.js：1162 行，零依赖
- base.css：257 行
- animations/animations.css：入场动画库

## 2026-09-18 质量验收补充

（patch 落在模板而非 runtime——runtime.js/base.css 保持上游原样。tech-sharing/style.css
追加的 `.slide>.deck-footer` 定位修复记入该模板的 ADAPTATION.md。）
