package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型，参考 PowerX 的设计
type BaseModel struct {
	ID        uint64         `gorm:"primarykey" json:"id"`
	TenantID  uint64         `json:"tenant_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// PowerXSchema 数据库模式名
const PowerXSchema = "public"

// 表名常量
const (
	TableTask         = "scrum_task"
	TableSprint       = "scrum_sprint"
	TableProject      = "scrum_project"
	TableTeamMember   = "scrum_team_member"
	TableTaskActivity = "scrum_task_activity"
	TableSprintReport = "scrum_sprint_report"
)
