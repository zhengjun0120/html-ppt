import { computed, ref, toValue } from 'vue'
import type { ComputedRef, MaybeRefOrGetter, Ref } from 'vue'

/**
 * 模板标签筛选 + 名字搜索（docs/template-filter-plan.md）。
 * 纯逻辑函数与 UI 解耦：TemplateMeta（/api/templates 全量返回）和画廊卡片
 * 结构上都满足 FilterableTemplate，两个界面共用同一套筛选。
 */

/** 能被筛选的模板最小形状；id/name 必有，其余缺省当空串/空数组处理 */
export interface FilterableTemplate {
  id: string
  name: string
  description?: string
  tags?: string[]
  scenario?: string[]
}

export interface TagVocabEntry {
  tag: string
  count: number
}

/** 词表从数据实时算出（不写死）：按出现次数降序、同次数按标签字典序。 */
export function templateTags(list: FilterableTemplate[]): TagVocabEntry[] {
  const counts = new Map<string, number>()
  for (const t of list) {
    // 单模板内去重 + trim（空白 tag 不成词）：一个模板写两个等价 tag 只算一次，不然词表计数虚高
    const seen = new Set<string>()
    for (const raw of t.tags ?? []) {
      const tag = raw.trim()
      if (tag) seen.add(tag)
    }
    for (const tag of seen) counts.set(tag, (counts.get(tag) ?? 0) + 1)
  }
  return [...counts.entries()]
    .map(([tag, count]) => ({ tag, count }))
    .sort((a, b) => b.count - a.count || (a.tag < b.tag ? -1 : a.tag > b.tag ? 1 : 0))
}

/** 单条匹配：tag 非空时先卡标签；query trim+lowercase 后对五个字段做整串子串匹配。 */
export function matchesTemplate(t: FilterableTemplate, query: string, tag: string): boolean {
  if (tag && !(t.tags ?? []).includes(tag)) return false
  const q = query.trim().toLowerCase()
  if (!q) return true
  if (t.id.toLowerCase().includes(q) || t.name.toLowerCase().includes(q)) return true
  if ((t.description ?? '').toLowerCase().includes(q)) return true
  return (
    (t.tags ?? []).some((x) => x.toLowerCase().includes(q)) ||
    (t.scenario ?? []).some((x) => x.toLowerCase().includes(q))
  )
}

/** 过滤：保持输入顺序（画廊「我的模板置顶」等上游排序不被打乱）。 */
export function filterTemplates(list: FilterableTemplate[], query: string, tag: string): FilterableTemplate[] {
  return list.filter((t) => matchesTemplate(t, query, tag))
}

export interface TemplateFilterState {
  query: Ref<string>
  activeTag: Ref<string>
  vocab: ComputedRef<TagVocabEntry[]>
  filtered: ComputedRef<FilterableTemplate[]>
  hasFilter: ComputedRef<boolean>
  /** 点 chip 选中，再点同一个取消（回到「全部」） */
  toggleTag: (tag: string) => void
  reset: () => void
}

/** 两处界面（画廊 / 我的模板派生区）共用的筛选状态；source 可以是任何响应式来源。 */
export function useTemplateFilter(source: MaybeRefOrGetter<FilterableTemplate[]>): TemplateFilterState {
  const query = ref('')
  const activeTag = ref('')
  const vocab = computed(() => templateTags(toValue(source)))
  const filtered = computed(() => filterTemplates(toValue(source), query.value, activeTag.value))
  const hasFilter = computed(() => query.value.trim() !== '' || activeTag.value !== '')
  function toggleTag(tag: string) {
    activeTag.value = activeTag.value === tag ? '' : tag
  }
  function reset() {
    query.value = ''
    activeTag.value = ''
  }
  return { query, activeTag, vocab, filtered, hasFilter, toggleTag, reset }
}
