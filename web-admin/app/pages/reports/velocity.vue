<template>
  <div class="space-y-6">
    <!-- 页面头部 -->
    <div class="border-b border-gray-200 dark:border-gray-700 pb-4">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
        {{ $t("navigation.velocityChart") }}
      </h1>
      <p class="text-gray-600 dark:text-gray-400 mt-1">
        {{ $t("reports.velocityDescription") }}
      </p>
    </div>

    <!-- 过滤器 -->
    <UCard>
      <div class="flex flex-wrap gap-4">
        <USelectMenu
          v-model="selectedTeam"
          :options="teamOptions"
          placeholder="选择团队"
          class="w-48"
        />
        <USelectMenu
          v-model="selectedPeriod"
          :options="periodOptions"
          placeholder="时间周期"
          class="w-48"
        />
        <UButton variant="outline" @click="exportChart">
          {{ $t("reports.exportChart") }}
        </UButton>
      </div>
    </UCard>

    <!-- 速度图表 -->
    <UCard>
      <template #header>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ $t("reports.teamVelocityTrend") }}
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

    <!-- 速度统计 -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <UCard>
        <div class="text-center">
          <div class="text-2xl font-bold text-blue-600 dark:text-blue-400">{{ averageVelocity }}</div>
          <div class="text-sm text-gray-600 dark:text-gray-400">{{ $t("reports.averageVelocity") }}</div>
        </div>
      </UCard>
      <UCard>
        <div class="text-center">
          <div class="text-2xl font-bold text-green-600 dark:text-green-400">{{ currentVelocity }}</div>
          <div class="text-sm text-gray-600 dark:text-gray-400">{{ $t("reports.currentVelocity") }}</div>
        </div>
      </UCard>
      <UCard>
        <div class="text-center">
          <div class="text-2xl font-bold text-yellow-600 dark:text-yellow-400">{{ velocityTrend }}%</div>
          <div class="text-sm text-gray-600 dark:text-gray-400">{{ $t("reports.velocityTrend") }}</div>
        </div>
      </UCard>
      <UCard>
        <div class="text-center">
          <div class="text-2xl font-bold text-purple-600 dark:text-purple-400">{{ predictedVelocity }}</div>
          <div class="text-sm text-gray-600 dark:text-gray-400">{{ $t("reports.predictedVelocity") }}</div>
        </div>
      </UCard>
    </div>

    <!-- 历史速度数据 -->
    <UCard>
      <template #header>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ $t("reports.velocityHistory") }}
        </h2>
      </template>

      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
          <thead class="bg-gray-50 dark:bg-gray-800">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                {{ $t("sprint.name") }}
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                {{ $t("reports.plannedPoints") }}
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                {{ $t("reports.completedPoints") }}
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                {{ $t("reports.velocity") }}
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                {{ $t("sprint.endDate") }}
              </th>
            </tr>
          </thead>
          <tbody class="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-700">
            <tr v-for="sprint in velocityHistory" :key="sprint.id">
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">
                {{ sprint.name }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-400">
                {{ sprint.plannedPoints }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-400">
                {{ sprint.completedPoints }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-400">
                <UBadge :color="getVelocityColor(sprint.velocity)" variant="soft">
                  {{ sprint.velocity }}
                </UBadge>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-400">
                {{ formatDate(sprint.endDate) }}
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
  title: 'velocityChart'
});

// 响应式数据
const selectedTeam = ref(null);
const selectedPeriod = ref('last6sprints');
const chartContainer = ref(null);
let chartInstance = null;

// 选项数据
const teamOptions = ref([
  { label: '开发团队A', value: 'team-a' },
  { label: '开发团队B', value: 'team-b' },
  { label: '全部团队', value: 'all' }
]);

const periodOptions = ref([
  { label: '最近6个冲刺', value: 'last6sprints' },
  { label: '最近12个冲刺', value: 'last12sprints' },
  { label: '本季度', value: 'thisQuarter' }
]);

// 统计数据
const averageVelocity = ref(45);
const currentVelocity = ref(52);
const velocityTrend = ref(15);
const predictedVelocity = ref(48);

// 历史数据
const velocityHistory = ref([
  { id: 1, name: 'Sprint 6', plannedPoints: 50, completedPoints: 52, velocity: 52, endDate: '2024-02-15' },
  { id: 2, name: 'Sprint 5', plannedPoints: 45, completedPoints: 43, velocity: 43, endDate: '2024-02-01' },
  { id: 3, name: 'Sprint 4', plannedPoints: 48, completedPoints: 47, velocity: 47, endDate: '2024-01-18' },
  { id: 4, name: 'Sprint 3', plannedPoints: 40, completedPoints: 38, velocity: 38, endDate: '2024-01-04' },
  { id: 5, name: 'Sprint 2', plannedPoints: 42, completedPoints: 41, velocity: 41, endDate: '2023-12-21' },
  { id: 6, name: 'Sprint 1', plannedPoints: 35, completedPoints: 33, velocity: 33, endDate: '2023-12-07' }
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
  const sprints = velocityHistory.value.slice().reverse();
  
  // 主题颜色配置
  const textColor = isDark.value ? '#ffffff' : '#374151';
  const backgroundColor = isDark.value ? '#1f2937' : '#ffffff';
  const gridColor = isDark.value ? '#374151' : '#e5e7eb';
  
  const option = {
    backgroundColor: backgroundColor,
    title: {
      text: '团队速度趋势',
      left: 'center',
      textStyle: { fontSize: 16, fontWeight: 'bold', color: textColor }
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'cross', crossStyle: { color: '#999' } },
      backgroundColor: isDark.value ? '#374151' : '#ffffff',
      borderColor: isDark.value ? '#6b7280' : '#d1d5db',
      textStyle: { color: textColor }
    },
    legend: {
      data: ['计划点数', '完成点数', '速度趋势'],
      top: 30,
      textStyle: { color: textColor }
    },
    xAxis: [{
      type: 'category',
      data: sprints.map(s => s.name),
      axisPointer: { type: 'shadow' },
      axisLabel: { color: textColor },
      axisLine: { lineStyle: { color: gridColor } }
    }],
    yAxis: [
      {
        type: 'value',
        name: '故事点数',
        nameTextStyle: { color: textColor },
        min: 0, max: 60, interval: 10,
        axisLabel: { formatter: '{value} 点', color: textColor },
        axisLine: { lineStyle: { color: gridColor } },
        splitLine: { lineStyle: { color: gridColor } }
      },
      {
        type: 'value',
        name: '速度',
        nameTextStyle: { color: textColor },
        min: 0, max: 60, interval: 10,
        axisLabel: { formatter: '{value}', color: textColor },
        axisLine: { lineStyle: { color: gridColor } },
        splitLine: { lineStyle: { color: gridColor } }
      }
    ],
    series: [
      {
        name: '计划点数',
        type: 'bar',
        data: sprints.map(s => s.plannedPoints),
        itemStyle: { color: '#93c5fd' }
      },
      {
        name: '完成点数',
        type: 'bar',
        data: sprints.map(s => s.completedPoints),
        itemStyle: { color: '#3b82f6' }
      },
      {
        name: '速度趋势',
        type: 'line',
        yAxisIndex: 1,
        data: sprints.map(s => s.velocity),
        lineStyle: { color: '#10b981', width: 3 },
        itemStyle: { color: '#10b981' },
        symbol: 'circle',
        symbolSize: 8
      }
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
    link.download = 'velocity-chart.png';
    link.click();
  }
};

const getVelocityColor = (velocity) => {
  if (velocity >= 50) return 'green';
  if (velocity >= 40) return 'yellow';
  return 'red';
};

const formatDate = (date) => new Date(date).toLocaleDateString();

// 生命周期
onMounted(() => {
  selectedTeam.value = teamOptions.value[0].value;
  
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