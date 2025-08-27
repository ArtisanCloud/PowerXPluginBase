package domain
package domain

import (
	"context"
	"time"

	"github.com/powerx-plugins/scrum/internal/db"
)

// TaskRepository 任务仓储接口
type TaskRepository interface {
	// CRUD 操作
	Create(ctx context.Context, tdb *db.TenantDB, task *Task) error
	GetByID(ctx context.Context, tdb *db.TenantDB, id uint) (*Task, error)
	Update(ctx context.Context, tdb *db.TenantDB, task *Task) error
	Delete(ctx context.Context, tdb *db.TenantDB, id uint) error
	
	// 查询操作
	List(ctx context.Context, tdb *db.TenantDB, opts *TaskListOptions) ([]*Task, int64, error)
	GetBySprintID(ctx context.Context, tdb *db.TenantDB, sprintID uint) ([]*Task, error)
	GetByAssignee(ctx context.Context, tdb *db.TenantDB, assignee int64) ([]*Task, error)
	GetByStatus(ctx context.Context, tdb *db.TenantDB, status TaskStatus) ([]*Task, error)
	Search(ctx context.Context, tdb *db.TenantDB, query string) ([]*Task, error)
	
	// 统计操作
	CountByStatus(ctx context.Context, tdb *db.TenantDB, status TaskStatus) (int64, error)
	CountBySprintID(ctx context.Context, tdb *db.TenantDB, sprintID uint) (int64, error)
	GetOverdueTasks(ctx context.Context, tdb *db.TenantDB) ([]*Task, error)
	
	// 批量操作
	BatchUpdateStatus(ctx context.Context, tdb *db.TenantDB, ids []uint, status TaskStatus) error
	BatchDelete(ctx context.Context, tdb *db.TenantDB, ids []uint) error
}

// SprintRepository Sprint 仓储接口
type SprintRepository interface {
	// CRUD 操作
	Create(ctx context.Context, tdb *db.TenantDB, sprint *Sprint) error
	GetByID(ctx context.Context, tdb *db.TenantDB, id uint) (*Sprint, error)
	Update(ctx context.Context, tdb *db.TenantDB, sprint *Sprint) error
	Delete(ctx context.Context, tdb *db.TenantDB, id uint) error
	
	// 查询操作
	List(ctx context.Context, tdb *db.TenantDB, opts *SprintListOptions) ([]*Sprint, int64, error)
	GetByStatus(ctx context.Context, tdb *db.TenantDB, status SprintStatus) ([]*Sprint, error)
	GetActive(ctx context.Context, tdb *db.TenantDB) (*Sprint, error)
	GetByDateRange(ctx context.Context, tdb *db.TenantDB, start, end time.Time) ([]*Sprint, error)
	
	// 关联操作
	GetWithTasks(ctx context.Context, tdb *db.TenantDB, id uint) (*Sprint, error)
	AddTask(ctx context.Context, tdb *db.TenantDB, sprintID, taskID uint) error
	RemoveTask(ctx context.Context, tdb *db.TenantDB, taskID uint) error
	
	// 统计操作
	CountByStatus(ctx context.Context, tdb *db.TenantDB, status SprintStatus) (int64, error)
	GetSprintStats(ctx context.Context, tdb *db.TenantDB, id uint) (*TaskStats, error)
}

// TaskListOptions 任务列表查询选项
type TaskListOptions struct {
	// 分页
	Page  int
	Limit int
	
	// 过滤
	Status    *TaskStatus
	Priority  *Priority
	Assignee  *int64
	SprintID  *uint
	Labels    []string
	Search    string
	DueBefore *time.Time
	DueAfter  *time.Time
	
	// 排序
	SortBy    string // created_at, updated_at, title, priority, due_date
	SortOrder string // asc, desc
	
	// 包含关联
	IncludeSprint bool
}

// SprintListOptions Sprint 列表查询选项
type SprintListOptions struct {
	// 分页
	Page  int
	Limit int
	
	// 过滤
	Status     *SprintStatus
	Search     string
	StartAfter *time.Time
	EndBefore  *time.Time
	
	// 排序
	SortBy    string // created_at, updated_at, start_date, end_date, name
	SortOrder string // asc, desc
	
	// 包含关联
	IncludeTasks bool
	IncludeStats bool
}

// TaskService 任务服务接口
type TaskService interface {
	// CRUD 操作
	CreateTask(ctx context.Context, tenantID int64, req *CreateTaskRequest) (*Task, error)
	GetTask(ctx context.Context, tenantID int64, id uint) (*Task, error)
	UpdateTask(ctx context.Context, tenantID int64, id uint, req *UpdateTaskRequest) (*Task, error)
	DeleteTask(ctx context.Context, tenantID int64, id uint) error
	
	// 列表和搜索
	ListTasks(ctx context.Context, tenantID int64, opts *TaskListOptions) ([]*Task, int64, error)
	SearchTasks(ctx context.Context, tenantID int64, query string) ([]*Task, error)
	
	// 状态管理
	UpdateTaskStatus(ctx context.Context, tenantID int64, id uint, status TaskStatus) error
	AssignTask(ctx context.Context, tenantID int64, id uint, assignee int64) error
	UnassignTask(ctx context.Context, tenantID int64, id uint) error
	
	// Sprint 管理
	AddTaskToSprint(ctx context.Context, tenantID int64, taskID, sprintID uint) error
	RemoveTaskFromSprint(ctx context.Context, tenantID int64, taskID uint) error
	
	// 标签管理
	AddTaskLabel(ctx context.Context, tenantID int64, id uint, label string) error
	RemoveTaskLabel(ctx context.Context, tenantID int64, id uint, label string) error
	
	// 批量操作
	BatchUpdateTaskStatus(ctx context.Context, tenantID int64, ids []uint, status TaskStatus) error
	BatchDeleteTasks(ctx context.Context, tenantID int64, ids []uint) error
	
	// 统计和报告
	GetTaskStatsByStatus(ctx context.Context, tenantID int64) (map[TaskStatus]int64, error)
	GetOverdueTasks(ctx context.Context, tenantID int64) ([]*Task, error)
	GetTasksByAssignee(ctx context.Context, tenantID int64, assignee int64) ([]*Task, error)
}

// SprintService Sprint 服务接口
type SprintService interface {
	// CRUD 操作
	CreateSprint(ctx context.Context, tenantID int64, req *CreateSprintRequest) (*Sprint, error)
	GetSprint(ctx context.Context, tenantID int64, id uint) (*Sprint, error)
	UpdateSprint(ctx context.Context, tenantID int64, id uint, req *UpdateSprintRequest) (*Sprint, error)
	DeleteSprint(ctx context.Context, tenantID int64, id uint) error
	
	// 列表和搜索
	ListSprints(ctx context.Context, tenantID int64, opts *SprintListOptions) ([]*Sprint, int64, error)
	GetActiveSprint(ctx context.Context, tenantID int64) (*Sprint, error)
	
	// Sprint 生命周期
	StartSprint(ctx context.Context, tenantID int64, id uint) error
	CompleteSprint(ctx context.Context, tenantID int64, id uint) error
	
	// 任务管理
	AddTaskToSprint(ctx context.Context, tenantID int64, sprintID, taskID uint) error
	RemoveTaskFromSprint(ctx context.Context, tenantID int64, taskID uint) error
	GetSprintTasks(ctx context.Context, tenantID int64, sprintID uint) ([]*Task, error)
	
	// 统计和报告
	GetSprintStats(ctx context.Context, tenantID int64, id uint) (*TaskStats, error)
	GetSprintVelocity(ctx context.Context, tenantID int64, id uint) (int, error)
	GetTeamVelocityHistory(ctx context.Context, tenantID int64, limit int) ([]int, error)
	
	// 计划和预测
	GenerateSprintPlan(ctx context.Context, tenantID int64, req *GenerateSprintPlanRequest) (*SprintPlan, error)
}

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	Status      TaskStatus             `json:"status,omitempty"`
	Priority    Priority               `json:"priority,omitempty"`
	Assignee    *int64                 `json:"assignee,omitempty"`
	SprintID    *uint                  `json:"sprint_id,omitempty"`
	Labels      []string               `json:"labels,omitempty"`
	DueDate     *time.Time             `json:"due_date,omitempty"`
	Estimate    *int                   `json:"estimate,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	Title       *string                `json:"title,omitempty"`
	Description *string                `json:"description,omitempty"`
	Status      *TaskStatus            `json:"status,omitempty"`
	Priority    *Priority              `json:"priority,omitempty"`
	Assignee    *int64                 `json:"assignee,omitempty"`
	SprintID    *uint                  `json:"sprint_id,omitempty"`
	Labels      []string               `json:"labels,omitempty"`
	DueDate     *time.Time             `json:"due_date,omitempty"`
	Estimate    *int                   `json:"estimate,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
}

// CreateSprintRequest 创建 Sprint 请求
type CreateSprintRequest struct {
	Name      string    `json:"name"`
	Goal      string    `json:"goal,omitempty"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Capacity  *int      `json:"capacity,omitempty"`
}

// UpdateSprintRequest 更新 Sprint 请求
type UpdateSprintRequest struct {
	Name      *string    `json:"name,omitempty"`
	Goal      *string    `json:"goal,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Capacity  *int       `json:"capacity,omitempty"`
	Status    *SprintStatus `json:"status,omitempty"`
}

// GenerateSprintPlanRequest 生成 Sprint 计划请求
type GenerateSprintPlanRequest struct {
	SprintDuration int     `json:"sprint_duration"` // 天数
	TeamCapacity   int     `json:"team_capacity"`   // 故事点
	PriorityTasks  []uint  `json:"priority_tasks,omitempty"`
}

// SprintPlan Sprint 计划
type SprintPlan struct {
	RecommendedTasks []uint `json:"recommended_tasks"`
	TotalPoints      int    `json:"total_points"`
	EstimatedVelocity int   `json:"estimated_velocity"`
	Confidence       float64 `json:"confidence"` // 0-1
	Reason           string `json:"reason"`
}