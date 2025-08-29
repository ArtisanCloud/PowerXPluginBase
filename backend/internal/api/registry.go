package api

import (
	"scrum-plugin/internal/api/http/admin"
	"scrum-plugin/internal/api/http/demo"
	"scrum-plugin/internal/api/http/health"
	powerxclient "scrum-plugin/internal/grpc/client"
	"scrum-plugin/internal/handlers"
	"scrum-plugin/internal/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Registry API 注册器
type Registry struct {
	engine *gin.Engine
	pxc    *powerxclient.PowerX // PowerX gRPC 客户端
	db     *gorm.DB             // 数据库连接

	// API 路由组
	adminRoutes  *admin.Routes
	healthRoutes *health.Routes
	demoRoutes   *demo.Routes // 新的 demo routes
}

// NewRegistry 创建 API 注册器
func NewRegistry(engine *gin.Engine, pxc *powerxclient.PowerX, db *gorm.DB) *Registry {
	return &Registry{
		engine: engine,
		pxc:    pxc,
		db:     db,
	}
}

// RegisterRoutes 注册所有路由
func (r *Registry) RegisterRoutes(healthHandler *handlers.HealthHandler) {
	// 初始化各个 API 路由 - 让它们自己创建所需的 handlers
	r.adminRoutes = admin.NewRoutes(r.db, r.pxc)     // 传入 PowerX 客户端
	r.healthRoutes = health.NewRoutes(healthHandler) // 健康检查 handler 从外部传入
	r.demoRoutes = demo.NewRoutes(r.pxc)             // 创建新的 demo routes

	// 注册健康检查路由（无需认证和中间件）
	r.healthRoutes.Register(r.engine)
	logger.Infof("Registered %s routes", r.healthRoutes.GetName())

	// 注册 gRPC 演示路由（用于测试和开发）- 使用新的分层架构
	if r.pxc != nil {
		// 使用新的分层架构的 demo routes
		apiGroup := r.engine.Group("/api/v1")
		r.demoRoutes.Register(apiGroup)
		logger.Info("Registered new demo routes with proper architecture")
	}

	// 注册需要认证的 API 路由
	r.registerAuthenticatedRoutes()
}

// registerAuthenticatedRoutes 注册需要认证的 API 路由
func (r *Registry) registerAuthenticatedRoutes() {
	// API v1 路由组（需要租户认证）
	v1 := r.engine.Group("/v1")
	// 注意：租户中间件在 router 层统一设置

	// 注册后台管理路由（包含 IAM、Sprint、Task）
	r.adminRoutes.Register(v1)
	logger.Infof("Registered %s routes at %s", r.adminRoutes.GetName(), r.adminRoutes.GetPrefix())
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
	if r.demoRoutes != nil {
		routes = append(routes, r.demoRoutes.GetName())
	}

	return routes
}
