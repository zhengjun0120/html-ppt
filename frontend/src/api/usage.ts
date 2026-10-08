import { request } from './client'

/** 与后端 internal/service/usage 的 json tag 对齐 */

export interface UsageTokenStat {
  tokens: number
  calls: number
  /** 缓存命中 ÷ 输入（0~1） */
  cached_rate: number
}

export interface UsageDay {
  date: string
  /** 输入里命中缓存的部分 */
  cached: number
  /** 输入里未命中的部分 */
  uncached_in: number
  output: number
  total: number
  calls: number
}

export interface UsageOverview {
  today: UsageTokenStat
  month: UsageTokenStat
  /** 固定近 30 天（含今天），空日补零 */
  daily: UsageDay[]
}

export const usageApi = {
  overview: () => request<UsageOverview>('/api/usage/overview'),
}
