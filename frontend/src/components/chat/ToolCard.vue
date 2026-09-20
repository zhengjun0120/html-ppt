<script setup lang="ts">
import { computed, ref } from 'vue'

import { prettyArgs, prettyResult, toolLabel, toolSummary } from '@/lib/tools'
import type { ViewTool } from '@/types/chat'

const props = defineProps<{ tool: ViewTool }>()

const expanded = ref(false)
// 长结果（版式骨架、搜索结果）默认限高滚动，点开看全文
const resultFull = ref(false)

const chip = computed(() =>
  ({ running: { cls: 'bg-accent-soft text-accent', label: '运行中' },
    ok: { cls: 'bg-success-soft text-success', label: '成功' },
    error: { cls: 'bg-danger-soft text-danger', label: '失败' } })[props.tool.state],
)

// 头部一句话摘要：优先按工具定制的中文摘要，退回"N 字结果"
const headDesc = computed(() => {
  if (props.tool.state === 'running') return '执行中…'
  if (props.tool.state === 'error') return '调用失败'
  const s = toolSummary(props.tool.name || '', props.tool.args)
  return s || `${props.tool.result.length} 字结果`
})

// 运行中参数是流式到达的（tool_delta 增量）：露出最新的一段，让"在生成"可见
const streamingTail = computed(() => {
  if (props.tool.state !== 'running' || !props.tool.args) return ''
  const raw = props.tool.args
  return raw.length > 96 ? '…' + raw.slice(-96) : raw
})

const shownArgs = computed(() => prettyArgs(props.tool.args))
const shownResult = computed(() => prettyResult(props.tool.result))
</script>

<template>
  <div class="overflow-hidden rounded-[10px] border border-line bg-surface-2 text-[12.5px]">
    <button
      class="flex w-full cursor-pointer items-center gap-2 px-3 py-2 text-left"
      :aria-expanded="expanded"
      @click="expanded = !expanded"
    >
      <span class="shrink-0 font-semibold">{{ toolLabel(tool.name || '') }}</span>
      <span v-if="tool.name" class="shrink-0 font-mono text-[10.5px] text-ink-3">{{ tool.name }}</span>
      <span class="min-w-0 truncate text-ink-3">{{ headDesc }}</span>
      <span class="ml-auto shrink-0 rounded-full px-2 py-0.5 text-[11px]" :class="chip.cls">{{ chip.label }}</span>
    </button>

    <!-- 运行中：参数流式增量可见（不用展开也能看到"在生成"） -->
    <div v-if="tool.state === 'running' && streamingTail" class="border-t border-line px-3 py-1.5">
      <p class="truncate font-mono text-[11px] text-ink-3 animate-pulse">{{ streamingTail }}</p>
    </div>

    <div v-if="expanded" class="border-t border-line">
      <template v-if="tool.args">
        <p class="px-3 pt-2 text-[10.5px] text-ink-3">
          参数
          <span v-if="tool.state === 'running'" class="text-accent">（流式接收中…）</span>
        </p>
        <pre class="max-h-40 overflow-auto whitespace-pre-wrap break-all bg-code-bg px-3 py-2 font-mono text-[11.5px] text-ink-2">{{ shownArgs }}</pre>
      </template>
      <template v-if="tool.result">
        <p class="px-3 pt-2 text-[10.5px] text-ink-3">
          结果
          <button
            class="cursor-pointer text-accent hover:underline"
            @click.stop="resultFull = !resultFull"
          >
            {{ resultFull ? '收起' : '展开全文' }}
          </button>
        </p>
        <pre
          class="overflow-auto whitespace-pre-wrap break-all bg-code-bg px-3 py-2 font-mono text-[11.5px] text-ink-2"
          :class="resultFull ? '' : 'max-h-40'"
        >{{ shownResult }}</pre>
      </template>
      <p v-if="!tool.result && tool.state === 'ok'" class="px-3 py-2 text-[11.5px] text-ink-3">（无文本结果）</p>
    </div>
  </div>
</template>
