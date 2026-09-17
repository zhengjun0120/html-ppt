// pitch-deck 契约规格
export default {
  id: 'pitch-deck',
  name: '融资路演',
  description: 'YC 风格 10 页路演结构：问题/方案/产品/市场/模式/ traction/团队/Ask，数字说话',
  tags: ['路演', '商业', '投资'],
  scenario: ['商业路演', '融资介绍', '创业比赛', '项目汇报'],
  source: 'pitch-deck',
  scope: 'tpl-pitch-deck',
  fonts: ['Inter', 'Noto Sans SC'],
  variants: [
    { id: 'default', name: '经典蓝' },
    { id: 'green', name: '生机绿', class: 'pd-green',
      css: '.tpl-pitch-deck.pd-green{--grad:linear-gradient(135deg,#0e9f6e,#3b82f6)}\n.tpl-pitch-deck.pd-green .ask-box{background:linear-gradient(135deg,#0e9f6e,#3b82f6)}' },
    { id: 'ember', name: '炽橙', class: 'pd-ember',
      css: '.tpl-pitch-deck.pd-ember{--grad:linear-gradient(135deg,#f97316,#db2777)}\n.tpl-pitch-deck.pd-ember .ask-box{background:linear-gradient(135deg,#f97316,#db2777)}' },
  ],
  demoLayouts: ['cover', 'problem-cards', 'pill-statement', 'feature-grid', 'market-metrics', 'business-model', 'traction-bars', 'team-cards', 'ask-box', 'thanks-mega'],
  layouts: [
    {
      id: 'cover', name: '路演封面', roles: ['cover'],
      use: '开场：项目一句话定位 + 讲者，渐变光斑背景',
      constraints: '主标题 ≤12 字；lede ≤24 字；kicker 是品牌短标签',
      skeleton: `<section class="slide" data-layout="cover">
  <p class="kicker">{{品牌标签}}</p>
  <h1 class="h1 anim-fade-up" data-anim="fade-up">{{主标题}}</h1>
  <p class="lede mt-m">{{一句话定位}}</p>
  <div class="deck-footer"><span class="mono">{{话题标签}}</span><span class="slide-number" data-current="1" data-total="{{总页数}}"></span></div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kicker', 'h1', 'anim-fade-up', 'lede', 'mt-m', 'deck-footer', 'mono', 'slide-number', 'notes'],
    },
    {
      id: 'problem-cards', name: '问题三卡', roles: ['content'],
      use: '三个并列痛点/问题，开头编号引导',
      constraints: '恰好 3 卡；卡标题 ≤10 字；说明 ≤36 字',
      skeleton: `<section class="slide" data-layout="problem-cards">
  <span class="section-num">{{章节数字}}</span>
  <p class="num-tag">{{引导语}}</p>
  <h2 class="h2 mt-s">{{问题陈述}}</h2>
  <div class="grid g3 mt-l">
    <div class="card"><h4>{{痛点1}}</h4><p class="dim">{{说明}}</p></div>
    <div class="card"><h4>{{痛点2}}</h4><p class="dim">{{说明}}</p></div>
    <div class="card"><h4>{{痛点3}}</h4><p class="dim">{{说明}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'section-num', 'num-tag', 'h2', 'mt-s', 'grid', 'g3', 'mt-l', 'card', 'h4', 'dim', 'notes'],
    },
    {
      id: 'pill-statement', name: '方案宣言', roles: ['content', 'divider'],
      use: '方案核心论点 + 关键词 pill 一排',
      constraints: '宣言 ≤28 字；pill 3-5 个、每个 ≤6 字',
      skeleton: `<section class="slide" data-layout="pill-statement">
  <span class="section-num">{{章节数字}}</span>
  <p class="num-tag">{{引导语}}</p>
  <h2 class="h2 mt-s">{{方案宣言}}</h2>
  <p class="lede mt-m">{{补充一句}}</p>
  <div class="row mt-l">
    <span class="pill pill-accent">{{关键词}}</span>
    <span class="pill pill-accent">{{关键词}}</span>
    <span class="pill">{{关键词}}</span>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'section-num', 'num-tag', 'h2', 'mt-s', 'lede', 'mt-m', 'row', 'mt-l', 'pill', 'pill-accent', 'notes'],
    },
    {
      id: 'feature-grid', name: '产品特性四宫格', roles: ['content'],
      use: '产品 4 个能力点（卡带悬浮态）',
      constraints: '恰好 4 卡；标题 ≤8 字；说明 ≤28 字',
      skeleton: `<section class="slide" data-layout="feature-grid">
  <span class="section-num">{{章节数字}}</span>
  <p class="num-tag">{{引导语}}</p>
  <h2 class="h2 mt-s">{{产品能力标题}}</h2>
  <div class="grid g2 mt-l">
    <div class="card card-hover"><h4>{{能力}}</h4><p class="dim">{{说明}}</p></div>
    <div class="card card-hover"><h4>{{能力}}</h4><p class="dim">{{说明}}</p></div>
    <div class="card card-hover"><h4>{{能力}}</h4><p class="dim">{{说明}}</p></div>
    <div class="card card-hover"><h4>{{能力}}</h4><p class="dim">{{说明}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'section-num', 'num-tag', 'h2', 'mt-s', 'grid', 'g2', 'mt-l', 'card', 'card-hover', 'h4', 'dim', 'notes'],
    },
    {
      id: 'market-metrics', name: '市场三指标', roles: ['data'],
      use: '市场规模 TAM/SAM/SOM 或三个核心数字',
      constraints: '3 个指标；指标名 ≤8 字；数字 ≤6 字符；说明 ≤24 字（无出处标“估算”）',
      skeleton: `<section class="slide" data-layout="market-metrics">
  <span class="section-num">{{章节数字}}</span>
  <p class="num-tag">{{引导语}}</p>
  <h2 class="h2 mt-s">{{市场标题}}</h2>
  <div class="grid g3 mt-l">
    <div class="metric"><span class="l">{{指标名}}</span><span class="n">{{数字}}</span><p class="dim">{{说明}}</p></div>
    <div class="metric"><span class="l">{{指标名}}</span><span class="n">{{数字}}</span><p class="dim">{{说明}}</p></div>
    <div class="metric"><span class="l">{{指标名}}</span><span class="n">{{数字}}</span><p class="dim">{{说明}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'section-num', 'num-tag', 'h2', 'mt-s', 'grid', 'g3', 'mt-l', 'metric', 'l', 'n', 'dim', 'notes'],
    },
    {
      id: 'business-model', name: '商业模式', roles: ['content'],
      use: '怎么赚钱：定价/模式三列或两列卡',
      constraints: '2-3 列；列标题 ≤8 字；说明 ≤32 字',
      skeleton: `<section class="slide" data-layout="business-model">
  <span class="section-num">{{章节数字}}</span>
  <p class="num-tag">{{引导语}}</p>
  <h2 class="h2 mt-s">{{模式标题}}</h2>
  <div class="grid g3 mt-l">
    <div class="card card-accent"><h4>{{模式1}}</h4><p class="dim">{{说明}}</p></div>
    <div class="card card-accent"><h4>{{模式2}}</h4><p class="dim">{{说明}}</p></div>
    <div class="card card-accent"><h4>{{模式3}}</h4><p class="dim">{{说明}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'section-num', 'num-tag', 'h2', 'mt-s', 'grid', 'g3', 'mt-l', 'card', 'card-accent', 'h4', 'dim', 'notes'],
    },
    {
      id: 'traction-bars', name: '增长势头', roles: ['data'],
      use: 'traction 横向条形（用 .bar 宽度表达相对大小）',
      constraints: '4-6 条；条目名 ≤10 字；宽度按真实比例给 style，数字标注在条目名后',
      skeleton: `<section class="slide" data-layout="traction-bars">
  <span class="section-num">{{章节数字}}</span>
  <p class="num-tag">{{引导语}}</p>
  <h2 class="h2 mt-s">{{增长标题}}</h2>
  <div class="traction-bar mt-l">
    <div class="bar" style="width:90%"><span>{{条目名 数字}}</span></div>
    <div class="bar" style="width:60%"><span>{{条目名 数字}}</span></div>
    <div class="bar" style="width:35%"><span>{{条目名 数字}}</span></div>
  </div>
  <p class="dim mt-l">{{口径说明（无来源标“估算”）}}</p>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'section-num', 'num-tag', 'h2', 'mt-s', 'traction-bar', 'mt-l', 'bar', 'dim', 'notes'],
    },
    {
      id: 'team-cards', name: '团队三卡', roles: ['content'],
      use: '核心成员 3 人（头像位 + 名字 + 背景）',
      constraints: '恰好 3 卡；姓名 ≤6 字；背景 ≤30 字',
      skeleton: `<section class="slide" data-layout="team-cards">
  <span class="section-num">{{章节数字}}</span>
  <p class="num-tag">{{引导语}}</p>
  <h2 class="h2 mt-s">{{团队标题}}</h2>
  <div class="grid g3 mt-l">
    <div class="card team-card"><div class="avatar"></div><h4>{{姓名}}</h4><p class="dim">{{背景}}</p></div>
    <div class="card team-card"><div class="avatar"></div><h4>{{姓名}}</h4><p class="dim">{{背景}}</p></div>
    <div class="card team-card"><div class="avatar"></div><h4>{{姓名}}</h4><p class="dim">{{背景}}</p></div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'section-num', 'num-tag', 'h2', 'mt-s', 'grid', 'g3', 'mt-l', 'card', 'team-card', 'avatar', 'h4', 'dim', 'notes'],
    },
    {
      id: 'ask-box', name: 'The Ask', roles: ['cta'],
      use: '要什么：金额/资源/支持，一页只说这一件事',
      constraints: 'Ask 一句话 ≤30 字；支持点 ≤3 条、每条 ≤14 字',
      skeleton: `<section class="slide" data-layout="ask-box">
  <p class="num-tag">{{引导语}}</p>
  <div class="ask-box mt-m">
    <h2 class="h2">{{Ask 一句话}}</h2>
    <p class="lede">{{用途说明}}</p>
  </div>
  <div class="row mt-l">
    <div class="dim">{{支持点}}</div>
    <div class="dim">{{支持点}}</div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'num-tag', 'ask-box', 'mt-m', 'h2', 'lede', 'row', 'mt-l', 'dim', 'notes'],
    },
    {
      id: 'thanks-mega', name: '致谢大字', roles: ['thanks', 'cta'],
      use: '收尾致谢 + 联系方式 pill',
      constraints: '大字 ≤8 字；lede ≤20 字；pill ≤2 个',
      skeleton: `<section class="slide center tc" data-layout="thanks-mega">
  <div class="mega">{{大字}}</div>
  <p class="mega-sub">{{一句话}}</p>
  <div class="row mt-l" style="justify-content:center">
    <span class="pill pill-accent">{{联系方式}}</span>
    <span class="pill">{{仓库/链接}}</span>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'center', 'tc', 'mega', 'mega-sub', 'row', 'mt-l', 'pill', 'pill-accent', 'notes'],
    },
  ],
  rules: `# pitch-deck · 质量规则

- **数字优先**：路演页每个论点尽量挂一个数字；没有来源就用量级并标“估算”，不编精确数。
- section-num/num-tag 是本模板的章节引导家具：每个内容页页首一组，数字用章节序（01-06），引导语 ≤8 字。
- traction-bars 的 .bar 宽度是数据可视化：宽度必须与数字成比例，别乱给。
- 团队卡 .avatar 是渐变圆位，不放真实照片；名字不虚构全名，用“创始人 / 技术合伙人”角色称谓即可。
- 禁 emoji；禁编造融资金额与估值——模板里的 Ask 数字必须来自用户输入或标“目标”。
- Ask 页（ask-box）整份 deck 最多 1 页，放在最后三分之一。
- 密度：单页可见文字 ≤180 字；卡内说明 1-2 句。
`,
}
