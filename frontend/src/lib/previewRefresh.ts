import { watch, type WatchStopHandle } from 'vue'

/**
 * 预览自动刷新：agent 成功执行写操作（chat.deckTouched 置位）后，
 * 防抖刷新预览；本轮结束（streaming → idle）时若有写操作，立即兜底刷新一次。
 * 返回清理函数（组件卸载时调用）。
 */
export function usePreviewAutoRefresh(
  chat: { deckTouched: boolean; status: string },
  bump: () => void,
): () => void {
  let timer: ReturnType<typeof setTimeout> | undefined

  const stopTouched: WatchStopHandle = watch(
    () => chat.deckTouched,
    (touched) => {
      if (!touched || timer) return
      timer = setTimeout(() => {
        timer = undefined
        bump()
      }, 800)
    },
  )

  const stopStatus: WatchStopHandle = watch(
    () => chat.status,
    (s, old) => {
      if (old === 'streaming' && s === 'idle' && chat.deckTouched) {
        if (timer) {
          clearTimeout(timer)
          timer = undefined
        }
        bump()
      }
    },
  )

  return () => {
    stopTouched()
    stopStatus()
    if (timer) clearTimeout(timer)
  }
}
