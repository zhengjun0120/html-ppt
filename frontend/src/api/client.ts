/**
 * 统一 fetch 封装。
 * 约定（与后端 response 包一致）：错误一律 `{"error": "中文信息"}` + HTTP 状态码；
 * 成功 200 + 数据本体（无业务码包装）；空列表是 [] 不是 null。
 * 401 统一触发 onUnauthorized（由 auth store 注册：清 token → 跳登录）。
 */

export class ApiError extends Error {
  readonly status: number
  constructor(status: number, message: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

const BASE: string = import.meta.env.VITE_API_BASE ?? ''

let getToken: () => string | null = () => localStorage.getItem('da-token')
let onUnauthorized: () => void = () => {}

export function configureClient(opts: {
  getToken?: () => string | null
  onUnauthorized?: () => void
}) {
  if (opts.getToken) getToken = opts.getToken
  if (opts.onUnauthorized) onUnauthorized = opts.onUnauthorized
}

export function currentToken(): string | null {
  return getToken()
}

function buildUrl(
  path: string,
  query?: Record<string, string | number | boolean | undefined>,
): string {
  let url = BASE + path
  if (query) {
    const qs = new URLSearchParams()
    for (const [k, v] of Object.entries(query)) {
      if (v !== undefined) qs.set(k, String(v))
    }
    const q = qs.toString()
    if (q) url += (url.includes('?') ? '&' : '?') + q
  }
  return url
}

async function readError(res: Response): Promise<string> {
  try {
    const data = (await res.json()) as { error?: string }
    if (data?.error) return data.error
  } catch {
    /* body 不是 JSON（如裸 404），走兜底文案 */
  }
  return `请求失败（${res.status}）`
}

export interface RequestOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  query?: Record<string, string | number | boolean | undefined>
  body?: unknown
  signal?: AbortSignal
}

export async function request<T>(path: string, opts: RequestOptions = {}): Promise<T> {
  const headers: Record<string, string> = {}
  const token = getToken()
  if (token) headers.Authorization = `Bearer ${token}`
  let body: string | undefined
  if (opts.body !== undefined) {
    headers['Content-Type'] = 'application/json'
    body = JSON.stringify(opts.body)
  }

  const res = await fetch(buildUrl(path, opts.query), {
    method: opts.method ?? 'GET',
    headers,
    body,
    signal: opts.signal,
  })

  if (res.status === 401) {
    onUnauthorized()
    throw new ApiError(401, await readError(res))
  }
  if (!res.ok) throw new ApiError(res.status, await readError(res))
  if (res.status === 204) return undefined as T
  return (await res.json()) as T
}

/**
 * 给 iframe / img 这类带不了 Authorization 头的场景拼 URL：
 * 走后端确立的 ?token= 回退（中间件对所有受保护路由通用）。
 * BASE 为空（dev 代理）时产出同源相对 URL，分离部署时产出后端绝对 URL。
 */
export function authedUrl(
  path: string,
  query: Record<string, string | number | undefined> = {},
): string {
  const token = getToken()
  return buildUrl(path, { ...query, token: token ?? undefined })
}
