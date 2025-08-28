<template>
  <div class="space-y-6">
    <!-- 页面头部 -->
    <div class="border-b border-gray-200 dark:border-gray-700 pb-4">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
        {{ $t("navigation.retrospective") }}
      </h1>
      <p class="text-gray-600 dark:text-gray-400 mt-1">
        {{ $t("retrospective.description") }}
      </p>
    </div>

    <!-- 回顾会议设置 -->
    <UCard>
      <template #header>
        <div class="flex justify-between items-center">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ $t("retrospective.currentRetrospective") }}
          </h2>
          <UButton
            color="primary"
            icon="i-heroicons-plus"
            @click="startRetrospective"
          >
            {{ $t("retrospective.startRetrospective") }}
          </UButton>
        </div>
      </template>

      <div v-if="!activeRetrospective" class="text-center py-12">
        <UIcon name="i-heroicons-chat-bubble-left-ellipsis" class="w-12 h-12 text-gray-400 mx-auto mb-4" />
        <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">
          {{ $t("retrospective.noActiveRetrospective") }}
        </h3>
        <p class="text-gray-600 dark:text-gray-400 mb-4">
          {{ $t("retrospective.noActiveRetrospectiveDescription") }}
        </p>
        <UButton
          color="primary"
          @click="startRetrospective"
        >
          {{ $t("retrospective.startFirstRetrospective") }}
        </UButton>
      </div>

      <div v-else class="space-y-6">
        <!-- 回顾状态 -->
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-3">
            <UBadge :color="retrospectiveStatusColor" variant="soft">
              {{ retrospectiveStatusText }}
            </UBadge>
            <span class="text-sm text-gray-600 dark:text-gray-400">
              {{ $t("retrospective.startedAt") }}: {{ formatTime(activeRetrospective.startedAt) }}
            </span>
          </div>
          <UButton
            v-if="activeRetrospective.status === 'inProgress'"
            color="red"
            variant="outline"
            @click="endRetrospective"
          >
            {{ $t("retrospective.endRetrospective") }}
          </UButton>
        </div>

        <!-- 三列布局：做得好的、需要改进的、行动项 -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <!-- 做得好的 -->
          <div class="space-y-4">
            <h3 class="text-lg font-medium text-green-600 dark:text-green-400 flex items-center">
              <UIcon name="i-heroicons-hand-thumb-up" class="w-5 h-5 mr-2" />
              {{ $t("retrospective.wentWell") }}
            </h3>
            <div class="space-y-2">
              <div
                v-for="item in wentWellItems"
                :key="item.id"
                class="p-3 bg-green-50 dark:bg-green-900/20 rounded-lg border border-green-200 dark:border-green-800"
              >
                <p class="text-sm text-gray-900 dark:text-white">{{ item.content }}</p>
                <div class="flex items-center justify-between mt-2">
                  <span class="text-xs text-gray-500">{{ item.author }}</span>
                  <span class="text-xs text-gray-400">{{ item.votes }} 👍</span>
                </div>
              </div>
            </div>
            <UTextarea
              v-model="newWentWellItem"
              :placeholder="$t('retrospective.addWentWellPlaceholder')"
              size="sm"
              @keydown.enter.shift.prevent="addWentWellItem"
            />
          </div>

          <!-- 需要改进的 -->
          <div class="space-y-4">
            <h3 class="text-lg font-medium text-yellow-600 dark:text-yellow-400 flex items-center">
              <UIcon name="i-heroicons-exclamation-triangle" class="w-5 h-5 mr-2" />
              {{ $t("retrospective.needsImprovement") }}
            </h3>
            <div class="space-y-2">
              <div
                v-for="item in improvementItems"
                :key="item.id"
                class="p-3 bg-yellow-50 dark:bg-yellow-900/20 rounded-lg border border-yellow-200 dark:border-yellow-800"
              >
                <p class="text-sm text-gray-900 dark:text-white">{{ item.content }}</p>
                <div class="flex items-center justify-between mt-2">
                  <span class="text-xs text-gray-500">{{ item.author }}</span>
                  <span class="text-xs text-gray-400">{{ item.votes }} 👍</span>
                </div>
              </div>
            </div>
            <UTextarea
              v-model="newImprovementItem"
              :placeholder="$t('retrospective.addImprovementPlaceholder')"
              size="sm"
              @keydown.enter.shift.prevent="addImprovementItem"
            />
          </div>

          <!-- 行动项 -->
          <div class="space-y-4">
            <h3 class="text-lg font-medium text-blue-600 dark:text-blue-400 flex items-center">
              <UIcon name="i-heroicons-bolt" class="w-5 h-5 mr-2" />
              {{ $t("retrospective.actionItems") }}
            </h3>
            <div class="space-y-2">
              <div
                v-for="item in actionItems"
                :key="item.id"
                class="p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800"
              >
                <p class="text-sm text-gray-900 dark:text-white">{{ item.content }}</p>
                <div class="flex items-center justify-between mt-2">
                  <span class="text-xs text-gray-500">{{ $t("retrospective.assignedTo") }}: {{ item.assignee }}</span>
                  <UBadge :color="item.priority === 'high' ? 'red' : item.priority === 'medium' ? 'yellow' : 'green'" variant="soft" size="xs">
                    {{ $t(`priority.${item.priority}`) }}
                  </UBadge>
                </div>
              </div>
            </div>
            <UTextarea
              v-model="newActionItem"
              :placeholder="$t('retrospective.addActionPlaceholder')"
              size="sm"
              @keydown.enter.shift.prevent="addActionItem"
            />
          </div>
        </div>
      </div>
    </UCard>

    <!-- 历史回顾 -->
    <UCard>
      <template #header>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ $t("retrospective.pastRetrospectives") }}
        </h2>
      </template>

      <div class="space-y-4">
        <div
          v-for="retrospective in pastRetrospectives"
          :key="retrospective.id"
          class="p-4 border border-gray-200 dark:border-gray-700 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors"
        >
          <div class="flex items-center justify-between mb-2">
            <h3 class="font-medium text-gray-900 dark:text-white">
              {{ retrospective.sprintName }}
            </h3>
            <span class="text-sm text-gray-500">
              {{ formatDate(retrospective.date) }}
            </span>
          </div>
          <div class="grid grid-cols-3 gap-4 text-sm">
            <div>
              <span class="text-green-600 dark:text-green-400">{{ $t("retrospective.wentWell") }}:</span>
              {{ retrospective.wentWellCount }}
            </div>
            <div>
              <span class="text-yellow-600 dark:text-yellow-400">{{ $t("retrospective.improvements") }}:</span>
              {{ retrospective.improvementCount }}
            </div>
            <div>
              <span class="text-blue-600 dark:text-blue-400">{{ $t("retrospective.actions") }}:</span>
              {{ retrospective.actionCount }}
            </div>
          </div>
        </div>

        <div v-if="pastRetrospectives.length === 0" class="text-center py-8">
          <UIcon name="i-heroicons-clock" class="w-8 h-8 text-gray-400 mx-auto mb-2" />
          <p class="text-gray-600 dark:text-gray-400">
            {{ $t("retrospective.noPastRetrospectives") }}
          </p>
        </div>
      </div>
    </UCard>
  </div>
</template>

<script setup>
// 页面元数据
definePageMeta({
  title: 'retrospective'
});

// 响应式数据
const activeRetrospective = ref(null);
const newWentWellItem = ref('');
const newImprovementItem = ref('');
const newActionItem = ref('');

// 示例数据
const wentWellItems = ref([
  {
    id: 1,
    content: '团队协作非常顺畅，每日站会很有效',
    author: '张三',
    votes: 5
  },
  {
    id: 2,
    content: '代码质量有明显提升，测试覆盖率达到80%',
    author: '李四',
    votes: 3
  }
]);

const improvementItems = ref([
  {
    id: 1,
    content: '需求变更太频繁，影响开发节奏',
    author: '王五',
    votes: 4
  },
  {
    id: 2,
    content: '测试环境不够稳定，影响测试效率',
    author: '赵六',
    votes: 2
  }
]);

const actionItems = ref([
  {
    id: 1,
    content: '建立需求变更流程，减少临时变更',
    assignee: '产品经理',
    priority: 'high'
  },
  {
    id: 2,
    content: '优化测试环境配置，提高稳定性',
    assignee: '运维团队',
    priority: 'medium'
  }
]);

const pastRetrospectives = ref([
  {
    id: 1,
    sprintName: 'Sprint 1',
    date: '2024-01-15',
    wentWellCount: 5,
    improvementCount: 3,
    actionCount: 2
  },
  {
    id: 2,
    sprintName: 'Sprint 2',
    date: '2024-01-29',
    wentWellCount: 4,
    improvementCount: 4,
    actionCount: 3
  }
]);

// 计算属性
const retrospectiveStatusColor = computed(() => {
  if (!activeRetrospective.value) return 'gray';
  return activeRetrospective.value.status === 'inProgress' ? 'green' : 'gray';
});

const retrospectiveStatusText = computed(() => {
  if (!activeRetrospective.value) return '';
  const { $t } = useNuxtApp();
  return $t(`retrospective.status.${activeRetrospective.value.status}`);
});

// 方法
const startRetrospective = () => {
  activeRetrospective.value = {
    id: Date.now(),
    status: 'inProgress',
    startedAt: new Date()
  };
};

const endRetrospective = () => {
  if (activeRetrospective.value) {
    activeRetrospective.value.status = 'completed';
    // 这里可以保存回顾结果
  }
};

const addWentWellItem = () => {
  if (newWentWellItem.value.trim()) {
    wentWellItems.value.push({
      id: Date.now(),
      content: newWentWellItem.value.trim(),
      author: '当前用户',
      votes: 0
    });
    newWentWellItem.value = '';
  }
};

const addImprovementItem = () => {
  if (newImprovementItem.value.trim()) {
    improvementItems.value.push({
      id: Date.now(),
      content: newImprovementItem.value.trim(),
      author: '当前用户',
      votes: 0
    });
    newImprovementItem.value = '';
  }
};

const addActionItem = () => {
  if (newActionItem.value.trim()) {
    actionItems.value.push({
      id: Date.now(),
      content: newActionItem.value.trim(),
      assignee: '待分配',
      priority: 'medium'
    });
    newActionItem.value = '';
  }
};

const formatTime = (date) => {
  return new Date(date).toLocaleTimeString();
};

const formatDate = (date) => {
  return new Date(date).toLocaleDateString();
};
</script>