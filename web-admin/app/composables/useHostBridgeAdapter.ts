// app/composables/useHostBridgeAdapter.ts
import { initPowerXBridge } from '~/bridge/powerx-bridge-client'
import { useI18n } from '#imports'
import { useTheme } from '~/composables/useTheme' // 你已有的封装

type BridgeOptions = {
  pluginId?: string
  instanceId?: string
  debug?: boolean
}

/** 将宿主广播适配到项目内现有的语言/主题切换实现 */
export function setupHostBridgeAdapter(opts: BridgeOptions = {}) {
  const { setLocale, locale } = useI18n()
  const { setTheme, currentTheme } = useTheme()
  // ↑ 假设 useTheme.ts 暴露了 setTheme('light'|'dark'|'system'|'auto')、currentTheme()

  // 宿主 -> 本地：'system' 与我们内部 'auto' 的对齐
  const fromHostTheme = (t: string) => (t === 'system' ? 'auto' : t)
  const toHostTheme   = (t: string) => (t === 'auto'   ? 'system' : t)

  const applyLocale = async (code: string) => {
    // 与 LanguageSelector.vue 的行为保持一致
    if (!code || code === String(locale.value)) return
    await setLocale(code)
    // 这里不主动 broadcast，避免回环；交给宿主侧做“源头广播”
  }

  const applyTheme = (t: string) => {
    const local = fromHostTheme(t)
    // 与 ThemeSelector.vue / useTheme.ts 的实现保持一致
    setTheme(local as any) // 由你的 useTheme 内部处理 class/data-attr/CSS 变量
  }

  const bridge = initPowerXBridge({
    debug: opts.debug ?? true,
    pluginId: opts.pluginId ?? 'com.demo.plugin',
    instanceId: opts.instanceId ?? 'dev-1',
    onLocale: (code) => { applyLocale(code) },
    onTheme:  (t)    => { applyTheme(t) },
    onSync:   ({ locale, theme }) => {
      applyLocale(locale)
      applyTheme(theme)
    }
  })

  return {
    bridge,
    /** 提供给你在本地切换后（比如 LanguageSelector/ThemeSelector 中）可选回传宿主用 */
    broadcastThemeToHost() {
      const hostTheme = toHostTheme(String(currentTheme()))
      bridge && window?.parent?.postMessage?.({ source: 'plugin', type: 'ping', ts: Date.now() }, '*')
      // 如果需要可上报当前 theme 给宿主，这里仅演示心跳；通常由宿主做单向广播即可
      return hostTheme
    }
  }
}
