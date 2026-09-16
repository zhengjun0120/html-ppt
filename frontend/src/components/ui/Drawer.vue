<script setup lang="ts">
import { PhX } from '@phosphor-icons/vue'
import { onBeforeUnmount, onMounted } from 'vue'

const props = defineProps<{ open: boolean; title: string; wide?: boolean }>()
const emit = defineEmits<{ close: [] }>()

function onKey(e: KeyboardEvent) {
  if (e.key === 'Escape' && props.open) emit('close')
}
onMounted(() => window.addEventListener('keydown', onKey))
onBeforeUnmount(() => window.removeEventListener('keydown', onKey))
</script>

<template>
  <Teleport to="body">
    <Transition name="drawer">
      <div
        v-if="open"
        class="fixed inset-0 z-50 bg-black/45"
        @click.self="emit('close')"
      >
        <aside
          class="absolute right-0 top-0 flex h-full flex-col border-l border-line bg-surface shadow-2xl"
          :class="wide ? 'w-full max-w-[560px]' : 'w-full max-w-[420px]'"
          role="dialog"
          :aria-label="title"
        >
          <header class="flex items-center justify-between border-b border-line px-4 py-3">
            <h2 class="text-[14px] font-bold">{{ title }}</h2>
            <button
              class="inline-flex h-7 w-7 cursor-pointer items-center justify-center rounded text-ink-2 hover:bg-surface-2 hover:text-ink"
              aria-label="关闭"
              @click="emit('close')"
            >
              <PhX :size="15" />
            </button>
          </header>
          <div class="min-h-0 flex-1 overflow-y-auto">
            <slot />
          </div>
        </aside>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.drawer-enter-active,
.drawer-leave-active {
  transition: opacity 0.2s ease;
}
.drawer-enter-active aside,
.drawer-leave-active aside {
  transition: transform 0.22s ease;
}
.drawer-enter-from,
.drawer-leave-to {
  opacity: 0;
}
.drawer-enter-from aside,
.drawer-leave-to aside {
  transform: translateX(24px);
}
</style>
