package api

import (
	"scrum-plugin/internal/api/http/admin"
	"scrum-plugin/internal/api/http/admin/iam"
	"scrum-plugin/internal/api/http/demo"
	"scrum-plugin/internal/api/http/health"
	"scrum-plugin/internal/api/http/sprint"
	"scrum-plugin/internal/api/http/task"
	powerxclient "scrum-plugin/internal/grpc/client"
	"scrum-plugin/internal/handlers"
	"scrum-plugin/internal/logger"
	"scrum-plugin/internal/services"
	iamhandler "scrum-plugin/internal/transport/http/admin/iam"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Registry API 注册器
type Registry struct {
	engine *gin.Engine
	pxc    *powerxclient.PowerX // PowerX gRPC 客户端
	db     *gorm.DB             // 数据库连接

	// API 路由组
	taskRoutes   *task.Routes
	adminRoutes  *admin.Routes
	healthRoutes *health.Routes
	sprintRoutes *sprint.Routes
	demoRoutes   *demo.Routes // 新的 demo routes

	// 新架构的 handlers
	memberHandler *iamhandler.MemberHandler
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
	r.taskRoutes = task.NewRoutes(r.db)
	r.adminRoutes = admin.NewRoutes(r.db)
	r.healthRoutes = health.NewRoutes(healthHandler) // 健康检查 handler 从外部传入
	r.sprintRoutes = sprint.NewRoutes(r.db)
	r.demoRoutes = demo.NewRoutes(r.pxc) // 创建新的 demo routes

	// 初始化新架构的 handlers
	r.initializeNewHandlers()

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

	// 注册管理端路由（无需 RBAC）
	r.adminRoutes.Register(v1)
	logger.Infof("Registered %s routes at %s", r.adminRoutes.GetName(), r.adminRoutes.GetPrefix())

	// 注册新架构的 IAM 路由
	adminGroup := r.engine.Group("/api/v1/admin")
	if r.memberHandler != nil {
		iam.RegisterIAMRoutes(adminGroup, r.memberHandler)
		logger.Info("Registered IAM routes at /api/v1/admin/iam")
	}

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
	if r.demoRoutes != nil {
		routes = append(routes, r.demoRoutes.GetName())
	}

	return routes
}

// initializeNewHandlers 初始化新架构的 handlers
func (r *Registry) initializeNewHandlers() {
	if r.pxc != nil && r.db != nil {
		// 初始化服务层
		memberService := services.NewMemberService(r.db)

		// 初始化 handler层
		r.memberHandler = iamhandler.NewMemberHandler(memberService, r.pxc)

		logger.Info("Initialized new architecture handlers")
	} else {
		logger.Warn("PowerX client or database not available, skipping new handlers initialization")
	}
}
