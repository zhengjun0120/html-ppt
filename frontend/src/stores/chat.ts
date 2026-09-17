import { defineStore } from 'pinia'

import { answerChat, startChat } from '@/api/chat'
import { fetchPending } from '@/api/chat'
import { iterateSse } from '@/lib/sse'
import type {
  AskQuestion,
  ChatStatus,
  StreamEvent,
  Transcript,
  TranscriptMessage,
  ViewEvent,
  ViewTool,
} from '@/types/chat'

let idSeq = 0
const nextId = () => ++idSeq

/** 会改动 deck 内容/外观的工具（与后端 func_tool.go 注册名对齐）：命中时预览需要刷新 */
export const DECK_MODIFYING_TOOLS = new Set([
  'write_deck',
  'update_slide',
  'insert_slide',
  'delete_slide',
  'update_theme',
  'update_custom_css',
])

/** 当前流的 abort 控制器（非序列化状态，不放 store.state） */
let abortCtl: AbortController | null = null

export const useChatStore = defineStore('chat', {
  state: () => ({
    sessionId: null as number | null,
    deckId: '',
    status: 'idle' as ChatStatus,
    events: [] as ViewEvent[],
    /** 暂停中的提问（ask_user）；非 null 时输入框锁定、提问卡可交互 */
    pendingAsk: null as { toolCallId: string; questions: AskQuestion[] } | null,
    /** 并行工具调用：tool_index → tool_call_id（流式期间） */
    toolIndexToId: {} as Record<number, string>,
    /** 本次 run 中 agent 是否已成功改动文稿（预览自动刷新的信号） */
    deckTouched: false,
    /** 正在载入历史会话（期间禁止发送，避免恢复流程吞掉新消息） */
    restoring: false,
  }),

  getters: {
    /** 输入框可发消息：不在流式/暂停态，也不在历史恢复中 */
    canSend: (s) => s.status === 'idle' && !s.restoring,
  },

  actions: {
    reset(deckId: string) {
      this.deckId = deckId
      this.sessionId = null
      this.status = 'idle'
      this.events = []
      this.pendingAsk = null
      this.toolIndexToId = {}
      this.deckTouched = false
      abortCtl?.abort()
      abortCtl = null
    },

    /** 流式文本块收尾 */
    finalizeStreaming() {
      for (const e of this.events) {
        if ('streaming' in e) e.streaming = false
      }
    },

    handleEvent(ev: StreamEvent) {
      switch (ev.type) {
        case 'session': {
          const id = Number(ev.content)
          if (Number.isFinite(id) && id > 0) this.sessionId = id
          break
        }
        case 'delta': {
          const last = this.events[this.events.length - 1]
          if (last && last.kind === 'agent' && last.streaming) last.text += ev.content
          else this.events.push({ kind: 'agent', id: nextId(), text: ev.content, streaming: true })
          break
        }
        case 'think': {
          const last = this.events[this.events.length - 1]
          if (last && last.kind === 'think' && last.streaming) last.text += ev.content
          else this.events.push({ kind: 'think', id: nextId(), text: ev.content, streaming: true })
          break
        }
        case 'sub_delta': {
          const last = this.events[this.events.length - 1]
          if (last && last.kind === 'sub' && last.streaming) last.text += ev.content
          else this.events.push({ kind: 'sub', id: nextId(), text: ev.content, streaming: true })
          break
        }
        case 'tool_start': {
          if (ev.tool_call_id) {
            const idx = ev.tool_index ?? 0
            this.toolIndexToId[idx] = ev.tool_call_id
          }
          this.events.push({
            kind: 'tool',
            id: nextId(),
            toolCallId: ev.tool_call_id ?? '',
            name: ev.tool_name ?? '',
            args: '',
            argsDone: false,
            result: '',
            state: 'running',
          })
          break
        }
        case 'tool_delta': {
          // 并行工具调用按 tool_index 归位到各自的卡片
          const toolCallId = this.toolIndexToId[ev.tool_index ?? 0]
          for (let i = this.events.length - 1; i >= 0; i--) {
            const e = this.events[i]
            if (e.kind === 'tool' && e.toolCallId === toolCallId && !e.argsDone) {
              e.args += ev.content
              break
            }
          }
          break
        }
        case 'tool_call': {
          const tool = this.findTool(ev.tool_call_id)
          if (tool) {
            tool.state = 'ok'
            tool.argsDone = true
            tool.result = ev.content
            // 写操作成功 = deck 已被改动，工作台据此刷新预览。
            // 名字优先取工具卡自己的（tool_start 时记录），事件缺 tool_name 也能命中
            if (DECK_MODIFYING_TOOLS.has(tool.name || ev.tool_name || '')) this.deckTouched = true
          }
          break
        }
        case 'tool_error': {
          const tool = this.findTool(ev.tool_call_id)
          if (tool) {
            tool.state = 'error'
            tool.argsDone = true
            tool.result = ev.content
          }
          break
        }
        case 'ask_user': {
          let questions: AskQuestion[] = []
          try {
            questions = (JSON.parse(ev.content) as { questions?: AskQuestion[] }).questions ?? []
          } catch {
            questions = [{ question: ev.content, options: [] }]
          }
          const toolCallId = ev.tool_call_id ?? ''
          this.events.push({
            kind: 'ask',
            id: nextId(),
            toolCallId,
            questions,
            answered: null,
            note: null,
            active: true,
          })
          this.pendingAsk = { toolCallId, questions }
          this.status = 'paused'
          break
        }
        case 'done': {
          this.finalizeStreaming()
          this.events.push({
            kind: 'done',
            id: nextId(),
            totalTokens: ev.total_tokens ?? 0,
            cachedTokens: ev.cached_tokens ?? 0,
          })
          this.status = 'idle'
          break
        }
        case 'error': {
          this.finalizeStreaming()
          this.events.push({ kind: 'error', id: nextId(), text: ev.content })
          this.status = 'idle'
          break
        }
        case 'trace':
        default:
          // 观测细节归观测台；未知类型静默忽略
          break
      }
    },

    findTool(toolCallId: string | undefined): ViewTool | null {
      if (!toolCallId) return null
      for (let i = this.events.length - 1; i >= 0; i--) {
        const e = this.events[i]
        if (e.kind === 'tool' && e.toolCallId === toolCallId) return e
      }
      return null
    },

    /** 发送用户消息（SSE 流式） */
    async send(text: string, webSearch = false) {
      const content = text.trim()
      if (!content || this.status !== 'idle') return
      this.deckTouched = false
      this.events.push({ kind: 'user', id: nextId(), text: content })
      this.status = 'streaming'
      await this.consume(() =>
        startChat({ sessionId: this.sessionId ?? 0, deckId: this.deckId, content, webSearch }),
      )
    },

    /** 回答 ask_user 提问 */
    async answer(answers: { question: string; answer: string }[], note: string) {
      if (this.sessionId == null || this.status !== 'paused') return
      this.deckTouched = false
      const ask = [...this.events].reverse().find((e) => e.kind === 'ask' && e.active)
      if (ask && ask.kind === 'ask') {
        ask.active = false
        ask.answered = answers.length ? answers : null
        ask.note = answers.length ? null : note
      }
      this.pendingAsk = null
      this.status = 'streaming'
      await this.consume(() => answerChat(this.sessionId!, answers, note))
    },

    /** 统一的 SSE 消费循环 */
    async consume(open: () => Promise<Response>) {
      const ctl = new AbortController()
      abortCtl = ctl
      try {
        const res = await open()
        for await (const frame of iterateSse(res)) {
          if (!frame.data) continue
          try {
            this.handleEvent(JSON.parse(frame.data) as StreamEvent)
          } catch {
            // 单帧 JSON 损坏不拖垮整条流
          }
        }
        // 流结束：服务端可能没发 done（连接中断），把流式块收尾
        this.finalizeStreaming()
        if (this.status === 'streaming') this.status = 'idle'
      } catch (e) {
        this.finalizeStreaming()
        if (this.status === 'streaming') this.status = 'idle'
        if (e instanceof DOMException && e.name === 'AbortError') {
          this.events.push({ kind: 'sys', id: nextId(), text: '已停止生成' })
        } else {
          this.events.push({ kind: 'error', id: nextId(), text: e instanceof Error ? e.message : '网络异常' })
        }
      } finally {
        abortCtl = null
      }
    },

    abort() {
      if (this.status === 'streaming' && abortCtl) abortCtl.abort()
    },

    /**
     * 回放：把历史会话投影成 ViewEvent（M5 进工作台时调用；
     * 与流式共用同一套渲染组件，回放即"把历史灌进同一个数组"）。
     */
    async loadSession(sessionId: number) {
      const { request } = await import('@/api/client')
      this.restoring = true
      try {
        const t = await request<Transcript>(`/api/chat/sessions/${sessionId}/messages`)
        abortCtl?.abort()
        abortCtl = null
        this.sessionId = sessionId
        this.deckId = t.session.deck_id
        this.status = t.session.pending ? 'paused' : 'idle'
        this.toolIndexToId = {}
        this.events = projectTranscript(t)
        const lastAsk = [...this.events].reverse().find(
          (e): e is Extract<ViewEvent, { kind: 'ask' }> => e.kind === 'ask',
        )
        this.pendingAsk =
          t.session.pending && lastAsk
            ? { toolCallId: lastAsk.toolCallId, questions: lastAsk.questions }
            : null
      } finally {
        this.restoring = false
      }
    },

    /** 兜底：仅知 session id 时从服务端拉 pending 重建提问卡 */
    async restorePendingFromServer() {
      if (this.sessionId == null) return
      const raw = await fetchPending(this.sessionId)
      if (!raw) {
        this.pendingAsk = null
        return
      }
      try {
        const questions = (JSON.parse(raw) as { questions?: AskQuestion[] }).questions ?? []
        this.pendingAsk = { toolCallId: '', questions }
      } catch { /* 忽略坏数据 */ }
    },
  },
})

/** TranscriptMessage[] → ViewEvent[]（纯函数，约定见 agent/transcript.go 的投影） */
export function projectTranscript(t: Transcript): ViewEvent[] {
  const out: ViewEvent[] = []
  let idSeq = 0
  const next = () => ++idSeq
  // tool_call_id → 视图事件（tool 卡 / ask 卡），tool 结果消息到达时回填
  const open: Record<string, ViewTool | Extract<ViewEvent, { kind: 'ask' }>> = {}
  let lastAskCallId = ''

  for (const m of t.messages as TranscriptMessage[]) {
    if (m.role === 'user') {
      out.push({ kind: 'user', id: next(), text: m.content })
      continue
    }
    if (m.role === 'assistant') {
      if (m.content) out.push({ kind: 'agent', id: next(), text: m.content, streaming: false })
      for (const tc of m.tool_calls ?? []) {
        if (tc.name === 'ask_user') {
          let questions: AskQuestion[] = []
          try {
            questions = (JSON.parse(tc.arguments) as { questions?: AskQuestion[] }).questions ?? []
          } catch { /* 参数损坏时保底空列表 */ }
          const ask: Extract<ViewEvent, { kind: 'ask' }> = {
            kind: 'ask',
            id: next(),
            toolCallId: tc.id,
            questions,
            answered: null,
            note: null,
            active: false,
          }
          out.push(ask)
          open[tc.id] = ask
          lastAskCallId = tc.id
          continue
        }
        const tool: ViewTool = {
          kind: 'tool',
          id: next(),
          toolCallId: tc.id,
          name: tc.name,
          args: tc.arguments,
          argsDone: true,
          result: '',
          state: 'running',
        }
        out.push(tool)
        open[tc.id] = tool
      }
      continue
    }
    // role = tool：回填结果
    if (m.tool_call_id) {
      const target = open[m.tool_call_id]
      if (!target) continue
      if (target.kind === 'tool') {
        target.result = m.content
        // 后端 transcript 无错误标记：按既有文案约定识别失败
        target.state = m.content.startsWith('工具调用失败') ? 'error' : 'ok'
      } else {
        // ask_user 的应答：{"answers":[...]} 或 {"note":"..."}
        try {
          const parsed = JSON.parse(m.content) as {
            answers?: { question: string; answer: string }[]
            note?: string
          }
          target.answered = parsed.answers ?? null
          target.note = parsed.note ?? null
        } catch { /* 保底：不展示应答 */ }
      }
    }
  }

  // 暂停中的会话：最后一个 ask_user 置为可交互（正在等回答）
  if (t.session.pending && lastAskCallId) {
    const ask = open[lastAskCallId]
    if (ask && ask.kind === 'ask') ask.active = true
  }
  return out
}
