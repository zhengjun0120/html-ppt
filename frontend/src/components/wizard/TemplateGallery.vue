<script setup lang="ts">
import { PhCheck, PhPresentationChart } from '@phosphor-icons/vue'
import { computed, ref } from 'vue'

import Button from '@/components/ui/Button.vue'
import { useChatStore } from '@/stores/chat'
import { useWizardStore } from '@/stores/wizard'
import { useToast } from '@/stores/toast'

/**
 * 模板画廊（gate 2 的主区域）：必选一个模板 + 变体，确认后触发生成。
 * live 预览 iframe 只在选中项挂载（demo deck 带动画与全尺寸画布，全部常驻太重）。
 */

const wizard = useWizardStore()
const chat = useChatStore()
const toast = useToast()

const selected = ref('')
const selectedVariant = ref('')
const starting = ref(false)

void wizard.loadTemplates().catch(() => toast.error('模板清单加载失败'))

const selectedMeta = computed(() => wizard.templates.find((t) => t.id === selected.value))
const canStart = computed(() => !!selected.value && !wizard.busy && !starting.value && chat.sessionId != null)

function pick(id: string) {
  selected.value = id
  const tpl = wizard.templates.find((t) => t.id === id)
  selectedVariant.value = tpl?.variants[0]?.id ?? 'default'
}

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
  <div class="mx-auto flex h-full w-full max-w-[1100px] flex-col gap-3 overflow-y-auto p-5">
    <div>
      <h2 class="text-[16px] font-semibold text-ink">选择模板</h2>
      <p class="mt-0.5 text-[12px] text-ink-3">
        模板决定整套视觉（配色、字体、版式）。必选一个；同一模板可换主题变体。
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
        <div class="relative aspect-video w-full overflow-hidden rounded-t-control bg-surface-2">
          <!-- 选中即实时预览：demo deck 按 1920×1080 设计，scale 适配卡片宽度 -->
          <iframe
            v-if="selected === t.id"
            :src="`/templates/${t.id}/index.html`"
            class="pointer-events-none h-[1080px] w-[1920px] origin-top-left border-0"
            style="transform: scale(var(--gallery-scale, 0.29))"
            sandbox="allow-scripts"
            loading="lazy"
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

    <div class="sticky bottom-0 mt-auto flex items-center gap-3 rounded-control border border-line bg-surface px-4 py-3">
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
      </template>
      <span v-else class="text-[12.5px] text-ink-3">先选中一个模板</span>
      <Button class="ml-auto" variant="primary" :disabled="!canStart" @click="start">
        {{ starting ? '正在启动…' : '用这个模板生成' }}
      </Button>
    </div>
  </div>
</template>
