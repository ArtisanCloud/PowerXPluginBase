<template>
  <aside
    class="w-64 min-w-64 max-w-64 bg-white dark:bg-gray-900 border-r border-gray-200 dark:border-gray-800 min-h-screen flex-shrink-0"
  >
    <nav class="p-4 space-y-6">
      <!-- 仪表盘入口 -->
      <div>
        <UButton
          to="/dashboard"
          variant="ghost"
          color="neutral"
          class="w-full justify-start"
          :class="{
            'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
              isExactActive('/dashboard'),
          }"
        >
          <UIcon name="i-heroicons-home" class="w-4 h-4 mr-3" />
          {{ t('navigation.dashboard') }}
        </UButton>
      </div>

      <!-- 笔记模块 -->
      <div>
        <div
          class="px-3 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider"
        >
          {{ t('navigation.notes') }}
        </div>
        <div class="space-y-1">
          <UButton
            to="/notes"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                isExactActive('/notes'),
            }"
          >
            <UIcon name="i-heroicons-document-text" class="w-4 h-4 mr-3" />
            {{ t('notes.overview') }}
          </UButton>

          <UButton
            to="/notes/active"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                isExactActive('/notes/active'),
            }"
          >
            <UIcon name="i-heroicons-play-circle" class="w-4 h-4 mr-3" />
            {{ t('navigation.activenote') }}
          </UButton>

          <UButton
            to="/notes/create"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                isExactActive('/notes/create'),
            }"
          >
            <UIcon name="i-heroicons-plus-circle" class="w-4 h-4 mr-3" />
            {{ t('notes.createNote') }}
          </UButton>
        </div>
      </div>

      <!-- 报表模块 -->
      <div>
        <div
          class="px-3 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider"
        >
          {{ t('navigation.reports') }}
        </div>
        <div class="space-y-1">
          <UButton
            to="/reports/daily"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                isExactActive('/reports/daily'),
            }"
          >
            <UIcon name="i-heroicons-calendar-days" class="w-4 h-4 mr-3" />
            {{ t('navigation.reportsDaily') }}
          </UButton>

          <UButton
            to="/reports/monthly"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                isExactActive('/reports/monthly'),
            }"
          >
            <UIcon name="i-heroicons-calendar" class="w-4 h-4 mr-3" />
            {{ t('navigation.reportsMonthly') }}
          </UButton>
        </div>
      </div>
    </nav>
  </aside>
</template>

<script setup lang="ts">
const { t } = useI18n();
const route = useRoute();

const normalizePath = (value: string) => {
  if (!value) {
    return '/';
  }
  if (value !== '/' && value.endsWith('/')) {
    return value.replace(/\/+$/, '');
  }
  return value;
};

const isExactActive = (target: string) => {
  const current = normalizePath(route.path);
  const normalizedTarget = normalizePath(target);

  if (normalizedTarget === '/dashboard') {
    return current === '/dashboard' || current === '/';
  }

  return current === normalizedTarget;
};
</script>
