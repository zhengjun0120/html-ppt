// 对比度核验：tokens.css 里所有"文字/背景"配对的 WCAG 对比值
// 用法：node tmp/check-contrast.mjs
const themes = {
  dark: {
    bg: '#07080A', surface: '#101111', surface2: '#1B1C1E',
    text1: '#F9F9F9', text2: '#B0B3B6', text3: '#82868C',
    accent: '#2DD4BF', accentContrast: '#04302A',
    success: '#5FC992', warning: '#FFBC33', danger: '#FF6363', info: '#55B3FF',
  },
  light: {
    bg: '#F5F7FB', surface: '#FFFFFF', surface2: '#EEF2F8',
    text1: '#1B2537', text2: '#4E5D77', text3: '#5E6E89',
    accent: '#0F766E', accentContrast: '#FFFFFF',
    success: '#047857', warning: '#B45309', danger: '#DC2626', info: '#2563EB',
  },
};

function lum(hex) {
  const [r, g, b] = [1, 3, 5].map((i) => parseInt(hex.slice(i, i + 2), 16) / 255)
    .map((c) => (c <= 0.03928 ? c / 12.92 : Math.pow((c + 0.055) / 1.055, 2.4)));
  return 0.2126 * r + 0.7152 * g + 0.0722 * b;
}
function ratio(fg, bg) {
  const [a, b] = [lum(fg), lum(bg)].sort((x, y) => y - x);
  return (a + 0.05) / (b + 0.05);
}

const pairs = [
  ['text-1 / bg（正文）', 'text1', 'bg', 4.5],
  ['text-1 / surface（卡片正文）', 'text1', 'surface', 4.5],
  ['text-2 / surface（次文本）', 'text2', 'surface', 4.5],
  ['text-3 / surface（弱化文本）', 'text3', 'surface', 4.5],
  ['text-3 / surface-2（输入框 placeholder）', 'text3', 'surface2', 4.5],
  ['accent / bg（链接、文字态主色）', 'accent', 'bg', 4.5],
  ['accent-contrast / accent（主按钮文字）', 'accentContrast', 'accent', 4.5],
  ['success / surface', 'success', 'surface', 4.5],
  ['warning / surface', 'warning', 'surface', 4.5],
  ['danger / surface', 'danger', 'surface', 4.5],
  ['info / surface', 'info', 'surface', 4.5],
];

let fail = 0;
for (const [name, t] of Object.entries(themes)) {
  console.log(`\n== ${name} ==`);
  for (const [label, fg, bg, min] of pairs) {
    const r = ratio(t[fg], t[bg]);
    const ok = r >= min;
    if (!ok) fail++;
    console.log(`${ok ? 'PASS' : 'FAIL'}  ${r.toFixed(2).padStart(5)} : 1  ${label}  (${t[fg]} on ${t[bg]})`);
  }
}
console.log(fail === 0 ? '\n全部通过 WCAG AA' : `\n${fail} 项未达标`);
process.exit(fail === 0 ? 0 : 1);
