package admin_console

import (
	"context"

	model "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models/admin_console"
	"github.com/ArtisanCloud/PowerXPlugin/internal/domain/repository"
	"gorm.io/gorm"
)

// AuditRepository persists admin console audit events.
type AuditRepository struct {
	*repository.BaseRepository[model.AuditEvent]
}

// NewAuditRepository constructs the audit repository.
func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{BaseRepository: repository.NewBaseRepository[model.AuditEvent](db)}
}

// Create inserts a new audit event record.
func (r *AuditRepository) Create(ctx context.Context, evt *model.AuditEvent) error {
	return r.DB.WithContext(ctx).Create(evt).Error
}

// LatestForAction returns the latest audit entry for a resource reference.
func (r *AuditRepository) LatestForAction(ctx context.Context, pluginID string, tenantID *string, resourceType, resourceRef string) (*model.AuditEvent, error) {
	if r.DB == nil {
		return nil, gorm.ErrInvalidDB
	}
	query := r.DB.WithContext(ctx).
		Where("plugin_id = ? AND resource_type = ?", pluginID, resourceType).
		Order("occurred_at DESC")
	if resourceRef != "" {
		query = query.Where("resource_ref = ?", resourceRef)
	}
	if tenantID == nil {
		query = query.Where("tenant_id IS NULL")
	} else {
		query = query.Where("tenant_id = ?", *tenantID)
	}
	var evt model.AuditEvent
	if err := query.First(&evt).Error; err != nil {
		return nil, err
	}
	return &evt, nil
}
