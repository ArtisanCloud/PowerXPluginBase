<template>
  <UCard class="task-card cursor-pointer hover:shadow-md transition-all">
    <div class="space-y-3">
      <!-- 任务标题和ID -->
      <div class="flex items-start justify-between">
        <div class="flex-1">
          <div class="flex items-center space-x-2 mb-1">
            <UBadge size="xs" variant="outline">
              {{ task.id }}
            </UBadge>
            <UBadge
              :color="getPriorityColor(task.priority)"
              size="xs"
              variant="soft"
            >
              {{ $t(`priority.${task.priority}`) }}
            </UBadge>
          </div>
          <h4 class="text-sm font-medium text-gray-900 dark:text-white line-clamp-2">
            {{ task.title }}
          </h4>
        </div>
        <UDropdownMenu :items="actionItems" :popper="{ placement: 'bottom-end' }">
          <UButton
            icon="i-heroicons-ellipsis-vertical"
            size="xs"
            variant="ghost"
            @click.stop
          />
        </UDropdownMenu>
      </div>

      <!-- 任务描述 -->
      <p class="text-xs text-gray-600 dark:text-gray-400 line-clamp-2">
        {{ task.description }}
      </p>

      <!-- 标签 -->
      <div v-if="task.labels && task.labels.length > 0" class="flex flex-wrap gap-1">
        <UBadge
          v-for="label in task.labels.slice(0, 2)"
          :key="label"
          size="xs"
          variant="soft"
          color="blue"
        >
          {{ label }}
        </UBadge>
        <UBadge
          v-if="task.labels.length > 2"
          size="xs"
          variant="outline"
        >
          +{{ task.labels.length - 2 }}
        </UBadge>
      </div>

      <!-- 底部信息 -->
      <div class="flex items-center justify-between pt-2 border-t border-gray-100 dark:border-gray-700">
        <div class="flex items-center space-x-2">
          <UAvatar
            v-if="task.assignee"
            :alt="task.assignee"
            size="xs"
            class="bg-primary-500"
          >
            <span class="text-xs font-medium text-white">
              {{ getInitials(task.assignee) }}
            </span>
          </UAvatar>
          <span v-else class="text-xs text-gray-400">
            {{ $t("common.unassigned") }}
          </span>
        </div>
        <div class="flex items-center space-x-2">
          <UIcon name="i-heroicons-chart-bar" class="h-3 w-3 text-gray-400" />
          <span class="text-xs text-gray-500 dark:text-gray-400">
            {{ task.storyPoints }}
          </span>
        </div>
      </div>
    </div>
  </UCard>
</template>

<script setup>
// Props
const props = defineProps({
  task: {
    type: Object,
    required: true
  }
});

// Emits
const emit = defineEmits(['edit', 'delete']);

// 国际化
const { t } = useI18n();

// 操作菜单项
const actionItems = computed(() => [
  [{
    label: t('common.edit'),
    icon: 'i-heroicons-pencil',
    click: () => emit('edit', props.task)
  }],
  [{
    label: t('common.delete'),
    icon: 'i-heroicons-trash',
    click: () => emit('delete', props.task)
  }]
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

// 获取姓名首字母
const getInitials = (name) => {
  return name
    .split(' ')
    .map(word => word.charAt(0))
    .join('')
    .toUpperCase()
    .slice(0, 2);
};
</script>