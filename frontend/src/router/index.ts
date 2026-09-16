import { createRouter, createWebHistory } from 'vue-router'

declare module 'vue-router' {
  interface RouteMeta {
    /** 公开页（登录/注册）：无 token 可访问；有 token 访问时跳去 /decks */
    public?: boolean
  }
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/login', component: () => import('@/views/LoginView.vue'), meta: { public: true } },
    { path: '/register', component: () => import('@/views/RegisterView.vue'), meta: { public: true } },
    { path: '/decks', component: () => import('@/views/DecksView.vue') },
    { path: '/decks/:id', component: () => import('@/views/WorkspaceView.vue') },
    { path: '/settings', component: () => import('@/views/SettingsView.vue') },
    { path: '/trace', component: () => import('@/views/TraceView.vue') },
    { path: '/', redirect: '/decks' },
    { path: '/:pathMatch(.*)*', redirect: '/decks' },
  ],
})

router.beforeEach((to) => {
  const authed = !!localStorage.getItem('da-token')
  if (to.meta.public !== true && !authed) {
    return { path: '/login', query: to.fullPath === '/decks' ? undefined : { next: to.fullPath } }
  }
  if (to.meta.public === true && authed) return { path: '/decks' }
})

export default router
