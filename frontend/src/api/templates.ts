import { authedUrl, request } from './client'

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
  /** 缩略图内容版本（8 位）：卡片 <img> 地址 /api/templates/:id/thumb?v= 的 cache-bust 键 */
  thumb?: string
  /** 清单接口（/api/templates）已瘦身不下发；单模板详情（GET /api/templates/:id）仍有 */
  layouts?: TemplateLayoutMeta[]
  fonts?: string[]
  source?: { derived_from: string; license: string }
  /** demo 页数（服务端统计顶层 section）；0/缺省 = 未知 */
  demo_pages?: number
}

export const templateApi = {
  /** 模板清单（公开接口）。gallery 进入时调用 */
  list: () => request<TemplateMeta[]>('/api/templates'),
  get: (id: string) => request<TemplateMeta>(`/api/templates/${id}`),
  /**
   * demo 预览 iframe 的 src（公开接口，无需鉴权）。variant 非空时服务端换肤
   * （把变体 class 挂到 body）；page > 0 时带 #/N 深链直接落在某一页；
   * slide > 0 时服务端只返回该页（缩略模式）——画廊一屏十几张卡，每张没必要
   * 为整本 demo 付解析+布局的 CPU 账。翻页用的完整预览不要传 slide。
   */
  /** 卡片缩略图地址（<img> 用；?token= 兼容图片标签带不了鉴权头）。
   *  首次访问触发后端 Chrome 渲染（1-3s），之后浏览器 immutable 长缓存；
   *  渲染失败 404 → 调用方回退活 iframe（onerror）。 */
  thumbUrl: (id: string, version: string) =>
    authedUrl(`/api/templates/${id}/thumb?v=${encodeURIComponent(version)}`),

  previewUrl: (id: string, variant = '', page = 0, slide = 0) => {
    const params = new URLSearchParams()
    if (variant) params.set('variant', variant)
    if (slide > 0) params.set('slide', String(slide))
    const q = params.size ? `?${params.toString()}` : ''
    const hash = page > 0 ? `#/${page}` : ''
    return `/api/templates/${id}/preview${q}${hash}`
  },
}
