import { fileURLToPath, URL } from 'node:url'

import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    port: 5173,
    proxy: {
      // dev 同源代理：免 CORS/凭据问题，SSE 流式透传
      '/api': { target: 'http://127.0.0.1:8080', changeOrigin: false },
      '/assets': { target: 'http://127.0.0.1:8080', changeOrigin: false },
      // deck-v2 模板库静态服务（画廊 live 预览 iframe 用）
      '/templates': { target: 'http://127.0.0.1:8080', changeOrigin: false },
    },
  },
})
