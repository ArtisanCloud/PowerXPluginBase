/**
 * Nuxt 组合式函数
 * 用于在组件中方便地使用 API
 */

import { createApiClient, type Task, type Sprint, type TaskListParams, type SprintListParams } from '~/utils/api'

/**
 * 使用任务 API
 */
export function useTasks() {
  const api = createApiClient()

  const tasks = ref<Task[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchTasks = async (params: TaskListParams = {}) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await api.getTasks(params)
      if (response.success && response.data) {
        tasks.value = response.data.data
        return response.data
      } else {
        throw new Error(response.error?.message || '获取任务失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '获取任务失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  const createTask = async (data: any) => {
    loading.value = true
    error.value = null

    try {
      const response = await api.createTask(data)
      if (response.success && response.data) {
        // 刷新任务列表
        await fetchTasks()
        return response.data
      } else {
        throw new Error(response.error?.message || '创建任务失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '创建任务失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateTask = async (id: number, data: any) => {
    loading.value = true
    error.value = null

    try {
      const response = await api.updateTask(id, data)
      if (response.success && response.data) {
        // 更新本地任务列表
        const index = tasks.value.findIndex(task => task.id === id)
        if (index !== -1) {
          tasks.value[index] = response.data
        }
        return response.data
      } else {
        throw new Error(response.error?.message || '更新任务失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '更新任务失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  const deleteTask = async (id: number) => {
    loading.value = true
    error.value = null

    try {
      const response = await api.deleteTask(id)
      if (response.success) {
        // 从本地列表中移除
        tasks.value = tasks.value.filter(task => task.id !== id)
      } else {
        throw new Error(response.error?.message || '删除任务失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '删除任务失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateTaskStatus = async (id: number, status: Task['status']) => {
    try {
      const response = await api.updateTaskStatus(id, status)
      if (response.success) {
        // 更新本地状态
        const task = tasks.value.find(t => t.id === id)
        if (task) {
          task.status = status
        }
      } else {
        throw new Error(response.error?.message || '更新状态失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '更新状态失败'
      throw err
    }
  }

  return {
    tasks: readonly(tasks),
    loading: readonly(loading),
    error: readonly(error),
    fetchTasks,
    createTask,
    updateTask,
    deleteTask,
    updateTaskStatus
  }
}

/**
 * 使用 Sprint API
 */
export function useSprints() {
  const api = createApiClient()

  const sprints = ref<Sprint[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchSprints = async (params: SprintListParams = {}) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await api.getSprints(params)
      if (response.success && response.data) {
        sprints.value = response.data.data
        return response.data
      } else {
        throw new Error(response.error?.message || '获取 Sprint 失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '获取 Sprint 失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  const createSprint = async (data: any) => {
    loading.value = true
    error.value = null

    try {
      const response = await api.createSprint(data)
      if (response.success && response.data) {
        // 刷新 Sprint 列表
        await fetchSprints()
        return response.data
      } else {
        throw new Error(response.error?.message || '创建 Sprint 失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '创建 Sprint 失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateSprint = async (id: number, data: any) => {
    loading.value = true
    error.value = null

    try {
      const response = await api.updateSprint(id, data)
      if (response.success && response.data) {
        // 更新本地 Sprint 列表
        const index = sprints.value.findIndex(sprint => sprint.id === id)
        if (index !== -1) {
          sprints.value[index] = response.data
        }
        return response.data
      } else {
        throw new Error(response.error?.message || '更新 Sprint 失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '更新 Sprint 失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  const deleteSprint = async (id: number) => {
    loading.value = true
    error.value = null

    try {
      const response = await api.deleteSprint(id)
      if (response.success) {
        // 从本地列表中移除
        sprints.value = sprints.value.filter(sprint => sprint.id !== id)
      } else {
        throw new Error(response.error?.message || '删除 Sprint 失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '删除 Sprint 失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  const startSprint = async (id: number) => {
    try {
      const response = await api.startSprint(id)
      if (response.success) {
        // 更新本地状态
        const sprint = sprints.value.find(s => s.id === id)
        if (sprint) {
          sprint.status = 'active'
        }
      } else {
        throw new Error(response.error?.message || '启动 Sprint 失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '启动 Sprint 失败'
      throw err
    }
  }

  const completeSprint = async (id: number) => {
    try {
      const response = await api.completeSprint(id)
      if (response.success) {
        // 更新本地状态
        const sprint = sprints.value.find(s => s.id === id)
        if (sprint) {
          sprint.status = 'completed'
        }
      } else {
        throw new Error(response.error?.message || '完成 Sprint 失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '完成 Sprint 失败'
      throw err
    }
  }

  return {
    sprints: readonly(sprints),
    loading: readonly(loading),
    error: readonly(error),
    fetchSprints,
    createSprint,
    updateSprint,
    deleteSprint,
    startSprint,
    completeSprint
  }
}

/**
 * 使用仪表板数据
 */
export function useDashboard() {
  const api = createApiClient()
  
  const stats = ref({
    todoTasks: 0,
    inProgressTasks: 0,
    doneTasks: 0,
    currentSprint: null as Sprint | null
  })
  
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchDashboardData = async () => {
    loading.value = true
    error.value = null

    try {
      // 并行获取任务统计和当前 Sprint
      const [tasksResponse, sprintsResponse] = await Promise.all([
        api.getTasks({ limit: 100 }), // 获取所有任务用于统计
        api.getSprints({ status: 'active', limit: 1 }) // 获取当前活跃的 Sprint
      ])

      if (tasksResponse.success && tasksResponse.data) {
        const tasks = tasksResponse.data.data
        stats.value.todoTasks = tasks.filter(t => t.status === 'todo').length
        stats.value.inProgressTasks = tasks.filter(t => t.status === 'in_progress').length
        stats.value.doneTasks = tasks.filter(t => t.status === 'done').length
      }

      if (sprintsResponse.success && sprintsResponse.data && sprintsResponse.data.data.length > 0) {
        stats.value.currentSprint = sprintsResponse.data.data[0]
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '获取仪表板数据失败'
      console.error('Dashboard data fetch error:', err)
    } finally {
      loading.value = false
    }
  }

  return {
    stats: readonly(stats),
    loading: readonly(loading),
    error: readonly(error),
    fetchDashboardData
  }
}