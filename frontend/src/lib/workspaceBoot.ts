import { ref } from 'vue'

import { ApiError } from '@/api/client'
import { listSessions } from '@/api/chat'
import { useChatStore } from '@/stores/chat'
import { useToast } from '@/stores/toast'

/**
 * 工作台/向导共享的启动逻辑（plan-v3 C1 从 WorkspaceView 抽出）：
 * 会话列表加载、?session=N 深链、最近会话自动恢复、deckId 变化跟随。
 * 调用方在 onMounted 里 await boot()，之后照常渲染。
 */
export function useWorkspaceBoot() {
  const chat = useChatStore()
  const toast = useToast()
  const sessions = ref<{ id: number; title: string; pending: boolean; updated_at: string }[]>([])
  const previewKey = ref(0)

  async function loadSessions(deckId: string) {
    if (!deckId) return
    try {
      sessions.value = await listSessions(deckId)
    } catch {
      /* 列表失败不阻塞工作台 */
    }
  }

  async function boot(routeDeckId: string, wizardSync: (deckId: string) => Promise<void>) {
    chat.reset(routeDeckId)
    const q = Number(new URLSearchParams(window.location.search).get('session'))
    if (Number.isFinite(q) && q > 0) {
      try {
        await chat.loadSession(q)
        await loadSessions(chat.deckId)
        await wizardSync(chat.deckId)
        return
      } catch (e) {
        toast.error(e instanceof ApiError ? e.message : '载入历史会话失败')
      }
    }
    await loadSessions(chat.deckId)
    if (sessions.value.length > 0) {
      try {
        await chat.loadSession(sessions.value[0].id)
      } catch (e) {
        toast.error(e instanceof ApiError ? e.message : '载入上次对话失败')
      }
    }
    await wizardSync(chat.deckId)
  }

  function switchSession(deckId: string, wizardSync: (deckId: string) => Promise<void>, id: number) {
    if (id === (chat.sessionId ?? 0)) return
    void chat
      .loadSession(id)
      .then(() => wizardSync(chat.deckId))
      .catch((e) => toast.error(e instanceof Error ? e.message : '切换会话失败'))
  }

  return { chat, sessions, previewKey, loadSessions, boot, switchSession }
}
