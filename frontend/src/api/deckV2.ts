import { request } from './client'

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
    const res = await fetchWithAuth(`/api/decks/${deckId}/outline`, {
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
