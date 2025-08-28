<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">
          {{ $t("navigation.activeSprint") }}
        </h1>
        <p class="mt-2 text-sm text-gray-600 dark:text-gray-400">
          {{ $t("sprint.activeDescription") }}
        </p>
      </div>
      <div class="flex items-center space-x-3">
        <UButton
          variant="outline"
          icon="i-heroicons-chart-bar"
          @click="viewBurndown"
        >
          {{ $t("navigation.burndownChart") }}
        </UButton>
        <UButton
          color="red"
          icon="i-heroicons-stop"
          @click="completeSprint"
        >
          {{ $t("sprint.completeSprint") }}
        </UButton>
      </div>
    </div>

    <!-- 冲刺信息卡片 -->
    <UCard v-if="activeSprint">
      <div class="flex items-start justify-between">
        <div class="flex-1">
          <div class="flex items-center space-x-3 mb-4">
            <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
              {{ activeSprint.name }}
            </h2>
            <UBadge color="green" variant="soft">
              {{ $t("sprint.active") }}
            </UBadge>
          </div>
          <p class="text-gray-600 dark:text-gray-400 mb-4">
            {{ activeSprint.goal }}
          </p>
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div>
              <div class="text-sm text-gray-500 dark:text-gray-400">
                {{ $t("sprint.startDate") }}
              </div>
              <div class="font-medium text-gray-900 dark:text-white">
                {{ formatDate(activeSprint.startDate) }}
              </div>
            </div>
            <div>
              <div class="text-sm text-gray-500 dark:text-gray-400">
                {{ $t("sprint.endDate") }}
              </div>
              <div class="font-medium text-gray-900 dark:text-white">
                {{ formatDate(activeSprint.endDate) }}
              </div>
            </div>
            <div>
              <div class="text-sm text-gray-500 dark:text-gray-400">
                {{ $t("sprint.daysRemaining") }}
              </div>
              <div class="font-medium text-gray-900 dark:text-white">
                {{ daysRemaining }} {{ $t("common.days") }}
              </div>
            </div>
            <div>
              <div class="text-sm text-gray-500 dark:text-gray-400">
                {{ $t("sprint.progress") }}
              </div>
              <div class="font-medium text-gray-900 dark:text-white">
                {{ sprintProgress }}%
              </div>
            </div>
          </div>
        </div>
      </div>
    </UCard>

    <!-- 冲刺统计 -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <UCard>
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <UIcon name="i-heroicons-clipboard-document-list" class="h-8 w-8 text-blue-500" />
          </div>
          <div class="ml-4">
            <div class="text-2xl font-bold text-gray-900 dark:text-white">
              {{ totalTasks }}
            </div>
            <div class="text-sm text-gray-500 dark:text-gray-400">
              {{ $t("sprint.totalTasks") }}
            </div>
          </div>
        </div>
      </UCard>

      <UCard>
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <UIcon name="i-heroicons-check-circle" class="h-8 w-8 text-green-500" />
          </div>
          <div class="ml-4">
            <div class="text-2xl font-bold text-gray-900 dark:text-white">
              {{ completedTasks }}
            </div>
            <div class="text-sm text-gray-500 dark:text-gray-400">
              {{ $t("sprint.completedTasks") }}
            </div>
          </div>
        </div>
      </UCard>

      <UCard>
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <UIcon name="i-heroicons-chart-bar" class="h-8 w-8 text-purple-500" />
          </div>
          <div class="ml-4">
            <div class="text-2xl font-bold text-gray-900 dark:text-white">
              {{ totalStoryPoints }}
            </div>
            <div class="text-sm text-gray-500 dark:text-gray-400">
              {{ $t("sprint.totalPoints") }}
            </div>
          </div>
        </div>
      </UCard>

      <UCard>
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <UIcon name="i-heroicons-bolt" class="h-8 w-8 text-orange-500" />
          </div>
          <div class="ml-4">
            <div class="text-2xl font-bold text-gray-900 dark:text-white">
              {{ velocity }}
            </div>
            <div class="text-sm text-gray-500 dark:text-gray-400">
              {{ $t("sprint.velocity") }}
            </div>
          </div>
        </div>
      </UCard>
    </div>

    <!-- 进度图表 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- 燃尽图 -->
      <UCard>
        <template #header>
          <h3 class="text-lg font-medium text-gray-900 dark:text-white">
            {{ $t("navigation.burndownChart") }}
          </h3>
        </template>
        <div class="h-64 flex items-center justify-center bg-gray-50 dark:bg-gray-800 rounded-lg">
          <div class="text-center text-gray-500 dark:text-gray-400">
            <UIcon name="i-heroicons-chart-bar" class="h-12 w-12 mx-auto mb-2" />
            <p>{{ $t("sprint.burndownPlaceholder") }}</p>
          </div>
        </div>
      </UCard>

      <!-- 任务分布 -->
      <UCard>
        <template #header>
          <h3 class="text-lg font-medium text-gray-900 dark:text-white">
            {{ $t("sprint.taskDistribution") }}
          </h3>
        </template>
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-2">
              <div class="w-3 h-3 bg-gray-400 rounded-full"/>
              <span class="text-sm text-gray-600 dark:text-gray-400">
                {{ $t("status.todo") }}
              </span>
            </div>
            <div class="flex items-center space-x-2">
              <span class="text-sm font-medium text-gray-900 dark:text-white">
                {{ todoTasksCount }}
              </span>
              <div class="w-20 bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                <div 
                  class="bg-gray-400 h-2 rounded-full" 
                  :style="{ width: `${(todoTasksCount / totalTasks) * 100}%` }"
                />
              </div>
            </div>
          </div>

          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-2">
              <div class="w-3 h-3 bg-blue-500 rounded-full"/>
              <span class="text-sm text-gray-600 dark:text-gray-400">
                {{ $t("status.inProgress") }}
              </span>
            </div>
            <div class="flex items-center space-x-2">
              <span class="text-sm font-medium text-gray-900 dark:text-white">
                {{ inProgressTasksCount }}
              </span>
              <div class="w-20 bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                <div 
                  class="bg-blue-500 h-2 rounded-full" 
                  :style="{ width: `${(inProgressTasksCount / totalTasks) * 100}%` }"
                />
              </div>
            </div>
          </div>

          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-2">
              <div class="w-3 h-3 bg-green-500 rounded-full"/>
              <span class="text-sm text-gray-600 dark:text-gray-400">
                {{ $t("status.done") }}
              </span>
            </div>
            <div class="flex items-center space-x-2">
              <span class="text-sm font-medium text-gray-900 dark:text-white">
                {{ completedTasks }}
              </span>
              <div class="w-20 bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                <div 
                  class="bg-green-500 h-2 rounded-full" 
                  :style="{ width: `${(completedTasks / totalTasks) * 100}%` }"
                />
              </div>
            </div>
          </div>
        </div>
      </UCard>
    </div>

    <!-- 冲刺任务列表 -->
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-medium text-gray-900 dark:text-white">
            {{ $t("sprint.sprintBacklog") }}
          </h3>
          <UButton
            icon="i-heroicons-plus"
            size="sm"
            @click="addTask"
          >
            {{ $t("task.createTask") }}
          </UButton>
        </div>
      </template>

      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
          <thead class="bg-gray-50 dark:bg-gray-800">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                {{ $t("task.title") }}
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                {{ $t("task.status") }}
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                {{ $t("task.assignee") }}
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                {{ $t("task.storyPoints") }}
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                {{ $t("task.priority") }}
              </th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                {{ $t("common.actions") }}
              </th>
            </tr>
          </thead>
          <tbody class="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-700">
            <tr
              v-for="task in sprintTasks"
              :key="task.id"
              class="hover:bg-gray-50 dark:hover:bg-gray-800"
            >
              <td class="px-6 py-4 whitespace-nowrap">
                <div>
                  <div class="text-sm font-medium text-gray-900 dark:text-white">
                    {{ task.title }}
                  </div>
                  <div class="text-sm text-gray-500 dark:text-gray-400">
                    {{ task.id }}
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <UBadge
                  :color="getStatusColor(task.status)"
                  variant="soft"
                  size="sm"
                >
                  {{ $t(`status.${task.status}`) }}
                </UBadge>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <UAvatar
                    v-if="task.assignee"
                    :alt="task.assignee"
                    size="xs"
                    class="bg-primary-500 mr-2"
                  >
                    <span class="text-xs font-medium text-white">
                      {{ getInitials(task.assignee) }}
                    </span>
                  </UAvatar>
                  <span class="text-sm text-gray-900 dark:text-white">
                    {{ task.assignee || $t("common.unassigned") }}
                  </span>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                {{ task.storyPoints }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <UBadge
                  :color="getPriorityColor(task.priority)"
                  variant="soft"
                  size="sm"
                >
                  {{ $t(`priority.${task.priority}`) }}
                </UBadge>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <UDropdownMenu :items="getTaskActions(task)" :popper="{ placement: 'bottom-end' }">
                  <UButton
                    icon="i-heroicons-ellipsis-vertical"
                    size="xs"
                    variant="ghost"
                  />
                </UDropdownMenu>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>

    <!-- 无活跃冲刺状态 -->
    <UCard v-if="!activeSprint">
      <div class="text-center py-12">
        <UIcon name="i-heroicons-clock" class="mx-auto h-12 w-12 text-gray-400" />
        <h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white">
          {{ $t("sprint.noActiveSprint") }}
        </h3>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
          {{ $t("sprint.noActiveSprintDescription") }}
        </p>
        <div class="mt-6">
          <UButton color="primary" @click="startSprint">
            {{ $t("sprint.startSprint") }}
          </UButton>
        </div>
      </div>
    </UCard>
  </div>
</template>

<script setup>
// 国际化
const { t } = useI18n();

// 页面元数据
useSeoMeta({
  title: () => `${t('common.appName')} - ${t('navigation.activeSprint')}`,
  description: () => t('sprint.activeDescription')
});

// Mock 活跃冲刺数据
const activeSprint = ref({
  id: 'SPRINT-001',
  name: 'Sprint 2 - 用户管理功能',
  goal: '完成用户注册、登录和个人资料管理功能',
  startDate: '2024-01-15',
  endDate: '2024-01-29',
  status: 'active'
});

// Mock 冲刺任务数据
const sprintTasks = ref([
  {
    id: 'TASK-001',
    title: '实现用户注册功能',
    status: 'done',
    assignee: 'John Doe',
    storyPoints: 5,
    priority: 'high'
  },
  {
    id: 'TASK-002',
    title: '设计登录界面',
    status: 'inProgress',
    assignee: 'Alice Smith',
    storyPoints: 3,
    priority: 'medium'
  },
  {
    id: 'TASK-003',
    title: '用户个人资料页面',
    status: 'todo',
    assignee: null,
    storyPoints: 8,
    priority: 'medium'
  },
  {
    id: 'TASK-004',
    title: '密码重置功能',
    status: 'todo',
    assignee: 'Bob Johnson',
    storyPoints: 5,
    priority: 'low'
  }
]);

// 计算属性
const totalTasks = computed(() => sprintTasks.value.length);
const completedTasks = computed(() => sprintTasks.value.filter(t => t.status === 'done').length);
const todoTasksCount = computed(() => sprintTasks.value.filter(t => t.status === 'todo').length);
const inProgressTasksCount = computed(() => sprintTasks.value.filter(t => t.status === 'inProgress').length);
const totalStoryPoints = computed(() => sprintTasks.value.reduce((sum, t) => sum + t.storyPoints, 0));
const completedStoryPoints = computed(() => 
  sprintTasks.value.filter(t => t.status === 'done').reduce((sum, t) => sum + t.storyPoints, 0)
);
const velocity = computed(() => completedStoryPoints.value);
const sprintProgress = computed(() => {
  if (totalTasks.value === 0) return 0;
  return Math.round((completedTasks.value / totalTasks.value) * 100);
});

const daysRemaining = computed(() => {
  if (!activeSprint.value) return 0;
  const endDate = new Date(activeSprint.value.endDate);
  const today = new Date();
  const timeDiff = endDate.getTime() - today.getTime();
  const daysDiff = Math.ceil(timeDiff / (1000 * 3600 * 24));
  return Math.max(0, daysDiff);
});

// 工具函数
const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('zh-CN');
};

const getInitials = (name) => {
  if (!name) return '';
  return name
    .split(' ')
    .map(word => word.charAt(0))
    .join('')
    .toUpperCase()
    .slice(0, 2);
};

const getStatusColor = (status) => {
  const colors = {
    todo: 'gray',
    inProgress: 'blue',
    review: 'yellow',
    done: 'green'
  };
  return colors[status] || 'gray';
};

const getPriorityColor = (priority) => {
  const colors = {
    critical: 'red',
    high: 'orange',
    medium: 'yellow',
    low: 'green'
  };
  return colors[priority] || 'gray';
};

const getTaskActions = (task) => {
  return [
    [{
      label: t('common.edit'),
      icon: 'i-heroicons-pencil',
      click: () => editTask(task)
    }],
    [{
      label: t('common.delete'),
      icon: 'i-heroicons-trash',
      click: () => deleteTask(task)
    }]
  ];
};

// 事件处理
const viewBurndown = () => {
  navigateTo('/reports/burndown');
};

const completeSprint = () => {
  // TODO: 实现完成冲刺功能
  console.log('完成冲刺');
};

const addTask = () => {
  // TODO: 实现添加任务功能
  console.log('添加任务');
};

const editTask = (task) => {
  // TODO: 实现编辑任务功能
  console.log('编辑任务:', task);
};

const deleteTask = (task) => {
  // TODO: 实现删除任务功能
  console.log('删除任务:', task);
};

const startSprint = () => {
  // TODO: 实现开始冲刺功能
  navigateTo('/sprints/planning');
};
</script>