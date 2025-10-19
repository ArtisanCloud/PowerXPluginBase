package config

import "time"

// IntegrationConfig 聚合协议适配、Webhook 与 Secrets 管理相关配置。
type IntegrationConfig struct {
	Idempotency IntegrationIdempotencyConfig `yaml:"idempotency" json:"idempotency"`
	Envelope    IntegrationEnvelopeConfig    `yaml:"envelope" json:"envelope"`
	Webhook     IntegrationWebhookConfig     `yaml:"webhook" json:"webhook"`
	Secrets     IntegrationSecretsConfig     `yaml:"secrets" json:"secrets"`
}

// IntegrationIdempotencyConfig 描述幂等后端配置。
type IntegrationIdempotencyConfig struct {
	Provider string `yaml:"provider" json:"provider"`
	RedisURL string `yaml:"redis_url" json:"redis_url"`
	TTLHours int    `yaml:"ttl_hours" json:"ttl_hours"`
}

// IntegrationEnvelopeConfig 描述 Envelope 负载约束。
type IntegrationEnvelopeConfig struct {
	PayloadThresholdBytes int64 `yaml:"payload_threshold_bytes" json:"payload_threshold_bytes"`
}

// IntegrationWebhookConfig 描述 Webhook 重试策略。
type IntegrationWebhookConfig struct {
	RetryPolicy []int  `yaml:"retry_policy" json:"retry_policy"`
	DLQTopic    string `yaml:"dlq_topic" json:"dlq_topic"`
}

// IntegrationSecretsConfig 描述外部凭证默认轮换策略。
type IntegrationSecretsConfig struct {
	RotationDaysDefault int `yaml:"rotation_days_default" json:"rotation_days_default"`
}

// IntegrationPayloadThreshold 返回 Envelope payload 阈值（字节）。
func (cfg *Config) IntegrationPayloadThreshold() int64 {
	const defaultThreshold = int64(1 << 20) // 1 MB
	if cfg == nil || cfg.Integration == nil {
		return defaultThreshold
	}
	if v := cfg.Integration.Envelope.PayloadThresholdBytes; v > 0 {
		return v
	}
	return defaultThreshold
}

// IntegrationIdempotencyTTL 返回幂等记录存活时间。
func (cfg *Config) IntegrationIdempotencyTTL() time.Duration {
	const defaultTTL = 24 * time.Hour
	if cfg == nil || cfg.Integration == nil {
		return defaultTTL
	}
	if hours := cfg.Integration.Idempotency.TTLHours; hours > 0 {
		return time.Duration(hours) * time.Hour
	}
	return defaultTTL
}

// IntegrationIdempotencyProvider 返回幂等后端类型与 Redis 地址。
func (cfg *Config) IntegrationIdempotencyProvider() (provider string, redisURL string) {
	if cfg == nil || cfg.Integration == nil {
		return "redis", "redis://localhost:6379"
	}
	prov := cfg.Integration.Idempotency.Provider
	if prov == "" {
		prov = "redis"
	}
	redisURL = cfg.Integration.Idempotency.RedisURL
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}
	return prov, redisURL
}

// IntegrationWebhookPolicy 返回 Webhook 重试策略与 DLQ Topic。
func (cfg *Config) IntegrationWebhookPolicy() ([]int, string) {
	defaultPolicy := []int{60, 300, 900}
	dlq := "plugin.webhook.dlq"
	if cfg == nil || cfg.Integration == nil {
		return defaultPolicy, dlq
	}
	policy := cfg.Integration.Webhook.RetryPolicy
	if len(policy) == 0 {
		policy = defaultPolicy
	}
	if topic := cfg.Integration.Webhook.DLQTopic; topic != "" {
		dlq = topic
	}
	return policy, dlq
}

// IntegrationSecretRotationDays 返回默认 Secrets 轮换天数。
func (cfg *Config) IntegrationSecretRotationDays() int {
	const defaultDays = 30
	if cfg == nil || cfg.Integration == nil {
		return defaultDays
	}
	if days := cfg.Integration.Secrets.RotationDaysDefault; days > 0 {
		return days
	}
	return defaultDays
}
