<script setup lang="ts">
import { PhArrowsOutSimple, PhArrowClockwise, PhClockCounterClockwise, PhDownload, PhExport } from '@phosphor-icons/vue'
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'

import { authedUrl } from '@/api/client'
import { exportDeck, thumbPages, thumbUrl } from '@/api/deckV2'
import Empty from '@/components/ui/Empty.vue'
import { useDeckStore } from '@/stores/deck'
import { useToast } from '@/stores/toast'

const props = defineProps<{ deckId: string }>()
const emit = defineEmits<{ history: [] }>()

const deckStore = useDeckStore()
const toast = useToast()
const reloadKey = ref(0)
const title = computed(() => (props.deckId ? deckStore.titleOf(props.deckId) : '新文稿'))

/** 翻页：iframe src 带 #/N 深链（runtime 监听 hashchange 落到对应页）。
 *  改 src 会整帧重载（本地文件，毫秒级），换来与缩略图/键盘的简单联动。 */
const page = ref(1)
const pages = ref<number[]>([])
const thumbsLoading = ref(false)
const thumbFailed = ref(false)

const iframeSrc = computed(() =>
  props.deckId ? authedUrl(`/api/decks/${props.deckId}/file`) + (page.value > 1 ? `#/${page.value}` : '') : '',
)

async function loadThumbs() {
  if (!props.deckId) return
  thumbsLoading.value = true
  thumbFailed.value = false
  try {
    pages.value = await thumbPages(props.deckId)
    if (pages.value.length && !pages.value.includes(page.value)) page.value = 1
  } catch {
    // 缩略图失败不阻塞预览主体（可能 Chrome 未配置）：隐藏侧栏，翻页条退化
    thumbFailed.value = true
    pages.value = []
  } finally {
    thumbsLoading.value = false
  }
}

function refresh() {
  reloadKey.value += 1
  void loadThumbs()
}

function openInNewTab() {
  window.open(authedUrl(`/api/decks/${props.deckId}/file`), '_blank', 'noopener')
}

function goto(p: number) {
  if (p < 1) return
  const max = pages.value.length
  page.value = max ? Math.min(p, Math.max(...pages.value)) : p
}

function onKey(e: KeyboardEvent) {
  if ((e.target as HTMLElement)?.tagName === 'INPUT' || (e.target as HTMLElement)?.tagName === 'TEXTAREA') return
  if (e.key === 'ArrowRight') goto(page.value + 1)
  else if (e.key === 'ArrowLeft') goto(page.value - 1)
}

onMounted(() => {
  void loadThumbs()
  window.addEventListener('keydown', onKey)
})
onBeforeUnmount(() => window.removeEventListener('keydown', onKey))
watch(() => props.deckId, () => {
  page.value = 1
  void loadThumbs()
})

const exporting = ref('')
async function doExport(format: 'pdf' | 'png' | 'html') {
  if (!props.deckId) return
  exporting.value = format
  try {
    const r = await exportDeck(props.deckId, format)
    window.open(authedUrl(`/api/decks/${props.deckId}/exports/${r.filename}`), '_blank', 'noopener')
    toast.info(`${r.filename} 已生成`)
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '导出失败')
  } finally {
    exporting.value = ''
  }
}
</script>

<template>
  <section class="flex min-w-0 flex-1 flex-col">
    <!-- 预览工具条 -->
    <div class="flex items-center gap-2 border-b border-line px-3 py-1.5">
      <span class="truncate text-[12.5px] font-semibold" :title="title">{{ title }}</span>
      <span v-if="props.deckId" class="truncate font-mono text-[11px] text-ink-3">{{ props.deckId }}</span>
      <span class="ml-auto flex items-center gap-1.5">
        <button
          v-if="props.deckId"
          class="inline-flex cursor-pointer items-center gap-1 rounded border border-line bg-surface-2 px-2 py-1 text-[11.5px] text-ink-2 transition-colors hover:border-line-strong hover:text-ink"
          :disabled="!!exporting"
          @click="doExport('pdf')"
        >
          <PhDownload :size="12" /> {{ exporting === 'pdf' ? '…' : 'PDF' }}
        </button>
        <button
          v-if="props.deckId"
          class="inline-flex cursor-pointer items-center gap-1 rounded border border-line bg-surface-2 px-2 py-1 text-[11.5px] text-ink-2 transition-colors hover:border-line-strong hover:text-ink"
          :disabled="!!exporting"
          @click="doExport('png')"
        >
          <PhExport :size="12" /> {{ exporting === 'png' ? '…' : 'PNG' }}
        </button>
        <button
          v-if="props.deckId"
          class="inline-flex cursor-pointer items-center gap-1 rounded border border-line bg-surface-2 px-2 py-1 text-[11.5px] text-ink-2 transition-colors hover:border-line-strong hover:text-ink"
          @click="emit('history')"
        >
          <PhClockCounterClockwise :size="12" /> 历史
        </button>
        <button
          v-if="props.deckId"
          class="inline-flex cursor-pointer items-center gap-1 rounded border border-line bg-surface-2 px-2 py-1 text-[11.5px] text-ink-2 transition-colors hover:border-line-strong hover:text-ink"
          @click="refresh"
        >
          <PhArrowClockwise :size="12" /> 刷新
        </button>
        <button
          v-if="props.deckId"
          class="inline-flex cursor-pointer items-center gap-1 rounded border border-line bg-surface-2 px-2 py-1 text-[11.5px] text-ink-2 transition-colors hover:border-line-strong hover:text-ink"
          @click="openInNewTab"
        >
          <PhArrowsOutSimple :size="12" /> 新标签
        </button>
      </span>
    </div>

    <div class="flex min-h-0 flex-1">
      <!-- 缩略图栏（渲染失败时隐藏，预览主体不受影响） -->
      <aside v-if="pages.length" class="w-[132px] shrink-0 overflow-y-auto border-r border-line bg-surface-2 p-2">
        <button
          v-for="p in pages"
          :key="p"
          class="mb-2 block w-full cursor-pointer overflow-hidden rounded border bg-surface transition-colors"
          :class="page === p ? 'border-accent shadow-[0_0_0_2px_var(--ring)]' : 'border-line hover:border-line-strong'"
          :title="`第 ${p} 页`"
          @click="goto(p)"
        >
          <img
            :src="thumbUrl(props.deckId, p)"
            :alt="`第 ${p} 页缩略图`"
            class="block w-full"
            loading="lazy"
          />
          <span class="block py-0.5 text-center font-mono text-[10px] text-ink-3">{{ p }}</span>
        </button>
      </aside>

      <!-- 画布 -->
      <div class="flex min-w-0 flex-1 flex-col">
        <div class="min-h-0 flex-1 bg-bg">
          <iframe
            v-if="props.deckId"
            :key="reloadKey + ':' + iframeSrc"
            :src="iframeSrc"
            class="h-full w-full border-0"
            sandbox="allow-scripts allow-popups"
            title="deck 预览"
          />
          <div v-else class="flex h-full items-center justify-center p-6">
            <Empty title="文稿将在 agent 创建后可预览" desc="回到对话里描述你想要的演示文稿，agent 会创建文稿并逐页生成。">
              <template #icon><PhDownload /></template>
            </Empty>
          </div>
        </div>
        <!-- 翻页条 -->
        <div v-if="props.deckId" class="flex items-center justify-center gap-2 border-t border-line bg-surface px-3 py-1.5 text-[12px] text-ink-2">
          <button
            class="cursor-pointer rounded border border-line px-2 py-0.5 transition-colors hover:border-line-strong disabled:cursor-not-allowed disabled:opacity-40"
            :disabled="page <= 1"
            @click="goto(page - 1)"
          >
            ← 上一页
          </button>
          <span class="font-mono text-[11.5px]">
            {{ page }}<span v-if="pages.length"> / {{ Math.max(...pages) }}</span>
          </span>
          <button
            class="cursor-pointer rounded border border-line px-2 py-0.5 transition-colors hover:border-line-strong disabled:cursor-not-allowed disabled:opacity-40"
            :disabled="pages.length > 0 && page >= Math.max(...pages)"
            @click="goto(page + 1)"
          >
            下一页 →
          </button>
          <span v-if="thumbsLoading" class="text-[10.5px] text-ink-3">缩略图渲染中…</span>
          <span class="hidden text-[10.5px] text-ink-3 sm:inline">方向键翻页</span>
        </div>
      </div>
    </div>
  </section>
</template>
