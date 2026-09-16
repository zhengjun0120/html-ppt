<script setup lang="ts">
import { computed, ref } from 'vue'

import type { ViewTool } from '@/types/chat'

const props = defineProps<{ tool: ViewTool }>()

const open = ref(false)

const chip = computed(() =>
  ({ running: { cls: 'bg-accent-soft text-accent', label: '运行中' },
    ok: { cls: 'bg-success-soft text-success', label: '成功' },
    error: { cls: 'bg-danger-soft text-danger', label: '失败' } })[props.tool.state],
)
</script>

<template>
  <div class="overflow-hidden rounded-[10px] border border-line bg-surface-2 text-[12.5px]">
    <button
      class="flex w-full cursor-pointer items-center gap-2 px-3 py-2 text-left"
      :aria-expanded="open"
      @click="open = !open"
    >
      <span class="font-mono font-semibold">{{ tool.name || '工具调用' }}</span>
      <span class="truncate text-ink-3">{{ tool.state === 'ok' ? `${tool.result.length} 字结果` : tool.state === 'running' ? '执行中…' : '调用失败' }}</span>
      <span class="ml-auto shrink-0 rounded-full px-2 py-0.5 text-[11px]" :class="chip.cls">{{ chip.label }}</span>
    </button>
    <div v-if="open" class="border-t border-line">
      <p v-if="tool.args" class="px-3 pt-2 text-[10.5px] text-ink-3">参数</p>
      <pre v-if="tool.args" class="max-h-32 overflow-auto whitespace-pre-wrap break-all bg-code-bg px-3 py-2 font-mono text-[11.5px] text-ink-2">{{ tool.args }}</pre>
      <p v-if="tool.result" class="px-3 pt-2 text-[10.5px] text-ink-3">结果</p>
      <pre v-if="tool.result" class="max-h-40 overflow-auto whitespace-pre-wrap break-all bg-code-bg px-3 py-2 font-mono text-[11.5px] text-ink-2">{{ tool.result }}</pre>
      <p v-if="!tool.result && tool.state === 'ok'" class="px-3 py-2 text-[11.5px] text-ink-3">（无文本结果）</p>
    </div>
  </div>
</template>
