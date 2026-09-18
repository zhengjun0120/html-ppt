import { request } from './client'

/** 模板变体：class 为 style.css 里预定义的 token 覆盖 class（空串 = 默认观感） */
export interface TemplateVariant {
  id: string
  name: string
  class: string
}

export interface TemplateLayoutMeta {
  id: string
  name: string
  use: string
  roles?: string[]
  constraints?: string
}

export interface TemplateMeta {
  id: string
  name: string
  description: string
  tags?: string[]
  scenario?: string[]
  canvas: { w: number; h: number }
  variants: TemplateVariant[]
  layouts: TemplateLayoutMeta[]
  fonts?: string[]
  source?: { derived_from: string; license: string }
}

export const templateApi = {
  /** 模板清单（公开接口）。gallery 进入时调用 */
  list: () => request<TemplateMeta[]>('/api/templates'),
  get: (id: string) => request<TemplateMeta>(`/api/templates/${id}`),
  /**
   * demo 预览 iframe 的 src（公开接口，无需鉴权）。variant 非空时服务端换肤
   * （把变体 class 挂到 body）；page > 0 时带 #/N 深链直接落在某一页。
   */
  previewUrl: (id: string, variant = '', page = 0) => {
    const q = variant ? `?variant=${encodeURIComponent(variant)}` : ''
    const hash = page > 0 ? `#/${page}` : ''
    return `/api/templates/${id}/preview${q}${hash}`
  },
}
