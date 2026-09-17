// knowledge-blueprint 契约规格
export default {
  id: 'knowledge-blueprint',
  name: '知识蓝图',
  description: '蓝图网格底纹 + 结构件：流程管线、双卡对照、大数字、代码盒，讲解复杂机制首选',
  tags: ['知识', '架构', '蓝图', '深色'],
  scenario: ['知识分享', '架构讲解', '机制解析', '技术深潜'],
  source: 'knowledge-arch-blueprint',
  scope: 'tpl-knowledge-arch-blueprint',
  fonts: ['Inter', 'Noto Sans SC', 'JetBrains Maple Mono'],
  variants: [
    { id: 'default', name: '蓝图青' },
    { id: 'amber', name: '蓝图琥珀', class: 'kb-amber',
      css: '.tpl-knowledge-arch-blueprint.kb-amber{--accent:#d97706}' },
  ],
  demoLayouts: ['pipeline-overview', 'section-open', 'dual-cards', 'codebox', 'statement', 'big-num', 'insight', 'legend'],
  layouts: [
    {
      id: 'pipeline-overview', name: '流程管线', roles: ['content', 'divider'],
      use: '整条链路的全景图（kb-step 逐站，hero 标当前站）',
      constraints: '3-5 站；每站 ≤8 字；hero 标 1 个',
      skeleton: `<section class="slide" data-layout="pipeline-overview">
  <div class="kb-kicker">{{引导语}}</div>
  <h1 class="kb-h1">{{链路名}}</h1>
  <p class="kb-sub">{{链路一句话}}</p>
  <div class="kb-pipeline mt-l">
    <div class="kb-step">{{站点 1}}</div>
    <div class="kb-step">{{站点 2}}</div>
    <div class="kb-step hero">{{当前站点}}</div>
    <div class="kb-step">{{站点 4}}</div>
  </div>
  <div class="kb-footer">{{页脚：模块名 · 页码}}</div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kb-kicker', 'kb-h1', 'kb-sub', 'kb-pipeline', 'mt-l', 'kb-step', 'hero', 'kb-footer', 'notes'],
    },
    {
      id: 'section-open', name: '章节开篇', roles: ['divider'],
      use: '章节过渡：大标题 + 一句话引导',
      constraints: '标题 ≤12 字；副句 ≤22 字',
      skeleton: `<section class="slide" data-layout="section-open">
  <div class="kb-kicker">{{章节进度}}</div>
  <h1 class="kb-h1">{{章节标题}}</h1>
  <p class="kb-sub">{{这一章讲什么}}</p>
  <div class="kb-footer">{{页脚}}</div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kb-kicker', 'kb-h1', 'kb-sub', 'kb-footer', 'notes'],
    },
    {
      id: 'dual-cards', name: '双卡对照', roles: ['content'],
      use: '两个概念/方案的对照卡（各带小标签行）',
      constraints: '恰好 2 卡；卡标题 ≤8 字；内容 ≤40 字',
      skeleton: `<section class="slide" data-layout="dual-cards">
  <div class="kb-kicker">{{引导语}}</div>
  <h1 class="kb-h1">{{对照标题}}</h1>
  <div class="kb-grid-2 mt-l">
    <div class="kb-card"><div class="kb-kicker">{{卡标签}}</div><p>{{内容}}</p></div>
    <div class="kb-card"><div class="kb-kicker">{{卡标签}}</div><p>{{内容}}</p></div>
  </div>
  <div class="kb-legend mt-l">{{图例/口径说明}}</div>
  <div class="kb-footer">{{页脚}}</div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kb-kicker', 'kb-h1', 'kb-grid-2', 'mt-l', 'kb-card', 'kb-legend', 'kb-footer', 'notes'],
    },
    {
      id: 'codebox', name: '代码盒', roles: ['code'],
      use: '一段核心代码/配置（等宽盒 + 高亮）',
      constraints: '代码 ≤18 行；kw/st 高亮 span 可用',
      skeleton: `<section class="slide" data-layout="codebox">
  <div class="kb-kicker">{{文件/场景}}</div>
  <h1 class="kb-h1">{{代码在做什么}}</h1>
  <pre class="kb-codebox mt-l"><span class="kw">{{代码（转义 &lt; &gt; &amp;）}}</span></pre>
  <div class="kb-footer">{{页脚}}</div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kb-kicker', 'kb-h1', 'kb-codebox', 'mt-l', 'kw', 'st', 'kb-footer', 'notes'],
    },
    {
      id: 'statement', name: '机制宣言', roles: ['quote', 'divider'],
      use: '整页一句话机制总结',
      constraints: '宣言 ≤26 字',
      skeleton: `<section class="slide" data-layout="statement">
  <div class="kb-kicker">{{引导语}}</div>
  <h1 class="kb-h1">{{机制宣言}}</h1>
  <div class="kb-footer">{{页脚}}</div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kb-kicker', 'kb-h1', 'kb-footer', 'notes'],
    },
    {
      id: 'big-num', name: '蓝图大数', roles: ['data'],
      use: '一个关键数字/量级（kb-big-num 承载）',
      constraints: '数字 ≤8 字符；来源写进正文，无出处标“估算”',
      skeleton: `<section class="slide" data-layout="big-num">
  <div class="kb-kicker">{{指标语境}}</div>
  <div class="kb-big-num mt-l">{{数字}}</div>
  <h1 class="kb-h1">{{指标名}}</h1>
  <div class="kb-footer">{{页脚}}</div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kb-kicker', 'kb-big-num', 'mt-l', 'kb-h1', 'kb-footer', 'notes'],
    },
    {
      id: 'insight', name: '洞见卡', roles: ['content'],
      use: '一条反直觉结论/关键洞见（kb-insight 承载）',
      constraints: '洞见 ≤32 字；支撑 ≤40 字',
      skeleton: `<section class="slide" data-layout="insight">
  <div class="kb-kicker">{{引导语}}</div>
  <div class="kb-insight mt-l">{{洞见一句话}}</div>
  <h1 class="kb-h1">{{洞见标题}}</h1>
  <p class="kb-sub">{{支撑说明}}</p>
  <div class="kb-footer">{{页脚}}</div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kb-kicker', 'kb-insight', 'mt-l', 'kb-h1', 'kb-sub', 'kb-footer', 'notes'],
    },
    {
      id: 'legend', name: '结构分解', roles: ['content'],
      use: '把一个结构拆成组成部件（kb-step 纵列 + 图例）',
      constraints: '3-5 个部件；每个 ≤12 字 + ≤16 字说明',
      skeleton: `<section class="slide" data-layout="legend">
  <div class="kb-kicker">{{引导语}}</div>
  <h1 class="kb-h1">{{结构名}}</h1>
  <div class="kb-pipeline mt-l">
    <div class="kb-step">{{部件 1}}</div>
    <div class="kb-step">{{部件 2}}</div>
    <div class="kb-step">{{部件 3}}</div>
  </div>
  <div class="kb-legend mt-l">{{部件间的衔接说明}}</div>
  <div class="kb-footer">{{页脚}}</div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kb-kicker', 'kb-h1', 'kb-pipeline', 'mt-l', 'kb-step', 'kb-legend', 'kb-footer', 'notes'],
    },
  ],
  rules: `# knowledge-blueprint · 质量规则

- kb-* 类是本模板的全部身份：不要混用 base 原语的 h1/h2/card 替代（骨架给什么用什么）。
- kb-step hero 一组管线里最多标 1 个（它是"当前讲到哪"的指针）。
- kb-codebox 里 HTML 必须转义；kw/st 高亮 span 只标关键字与字符串。
- 本模板是深色蓝图风：不给页面加背景色/渐变，底纹由模板自带。
- 章节开篇（section-open）页文字极少是刻意的，别塞内容。
- 禁 emoji；单页可见文字 ≤170 字。
`,
}
