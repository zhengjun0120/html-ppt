<script setup lang="ts">
import { PhPlus } from '@phosphor-icons/vue'
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import Button from '@/components/ui/Button.vue'
import ChatInput from '@/components/chat/ChatInput.vue'
import ChatMessages from '@/components/chat/ChatMessages.vue'
import HistoryDrawer from '@/components/preview/HistoryDrawer.vue'
import DeckPreview from '@/components/preview/DeckPreview.vue'
import WizardStepper from '@/components/wizard/WizardStepper.vue'
import { usePreviewAutoRefresh } from '@/lib/previewRefresh'
import { useWorkspaceBoot } from '@/lib/workspaceBoot'
import { useToast } from '@/stores/toast'
import { useWizardStore, STEP_ROUTES } from '@/stores/wizard'

/**
 * 第 5 步 · 成品预览 + 迭代（/decks/:id，plan-v3 C1）。
 * 双栏：左预览（缩略图/翻页/导出）+ 右对话（迭代修改）。
 * v2 deck 还没走到迭代阶段时自动跳回对应向导步骤；?session=N 深链保留。
 */

const route = useRoute()
const router = useRouter()
const wizard = useWizardStore()
const toast = useToast()
const isNew = String(route.params.id) === 'new'
const routeDeckId = isNew ? '' : String(route.params.id)
const { chat, sessions, previewKey, loadSessions, boot, switchSession } = useWorkspaceBoot()
const historyOpen = ref(false)

const deckId = computed(() => chat.deckId || routeDeckId)
const stopPreviewRefresh = usePreviewAutoRefresh(chat, () => {
  previewKey.value += 1
})
onBeforeUnmount(stopPreviewRefresh)

function onRestored() {
  previewKey.value += 1
}

async function wizardSync(id: string) {
  await wizard.syncFromChat(id)
  const step = wizard.step
  if (step && step !== 'iterate' && id) {
    void router.replace(STEP_ROUTES[step](id))
  }
}

onMounted(async () => {
  await boot(routeDeckId, (id) => wizardSync(id))
  // v2 deck 尚未到迭代阶段：回向导对应步骤
  const step = wizard.step
  if (step && step !== 'iterate' && deckId.value) {
    void router.replace(STEP_ROUTES[step](deckId.value))
  }
})

// 冷启动？本页不该出现无 deck 的冷启动（/new 负责）；deckId 空时仅展示提示
watch(
  () => chat.deckId,
  (id, old) => {
    if (id && id !== old) void loadSessions(id)
  },
)
watch(
  () => chat.status,
  (s, old) => {
    if (old === 'streaming' && s === 'idle') {
      void loadSessions(deckId.value)
      void wizard.refresh()
    }
  },
)

function newConversation() {
  chat.reset(deckId.value)
  toast.info('已开启新对话')
}
</script>

<template>
  <div class="flex h-full overflow-hidden">
    <DeckPreview
      v-if="deckId"
      :key="previewKey + ':' + deckId"
      :deck-id="deckId"
      class="h-full"
      @history="historyOpen = true"
    />
    <div v-else class="flex min-w-0 flex-1 items-center justify-center p-6 text-[13px] text-ink-3">
      文稿不存在或尚未创建。
    </div>

    <!-- 对话侧栏：迭代修改通道 -->
    <aside class="min-h-0 w-full shrink-0 flex-col border-l border-line bg-surface md:flex md:w-[380px]">
      <div class="flex items-center justify-between border-b border-line px-3 py-1.5">
        <WizardStepper v-if="wizard.isV2" />
        <span v-else class="text-[12.5px] font-semibold text-ink-2">对话</span>
        <Button size="sm" @click="newConversation">
          <PhPlus :size="12" />
          新对话
        </Button>
      </div>
      <div v-if="deckId && sessions.length" class="border-b border-line px-3 py-1.5">
        <select
          class="w-full cursor-pointer rounded-control border border-line bg-surface-2 px-2 py-1 text-[12px] text-ink-2 outline-none focus-visible:border-accent"
          :value="chat.sessionId ?? 0"
          @change="switchSession(deckId, (id) => wizardSync(id), Number(($event.target as HTMLSelectElement).value))"
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

    <HistoryDrawer
      v-if="deckId"
      :open="historyOpen"
      :deck-id="deckId"
      @close="historyOpen = false"
      @restored="onRestored"
    />
  </div>
</template>
