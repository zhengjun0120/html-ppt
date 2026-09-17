#!/usr/bin/env node
/**
 * materialize.mjs <spec.json> — 从紧凑规格生成模板契约文件：
 * template.json / layouts.md / rules.md / ADAPTATION.md + demo 页 data-layout + 变体 CSS。
 * 规格字段见 docs/refactor-plan.md §5；layouts 条目与注册表解析器逐字对齐。
 */
import { readFileSync, writeFileSync, existsSync, appendFileSync } from 'node:fs';
import { join } from 'node:path';
import { pathToFileURL as toURL } from 'node:url';

const specArg = process.argv[2];
if (!specArg) { console.error('用法: node materialize.mjs <spec.mjs|spec.json>'); process.exit(1); }
const spec = specArg.endsWith('.mjs')
  ? (await import(toURL(specArg))).default
  : JSON.parse(readFileSync(specPath, 'utf8'));
const dest = join('D:/go_files/html-ppt/backend/templates', spec.id);

// ---- template.json ----
const meta = {
  id: spec.id,
  name: spec.name,
  description: spec.description,
  tags: spec.tags,
  scenario: spec.scenario,
  canvas: { w: 1920, h: 1080 },
  variants: spec.variants.map(v => ({ id: v.id, name: v.name, class: v.class || '' })),
  layouts: spec.layouts.map(l => ({
    id: l.id, name: l.name, use: l.use,
    ...(l.roles ? { roles: l.roles } : {}),
    ...(l.constraints ? { constraints: l.constraints } : {}),
  })),
  ...(spec.fonts ? { fonts: spec.fonts } : {}),
  source: { derived_from: `html-ppt-skill/templates/full-decks/${spec.source}`, license: 'MIT' },
};
writeFileSync(join(dest, 'template.json'), JSON.stringify(meta, null, 2) + '\n');

// ---- layouts.md ----
let md = `# ${spec.name} · 版式登记簿（LLM 唯一契约）

> **版式锁**：每个 XXX> 必须带 \`data-layout="<本文件登记的 id>"\`，只准使用骨架里出现过的类名
> （base 原语 + 本模板专属类）。未登记的版式和未知类名会被服务端直接拒收。
> 先 \`read_layout\` 拿骨架，再替换 {{占位符}}。不要写 data-id / data-title。
>
> 本模板所有专属类都以 \`.${spec.scope}\` 作用域前缀生效，骨架里已写全，照抄结构即可。
`;
for (const l of spec.layouts) {
  md += `\n---\n\n## ${l.id}（${l.name}）\n\n`;
  md += `用途：${l.use}\n`;
  md += `适用 role：${(l.roles || ['content']).join(' / ')}。\n`;
  if (l.constraints) md += `内容约束：${l.constraints}\n`;
  md += '\n```html\n' + l.skeleton.trim() + '\n```\n\n';
  md += `合法类名：${l.classes.join(', ')}\n`;
}
writeFileSync(join(dest, 'layouts.md'), md);

// ---- rules.md ----
writeFileSync(join(dest, 'rules.md'), spec.rules.trim() + '\n');

// ---- 变体 CSS 追加 ----
for (const v of spec.variants) {
  if (v.css) appendFileSync(join(dest, 'style.css'), '\n' + v.css.trim() + '\n');
}

// ---- demo 页 data-layout ----
let html = readFileSync(join(dest, 'index.html'), 'utf8');
let idx = 0;
html = html.replace(/<section([^>]*)>/g, (m, attrs) => {
  if (!attrs.includes('<section') && (m.includes(' SLIDES') || attrs.includes('data-layout'))) return m;
  // 只处理 SLIDES 区间内的 section
  return m;
});
// 简单做法：按顺序给 SLIDES 区间内的 section 补 data-layout
{
  const [pre, seg, post] = split3(html, '<!-- SLIDES:START -->', '<!-- SLIDES:END -->');
  let n = 0;
  const patched = seg
    .replace(/ data-layout="[^"]*"/g, '')           // 幂等：先剥掉旧的
    .replace(/<section([^>]*)>/g, (m, attrs) => {
    const id = spec.demoLayouts[n];
    n++;
    return id ? `<section${attrs} data-layout="${id}">` : m;
  });
  html = pre + '<!-- SLIDES:START -->' + patched + '<!-- SLIDES:END -->' + post;
  if (n !== spec.demoLayouts.length) {
    console.error(`✗ demo 页数 ${n} 与 demoLayouts ${spec.demoLayouts.length} 不一致`);
    process.exit(1);
  }
}
writeFileSync(join(dest, 'index.html'), html);

// ---- ADAPTATION.md ----
if (!existsSync(join(dest, 'ADAPTATION.md'))) {
  writeFileSync(join(dest, 'ADAPTATION.md'),
`# ${spec.id} 入库适配记录

来源：html-ppt-skill/templates/full-decks/${spec.source}（MIT，见 UPSTREAM-README.md）
适配日期：2026-09-18（scaffold + materialize 流水线）

1. 共享资产路径改绝对路径（/assets/deck-v2/*）；
2. 插入 SLIDES:START/END 挂载标记；demo 页按序补 data-layout；
3. style.css 末尾追加主题变体槽（template.json variants 登记）；
4. 契约文件：template.json / layouts.md（${spec.layouts.length} 版式）/ rules.md / 本文件。
`);
}

console.log(`✓ ${spec.id}: ${spec.layouts.length} 版式契约材料化完成`);

function split3(s, a, b) {
  const i = s.indexOf(a), j = s.indexOf(b);
  return [s.slice(0, i + a.length), s.slice(i + a.length, j), s.slice(j)];
}
