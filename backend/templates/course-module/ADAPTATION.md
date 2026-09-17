# course-module 入库适配记录

来源：html-ppt-skill/templates/full-decks/course-module（MIT，见 UPSTREAM-README.md）
适配日期：2026-09-18（scaffold + materialize 流水线）

1. 共享资产路径改绝对路径（/assets/deck-v2/*）；
2. 插入 SLIDES:START/END 挂载标记；demo 页按序补 data-layout；
3. style.css 末尾追加主题变体槽（template.json variants 登记）；
4. 契约文件：template.json / layouts.md（7 版式）/ rules.md / 本文件。
