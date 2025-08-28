<template>
  <div class="space-y-6">
    <!-- 页面头部 -->
    <div class="border-b border-gray-200 dark:border-gray-700 pb-4">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
        {{ $t("navigation.cumulativeFlow") }}
      </h1>
      <p class="text-gray-600 dark:text-gray-400 mt-1">
        {{ $t("reports.cumulativeFlowDescription") }}
      </p>
    </div>

    <!-- 过滤器 -->
    <UCard>
      <div class="flex flex-wrap gap-4">
        <USelectMenu
          v-model="selectedProject"
          :options="projectOptions"
          placeholder="选择项目"
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

    <!-- 累积流图 -->
    <UCard>
      <template #header>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ $t("reports.cumulativeFlowDiagram") }}
        </h2>
      </template>

      <div class="h-96 flex items-center justify-center bg-gray-50 dark:bg-gray-800 rounded-lg">
        <div class="text-center">
          <UIcon name="i-heroicons-chart-bar" class="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">
            {{ $t("reports.cumulativeFlowPlaceholder") }}
          </h3>
          <p class="text-gray-600 dark:text-gray-400">
            {{ $t("reports.cumulativeFlowChartDescription") }}
          </p>
        </div>
      </div>
    </UCard>

    <!-- 状态分布统计 -->
    <div class="grid grid-cols-1 md:grid-cols-5 gap-6">
      <UCard v-for="(status, index) in statusDistribution" :key="index">
        <div class="text-center">
          <div class="text-2xl font-bold" :class="status.color">{{ status.count }}</div>
          <div class="text-sm text-gray-600 dark:text-gray-400">{{ status.name }}</div>
        </div>
      </UCard>
    </div>

    <!-- 流量指标 -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <UCard>
        <template #header>
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ $t("reports.flowMetrics") }}
          </h3>
        </template>

        <div class="space-y-4">
          <div class="flex justify-between items-center">
            <span class="text-gray-600 dark:text-gray-400">{{ $t("reports.averageCycleTime") }}</span>
            <span class="font-semibold text-gray-900 dark:text-white">{{ averageCycleTime }} {{ $t("common.days") }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-gray-600 dark:text-gray-400">{{ $t("reports.averageLeadTime") }}</span>
            <span class="font-semibold text-gray-900 dark:text-white">{{ averageLeadTime }} {{ $t("common.days") }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-gray-600 dark:text-gray-400">{{ $t("reports.throughput") }}</span>
            <span class="font-semibold text-gray-900 dark:text-white">{{ throughput }} {{ $t("reports.itemsPerWeek") }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-gray-600 dark:text-gray-400">{{ $t("reports.workInProgress") }}</span>
            <span class="font-semibold text-gray-900 dark:text-white">{{ workInProgress }}</span>
          </div>
        </div>
      </UCard>

      <UCard>
        <template #header>
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ $t("reports.bottleneckAnalysis") }}
          </h3>
        </template>

        <div class="space-y-4">
          <div v-for="bottleneck in bottlenecks" :key="bottleneck.stage" class="flex justify-between items-center">
            <span class="text-gray-600 dark:text-gray-400">{{ bottleneck.stage }}</span>
            <div class="flex items-center space-x-2">
              <UBadge :color="bottleneck.severity" variant="soft">
                {{ bottleneck.level }}
              </UBadge>
              <span class="text-sm text-gray-500">{{ bottleneck.avgTime }}天</span>
            </div>
          </div>
        </div>
      </UCard>
    </div>

    <!-- 详细数据表格 -->
    <UCard>
      <template #header>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ $t("reports.flowDetails") }}
        </h2>
      </template>

      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
          <thead class="bg-gray-50 dark:bg-gray-800">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                {{ $t("common.date") }}
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                {{ $t("status.todo") }}
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                {{ $t("status.inProgress") }}
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                {{ $t("status.review") }}
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                {{ $t("status.done") }}
              </th>
            </tr>
          </thead>
          <tbody class="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-700">
            <tr v-for="data in flowData" :key="data.date">
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">
                {{ formatDate(data.date) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-400">
                {{ data.todo }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-400">
                {{ data.inProgress }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-400">
                {{ data.review }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-400">
                {{ data.done }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>
  </div>
</template>

<script setup>
// 页面元数据
definePageMeta({
  title: 'cumulativeFlow'
});

// 响应式数据
const selectedProject = ref(null);
const selectedTimeRange = ref('last30days');

// 选项数据
const projectOptions = ref([
  { label: '项目A', value: 'project-a' },
  { label: '项目B', value: 'project-b' },
  { label: '全部项目', value: 'all' }
]);

const timeRangeOptions = ref([
  { label: '最近30天', value: 'last30days' },
  { label: '最近60天', value: 'last60days' },
  { label: '最近90天', value: 'last90days' }
]);

// 统计数据
const statusDistribution = ref([
  { name: '待办', count: 15, color: 'text-gray-600' },
  { name: '进行中', count: 8, color: 'text-blue-600' },
  { name: '待审核', count: 5, color: 'text-yellow-600' },
  { name: '测试中', count: 3, color: 'text-purple-600' },
  { name: '已完成', count: 42, color: 'text-green-600' }
]);

const averageCycleTime = ref(5.2);
const averageLeadTime = ref(8.1);
const throughput = ref(12);
const workInProgress = ref(16);

const bottlenecks = ref([
  {
    stage: '开发阶段',
    level: '轻微',
    severity: 'yellow',
    avgTime: 3.2
  },
  {
    stage: '测试阶段',
    level: '严重',
    severity: 'red',
    avgTime: 4.8
  },
  {
    stage: '代码审查',
    level: '正常',
    severity: 'green',
    avgTime: 1.5
  }
]);

const flowData = ref([
  {
    date: '2024-02-15',
    todo: 20,
    inProgress: 8,
    review: 5,
    done: 35
  },
  {
    date: '2024-02-14',
    todo: 22,
    inProgress: 7,
    review: 6,
    done: 33
  },
  {
    date: '2024-02-13',
    todo: 25,
    inProgress: 6,
    review: 4,
    done: 31
  }
]);

// 方法
const exportChart = () => {
  console.log('导出累积流图');
};

const formatDate = (date) => {
  return new Date(date).toLocaleDateString();
};

// 初始化
onMounted(() => {
  selectedProject.value = projectOptions.value[0].value;
});
</script>