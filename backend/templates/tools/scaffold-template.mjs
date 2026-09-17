#!/usr/bin/env node
/**
 * scaffold-template.mjs — 模板入库脚手架（一次性机械活自动化）：
 *   1. 从 html-ppt-skill 复制模板目录（index.html + style.css + README）
 *   2. 共享资产路径改写为 /assets/deck-v2/*
 *   3. 插入 SLIDES:START/END 挂载标记
 *   4. 抽取 demo 每页的 class 结构清单（供人工写 layouts.md 参考）
 * 用法：node scaffold-template.mjs <源模板名> <目标模板id>
 */
import { readFileSync, writeFileSync, mkdirSync, existsSync, appendFileSync } from 'node:fs';
import { join } from 'node:path';

const SKILL = 'D:/go_files/html-ppt-skills/html-ppt-skill-main/html-ppt-skill-main/templates/full-decks';
const DEST = 'D:/go_files/html-ppt/backend/templates';

const [srcName, destId] = process.argv.slice(2);
if (!srcName || !destId) {
  console.error('用法: node scaffold-template.mjs <源模板名> <目标id>');
  process.exit(1);
}

const src = join(SKILL, srcName);
const dest = join(DEST, destId);
if (!existsSync(join(src, 'index.html'))) { console.error('源模板不存在:', src); process.exit(1); }
if (existsSync(dest)) { console.error('目标已存在:', dest); process.exit(1); }
mkdirSync(dest, { recursive: true });
mkdirSync(join(dest, 'preview'), { recursive: true });

// 1) index.html：路径改写 + 挂载标记
let html = readFileSync(join(src, 'index.html'), 'utf8');
html = html
  .replaceAll('../../../assets/fonts.css', '/assets/deck-v2/fonts.css')
  .replaceAll('../../../assets/base.css', '/assets/deck-v2/base.css')
  .replaceAll('../../../assets/animations/animations.css', '/assets/deck-v2/animations.css')
  .replaceAll('../../../assets/runtime.js', '/assets/deck-v2/runtime.js');
// SLIDES 标记：<div class="deck"> 之后 / 其对应闭合之前
const deckOpen = html.indexOf('<div class="deck">');
if (deckOpen < 0) { console.error('源模板没有 <div class="deck">'); process.exit(1); }
const afterOpen = deckOpen + '<div class="deck">'.length;
html = html.slice(0, afterOpen) + '\n<!-- SLIDES:START -->' + html.slice(afterOpen);
// 结尾标记：deck 容器的闭合（最后一个 </div> 在 <script 之前）
const scriptIdx = html.lastIndexOf('<script');
const closeIdx = html.lastIndexOf('</div>', scriptIdx);
html = html.slice(0, closeIdx) + '<!-- SLIDES:END -->\n' + html.slice(closeIdx);
writeFileSync(join(dest, 'index.html'), html);

// 2) style.css / README
writeFileSync(join(dest, 'style.css'), readFileSync(join(src, 'style.css'), 'utf8'));
if (existsSync(join(src, 'README.md'))) {
  writeFileSync(join(dest, 'UPSTREAM-README.md'), readFileSync(join(src, 'README.md'), 'utf8'));
}

// 3) demo 页面 class 结构清单（写 layouts.md 的原料）
const seg = html.split('<!-- SLIDES:START -->')[1].split('<!-- SLIDES:END -->')[0];
const sections = [...seg.matchAll(/<section[^>]*>[\s\S]*?<\/section>/g)].map(m => m[0]);
let report = `# ${destId} demo 结构清单（脚手架生成，写 layouts.md 用后即弃）\n\n`;
sections.forEach((s, i) => {
  const open = s.match(/<section[^>]*>/)[0];
  report += `## demo 第 ${i + 1} 页\n${open}\n`;
  // 直接子元素（tag.class 一行一个）
  const body = s.slice(open.length, s.lastIndexOf('</section>'));
  for (const line of body.split('\n')) {
    const el = line.match(/<([a-z]+)[^>]*class="([^"]*)"/);
    if (el) report += `  <${el[1]} class="${el[2]}">\n`;
  }
  report += '\n';
});
writeFileSync(join(dest, 'DEMO-STRUCTURE.md'), report);

console.log(`✓ 脚手架完成: ${dest}`);
console.log('  待人工完成: template.json / layouts.md / rules.md / ADAPTATION.md / demo 页 data-layout / variants');
