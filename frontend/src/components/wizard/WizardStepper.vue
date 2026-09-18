<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'

import { STEP_LABELS, STEP_ROUTES, stepOrder, useWizardStore, type WizardStep } from '@/stores/wizard'

const props = defineProps<{ clickable?: boolean }>()
const wizard = useWizardStore()
const router = useRouter()

const current = computed(() => wizard.step)
// 当前步之前的都算完成（生成步进行中不算完成）
function isDone(key: WizardStep): boolean {
  if (!current.value) return false
  return stepOrder(key) < stepOrder(current.value)
}
/** 可点击回退：已完成的步骤且向导有 deck 时可跳回查看 */
function go(key: WizardStep) {
  if (!props.clickable || !wizard.deckId) return
  if (!isDone(key) && key !== current.value) return
  void router.push(STEP_ROUTES[key](wizard.deckId))
}
</script>

<template>
  <ol class="flex items-center gap-1 text-[11.5px]" aria-label="生成流程">
    <template v-for="(s, i) in STEP_LABELS" :key="s.key">
      <li v-if="i > 0" class="h-px w-4 bg-line" aria-hidden="true" />
      <li
        class="flex items-center gap-1.5 rounded-full px-2 py-0.5 transition-colors"
        :class="[
          current === s.key
            ? 'bg-accent-soft font-semibold text-accent'
            : isDone(s.key)
              ? 'text-ink-2'
              : 'text-ink-3',
          clickable && wizard.deckId && (isDone(s.key) || current === s.key) ? 'cursor-pointer hover:text-ink' : '',
        ]"
        @click="go(s.key)"
      >
        <span
          class="inline-flex h-4 w-4 items-center justify-center rounded-full border text-[9.5px] leading-none"
          :class="current === s.key ? 'border-accent bg-accent text-on-accent' : isDone(s.key) ? 'border-line-strong' : 'border-line'"
        >
          {{ isDone(s.key) && current !== s.key ? '✓' : i + 1 }}
        </span>
        {{ s.label }}
      </li>
    </template>
  </ol>
</template>
