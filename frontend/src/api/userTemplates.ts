import { authedUrl, currentToken, request } from './client'
import { iterateSse } from '@/lib/sse'
import type { TemplateVariant } from './templates'

/**
 * 用户自定义模板（plan-v3 B）：fork → 定制 → 门禁发布 → 社区使用。
 * 行类型与 store.UserTemplate 对齐；canvas/variants 由后端从注册表补齐
 * （画廊预览卡需要画布定比例、变体做换肤）。
 */
export interface UserTemplateRow {
  id: string
  user_id: number
  base_id: string
  name: string
  description: string
  visibility: 'private' | 'public'
  status: 'draft' | 'publishing' | 'published' | 'failed'
  publish_error?: string
  publish_report?: string
  /** 缩略图内容版本（8 位）：卡片 <img> 地址 /api/templates/:id/thumb?v= 的 cache-bust 键 */
  thumb?: string
  created_at: string
  updated_at: string
  canvas?: { w: number; h: number }
  variants?: TemplateVariant[]
}

export interface CommunityTemplate {
  id: string
  name: string
  description: string
  base_id: string
  author: string
  updated_at: string
  canvas?: { w: number; h: number }
  variants?: TemplateVariant[]
  /** 封面缩略图的内容版本前缀，拼 templateApi.thumbUrl 用 */
  thumb?: string
  /** demo 页数（预览弹窗的"第 X / N 页"） */
  demo_pages?: number
}

export interface PublishReport {
  structure: string
  render?: { pages: number; min_fill: number; max_flag_font: number; flaws?: string[] }
  note?: string
}

/** 一条模板历史版本（后端 usertpl.UTVersionMeta；changed=与上一版相比变动的文件） */
export interface UTVersionMeta {
  version: string
  time: number
  operation: 'fork' | 'chat' | 'edit' | 'meta' | 'restore'
  detail: string
  changed?: string[]
}

export const userTemplateApi = {
  fork: (baseId: string, name = '') =>
    request<UserTemplateRow>(`/api/templates/${baseId}/fork`, {
      method: 'POST',
      body: { name },
    }),
  /** 从空白脚手架新建（后端约定 baseId="_blank" 走 CreateBlank，不派生任何内置模板） */
  blank: (name = '') => request<UserTemplateRow>(`/api/templates/_blank/fork`, {
    method: 'POST',
    body: { name },
  }),
  list: () => request<UserTemplateRow[]>('/api/user-templates'),
  get: (id: string) => request<UserTemplateRow>(`/api/user-templates/${id}`),
  updateMeta: (id: string, name: string, description = '') =>
    request<{ ok: boolean }>(`/api/user-templates/${id}`, {
      method: 'PUT',
      body: { name, description },
    }),
  remove: (id: string) => request<{ ok: boolean }>(`/api/user-templates/${id}`, { method: 'DELETE' }),
  /** 编辑器全量保存 index.html（raw HTML；服务端记 edit 版本，滚动备份已退役） */
  saveFile: (id: string, html: string) =>
    request<{ ok: boolean }>(`/api/user-templates/${id}/file`, {
      method: 'PUT',
      body: html,
      raw: true,
    }),
  /** 手动保存 style.css（raw CSS；安全预检 + 记 edit 版本；published 409） */
  saveStyle: (id: string, css: string) =>
    request<{ ok: boolean }>(`/api/user-templates/${id}/style`, {
      method: 'PUT',
      body: css,
      raw: true,
    }),
  // —— 历史版本（docs/user-template-history-plan.md §4）——//
  history: (id: string) => request<UTVersionMeta[]>(`/api/user-templates/${id}/history`),
  restoreVersion: (id: string, version: string) =>
    request<{ ok: boolean }>(`/api/user-templates/${id}/history/${version}/restore`, { method: 'POST' }),
  deleteVersion: (id: string, version: string) =>
    request<{ ok: boolean }>(`/api/user-templates/${id}/history/${version}`, { method: 'DELETE' }),
  clearHistory: (id: string) =>
    request<{ deleted: number }>(`/api/user-templates/${id}/history`, { method: 'DELETE' }),
  publish: (id: string) => request<PublishReport>(`/api/user-templates/${id}/publish`, { method: 'POST' }),
  unpublish: (id: string) => request<{ ok: boolean }>(`/api/user-templates/${id}/unpublish`, { method: 'POST' }),
  /** 质量体检（渲染量测：溢出/填充率/最小字号）。只报告不拦发布，draft/published 都能跑 */
  checkup: (id: string) => request<PublishReport>(`/api/user-templates/${id}/checkup`, { method: 'POST' }),
  community: () => request<CommunityTemplate[]>('/api/community-templates'),
}

/** 读当前 style.css 全文（样式面板；raw 文本，走 request 之外的白名单资产端点） */
export async function fetchTemplateStyle(id: string): Promise<string> {
  const resp = await fetch(authedUrl(`/api/user-templates/${id}/assets/style.css`))
  if (!resp.ok) throw new Error('样式读取失败')
  return resp.text()
}

// —— 定制对话（plan-v3 B2；2026-09-28 起走 SSE 流式）——//

/** 定制对话的 SSE 事件（后端 usertpl.CustEvent，窄集合） */
export interface CustomizeEvent {
  type: 'tool_start' | 'tool_progress' | 'tool_done' | 'delta' | 'done' | 'error'
  content?: string // delta 文本 / tool_done 结果摘要 / error 消息
  tool_name?: string
  tool_call_id?: string
  tool_index?: number
  bytes?: number // tool_progress：该工具已生成的参数字节数
  dirty?: boolean // done：本轮是否有文件写入
  reply?: string // done：最终答复全文
}

const CUSTOMIZE_BASE: string = import.meta.env.VITE_API_BASE ?? ''

/** 流式定制对话。onEvent 逐帧回调；调用方负责错误展示与 finally 复位。 */
export async function customizeChatStream(
  id: string,
  message: string,
  onEvent: (ev: CustomizeEvent) => void,
  signal?: AbortSignal,
): Promise<void> {
  const token = currentToken()
  const res = await fetch(`${CUSTOMIZE_BASE}/api/user-templates/${id}/chat/stream`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    body: JSON.stringify({ message }),
    signal,
  })
  if (!res.ok) {
    let msg = `请求失败（${res.status}）`
    try {
      const data = (await res.json()) as { error?: string }
      if (data?.error) msg = data.error
    } catch { /* 非 JSON body */ }
    throw new Error(msg)
  }
  for await (const frame of iterateSse(res)) {
    try {
      onEvent(JSON.parse(frame.data) as CustomizeEvent)
    } catch {
      // 单帧坏了跳过，不让一行脏数据打断整轮展示
    }
  }
}

/** 同步版定制对话（保留给降级路径；前端正常走 customizeChatStream） */
export function customizeChat(id: string, message: string) {
  return request<{ reply: string }>(`/api/user-templates/${id}/chat`, {
    method: 'POST',
    body: { message },
  })
}

/** demo 预览地址：鉴权 demo 端点（草稿收口后不再走公开静态；token 由 authedUrl 带） */
export function userTemplatePreviewUrl(id: string, page = 0) {
  return authedUrl(`/api/user-templates/${id}/demo`) + (page > 0 ? `#/${page}` : '')
}
