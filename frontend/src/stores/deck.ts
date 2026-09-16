import { defineStore } from 'pinia'

import { deckApi, type DeckMeta } from '@/api/decks'

export const useDeckStore = defineStore('deck', {
  state: () => ({
    list: [] as DeckMeta[],
    loaded: false,
    loading: false,
  }),
  getters: {
    titleOf: (s) => (id: string) => s.list.find((d) => d.id === id)?.title ?? id,
  },
  actions: {
    async ensureList() {
      if (this.loaded) return
      await this.refresh()
    },
    async refresh() {
      this.loading = true
      try {
        this.list = await deckApi.list()
        this.loaded = true
      } finally {
        this.loading = false
      }
    },
  },
})
