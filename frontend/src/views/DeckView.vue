<script setup lang="ts">
import { PhArrowLeft, PhPlus } from '@phosphor-icons/vue'
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import Button from '@/components/ui/Button.vue'
import ChatInput from '@/components/chat/ChatInput.vue'
import ChatMessages from '@/components/chat/ChatMessages.vue'
import HistoryDrawer from '@/components/preview/HistoryDrawer.vue'
import DeckPreview from '@/components/preview/DeckPreview.vue'
import DeckEditModal from '@/components/editor/DeckEditModal.vue'
import WizardStepper from '@/components/wizard/WizardStepper.vue'
import { usePreviewAutoRefresh } from '@/lib/previewRefresh'
import { useWorkspaceBoot } from '@/lib/workspaceBoot'
import { useDeckStore } from '@/stores/deck'
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
const editOpen = ref(false)
const deckStore = useDeckStore()
const deckStoreTitle = computed(() => deckStore.titleOf(deckId.value) || '编辑文稿')

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

/** 返回文稿列表：从列表进来走浏览器历史（页码在 ?page、滚动由 router.scrollBehavior 原生落回）；
 *  直链/刷新进来没有可退的历史，就带上离开时记下的页码直达（滚动由列表页 restoreScroll 兜底）。 */
function backToDecks() {
  const back = window.history.state?.back as string | undefined
  if (back && back.startsWith('/decks')) {
    router.back()
    return
  }
  const savedPage = Number(sessionStorage.getItem('decks-return-page') || 0)
  void router.push({ path: '/decks', query: savedPage > 1 ? { page: String(savedPage) } : undefined })
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
      @edit="editOpen = true"
    >
      <template #toolbar-start>
        <button
          class="inline-flex cursor-pointer items-center gap-1 rounded border border-line bg-surface-2 px-2 py-1 text-[11.5px] font-semibold text-ink-2 transition-colors hover:border-accent hover:text-ink"
          @click="backToDecks"
        >
          <PhArrowLeft :size="12" /> 返回文稿
        </button>
      </template>
    </DeckPreview>
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

    <!-- 手动编辑弹窗：保存后刷新预览（deck 文件是 no-store，重挂 iframe 即拿到新版） -->
    <DeckEditModal
      v-if="deckId"
      :open="editOpen"
      kind="deck"
      :id="deckId"
      :title="deckStoreTitle"
      :canvas="{ w: 1920, h: 1080 }"
      @close="editOpen = false"
      @saved="previewKey += 1"
    />
  </div>
</template>
