<script setup lang="ts">
import { PhPlus } from '@phosphor-icons/vue'
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

import Button from '@/components/ui/Button.vue'
import ChatInput from '@/components/chat/ChatInput.vue'
import ChatMessages from '@/components/chat/ChatMessages.vue'
import HistoryDrawer from '@/components/preview/HistoryDrawer.vue'
import PreviewPane from '@/components/preview/PreviewPane.vue'
import GeneratingProgress from '@/components/wizard/GeneratingProgress.vue'
import OutlinePanel from '@/components/wizard/OutlinePanel.vue'
import TemplateGallery from '@/components/wizard/TemplateGallery.vue'
import WizardStepper from '@/components/wizard/WizardStepper.vue'
import { ApiError } from '@/api/client'
import { listSessions } from '@/api/chat'
import { useChatStore } from '@/stores/chat'
import { useDeckStore } from '@/stores/deck'
import { useWizardStore } from '@/stores/wizard'
import { usePreviewAutoRefresh } from '@/lib/previewRefresh'
import { useToast } from '@/stores/toast'

const route = useRoute()
// /decks/new = 无文稿冷启动模式：不带 deck_id 发起对话，agent 会创建新文稿
const isNew = String(route.params.id) === 'new'
const routeDeckId = isNew ? '' : String(route.params.id)
const chat = useChatStore()
const deckStore = useDeckStore()
const wizard = useWizardStore()
const toast = useToast()

// 工作台关注的 deck id：冷启动时由管线事件（outline_updated/stage）带回
const deckId = computed(() => chat.deckId || routeDeckId)

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
  if (!deckId.value) return
  try {
    sessions.value = await listSessions(deckId.value)
  } catch { /* 列表失败不阻塞工作台 */ }
}

onMounted(async () => {
  chat.reset(routeDeckId)
  void deckStore.ensureList().catch(() => {})
  const q = Number(route.query.session)
  if (Number.isFinite(q) && q > 0) {
    // 深链钩子：?session=N 直接载入指定历史会话
    try {
      await chat.loadSession(q)
      void loadSessions()
      await wizard.syncFromChat(chat.deckId)
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
  // v2 向导状态随会话恢复（载入的会话可能正停在任一阶段）
  await wizard.syncFromChat(chat.deckId)
})

// 冷启动：outline_updated / stage 事件把 deckId 送回 chat store，向导跟着拉元数据
watch(
  () => chat.deckId,
  (id, old) => {
    if (id && id !== old && id !== routeDeckId) {
      void wizard.syncFromChat(id).catch(() => {})
      void loadSessions()
    }
  },
)

// 一轮对话结束后：向导阶段可能已被 agent 推进（如 outlining → outline_review）
watch(
  () => chat.status,
  (s, old) => {
    if (old === 'streaming' && s === 'idle') {
      void loadSessions()
      void wizard.refresh()
    }
  },
)

// 生成结束后回到本页（刷新会话列表即可，阶段由 wizard.refresh 拿到）
function switchSession(id: number) {
  if (id === (chat.sessionId ?? 0)) return
  void chat
    .loadSession(id)
    .then(() => wizard.syncFromChat(chat.deckId))
    .catch((e) => toast.error(e instanceof ApiError ? e.message : '切换会话失败'))
}

function newConversation() {
  chat.reset(deckId.value)
  wizard.reset()
  toast.info('已开启新对话')
}

function onRestored() {
  // 恢复/改动后刷新预览
  previewKey.value += 1
  void deckStore.refresh().catch(() => {})
}

/** 主区域形态：v2 向导按阶段切换；v1/无向导保持预览 */
const mainMode = computed<'preview' | 'wait' | 'outline' | 'gallery' | 'generating'>(() => {
  if (wizard.showOutline) return 'outline'
  if (wizard.showGallery) return 'gallery'
  if (wizard.showGenerating) return 'generating'
  if (wizard.isV2 && (wizard.stage === 'outlining' || wizard.stage === '')) return 'wait'
  return 'preview'
})
const inputLockedHint = computed(() =>
  wizard.locksInput ? '请先在主区域选择模板（这一步确定整套视觉，对话里做不了）' : '',
)
const chatAreaClass = computed(() =>
  mainMode.value === 'gallery' ? 'hidden md:flex' : 'flex',
)
</script>

<template>
  <div class="flex h-full overflow-hidden">
    <!-- 主区域：预览为主；v2 向导阶段切换为 大纲面板 / 模板画廊 / 生成进度 -->
    <div class="min-h-0 min-w-0 flex-1" :class="mobileView === 'chat' ? 'hidden md:block' : 'block'">
      <template v-if="mainMode === 'preview'">
        <PreviewPane :key="previewKey + ':' + deckId" :deck-id="deckId" class="h-full" @history="historyOpen = true" />
      </template>
      <OutlinePanel v-else-if="mainMode === 'outline'" class="h-full" />
      <TemplateGallery v-else-if="mainMode === 'gallery'" class="h-full" />
      <GeneratingProgress v-else-if="mainMode === 'generating'" class="h-full" />
      <div v-else class="flex h-full flex-col items-center justify-center gap-3 p-6 text-center">
        <WizardStepper />
        <p class="max-w-[420px] text-[13px] text-ink-2">
          在右侧对话里描述你想要的演示文稿：主题、受众、时长。
          agent 会先和你对齐大纲，再由你挑选模板，最后逐页生成。
        </p>
      </div>
    </div>

    <!-- 对话侧栏：可收起（窄屏自动转 tab 切换）；画廊全屏时隐藏 -->
    <aside
      class="min-h-0 w-full shrink-0 flex-col border-l border-line bg-surface md:flex md:w-[380px]"
      :class="chatAreaClass"
    >
      <div class="flex items-center justify-between border-b border-line px-3 py-1.5">
        <WizardStepper v-if="wizard.isV2" />
        <span v-else class="text-[12.5px] font-semibold text-ink-2">对话</span>
        <Button size="sm" @click="newConversation">
          <PhPlus :size="12" />
          新对话
        </Button>
      </div>
      <!-- 会话切换器：自动恢复最近对话（M5） -->
      <div v-if="deckId && sessions.length" class="border-b border-line px-3 py-1.5">
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
        <ChatInput :locked-hint="inputLockedHint" />
      </div>
    </aside>

    <!-- 历史版本抽屉 -->
    <HistoryDrawer
      v-if="deckId"
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
