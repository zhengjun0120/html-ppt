<script setup lang="ts">
// 骨架屏：加载态用"最终布局的形状"占位（DESIGN.md §5），shimmer 尊重 reduced-motion
withDefaults(defineProps<{ rounded?: string }>(), { rounded: 'rounded' })
</script>

<template>
  <div class="skel relative overflow-hidden bg-surface-2" :class="rounded" />
</template>

<style scoped>
.skel::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(
    90deg,
    transparent,
    color-mix(in srgb, var(--text-1) 8%, transparent),
    transparent
  );
  animation: sk-shimmer 1.4s infinite;
}
@media (prefers-reduced-motion: reduce) {
  .skel::after {
    animation: none;
  }
}
@keyframes sk-shimmer {
  from {
    transform: translateX(-100%);
  }
  to {
    transform: translateX(100%);
  }
}
</style>
