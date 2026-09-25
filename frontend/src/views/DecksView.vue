<script setup lang="ts">
import { PhArrowClockwise, PhCards, PhPlus, PhPresentation } from '@phosphor-icons/vue'
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

import { ApiError } from '@/api/client'
import { thumbUrl } from '@/api/deckV2'
import Button from '@/components/ui/Button.vue'
import Empty from '@/components/ui/Empty.vue'
import Pagination from '@/components/ui/Pagination.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import { useDeckStore } from '@/stores/deck'
import { useToast } from '@/stores/toast'

const deckStore = useDeckStore()
const router = useRouter()
const toast = useToast()

/** 封面加载失败的 deck id（无缩略图/未渲染时回退图标位） */
const coverFailed = ref<Record<string, boolean>>({})

// —— 纯前端分页：/api/decks 本就全量返回，切片即可；缩略图 loading=lazy，
// 每页只挂 12 张封面，首屏不用再等 40+ 个 iframe 式的 img 解析。——//
const PER_PAGE = 12
const page = ref(1)
const pageCount = computed(() => Math.max(1, Math.ceil(deckStore.list.length / PER_PAGE)))
const pagedDecks = computed(() => deckStore.list.slice((page.value - 1) * PER_PAGE, page.value * PER_PAGE))
// 删除/刷新后列表变短，当前页可能越界
watch(pageCount, (n) => {
  if (page.value > n) page.value = n
})
function setPage(p: number) {
  page.value = p
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(async () => {
  try {
    await deckStore.refresh()
  } catch (e) {
    toast.error(e instanceof ApiError ? e.message : '文稿列表加载失败')
  }
})
</script>

<template>
  <div class="mx-auto max-w-[1100px] p-6">
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
        @click="router.push(`/decks/${d.id}`)"
      >
        <div class="relative flex h-[130px] items-center justify-center overflow-hidden bg-linear-to-br from-surface-3 to-surface-2 text-ink-3 transition-colors group-hover:text-ink-2">
          <img
            v-if="!coverFailed[d.id]"
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
