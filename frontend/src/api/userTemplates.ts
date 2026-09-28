import { authedUrl, request } from './client'
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
  community: () => request<CommunityTemplate[]>('/api/community-templates'),
}

/** 读当前 style.css 全文（样式面板；raw 文本，走 request 之外的白名单资产端点） */
export async function fetchTemplateStyle(id: string): Promise<string> {
  const resp = await fetch(authedUrl(`/api/user-templates/${id}/assets/style.css`))
  if (!resp.ok) throw new Error('样式读取失败')
  return resp.text()
}

// —— 定制对话（plan-v3 B2）——//
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
