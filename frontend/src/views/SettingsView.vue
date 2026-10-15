<script setup lang="ts">
import { PhCube, PhKey, PhSignOut, PhUserCircle } from '@phosphor-icons/vue'
import { onMounted, ref } from 'vue'

import { ApiError } from '@/api/client'
import { modelApi, type UserModelItem, type UserModelListResp } from '@/api/models'
import Badge from '@/components/ui/Badge.vue'
import Button from '@/components/ui/Button.vue'
import Dialog from '@/components/ui/Dialog.vue'
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
  void loadModels()
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

// —— 模型管理 ——

const models = ref<UserModelListResp | null>(null)
const modelsLoading = ref(true)
const formOpen = ref(false)
const editId = ref<number | null>(null) // null = 新增
const form = ref({ name: '', model_id: '', base_url: '', api_key: '' })
const formErr = ref('')
const formSaving = ref(false)
const formTesting = ref(false)
const busyId = ref<number | null>(null)
const removeTarget = ref<UserModelItem | null>(null)

async function loadModels() {
  modelsLoading.value = true
  try {
    models.value = await modelApi.list()
  } catch (e) {
    toast.error(e instanceof ApiError ? e.message : '模型列表加载失败')
  } finally {
    modelsLoading.value = false
  }
}

function openAdd() {
  editId.value = null
  form.value = { name: '', model_id: '', base_url: '', api_key: '' }
  formErr.value = ''
  formOpen.value = true
}

function openEdit(m: UserModelItem) {
  editId.value = m.id
  // key 不回填：留空 = 保留原 key
  form.value = { name: m.name, model_id: m.model_id, base_url: m.base_url, api_key: '' }
  formErr.value = ''
  formOpen.value = true
}

function validateForm(): string {
  if (!form.value.model_id.trim()) return '模型 ID 不能为空'
  const u = form.value.base_url.trim()
  if (!/^https?:\/\//.test(u)) return '基础地址必须以 http:// 或 https:// 开头'
  if (!editId.value && !form.value.api_key.trim()) return '' // key 可空（本地网关），不拦
  return ''
}

async function saveModel() {
  formErr.value = validateForm()
  if (formErr.value) return
  formSaving.value = true
  const payload = {
    name: form.value.name.trim(),
    model_id: form.value.model_id.trim(),
    base_url: form.value.base_url.trim(),
    api_key: form.value.api_key.trim(),
  }
  try {
    if (editId.value != null) {
      await modelApi.update(editId.value, payload)
      toast.success('模型已更新')
    } else {
      await modelApi.create(payload)
      toast.success('模型已添加')
    }
    formOpen.value = false
    void loadModels()
  } catch (e) {
    formErr.value = e instanceof ApiError ? e.message : '保存失败，请稍后重试'
  } finally {
    formSaving.value = false
  }
}

async function testForm() {
  formErr.value = ''
  if (!form.value.model_id.trim() || !/^https?:\/\//.test(form.value.base_url.trim())) {
    formErr.value = '测试前先填好模型 ID 和基础地址'
    return
  }
  formTesting.value = true
  try {
    // 编辑且没重填 key 时按库存配置测；其余按表单现值测
    await modelApi.test({
      id: editId.value != null && !form.value.api_key.trim() ? editId.value : undefined,
      base_url: form.value.base_url.trim(),
      model_id: form.value.model_id.trim(),
      api_key: form.value.api_key.trim() || undefined,
    })
    toast.success('连通成功，模型可用')
  } catch (e) {
    toast.error(e instanceof ApiError ? `连通失败：${e.message}` : '连通失败，请稍后重试')
  } finally {
    formTesting.value = false
  }
}

async function setActive(m: UserModelItem | null) {
  const id = m?.id ?? 0
  busyId.value = id
  try {
    await modelApi.setActive(id)
    toast.success(m ? `已切换到「${m.name}」` : '已切回平台模型')
    void loadModels()
  } catch (e) {
    toast.error(e instanceof ApiError ? e.message : '切换失败，请稍后重试')
  } finally {
    busyId.value = null
  }
}

async function testStored(m: UserModelItem) {
  busyId.value = m.id
  try {
    await modelApi.test({ id: m.id })
    toast.success(`「${m.name}」连通正常`)
  } catch (e) {
    toast.error(e instanceof ApiError ? `连通失败：${e.message}` : '连通失败，请稍后重试')
  } finally {
    busyId.value = null
  }
}

async function removeModel() {
  if (!removeTarget.value) return
  const target = removeTarget.value
  removeTarget.value = null
  busyId.value = target.id
  try {
    await modelApi.remove(target.id)
    toast.info(`已删除「${target.name}」`)
    void loadModels()
  } catch (e) {
    toast.error(e instanceof ApiError ? e.message : '删除失败，请稍后重试')
  } finally {
    busyId.value = null
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
      <h2 class="mb-4 flex items-center gap-2 text-[13px] font-semibold text-ink-2">
        <PhCube :size="16" />
        模型
      </h2>
      <p class="mb-3 text-[12.5px] text-ink-3">
        当前使用：
        <Badge :tone="models && models.active_id ? 'success' : 'neutral'">
          {{ modelsLoading
            ? '…'
            : models && models.active_id
              ? (models.models.find((m) => m.id === models!.active_id)?.name ?? '自有模型')
              : `平台模型（${models?.platform_model ?? '…'}）` }}
        </Badge>
      </p>

      <div class="mb-3 flex items-center justify-between">
        <p class="text-[12.5px] text-ink-3">
          添加自己的 OpenAI 兼容模型（官方 API、聚合网关、本地服务都行）；
          选中后对话与生成会改用它，费用记你自己的额度，用量页可按模型筛选。
        </p>
        <Button v-if="!formOpen" size="sm" @click="openAdd">添加模型</Button>
      </div>

      <!-- 新增/编辑表单（行内，不做弹窗——字段少且要测试按钮） -->
      <div v-if="formOpen" class="mb-4 flex flex-col gap-2 rounded-control border border-line bg-surface-2 p-3">
        <div class="grid grid-cols-1 gap-2 sm:grid-cols-2">
          <InputField v-model="form.model_id" label="模型 ID" placeholder="例如 gpt-4o、deepseek-chat" />
          <InputField v-model="form.name" label="名称（可选，默认用模型 ID）" placeholder="自定义显示名" />
        </div>
        <InputField v-model="form.base_url" label="基础地址" placeholder="例如 https://api.openai.com/v1" />
        <InputField
          v-model="form.api_key"
          label="API Key"
          type="password"
          :placeholder="editId != null ? '留空 = 保留原 Key' : 'sk-…（本地服务可不填，密文存储不回显）'"
        />
        <p v-if="formErr" class="text-xs text-danger">{{ formErr }}</p>
        <div class="flex items-center gap-2">
          <Button :loading="formSaving" @click="saveModel">{{ editId != null ? '保存修改' : '添加' }}</Button>
          <Button variant="ghost" :loading="formTesting" @click="testForm">测试连通</Button>
          <Button variant="ghost" :disabled="formSaving || formTesting" @click="formOpen = false">取消</Button>
        </div>
      </div>

      <div v-if="modelsLoading" class="py-2 text-[12.5px] text-ink-3">加载中…</div>
      <div v-else-if="models && models.models.length" class="flex flex-col gap-2">
        <div
          v-for="m in models.models"
          :key="m.id"
          class="flex items-center gap-2 rounded-control border border-line px-3 py-2"
        >
          <div class="min-w-0 flex-1">
            <p class="flex items-center gap-1.5 text-[13px]">
              <span class="truncate">{{ m.name }}</span>
              <Badge v-if="m.id === models?.active_id" tone="success">使用中</Badge>
            </p>
            <p class="mt-0.5 truncate text-[11.5px] text-ink-3">
              {{ m.model_id }} · {{ m.base_url }}<template v-if="!m.has_key"> · 无 Key</template>
            </p>
          </div>
          <div class="flex shrink-0 items-center gap-1">
            <Button
              v-if="m.id !== models?.active_id"
              size="sm"
              variant="ghost"
              :disabled="busyId !== null"
              @click="setActive(m)"
            >设为当前</Button>
            <Button size="sm" variant="ghost" :disabled="busyId !== null" @click="testStored(m)">测试</Button>
            <Button size="sm" variant="ghost" :disabled="busyId !== null" @click="openEdit(m)">编辑</Button>
            <Button size="sm" variant="ghost" :disabled="busyId !== null" @click="removeTarget = m">删除</Button>
          </div>
        </div>
        <p class="text-[11.5px] text-ink-3">
          想切回平台模型：
          <button
            class="cursor-pointer text-accent underline-offset-2 hover:underline"
            :disabled="busyId !== null"
            @click="setActive(null)"
          >点这里</button>。
          看图审查要求所选模型支持图片输入，不支持时自动降级为只看版面数字。
        </p>
      </div>
      <div v-else class="text-[12.5px] text-ink-3">还没有添加模型，点上方「添加模型」。</div>

      <!-- 删除确认：danger 态，替代原生 confirm（MyTemplatesView 同款） -->
      <Dialog
        :open="!!removeTarget"
        title="删除模型"
        :desc="`删除「${removeTarget?.name ?? ''}」后，正在使用它会自动切回平台模型。`"
        confirm-text="删除"
        danger
        @confirm="removeModel"
        @close="removeTarget = null"
      />
    </section>

    <section class="rounded-card border border-line bg-surface p-5">
      <h2 class="mb-1 flex items-center gap-2 text-[13px] font-semibold text-ink-2">
        <PhKey :size="16" />
        平台模型 API Key（BYOK）
      </h2>
      <p class="mb-4 text-[12.5px] text-ink-3">
        给平台模型配自己的 Key；使用上方自选模型时本项不参与计费。
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
