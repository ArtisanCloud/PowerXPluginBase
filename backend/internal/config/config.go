package config
package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Config 插件配置结构
type Config struct {
	// 服务配置
	BindAddr string `json:"bind_addr"`
	LogLevel string `json:"log_level"`
	DevMode  bool   `json:"dev_mode"`

	// 数据库配置
	DBDSN    string `json:"db_dsn"`
	DBSchema string `json:"db_schema"`

	// 运行模式
	RunMigrate bool `json:"run_migrate"`

	// PowerX 上下文配置
	Context ContextConfig `json:"context"`
}

// ContextConfig PowerX 上下文相关配置
type ContextConfig struct {
	// HMAC 模式配置
	HMACSecret string `json:"hmac_secret"`
	KeyID      string `json:"key_id"`

	// JWT 模式配置
	JWKSURL  string        `json:"jwks_url"`
	Issuer   string        `json:"issuer"`
	Audience string        `json:"audience"`
	TTL      time.Duration `json:"ttl"`
}

// Load 加载配置
func Load() (*Config, error) {
	// 尝试加载 .env 文件（开发环境）
	if err := godotenv.Load(); err != nil {
		// 生产环境可能没有 .env 文件，这是正常的
		logrus.Debug("No .env file found, using environment variables")
	}

	cfg := &Config{
		// 默认值
		BindAddr: ":8091",
		LogLevel: "info",
		DevMode:  false,
		DBSchema: "scrum",
		Context: ContextConfig{
			TTL: 300 * time.Second, // 5分钟
		},
	}

	// 从环境变量读取配置
	if addr := os.Getenv("PX_BIND_ADDR"); addr != "" {
		cfg.BindAddr = addr
	}

	if level := os.Getenv("PX_LOG_LEVEL"); level != "" {
		cfg.LogLevel = level
	}

	if dsn := os.Getenv("PX_DB_DSN"); dsn != "" {
		cfg.DBDSN = dsn
	}

	if schema := os.Getenv("PX_DB_SCHEMA"); schema != "" {
		cfg.DBSchema = schema
	}

	// 布尔型配置
	if devMode := os.Getenv("PX_DEV_MODE"); devMode == "1" || devMode == "true" {
		cfg.DevMode = true
	}

	if runMigrate := os.Getenv("PX_RUN_MIGRATE"); runMigrate == "true" {
		cfg.RunMigrate = true
	}

	// 上下文配置
	cfg.Context.HMACSecret = os.Getenv("PLUGIN_CTX_HMAC_SECRET")
	cfg.Context.KeyID = os.Getenv("PLUGIN_CTX_KID")
	cfg.Context.JWKSURL = os.Getenv("PX_CTX_JWKS_URL")
	cfg.Context.Issuer = os.Getenv("PX_CTX_ISSUER")
	cfg.Context.Audience = os.Getenv("PX_CTX_AUDIENCE")

	if ttlStr := os.Getenv("PX_CTX_TTL"); ttlStr != "" {
		if ttl, err := time.ParseDuration(ttlStr); err == nil {
			cfg.Context.TTL = ttl
		}
	}

	return cfg, nil
}

// GetString 获取字符串配置，支持默认值
func GetString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetInt 获取整数配置，支持默认值
func GetInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// GetBool 获取布尔配置，支持默认值
func GetBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "true" || value == "1" {
		return true
	}
	if value == "false" || value == "0" {
		return false
	}
	return defaultValue
}

// GetDuration 获取时间间隔配置，支持默认值
func GetDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// IsProduction 判断是否为生产环境
func (c *Config) IsProduction() bool {
	return !c.DevMode
}

// IsHMACMode 判断是否使用 HMAC 模式
func (c *Config) IsHMACMode() bool {
	return c.Context.HMACSecret != ""
}

// IsJWTMode 判断是否使用 JWT 模式
func (c *Config) IsJWTMode() bool {
	return c.Context.JWKSURL != ""
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.DBDSN == "" {
		return NewConfigError("PX_DB_DSN is required")
	}

	if !c.DevMode && !c.IsHMACMode() && !c.IsJWTMode() {
		return NewConfigError("Either HMAC or JWT mode must be configured in production")
	}

	return nil
}

// ConfigError 配置错误类型
type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return "config error: " + e.Message
}

func NewConfigError(message string) *ConfigError {
	return &ConfigError{Message: message}
}