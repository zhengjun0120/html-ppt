<script setup lang="ts">
import { PhClockCounterClockwise } from '@phosphor-icons/vue'
import { ref, watch } from 'vue'

import { deckApi, type VersionMeta } from '@/api/decks'
import { ApiError } from '@/api/client'
import Badge from '@/components/ui/Badge.vue'
import Button from '@/components/ui/Button.vue'
import Dialog from '@/components/ui/Dialog.vue'
import Drawer from '@/components/ui/Drawer.vue'
import Empty from '@/components/ui/Empty.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import { useToast } from '@/stores/toast'

const props = defineProps<{ open: boolean; deckId: string }>()
const emit = defineEmits<{ close: []; restored: [] }>()

const toast = useToast()
const versions = ref<VersionMeta[]>([])
const loading = ref(false)
const loaded = ref(false)

// 确认弹窗状态：restore=恢复到某版本 / delete=删除某版本 / clear=清空全部
const confirm = ref<null | { kind: 'restore' | 'delete' | 'clear'; version?: string }>(null)

async function load() {
  loading.value = true
  try {
    versions.value = await deckApi.history(props.deckId)
    loaded.value = true
  } catch (e) {
    toast.error(e instanceof ApiError ? e.message : '历史加载失败')
  } finally {
    loading.value = false
  }
}

watch(
  () => [props.open, props.deckId],
  ([o]) => {
    if (o) void load()
  },
)

function fmtTime(unix: number): string {
  const d = new Date(unix * 1000)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

async function doConfirm() {
  const c = confirm.value
  if (!c) return
  try {
    if (c.kind === 'restore' && c.version) {
      const r = await deckApi.restore(props.deckId, c.version)
      toast.success(`已恢复到 ${r.restored}，预览已刷新`)
      emit('restored')
    } else if (c.kind === 'delete' && c.version) {
      await deckApi.deleteVersion(props.deckId, c.version)
      toast.success(`已删除 ${c.version}`)
    } else if (c.kind === 'clear') {
      const r = await deckApi.clearHistory(props.deckId)
      toast.info(`已清空 ${r.deleted} 个历史版本`)
    }
    confirm.value = null
    await load()
  } catch (e) {
    toast.error(e instanceof ApiError ? e.message : '操作失败')
  }
}
</script>

<template>
  <Drawer :open="open" title="历史版本" wide @close="emit('close')">
    <div class="p-4">
      <div class="mb-3 flex items-center justify-between">
        <p class="text-[12px] text-ink-3">共 {{ versions.length }} 个版本 · 恢复会以新版本形式追加，不会丢失当前内容</p>
        <Button
          v-if="versions.length"
          variant="danger"
          size="sm"
          @click="confirm = { kind: 'clear' }"
        >
          清空历史
        </Button>
      </div>

      <template v-if="loading">
        <div v-for="i in 3" :key="i" class="mb-2 flex items-center gap-3 rounded-control border border-line bg-surface p-3">
          <Skeleton class="h-4 w-20" />
          <Skeleton class="h-4 flex-1" />
          <Skeleton class="h-4 w-14" />
        </div>
      </template>

      <Empty
        v-else-if="loaded && versions.length === 0"
        title="还没有历史版本"
        desc="agent 每次改动文稿后都会自动留一个版本快照，可随时回看或恢复。"
      />

      <ul v-else class="flex flex-col gap-2">
        <li
          v-for="v in versions"
          :key="v.version"
          class="rounded-control border border-line bg-surface p-3 transition-colors hover:border-line-strong"
        >
          <div class="flex items-center gap-2">
            <span class="font-mono text-[12.5px] font-semibold">{{ v.version }}</span>
            <Badge :tone="v.operation === 'restore' ? 'info' : 'accent'">
              {{ v.operation === 'restore' ? '恢复' : 'run' }}
            </Badge>
            <span class="text-[11.5px] text-ink-3">{{ v.slides }} 页 · {{ fmtTime(v.time) }}</span>
            <span class="ml-auto flex shrink-0 gap-1.5">
              <Button size="sm" @click="confirm = { kind: 'restore', version: v.version }">
                <PhClockCounterClockwise :size="12" />
                恢复
              </Button>
              <Button variant="danger" size="sm" @click="confirm = { kind: 'delete', version: v.version }">
                删除
              </Button>
            </span>
          </div>
          <p class="mt-1 line-clamp-2 text-[12px] text-ink-2" :title="v.detail">{{ v.detail }}</p>
        </li>
      </ul>
    </div>

    <Dialog
      :open="!!confirm"
      :title="confirm?.kind === 'restore' ? `恢复到 ${confirm?.version}？` : confirm?.kind === 'delete' ? `删除版本 ${confirm?.version}？` : '清空全部历史版本？'"
      :desc="confirm?.kind === 'restore' ? '恢复会以一个新版本追加，当前内容仍可在历史中找到。' : confirm?.kind === 'delete' ? '该版本删除后不可找回。' : `将删除 ${versions.length} 个版本，不可找回。`"
      :confirm-text="confirm?.kind === 'restore' ? '恢复' : confirm?.kind === 'delete' ? '删除' : '全部删除'"
      :danger="confirm?.kind !== 'restore'"
      @confirm="doConfirm"
      @close="confirm = null"
    />
  </Drawer>
</template>
