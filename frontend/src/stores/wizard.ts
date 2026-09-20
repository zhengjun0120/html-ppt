import { defineStore } from 'pinia'

import { deckV2Api, OutlineConflictError, type DeckFile, type DeckStage, type Outline } from '@/api/deckV2'
import { templateApi, type TemplateMeta, type TemplateVariant } from '@/api/templates'

/**
 * deck-v2 向导状态：阶段、大纲、模板。
 *
 * 数据流向：chat 流里的事件（outline_updated / stage / gate_waiting）负责把
 * deckId 送进 chat store；本 store 在 deckId 变化时拉取 meta 作为权威状态，
 * REST 过渡（确认大纲/选模板）成功后本地同步更新 + 重拉。
 */

export type WizardStep = 'clarify' | 'outline' | 'template' | 'generate' | 'iterate'

const STAGE_TO_STEP: Record<DeckStage, WizardStep> = {
  draft: 'clarify',
  outlining: 'outline',
  outline_review: 'outline',
  selecting_template: 'template',
  generating: 'generate',
  iterating: 'iterate',
}

export const STEP_LABELS: { key: WizardStep; label: string }[] = [
  { key: 'clarify', label: '澄清' },
  { key: 'outline', label: '大纲' },
  { key: 'template', label: '模板' },
  { key: 'generate', label: '生成' },
  { key: 'iterate', label: '迭代' },
]

/** 各步骤的独立路由（plan-v3 C1：五步五页） */
export const STEP_ROUTES: Record<WizardStep, (deckId: string) => string> = {
  clarify: () => '/new',
  outline: () => '/new/outline',
  template: () => '/new/template',
  generate: () => '/new/generating',
  iterate: (deckId) => `/decks/${deckId}`,
}

/** 步骤先后序（守卫用：不允许跳到尚未到达的步骤） */
export function stepOrder(key: WizardStep): number {
  return STEP_LABELS.findIndex((s) => s.key === key)
}

export const useWizardStore = defineStore('wizard', {
  state: () => ({
    deckId: '',
    /** 空 = 还没有 v2 deck（冷启动阶段）；'v1' 表示旧文稿（向导不接管主区域） */
    stage: '' as '' | DeckStage,
    format: '',
    title: '',
    canvas: null as { w: number; h: number } | null,
    templateId: '',
    variant: '',
    variants: [] as TemplateVariant[],
    outline: null as Outline | null,
    templates: [] as TemplateMeta[],
    templatesLoaded: false,
    /** 正在保存大纲 / 确认 / 选模板（按钮禁用态） */
    busy: false,
  }),

  getters: {
    isV2: (s) => s.format === 'v2',
    // 跨 getter 引用必须用 this（store 实例）：箭头写法 s.isV2 取的是 state，
    // 恒为 undefined——会让整台阶段机哑火（step 恒 null，向导永远不跳页）。
    step(): WizardStep | null {
      if (!this.isV2 || !this.stage) return null
      if (this.stage === 'outlining' || this.stage === 'outline_review') {
        // 有大纲可看 = 大纲步骤；还在聊需求 = 澄清步骤
        return this.outline ? 'outline' : 'clarify'
      }
      return STAGE_TO_STEP[this.stage] ?? null
    },
    /** 向导接管输入框：选模板期间聊天锁定（D 系决策 R3：这一步在对话里做不了） */
    locksInput(): boolean {
      return this.isV2 && this.stage === 'selecting_template'
    },
    /** 主区域显示大纲面板 */
    showOutline(): boolean {
      return this.isV2 && this.stage === 'outline_review'
    },
    /** 主区域显示模板画廊 */
    showGallery(): boolean {
      return this.isV2 && this.stage === 'selecting_template'
    },
    /** 主区域显示生成进度 */
    showGenerating(): boolean {
      return this.isV2 && this.stage === 'generating'
    },
    defaultVariantId: (s) => s.variants[0]?.id ?? 'default',
  },

  actions: {
    reset() {
      this.deckId = ''
      this.stage = ''
      this.format = ''
      this.title = ''
      this.canvas = null
      this.templateId = ''
      this.variant = ''
      this.variants = []
      this.outline = null
      this.busy = false
    },

    /** chat.deckId 变化 / 会话恢复后调用：拉元数据同步阶段 */
    async syncFromChat(deckId: string) {
      if (!deckId) {
        // 无 deck 上下文 = 冷启动/新建意图：必须清掉残留状态。store 是全局单例，
        // 上一个文稿的阶段（比如停在选模板）会泄漏给新会话——锁输入框、
        // stepper 高亮旧阶段、守卫把人带进旧文稿的步骤页。
        this.reset()
        return
      }
      if (deckId === this.deckId) {
        // 同一 deck：refresh 但不折腾（阶段由后端守卫，本地只展示）
        await this.refresh()
        return
      }
      this.reset()
      await this.loadMeta(deckId)
    },

    async loadMeta(deckId: string) {
      const meta = await deckV2Api.meta(deckId)
      this.applyMeta(meta, deckId)
    },

    applyMeta(meta: { deck: DeckFile; outline?: Outline; variants?: TemplateVariant[] }, deckId?: string) {
      this.deckId = deckId ?? meta.deck.id
      this.stage = meta.deck.stage
      this.format = meta.deck.format
      this.title = meta.deck.title
      this.templateId = meta.deck.template_id ?? ''
      this.variant = meta.deck.variant ?? ''
      this.canvas = meta.deck.canvas ?? null
      this.variants = meta.variants ?? []
      this.outline = meta.outline ?? null
    },

    async refresh() {
      if (!this.deckId) return
      try {
        await this.loadMeta(this.deckId)
      } catch { /* 元数据刷新失败不阻塞（deck 可能刚被删） */ }
    },

    async loadTemplates(force = false) {
      if (this.templatesLoaded && !force) return
      this.templates = await templateApi.list()
      this.templatesLoaded = true
    },

    /** 大纲面板保存：整份替换。版本冲突抛 OutlineConflictError 给组件提示 */
    async saveOutline(outline: Outline): Promise<{ version: number }> {
      if (!this.deckId) throw new Error('deck 未就绪')
      this.busy = true
      try {
        const r = await deckV2Api.putOutline(this.deckId, outline.version, outline)
        if (this.outline) this.outline.version = r.version
        return r
      } catch (e) {
        if (e instanceof OutlineConflictError) {
          // 冲突即以服务端为准重拉，组件提示"已加载最新版"
          await this.refresh()
        }
        throw e
      } finally {
        this.busy = false
      }
    },

    /** gate 1：确认大纲（outline_review → selecting_template），顺带加载模板清单 */
    async confirmOutline() {
      if (!this.deckId) return
      this.busy = true
      try {
        const r = await deckV2Api.confirmOutline(this.deckId)
        this.stage = r.stage
        await this.loadTemplates()
      } finally {
        this.busy = false
      }
    },

    /** gate 2：选模板（实例化）。返回后由调用方触发 chat.runGeneration */
    async selectTemplate(templateId: string, variant: string): Promise<DeckStage> {
      if (!this.deckId) throw new Error('deck 未就绪')
      this.busy = true
      try {
        const r = await deckV2Api.selectTemplate(this.deckId, templateId, variant || 'default')
        this.stage = r.stage
        this.templateId = templateId
        this.variant = variant || 'default'
        const tpl = this.templates.find((t) => t.id === templateId)
        if (tpl) this.canvas = tpl.canvas
        return r.stage
      } finally {
        this.busy = false
      }
    },
  },
})
