import { request } from './client'

/** 与后端 internal/handler/usermodel.go、internal/service/usermodel 的 json tag 对齐 */

export interface UserModelItem {
  id: number
  /** 显示名（默认 = model_id） */
  name: string
  /** 发给供应商的 model 参数 */
  model_id: string
  base_url: string
  /** key 已配置（密文永不回传） */
  has_key: boolean
  created_at: string
}

export interface UserModelListResp {
  models: UserModelItem[]
  /** 当前使用的自选模型 id；0 = 平台模型 */
  active_id: number
  platform_model: string
}

export interface UserModelInput {
  name?: string
  model_id: string
  base_url: string
  /** 新增必填语义上可空（本地网关）；编辑传空 = 保留原 key */
  api_key?: string
}

export const modelApi = {
  list: () => request<UserModelListResp>('/api/models'),
  create: (input: UserModelInput) =>
    request<UserModelItem>('/api/models', { method: 'POST', body: input }),
  update: (id: number, input: UserModelInput) =>
    request<Record<string, never>>(`/api/models/${id}`, { method: 'PUT', body: input }),
  remove: (id: number) =>
    request<Record<string, never>>(`/api/models/${id}`, { method: 'DELETE' }),
  /** model_id 传 0 = 切回平台模型 */
  setActive: (modelId: number) =>
    request<Record<string, never>>('/api/models/active', { method: 'POST', body: { model_id: modelId } }),
  /** id 给了按已存配置测（api_key 可选覆盖）；否则按传入值测 */
  test: (input: { id?: number; base_url?: string; model_id?: string; api_key?: string }) =>
    request<{ ok: boolean }>('/api/models/test', { method: 'POST', body: input }),
}
