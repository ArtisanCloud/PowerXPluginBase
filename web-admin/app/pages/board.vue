<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">
          {{ $t("navigation.taskBoard") }}
        </h1>
        <p class="mt-2 text-sm text-gray-600 dark:text-gray-400">
          {{ $t("board.description") }}
        </p>
      </div>
      <div class="flex items-center space-x-3">
        <USelect
          v-model="selectedSprint"
          :options="sprintOptions"
          :placeholder="$t('sprint.selectSprint')"
          size="sm"
        />
        <UButton
          color="primary"
          icon="i-heroicons-plus"
          @click="createTask"
        >
          {{ $t("task.createTask") }}
        </UButton>
      </div>
    </div>

    <!-- 看板列 -->
    <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-6 min-h-[600px]">
      <!-- 待办列 -->
      <div class="bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
        <div class="flex items-center justify-between mb-4">
          <h3 class="font-medium text-gray-900 dark:text-white">
            {{ $t("status.todo") }}
          </h3>
          <UBadge size="xs" variant="outline">
            {{ todoTasks.length }}
          </UBadge>
        </div>
        <div class="space-y-3">
          <TaskCard
            v-for="task in todoTasks"
            :key="task.id"
            :task="task"
            @edit="editTask"
            @delete="deleteTask"
          />
        </div>
      </div>

      <!-- 进行中列 -->
      <div class="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4">
        <div class="flex items-center justify-between mb-4">
          <h3 class="font-medium text-gray-900 dark:text-white">
            {{ $t("status.inProgress") }}
          </h3>
          <UBadge size="xs" variant="outline">
            {{ inProgressTasks.length }}
          </UBadge>
        </div>
        <div class="space-y-3">
          <TaskCard
            v-for="task in inProgressTasks"
            :key="task.id"
            :task="task"
            @edit="editTask"
            @delete="deleteTask"
          />
        </div>
      </div>

      <!-- 审核中列 -->
      <div class="bg-yellow-50 dark:bg-yellow-900/20 rounded-lg p-4">
        <div class="flex items-center justify-between mb-4">
          <h3 class="font-medium text-gray-900 dark:text-white">
            {{ $t("status.review") }}
          </h3>
          <UBadge size="xs" variant="outline">
            {{ reviewTasks.length }}
          </UBadge>
        </div>
        <div class="space-y-3">
          <TaskCard
            v-for="task in reviewTasks"
            :key="task.id"
            :task="task"
            @edit="editTask"
            @delete="deleteTask"
          />
        </div>
      </div>

      <!-- 已完成列 -->
      <div class="bg-green-50 dark:bg-green-900/20 rounded-lg p-4">
        <div class="flex items-center justify-between mb-4">
          <h3 class="font-medium text-gray-900 dark:text-white">
            {{ $t("status.done") }}
          </h3>
          <UBadge size="xs" variant="outline">
            {{ doneTasks.length }}
          </UBadge>
        </div>
        <div class="space-y-3">
          <TaskCard
            v-for="task in doneTasks"
            :key="task.id"
            :task="task"
            @edit="editTask"
            @delete="deleteTask"
          />
        </div>
      </div>
    </div>

    <!-- 统计信息 -->
    <UCard>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div class="text-center">
          <div class="text-2xl font-bold text-gray-900 dark:text-white">
            {{ totalTasks }}
          </div>
          <div class="text-sm text-gray-500 dark:text-gray-400">
            {{ $t("board.totalTasks") }}
          </div>
        </div>
        <div class="text-center">
          <div class="text-2xl font-bold text-blue-600">
            {{ inProgressTasks.length }}
          </div>
          <div class="text-sm text-gray-500 dark:text-gray-400">
            {{ $t("board.activeTasks") }}
          </div>
        </div>
        <div class="text-center">
          <div class="text-2xl font-bold text-green-600">
            {{ doneTasks.length }}
          </div>
          <div class="text-sm text-gray-500 dark:text-gray-400">
            {{ $t("board.completedTasks") }}
          </div>
        </div>
        <div class="text-center">
          <div class="text-2xl font-bold text-purple-600">
            {{ completionRate }}%
          </div>
          <div class="text-sm text-gray-500 dark:text-gray-400">
            {{ $t("board.completionRate") }}
          </div>
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
  title: () => `${t('common.appName')} - ${t('navigation.taskBoard')}`,
  description: () => t('board.description')
});

// 选中的冲刺
const selectedSprint = ref('sprint-1');

// 冲刺选项
const sprintOptions = computed(() => [
  { label: 'Sprint 1', value: 'sprint-1' },
  { label: 'Sprint 2', value: 'sprint-2' },
  { label: 'Sprint 3', value: 'sprint-3' }
]);

// Mock 任务数据
const tasks = ref([
  {
    id: 'TASK-001',
    title: t('board.sampleTask1'),
    description: t('board.sampleTaskDescription1'),
    status: 'todo',
    priority: 'high',
    assignee: 'John Doe',
    storyPoints: 5,
    labels: ['frontend', 'urgent']
  },
  {
    id: 'TASK-002',
    title: t('board.sampleTask2'),
    description: t('board.sampleTaskDescription2'),
    status: 'inProgress',
    priority: 'medium',
    assignee: 'Alice Smith',
    storyPoints: 3,
    labels: ['backend']
  },
  {
    id: 'TASK-003',
    title: t('board.sampleTask3'),
    description: t('board.sampleTaskDescription3'),
    status: 'review',
    priority: 'low',
    assignee: 'Bob Johnson',
    storyPoints: 2,
    labels: ['testing']
  },
  {
    id: 'TASK-004',
    title: t('board.sampleTask4'),
    description: t('board.sampleTaskDescription4'),
    status: 'done',
    priority: 'high',
    assignee: 'Carol Brown',
    storyPoints: 8,
    labels: ['database', 'migration']
  }
]);

// 按状态分组任务
const todoTasks = computed(() => tasks.value.filter(task => task.status === 'todo'));
const inProgressTasks = computed(() => tasks.value.filter(task => task.status === 'inProgress'));
const reviewTasks = computed(() => tasks.value.filter(task => task.status === 'review'));
const doneTasks = computed(() => tasks.value.filter(task => task.status === 'done'));

// 统计信息
const totalTasks = computed(() => tasks.value.length);
const completionRate = computed(() => {
  if (totalTasks.value === 0) return 0;
  return Math.round((doneTasks.value.length / totalTasks.value) * 100);
});

// 事件处理
const createTask = () => {
  // TODO: 实现创建任务功能
  console.log('创建任务');
};

const editTask = (task) => {
  // TODO: 实现编辑任务功能
  console.log('编辑任务:', task);
};

const deleteTask = (task) => {
  // TODO: 实现删除任务功能
  console.log('删除任务:', task);
};
</script>

<style scoped>
/* 任务卡片拖拽样式 */
.task-card {
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.task-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}
</style>