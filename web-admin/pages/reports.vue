<template>
  <div class=\"min-h-screen bg-gray-50\">
    <!-- 页面头部 -->
    <header class=\"bg-white shadow-sm border-b\">
      <div class=\"max-w-7xl mx-auto px-4 sm:px-6 lg:px-8\">
        <div class=\"flex justify-between items-center py-6\">
          <div class=\"flex items-center\">
            <NuxtLink to=\"/\" class=\"text-gray-500 hover:text-gray-900 mr-4\">
              ← 返回仪表板
            </NuxtLink>
            <h1 class=\"text-2xl font-bold text-gray-900\">
              报告与分析
            </h1>
          </div>
          <div class=\"flex items-center space-x-4\">
            <select class=\"block rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm\">
              <option value=\"current\">当前Sprint</option>
              <option value=\"last\">上个Sprint</option>
              <option value=\"all\">全部Sprint</option>
            </select>
            <button class=\"bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-md text-sm font-medium\">
              导出报告
            </button>
          </div>
        </div>
      </div>
    </header>

    <!-- 报告内容 -->
    <div class=\"max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6\">
      <!-- 关键指标 -->
      <div class=\"mb-8\">
        <h2 class=\"text-lg font-medium text-gray-900 mb-4\">关键指标</h2>
        <div class=\"grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4\">
          <div class=\"bg-white overflow-hidden shadow rounded-lg\">
            <div class=\"p-5\">
              <div class=\"flex items-center\">
                <div class=\"flex-shrink-0\">
                  <svg class=\"h-6 w-6 text-blue-400\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\">
                    <path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z\" />
                  </svg>
                </div>
                <div class=\"ml-5 w-0 flex-1\">
                  <dl>
                    <dt class=\"text-sm font-medium text-gray-500 truncate\">团队速度</dt>
                    <dd class=\"text-lg font-medium text-gray-900\">38 故事点</dd>
                    <dd class=\"text-xs text-green-600\">+12% 环比</dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>

          <div class=\"bg-white overflow-hidden shadow rounded-lg\">
            <div class=\"p-5\">
              <div class=\"flex items-center\">
                <div class=\"flex-shrink-0\">
                  <svg class=\"h-6 w-6 text-green-400\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\">
                    <path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z\" />
                  </svg>
                </div>
                <div class=\"ml-5 w-0 flex-1\">
                  <dl>
                    <dt class=\"text-sm font-medium text-gray-500 truncate\">Sprint完成率</dt>
                    <dd class=\"text-lg font-medium text-gray-900\">85%</dd>
                    <dd class=\"text-xs text-green-600\">+5% 环比</dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>

          <div class=\"bg-white overflow-hidden shadow rounded-lg\">
            <div class=\"p-5\">
              <div class=\"flex items-center\">
                <div class=\"flex-shrink-0\">
                  <svg class=\"h-6 w-6 text-yellow-400\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\">
                    <path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z\" />
                  </svg>
                </div>
                <div class=\"ml-5 w-0 flex-1\">
                  <dl>
                    <dt class=\"text-sm font-medium text-gray-500 truncate\">平均周期时间</dt>
                    <dd class=\"text-lg font-medium text-gray-900\">4.2 天</dd>
                    <dd class=\"text-xs text-red-600\">+0.5天 环比</dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>

          <div class=\"bg-white overflow-hidden shadow rounded-lg\">
            <div class=\"p-5\">
              <div class=\"flex items-center\">
                <div class=\"flex-shrink-0\">
                  <svg class=\"h-6 w-6 text-purple-400\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\">
                    <path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M16.5 18.75h-9m9 0a3 3 0 013 3h-15a3 3 0 013-3m9 0v-3.375c0-.621-.503-1.125-1.125-1.125h-.871M7.5 18.75v-3.375c0-.621.504-1.125 1.125-1.125h.872m5.007 0H9.497m5.007 0a7.454 7.454 0 01-.982-3.172M9.497 14.25a7.454 7.454 0 00.981-3.172M5.25 4.236c-.982.143-1.954.317-2.916.52A6.003 6.003 0 007.73 9.728M5.25 4.236V4.5c0 2.108.966 3.99 2.48 5.228M5.25 4.236V2.721C7.456 2.41 9.71 2.25 12 2.25c2.291 0 4.545.16 6.75.47v1.516M7.73 9.728a6.726 6.726 0 002.748 1.35m8.272-6.842V4.5c0 2.108-.966 3.99-2.48 5.228m2.48-5.228a25.628 25.628 0 012.916.52 6.003 6.003 0 00-5.395 4.972m0 0a6.726 6.726 0 01-2.749 1.35m0 0a6.772 6.772 0 01-3.044 0\" />
                  </svg>
                </div>
                <div class=\"ml-5 w-0 flex-1\">
                  <dl>
                    <dt class=\"text-sm font-medium text-gray-500 truncate\">质量得分</dt>
                    <dd class=\"text-lg font-medium text-gray-900\">92/100</dd>
                    <dd class=\"text-xs text-green-600\">+3 环比</dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 图表区域 -->
      <div class=\"grid grid-cols-1 gap-6 lg:grid-cols-2 mb-8\">
        <!-- 燃尽图 -->
        <div class=\"bg-white shadow rounded-lg p-6\">
          <h3 class=\"text-lg font-medium text-gray-900 mb-4\">Sprint 燃尽图</h3>
          <div class=\"h-64 flex items-center justify-center bg-gray-50 rounded\">
            <div class=\"text-center\">
              <svg class=\"mx-auto h-12 w-12 text-gray-400\" fill=\"none\" viewBox=\"0 0 24 24\" stroke=\"currentColor\">
                <path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z\" />
              </svg>
              <p class=\"mt-2 text-sm text-gray-500\">燃尽图显示区域</p>
              <p class=\"text-xs text-gray-400\">（需要图表库支持）</p>
            </div>
          </div>
        </div>

        <!-- 速度趋势图 -->
        <div class=\"bg-white shadow rounded-lg p-6\">
          <h3 class=\"text-lg font-medium text-gray-900 mb-4\">团队速度趋势</h3>
          <div class=\"h-64 flex items-center justify-center bg-gray-50 rounded\">
            <div class=\"text-center\">
              <svg class=\"mx-auto h-12 w-12 text-gray-400\" fill=\"none\" viewBox=\"0 0 24 24\" stroke=\"currentColor\">
                <path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M7 12l3-3 3 3 4-4\" />
                <path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M8 21l4-4 4 4\" />
                <path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M3 4h18\" />
                <path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M21 12H3\" />
              </svg>
              <p class=\"mt-2 text-sm text-gray-500\">速度趋势图显示区域</p>
              <p class=\"text-xs text-gray-400\">（需要图表库支持）</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 累积流图 -->
      <div class=\"bg-white shadow rounded-lg p-6 mb-8\">
        <h3 class=\"text-lg font-medium text-gray-900 mb-4\">累积流图</h3>
        <div class=\"h-80 flex items-center justify-center bg-gray-50 rounded\">
          <div class=\"text-center\">
            <svg class=\"mx-auto h-16 w-16 text-gray-400\" fill=\"none\" viewBox=\"0 0 24 24\" stroke=\"currentColor\">
              <path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M4 6h16M4 12h16M4 18h16\" />
            </svg>
            <p class=\"mt-2 text-sm text-gray-500\">累积流图显示区域</p>
            <p class=\"text-xs text-gray-400\">显示任务在不同状态的流动情况</p>
          </div>
        </div>
      </div>

      <!-- 详细统计表格 -->
      <div class=\"bg-white shadow rounded-lg\">
        <div class=\"px-6 py-4 border-b border-gray-200\">
          <h3 class=\"text-lg font-medium text-gray-900\">Sprint 详细统计</h3>
        </div>
        <div class=\"overflow-x-auto\">
          <table class=\"min-w-full divide-y divide-gray-200\">
            <thead class=\"bg-gray-50\">
              <tr>
                <th class=\"px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider\">Sprint</th>
                <th class=\"px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider\">状态</th>
                <th class=\"px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider\">计划故事点</th>
                <th class=\"px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider\">完成故事点</th>
                <th class=\"px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider\">完成率</th>
                <th class=\"px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider\">任务数</th>
                <th class=\"px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider\">团队速度</th>
              </tr>
            </thead>
            <tbody class=\"bg-white divide-y divide-gray-200\">
              <tr v-for=\"sprint in sprintStats\" :key=\"sprint.id\">
                <td class=\"px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900\">
                  {{ sprint.name }}
                </td>
                <td class=\"px-6 py-4 whitespace-nowrap\">
                  <span class=\"inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium\" :class=\"getStatusColor(sprint.status)\">
                    {{ getStatusText(sprint.status) }}
                  </span>
                </td>
                <td class=\"px-6 py-4 whitespace-nowrap text-sm text-gray-900\">{{ sprint.planned_points }}</td>
                <td class=\"px-6 py-4 whitespace-nowrap text-sm text-gray-900\">{{ sprint.completed_points }}</td>
                <td class=\"px-6 py-4 whitespace-nowrap text-sm text-gray-900\">
                  <div class=\"flex items-center\">
                    <div class=\"w-16 bg-gray-200 rounded-full h-2 mr-2\">
                      <div class=\"bg-green-600 h-2 rounded-full\" :style=\"{ width: sprint.completion_rate + '%' }\"></div>
                    </div>
                    <span>{{ sprint.completion_rate }}%</span>
                  </div>
                </td>
                <td class=\"px-6 py-4 whitespace-nowrap text-sm text-gray-900\">{{ sprint.total_tasks }}</td>
                <td class=\"px-6 py-4 whitespace-nowrap text-sm text-gray-900\">{{ sprint.velocity }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
// 页面元数据
useSeoMeta({
  title: 'PowerX Scrum - 报告分析',
  description: 'Scrum项目报告和数据分析'
})

// 获取运行时配置
const config = useRuntimeConfig()

// Sprint 统计数据
const sprintStats = ref([
  {
    id: 2,
    name: 'Sprint 2',
    status: 'active',
    planned_points: 47,
    completed_points: 21,
    completion_rate: 45,
    total_tasks: 4,
    velocity: 21
  },
  {
    id: 1,
    name: 'Sprint 1',
    status: 'completed',
    planned_points: 26,
    completed_points: 26,
    completion_rate: 100,
    total_tasks: 3,
    velocity: 26
  },
  {
    id: 0,
    name: 'Sprint 0',
    status: 'completed',
    planned_points: 35,
    completed_points: 32,
    completion_rate: 91,
    total_tasks: 5,
    velocity: 32
  }
])

// 工具函数
const getStatusColor = (status) => {
  const colors = {
    'planning': 'bg-gray-100 text-gray-800',
    'active': 'bg-blue-100 text-blue-800',
    'completed': 'bg-green-100 text-green-800'
  }
  return colors[status] || 'bg-gray-100 text-gray-800'
}

const getStatusText = (status) => {
  const texts = {
    'planning': '计划中',
    'active': '进行中',
    'completed': '已完成'
  }
  return texts[status] || status
}

// 这里可以添加API调用逻辑
// const { data: reportData } = await $fetch(`${config.public.apiBaseUrl}/reports`)
</script>