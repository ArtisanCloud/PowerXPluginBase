package main

import (
	"fmt"
	"os"

	"github.com/powerx-plugins/scrum/internal/config"
	"github.com/powerx-plugins/scrum/internal/db"
	"github.com/powerx-plugins/scrum/internal/domain"
	"github.com/powerx-plugins/scrum/internal/logger"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	logger.Init(cfg.LogLevel)
	logger.Info("Starting database migration...")

	// 连接数据库
	if err := db.Connect(cfg); err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.WithError(err).Error("Failed to close database connection")
		}
	}()

	// 检查是否重置数据库
	if os.Getenv("PX_RESET_DB") == "true" {
		logger.Warn("Resetting database...")
		if err := resetDatabase(cfg); err != nil {
			logger.WithError(err).Fatal("Failed to reset database")
		}
		logger.Info("Database reset completed")
	}

	// 运行迁移
	if err := runMigrations(cfg); err != nil {
		logger.WithError(err).Fatal("Failed to run migrations")
	}

	logger.Info("Database migration completed successfully")
}

// resetDatabase 重置数据库（危险操作）
func resetDatabase(cfg *config.Config) error {
	logger.Warn("Dropping all tables in schema...")

	// 删除表（注意顺序，先删除有外键的表）
	tables := []string{
		"task",
		"sprint",
	}

	for _, table := range tables {
		sql := fmt.Sprintf("DROP TABLE IF EXISTS %s.%s CASCADE", cfg.DBSchema, table)
		if err := db.DB.Exec(sql).Error; err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
		logger.WithField("table", table).Info("Table dropped")
	}

	return nil
}

// runMigrations 运行数据库迁移
func runMigrations(cfg *config.Config) error {
	logger.Info("Creating schema if not exists...")
	
	// 创建 schema
	sql := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", cfg.DBSchema)
	if err := db.DB.Exec(sql).Error; err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	// 设置搜索路径
	sql = fmt.Sprintf("SET search_path TO %s", cfg.DBSchema)
	if err := db.DB.Exec(sql).Error; err != nil {
		return fmt.Errorf("failed to set search path: %w", err)
	}

	logger.Info("Running auto migrations...")

	// 定义要迁移的模型
	models := []interface{}{
		&domain.Sprint{}, // 先创建 Sprint 表，因为 Task 表引用它
		&domain.Task{},
	}

	// 执行 GORM 自动迁移
	if err := db.DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("auto migrate failed: %w", err)
	}

	logger.Info("Creating indexes...")
	if err := createIndexes(cfg); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	logger.Info("Enabling Row Level Security...")
	if err := enableRLS(cfg); err != nil {
		return fmt.Errorf("failed to enable RLS: %w", err)
	}

	logger.Info("Creating RLS policies...")
	if err := createRLSPolicies(cfg); err != nil {
		return fmt.Errorf("failed to create RLS policies: %w", err)
	}

	return nil
}

// createIndexes 创建索引
func createIndexes(cfg *config.Config) error {
	indexes := []string{
		// Task 表索引
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_task_tenant_id ON %s.task(tenant_id)", cfg.DBSchema, cfg.DBSchema),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_task_status ON %s.task(status)", cfg.DBSchema, cfg.DBSchema),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_task_priority ON %s.task(priority)", cfg.DBSchema, cfg.DBSchema),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_task_assignee ON %s.task(assignee)", cfg.DBSchema, cfg.DBSchema),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_task_sprint_id ON %s.task(sprint_id)", cfg.DBSchema, cfg.DBSchema),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_task_due_date ON %s.task(due_date)", cfg.DBSchema, cfg.DBSchema),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_task_created_at ON %s.task(created_at)", cfg.DBSchema, cfg.DBSchema),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_task_labels ON %s.task USING GIN(labels)", cfg.DBSchema, cfg.DBSchema),
		
		// Sprint 表索引
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_sprint_tenant_id ON %s.sprint(tenant_id)", cfg.DBSchema, cfg.DBSchema),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_sprint_status ON %s.sprint(status)", cfg.DBSchema, cfg.DBSchema),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_sprint_start_date ON %s.sprint(start_date)", cfg.DBSchema, cfg.DBSchema),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_sprint_end_date ON %s.sprint(end_date)", cfg.DBSchema, cfg.DBSchema),
		
		// 复合索引
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_task_tenant_status ON %s.task(tenant_id, status)", cfg.DBSchema, cfg.DBSchema),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_task_tenant_assignee ON %s.task(tenant_id, assignee)", cfg.DBSchema, cfg.DBSchema),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_sprint_tenant_status ON %s.sprint(tenant_id, status)", cfg.DBSchema, cfg.DBSchema),
	}

	for _, index := range indexes {
		if err := db.DB.Exec(index).Error; err != nil {
			logger.WithError(err).WithField("index", index).Warn("Failed to create index (may already exist)")
		}
	}

	return nil
}

// enableRLS 启用行级安全
func enableRLS(cfg *config.Config) error {
	tables := []string{"task", "sprint"}

	for _, table := range tables {
		fullTableName := fmt.Sprintf("%s.%s", cfg.DBSchema, table)
		sql := fmt.Sprintf("ALTER TABLE %s ENABLE ROW LEVEL SECURITY", fullTableName)
		
		if err := db.DB.Exec(sql).Error; err != nil {
			logger.WithError(err).WithField("table", table).Warn("Failed to enable RLS (may already be enabled)")
		} else {
			logger.WithField("table", table).Info("RLS enabled")
		}
	}

	return nil
}

// createRLSPolicies 创建 RLS 策略
func createRLSPolicies(cfg *config.Config) error {
	policies := []struct {
		table  string
		policy string
		sql    string
	}{
		{
			table:  "task",
			policy: "p_tenant_isolation",
			sql: fmt.Sprintf(`
				CREATE POLICY IF NOT EXISTS p_tenant_isolation ON %s.task
				USING (tenant_id::text = current_setting('app.tenant_id', true))
			`, cfg.DBSchema),
		},
		{
			table:  "sprint",
			policy: "p_tenant_isolation",
			sql: fmt.Sprintf(`
				CREATE POLICY IF NOT EXISTS p_tenant_isolation ON %s.sprint
				USING (tenant_id::text = current_setting('app.tenant_id', true))
			`, cfg.DBSchema),
		},
	}

	for _, policy := range policies {
		if err := db.DB.Exec(policy.sql).Error; err != nil {
			logger.WithError(err).WithFields(logger.Logger.Fields{
				"table":  policy.table,
				"policy": policy.policy,
			}).Warn("Failed to create RLS policy (may already exist)")
		} else {
			logger.WithFields(logger.Logger.Fields{
				"table":  policy.table,
				"policy": policy.policy,
			}).Info("RLS policy created")
		}
	}

	return nil
}