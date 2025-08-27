<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 页面头部 -->
    <header class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center py-6">
          <div class="flex items-center">
            <NuxtLink to="/" class="text-gray-500 hover:text-gray-900 mr-4">
              ← 返回仪表板
            </NuxtLink>
            <h1 class="text-2xl font-bold text-gray-900">
              任务管理
            </h1>
          </div>
          <div class="flex items-center space-x-4">
            <button class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-md text-sm font-medium">
              创建任务
            </button>
          </div>
        </div>
      </div>
    </header>

    <!-- 过滤器和搜索 -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <div class="bg-white p-4 rounded-lg shadow-sm mb-6">
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">状态</label>
            <select class="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
              <option value="">全部状态</option>
              <option value="todo">待办</option>
              <option value="in_progress">进行中</option>
              <option value="done">已完成</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">优先级</label>
            <select class="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
              <option value="">全部优先级</option>
              <option value="urgent">紧急</option>
              <option value="high">高</option>
              <option value="medium">中</option>
              <option value="low">低</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">分配人</label>
            <select class="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
              <option value="">全部成员</option>
              <option value="101">John Doe</option>
              <option value="102">Alice Smith</option>
              <option value="103">Mike Brown</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">搜索</label>
            <input type="text" placeholder="搜索任务..." class="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm">
          </div>
        </div>
      </div>

      <!-- 任务列表 -->
      <div class="bg-white shadow overflow-hidden sm:rounded-md">
        <ul role="list" class="divide-y divide-gray-200">
          <!-- 任务项目 -->
          <li v-for="task in tasks" :key="task.id" class="px-6 py-4 hover:bg-gray-50">
            <div class="flex items-center justify-between">
              <div class="flex items-center flex-1">
                <div class="flex-shrink-0">
                  <div class="h-2 w-2 rounded-full" :class="getStatusColor(task.status)"></div>
                </div>
                <div class="ml-4 flex-1">
                  <div class="flex items-center justify-between">
                    <div class="flex-1">
                      <h3 class="text-sm font-medium text-gray-900">
                        {{ task.title }}
                      </h3>
                      <p class="text-sm text-gray-500 mt-1">
                        {{ task.description }}
                      </p>
                      <div class="flex items-center mt-2 space-x-4">
                        <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium" :class="getPriorityColor(task.priority)">
                          {{ getPriorityText(task.priority) }}
                        </span>
                        <span v-if="task.estimate" class="text-xs text-gray-500">
                          {{ task.estimate }} 故事点
                        </span>
                        <span v-if="task.assignee" class="text-xs text-gray-500">
                          分配给: {{ getAssigneeName(task.assignee) }}
                        </span>
                        <span v-if="task.due_date" class="text-xs text-gray-500">
                          截止: {{ formatDate(task.due_date) }}
                        </span>
                      </div>
                      <div v-if="task.labels && task.labels.length" class="flex flex-wrap gap-1 mt-2">
                        <span v-for="label in task.labels" :key="label" class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800">
                          {{ label }}
                        </span>
                      </div>
                    </div>
                    <div class="flex items-center space-x-2">
                      <button class="text-indigo-600 hover:text-indigo-900 text-sm font-medium">
                        编辑
                      </button>
                      <button class="text-gray-400 hover:text-gray-600">
                        <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.75a.75.75 0 110-1.5.75.75 0 010 1.5zM12 12.75a.75.75 0 110-1.5.75.75 0 010 1.5zM12 18.75a.75.75 0 110-1.5.75.75 0 010 1.5z" />
                        </svg>
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </li>
        </ul>
      </div>

      <!-- 分页 -->
      <div class="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6 mt-6">
        <div class="flex-1 flex justify-between sm:hidden">
          <a href="#" class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50">
            上一页
          </a>
          <a href="#" class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50">
            下一页
          </a>
        </div>
        <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
          <div>
            <p class="text-sm text-gray-700">
              显示第 <span class="font-medium">1</span> 到 <span class="font-medium">10</span> 项，共 <span class="font-medium">25</span> 项结果
            </p>
          </div>
          <div>
            <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
              <a href="#" class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
                <span class="sr-only">上一页</span>
                <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
              </a>
              <a href="#" class="relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-700 hover:bg-gray-50">
                1
              </a>
              <a href="#" class="relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-700 hover:bg-gray-50">
                2
              </a>
              <a href="#" class="relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-700 hover:bg-gray-50">
                3
              </a>
              <a href="#" class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
                <span class="sr-only">下一页</span>
                <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                </svg>
              </a>
            </nav>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
// 页面元数据
useSeoMeta({
  title: 'PowerX Scrum - 任务管理',
  description: 'Scrum任务管理页面'
})

// 获取运行时配置
const config = useRuntimeConfig()

// 示例任务数据（实际应该从API获取）
const tasks = ref([
  {
    id: 1,
    title: '实现 Sprint 管理功能',
    description: '开发 Sprint 的创建、启动、完成等生命周期管理',
    status: 'in_progress',
    priority: 'high',
    assignee: 102,
    estimate: 21,
    due_date: '2025-01-03',
    labels: ['sprint', 'management', 'backend']
  },
  {
    id: 2,
    title: '设计任务看板界面',
    description: '创建拖拽式的任务看板，支持状态变更',
    status: 'todo',
    priority: 'medium',
    assignee: 101,
    estimate: 13,
    due_date: '2025-01-06',
    labels: ['kanban', 'ui', 'frontend']
  },
  {
    id: 3,
    title: '开发燃尽图功能',
    description: '实现 Sprint 进度的可视化燃尽图',
    status: 'todo',
    priority: 'low',
    assignee: 104,
    estimate: 8,
    labels: ['charts', 'analytics', 'frontend']
  },
  {
    id: 4,
    title: '设计用户注册登录界面',
    description: '创建用户友好的注册和登录页面',
    status: 'done',
    priority: 'high',
    assignee: 101,
    estimate: 8,
    labels: ['ui', 'frontend', 'user-auth']
  },
  {
    id: 5,
    title: '实现用户认证API',
    description: '开发用户注册、登录、注销的后端API',
    status: 'done',
    priority: 'high',
    assignee: 102,
    estimate: 13,
    labels: ['api', 'backend', 'auth']
  }
])

// 工具函数
const getStatusColor = (status) => {
  const colors = {
    'todo': 'bg-gray-400',
    'in_progress': 'bg-blue-400',
    'done': 'bg-green-400'
  }
  return colors[status] || 'bg-gray-400'
}

const getPriorityColor = (priority) => {
  const colors = {
    'urgent': 'bg-red-100 text-red-800',
    'high': 'bg-orange-100 text-orange-800',
    'medium': 'bg-yellow-100 text-yellow-800',
    'low': 'bg-green-100 text-green-800'
  }
  return colors[priority] || 'bg-gray-100 text-gray-800'
}

const getPriorityText = (priority) => {
  const texts = {
    'urgent': '紧急',
    'high': '高',
    'medium': '中',
    'low': '低'
  }
  return texts[priority] || priority
}

const getAssigneeName = (assigneeId) => {
  const assignees = {
    101: 'John Doe',
    102: 'Alice Smith',
    103: 'Mike Brown',
    104: 'Sarah Wilson'
  }
  return assignees[assigneeId] || `用户${assigneeId}`
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('zh-CN')
}

// 这里可以添加API调用逻辑
// const { data: tasksData } = await $fetch(`${config.public.apiBaseUrl}/tasks`)
</script>