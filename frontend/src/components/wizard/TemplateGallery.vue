<script setup lang="ts">
import { PhCheck, PhCaretLeft, PhCaretRight, PhPresentationChart } from '@phosphor-icons/vue'
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

import Button from '@/components/ui/Button.vue'
import { templateApi } from '@/api/templates'
import { useChatStore } from '@/stores/chat'
import { useWizardStore } from '@/stores/wizard'
import { useToast } from '@/stores/toast'

/**
 * 模板画廊（gate 2 的主区域）：必选一个模板 + 变体，确认后触发生成。
 * plan-v3 C2 升级：
 *   - 选中项的 live 预览走 /api/templates/:id/preview?variant=（服务端换肤，
 *     切变体只换 query，即时看效果）；
 *   - demo 支持 #/N 翻页（runtime 深链），预览不再只有封面；
 *   - 画布比例自适应（竖版模板如 xhs-post 也能正确缩放）。
 */

const wizard = useWizardStore()
const chat = useChatStore()
const toast = useToast()

const selected = ref('')
const selectedVariant = ref('')
const demoPage = ref(1)
const starting = ref(false)

void wizard.loadTemplates().catch(() => toast.error('模板清单加载失败'))

const selectedMeta = computed(() => wizard.templates.find((t) => t.id === selected.value))
const canStart = computed(() => !!selected.value && !wizard.busy && !starting.value && chat.sessionId != null)

function pick(id: string) {
  if (selected.value === id) return
  selected.value = id
  demoPage.value = 1
  const tpl = wizard.templates.find((t) => t.id === id)
  selectedVariant.value = tpl?.variants[0]?.id ?? 'default'
  // 重新测量预览盒宽度（不同卡片宽度不同）
  requestAnimationFrame(measure)
}

const previewBox = ref<HTMLElement | null>(null)
const boxW = ref(0)
let ro: ResizeObserver | null = null
function measure() {
  if (previewBox.value) boxW.value = previewBox.value.clientWidth
}
onMounted(() => {
  ro = new ResizeObserver((es) => {
    for (const e of es) boxW.value = e.contentRect.width
  })
  if (previewBox.value) ro.observe(previewBox.value)
})
onBeforeUnmount(() => ro?.disconnect())

/** iframe 按 canvas 设计像素渲染，scale 适配预览盒宽度（竖版/横版通吃） */
const previewScale = computed(() => {
  const w = selectedMeta.value?.canvas.w ?? 1920
  return boxW.value > 0 ? boxW.value / w : 0.29
})
const previewSrc = computed(() =>
  selected.value ? templateApi.previewUrl(selected.value, selectedVariant.value, demoPage.value) : '',
)

async function start() {
  if (!canStart.value || !selectedMeta.value) return
  starting.value = true
  try {
    await wizard.selectTemplate(selected.value, selectedVariant.value)
    await chat.runGeneration(wizard.deckId)
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '选择模板失败')
  } finally {
    starting.value = false
  }
}
</script>

<template>
  <div class="mx-auto flex h-full w-full max-w-[1200px] flex-col gap-3 overflow-y-auto p-5">
    <div>
      <h2 class="text-[16px] font-semibold text-ink">选择模板</h2>
      <p class="mt-0.5 text-[12px] text-ink-3">
        模板决定整套视觉（配色、字体、版式）。必选一个；选中后可切换主题变体、翻页预览整本 demo。
      </p>
    </div>

    <div class="grid grid-cols-1 gap-3 md:grid-cols-2 xl:grid-cols-3">
      <button
        v-for="t in wizard.templates"
        :key="t.id"
        type="button"
        class="group flex cursor-pointer flex-col rounded-control border bg-surface text-left transition-all hover:border-accent"
        :class="selected === t.id ? 'border-accent shadow-[0_0_0_3px_var(--ring)]' : 'border-line'"
        @click="pick(t.id)"
      >
        <div
          ref="previewBox"
          class="relative w-full overflow-hidden rounded-t-control bg-surface-2"
          :style="{ aspectRatio: `${t.canvas.w} / ${t.canvas.h}` }"
        >
          <!-- 选中即实时预览：服务端按变体换肤；iframe 按 canvas 设计像素等比缩放 -->
          <iframe
            v-if="selected === t.id"
            :key="previewSrc"
            :src="previewSrc"
            class="pointer-events-none absolute left-0 top-0 border-0"
            :style="{
              width: `${t.canvas.w}px`,
              height: `${t.canvas.h}px`,
              transform: `scale(${previewScale})`,
              transformOrigin: 'top left',
            }"
            sandbox="allow-scripts"
            title="模板实时预览"
          />
          <div v-else class="flex h-full items-center justify-center text-ink-3">
            <PhPresentationChart :size="28" />
          </div>
          <span
            v-if="selected === t.id"
            class="absolute right-2 top-2 inline-flex items-center gap-1 rounded-full bg-accent px-2 py-0.5 text-[10.5px] font-semibold text-on-accent"
          >
            <PhCheck :size="10" /> 已选
          </span>
        </div>
        <div class="flex flex-col gap-1 p-3">
          <div class="flex items-center gap-2">
            <span class="text-[13.5px] font-semibold text-ink">{{ t.name }}</span>
            <span class="font-mono text-[10.5px] text-ink-3">{{ t.id }}</span>
          </div>
          <p class="line-clamp-2 text-[11.5px] text-ink-2">{{ t.description }}</p>
          <div v-if="t.scenario?.length" class="flex flex-wrap gap-1">
            <span v-for="s in t.scenario.slice(0, 3)" :key="s" class="rounded bg-surface-2 px-1.5 py-0.5 text-[10px] text-ink-3">
              {{ s }}
            </span>
          </div>
        </div>
      </button>
    </div>

    <div class="sticky bottom-0 mt-auto flex flex-wrap items-center gap-3 rounded-control border border-line bg-surface px-4 py-3">
      <template v-if="selectedMeta">
        <span class="text-[12.5px] text-ink-2">主题变体：</span>
        <label
          v-for="v in selectedMeta.variants"
          :key="v.id"
          class="flex cursor-pointer items-center gap-1.5 text-[12.5px]"
          :class="selectedVariant === v.id ? 'font-semibold text-ink' : 'text-ink-3'"
        >
          <input v-model="selectedVariant" type="radio" :value="v.id" class="accent-[var(--accent)]" />
          {{ v.name }}
        </label>
        <span class="mx-1 hidden h-4 w-px bg-line sm:inline-block" />
        <span class="inline-flex items-center gap-1">
          <button
            class="cursor-pointer rounded border border-line px-1.5 py-0.5 transition-colors hover:border-line-strong disabled:cursor-not-allowed disabled:opacity-40"
            :disabled="demoPage <= 1"
            title="上一页"
            @click="demoPage = Math.max(1, demoPage - 1)"
          >
            <PhCaretLeft :size="11" />
          </button>
          <span class="font-mono text-[11.5px] text-ink-2">第 {{ demoPage }} 页</span>
          <button
            class="cursor-pointer rounded border border-line px-1.5 py-0.5 transition-colors hover:border-line-strong"
            title="下一页"
            @click="demoPage = demoPage + 1"
          >
            <PhCaretRight :size="11" />
          </button>
        </span>
      </template>
      <span v-else class="text-[12.5px] text-ink-3">先选中一个模板</span>
      <Button class="ml-auto" variant="primary" :disabled="!canStart" @click="start">
        {{ starting ? '正在启动…' : '用这个模板生成' }}
      </Button>
    </div>
  </div>
</template>
