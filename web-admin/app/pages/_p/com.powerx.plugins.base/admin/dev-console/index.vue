<template>
  <div class="space-y-6">
    <header class="space-y-2">
      <div class="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
        <UIcon name="i-heroicons-command-line" class="text-primary" />
        <span class="uppercase tracking-wide">Dev Console</span>
      </div>
      <div>
        <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">配置与治理</h1>
        <p class="text-sm text-gray-600 dark:text-gray-300">管理 Dev Console 的默认参数、导出策略和安全操作。</p>
      </div>
    </header>

    <UAlert
      v-if="store.error"
      color="red"
      variant="soft"
      :title="store.error"
    />

    <USkeleton v-if="store.loading" class="h-32 w-full" repeat="3" />

    <div v-else class="space-y-6">
      <ConfigSectionCard
        v-for="section in store.sections"
        :key="section.key"
        :section="section"
        :pending="savingKey === section.key"
        @submit="values => saveSection(section.key, values)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useDevConsoleConfigStore } from '~/app/stores/dev-console/config'
import ConfigSectionCard from '~/app/components/dev-console/ConfigSectionCard.vue'

const store = useDevConsoleConfigStore()
const { sections } = storeToRefs(store)
const toast = useToast()
const savingKey = ref('')

onMounted(async () => {
  if (!sections.value.length) {
    await store.fetchSections().catch(() => {})
  }
})

async function saveSection(key: string, values: Record<string, any>) {
  savingKey.value = key
  try {
    await store.updateSection(key, { values })
    toast.add({ title: '配置已更新', color: 'green' })
  } catch (error: any) {
    toast.add({ title: '保存失败', description: error?.message ?? '请稍后再试', color: 'red' })
  } finally {
    savingKey.value = ''
  }
}
</script>
