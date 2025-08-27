package contracts

import "time"

// APIResponse 标准 API 响应格式
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     *APIError   `json:"error,omitempty"`
	Message   string      `json:"message,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// APIError API 错误信息
type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// PaginationRequest 分页请求
type PaginationRequest struct {
	Page  int `json:"page" form:"page" binding:"min=1"`
	Limit int `json:"limit" form:"limit" binding:"min=1,max=100"`
}

// PaginationResponse 分页响应
type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// ListResponse 列表响应
type ListResponse struct {
	Data       interface{}         `json:"data"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string            `json:"status"`
	Service   string            `json:"service"`
	Version   string            `json:"version"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    map[string]string `json:"checks,omitempty"`
}

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	Title       string                 `json:"title" binding:"required,max=200"`
	Description string                 `json:"description,omitempty" binding:"max=2000"`
	Status      string                 `json:"status,omitempty" binding:"oneof=todo in_progress done"`
	Priority    string                 `json:"priority,omitempty" binding:"oneof=low medium high urgent"`
	Assignee    *int64                 `json:"assignee,omitempty"`
	SprintID    *int64                 `json:"sprint_id,omitempty"`
	Labels      []string               `json:"labels,omitempty"`
	DueDate     *time.Time             `json:"due_date,omitempty"`
	Estimate    *int                   `json:"estimate,omitempty"` // 故事点
	Meta        map[string]interface{} `json:"meta,omitempty"`
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	Title       *string                `json:"title,omitempty" binding:"omitempty,max=200"`
	Description *string                `json:"description,omitempty" binding:"omitempty,max=2000"`
	Status      *string                `json:"status,omitempty" binding:"omitempty,oneof=todo in_progress done"`
	Priority    *string                `json:"priority,omitempty" binding:"omitempty,oneof=low medium high urgent"`
	Assignee    *int64                 `json:"assignee,omitempty"`
	SprintID    *int64                 `json:"sprint_id,omitempty"`
	Labels      []string               `json:"labels,omitempty"`
	DueDate     *time.Time             `json:"due_date,omitempty"`
	Estimate    *int                   `json:"estimate,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
}

// TaskResponse 任务响应
type TaskResponse struct {
	ID          uint                   `json:"id"`
	TenantID    int64                  `json:"tenant_id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	Status      string                 `json:"status"`
	Priority    string                 `json:"priority"`
	Assignee    *int64                 `json:"assignee,omitempty"`
	SprintID    *int64                 `json:"sprint_id,omitempty"`
	Labels      []string               `json:"labels,omitempty"`
	DueDate     *time.Time             `json:"due_date,omitempty"`
	Estimate    *int                   `json:"estimate,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// TaskListRequest 任务列表请求
type TaskListRequest struct {
	PaginationRequest
	Status    string   `json:"status" form:"status" binding:"omitempty,oneof=todo in_progress done"`
	Priority  string   `json:"priority" form:"priority" binding:"omitempty,oneof=low medium high urgent"`
	Assignee  *int64   `json:"assignee" form:"assignee"`
	SprintID  *int64   `json:"sprint_id" form:"sprint_id"`
	Labels    []string `json:"labels" form:"labels"`
	Search    string   `json:"search" form:"search" binding:"max=100"`
	SortBy    string   `json:"sort_by" form:"sort_by" binding:"omitempty,oneof=created_at updated_at title priority due_date"`
	SortOrder string   `json:"sort_order" form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// CreateSprintRequest 创建 Sprint 请求
type CreateSprintRequest struct {
	Name      string    `json:"name" binding:"required,max=100"`
	Goal      string    `json:"goal,omitempty" binding:"max=500"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
	Capacity  *int      `json:"capacity,omitempty"` // 团队容量（故事点）
	Status    string    `json:"status,omitempty" binding:"omitempty,oneof=planning active completed"`
}

// UpdateSprintRequest 更新 Sprint 请求
type UpdateSprintRequest struct {
	Name      *string    `json:"name,omitempty" binding:"omitempty,max=100"`
	Goal      *string    `json:"goal,omitempty" binding:"omitempty,max=500"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Capacity  *int       `json:"capacity,omitempty"`
	Status    *string    `json:"status,omitempty" binding:"omitempty,oneof=planning active completed"`
}

// SprintResponse Sprint 响应
type SprintResponse struct {
	ID        uint      `json:"id"`
	TenantID  int64     `json:"tenant_id"`
	Name      string    `json:"name"`
	Goal      string    `json:"goal,omitempty"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Capacity  *int      `json:"capacity,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 统计信息
	Stats *SprintStats `json:"stats,omitempty"`
}

// SprintStats Sprint 统计信息
type SprintStats struct {
	TotalTasks      int     `json:"total_tasks"`
	CompletedTasks  int     `json:"completed_tasks"`
	TotalPoints     int     `json:"total_points"`
	CompletedPoints int     `json:"completed_points"`
	Progress        float64 `json:"progress"` // 百分比
}

// SprintListRequest Sprint 列表请求
type SprintListRequest struct {
	PaginationRequest
	Status    string `json:"status" form:"status" binding:"omitempty,oneof=planning active completed"`
	Search    string `json:"search" form:"search" binding:"max=100"`
	SortBy    string `json:"sort_by" form:"sort_by" binding:"omitempty,oneof=created_at updated_at start_date end_date name"`
	SortOrder string `json:"sort_order" form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// AgentToolRequest Agent 工具调用请求
type AgentToolRequest struct {
	ToolID  string                 `json:"tool_id" binding:"required"`
	Input   map[string]interface{} `json:"input,omitempty"`
	Context map[string]interface{} `json:"context,omitempty"`
}

// AgentToolResponse Agent 工具调用响应
type AgentToolResponse struct {
	Success bool                   `json:"success"`
	Output  map[string]interface{} `json:"output,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// WorkflowExecuteRequest 工作流执行请求
type WorkflowExecuteRequest struct {
	WorkflowID string                 `json:"workflow_id" binding:"required"`
	Input      map[string]interface{} `json:"input,omitempty"`
	Context    map[string]interface{} `json:"context,omitempty"`
}

// WorkflowExecuteResponse 工作流执行响应
type WorkflowExecuteResponse struct {
	ExecutionID string                 `json:"execution_id"`
	Status      string                 `json:"status"` // running, completed, failed
	Result      map[string]interface{} `json:"result,omitempty"`
	Error       string                 `json:"error,omitempty"`
	Steps       []WorkflowStepResult   `json:"steps,omitempty"`
}

// WorkflowStepResult 工作流步骤结果
type WorkflowStepResult struct {
	StepID    string                 `json:"step_id"`
	Status    string                 `json:"status"` // pending, running, completed, failed, skipped
	Input     map[string]interface{} `json:"input,omitempty"`
	Output    map[string]interface{} `json:"output,omitempty"`
	Error     string                 `json:"error,omitempty"`
	StartTime *time.Time             `json:"start_time,omitempty"`
	EndTime   *time.Time             `json:"end_time,omitempty"`
	Duration  *time.Duration         `json:"duration,omitempty"`
}

// ErrorCode 错误码常量
const (
	// 通用错误码
	ErrCodeInternalError    = "INTERNAL_ERROR"
	ErrCodeInvalidRequest   = "INVALID_REQUEST"
	ErrCodeUnauthorized     = "UNAUTHORIZED"
	ErrCodeForbidden        = "FORBIDDEN"
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodeConflict         = "CONFLICT"
	ErrCodeValidationFailed = "VALIDATION_FAILED"

	// 业务错误码
	ErrCodeTaskNotFound      = "TASK_NOT_FOUND"
	ErrCodeSprintNotFound    = "SPRINT_NOT_FOUND"
	ErrCodeInvalidTaskStatus = "INVALID_TASK_STATUS"
	ErrCodeSprintClosed      = "SPRINT_CLOSED"
	ErrCodePermissionDenied  = "PERMISSION_DENIED"

	// 租户相关错误码
	ErrCodeTenantNotFound = "TENANT_NOT_FOUND"
	ErrCodeTenantMismatch = "TENANT_MISMATCH"

	// Agent 相关错误码
	ErrCodeAgentToolNotFound = "AGENT_TOOL_NOT_FOUND"
	ErrCodeWorkflowNotFound  = "WORKFLOW_NOT_FOUND"
	ErrCodeWorkflowFailed    = "WORKFLOW_FAILED"
)
