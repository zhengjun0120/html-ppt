<script setup lang="ts">
import { computed } from 'vue'

import { useChatStore } from '@/stores/chat'
import { useWizardStore } from '@/stores/wizard'

/**
 * 生成进度（generating 阶段的主区域）：按大纲逐页展示写入状态。
 * 进度来自 SSE page_generated 事件（chat.pageProgress）；预览 iframe
 * 由 deckTouched 信号自动刷新，能看着页面逐批长出来。
 */

const chat = useChatStore()
const wizard = useWizardStore()

const pages = computed(() => wizard.outline?.pages ?? [])
const doneCount = computed(() => Object.values(chat.pageProgress).filter(Boolean).length)
const allDone = computed(() => pages.value.length > 0 && doneCount.value >= pages.value.length)

function stateOf(no: number): 'ok' | 'fail' | 'pending' {
  const v = chat.pageProgress[no]
  if (v === true) return 'ok'
  if (v === false) return 'fail'
  return 'pending'
}
</script>

<template>
  <div class="mx-auto flex h-full w-full max-w-[760px] flex-col items-center justify-center gap-5 p-6">
    <div class="text-center">
      <h2 class="text-[16px] font-semibold text-ink">{{ allDone ? '收尾与自查中…' : '正在按大纲生成页面' }}</h2>
      <p class="mt-1 text-[12.5px] text-ink-3">
        {{ doneCount }} / {{ pages.length }} 页已写入 · 量测与审查自动进行，右侧对话可看到全过程
      </p>
    </div>

    <div class="flex w-full flex-col gap-1.5">
      <div
        v-for="p in pages"
        :key="p.no"
        class="flex items-center gap-2.5 rounded-control border bg-surface px-3 py-1.5 text-[12.5px] transition-colors"
        :class="stateOf(p.no) === 'ok' ? 'border-line-strong' : stateOf(p.no) === 'fail' ? 'border-danger' : 'border-line opacity-70'"
      >
        <span class="w-8 shrink-0 font-mono text-[10.5px] text-ink-3">{{ String(p.no).padStart(2, '0') }}</span>
        <span
          class="inline-flex h-4 w-4 shrink-0 items-center justify-center rounded-full border text-[9px]"
          :class="stateOf(p.no) === 'ok' ? 'border-accent bg-accent text-on-accent' : stateOf(p.no) === 'fail' ? 'border-danger text-danger' : 'border-line text-transparent'"
        >
          {{ stateOf(p.no) === 'ok' ? '✓' : stateOf(p.no) === 'fail' ? '!' : '·' }}
        </span>
        <span class="min-w-0 flex-1 truncate text-ink" :class="stateOf(p.no) === 'pending' ? 'text-ink-3' : ''">{{ p.title }}</span>
        <span class="shrink-0 text-[10.5px] text-ink-3">
          {{ stateOf(p.no) === 'ok' ? '已写入' : stateOf(p.no) === 'fail' ? '写入失败' : '待生成' }}
        </span>
      </div>
    </div>

    <p class="text-[11.5px] text-ink-3">下方预览跟随写入实时生长；完成后自动进入「迭代」页继续修改。</p>
  </div>
</template>
