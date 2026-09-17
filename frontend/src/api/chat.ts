import { currentToken } from './client'

/**
 * chat 相关接口。start/answer 返回原始 Response（SSE 流，由 lib/sse 消费）；
 * 错误处理与 client.request 一致：非 2xx 读 {"error"}。
 */
const BASE: string = import.meta.env.VITE_API_BASE ?? ''

async function sseFetch(path: string, body: unknown, signal?: AbortSignal): Promise<Response> {
  const token = currentToken()
  const res = await fetch(BASE + path, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    body: JSON.stringify(body),
    signal,
  })
  if (!res.ok) {
    let msg = `请求失败（${res.status}）`
    try {
      const data = (await res.json()) as { error?: string }
      if (data?.error) msg = data.error
    } catch { /* 非 JSON body */ }
    throw new Error(msg)
  }
  return res
}

export interface StartChatArgs {
  sessionId: number
  deckId: string
  content: string
  webSearch?: boolean
}

export function startChat(a: StartChatArgs, signal?: AbortSignal): Promise<Response> {
  return sseFetch(
    '/api/chat',
    {
      session_id: a.sessionId,
      user_content: a.content,
      deck_id: a.deckId,
      enable_web_search: a.webSearch ?? false,
    },
    signal,
  )
}

export function answerChat(
  sessionId: number,
  answers: { question: string; answer: string }[],
  note: string,
  signal?: AbortSignal,
): Promise<Response> {
  return sseFetch('/api/chat/answer', { session_id: sessionId, answers, note }, signal)
}

export interface PendingResp {
  questions: string
}

export async function fetchPending(sessionId: number): Promise<string> {
  const { request } = await import('./client')
  const r = await request<PendingResp>('/api/chat/pending', { query: { session_id: sessionId } })
  return r.questions
}

/** 生成 run 的 SSE 入口（POST /api/decks/:id/generate）。resume=1 断线续跑 */
export function generateDeck(
  deckId: string,
  sessionId: number,
  resume = false,
  signal?: AbortSignal,
): Promise<Response> {
  return sseFetch(
    `/api/decks/${deckId}/generate?session_id=${sessionId}${resume ? '&resume=1' : ''}`,
    {},
    signal,
  )
}

/** 某个 deck 名下的会话列表（按最近活跃倒序） */
export async function listSessions(deckId: string): Promise<{ id: number; title: string; pending: boolean; updated_at: string }[]> {
  const { request } = await import('./client')
  return request(`/api/decks/${deckId}/chat/sessions`)
}
