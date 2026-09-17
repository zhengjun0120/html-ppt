// course-module 契约规格
export default {
  id: 'course-module',
  name: '课程模块',
  description: '教学结构：学习目标、概念讲解、示例、随堂练习、选择自测、小结，培训首选',
  tags: ['课程', '培训', '教学'],
  scenario: ['课程', '培训', '教学', '工作坊'],
  source: 'course-module',
  scope: 'tpl-course-module',
  fonts: ['Inter', 'Noto Sans SC', 'JetBrains Maple Mono'],
  variants: [
    { id: 'default', name: '学院蓝' },
    { id: 'warm', name: '暖讲台', class: 'cm-warm',
      css: '.tpl-course-module.cm-warm{--accent:#ea580c}' },
    { id: 'forest', name: '森野', class: 'cm-forest',
      css: '.tpl-course-module.cm-forest{--accent:#16a34a}' },
  ],
  demoLayouts: ['cover', 'objectives', 'concept', 'example', 'exercise', 'check', 'summary'],
  layouts: [
    {
      id: 'cover', name: '课程封面', roles: ['cover'],
      use: '课程标题 + 侧栏课程导航出现',
      constraints: '标题 ≤14 字；副标 ≤20 字',
      skeleton: `<section class="slide" data-layout="cover">
  <p class="kicker">{{课程模块标签}}</p>
  <h1 class="h1 mt-s">{{课程标题}}</h1>
  <p class="dim mt-m">{{面向谁、学完能做什么}}</p>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kicker', 'h1', 'mt-s', 'dim', 'mt-m', 'notes'],
    },
    {
      id: 'objectives', name: '学习目标', roles: ['toc'],
      use: '本模块学完后能做到的 3-4 条目标',
      constraints: '3-4 条、每条动词开头 ≤20 字',
      skeleton: `<section class="slide" data-layout="objectives">
  <p class="kicker">objectives</p>
  <h2 class="h2 mt-s">学完这一模块，你会</h2>
  <div class="stack mt-l">
    <div class="callout">{{目标 1（动词开头）}}</div>
    <div class="callout">{{目标 2}}</div>
    <div class="callout">{{目标 3}}</div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kicker', 'h2', 'mt-s', 'stack', 'mt-l', 'callout', 'notes'],
    },
    {
      id: 'concept', name: '概念讲解', roles: ['content'],
      use: '核心概念 + 要点解释（concept-box 承载定义）',
      constraints: '定义 ≤40 字；要点 2-4 条、每条 ≤24 字',
      skeleton: `<section class="slide" data-layout="concept">
  <p class="kicker">concept</p>
  <h2 class="h2 mt-s">{{概念名}}</h2>
  <div class="concept-box mt-l">{{概念定义（一句话）}}</div>
  <div class="stack mt-l">
    <div class="callout">{{要点解释}}</div>
    <div class="callout">{{要点解释}}</div>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kicker', 'h2', 'mt-s', 'concept-box', 'mt-l', 'stack', 'callout', 'notes'],
    },
    {
      id: 'example', name: '示例讲解', roles: ['code', 'content'],
      use: '代码/操作示例（.code 块承载，语法高亮可选）',
      constraints: '代码 ≤16 行；配 1-2 句解说',
      skeleton: `<section class="slide" data-layout="example">
  <p class="kicker">example</p>
  <h2 class="h2 mt-s">{{示例标题}}</h2>
  <div class="code mt-m">
{{示例代码（可转义 &lt; &gt;）}}
  </div>
  <p class="dim mt-m">{{这段代码在做什么}}</p>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kicker', 'h2', 'mt-s', 'code', 'mt-m', 'dim', 'notes'],
    },
    {
      id: 'exercise', name: '随堂练习', roles: ['content', 'cta'],
      use: '学员动手的任务说明（exercise 容器）',
      constraints: '任务 ≤30 字；要求 2-3 条',
      skeleton: `<section class="slide" data-layout="exercise">
  <p class="kicker">exercise</p>
  <h2 class="h2 mt-s">动手练一练</h2>
  <div class="exercise mt-l">
    <p>{{任务描述}}</p>
    <p class="dim">{{要求/提示}}</p>
  </div>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kicker', 'h2', 'mt-s', 'exercise', 'mt-l', 'dim', 'notes'],
    },
    {
      id: 'check', name: '选择自测', roles: ['content'],
      use: '一道选择题（mcq 选项，correct 标正确项）',
      constraints: '恰好 1 题；题干 ≤30 字；4 个选项、每项 ≤20 字',
      skeleton: `<section class="slide" data-layout="check">
  <p class="kicker">quiz</p>
  <h2 class="h2 mt-s">{{题干}}</h2>
  <div class="stack mt-l">
    <div class="mcq correct">{{正确选项}}</div>
    <div class="mcq">{{干扰项}}</div>
    <div class="mcq">{{干扰项}}</div>
    <div class="mcq">{{干扰项}}</div>
  </div>
  <div class="notes">{{讲稿：为什么对、其余错在哪}}</div>
</section>`,
      classes: ['slide', 'kicker', 'h2', 'mt-s', 'stack', 'mt-l', 'mcq', 'correct', 'notes'],
    },
    {
      id: 'summary', name: '模块小结', roles: ['content', 'thanks'],
      use: '小结 + 下一模块预告',
      constraints: '3 条记忆点、每条 ≤20 字；预告 ≤16 字',
      skeleton: `<section class="slide" data-layout="summary">
  <p class="kicker">summary</p>
  <h2 class="h2 mt-s">{{小结标题}}</h2>
  <div class="stack mt-l">
    <div class="callout">{{记忆点 1}}</div>
    <div class="callout">{{记忆点 2}}</div>
    <div class="callout">{{记忆点 3}}</div>
  </div>
  <p class="dim mt-l">下一模块：{{预告}}</p>
  <div class="notes">{{讲稿}}</div>
</section>`,
      classes: ['slide', 'kicker', 'h2', 'mt-s', 'stack', 'mt-l', 'callout', 'dim', 'notes'],
    },
  ],
  rules: `# course-module · 质量规则

- 教学结构固定：目标 → 概念 → 示例 → 练习 → 自测 → 小结，顺序不可乱。
- concept-box 里只放定义（一句话），解释进 callout，别把定义写成一整段。
- mcq 的正确项用 correct 类标注（课堂揭晓态）；干扰项要“似是而非”，不出送分题。
- 代码块里的 HTML 必须转义（&lt; &gt; &amp;），行数 ≤16。
- 侧栏（sidebar）是模板骨架的一部分，页面内容不写进去。
- 禁 emoji；每页文字 ≤180 字。
`,
}
