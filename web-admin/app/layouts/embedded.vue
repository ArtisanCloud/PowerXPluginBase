<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from '#imports'
import { setupHostBridgeAdapter } from '~/composables/useHostBridgeAdapter'
// 如果你的主题需要在挂载时同步一次到 DOM（data-theme），可引入：
import { useTheme } from '~/composables/useTheme'

const route = useRoute()

onMounted(() => {
  if (!process.client) return
  // 防重复初始化（HMR/切换页面时）
  if ((window as any).__PX_ADAPTER_BOUND__) return

  // 可选：先把当前主题同步一遍（避免首次渲染颜色不一致）
  try { useTheme().initTheme?.() } catch {}

  // 从路由里推断 pluginId/instanceId（按你需要调整）
  const pluginId = (route.query.pluginId as string) || 'com.powerx.plugin'
  const instanceId = (route.query.instanceId as string) || route.fullPath

  const { bridge } = setupHostBridgeAdapter({
      // debug: import.meta.dev,
      debug: true,
      pluginId,
      instanceId,
    })

  // 确认启动
  bridge.start?.()
  console.info('[Bridge][Plugin] adapter mounted.')

    // 标记已初始化 + 暴露便于调试
  ;(window as any).__PX_ADAPTER_BOUND__ = true
  ;(window as any).__PX_ADAPTER__ = {bridge}
  console.info('[embedded] Host bridge adapter mounted.', { pluginId, instanceId, ret })
})
</script>

<template>
  <!-- 保持布局极简，不要额外外壳，slot 里的页面就是“插件视图” -->
  <div class="embedded-wrap">
    <slot />
  </div>
</template>

<style scoped>
.embedded-wrap { min-height: 100dvh; }
</style>
