/**
 * agent 工具的展示元数据：中文名、参数的一句中文摘要、参数/结果的美化。
 * 工具卡（ToolCard）与相邻同类折叠组（ToolGroupCard）共用。
 */

export const TOOL_LABELS: Record<string, string> = {
  plan_pages: '全局规划',
  write_pages: '写入页面',
  update_slide: '修改页面',
  insert_slide: '插入页面',
  delete_slide: '删除页面',
  list_slides: '列出页面',
  read_slide: '读取页面',
  read_guidelines: '读取设计规范',
  read_layout: '读取版式',
  review_slides: '看图审查',
  list_history: '版本历史',
  read_history_diff: '版本对比',
  read_outline: '读取大纲',
  update_outline: '修改大纲',
  submit_outline: '提交大纲',
  web_search: '联网搜索',
  ask_user: '向用户提问',
  write_tokens: '写入样式令牌',
  set_meta: '更新模板信息',
  finish: '结束定制',
}

export function toolLabel(name: string): string {
  return TOOL_LABELS[name] ?? name
}

function clip(s: string, n: number): string {
  const r = Array.from(s)
  return r.length <= n ? s : r.slice(0, n).join('') + '…'
}

/**
 * 参数的一句中文摘要（显示在工具卡头部）。JSON 没解析出来（流式中途 /
 * 形状不符）就返回空串，卡片自己退回通用展示。
 */
export function toolSummary(name: string, argsJson: string): string {
  let a: Record<string, unknown>
  try {
    a = JSON.parse(argsJson) as Record<string, unknown>
  } catch {
    return ''
  }
  try {
    switch (name) {
      case 'write_pages': {
        const pages = Array.isArray(a.pages) ? (a.pages as { title?: string }[]) : []
        const titles = pages.map((p) => p?.title).filter((t): t is string => !!t)
        const head = titles.slice(0, 3).join('、')
        return `写入 ${pages.length} 页${head ? `：${head}${titles.length > 3 ? '…' : ''}` : ''}`
      }
      case 'update_slide':
        return `修改页 ${String(a.slide_id ?? '?')}${typeof a.reason === 'string' && a.reason ? `：${clip(a.reason, 40)}` : ''}`
      case 'insert_slide':
        return `在 ${String(a.after_slide_id ?? '末尾')} 后插入新页`
      case 'delete_slide':
        return `删除页 ${String(a.slide_id ?? '?')}`
      case 'read_slide':
        return `读取页 ${String(a.slide_id ?? '?')}`
      case 'read_layout':
        return `版式 ${String(a.layout_id ?? '?')}`
      case 'web_search':
        return typeof a.query === 'string' ? clip(a.query, 50) : ''
      case 'read_history_diff':
        return `对比 ${String(a.from_version || '最新')} → ${String(a.to_version || '当前')}`
      case 'plan_pages':
        return '全局版式规划'
      default:
        return ''
    }
  } catch {
    return ''
  }
}

/**
 * 参数美化：完整 JSON 缩进两格展示（超长字符串值截断留字数）；
 * 流式中的半截 JSON 解析失败就原样返回——正好能看到增量长大。
 */
export function prettyArgs(argsJson: string): string {
  try {
    const parsed: unknown = JSON.parse(argsJson)
    return JSON.stringify(parsed, (_k, v) => {
      if (typeof v === 'string' && v.length > 300) {
        return v.slice(0, 300) + `…（共 ${v.length} 字）`
      }
      return v
    }, 2)
  } catch {
    return argsJson
  }
}

/** 结果美化：JSON 形状的结果（版本历史/版本对比等）格式化缩进，其余原样 */
export function prettyResult(result: string): string {
  const head = result.trimStart()
  if (head.startsWith('{') || head.startsWith('[')) {
    try {
      return JSON.stringify(JSON.parse(result), null, 2)
    } catch {
      /* 非 JSON，原样 */
    }
  }
  return result
}
