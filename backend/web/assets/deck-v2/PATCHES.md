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

## 上游版本

- 导入日期：2026-09-17
- 来源仓库：html-ppt-skill-main（本地副本 D:\go_files\html-ppt-skills\html-ppt-skill-main）
- runtime.js：1162 行，零依赖
- base.css：257 行
- animations/animations.css：入场动画库

## 2026-09-18 质量验收补充

（patch 落在模板而非 runtime——runtime.js/base.css 保持上游原样。tech-sharing/style.css
追加的 `.slide>.deck-footer` 定位修复记入该模板的 ADAPTATION.md。）
