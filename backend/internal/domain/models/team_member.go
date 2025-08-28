package models

import (
	"fmt"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// TeamMember 团队成员模型
type TeamMember struct {
	BaseModel

	// 关联信息
	ProjectID uint  `gorm:"column:project_id;not null;index" json:"project_id"`
	UserID    int64 `gorm:"column:user_id;not null;index" json:"user_id"`

	// 角色和权限
	Role        TeamRole `gorm:"column:role;type:varchar(20);not null;default:'developer';index" json:"role"`
	Permissions string   `gorm:"column:permissions;type:text" json:"permissions,omitempty"`

	// 工作安排
	Status       MemberStatus `gorm:"column:status;type:varchar(20);not null;default:'active';index" json:"status"`
	Capacity     *int         `gorm:"column:capacity" json:"capacity,omitempty"`         // 工作容量（故事点/冲刺）
	Availability *float64     `gorm:"column:availability" json:"availability,omitempty"` // 可用性百分比（0-100）

	// 时间管理
	JoinedAt *time.Time `gorm:"column:joined_at" json:"joined_at,omitempty"`
	LeftAt   *time.Time `gorm:"column:left_at" json:"left_at,omitempty"`

	// 元数据
	Meta datatypes.JSON `gorm:"column:meta;type:jsonb" json:"meta,omitempty"`

	// 关联模型
	Project *Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
}

// TeamRole 团队角色枚举
type TeamRole string

const (
	TeamRoleScrumMaster  TeamRole = "scrum_master"
	TeamRoleProductOwner TeamRole = "product_owner"
	TeamRoleDeveloper    TeamRole = "developer"
	TeamRoleTester       TeamRole = "tester"
	TeamRoleDesigner     TeamRole = "designer"
	TeamRoleAnalyst      TeamRole = "analyst"
	TeamRoleArchitect    TeamRole = "architect"
)

// MemberStatus 成员状态枚举
type MemberStatus string

const (
	MemberStatusActive     MemberStatus = "active"
	MemberStatusInactive   MemberStatus = "inactive"
	MemberStatusOnVacation MemberStatus = "on_vacation"
	MemberStatusLeft       MemberStatus = "left"
)

// TableName 返回表名
func (tm *TeamMember) TableName() string {
	return PowerXSchema + "." + TableTeamMember
}

// GetTableName 获取表名
func (tm *TeamMember) GetTableName(needFull bool) string {
	if needFull {
		return tm.TableName()
	}
	return TableTeamMember
}

// BeforeCreate GORM 钩子
func (tm *TeamMember) BeforeCreate(tx *gorm.DB) error {
	// 设置默认值
	if tm.Role == "" {
		tm.Role = TeamRoleDeveloper
	}
	if tm.Status == "" {
		tm.Status = MemberStatusActive
	}
	if tm.JoinedAt == nil {
		now := time.Now()
		tm.JoinedAt = &now
	}
	if tm.Availability == nil {
		full := 100.0
		tm.Availability = &full
	}

	return nil
}

// 业务方法

// IsActive 检查成员是否活跃
func (tm *TeamMember) IsActive() bool {
	return tm.Status == MemberStatusActive
}

// IsOnVacation 检查成员是否在休假
func (tm *TeamMember) IsOnVacation() bool {
	return tm.Status == MemberStatusOnVacation
}

// HasLeft 检查成员是否已离开
func (tm *TeamMember) HasLeft() bool {
	return tm.Status == MemberStatusLeft
}

// GetEffectiveCapacity 获取有效工作容量
func (tm *TeamMember) GetEffectiveCapacity() int {
	if tm.Capacity == nil || tm.Availability == nil {
		return 0
	}

	if !tm.IsActive() {
		return 0
	}

	return int(float64(*tm.Capacity) * (*tm.Availability / 100.0))
}

// GetWorkDuration 获取工作时长
func (tm *TeamMember) GetWorkDuration() *time.Duration {
	if tm.JoinedAt == nil {
		return nil
	}

	var endTime time.Time
	if tm.LeftAt != nil {
		endTime = *tm.LeftAt
	} else {
		endTime = time.Now()
	}

	duration := endTime.Sub(*tm.JoinedAt)
	return &duration
}

// IsScrumMaster 检查是否为 Scrum Master
func (tm *TeamMember) IsScrumMaster() bool {
	return tm.Role == TeamRoleScrumMaster
}

// IsProductOwner 检查是否为产品负责人
func (tm *TeamMember) IsProductOwner() bool {
	return tm.Role == TeamRoleProductOwner
}

// IsDeveloper 检查是否为开发人员
func (tm *TeamMember) IsDeveloper() bool {
	return tm.Role == TeamRoleDeveloper
}

// IsTester 检查是否为测试人员
func (tm *TeamMember) IsTester() bool {
	return tm.Role == TeamRoleTester
}

// CanManageSprint 检查是否可以管理冲刺
func (tm *TeamMember) CanManageSprint() bool {
	return tm.IsScrumMaster() || tm.IsProductOwner()
}

// CanManageBacklog 检查是否可以管理待办列表
func (tm *TeamMember) CanManageBacklog() bool {
	return tm.IsProductOwner() || tm.IsScrumMaster()
}

// CanAssignTasks 检查是否可以分配任务
func (tm *TeamMember) CanAssignTasks() bool {
	return tm.IsScrumMaster() || tm.IsProductOwner()
}

// CanTransitionTo 检查是否可以转换到指定状态
func (tm *TeamMember) CanTransitionTo(status MemberStatus) bool {
	switch tm.Status {
	case MemberStatusActive:
		return status == MemberStatusInactive || status == MemberStatusOnVacation || status == MemberStatusLeft
	case MemberStatusInactive:
		return status == MemberStatusActive || status == MemberStatusLeft
	case MemberStatusOnVacation:
		return status == MemberStatusActive || status == MemberStatusInactive || status == MemberStatusLeft
	case MemberStatusLeft:
		// 通常离开的成员不能再转换状态，但允许重新加入
		return status == MemberStatusActive
	default:
		return false
	}
}

// Leave 成员离开项目
func (tm *TeamMember) Leave() {
	tm.Status = MemberStatusLeft
	now := time.Now()
	tm.LeftAt = &now
}

// Validate 验证团队成员数据
func (tm *TeamMember) Validate() error {
	if tm.ProjectID == 0 {
		return fmt.Errorf("project_id is required")
	}

	if tm.UserID == 0 {
		return fmt.Errorf("user_id is required")
	}

	if !tm.Role.IsValid() {
		return fmt.Errorf("invalid role: %s", tm.Role)
	}

	if !tm.Status.IsValid() {
		return fmt.Errorf("invalid status: %s", tm.Status)
	}

	if tm.Capacity != nil && *tm.Capacity < 0 {
		return fmt.Errorf("capacity cannot be negative")
	}

	if tm.Availability != nil && (*tm.Availability < 0 || *tm.Availability > 100) {
		return fmt.Errorf("availability must be between 0 and 100")
	}

	if tm.JoinedAt != nil && tm.LeftAt != nil && tm.LeftAt.Before(*tm.JoinedAt) {
		return fmt.Errorf("left date cannot be before joined date")
	}

	return nil
}

// 枚举验证方法

// IsValid 验证团队角色是否有效
func (tr TeamRole) IsValid() bool {
	switch tr {
	case TeamRoleScrumMaster, TeamRoleProductOwner, TeamRoleDeveloper,
		TeamRoleTester, TeamRoleDesigner, TeamRoleAnalyst, TeamRoleArchitect:
		return true
	default:
		return false
	}
}

// IsValid 验证成员状态是否有效
func (ms MemberStatus) IsValid() bool {
	switch ms {
	case MemberStatusActive, MemberStatusInactive, MemberStatusOnVacation, MemberStatusLeft:
		return true
	default:
		return false
	}
}

// String 实现 Stringer 接口
func (tr TeamRole) String() string {
	return string(tr)
}

func (ms MemberStatus) String() string {
	return string(ms)
}

// GetRoleDisplayName 获取角色显示名称
func (tr TeamRole) GetRoleDisplayName() string {
	switch tr {
	case TeamRoleScrumMaster:
		return "Scrum Master"
	case TeamRoleProductOwner:
		return "Product Owner"
	case TeamRoleDeveloper:
		return "Developer"
	case TeamRoleTester:
		return "Tester"
	case TeamRoleDesigner:
		return "Designer"
	case TeamRoleAnalyst:
		return "Business Analyst"
	case TeamRoleArchitect:
		return "Solution Architect"
	default:
		return string(tr)
	}
}
