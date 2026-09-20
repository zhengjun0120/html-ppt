<script setup lang="ts">
import { PhChatCircleDots, PhPlus, PhX } from '@phosphor-icons/vue'
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import Button from '@/components/ui/Button.vue'
import ChatInput from '@/components/chat/ChatInput.vue'
import ChatMessages from '@/components/chat/ChatMessages.vue'
import { listRecentSessions, type RecentSession } from '@/api/chat'
import { useToast } from '@/stores/toast'
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
const toast = useToast()
const { chat, sessions, previewKey, loadSessions, boot, switchSession } = useWorkspaceBoot()

const ROUTE_STEP: Record<string, WizardStep> = {
  'wizard-clarify': 'clarify',
  'wizard-outline': 'outline',
  'wizard-template': 'template',
  'wizard-generate': 'generate',
}

const routeStep = computed<WizardStep>(() => ROUTE_STEP[route.name as string] ?? 'clarify')
const deckId = computed(() => chat.deckId)
// 抽屉默认态跟屏幕宽：桌面常驻右栏，移动端收起（浮动唤出钮）
const chatOpen = ref(window.matchMedia('(min-width: 768px)').matches)
const chatLocked = computed(() => routeStep.value === 'clarify') // 对话已是主区域

const stopPreviewRefresh = usePreviewAutoRefresh(chat, () => {
  previewKey.value += 1
})
onBeforeUnmount(stopPreviewRefresh)

// 启动：恢复会话 + 同步向导状态，然后做一次守卫校正
onMounted(async () => {
  await boot('', (id) => wizard.syncFromChat(id))
  await enforce()
  void probeResumable()
})

// —— 续接横幅：澄清对话在 write_outline 之前没有 deck 行，文稿列表里
// 看不见，离开 /new 后唯一入口。用户确认后才恢复，不自动抢占新对话。——

const resumeBanner = ref<RecentSession | null>(null)

/** 未走完向导的会话：还没落 deck（澄清中），或 deck 阶段还没到迭代 */
const WIZARD_STAGES = new Set(['draft', 'outlining', 'outline_review', 'selecting_template', 'generating'])
function resumable(s: RecentSession): boolean {
  if (!s.deck_id) return true
  return s.deck_format === 'v2' && WIZARD_STAGES.has(s.deck_stage)
}

async function probeResumable() {
  if (route.query.session || chat.events.length > 0) return // 深链已恢复 / 已在聊
  try {
    const recent = await listRecentSessions(10)
    const hit = recent.find(resumable)
    const dismissed = sessionStorage.getItem('wizard.resumeDismissed')
    if (hit && dismissed !== String(hit.id)) resumeBanner.value = hit
  } catch { /* 拉不到就不出横幅 */ }
}

async function resumeLast() {
  const s = resumeBanner.value
  if (!s || chat.status !== 'idle') return
  resumeBanner.value = null
  sessionStorage.removeItem('wizard.resumeDismissed')
  try {
    await chat.loadSession(s.id)
    // 有 deck 的会被阶段 watcher 带到对应步骤页；无 deck 的留在本页续聊
    await wizard.syncFromChat(chat.deckId)
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '恢复会话失败')
  }
}

function dismissResume() {
  if (resumeBanner.value) sessionStorage.setItem('wizard.resumeDismissed', String(resumeBanner.value.id))
  resumeBanner.value = null
}

// 阶段推进 → 自动前跳到规范路由（迭代 → /decks/:id）。
// 面板组件不做路由，这里是"大纲产出/确认大纲/选模板后换页"的唯一驱动。
watch(
  () => wizard.step,
  (step) => {
    if (!step) return
    if (step === 'iterate' && wizard.deckId) {
      void router.replace(STEP_ROUTES.iterate(wizard.deckId))
      return
    }
    if (stepOrder(step) > stepOrder(routeStep.value)) {
      void router.replace(STEP_ROUTES[step](wizard.deckId))
    }
  },
)
// 路由切换（stepper 点击回看）→ 只拦"超前"（跳到尚未到达的步骤），允许回看旧步骤
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

const showOutline = computed(() => routeStep.value === 'outline')
const showGallery = computed(() => routeStep.value === 'template')
const showGenerating = computed(() => routeStep.value === 'generate')
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
        <!-- 澄清步：对话即主区域。外层必须是 flex 列：ChatMessages 根节点的
             flex-1/overflow-y-auto 依赖 flex 父级定高，否则消息区按内容长高、
             被裁剪且无法滚动 -->
        <div v-else class="mx-auto flex h-full w-full max-w-[760px] flex-col">
          <!-- 续接横幅：离开 /new 后未走完向导的对话从这里找回 -->
          <div v-if="resumeBanner" class="mx-4 mt-3 flex items-center gap-3 rounded-control border border-accent-border bg-accent-soft px-3 py-2">
            <span class="min-w-0 flex-1 truncate text-[12.5px] text-ink-2">
              上次有进行中的对话：《{{ resumeBanner.title || '未命名对话' }}》{{ resumeBanner.pending ? '，还有一个提问等你回答' : '' }}
            </span>
            <Button size="sm" variant="primary" @click="resumeLast">继续对话</Button>
            <button class="shrink-0 text-ink-3 transition-colors hover:text-ink-2" title="这次不续" @click="dismissResume">
              <PhX :size="14" />
            </button>
          </div>
          <div class="flex min-h-0 flex-1 flex-col overflow-hidden px-4 pt-4">
            <ChatMessages class="min-h-0 flex-1" />
          </div>
          <div class="border-t border-line p-3">
            <ChatInput :locked-hint="inputLockedHint" />
          </div>
        </div>
      </div>
    </div>

    <!-- 对话抽屉：非澄清步展开（大纲修订通道 / 生成过程监控）。移动端全屏浮层，桌面右侧定宽栏 -->
    <aside
      v-if="!chatLocked && chatOpen"
      class="fixed inset-0 z-30 flex min-h-0 w-full flex-col bg-surface md:relative md:inset-auto md:z-auto md:w-[380px] md:shrink-0 md:border-l md:border-line"
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
