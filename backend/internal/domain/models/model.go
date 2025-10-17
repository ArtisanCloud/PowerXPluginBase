package models

// backend/internal/domain/models/model.go

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型，包含通用字段
type BaseModel struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement;comment:自增ID" json:"id"`
	TenantID  uint64         `gorm:"not null;index;comment:租户ID"           json:"tenant_id"`
	CreatedAt time.Time      `gorm:"autoCreateTime;comment:创建时间"          json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;comment:更新时间"          json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:软删除时间"                  json:"deleted_at,omitempty"`
}

type BaseNoTenantModel struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

const (
	TablePluginTenantExt            = "plugin_tenant_ext"
	TableTemplate                   = "template"
	TablePluginCredentials          = "plugin_credentials"
	TablePrivacyDataClassifications = "privacy_data_classifications"
	TablePrivacyConsentTokens       = "privacy_consent_tokens"
	TablePrivacyLifecycleEvents     = "privacy_lifecycle_events"
	TableSecurityBaselineChecklists = "security_baseline_checklists"
	TableSecurityAuditReports       = "security_audit_reports"
	TableToolGrantRevocations       = "tool_grant_revocations"
	TableToolGrantUsageEvents       = "tool_grant_usage_events"
)
