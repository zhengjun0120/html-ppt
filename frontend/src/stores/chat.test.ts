import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'

import { projectTranscript, useChatStore } from './chat'
import type { Transcript } from '@/types/chat'

describe('chat store 事件归一（handleEvent）', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  function store() {
    const s = useChatStore()
    s.reset('deck-0001')
    return s
  }

  it('delta 增量并入同一个 agent 块', () => {
    const s = store()
    s.handleEvent({ type: 'session', content: '5' })
    s.handleEvent({ type: 'delta', content: '你好' })
    s.handleEvent({ type: 'delta', content: '，世界' })
    expect(s.sessionId).toBe(5)
    expect(s.events.filter((e) => e.kind === 'agent')).toHaveLength(1)
    const agent = s.events.find((e) => e.kind === 'agent')
    expect(agent && 'text' in agent && agent.text).toBe('你好，世界')
  })

  it('并行工具按 tool_index 归位，结果按 tool_call_id 回填', () => {
    const s = store()
    s.handleEvent({ type: 'tool_start', tool_index: 0, tool_call_id: 'c1', tool_name: 'write_page' })
    s.handleEvent({ type: 'tool_start', tool_index: 1, tool_call_id: 'c2', tool_name: 'vision_review' })
    s.handleEvent({ type: 'tool_delta', tool_index: 1, content: '{"page":2}' })
    s.handleEvent({ type: 'tool_call', tool_call_id: 'c1', content: '已写入' })
    s.handleEvent({ type: 'tool_error', tool_call_id: 'c2', content: '工具调用失败 err: timeout' })
    const tools = s.events.filter((e): e is Extract<(typeof s.events)[number], { kind: 'tool' }> => e.kind === 'tool')
    expect(tools).toHaveLength(2)
    expect(tools[0]).toMatchObject({ name: 'write_page', state: 'ok', result: '已写入' })
    expect(tools[1]).toMatchObject({ name: 'vision_review', state: 'error', args: '{"page":2}' })
  })

  it('写操作成功置 deckTouched，只读工具不影响', () => {
    const s = store()
    s.handleEvent({ type: 'tool_start', tool_index: 0, tool_call_id: 'r1', tool_name: 'review_slides' })
    s.handleEvent({ type: 'tool_call', tool_call_id: 'r1', content: '报告' })
    expect(s.deckTouched).toBe(false)

    s.handleEvent({ type: 'tool_start', tool_index: 1, tool_call_id: 'w1', tool_name: 'update_slide' })
    s.handleEvent({ type: 'tool_call', tool_call_id: 'w1', content: '已更新' })
    expect(s.deckTouched).toBe(true)
  })

  it('ask_user 进入暂停态，done/error 收尾', () => {
    const s = store()
    s.handleEvent({ type: 'delta', content: '确认一下' })
    s.handleEvent({ type: 'ask_user', tool_call_id: 'c9', content: '{"questions":[{"question":"几页？","options":["3页","5页"]}]}' })
    expect(s.status).toBe('paused')
    expect(s.pendingAsk).toEqual({ toolCallId: 'c9', questions: [{ question: '几页？', options: ['3页', '5页'] }] })
    s.handleEvent({ type: 'done', content: '', total_tokens: 100, cached_tokens: 80 })
    expect(s.status).toBe('idle')
    const done = s.events.find((e) => e.kind === 'done')
    expect(done && 'totalTokens' in done && done.totalTokens).toBe(100)
  })
})

describe('projectTranscript（历史回放投影）', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('标准对话投影：system 跳过、工具结果回填、ask_user 应答展示', () => {
    const t: Transcript = {
      session: { id: 21, deck_id: 'deck-0020', title: 't', pending: false, created_at: '', updated_at: '' },
      messages: [
        { seq: 1, role: 'user', content: '做一份 Go 介绍' },
        {
          seq: 2,
          role: 'assistant',
          content: '',
          tool_calls: [
            { id: 'c1', name: 'write_page', arguments: '{"page":1}' },
            { id: 'c2', name: 'ask_user', arguments: '{"questions":[{"question":"要目录吗？"}]}' },
          ],
        },
        { seq: 3, role: 'tool', tool_call_id: 'c1', content: '已写入第 1 页' },
        { seq: 4, role: 'tool', tool_call_id: 'c2', content: '{"answers":[{"question":"要目录吗？","answer":"要"}]}' },
        { seq: 5, role: 'assistant', content: '目录已加入' },
      ],
    }
    const out = projectTranscript(t)
    expect(out.map((e) => e.kind)).toEqual(['user', 'tool', 'ask', 'agent'])
    const tool = out[1]
    expect(tool.kind === 'tool' && tool.result).toBe('已写入第 1 页')
    const ask = out[2]
    expect(ask.kind === 'ask' && ask.answered?.[0].answer).toBe('要')
    expect(ask.kind === 'ask' && ask.active).toBe(false)
  })

  it('暂停中的会话：最后一个 ask_user 保持可交互', () => {
    const t: Transcript = {
      session: { id: 16, deck_id: '', title: 't', pending: true, created_at: '', updated_at: '' },
      messages: [
        { seq: 1, role: 'user', content: '做 PPT' },
        { seq: 2, role: 'assistant', content: '', tool_calls: [{ id: 'cp', name: 'ask_user', arguments: '{"questions":[{"question":"风格？"}]}' }] },
      ],
    }
    const out = projectTranscript(t)
    const ask = out.find((e) => e.kind === 'ask')
    expect(ask && 'active' in ask && ask.active).toBe(true)
  })
})
