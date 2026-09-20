<script setup lang="ts">
import { computed, ref } from 'vue'

import ToolCard from './ToolCard.vue'
import { toolLabel } from '@/lib/tools'
import type { ViewTool } from '@/types/chat'

/**
 * 相邻同类型工具的折叠组：一轮 run 里连读 N 个版式/连写 N 页时，
 * 收成一张"中文名 ×N"的卡，展开看每一笔的参数与结果。
 * 最新一笔的状态代表整组（有失败显示失败，全成功显示成功）。
 */
const props = defineProps<{ name: string; tools: ViewTool[] }>()

const expanded = ref(false)

const groupChip = computed(() => {
  if (props.tools.some((t) => t.state === 'running'))
    return { cls: 'bg-accent-soft text-accent', label: '运行中' }
  if (props.tools.some((t) => t.state === 'error'))
    return { cls: 'bg-danger-soft text-danger', label: '有失败' }
  return { cls: 'bg-success-soft text-success', label: '成功' }
})
</script>

<template>
  <div class="overflow-hidden rounded-[10px] border border-line bg-surface-2 text-[12.5px]">
    <button
      class="flex w-full cursor-pointer items-center gap-2 px-3 py-2 text-left"
      :aria-expanded="expanded"
      @click="expanded = !expanded"
    >
      <span class="shrink-0 font-semibold">{{ toolLabel(name) }}</span>
      <span class="font-mono text-[10.5px] text-ink-3">×{{ tools.length }}</span>
      <span class="text-ink-3">{{ expanded ? '收起这批调用' : '展开查看每次调用' }}</span>
      <span class="ml-auto shrink-0 rounded-full px-2 py-0.5 text-[11px]" :class="groupChip.cls">{{ groupChip.label }}</span>
    </button>
    <div v-if="expanded" class="flex flex-col gap-1.5 border-t border-line p-2">
      <ToolCard v-for="t in tools" :key="t.id" :tool="t" />
    </div>
  </div>
</template>
