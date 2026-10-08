# 模板标签筛选 + 名字搜索 施工方案

状态：已批准（2026-09-30）。纯前端实现，不改后端任何东西。

## 背景与数据事实

- 后端 101 个内置模板每份 `template.json` 都有 `tags`（风格标签，约 30 个词：深色/国风/杂志/科技…）与 `scenario`（适用场景）两个维度；`GET /api/templates` 公开接口一次全量返回（后端注释明确不分页）。
- 画廊此前只展示 scenario 前 3 个 chip，tags 在数据里存在但未露出。
- 用户模板（ut-*）没有 tags 字段：标签筛选天然只覆盖内置模板，用户模板只参与名字搜索。已接受。

## 已拍板的决策

1. 标签筛选 = 单选 chip +「全部」，再点同 chip 取消。
2. 搜索范围 = 名称 + id + tags 成员 + scenario 成员 + 描述，大小写不敏感子串；多词不拆分（整串匹配）。
3. 卡片上加 tags 行（最多 3 个，静态 chip，accent 软底与 scenario 中性底区分）。
4. 两处界面都做：向导画廊 + 我的模板页「从内置模板派生」区。

## 施工内容

### 新文件

| 文件 | 职责 |
|---|---|
| `frontend/src/lib/templateFilter.ts` | 纯逻辑（templateTags 词表 / matchesTemplate 匹配 / filterTemplates 过滤）+ `useTemplateFilter` 薄 composable（query/activeTag/vocab/filtered/hasFilter/toggleTag/reset） |
| `frontend/src/components/templates/TemplateFilterBar.vue` | 两处共用的筛选工具条：搜索框 + 词表 chip 行 + 筛出计数/清空 |
| `frontend/src/lib/templateFilter.test.ts` | vitest 详细单测（含 backend/templates 数据契约守卫） |

### 改动文件

**`frontend/src/components/wizard/TemplateGallery.vue`**
- `GalleryCard` 加 `tags?: string[]`；`cards` computed 带出 `tags: ut ? undefined : t.tags`。
- `useTemplateFilter(cards)`；瀑布流 `columns` 改遍历 `filtered`（我的模板置顶顺序保留）。
- 筛选条插在头部说明（原 166 行）与瀑布流（原 168 行）之间；瀑布流 `v-if="filtered.length"`，v-else 零结果空态 + 清空按钮。
- 卡片 scenario 行上方加 tags 行。

**`frontend/src/views/MyTemplatesView.vue`**
- `builtins`（原 39 行）接入 composable（解构改名 builtin 前缀，避免歧义）。
- `builtinPageCount`/`pagedBuiltins` 改基于过滤后列表；`watch([query, tag])` → `builtinPage = 1`；`load()` 越界钳制（原 58 行）改用过滤后长度。
- 筛选条插在标题行结束（原 392 行）与网格（原 393 行）之间；分页 `v-if` 改过滤后列表；补零结果空态。

### 关键设计取舍

- 词表从数据实时算出（次数降序、同次数字典序），不写死——以后加模板自动跟上。
- **词表折叠（施工中实测补全）**：实际词表 227 个标签、长尾大量 count=1，全部平铺约 15 行会把卡片区淹没（截图证据）。筛选条只平铺前 12 个 +「更多 N ▸」展开/「收起」，折叠时被选中的长尾标签保持可见。
- 筛选作用于喂给列表/分页的数组，卡片渲染、选中、变体、翻页逻辑零改动；被筛掉的卡片 iframe 直接卸载，画廊反而更轻。
- 不 debounce：百来条数据的 computed 直过滤零延迟。
- 契约守卫测试读 `../backend/templates/*/template.json`（node:fs），目录缺失自动 skip——防将来导入新批次漏 tags 导致筛选漏卡片；templates/ 下的 tools/ 等非模板目录以「存在 template.json」为准。

## 测试清单

1. templateTags：计数、单模板内去重、缺/空 tags 跳过、排序（次数降序→字典序）、空输入。
2. matchesTemplate：空条件放行；tag 命中/不命中；query 命中 name/id/描述/tags/scenario；大小写不敏感；trim；多词整串不匹配（既定行为）；无 tags 模板在 tag 筛选下排除。
3. filterTemplates：tag AND query 组合、顺序保持、零命中、空条件等价输入。
4. useTemplateFilter：toggleTag 选中/取消、reset、source 响应式重算、hasFilter 联动。
5. 数据契约：全部内置模板 tags 非空。

## 验收标准

`npm run typecheck` 0 错、`npm run test` 全绿、`npm run build` 成功；浏览器点通：画廊 chip 筛选/中文搜索/组合/清空/零结果空态/卡片 tags 显示/置顶保持；我的模板页派生区筛选 + 分页回 1 + 切页滚动定位正常。截图证据落 `tmp/template-filter-evidence/`。

## 收尾

全部文件 cp 镜像到 backend checkout（cmp 校验）；本地提交不 push。
