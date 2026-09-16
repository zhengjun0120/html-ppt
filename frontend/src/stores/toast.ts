import { defineStore } from 'pinia'

export type ToastKind = 'success' | 'error' | 'info' | 'warning'

export interface ToastItem {
  id: number
  kind: ToastKind
  text: string
}

let seq = 0

/** 瞬态通知（DESIGN.md §5：toast 只用于瞬态，表单错误必须内联） */
export const useToastStore = defineStore('toast', {
  state: () => ({ items: [] as ToastItem[] }),
  actions: {
    push(kind: ToastKind, text: string) {
      const id = ++seq
      this.items.push({ id, kind, text })
      setTimeout(() => this.dismiss(id), 4200)
    },
    dismiss(id: number) {
      this.items = this.items.filter((t) => t.id !== id)
    },
  },
})

export function useToast() {
  const s = useToastStore()
  return {
    success: (text: string) => s.push('success', text),
    error: (text: string) => s.push('error', text),
    info: (text: string) => s.push('info', text),
    warning: (text: string) => s.push('warning', text),
  }
}
