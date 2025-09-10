// internal/domain/models/iam/iam_cache_ext.go
package iam

import (
	"github.com/ArtisanCloud/PowerXPlugin/internal/domain/models"
	"time"

	"gorm.io/datatypes"
)

// =========================
// TeamCache：团队轻量缓存
// 说明：使用 BaseModel 的自增 ID 作为主键；业务唯一性依赖 (tenant_id, team_id)
// 的唯一索引，请在迁移里用 SQL 创建（见文末）。
// =========================
type TeamCache struct {
	models.BaseModel

	TeamID      uint64     `gorm:"not null;index;comment:底座Team主键"            json:"team_id"`
	Name        string     `gorm:"type:varchar(128);comment:团队名称"              json:"name"`
	Path        string     `gorm:"type:text;comment:层级路径/展示用"                json:"path"`
	Ver         int64      `gorm:"type:bigint;not null;default:0;comment:逻辑版本/etag" json:"ver"`
	RefreshedAt *time.Time `gorm:"type:timestamptz;comment:最近一次刷新时间"          json:"refreshed_at"`
}

func (TeamCache) TableName() string { return models.S(models.TableTeamCache) }

// =========================
// TeamExt：团队插件扩展表
// 说明：使用 BaseModel.ID 为主键；业务唯一 (tenant_id, team_id) 建唯一索引。
// =========================
type TeamExt struct {
	models.BaseModel

	TeamID uint64         `gorm:"not null;index;comment:底座Team主键"  json:"team_id"`
	Alias  string         `gorm:"type:varchar(128);comment:插件内别名" json:"alias"`
	Meta   datatypes.JSON `gorm:"type:jsonb;comment:插件内配置/元数据"  json:"meta"`
	Tags   datatypes.JSON `gorm:"type:jsonb;comment:标签(插件维度)"      json:"tags"`
	Status int16          `gorm:"type:smallint;not null;default:1;comment:插件视角状态" json:"status"`
}

func (TeamExt) TableName() string { return models.S(models.TableTeamExt) }

// =========================
// MemberExt：成员插件扩展表
// 说明：使用 BaseModel.ID 为主键；业务唯一 (tenant_id, member_id) 建唯一索引。
// =========================
type MemberExt struct {
	models.BaseModel

	MemberID uint64         `gorm:"not null;index;comment:底座Member主键" json:"member_id"`
	Alias    string         `gorm:"type:varchar(128);comment:插件内别名"   json:"alias"`
	Prefs    datatypes.JSON `gorm:"type:jsonb;comment:个性化偏好"          json:"prefs"`
	Tags     datatypes.JSON `gorm:"type:jsonb;comment:插件内标签"          json:"tags"`
	Status   int16          `gorm:"type:smallint;not null;default:1;comment:插件视角状态" json:"status"`
}

func (MemberExt) TableName() string { return models.S(models.TableMemberExt) }

// =========================
// MemberCache：成员轻量缓存
// 说明：使用 BaseModel.ID 为主键；业务唯一 (tenant_id, member_id) 建唯一索引。
// =========================
type MemberCache struct {
	models.BaseModel

	MemberID    uint64         `gorm:"not null;index;comment:底座Member主键" json:"member_id"`
	Username    string         `gorm:"type:varchar(128);comment:用户名"       json:"username"`
	DisplayName string         `gorm:"type:varchar(128);comment:昵称/显示名"  json:"display_name"`
	AvatarURL   string         `gorm:"type:text;comment:头像URL"              json:"avatar_url"`
	TeamIDs     datatypes.JSON `gorm:"type:jsonb;comment:团队ID快照(可选)"    json:"team_ids"`

	Ver         int64      `gorm:"type:bigint;not null;default:0;comment:逻辑版本/etag" json:"ver"`
	RefreshedAt *time.Time `gorm:"type:timestamptz;comment:最近一次刷新时间"              json:"refreshed_at"`
}

func (MemberCache) TableName() string { return models.S(models.TableMemberCache) }
