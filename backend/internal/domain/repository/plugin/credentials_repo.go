package plugin

import (
    "context"
    "errors"
    "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

type CredentialsRepo struct {
    db *gorm.DB
}

func NewCredentialsRepo(db *gorm.DB) *CredentialsRepo { return &CredentialsRepo{db: db} }

// Upsert 保存/更新某租户-插件的凭证（密文）
func (r *CredentialsRepo) Upsert(ctx context.Context, pc *models.PluginCredential) error {
    if pc == nil {
        return errors.New("nil model")
    }
    return r.db.WithContext(ctx).Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "tenant_id"}, {Name: "plugin_id"}},
        DoUpdates: clause.AssignmentColumns([]string{"client_id", "secret_ciphertext", "iv_nonce", "key_version", "updated_at"}),
    }).Create(pc).Error
}

// GetByTenantPlugin 获取某租户-插件的凭证
func (r *CredentialsRepo) GetByTenantPlugin(ctx context.Context, tenantID int64, pluginID string) (*models.PluginCredential, error) {
    var pc models.PluginCredential
    if err := r.db.WithContext(ctx).Where("tenant_id = ? AND plugin_id = ?", tenantID, pluginID).First(&pc).Error; err != nil {
        return nil, err
    }
    return &pc, nil
}

