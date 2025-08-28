<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold mb-4">ECharts 测试页面</h1>
    
    <div class="mb-4">
      <button class="bg-blue-500 text-white px-4 py-2 rounded" @click="createChart">
        创建图表
      </button>
    </div>
    
    <div ref="chartContainer" class="w-full h-96 border border-gray-300 rounded-lg"/>
    
    <div class="mt-4">
      <p class="text-sm text-gray-600">
        状态: {{ chartStatus }}
      </p>
    </div>
  </div>
</template>

<script setup>
import * as echarts from 'echarts';

const chartContainer = ref(null);
const chartStatus = ref('未初始化');
let chartInstance = null;

const createChart = () => {
  try {
    chartStatus.value = '正在创建图表...';
    
    if (!chartContainer.value) {
      chartStatus.value = '错误: 容器元素未找到';
      return;
    }
    
    if (chartInstance) {
      chartInstance.dispose();
    }
    
    chartInstance = echarts.init(chartContainer.value);
    
    const option = {
      title: {
        text: 'ECharts 测试图表'
      },
      tooltip: {},
      xAxis: {
        data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
      },
      yAxis: {},
      series: [{
        name: '测试数据',
        type: 'bar',
        data: [120, 200, 150, 80, 70, 110, 130]
      }]
    };
    
    chartInstance.setOption(option);
    chartStatus.value = '图表创建成功!';
    
  } catch (error) {
    console.error('创建图表时出错:', error);
    chartStatus.value = `错误: ${error.message}`;
  }
};

onMounted(() => {
  chartStatus.value = 'ECharts 已导入，等待创建图表';
});

onUnmounted(() => {
  if (chartInstance) {
    chartInstance.dispose();
  }
});
</script>