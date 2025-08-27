package db

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型，包含通用字段
type BaseModel struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	TenantID  int64          `json:"tenant_id" gorm:"not null;index:idx_tenant"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// SetTenantID 设置租户 ID
func (bm *BaseModel) SetTenantID(tenantID int64) {
	bm.TenantID = tenantID
}

// GetTenantID 获取租户 ID
func (bm *BaseModel) GetTenantID() int64 {
	return bm.TenantID
}

// BeforeCreate GORM 钩子：创建前自动设置租户 ID
func (bm *BaseModel) BeforeCreate(tx *gorm.DB) error {
	// 如果租户 ID 未设置，尝试从会话中获取
	if bm.TenantID == 0 {
		if tenantID, err := GetCurrentTenantID(tx); err == nil {
			bm.TenantID = tenantID
		}
	}
	return nil
}

// TenantModel 租户模型接口
type TenantModel interface {
	SetTenantID(tenantID int64)
	GetTenantID() int64
}

// Ensure BaseModel implements TenantModel
var _ TenantModel = (*BaseModel)(nil)
