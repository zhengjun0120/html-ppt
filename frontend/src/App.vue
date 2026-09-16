<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

import AppShell from '@/components/layout/AppShell.vue'
import ToastHost from '@/components/ui/ToastHost.vue'

const route = useRoute()
const blank = computed(() => route.meta.public === true)
</script>

<template>
  <AppShell v-if="!blank">
    <RouterView />
  </AppShell>
  <RouterView v-else />
  <ToastHost />
</template>
