package models

import (
	"fmt"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Sprint 冲刺模型
type Sprint struct {
	BaseModel

	// 基本信息
	Name        string       `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Goal        string       `gorm:"column:goal;type:text" json:"goal,omitempty"`
	Description string       `gorm:"column:description;type:text" json:"description,omitempty"`
	Status      SprintStatus `gorm:"column:status;type:varchar(20);not null;default:'planning';index" json:"status"`
	ProjectID   uint         `gorm:"column:project_id;not null;index" json:"project_id"`

	// 时间管理
	StartDate   time.Time  `gorm:"column:start_date;not null" json:"start_date"`
	EndDate     time.Time  `gorm:"column:end_date;not null" json:"end_date"`
	ActualStart *time.Time `gorm:"column:actual_start" json:"actual_start,omitempty"`
	ActualEnd   *time.Time `gorm:"column:actual_end" json:"actual_end,omitempty"`

	// 容量和目标
	Capacity          *int `gorm:"column:capacity" json:"capacity,omitempty"`                     // 团队容量（故事点）
	CommittedPoints   *int `gorm:"column:committed_points" json:"committed_points,omitempty"`     // 承诺的故事点
	CompletedPoints   *int `gorm:"column:completed_points;default:0" json:"completed_points"`     // 完成的故事点
	EstimatedVelocity *int `gorm:"column:estimated_velocity" json:"estimated_velocity,omitempty"` // 预估速度
	ActualVelocity    *int `gorm:"column:actual_velocity" json:"actual_velocity,omitempty"`       // 实际速度

	// 元数据
	Meta datatypes.JSON `gorm:"column:meta;type:jsonb" json:"meta,omitempty"`

	// 关联模型
	Project *Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
	Tasks   []Task   `json:"tasks,omitempty" gorm:"foreignKey:SprintID"`
}

// SprintStatus 冲刺状态枚举
type SprintStatus string

const (
	SprintStatusPlanning  SprintStatus = "planning"  // 规划中
	SprintStatusActive    SprintStatus = "active"    // 进行中
	SprintStatusCompleted SprintStatus = "completed" // 已完成
	SprintStatusCancelled SprintStatus = "cancelled" // 已取消
)

// TableName 返回表名
func (s *Sprint) TableName() string {
	return PowerXSchema + "." + TableSprint
}

// GetTableName 获取表名
func (s *Sprint) GetTableName(needFull bool) string {
	if needFull {
		return s.TableName()
	}
	return TableSprint
}

// BeforeCreate GORM 钩子
func (s *Sprint) BeforeCreate(tx *gorm.DB) error {
	// 设置默认值
	if s.Status == "" {
		s.Status = SprintStatusPlanning
	}
	if s.CompletedPoints == nil {
		zero := 0
		s.CompletedPoints = &zero
	}

	return nil
}

// 业务方法

// IsActive 检查冲刺是否处于活跃状态
func (s *Sprint) IsActive() bool {
	return s.Status == SprintStatusActive
}

// IsCompleted 检查冲刺是否已完成
func (s *Sprint) IsCompleted() bool {
	return s.Status == SprintStatusCompleted
}

// IsPlanning 检查冲刺是否处于规划状态
func (s *Sprint) IsPlanning() bool {
	return s.Status == SprintStatusPlanning
}

// IsCancelled 检查冲刺是否被取消
func (s *Sprint) IsCancelled() bool {
	return s.Status == SprintStatusCancelled
}

// Duration 计算冲刺持续时间
func (s *Sprint) Duration() time.Duration {
	return s.EndDate.Sub(s.StartDate)
}

// DurationDays 计算冲刺持续天数
func (s *Sprint) DurationDays() int {
	return int(s.Duration().Hours() / 24)
}

// IsStarted 检查冲刺是否已开始
func (s *Sprint) IsStarted() bool {
	return time.Now().After(s.StartDate) || time.Now().Equal(s.StartDate)
}

// IsEnded 检查冲刺是否已结束
func (s *Sprint) IsEnded() bool {
	return time.Now().After(s.EndDate)
}

// DaysRemaining 计算剩余天数
func (s *Sprint) DaysRemaining() int {
	if s.IsEnded() {
		return 0
	}
	remaining := time.Until(s.EndDate).Hours() / 24
	if remaining < 0 {
		return 0
	}
	return int(remaining)
}

// DaysElapsed 计算已过天数
func (s *Sprint) DaysElapsed() int {
	if !s.IsStarted() {
		return 0
	}
	elapsed := time.Since(s.StartDate).Hours() / 24
	if elapsed < 0 {
		return 0
	}
	return int(elapsed)
}

// Progress 计算进度百分比（基于时间）
func (s *Sprint) Progress() float64 {
	if !s.IsStarted() {
		return 0
	}
	if s.IsEnded() {
		return 100
	}

	total := s.Duration().Hours()
	elapsed := time.Since(s.StartDate).Hours()

	if total <= 0 {
		return 0
	}

	progress := (elapsed / total) * 100
	if progress > 100 {
		return 100
	}
	if progress < 0 {
		return 0
	}

	return progress
}

// TaskProgress 计算任务完成进度百分比
func (s *Sprint) TaskProgress() float64 {
	if s.CommittedPoints == nil || *s.CommittedPoints == 0 {
		return 0
	}

	completedPoints := 0
	if s.CompletedPoints != nil {
		completedPoints = *s.CompletedPoints
	}

	return float64(completedPoints) / float64(*s.CommittedPoints) * 100
}

// CanTransitionTo 检查是否可以转换到指定状态
func (s *Sprint) CanTransitionTo(status SprintStatus) bool {
	switch s.Status {
	case SprintStatusPlanning:
		return status == SprintStatusActive || status == SprintStatusCancelled
	case SprintStatusActive:
		return status == SprintStatusCompleted || status == SprintStatusCancelled
	case SprintStatusCompleted:
		// 通常完成的冲刺不能再转换状态，但允许重新激活用于修正
		return status == SprintStatusActive
	case SprintStatusCancelled:
		// 取消的冲刺可以重新规划
		return status == SprintStatusPlanning
	default:
		return false
	}
}

// CanAddTask 检查是否可以添加任务
func (s *Sprint) CanAddTask() bool {
	// 只有规划中或活跃的冲刺可以添加任务
	return s.Status == SprintStatusPlanning || s.Status == SprintStatusActive
}

// CanRemoveTask 检查是否可以移除任务
func (s *Sprint) CanRemoveTask() bool {
	// 只有规划中或活跃的冲刺可以移除任务
	return s.Status == SprintStatusPlanning || s.Status == SprintStatusActive
}

// GetTaskStats 获取任务统计信息
func (s *Sprint) GetTaskStats() *TaskStats {
	stats := &TaskStats{}

	for _, task := range s.Tasks {
		stats.Total++

		switch task.Status {
		case TaskStatusTodo:
			stats.Todo++
		case TaskStatusInProgress:
			stats.InProgress++
		case TaskStatusReview:
			stats.Review++
		case TaskStatusTesting:
			stats.Testing++
		case TaskStatusBlocked:
			stats.Blocked++
		case TaskStatusDone:
			stats.Done++
		}

		if task.StoryPoints != nil {
			stats.TotalPoints += *task.StoryPoints
			if task.IsCompleted() {
				stats.CompletedPoints += *task.StoryPoints
			}
		}
	}

	// 计算完成百分比
	if stats.Total > 0 {
		stats.CompletionRate = float64(stats.Done) / float64(stats.Total) * 100
	}

	if stats.TotalPoints > 0 {
		stats.PointsCompletionRate = float64(stats.CompletedPoints) / float64(stats.TotalPoints) * 100
	}

	return stats
}

// GetVelocity 计算团队速度（已完成的故事点）
func (s *Sprint) GetVelocity() int {
	if s.CompletedPoints != nil {
		return *s.CompletedPoints
	}

	velocity := 0
	for _, task := range s.Tasks {
		if task.IsCompleted() && task.StoryPoints != nil {
			velocity += *task.StoryPoints
		}
	}
	return velocity
}

// GetCommittedPoints 获取承诺的故事点总数
func (s *Sprint) GetCommittedPoints() int {
	if s.CommittedPoints != nil {
		return *s.CommittedPoints
	}

	points := 0
	for _, task := range s.Tasks {
		if task.StoryPoints != nil {
			points += *task.StoryPoints
		}
	}
	return points
}

// Validate 验证冲刺数据
func (s *Sprint) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("name is required")
	}

	if len(s.Name) > 100 {
		return fmt.Errorf("name too long (max 100 characters)")
	}

	if !s.Status.IsValid() {
		return fmt.Errorf("invalid status: %s", s.Status)
	}

	if s.EndDate.Before(s.StartDate) {
		return fmt.Errorf("end date cannot be before start date")
	}

	if s.Capacity != nil && *s.Capacity < 0 {
		return fmt.Errorf("capacity cannot be negative")
	}

	if s.CommittedPoints != nil && *s.CommittedPoints < 0 {
		return fmt.Errorf("committed points cannot be negative")
	}

	return nil
}

// TaskStats 任务统计结构
type TaskStats struct {
	Total                int     `json:"total"`
	Todo                 int     `json:"todo"`
	InProgress           int     `json:"in_progress"`
	Review               int     `json:"review"`
	Testing              int     `json:"testing"`
	Blocked              int     `json:"blocked"`
	Done                 int     `json:"done"`
	TotalPoints          int     `json:"total_points"`
	CompletedPoints      int     `json:"completed_points"`
	CompletionRate       float64 `json:"completion_rate"`        // 任务完成率
	PointsCompletionRate float64 `json:"points_completion_rate"` // 故事点完成率
}

// 枚举验证方法

// IsValid 验证冲刺状态是否有效
func (ss SprintStatus) IsValid() bool {
	switch ss {
	case SprintStatusPlanning, SprintStatusActive, SprintStatusCompleted, SprintStatusCancelled:
		return true
	default:
		return false
	}
}

// String 实现 Stringer 接口
func (ss SprintStatus) String() string {
	return string(ss)
}
