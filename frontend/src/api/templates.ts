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
}
