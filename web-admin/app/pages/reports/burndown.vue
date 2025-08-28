<template>
  <div class="space-y-6">
    <!-- 页面头部 -->
    <div class="border-b border-gray-200 dark:border-gray-700 pb-4">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
        {{ $t("navigation.burndownChart") }}
      </h1>
      <p class="text-gray-600 dark:text-gray-400 mt-1">
        {{ $t("reports.burndownDescription") }}
      </p>
    </div>

    <!-- 过滤器 -->
    <UCard>
      <div class="flex flex-wrap gap-4">
        <USelectMenu
          v-model="selectedSprint"
          :options="sprintOptions"
          placeholder="选择冲刺"
          class="w-48"
        />
        <USelectMenu
          v-model="selectedTimeRange"
          :options="timeRangeOptions"
          placeholder="时间范围"
          class="w-48"
        />
        <UButton variant="outline" @click="exportChart">
          {{ $t("reports.exportChart") }}
        </UButton>
      </div>
    </UCard>

    <!-- 燃尽图 -->
    <UCard>
      <template #header>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ $t("reports.burndownChart") }}
        </h2>
      </template>

      <div class="h-96 flex items-center justify-center bg-gray-50 dark:bg-gray-800 rounded-lg">
        <div class="text-center">
          <UIcon name="i-heroicons-chart-bar" class="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">
            {{ $t("reports.chartPlaceholder") }}
          </h3>
          <p class="text-gray-600 dark:text-gray-400">
            {{ $t("reports.chartPlaceholderDescription") }}
          </p>
        </div>
      </div>
    </UCard>

    <!-- 统计信息 -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <UCard>
        <div class="text-center">
          <div class="text-2xl font-bold text-blue-600 dark:text-blue-400">{{ totalStoryPoints }}</div>
          <div class="text-sm text-gray-600 dark:text-gray-400">{{ $t("reports.totalStoryPoints") }}</div>
        </div>
      </UCard>
      <UCard>
        <div class="text-center">
          <div class="text-2xl font-bold text-green-600 dark:text-green-400">{{ completedStoryPoints }}</div>
          <div class="text-sm text-gray-600 dark:text-gray-400">{{ $t("reports.completedStoryPoints") }}</div>
        </div>
      </UCard>
      <UCard>
        <div class="text-center">
          <div class="text-2xl font-bold text-yellow-600 dark:text-yellow-400">{{ remainingStoryPoints }}</div>
          <div class="text-sm text-gray-600 dark:text-gray-400">{{ $t("reports.remainingStoryPoints") }}</div>
        </div>
      </UCard>
      <UCard>
        <div class="text-center">
          <div class="text-2xl font-bold text-purple-600 dark:text-purple-400">{{ burndownRate }}%</div>
          <div class="text-sm text-gray-600 dark:text-gray-400">{{ $t("reports.burndownRate") }}</div>
        </div>
      </UCard>
    </div>
  </div>
</template>

<script setup>
// 页面元数据
definePageMeta({
  title: 'burndownChart'
});

// 响应式数据
const selectedSprint = ref(null);
const selectedTimeRange = ref('currentSprint');

// 选项数据
const sprintOptions = ref([
  { label: 'Sprint 3 (当前)', value: 'sprint-3' },
  { label: 'Sprint 2', value: 'sprint-2' },
  { label: 'Sprint 1', value: 'sprint-1' }
]);

const timeRangeOptions = ref([
  { label: '当前冲刺', value: 'currentSprint' },
  { label: '最近30天', value: 'last30days' },
  { label: '最近3个月', value: 'last3months' }
]);

// 统计数据
const totalStoryPoints = ref(120);
const completedStoryPoints = ref(85);
const remainingStoryPoints = computed(() => totalStoryPoints.value - completedStoryPoints.value);
const burndownRate = computed(() => Math.round((completedStoryPoints.value / totalStoryPoints.value) * 100));

// 方法
const exportChart = () => {
  // 导出图表逻辑
  console.log('导出燃尽图');
};

// 初始化
onMounted(() => {
  selectedSprint.value = sprintOptions.value[0].value;
});
</script>