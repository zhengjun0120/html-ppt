#!/usr/bin/env node
/**
 * validate-template.mjs — 模板入库校验（Node，无依赖）。
 *
 * 用法：node validate-template.mjs <模板目录> [...更多目录]
 * 不传目录 = 校验同目录下所有含 template.json 的子目录。
 *
 * 校验项（与后端 registry 的 Go 版校验同一套规则，这里供入库/贡献者本地自检）：
 *   1. 文件齐全：template.json / index.html / style.css / layouts.md / rules.md
 *   2. template.json 可解析，关键字段齐全，id 与目录名一致
 *   3. index.html 含 SLIDES:START / SLIDES:END 挂载标记
 *   4. index.html 引用的 /assets/deck-v2/* 文件真实存在（相对 assets 根解析）
 *   5. layouts.md 解析出的版式 id 集合 == template.json layouts[].id 集合
 *   6. 每个 layouts.md 条目有"合法类名："行，且其中的类名 ⊆ index.html+style.css 出现过的类名
 *   7. variants[].class（非空时）在 style.css 中有对应选择器
 *   8. index.html 的 SLIDES 区间内每个 section 都带 data-layout，且 id 已登记
 */
import { readFileSync, existsSync, readdirSync, statSync } from 'node:fs';
import { join, dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const here = dirname(fileURLToPath(import.meta.url));
const templatesRoot = resolve(here, '..');                    // 本脚本在 templates/tools/ 下
const assetsRoot = resolve(here, '../../web/assets');         // backend/web/assets
const baseCssPath = join(assetsRoot, 'deck-v2', 'base.css');

let failures = 0;
const fail = (dir, msg) => { console.error(`  ✗ ${msg}`); failures++; };
const ok = (msg) => console.log(`  ✓ ${msg}`);

// 从 CSS/HTML 全文里收集出现过的类名 token（\.[A-Za-z][\w-]* 以及 class="..." 拆词）
function collectClassNames(...texts) {
  const set = new Set();
  for (const text of texts) {
    for (const m of text.matchAll(/\.([A-Za-z][A-Za-z0-9_-]*)/g)) set.add(m[1]);
    for (const m of text.matchAll(/class="([^"]*)"/g)) {
      for (const c of m[1].split(/\s+/)) if (/^[A-Za-z][A-Za-z0-9_-]*$/.test(c)) set.add(c);
    }
  }
  return set;
}

// 解析 layouts.md：返回 Map<layoutId, {classes: string[]}>
function parseLayouts(md) {
  const out = new Map();
  const re = /^##\s+([A-Za-z][A-Za-z0-9_-]*)/gm;
  let m;
  const heads = [];
  while ((m = re.exec(md))) heads.push({ id: m[1], start: m.index });
  for (let i = 0; i < heads.length; i++) {
    const body = md.slice(heads[i].start, i + 1 < heads.length ? heads[i + 1].start : md.length);
    const cm = /合法类名[：:]\s*(.+)/.exec(body);
    const classes = cm ? cm[1].split(/[，,\s]+/).map(s => s.trim()).filter(Boolean) : [];
    out.set(heads[i].id, { classes, body });
  }
  return out;
}

function validateTemplate(dir) {
  console.log(`\n== ${dir} ==`);
  const need = ['template.json', 'index.html', 'style.css', 'layouts.md', 'rules.md'];
  for (const f of need) {
    if (!existsSync(join(dir, f))) return fail(dir, `缺少 ${f}`);
  }
  ok('五个必需文件齐全');

  let meta;
  try { meta = JSON.parse(readFileSync(join(dir, 'template.json'), 'utf8')); }
  catch (e) { return fail(dir, `template.json 解析失败: ${e.message}`); }
  for (const k of ['id', 'name', 'description', 'canvas', 'variants', 'layouts']) {
    if (meta[k] === undefined) return fail(dir, `template.json 缺字段 ${k}`);
  }
  if (meta.id !== dir.split(/[\\/]/).pop()) fail(dir, `template.json.id (${meta.id}) 与目录名不一致`);
  ok(`template.json 合法（${meta.layouts.length} 个版式，${meta.variants.length} 个变体）`);

  const html = readFileSync(join(dir, 'index.html'), 'utf8');
  const css = readFileSync(join(dir, 'style.css'), 'utf8');
  const md = readFileSync(join(dir, 'layouts.md'), 'utf8');

  if (!html.includes('<!-- SLIDES:START -->') || !html.includes('<!-- SLIDES:END -->')) {
    fail(dir, 'index.html 缺 SLIDES:START/END 挂载标记');
  } else ok('挂载标记就绪');

  // /assets/deck-v2/* 引用必须存在
  for (const m of html.matchAll(/(?:href|src)="(\/assets\/[^"]+)"/g)) {
    const p = join(assetsRoot, m[1].replace('/assets/', ''));
    if (!existsSync(p)) fail(dir, `index.html 引用不存在的资产: ${m[1]}`);
  }
  ok('index.html 资产引用可解析');

  const layouts = parseLayouts(md);
  const declared = new Set(meta.layouts.map(l => l.id));
  for (const id of declared) if (!layouts.has(id)) fail(dir, `template.json 登记的版式 ${id} 在 layouts.md 无条目`);
  for (const id of layouts.keys()) if (!declared.has(id)) fail(dir, `layouts.md 条目 ${id} 未在 template.json 登记`);
  ok(`layouts.md 与 template.json 版式一致`);

  for (const [id, l] of layouts) {
    if (l.classes.length === 0) fail(dir, `版式 ${id} 缺"合法类名："行`);
  }

  // 已知类名 = base.css 原语 ∪ 模板 index.html/style.css 出现过的类
  const known = collectClassNames(readFileSync(baseCssPath, 'utf8'), html, css);
  for (const [id, l] of layouts) {
    for (const c of l.classes) {
      if (!known.has(c)) fail(dir, `版式 ${id} 声明的类名 .${c} 在 index.html/style.css 中不存在`);
    }
  }
  ok('各版式合法类名全部存在');

  for (const v of meta.variants) {
    if (v.class && !css.includes('.' + v.class)) fail(dir, `变体 ${v.id} 的 class .${v.class} 在 style.css 中无定义`);
  }
  ok('变体 class 检查通过');

  // SLIDES 区间内的 section 都要带已登记的 data-layout
  const seg = html.split('<!-- SLIDES:START -->')[1]?.split('<!-- SLIDES:END -->')[0] ?? '';
  for (const m of seg.matchAll(/<section[^>]*>/g)) {
    const dm = /data-layout="([^"]*)"/.exec(m[0]);
    if (!dm) fail(dir, `demo section 缺 data-layout: ${m[0].slice(0, 60)}`);
    else if (!declared.has(dm[1])) fail(dir, `demo section 的 data-layout=${dm[1]} 未登记`);
  }
  ok('demo 页面版式登记检查通过');
}

// 入口
const args = process.argv.slice(2);
let dirs = args;
if (dirs.length === 0) {
  dirs = readdirSync(templatesRoot).filter(n => {
    const p = join(templatesRoot, n);
    return statSync(p).isDirectory() && n !== 'tools' && existsSync(join(p, 'template.json'));
  }).map(n => join(templatesRoot, n));
}
if (dirs.length === 0) { console.log('没有找到待校验的模板目录'); process.exit(0); }
for (const d of dirs) validateTemplate(resolve(d));
console.log(failures === 0 ? `\n全部通过 ✓` : `\n${failures} 处失败 ✗`);
process.exit(failures === 0 ? 0 : 1);
