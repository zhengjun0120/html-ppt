// weekly-report 契约规格
export default {
  id: 'weekly-report',
  name: '工作周报',
  description: '亮色高效风：KPI 面板、本周交付、数据图表、阻塞与下周计划，汇报用',
  tags: ['周报', '汇报', '团队'],
  scenario: ['工作汇报', '周报', '迭代回顾', '项目同步'],
  source: 'weekly-report',
  scope: 'tpl-weekly-report',
  fonts: ['Inter', 'Noto Sans SC', 'JetBrains Maple Mono'],
  variants: [
    { id: 'default', name: '商务蓝' },
    { id: 'violet', name: '紫罗兰', class: 'wr-violet',
      css: '.tpl-weekly-report.wr-violet{--accent:#7c5cff}\n.tpl-weekly-report.wr-violet .week-chip{background:#7c5cff}' },
    { id: 'teal', name: '青碧', class: 'wr-teal',
      css: '.tpl-weekly-report.wr-teal{--accent:#0e9488}\n.tpl-weekly-report.wr-teal .week-chip{background:#0e9488}' },
  ],
  demoLayouts: ['cover', 'kpi-grid', 'shipped-list', 'metrics-chart', 'blockers', 'next-week', 'thanks'],
  layouts: [
    {
      id: 'cover', name: '周报封面', roles: ['cover'],
      use: '周报开头：周期徽章 + 大标题 + 汇报人',
      constraints: '标题 ≤14 字；kicker 是周期（如“2026 年第 38 周”）',
      skeleton: `<section class="slide" data-layout="cover">
  <div class="cover-head"><div class="logo">{{团队/项目名}}</div><div class="week-chip">{{周期}}</div></div>
  <p class="kicker">{{汇报主题}}</p>
  <h1 class="h1 mt-s">{{本周标题}}</h1>
  <p class="lede mt-m">{{一句话总览}}</p>
  <div class="deck-footer"><span class="meta">{{汇报人 · 日期}}</span></div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'cover-head', 'logo', 'week-chip', 'kicker', 'h1', 'mt-s', 'lede', 'mt-m', 'deck-footer', 'meta', 'notes'],
    },
    {
      id: 'kpi-grid', name: 'KPI 面板', roles: ['data'],
      use: '本周 4-8 个关键指标（good/warn/bad 状态色）',
      constraints: '4-8 卡；每卡 = 指标名 + 数字 + 环比；状态色如实标注',
      skeleton: `<section class="slide" data-layout="kpi-grid">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{指标面板标题}}</h2>
  <div class="grid g4 mt-l">
    <div class="kpi good"><span class="value">{{数字}}</span><span class="label">{{指标名}} <span class="delta">{{环比}}</span></span></div>
    <div class="kpi warn"><span class="value">{{数字}}</span><span class="label">{{指标名}} <span class="delta">{{环比}}</span></span></div>
    <div class="kpi bad"><span class="value">{{数字}}</span><span class="label">{{指标名}} <span class="delta">{{环比}}</span></span></div>
    <div class="kpi"><span class="value">{{数字}}</span><span class="label">{{指标名}} <span class="delta">{{环比}}</span></span></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kicker', 'h2', 'grid', 'g4', 'mt-l', 'kpi', 'good', 'warn', 'bad', 'value', 'label', 'delta', 'notes'],
    },
    {
      id: 'shipped-list', name: '本周交付', roles: ['content'],
      use: '本周完成的事项清单（ship-item 逐条）',
      constraints: '4-8 条；每条一句话 ≤30 字，可带状态 pill',
      skeleton: `<section class="slide" data-layout="shipped-list">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{交付标题}}</h2>
  <div class="mt-l">
    <div class="ship-item">{{完成事项（可带 <span class="pill">状态</span>）}}</div>
    <div class="ship-item">{{完成事项}}</div>
    <div class="ship-item">{{完成事项}}</div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kicker', 'h2', 'mt-l', 'ship-item', 'pill', 'notes'],
    },
    {
      id: 'metrics-chart', name: '数据图表', roles: ['data'],
      use: '趋势柱状（chart-bars，高度/宽度按真实比例）',
      constraints: '3-6 根柱；柱高按真实比例给 style；每组 = 名称 + 数值',
      skeleton: `<section class="slide" data-layout="metrics-chart">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{趋势标题}}</h2>
  <div class="chart mt-l">
    <div class="chart-bars">
      <div class="col"><div class="c" style="height:80%"><span>{{数值}}</span></div><span class="l">{{名称}}</span></div>
      <div class="col"><div class="c" style="height:80%"><span>{{数值}}</span></div><span class="l">{{名称}}</span></div>
      <div class="col"><div class="c" style="height:80%"><span>{{数值}}</span></div><span class="l">{{名称}}</span></div>
    </div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kicker', 'h2', 'chart', 'mt-l', 'chart-bars', 'col', 'b', 'lbl', 'notes'],
    },
    {
      id: 'blockers', name: '阻塞与风险', roles: ['content'],
      use: '卡住的事、需要的支持（blocker 逐条）',
      constraints: '1-4 条；每条说清“卡在哪 + 需要谁做什么”',
      skeleton: `<section class="slide" data-layout="blockers">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{阻塞标题}}</h2>
  <div class="mt-l">
    <div class="blocker">{{阻塞描述：卡在哪 + 需要什么支持}}</div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kicker', 'h2', 'mt-l', 'blocker', 'notes'],
    },
    {
      id: 'next-week', name: '下周计划', roles: ['content', 'cta'],
      use: '下周要做的事（next-row 逐条）',
      constraints: '3-6 条；每条 ≤24 字、动词开头',
      skeleton: `<section class="slide" data-layout="next-week">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{下周标题}}</h2>
  <div class="mt-l">
    <div class="next-row">{{计划事项}}</div>
    <div class="next-row">{{计划事项}}</div>
    <div class="next-row">{{计划事项}}</div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kicker', 'h2', 'mt-l', 'next-row', 'notes'],
    },
    {
      id: 'thanks', name: '收尾', roles: ['thanks'],
      use: '周报收尾：一句话 + 联系方式',
      constraints: '大字 ≤10 字；补充 ≤16 字',
      skeleton: `<section class="slide" data-layout="thanks">
  <p class="kicker">{{周期}}</p>
  <h2 class="h2">{{收尾一句话}}</h2>
  <p class="lede mt-m">{{补充}}</p>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kicker', 'h2', 'lede', 'mt-m', 'notes'],
    },
  ],
  rules: `# weekly-report · 质量规则

- KPI 状态色（good/warn/bad）必须与数据同向：坏消息用 bad，不要粉饰。
- chart-bars 的柱高必须与数值成比例；数值直接标在柱内。
- 每条交付/计划写“做了什么/要做什么”，不写形容词；无数据支撑的判断删掉。
- blockers 页没有阻塞就整页省略，不要写“无阻塞”凑页。
- 周期信息（week-chip）保持模板原样格式（如 W38 · 09.14-09.20）。
- 禁 emoji；单页文字 ≤200 字。
`,
}
