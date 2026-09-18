<script setup lang="ts">
import { PhSignOut } from '@phosphor-icons/vue'
import { onMounted } from 'vue'

import ThemeToggle from './ThemeToggle.vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()

onMounted(() => {
  // 外壳只挂在受保护布局上；有 token 才拉用户摘要（失败静默，401 由 client 全局处理）
  if (localStorage.getItem('da-token')) void auth.fetchMe().catch(() => {})
})
</script>

<template>
  <div class="flex h-dvh flex-col">
    <header class="sticky top-0 z-40 border-b border-line bg-surface/85 backdrop-blur">
      <div class="flex h-12 items-center gap-5 px-5">
        <RouterLink to="/decks" class="text-[15px] font-bold tracking-tight">
          Deck<span class="text-accent">Agent</span>
        </RouterLink>
        <nav class="flex items-center gap-1">
          <RouterLink to="/decks" class="nav-link">文稿</RouterLink>
          <RouterLink to="/my-templates" class="nav-link">我的模板</RouterLink>
          <RouterLink to="/trace" class="nav-link">观测台</RouterLink>
          <RouterLink to="/settings" class="nav-link">设置</RouterLink>
        </nav>
        <div class="ml-auto flex items-center gap-2">
          <span v-if="auth.email" class="hidden text-[12.5px] text-ink-2 md:inline">{{ auth.email }}</span>
          <button
            class="inline-flex h-8 w-8 cursor-pointer items-center justify-center rounded-control text-ink-2 transition-colors hover:bg-surface-2 hover:text-ink"
            aria-label="退出登录"
            @click="auth.logout()"
          >
            <PhSignOut :size="15" />
          </button>
          <ThemeToggle />
        </div>
      </div>
    </header>
    <main class="min-h-0 flex-1 overflow-y-auto">
      <RouterView />
    </main>
  </div>
</template>
