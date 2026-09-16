import { createPinia } from 'pinia'
import { createApp } from 'vue'

import { configureClient } from './api/client'
import App from './App.vue'
import router from './router'
import { TOKEN_KEY, useAuthStore } from './stores/auth'
import './styles/main.css'

const app = createApp(App)
app.use(createPinia())
app.use(router)

// 401 全局出口：清 token 回登录页（auth store 的 logout 走同一入口）
const auth = useAuthStore()
configureClient({
  getToken: () => localStorage.getItem(TOKEN_KEY),
  onUnauthorized: () => {
    localStorage.removeItem(TOKEN_KEY)
    void router.push('/login')
  },
})

app.mount('#app')

// 首帧不被阻塞：路由就绪后静默续期（当天已刷新则 429 静默忽略）
void router.isReady().then(() => auth.silentRefresh())
