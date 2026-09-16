<script setup lang="ts">
import { computed } from 'vue'
import Spinner from './Spinner.vue'

const props = withDefaults(
  defineProps<{
    // type 放首位：@babel/parser 7.29.8 对「多个字符串联合属性后跟 type 属性」的
    // 多行泛型有解析回归（错误报在前一行 `?` 处），调整顺序可绕过
    type?: 'button' | 'submit'
    variant?: 'primary' | 'secondary' | 'ghost' | 'danger'
    size?: 'sm' | 'md'
    disabled?: boolean
    loading?: boolean
  }>(),
  { type: 'button', variant: 'secondary', size: 'md', disabled: false, loading: false },
)

const emit = defineEmits<{ click: [e: MouseEvent] }>()

const variantClass = computed(
  () =>
    ({
      primary:
        'bg-accent text-accent-contrast hover:opacity-85 disabled:hover:opacity-100 shadow-none',
      secondary:
        'bg-surface-2 text-ink border border-line hover:border-line-strong dark:shadow-[inset_0_1px_0_rgba(255,255,255,0.05),inset_0_-1px_0_rgba(0,0,0,0.18)]',
      ghost: 'bg-transparent text-ink-2 hover:bg-surface-2 hover:text-ink',
      danger:
        'bg-danger-soft text-danger border border-danger/30 hover:border-danger/55',
    })[props.variant],
)

const sizeClass = computed(() => ({ sm: 'h-7 px-2.5 text-xs', md: 'h-9 px-4 text-[13px]' })[props.size])
</script>

<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    class="inline-flex cursor-pointer items-center justify-center gap-1.5 rounded-control font-semibold transition-[opacity,border-color,background-color,transform] select-none active:translate-y-px active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-55"
    :class="[variantClass, sizeClass]"
    @click="emit('click', $event)"
  >
    <Spinner v-if="loading" :size="13" />
    <slot />
  </button>
</template>
