// editorial-white 契约规格
export default {
  id: 'editorial-white',
  name: '杂志白',
  description: '亮色杂志编辑风：细线分栏、马卡龙卡片、金句排版，轻内容与知识科普首选',
  tags: ['杂志', '亮色', '科普', '轻内容'],
  scenario: ['知识科普', '轻内容分享', '读书笔记', '生活方式', '小红书图文'],
  source: 'xhs-white-editorial',
  scope: 'tpl-xhs-white-editorial',
  fonts: ['Playfair Display', 'Noto Serif SC', 'Noto Sans SC'],
  variants: [
    { id: 'default', name: '经典白' },
    { id: 'mint', name: '薄荷', class: 'xw-mint',
      css: '.tpl-xhs-white-editorial.xw-mint{--accent:#0d9488}' },
  ],
  demoLayouts: ['hero-quote', 'statement', 'quad-cards', 'steps', 'big-stat', 'two-column', 'quote', 'grid-3'],
  layouts: [
    {
      id: 'hero-quote', name: '杂志封面', roles: ['cover'],
      use: '刊头 + 大标题 + hero 区（可放一张图或纯排版）',
      constraints: '标题 ≤12 字；副标 ≤22 字；kicker 是栏目名',
      skeleton: `<section class="slide" data-layout="hero-quote">
  <div class="xw-topline"></div>
  <div class="xw-topbar"><span>{{刊名}}</span><span>{{期号/日期}}</span></div>
  <div class="xw-page">
    <p class="xw-kicker">{{栏目名}}</p>
    <h1 class="xw-title">{{标题}}</h1>
    <p class="xw-sub">{{副标}}</p>
    <div class="xw-hero">{{hero 区（一句话或配图说明）}}</div>
  </div>
  <div class="xw-footer">{{页脚}}</div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'xw-topline', 'xw-topbar', 'xw-page', 'xw-kicker', 'xw-title', 'xw-sub', 'xw-hero', 'xw-footer', 'notes'],
    },
    {
      id: 'statement', name: '单句页', roles: ['divider', 'quote'],
      use: '整页只放一句转场/观点，杂志翻页感',
      constraints: '单句 ≤24 字',
      skeleton: `<section class="slide" data-layout="statement">
  <div class="xw-topline"></div>
  <div class="xw-topbar"><span>{{刊名}}</span><span>{{页码}}</span></div>
  <div class="xw-page">
    <p class="xw-kicker">{{栏目}}</p>
    <h1 class="xw-title">{{单句}}</h1>
  </div>
  <div class="xw-footer">{{页脚}}</div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'xw-topline', 'xw-topbar', 'xw-page', 'xw-kicker', 'xw-title', 'xw-footer', 'notes'],
    },
    {
      id: 'quad-cards', name: '马卡龙四卡', roles: ['content'],
      use: '四张浅色卡（soft-pink/blue/green/orange 各一，颜色即分类）',
      constraints: '恰好 4 卡；卡标题 ≤8 字；内容 ≤26 字；颜色与内容情感匹配',
      skeleton: `<section class="slide" data-layout="quad-cards">
  <div class="xw-topline"></div>
  <div class="xw-page">
    <h2 class="xw-title-md">{{主题}}</h2>
    <div class="xw-grid-2 mt-l">
      <div class="xw-card soft-pink"><h4>{{卡标题}}</h4><p>{{内容}}</p></div>
      <div class="xw-card soft-blue"><h4>{{卡标题}}</h4><p>{{内容}}</p></div>
      <div class="xw-card soft-green"><h4>{{卡标题}}</h4><p>{{内容}}</p></div>
      <div class="xw-card soft-orange"><h4>{{卡标题}}</h4><p>{{内容}}</p></div>
    </div>
  </div>
  <div class="xw-footer">{{页脚}}</div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'xw-topline', 'xw-page', 'xw-title-md', 'xw-grid-2', 'mt-l', 'xw-card', 'soft-pink', 'soft-blue', 'soft-green', 'soft-orange', 'h4', 'xw-footer', 'notes'],
    },
    {
      id: 'steps', name: '步骤指南', roles: ['content'],
      use: '操作/方法步骤（xw-step 编号步进）',
      constraints: '3-5 步；每步 ≤22 字、动词开头',
      skeleton: `<section class="slide" data-layout="steps">
  <div class="xw-topline"></div>
  <div class="xw-page">
    <h2 class="xw-title-md">{{指南标题}}</h2>
    <div class="xw-steps mt-l">
      <div class="xw-step">{{第 1 步}}</div>
      <div class="xw-step">{{第 2 步}}</div>
      <div class="xw-step">{{第 3 步}}</div>
    </div>
  </div>
  <div class="xw-footer">{{页脚}}</div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'xw-topline', 'xw-page', 'xw-title-md', 'xw-steps', 'mt-l', 'xw-step', 'xw-footer', 'notes'],
    },
    {
      id: 'big-stat', name: '大数字', roles: ['data'],
      use: '一个关键数字/结论的杂志式呈现（渐变大字）',
      constraints: '数字 ≤8 字符；来源说明 ≤20 字，无出处标“估算”',
      skeleton: `<section class="slide" data-layout="big-stat">
  <div class="xw-topline"></div>
  <div class="xw-page">
    <p class="xw-kicker">{{指标语境}}</p>
    <div class="xw-big-stat xw-grad mt-l">{{数字}}</div>
    <h2 class="xw-title-md">{{指标名}}</h2>
    <p class="xw-sub">{{来源/口径说明}}</p>
  </div>
  <div class="xw-footer">{{页脚}}</div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'xw-topline', 'xw-page', 'xw-kicker', 'xw-big-stat', 'xw-grad', 'mt-l', 'xw-title-md', 'xw-sub', 'xw-footer', 'notes'],
    },
    {
      id: 'two-column', name: '双栏对照', roles: ['content'],
      use: '两个角度/方案的双栏卡片',
      constraints: '恰好 2 卡；内容 ≤34 字',
      skeleton: `<section class="slide" data-layout="two-column">
  <div class="xw-topline"></div>
  <div class="xw-page">
    <h2 class="xw-title-md">{{对照标题}}</h2>
    <div class="xw-grid-2 mt-l">
      <div class="xw-card soft-blue"><h4>{{栏 A}}</h4><p>{{内容}}</p></div>
      <div class="xw-card soft-purple"><h4>{{栏 B}}</h4><p>{{内容}}</p></div>
    </div>
  </div>
  <div class="xw-footer">{{页脚}}</div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'xw-topline', 'xw-page', 'xw-title-md', 'xw-grid-2', 'mt-l', 'xw-card', 'soft-blue', 'soft-purple', 'h4', 'xw-footer', 'notes'],
    },
    {
      id: 'quote', name: '金句页', roles: ['quote'],
      use: '一句金句/书中原话（xw-quote 衬线排版）',
      constraints: '引文 ≤34 字；出处真实，没有就不写出处行',
      skeleton: `<section class="slide" data-layout="quote">
  <div class="xw-topline"></div>
  <div class="xw-page">
    <div class="xw-quote">{{引文}}</div>
    <p class="xw-sub">{{—— 出处}}</p>
  </div>
  <div class="xw-footer">{{页脚}}</div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'xw-topline', 'xw-page', 'xw-quote', 'xw-sub', 'xw-footer', 'notes'],
    },
    {
      id: 'grid-3', name: '三栏要点', roles: ['content'],
      use: '三个并列要点的三栏卡（可混用马卡龙色）',
      constraints: '恰好 3 卡；内容 ≤26 字',
      skeleton: `<section class="slide" data-layout="grid-3">
  <div class="xw-topline"></div>
  <div class="xw-page">
    <h2 class="xw-title-md">{{主题}}</h2>
    <div class="xw-grid-3 mt-l">
      <div class="xw-card soft-blue"><h4>{{要点}}</h4><p>{{内容}}</p></div>
      <div class="xw-card soft-green"><h4>{{要点}}</h4><p>{{内容}}</p></div>
      <div class="xw-card soft-orange"><h4>{{要点}}</h4><p>{{内容}}</p></div>
    </div>
  </div>
  <div class="xw-footer">{{页脚}}</div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'xw-topline', 'xw-page', 'xw-title-md', 'xw-grid-3', 'mt-l', 'xw-card', 'soft-blue', 'soft-green', 'soft-orange', 'h4', 'xw-footer', 'notes'],
    },
  ],
  rules: `# editorial-white · 质量规则

- 刊头家具（xw-topline/xw-topbar/xw-footer）每页保留，页码写真实序号——这是"杂志"的仪式感。
- 马卡龙色卡的颜色是语义：pink=感受、blue=事实、green=行动、orange=注意，不要乱配。
- 标题用 xw-title（衬线大字），内容页标题用 xw-title-md，层级别混。
- 本模板亮色底：不放深色大色块；文字密度天然低，单页 ≤150 字。
- 引文必须真实（书/人说过的），编造的金句比没有金句更糟。
- 禁 emoji；图片位留 xw-hero 并在 notes 里说明需要什么图。
`,
}
