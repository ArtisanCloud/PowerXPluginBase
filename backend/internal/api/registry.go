package api

import (
	"github.com/gin-gonic/gin"
	"scrum-plugin/internal/api/http/admin"
	"scrum-plugin/internal/api/http/health"
	"scrum-plugin/internal/api/http/sprint"
	"scrum-plugin/internal/api/http/task"
	"scrum-plugin/internal/handlers"
	"scrum-plugin/internal/logger"
)

// Registry API 注册器
type Registry struct {
	engine *gin.Engine

	// API 路由组
	taskRoutes   *task.Routes
	adminRoutes  *admin.Routes
	healthRoutes *health.Routes
	sprintRoutes *sprint.Routes
}

// NewRegistry 创建 API 注册器
func NewRegistry(engine *gin.Engine) *Registry {
	return &Registry{
		engine: engine,
	}
}

// RegisterRoutes 注册所有路由
func (r *Registry) RegisterRoutes(
	taskHandler *handlers.TaskHandler,
	adminHandler *handlers.AdminHandler,
	healthHandler *handlers.HealthHandler,
) {
	// 初始化各个 API 路由
	r.taskRoutes = task.NewRoutes(taskHandler)
	r.adminRoutes = admin.NewRoutes(adminHandler)
	r.healthRoutes = health.NewRoutes(healthHandler)
	r.sprintRoutes = sprint.NewRoutes()

	// 注册健康检查路由（无需认证和中间件）
	r.healthRoutes.Register(r.engine)
	logger.Infof("Registered %s routes", r.healthRoutes.GetName())

	// 注册需要认证的 API 路由
	r.registerAuthenticatedRoutes()
}

// registerAuthenticatedRoutes 注册需要认证的 API 路由
func (r *Registry) registerAuthenticatedRoutes() {
	// API v1 路由组（需要租户认证）
	v1 := r.engine.Group("/v1")
	// 注意：租户中间件在 router 层统一设置

	// 注册管理端路由（无需 RBAC）
	r.adminRoutes.Register(v1)
	logger.Infof("Registered %s routes at %s", r.adminRoutes.GetName(), r.adminRoutes.GetPrefix())

	// 业务路由组（需要 RBAC）
	api := r.engine.Group("/api/v1")
	// 注意：RBAC 中间件在 router 层统一设置

	// 注册业务 API 路由
	r.taskRoutes.Register(api)
	logger.Infof("Registered %s routes at %s", r.taskRoutes.GetName(), r.taskRoutes.GetPrefix())

	r.sprintRoutes.Register(api)
	logger.Infof("Registered %s routes at %s", r.sprintRoutes.GetName(), r.sprintRoutes.GetPrefix())
}

// GetRegisteredRoutes 获取已注册的路由信息
func (r *Registry) GetRegisteredRoutes() []string {
	routes := []string{}

	if r.healthRoutes != nil {
		routes = append(routes, r.healthRoutes.GetName())
	}
	if r.adminRoutes != nil {
		routes = append(routes, r.adminRoutes.GetName())
	}
	if r.taskRoutes != nil {
		routes = append(routes, r.taskRoutes.GetName())
	}
	if r.sprintRoutes != nil {
		routes = append(routes, r.sprintRoutes.GetName())
	}

	return routes
}
