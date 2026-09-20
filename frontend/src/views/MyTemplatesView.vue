<script setup lang="ts">
import { PhArrowClockwise, PhGlobe, PhGlobeHemisphereWest, PhPencilSimple, PhPlus, PhSpinner, PhTrash } from '@phosphor-icons/vue'
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import { templateApi, type TemplateMeta } from '@/api/templates'
import { userTemplateApi, userTemplatePreviewUrl, type CommunityTemplate, type UserTemplateRow } from '@/api/userTemplates'
import Button from '@/components/ui/Button.vue'
import Empty from '@/components/ui/Empty.vue'
import { useToast } from '@/stores/toast'

/**
 * 我的模板（plan-v3 B3）：克隆/定制/发布/下架/删除的管理页。
 * 发布跑自动门禁（结构校验 + demo 渲染量测），结果与失败原因就地展示。
 * 「从内置模板派生」直接展示内置模板卡片（真实缩略 + 中文名 + 适用场景标签）——
 * 派生是挑"视觉起点"，看不见起点就没法挑。
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
/** base_id → 中文名（内置模板都有中文名；查不到兜底原 id） */
function baseName(id: string): string {
  return metaOf.value.get(id)?.name ?? id
}

async function load() {
  loading.value = true
  try {
    const [m, c, t] = await Promise.all([userTemplateApi.list(), userTemplateApi.community(), templateApi.list()])
    mine.value = m
    community.value = c
    templates.value = t
  } catch (e) {
    toast.error(e instanceof Error ? e.message : '模板加载失败')
  } finally {
    loading.value = false
  }
}

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
    <div v-else class="columns-1 gap-4 sm:columns-2 lg:columns-3">
      <div v-for="row in mine" :key="row.id" class="mb-4 break-inside-avoid overflow-hidden rounded-card border border-line bg-surface">
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
    <div class="mt-8">
      <h2 class="text-[14px] font-bold">从内置模板派生</h2>
      <p class="mt-0.5 text-[12px] text-ink-3">选一个接近的起点，结构契约继承内置模板，定制只动视觉 token，质量有底。</p>
      <div class="mt-3 columns-1 gap-4 sm:columns-2 lg:columns-3">
        <div v-for="b in builtins" :key="b.id" class="mb-4 break-inside-avoid overflow-hidden rounded-card border border-line bg-surface transition-colors hover:border-accent">
          <div
            :ref="setBoxRef('base-' + b.id)"
            :data-card-id="'base-' + b.id"
            class="relative w-full overflow-hidden border-b border-line bg-surface-2"
            :style="{ aspectRatio: `${b.canvas.w} / ${b.canvas.h}` }"
          >
            <iframe
              :src="templateApi.previewUrl(b.id, '', 1)"
              class="pointer-events-none absolute left-0 top-0 border-0"
              :style="{
                width: `${b.canvas.w}px`,
                height: `${b.canvas.h}px`,
                transform: `scale(${scaleFor('base-' + b.id, b.canvas)})`,
                transformOrigin: 'top left',
              }"
              sandbox="allow-scripts"
              loading="lazy"
              :title="b.name + ' 预览'"
            />
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
    </div>
  </div>
</template>
