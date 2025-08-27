<template>
  <div>
    <!-- 主要内容 -->
    <main class="space-y-8">
      <!-- 页面标题 -->
      <div>
        <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">
          {{ $t("navigation.dashboard") }}
        </h1>
        <p class="mt-2 text-sm text-gray-600 dark:text-gray-400">
          {{ $t("dashboard.welcome") }}
        </p>
      </div>

      <!-- 仪表板概览 -->
      <div>
        <h2 class="text-lg font-medium text-gray-900 dark:text-white mb-4">
          {{ $t("dashboard.projectOverview") }}
        </h2>
        <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
          <!-- 待办任务统计卡片 -->
          <UCard class="hover:shadow-md transition-shadow">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <UIcon name="i-heroicons-list-bullet" class="h-6 w-6 text-gray-400" />
              </div>
              <div class="ml-5 w-0 flex-1">
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
                  {{ $t("dashboard.todoTasks") }}
                </dt>
                <dd class="text-lg font-medium text-gray-900 dark:text-white">12</dd>
              </div>
            </div>
          </UCard>

          <!-- 进行中统计卡片 -->
          <UCard class="hover:shadow-md transition-shadow">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <UIcon name="i-heroicons-clock" class="h-6 w-6 text-blue-500" />
              </div>
              <div class="ml-5 w-0 flex-1">
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
                  {{ $t("status.inProgress") }}
                </dt>
                <dd class="text-lg font-medium text-gray-900 dark:text-white">8</dd>
              </div>
            </div>
          </UCard>

          <!-- 已完成统计卡片 -->
          <UCard class="hover:shadow-md transition-shadow">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <UIcon name="i-heroicons-check-circle" class="h-6 w-6 text-green-500" />
              </div>
              <div class="ml-5 w-0 flex-1">
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
                  {{ $t("status.done") }}
                </dt>
                <dd class="text-lg font-medium text-gray-900 dark:text-white">24</dd>
              </div>
            </div>
          </UCard>

          <!-- 当前Sprint统计卡片 -->
          <UCard class="hover:shadow-md transition-shadow">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <UIcon name="i-heroicons-bolt" class="h-6 w-6 text-purple-500" />
              </div>
              <div class="ml-5 w-0 flex-1">
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400 truncate">
                  {{ $t("dashboard.currentSprint") }}
                </dt>
                <dd class="text-lg font-medium text-gray-900 dark:text-white">Sprint 2</dd>
              </div>
            </div>
          </UCard>
        </div>
      </div>

      <!-- 快速操作 -->
      <div>
        <h2 class="text-lg font-medium text-gray-900 dark:text-white mb-4">
          {{ $t("dashboard.quickActions") }}
        </h2>
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <UCard class="cursor-pointer hover:shadow-md transition-all hover:scale-[1.02]" @click="createNewTask">
            <div class="text-center py-8">
              <UIcon name="i-heroicons-plus-circle" class="mx-auto h-12 w-12 text-primary-500 mb-4" />
              <span class="block text-sm font-medium text-gray-900 dark:text-white">
                {{ $t("dashboard.createNewTask") }}
              </span>
            </div>
          </UCard>

          <UCard class="cursor-pointer hover:shadow-md transition-all hover:scale-[1.02]" @click="createSprint">
            <div class="text-center py-8">
              <UIcon name="i-heroicons-rocket-launch" class="mx-auto h-12 w-12 text-primary-500 mb-4" />
              <span class="block text-sm font-medium text-gray-900 dark:text-white">
                {{ $t("dashboard.createSprint") }}
              </span>
            </div>
          </UCard>

          <UCard class="cursor-pointer hover:shadow-md transition-all hover:scale-[1.02]" @click="viewReports">
            <div class="text-center py-8">
              <UIcon name="i-heroicons-chart-bar" class="mx-auto h-12 w-12 text-primary-500 mb-4" />
              <span class="block text-sm font-medium text-gray-900 dark:text-white">
                {{ $t("dashboard.viewReports") }}
              </span>
            </div>
          </UCard>
        </div>
      </div>

      <!-- 最近活动 -->
      <div>
        <h2 class="text-lg font-medium text-gray-900 dark:text-white mb-4">
          {{ $t("dashboard.recentActivity") }}
        </h2>
        <UCard>
          <div class="divide-y divide-gray-200 dark:divide-gray-700">
            <!-- 活动项 1 -->
            <div class="py-4 first:pt-0 last:pb-0">
              <div class="flex items-center justify-between">
                <div class="flex items-center">
                  <UAvatar size="sm" alt="John Doe" class="bg-blue-500">
                    <span class="text-xs font-medium text-white">JD</span>
                  </UAvatar>
                  <div class="ml-4">
                    <div class="text-sm font-medium text-gray-900 dark:text-white">
                      John Doe {{ $t("dashboard.createdTask") }}
                    </div>
                    <div class="text-sm text-gray-500 dark:text-gray-400">
                      "{{ $t("dashboard.implementUserAuth") }}" - {{ $t("dashboard.hoursAgo", { hours: 2 }) }}
                    </div>
                  </div>
                </div>
                <div class="flex-shrink-0">
                  <UBadge color="green" variant="soft">
                    {{ $t("dashboard.new") }}
                  </UBadge>
                </div>
              </div>
            </div>

            <!-- 活动项 2 -->
            <div class="py-4 first:pt-0 last:pb-0">
              <div class="flex items-center justify-between">
                <div class="flex items-center">
                  <UAvatar size="sm" alt="Alice Smith" class="bg-purple-500">
                    <span class="text-xs font-medium text-white">AS</span>
                  </UAvatar>
                  <div class="ml-4">
                    <div class="text-sm font-medium text-gray-900 dark:text-white">
                      Alice Smith {{ $t("dashboard.completedTask") }}
                    </div>
                    <div class="text-sm text-gray-500 dark:text-gray-400">
                      "{{ $t("dashboard.designUserInterface") }}" - {{ $t("dashboard.hoursAgo", { hours: 4 }) }}
                    </div>
                  </div>
                </div>
                <div class="flex-shrink-0">
                  <UBadge color="blue" variant="soft">
                    {{ $t("dashboard.completed") }}
                  </UBadge>
                </div>
              </div>
            </div>

            <!-- 活动项 3 -->
            <div class="py-4 first:pt-0 last:pb-0">
              <div class="flex items-center justify-between">
                <div class="flex items-center">
                  <UAvatar size="sm" alt="Mike Brown" class="bg-orange-500">
                    <span class="text-xs font-medium text-white">MB</span>
                  </UAvatar>
                  <div class="ml-4">
                    <div class="text-sm font-medium text-gray-900 dark:text-white">
                      Mike Brown {{ $t("dashboard.startedSprint") }}
                    </div>
                    <div class="text-sm text-gray-500 dark:text-gray-400">
                      "Sprint 2 - {{ $t("dashboard.advancedFeatures") }}" - {{ $t("dashboard.daysAgo", { days: 1 }) }}
                    </div>
                  </div>
                </div>
                <div class="flex-shrink-0">
                  <UBadge color="yellow" variant="soft">
                    Sprint
                  </UBadge>
                </div>
              </div>
            </div>
          </div>
        </UCard>
      </div>
    </main>
  </div>
</template>

<script setup>
// 国际化
const { t } = useI18n();

// 页面元数据
useSeoMeta({
  title: () => `${t('common.appName')} - ${t('navigation.dashboard')}`,
  description: () => t('dashboard.description')
})

// 获取运行时配置
const _config = useRuntimeConfig()

// 快速操作事件处理函数
const createNewTask = () => {
  // TODO: 实现创建新任务功能
  navigateTo('/tasks/create')
}

const createSprint = () => {
  // TODO: 实现创建冲刺功能
  navigateTo('/sprints/create')
}

const viewReports = () => {
  // TODO: 跳转到报告页面
  navigateTo('/reports')
}

// 这里可以添加数据获取逻辑
// const { data: stats } = await $fetch(`${config.public.apiBaseUrl}/stats`)
</script>

<style scoped>
/* 如果需要自定义样式可以在这里添加 */
</style>