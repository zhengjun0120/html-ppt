<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

import Button from './Button.vue'

const props = withDefaults(
  defineProps<{
    open: boolean
    title: string
    desc?: string
    confirmText?: string
    danger?: boolean
  }>(),
  { confirmText: '确认', danger: false, desc: '' },
)
const emit = defineEmits<{ confirm: []; close: [] }>()

const confirmBtn = ref<InstanceType<typeof Button> | null>(null)

function onKey(e: KeyboardEvent) {
  if (e.key === 'Escape' && props.open) emit('close')
}
onMounted(() => window.addEventListener('keydown', onKey))
onBeforeUnmount(() => window.removeEventListener('keydown', onKey))

watch(
  () => props.open,
  async (o) => {
    if (o) {
      await nextTick()
      // 焦点进对话框，ESC/回车才有稳定落点
      confirmBtn.value?.$el?.focus()
    }
  },
)
</script>

<template>
  <Teleport to="body">
    <Transition name="dlg">
      <div
        v-if="open"
        class="fixed inset-0 z-[60] flex items-center justify-center bg-black/45 p-6"
        @click.self="emit('close')"
      >
        <div
          class="w-full max-w-[400px] rounded-card border border-line bg-surface p-5 shadow-2xl"
          role="alertdialog"
          :aria-label="title"
        >
          <h3 class="text-[14.5px] font-bold">{{ title }}</h3>
          <p v-if="desc" class="mt-1.5 text-[12.5px] text-ink-2">{{ desc }}</p>
          <div class="mt-5 flex justify-end gap-2">
            <Button @click="emit('close')">取消</Button>
            <Button ref="confirmBtn" :variant="danger ? 'danger' : 'primary'" @click="emit('confirm')">
              {{ confirmText }}
            </Button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.dlg-enter-active,
.dlg-leave-active {
  transition: opacity 0.18s ease;
}
.dlg-enter-from,
.dlg-leave-to {
  opacity: 0;
}
</style>
