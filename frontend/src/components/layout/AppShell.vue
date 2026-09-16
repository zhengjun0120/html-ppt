<script setup lang="ts">
import ThemeToggle from './ThemeToggle.vue'
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
          <RouterLink to="/trace" class="nav-link">观测台</RouterLink>
          <RouterLink to="/settings" class="nav-link">设置</RouterLink>
        </nav>
        <div class="ml-auto flex items-center gap-2">
          <!-- M1：用户邮箱 + 退出登录 -->
          <slot name="actions" />
          <ThemeToggle />
        </div>
      </div>
    </header>
    <main class="min-h-0 flex-1 overflow-y-auto">
      <RouterView />
    </main>
  </div>
</template>
