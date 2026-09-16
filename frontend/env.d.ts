/// <reference types="vite/client" />

interface ImportMetaEnv {
  /** 后端 API 源；dev 为空串（走 vite 同源代理），生产填后端绝对地址 */
  readonly VITE_API_BASE: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
