package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Config 插件配置结构
type Config struct {
	// 服务配置
	Server ServerConfig `yaml:"server" json:"server"`

	// 数据库配置
	Database DatabaseConfig `yaml:"database" json:"database"`

	// 运行时配置
	Runtime RuntimeConfig `yaml:"runtime" json:"runtime"`

	// PowerX 上下文配置
	Context ContextConfig `yaml:"context" json:"context"`

	// 业务配置
	Business BusinessConfig `yaml:"business" json:"business"`

	// 安全配置
	Security SecurityConfig `yaml:"security" json:"security"`

	// 监控配置
	Monitoring MonitoringConfig `yaml:"monitoring" json:"monitoring"`

	// 日志配置
	Logging LoggingConfig `yaml:"logging" json:"logging"`

	// 向后兼容的字段（从环境变量或旧配置中填充）
	BindAddr   string `yaml:"-" json:"bind_addr,omitempty"`
	LogLevel   string `yaml:"-" json:"log_level,omitempty"`
	DevMode    bool   `yaml:"-" json:"dev_mode,omitempty"`
	DBDSN      string `yaml:"-" json:"db_dsn,omitempty"`
	DBSchema   string `yaml:"-" json:"db_schema,omitempty"`
	RunMigrate bool   `yaml:"-" json:"run_migrate,omitempty"`
}

// ServerConfig 服务配置
type ServerConfig struct {
	BindAddr string `yaml:"bind_addr" json:"bind_addr"`
	LogLevel string `yaml:"log_level" json:"log_level"`
	DevMode  bool   `yaml:"dev_mode" json:"dev_mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	DSN    string `yaml:"dsn" json:"dsn"`
	Schema string `yaml:"schema" json:"schema"`
}

// RuntimeConfig 运行时配置
type RuntimeConfig struct {
	RunMigrate bool `yaml:"run_migrate" json:"run_migrate"`
}

// BusinessConfig 业务配置
type BusinessConfig struct {
	Sprint        SprintConfig        `yaml:"sprint" json:"sprint"`
	Task          TaskConfig          `yaml:"task" json:"task"`
	Notifications NotificationsConfig `yaml:"notifications" json:"notifications"`
	Cache         CacheConfig         `yaml:"cache" json:"cache"`
}

// SprintConfig Sprint 相关配置
type SprintConfig struct {
	DefaultCapacity int `yaml:"default_capacity" json:"default_capacity"`
	MaxDurationDays int `yaml:"max_duration_days" json:"max_duration_days"`
}

// TaskConfig 任务相关配置
type TaskConfig struct {
	MaxEstimatePoints int    `yaml:"max_estimate_points" json:"max_estimate_points"`
	DefaultPriority   string `yaml:"default_priority" json:"default_priority"`
}

// NotificationsConfig 通知配置
type NotificationsConfig struct {
	Enabled bool        `yaml:"enabled" json:"enabled"`
	Email   EmailConfig `yaml:"email" json:"email"`
	Slack   SlackConfig `yaml:"slack" json:"slack"`
}

// EmailConfig 邮件配置
type EmailConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	SMTPHost string `yaml:"smtp_host" json:"smtp_host"`
	SMTPPort int    `yaml:"smtp_port" json:"smtp_port"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	From     string `yaml:"from" json:"from"`
}

// SlackConfig Slack 配置
type SlackConfig struct {
	Enabled    bool   `yaml:"enabled" json:"enabled"`
	WebhookURL string `yaml:"webhook_url" json:"webhook_url"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Enabled  bool          `yaml:"enabled" json:"enabled"`
	RedisURL string        `yaml:"redis_url" json:"redis_url"`
	TTL      time.Duration `yaml:"ttl" json:"ttl"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	EnableCORS  bool            `yaml:"enable_cors" json:"enable_cors"`
	CORSOrigins []string        `yaml:"cors_origins" json:"cors_origins"`
	RateLimit   RateLimitConfig `yaml:"rate_limit" json:"rate_limit"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled           bool `yaml:"enabled" json:"enabled"`
	RequestsPerMinute int  `yaml:"requests_per_minute" json:"requests_per_minute"`
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	Metrics     MetricsConfig     `yaml:"metrics" json:"metrics"`
	HealthCheck HealthCheckConfig `yaml:"health_check" json:"health_check"`
}

// MetricsConfig 指标配置
type MetricsConfig struct {
	Enabled bool   `yaml:"enabled" json:"enabled"`
	Path    string `yaml:"path" json:"path"`
}

// HealthCheckConfig 健康检查配置
type HealthCheckConfig struct {
	Enabled bool   `yaml:"enabled" json:"enabled"`
	Path    string `yaml:"path" json:"path"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level      string `yaml:"level" json:"level"`
	Format     string `yaml:"format" json:"format"`
	Output     string `yaml:"output" json:"output"`
	FilePath   string `yaml:"file_path" json:"file_path"`
	MaxSize    int    `yaml:"max_size" json:"max_size"`
	MaxBackups int    `yaml:"max_backups" json:"max_backups"`
	MaxAge     int    `yaml:"max_age" json:"max_age"`
}

// ContextConfig PowerX 上下文相关配置
type ContextConfig struct {
	// HMAC 模式配置
	HMACSecret string `yaml:"hmac_secret" json:"hmac_secret"`
	KeyID      string `yaml:"key_id" json:"key_id"`

	// JWT 模式配置
	JWKSURL  string        `yaml:"jwks_url" json:"jwks_url"`
	Issuer   string        `yaml:"issuer" json:"issuer"`
	Audience string        `yaml:"audience" json:"audience"`
	TTL      time.Duration `yaml:"ttl" json:"ttl"`
}

// Load 加载配置，优先级：YAML 文件 > 环境变量 > 默认值
func Load() (*Config, error) {
	// 尝试加载 .env 文件（开发环境）
	if err := godotenv.Load(); err != nil {
		// 生产环境可能没有 .env 文件，这是正常的
		logrus.Debug("No .env file found, using environment variables")
	}

	// 设置默认配置
	cfg := getDefaultConfig()

	// 尝试加载 YAML 配置文件
	if err := loadYAMLConfig(cfg); err != nil {
		logrus.WithError(err).Warn("Failed to load YAML config, using defaults and environment variables")
	}

	// 从环境变量覆盖配置（最高优先级）
	loadEnvConfig(cfg)

	// 同步向后兼容字段
	syncBackwardCompatibility(cfg)

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// getDefaultConfig 获取默认配置
func getDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			BindAddr: ":8086",
			LogLevel: "info",
			DevMode:  false,
		},
		Database: DatabaseConfig{
			Schema: "scrum",
		},
		Runtime: RuntimeConfig{
			RunMigrate: false,
		},
		Context: ContextConfig{
			TTL: 300 * time.Second, // 5分钟
		},
		Business: BusinessConfig{
			Sprint: SprintConfig{
				DefaultCapacity: 40,
				MaxDurationDays: 30,
			},
			Task: TaskConfig{
				MaxEstimatePoints: 100,
				DefaultPriority:   "medium",
			},
			Notifications: NotificationsConfig{
				Enabled: false,
				Email: EmailConfig{
					SMTPPort: 587,
				},
			},
			Cache: CacheConfig{
				Enabled: false,
				TTL:     time.Hour,
			},
		},
		Security: SecurityConfig{
			EnableCORS: true,
			CORSOrigins: []string{
				"http://localhost:3036",
				"http://localhost:3000",
			},
			RateLimit: RateLimitConfig{
				Enabled:           true,
				RequestsPerMinute: 60,
			},
		},
		Monitoring: MonitoringConfig{
			Metrics: MetricsConfig{
				Enabled: false,
				Path:    "/metrics",
			},
			HealthCheck: HealthCheckConfig{
				Enabled: true,
				Path:    "/health",
			},
		},
		Logging: LoggingConfig{
			Level:      "info",
			Format:     "json",
			Output:     "stdout",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
		},
	}
}

// loadYAMLConfig 加载 YAML 配置文件
func loadYAMLConfig(cfg *Config) error {
	// 查找配置文件
	configPaths := []string{
		"./etc/config.yaml",
		"./config.yaml",
		"../etc/config.yaml",
		filepath.Join(os.Getenv("CONFIG_PATH"), "config.yaml"),
	}

	var configFile string
	for _, path := range configPaths {
		if path == "" {
			continue
		}
		if _, err := os.Stat(path); err == nil {
			configFile = path
			break
		}
	}

	if configFile == "" {
		return fmt.Errorf("config file not found")
	}

	// 读取文件
	file, err := os.Open(configFile)
	if err != nil {
		return fmt.Errorf("failed to open config file %s: %w", configFile, err)
	}
	defer file.Close()

	// 读取文件内容
	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析 YAML
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to parse YAML config: %w", err)
	}

	logrus.WithField("config_file", configFile).Info("YAML config loaded successfully")
	return nil
}

// loadEnvConfig 从环境变量加载配置
func loadEnvConfig(cfg *Config) {
	// 服务配置
	if addr := os.Getenv("PX_BIND_ADDR"); addr != "" {
		cfg.Server.BindAddr = addr
	}
	if level := os.Getenv("PX_LOG_LEVEL"); level != "" {
		cfg.Server.LogLevel = level
		cfg.Logging.Level = level
	}
	if devMode := os.Getenv("PX_DEV_MODE"); devMode != "" {
		cfg.Server.DevMode = (devMode == "1" || devMode == "true")
	}

	// 数据库配置
	if dsn := os.Getenv("PX_DB_DSN"); dsn != "" {
		cfg.Database.DSN = dsn
	}
	if schema := os.Getenv("PX_DB_SCHEMA"); schema != "" {
		cfg.Database.Schema = schema
	}

	// 运行时配置
	if runMigrate := os.Getenv("PX_RUN_MIGRATE"); runMigrate == "true" {
		cfg.Runtime.RunMigrate = true
	}

	// 上下文配置
	if hmacSecret := os.Getenv("PLUGIN_CTX_HMAC_SECRET"); hmacSecret != "" {
		cfg.Context.HMACSecret = hmacSecret
	}
	if keyID := os.Getenv("PLUGIN_CTX_KID"); keyID != "" {
		cfg.Context.KeyID = keyID
	}
	if jwksURL := os.Getenv("PX_CTX_JWKS_URL"); jwksURL != "" {
		cfg.Context.JWKSURL = jwksURL
	}
	if issuer := os.Getenv("PX_CTX_ISSUER"); issuer != "" {
		cfg.Context.Issuer = issuer
	}
	if audience := os.Getenv("PX_CTX_AUDIENCE"); audience != "" {
		cfg.Context.Audience = audience
	}
	if ttlStr := os.Getenv("PX_CTX_TTL"); ttlStr != "" {
		if ttl, err := time.ParseDuration(ttlStr); err == nil {
			cfg.Context.TTL = ttl
		}
	}
}

// syncBackwardCompatibility 同步向后兼容字段
func syncBackwardCompatibility(cfg *Config) {
	cfg.BindAddr = cfg.Server.BindAddr
	cfg.LogLevel = cfg.Server.LogLevel
	cfg.DevMode = cfg.Server.DevMode
	cfg.DBDSN = cfg.Database.DSN
	cfg.DBSchema = cfg.Database.Schema
	cfg.RunMigrate = cfg.Runtime.RunMigrate
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
	return !c.Server.DevMode && !c.DevMode
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
	// 数据库配置验证
	if c.Database.DSN == "" && c.DBDSN == "" {
		return NewConfigError("database DSN is required (set PX_DB_DSN or configure in YAML)")
	}

	// 认证模式验证
	if !c.Server.DevMode && !c.IsHMACMode() && !c.IsJWTMode() {
		return NewConfigError("either HMAC or JWT mode must be configured in production")
	}

	// 业务配置验证
	if c.Business.Sprint.DefaultCapacity <= 0 {
		return NewConfigError("sprint default capacity must be positive")
	}

	if c.Business.Sprint.MaxDurationDays <= 0 || c.Business.Sprint.MaxDurationDays > 365 {
		return NewConfigError("sprint max duration days must be between 1 and 365")
	}

	if c.Business.Task.MaxEstimatePoints <= 0 {
		return NewConfigError("task max estimate points must be positive")
	}

	validPriorities := map[string]bool{
		"low": true, "medium": true, "high": true, "urgent": true,
	}
	if !validPriorities[c.Business.Task.DefaultPriority] {
		return NewConfigError("task default priority must be one of: low, medium, high, urgent")
	}

	// 安全配置验证
	if c.Security.RateLimit.Enabled && c.Security.RateLimit.RequestsPerMinute <= 0 {
		return NewConfigError("rate limit requests per minute must be positive when enabled")
	}

	// 日志配置验证
	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLogLevels[c.Logging.Level] {
		return NewConfigError("logging level must be one of: debug, info, warn, error")
	}

	validLogFormats := map[string]bool{
		"json": true, "text": true,
	}
	if !validLogFormats[c.Logging.Format] {
		return NewConfigError("logging format must be one of: json, text")
	}

	validLogOutputs := map[string]bool{
		"stdout": true, "stderr": true, "file": true,
	}
	if !validLogOutputs[c.Logging.Output] {
		return NewConfigError("logging output must be one of: stdout, stderr, file")
	}

	if c.Logging.Output == "file" && c.Logging.FilePath == "" {
		return NewConfigError("logging file path must be specified when output is 'file'")
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
