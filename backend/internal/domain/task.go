package domain
package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/powerx-plugins/scrum/internal/db"
	"gorm.io/gorm"
)

// Task 任务模型
type Task struct {
	db.BaseModel
	Title       string     `json:"title" gorm:"not null;size:200"`
	Description string     `json:"description" gorm:"type:text"`
	Status      TaskStatus `json:"status" gorm:"not null;default:'todo'"`
	Priority    Priority   `json:"priority" gorm:"not null;default:'medium'"`
	Assignee    *int64     `json:"assignee,omitempty" gorm:"index"`
	SprintID    *uint      `json:"sprint_id,omitempty" gorm:"index"`
	Labels      Labels     `json:"labels" gorm:"type:jsonb"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Estimate    *int       `json:"estimate,omitempty"` // 故事点估算
	Meta        Meta       `json:"meta" gorm:"type:jsonb"`

	// 关联
	Sprint *Sprint `json:"sprint,omitempty" gorm:"foreignKey:SprintID"`
}

// TaskStatus 任务状态枚举
type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

// Priority 优先级枚举
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)

// Labels 标签类型（JSON 数组）
type Labels []string

// Meta 元数据类型（JSON 对象）
type Meta map[string]interface{}

// TableName 返回表名
func (Task) TableName() string {
	return "task"
}

// BeforeCreate GORM 钩子
func (t *Task) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := t.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 设置默认值
	if t.Status == "" {
		t.Status = TaskStatusTodo
	}
	if t.Priority == "" {
		t.Priority = PriorityMedium
	}
	if t.Labels == nil {
		t.Labels = Labels{}
	}
	if t.Meta == nil {
		t.Meta = Meta{}
	}

	return nil
}

// IsValidStatus 验证状态是否有效
func (ts TaskStatus) IsValid() bool {
	switch ts {
	case TaskStatusTodo, TaskStatusInProgress, TaskStatusDone:
		return true
	default:
		return false
	}
}

// String 实现 Stringer 接口
func (ts TaskStatus) String() string {
	return string(ts)
}

// IsValidPriority 验证优先级是否有效
func (p Priority) IsValid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh, PriorityUrgent:
		return true
	default:
		return false
	}
}

// String 实现 Stringer 接口
func (p Priority) String() string {
	return string(p)
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

// Value 实现 driver.Valuer 接口，用于数据库存储
func (m Meta) Value() (driver.Value, error) {
	if m == nil {
		return "{}", nil
	}
	return json.Marshal(m)
}

// Scan 实现 sql.Scanner 接口，用于数据库读取
func (m *Meta) Scan(value interface{}) error {
	if value == nil {
		*m = Meta{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into Meta", value)
	}

	return json.Unmarshal(bytes, m)
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

// SetMeta 设置元数据
func (t *Task) SetMeta(key string, value interface{}) {
	if t.Meta == nil {
		t.Meta = Meta{}
	}
	t.Meta[key] = value
}

// GetMeta 获取元数据
func (t *Task) GetMeta(key string) (interface{}, bool) {
	if t.Meta == nil {
		return nil, false
	}
	value, exists := t.Meta[key]
	return value, exists
}

// IsCompleted 检查任务是否已完成
func (t *Task) IsCompleted() bool {
	return t.Status == TaskStatusDone
}

// IsInProgress 检查任务是否正在进行中
func (t *Task) IsInProgress() bool {
	return t.Status == TaskStatusInProgress
}

// IsTodo 检查任务是否待办
func (t *Task) IsTodo() bool {
	return t.Status == TaskStatusTodo
}

// IsOverdue 检查任务是否过期
func (t *Task) IsOverdue() bool {
	if t.DueDate == nil {
		return false
	}
	return time.Now().After(*t.DueDate) && !t.IsCompleted()
}

// DaysUntilDue 计算距离截止日期的天数
func (t *Task) DaysUntilDue() *int {
	if t.DueDate == nil {
		return nil
	}
	days := int(time.Until(*t.DueDate).Hours() / 24)
	return &days
}

// CanTransitionTo 检查是否可以转换到指定状态
func (t *Task) CanTransitionTo(status TaskStatus) bool {
	switch t.Status {
	case TaskStatusTodo:
		return status == TaskStatusInProgress || status == TaskStatusDone
	case TaskStatusInProgress:
		return status == TaskStatusTodo || status == TaskStatusDone
	case TaskStatusDone:
		return status == TaskStatusTodo || status == TaskStatusInProgress
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
	
	if t.Estimate != nil && *t.Estimate < 0 {
		return fmt.Errorf("estimate cannot be negative")
	}
	
	return nil
}