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

      <div class="h-96">
        <ClientOnly>
          <div ref="chartContainer" class="h-full w-full" />
          <template #fallback>
            <div class="h-full flex items-center justify-center bg-gray-50 dark:bg-gray-800 rounded-lg">
              <div class="text-center">
                <UIcon name="i-heroicons-chart-bar" class="w-16 h-16 text-gray-400 mx-auto mb-4" />
                <p class="text-gray-600 dark:text-gray-400">{{ $t("common.loading") }}</p>
              </div>
            </div>
          </template>
        </ClientOnly>
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
import * as echarts from 'echarts';

// 页面元数据
definePageMeta({
  title: 'cumulativeFlow'
});

// 响应式数据
const selectedProject = ref(null);
const selectedTimeRange = ref('last30days');
const chartContainer = ref(null);
let chartInstance = null;

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
  { stage: '开发阶段', level: '轻微', severity: 'yellow', avgTime: 3.2 },
  { stage: '测试阶段', level: '严重', severity: 'red', avgTime: 4.8 },
  { stage: '代码审查', level: '正常', severity: 'green', avgTime: 1.5 }
]);

// 累积流图数据
const flowData = ref([
  { date: '2024-02-01', todo: 25, inProgress: 5, review: 3, testing: 2, done: 15 },
  { date: '2024-02-02', todo: 27, inProgress: 6, review: 4, testing: 2, done: 17 },
  { date: '2024-02-03', todo: 26, inProgress: 7, review: 5, testing: 3, done: 19 },
  { date: '2024-02-04', todo: 24, inProgress: 8, review: 4, testing: 4, done: 21 },
  { date: '2024-02-05', todo: 23, inProgress: 9, review: 5, testing: 3, done: 23 },
  { date: '2024-02-06', todo: 22, inProgress: 8, review: 6, testing: 4, done: 25 },
  { date: '2024-02-07', todo: 21, inProgress: 7, review: 5, testing: 5, done: 27 },
  { date: '2024-02-08', todo: 20, inProgress: 8, review: 4, testing: 4, done: 29 },
  { date: '2024-02-09', todo: 19, inProgress: 9, review: 5, testing: 3, done: 31 },
  { date: '2024-02-10', todo: 18, inProgress: 8, review: 6, testing: 4, done: 33 }
]);

// 主题配置
const colorMode = useColorMode();
const isDark = computed(() => colorMode.value === 'dark');

// 创建图表
const createChart = () => {
  if (!chartContainer.value) return;
  
  if (chartInstance) {
    chartInstance.dispose();
  }
  
  chartInstance = echarts.init(chartContainer.value);
  
  const dates = flowData.value.map(item => {
    return new Date(item.date).toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' });
  });

  // 主题颜色配置
  const textColor = isDark.value ? '#ffffff' : '#374151';
  const backgroundColor = isDark.value ? '#1f2937' : '#ffffff';
  const gridColor = isDark.value ? '#374151' : '#e5e7eb';

  // 计算累积数据
  const cumulativeData = flowData.value.map(item => ({
    todo: item.todo,
    inProgress: item.todo + item.inProgress,
    review: item.todo + item.inProgress + item.review,
    testing: item.todo + item.inProgress + item.review + item.testing,
    done: item.todo + item.inProgress + item.review + item.testing + item.done
  }));
  
  const option = {
    backgroundColor: backgroundColor,
    title: {
      text: '累积流图',
      left: 'center',
      textStyle: { fontSize: 16, fontWeight: 'bold', color: textColor }
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'line' },
      backgroundColor: isDark.value ? '#374151' : '#ffffff',
      borderColor: isDark.value ? '#6b7280' : '#d1d5db',
      textStyle: { color: textColor },
      formatter: function(params) {
        let result = params[0].axisValue + '<br/>';
        const originalData = flowData.value[params[0].dataIndex];
        result += `待办: ${originalData.todo}<br/>`;
        result += `进行中: ${originalData.inProgress}<br/>`;
        result += `待审核: ${originalData.review}<br/>`;
        result += `测试中: ${originalData.testing}<br/>`;
        result += `已完成: ${originalData.done}<br/>`;
        return result;
      }
    },
    legend: {
      data: ['待办', '进行中', '待审核', '测试中', '已完成'],
      top: 30,
      textStyle: { color: textColor }
    },
    grid: {
      left: '3%', right: '4%', bottom: '3%', top: '15%', containLabel: true
    },
    xAxis: {
      type: 'category', boundaryGap: false, data: dates,
      axisLabel: { rotate: 45, color: textColor },
      axisLine: { lineStyle: { color: gridColor } }
    },
    yAxis: {
      type: 'value', name: '任务数量',
      nameTextStyle: { color: textColor }, min: 0,
      axisLabel: { color: textColor },
      axisLine: { lineStyle: { color: gridColor } },
      splitLine: { lineStyle: { color: gridColor } }
    },
    series: [
      { name: '已完成', type: 'line', stack: 'total', data: cumulativeData.map(d => d.done), areaStyle: { color: '#10b981' }, lineStyle: { width: 0 }, symbol: 'none' },
      { name: '测试中', type: 'line', stack: 'total', data: cumulativeData.map(d => d.testing), areaStyle: { color: '#8b5cf6' }, lineStyle: { width: 0 }, symbol: 'none' },
      { name: '待审核', type: 'line', stack: 'total', data: cumulativeData.map(d => d.review), areaStyle: { color: '#f59e0b' }, lineStyle: { width: 0 }, symbol: 'none' },
      { name: '进行中', type: 'line', stack: 'total', data: cumulativeData.map(d => d.inProgress), areaStyle: { color: '#3b82f6' }, lineStyle: { width: 0 }, symbol: 'none' },
      { name: '待办', type: 'line', stack: 'total', data: cumulativeData.map(d => d.todo), areaStyle: { color: '#6b7280' }, lineStyle: { width: 0 }, symbol: 'none' }
    ]
  };
  
  chartInstance.setOption(option);
  
  // 响应窗口大小变化
  window.addEventListener('resize', () => {
    if (chartInstance) {
      chartInstance.resize();
    }
  });
};

// 监听主题变化
watch(isDark, () => {
  if (chartInstance) {
    createChart();
  }
});

// 方法
const exportChart = () => {
  if (chartInstance) {
    const url = chartInstance.getDataURL({ type: 'png', pixelRatio: 2, backgroundColor: '#fff' });
    const link = document.createElement('a');
    link.href = url;
    link.download = 'cumulative-flow-chart.png';
    link.click();
  }
};

const formatDate = (date) => new Date(date).toLocaleDateString();

// 生命周期
onMounted(() => {
  selectedProject.value = projectOptions.value[0].value;
  
  // 确保在下一个 tick 中创建图表，给 DOM 足够时间渲染
  nextTick(() => {
    setTimeout(() => {
      createChart();
    }, 100);
  });
});

onUnmounted(() => {
  if (chartInstance) {
    chartInstance.dispose();
    chartInstance = null;
  }
  window.removeEventListener('resize', () => {});
});
</script>