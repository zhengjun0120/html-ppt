<script setup lang="ts">
import { PhChatCircleDots, PhPlus, PhX } from '@phosphor-icons/vue'
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import Button from '@/components/ui/Button.vue'
import ChatInput from '@/components/chat/ChatInput.vue'
import ChatMessages from '@/components/chat/ChatMessages.vue'
import GeneratingProgress from '@/components/wizard/GeneratingProgress.vue'
import OutlinePanel from '@/components/wizard/OutlinePanel.vue'
import TemplateGallery from '@/components/wizard/TemplateGallery.vue'
import WizardStepper from '@/components/wizard/WizardStepper.vue'
import { usePreviewAutoRefresh } from '@/lib/previewRefresh'
import { useWorkspaceBoot } from '@/lib/workspaceBoot'
import { useWizardStore, STEP_ROUTES, stepOrder, type WizardStep } from '@/stores/wizard'

/**
 * 向导步骤页（plan-v3 C1）：/new、/new/outline、/new/template、/new/generating
 * 四条路由共用本容器。职责：
 *   1. 启动时的会话恢复（useWorkspaceBoot）；
 *   2. 步骤守卫——路由 step 不允许超前于 deck 当前阶段（stage 推进时自动前跳，
 *      迭代阶段自动跳去 /decks/:id）；
 *   3. 按步骤渲染主区域。澄清步的对话就是主区域；其余步对话收进抽屉
 *      （大纲修订通道与生成过程监控）。
 */

const route = useRoute()
const router = useRouter()
const wizard = useWizardStore()
const { chat, sessions, previewKey, loadSessions, boot, switchSession } = useWorkspaceBoot()

const ROUTE_STEP: Record<string, WizardStep> = {
  'wizard-clarify': 'clarify',
  'wizard-outline': 'outline',
  'wizard-template': 'template',
  'wizard-generate': 'generate',
}

const routeStep = computed<WizardStep>(() => ROUTE_STEP[route.name as string] ?? 'clarify')
const deckId = computed(() => chat.deckId)
const chatOpen = ref(false)
const chatLocked = computed(() => routeStep.value === 'wizard-clarify') // 对话已是主区域

const stopPreviewRefresh = usePreviewAutoRefresh(chat, () => {
  previewKey.value += 1
})
onBeforeUnmount(stopPreviewRefresh)

// 启动：恢复会话 + 同步向导状态，然后做一次守卫校正
onMounted(async () => {
  await boot('', (id) => wizard.syncFromChat(id))
  await enforce()
})

// 阶段推进 → 自动前跳到规范路由（迭代 → /decks/:id）
watch(
  () => wizard.step,
  () => void enforce(),
)
// 路由切换（stepper 点击）→ 校验目标步合法
watch(routeStep, () => void enforce())

/** 守卫：目标步超前于 deck 实际阶段时，拉回当前阶段的规范路由 */
async function enforce() {
  const canonical = wizard.step // null = 还没有 v2 deck（冷启动）
  const cur = routeStep.value
  if (canonical === 'iterate' && wizard.deckId) {
    void router.replace(STEP_ROUTES.iterate(wizard.deckId))
    return
  }
  if (!canonical) {
    if (cur !== 'clarify') void router.replace('/new')
    return
  }
  if (stepOrder(cur) > stepOrder(canonical)) {
    void router.replace(STEP_ROUTES[canonical](wizard.deckId))
  }
}

function newConversation() {
  chat.reset(deckId.value)
  wizard.reset()
  void router.replace('/new')
}

function onSessionChange(e: Event) {
  switchSession(deckId.value, (id) => wizard.syncFromChat(id), Number((e.target as HTMLSelectElement).value))
}

// 冷启动时 deckId 由管线事件送回 → 加载会话列表
watch(
  () => chat.deckId,
  (id, old) => {
    if (id && id !== old) void loadSessions(id)
  },
)
// 一轮对话结束：阶段可能被 agent 推进 → 刷新 + 守卫校正
watch(
  () => chat.status,
  (s, old) => {
    if (old === 'streaming' && s === 'idle') {
      void loadSessions(deckId.value)
      void wizard.refresh().then(() => enforce())
    }
  },
)

const showOutline = computed(() => routeStep.value === 'wizard-outline')
const showGallery = computed(() => routeStep.value === 'wizard-template')
const showGenerating = computed(() => routeStep.value === 'wizard-generate')
const inputLockedHint = computed(() =>
  wizard.locksInput ? '请在模板页选择一个模板（这一步确定整套视觉，对话里做不了）' : '',
)
const headTitle = computed(
  () =>
    ({
      'wizard-clarify': '',
      'wizard-outline': '确认大纲：可直接编辑，也可在对话里让我改',
      'wizard-template': '选择模板：这一步确定整套视觉',
      'wizard-generate': '生成中',
    })[route.name as string] ?? '',
)
</script>

<template>
  <div class="flex h-full overflow-hidden">
    <!-- 主区域 -->
    <div class="flex min-h-0 min-w-0 flex-1 flex-col">
      <!-- 步骤头：stepper（可点击回退）+ 标题 + 新对话 -->
      <header class="flex items-center gap-3 border-b border-line bg-surface px-4 py-2">
        <WizardStepper :clickable="true" />
        <span v-if="headTitle" class="hidden truncate text-[12.5px] font-semibold text-ink lg:inline">{{ headTitle }}</span>
        <Button class="ml-auto" size="sm" @click="newConversation">
          <PhPlus :size="12" />
          新对话
        </Button>
      </header>

      <div class="min-h-0 flex-1">
        <OutlinePanel v-if="showOutline" class="h-full" />
        <TemplateGallery v-else-if="showGallery" class="h-full" />
        <GeneratingProgress v-else-if="showGenerating" class="h-full" />
        <!-- 澄清步：对话即主区域 -->
        <div v-else class="mx-auto flex h-full w-full max-w-[760px] flex-col">
          <div class="min-h-0 flex-1 overflow-hidden px-4 pt-4">
            <ChatMessages />
          </div>
          <div class="border-t border-line p-3">
            <ChatInput :locked-hint="inputLockedHint" />
          </div>
        </div>
      </div>
    </div>

    <!-- 对话抽屉：非澄清步可展开（大纲修订通道 / 生成过程监控） -->
    <aside
      v-if="!chatLocked"
      class="min-h-0 w-full shrink-0 flex-col border-l border-line bg-surface md:flex md:w-[380px]"
      :class="chatOpen ? 'fixed inset-0 z-30 flex md:relative' : 'hidden'"
    >
      <div class="flex items-center justify-between border-b border-line px-3 py-1.5">
        <span class="text-[12.5px] font-semibold text-ink-2">对话</span>
        <Button size="sm" variant="ghost" @click="chatOpen = false">
          <PhX :size="12" />
          关闭
        </Button>
      </div>
      <div v-if="deckId && sessions.length" class="border-b border-line px-3 py-1.5">
        <select
          class="w-full cursor-pointer rounded-control border border-line bg-surface-2 px-2 py-1 text-[12px] text-ink-2 outline-none focus-visible:border-accent"
          :value="chat.sessionId ?? 0"
          @change="onSessionChange"
        >
          <option v-if="chat.sessionId == null" :value="0">本次对话</option>
          <option v-for="s in sessions" :key="s.id" :value="s.id">
            {{ s.title || `会话 ${s.id}` }}{{ s.pending ? '（待回答）' : '' }}
          </option>
        </select>
      </div>
      <ChatMessages class="min-h-0 flex-1 overflow-y-auto" />
      <div class="border-t border-line p-3">
        <ChatInput :locked-hint="inputLockedHint" />
      </div>
    </aside>

    <!-- 抽屉关闭时的浮动唤出钮（仅非澄清步） -->
    <button
      v-if="!chatLocked && !chatOpen"
      class="fixed bottom-24 right-4 z-30 inline-flex h-10 w-10 items-center justify-center rounded-full border border-line bg-surface text-ink-2 shadow-sm hover:text-ink md:bottom-6"
      title="打开对话"
      @click="chatOpen = true"
    >
      <PhChatCircleDots :size="16" />
    </button>
  </div>
</template>
