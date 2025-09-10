package admin

import (
	"time"

	"github.com/ArtisanCloud/PowerXPlugin/internal/contracts"
	"github.com/ArtisanCloud/PowerXPlugin/internal/logger"
	"github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"

	"github.com/gin-gonic/gin"
)

// AdminHandler 管理端处理器
type AdminHandler struct {
	deps *app.Deps
}

// NewAdminHandler 创建管理端处理器
func NewAdminHandler(deps *app.Deps) *AdminHandler {
	return &AdminHandler{
		deps: deps,
	}
}

// GetManifest 获取插件清单
func (h *AdminHandler) GetManifest(c *gin.Context) {
	log := logger.HandlerLogger("admin").WithContext(c.Request.Context())

	manifest := &contracts.PluginManifest{
		ID:          "com.powerx.plugins.base",
		Name:        "Base Note Plugin",
		Version:     "0.1.0",
		Description: "A comprehensive Base note management plugin for PowerX",
		Author:      "PowerX Team",
		Homepage:    "https://base-plugin",
		Repository:  "https://base-plugin.git",
		License:     "MIT",
		Tags:        []string{"base", "agile", "note-management", "project-management"},

		Backend: contracts.BackendConfig{
			Entry:  "backend/bin/plugin",
			Port:   8091,
			Health: "/healthz",
		},

		Frontend: &contracts.FrontendConfig{
			Entry: "web-admin/.output",
			Routes: map[string]string{
				"/admin/*": "index.html",
			},
			PublicPath: "/_p/com.powerx.plugins.base/admin/",
		},

		Menus: []contracts.MenuConfig{
			{
				ID:    "base",
				Title: "menu.base",
				Icon:  "i-heroicons-clipboard-document-check",
				Path:  "/plugins/base",
				Order: 20,
				Children: []contracts.MenuConfig{
					{
						ID:    "base.dashboard",
						Title: "menu.base.dashboard",
						Icon:  "i-heroicons-squares-2x2",
						Path:  "/plugins/base/dashboard",
						Order: 1,
					},
					{
						ID:    "base.notes",
						Title: "menu.base.notes",
						Icon:  "i-heroicons-list-bullet",
						Path:  "/plugins/base/notes",
						Order: 2,
					},
				},
				RequiredPermissions: []string{"base:note:read"},
			},
		},

		Permissions: []contracts.PermissionConfig{
			{
				Resource:    "base:note",
				Actions:     []string{"read", "create", "update", "delete"},
				Description: "Note management permissions",
			},
		},

		Agents: []contracts.AgentConfig{
			{
				ID:           "base.assistant",
				PluginID:     "com.powerx.plugins.base",
				Name:         "Base 助理",
				Description:  "智能的 Base 管理助手，可以帮助创建任务、管理 Sprint、生成报告",
				Model:        "gpt-4",
				Instructions: "你是一个专业的 Base 管理助手。你可以帮助用户创建和管理任务、规划 Sprint、跟踪进度并生成报告。请始终以友好、专业的方式回应用户的请求。",
				DefaultTools: []string{
					"base.note.create",
					"base.note.query",
					"base.note.update",
					"base.report.generate",
				},
				RequiredPermissions: []string{"base:note:read"},
			},
		},

		Tools: []contracts.ToolConfig{
			{
				ID:           "base.note.create",
				PluginID:     "com.powerx.plugins.base",
				Name:         "创建任务",
				Description:  "创建一个新的 Base 任务",
				Transport:    "http",
				Endpoint:     "/api/v1/notes",
				Method:       "POST",
				RBACResource: "base:note",
				InputSchema: &contracts.JSONSchema{
					Type: "object",
					Properties: map[string]*contracts.JSONSchemaProperty{
						"title": {
							Type:        "string",
							Description: "任务标题",
						},
						"description": {
							Type:        "string",
							Description: "任务描述",
						},
						"priority": {
							Type:        "string",
							Enum:        []interface{}{"low", "medium", "high", "urgent"},
							Description: "任务优先级",
							Default:     "medium",
						},
						"assignee": {
							Type:        "integer",
							Description: "分配给的用户ID",
						},
						"estimate": {
							Type:        "integer",
							Minimum:     func() *float64 { v := 1.0; return &v }(),
							Maximum:     func() *float64 { v := 100.0; return &v }(),
							Description: "故事点估算",
						},
					},
					Required: []string{"title"},
				},
				OutputSchema: &contracts.JSONSchema{
					Type: "object",
					Properties: map[string]*contracts.JSONSchemaProperty{
						"id": {
							Type:        "integer",
							Description: "任务ID",
						},
						"title": {
							Type:        "string",
							Description: "任务标题",
						},
						"status": {
							Type:        "string",
							Description: "任务状态",
						},
					},
				},
				Timeout: 30,
			},
			{
				ID:           "base.note.query",
				PluginID:     "com.powerx.plugins.base",
				Name:         "查询任务",
				Description:  "查询 Base 任务列表",
				Transport:    "http",
				Endpoint:     "/api/v1/notes",
				Method:       "GET",
				RBACResource: "base:note",
				InputSchema: &contracts.JSONSchema{
					Type: "object",
					Properties: map[string]*contracts.JSONSchemaProperty{
						"status": {
							Type:        "string",
							Enum:        []interface{}{"todo", "in_progress", "done"},
							Description: "按状态过滤",
						},
						"assignee": {
							Type:        "integer",
							Description: "按分配人过滤",
						},
						"page": {
							Type:        "integer",
							Minimum:     func() *float64 { v := 1.0; return &v }(),
							Default:     1,
							Description: "页码",
						},
						"limit": {
							Type:        "integer",
							Minimum:     func() *float64 { v := 1.0; return &v }(),
							Maximum:     func() *float64 { v := 100.0; return &v }(),
							Default:     20,
							Description: "每页数量",
						},
					},
				},
				Timeout: 30,
			},
			// Sprint 相关遗留已移除
		},

		// Workflows 中与 Sprint 相关的配置已移除

		Dependencies: []contracts.DependencyConfig{
			{
				Name:    "postgresql",
				Version: ">=12.0",
				Type:    "database",
			},
		},

		ConfigSchema: &contracts.ConfigSchema{
			Type: "object",
			Properties: map[string]*contracts.JSONSchemaProperty{
				"note_auto_close": {
					Type:        "boolean",
					Description: "是否自动关闭过期任务",
					Default:     false,
				},
				"email_notifications": {
					Type:        "boolean",
					Description: "是否启用邮件通知",
					Default:     true,
				},
			},
		},

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	log.Info("Plugin manifest requested")

	contracts.ResponseSuccess(c, manifest)
}

// GetRBACInfo 获取 RBAC 信息
func (h *AdminHandler) GetRBACInfo(c *gin.Context) {
	log := logger.HandlerLogger("admin").WithContext(c.Request.Context())

	rbacInfo := &contracts.RBACInfo{
		Resources: []contracts.Resource{
			{
				Name:        "base:note",
				Description: "Base 任务管理",
				Actions: []contracts.Action{
					{Name: "read", Description: "查看任务"},
					{Name: "create", Description: "创建任务"},
					{Name: "update", Description: "更新任务"},
					{Name: "delete", Description: "删除任务"},
				},
			},
			{
				Name:        "base:report",
				Description: "Base 报告",
				Actions: []contracts.Action{
					{Name: "read", Description: "查看报告"},
				},
			},
		},
		Roles: []contracts.Role{
			{
				Name:        "base_master",
				Description: "Base Master 角色",
				Permissions: []string{
					"base:note:*",
					"base:report:read",
				},
			},
			{
				Name:        "product_owner",
				Description: "Product Owner 角色",
				Permissions: []string{
					"base:note:read",
					"base:note:create",
					"base:note:update",
					"base:report:read",
				},
			},
			{
				Name:        "developer",
				Description: "开发者角色",
				Permissions: []string{
					"base:note:read",
					"base:note:update",
				},
			},
		},
		Permissions: []contracts.Permission{
			{Resource: "base:note", Action: "read"},
			{Resource: "base:note", Action: "create"},
			{Resource: "base:note", Action: "update"},
			{Resource: "base:note", Action: "delete"},
			{Resource: "base:report", Action: "read"},
		},
	}

	log.Info("RBAC info requested")

	contracts.ResponseSuccess(c, rbacInfo)
}
