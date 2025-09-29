<script setup lang="ts">
import { initPowerXBridge } from '~/bridge/powerx-bridge-client'

const locale = ref<string>('-')
const theme  = ref<string>('-')

const applyLocale = (loc: string) => {
  locale.value = loc
  document.documentElement.setAttribute('data-locale', loc)
}
const applyTheme = (thm: string) => {
  theme.value = thm
  document.documentElement.setAttribute('data-theme', thm)
}

onMounted(() => {
  const bridge = initPowerXBridge({
    debug: true,
    pluginId: 'com.demo.plugin',
    instanceId: 'dev-1',
    onLocale: applyLocale,
    onTheme: applyTheme,
    onSync: ({ locale, theme }) => { applyLocale(locale); applyTheme(theme) }
  })
  // 暴露便于控制台调试
  // @ts-expect-error
  window.requestSync = () => bridge.requestSync()
  // @ts-expect-error
  window.ping = () => bridge.ping()
})
</script>

<template>
  <div class="wrap">
    <h3>Plugin (iframe)</h3>
    <div class="badges">
      <span class="badge">locale: <b>{{ locale }}</b></span>
      <span class="badge">theme:  <b>{{ theme }}</b></span>
    </div>
    <p>可在控制台执行 <code>requestSync()</code> / <code>ping()</code>。</p>
  </div>
</template>

<style scoped>
.wrap { padding: 16px; font-family: Inter, system-ui, -apple-system, Segoe UI, Roboto, sans-serif; }
.badges { margin: 8px 0 12px; display:flex; gap:8px; }
.badge { display:inline-block; padding:6px 10px; border-radius:8px; background: #eee; }
:root { --bg:#fff; --fg:#111; --chip:#eee; }
:global([data-theme="dark"]) .wrap { background:#111; color:#eee; }
</style>
