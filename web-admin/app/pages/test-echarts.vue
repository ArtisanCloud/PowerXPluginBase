<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold mb-4">ECharts 测试页面</h1>
    
    <UCard class="mb-6">
      <template #header>
        <h2>简单图表测试</h2>
      </template>
      
      <div class="h-64">
        <ClientOnly>
          <VChart 
            :option="simpleOption" 
            class="h-full w-full"
            autoresize
          />
          <template #fallback>
            <div class="h-full flex items-center justify-center bg-gray-100 dark:bg-gray-800">
              <p>图表加载中...</p>
            </div>
          </template>
        </ClientOnly>
      </div>
    </UCard>

    <div class="text-sm text-gray-600 dark:text-gray-400">
      <p>如果您能看到上方的柱状图，说明 ECharts 已正常工作。</p>
      <p>当前主题: {{ colorMode.value }}</p>
    </div>
  </div>
</template>

<script setup>
import { use } from 'echarts/core';
import { CanvasRenderer } from 'echarts/renderers';
import { BarChart } from 'echarts/charts';
import {
  TitleComponent,
  TooltipComponent,
  GridComponent
} from 'echarts/components';
import VChart from 'vue-echarts';

use([
  CanvasRenderer,
  BarChart,
  TitleComponent,
  TooltipComponent,
  GridComponent
]);

const colorMode = useColorMode();

const simpleOption = computed(() => {
  const isDark = colorMode.value === 'dark';
  const textColor = isDark ? '#ffffff' : '#374151';
  const backgroundColor = isDark ? '#1f2937' : '#ffffff';
  
  return {
    backgroundColor: backgroundColor,
    title: {
      text: 'ECharts 测试图表',
      textStyle: {
        color: textColor
      }
    },
    tooltip: {
      backgroundColor: isDark ? '#374151' : '#ffffff',
      textStyle: {
        color: textColor
      }
    },
    xAxis: {
      type: 'category',
      data: ['周一', '周二', '周三', '周四', '周五'],
      axisLabel: {
        color: textColor
      }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        color: textColor
      }
    },
    series: [{
      data: [120, 200, 150, 80, 70],
      type: 'bar',
      itemStyle: {
        color: '#3b82f6'
      }
    }]
  };
});
</script>