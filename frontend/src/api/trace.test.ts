import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { configureClient } from './client'
import { traceApi } from './trace'

// traceApi.list 的请求形状契约：分页/类型/状态/搜索全部走 query 表达
// （undefined 项不出现），响应 {runs,total} 原样透传。
// fetch stub（deckV2.test.ts 同款），不发真请求。
describe('traceApi.list', () => {
  beforeEach(() => {
    configureClient({ getToken: () => null, onUnauthorized: () => {} })
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('分页与筛选参数进 query，响应原样透传', async () => {
    const payload = { runs: [], total: 42 }
    const fetchMock = vi.fn(async () => new Response(JSON.stringify(payload), { status: 200 }))
    vi.stubGlobal('fetch', fetchMock)

    await expect(traceApi.list({ limit: 25, offset: 50, kind: 'tplsugg', status: 'error', q: '讲稿' })).resolves.toEqual(payload)
    const [url] = fetchMock.mock.calls[0] as [string, RequestInit]
    expect(url).toBe('/api/traces?limit=25&offset=50&kind=tplsugg&status=error&q=%E8%AE%B2%E7%A8%BF')
  })

  it('undefined 项不产生空 query', async () => {
    const fetchMock = vi.fn(async () => new Response(JSON.stringify({ runs: [], total: 0 }), { status: 200 }))
    vi.stubGlobal('fetch', fetchMock)

    await traceApi.list({ limit: 25, kind: undefined, q: undefined })
    const [url] = fetchMock.mock.calls[0] as [string, RequestInit]
    expect(url).toBe('/api/traces?limit=25')
  })
})
