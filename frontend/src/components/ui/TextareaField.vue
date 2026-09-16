<script setup lang="ts">
withDefaults(
  defineProps<{
    modelValue: string
    label?: string
    placeholder?: string
    rows?: number
    disabled?: boolean
  }>(),
  { rows: 3, disabled: false },
)

const emit = defineEmits<{ 'update:modelValue': [v: string] }>()
</script>

<template>
  <label class="block">
    <span v-if="label" class="mb-1.5 block text-[13px] font-medium text-ink-2">{{ label }}</span>
    <textarea
      :value="modelValue"
      :placeholder="placeholder"
      :rows="rows"
      :disabled="disabled"
      class="w-full resize-none rounded-control border border-line bg-surface-2 px-3 py-2 text-[13.5px] text-ink outline-none transition-[border-color,box-shadow] placeholder:text-ink-3 focus-visible:border-accent focus-visible:shadow-[0_0_0_3px_var(--ring)] disabled:cursor-not-allowed disabled:opacity-55"
      @input="emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
    />
  </label>
</template>
