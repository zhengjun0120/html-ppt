<script setup lang="ts">
import { PhArrowClockwise, PhArrowCounterClockwise, PhCaretDown, PhCaretLeft, PhCaretRight, PhFloppyDisk, PhX } from '@phosphor-icons/vue'
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'

import { authedUrl } from '@/api/client'
import { deckApi } from '@/api/decks'
import { userTemplateApi } from '@/api/userTemplates'
import { useToast } from '@/stores/toast'

/**
 * 编辑弹窗（docs/deck-editor-plan.md §4.6）：预览弹窗里直接改 deck / 用户模板。
 *
 * 与只读的 TemplatePreviewModal 相反的交互前提：iframe 必须可点（选中/拖拽/双击
 * 进文本编辑都在 iframe 里发生，键盘事件随之进 iframe），所以父页不装全局翻页
 * 键盘、不拦 Esc 关闭——关闭走工具栏按钮 + 脏更改内联确认。所有编辑器本体
 * （editor.js）跑在 iframe 内，父页只做外壳：保存流、翻页（preview-goto）、
 * 撤销/重做按钮、脏标记。
 *
 * - src：deck 走 /api/decks/:id/file，用户模板走 /api/user-templates/:id/editor
 *   （静态路由没有服务端钩子，注入 editor.js 只能走受保护端点）。
 *   ?preview=1&edit=1：preview 让 runtime 进预览模式（单页 + preview-goto），
 *   edit 让后端注入 editor.js。
 * - 握手：editor-ready {pages}；8s 未到 → 该对象不支持在线编辑（v1 deck 等）。
 * - 保存：editor-save → editor-serialize {html} → PUT 全量 → editor-saved。
 *   文稿写 history 版本；模板覆盖 + 服务端滚动备份。
 */
const props = defineProps<{
  open: boolean
  kind: 'deck' | 'usertpl'
  id: string
  title: string
  canvas: { w: number; h: number }
  /** 页数兜底（editor-ready 之前显示用）；到了 ready 以 ready 的为准 */
  pages?: number
}>()
const emit = defineEmits<{ close: []; saved: [] }>()

const toast = useToast()
const frameEl = ref<HTMLIFrameElement | null>(null)
const ready = ref(false)
const failed = ref(false)
const dirty = ref(false)
const saving = ref(false)
const page = ref(1)
const totalPages = ref(0)
const confirmDiscard = ref(false)
/** 选中态（editor-selection 上报）：字号步进器只在有选中时出现 */
const hasSel = ref(false)
const fontPx = ref<number | null>(null)

const src = computed(() => {
  if (!props.id) return ''
  return props.kind === 'deck'
    ? authedUrl(`/api/decks/${props.id}/file`, { preview: 1, edit: 1 })
    : authedUrl(`/api/user-templates/${props.id}/editor`, { preview: 1, edit: 1 })
})
const total = () => (totalPages.value > 0 ? totalPages.value : props.pages && props.pages > 0 ? props.pages : 0)

watch(
  () => props.open,
  (v) => {
    if (v) {
      page.value = 1
      ready.value = false
      failed.value = false
      dirty.value = false
      confirmDiscard.value = false
      totalPages.value = 0
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = ''
    }
  },
)

function goto(i: number) {
  frameEl.value?.contentWindow?.postMessage({ type: 'preview-goto', idx: i }, '*')
}
function prev() {
  if (page.value <= 1) return
  page.value -= 1
  if (ready.value) goto(page.value - 1)
}
function next() {
  const n = total()
  if (n && page.value >= n) return
  page.value += 1
  if (ready.value) goto(page.value - 1)
}

function sendEditor(msg: Record<string, unknown>) {
  frameEl.value?.contentWindow?.postMessage(msg, '*')
}
function doUndo() {
  sendEditor({ type: 'editor-undo' })
}
function doRedo() {
  sendEditor({ type: 'editor-redo' })
}
/** 字号步进：等比缩放选中块及其内部文字（editor.js applyFontFactor） */
function fontStep(factor: number) {
  sendEditor({ type: 'editor-font', factor })
}

async function save() {
  if (!ready.value || saving.value || failed.value) return
  saving.value = true
  try {
    const html = await requestSerialize()
    if (props.kind === 'deck') {
      await deckApi.saveFile(props.id, html)
    } else {
      await userTemplateApi.saveFile(props.id, html)
    }
    sendEditor({ type: 'editor-saved' })
    dirty.value = false
    toast.info('已保存')
    emit('saved')
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    saving.value = false
  }
}

/** editor-save → editor-serialize 的应答，2.5s 超时防悬挂 */
function requestSerialize(): Promise<string> {
  return new Promise((resolve, reject) => {
    const frame = frameEl.value
    if (!frame) return reject(new Error('画布未就绪'))
    const onMsg = (e: MessageEvent) => {
      if (e.source !== frame.contentWindow) return
      const d = e.data as { type?: string; html?: string } | null
      if (d?.type === 'editor-serialize') {
        window.removeEventListener('message', onMsg)
        if (typeof d.html === 'string' && d.html.length > 0) resolve(d.html)
        else reject(new Error('编辑器返回了空内容'))
      }
    }
    window.addEventListener('message', onMsg)
    sendEditor({ type: 'editor-save' })
    setTimeout(() => {
      window.removeEventListener('message', onMsg)
      reject(new Error('编辑器无响应'))
    }, 2500)
  })
}

function onCloseClick() {
  if (dirty.value && !confirmDiscard.value) {
    confirmDiscard.value = true
    return
  }
  emit('close')
}

function onMessage(e: MessageEvent) {
  if (e.source !== frameEl.value?.contentWindow) return
  const d = e.data as { type?: string; pages?: number; dirty?: boolean; html?: string; selected?: boolean; fontSize?: number | null } | null
  if (!d || typeof d !== 'object') return
  if (d.type === 'editor-ready') {
    ready.value = true
    failed.value = false
    if (typeof d.pages === 'number' && d.pages > 0) totalPages.value = d.pages
    goto(page.value - 1) // 就绪前点的翻页在这里补发
  } else if (d.type === 'editor-dirty') {
    dirty.value = !!d.dirty
    if (d.dirty) confirmDiscard.value = false
  } else if (d.type === 'editor-selection') {
    hasSel.value = !!d.selected
    fontPx.value = typeof d.fontSize === 'number' ? d.fontSize : null
  }
}

// 父页聚焦时的 Ctrl+S（iframe 聚焦时按键不经过这里，由 editor.js 自己拦）
function onKey(e: KeyboardEvent) {
  if (!props.open) return
  if ((e.ctrlKey || e.metaKey) && (e.key === 's' || e.key === 'S')) {
    e.preventDefault()
    void save()
  } else if (e.key === 'Escape') {
    onCloseClick()
  }
}

onMounted(() => {
  window.addEventListener('message', onMessage)
  window.addEventListener('keydown', onKey)
})
onBeforeUnmount(() => {
  window.removeEventListener('message', onMessage)
  window.removeEventListener('keydown', onKey)
  document.body.style.overflow = ''
})

// editor-ready 超时兜底：v1 deck / 后端未部署 / 脚本 404
let readyTimer: ReturnType<typeof setTimeout> | null = null
watch(
  () => [props.open, src.value] as const,
  ([o]) => {
    if (readyTimer) clearTimeout(readyTimer)
    if (!o) return
    readyTimer = setTimeout(() => {
      if (props.open && !ready.value) failed.value = true
    }, 8000)
  },
  { immediate: true },
)
</script>

<template>
  <Teleport to="body">
    <div v-if="open" class="fixed inset-0 z-50 flex flex-col bg-black/85 p-4 sm:p-6">
      <!-- 工具栏 -->
      <div class="flex items-center gap-2.5">
        <p class="truncate text-[14px] font-semibold text-white/90" :title="title">{{ title }}</p>
        <span
          v-if="dirty"
          class="shrink-0 rounded-full bg-amber-400/90 px-2 py-0.5 text-[10.5px] font-semibold text-amber-950"
          title="有未保存的修改"
        >
          未保存
        </span>
        <span class="ml-auto flex shrink-0 items-center gap-1.5 text-white/80">
          <template v-if="ready && !failed">
            <button
              type="button"
              class="cursor-pointer rounded-full p-1.5 transition-colors hover:bg-white/10 hover:text-white"
              title="撤销（Ctrl+Z）"
              @click="doUndo"
            >
              <PhArrowCounterClockwise :size="15" />
            </button>
            <button
              type="button"
              class="cursor-pointer rounded-full p-1.5 transition-colors hover:bg-white/10 hover:text-white"
              title="重做（Ctrl+Shift+Z）"
              @click="doRedo"
            >
              <PhArrowClockwise :size="15" />
            </button>
            <!-- 字号步进：选中元素后出现，作用于选中块及其内部文字（等比缩放） -->
            <template v-if="hasSel">
              <span class="mx-0.5 hidden h-4 w-px bg-white/20 sm:block"></span>
              <button
                type="button"
                class="cursor-pointer rounded-full px-1.5 py-0.5 text-[11px] font-bold transition-colors hover:bg-white/10 hover:text-white"
                title="减小字号（Ctrl+-）"
                @click="fontStep(0.9)"
              >
                A−
              </button>
              <span class="w-8 text-center font-mono text-[11.5px] text-white/70" title="选中文字的当前字号（px）">
                {{ fontPx != null ? Math.round(fontPx) : '-' }}
              </span>
              <button
                type="button"
                class="cursor-pointer rounded-full px-1.5 py-0.5 text-[13px] font-bold transition-colors hover:bg-white/10 hover:text-white"
                title="增大字号（Ctrl+=）"
                @click="fontStep(1.1)"
              >
                A+
              </button>
            </template>
            <button
              type="button"
              class="inline-flex cursor-pointer items-center gap-1.5 rounded-full border border-white/25 px-3.5 py-1.5 text-[12.5px] font-semibold transition-colors hover:border-white/50 hover:bg-white/10 disabled:cursor-not-allowed disabled:opacity-40"
              :class="dirty ? 'border-amber-300/70 text-amber-200 hover:border-amber-200 hover:bg-amber-400/10' : ''"
              :disabled="saving"
              title="保存（Ctrl+S）"
              @click="save"
            >
              <PhFloppyDisk :size="14" />
              {{ saving ? '保存中…' : '保存' }}
            </button>
            <span class="mx-1 hidden h-4 w-px bg-white/20 sm:block"></span>
            <button
              type="button"
              class="cursor-pointer rounded-full p-1.5 transition-colors hover:bg-white/10 hover:text-white disabled:cursor-not-allowed disabled:opacity-35"
              :disabled="page <= 1"
              title="上一页"
              @click="prev"
            >
              <PhCaretLeft :size="14" />
            </button>
            <span class="font-mono text-[12px] text-white/70">
              {{ page }}<span v-if="total()"> / {{ total() }}</span>
            </span>
            <button
              type="button"
              class="cursor-pointer rounded-full p-1.5 transition-colors hover:bg-white/10 hover:text-white disabled:cursor-not-allowed disabled:opacity-35"
              :disabled="total() != null && page >= total()!"
              title="下一页"
              @click="next"
            >
              <PhCaretRight :size="14" />
            </button>
            <span class="mx-1 hidden h-4 w-px bg-white/20 sm:block"></span>
          </template>
          <!-- 脏更改关闭确认（内联条，不弹原生 confirm） -->
          <template v-if="confirmDiscard">
            <span class="text-[12.5px] text-amber-200">未保存的修改将丢失</span>
            <button
              type="button"
              class="cursor-pointer rounded-full bg-amber-400/90 px-3 py-1 text-[12px] font-semibold text-amber-950 transition-colors hover:bg-amber-300"
              @click="emit('close')"
            >
              放弃修改
            </button>
            <button
              type="button"
              class="cursor-pointer rounded-full border border-white/25 px-3 py-1 text-[12px] text-white/85 transition-colors hover:bg-white/10"
              @click="confirmDiscard = false"
            >
              继续编辑
            </button>
          </template>
          <button
            v-else
            type="button"
            class="cursor-pointer rounded-full p-1.5 text-white/70 transition-colors hover:bg-white/10 hover:text-white"
            title="关闭"
            @click="onCloseClick"
          >
            <PhX :size="18" />
          </button>
        </span>
      </div>

      <!-- 画布 -->
      <div class="flex min-h-0 flex-1 items-center justify-center py-3">
        <div
          v-if="!failed"
          class="relative overflow-hidden rounded-control bg-surface-2 shadow-2xl"
          :style="{
            width: `min(94vw, calc((100vh - 130px) * ${canvas.w / canvas.h}))`,
            aspectRatio: `${canvas.w} / ${canvas.h}`,
          }"
        >
          <iframe
            ref="frameEl"
            :key="src"
            :src="src"
            sandbox="allow-scripts"
            class="absolute inset-0 h-full w-full border-0"
            :title="title + ' 编辑'"
          />
        </div>
        <div v-else class="max-w-md text-center">
          <p class="text-[14px] font-semibold text-white/90">这份内容不支持在线编辑</p>
          <p class="mt-2 text-[12.5px] leading-relaxed text-white/60">
            仅 deck-v2 文稿和自己的用户模板可以进入编辑器。请检查后刷新重试，或回到对话里让 agent 修改。
          </p>
        </div>
      </div>

      <!-- 底部提示 -->
      <div class="flex items-center justify-center gap-4 text-[11.5px] text-white/45">
        <template v-if="ready && !failed">
          <span>双击文字直接修改</span>
          <span>·</span>
          <span>点选卡片后拖动 / 缩放</span>
          <span>·</span>
          <span>Esc 取消选中</span>
          <span>·</span>
          <span>Ctrl+S 保存</span>
        </template>
        <template v-else-if="!failed">
          <span class="inline-flex items-center gap-1.5">
            <PhCaretDown :size="12" class="animate-bounce" />
            编辑器加载中…
          </span>
        </template>
      </div>
    </div>
  </Teleport>
</template>
