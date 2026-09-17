// data-dark 契约规格
export default {
  id: 'data-dark',
  name: '数据深空',
  description: '玻璃拟态暗色风：大数字、发光卡片、终端命令行，数据与研究叙事首选',
  tags: ['数据', '研究', '暗色', '玻璃拟态'],
  scenario: ['数据分享', '研究报告', '年度总结', '趋势解读'],
  source: 'graphify-dark-graph',
  scope: 'tpl-graphify-dark-graph',
  fonts: ['Inter', 'Noto Sans SC', 'JetBrains Maple Mono'],
  variants: [
    { id: 'default', name: '极光' },
    { id: 'sunset', name: '落日', class: 'gd-sunset',
      css: '.tpl-graphify-dark-graph.gd-sunset .gd-big.gd-rainbow{background:linear-gradient(90deg,#fb923c,#f472b6,#a78bfa);-webkit-background-clip:text;background-clip:text;color:transparent}' },
  ],
  demoLayouts: ['opening', 'context', 'glass-grid', 'big-num', 'terminal', 'compare', 'trend', 'closing'],
  layouts: [
    {
      id: 'opening', name: '数据开场', roles: ['cover'],
      use: '开场：编号 + 大标题 + 一句话钩子',
      constraints: '标题 ≤14 字；lede ≤24 字；kicker 是数据集/报告名',
      skeleton: `<section class="slide" data-layout="opening">
  <div class="gd-ambient"></div>
  <div class="gd-snum">{{页码/编号}}</div>
  <p class="gd-eyebrow">{{报告名}}</p>
  <h1 class="gd-h1">{{标题}}</h1>
  <p class="gd-lede">{{钩子一句话}}</p>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'gd-ambient', 'gd-snum', 'gd-eyebrow', 'gd-h1', 'gd-lede', 'notes'],
    },
    {
      id: 'context', name: '背景铺垫', roles: ['divider', 'content'],
      use: '研究背景/问题定义（居中叙事）',
      constraints: '标题 ≤14 字；lede ≤36 字',
      skeleton: `<section class="slide" data-layout="context">
  <div class="gd-ambient"></div>
  <div class="gd-snum">{{编号}}</div>
  <div class="gd-eyebrow">{{章节标签}}</div>
  <h1 class="gd-h1">{{问题/背景}}</h1>
  <p class="gd-lede">{{展开一句}}</p>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'gd-ambient', 'gd-snum', 'gd-eyebrow', 'gd-h1', 'gd-lede', 'notes'],
    },
    {
      id: 'glass-grid', name: '玻璃数据卡', roles: ['data', 'content'],
      use: '3-4 张玻璃卡并列（每卡一个维度：数字或事实）',
      constraints: '3-4 卡；卡标题 ≤8 字；内容 ≤28 字；数字无出处标“估算”',
      skeleton: `<section class="slide" data-layout="glass-grid">
  <div class="gd-ambient"></div>
  <div class="gd-snum">{{编号}}</div>
  <p class="gd-eyebrow">{{引导语}}</p>
  <h2 class="gd-h2">{{维度标题}}</h2>
  <div class="gd-grid-3 mt-l">
    <div class="gd-glass gd-glass-blue"><span class="gd-tag">{{维度}}</span><p>{{内容}}</p></div>
    <div class="gd-glass gd-glass-green"><span class="gd-tag">{{维度}}</span><p>{{内容}}</p></div>
    <div class="gd-glass gd-glass-warm"><span class="gd-tag">{{维度}}</span><p>{{内容}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'gd-ambient', 'gd-snum', 'gd-eyebrow', 'gd-h2', 'gd-grid-3', 'mt-l', 'gd-glass', 'gd-glass-blue', 'gd-glass-green', 'gd-glass-warm', 'gd-tag', 'notes'],
    },
    {
      id: 'big-num', name: '彩虹大数', roles: ['data'],
      use: '全页一个核心数字（gd-big 彩虹渐变字）',
      constraints: '数字 ≤8 字符；指标名 ≤12 字；口径说明 ≤20 字',
      skeleton: `<section class="slide" data-layout="big-num">
  <div class="gd-ambient"></div>
  <div class="gd-snum">{{编号}}</div>
  <p class="gd-eyebrow">{{指标语境}}</p>
  <div class="gd-big gd-rainbow mt-l">{{数字}}</div>
  <h2 class="gd-h2">{{指标名}}</h2>
  <p class="gd-lede">{{口径说明（无来源标“估算”）}}</p>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'gd-ambient', 'gd-snum', 'gd-eyebrow', 'gd-big', 'gd-rainbow', 'mt-l', 'gd-h2', 'gd-lede', 'notes'],
    },
    {
      id: 'terminal', name: '命令行实证', roles: ['code'],
      use: '终端命令 + 输出（gd-cmd 行）',
      constraints: '命令 1-2 条 + 输出 ≤6 行；每行 ≤64 字符',
      skeleton: `<section class="slide" data-layout="terminal">
  <div class="gd-ambient"></div>
  <div class="gd-snum">{{编号}}</div>
  <p class="gd-eyebrow">{{实证语境}}</p>
  <h2 class="gd-h2">{{这条命令证明什么}}</h2>
  <pre class="gd-codebox mt-l"><span class="gd-cmd">$ {{命令}}</span>
{{输出行（转义 &lt; &gt;）}}</pre>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'gd-ambient', 'gd-snum', 'gd-eyebrow', 'gd-h2', 'gd-codebox', 'mt-l', 'gd-cmd', 'kw', 'st', 'fn', 'notes'],
    },
    {
      id: 'compare', name: '对照卡', roles: ['content'],
      use: '两列对照（如方案 A/B、前后、国内外）',
      constraints: '恰好 2 列；每列 ≤34 字',
      skeleton: `<section class="slide" data-layout="compare">
  <div class="gd-ambient"></div>
  <div class="gd-snum">{{编号}}</div>
  <p class="gd-eyebrow">{{引导语}}</p>
  <h2 class="gd-h2">{{对照标题}}</h2>
  <div class="gd-grid-3 mt-l" style="grid-template-columns:1fr 1fr">
    <div class="gd-glass"><span class="gd-tag">{{列 A}}</span><p>{{内容}}</p></div>
    <div class="gd-glass"><span class="gd-tag">{{列 B}}</span><p>{{内容}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'gd-ambient', 'gd-snum', 'gd-eyebrow', 'gd-h2', 'gd-grid-3', 'mt-l', 'gd-glass', 'gd-tag', 'notes'],
    },
    {
      id: 'trend', name: '四格趋势', roles: ['data'],
      use: '四个趋势/发现并列（grid-4 玻璃卡）',
      constraints: '恰好 4 卡；每卡 ≤26 字',
      skeleton: `<section class="slide" data-layout="trend">
  <div class="gd-ambient"></div>
  <div class="gd-snum">{{编号}}</div>
  <p class="gd-eyebrow">{{引导语}}</p>
  <h2 class="gd-h2">{{趋势标题}}</h2>
  <div class="gd-grid-4 mt-l">
    <div class="gd-glass"><span class="gd-tag">{{发现}}</span><p>{{内容}}</p></div>
    <div class="gd-glass"><span class="gd-tag">{{发现}}</span><p>{{内容}}</p></div>
    <div class="gd-glass"><span class="gd-tag">{{发现}}</span><p>{{内容}}</p></div>
    <div class="gd-glass"><span class="gd-tag">{{发现}}</span><p>{{内容}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'gd-ambient', 'gd-snum', 'gd-eyebrow', 'gd-h2', 'gd-grid-4', 'mt-l', 'gd-glass', 'gd-tag', 'notes'],
    },
    {
      id: 'closing', name: '收束', roles: ['thanks', 'cta'],
      use: '研究收束：结论一句 + 后续动作',
      constraints: '结论 ≤20 字；lede ≤24 字',
      skeleton: `<section class="slide" data-layout="closing">
  <div class="gd-ambient"></div>
  <div class="gd-snum">{{编号}}</div>
  <p class="gd-eyebrow">{{章节}}</p>
  <h1 class="gd-h1">{{结论一句}}</h1>
  <p class="gd-lede">{{后续动作/致谢}}</p>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'gd-ambient', 'gd-snum', 'gd-eyebrow', 'gd-h1', 'gd-lede', 'notes'],
    },
  ],
  rules: `# data-dark · 质量规则

- gd-ambient/gd-snum 是模板家具：ambient 背景层每页保留原样，snum 写章节编号（01-08）。
- 玻璃卡色（blue/green/warm）是语义：冷数据用 blue、正向用 green、风险用 warm，不要随机选。
- gd-big 彩虹大字一页只允许 1 个；数字必须真实或标“估算”。
- gd-cmd 行以 $ 开头；命令与输出必须真实可复现，不编造运行结果。
- 玻璃卡内容天然少字：单页 ≤150 字，塞满就毁了玻璃感。
- 禁 emoji。
`,
}
