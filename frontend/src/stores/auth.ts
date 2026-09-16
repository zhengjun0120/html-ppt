import { defineStore } from 'pinia'

import { authApi } from '@/api/auth'
import { ApiError } from '@/api/client'
import router from '@/router'

export const TOKEN_KEY = 'da-token'

/** 静默续期的豁免名单：这些路径上的 401 不再重复跳转（守卫已处理） */
function isPublicRoute(): boolean {
  return router.currentRoute.value.meta.public === true
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    email: '',
    hasKey: false,
  }),
  actions: {
    saveToken(t: string) {
      localStorage.setItem(TOKEN_KEY, t)
    },
    clearToken() {
      localStorage.removeItem(TOKEN_KEY)
      this.email = ''
      this.hasKey = false
    },
    async login(email: string, password: string) {
      const r = await authApi.login(email, password)
      this.saveToken(r.token)
    },
    async register(email: string, code: string, password: string) {
      const r = await authApi.register(email, code, password)
      this.saveToken(r.token)
    },
    async fetchMe() {
      const r = await authApi.me()
      this.email = r.email
      this.hasKey = r.has_api_key
    },
    async saveApiKey(key: string) {
      await authApi.setApiKey(key)
      await this.fetchMe()
    },
    logout() {
      this.clearToken()
      void router.push('/login')
    },
    /**
     * 每天首次上线静默换发新 token（滑续 30 天）。
     * 429 = 今天已经刷过（本设备或别的设备），继续用旧 token，不是错误；
     * 401 = token 已过期，清掉让守卫送去登录。
     */
    async silentRefresh() {
      if (!localStorage.getItem(TOKEN_KEY)) return
      try {
        const r = await authApi.refresh()
        this.saveToken(r.token)
      } catch (e) {
        if (e instanceof ApiError && e.status === 401 && !isPublicRoute()) {
          this.logout()
        }
        // 429 与网络错误：静默忽略
      }
    },
  },
})
