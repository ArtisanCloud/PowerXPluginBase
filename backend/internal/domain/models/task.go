package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Task 任务模型
type Task struct {
	BaseModel

	// 基本信息
	Title       string     `gorm:"column:title;type:varchar(200);not null" json:"title"`
	Description string     `gorm:"column:description;type:text" json:"description,omitempty"`
	TaskType    TaskType   `gorm:"column:task_type;type:varchar(20);not null;default:'user_story'" json:"task_type"`
	Status      TaskStatus `gorm:"column:status;type:varchar(20);not null;default:'todo';index" json:"status"`
	Priority    Priority   `gorm:"column:priority;type:varchar(20);not null;default:'medium';index" json:"priority"`

	// 分配和关联
	AssigneeID *uint64 `gorm:"column:assignee_id;index" json:"assignee_id,omitempty"`
	ReporterID *uint64 `gorm:"column:reporter_id;index" json:"reporter_id,omitempty"`
	SprintID   *uint64 `gorm:"column:sprint_id;index" json:"sprint_id,omitempty"`
	ProjectID  uint64  `gorm:"column:project_id;not null;index" json:"project_id"`
	ParentID   *uint64 `gorm:"column:parent_id;index" json:"parent_id,omitempty"` // 用于子任务

	// 时间管理
	DueDate     *time.Time `gorm:"column:due_date" json:"due_date,omitempty"`
	StartDate   *time.Time `gorm:"column:start_date" json:"start_date,omitempty"`
	CompletedAt *time.Time `gorm:"column:completed_at" json:"completed_at,omitempty"`

	// 估算和跟踪
	Estimate       *int `gorm:"column:estimate" json:"estimate,omitempty"`
	StoryPoints    *int `gorm:"column:story_points" json:"story_points,omitempty"`
	OriginalHours  *int `gorm:"column:original_hours" json:"original_hours,omitempty"`
	RemainingHours *int `gorm:"column:remaining_hours" json:"remaining_hours,omitempty"`
	LoggedHours    *int `gorm:"column:logged_hours;default:0" json:"logged_hours,omitempty"`

	// 验收标准和业务价值
	AcceptanceCriteria string  `gorm:"column:acceptance_criteria;type:text" json:"acceptance_criteria,omitempty"`
	BusinessValue      *int    `gorm:"column:business_value" json:"business_value,omitempty"`
	EpicColor          *string `gorm:"column:epic_color;type:varchar(7)" json:"epic_color,omitempty"` // Epic的颜色标识`

	// 标签和元数据
	Labels Labels         `gorm:"column:labels;type:jsonb" json:"labels,omitempty"`
	Meta   datatypes.JSON `gorm:"column:meta;type:jsonb" json:"meta,omitempty"`

	// 关联模型
	Sprint   *Sprint  `json:"sprint,omitempty" gorm:"foreignKey:SprintID"`
	Project  *Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
	Parent   *Task    `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Task   `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

// TaskType 任务类型枚举
type TaskType string

const (
	TaskTypeUserStory TaskType = "user_story"
	TaskTypeTask      TaskType = "task"
	TaskTypeBug       TaskType = "bug"
	TaskTypeEpic      TaskType = "epic"
	TaskTypeSubtask   TaskType = "subtask"
)

// TaskStatus 任务状态枚举
type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusReview     TaskStatus = "review"
	TaskStatusTesting    TaskStatus = "testing"
	TaskStatusBlocked    TaskStatus = "blocked"
	TaskStatusDone       TaskStatus = "done"
)

// Priority 优先级枚举
type Priority string

const (
	PriorityLowest   Priority = "lowest"
	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityUrgent   Priority = "urgent"
	PriorityCritical Priority = "critical"
)

// Labels 标签类型（JSON 数组）
type Labels []string

// TableName 返回表名
func (t *Task) TableName() string {
	return PowerXSchema + "." + TableTask
}

// GetTableName 获取表名
func (t *Task) GetTableName(needFull bool) string {
	if needFull {
		return t.TableName()
	}
	return TableTask
}

// BeforeCreate GORM 钩子
func (t *Task) BeforeCreate(tx *gorm.DB) error {
	// 设置默认值
	if t.TaskType == "" {
		t.TaskType = TaskTypeUserStory
	}
	if t.Status == "" {
		t.Status = TaskStatusTodo
	}
	if t.Priority == "" {
		t.Priority = PriorityMedium
	}
	if t.Labels == nil {
		t.Labels = Labels{}
	}
	if t.LoggedHours == nil {
		zero := 0
		t.LoggedHours = &zero
	}

	return nil
}

// Value 实现 driver.Valuer 接口，用于数据库存储
func (l Labels) Value() (driver.Value, error) {
	if l == nil {
		return "[]", nil
	}
	return json.Marshal(l)
}

// Scan 实现 sql.Scanner 接口，用于数据库读取
func (l *Labels) Scan(value interface{}) error {
	if value == nil {
		*l = Labels{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into Labels", value)
	}

	return json.Unmarshal(bytes, l)
}

// 业务方法

// IsCompleted 检查任务是否已完成
func (t *Task) IsCompleted() bool {
	return t.Status == TaskStatusDone
}

// IsInProgress 检查任务是否正在进行中
func (t *Task) IsInProgress() bool {
	return t.Status == TaskStatusInProgress
}

// IsBlocked 检查任务是否被阻塞
func (t *Task) IsBlocked() bool {
	return t.Status == TaskStatusBlocked
}

// IsOverdue 检查任务是否过期
func (t *Task) IsOverdue() bool {
	if t.DueDate == nil {
		return false
	}
	return time.Now().After(*t.DueDate) && !t.IsCompleted()
}

// HasLabel 检查是否包含指定标签
func (t *Task) HasLabel(label string) bool {
	for _, l := range t.Labels {
		if l == label {
			return true
		}
	}
	return false
}

// AddLabel 添加标签
func (t *Task) AddLabel(label string) {
	if !t.HasLabel(label) {
		t.Labels = append(t.Labels, label)
	}
}

// RemoveLabel 移除标签
func (t *Task) RemoveLabel(label string) {
	for i, l := range t.Labels {
		if l == label {
			t.Labels = append(t.Labels[:i], t.Labels[i+1:]...)
			break
		}
	}
}

// IsEpic 检查是否为Epic
func (t *Task) IsEpic() bool {
	return t.TaskType == TaskTypeEpic
}

// IsUserStory 检查是否为用户故事
func (t *Task) IsUserStory() bool {
	return t.TaskType == TaskTypeUserStory
}

// GetChildrenTasks 获取子任务
func (t *Task) GetChildrenTasks() []Task {
	return t.Children
}

// GetEpic 获取所属的Epic（如果是用户故事或任务）
func (t *Task) GetEpic() *Task {
	if t.Parent != nil && t.Parent.IsEpic() {
		return t.Parent
	}
	return nil
}

// GetCompletionPercentage 计算完成百分比（基于子任务）
func (t *Task) GetCompletionPercentage() float64 {
	if len(t.Children) == 0 {
		if t.IsCompleted() {
			return 100.0
		}
		return 0.0
	}

	completedCount := 0
	for _, child := range t.Children {
		if child.IsCompleted() {
			completedCount++
		}
	}

	return float64(completedCount) / float64(len(t.Children)) * 100.0
}

// GetTotalStoryPoints 计算总故事点数（包括子任务）
func (t *Task) GetTotalStoryPoints() int {
	total := 0
	if t.StoryPoints != nil {
		total += *t.StoryPoints
	}

	for _, child := range t.Children {
		total += child.GetTotalStoryPoints()
	}

	return total
}

// CanTransitionTo 检查是否可以转换到指定状态
func (t *Task) CanTransitionTo(status TaskStatus) bool {
	switch t.Status {
	case TaskStatusTodo:
		return status == TaskStatusInProgress || status == TaskStatusBlocked
	case TaskStatusInProgress:
		return status == TaskStatusTodo || status == TaskStatusReview ||
			status == TaskStatusTesting || status == TaskStatusBlocked || status == TaskStatusDone
	case TaskStatusReview:
		return status == TaskStatusInProgress || status == TaskStatusTesting ||
			status == TaskStatusBlocked || status == TaskStatusDone
	case TaskStatusTesting:
		return status == TaskStatusReview || status == TaskStatusInProgress ||
			status == TaskStatusBlocked || status == TaskStatusDone
	case TaskStatusBlocked:
		return status == TaskStatusTodo || status == TaskStatusInProgress
	case TaskStatusDone:
		return status == TaskStatusInProgress || status == TaskStatusReview
	default:
		return false
	}
}

// Validate 验证任务数据
func (t *Task) Validate() error {
	if t.Title == "" {
		return fmt.Errorf("title is required")
	}

	if len(t.Title) > 200 {
		return fmt.Errorf("title too long (max 200 characters)")
	}

	if !t.Status.IsValid() {
		return fmt.Errorf("invalid status: %s", t.Status)
	}

	if !t.Priority.IsValid() {
		return fmt.Errorf("invalid priority: %s", t.Priority)
	}

	if !t.TaskType.IsValid() {
		return fmt.Errorf("invalid task type: %s", t.TaskType)
	}

	if t.StoryPoints != nil && *t.StoryPoints < 0 {
		return fmt.Errorf("story points cannot be negative")
	}

	if t.BusinessValue != nil && *t.BusinessValue < 0 {
		return fmt.Errorf("business value cannot be negative")
	}

	if t.EpicColor != nil && len(*t.EpicColor) > 0 {
		// 验证颜色值格式 (#RRGGBB)
		if len(*t.EpicColor) != 7 || (*t.EpicColor)[0] != '#' {
			return fmt.Errorf("epic color must be in format #RRGGBB")
		}
	}

	if t.DueDate != nil && t.StartDate != nil && t.DueDate.Before(*t.StartDate) {
		return fmt.Errorf("due date cannot be before start date")
	}

	return nil
}

// 枚举验证方法

// IsValid 验证任务类型是否有效
func (tt TaskType) IsValid() bool {
	switch tt {
	case TaskTypeUserStory, TaskTypeTask, TaskTypeBug, TaskTypeEpic, TaskTypeSubtask:
		return true
	default:
		return false
	}
}

// IsValid 验证状态是否有效
func (ts TaskStatus) IsValid() bool {
	switch ts {
	case TaskStatusTodo, TaskStatusInProgress, TaskStatusReview,
		TaskStatusTesting, TaskStatusBlocked, TaskStatusDone:
		return true
	default:
		return false
	}
}

// IsValid 验证优先级是否有效
func (p Priority) IsValid() bool {
	switch p {
	case PriorityLowest, PriorityLow, PriorityMedium,
		PriorityHigh, PriorityUrgent, PriorityCritical:
		return true
	default:
		return false
	}
}

// String 实现 Stringer 接口
func (tt TaskType) String() string {
	return string(tt)
}

func (ts TaskStatus) String() string {
	return string(ts)
}

func (p Priority) String() string {
	return string(p)
}
