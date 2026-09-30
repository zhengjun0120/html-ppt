import { existsSync, readdirSync, readFileSync } from 'node:fs'
import { dirname, join, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { describe, expect, it } from 'vitest'
import { nextTick, ref } from 'vue'

import {
  filterTemplates,
  matchesTemplate,
  templateTags,
  useTemplateFilter,
  type FilterableTemplate,
} from './templateFilter'

/** 造数工厂：只关心筛选相关字段，其余给合法缺省 */
function tpl(p: Partial<FilterableTemplate> & Pick<FilterableTemplate, 'id' | 'name'>): FilterableTemplate {
  return { description: '', ...p }
}

const LIB_DIR = dirname(fileURLToPath(import.meta.url))
// test 文件在 frontend/src/lib/ → 往上三级是仓库根，内置模板数据源在 backend/templates/
const TEMPLATES_DIR = resolve(LIB_DIR, '../../../backend/templates')

describe('templateTags（标签词表）', () => {
  it('统计每个标签的出现次数', () => {
    const list = [
      tpl({ id: 'a', name: 'A', tags: ['深色', '科技'] }),
      tpl({ id: 'b', name: 'B', tags: ['深色'] }),
      tpl({ id: 'c', name: 'C', tags: ['国风'] }),
    ]
    expect(templateTags(list)).toEqual([
      { tag: '深色', count: 2 },
      { tag: '国风', count: 1 }, // 同次数按码点字典序：国(U+56FD) < 科(U+79D1)
      { tag: '科技', count: 1 },
    ])
  })

  it('单个模板内重复 tag 只计一次', () => {
    const list = [tpl({ id: 'a', name: 'A', tags: ['深色', '深色', '深色'] })]
    expect(templateTags(list)).toEqual([{ tag: '深色', count: 1 }])
  })

  it('跳过缺 tags / 空 tags / 空白 tag 的模板', () => {
    const list = [
      tpl({ id: 'a', name: 'A' }), // 缺 tags 字段
      tpl({ id: 'b', name: 'B', tags: [] }), // 空数组
      tpl({ id: 'c', name: 'C', tags: ['', '  '] }), // 全是空串/空白
    ]
    expect(templateTags(list)).toEqual([])
  })

  it('等价 tag（仅空白差异）合并计数', () => {
    const list = [tpl({ id: 'a', name: 'A', tags: ['深色', ' 深色 ', '深色'] })]
    expect(templateTags(list)).toEqual([{ tag: '深色', count: 1 }])
  })

  it('排序：次数降序，同次数按标签字典序（稳定可预期）', () => {
    const list = [
      tpl({ id: 'a', name: 'A', tags: ['国风', '深色', '极简'] }),
      tpl({ id: 'b', name: 'B', tags: ['深色'] }),
    ]
    const tags = templateTags(list).map((v) => v.tag)
    expect(tags[0]).toBe('深色') // 2 次
    expect(tags.slice(1)).toEqual(['国风', '极简']) // 各 1 次，字典序
  })

  it('空输入给空词表', () => {
    expect(templateTags([])).toEqual([])
  })
})

describe('matchesTemplate（单条匹配）', () => {
  const t = tpl({
    id: 'jw-dark-terminal',
    name: '暗色终端',
    description: '开发者风格的技术分享模板',
    tags: ['深色', '开发者'],
    scenario: ['技术分享', '开源项目介绍'],
  })

  it('空 query + 空 tag 放行一切', () => {
    expect(matchesTemplate(t, '', '')).toBe(true)
    expect(matchesTemplate(t, '   ', '')).toBe(true)
  })

  it('tag 命中 / 不命中', () => {
    expect(matchesTemplate(t, '', '深色')).toBe(true)
    expect(matchesTemplate(t, '', '开发者')).toBe(true)
    expect(matchesTemplate(t, '', '国风')).toBe(false)
  })

  it('query 命中名称', () => {
    expect(matchesTemplate(t, '终端', '')).toBe(true)
    expect(matchesTemplate(t, '暗色终端', '')).toBe(true)
  })

  it('query 命中 id', () => {
    expect(matchesTemplate(t, 'terminal', '')).toBe(true)
    expect(matchesTemplate(t, 'jw-dark', '')).toBe(true)
  })

  it('query 命中描述', () => {
    expect(matchesTemplate(t, '开发者风格', '')).toBe(true)
  })

  it('query 命中 tags 成员', () => {
    expect(matchesTemplate(t, '深色', '')).toBe(true)
    // 子串即可："深" 应命中 tag「深色」
    expect(matchesTemplate(t, '深', '')).toBe(true)
  })

  it('query 命中 scenario 成员', () => {
    expect(matchesTemplate(t, '技术分享', '')).toBe(true)
    expect(matchesTemplate(t, '开源', '')).toBe(true)
  })

  it('大小写不敏感（英文字段）', () => {
    expect(matchesTemplate(t, 'TERMINAL', '')).toBe(true)
    expect(matchesTemplate(t, 'JW', '')).toBe(true)
  })

  it('query 前后空白被 trim', () => {
    expect(matchesTemplate(t, '  终端  ', '')).toBe(true)
  })

  it('多词 query 不拆分、整串匹配（既定行为：搜不到就搜不到，不做分词）', () => {
    expect(matchesTemplate(t, '深色 国风', '')).toBe(false)
    expect(matchesTemplate(t, '暗色 终端', '')).toBe(false)
  })

  it('什么都碰不到的 query 返回 false', () => {
    expect(matchesTemplate(t, '不存在的东西', '')).toBe(false)
  })

  it('tag 与 query 同时给：取交集', () => {
    // 命中描述但 tag 不符 → 淘汰
    expect(matchesTemplate(t, '开发者风格', '国风')).toBe(false)
    // tag 符合 + query 命中 → 通过
    expect(matchesTemplate(t, '技术分享', '深色')).toBe(true)
  })

  it('没有 tags 字段的用户模板：tag 筛选下被排除，纯搜索仍可命中', () => {
    const ut = tpl({ id: 'ut-abc', name: '我的定制模板' })
    expect(matchesTemplate(ut, '', '深色')).toBe(false)
    expect(matchesTemplate(ut, '定制', '')).toBe(true)
    expect(matchesTemplate(ut, '', '')).toBe(true)
  })

  it('缺 description / scenario 字段不报错', () => {
    const bare = tpl({ id: 'x', name: '极简' })
    expect(matchesTemplate(bare, '极简', '')).toBe(true)
    expect(matchesTemplate(bare, '不存在', '')).toBe(false)
  })
})

describe('filterTemplates（列表过滤）', () => {
  const list = [
    tpl({ id: 'ut-mine', name: '我的模板', description: '' }),
    tpl({ id: 'a-dark', name: '暗色甲', tags: ['深色'], scenario: ['技术分享'] }),
    tpl({ id: 'b-warm', name: '暖色乙', tags: ['治愈'], scenario: ['教学科普'] }),
    tpl({ id: 'c-dark', name: '暗色丙', tags: ['深色', '国风'], description: '水墨质感' }),
  ]

  it('空条件等价于输入（同一批对象、同顺序）', () => {
    expect(filterTemplates(list, '', '')).toEqual(list)
  })

  it('保持输入顺序（上游「我的模板置顶」排序不被打乱）', () => {
    const out = filterTemplates(list, '暗色', '')
    expect(out.map((t) => t.id)).toEqual(['a-dark', 'c-dark'])
  })

  it('tag 筛选 + query 组合取交集', () => {
    expect(filterTemplates(list, '', '深色').map((t) => t.id)).toEqual(['a-dark', 'c-dark'])
    expect(filterTemplates(list, '水墨', '深色').map((t) => t.id)).toEqual(['c-dark'])
    expect(filterTemplates(list, '水墨', '治愈')).toEqual([])
  })

  it('零命中给空数组', () => {
    expect(filterTemplates(list, '查无此模板', '')).toEqual([])
  })

  it('用户模板（无 tags）只被 query 命中、不被 tag 命中', () => {
    expect(filterTemplates(list, '我的', '')).toEqual([list[0]])
    expect(filterTemplates(list, '', '治愈').map((t) => t.id)).toEqual(['b-warm'])
  })
})

describe('useTemplateFilter（组合式状态）', () => {
  it('toggleTag 选中，再点同一个取消；点别的切换', () => {
    const source = ref<FilterableTemplate[]>([tpl({ id: 'a', name: 'A', tags: ['深色'] })])
    const f = useTemplateFilter(source)
    expect(f.activeTag.value).toBe('')
    f.toggleTag('深色')
    expect(f.activeTag.value).toBe('深色')
    f.toggleTag('深色')
    expect(f.activeTag.value).toBe('')
    f.toggleTag('深色')
    f.toggleTag('国风')
    expect(f.activeTag.value).toBe('国风')
  })

  it('reset 同时清掉 query 和 activeTag；hasFilter 联动', () => {
    const source = ref<FilterableTemplate[]>([])
    const f = useTemplateFilter(source)
    expect(f.hasFilter.value).toBe(false)
    f.query.value = '终端'
    expect(f.hasFilter.value).toBe(true)
    f.query.value = '   ' // 纯空白不算筛选中
    expect(f.hasFilter.value).toBe(false)
    f.toggleTag('深色')
    f.query.value = 'x'
    expect(f.hasFilter.value).toBe(true)
    f.reset()
    expect(f.hasFilter.value).toBe(false)
    expect(f.query.value).toBe('')
    expect(f.activeTag.value).toBe('')
  })

  it('filtered 随条件变化实时重算', async () => {
    const source = ref<FilterableTemplate[]>([
      tpl({ id: 'a', name: '暗色终端', tags: ['深色'] }),
      tpl({ id: 'b', name: '暖色卡片', tags: ['治愈'] }),
    ])
    const f = useTemplateFilter(source)
    expect(f.filtered.value.map((t) => t.id)).toEqual(['a', 'b'])
    f.query.value = '终端'
    expect(f.filtered.value.map((t) => t.id)).toEqual(['a'])
    f.query.value = ''
    f.toggleTag('治愈')
    expect(f.filtered.value.map((t) => t.id)).toEqual(['b'])
    await nextTick()
    expect(f.filtered.value.map((t) => t.id)).toEqual(['b'])
  })

  it('source 是响应式 ref：增删条目后 vocab 与 filtered 跟着重算', () => {
    const source = ref<FilterableTemplate[]>([tpl({ id: 'a', name: 'A', tags: ['深色'] })])
    const f = useTemplateFilter(source)
    expect(f.vocab.value).toEqual([{ tag: '深色', count: 1 }])
    source.value = [
      tpl({ id: 'a', name: 'A', tags: ['深色'] }),
      tpl({ id: 'b', name: 'B', tags: ['深色', '国风'] }),
    ]
    expect(f.vocab.value).toEqual([
      { tag: '深色', count: 2 },
      { tag: '国风', count: 1 },
    ])
    expect(f.filtered.value).toHaveLength(2)
    source.value = []
    expect(f.vocab.value).toEqual([])
    expect(f.filtered.value).toEqual([])
  })

  it('source 用 getter 传 computed 也行（画廊的 cards 是 computed）', () => {
    const source = ref<FilterableTemplate[]>([tpl({ id: 'a', name: '甲', tags: ['深色'] })])
    const f = useTemplateFilter(() => source.value)
    f.query.value = '甲'
    expect(f.filtered.value).toHaveLength(1)
  })
})

describe.skipIf(!existsSync(TEMPLATES_DIR))('内置模板数据契约（backend/templates）', () => {
  it('每个内置模板的 template.json 都带非空 tags——筛选词表的数据源不能漏', () => {
    // templates/ 下混着 tools/ 之类的非模板目录：以「有 template.json」为准
    const dirs = readdirSync(TEMPLATES_DIR, { withFileTypes: true })
      .filter((d) => d.isDirectory() && existsSync(join(TEMPLATES_DIR, d.name, 'template.json')))
      .map((d) => d.name)
    // 目录存在但规模异常（比如 checkout 不完整）时别静默通过
    expect(dirs.length).toBeGreaterThan(50)

    const missing: string[] = []
    for (const id of dirs) {
      const meta = JSON.parse(readFileSync(join(TEMPLATES_DIR, id, 'template.json'), 'utf-8')) as {
        tags?: unknown
      }
      if (!Array.isArray(meta.tags) || meta.tags.length === 0) missing.push(id)
    }
    expect(missing).toEqual([])
  })
})
