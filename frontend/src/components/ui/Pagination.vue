<script setup lang="ts">
import { PhCaretLeft, PhCaretRight } from '@phosphor-icons/vue'
import { computed } from 'vue'

/**
 * 分页控件（配纯前端切片用）：给 page-count，出「‹ 1 … 4 5 6 … 20 ›」。
 * 当前页两侧各留 1 页，首尾页常驻，中间折叠成省略号；只有一页时不渲染——
 * 列表小的时候调用方不需要任何条件判断。
 */
const props = defineProps<{
  modelValue: number
  pageCount: number
}>()
const emit = defineEmits<{ 'update:modelValue': [page: number] }>()

const pages = computed<(number | '…')[]>(() => {
  const n = props.pageCount
  if (n <= 7) return Array.from({ length: n }, (_, i) => i + 1)
  const cur = props.modelValue
  const lo = Math.max(2, cur - 1)
  const hi = Math.min(n - 1, cur + 1)
  const out: (number | '…')[] = [1]
  if (lo > 2) out.push('…')
  for (let i = lo; i <= hi; i++) out.push(i)
  if (hi < n - 1) out.push('…')
  out.push(n)
  return out
})

function go(p: number) {
  if (p < 1 || p > props.pageCount || p === props.modelValue) return
  emit('update:modelValue', p)
}
</script>

<template>
  <nav v-if="pageCount > 1" class="flex items-center justify-center gap-1.5" aria-label="分页">
    <button
      type="button"
      class="flex h-7 w-7 cursor-pointer items-center justify-center rounded border border-line text-ink-2 transition-colors hover:border-line-strong disabled:cursor-not-allowed disabled:opacity-40"
      title="上一页"
      aria-label="上一页"
      :disabled="modelValue <= 1"
      @click="go(modelValue - 1)"
    >
      <PhCaretLeft :size="12" />
    </button>
    <template v-for="(p, i) in pages" :key="i">
      <span v-if="p === '…'" class="px-0.5 text-[12px] text-ink-3">…</span>
      <button
        v-else
        type="button"
        class="h-7 min-w-[28px] cursor-pointer rounded border px-1.5 text-[12px] transition-colors"
        :class="p === modelValue ? 'border-accent bg-accent font-semibold text-on-accent' : 'border-line text-ink-2 hover:border-line-strong'"
        :aria-current="p === modelValue ? 'page' : undefined"
        @click="go(p)"
      >
        {{ p }}
      </button>
    </template>
    <button
      type="button"
      class="flex h-7 w-7 cursor-pointer items-center justify-center rounded border border-line text-ink-2 transition-colors hover:border-line-strong disabled:cursor-not-allowed disabled:opacity-40"
      title="下一页"
      aria-label="下一页"
      :disabled="modelValue >= pageCount"
      @click="go(modelValue + 1)"
    >
      <PhCaretRight :size="12" />
    </button>
  </nav>
</template>
