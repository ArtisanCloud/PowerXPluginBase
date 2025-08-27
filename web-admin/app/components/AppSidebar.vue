<template>
  <aside
    class="w-64 min-w-64 max-w-64 bg-white dark:bg-gray-900 border-r border-gray-200 dark:border-gray-800 min-h-screen flex-shrink-0"
  >
    <nav class="p-4 space-y-6">
      <!-- 仪表盘 -->
      <div>
        <UButton
          to="/dashboard"
          variant="ghost"
          color="neutral"
          class="w-full justify-start"
          :class="{
            'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
              $route.path === '/dashboard',
          }"
        >
          <UIcon name="i-heroicons-chart-bar" class="w-4 h-4 mr-3" />
          {{ $t("navigation.dashboard") }}
        </UButton>
      </div>

      <!-- Scrum 管理 -->
      <div>
        <div
          class="px-3 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider"
        >
          {{ $t("navigation.scrumManagement") }}
        </div>
        <div class="space-y-1">
          <!-- 产品待办列表 -->
          <UButton
            to="/backlog"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/backlog',
            }"
          >
            <UIcon name="i-heroicons-list-bullet" class="w-4 h-4 mr-3" />
            {{ $t("navigation.backlog") }}
          </UButton>

          <!-- 冲刺管理 -->
          <div>
            <UButton
              variant="ghost"
              color="neutral"
              class="w-full justify-start"
              @click="toggleSprintMenu"
            >
              <UIcon name="i-heroicons-clock" class="w-4 h-4 mr-3" />
              {{ $t("navigation.sprints") }}
              <UIcon
                :name="
                  showSprintMenu
                    ? 'i-heroicons-chevron-down'
                    : 'i-heroicons-chevron-right'
                "
                class="w-4 h-4 ml-auto"
              />
            </UButton>

            <div v-show="showSprintMenu" class="ml-6 mt-1 space-y-1">
              <UButton
                to="/sprints/active"
                variant="ghost"
                color="neutral"
                size="sm"
                class="w-full justify-start text-sm"
                :class="{
                  'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                    $route.path === '/sprints/active',
                }"
              >
                <UIcon name="i-heroicons-play" class="w-3 h-3 mr-2" />
                {{ $t("navigation.activeSprint") }}
              </UButton>

              <UButton
                to="/sprints/planning"
                variant="ghost"
                color="neutral"
                size="sm"
                class="w-full justify-start text-sm"
                :class="{
                  'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                    $route.path === '/sprints/planning',
                }"
              >
                <UIcon name="i-heroicons-calendar" class="w-3 h-3 mr-2" />
                {{ $t("navigation.sprintPlanning") }}
              </UButton>

              <UButton
                to="/sprints/retrospective"
                variant="ghost"
                color="neutral"
                size="sm"
                class="w-full justify-start text-sm"
                :class="{
                  'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                    $route.path === '/sprints/retrospective',
                }"
              >
                <UIcon name="i-heroicons-arrow-path-rounded-square" class="w-3 h-3 mr-2" />
                {{ $t("navigation.retrospective") }}
              </UButton>
            </div>
          </div>

          <!-- 任务看板 -->
          <UButton
            to="/board"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/board',
            }"
          >
            <UIcon name="i-heroicons-view-columns" class="w-4 h-4 mr-3" />
            {{ $t("navigation.taskBoard") }}
          </UButton>

          <!-- 任务管理 -->
          <UButton
            to="/tasks"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/tasks',
            }"
          >
            <UIcon name="i-heroicons-check-circle" class="w-4 h-4 mr-3" />
            {{ $t("navigation.tasks") }}
          </UButton>
        </div>
      </div>

      <!-- 团队协作 -->
      <div>
        <div
          class="px-3 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider"
        >
          {{ $t("navigation.teamCollaboration") }}
        </div>
        <div class="space-y-1">
          <!-- 团队成员 -->
          <UButton
            to="/team"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/team',
            }"
          >
            <UIcon name="i-heroicons-users" class="w-4 h-4 mr-3" />
            {{ $t("navigation.team") }}
          </UButton>

          <!-- 每日站会 -->
          <UButton
            to="/daily-standup"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/daily-standup',
            }"
          >
            <UIcon name="i-heroicons-microphone" class="w-4 h-4 mr-3" />
            {{ $t("navigation.dailyStandup") }}
          </UButton>
        </div>
      </div>

      <!-- 报告分析 -->
      <div>
        <div
          class="px-3 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider"
        >
          {{ $t("navigation.reportsAnalytics") }}
        </div>
        <div class="space-y-1">
          <!-- 燃尽图 -->
          <UButton
            to="/reports/burndown"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/reports/burndown',
            }"
          >
            <UIcon name="i-heroicons-chart-bar" class="w-4 h-4 mr-3" />
            {{ $t("navigation.burndownChart") }}
          </UButton>

          <!-- 速度图表 -->
          <UButton
            to="/reports/velocity"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/reports/velocity',
            }"
          >
            <UIcon name="i-heroicons-arrow-trending-up" class="w-4 h-4 mr-3" />
            {{ $t("navigation.velocityChart") }}
          </UButton>

          <!-- 累积流图 -->
          <UButton
            to="/reports/cumulative-flow"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/reports/cumulative-flow',
            }"
          >
            <UIcon name="i-heroicons-chart-pie" class="w-4 h-4 mr-3" />
            {{ $t("navigation.cumulativeFlow") }}
          </UButton>
        </div>
      </div>

      <!-- 系统设置 -->
      <div>
        <div
          class="px-3 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider"
        >
          {{ $t("navigation.systemSettings") }}
        </div>
        <div class="space-y-1">
          <!-- 项目设置 -->
          <UButton
            to="/settings/project"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/settings/project',
            }"
          >
            <UIcon name="i-heroicons-cog-6-tooth" class="w-4 h-4 mr-3" />
            {{ $t("navigation.projectSettings") }}
          </UButton>

          <!-- 权限管理 -->
          <UButton
            to="/settings/permissions"
            variant="ghost"
            color="neutral"
            class="w-full justify-start"
            :class="{
              'bg-primary-50 text-primary-600 dark:bg-primary-950 dark:text-primary-400':
                $route.path === '/settings/permissions',
            }"
          >
            <UIcon name="i-heroicons-key" class="w-4 h-4 mr-3" />
            {{ $t("navigation.permissions") }}
          </UButton>
        </div>
      </div>
    </nav>
  </aside>
</template>

<script setup lang="ts">
const { t } = useI18n();

// 控制冲刺管理子菜单的展开状态
const showSprintMenu = ref(false);

// 切换冲刺管理子菜单
const toggleSprintMenu = () => {
  showSprintMenu.value = !showSprintMenu.value;
};

// 监听路由变化，如果当前路由在冲刺模块下，自动展开子菜单
const route = useRoute();
watch(
  () => route.path,
  (newPath) => {
    if (newPath.startsWith("/sprints/")) {
      showSprintMenu.value = true;
    }
  },
  { immediate: true }
);
</script>
