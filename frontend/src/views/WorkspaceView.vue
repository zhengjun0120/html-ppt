<script setup lang="ts">
import { PhPlus } from '@phosphor-icons/vue'
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

import Button from '@/components/ui/Button.vue'
import ChatInput from '@/components/chat/ChatInput.vue'
import ChatMessages from '@/components/chat/ChatMessages.vue'
import HistoryDrawer from '@/components/preview/HistoryDrawer.vue'
import PreviewPane from '@/components/preview/PreviewPane.vue'
import { ApiError } from '@/api/client'
import { listSessions } from '@/api/chat'
import { useChatStore } from '@/stores/chat'
import { useDeckStore } from '@/stores/deck'
import { usePreviewAutoRefresh } from '@/lib/previewRefresh'
import { useToast } from '@/stores/toast'

const route = useRoute()
// /decks/new = 无文稿冷启动模式：不带 deck_id 发起对话，agent 会创建新文稿
const isNew = String(route.params.id) === 'new'
const deckId = isNew ? '' : String(route.params.id)
const chat = useChatStore()
const deckStore = useDeckStore()
const toast = useToast()

// 窄屏（<md）：对话/预览 二选一；桌面双栏
const mobileView = ref<'preview' | 'chat'>('preview')
const historyOpen = ref(false)
// agent 写操作后自动刷新预览（组件卸载时停止）
const stopPreviewRefresh = usePreviewAutoRefresh(chat, () => { previewKey.value += 1 })
onBeforeUnmount(stopPreviewRefresh)
const previewKey = ref(0)

// 会话历史（M5）：deck 名下的历史会话，供切换器与自动恢复
const sessions = ref<{ id: number; title: string; pending: boolean; updated_at: string }[]>([])

async function loadSessions() {
  if (isNew) return
  try {
    sessions.value = await listSessions(deckId)
  } catch { /* 列表失败不阻塞工作台 */ }
}

onMounted(async () => {
  chat.reset(deckId)
  void deckStore.ensureList().catch(() => {})
  const q = Number(route.query.session)
  if (Number.isFinite(q) && q > 0) {
    // 深链钩子：?session=N 直接载入指定历史会话
    try {
      await chat.loadSession(q)
      void loadSessions()
      return
    } catch (e) {
      toast.error(e instanceof ApiError ? e.message : '载入历史会话失败')
    }
  }
  // 默认恢复最近一次对话（M5）
  await loadSessions()
  if (sessions.value.length > 0) {
    try {
      await chat.loadSession(sessions.value[0].id)
    } catch (e) {
      toast.error(e instanceof ApiError ? e.message : '载入上次对话失败')
    }
  }
})

watch(
  () => chat.status,
  (s, old) => {
    // 一轮对话结束后刷新会话列表（新会话已入库、旧会话时间更新）
    if (old === 'streaming' && s === 'idle') void loadSessions()
  },
)

function switchSession(id: number) {
  if (id === (chat.sessionId ?? 0)) return
  void chat
    .loadSession(id)
    .catch((e) => toast.error(e instanceof ApiError ? e.message : '切换会话失败'))
}

function newConversation() {
  chat.reset(deckId)
  toast.info('已开启新对话')
}

function onRestored() {
  // 恢复/改动后刷新预览
  previewKey.value += 1
  void deckStore.refresh().catch(() => {})
}
</script>

<template>
  <div class="flex h-full overflow-hidden">
    <!-- 预览为主：占绝大部分宽度；key 变化强制重建 iframe（历史恢复后刷新） -->
    <div class="min-h-0 min-w-0 flex-1" :class="mobileView === 'chat' ? 'hidden md:block' : 'block'">
      <PreviewPane :key="previewKey" :deck-id="deckId" class="h-full" @history="historyOpen = true" />
    </div>

    <!-- 对话侧栏：可收起（窄屏自动转 tab 切换） -->
    <aside
      class="min-h-0 w-full shrink-0 flex-col border-l border-line bg-surface md:flex md:w-[380px]"
      :class="mobileView === 'preview' ? 'hidden md:flex' : 'flex'"
    >
      <div class="flex items-center justify-between border-b border-line px-3 py-1.5">
        <span class="text-[12.5px] font-semibold text-ink-2">对话</span>
        <Button size="sm" @click="newConversation">
          <PhPlus :size="12" />
          新对话
        </Button>
      </div>
      <!-- 会话切换器：自动恢复最近对话（M5） -->
      <div v-if="!isNew && sessions.length" class="border-b border-line px-3 py-1.5">
        <select
          class="w-full cursor-pointer rounded-control border border-line bg-surface-2 px-2 py-1 text-[12px] text-ink-2 outline-none focus-visible:border-accent"
          :value="chat.sessionId ?? 0"
          @change="switchSession(Number(($event.target as HTMLSelectElement).value))"
        >
          <option v-if="chat.sessionId == null" :value="0">本次对话</option>
          <option v-for="s in sessions" :key="s.id" :value="s.id">
            {{ s.title || `会话 ${s.id}` }}{{ s.pending ? '（待回答）' : '' }}
          </option>
          <option v-if="chat.sessionId != null && !sessions.some((s) => s.id === chat.sessionId)" :value="chat.sessionId">
            本次对话（进行中）
          </option>
        </select>
      </div>
      <ChatMessages />
      <div class="border-t border-line p-3 pb-14 md:pb-3">
        <ChatInput />
      </div>
    </aside>

    <!-- 历史版本抽屉 -->
    <HistoryDrawer
      v-if="!isNew"
      :open="historyOpen"
      :deck-id="deckId"
      @close="historyOpen = false"
      @restored="onRestored"
    />

    <!-- 窄屏 tab 栏 -->
    <nav class="fixed bottom-0 left-0 right-0 z-30 flex border-t border-line bg-surface md:hidden">
      <button
        class="flex-1 cursor-pointer py-2.5 text-[13px]"
        :class="mobileView === 'preview' ? 'font-semibold text-accent' : 'text-ink-2'"
        @click="mobileView = 'preview'"
      >
        预览
      </button>
      <button
        class="flex-1 cursor-pointer border-l border-line py-2.5 text-[13px]"
        :class="mobileView === 'chat' ? 'font-semibold text-accent' : 'text-ink-2'"
        @click="mobileView = 'chat'"
      >
        对话
      </button>
    </nav>
  </div>
</template>
