package config

import (
	"strings"
	"time"
)

// MarketplaceConfig 聚合 Marketplace 模块相关配置。
type MarketplaceConfig struct {
	Checklist      MarketplaceChecklistConfig `yaml:"checklist" json:"checklist"`
	Documentation  MarketplaceDocsConfig      `yaml:"documentation" json:"documentation"`
	Recommendation MarketplaceRecommendation  `yaml:"recommendation" json:"recommendation"`
}

// MarketplaceChecklistConfig 描述 Ready Checklist GraphQL 入口与文档。
type MarketplaceChecklistConfig struct {
	GraphQLPath string `yaml:"graphql_path" json:"graphql_path"`
	DocsURL     string `yaml:"docs_url" json:"docs_url"`
}

// MarketplaceDocsConfig 引用 Marketplace 相关说明文档。
type MarketplaceDocsConfig struct {
	Overview        string `yaml:"overview" json:"overview"`
	ReadyChecklist  string `yaml:"ready_checklist" json:"ready_checklist"`
	ListingPlaybook string `yaml:"listing_playbook" json:"listing_playbook"`
}

// MarketplaceRecommendation 描述推荐服务的默认配置。
type MarketplaceRecommendation struct {
	Enabled          bool    `yaml:"enabled" json:"enabled"`
	DefaultWeight    float64 `yaml:"default_weight" json:"default_weight"`
	ExperimentTopic  string  `yaml:"experiment_topic" json:"experiment_topic"`
	FrequencyMinutes int     `yaml:"frequency_minutes" json:"frequency_minutes"`
}

// MarketplaceChecklistGraphQLPath 返回 checklist GraphQL 路径。
func (cfg *Config) MarketplaceChecklistGraphQLPath() string {
	const defaultPath = "/api/v1/admin/marketplace/checklist/graphql"
	if cfg == nil || cfg.Marketplace == nil {
		return defaultPath
	}
	path := strings.TrimSpace(cfg.Marketplace.Checklist.GraphQLPath)
	if path == "" {
		return defaultPath
	}
	return path
}

// MarketplaceChecklistDocs 提供 checklist 文档地址。
func (cfg *Config) MarketplaceChecklistDocs() string {
	if cfg == nil || cfg.Marketplace == nil {
		return ""
	}
	return strings.TrimSpace(cfg.Marketplace.Checklist.DocsURL)
}

// RecommendationFrequency returns refresh frequency for recommendation sync job.
func (cfg *Config) RecommendationFrequency() time.Duration {
	if cfg == nil || cfg.Marketplace == nil {
		return time.Hour
	}
	minutes := cfg.Marketplace.Recommendation.FrequencyMinutes
	if minutes <= 0 {
		return time.Hour
	}
	return time.Duration(minutes) * time.Minute
}
