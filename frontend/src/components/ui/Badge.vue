<script setup lang="ts">
withDefaults(
  defineProps<{ tone?: 'success' | 'warning' | 'danger' | 'info' | 'accent' | 'neutral' }>(),
  { tone: 'neutral' },
)

const toneClass: Record<string, string> = {
  success: 'bg-success-soft text-success',
  warning: 'bg-warning-soft text-warning',
  danger: 'bg-danger-soft text-danger',
  info: 'bg-surface-2 text-info',
  accent: 'bg-accent-soft text-accent',
  neutral: 'bg-surface-2 text-ink-2',
}
</script>

<template>
  <span
    class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[11px] font-medium"
    :class="toneClass[tone]"
  >
    <slot />
  </span>
</template>
