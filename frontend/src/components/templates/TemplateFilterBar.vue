<script setup lang="ts">
import { PhCaretDown, PhCaretUp, PhMagnifyingGlass, PhX } from '@phosphor-icons/vue'
import { computed, ref } from 'vue'

import type { TagVocabEntry } from '@/lib/templateFilter'

/**
 * 模板筛选工具条（docs/template-filter-plan.md）：搜索框 + 标签 chip 行。
 * 画廊（向导选模板）与模板页「内置模板」共用；筛选状态在父组件
 * （useTemplateFilter），这里只回传 v-model，词表与计数都由父级算好传入。
 * 词表实测 220+ 个（长尾全是 count=1）：只平铺前 LIMIT 个，其余进「更多」展开，
 * 不然十几行 chip 会把卡片区整个淹没。选中的长尾标签在折叠时也保持可见。
 */

const props = defineProps<{
  query: string
  activeTag: string
  vocab: TagVocabEntry[]
  total: number
  shown: number
}>()

const emit = defineEmits<{
  'update:query': [v: string]
  'update:activeTag': [v: string]
}>()

const LIMIT = 12
const expanded = ref(false)
const visibleVocab = computed(() => {
  const list = expanded.value ? props.vocab : props.vocab.slice(0, LIMIT)
  // 折叠状态下被选中的长尾标签保持可见（不然选中了却看不到哪个亮着）
  if (!expanded.value && props.activeTag && !list.some((v) => v.tag === props.activeTag)) {
    const hit = props.vocab.find((v) => v.tag === props.activeTag)
    if (hit) return [...list, hit]
  }
  return list
})
const overflowCount = computed(() => Math.max(0, props.vocab.length - LIMIT))

const filtering = () => props.query.trim() !== '' || props.activeTag !== ''

function onInput(e: Event) {
  emit('update:query', (e.target as HTMLInputElement).value)
}

function clearAll() {
  emit('update:query', '')
  emit('update:activeTag', '')
}

function chipClass(active: boolean): string {
  return `cursor-pointer rounded-full border px-3 py-1 text-[12.5px] transition-colors ${
    active ? 'border-accent bg-accent-soft text-ink' : 'border-line bg-surface-2 text-ink-2 hover:border-accent'
  }`
}
</script>

<template>
  <div class="flex flex-col gap-2">
    <div class="flex items-center gap-3">
      <div class="relative min-w-0 flex-1">
        <PhMagnifyingGlass
          :size="13"
          class="pointer-events-none absolute left-2.5 top-1/2 -translate-y-1/2 text-ink-3"
        />
        <input
          :value="query"
          type="text"
          placeholder="搜索名称 / 标签 / 场景…"
          class="w-full rounded-control border border-line bg-surface-2 py-1.5 pl-8 text-[12.5px] text-ink outline-none transition-[border-color,box-shadow] placeholder:text-ink-3 focus-visible:border-accent focus-visible:shadow-[0_0_0_3px_var(--ring)]"
          :class="query ? 'pr-8' : 'pr-3'"
          @input="onInput"
        />
        <button
          v-if="query"
          type="button"
          class="absolute right-2 top-1/2 -translate-y-1/2 cursor-pointer text-ink-3 hover:text-ink"
          title="清空搜索"
          @click="emit('update:query', '')"
        >
          <PhX :size="12" />
        </button>
      </div>
      <span v-if="filtering()" class="shrink-0 text-[11.5px] text-ink-3">筛出 {{ shown }} / {{ total }}</span>
      <button
        v-if="filtering()"
        type="button"
        class="shrink-0 cursor-pointer text-[12.5px] font-semibold text-accent hover:underline"
        @click="clearAll"
      >
        清空
      </button>
    </div>

    <div class="flex flex-wrap items-center gap-1.5">
      <button type="button" :class="chipClass(!activeTag)" @click="emit('update:activeTag', '')">全部</button>
      <button
        v-for="v in visibleVocab"
        :key="v.tag"
        type="button"
        :class="chipClass(activeTag === v.tag)"
        @click="emit('update:activeTag', activeTag === v.tag ? '' : v.tag)"
      >
        {{ v.tag }}
        <span class="font-mono text-[10px] opacity-60">{{ v.count }}</span>
      </button>
      <button
        v-if="overflowCount > 0 && !expanded"
        type="button"
        class="inline-flex cursor-pointer items-center gap-0.5 rounded-full border border-dashed border-line px-3 py-1 text-[12.5px] text-ink-3 transition-colors hover:border-accent hover:text-ink-2"
        @click="expanded = true"
      >
        更多 {{ overflowCount }} <PhCaretDown :size="10" />
      </button>
      <button
        v-else-if="overflowCount > 0"
        type="button"
        class="inline-flex cursor-pointer items-center gap-0.5 rounded-full border border-dashed border-line px-3 py-1 text-[12.5px] text-ink-3 transition-colors hover:border-accent hover:text-ink-2"
        @click="expanded = false"
      >
        收起 <PhCaretUp :size="10" />
      </button>
    </div>
  </div>
</template>
