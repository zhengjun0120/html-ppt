<script setup lang="ts">
import { PhCaretDoubleLeft, PhCaretLeft, PhCaretRight, PhCaretDoubleRight } from '@phosphor-icons/vue'
import { computed } from 'vue'

/**
 * 分页控件（配纯前端切片用）：给 page-count，出
 * 「« ‹ 1 … 4 5 6 7 8 … 20 › »」。
 * 当前页两侧各留 2 页、首尾页常驻；省略号可点，向前/向后快翻 5 页；
 * 两端双箭头直达首页/末页——页数多的时候不用一页一页点。只有一页时不渲染。
 */
const props = defineProps<{
  modelValue: number
  pageCount: number
}>()
const emit = defineEmits<{ 'update:modelValue': [page: number] }>()

/** 页码窗口：总页数 ≤9 全量铺开；否则首尾常驻 + 当前页 ±2，中间折叠成可点的省略号 */
const pages = computed<(number | 'skip-l' | 'skip-r')[]>(() => {
  const n = props.pageCount
  if (n <= 9) return Array.from({ length: n }, (_, i) => i + 1)
  const cur = props.modelValue
  const lo = Math.max(2, cur - 2)
  const hi = Math.min(n - 1, cur + 2)
  const out: (number | 'skip-l' | 'skip-r')[] = [1]
  if (lo > 2) out.push('skip-l')
  for (let i = lo; i <= hi; i++) out.push(i)
  if (hi < n - 1) out.push('skip-r')
  out.push(n)
  return out
})

function go(p: number) {
  if (p < 1 || p > props.pageCount || p === props.modelValue) return
  emit('update:modelValue', p)
}
function skip(dir: 'l' | 'r') {
  go(dir === 'l' ? Math.max(1, props.modelValue - 5) : Math.min(props.pageCount, props.modelValue + 5))
}
</script>

<template>
  <nav v-if="pageCount > 1" class="flex flex-wrap items-center justify-center gap-1.5" aria-label="分页">
    <button
      type="button"
      class="flex h-7 w-7 cursor-pointer items-center justify-center rounded border border-line text-ink-2 transition-colors hover:border-line-strong disabled:cursor-not-allowed disabled:opacity-40"
      title="第一页"
      aria-label="第一页"
      :disabled="modelValue <= 1"
      @click="go(1)"
    >
      <PhCaretDoubleLeft :size="12" />
    </button>
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
      <button
        v-if="p !== 'skip-l' && p !== 'skip-r'"
        type="button"
        class="h-7 min-w-[28px] cursor-pointer rounded border px-1.5 text-[12px] transition-colors"
        :class="p === modelValue ? 'border-accent bg-accent font-semibold text-accent-contrast' : 'border-line text-ink-2 hover:border-line-strong'"
        :aria-current="p === modelValue ? 'page' : undefined"
        @click="go(p)"
      >
        {{ p }}
      </button>
      <button
        v-else
        type="button"
        class="cursor-pointer px-1 text-[12px] text-ink-3 transition-colors hover:text-ink"
        :title="p === 'skip-l' ? '向前快翻 5 页' : '向后快翻 5 页'"
        :aria-label="p === 'skip-l' ? '向前快翻 5 页' : '向后快翻 5 页'"
        @click="skip(p === 'skip-l' ? 'l' : 'r')"
      >
        …
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
    <button
      type="button"
      class="flex h-7 w-7 cursor-pointer items-center justify-center rounded border border-line text-ink-2 transition-colors hover:border-line-strong disabled:cursor-not-allowed disabled:opacity-40"
      title="最后一页"
      aria-label="最后一页"
      :disabled="modelValue >= pageCount"
      @click="go(pageCount)"
    >
      <PhCaretDoubleRight :size="12" />
    </button>
  </nav>
</template>
