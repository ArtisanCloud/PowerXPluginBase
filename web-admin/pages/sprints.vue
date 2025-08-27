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
              Sprint 管理
            </h1>
          </div>
          <div class=\"flex items-center space-x-4\">
            <button class=\"bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-md text-sm font-medium\">
              创建 Sprint
            </button>
          </div>
        </div>
      </div>
    </header>

    <!-- Sprint 列表 -->
    <div class=\"max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6\">
      <div class=\"space-y-6\">
        <!-- Sprint 卡片 -->
        <div v-for=\"sprint in sprints\" :key=\"sprint.id\" class=\"bg-white shadow rounded-lg overflow-hidden\">
          <div class=\"px-6 py-4 border-b border-gray-200\">
            <div class=\"flex items-center justify-between\">
              <div class=\"flex items-center\">
                <h3 class=\"text-lg font-medium text-gray-900\">{{ sprint.name }}</h3>
                <span class=\"ml-3 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium\" :class=\"getStatusColor(sprint.status)\">
                  {{ getStatusText(sprint.status) }}
                </span>
              </div>
              <div class=\"flex items-center space-x-2\">
                <button v-if=\"sprint.status === 'planning'\" class=\"bg-green-600 hover:bg-green-700 text-white px-3 py-1 rounded text-sm\">
                  启动 Sprint
                </button>
                <button v-if=\"sprint.status === 'active'\" class=\"bg-blue-600 hover:bg-blue-700 text-white px-3 py-1 rounded text-sm\">
                  完成 Sprint
                </button>
                <button class=\"text-gray-400 hover:text-gray-600\">
                  <svg class=\"h-5 w-5\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\">
                    <path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M12 6.75a.75.75 0 110-1.5.75.75 0 010 1.5zM12 12.75a.75.75 0 110-1.5.75.75 0 010 1.5zM12 18.75a.75.75 0 110-1.5.75.75 0 010 1.5z\" />
                  </svg>
                </button>
              </div>
            </div>
            <p class=\"mt-1 text-sm text-gray-600\">{{ sprint.goal }}</p>
            <div class=\"mt-2 flex items-center text-sm text-gray-500 space-x-6\">
              <span>{{ formatDate(sprint.start_date) }} - {{ formatDate(sprint.end_date) }}</span>
              <span v-if=\"sprint.capacity\">容量: {{ sprint.capacity }} 故事点</span>
              <span v-if=\"sprint.duration\">持续时间: {{ sprint.duration }} 天</span>
            </div>
          </div>

          <!-- Sprint 统计信息 -->
          <div class=\"px-6 py-4\">
            <div class=\"grid grid-cols-1 gap-4 sm:grid-cols-4\">
              <div class=\"text-center\">
                <div class=\"text-2xl font-bold text-gray-900\">{{ sprint.stats.total_tasks }}</div>
                <div class=\"text-sm text-gray-500\">总任务数</div>
              </div>
              <div class=\"text-center\">
                <div class=\"text-2xl font-bold text-green-600\">{{ sprint.stats.completed_tasks }}</div>
                <div class=\"text-sm text-gray-500\">已完成</div>
              </div>
              <div class=\"text-center\">
                <div class=\"text-2xl font-bold text-blue-600\">{{ sprint.stats.total_points }}</div>
                <div class=\"text-sm text-gray-500\">总故事点</div>
              </div>
              <div class=\"text-center\">
                <div class=\"text-2xl font-bold text-purple-600\">{{ Math.round(sprint.stats.progress) }}%</div>
                <div class=\"text-sm text-gray-500\">完成进度</div>
              </div>
            </div>

            <!-- 进度条 -->
            <div class=\"mt-4\">
              <div class=\"flex justify-between text-sm text-gray-600 mb-1\">
                <span>Sprint 进度</span>
                <span>{{ Math.round(sprint.stats.progress) }}%</span>
              </div>
              <div class=\"w-full bg-gray-200 rounded-full h-2\">
                <div class=\"bg-blue-600 h-2 rounded-full\" :style=\"{ width: sprint.stats.progress + '%' }\"></div>
              </div>
            </div>

            <!-- 任务进度条 -->
            <div class=\"mt-3\">
              <div class=\"flex justify-between text-sm text-gray-600 mb-1\">
                <span>任务完成</span>
                <span>{{ sprint.stats.completed_tasks }} / {{ sprint.stats.total_tasks }}</span>
              </div>
              <div class=\"w-full bg-gray-200 rounded-full h-2\">
                <div class=\"bg-green-600 h-2 rounded-full\" :style=\"{ width: (sprint.stats.completed_tasks / sprint.stats.total_tasks * 100) + '%' }\"></div>
              </div>
            </div>
          </div>

          <!-- Sprint 任务 -->
          <div v-if=\"sprint.tasks && sprint.tasks.length\" class=\"border-t border-gray-200\">
            <div class=\"px-6 py-3 bg-gray-50\">
              <h4 class=\"text-sm font-medium text-gray-900\">Sprint 任务 ({{ sprint.tasks.length }})</h4>
            </div>
            <div class=\"divide-y divide-gray-200\">
              <div v-for=\"task in sprint.tasks.slice(0, 3)\" :key=\"task.id\" class=\"px-6 py-3\">
                <div class=\"flex items-center justify-between\">
                  <div class=\"flex items-center\">
                    <div class=\"h-2 w-2 rounded-full mr-3\" :class=\"getTaskStatusColor(task.status)\"></div>
                    <span class=\"text-sm text-gray-900\">{{ task.title }}</span>
                  </div>
                  <div class=\"flex items-center space-x-2\">
                    <span v-if=\"task.estimate\" class=\"text-xs text-gray-500\">{{ task.estimate }}pts</span>
                    <span class=\"inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium\" :class=\"getTaskPriorityColor(task.priority)\">
                      {{ getTaskPriorityText(task.priority) }}
                    </span>
                  </div>
                </div>
              </div>
              <div v-if=\"sprint.tasks.length > 3\" class=\"px-6 py-3 text-center\">
                <button class=\"text-sm text-indigo-600 hover:text-indigo-900\">
                  查看全部 {{ sprint.tasks.length }} 个任务
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
// 页面元数据
useSeoMeta({
  title: 'PowerX Scrum - Sprint管理',
  description: 'Sprint管理和规划'
})

// 获取运行时配置
const config = useRuntimeConfig()

// 示例Sprint数据
const sprints = ref([
  {
    id: 2,
    name: 'Sprint 2 - 高级功能',
    goal: '实现 Sprint 管理和报告功能',
    start_date: '2024-12-27',
    end_date: '2025-01-10',
    capacity: 45,
    status: 'active',
    duration: 14,
    stats: {
      total_tasks: 4,
      completed_tasks: 1,
      total_points: 47,
      completed_points: 21,
      progress: 65
    },
    tasks: [
      {
        id: 1,
        title: '实现 Sprint 管理功能',
        status: 'in_progress',
        priority: 'high',
        estimate: 21
      },
      {
        id: 2,
        title: '设计任务看板界面',
        status: 'todo',
        priority: 'medium',
        estimate: 13
      },
      {
        id: 3,
        title: '开发燃尽图功能',
        status: 'todo',
        priority: 'low',
        estimate: 8
      },
      {
        id: 4,
        title: '编写API文档',
        status: 'todo',
        priority: 'low',
        estimate: 5
      }
    ]
  },
  {
    id: 1,
    name: 'Sprint 1 - MVP 基础功能',
    goal: '完成用户管理和基础任务功能',
    start_date: '2024-12-13',
    end_date: '2024-12-27',
    capacity: 40,
    status: 'completed',
    duration: 14,
    stats: {
      total_tasks: 3,
      completed_tasks: 3,
      total_points: 26,
      completed_points: 26,
      progress: 100
    },
    tasks: [
      {
        id: 4,
        title: '设计用户注册登录界面',
        status: 'done',
        priority: 'high',
        estimate: 8
      },
      {
        id: 5,
        title: '实现用户认证API',
        status: 'done',
        priority: 'high',
        estimate: 13
      },
      {
        id: 6,
        title: '创建任务基础CRUD功能',
        status: 'done',
        priority: 'medium',
        estimate: 5
      }
    ]
  },
  {
    id: 3,
    name: 'Sprint 3 - 优化和集成',
    goal: '性能优化和第三方集成',
    start_date: '2025-01-10',
    end_date: '2025-01-24',
    capacity: 50,
    status: 'planning',
    duration: 14,
    stats: {
      total_tasks: 2,
      completed_tasks: 0,
      total_points: 21,
      completed_points: 0,
      progress: 0
    },
    tasks: [
      {
        id: 7,
        title: '性能优化和缓存实现',
        status: 'todo',
        priority: 'medium',
        estimate: 13
      },
      {
        id: 8,
        title: '集成第三方通知服务',
        status: 'todo',
        priority: 'low',
        estimate: 8
      }
    ]
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

const getTaskStatusColor = (status) => {
  const colors = {
    'todo': 'bg-gray-400',
    'in_progress': 'bg-blue-400',
    'done': 'bg-green-400'
  }
  return colors[status] || 'bg-gray-400'
}

const getTaskPriorityColor = (priority) => {
  const colors = {
    'urgent': 'bg-red-100 text-red-800',
    'high': 'bg-orange-100 text-orange-800',
    'medium': 'bg-yellow-100 text-yellow-800',
    'low': 'bg-green-100 text-green-800'
  }
  return colors[priority] || 'bg-gray-100 text-gray-800'
}

const getTaskPriorityText = (priority) => {
  const texts = {
    'urgent': '紧急',
    'high': '高',
    'medium': '中',
    'low': '低'
  }
  return texts[priority] || priority
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('zh-CN')
}

// 这里可以添加API调用逻辑
// const { data: sprintsData } = await $fetch(`${config.public.apiBaseUrl}/sprints`)
</script>