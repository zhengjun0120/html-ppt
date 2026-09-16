<script setup lang="ts">
withDefaults(
  defineProps<{
    modelValue: string
    label?: string
    type?: string
    placeholder?: string
    disabled?: boolean
    error?: string
    autocomplete?: string
  }>(),
  { type: 'text', disabled: false },
)

const emit = defineEmits<{ 'update:modelValue': [v: string] }>()
</script>

<template>
  <label class="block">
    <!-- label 在输入框上方（DESIGN.md §5），placeholder 永不充当字段说明 -->
    <span v-if="label" class="mb-1.5 block text-[13px] font-medium text-ink-2">{{ label }}</span>
    <input
      :type="type"
      :value="modelValue"
      :placeholder="placeholder"
      :disabled="disabled"
      :autocomplete="autocomplete"
      class="h-9 w-full rounded-control border bg-surface-2 px-3 text-[13.5px] text-ink outline-none transition-[border-color,box-shadow] placeholder:text-ink-3 focus-visible:border-accent focus-visible:shadow-[0_0_0_3px_var(--ring)] disabled:cursor-not-allowed disabled:opacity-55"
      :class="error ? 'border-danger' : 'border-line'"
      @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    />
    <span v-if="error" class="mt-1 block text-xs text-danger">{{ error }}</span>
  </label>
</template>
