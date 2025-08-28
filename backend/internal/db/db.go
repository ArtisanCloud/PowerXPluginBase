package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"scrum-plugin/internal/config"
	"scrum-plugin/internal/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// db 私有数据库实例
var db *gorm.DB

// GetGlobalDB 获取全局数据库实例
func GetGlobalDB() *gorm.DB {
	return db
}

// TenantDB 租户作用域的数据库实例
type TenantDB struct {
	*gorm.DB
	TenantID int64
	ctx      context.Context
}

// Connect 连接数据库
func Connect(cfg *config.Config) error {
	var err error

	// 配置 GORM 日志 - 使用 Silent 模式以减少数据库日志输出
	// 只有在明确设置 debug 模式时才显示数据库详细日志
	gormLogLevel := gormLogger.Silent
	if cfg.Server.DevMode && cfg.LogLevel == "debug" {
		gormLogLevel = gormLogger.Info
	}

	gormConfig := &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// 连接数据库
	db, err = gorm.Open(postgres.Open(cfg.DBDSN), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 获取底层 SQL DB 实例
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB instance: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// 创建 schema（如果不存在）
	if err := createSchema(cfg.DBSchema); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	// 设置默认 schema
	if err := setDefaultSchema(cfg.DBSchema); err != nil {
		return fmt.Errorf("failed to set default schema: %w", err)
	}

	logger.Infof("Database connected successfully with schema: %s", cfg.DBSchema)
	return nil
}

// createSchema 创建 schema
func createSchema(schema string) error {
	sql := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema)
	return db.Exec(sql).Error
}

// setDefaultSchema 设置默认 schema
func setDefaultSchema(schema string) error {
	sql := fmt.Sprintf("SET search_path TO %s", schema)
	return db.Exec(sql).Error
}

// BeginTenantTx 开始租户作用域的事务
func BeginTenantTx(ctx context.Context, tenantID int64) (*TenantDB, error) {
	// 开始事务
	tx := db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// 设置租户 ID 到会话变量
	if err := setTenantID(tx, tenantID); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to set tenant id: %w", err)
	}

	return &TenantDB{
		DB:       tx,
		TenantID: tenantID,
		ctx:      ctx,
	}, nil
}

// setTenantID 设置租户 ID 到会话变量
func setTenantID(db *gorm.DB, tenantID int64) error {
	sql := fmt.Sprintf("SET LOCAL app.tenant_id = %d", tenantID)
	return db.Exec(sql).Error
}

// Commit 提交事务
func (tdb *TenantDB) Commit() error {
	return tdb.DB.Commit().Error
}

// Rollback 回滚事务
func (tdb *TenantDB) Rollback() error {
	return tdb.DB.Rollback().Error
}

// WithContext 设置上下文
func (tdb *TenantDB) WithContext(ctx context.Context) *TenantDB {
	return &TenantDB{
		DB:       tdb.DB.WithContext(ctx),
		TenantID: tdb.TenantID,
		ctx:      ctx,
	}
}

// GetTenantID 获取当前租户 ID
func (tdb *TenantDB) GetTenantID() int64 {
	return tdb.TenantID
}

// Close 关闭数据库连接
func Close() error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// Health 检查数据库健康状态
func Health() error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx)
}

// Migrate 执行数据库迁移
func Migrate(models ...interface{}) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	return db.AutoMigrate(models...)
}

// EnableRLS 为表启用行级安全
func EnableRLS(tableName string) error {
	sql := fmt.Sprintf("ALTER TABLE %s ENABLE ROW LEVEL SECURITY", tableName)
	return db.Exec(sql).Error
}

// CreateRLSPolicy 创建 RLS 策略
func CreateRLSPolicy(tableName, policyName string) error {
	sql := fmt.Sprintf(`
		CREATE POLICY IF NOT EXISTS %s ON %s
		USING (tenant_id::text = current_setting('app.tenant_id', true))
	`, policyName, tableName)
	return db.Exec(sql).Error
}

// DropRLSPolicy 删除 RLS 策略
func DropRLSPolicy(tableName, policyName string) error {
	sql := fmt.Sprintf("DROP POLICY IF EXISTS %s ON %s", policyName, tableName)
	return db.Exec(sql).Error
}

// WithTenant 在开发模式下模拟租户上下文
func WithTenant(tenantID int64) *gorm.DB {
	return db.Session(&gorm.Session{}).Where("tenant_id = ?", tenantID)
}

// GetCurrentTenantID 从数据库会话中获取当前租户 ID
func GetCurrentTenantID(db *gorm.DB) (int64, error) {
	var tenantID sql.NullString
	err := db.Raw("SELECT current_setting('app.tenant_id', true)").Scan(&tenantID).Error
	if err != nil {
		return 0, err
	}

	if !tenantID.Valid || tenantID.String == "" {
		return 0, fmt.Errorf("tenant_id not set in session")
	}

	var result int64
	_, err = fmt.Sscanf(tenantID.String, "%d", &result)
	if err != nil {
		return 0, fmt.Errorf("invalid tenant_id format: %s", tenantID.String)
	}

	return result, nil
}

// ExecInTenant 在指定租户上下文中执行操作
func ExecInTenant(ctx context.Context, tenantID int64, fn func(*TenantDB) error) error {
	tdb, err := BeginTenantTx(ctx, tenantID)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tdb.Rollback()
			panic(r)
		}
	}()

	if err := fn(tdb); err != nil {
		tdb.Rollback()
		return err
	}

	return tdb.Commit()
}
