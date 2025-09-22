<template>
  <div class="p-6 space-y-6">
    <header class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
          {{ t('navigation.reportsDaily') }}
        </h1>
        <p class="text-gray-600 dark:text-gray-400 mt-1">
          {{ t('reports.daily.overviewHint') }}
        </p>
      </div>
    </header>

    <UCard>
      <div class="text-gray-600 dark:text-gray-300">
        {{ t('reports.placeholder') }}
      </div>
    </UCard>
  </div>
</template>

<script setup lang="ts">
const { t } = useI18n();

definePageMeta({
  title: 'reportsDaily',
  layout: 'default',
});
</script>
