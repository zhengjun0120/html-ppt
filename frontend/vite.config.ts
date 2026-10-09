import { fileURLToPath, URL } from 'node:url'

import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'

// dev 同源代理的后端地址。默认 8080；与其他分支的实例并行跑时用
// BACKEND_ADDR=http://127.0.0.1:8081 覆盖（配合后端 SERVER_ADDR）。
const backendAddr = process.env.BACKEND_ADDR || 'http://127.0.0.1:8080'

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  build: {
    // 构建产物挪到 /app/：/assets/ 是 deck-v2 运行时（翻页脚本/字体）的命名空间，
    // 生产 nginx 把 /assets/ 代理给后端——Vite 默认的 assets 目录会跟它撞车。
    assetsDir: 'app',
  },
  server: {
    port: 5173,
    proxy: {
      // dev 同源代理：免 CORS/凭据问题，SSE 流式透传
      '/api': { target: backendAddr, changeOrigin: false },
      '/assets': { target: backendAddr, changeOrigin: false },
      // deck-v2 模板库静态服务（画廊 live 预览 iframe 用）
      '/templates': { target: backendAddr, changeOrigin: false },
      // 用户自定义模板 demo 静态（定制工作台预览 iframe 用）
      '/user-templates': { target: backendAddr, changeOrigin: false },
    },
  },
})
