<script setup lang="ts">
import { PhArrowClockwise, PhCheck, PhCaretLeft, PhCaretRight, PhSparkle } from '@phosphor-icons/vue'
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

import Button from '@/components/ui/Button.vue'
import TemplateFilterBar from '@/components/templates/TemplateFilterBar.vue'
import { deckV2Api, type TplSuggestion } from '@/api/deckV2'
import { templateApi, type TemplateVariant } from '@/api/templates'
import { userTemplateApi, type CommunityTemplate, type UserTemplateRow } from '@/api/userTemplates'
import { useTemplateFilter } from '@/lib/templateFilter'
import { resolveSuggestion, variantFor } from '@/lib/templateSuggest'
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
  tags?: string[]
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
      // 用户模板没有 tags 数据：标签筛选只覆盖内置模板（方案已拍板接受）
      tags: ut ? undefined : t.tags,
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

// —— 标签筛选 + 名字搜索：词表从数据算出；过滤作用在喂瀑布流的数组上，
// 被筛掉的卡片 iframe 直接卸载，选中/变体/翻页逻辑零改动。——//
const { query, activeTag, vocab, filtered, toggleTag, reset } = useTemplateFilter(cards)

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

// —— AI 推荐（tplsuggest）：进页自动算一次（LLM 秒级，异步不阻塞选模板），
// 结果后端按 deck 缓存，重进页面/刷新都命中缓存；「重新推荐」才带 refresh 重算。
// 点击小卡 = 预选 + 滚动定位到主网格，大预览/变体/开始生成全走既有链路——
// 推荐只做引导，确认权仍在用户手里。失败/空结果收成一行弱化重试，不占版面。——//
const rootEl = ref<HTMLElement | null>(null)
const suggests = ref<TplSuggestion[]>([])
const suggestState = ref<'idle' | 'loading' | 'ready' | 'empty' | 'failed'>('idle')
let suggestAbort: AbortController | null = null

// 推荐条目 → 画廊卡片（候选外 id / 已删模板静默丢弃，不渲染空壳）
const suggestCards = computed(() =>
  suggests.value
    .map((s) => ({ s, card: resolveSuggestion(cards.value, s) }))
    .filter((x): x is { s: TplSuggestion; card: GalleryCard } => !!x.card),
)

async function loadSuggestions(refresh = false) {
  if (!wizard.deckId) return
  suggestAbort?.abort()
  const ac = new AbortController()
  suggestAbort = ac
  if (!suggests.value.length) suggestState.value = 'loading'
  try {
    const r = await deckV2Api.suggestTemplates(wizard.deckId, refresh)
    if (ac.signal.aborted) return
    suggests.value = r.suggestions ?? []
    suggestState.value = suggests.value.length ? 'ready' : 'empty'
  } catch {
    if (ac.signal.aborted) return
    // 已有推荐在展示时刷新失败：保住旧结果，只用 toast 提示；首次失败才收成重试行
    if (!suggests.value.length) suggestState.value = 'failed'
    else toast.error('重新推荐失败')
  }
}

watch(
  () => wizard.deckId,
  (id, old) => {
    if (id === old) return
    suggests.value = []
    if (id) void loadSuggestions()
    else suggestState.value = 'idle'
  },
  { immediate: true },
)

function pickSuggestion(s: TplSuggestion) {
  const card = resolveSuggestion(cards.value, s)
  if (!card) return
  // 目标卡片可能正被筛掉：先清筛选回到全量，再预选、再定位
  if (query.value || activeTag.value) reset()
  pick(card.id)
  const v = variantFor(card, s)
  if (v) selectedVariant.value = v
  void nextTick(() => {
    rootEl.value
      ?.querySelector(`[data-card-id="${CSS.escape(card.id)}"]`)
      ?.scrollIntoView({ behavior: 'smooth', block: 'center' })
  })
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
  for (const c of filtered.value) {
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
  suggestAbort?.abort()
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
  <div ref="rootEl" class="mx-auto flex h-full w-full max-w-[1200px] flex-col gap-3 overflow-y-auto p-5">
    <div>
      <h2 class="text-[16px] font-semibold text-ink">选择模板</h2>
      <p class="mt-0.5 text-[12px] text-ink-3">
        模板决定整套视觉（配色、字体、版式）。必选一个；选中后可切换主题变体、翻页预览整本 demo。
      </p>
    </div>

    <TemplateFilterBar
      v-model:query="query"
      :active-tag="activeTag"
      :vocab="vocab"
      :total="cards.length"
      :shown="filtered.length"
      class="mt-3"
      @update:active-tag="toggleTag"
    />

    <!-- AI 推荐：首算时的骨架 / 有结果的小卡横排 / 失败或空结果的一行弱化重试。
         小卡点击=预选+定位主网格，不用小卡发起生成——确认仍在底部主链路。 -->
    <div v-if="suggestState === 'loading' && !suggests.length" class="flex gap-3" aria-hidden="true">
      <div v-for="i in 4" :key="i" class="h-[168px] flex-1 animate-pulse rounded-control bg-surface-2" />
    </div>
    <div v-else-if="suggests.length" class="rounded-card border border-line bg-surface p-3">
      <div class="flex items-center gap-1.5">
        <PhSparkle :size="14" class="shrink-0 text-accent" />
        <span class="text-[13px] font-semibold text-ink">AI 推荐</span>
        <span class="hidden truncate text-[11.5px] text-ink-3 sm:inline">按大纲和你的对话挑的，点一张直接预选</span>
        <Button
          class="ml-auto shrink-0"
          size="sm"
          variant="ghost"
          :loading="suggestState === 'loading'"
          @click="loadSuggestions(true)"
        >
          <PhArrowClockwise :size="12" />
          重新推荐
        </Button>
      </div>
      <div class="mt-2 flex gap-3 overflow-x-auto pb-1">
        <button
          v-for="sc in suggestCards"
          :key="sc.s.template_id"
          type="button"
          class="w-[184px] shrink-0 cursor-pointer flex-col overflow-hidden rounded-control border bg-surface text-left transition-all hover:border-accent"
          :class="selected === sc.s.template_id ? 'border-accent shadow-[0_0_0_3px_var(--ring)]' : 'border-line'"
          :title="'选用 ' + sc.card.name"
          @click="pickSuggestion(sc.s)"
        >
          <div
            class="relative w-full overflow-hidden bg-surface-2"
            :style="{ aspectRatio: `${sc.card.canvas.w} / ${sc.card.canvas.h}` }"
          >
            <iframe
              :src="templateApi.previewUrl(sc.s.template_id, '', 1, 1)"
              loading="lazy"
              class="pointer-events-none absolute left-0 top-0 border-0"
              :style="{
                width: `${sc.card.canvas.w}px`,
                height: `${sc.card.canvas.h}px`,
                transform: `scale(${184 / sc.card.canvas.w})`,
                transformOrigin: 'top left',
              }"
              sandbox="allow-scripts allow-same-origin"
              :title="sc.card.name + ' 缩略预览'"
              aria-hidden="true"
            />
            <span
              v-if="selected === sc.s.template_id"
              class="absolute right-1.5 top-1.5 rounded-full bg-accent px-1.5 py-0.5 text-[10px] font-semibold text-accent-contrast"
            >
              已选
            </span>
          </div>
          <div class="p-2">
            <p class="truncate text-[12px] font-semibold text-ink">{{ sc.card.name }}</p>
            <p class="mt-0.5 line-clamp-2 text-[11px] leading-snug text-ink-2">{{ sc.s.reason }}</p>
          </div>
        </button>
      </div>
    </div>
    <div v-else-if="suggestState === 'failed' || suggestState === 'empty'" class="flex items-center gap-2 text-[12px] text-ink-3">
      <span>{{ suggestState === 'failed' ? '推荐生成失败' : '这次没给出合适的推荐' }}</span>
      <button class="cursor-pointer font-semibold text-accent hover:underline" @click="loadSuggestions(true)">
        重试
      </button>
    </div>

    <div v-if="filtered.length" class="mt-3 flex items-start gap-3">
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
              class="absolute right-2 top-2 inline-flex items-center gap-1 rounded-full bg-accent px-2 py-0.5 text-[10.5px] font-semibold text-accent-contrast"
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
            <div v-if="c.tags?.length" class="flex flex-wrap gap-1">
              <span v-for="s in c.tags.slice(0, 3)" :key="s" class="rounded bg-accent-soft px-1.5 py-0.5 text-[10px] text-ink-2">
                {{ s }}
              </span>
            </div>
            <div v-if="c.scenario?.length" class="flex flex-wrap gap-1">
              <span v-for="s in c.scenario.slice(0, 3)" :key="s" class="rounded bg-surface-2 px-1.5 py-0.5 text-[10px] text-ink-3">
                {{ s }}
              </span>
            </div>
          </div>
        </button>
      </div>
    </div>

    <div
      v-else
      class="mt-3 flex flex-col items-center gap-2 rounded-card border border-dashed border-line px-6 py-12 text-center"
    >
      <p class="text-[13px] text-ink-2">没有匹配的模板</p>
      <button class="cursor-pointer text-[12.5px] font-semibold text-accent hover:underline" @click="reset">
        清空筛选条件
      </button>
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
            aria-label="上一页"
            @click="demoPage = Math.max(1, demoPage - 1)"
          >
            <PhCaretLeft :size="11" />
          </button>
          <span class="font-mono text-[11.5px] text-ink-2">第 {{ demoPage }} 页</span>
          <button
            class="cursor-pointer rounded border border-line px-1.5 py-0.5 transition-colors hover:border-line-strong"
            title="下一页"
            aria-label="下一页"
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
