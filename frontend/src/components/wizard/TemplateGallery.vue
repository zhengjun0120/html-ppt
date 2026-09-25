<script setup lang="ts">
import { PhCheck, PhCaretLeft, PhCaretRight } from '@phosphor-icons/vue'
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

import Button from '@/components/ui/Button.vue'
import { templateApi, type TemplateVariant } from '@/api/templates'
import { userTemplateApi, type CommunityTemplate, type UserTemplateRow } from '@/api/userTemplates'
import { useChatStore } from '@/stores/chat'
import { useWizardStore } from '@/stores/wizard'
import { useToast } from '@/stores/toast'

/**
 * 模板画廊（gate 2 的主区域）：必选一个模板 + 变体，确认后触发生成。
 * - 每张卡直接渲染模板 demo 第 1 页的缩略 iframe（loading=lazy 懒加载，
 *   按 canvas 设计像素等比缩放），选中后切变体/翻页预览整本；
 *   ——此前的占位图标让画廊没有任何一眼可见的视觉预览。
 * - /api/templates 是"内置+用户层"合并视图：用户模板排在最前并带"我的"徽标，
 *   其余按 id 排序。（此前用户模板是颗小 chip，点了之后预览渲染在屏幕外的
 *   卡片里，看起来就像"没有预览"。）
 */

const wizard = useWizardStore()
const chat = useChatStore()
const toast = useToast()

const userTemplates = ref<UserTemplateRow[]>([])
const community = ref<CommunityTemplate[]>([])

const selected = ref('')
const selectedVariant = ref('')
const demoPage = ref(1)
const starting = ref(false)

void wizard.loadTemplates().catch(() => toast.error('模板清单加载失败'))
void userTemplateApi.list().then((r) => (userTemplates.value = r)).catch(() => {})
void userTemplateApi.community().then((r) => (community.value = r)).catch(() => {})

interface GalleryCard {
  id: string
  name: string
  description: string
  scenario?: string[]
  canvas: { w: number; h: number }
  variants: TemplateVariant[]
  mine?: boolean
  published?: boolean
}

const mineInfo = computed(() => new Map(userTemplates.value.map((u) => [u.id, u])))
const cards = computed<GalleryCard[]>(() => {
  const all: GalleryCard[] = wizard.templates.map((t) => {
    const ut = mineInfo.value.get(t.id)
    return {
      id: t.id,
      name: ut ? ut.name : t.name,
      description: ut ? ut.description : t.description,
      scenario: ut ? undefined : t.scenario,
      canvas: t.canvas,
      variants: t.variants,
      mine: !!ut,
      published: ut?.status === 'published',
    }
  })
  // 我的模板置顶，其余按 id 稳定排序
  return [...all.filter((c) => c.mine), ...all.filter((c) => !c.mine)]
})

const selectedCard = computed(() => cards.value.find((c) => c.id === selected.value))
const canStart = computed(
  () => !!selectedCard.value && !wizard.busy && !starting.value && chat.sessionId != null,
)

function pick(id: string) {
  if (selected.value === id) return
  selected.value = id
  demoPage.value = 1
  selectedVariant.value = selectedCard.value?.variants[0]?.id ?? 'default'
}

// —— 瀑布流布局：卡片按"最短列"分配到 N 个纵向列（N 跟随视口宽度）。——//
// 不用 CSS columns：它按列填充，会把「我的模板」沉进第一列；按列高分配
// 保住从左到右的阅读顺序（我的模板仍在最上面一排）。列高估算 = 预览相对高
//（画布高宽比）+ 文本块常数，零测量、够用。
const columnCount = ref(1)
function updateColumns() {
  const w = window.innerWidth
  columnCount.value = w >= 1280 ? 3 : w >= 768 ? 2 : 1
}
const columns = computed<GalleryCard[][]>(() => {
  const cols: GalleryCard[][] = Array.from({ length: columnCount.value }, () => [])
  const heights = new Array<number>(columnCount.value).fill(0)
  for (const c of cards.value) {
    let i = 0
    for (let k = 1; k < heights.length; k++) {
      if (heights[k] < heights[i]) i = k
    }
    cols[i].push(c)
    heights[i] += c.canvas.h / c.canvas.w + 0.6
  }
  return cols
})

// —— 每张卡的预览盒宽度测量：iframe 按 canvas 设计像素等比缩到盒宽，盒按画布比例出盒
const boxWidths = ref<Record<string, number>>({})
const boxEls = new Map<string, HTMLElement>()
const ro = new ResizeObserver((es) => {
  for (const e of es) {
    const id = (e.target as HTMLElement).dataset.cardId
    if (id) boxWidths.value[id] = e.contentRect.width
  }
})
function setBoxRef(id: string) {
  return (el: unknown) => {
    if (el) {
      const e = el as HTMLElement
      boxEls.set(id, e)
      ro.observe(e)
    } else {
      const old = boxEls.get(id)
      if (old) ro.unobserve(old)
      boxEls.delete(id)
    }
  }
}
onMounted(() => {
  updateColumns()
  window.addEventListener('resize', updateColumns)
})
onBeforeUnmount(() => {
  window.removeEventListener('resize', updateColumns)
  ro.disconnect()
})

function scaleFor(c: GalleryCard): number {
  const w = boxWidths.value[c.id]
  return w && w > 0 ? w / c.canvas.w : 0.29
}

const previewSrc = computed(() =>
  selected.value ? templateApi.previewUrl(selected.value, selectedVariant.value, demoPage.value) : '',
)

async function start() {
  if (!canStart.value || !selectedCard.value) return
  starting.value = true
  try {
    await wizard.selectTemplate(selected.value, selectedVariant.value)
    // 用户模板不在 store 的内置清单里，deck 元数据就绪前预览要用画布
    wizard.canvas = selectedCard.value.canvas
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

    <div class="mt-3 flex items-start gap-3">
      <div v-for="(col, ci) in columns" :key="ci" class="flex min-w-0 flex-1 flex-col gap-3">
        <button
          v-for="c in col"
          :key="c.id"
          type="button"
          class="group flex cursor-pointer flex-col rounded-control border bg-surface text-left transition-all hover:border-accent"
          :class="selected === c.id ? 'border-accent shadow-[0_0_0_3px_var(--ring)]' : 'border-line'"
          @click="pick(c.id)"
        >
          <div
            :ref="setBoxRef(c.id)"
            :data-card-id="c.id"
            class="relative w-full overflow-hidden rounded-t-control bg-surface-2"
            :style="{ aspectRatio: `${c.canvas.w} / ${c.canvas.h}` }"
          >
            <!-- 未选中：demo 第 1 页缩略；选中后：换肤 + 翻页的实时预览 -->
            <iframe
              :src="selected === c.id ? previewSrc : templateApi.previewUrl(c.id, '', 1, 1)"
              :key="selected === c.id ? previewSrc : `thumb-${c.id}`"
              loading="lazy"
              class="pointer-events-none absolute left-0 top-0 border-0"
              :style="{
                width: `${c.canvas.w}px`,
                height: `${c.canvas.h}px`,
                transform: `scale(${scaleFor(c)})`,
                transformOrigin: 'top left',
              }"
              :sandbox="c.mine ? 'allow-scripts' : 'allow-scripts allow-same-origin'"
              :title="c.name + ' 预览'"
              aria-hidden="true"
            />
            <span
              v-if="selected === c.id"
              class="absolute right-2 top-2 inline-flex items-center gap-1 rounded-full bg-accent px-2 py-0.5 text-[10.5px] font-semibold text-on-accent"
            >
              <PhCheck :size="10" /> 已选
            </span>
          </div>
          <div class="flex flex-col gap-1 p-3">
            <div class="flex items-center gap-2">
              <span class="text-[13.5px] font-semibold text-ink">{{ c.name }}</span>
              <span v-if="c.mine" class="rounded bg-surface-2 px-1 text-[10px] text-ink-3">我的模板</span>
              <span v-if="c.mine && c.published" class="rounded bg-surface-2 px-1 text-[10px] text-ink-3">已发布</span>
              <span v-if="!c.mine" class="font-mono text-[10.5px] text-ink-3">{{ c.id }}</span>
            </div>
            <p class="line-clamp-2 text-[11.5px] text-ink-2">{{ c.description }}</p>
            <div v-if="c.scenario?.length" class="flex flex-wrap gap-1">
              <span v-for="s in c.scenario.slice(0, 3)" :key="s" class="rounded bg-surface-2 px-1.5 py-0.5 text-[10px] text-ink-3">
                {{ s }}
              </span>
            </div>
          </div>
        </button>
      </div>
    </div>

    <div class="sticky bottom-0 mt-auto flex flex-wrap items-center gap-3 rounded-control border border-line bg-surface px-4 py-3">
      <template v-if="selectedCard">
        <template v-if="selectedCard.variants.length">
          <span class="text-[12.5px] text-ink-2">主题变体：</span>
          <label
            v-for="v in selectedCard.variants"
            :key="v.id"
            class="flex cursor-pointer items-center gap-1.5 text-[12.5px]"
            :class="selectedVariant === v.id ? 'font-semibold text-ink' : 'text-ink-3'"
          >
            <input v-model="selectedVariant" type="radio" :value="v.id" class="accent-[var(--accent)]" />
            {{ v.name }}
          </label>
          <span class="mx-1 hidden h-4 w-px bg-line sm:inline-block" />
        </template>
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
