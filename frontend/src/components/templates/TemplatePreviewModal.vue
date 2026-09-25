<script setup lang="ts">
import { PhCaretDown, PhCaretUp, PhX } from '@phosphor-icons/vue'
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'

/**
 * 模板预览弹窗：整本 demo 放大到接近满屏，逐页翻看。
 *
 * 翻页走 runtime.js 的预览模式协议（src 里 ?preview=1 激活）：iframe 只加载
 * 一次，之后每次翻页向 iframe postMessage {type:'preview-goto', idx}，由内部
 * 无刷新切换——带着 base.css 的 .5s 淡入+位移过渡。此前"每页换 src 整个重建
 * iframe"的方案会白屏一闪、内容硬切，就是它被换掉的原因。
 *
 * - 握手：iframe 就绪后回报 {type:'preview-ready'}；就绪前点的翻页记在 page
 *   里，就绪时补发，不会丢。
 * - 键盘 ↑/↓、←/→、PgUp/PgDn 翻页，Esc 关闭。iframe pointer-events-none，
 *   焦点留在父页，按键不会被 demo 内部吞掉。
 * - 页数由调用方给（内置 = meta.demo_pages；用户模板 = 基模板的页数——fork
 *   时 index.html 原样拷贝，定制只改 style.css）。页数未知时"下一页"不设限，
 *   preview-goto 的 idx 由 runtime 钳在最后一页。
 */
const props = defineProps<{
  open: boolean
  title: string
  canvas: { w: number; h: number }
  pages?: number
  sandbox: string
  /** demo 首页地址，需带 ?preview=1 激活 runtime 的预览模式协议 */
  src: string
}>()
const emit = defineEmits<{ close: [] }>()

const page = ref(1)
const ready = ref(false)
const frameEl = ref<HTMLIFrameElement | null>(null)
const total = () => (props.pages && props.pages > 0 ? props.pages : undefined)

watch(
  () => props.open,
  (v) => {
    if (v) {
      page.value = 1
      ready.value = false
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

function onMessage(e: MessageEvent) {
  if (e.source !== frameEl.value?.contentWindow) return
  if ((e.data as { type?: string } | null)?.type === 'preview-ready') {
    ready.value = true
    goto(page.value - 1) // 就绪前的翻页在这里补发
  }
}
function onKey(e: KeyboardEvent) {
  if (!props.open) return
  if (e.key === 'Escape') {
    emit('close')
    return
  }
  if (e.key === 'ArrowUp' || e.key === 'ArrowLeft' || e.key === 'PageUp') {
    e.preventDefault()
    prev()
  } else if (e.key === 'ArrowDown' || e.key === 'ArrowRight' || e.key === 'PageDown' || e.key === ' ') {
    e.preventDefault()
    next()
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
</script>

<template>
  <Teleport to="body">
    <div
      v-if="open"
      class="fixed inset-0 z-50 flex flex-col bg-black/85 p-4 sm:p-6"
      @click.self="emit('close')"
    >
      <div class="flex items-center gap-3">
        <p class="truncate text-[14px] font-semibold text-white/90" :title="title">{{ title }}</p>
        <span class="ml-auto shrink-0 font-mono text-[12px] text-white/60">
          {{ total() ? `第 ${page} / ${total()} 页` : `第 ${page} 页` }}
        </span>
        <button
          type="button"
          class="shrink-0 cursor-pointer rounded-full p-1.5 text-white/70 transition-colors hover:bg-white/10 hover:text-white"
          title="关闭（Esc）"
          @click="emit('close')"
        >
          <PhX :size="18" />
        </button>
      </div>

      <div class="flex min-h-0 flex-1 items-center justify-center py-3" @click.self="emit('close')">
        <div
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
            :sandbox="sandbox"
            class="pointer-events-none absolute inset-0 h-full w-full border-0"
            :title="title + ' 预览'"
          />
        </div>
      </div>

      <div class="flex items-center justify-center gap-3">
        <button
          type="button"
          class="inline-flex cursor-pointer items-center gap-1.5 rounded-full border border-white/20 px-4 py-1.5 text-[13px] text-white/85 transition-colors hover:border-white/40 hover:bg-white/10 disabled:cursor-not-allowed disabled:opacity-35"
          :disabled="page <= 1"
          title="上一页（↑）"
          @click="prev"
        >
          <PhCaretUp :size="13" />
          上一页
        </button>
        <button
          type="button"
          class="inline-flex cursor-pointer items-center gap-1.5 rounded-full border border-white/20 px-4 py-1.5 text-[13px] text-white/85 transition-colors hover:border-white/40 hover:bg-white/10 disabled:cursor-not-allowed disabled:opacity-35"
          :disabled="total() != null && page >= total()!"
          title="下一页（↓）"
          @click="next"
        >
          下一页
          <PhCaretDown :size="13" />
        </button>
      </div>
    </div>
  </Teleport>
</template>
