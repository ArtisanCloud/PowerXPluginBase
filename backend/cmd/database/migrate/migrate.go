package migrate

import (
	"context"

	"github.com/ArtisanCloud/PowerXPlugin/internal/config"
	"github.com/ArtisanCloud/PowerXPlugin/internal/domain/models"
	runtimeopsmodel "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models/runtime_ops"
	templatemodel "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models/template"
	"gorm.io/gorm"
)

// MigratePluginModels 只做 AutoMigrate（最小实现）
func MigratePluginModels(ctx context.Context, db *gorm.DB) error {
	return db.AutoMigrate(
		&models.PluginCredential{},
		&models.PluginTenantExt{},
		&templatemodel.Template{},
		&runtimeopsmodel.MCPSession{},
		&runtimeopsmodel.RuntimeAuditEvent{},
		&runtimeopsmodel.QuotaLedger{},
		&runtimeopsmodel.MarketplaceOverage{},
	)
}

func ResetDatabase(ctx context.Context, db *gorm.DB, cfg *config.DatabaseConfig) error {
	// 如果你用 GORM，可以直接 drop 所有表
	// 或者先获取表名，再循环 drop
	// 这里举例简单版本：
	err := db.Exec("DROP SCHEMA " + cfg.Schema + " CASCADE; CREATE SCHEMA " + cfg.Schema + ";").Error
	if err != nil {
		return err
	}
	return nil
}
