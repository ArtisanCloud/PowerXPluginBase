package bootstrap

import (
	"context"

	"github.com/ArtisanCloud/PowerXPlugin/internal/config"
	"github.com/ArtisanCloud/PowerXPlugin/internal/db"
	"github.com/ArtisanCloud/PowerXPlugin/internal/domain/models"
	"github.com/ArtisanCloud/PowerXPlugin/internal/logger"
	"gorm.io/gorm"
)

func BootstrapPlugin(ctx context.Context, cfg *config.Config) (*gorm.DB, error) {
	// 初始化日志
	logger.Init(cfg.LogLevel)
	logger.Info("Starting PowerX Note Plugin...")

	// 初始化 schema
	models.InitSchemaFrom(cfg.Database.Schema)

	// 连接数据库
	queryDB, err := db.Connect(cfg.Database)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.WithError(err).Error("Failed to close database connection")
		}
	}()

	return queryDB, nil
}
