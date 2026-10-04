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
import { fieldNotes, kindNotes, roleNotes, subNameZh, subStageZh, toolNotes } from './traceNotes'

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

/** llm_request 展开后的完整上下文。按消息拆块渲染（role 徽标 + 正文按真实
 * 换行显示）——此前是整包 pretty JSON：几 KB 的 content 字符串被 JSON 转义成
 * 一行，把缩进结构和换行全部淹没，用户看到的就是一坨带字面 \n 的文本。 */
type ExpandedMessage = { role: string; blocks: KvBlock[] }
type ExpandedCtx =
  | { state: 'loading' }
  | { state: 'error'; message: string }
  | { state: 'msgs'; msgs: ExpandedMessage[] }
  | { state: 'json'; json: string }
const expandedCtx = ref<Record<number, ExpandedCtx | undefined>>({})

// 消息正文三种形态：散文（JSON 解析失败 → 原样）、JSON 字符串（tool 结果，
// 如 {"rules":…} → 按 key 拆块、真换行）、非字符串（多模态等 → pretty JSON）。
// assistant 发起工具调用的消息 content 为 null，列出工具名单即可——参数全文
// 在时间线的 tool_call 卡片里，重复展示只会变成第二个转义沼泽。
function splitMessages(raw: unknown): ExpandedCtx {
  if (!Array.isArray(raw)) {
    return { state: 'json', json: pretty(JSON.stringify(raw ?? null)) }
  }
  const msgs = raw.map((m): ExpandedMessage => {
    const obj = (typeof m === 'object' && m !== null ? m : {}) as Record<string, unknown>
    const role = typeof obj.role === 'string' ? obj.role : '?'
    const blocks: KvBlock[] = []
    if (typeof obj.content === 'string') {
      if (obj.content.trim() !== '') blocks.push(...kvBlocks(obj.content))
    } else if (obj.content != null) {
      blocks.push({ key: '', text: pretty(JSON.stringify(obj.content)), isJson: true })
    }
    if (Array.isArray(obj.tool_calls)) {
      const names = (obj.tool_calls as unknown[])
        .map((tc) => {
          const fn = (typeof tc === 'object' && tc !== null ? (tc as Record<string, unknown>).function : null) as Record<string, unknown> | null
          return fn && typeof fn.name === 'string' ? fn.name : ''
        })
        .filter(Boolean)
      if (names.length) {
        blocks.push({
          key: '发起工具调用',
          text: names.map((n) => `→ ${n}${toolNotes[n] ? ' · ' + toolNotes[n].zh : ''}`).join('\n'),
        })
      }
    }
    return { role, blocks }
  })
  return { state: 'msgs', msgs }
}

async function toggleEvent(e: TraceEvent) {
  expanded.value[e.seq] = !expanded.value[e.seq]
  if (expanded.value[e.seq] && e.kind === 'llm_request' && selected.value && e.messages == null) {
    expandedCtx.value[e.seq] = { state: 'loading' }
    try {
      const full = await traceApi.event(selected.value.session_id, selected.value.run_id, e.seq)
      expandedCtx.value[e.seq] = splitMessages(full.messages ?? null)
    } catch (err) {
      expandedCtx.value[e.seq] = {
        state: 'error',
        message: `加载失败：${err instanceof ApiError ? err.message : '未知错误'}`,
      }
    }
  }
}
function ctxState(seq: number): ExpandedCtx['state'] | 'none' {
  return expandedCtx.value[seq]?.state ?? 'none'
}
function ctxMsgs(seq: number): ExpandedMessage[] {
  const v = expandedCtx.value[seq]
  return v && v.state === 'msgs' ? v.msgs : []
}
function ctxJson(seq: number): string {
  const v = expandedCtx.value[seq]
  return v && v.state === 'json' ? v.json : ''
}
function ctxError(seq: number): string {
  const v = expandedCtx.value[seq]
  return v && v.state === 'error' ? v.message : ''
}

/** 工具参数/结果展开：对象按 key 拆块——string 值按真实换行渲染。此前整包
 * pretty JSON：长文本值（模板规则、大纲全文等）被重新转义成一行，\n 全是
 * 字面量、换行与结构全被淹没。非对象/解析失败退回 pretty JSON 兜底。 */
type KvBlock = { key: string; text: string; isJson?: boolean }
function kvBlocks(raw: string): KvBlock[] {
  let obj: unknown
  try {
    obj = JSON.parse(raw)
  } catch {
    return [{ key: '', text: raw }]
  }
  if (typeof obj !== 'object' || obj === null || Array.isArray(obj)) {
    return [{ key: '', text: pretty(raw) }]
  }
  const entries = Object.entries(obj as Record<string, unknown>)
  if (entries.length === 0) {
    return [{ key: '', text: '{}' }]
  }
  return entries.map(([k, v]) =>
    typeof v === 'string'
      ? { key: k, text: v }
      : { key: k, text: pretty(JSON.stringify(v)), isJson: true },
  )
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
          <p class="mt-2 whitespace-pre-wrap break-words text-[13px] text-ink">{{ selected.user_content || '（无输入）' }}</p>
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
                <span v-if="kindNotes[e.kind]" class="text-[11px] text-ink-3">{{ kindNotes[e.kind] }}</span>
                <span v-if="e.tool_name" class="font-mono font-semibold">{{ e.tool_name }}</span>
                <span v-if="e.tool_name && toolNotes[e.tool_name]" class="text-[11px] text-ink-3" :title="toolNotes[e.tool_name].desc">{{ toolNotes[e.tool_name].zh }}</span>
                <!-- sub_step 的来源标签（vision / web_search · 阶段）——后端从不写 component，名字只在 sub 里 -->
                <span v-if="e.sub" class="font-mono font-semibold">{{ e.sub.name }}<template v-if="e.sub.stage"> · {{ e.sub.stage }}</template></span>
                <span v-if="e.sub && (subNameZh(e.sub.name) || subStageZh(e.sub.name, e.sub.stage))" class="text-[11px] text-ink-3">{{ [subNameZh(e.sub.name), subStageZh(e.sub.name, e.sub.stage)].filter(Boolean).join(' · ') }}</span>
                <span v-if="e.finish_reason" class="text-ink-3">finish: {{ e.finish_reason }}</span>
                <span v-if="e.message_count != null" class="text-ink-3">{{ e.message_count }} 条上下文</span>
                <span v-if="e.component" class="text-ink-3">{{ e.component }}</span>
                <span v-if="e.duration_ms != null" class="ml-auto text-ink-3">{{ fmtDuration(e.duration_ms) }}</span>
                <span v-else-if="e.bytes" class="ml-auto text-ink-3">{{ fmtTokens(e.bytes) }}B</span>
              </div>

              <p v-if="e.user_content" class="mt-1 line-clamp-2 whitespace-pre-wrap break-words text-ink-2">{{ e.user_content }}</p>
              <!-- 折叠行是 3 行预览（结尾的 … 就是它），展开后由展开体里的全文块接管 -->
              <p v-if="e.content && !expanded[e.seq]" class="mt-1 line-clamp-3 whitespace-pre-wrap break-words text-ink-2">{{ e.content }}</p>
              <p v-if="e.args" class="mt-1 truncate font-mono text-[11px] text-ink-3">{{ e.args }}</p>
              <p v-if="e.error" class="mt-1 whitespace-pre-wrap break-words text-danger">{{ e.error }}</p>

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
                <!-- 工具参数/结果：对象按 key 拆块，string 值（规则/大纲/页面全文等）按真实换行渲染 -->
                <template v-if="e.args">
                  <div class="rounded bg-code p-2">
                    <p class="mb-1.5 text-[10.5px] font-semibold tracking-wider text-ink-3">
                      参数<template v-if="e.tool_name"> · {{ e.tool_name }}<template v-if="toolNotes[e.tool_name]">（{{ toolNotes[e.tool_name].zh }}）</template></template>
                    </p>
                    <div v-for="(b, bi) in kvBlocks(e.args)" :key="bi" :class="bi > 0 ? 'mt-2 border-t border-line/60 pt-2' : ''">
                      <template v-if="b.key">
                        <p class="mb-0.5 font-mono text-[10.5px] font-semibold text-accent">
                          {{ b.key }}<template v-if="fieldNotes[b.key]"><span class="ml-1.5 font-normal text-ink-2">{{ fieldNotes[b.key].zh }}</span></template><span v-if="b.isJson" class="ml-1.5 font-normal text-ink-3">非文本，按 JSON</span>
                        </p>
                        <p v-if="fieldNotes[b.key]?.desc" class="mb-1 text-[10.5px] leading-snug text-ink-3">{{ fieldNotes[b.key].desc }}</p>
                      </template>
                      <pre class="max-h-96 overflow-auto whitespace-pre-wrap break-words font-mono text-[11px] leading-relaxed text-ink-2">{{ b.text }}</pre>
                    </div>
                  </div>
                </template>
                <template v-if="e.result">
                  <div class="rounded bg-code p-2">
                    <p class="mb-1.5 text-[10.5px] font-semibold tracking-wider text-ink-3">
                      结果<template v-if="e.tool_name"> · {{ e.tool_name }}<template v-if="toolNotes[e.tool_name]">（{{ toolNotes[e.tool_name].zh }}）</template></template>
                    </p>
                    <div v-for="(b, bi) in kvBlocks(e.result)" :key="bi" :class="bi > 0 ? 'mt-2 border-t border-line/60 pt-2' : ''">
                      <template v-if="b.key">
                        <p class="mb-0.5 font-mono text-[10.5px] font-semibold text-accent">
                          {{ b.key }}<template v-if="fieldNotes[b.key]"><span class="ml-1.5 font-normal text-ink-2">{{ fieldNotes[b.key].zh }}</span></template><span v-if="b.isJson" class="ml-1.5 font-normal text-ink-3">非文本，按 JSON</span>
                        </p>
                        <p v-if="fieldNotes[b.key]?.desc" class="mb-1 text-[10.5px] leading-snug text-ink-3">{{ fieldNotes[b.key].desc }}</p>
                      </template>
                      <pre class="max-h-96 overflow-auto whitespace-pre-wrap break-words font-mono text-[11px] leading-relaxed text-ink-2">{{ b.text }}</pre>
                    </div>
                  </div>
                </template>
                <!-- 模型返回/sub_step 等的 content 全文（llm_request 没有 content，走下面的完整上下文块）；
                     JSON 形态的返回自动缩进，纯文本原样保留 -->
                <pre v-if="e.content && e.kind !== 'llm_request'" class="max-h-96 overflow-auto whitespace-pre-wrap break-words rounded bg-code p-2 font-mono text-[11px] leading-relaxed text-ink-2">{{ pretty(e.content) }}</pre>
                <!-- sub_step 的正文与结构化数据（搜索词、审查报告、逐页量测等此前完全不可见） -->
                <div v-if="e.sub && (e.sub.text || e.sub.data != null)" class="rounded bg-code p-2">
                  <p class="mb-1.5 text-[10.5px] font-semibold tracking-wider text-ink-3">
                    {{ e.sub.name }}<template v-if="e.sub.stage"> · {{ e.sub.stage }}</template><template v-if="subNameZh(e.sub.name) || subStageZh(e.sub.name, e.sub.stage)"><span class="ml-1 font-normal text-ink-2">{{ [subNameZh(e.sub.name), subStageZh(e.sub.name, e.sub.stage)].filter(Boolean).join(' · ') }}</span></template>
                  </p>
                  <pre v-if="e.sub.text" class="max-h-96 overflow-auto whitespace-pre-wrap break-words font-mono text-[11px] leading-relaxed text-ink-2">{{ e.sub.text }}</pre>
                  <pre v-if="e.sub.data != null" class="mt-2 max-h-80 overflow-auto whitespace-pre-wrap break-all border-t border-line/60 pt-2 font-mono text-[11px] text-ink-2">{{ pretty(JSON.stringify(e.sub.data)) }}</pre>
                </div>
                <template v-if="e.kind === 'llm_request' && e.messages == null">
                  <p v-if="ctxState(e.seq) === 'loading'" class="text-[11px] text-ink-3">加载完整上下文…</p>
                  <p v-else-if="ctxState(e.seq) === 'error'" class="text-[11px] text-danger">{{ ctxError(e.seq) }}</p>
                  <!-- 每条消息一块：role 徽标（带中文说明）+ 正文按真实换行渲染（JSON 字符串正文按 key 再拆） -->
                  <div v-else-if="ctxState(e.seq) === 'msgs'" class="flex flex-col gap-2">
                    <div v-for="(m, mi) in ctxMsgs(e.seq)" :key="mi" class="rounded bg-code p-2">
                      <div class="flex flex-wrap items-center gap-1.5">
                        <Badge tone="neutral">{{ m.role }}</Badge>
                        <span v-if="roleNotes[m.role]" class="text-[10.5px] text-ink-3">{{ roleNotes[m.role] }}</span>
                      </div>
                      <div v-if="m.blocks.length" class="mt-1.5 flex flex-col gap-2">
                        <div v-for="(b, bi) in m.blocks" :key="bi">
                          <template v-if="b.key">
                            <p class="mb-0.5 font-mono text-[10.5px] font-semibold text-accent">
                              {{ b.key }}<template v-if="fieldNotes[b.key]"><span class="ml-1.5 font-normal text-ink-2">{{ fieldNotes[b.key].zh }}</span></template><span v-if="b.isJson" class="ml-1.5 font-normal text-ink-3">非文本，按 JSON</span>
                            </p>
                            <p v-if="fieldNotes[b.key]?.desc" class="mb-1 text-[10.5px] leading-snug text-ink-3">{{ fieldNotes[b.key].desc }}</p>
                          </template>
                          <pre class="max-h-96 overflow-auto whitespace-pre-wrap break-words font-mono text-[11px] leading-relaxed text-ink-2">{{ b.text }}</pre>
                        </div>
                      </div>
                      <p v-else class="mt-1.5 text-[11px] text-ink-3">（无文本正文）</p>
                    </div>
                  </div>
                  <pre v-else-if="ctxState(e.seq) === 'json'" class="max-h-80 overflow-auto whitespace-pre-wrap break-all rounded bg-code p-2 font-mono text-[11px] text-ink-2">{{ ctxJson(e.seq) }}</pre>
                  <p v-else class="text-[11px] text-ink-3">点击已展开；完整上下文在再次点击时加载</p>
                </template>
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
