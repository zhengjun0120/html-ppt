<script setup lang="ts">
import { PhArrowClockwise, PhCards, PhPlus, PhPresentation } from '@phosphor-icons/vue'
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { ApiError } from '@/api/client'
import type { DeckMeta } from '@/api/decks'
import { thumbUrl } from '@/api/deckV2'
import Button from '@/components/ui/Button.vue'
import Empty from '@/components/ui/Empty.vue'
import Pagination from '@/components/ui/Pagination.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import { useDeckStore } from '@/stores/deck'
import { useToast } from '@/stores/toast'

const deckStore = useDeckStore()
const router = useRouter()
const route = useRoute()
const toast = useToast()

/** 封面加载失败的 deck id（无缩略图/未渲染时回退图标位） */
const coverFailed = ref<Record<string, boolean>>({})

// —— 纯前端分页：/api/decks 本就全量返回，切片即可；缩略图 loading=lazy，
// 每页只挂 12 张封面，首屏不用再等 40+ 个 iframe 式的 img 解析。——//
const PER_PAGE = 12
/** 页码 ↔ URL query：翻页写回 ?page=N，从文稿页返回/浏览器后退/刷新都落回原页 */
function pageFromQuery(): number {
  const n = Number(route.query.page)
  return Number.isInteger(n) && n >= 1 ? n : 1
}
const page = ref(pageFromQuery())
const pageCount = computed(() => Math.max(1, Math.ceil(deckStore.list.length / PER_PAGE)))
const pagedDecks = computed(() => deckStore.list.slice((page.value - 1) * PER_PAGE, page.value * PER_PAGE))
// 删除/刷新后列表变短，当前页可能越界；页码变化统一写回 query（翻页/回钳都走这里）
watch(pageCount, (n) => {
  if (page.value > n) page.value = n
})
watch(page, (p) => {
  void router.replace({ query: p > 1 ? { page: String(p) } : undefined }).catch(() => {})
})
// 列表滚在壳层 <main class="overflow-y-auto"> 里，window 本身不滚——记位置、回顶、
// 恢复都得作用到这个容器上，从视图根向上找最近的滚动祖先拿它。
const rootEl = ref<HTMLElement | null>(null)
function scroller(): HTMLElement | null {
  let el: HTMLElement | null = rootEl.value
  while (el && el !== document.body) {
    const oy = getComputedStyle(el).overflowY
    if (oy === 'auto' || oy === 'scroll') return el
    el = el.parentElement
  }
  return null
}

function setPage(p: number) {
  page.value = p
  const el = scroller()
  if (el) el.scrollTo({ top: 0, behavior: 'smooth' })
}

// —— 返回文稿：点开文稿前记下当前位置——页码进 URL query、滚动高度进 sessionStorage。
// 列表滚在壳层 <main> 里而非 window，savedPosition 管不到它，所以 router.back()
// 与 fallback push 两条返回路径都靠 restoreScroll 在 onMounted 落回。——//
function openDeck(id: string) {
  try {
    sessionStorage.setItem('decks-return-page', String(page.value))
    sessionStorage.setItem('decks-return-y', String(scroller()?.scrollTop ?? 0))
  } catch {
    /* 隐私模式等存不进就算了：只损失滚动恢复，不影响导航 */
  }
  router.push(`/decks/${id}`)
}

/** 只有生成过页面的文稿才有封面：未实例化（澄清/大纲/选模板阶段）的 deck
 * 没有 index.html，发封面请求只会换来后端一次必然失败的渲染（404）。 */
function hasCover(d: DeckMeta): boolean {
  return d.stage === 'generating' || d.stage === 'iterating'
}

/** 滚动落回：仅当 URL ?page 与离开时记录一致才回滚；恢复即清，避免之后的普通进入被旧位置拽走 */
async function restoreScroll() {
  const savedPage = Number(sessionStorage.getItem('decks-return-page') || 0)
  if (!savedPage || savedPage !== page.value) return
  const y = Number(sessionStorage.getItem('decks-return-y') || 0)
  sessionStorage.removeItem('decks-return-page')
  sessionStorage.removeItem('decks-return-y')
  if (y > 0) {
    await nextTick()
    scroller()?.scrollTo({ top: y })
  }
}

onMounted(async () => {
  try {
    await deckStore.refresh()
  } catch (e) {
    toast.error(e instanceof ApiError ? e.message : '文稿列表加载失败')
  }
  await restoreScroll()
})
</script>

<template>
  <div ref="rootEl" class="mx-auto max-w-[1100px] p-6">
    <div class="mb-5 flex items-center justify-between">
      <h1 class="text-[16px] font-bold">我的文稿</h1>
      <div class="flex items-center gap-2">
        <Button :loading="deckStore.loading" @click="deckStore.refresh()">
          <PhArrowClockwise :size="13" />
          刷新
        </Button>
        <Button variant="primary" @click="router.push('/new')">
          <PhPlus :size="13" />
          新文稿
        </Button>
      </div>
    </div>

    <!-- 加载骨架屏：形状对齐最终布局 -->
    <div v-if="deckStore.loading && deckStore.list.length === 0" class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <div v-for="i in 6" :key="i" class="overflow-hidden rounded-card border border-line bg-surface">
        <Skeleton class="h-[130px] rounded-none" />
        <div class="space-y-2 p-3">
          <Skeleton class="h-3.5 w-3/5" />
          <Skeleton class="h-3 w-2/5" />
        </div>
      </div>
    </div>

    <!-- 空态 -->
    <Empty
      v-else-if="deckStore.list.length === 0"
      title="还没有文稿"
      desc="点「新文稿」，在对话里描述你想要的演示文稿，agent 会先对齐大纲、再由你挑模板，逐页生成后自动存到这里。"
    >
      <template #icon><PhPresentation /></template>
    </Empty>

    <!-- 卡片网格：封面用第一页缩略图（首次访问会触发后端渲染，之后走缓存） -->
    <div v-else class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <button
        v-for="d in pagedDecks"
        :key="d.id"
        class="group cursor-pointer overflow-hidden rounded-card border border-line bg-surface text-left transition-[transform,border-color] hover:-translate-y-0.5 hover:border-line-strong"
        @click="openDeck(d.id)"
      >
        <div class="relative flex h-[130px] items-center justify-center overflow-hidden bg-linear-to-br from-surface-3 to-surface-2 text-ink-3 transition-colors group-hover:text-ink-2">
          <img
            v-if="hasCover(d) && !coverFailed[d.id]"
            :src="thumbUrl(d.id, 1)"
            :alt="`${d.title} 封面`"
            class="absolute inset-0 h-full w-full object-cover object-top"
            loading="lazy"
            @error="coverFailed[d.id] = true"
          />
          <PhCards v-else :size="28" />
        </div>
        <div class="p-3">
          <p class="truncate text-[13.5px] font-semibold" :title="d.title">{{ d.title }}</p>
          <p class="mt-0.5 font-mono text-[11px] text-ink-3">{{ d.id }}</p>
        </div>
      </button>
    </div>

    <Pagination
      v-if="deckStore.list.length"
      class="mt-6"
      :model-value="page"
      :page-count="pageCount"
      @update:model-value="setPage"
    />

    <p v-if="deckStore.list.length" class="mt-4 text-center text-[11.5px] text-ink-3">
      共 {{ deckStore.list.length }} 份文稿 · 点开任意文稿继续迭代；封面是第一页实时快照，内容更新后自动刷新。
    </p>
  </div>
</template>
