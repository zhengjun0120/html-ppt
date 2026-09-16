import { request } from './client'

/** 与 trace/store.go、trace/event.go 的 json tag 对齐 */
export interface TraceUsage {
  prompt: number
  completion: number
  total: number
  cached: number
  reasoning?: number
  image_tokens?: number
  calls: number
}

export interface RunUsageSummary {
  duration_ms: number
  turns: number
  tool_calls: number
  usage: Record<string, TraceUsage | undefined>
  total: TraceUsage
}

export interface RunMeta {
  run_id: string
  parent_run_id?: string
  session_id: number
  deck_id?: string
  user_content?: string
  model?: string
  started_at: string
  duration_ms: number
  status: 'ok' | 'paused' | 'error' | 'running'
  turns?: number
  tool_calls?: number
  usage?: RunUsageSummary
  bytes?: number
}

export interface SessionMeta {
  session_id: number
  deck_id?: string
  runs: number
  last_at: string
  latest_run_id: string
  usage: RunUsageSummary
  running_runs: number
}

export interface TraceListResp {
  sessions: SessionMeta[]
  runs: RunMeta[]
}

export interface TraceEventImage {
  name: string
  label?: string
  url?: string
}

export interface TraceEvent {
  seq: number
  ts: string
  kind: 'run_start' | 'llm_request' | 'llm_response' | 'tool_call' | 'tool_result' | 'sub_step' | 'usage' | 'error' | 'run_end'
  run_id?: string
  parent_run_id?: string
  session_id?: number
  deck_id?: string
  user_content?: string
  model?: string
  tools?: string[]
  turn?: number
  tool_name?: string
  tool_call_id?: string
  args?: string
  result?: string
  duration_ms?: number
  error?: string
  finish_reason?: string
  messages?: unknown
  message_count?: number
  bytes?: number
  content?: string
  tool_calls?: { id: string; name: string; arguments: string }[]
  component?: string
  usage?: TraceUsage
  sub?: { name: string; stage: string; text?: string; data?: unknown }
  images?: TraceEventImage[]
  status?: string
  summary?: RunUsageSummary
}

export interface RunReadResp {
  events: TraceEvent[]
  next_offset: number
  running: boolean
}

export const traceApi = {
  list: (query?: { session_id?: number; limit?: number; offset?: number }) =>
    request<TraceListResp>('/api/traces', { query }),
  run: (
    sessionId: number | string,
    runId: string,
    query?: { from_offset?: number; from_seq?: number; include?: string; limit?: number },
  ) => request<RunReadResp>(`/api/traces/${sessionId}/${runId}`, { query }),
  event: (sessionId: number | string, runId: string, seq: number, includeMessages = true) =>
    request<TraceEvent>(`/api/traces/${sessionId}/${runId}/events/${seq}`, {
      query: includeMessages ? { include: 'messages' } : undefined,
    }),
}
