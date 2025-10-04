import { useColorMode } from '@vueuse/core'
import { useI18n } from '#imports'

export const useTheme = () => {
  const colorMode = useColorMode() // ref<'light'|'dark'|'auto'>

  // 只读：当前主题
  const theme = computed(() => colorMode.value)

  // 仅列出 light/dark；如需“跟随系统”，仍然用 colorMode 的 'auto'
  const themes = [
    { value: 'light', label: '浅色主题' },
    { value: 'dark',  label: '深色主题' },
  ]

  // 设置主题：统一改 .value；支持传入 'light' | 'dark' | 'auto' | 'system'
  const setTheme = (newTheme: string) => {
    const v = newTheme === 'system' ? 'auto' : newTheme
    colorMode.value = v as any
    // 可选：同步到 DOM，便于样式选择器
    const host = v === 'auto' ? (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light') : v
    document.documentElement.dataset.theme = host // 'light' | 'dark'

  }

  // 初始化主题：只设置 colorMode.value；不要写 theme.value（它是 computed）
  const initTheme = () => {
    if (import.meta.client) {
      const runtimeConfig = useRuntimeConfig()
      const envTheme = runtimeConfig.public.defaultTheme
      const savedTheme = localStorage.getItem('theme')
      const initialTheme = (envTheme || savedTheme || 'light') as string
      setTheme(initialTheme) // 内部会统一落到 colorMode.value / data-theme
    }
  }

  // （可选）兼容你旧的 window message 方案，若已改为 bridge，可保留或删除
  const listenToParentMessages = () => {
    if (import.meta.client) {
      window.addEventListener('message', (event) => {
        if (event.data?.type === 'THEME_CHANGE') setTheme(event.data.theme)
        if (event.data?.type === 'LANGUAGE_CHANGE') {
          const { setLocale } = useI18n()
          setLocale(event.data.locale)
        }
      })
    }
  }

  return {
    theme: readonly(theme),
    themes,
    setTheme,
    initTheme,
    listenToParentMessages,
  }
}
