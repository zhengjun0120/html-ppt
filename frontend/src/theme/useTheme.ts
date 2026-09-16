import { ref, watchEffect } from 'vue'

export type Theme = 'light' | 'dark'

const THEME_KEY = 'da-theme'
const theme = ref<Theme>(
  document.documentElement.dataset.theme === 'light' ? 'light' : 'dark',
)

function apply(t: Theme) {
  document.documentElement.dataset.theme = t
  theme.value = t
}

/** 手动选择持久化；一旦手动选择，不再跟随系统 */
function set(t: Theme) {
  localStorage.setItem(THEME_KEY, t)
  apply(t)
}

function toggle() {
  set(theme.value === 'dark' ? 'light' : 'dark')
}

// 用户从未手动选择过时，跟随系统变化
if (!localStorage.getItem(THEME_KEY)) {
  const media = matchMedia('(prefers-color-scheme: light)')
  media.addEventListener('change', (e) => {
    apply(e.matches ? 'light' : 'dark')
  })
}

watchEffect(() => {
  document.documentElement.dataset.theme = theme.value
})

export function useTheme() {
  return { theme, set, toggle }
}
