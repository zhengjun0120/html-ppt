<script setup lang="ts">
import { PhDownloadSimple, PhMagnifyingGlass } from '@phosphor-icons/vue'
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { ApiError, authedUrl } from '@/api/client'
import { traceApi, type RunMeta, type TraceEvent } from '@/api/trace'
import Badge from '@/components/ui/Badge.vue'
import Button from '@/components/ui/Button.vue'
import Empty from '@/components/ui/Empty.vue'
import Pagination from '@/components/ui/Pagination.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import { useToast } from '@/stores/toast'

const toast = useToast()
const router = useRouter()
const route = useRoute()

// —— 列表：服务端分页 + 筛选。观测数据只增不减，旧版一次拉 100 条且没有翻页
// UI——既慢、老记录又永远看不到。筛选与搜索都在服务端做：q 匹配的是后端索引
// 里的全文，不受列表响应里 user_content 预览截断的影响。——//
const PAGE_SIZE = 25
const runs = ref<RunMeta[]>([])
const total = ref(0)
const listLoading = ref(false)
const search = ref('')
const kind = ref<'' | 'deck' | 'customize' | 'tplsugg'>('')

const kindChips: { value: '' | 'deck' | 'customize' | 'tplsugg'; label: string }[] = [
  { value: '', label: '全部' },
  { value: 'deck', label: '文稿对话' },
  { value: 'customize', label: '模板定制' },
  { value: 'tplsugg', label: '模板推荐' },
]

/** 页码 ↔ URL query（DecksView 同款）：翻页写回 ?page=N，刷新/后退落回原页 */
function pageFromQuery(): number {
  const n = Number(route.query.page)
  return Number.isInteger(n) && n >= 1 ? n : 1
}
const page = ref(pageFromQuery())
const pageCount = computed(() => Math.max(1, Math.ceil(total.value / PAGE_SIZE)))
// 删除/筛选后列表变短，当前页可能越界；页码变化统一写回 query（翻页/回钳都走这里）
watch(pageCount, (n) => {
  if (page.value > n) page.value = n
})
watch(page, (p) => {
  void router.replace({ query: p > 1 ? { page: String(p) } : undefined }).catch(() => {})
  void loadList()
})

let listAbort: AbortController | null = null
async function loadList() {
  // 竞态守卫：慢的旧响应回来时已被更新的请求取代，不许覆盖新数据
  listAbort?.abort()
  const ac = new AbortController()
  listAbort = ac
  listLoading.value = true
  try {
    const r = await traceApi.list({
      limit: PAGE_SIZE,
      offset: (page.value - 1) * PAGE_SIZE,
      kind: kind.value || undefined,
      q: search.value.trim() || undefined,
    })
    if (ac.signal.aborted) return
    runs.value = r.runs
    total.value = r.total
  } catch (e) {
    if (ac.signal.aborted) return
    toast.error(e instanceof ApiError ? e.message : '观测数据加载失败')
  } finally {
    if (!ac.signal.aborted) listLoading.value = false
  }
}

// 筛选/搜索变化一律回第 1 页再加载；已在第 1 页就显式拉（page watch 不会触发）
function applyFiltersReset() {
  if (page.value !== 1) page.value = 1
  else void loadList()
}
watch(kind, applyFiltersReset)
function setKind(k: '' | 'deck' | 'customize' | 'tplsugg') {
  // watch(kind) 统一触发回第 1 页 + 加载，这里只改值——显式再调一次会双发请求
  kind.value = k
}
// 搜索 300ms 防抖：索引过滤虽快，敲一个字发一次请求仍然没必要
let searchTimer: ReturnType<typeof setTimeout> | undefined
watch(search, () => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(applyFiltersReset, 300)
})
onBeforeUnmount(() => clearTimeout(searchTimer))

const listEl = ref<HTMLElement | null>(null)
function setPage(p: number) {
  page.value = p // watch(page) 负责加载
  listEl.value?.scrollTo({ top: 0, behavior: 'smooth' })
}
void loadList()

// —— 详情（选中 run）——
const selected = ref<RunMeta | null>(null)
const events = ref<TraceEvent[]>([])
const detailLoading = ref(false)
const expanded = ref<Record<number, boolean>>({})

let pollTimer: ReturnType<typeof setInterval> | undefined
let offset = 0

function stopPoll() {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = undefined
  }
}

async function selectRun(r: RunMeta) {
  selected.value = r
  events.value = []
  expanded.value = {}
  offset = 0
  stopPoll()
  detailLoading.value = true
  try {
    await fetchEvents()
    // running 的 run 用 from_offset 增量轮询，直到结束
    if (selected.value && isRunning(selected.value)) {
      pollTimer = setInterval(async () => {
        try {
          if (selected.value && !isRunning(selected.value)) {
            stopPoll()
            return
          }
          await fetchEvents()
        } catch { /* 轮询失败静默，下轮再试 */ }
      }, 2500)
    }
  } finally {
    detailLoading.value = false
  }
}

async function fetchEvents() {
  const r = selected.value
  if (!r) return
  const res = await traceApi.run(r.session_id, r.run_id, offset > 0 ? { from_offset: offset, limit: 500 } : { limit: 500 })
  if (res.events.length) events.value = [...events.value, ...res.events]
  offset = res.next_offset
  if (!res.running) stopPoll()
}

function isRunning(r: RunMeta): boolean {
  return r.status === 'running'
}

onBeforeUnmount(stopPoll)

// —— 展示辅助 ——
function fmtDuration(ms?: number): string {
  if (ms == null) return ''
  return ms < 1000 ? `${ms}ms` : `${(ms / 1000).toFixed(1)}s`
}
function fmtTime(iso: string): string {
  const d = new Date(iso)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}
function fmtTokens(n?: number): string {
  return n == null ? '' : n >= 10000 ? `${(n / 1000).toFixed(1)}k` : String(n)
}
function statusTone(s: RunMeta['status']): 'success' | 'warning' | 'danger' | 'info' {
  return s === 'ok' ? 'success' : s === 'paused' ? 'warning' : s === 'error' ? 'danger' : 'info'
}
function pretty(raw: string | undefined): string {
  if (!raw) return ''
  try {
    return JSON.stringify(JSON.parse(raw), null, 2)
  } catch {
    return raw
  }
}

const kindTone: Record<string, 'success' | 'warning' | 'danger' | 'info' | 'accent' | 'neutral'> = {
  run_start: 'neutral',
  llm_request: 'info',
  llm_response: 'success',
  tool_call: 'accent',
  tool_result: 'neutral',
  sub_step: 'warning',
  usage: 'neutral',
  error: 'danger',
  run_end: 'neutral',
}

/** llm_request 单事件展开时按需拉完整上下文（messages 很大，默认不下发） */
const expandedMessages = ref<Record<number, string | 'loading'>>({})
async function toggleEvent(e: TraceEvent) {
  expanded.value[e.seq] = !expanded.value[e.seq]
  if (expanded.value[e.seq] && e.kind === 'llm_request' && selected.value && e.messages == null) {
    expandedMessages.value[e.seq] = 'loading'
    try {
      const full = await traceApi.event(selected.value.session_id, selected.value.run_id, e.seq)
      expandedMessages.value[e.seq] = pretty(JSON.stringify(full.messages ?? null))
    } catch (err) {
      expandedMessages.value[e.seq] = `加载失败：${err instanceof ApiError ? err.message : '未知错误'}`
    }
  }
}

function exportUrl(r: RunMeta): string {
  return authedUrl(`/api/traces/${r.session_id}/${r.run_id}/export`)
}
function imgUrl(r: RunMeta, name: string): string {
  return authedUrl(`/api/traces/${r.session_id}/${r.run_id}/img/${name}`)
}
</script>

<template>
  <div class="flex h-full overflow-hidden">
    <!-- 左：run 列表（服务端分页 + 类型筛选 + 全文搜索） -->
    <aside class="flex min-h-0 w-[400px] shrink-0 flex-col border-r border-line bg-surface">
      <div class="border-b border-line p-3">
        <div class="mb-2 flex items-center justify-between">
          <h1 class="text-[14px] font-bold">观测台</h1>
          <Button size="sm" :loading="listLoading" @click="loadList()">刷新</Button>
        </div>
        <div class="relative">
          <PhMagnifyingGlass class="absolute left-2.5 top-2.5 text-ink-3" :size="14" />
        <input
          v-model="search"
          placeholder="搜 run_id / 用户输入 / deck / 模板…"
          class="w-full rounded-control border border-line bg-surface-2 py-1.5 pl-8 pr-2 text-[12.5px] text-ink outline-none placeholder:text-ink-3 focus-visible:border-accent"
        />
        </div>
        <div class="mt-2 flex flex-wrap gap-1">
          <button
            v-for="kc in kindChips"
            :key="kc.value"
            type="button"
            class="cursor-pointer rounded-full px-2.5 py-0.5 text-[11px] font-semibold transition-colors"
            :class="kind === kc.value ? 'bg-accent text-accent-contrast' : 'bg-surface-2 text-ink-2 hover:text-ink'"
            @click="setKind(kc.value)"
          >
            {{ kc.label }}
          </button>
        </div>
      </div>

      <div ref="listEl" class="min-h-0 flex-1 overflow-y-auto p-2">
        <Skeleton v-if="listLoading && runs.length === 0" v-for="i in 5" :key="i" class="mb-2 h-16 rounded-control" />
        <Empty v-else-if="runs.length === 0" title="没有观测记录" desc="agent 跑过对话后，这里会出现每次 run 的完整轨迹。" />
        <button
          v-for="r in runs"
          :key="r.run_id"
          class="mb-1.5 w-full cursor-pointer rounded-control border p-2.5 text-left transition-colors"
          :class="selected?.run_id === r.run_id ? 'border-accent bg-accent-soft' : 'border-line bg-surface hover:border-line-strong'"
          @click="selectRun(r)"
        >
          <div class="flex items-center gap-2">
            <Badge :tone="statusTone(r.status)">{{ r.status }}</Badge>
            <Badge v-if="r.run_kind === 'customize'" tone="info">定制</Badge>
            <Badge v-if="r.run_kind === 'tplsugg'" tone="info">模板推荐</Badge>
            <span class="truncate font-mono text-[11px] text-ink-3">{{ r.run_id }}</span>
            <span class="ml-auto flex shrink-0 items-center gap-1.5 text-[11px] text-ink-3">
              <span v-if="r.usage?.total?.total" :title="`输入 ${r.usage.total.prompt ?? 0} / 输出 ${r.usage.total.completion ?? 0} / 缓存 ${r.usage.total.cached ?? 0}`" class="font-mono">{{ fmtTokens(r.usage.total.total) }} tok</span>
              <span>{{ fmtDuration(r.duration_ms) }}</span>
            </span>
          </div>
          <p class="mt-1 line-clamp-2 text-[12px] text-ink-2">{{ r.user_content || '（无输入）' }}</p>
          <p class="mt-0.5 truncate font-mono text-[10.5px] text-ink-3">
            {{ fmtTime(r.started_at) }}<template v-if="r.deck_id"> · {{ r.deck_id }}</template><template v-if="r.turns"> · {{ r.turns }} 轮</template>
          </p>
        </button>
      </div>

      <!-- 分页脚注：Pagination 在单页时自己不渲染，共 N 条也只在有数据时露 -->
      <div v-if="total > 0" class="border-t border-line p-2">
        <p class="mb-1.5 text-center text-[11px] text-ink-3">共 {{ total }} 条记录</p>
        <Pagination :model-value="page" :page-count="pageCount" @update:model-value="setPage" />
      </div>
    </aside>

    <!-- 右：run 详情 -->
    <section class="min-h-0 min-w-0 flex-1 overflow-y-auto">
      <Empty v-if="!selected" title="选择一个 run 查看轨迹" desc="左侧列表按时间排列；选中的 run 会展示完整事件时间线、用量与审查截图。" />

      <template v-else>
        <!-- run 摘要头 -->
        <div class="border-b border-line bg-surface p-4">
          <div class="flex flex-wrap items-center gap-2">
            <Badge :tone="statusTone(selected.status)">{{ selected.status }}</Badge>
            <Badge v-if="selected.run_kind === 'customize'" tone="info">定制</Badge>
            <Badge v-if="selected.run_kind === 'tplsugg'" tone="info">模板推荐</Badge>
            <span class="font-mono text-[12px]">{{ selected.run_id }}</span>
            <span v-if="selected.parent_run_id" class="font-mono text-[10.5px] text-ink-3">← {{ selected.parent_run_id }}</span>
            <a
              class="ml-auto inline-flex cursor-pointer items-center gap-1 rounded border border-line bg-surface-2 px-2 py-1 text-[11.5px] text-ink-2 hover:border-line-strong hover:text-ink"
              :href="exportUrl(selected)"
            >
              <PhDownloadSimple :size="12" /> 导出 JSON
            </a>
          </div>
          <p class="mt-2 text-[13px] text-ink">{{ selected.user_content || '（无输入）' }}</p>
          <div class="mt-2 flex flex-wrap gap-x-5 gap-y-1 text-[11.5px] text-ink-3">
            <span>{{ fmtTime(selected.started_at) }}</span>
            <span>耗时 {{ fmtDuration(selected.duration_ms) }}</span>
            <span v-if="selected.turns != null">{{ selected.turns }} 轮</span>
            <span v-if="selected.tool_calls != null">{{ selected.tool_calls }} 次工具</span>
            <span v-if="selected.model" class="font-mono">{{ selected.model }}</span>
            <span v-if="selected.deck_id" class="font-mono">{{ selected.deck_id }}</span>
            <span v-if="selected.usage?.total" title="输入含缓存命中；输出含推理 token">tokens {{ fmtTokens(selected.usage.total.total) }}（输入 {{ fmtTokens(selected.usage.total.prompt) }} / 输出 {{ fmtTokens(selected.usage.total.completion) }}<template v-if="selected.usage.total.cached">，缓存 {{ fmtTokens(selected.usage.total.cached) }}</template><template v-if="selected.usage.total.reasoning">，推理 {{ fmtTokens(selected.usage.total.reasoning) }}</template>）</span>
          </div>
          <!-- 分项用量（主循环 / 视觉审查 / 联网搜索） -->
          <div v-if="selected.usage?.usage" class="mt-2 flex flex-wrap gap-2">
            <Badge v-for="(u, name) in selected.usage.usage" v-show="u" :key="name" tone="neutral">
              <span class="font-mono">{{ name }}</span>
              ：{{ fmtTokens(u!.total) }} tok / {{ u!.calls }} 次
            </Badge>
          </div>
        </div>

        <!-- 事件时间线 -->
        <div class="p-3">
          <Skeleton v-if="detailLoading && events.length === 0" v-for="i in 4" :key="i" class="mb-2 h-10 rounded-control" />
          <template v-for="e in events" :key="e.seq">
            <!-- 分轮标题 -->
            <p v-if="e.turn != null && (e.seq === 0 || events.find(x => x.seq === e.seq - 1)?.turn !== e.turn)" class="mb-1.5 mt-3 text-[10.5px] font-semibold tracking-wider text-ink-3">
              第 {{ e.turn }} 轮
            </p>
            <div
              class="mb-1.5 cursor-pointer rounded-control border bg-surface p-2.5 text-[12.5px] transition-colors"
              :class="e.kind === 'error' ? 'border-danger/40' : 'border-line hover:border-line-strong'"
              @click="toggleEvent(e)"
            >
              <div class="flex flex-wrap items-center gap-2">
                <Badge :tone="kindTone[e.kind] ?? 'neutral'">{{ e.kind }}</Badge>
                <span v-if="e.tool_name" class="font-mono font-semibold">{{ e.tool_name }}</span>
                <span v-if="e.finish_reason" class="text-ink-3">finish: {{ e.finish_reason }}</span>
                <span v-if="e.message_count != null" class="text-ink-3">{{ e.message_count }} 条上下文</span>
                <span v-if="e.component" class="text-ink-3">{{ e.component }}</span>
                <span v-if="e.duration_ms != null" class="ml-auto text-ink-3">{{ fmtDuration(e.duration_ms) }}</span>
                <span v-else-if="e.bytes" class="ml-auto text-ink-3">{{ fmtTokens(e.bytes) }}B</span>
              </div>

              <p v-if="e.user_content" class="mt-1 line-clamp-2 text-ink-2">{{ e.user_content }}</p>
              <p v-if="e.content" class="mt-1 line-clamp-3 whitespace-pre-wrap break-words text-ink-2">{{ e.content }}</p>
              <p v-if="e.args" class="mt-1 truncate font-mono text-[11px] text-ink-3">{{ e.args }}</p>
              <p v-if="e.error" class="mt-1 text-danger">{{ e.error }}</p>

              <!-- 审查截图 -->
              <div v-if="e.images?.length" class="mt-2 flex flex-wrap gap-2" @click.stop>
                <img
                  v-for="img in e.images"
                  :key="img.name"
                  :src="imgUrl(selected, img.name)"
                  :alt="img.label || img.name"
                  class="h-20 rounded border border-line"
                  loading="lazy"
                />
              </div>

              <!-- 展开体：参数 / 结果 / 完整上下文 -->
              <div v-if="expanded[e.seq]" class="mt-2 flex flex-col gap-2" @click.stop>
                <pre v-if="e.args" class="max-h-48 overflow-auto whitespace-pre-wrap break-all rounded bg-code p-2 font-mono text-[11px] text-ink-2">{{ pretty(e.args) }}</pre>
                <pre v-if="e.result" class="max-h-64 overflow-auto whitespace-pre-wrap break-all rounded bg-code p-2 font-mono text-[11px] text-ink-2">{{ pretty(e.result) }}</pre>
                <template v-if="e.kind === 'llm_request' && e.messages == null">
                  <p class="text-[11px] text-ink-3">
                    {{ expandedMessages[e.seq] === 'loading' ? '加载完整上下文…' : (expandedMessages[e.seq] ?? '点击已展开；完整上下文在再次点击时加载') }}
                  </p>
                </template>
                <pre v-if="e.kind === 'llm_request' && expandedMessages[e.seq] && expandedMessages[e.seq] !== 'loading'" class="max-h-80 overflow-auto whitespace-pre-wrap break-all rounded bg-code p-2 font-mono text-[11px] text-ink-2">{{ expandedMessages[e.seq] }}</pre>
                <div v-if="e.usage" class="text-[11px] text-ink-3">
                  tokens：{{ e.usage.total }}（输入 {{ e.usage.prompt }} / 输出 {{ e.usage.completion }}，缓存 {{ e.usage.cached }}<template v-if="e.usage.reasoning">，推理 {{ e.usage.reasoning }}</template>）
                </div>
                <div v-if="e.summary" class="text-[11px] text-ink-3">
                  run 汇总：{{ e.summary.turns }} 轮 / {{ e.summary.tool_calls }} 次工具 / {{ fmtDuration(e.summary.duration_ms) }}
                </div>
              </div>
            </div>
          </template>
        </div>
      </template>
    </section>
  </div>
</template>
