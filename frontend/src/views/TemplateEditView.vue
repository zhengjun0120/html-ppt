<script setup lang="ts">
import { PhArrowClockwise, PhGlobe, PhPaperPlaneTilt } from '@phosphor-icons/vue'
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { customizeChat, userTemplateApi, userTemplatePreviewUrl, type UserTemplateRow } from '@/api/userTemplates'
import Button from '@/components/ui/Button.vue'
import { useToast } from '@/stores/toast'

/**
 * 定制工作台（plan-v3 B3）：/templates/:id/edit
 * 左 = demo 实时预览（每次对话改完 token 自动刷新）；
 * 右 = 与 agent 的定制对话（改色板/圆角/观感，token 级）；
 * 顶栏 = 名称编辑 + 发布门禁（两关自动跑，结果就地展示）。
 */

const route = useRoute()
const router = useRouter()
const toast = useToast()

const utId = String(route.params.id)
const row = ref<UserTemplateRow | null>(null)
const name = ref('')
const messages = ref<{ role: 'user' | 'agent'; text: string }[]>([])
const input = ref('')
const sending = ref(false)
const publishing = ref(false)
const previewKey = ref(0)
const demoPage = ref(1)

const previewSrc = computed(() => (row.value ? userTemplatePreviewUrl(utId, demoPage.value) + (previewKey.value ? `?v=${previewKey.value}` : '') : ''))
const previewKeyedSrc = computed(() => previewKey.value + ':' + previewSrc.value)

async function load() {
  try {
    row.value = await userTemplateApi.get(utId)
    name.value = row.value.name
    if (row.value.publish_error) {
      messages.value.push({ role: 'agent', text: `上次发布未过门禁：${row.value.publish_error}` })
    }
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '模板加载失败')
    router.push('/my-templates')
  }
}

onMounted(load)

async function send() {
  const text = input.value.trim()
  if (!text || sending.value) return
  input.value = ''
  messages.value.push({ role: 'user', text })
  sending.value = true
  try {
    const { reply } = await customizeChat(utId, text)
    messages.value.push({ role: 'agent', text: reply })
    previewKey.value += 1 // token 已写入，刷新预览
  } catch (e) {
    messages.value.push({ role: 'agent', text: e instanceof Error ? e.message : '定制失败' })
  } finally {
    sending.value = false
  }
}

async function rename() {
  if (!name.value.trim() || name.value === row.value?.name) return
  try {
    await userTemplateApi.updateMeta(utId, name.value.trim())
    toast.info('已保存名称')
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  }
}

async function publish() {
  publishing.value = true
  try {
    await userTemplateApi.publish(utId)
    toast.info('发布成功，已进入社区模板')
    await load()
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '发布门禁未过')
    await load()
  } finally {
    publishing.value = false
  }
}

async function unpublish() {
  try {
    await userTemplateApi.unpublish(utId)
    toast.info('已下架')
    await load()
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '下架失败')
  }
}
</script>

<template>
  <div class="mx-auto max-w-[1200px] p-6">
    <div class="mb-4 flex flex-wrap items-center gap-2">
      <input
        v-model="name"
        class="min-w-0 flex-1 rounded-control border border-line bg-surface px-3 py-1.5 text-[14px] font-semibold text-ink outline-none focus-visible:border-accent"
        placeholder="模板名称"
        @change="rename"
      />
      <Button @click="router.push('/my-templates')">返回</Button>
      <Button v-if="row?.visibility === 'public'" @click="unpublish">
        <PhGlobe :size="12" />
        下架
      </Button>
      <Button variant="primary" :loading="publishing" @click="publish">
        <PhGlobe :size="12" />
        {{ row?.visibility === 'public' ? '重新过门禁' : '过门禁并公开' }}
      </Button>
    </div>
    <p v-if="row?.publish_error" class="mb-3 rounded-control border border-[#FDEBEC] bg-[#FDEBEC]/40 p-2 text-[12px] text-[#9F2F2D]">
      上次门禁未过：{{ row.publish_error }}
    </p>

    <div class="grid grid-cols-1 gap-4 lg:grid-cols-3">
      <!-- 预览 -->
      <div class="lg:col-span-2">
        <div class="overflow-hidden rounded-card border border-line bg-surface">
          <div class="flex items-center justify-between border-b border-line px-3 py-1.5 text-[12px] text-ink-2">
            <span>实时预览（demo）</span>
            <span class="flex items-center gap-1.5">
              <button
                class="cursor-pointer rounded border border-line px-2 py-0.5 transition-colors hover:border-line-strong"
                @click="demoPage = Math.max(1, demoPage - 1)"
              >
                ←
              </button>
              <span class="font-mono text-[11px]">第 {{ demoPage }} 页</span>
              <button
                class="cursor-pointer rounded border border-line px-2 py-0.5 transition-colors hover:border-line-strong"
                @click="demoPage = demoPage + 1"
              >
                →
              </button>
              <button class="cursor-pointer rounded border border-line px-2 py-0.5 transition-colors hover:border-line-strong" title="刷新预览" @click="previewKey += 1">
                <PhArrowClockwise :size="11" />
              </button>
            </span>
          </div>
          <div class="h-[460px] overflow-hidden bg-surface-2">
            <iframe
              v-if="row"
              :key="previewKeyedSrc"
              :src="previewSrc"
              class="h-full w-full border-0"
              sandbox="allow-scripts"
              title="模板定制预览"
            />
          </div>
        </div>
        <p class="mt-2 text-[11.5px] text-ink-3">
          对话里的每次改动都会重写 token 并刷新这里的预览；发布前门禁会再整本量测一遍。
        </p>
      </div>

      <!-- 定制对话 -->
      <div class="flex min-h-[520px] flex-col rounded-card border border-line bg-surface">
        <div class="border-b border-line px-3 py-1.5 text-[12.5px] font-semibold text-ink-2">定制对话</div>
        <div class="min-h-0 flex-1 space-y-3 overflow-y-auto p-3">
          <div v-if="messages.length === 0" class="rounded-control bg-surface-2 p-3 text-[12px] text-ink-3">
            试试：「主色换成暖橙色，整体更圆一点」「描述文案的弱色再淡一些，底色换成暖白」。
          </div>
          <div
            v-for="(m, i) in messages"
            :key="i"
            class="rounded-control px-3 py-2 text-[12.5px] leading-relaxed"
            :class="m.role === 'user' ? 'ml-6 bg-accent-soft text-ink' : 'mr-2 bg-surface-2 text-ink'"
          >
            {{ m.text }}
          </div>
        </div>
        <div class="border-t border-line p-2">
          <div class="flex gap-2">
            <input
              v-model="input"
              class="min-w-0 flex-1 rounded-control border border-line bg-surface-2 px-3 py-1.5 text-[12.5px] outline-none focus-visible:border-accent"
              placeholder="描述想改的视觉（回车发送）"
              @keydown.enter="send"
            />
            <Button variant="primary" :loading="sending" @click="send">
              <PhPaperPlaneTilt :size="12" />
            </Button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
