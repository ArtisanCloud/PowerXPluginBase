package handlers

import (
	"net/http"
	"time"

	"scrum-plugin/internal/contracts"
	"scrum-plugin/internal/logger"

	"github.com/gin-gonic/gin"
)

// AdminHandler 管理端处理器
type AdminHandler struct{}

// NewAdminHandler 创建管理端处理器
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// GetManifest 获取插件清单
func (h *AdminHandler) GetManifest(c *gin.Context) {
	log := logger.HandlerLogger("admin").WithContext(c.Request.Context())

	manifest := &contracts.PluginManifest{
		ID:          "com.powerx.plugins.scrum",
		Name:        "Scrum Task Plugin",
		Version:     "0.1.0",
		Description: "A comprehensive Scrum task management plugin for PowerX",
		Author:      "PowerX Team",
		Homepage:    "https://scrum-plugin",
		Repository:  "https://scrum-plugin.git",
		License:     "MIT",
		Tags:        []string{"scrum", "agile", "task-management", "project-management"},

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
			PublicPath: "/_p/com.powerx.plugins.scrum/admin/",
		},

		Menus: []contracts.MenuConfig{
			{
				ID:    "scrum",
				Title: "menu.scrum",
				Icon:  "i-heroicons-clipboard-document-check",
				Path:  "/plugins/scrum",
				Order: 20,
				Children: []contracts.MenuConfig{
					{
						ID:    "scrum.dashboard",
						Title: "menu.scrum.dashboard",
						Icon:  "i-heroicons-squares-2x2",
						Path:  "/plugins/scrum/dashboard",
						Order: 1,
					},
					{
						ID:    "scrum.tasks",
						Title: "menu.scrum.tasks",
						Icon:  "i-heroicons-list-bullet",
						Path:  "/plugins/scrum/tasks",
						Order: 2,
					},
					{
						ID:    "scrum.sprints",
						Title: "menu.scrum.sprints",
						Icon:  "i-heroicons-calendar-days",
						Path:  "/plugins/scrum/sprints",
						Order: 3,
					},
					{
						ID:    "scrum.reports",
						Title: "menu.scrum.reports",
						Icon:  "i-heroicons-chart-bar",
						Path:  "/plugins/scrum/reports",
						Order: 4,
					},
				},
				RequiredPermissions: []string{"scrum:task:read"},
			},
		},

		Permissions: []contracts.PermissionConfig{
			{
				Resource:    "scrum:task",
				Actions:     []string{"read", "create", "update", "delete"},
				Description: "Task management permissions",
			},
			{
				Resource:    "scrum:sprint",
				Actions:     []string{"read", "create", "update", "delete"},
				Description: "Sprint management permissions",
			},
			{
				Resource:    "scrum:report",
				Actions:     []string{"read"},
				Description: "Scrum reports access",
			},
		},

		Agents: []contracts.AgentConfig{
			{
				ID:           "scrum.assistant",
				PluginID:     "com.powerx.plugins.scrum",
				Name:         "Scrum 助理",
				Description:  "智能的 Scrum 管理助手，可以帮助创建任务、管理 Sprint、生成报告",
				Model:        "gpt-4",
				Instructions: "你是一个专业的 Scrum 管理助手。你可以帮助用户创建和管理任务、规划 Sprint、跟踪进度并生成报告。请始终以友好、专业的方式回应用户的请求。",
				DefaultTools: []string{
					"scrum.task.create",
					"scrum.task.query",
					"scrum.task.update",
					"scrum.sprint.create",
					"scrum.sprint.query",
					"scrum.report.generate",
				},
				RequiredPermissions: []string{"scrum:task:read", "scrum:sprint:read"},
			},
		},

		Tools: []contracts.ToolConfig{
			{
				ID:           "scrum.task.create",
				PluginID:     "com.powerx.plugins.scrum",
				Name:         "创建任务",
				Description:  "创建一个新的 Scrum 任务",
				Transport:    "http",
				Endpoint:     "/api/v1/tasks",
				Method:       "POST",
				RBACResource: "scrum:task",
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
				ID:           "scrum.task.query",
				PluginID:     "com.powerx.plugins.scrum",
				Name:         "查询任务",
				Description:  "查询 Scrum 任务列表",
				Transport:    "http",
				Endpoint:     "/api/v1/tasks",
				Method:       "GET",
				RBACResource: "scrum:task",
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
			{
				ID:           "scrum.sprint.create",
				PluginID:     "com.powerx.plugins.scrum",
				Name:         "创建 Sprint",
				Description:  "创建一个新的 Sprint",
				Transport:    "http",
				Endpoint:     "/api/v1/sprints",
				Method:       "POST",
				RBACResource: "scrum:sprint",
				InputSchema: &contracts.JSONSchema{
					Type: "object",
					Properties: map[string]*contracts.JSONSchemaProperty{
						"name": {
							Type:        "string",
							Description: "Sprint 名称",
						},
						"goal": {
							Type:        "string",
							Description: "Sprint 目标",
						},
						"start_date": {
							Type:        "string",
							Format:      "date-time",
							Description: "开始日期",
						},
						"end_date": {
							Type:        "string",
							Format:      "date-time",
							Description: "结束日期",
						},
						"capacity": {
							Type:        "integer",
							Minimum:     func() *float64 { v := 1.0; return &v }(),
							Description: "团队容量（故事点）",
						},
					},
					Required: []string{"name", "start_date", "end_date"},
				},
				Timeout: 30,
			},
		},

		Workflows: []contracts.WorkflowConfig{
			{
				ID:          "scrum.plan.generate",
				PluginID:    "com.powerx.plugins.scrum",
				Name:        "生成 Sprint 计划",
				Description: "基于历史数据和团队能力自动生成 Sprint 计划",
				Endpoint:    "/api/v1/workflows/plan-generate",
				InputSchema: &contracts.JSONSchema{
					Type: "object",
					Properties: map[string]*contracts.JSONSchemaProperty{
						"sprint_duration": {
							Type:        "integer",
							Description: "Sprint 持续天数",
							Default:     14,
						},
						"team_capacity": {
							Type:        "integer",
							Description: "团队总容量（故事点）",
						},
						"priority_tasks": {
							Type: "array",
							Items: &contracts.JSONSchema{
								Type: "integer",
							},
							Description: "优先任务ID列表",
						},
					},
					Required: []string{"team_capacity"},
				},
				RequiredPermissions: []string{"scrum:task:read", "scrum:sprint:create"},
			},
		},

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
				"default_sprint_duration": {
					Type:        "integer",
					Description: "默认 Sprint 持续天数",
					Default:     14,
					Minimum:     func() *float64 { v := 1.0; return &v }(),
					Maximum:     func() *float64 { v := 30.0; return &v }(),
				},
				"task_auto_close": {
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

	c.JSON(http.StatusOK, contracts.APIResponse{
		Success:   true,
		Data:      manifest,
		Timestamp: time.Now(),
	})
}

// GetRBACInfo 获取 RBAC 信息
func (h *AdminHandler) GetRBACInfo(c *gin.Context) {
	log := logger.HandlerLogger("admin").WithContext(c.Request.Context())

	rbacInfo := &contracts.RBACInfo{
		Resources: []contracts.Resource{
			{
				Name:        "scrum:task",
				Description: "Scrum 任务管理",
				Actions: []contracts.Action{
					{Name: "read", Description: "查看任务"},
					{Name: "create", Description: "创建任务"},
					{Name: "update", Description: "更新任务"},
					{Name: "delete", Description: "删除任务"},
				},
			},
			{
				Name:        "scrum:sprint",
				Description: "Sprint 管理",
				Actions: []contracts.Action{
					{Name: "read", Description: "查看 Sprint"},
					{Name: "create", Description: "创建 Sprint"},
					{Name: "update", Description: "更新 Sprint"},
					{Name: "delete", Description: "删除 Sprint"},
				},
			},
			{
				Name:        "scrum:report",
				Description: "Scrum 报告",
				Actions: []contracts.Action{
					{Name: "read", Description: "查看报告"},
				},
			},
		},
		Roles: []contracts.Role{
			{
				Name:        "scrum_master",
				Description: "Scrum Master 角色",
				Permissions: []string{
					"scrum:task:*",
					"scrum:sprint:*",
					"scrum:report:read",
				},
			},
			{
				Name:        "product_owner",
				Description: "Product Owner 角色",
				Permissions: []string{
					"scrum:task:read",
					"scrum:task:create",
					"scrum:task:update",
					"scrum:sprint:read",
					"scrum:report:read",
				},
			},
			{
				Name:        "developer",
				Description: "开发者角色",
				Permissions: []string{
					"scrum:task:read",
					"scrum:task:update",
					"scrum:sprint:read",
				},
			},
		},
		Permissions: []contracts.Permission{
			{Resource: "scrum:task", Action: "read"},
			{Resource: "scrum:task", Action: "create"},
			{Resource: "scrum:task", Action: "update"},
			{Resource: "scrum:task", Action: "delete"},
			{Resource: "scrum:sprint", Action: "read"},
			{Resource: "scrum:sprint", Action: "create"},
			{Resource: "scrum:sprint", Action: "update"},
			{Resource: "scrum:sprint", Action: "delete"},
			{Resource: "scrum:report", Action: "read"},
		},
	}

	log.Info("RBAC info requested")

	c.JSON(http.StatusOK, contracts.APIResponse{
		Success:   true,
		Data:      rbacInfo,
		Timestamp: time.Now(),
	})
}
