<script setup lang="ts">
defineProps<{ title: string; desc?: string }>()
</script>

<template>
  <div class="flex flex-col items-center justify-center gap-2 py-16 text-center">
    <div class="text-ink-3 [&_svg]:h-10 [&_svg]:w-10">
      <slot name="icon" />
    </div>
    <p class="text-[14px] font-semibold">{{ title }}</p>
    <p v-if="desc" class="max-w-[36ch] text-[12.5px] text-ink-3">{{ desc }}</p>
    <div v-if="$slots.action" class="mt-2">
      <slot name="action" />
    </div>
  </div>
</template>
