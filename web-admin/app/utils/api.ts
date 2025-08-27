/**
 * API 工具类
 * 用于与 PowerX Scrum Plugin 后端 API 通信
 */

interface ApiResponse<T = any> {
  success: boolean
  data?: T
  error?: {
    code: string
    message: string
    details?: any
  }
  message?: string
  timestamp: string
}

interface PaginationResponse<T = any> {
  data: T[]
  pagination: {
    page: number
    limit: number
    total: number
    total_pages: number
  }
}

interface Task {
  id: number
  tenant_id: number
  title: string
  description?: string
  status: 'todo' | 'in_progress' | 'done'
  priority: 'low' | 'medium' | 'high' | 'urgent'
  assignee?: number
  sprint_id?: number
  labels?: string[]
  due_date?: string
  estimate?: number
  meta?: Record<string, any>
  created_at: string
  updated_at: string
}

interface Sprint {
  id: number
  tenant_id: number
  name: string
  goal?: string
  start_date: string
  end_date: string
  capacity?: number
  status: 'planning' | 'active' | 'completed'
  created_at: string
  updated_at: string
  stats?: {
    total_tasks: number
    completed_tasks: number
    total_points: number
    completed_points: number
    progress: number
  }
}

interface CreateTaskRequest {
  title: string
  description?: string
  status?: Task['status']
  priority?: Task['priority']
  assignee?: number
  sprint_id?: number
  labels?: string[]
  due_date?: string
  estimate?: number
  meta?: Record<string, any>
}

interface UpdateTaskRequest extends Partial<CreateTaskRequest> {}

interface CreateSprintRequest {
  name: string
  goal?: string
  start_date: string
  end_date: string
  capacity?: number
}

interface UpdateSprintRequest extends Partial<CreateSprintRequest> {
  status?: Sprint['status']
}

interface TaskListParams {
  page?: number
  limit?: number
  status?: Task['status']
  priority?: Task['priority']
  assignee?: number
  sprint_id?: number
  labels?: string[]
  search?: string
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

interface SprintListParams {
  page?: number
  limit?: number
  status?: Sprint['status']
  search?: string
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

export class ApiClient {
  private baseUrl: string

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl
  }

  private async request<T = any>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    const url = `${this.baseUrl}${endpoint}`
    
    const defaultHeaders = {
      'Content-Type': 'application/json',
      // 在生产环境中，这里需要添加 PowerX 的认证头
      // 'X-PowerX-CTX': getAuthContext(),
      // 'X-PowerX-CTX-JWT': getJWTToken(),
    }

    const config: RequestInit = {
      ...options,
      headers: {
        ...defaultHeaders,
        ...options.headers,
      },
    }

    try {
      const response = await fetch(url, config)
      const data = await response.json()

      if (!response.ok) {
        throw new Error(data.error?.message || `HTTP ${response.status}`)
      }

      return data
    } catch (error) {
      console.error('API Request failed:', error)
      throw error
    }
  }

  // ========== 任务相关 API ==========

  /**
   * 获取任务列表
   */
  async getTasks(params: TaskListParams = {}): Promise<ApiResponse<PaginationResponse<Task>>> {
    const searchParams = new URLSearchParams()
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined) {
        if (Array.isArray(value)) {
          value.forEach(v => searchParams.append(key, v.toString()))
        } else {
          searchParams.append(key, value.toString())
        }
      }
    })

    return this.request(`/tasks?${searchParams.toString()}`)
  }

  /**
   * 获取任务详情
   */
  async getTask(id: number): Promise<ApiResponse<Task>> {
    return this.request(`/tasks/${id}`)
  }

  /**
   * 创建任务
   */
  async createTask(data: CreateTaskRequest): Promise<ApiResponse<Task>> {
    return this.request('/tasks', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  /**
   * 更新任务
   */
  async updateTask(id: number, data: UpdateTaskRequest): Promise<ApiResponse<Task>> {
    return this.request(`/tasks/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  /**
   * 删除任务
   */
  async deleteTask(id: number): Promise<ApiResponse<void>> {
    return this.request(`/tasks/${id}`, {
      method: 'DELETE',
    })
  }

  /**
   * 更新任务状态
   */
  async updateTaskStatus(id: number, status: Task['status']): Promise<ApiResponse<void>> {
    return this.request(`/tasks/${id}/status`, {
      method: 'PATCH',
      body: JSON.stringify({ status }),
    })
  }

  // ========== Sprint 相关 API ==========

  /**
   * 获取 Sprint 列表
   */
  async getSprints(params: SprintListParams = {}): Promise<ApiResponse<PaginationResponse<Sprint>>> {
    const searchParams = new URLSearchParams()
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined) {
        searchParams.append(key, value.toString())
      }
    })

    return this.request(`/sprints?${searchParams.toString()}`)
  }

  /**
   * 获取 Sprint 详情
   */
  async getSprint(id: number): Promise<ApiResponse<Sprint>> {
    return this.request(`/sprints/${id}`)
  }

  /**
   * 创建 Sprint
   */
  async createSprint(data: CreateSprintRequest): Promise<ApiResponse<Sprint>> {
    return this.request('/sprints', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  /**
   * 更新 Sprint
   */
  async updateSprint(id: number, data: UpdateSprintRequest): Promise<ApiResponse<Sprint>> {
    return this.request(`/sprints/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  /**
   * 删除 Sprint
   */
  async deleteSprint(id: number): Promise<ApiResponse<void>> {
    return this.request(`/sprints/${id}`, {
      method: 'DELETE',
    })
  }

  /**
   * 启动 Sprint
   */
  async startSprint(id: number): Promise<ApiResponse<void>> {
    return this.request(`/sprints/${id}/start`, {
      method: 'POST',
    })
  }

  /**
   * 完成 Sprint
   */
  async completeSprint(id: number): Promise<ApiResponse<void>> {
    return this.request(`/sprints/${id}/complete`, {
      method: 'POST',
    })
  }

  /**
   * 获取 Sprint 任务
   */
  async getSprintTasks(id: number): Promise<ApiResponse<Task[]>> {
    return this.request(`/sprints/${id}/tasks`)
  }

  /**
   * 获取 Sprint 统计
   */
  async getSprintStats(id: number): Promise<ApiResponse<Sprint['stats']>> {
    return this.request(`/sprints/${id}/stats`)
  }

  // ========== 管理相关 API ==========

  /**
   * 获取插件清单
   */
  async getManifest(): Promise<ApiResponse<any>> {
    return this.request('/admin/manifest')
  }

  /**
   * 获取 RBAC 信息
   */
  async getRBACInfo(): Promise<ApiResponse<any>> {
    return this.request('/admin/rbac')
  }

  /**
   * 健康检查
   */
  async healthCheck(): Promise<ApiResponse<any>> {
    return this.request('/healthz')
  }
}

// 创建 API 客户端实例
export function createApiClient(): ApiClient {
  const config = useRuntimeConfig()
  return new ApiClient(config.public.apiBaseUrl as string)
}

// 导出类型
export type {
  ApiResponse,
  PaginationResponse,
  Task,
  Sprint,
  CreateTaskRequest,
  UpdateTaskRequest,
  CreateSprintRequest,
  UpdateSprintRequest,
  TaskListParams,
  SprintListParams,
}