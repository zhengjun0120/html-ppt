import { request } from './client'

/**
 * 用户自定义模板（plan-v3 B）：fork → 定制 → 门禁发布 → 社区使用。
 * 行类型与 store.UserTemplate 对齐。
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
  created_at: string
  updated_at: string
}

export interface CommunityTemplate {
  id: string
  name: string
  description: string
  base_id: string
  author: string
  updated_at: string
}

export interface PublishReport {
  structure: string
  render?: { pages: number; min_fill: number; max_flag_font: number; flaws?: string[] }
  note?: string
}

export const userTemplateApi = {
  fork: (baseId: string, name = '') =>
    request<UserTemplateRow>(`/api/templates/${baseId}/fork`, {
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
  publish: (id: string) => request<PublishReport>(`/api/user-templates/${id}/publish`, { method: 'POST' }),
  unpublish: (id: string) => request<{ ok: boolean }>(`/api/user-templates/${id}/unpublish`, { method: 'POST' }),
  community: () => request<CommunityTemplate[]>('/api/community-templates'),
}

// —— 定制对话（plan-v3 B2）——//
export function customizeChat(id: string, message: string) {
  return request<{ reply: string }>(`/api/user-templates/${id}/chat`, {
    method: 'POST',
    body: { message },
  })
}

/** demo 预览地址（公开静态；定制/预览共用） */
export function userTemplatePreviewUrl(id: string, page = 0) {
  return `/user-templates/${id}/index.html${page > 0 ? `#/${page}` : ''}`
}
