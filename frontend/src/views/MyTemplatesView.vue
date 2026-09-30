<script setup lang="ts">
import { PhArrowClockwise, PhArrowsOutSimple, PhGlobe, PhGlobeHemisphereWest, PhPencilSimple, PhPlus, PhSpinner, PhTrash } from '@phosphor-icons/vue'
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import { templateApi, type TemplateMeta } from '@/api/templates'
import { userTemplateApi, userTemplatePreviewUrl, type CommunityTemplate, type UserTemplateRow } from '@/api/userTemplates'
import TemplateFilterBar from '@/components/templates/TemplateFilterBar.vue'
import TemplatePreviewModal from '@/components/templates/TemplatePreviewModal.vue'
import DeckEditModal from '@/components/editor/DeckEditModal.vue'
import Button from '@/components/ui/Button.vue'
import Dialog from '@/components/ui/Dialog.vue'
import Empty from '@/components/ui/Empty.vue'
import Pagination from '@/components/ui/Pagination.vue'
import { useTemplateFilter } from '@/lib/templateFilter'
import { useToast } from '@/stores/toast'
import { watch } from 'vue'

/**
 * 我的模板（plan-v3 B3）：克隆/定制/发布/下架/删除的管理页。
 * 发布 = 挂载校验秒级完成（2026-09-28 门禁降级），视觉质量用工作台的「质量体检」随时量测。
 * 「从内置模板派生」直接展示内置模板卡片（真实缩略 + 中文名 + 适用场景标签）——
 * 派生是挑"视觉起点"，看不见起点就没法挑。从空白新建走独立命名弹窗。
 */

const router = useRouter()
const toast = useToast()
const mine = ref<UserTemplateRow[]>([])
const community = ref<CommunityTemplate[]>([])
const templates = ref<TemplateMeta[]>([])
const loading = ref(true)
const busyId = ref('')
const publishReport = ref('')

const STATUS_LABELS: Record<string, string> = {
  draft: '草稿', publishing: '门禁运行中', published: '已公开', failed: '门禁未过',
}

const metaOf = computed(() => new Map(templates.value.map((t) => [t.id, t])))
/** 可派生的内置模板（注册表合并视图里的用户模板不算"内置"） */
const builtins = computed(() => templates.value.filter((t) => !t.id.startsWith('ut-')))
// —— 派生选起点也要筛/搜：百来个内置模板靠翻页找太费劲（docs/template-filter-plan.md）——//
const {
  query: builtinQuery,
  activeTag: builtinTag,
  vocab: builtinVocab,
  filtered: filteredBuiltins,
  toggleTag: toggleBuiltinTag,
  reset: resetBuiltinFilter,
} = useTemplateFilter(builtins)
/** base_id → 中文名（内置模板都有中文名；查不到兜底原 id） */
function baseName(id: string): string {
  if (id === '_blank') return '空白起点'
  return metaOf.value.get(id)?.name ?? id
}

async function load() {
  loading.value = true
  try {
    const [m, c, t] = await Promise.all([userTemplateApi.list(), userTemplateApi.community(), templateApi.list()])
    mine.value = m
    community.value = c
    templates.value = t
    // 刷新后列表变短，当前页可能越界
    const clamp = (p: { value: number }, n: number) => {
      if (p.value > n) p.value = n
    }
    clamp(minePage, Math.max(1, Math.ceil(m.length / MINE_PER_PAGE)))
    clamp(builtinPage, Math.max(1, Math.ceil(filteredBuiltins.value.length / BUILTIN_PER_PAGE)))
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '模板加载失败')
  } finally {
    loading.value = false
  }
}

// —— 两个卡片区各自纯前端分页（列表本就全量拉回，切片即可）：
// 9/页对齐三列瀑布流的 3 行；切页后滚回区块顶部，不然新页卡片在视口外。——//
const MINE_PER_PAGE = 9
const BUILTIN_PER_PAGE = 9
const minePage = ref(1)
const builtinPage = ref(1)
const minePageCount = computed(() => Math.max(1, Math.ceil(mine.value.length / MINE_PER_PAGE)))
const builtinPageCount = computed(() => Math.max(1, Math.ceil(filteredBuiltins.value.length / BUILTIN_PER_PAGE)))
const pagedMine = computed(() => mine.value.slice((minePage.value - 1) * MINE_PER_PAGE, minePage.value * MINE_PER_PAGE))
const pagedBuiltins = computed(() => filteredBuiltins.value.slice((builtinPage.value - 1) * BUILTIN_PER_PAGE, builtinPage.value * BUILTIN_PER_PAGE))
const mineSection = ref<HTMLElement>()
const builtinSection = ref<HTMLElement>()
function setMinePage(p: number) {
  minePage.value = p
  nextTick(() => mineSection.value?.scrollIntoView({ behavior: 'smooth', block: 'start' }))
}
function setBuiltinPage(p: number) {
  builtinPage.value = p
  nextTick(() => builtinSection.value?.scrollIntoView({ behavior: 'smooth', block: 'start' }))
}

// 筛选/搜索条件变化后过滤列表变短，停在原页可能直接越界——统一回第 1 页
watch([builtinQuery, builtinTag], () => {
  builtinPage.value = 1
})

async function fork(baseId: string) {
  busyId.value = 'fork:' + baseId
  try {
    const row = await userTemplateApi.fork(baseId)
    toast.info(`已从「${baseName(row.base_id)}」派生，去定制工作台改出你的风格`)
    await router.push(`/my-templates/${row.id}/edit`)
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '派生失败')
  } finally {
    busyId.value = ''
  }
}

// 从空白新建：中性灰阶脚手架（两个最小说明版式），结构与视觉在工作台里从零长出来。
// 命名走独立弹窗（原生 prompt 在部分环境抓不到焦点且样式割裂）；留空用默认名。
const blankOpen = ref(false)
const blankName = ref('')
const blankCreating = ref(false)
const blankInput = ref<HTMLInputElement | null>(null)

function openBlankDialog() {
  blankName.value = ''
  blankOpen.value = true
  // Dialog 自己会把焦点钉在确认按钮上（另一个 watcher），宏任务里抢回来给输入框
  setTimeout(() => blankInput.value?.focus(), 0)
}

async function createBlank() {
  if (blankCreating.value) return
  blankCreating.value = true
  try {
    const row = await userTemplateApi.blank(blankName.value.trim())
    blankOpen.value = false
    toast.info('空白模板已创建，去工作台把它长成你的样子')
    await router.push(`/my-templates/${row.id}/edit`)
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '创建失败')
  } finally {
    blankCreating.value = false
  }
}

watch(blankOpen, (o) => {
  if (!o) blankName.value = ''
})

async function togglePublish(row: UserTemplateRow) {
  busyId.value = row.id + ':pub'
  publishReport.value = ''
  try {
    if (row.visibility === 'public') {
      await userTemplateApi.unpublish(row.id)
      toast.info('已下架（已生成的文稿不受影响）')
    } else {
      const report = await userTemplateApi.publish(row.id)
      toast.info(`发布成功：${report.render?.pages ?? 0} 页 demo 全过量测`)
    }
    await load()
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '发布失败')
    publishReport.value = e instanceof Error ? e.message : ''
    await load()
  } finally {
    busyId.value = ''
  }
}

async function remove(row: UserTemplateRow) {
  if (!window.confirm(`删除模板「${row.name}」？已生成的文稿不受影响。`)) return
  busyId.value = row.id + ':del'
  try {
    await userTemplateApi.remove(row.id)
    toast.info('已删除')
    await load()
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '删除失败')
  } finally {
    busyId.value = ''
  }
}

// —— 预览缩放的按卡测量：画布设计像素等比缩到盒宽（预览盒按画布比例出盒）。
// 布局是 CSS 多列瀑布流（columns + break-inside-avoid）：卡片天然不等高、
// 按列紧密排布，竖版模板不会拉伸同行卡片，也不需要信箱式裁边。
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
onBeforeUnmount(() => ro.disconnect())
function scaleFor(id: string, canvas: { w: number; h: number }): number {
  const w = boxWidths.value[id]
  return w && w > 0 ? w / canvas.w : 0.29
}

onMounted(load)

// —— 预览弹窗：放大查看整本 demo，逐页翻看。——//
// 页数：内置模板取 meta.demo_pages；用户模板取基模板的（fork 时 index.html
// 原样拷贝、定制只改 style.css，页数必然一致；基模板查不到就未知，翻页由
// runtime.js 的 go() 钳在最后一页兜底）。沙箱：内置模板内容是仓库静态文件，
// 放开 same-origin 共享缓存；用户模板是用户内容，保持 scripts-only。
interface PreviewTarget {
  kind: 'builtin' | 'user'
  id: string
  name: string
  canvas: { w: number; h: number }
  pages?: number
}
const preview = ref<PreviewTarget | null>(null)
const previewSandbox = computed(() =>
  preview.value?.kind === 'user' ? 'allow-scripts' : 'allow-scripts allow-same-origin',
)
/** 预览模式 src：?preview=1 激活 runtime 协议，之后翻页走 preview-goto
 * postMessage（无刷新、带过渡）。不再需要 #/N 深链。 */
const previewSrc = computed(() => {
  const t = preview.value
  if (!t) return ''
  return t.kind === 'builtin'
    ? `/api/templates/${t.id}/preview?preview=1`
    : `/user-templates/${t.id}/index.html?preview=1`
})
function openBuiltinPreview(b: TemplateMeta) {
  preview.value = { kind: 'builtin', id: b.id, name: b.name, canvas: b.canvas, pages: b.demo_pages }
}
function openUserPreview(row: UserTemplateRow) {
  const base = row.base_id ? metaOf.value.get(row.base_id) : undefined
  preview.value = {
    kind: 'user',
    id: row.id,
    name: row.name,
    canvas: row.canvas ?? { w: 1920, h: 1080 },
    pages: base?.demo_pages,
  }
}

// —— 编辑弹窗（deck-editor-plan §4.6）：仅用户模板可编辑，内置卡只读。——//
const editing = ref<UserTemplateRow | null>(null)
</script>

<template>
  <div class="mx-auto max-w-[1100px] p-6">
    <div class="mb-5 flex items-center justify-between">
      <div>
        <h1 class="text-[16px] font-bold">我的模板</h1>
        <p class="mt-0.5 text-[12px] text-ink-3">从内置模板派生，对话里定制视觉；过自动门禁后公开给所有人用。</p>
      </div>
      <Button :loading="loading" @click="load">
        <PhArrowClockwise :size="13" />
        刷新
      </Button>
    </div>

    <!-- 我的模板 -->
    <div v-if="loading" class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <div v-for="i in 3" :key="i" class="overflow-hidden rounded-card border border-line bg-surface">
        <div class="flex h-[130px] items-center justify-center"><PhSpinner :size="20" class="animate-spin text-ink-3" /></div>
      </div>
    </div>
    <Empty
      v-else-if="mine.length === 0"
      title="还没有自己的模板"
      desc="在下方「从内置模板派生」里挑一个起点，克隆后在定制工作台里和 agent 一起改出你的风格。"
    >
      <template #icon><PhGlobeHemisphereWest /></template>
    </Empty>
    <div v-else ref="mineSection" class="scroll-mt-4">
      <div class="columns-1 gap-4 sm:columns-2 lg:columns-3">
        <div v-for="row in pagedMine" :key="row.id" class="mb-4 break-inside-avoid overflow-hidden rounded-card border border-line bg-surface">
        <div
          :ref="setBoxRef(row.id)"
          :data-card-id="row.id"
          class="relative w-full overflow-hidden border-b border-line bg-surface-2"
          :style="{ aspectRatio: `${(row.canvas?.w ?? 1920)} / ${(row.canvas?.h ?? 1080)}` }"
        >
          <iframe
            :src="userTemplatePreviewUrl(row.id)"
            class="pointer-events-none absolute left-0 top-0 border-0"
            :style="{
              width: `${row.canvas?.w ?? 1920}px`,
              height: `${row.canvas?.h ?? 1080}px`,
              transform: `scale(${scaleFor(row.id, row.canvas ?? { w: 1920, h: 1080 })})`,
              transformOrigin: 'top left',
            }"
            sandbox="allow-scripts"
            loading="lazy"
            title="模板预览"
          />
          <span
            class="absolute right-2 top-2 rounded-full px-2 py-0.5 text-[10.5px] font-semibold"
            :class="row.visibility === 'public' ? 'bg-accent text-on-accent' : 'bg-surface-3 text-ink-2'"
          >
            {{ row.visibility === 'public' ? '已公开' : '私有' }}
          </span>
          <button
            type="button"
            class="absolute bottom-2 right-2 inline-flex cursor-pointer items-center gap-1 rounded-full bg-surface/90 px-2.5 py-1 text-[11px] font-semibold text-ink-2 shadow-sm backdrop-blur transition-colors hover:bg-surface hover:text-ink"
            title="放大预览"
            @click="openUserPreview(row)"
          >
            <PhArrowsOutSimple :size="11" />
            预览
          </button>
          <!-- 编辑入口：改的是模板 demo 本身（保存滚动备份 5 版；已发布模板需先下架） -->
          <button
            type="button"
            class="absolute bottom-2 right-[74px] inline-flex cursor-pointer items-center gap-1 rounded-full bg-surface/90 px-2.5 py-1 text-[11px] font-semibold text-ink-2 shadow-sm backdrop-blur transition-colors hover:bg-surface hover:text-ink"
            title="手动编辑 demo"
            @click="editing = row"
          >
            <PhPencilSimple :size="11" />
            编辑
          </button>
        </div>
        <div class="space-y-2 p-3">
          <div class="flex items-center justify-between gap-2">
            <p class="truncate text-[13.5px] font-semibold" :title="row.name">{{ row.name }}</p>
            <span
              class="shrink-0 rounded px-1.5 py-0.5 text-[10px]"
              :class="row.status === 'failed' ? 'bg-[#FDEBEC] text-[#9F2F2D]' : 'bg-surface-2 text-ink-3'"
            >
              {{ STATUS_LABELS[row.status] ?? row.status }}
            </span>
          </div>
          <p class="text-[11px] text-ink-3">派生自「{{ baseName(row.base_id) }}」</p>
          <p v-if="row.publish_error" class="line-clamp-3 rounded bg-[#FDEBEC] p-1.5 text-[11px] text-[#9F2F2D]" :title="row.publish_error">
            {{ row.publish_error }}
          </p>
          <div class="flex flex-wrap gap-1.5 pt-1">
            <Button size="sm" @click="router.push(`/my-templates/${row.id}/edit`)">
              <PhPencilSimple :size="12" />
              定制
            </Button>
            <Button
              size="sm"
              :loading="busyId === row.id + ':pub'"
              @click="togglePublish(row)"
            >
              <PhGlobe :size="12" />
              {{ row.visibility === 'public' ? '下架' : '发布' }}
            </Button>
            <Button size="sm" variant="ghost" :loading="busyId === row.id + ':del'" @click="remove(row)">
              <PhTrash :size="12" />
              删除
            </Button>
          </div>
        </div>
      </div>
      </div>
      <Pagination
        v-if="mine.length > MINE_PER_PAGE"
        class="mt-6"
        :model-value="minePage"
        :page-count="minePageCount"
        @update:model-value="setMinePage"
      />
    </div>

    <!-- 社区模板 -->
    <div class="mt-8">
      <h2 class="text-[14px] font-bold">社区模板</h2>
      <p class="mt-0.5 text-[12px] text-ink-3">其他用户发布并通过门禁的模板，可以直接选用或再派生。</p>
      <div v-if="community.length" class="mt-3 flex flex-wrap gap-2">
        <span
          v-for="c in community"
          :key="c.id"
          class="inline-flex items-center gap-2 rounded-full border border-line bg-surface px-3 py-1 text-[12px] text-ink-2"
        >
          {{ c.name }}
          <span class="text-ink-3">来自 {{ c.author }}</span>
          <button class="cursor-pointer font-semibold text-accent hover:underline" @click="fork(c.id)">
            <PhPlus :size="11" class="inline" />
            派生
          </button>
        </span>
      </div>
      <p v-else class="mt-2 text-[12px] text-ink-3">还没有公开的社区模板——发布第一个吧。</p>
    </div>

    <!-- 从内置模板派生 -->
    <div ref="builtinSection" class="mt-8 scroll-mt-4">
      <div class="flex flex-wrap items-center justify-between gap-2">
        <div>
          <h2 class="text-[14px] font-bold">从内置模板派生</h2>
          <p class="mt-0.5 text-[12px] text-ink-3">选一个接近的起点，结构契约继承内置模板，定制只动视觉 token，质量有底。</p>
        </div>
        <Button @click="openBlankDialog">
          从空白新建
        </Button>
      </div>
      <TemplateFilterBar
        v-model:query="builtinQuery"
        :active-tag="builtinTag"
        :vocab="builtinVocab"
        :total="builtins.length"
        :shown="filteredBuiltins.length"
        class="mt-4"
        @update:active-tag="toggleBuiltinTag"
      />
      <div v-if="filteredBuiltins.length" class="mt-3 columns-1 gap-4 sm:columns-2 lg:columns-3">
        <div v-for="b in pagedBuiltins" :key="b.id" class="mb-4 break-inside-avoid overflow-hidden rounded-card border border-line bg-surface transition-colors hover:border-accent">
          <div
            :ref="setBoxRef('base-' + b.id)"
            :data-card-id="'base-' + b.id"
            class="relative w-full overflow-hidden border-b border-line bg-surface-2"
            :style="{ aspectRatio: `${b.canvas.w} / ${b.canvas.h}` }"
          >
            <iframe
              :src="templateApi.previewUrl(b.id, '', 1, 1)"
              class="pointer-events-none absolute left-0 top-0 border-0"
              :style="{
                width: `${b.canvas.w}px`,
                height: `${b.canvas.h}px`,
                transform: `scale(${scaleFor('base-' + b.id, b.canvas)})`,
                transformOrigin: 'top left',
              }"
              sandbox="allow-scripts allow-same-origin"
              loading="lazy"
              :title="b.name + ' 预览'"
            />
            <button
              type="button"
              class="absolute bottom-2 right-2 inline-flex cursor-pointer items-center gap-1 rounded-full bg-surface/90 px-2.5 py-1 text-[11px] font-semibold text-ink-2 shadow-sm backdrop-blur transition-colors hover:bg-surface hover:text-ink"
              title="放大预览"
              @click="openBuiltinPreview(b)"
            >
              <PhArrowsOutSimple :size="11" />
              预览
            </button>
          </div>
          <div class="space-y-2 p-3">
            <p class="text-[13.5px] font-semibold">{{ b.name }}</p>
            <p class="line-clamp-2 text-[11.5px] text-ink-2">{{ b.description }}</p>
            <div v-if="b.scenario?.length" class="flex flex-wrap gap-1">
              <span v-for="s in b.scenario.slice(0, 3)" :key="s" class="rounded bg-surface-2 px-1.5 py-0.5 text-[10px] text-ink-3">
                {{ s }}
              </span>
            </div>
            <div class="pt-1">
              <Button size="sm" :loading="busyId === 'fork:' + b.id" @click="fork(b.id)">
                <PhPlus :size="12" />
                从这个派生
              </Button>
            </div>
          </div>
        </div>
      </div>
      <div
        v-else
        class="mt-4 flex flex-col items-center gap-2 rounded-card border border-dashed border-line px-6 py-10 text-center"
      >
        <p class="text-[13px] text-ink-2">没有匹配的内置模板</p>
        <button class="cursor-pointer text-[12.5px] font-semibold text-accent hover:underline" @click="resetBuiltinFilter">
          清空筛选条件
        </button>
      </div>
      <Pagination
        v-if="filteredBuiltins.length > BUILTIN_PER_PAGE"
        class="mt-6"
        :model-value="builtinPage"
        :page-count="builtinPageCount"
        @update:model-value="setBuiltinPage"
      />
    </div>

    <!-- 模板预览弹窗：放大 + 逐页翻看（单 iframe，postMessage 无刷新翻页） -->
    <TemplatePreviewModal
      :open="preview != null"
      :title="preview?.name ?? ''"
      :canvas="preview?.canvas ?? { w: 1920, h: 1080 }"
      :pages="preview?.pages"
      :sandbox="previewSandbox"
      :src="previewSrc"
      @close="preview = null"
    />

    <!-- 编辑弹窗：用户模板 demo 手改；保存后刷新列表（模板卡 iframe 是 no-cache，重挂即新版） -->
    <DeckEditModal
      :open="editing != null"
      kind="usertpl"
      :id="editing?.id ?? ''"
      :title="editing?.name ?? ''"
      :canvas="editing?.canvas ?? { w: 1920, h: 1080 }"
      @close="editing = null"
      @saved="load()"
    />

    <!-- 从空白新建：命名弹窗（回车=创建，留空用默认名「空白模板」） -->
    <Dialog
      :open="blankOpen"
      title="从空白新建"
      desc="中性灰阶的空白起点：先给两个最小说明版式搭好结构契约，配色、字体、版式都在工作台里从零定制。"
      confirm-text="创建"
      @confirm="createBlank"
      @close="blankOpen = false"
    >
      <input
        ref="blankInput"
        v-model="blankName"
        class="mt-4 w-full rounded-control border border-line bg-surface-2 px-3 py-2 text-[13px] text-ink outline-none placeholder:text-ink-3 focus-visible:border-accent"
        placeholder="模板名称（留空则叫「空白模板」）"
        maxlength="40"
        :disabled="blankCreating"
        @keydown.enter.prevent="createBlank"
      />
    </Dialog>
  </div>
</template>
