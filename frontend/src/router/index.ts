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
    // 向导五步（plan-v3 C1）：每步一个 URL，可刷新可回退。
    // WizardView 内部做步骤守卫（不允许跳到尚未到达的步骤；stage 推进时自动前跳）。
    { path: '/new', name: 'wizard-clarify', component: () => import('@/views/wizard/WizardView.vue') },
    { path: '/new/outline', name: 'wizard-outline', component: () => import('@/views/wizard/WizardView.vue') },
    { path: '/new/template', name: 'wizard-template', component: () => import('@/views/wizard/WizardView.vue') },
    { path: '/new/generating', name: 'wizard-generate', component: () => import('@/views/wizard/WizardView.vue') },
    // 第 5 步 · 成品预览 + 迭代（旧工作台瘦身版）
    { path: '/my-templates', component: () => import('@/views/MyTemplatesView.vue') },
    { path: '/my-templates/:id/edit', component: () => import('@/views/TemplateEditView.vue') },
    { path: '/decks', component: () => import('@/views/DecksView.vue') },
    { path: '/decks/new', redirect: '/new' },
    { path: '/decks/:id', name: 'deck', component: () => import('@/views/DeckView.vue') },
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
