/**
 * 对话域模型。
 * StreamEvent 与后端 agent/stream_event.go 逐字段对齐；
 * ViewEvent 是渲染层视图模型——流式与回放（历史）共用同一条渲染管线，
 * 后端历史接口产出的 TranscriptMessage 也投影成 ViewEvent。
 */

/** 后端 SSE 事件（agent/stream_event.go） */
export interface StreamEvent {
  type: string
  content: string
  tool_call_id?: string
  tool_name?: string
  tool_index?: number
  prompt_tokens?: number
  completion_tokens?: number
  total_tokens?: number
  cached_tokens?: number
  usage?: unknown
}

export interface AskQuestion {
  question: string
  options?: string[]
}

export interface ViewToolCall {
  id: string
  name: string
  arguments: string
}

export interface ViewUserMsg {
  kind: 'user'
  id: number
  text: string
}
export interface ViewAgentMsg {
  kind: 'agent'
  id: number
  text: string
  streaming: boolean
}
export interface ViewThink {
  kind: 'think'
  id: number
  text: string
  streaming: boolean
}
export interface ViewSub {
  kind: 'sub'
  id: number
  text: string
  streaming: boolean
}
export type ToolState = 'running' | 'ok' | 'error'
export interface ViewTool {
  kind: 'tool'
  id: number
  toolCallId: string
  name: string
  args: string
  argsDone: boolean
  result: string
  state: ToolState
}
export interface ViewAsk {
  kind: 'ask'
  id: number
  toolCallId: string
  questions: AskQuestion[]
  /** 已回答的内容（回放/提交后）；pending 态为 null */
  answered: { question: string; answer: string }[] | null
  note: string | null
  /** true = 正在等用户回答（可交互） */
  active: boolean
}
export interface ViewDone {
  kind: 'done'
  id: number
  totalTokens: number
  cachedTokens: number
}
export interface ViewError {
  kind: 'error'
  id: number
  text: string
}
export interface ViewSys {
  kind: 'sys'
  id: number
  text: string
}

export type ViewEvent =
  | ViewUserMsg
  | ViewAgentMsg
  | ViewThink
  | ViewSub
  | ViewTool
  | ViewAsk
  | ViewDone
  | ViewError
  | ViewSys

export type ChatStatus = 'idle' | 'streaming' | 'paused'

/** GET /api/decks/:id/chat/sessions 的元素（agent/transcript.go SessionSummary） */
export interface SessionSummary {
  id: number
  deck_id: string
  title: string
  pending: boolean
  created_at: string
  updated_at: string
}

/** GET /api/chat/sessions/:id/messages 的消息（TranscriptMessage） */
export interface TranscriptMessage {
  seq: number
  role: 'user' | 'assistant' | 'tool'
  content: string
  tool_calls?: ViewToolCall[]
  tool_call_id?: string
  tool_name?: string
}

export interface Transcript {
  session: SessionSummary
  messages: TranscriptMessage[]
}
