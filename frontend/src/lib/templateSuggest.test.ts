import { describe, expect, it } from 'vitest'

import type { TplSuggestion } from '@/api/deckV2'
import { resolveSuggestion, variantFor, type SuggestableCard } from './templateSuggest'

const cards: SuggestableCard[] = [
  { id: 'tech-sharing', variants: [{ id: 'default' }, { id: 'blue' }, { id: 'ember' }] },
  { id: 'minimal-white', variants: [{ id: 'default' }] },
]

const sug = (partial: Partial<TplSuggestion>): TplSuggestion => ({
  template_id: 'tech-sharing',
  reason: '科技感契合',
  ...partial,
})

describe('resolveSuggestion', () => {
  it('按 template_id 找到卡片', () => {
    expect(resolveSuggestion(cards, sug({}))?.id).toBe('tech-sharing')
  })

  it('推荐了不存在的模板返回 null（组件静默忽略）', () => {
    expect(resolveSuggestion(cards, sug({ template_id: 'deleted-tpl' }))).toBeNull()
  })

  it('空卡片列表返回 null', () => {
    expect(resolveSuggestion([], sug({}))).toBeNull()
  })

  it('ut- 前缀的用户模板不在画廊候选时返回 null', () => {
    expect(resolveSuggestion(cards, sug({ template_id: 'ut-mine' }))).toBeNull()
  })
})

describe('variantFor', () => {
  it('推荐变体存在时原样返回', () => {
    expect(variantFor(cards[0], sug({ variant_id: 'blue' }))).toBe('blue')
  })

  it('推荐变体不存在返回 null（回退默认变体）', () => {
    expect(variantFor(cards[0], sug({ variant_id: 'no-such' }))).toBeNull()
  })

  it('没给变体返回 null', () => {
    expect(variantFor(cards[0], sug({}))).toBeNull()
  })

  it('单变体模板推荐 default 也认', () => {
    expect(variantFor(cards[1], sug({ template_id: 'minimal-white', variant_id: 'default' }))).toBe('default')
  })

  it('变体属于别的模板不算数', () => {
    // ember 只在 tech-sharing 上；对 minimal-white 是非法变体
    expect(variantFor(cards[1], sug({ template_id: 'minimal-white', variant_id: 'ember' }))).toBeNull()
  })
})
