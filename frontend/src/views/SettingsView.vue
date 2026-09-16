<script setup lang="ts">
import { PhKey, PhSignOut, PhUserCircle } from '@phosphor-icons/vue'
import { onMounted, ref } from 'vue'

import { ApiError } from '@/api/client'
import Badge from '@/components/ui/Badge.vue'
import Button from '@/components/ui/Button.vue'
import InputField from '@/components/ui/InputField.vue'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/stores/toast'

const auth = useAuthStore()
const toast = useToast()

const apiKey = ref('')
const saving = ref(false)
const keyError = ref('')

onMounted(() => {
  void auth.fetchMe().catch(() => {})
})

async function saveKey() {
  keyError.value = ''
  const v = apiKey.value.trim()
  if (v === '') {
    keyError.value = '要清除请用「清除」按钮；保存需要非空 Key'
    return
  }
  saving.value = true
  try {
    await auth.saveApiKey(v)
    apiKey.value = ''
    toast.success('API Key 已保存')
  } catch (e) {
    keyError.value = e instanceof ApiError ? e.message : '保存失败，请稍后重试'
  } finally {
    saving.value = false
  }
}

async function clearKey() {
  saving.value = true
  keyError.value = ''
  try {
    await auth.saveApiKey('')
    apiKey.value = ''
    toast.info('已清除自备 API Key')
  } catch (e) {
    keyError.value = e instanceof ApiError ? e.message : '操作失败，请稍后重试'
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="mx-auto flex max-w-[640px] flex-col gap-4 p-6">
    <h1 class="text-[16px] font-bold">设置</h1>

    <section class="rounded-card border border-line bg-surface p-5">
      <h2 class="mb-4 flex items-center gap-2 text-[13px] font-semibold text-ink-2">
        <PhUserCircle :size="16" />
        账号
      </h2>
      <div class="flex items-center justify-between">
        <div>
          <p class="text-[13.5px]">{{ auth.email || '…' }}</p>
          <p class="mt-0.5 text-[11.5px] text-ink-3">登录凭据 30 天有效，每天上线自动续期</p>
        </div>
        <Button variant="danger" size="sm" @click="auth.logout()">
          <PhSignOut :size="13" />
          退出登录
        </Button>
      </div>
    </section>

    <section class="rounded-card border border-line bg-surface p-5">
      <h2 class="mb-1 flex items-center gap-2 text-[13px] font-semibold text-ink-2">
        <PhKey :size="16" />
        自备 LLM API Key（BYOK）
      </h2>
      <p class="mb-4 text-[12.5px] text-ink-3">
        配置后，对话将使用你自己的模型额度；
        <Badge :tone="auth.hasKey ? 'success' : 'neutral'">{{ auth.hasKey ? '已配置' : '未配置' }}</Badge>
      </p>
      <div class="flex items-end gap-2">
        <div class="flex-1">
          <InputField
            v-model="apiKey"
            label="API Key"
            type="password"
            placeholder="sk-…（保存后不会回显）"
            :error="keyError"
          />
        </div>
        <Button :loading="saving" @click="saveKey">保存</Button>
        <Button variant="danger" :disabled="!auth.hasKey || saving" @click="clearKey">清除</Button>
      </div>
      <p class="mt-2 text-[11.5px] text-ink-3">
        Key 以密文存储于服务端，仅用于替你调用模型；未配置时使用系统免费额度。
      </p>
    </section>
  </div>
</template>
