<template>
  <div class="space-y-6">
    <!-- 页面标题 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">
          {{ $t("navigation.backlog") }}
        </h1>
        <p class="mt-2 text-sm text-gray-600 dark:text-gray-400">
          {{ $t("backlog.description") }}
        </p>
      </div>
      <UButton
        color="primary"
        icon="i-heroicons-plus"
        @click="createUserStory"
      >
        {{ $t("backlog.createUserStory") }}
      </UButton>
    </div>

    <!-- 过滤和排序 -->
    <UCard>
      <div class="flex flex-col sm:flex-row gap-4">
        <div class="flex-1">
          <UInput
            v-model="searchQuery"
            :placeholder="$t('common.search')"
            icon="i-heroicons-magnifying-glass"
          />
        </div>
        <div class="flex gap-2">
          <USelect
            v-model="selectedPriority"
            :options="priorityOptions"
            :placeholder="$t('task.priority')"
            size="sm"
          />
          <USelect
            v-model="selectedStatus"
            :options="statusOptions"
            :placeholder="$t('task.status')"
            size="sm"
          />
        </div>
      </div>
    </UCard>

    <!-- Epic 分组 -->
    <div v-for="epic in epics" :key="epic.id" class="space-y-4">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <UIcon name="i-heroicons-bookmark" class="h-5 w-5 text-purple-500" />
              <div>
                <h3 class="text-lg font-medium text-gray-900 dark:text-white">
                  {{ epic.title }}
                </h3>
                <p class="text-sm text-gray-500 dark:text-gray-400">
                  {{ epic.description }}
                </p>
              </div>
            </div>
            <UBadge :color="epic.status === 'active' ? 'green' : 'gray'" variant="soft">
              {{ $t(`status.${epic.status}`) }}
            </UBadge>
          </div>
        </template>

        <!-- 用户故事列表 -->
        <div class="space-y-3">
          <div
            v-for="story in epic.userStories"
            :key="story.id"
            class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:bg-gray-50 dark:hover:bg-gray-800 cursor-pointer transition-colors"
            @click="editUserStory(story)"
          >
            <div class="flex items-start justify-between">
              <div class="flex-1">
                <div class="flex items-center space-x-2 mb-2">
                  <UBadge size="xs" variant="outline">
                    {{ story.id }}
                  </UBadge>
                  <UBadge
                    :color="getPriorityColor(story.priority)"
                    size="xs"
                    variant="soft"
                  >
                    {{ $t(`priority.${story.priority}`) }}
                  </UBadge>
                </div>
                <h4 class="text-sm font-medium text-gray-900 dark:text-white mb-1">
                  {{ story.title }}
                </h4>
                <p class="text-sm text-gray-600 dark:text-gray-400 mb-2">
                  {{ story.description }}
                </p>
                <div class="flex items-center space-x-4 text-xs text-gray-500 dark:text-gray-400">
                  <span>{{ $t("task.storyPoints") }}: {{ story.storyPoints }}</span>
                  <span>{{ $t("task.assignee") }}: {{ story.assignee || $t("common.unassigned") }}</span>
                </div>
              </div>
              <div class="flex items-center space-x-2">
                <UButton
                  icon="i-heroicons-pencil"
                  size="xs"
                  variant="ghost"
                  @click.stop="editUserStory(story)"
                />
                <UButton
                  icon="i-heroicons-trash"
                  size="xs"
                  variant="ghost"
                  color="red"
                  @click.stop="deleteUserStory(story)"
                />
              </div>
            </div>
          </div>
        </div>
      </UCard>
    </div>

    <!-- 空状态 -->
    <UCard v-if="epics.length === 0">
      <div class="text-center py-12">
        <UIcon name="i-heroicons-document-text" class="mx-auto h-12 w-12 text-gray-400" />
        <h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white">
          {{ $t("backlog.noItems") }}
        </h3>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
          {{ $t("backlog.noItemsDescription") }}
        </p>
        <div class="mt-6">
          <UButton color="primary" @click="createUserStory">
            {{ $t("backlog.createFirstStory") }}
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
  title: () => `${t('common.appName')} - ${t('navigation.backlog')}`,
  description: () => t('backlog.description')
});

// 搜索和过滤
const searchQuery = ref('');
const selectedPriority = ref('');
const selectedStatus = ref('');

// 优先级选项
const priorityOptions = computed(() => [
  { label: t('common.all'), value: '' },
  { label: t('priority.critical'), value: 'critical' },
  { label: t('priority.high'), value: 'high' },
  { label: t('priority.medium'), value: 'medium' },
  { label: t('priority.low'), value: 'low' }
]);

// 状态选项
const statusOptions = computed(() => [
  { label: t('common.all'), value: '' },
  { label: t('status.todo'), value: 'todo' },
  { label: t('status.inProgress'), value: 'inProgress' },
  { label: t('status.done'), value: 'done' }
]);

// Mock 数据
const epics = ref([
  {
    id: 'EPIC-001',
    title: t('backlog.sampleEpic'),
    description: t('backlog.sampleEpicDescription'),
    status: 'active',
    userStories: [
      {
        id: 'US-001',
        title: t('backlog.sampleUserStory1'),
        description: t('backlog.sampleUserStoryDescription1'),
        priority: 'high',
        storyPoints: 8,
        assignee: 'John Doe',
        status: 'todo'
      },
      {
        id: 'US-002',
        title: t('backlog.sampleUserStory2'),
        description: t('backlog.sampleUserStoryDescription2'),
        priority: 'medium',
        storyPoints: 5,
        assignee: null,
        status: 'todo'
      }
    ]
  }
]);

// 获取优先级颜色
const getPriorityColor = (priority) => {
  const colors = {
    critical: 'red',
    high: 'orange',
    medium: 'yellow',
    low: 'green'
  };
  return colors[priority] || 'gray';
};

// 事件处理
const createUserStory = () => {
  // TODO: 实现创建用户故事功能
  console.log('创建用户故事');
};

const editUserStory = (story) => {
  // TODO: 实现编辑用户故事功能
  console.log('编辑用户故事:', story);
};

const deleteUserStory = (story) => {
  // TODO: 实现删除用户故事功能
  console.log('删除用户故事:', story);
};
</script>