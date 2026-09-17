import { nextTick, reactive } from 'vue'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { usePreviewAutoRefresh } from './previewRefresh'

describe('usePreviewAutoRefresh（预览自动刷新）', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })
  afterEach(() => {
    vi.useRealTimers()
  })

  async function make() {
    const chat = reactive({ deckTouched: false, status: 'idle' })
    let bumps = 0
    const stop = usePreviewAutoRefresh(chat, () => bumps++)
    return {
      chat,
      bump: () => bumps++,
      get bumps() {
        return bumps
      },
      stop,
      async touch() {
        chat.deckTouched = true
        await nextTick()
      },
      async settle() {
        await nextTick()
      },
    }
  }

  it('写操作置位后防抖 800ms 刷新一次', async () => {
    const h = await make()
    h.chat.status = 'streaming'
    await h.touch()
    vi.advanceTimersByTime(799)
    expect(h.bumps).toBe(0)
    vi.advanceTimersByTime(1)
    expect(h.bumps).toBe(1)
    // 持续为 true 不重复刷
    vi.advanceTimersByTime(5000)
    expect(h.bumps).toBe(1)
  })

  it('本轮结束（streaming→idle）且有写操作：立即兜底刷新', async () => {
    const h = await make()
    h.chat.status = 'streaming'
    await h.touch()
    vi.advanceTimersByTime(100) // 防抖挂起中
    h.chat.status = 'idle' // 兜底立即刷
    await h.settle()
    expect(h.bumps).toBe(1)
    vi.advanceTimersByTime(2000) // 挂起的防抖已被取消
    expect(h.bumps).toBe(1)
  })

  it('只读对话（无写操作）不刷新', async () => {
    const h = await make()
    h.chat.status = 'streaming'
    await h.settle()
    h.chat.status = 'idle'
    await h.settle()
    vi.advanceTimersByTime(5000)
    expect(h.bumps).toBe(0)
  })

  it('stop 后不再刷新', async () => {
    const h = await make()
    await h.touch()
    h.stop()
    vi.advanceTimersByTime(2000)
    expect(h.bumps).toBe(0)
  })
})
