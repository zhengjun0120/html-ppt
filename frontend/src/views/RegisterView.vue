<script setup lang="ts">
import { onBeforeUnmount, ref } from 'vue'
import { useRouter } from 'vue-router'

import { authApi } from '@/api/auth'
import { ApiError } from '@/api/client'
import Button from '@/components/ui/Button.vue'
import InputField from '@/components/ui/InputField.vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()

const email = ref('')
const code = ref('')
const password = ref('')
const formError = ref('')
const emailError = ref('')
const codeError = ref('')
const passwordError = ref('')
const submitting = ref(false)
const sending = ref(false)
const countdown = ref(0)

let timer: ReturnType<typeof setInterval> | undefined

function validEmail(v: string): boolean {
  return /^[^@\s]+@[^@\s]+\.[^@\s]+$/.test(v)
}

function startCountdown() {
  countdown.value = 60
  timer = setInterval(() => {
    countdown.value -= 1
    if (countdown.value <= 0 && timer) {
      clearInterval(timer)
      timer = undefined
    }
  }, 1000)
}

onBeforeUnmount(() => {
  if (timer) clearInterval(timer)
})

async function sendCode() {
  emailError.value = ''
  if (!validEmail(email.value)) {
    emailError.value = '邮箱格式不正确'
    return
  }
  sending.value = true
  try {
    await authApi.requestCode(email.value)
    startCountdown()
  } catch (e) {
    emailError.value = e instanceof ApiError ? e.message : '发送失败，请稍后重试'
  } finally {
    sending.value = false
  }
}

async function submit() {
  formError.value = emailError.value = codeError.value = passwordError.value = ''
  if (!validEmail(email.value)) emailError.value = '邮箱格式不正确'
  if (code.value.trim() === '') codeError.value = '请输入邮箱验证码'
  if (password.value.length < 8 || password.value.length > 72) passwordError.value = '密码需 8~72 位'
  if (emailError.value || codeError.value || passwordError.value) return

  submitting.value = true
  try {
    await auth.register(email.value, code.value.trim(), password.value)
    void router.push('/decks')
  } catch (e) {
    formError.value = e instanceof ApiError ? e.message : '网络异常，请稍后重试'
  } finally {
    submitting.value = false
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
          <InputField v-model="email" label="邮箱" type="email" placeholder="you@example.com" :error="emailError" autocomplete="email" />
          <div>
            <span class="mb-1.5 block text-[13px] font-medium text-ink-2">验证码</span>
            <div class="flex gap-2">
              <input
                v-model="code"
                inputmode="numeric"
                maxlength="6"
                placeholder="6 位邮箱验证码"
                class="h-9 w-full min-w-0 rounded-control border bg-surface-2 px-3 text-[13.5px] text-ink outline-none transition-[border-color,box-shadow] placeholder:text-ink-3 focus-visible:border-accent focus-visible:shadow-[0_0_0_3px_var(--ring)]"
                :class="codeError ? 'border-danger' : 'border-line'"
              />
              <Button size="md" :disabled="countdown > 0" :loading="sending" class="shrink-0" @click="sendCode">
                {{ countdown > 0 ? `${countdown}s` : '发送验证码' }}
              </Button>
            </div>
            <span v-if="codeError" class="mt-1 block text-xs text-danger">{{ codeError }}</span>
          </div>
          <InputField
            v-model="password"
            label="密码"
            type="password"
            placeholder="8~72 位"
            :error="passwordError"
            autocomplete="new-password"
          />
          <p v-if="formError" class="text-xs text-danger">{{ formError }}</p>
          <Button type="submit" variant="primary" :loading="submitting" class="w-full">注册并登录</Button>
        </div>
      </form>
      <p class="mt-4 text-center text-[12.5px] text-ink-3">
        已有账号？
        <RouterLink to="/login" class="text-accent hover:underline">直接登录</RouterLink>
      </p>
    </div>
  </div>
</template>
