<script setup lang="ts">
import { PhArrowRight } from '@phosphor-icons/vue'
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { ApiError } from '@/api/client'
import Button from '@/components/ui/Button.vue'
import InputField from '@/components/ui/InputField.vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()

const email = ref('')
const password = ref('')
const formError = ref('')
const loading = ref(false)

async function submit() {
  formError.value = ''
  if (!email.value.trim() || !password.value) {
    formError.value = '请输入邮箱和密码'
    return
  }
  loading.value = true
  try {
    await auth.login(email.value.trim(), password.value)
    const next = typeof route.query.next === 'string' ? route.query.next : '/decks'
    void router.push(next)
  } catch (e) {
    formError.value = e instanceof ApiError ? e.message : '网络异常，请稍后重试'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="flex min-h-dvh flex-col items-center justify-center p-6">
    <div class="w-full max-w-[360px]">
      <h1 class="mb-8 text-center text-[17px] font-bold">
        Deck<span class="text-accent">Agent</span>
      </h1>
      <form class="rounded-card border border-line bg-surface p-6" novalidate @submit.prevent="submit">
        <div class="flex flex-col gap-4">
          <InputField v-model="email" label="邮箱" type="email" placeholder="you@example.com" autocomplete="email" />
          <InputField v-model="password" label="密码" type="password" placeholder="••••••••" autocomplete="current-password" />
          <p v-if="formError" class="text-xs text-danger">{{ formError }}</p>
          <Button type="submit" variant="primary" :loading="loading" class="w-full">
            登录
            <PhArrowRight :size="14" />
          </Button>
        </div>
      </form>
      <p class="mt-4 text-center text-[12.5px] text-ink-3">
        还没有账号？
        <RouterLink to="/register" class="text-accent hover:underline">注册</RouterLink>
      </p>
    </div>
  </div>
</template>
