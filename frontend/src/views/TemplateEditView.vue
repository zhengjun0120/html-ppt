<script setup lang="ts">
import { PhArrowClockwise, PhCheck, PhCircleNotch, PhClockCounterClockwise, PhGlobe, PhImage, PhPaperPlaneTilt, PhPulse, PhX } from '@phosphor-icons/vue'
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import TemplateHistoryDrawer from '@/components/usertpl/TemplateHistoryDrawer.vue'
import { customizeChatStream, fetchTemplateStyle, userTemplateApi, userTemplatePreviewUrl, type LayoutMetaPatch, type PublishReport, type StructureContract, type StructureLayout, type UserTemplateRow } from '@/api/userTemplates'
import Button from '@/components/ui/Button.vue'
import { MAX_CHAT_IMAGES, fileToChatImage } from '@/lib/image'
import { useToast } from '@/stores/toast'

/**
 * 定制工作台（plan-v3 B3 + user-template-history-plan §3.6）：
 * 左 = demo 实时预览 + style.css 手动编辑面板 + 质量体检报告；
 * 右 = 与 agent 的定制对话（SSE 流式：工具气泡/生成进度实时可见）；
 * 顶栏 = 名称编辑 + 历史 + 体检/发布（2026-09-28 门禁降级：发布秒级，量测独立为体检）。
 */

const route = useRoute()
const router = useRouter()
const toast = useToast()

const utId = String(route.params.id)
const row = ref<UserTemplateRow | null>(null)
const name = ref('')
const messages = ref<{ role: 'user' | 'agent'; text: string; think?: string; images?: string[] }[]>([])
const input = ref('')
const sending = ref(false)
const publishing = ref(false)
const previewKey = ref(0)
const demoPage = ref(1)

// —— 用户附图（参考图/配色灵感；压缩后随消息发，服务端校验 ≤3 张）——//
const pendingImages = ref<string[]>([])
const attaching = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

async function onPickImages(e: Event) {
  const files = (e.target as HTMLInputElement).files
  if (!files?.length) return
  attaching.value = true
  try {
    for (const f of Array.from(files)) {
      if (pendingImages.value.length >= MAX_CHAT_IMAGES) {
        toast.error(`一次最多附 ${MAX_CHAT_IMAGES} 张图`)
        break
      }
      try {
        pendingImages.value.push(await fileToChatImage(f))
      } catch (err) {
        toast.error(err instanceof Error ? err.message : '图片处理失败')
      }
    }
  } finally {
    attaching.value = false
    e.target.value = '' // 允许重选同一个文件
  }
}

function removePendingImage(i: number) {
  pendingImages.value.splice(i, 1)
}

// —— 实时预览（runtime 预览模式协议）——//
// src 带 ?preview=1 让 demo 内的 runtime.js 进单页锁定模式；翻页向 iframe
// postMessage {type:'preview-goto', idx}（0 基），iframe 只加载一次、内部按
// base.css 的 .5s 淡入+位移过渡换页——与文稿编辑器/模板预览弹窗同一套协议。
// 此前"页码进 src + :key 绑 src"每次翻页整个重建 iframe：白屏一闪、无过渡、
// 且页码无上界（2026-10-06 用户反馈）。
// v= 是 cache-bust：刷新按钮/对话改动后递增，:key 只绑它——只有这两类事件
// 才整帧重建，纯翻页不动 iframe。
const previewReady = ref(false)
const previewFrame = ref<HTMLIFrameElement | null>(null)
const previewSrc = computed(() =>
  row.value ? userTemplatePreviewUrl(utId) + `&preview=1&v=${previewKey.value}` : '',
)

function gotoPreviewPage(i: number) {
  previewFrame.value?.contentWindow?.postMessage({ type: 'preview-goto', idx: i }, '*')
}
function onPreviewMessage(e: MessageEvent) {
  if (e.source !== previewFrame.value?.contentWindow) return
  if ((e.data as { type?: string } | null)?.type === 'preview-ready') {
    previewReady.value = true
    gotoPreviewPage(demoPage.value - 1) // 就绪前点的翻页在这里补发
  }
}
function prevPage() {
  if (demoPage.value <= 1) return
  demoPage.value -= 1
  if (previewReady.value) gotoPreviewPage(demoPage.value - 1)
}
function nextPage() {
  const n = totalPages.value
  if (n && demoPage.value >= n) return
  demoPage.value += 1
  if (previewReady.value) gotoPreviewPage(demoPage.value - 1)
}
// 整帧重建（手动刷新 / 对话改动后）必定回到第 1 页：新加载的预览从第一页开始
function refreshPreview() {
  previewKey.value += 1
  demoPage.value = 1
  previewReady.value = false
}

// —— 质量体检（渲染量测：溢出/填充率/最小字号；只报告不拦发布）——//
const checking = ref(false)
const checkupReport = ref<PublishReport | null>(null)

// —— 对话进行中的实时活动（思考流 + 工具气泡 + 字节数进度 + 流式草稿）——//
interface ActivityStep {
  name: string
  brief: string
  bytes: number
  running: boolean
}
const activity = ref<{ steps: ActivityStep[]; draft: string; think: string } | null>(null)
const chatBody = ref<HTMLElement | null>(null)

// 思考流预览：折叠行显示最新一行（与文稿对话 ChatMessages 的 think 块同一交互）
function thinkPreview(text: string): string {
  const lines = text.split('\n').map((l) => l.trim()).filter(Boolean)
  return lines.length ? lines[lines.length - 1] : ''
}

const TOOL_LABEL: Record<string, string> = {
  write_tokens: '改设计 token',
  write_style: '重写 style.css',
  write_demo: '重写页面结构',
  set_meta: '改名称/描述',
  set_layout_roles: '改版式适用场景',
  set_layout_meta: '改版式名称/用途',
  add_layout: '新增版式',
  remove_layout: '删除版式',
  finish: '汇报总结',
}

watch(
  () => [messages.value.length, activity.value?.steps.length, activity.value?.draft.length, activity.value?.think.length],
  async () => {
    await nextTick()
    chatBody.value?.scrollTo({ top: chatBody.value.scrollHeight })
  },
)

// —— 历史版本（每次用户动作一版，可回滚）——//
const showHistory = ref(false)

// —— 样式面板：style.css 手动编辑（与对话 write_style 同级，写入都过历史）——//
const styleOpen = ref(false)
const styleText = ref('')
const styleLoadedText = ref('')
const styleSaving = ref(false)
const styleDirty = computed(() => styleText.value !== styleLoadedText.value)

async function loadStyle() {
  try {
    styleLoadedText.value = await fetchTemplateStyle(utId)
    styleText.value = styleLoadedText.value
  } catch {
    // 读不到（如 published 前的瞬态）保持现状，不打断主流程
  }
}

async function saveStyle() {
  if (!styleDirty.value || styleSaving.value) return
  styleSaving.value = true
  try {
    await userTemplateApi.saveStyle(utId, styleText.value)
    styleLoadedText.value = styleText.value
    previewKey.value += 1
    toast.success('样式已保存，并记入历史')
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '样式保存失败')
  } finally {
    styleSaving.value = false
  }
}

function toggleStyle() {
  styleOpen.value = !styleOpen.value
  if (styleOpen.value) void loadStyle()
}

function onRestored() {
  refreshPreview()
  if (styleOpen.value) void loadStyle()
  void loadContract() // 回滚可能改了结构契约
  void load()
}

// —— 版式面板（结构契约）：预览 | 版式 双视图切换 ——//
// 面板展示"生成侧实际生效的契约"（后端读注册表挂载快照）；role 点选即保存。

const ROLE_LABELS: Record<string, string> = {
  cover: '封面', toc: '目录', divider: '章节', content: '正文',
  data: '数据', quote: '引用', code: '代码', cta: '行动', thanks: '收尾',
}
const ROLE_ORDER = ['cover', 'toc', 'divider', 'content', 'data', 'quote', 'code', 'cta', 'thanks']

const mainView = ref<'preview' | 'layouts'>('preview')
const contract = ref<StructureContract | null>(null)
const contractLoading = ref(false)
const savingLayout = ref('')
const expandedSkeleton = ref('')
const editingLayout = ref('')
const editName = ref('')
const editUse = ref('')

const published = computed(() => row.value?.status === 'published' || row.value?.status === 'publishing')

/** 版式 → 演示它的 demo 页码列表（徽标用） */
const demoByLayout = computed(() => {
  const m = new Map<string, number[]>()
  for (const p of contract.value?.demo_pages ?? []) {
    const arr = m.get(p.layout) ?? []
    arr.push(p.no)
    m.set(p.layout, arr)
  }
  return m
})

/** 总页数（翻页上界与「第 N / M 页」）——必须在 contract 声明之后（TDZ） */
const totalPages = computed(() => contract.value?.demo_pages.length ?? 0)
// 对话可能增删 demo 页：契约重拉后页数收缩时把当前页钳回界内
watch(totalPages, (n) => {
  if (n && demoPage.value > n) {
    demoPage.value = n
    if (previewReady.value) gotoPreviewPage(n - 1)
  }
})

function switchToLayouts() {
  mainView.value = 'layouts'
  if (!contract.value && !contractLoading.value) void loadContract()
}

async function loadContract() {
  contractLoading.value = true
  try {
    contract.value = await userTemplateApi.structure(utId)
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '结构契约加载失败')
  } finally {
    contractLoading.value = false
  }
}

async function saveLayout(l: StructureLayout, patch: LayoutMetaPatch) {
  if (savingLayout.value || published.value) return
  savingLayout.value = l.id
  try {
    const res = await userTemplateApi.updateLayout(utId, l.id, patch)
    if (res.warning) toast.error(res.warning)
    else toast.success('已保存并记入历史')
    await loadContract()
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '版式保存失败')
  } finally {
    savingLayout.value = ''
  }
}

async function toggleRole(l: StructureLayout, role: string) {
  if (savingLayout.value || published.value) return
  const cur = l.roles ?? []
  const next = cur.includes(role) ? cur.filter((r) => r !== role) : [...cur, role]
  if (next.length > 3) {
    toast.error('一个版式最多 3 个适用场景')
    return
  }
  await saveLayout(l, { roles: next })
}

function startEditLayout(l: StructureLayout) {
  editingLayout.value = l.id
  editName.value = l.name
  editUse.value = l.use ?? ''
}

async function saveLayoutMeta(l: StructureLayout) {
  const patch: LayoutMetaPatch = {}
  if (editName.value.trim() !== l.name) patch.name = editName.value.trim()
  if (editUse.value.trim() !== (l.use ?? '')) patch.use = editUse.value.trim()
  if (!('name' in patch) && !('use' in patch)) {
    editingLayout.value = ''
    return
  }
  await saveLayout(l, patch)
  if (!savingLayout.value) editingLayout.value = ''
}

async function load() {
  try {
    row.value = await userTemplateApi.get(utId)
    name.value = row.value.name
    if (row.value.publish_error) {
      messages.value.push({ role: 'agent', text: `上次发布失败：${row.value.publish_error}` })
    }
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '模板加载失败')
    router.push('/my-templates')
  }
}

onMounted(() => {
  window.addEventListener('message', onPreviewMessage)
  void load()
  void loadContract() // 总页数与版式面板共用，挂载即拉
})
onBeforeUnmount(() => window.removeEventListener('message', onPreviewMessage))

async function send() {
  const text = input.value.trim()
  const imgs = [...pendingImages.value]
  if ((!text && !imgs.length) || sending.value || attaching.value) return
  input.value = ''
  pendingImages.value = []
  messages.value.push({ role: 'user', text, images: imgs.length ? imgs : undefined })
  sending.value = true
  activity.value = { steps: [], draft: '', think: '' }
  const stepByIndex = new Map<number, ActivityStep>()
  const stepByCall = new Map<string, ActivityStep>()
  try {
    await customizeChatStream(utId, text, imgs, (ev) => {
      const act = activity.value
      if (!act) return
      switch (ev.type) {
        case 'tool_start': {
          const step: ActivityStep = { name: ev.tool_name ?? '', brief: '', bytes: 0, running: true }
          act.steps.push(step)
          if (ev.tool_index != null) stepByIndex.set(ev.tool_index, step)
          if (ev.tool_call_id) stepByCall.set(ev.tool_call_id, step)
          act.draft = '' // 新工具开始 = 前面那段正文只是过渡话，收起
          break
        }
        case 'tool_progress': {
          const step = ev.tool_index != null ? stepByIndex.get(ev.tool_index) : undefined
          if (step) step.bytes = ev.bytes ?? 0
          break
        }
        case 'tool_done': {
          const step = ev.tool_call_id ? stepByCall.get(ev.tool_call_id) : undefined
          if (step) {
            step.running = false
            step.brief = ev.content ?? ''
          }
          break
        }
        case 'delta':
          act.draft += ev.content ?? ''
          break
        case 'think':
          act.think += ev.content ?? ''
          break
        case 'done':
          messages.value.push({
            role: 'agent',
            text: ev.reply || '（本轮无回复）',
            think: act.think || undefined, // 思考过程随消息留存（与文稿对话一致，可回看）
          })
          if (ev.dirty) {
            refreshPreview() // 已写入磁盘，整帧重建回第 1 页
            if (styleOpen.value) void loadStyle() // 对话可能整体重写过 style.css
            void loadContract() // 对话可能改了版式/页数——总页数与版式面板共用
          }
          break
        case 'error':
          messages.value.push({ role: 'agent', text: ev.content || '定制失败' })
          break
      }
    })
  } catch (e) {
    messages.value.push({ role: 'agent', text: e instanceof Error ? e.message : '定制失败' })
  } finally {
    sending.value = false
    activity.value = null
  }
}

async function runCheckup() {
  if (checking.value) return
  checking.value = true
  try {
    checkupReport.value = await userTemplateApi.checkup(utId)
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '体检失败')
  } finally {
    checking.value = false
  }
}

async function rename() {
  if (!name.value.trim() || name.value === row.value?.name) return
  try {
    await userTemplateApi.updateMeta(utId, name.value.trim())
    toast.info('已保存名称')
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  }
}

async function publish() {
  publishing.value = true
  try {
    await userTemplateApi.publish(utId)
    toast.info('发布成功，已进入社区模板')
    await load()
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '发布门禁未过')
    await load()
  } finally {
    publishing.value = false
  }
}

async function unpublish() {
  try {
    await userTemplateApi.unpublish(utId)
    toast.info('已下架')
    await load()
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '下架失败')
  }
}
</script>

<template>
  <div class="mx-auto max-w-[1200px] p-6">
    <div class="mb-4 flex flex-wrap items-center gap-2">
      <input
        v-model="name"
        class="min-w-0 flex-1 rounded-control border border-line bg-surface px-3 py-1.5 text-[14px] font-semibold text-ink outline-none focus-visible:border-accent"
        placeholder="模板名称"
        @change="rename"
      />
      <Button @click="router.push('/my-templates')">返回</Button>
      <Button @click="showHistory = true">
        <PhClockCounterClockwise :size="12" />
        历史
      </Button>
      <Button v-if="row?.visibility === 'public'" @click="unpublish">
        <PhGlobe :size="12" />
        下架
      </Button>
      <Button :loading="checking" title="渲染量测：溢出/填充率/最小字号；只报告，不拦发布" @click="runCheckup">
        <PhPulse :size="12" />
        质量体检
      </Button>
      <Button variant="primary" :loading="publishing" title="挂载校验 + 公开，秒级完成" @click="publish">
        <PhGlobe :size="12" />
        {{ row?.visibility === 'public' ? '重新发布' : '发布到社区' }}
      </Button>
    </div>
    <p v-if="row?.publish_error" class="mb-3 rounded-control border border-danger-soft bg-danger-soft p-2 text-[12px] text-danger">
      上次发布失败：{{ row.publish_error }}
    </p>

    <!-- 质量体检报告（只报告不改状态；发布同样秒级可重试） -->
    <div
      v-if="checkupReport"
      class="mb-3 rounded-card border p-3 text-[12.5px]"
      :class="checkupReport.render?.flaws?.length ? 'border-warning-soft bg-warning-soft' : 'border-success-soft bg-success-soft'"
    >
      <div class="flex items-center justify-between">
        <span class="font-semibold">
          质量体检{{ checkupReport.render?.flaws?.length ? `：发现 ${checkupReport.render.flaws.length} 个问题` : '：未发现问题' }}
        </span>
        <button class="cursor-pointer text-[11px] text-ink-3 hover:text-ink" @click="checkupReport = null">收起</button>
      </div>
      <p v-if="checkupReport.render" class="mt-1 text-ink-2">
        共 {{ checkupReport.render.pages }} 页 · 最低填充率 {{ checkupReport.render.min_fill }}% · 最小字号
        {{ checkupReport.render.max_flag_font }}px（报告不影响发布，存档在模板上）
      </p>
      <p v-if="checkupReport.structure && checkupReport.structure !== 'ok'" class="mt-1 text-danger">
        结构/安全：{{ checkupReport.structure }}
      </p>
      <ul v-if="checkupReport.render?.flaws?.length" class="mt-1.5 list-disc pl-5 text-ink-2">
        <li v-for="f in checkupReport.render.flaws" :key="f">{{ f }}</li>
      </ul>
      <div v-if="checkupReport.taste?.length" class="mt-1.5">
        <span class="font-semibold text-ink-2">AI 味提示（{{ checkupReport.taste.length }} 处）：示例文案是生成范本，建议修掉</span>
        <ul class="mt-1 list-disc pl-5 text-ink-2">
          <li v-for="t in checkupReport.taste" :key="t">{{ t }}</li>
        </ul>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-4 lg:grid-cols-3">
      <!-- 预览 / 版式 双视图 -->
      <div class="lg:col-span-2">
        <div class="overflow-hidden rounded-card border border-line bg-surface">
          <div class="flex items-center justify-between border-b border-line px-3 py-1.5 text-[12px] text-ink-2">
            <div class="flex items-center gap-1">
              <button
                class="cursor-pointer rounded px-2 py-0.5 transition-colors"
                :class="mainView === 'preview' ? 'bg-accent-soft font-semibold text-accent' : 'hover:text-ink'"
                @click="mainView = 'preview'"
              >
                实时预览
              </button>
              <button
                class="cursor-pointer rounded px-2 py-0.5 transition-colors"
                :class="mainView === 'layouts' ? 'bg-accent-soft font-semibold text-accent' : 'hover:text-ink'"
                @click="switchToLayouts"
              >
                版式
              </button>
            </div>
            <span v-if="mainView === 'preview'" class="flex items-center gap-1.5">
              <button
                class="cursor-pointer rounded border border-line px-2 py-0.5 transition-colors hover:border-line-strong disabled:cursor-not-allowed disabled:opacity-40"
                :disabled="demoPage <= 1"
                @click="prevPage"
              >
                ←
              </button>
              <span class="font-mono text-[11px]">第 {{ demoPage }}{{ totalPages ? ' / ' + totalPages : '' }} 页</span>
              <button
                class="cursor-pointer rounded border border-line px-2 py-0.5 transition-colors hover:border-line-strong disabled:cursor-not-allowed disabled:opacity-40"
                :disabled="totalPages > 0 && demoPage >= totalPages"
                @click="nextPage"
              >
                →
              </button>
              <button class="cursor-pointer rounded border border-line px-2 py-0.5 transition-colors hover:border-line-strong" title="刷新预览" @click="refreshPreview">
                <PhArrowClockwise :size="11" />
              </button>
              <button
                class="cursor-pointer rounded border px-2 py-0.5 transition-colors"
                :class="styleOpen ? 'border-accent text-accent' : 'border-line hover:border-line-strong'"
                title="编辑 style.css"
                @click="toggleStyle"
              >
                样式
              </button>
            </span>
            <span v-else class="text-[11px] text-ink-3">生成时模型按「适用场景 → 版式」选页型；点标签即改</span>
          </div>
          <div v-if="mainView === 'preview'" class="h-[460px] overflow-hidden bg-surface-2">
            <iframe
              v-if="row"
              ref="previewFrame"
              :key="previewKey"
              :src="previewSrc"
              class="h-full w-full border-0"
              sandbox="allow-scripts"
              title="模板定制预览"
            />
          </div>
          <!-- 版式面板：结构契约（只读部分）+ role 点选 / 名称用途编辑（写路径） -->
          <div v-else class="max-h-[620px] overflow-y-auto p-3">
            <div v-if="contractLoading" class="py-8 text-center text-[12px] text-ink-3">加载结构契约…</div>
            <div v-else-if="!contract" class="py-8 text-center text-[12px] text-ink-3">结构契约加载失败，请刷新重试</div>
            <template v-else>
              <p class="mb-3 text-[11.5px] leading-relaxed text-ink-3">
                共 {{ contract.layouts.length }} 个版式。「适用场景」决定生成时哪些页面角色会选中它（同一版式最多 3 个）。骨架与类名是生成侧的硬契约；增删版式走右侧对话（模型写骨架后自动过注册表校验，失败自动回滚）。
              </p>
              <div class="space-y-3">
                <div v-for="l in contract.layouts" :key="l.id" class="rounded-card border border-line bg-surface-2 p-3">
                  <div class="flex flex-wrap items-center gap-2">
                    <span class="text-[13.5px] font-semibold text-ink">{{ l.name }}</span>
                    <span class="font-mono text-[11px] text-ink-3">{{ l.id }}</span>
                    <span v-if="l.pattern" class="rounded-full border border-line px-2 py-0.5 font-mono text-[10.5px] text-ink-3" title="视觉指纹：节奏校验按它判「版面长得一样」的假多样性">{{ l.pattern }}</span>
                    <span
                      v-for="n in demoByLayout.get(l.id) ?? []"
                      :key="n"
                      class="rounded-full bg-accent-soft px-2 py-0.5 text-[10.5px] text-ink-2"
                      title="演示该版式的 demo 页"
                    >
                      示例 {{ n }}
                    </span>
                    <button
                      class="ml-auto cursor-pointer rounded border border-line px-2 py-0.5 text-[11px] transition-colors hover:border-line-strong disabled:cursor-not-allowed disabled:opacity-50"
                      :disabled="published"
                      :title="published ? '已发布模板先下架再改' : '编辑名称与用途'"
                      @click="startEditLayout(l)"
                    >
                      编辑
                    </button>
                  </div>
                  <p v-if="l.use" class="mt-1.5 text-[12px] leading-relaxed text-ink-2">{{ l.use }}</p>
                  <p v-if="l.constraints" class="mt-1 text-[11.5px] text-ink-3">约束：{{ l.constraints }}</p>
                  <p v-if="l.repeats && Object.keys(l.repeats).length" class="mt-1 text-[11.5px] text-ink-3">
                    数量契约：
                    <span v-for="(n, cls) in l.repeats" :key="cls" class="mr-2 font-mono">{{ cls }}={{ n }}</span>
                  </p>
                  <div class="mt-2 flex flex-wrap items-center gap-1.5">
                    <span class="text-[11px] text-ink-3">适用场景：</span>
                    <button
                      v-for="r in ROLE_ORDER"
                      :key="r"
                      class="cursor-pointer rounded-full border px-2.5 py-0.5 text-[11px] transition-colors disabled:cursor-not-allowed disabled:opacity-50"
                      :class="(l.roles ?? []).includes(r) ? 'border-accent bg-accent-soft font-semibold text-accent' : 'border-line text-ink-3 hover:border-line-strong'"
                      :disabled="published || savingLayout === l.id"
                      :title="published ? '已发布模板先下架再改' : (l.roles ?? []).includes(r) ? '点击取消该场景' : '点击加上该场景'"
                      @click="toggleRole(l, r)"
                    >
                      {{ ROLE_LABELS[r] }}
                    </button>
                  </div>
                  <div v-if="editingLayout === l.id" class="mt-2 space-y-2 rounded-control border border-line bg-surface p-2">
                    <input
                      v-model="editName"
                      class="w-full rounded-control border border-line bg-surface-2 px-2.5 py-1.5 text-[12.5px] text-ink outline-none focus-visible:border-accent"
                      placeholder="版式名称（≤40 字）"
                      maxlength="40"
                    />
                    <textarea
                      v-model="editUse"
                      class="h-[64px] w-full resize-y rounded-control border border-line bg-surface-2 px-2.5 py-1.5 text-[12.5px] leading-relaxed text-ink outline-none focus-visible:border-accent"
                      placeholder="用途：什么内容适合用这个版式（写具体一点，生成时模型按它选页型；≤200 字）"
                      maxlength="200"
                    />
                    <div class="flex gap-2">
                      <Button size="sm" variant="primary" :loading="savingLayout === l.id" @click="saveLayoutMeta(l)">保存</Button>
                      <Button size="sm" @click="editingLayout = ''">取消</Button>
                    </div>
                  </div>
                  <button
                    v-if="l.skeleton"
                    class="mt-2 cursor-pointer text-[11px] text-ink-3 hover:text-ink"
                    @click="expandedSkeleton = expandedSkeleton === l.id ? '' : l.id"
                  >
                    {{ expandedSkeleton === l.id ? '收起骨架' : '查看骨架' }}
                  </button>
                  <pre
                    v-if="expandedSkeleton === l.id"
                    class="mt-1 overflow-x-auto rounded-control bg-surface-2 p-2 font-mono text-[11px] leading-relaxed text-ink-2"
                  >{{ l.skeleton }}</pre>
                </div>
              </div>
              <div class="mt-4">
                <div class="mb-1 text-[11.5px] font-semibold text-ink-2">模板质量规则（rules.md，生成时注入模型）</div>
                <pre class="whitespace-pre-wrap rounded-control bg-surface-2 p-2 font-mono text-[11px] leading-relaxed text-ink-2">{{ contract.rules_md }}</pre>
              </div>
            </template>
          </div>
        </div>
        <template v-if="mainView === 'preview'">
        <!-- 样式面板：style.css 手动编辑（安全预检与历史同对话路径） -->
        <div v-if="styleOpen" class="mt-3 overflow-hidden rounded-card border border-line bg-surface">
          <div class="flex items-center justify-between border-b border-line px-3 py-1.5 text-[12px] text-ink-2">
            <span>style.css 手动编辑<span v-if="styleDirty" class="ml-2 text-warning">有未保存改动</span></span>
            <Button
              size="sm"
              variant="primary"
              :loading="styleSaving"
              :disabled="!styleDirty || row?.status === 'published'"
              :title="row?.status === 'published' ? '已发布模板先下架再改' : undefined"
              @click="saveStyle"
            >
              保存样式
            </Button>
          </div>
          <textarea
            v-model="styleText"
            class="h-[280px] w-full resize-y border-0 bg-surface p-3 font-mono text-[12px] leading-relaxed text-ink outline-none"
            spellcheck="false"
          />
        </div>
        <p class="mt-2 text-[11.5px] text-ink-3">
          对话里的每次改动都会实时落盘并记入历史；「质量体检」可随时量测溢出/填充率/字号，发布秒级完成。
        </p>
        </template>
      </div>

      <!-- 定制对话 -->
      <div class="flex min-h-[520px] flex-col rounded-card border border-line bg-surface">
        <div class="border-b border-line px-3 py-1.5 text-[12.5px] font-semibold text-ink-2">定制对话</div>
        <div ref="chatBody" class="min-h-0 flex-1 space-y-3 overflow-y-auto p-3">
          <div v-if="messages.length === 0 && !activity" class="rounded-control bg-surface-2 p-3 text-[12px] text-ink-3">
            试试：「主色换成暖橙色，整体更圆一点」「描述文案的弱色再淡一些，底色换成暖白」。
          </div>
          <div
            v-for="(m, i) in messages"
            :key="i"
            class="rounded-control px-3 py-2 text-[12.5px] leading-relaxed"
            :class="m.role === 'user' ? 'ml-6 bg-accent-soft text-ink' : 'mr-2 bg-surface-2 text-ink'"
          >
            <details v-if="m.think" class="mb-1 border-l-2 border-line-strong pl-2.5 text-[11.5px] text-ink-3">
              <summary class="cursor-pointer select-none hover:text-ink-2">已展开思考过程</summary>
              <div class="mt-1 max-h-48 overflow-y-auto whitespace-pre-wrap break-words">{{ m.think }}</div>
            </details>
            <div v-if="m.images?.length" class="mb-1.5 flex flex-wrap gap-1.5">
              <img v-for="src in m.images" :key="src" :src="src" loading="lazy" class="h-20 rounded-control border border-line" alt="附图" />
            </div>
            <template v-if="m.text">{{ m.text }}</template>
          </div>
          <!-- 进行中的工具气泡 + 生成进度 + 流式草稿 -->
          <div v-if="activity" class="mr-2 rounded-control bg-surface-2 px-3 py-2 text-[12.5px]">
            <!-- 思考流：折叠行实时显示最新一行（与文稿对话的思考块同一交互） -->
            <details v-if="activity.think" class="border-l-2 border-line-strong py-0.5 pl-2.5 text-[12px] text-ink-3">
              <summary class="cursor-pointer select-none truncate hover:text-ink-2" :title="activity.think">
                {{ thinkPreview(activity.think) || '思考中…' }}
              </summary>
              <div class="mt-1 max-h-48 overflow-y-auto whitespace-pre-wrap break-words">{{ activity.think }}</div>
            </details>
            <div v-for="(s, i) in activity.steps" :key="i" class="flex items-center gap-2 py-0.5">
              <PhCircleNotch v-if="s.running" :size="12" class="shrink-0 animate-spin text-accent" />
              <PhCheck v-else :size="12" class="shrink-0 text-success" />
              <span class="shrink-0 font-mono text-[11.5px] font-semibold">{{ s.name }}</span>
              <span class="shrink-0 text-[11px] text-ink-3">{{ TOOL_LABEL[s.name] ?? '执行中' }}</span>
              <span v-if="s.running && s.bytes" class="ml-auto shrink-0 font-mono text-[11px] text-ink-3">
                已生成 {{ (s.bytes / 1024).toFixed(0) }}KB…
              </span>
              <span v-else-if="!s.running && s.brief" class="line-clamp-1 ml-auto max-w-[55%] text-[11px] text-ink-3" :title="s.brief">
                {{ s.brief }}
              </span>
            </div>
            <p v-if="activity.draft" class="mt-1 whitespace-pre-wrap break-words text-ink-2">{{ activity.draft }}</p>
            <p v-if="!activity.steps.length && !activity.draft && !activity.think" class="animate-pulse text-ink-3">正在思考…</p>
          </div>
        </div>
        <div class="border-t border-line p-2">
          <!-- 待发附图条 -->
          <div v-if="pendingImages.length" class="mb-2 flex flex-wrap gap-1.5">
            <div v-for="(src, i) in pendingImages" :key="src" class="relative">
              <img :src="src" class="h-14 rounded-control border border-line" alt="待发送附图" />
              <button
                class="absolute -right-1.5 -top-1.5 flex h-4.5 w-4.5 items-center justify-center rounded-full bg-ink text-surface hover:bg-danger"
                :aria-label="`移除第 ${i + 1} 张附图`"
                @click="removePendingImage(i)"
              >
                <PhX :size="10" weight="bold" />
              </button>
            </div>
          </div>
          <div class="flex gap-2">
            <input
              ref="fileInput"
              type="file"
              accept="image/png,image/jpeg,image/webp"
              multiple
              class="hidden"
              @change="onPickImages"
            />
            <button
              class="flex shrink-0 items-center justify-center rounded-control border border-line px-2.5 text-ink-2 hover:border-accent hover:text-accent disabled:opacity-50"
              :disabled="sending || attaching || pendingImages.length >= MAX_CHAT_IMAGES"
              :title="pendingImages.length >= MAX_CHAT_IMAGES ? `一次最多附 ${MAX_CHAT_IMAGES} 张图` : '附图（最多 3 张，自动压缩）'"
              aria-label="附图"
              @click="fileInput?.click()"
            >
              <PhCircleNotch v-if="attaching" :size="14" class="animate-spin" />
              <PhImage v-else :size="14" />
            </button>
            <input
              v-model="input"
              class="min-w-0 flex-1 rounded-control border border-line bg-surface-2 px-3 py-1.5 text-[12.5px] outline-none focus-visible:border-accent"
              placeholder="描述想改的视觉，或附一张参考图（回车发送）"
              @keydown.enter="send"
            />
            <Button variant="primary" :loading="sending" aria-label="发送" @click="send">
              <PhPaperPlaneTilt :size="12" />
            </Button>
          </div>
        </div>
      </div>
    </div>

    <TemplateHistoryDrawer
      :open="showHistory"
      :template-id="utId"
      :published="row?.status === 'published'"
      @close="showHistory = false"
      @restored="onRestored"
    />
  </div>
</template>
