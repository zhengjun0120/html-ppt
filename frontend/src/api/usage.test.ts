import { afterEach, describe, expect, it, vi } from 'vitest'

import { configureClient } from './client'
import { usageApi } from './usage'

// 与 trace.test.ts 同款契约测试：不真发请求，只认 URL/方法/参数形状。
afterEach(() => {
  vi.unstubAllGlobals()
})

describe('usageApi.overview', () => {
  it('GET /api/usage/overview 并透传聚合体', async () => {
    const payload = {
      today: { tokens: 2000, calls: 2, cached_rate: 0.5 },
      month: { tokens: 12345, calls: 9, cached_rate: 0.4 },
      daily: [{ date: '2026-10-08', cached: 800, uncached_in: 700, output: 500, total: 2000, calls: 2 }],
    }
    const fetchMock = vi.fn(async () => new Response(JSON.stringify(payload), { status: 200 }))
    vi.stubGlobal('fetch', fetchMock)
    configureClient({ getToken: () => 'tok-abc' })

    const got = await usageApi.overview()

    expect(got).toEqual(payload)
    const [url, init] = fetchMock.mock.calls[0] as [string, RequestInit]
    expect(url).toBe('/api/usage/overview')
    expect(init.method).toBe('GET')
    expect((init.headers as Record<string, string>)['Authorization']).toBe('Bearer tok-abc')
  })
})
