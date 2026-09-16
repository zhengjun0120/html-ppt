import { request } from './client'

export interface DeckMeta {
  id: string
  title: string
}

export interface VersionMeta {
  version: string
  time: number
  operation: 'run' | 'restore'
  detail: string
  slides: number
}

export const deckApi = {
  list: () => request<DeckMeta[]>('/api/decks'),
  history: (deckId: string) => request<VersionMeta[]>(`/api/decks/${deckId}/history`),
  restore: (deckId: string, version: string) =>
    request<{ restored: string }>(`/api/decks/${deckId}/history/${version}/restore`, { method: 'POST' }),
  deleteVersion: (deckId: string, version: string) =>
    request<Record<string, never>>(`/api/decks/${deckId}/history/${version}`, { method: 'DELETE' }),
  clearHistory: (deckId: string) =>
    request<{ deleted: number }>(`/api/decks/${deckId}/history`, { method: 'DELETE' }),
}
