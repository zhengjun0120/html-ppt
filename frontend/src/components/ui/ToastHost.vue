<script setup lang="ts">
import { PhCheckCircle, PhInfo, PhWarningCircle, PhXCircle } from '@phosphor-icons/vue'
import { useToastStore, type ToastKind } from '@/stores/toast'

const store = useToastStore()

const icon: Record<ToastKind, unknown> = {
  success: PhCheckCircle,
  error: PhXCircle,
  warning: PhWarningCircle,
  info: PhInfo,
}
const tone: Record<ToastKind, string> = {
  success: 'bg-success-soft text-success',
  error: 'bg-danger-soft text-danger',
  warning: 'bg-warning-soft text-warning',
  info: 'bg-surface-2 text-info',
}
</script>

<template>
  <div class="pointer-events-none fixed top-14 left-1/2 z-50 flex -translate-x-1/2 flex-col items-center gap-2">
    <TransitionGroup name="toast">
      <div
        v-for="t in store.items"
        :key="t.id"
        class="pointer-events-auto flex cursor-pointer items-center gap-2 rounded-control px-3.5 py-2 text-[13px] font-medium shadow-[0_8px_28px_rgba(0,0,0,0.35)]"
        :class="tone[t.kind]"
        role="status"
        @click="store.dismiss(t.id)"
      >
        <component :is="icon[t.kind]" :size="15" />
        {{ t.text }}
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition:
    opacity 0.2s ease,
    transform 0.2s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
</style>
