// product-launch 契约规格
export default {
  id: 'product-launch',
  name: '产品发布',
  description: '发布会节奏：暗色 hero 开场、产品特性三连、定价卡、发布宣言',
  tags: ['产品', '发布', '营销'],
  scenario: ['产品发布', '版本发布', '功能介绍', '客户演示'],
  source: 'product-launch',
  scope: 'tpl-product-launch',
  fonts: ['Inter', 'Noto Sans SC'],
  variants: [
    { id: 'default', name: '经典暗色' },
    { id: 'violet', name: '幻紫', class: 'pl-violet',
      css: '.tpl-product-launch.pl-violet{--accent:#a78bfa}' },
    { id: 'cyan', name: '电光青', class: 'pl-cyan',
      css: '.tpl-product-launch.pl-cyan{--accent:#22d3ee}' },
  ],
  demoLayouts: ['cover', 'introducing', 'feature-trio', 'fit-cards', 'feature-duo', 'how-it-works', 'pricing', 'ship'],
  layouts: [
    {
      id: 'cover', name: '发布封面', roles: ['cover'],
      use: '暗色 hero：产品名大标题 + 一句话主张',
      constraints: '标题 ≤10 字；lede ≤20 字；kicker 是发布主题标签',
      skeleton: `<section class="slide dark" data-layout="cover">
  <p class="kicker">{{发布标签}}</p>
  <h1 class="h1 anim-fade-up" data-anim="fade-up">{{产品名/主张}}</h1>
  <p class="lede mt-m">{{一句话主张}}</p>
  <div class="deck-footer"><span class="brand">{{品牌}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'dark', 'kicker', 'h1', 'anim-fade-up', 'lede', 'mt-m', 'deck-footer', 'brand', 'slide-number', 'notes'],
    },
    {
      id: 'introducing', name: '发布宣言', roles: ['divider'],
      use: '居中宣言页：新版本来袭，一页只说一句话',
      constraints: '宣言 ≤12 字；副句 ≤20 字',
      skeleton: `<section class="slide center tc" data-layout="introducing">
  <p class="kicker">{{引导语}}</p>
  <h1 class="h1">{{宣言}}</h1>
  <p class="lede">{{副句}}</p>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'center', 'tc', 'kicker', 'h1', 'lede', 'notes'],
    },
    {
      id: 'feature-trio', name: '特性三连', roles: ['content'],
      use: '三个核心特性（feature-card 逐卡）',
      constraints: '恰好 3 卡；卡标题 ≤8 字；说明 ≤30 字',
      skeleton: `<section class="slide" data-layout="feature-trio">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{特性主题}}</h2>
  <div class="grid g3 mt-l">
    <div class="feature-card"><h4>{{特性}}</h4><p class="dim">{{说明}}</p></div>
    <div class="feature-card"><h4>{{特性}}</h4><p class="dim">{{说明}}</p></div>
    <div class="feature-card"><h4>{{特性}}</h4><p class="dim">{{说明}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kicker', 'h2', 'grid', 'g3', 'mt-l', 'feature-card', 'h4', 'dim', 'notes'],
    },
    {
      id: 'fit-cards', name: '场景适配', roles: ['content'],
      use: '适用场景/人群三卡（暗色页变体）',
      constraints: '恰好 3 卡；说明 ≤28 字',
      skeleton: `<section class="slide dark" data-layout="fit-cards">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{场景标题}}</h2>
  <div class="grid g3 mt-l">
    <div class="card"><h4>{{场景}}</h4><p class="dim">{{说明}}</p></div>
    <div class="card"><h4>{{场景}}</h4><p class="dim">{{说明}}</p></div>
    <div class="card"><h4>{{场景}}</h4><p class="dim">{{说明}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'dark', 'kicker', 'h2', 'grid', 'g3', 'mt-l', 'card', 'h4', 'dim', 'notes'],
    },
    {
      id: 'feature-duo', name: '深度两卡', roles: ['content'],
      use: '两个重点能力各占一卡，讲深一点',
      constraints: '恰好 2 卡；卡标题 ≤10 字；说明 ≤48 字',
      skeleton: `<section class="slide" data-layout="feature-duo">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{能力标题}}</h2>
  <div class="grid g2 mt-l">
    <div class="feature-card"><h4>{{能力}}</h4><p class="dim">{{说明}}</p></div>
    <div class="feature-card"><h4>{{能力}}</h4><p class="dim">{{说明}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kicker', 'h2', 'grid', 'g2', 'mt-l', 'feature-card', 'h4', 'dim', 'notes'],
    },
    {
      id: 'how-it-works', name: '三步上手', roles: ['content'],
      use: '使用流程三步（step 逐行）',
      constraints: '恰好 3 步；每步 ≤20 字、动词开头',
      skeleton: `<section class="slide" data-layout="how-it-works">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{流程标题}}</h2>
  <div class="stack mt-l">
    <div class="step">{{第 1 步}}</div>
    <div class="step">{{第 2 步}}</div>
    <div class="step">{{第 3 步}}</div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kicker', 'h2', 'stack', 'mt-l', 'step', 'notes'],
    },
    {
      id: 'pricing', name: '定价卡', roles: ['content', 'cta'],
      use: '定价三档（price-card，可标推荐档）',
      constraints: '恰好 3 档；档名 ≤6 字；价格真实，没有定价就写“联系我们”',
      skeleton: `<section class="slide" data-layout="pricing">
  <p class="kicker">{{引导语}}</p>
  <h2 class="h2">{{定价标题}}</h2>
  <div class="grid g3 mt-l">
    <div class="price-card"><span class="amount">{{价格}}</span><h4>{{档名}}</h4><ul><li>{{包含内容}}</li></ul></div>
    <div class="price-card"><span class="amount">{{价格}}</span><h4>{{档名}}</h4><ul><li>{{包含内容}}</li></ul></div>
    <div class="price-card"><span class="amount">{{价格}}</span><h4>{{档名}}</h4><ul><li>{{包含内容}}</li></ul></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kicker', 'h2', 'grid', 'g3', 'mt-l', 'price-card', 'amount', 'h4', 'notes'],
    },
    {
      id: 'ship', name: '发布号召', roles: ['cta', 'thanks'],
      use: '收尾号召：现在就用/扫码/链接',
      constraints: '号召 ≤10 字；lede ≤20 字',
      skeleton: `<section class="slide center tc" data-layout="ship">
  <p class="kicker">{{发布日期}}</p>
  <h1 class="h1">{{号召语}}</h1>
  <p class="lede">{{获取方式}}</p>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'center', 'tc', 'kicker', 'h1', 'lede', 'notes'],
    },
  ],
  rules: `# product-launch · 质量规则

- .dark 页是节奏锚点：cover/宣言/场景页用，连续不超过 2 页。
- 定价必须是真实价格或明确写“发布后公布”，不编造。
- 特性说明写用户收益（“上传即转写”），不写技术参数堆砌；参数放讲稿。
- hero 页（cover/introducing/ship）文字极少是大胆的正确的，别往里塞内容。
- 禁 emoji；单页文字 ≤140 字（本模板刻意轻密度）。
`,
}
