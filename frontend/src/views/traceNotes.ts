// 观测台注解表：给事件里的工具、字段、消息角色、事件类型、子步骤阶段配
// 中文名与一句话说明。口径来自后端工具定义（func_tool_v2.go 的 mountTool
// 描述与 jsonschema 标签）——改工具定义时这里要跟着改；没收录的键只显示
// 原名，不影响渲染。

export interface FieldNote {
  zh: string
  desc?: string
}

// 工具名 → 中文名 + 用途。口径：func_tool_v2.go / func_tool.go 的 mountTool 描述。
export const toolNotes: Record<string, FieldNote> = {
  ask_user: { zh: '向用户提问', desc: '一次问 1-6 个问题，问完本轮暂停等用户回答' },
  web_search: { zh: '联网搜索', desc: 'DeepSeek 服务端搜索，会真的抓网页；分开计费' },
  review_slides: { zh: '看图审查', desc: '渲染整份 deck，可点名几页真正看图，返回量测数字与审查报告' },
  list_slides: { zh: '查看页面目录', desc: '列已写入的页面与版式，自查页序用' },
  read_slide: { zh: '读取页面', desc: '读某页当前 HTML 与指纹，改页前先读' },
  update_slide: { zh: '修改页面', desc: '整页替换某页 HTML，需带上一次的指纹防并发覆盖' },
  submit_outline: { zh: '提交大纲', desc: '提交结构化大纲并创建文稿' },
  read_outline: { zh: '读取大纲', desc: '读当前大纲全文与版本号，改大纲前必须先读' },
  update_outline: { zh: '修改大纲', desc: '整份替换大纲，带版本号做冲突检测' },
  plan_pages: { zh: '全局规划', desc: '为大纲每一页分配版式，节奏违规会被打回重排' },
  read_layout: { zh: '取版式骨架', desc: '取某版式的 HTML 骨架与合法类名，写页前必读' },
  read_guidelines: { zh: '取模板质量规则', desc: '取模板专属的写作与排版规约全文' },
  write_pages: { zh: '批量写页', desc: '分批写入页面（每批 2-4 页），返回量测数字' },
  insert_slide: { zh: '插入新页', desc: '在某页之后插入新页' },
  delete_slide: { zh: '删除页面', desc: '删除一页（至少保留一页）' },
  list_history: { zh: '查看版本历史', desc: '列版本号（新→旧），供比较差异用' },
  read_history_diff: { zh: '比较版本差异', desc: '比较两份状态的页面差异，自查改动用' },
}

// 字段名 → 中文名 + 一句话说明。覆盖工具参数与结果里出现过的键；
// 同名键在不同工具里含义相近时取通用口径。
export const fieldNotes: Record<string, FieldNote> = {
  // —— 通用 ——
  deck_id: { zh: '文稿 ID', desc: '操作目标文稿' },
  slide_id: { zh: '页面 ID', desc: '来自 list_slides，形如 s3' },
  title: { zh: '标题' },
  version: { zh: '版本号', desc: '乐观锁：之前读到的版本号原样传入，过期会报冲突' },
  limit: { zh: '条数上限' },
  offset: { zh: '翻页偏移', desc: '跳过最近 N 条，看更早的记录' },
  query: { zh: '搜索问题', desc: '写成一句完整的话最准（带年份、地区、主体）' },
  max_uses: { zh: '搜索轮数上限', desc: '默认 3、最大 8；子模型可能不严格遵守' },
  focus: { zh: '审查关注点', desc: '这次看图特别想验证的问题，一句话写具体' },
  // —— 大纲字段 ——
  audience: { zh: '受众与场合' },
  duration_min: { zh: '预计时长（分钟）' },
  tone: { zh: '语气基调' },
  hook: { zh: '叙事钩子', desc: '一句话说清为什么值得听' },
  arcs: { zh: '叙事分段', desc: '主体推进的各段名' },
  pages: { zh: '页面清单', desc: '逐页条目（大纲 / 本批写入 / 点名看图）' },
  no: { zh: '页码', desc: '1 基' },
  role: { zh: '页面角色', desc: 'cover/toc/divider/content/data/quote/code/cta/thanks' },
  points: { zh: '内容要点', desc: '2-5 条、每条一句话' },
  layout_hint: { zh: '版式建议' },
  materials: { zh: '素材清单', desc: '需要用户提供的图或数据' },
  notes: { zh: '备注' },
  assignments: { zh: '版式分配', desc: '每页一条，必须覆盖大纲全部页' },
  reason: { zh: '选择理由' },
  // —— 写页 / 改页 ——
  html: { zh: '页面 HTML', desc: '基于版式骨架填充的完整 <section>' },
  new_html: { zh: '新页面 HTML', desc: '替换后的完整 <section>，保留 data-layout' },
  fingerprint: { zh: '页面指纹', desc: 'read_slide 返回原样传入，防并发覆盖' },
  after_slide_id: { zh: '插入位置', desc: '某页 slide_id 或 end（末尾）' },
  questions: { zh: '问题清单', desc: '要问用户的问题，1-6 个' },
  from_version: { zh: '基线版本号', desc: '不传 = 最新记档版本' },
  to_version: { zh: '目标版本号', desc: '不传 = 当前使用中的内容' },
  // —— 工具结果字段 ——
  rules: { zh: '质量规则', desc: '模板专属的写作与排版规约，生成阶段注入' },
  classes: { zh: '可用类名', desc: '该版式允许使用的 CSS 类清单，只准用这些' },
  skeleton: { zh: '版式骨架', desc: '版式的 HTML 骨架，把占位符换成真实内容' },
  repeats: { zh: '数量契约', desc: '类名 → 应出现的个数（如 ac-item: 4）' },
  template: { zh: '模板 ID' },
  layout: { zh: '版式 ID' },
  answer: { zh: '搜索答案', desc: '读过网页后写的答案，与既有知识冲突时以它为准' },
  sources: { zh: '来源链接', desc: '标题与 URL，写进页面的小字引用' },
  searches: { zh: '计费搜索次数', desc: '这次调用实际发生的搜索轮数' },
  truncated: { zh: '是否被截断', desc: 'true 表示答案或来源被截断过，别当完整信息' },
  tool_error: { zh: '工具自身报错', desc: '有它说明这次搜索没成功' },
}

// 消息角色 → 说明。
export const roleNotes: Record<string, string> = {
  system: '系统提示（任务规则与背景，每轮都带）',
  user: '用户输入',
  assistant: '模型回复',
  tool: '工具结果（回填给模型）',
}

// 事件类型 → 中文名。
export const kindNotes: Record<string, string> = {
  run_start: '开始',
  llm_request: '模型请求',
  llm_response: '模型返回',
  tool_call: '工具调用',
  tool_result: '工具结果',
  sub_step: '子步骤',
  usage: '用量',
  error: '出错',
  run_end: '结束',
}

// 子步骤的来源与阶段 → 中文。口径：vision_review_v2.go / web_search.go 的 emit 点。
const subNameMap: Record<string, string> = {
  vision: '视觉审查',
  web_search: '联网搜索',
}

const subStageMap: Record<string, string> = {
  'vision.grant': '发票据',
  'vision.capture': '实拍量测',
  'vision.review_prompt': '审查提示词',
  'vision.review_response': '审查报告',
  'vision.error': '出错',
  'web_search.request': '搜索请求',
  'web_search.queries': '实际搜索词',
  'web_search.response': '响应',
  'web_search.error': '出错',
}

export function subNameZh(name: string): string {
  return subNameMap[name] ?? ''
}

export function subStageZh(name: string, stage: string): string {
  return subStageMap[`${name}.${stage}`] ?? ''
}
