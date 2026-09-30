import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { configureClient } from './client'
import { deckV2Api } from './deckV2'

// suggestTemplates 的请求形状契约：路径、方法、refresh 的 query 表达。
// 走 fetch stub（client.test.ts 同款），不发真请求；node 环境没有 localStorage，
// 用 configureClient 注入空 token 源绕开。
describe('deckV2Api.suggestTemplates', () => {
  beforeEach(() => {
    configureClient({ getToken: () => null, onUnauthorized: () => {} })
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('默认无 query、POST、无 body', async () => {
    const fetchMock = vi.fn(async () => new Response(JSON.stringify({ suggestions: [] }), { status: 200 }))
    vi.stubGlobal('fetch', fetchMock)

    await deckV2Api.suggestTemplates('deck-9')
    const [url, init] = fetchMock.mock.calls[0] as [string, RequestInit]
    expect(url).toBe('/api/decks/deck-9/template-suggestions')
    expect(init.method).toBe('POST')
    expect(init.body).toBeUndefined()
  })

  it('refresh=true 走 ?refresh=1', async () => {
    const fetchMock = vi.fn(async () => new Response(JSON.stringify({ suggestions: [] }), { status: 200 }))
    vi.stubGlobal('fetch', fetchMock)

    await deckV2Api.suggestTemplates('deck-9', true)
    const [url] = fetchMock.mock.calls[0] as [string, RequestInit]
    expect(url).toBe('/api/decks/deck-9/template-suggestions?refresh=1')
  })

  it('返回体原样透传 suggestions', async () => {
    const payload = { suggestions: [{ template_id: 'tech-sharing', variant_id: 'blue', reason: '契合' }] }
    vi.stubGlobal('fetch', vi.fn(async () => new Response(JSON.stringify(payload), { status: 200 })))

    const r = await deckV2Api.suggestTemplates('deck-9')
    expect(r.suggestions[0].template_id).toBe('tech-sharing')
  })

  it('非 2xx 抛 ApiError（阶段守卫 409 场景）', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn(async () => new Response(JSON.stringify({ error: '阶段不对' }), { status: 409 })),
    )
    await expect(deckV2Api.suggestTemplates('deck-9')).rejects.toThrow('阶段不对')
  })
})
