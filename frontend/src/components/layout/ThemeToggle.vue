<script setup lang="ts">
import { PhMoon, PhSun } from '@phosphor-icons/vue'
import { useTheme } from '@/theme/useTheme'

const { theme, toggle } = useTheme()
</script>

<template>
  <button
    class="inline-flex h-8 w-8 cursor-pointer items-center justify-center rounded-control border border-line bg-surface-2 text-ink-1 transition-colors hover:border-line-strong"
    :aria-label="theme === 'dark' ? '切换到亮色主题' : '切换到暗色主题'"
    @click="toggle"
  >
    <PhSun v-if="theme === 'dark'" :size="15" />
    <PhMoon v-else :size="15" />
  </button>
</template>
