import { request } from './client'

export interface MeResp {
  email: string
  has_api_key: boolean
}

export const authApi = {
  login: (email: string, password: string) =>
    request<{ token: string }>('/api/auth/login', { method: 'POST', body: { email, password } }),

  requestCode: (email: string) =>
    request<{ sent: boolean }>('/api/auth/code', { method: 'POST', body: { email } }),

  register: (email: string, code: string, password: string) =>
    request<{ token: string }>('/api/auth/register', {
      method: 'POST',
      body: { email, code, password },
    }),

  /** 每用户每天一次；429 = 今日已刷新（继续用旧 token，不算错误） */
  refresh: () => request<{ token: string }>('/api/auth/refresh', { method: 'POST' }),

  me: () => request<MeResp>('/api/auth/me'),

  setApiKey: (apiKey: string) =>
    request<Record<string, never>>('/api/auth/apikey', { method: 'POST', body: { api_key: apiKey } }),
}
