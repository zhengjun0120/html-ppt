import type { TplSuggestion } from '@/api/deckV2'

/** 模板推荐的纯逻辑（TemplateGallery 用）：推荐条目 → 画廊卡片的映射决策。
 *  抽出来是为了 vitest 可测——组件里只剩 DOM 定位与状态机。 */

/** 画廊卡片的最小结构约束（GalleryCard 的形状子集） */
export interface SuggestableCard {
  id: string
  variants: { id: string }[]
}

/** resolveSuggestion 按 template_id 找推荐对应的画廊卡片。
 *  找不到（模板被删/候选外 id 漏网）返回 null——调用方静默忽略该条推荐。 */
export function resolveSuggestion<T extends SuggestableCard>(cards: T[], s: TplSuggestion): T | null {
  return cards.find((c) => c.id === s.template_id) ?? null
}

/** variantFor 推荐里给的变体是否真实存在于卡片；不存在返回 null，
 *  调用方回退到卡片自己的默认变体逻辑（pick() 的首变体/default）。 */
export function variantFor(card: SuggestableCard, s: TplSuggestion): string | null {
  if (!s.variant_id) return null
  return card.variants.some((v) => v.id === s.variant_id) ? s.variant_id : null
}
