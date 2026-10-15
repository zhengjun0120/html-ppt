<script setup lang="ts">
import { computed } from 'vue'

import Badge from '@/components/ui/Badge.vue'
import Drawer from '@/components/ui/Drawer.vue'
import Empty from '@/components/ui/Empty.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import type { SessionSummary } from '@/types/chat'

/**
 * 会话列表抽屉：本文稿的全部 agent 对话，点选切换（plan-v3 的会话选择 UI）。
 * 数据由调用方传入（workspaceBoot 的 sessions ref），切换走现成的
 * switchSession → chat.loadSession 链路，这里只管展示与选中。
 */
const props = defineProps<{
  open: boolean
  sessions: SessionSummary[]
  currentId: number | null
  loading?: boolean
}>()
const emit = defineEmits<{ close: []; switch: [id: number] }>()

function fmtTime(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return ''
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

/** 相对时间做副标题：今天/昨天直接说，更早给日期 */
const relTime = computed(() => {
  const dayStart = new Date()
  dayStart.setHours(0, 0, 0, 0)
  return (iso: string) => {
    const d = new Date(iso)
    if (Number.isNaN(d.getTime())) return ''
    const diffDays = Math.floor((dayStart.getTime() - new Date(d).setHours(0, 0, 0, 0)) / 86400000)
    if (diffDays <= 0) return `今天 ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
    if (diffDays === 1) return `昨天 ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
    return fmtTime(iso)
  }
})

function pick(id: number) {
  if (id === (props.currentId ?? 0)) return
  emit('switch', id)
}
</script>

<template>
  <Drawer :open="open" title="会话列表" @close="emit('close')">
    <div class="p-4">
      <template v-if="loading">
        <div v-for="i in 3" :key="i" class="mb-2 rounded-control border border-line bg-surface p-3">
          <Skeleton class="h-4 w-2/3" />
          <Skeleton class="mt-2 h-3 w-1/3" />
        </div>
      </template>

      <Empty
        v-else-if="sessions.length === 0"
        title="还没有对话"
        desc="在对话面板里发出的第一句话会创建一个会话，之后可在这里来回切换。"
      />

      <ul v-else class="flex flex-col gap-2">
        <li v-for="s in sessions" :key="s.id">
          <button
            class="block w-full cursor-pointer rounded-control border bg-surface p-3 text-left transition-colors"
            :class="
              s.id === currentId
                ? 'border-accent bg-accent-soft'
                : 'border-line hover:border-line-strong'
            "
            @click="pick(s.id)"
          >
            <div class="flex items-center gap-2">
              <span class="min-w-0 flex-1 truncate text-[12.5px] font-semibold" :title="s.title">
                {{ s.title || `会话 ${s.id}` }}
              </span>
              <Badge v-if="s.pending" tone="warning">待回答</Badge>
              <Badge v-else-if="!s.has_messages" tone="neutral">空对话</Badge>
              <Badge v-if="s.id === currentId" tone="accent">当前</Badge>
            </div>
            <p class="mt-1 text-[11px] text-ink-3">{{ relTime(s.updated_at) }}</p>
          </button>
        </li>
      </ul>
    </div>
  </Drawer>
</template>
