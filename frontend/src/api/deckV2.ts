import { authedUrl, request } from './client'
import type { TemplateVariant } from './templates'

/** deck-v2 管线类型与接口（与 service/deck/v2.go 的 JSON 形状逐字段对齐） */

export type DeckStage =
  | 'draft'
  | 'outlining'
  | 'outline_review'
  | 'selecting_template'
  | 'generating'
  | 'iterating'

export interface DeckFile {
  id: string
  format: string
  stage: DeckStage
  title: string
  template_id?: string
  variant?: string
  canvas: { w: number; h: number }
  page_plan?: { no: number; layout: string; reason?: string }[]
}

export interface OutlineMaterial {
  type?: string
  desc: string
}

export interface OutlinePage {
  no: number
  role: string
  title: string
  points?: string[]
  layout_hint?: string
  materials?: OutlineMaterial[]
  notes?: string
}

export interface Outline {
  version: number
  title: string
  meta?: { audience?: string; duration_min?: number; page_count?: number; tone?: string }
  narrative?: { hook?: string; arcs?: string[] }
  pages: OutlinePage[]
}

export interface DeckV2Meta {
  deck: DeckFile
  outline?: Outline
  variants?: TemplateVariant[]
}

export class OutlineConflictError extends Error {
  readonly latestVersion: number
  constructor(latestVersion: number, message: string) {
    super(message)
    this.name = 'OutlineConflictError'
    this.latestVersion = latestVersion
  }
}

export const deckV2Api = {
  meta: (deckId: string) => request<DeckV2Meta>(`/api/decks/${deckId}/meta`),

  /** 面板直改大纲（整份替换 + 乐观锁）。409 = 版本冲突 */
  async putOutline(deckId: string, version: number, outline: Outline): Promise<{ version: number }> {
    const res = await fetchWithAuth<{ version: number }>(`/api/decks/${deckId}/outline`, {
      method: 'PUT',
      body: JSON.stringify({ version, outline }),
    })
    return res
  },

  confirmOutline: (deckId: string) =>
    request<{ stage: DeckStage }>(`/api/decks/${deckId}/outline/confirm`, { method: 'POST' }),

  selectTemplate: (deckId: string, templateId: string, variant: string) =>
    request<{ stage: DeckStage }>(`/api/decks/${deckId}/template`, {
      method: 'POST',
      body: { template_id: templateId, variant },
    }),
}

/** 带 Bearer 的裸 fetch（PUT 需要 409 冲突语义，走 request 会把非 2xx 直接抛掉） */
async function fetchWithAuth<T>(path: string, init: RequestInit = {}): Promise<T> {
  const { currentToken } = await import('./client')
  const token = currentToken()
  const base: string = import.meta.env.VITE_API_BASE ?? ''
  const res = await fetch(base + path, {
    ...init,
    headers: { 'Content-Type': 'application/json', ...(token ? { Authorization: `Bearer ${token}` } : {}) },
  })
  const text = await res.text()
  let body: unknown = {}
  try {
    body = text ? JSON.parse(text) : {}
  } catch { /* 非 JSON */ }
  if (!res.ok) {
    const data = body as { error?: string; latest_version?: number; message?: string; code?: string }
    if (res.status === 409 && data.code === 'outline_conflict') {
      throw new OutlineConflictError(data.latest_version ?? 0, data.message ?? '大纲版本冲突')
    }
    throw new Error(data.error ?? `请求失败（${res.status}）`)
  }
  return body as T
}

// —— 缩略图（plan-v3 C3）：预览栏翻页与文稿列表封面 ——//

/** 单页缩略图地址（<img> 用：?token= 兼容图片标签带不了鉴权头）。
 *  首次访问会触发后端整本渲染（10-20s），之后按内容版本缓存。 */
export function thumbUrl(deckId: string, no: number): string {
  return authedUrl(`/api/decks/${deckId}/thumbs/${no}`)
}

/** 缩略图清单（pages 为 1 基页码升序） */
export async function thumbPages(deckId: string): Promise<number[]> {
  const { request } = await import('./client')
  const r = await request<{ pages: number[]; count: number }>(`/api/decks/${deckId}/thumbs`)
  return r.pages
}

/** 触发导出（pdf/png/html），返回产物文件名；下载用 thumbUrl 风格的 authedUrl 拼 exports 路径 */
export function exportDeck(deckId: string, format: 'pdf' | 'png' | 'html') {
  return request<{ path: string; filename: string; size: number }>(`/api/decks/${deckId}/export`, {
    method: 'POST',
    body: JSON.stringify({ format }),
  })
}
