<script setup lang="ts">
import { PhPlus } from '@phosphor-icons/vue'
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

import Button from '@/components/ui/Button.vue'
import ChatInput from '@/components/chat/ChatInput.vue'
import ChatMessages from '@/components/chat/ChatMessages.vue'
import HistoryDrawer from '@/components/preview/HistoryDrawer.vue'
import PreviewPane from '@/components/preview/PreviewPane.vue'
import { ApiError } from '@/api/client'
import { useChatStore } from '@/stores/chat'
import { useDeckStore } from '@/stores/deck'
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
const previewKey = ref(0)

onMounted(async () => {
  chat.reset(deckId)
  void deckStore.ensureList().catch(() => {})
  // 深链钩子：?session=N 直接载入历史会话（也是 pending 恢复的测试入口）
  const q = Number(route.query.session)
  if (Number.isFinite(q) && q > 0) {
    try {
      await chat.loadSession(q)
    } catch (e) {
      toast.error(e instanceof ApiError ? e.message : '载入历史会话失败')
    }
  }
})

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
