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
    // :tab? = 二级导航（mine/community/builtin），缺省与非法值在视图内回退「我的模板」
    { path: '/my-templates/:tab?', component: () => import('@/views/MyTemplatesView.vue') },
    { path: '/my-templates/:id/edit', component: () => import('@/views/TemplateEditView.vue') },
    { path: '/decks', component: () => import('@/views/DecksView.vue') },
    { path: '/decks/new', redirect: '/new' },
    { path: '/decks/:id', name: 'deck', component: () => import('@/views/DeckView.vue') },
    { path: '/settings', component: () => import('@/views/SettingsView.vue') },
    { path: '/trace', component: () => import('@/views/TraceView.vue') },
    { path: '/usage', component: () => import('@/views/UsageView.vue') },
    { path: '/', redirect: '/decks' },
    { path: '/:pathMatch(.*)*', redirect: '/decks' },
  ],
  // 滚动恢复：浏览器前进/后退回到原生滚动位置（文稿列表返回时落回原处）；
  // 新导航回顶；同路径仅 query 变化（列表翻页）不劫持滚动，交给 setPage 自己平滑回顶。
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) return savedPosition
    if (to.path === from.path) return false
    return { top: 0 }
  },
})

router.beforeEach((to) => {
  const authed = !!localStorage.getItem('da-token')
  if (to.meta.public !== true && !authed) {
    return { path: '/login', query: to.fullPath === '/decks' ? undefined : { next: to.fullPath } }
  }
  if (to.meta.public === true && authed) return { path: '/decks' }
})

export default router
