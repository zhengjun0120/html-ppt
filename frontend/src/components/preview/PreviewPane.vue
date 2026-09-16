<script setup lang="ts">
import { PhArrowsOutSimple, PhArrowClockwise, PhClockCounterClockwise, PhPresentationChart } from '@phosphor-icons/vue'
import { computed, ref } from 'vue'

import { authedUrl } from '@/api/client'
import Empty from '@/components/ui/Empty.vue'
import { useDeckStore } from '@/stores/deck'

const props = defineProps<{ deckId: string }>()
const emit = defineEmits<{ history: [] }>()

const deckStore = useDeckStore()
const reloadKey = ref(0)
const title = computed(() => (props.deckId ? deckStore.titleOf(props.deckId) : '新文稿'))

const iframeSrc = computed(() =>
  props.deckId ? authedUrl(`/api/decks/${props.deckId}/file`) : '',
)

function refresh() {
  reloadKey.value += 1
}

function openInNewTab() {
  window.open(authedUrl(`/api/decks/${props.deckId}/file`), '_blank', 'noopener')
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
    <!-- 画布：deck 用自己的主题（内容，不是应用界面），深色外壳衬亮色片子 -->
    <div class="min-h-0 flex-1 bg-bg">
      <iframe
        v-if="props.deckId"
        :key="reloadKey"
        :src="iframeSrc"
        class="h-full w-full border-0"
        sandbox="allow-scripts allow-popups"
        title="deck 预览"
      />
      <div v-else class="flex h-full items-center justify-center p-6">
        <Empty
          title="文稿将在 agent 创建后可预览"
          desc="在右侧对话里描述你想要的演示文稿，agent 会创建文稿并逐页生成。生成后可从「文稿」页进入继续编辑。"
        >
          <template #icon><PhPresentationChart /></template>
        </Empty>
      </div>
    </div>
  </section>
</template>
