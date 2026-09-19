import { afterEach, describe, expect, it, vi } from 'vitest'

import { ApiError, configureClient, request } from './client'

describe('api client', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('401 时触发 onUnauthorized 并抛 ApiError', async () => {
    const unauthorized = vi.fn()
    configureClient({ getToken: () => 'token-x', onUnauthorized: unauthorized })
    vi.stubGlobal(
      'fetch',
      vi.fn(async () => new Response(JSON.stringify({ error: '登录已过期，请重新登录' }), { status: 401 })),
    )

    await expect(request('/api/auth/me')).rejects.toMatchObject({
      status: 401,
      message: '登录已过期，请重新登录',
    })
    expect(unauthorized).toHaveBeenCalledOnce()
  })

  it('其他非 2xx 抛出后端 error 文案', async () => {
    const unauthorized = vi.fn()
    configureClient({ getToken: () => null, onUnauthorized: unauthorized })
    vi.stubGlobal(
      'fetch',
      vi.fn(async () => new Response(JSON.stringify({ error: '会话不存在' }), { status: 404 })),
    )

    await expect(request('/api/chat/sessions/1/messages')).rejects.toMatchObject({
      status: 404,
      message: '会话不存在',
    })
    expect(unauthorized).not.toHaveBeenCalled()
  })

  it('成功返回数据本体，并带 Authorization 头', async () => {
    configureClient({ getToken: () => 'token-y' })
    const fetchMock = vi.fn(async (_input: string | URL | Request, _init?: RequestInit) =>
      new Response(JSON.stringify([{ id: 'deck-0001', title: 'T' }]), { status: 200 }),
    )
    vi.stubGlobal('fetch', fetchMock)

    await expect(request('/api/decks')).resolves.toEqual([{ id: 'deck-0001', title: 'T' }])
    const init = fetchMock.mock.calls[0][1] as RequestInit
    expect((init.headers as Record<string, string>).Authorization).toBe('Bearer token-y')
  })

  it('authedUrl 拼接 ?token= 回退参数', async () => {
    configureClient({ getToken: () => 'tk' })
    const { authedUrl } = await import('./client')
    expect(authedUrl('/api/decks/deck-0001/file')).toContain('token=tk')
    expect(authedUrl('/api/decks/deck-0001/file', { foo: 1 })).toContain('foo=1')
  })
})
