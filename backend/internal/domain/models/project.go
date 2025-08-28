package models

import (
	"fmt"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Project 项目模型
type Project struct {
	BaseModel

	// 基本信息
	Key         string        `gorm:"column:key;type:varchar(20);uniqueIndex;not null" json:"key"`
	Name        string        `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Description string        `gorm:"column:description;type:text" json:"description,omitempty"`
	Status      ProjectStatus `gorm:"column:status;type:varchar(20);not null;default:'active';index" json:"status"`

	// 项目管理
	LeadID     *int64     `gorm:"column:lead_id;index" json:"lead_id,omitempty"`   // 项目负责人
	StartDate  *time.Time `gorm:"column:start_date" json:"start_date,omitempty"`   // 项目开始日期
	EndDate    *time.Time `gorm:"column:end_date" json:"end_date,omitempty"`       // 项目结束日期
	Budget     *float64   `gorm:"column:budget" json:"budget,omitempty"`           // 项目预算
	ActualCost *float64   `gorm:"column:actual_cost;default:0" json:"actual_cost"` // 实际成本

	// 默认设置
	DefaultSprintLength int `gorm:"column:default_sprint_length;default:14" json:"default_sprint_length"` // 默认冲刺长度（天）
	SprintStartDay      int `gorm:"column:sprint_start_day;default:1" json:"sprint_start_day"`            // 冲刺开始日（1=周一）

	// 权限和可见性
	IsPublic   bool `gorm:"column:is_public;default:false" json:"is_public"`
	IsArchived bool `gorm:"column:is_archived;default:false;index" json:"is_archived"`

	// 元数据
	Meta datatypes.JSON `gorm:"column:meta;type:jsonb" json:"meta,omitempty"`

	// 关联模型
	Sprints     []Sprint     `json:"sprints,omitempty" gorm:"foreignKey:ProjectID"`
	Tasks       []Task       `json:"tasks,omitempty" gorm:"foreignKey:ProjectID"`
	TeamMembers []TeamMember `json:"team_members,omitempty" gorm:"foreignKey:ProjectID"`
}

// ProjectStatus 项目状态枚举
type ProjectStatus string

const (
	ProjectStatusActive    ProjectStatus = "active"    // 活跃
	ProjectStatusCompleted ProjectStatus = "completed" // 已完成
	ProjectStatusOnHold    ProjectStatus = "on_hold"   // 暂停
	ProjectStatusCancelled ProjectStatus = "cancelled" // 已取消
)

// TableName 返回表名
func (p *Project) TableName() string {
	return PowerXSchema + "." + TableProject
}

// GetTableName 获取表名
func (p *Project) GetTableName(needFull bool) string {
	if needFull {
		return p.TableName()
	}
	return TableProject
}

// BeforeCreate GORM 钩子
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	// 设置默认值
	if p.Status == "" {
		p.Status = ProjectStatusActive
	}
	if p.DefaultSprintLength == 0 {
		p.DefaultSprintLength = 14 // 默认 2 周
	}
	if p.SprintStartDay == 0 {
		p.SprintStartDay = 1 // 默认周一开始
	}
	if p.ActualCost == nil {
		zero := 0.0
		p.ActualCost = &zero
	}

	return nil
}

// 业务方法

// IsActive 检查项目是否活跃
func (p *Project) IsActive() bool {
	return p.Status == ProjectStatusActive && !p.IsArchived
}

// IsCompleted 检查项目是否已完成
func (p *Project) IsCompleted() bool {
	return p.Status == ProjectStatusCompleted
}

// IsOnHold 检查项目是否暂停
func (p *Project) IsOnHold() bool {
	return p.Status == ProjectStatusOnHold
}

// IsCancelled 检查项目是否被取消
func (p *Project) IsCancelled() bool {
	return p.Status == ProjectStatusCancelled
}

// Duration 计算项目持续时间
func (p *Project) Duration() *time.Duration {
	if p.StartDate == nil || p.EndDate == nil {
		return nil
	}
	duration := p.EndDate.Sub(*p.StartDate)
	return &duration
}

// DaysRemaining 计算剩余天数
func (p *Project) DaysRemaining() *int {
	if p.EndDate == nil {
		return nil
	}
	remaining := time.Until(*p.EndDate).Hours() / 24
	if remaining < 0 {
		zero := 0
		return &zero
	}
	days := int(remaining)
	return &days
}

// Progress 计算项目进度（基于时间）
func (p *Project) Progress() float64 {
	if p.StartDate == nil || p.EndDate == nil {
		return 0
	}

	now := time.Now()
	if now.Before(*p.StartDate) {
		return 0
	}
	if now.After(*p.EndDate) {
		return 100
	}

	total := p.EndDate.Sub(*p.StartDate).Hours()
	elapsed := now.Sub(*p.StartDate).Hours()

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

// GetActiveSprint 获取当前活跃的冲刺
func (p *Project) GetActiveSprint() *Sprint {
	for _, sprint := range p.Sprints {
		if sprint.IsActive() {
			return &sprint
		}
	}
	return nil
}

// GetCompletedSprints 获取已完成的冲刺
func (p *Project) GetCompletedSprints() []Sprint {
	var completed []Sprint
	for _, sprint := range p.Sprints {
		if sprint.IsCompleted() {
			completed = append(completed, sprint)
		}
	}
	return completed
}

// GetTaskStats 获取项目任务统计
func (p *Project) GetTaskStats() *TaskStats {
	stats := &TaskStats{}

	for _, task := range p.Tasks {
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

// GetVelocityHistory 获取团队速度历史记录
func (p *Project) GetVelocityHistory(limit int) []int {
	var velocities []int
	completedSprints := p.GetCompletedSprints()

	// 按结束时间排序
	for i := len(completedSprints) - 1; i >= 0 && len(velocities) < limit; i-- {
		velocities = append(velocities, completedSprints[i].GetVelocity())
	}

	return velocities
}

// CanTransitionTo 检查是否可以转换到指定状态
func (p *Project) CanTransitionTo(status ProjectStatus) bool {
	switch p.Status {
	case ProjectStatusActive:
		return status == ProjectStatusCompleted || status == ProjectStatusOnHold || status == ProjectStatusCancelled
	case ProjectStatusOnHold:
		return status == ProjectStatusActive || status == ProjectStatusCancelled
	case ProjectStatusCompleted:
		return status == ProjectStatusActive // 允许重新激活已完成的项目
	case ProjectStatusCancelled:
		return status == ProjectStatusActive // 允许重新激活已取消的项目
	default:
		return false
	}
}

// Validate 验证项目数据
func (p *Project) Validate() error {
	if p.Key == "" {
		return fmt.Errorf("key is required")
	}

	if len(p.Key) > 20 {
		return fmt.Errorf("key too long (max 20 characters)")
	}

	if p.Name == "" {
		return fmt.Errorf("name is required")
	}

	if len(p.Name) > 100 {
		return fmt.Errorf("name too long (max 100 characters)")
	}

	if !p.Status.IsValid() {
		return fmt.Errorf("invalid status: %s", p.Status)
	}

	if p.StartDate != nil && p.EndDate != nil && p.EndDate.Before(*p.StartDate) {
		return fmt.Errorf("end date cannot be before start date")
	}

	if p.DefaultSprintLength < 1 || p.DefaultSprintLength > 30 {
		return fmt.Errorf("default sprint length must be between 1 and 30 days")
	}

	if p.SprintStartDay < 1 || p.SprintStartDay > 7 {
		return fmt.Errorf("sprint start day must be between 1 (Monday) and 7 (Sunday)")
	}

	if p.Budget != nil && *p.Budget < 0 {
		return fmt.Errorf("budget cannot be negative")
	}

	return nil
}

// 枚举验证方法

// IsValid 验证项目状态是否有效
func (ps ProjectStatus) IsValid() bool {
	switch ps {
	case ProjectStatusActive, ProjectStatusCompleted, ProjectStatusOnHold, ProjectStatusCancelled:
		return true
	default:
		return false
	}
}

// String 实现 Stringer 接口
func (ps ProjectStatus) String() string {
	return string(ps)
}
