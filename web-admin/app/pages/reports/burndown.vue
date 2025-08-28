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
import * as echarts from 'echarts';

// 页面元数据
definePageMeta({
  title: 'burndownChart'
});

// 响应式数据
const selectedSprint = ref(null);
const selectedTimeRange = ref('currentSprint');
const chartContainer = ref(null);
let chartInstance = null;

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

// 燃尽图数据
const burndownData = ref({
  dates: ['2024-02-01', '2024-02-02', '2024-02-03', '2024-02-04', '2024-02-05', '2024-02-06', '2024-02-07', '2024-02-08', '2024-02-09', '2024-02-10'],
  ideal: [120, 108, 96, 84, 72, 60, 48, 36, 24, 12],
  actual: [120, 115, 105, 95, 88, 78, 65, 52, 38, 25],
  remaining: [120, 115, 105, 95, 88, 78, 65, 52, 38, 25]
});

// 主题配置
const colorMode = useColorMode();
const isDark = computed(() => colorMode.value === 'dark');

// 创建图表
const createChart = () => {
  if (!chartContainer.value) return;
  
  // 销毁现有图表
  if (chartInstance) {
    chartInstance.dispose();
  }
  
  // 创建新图表
  chartInstance = echarts.init(chartContainer.value);
  
  // 主题颜色配置
  const textColor = isDark.value ? '#ffffff' : '#374151';
  const backgroundColor = isDark.value ? '#1f2937' : '#ffffff';
  const gridColor = isDark.value ? '#374151' : '#e5e7eb';
  
  const option = {
    backgroundColor: backgroundColor,
    title: {
      text: '燃尽图',
      left: 'center',
      textStyle: {
        fontSize: 16,
        fontWeight: 'bold',
        color: textColor
      }
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      },
      backgroundColor: isDark.value ? '#374151' : '#ffffff',
      borderColor: isDark.value ? '#6b7280' : '#d1d5db',
      textStyle: {
        color: textColor
      },
      formatter: function(params) {
        let result = params[0].axisValue + '<br/>';
        params.forEach(param => {
          result += `${param.marker}${param.seriesName}: ${param.value} 点<br/>`;
        });
        return result;
      }
    },
    legend: {
      data: ['理想燃尽', '实际燃尽', '剩余工作'],
      top: 30,
      textStyle: {
        color: textColor
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      top: '15%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: burndownData.value.dates.map(date => {
        return new Date(date).toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' });
      }),
      axisLabel: {
        rotate: 45,
        color: textColor
      },
      axisLine: {
        lineStyle: {
          color: gridColor
        }
      }
    },
    yAxis: {
      type: 'value',
      name: '故事点数',
      nameTextStyle: {
        color: textColor
      },
      min: 0,
      max: Math.max(...burndownData.value.ideal) * 1.1,
      axisLabel: {
        color: textColor
      },
      axisLine: {
        lineStyle: {
          color: gridColor
        }
      },
      splitLine: {
        lineStyle: {
          color: gridColor
        }
      }
    },
    series: [
      {
        name: '理想燃尽',
        type: 'line',
        data: burndownData.value.ideal,
        lineStyle: {
          color: '#10b981',
          type: 'dashed'
        },
        itemStyle: {
          color: '#10b981'
        },
        symbol: 'circle',
        symbolSize: 6
      },
      {
        name: '实际燃尽',
        type: 'line',
        data: burndownData.value.actual,
        lineStyle: {
          color: '#3b82f6',
          width: 3
        },
        itemStyle: {
          color: '#3b82f6'
        },
        symbol: 'circle',
        symbolSize: 8,
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [{
              offset: 0, color: 'rgba(59, 130, 246, 0.3)'
            }, {
              offset: 1, color: 'rgba(59, 130, 246, 0.05)'
            }]
          }
        }
      },
      {
        name: '剩余工作',
        type: 'line',
        data: burndownData.value.remaining,
        lineStyle: {
          color: '#f59e0b'
        },
        itemStyle: {
          color: '#f59e0b'
        },
        symbol: 'diamond',
        symbolSize: 6
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
    const url = chartInstance.getDataURL({
      type: 'png',
      pixelRatio: 2,
      backgroundColor: '#fff'
    });
    
    const link = document.createElement('a');
    link.href = url;
    link.download = 'burndown-chart.png';
    link.click();
  }
};

// 生命周期
onMounted(() => {
  selectedSprint.value = sprintOptions.value[0].value;
  
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