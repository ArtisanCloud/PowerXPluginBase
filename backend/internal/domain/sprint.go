package domain

import (
	"fmt"
	"time"

	"github.com/powerx-plugins/scrum/internal/db"
	"gorm.io/gorm"
)

// Sprint Sprint 模型
type Sprint struct {
	db.BaseModel
	Name      string       `json:"name" gorm:"not null;size:100"`
	Goal      string       `json:"goal" gorm:"type:text"`
	StartDate time.Time    `json:"start_date" gorm:"not null"`
	EndDate   time.Time    `json:"end_date" gorm:"not null"`
	Capacity  *int         `json:"capacity,omitempty"` // 团队容量（故事点）
	Status    SprintStatus `json:"status" gorm:"not null;default:'planning'"`

	// 关联
	Tasks []Task `json:"tasks,omitempty" gorm:"foreignKey:SprintID"`
}

// SprintStatus Sprint 状态枚举
type SprintStatus string

const (
	SprintStatusPlanning  SprintStatus = "planning"
	SprintStatusActive    SprintStatus = "active"
	SprintStatusCompleted SprintStatus = "completed"
)

// TableName 返回表名
func (Sprint) TableName() string {
	return "sprint"
}

// BeforeCreate GORM 钩子
func (s *Sprint) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的钩子
	if err := s.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 设置默认值
	if s.Status == "" {
		s.Status = SprintStatusPlanning
	}

	return nil
}

// IsValidStatus 验证状态是否有效
func (ss SprintStatus) IsValid() bool {
	switch ss {
	case SprintStatusPlanning, SprintStatusActive, SprintStatusCompleted:
		return true
	default:
		return false
	}
}

// String 实现 Stringer 接口
func (ss SprintStatus) String() string {
	return string(ss)
}

// IsActive 检查 Sprint 是否处于活跃状态
func (s *Sprint) IsActive() bool {
	return s.Status == SprintStatusActive
}

// IsCompleted 检查 Sprint 是否已完成
func (s *Sprint) IsCompleted() bool {
	return s.Status == SprintStatusCompleted
}

// IsPlanning 检查 Sprint 是否处于规划状态
func (s *Sprint) IsPlanning() bool {
	return s.Status == SprintStatusPlanning
}

// Duration 计算 Sprint 持续时间
func (s *Sprint) Duration() time.Duration {
	return s.EndDate.Sub(s.StartDate)
}

// DurationDays 计算 Sprint 持续天数
func (s *Sprint) DurationDays() int {
	return int(s.Duration().Hours() / 24)
}

// IsStarted 检查 Sprint 是否已开始
func (s *Sprint) IsStarted() bool {
	return time.Now().After(s.StartDate) || time.Now().Equal(s.StartDate)
}

// IsEnded 检查 Sprint 是否已结束
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

// CanTransitionTo 检查是否可以转换到指定状态
func (s *Sprint) CanTransitionTo(status SprintStatus) bool {
	switch s.Status {
	case SprintStatusPlanning:
		return status == SprintStatusActive
	case SprintStatusActive:
		return status == SprintStatusCompleted
	case SprintStatusCompleted:
		// 通常完成的 Sprint 不能再转换状态，但这里允许重新激活
		return status == SprintStatusActive
	default:
		return false
	}
}

// Validate 验证 Sprint 数据
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

	return nil
}

// CanAddTask 检查是否可以添加任务
func (s *Sprint) CanAddTask() bool {
	// 只有规划中或活跃的 Sprint 可以添加任务
	return s.Status == SprintStatusPlanning || s.Status == SprintStatusActive
}

// CanRemoveTask 检查是否可以移除任务
func (s *Sprint) CanRemoveTask() bool {
	// 只有规划中或活跃的 Sprint 可以移除任务
	return s.Status == SprintStatusPlanning || s.Status == SprintStatusActive
}

// GetVelocity 计算团队速度（已完成的故事点）
func (s *Sprint) GetVelocity() int {
	velocity := 0
	for _, task := range s.Tasks {
		if task.IsCompleted() && task.Estimate != nil {
			velocity += *task.Estimate
		}
	}
	return velocity
}

// GetCommittedPoints 获取承诺的故事点总数
func (s *Sprint) GetCommittedPoints() int {
	points := 0
	for _, task := range s.Tasks {
		if task.Estimate != nil {
			points += *task.Estimate
		}
	}
	return points
}

// GetTaskStats 获取任务统计
func (s *Sprint) GetTaskStats() TaskStats {
	stats := TaskStats{}

	for _, task := range s.Tasks {
		stats.Total++

		switch task.Status {
		case TaskStatusTodo:
			stats.Todo++
		case TaskStatusInProgress:
			stats.InProgress++
		case TaskStatusDone:
			stats.Done++
		}

		if task.Estimate != nil {
			stats.TotalPoints += *task.Estimate
			if task.IsCompleted() {
				stats.CompletedPoints += *task.Estimate
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

// TaskStats 任务统计结构
type TaskStats struct {
	Total                int     `json:"total"`
	Todo                 int     `json:"todo"`
	InProgress           int     `json:"in_progress"`
	Done                 int     `json:"done"`
	TotalPoints          int     `json:"total_points"`
	CompletedPoints      int     `json:"completed_points"`
	CompletionRate       float64 `json:"completion_rate"`        // 任务完成率
	PointsCompletionRate float64 `json:"points_completion_rate"` // 故事点完成率
}
