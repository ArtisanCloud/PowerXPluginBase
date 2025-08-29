package router

import (
	"scrum-plugin/internal/api"
	"time"

	"scrum-plugin/internal/config"
	"scrum-plugin/internal/db"
	powerxclient "scrum-plugin/internal/grpc/client"
	"scrum-plugin/internal/handlers"
	"scrum-plugin/internal/logger"
	"scrum-plugin/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Router 路由器结构
type Router struct {
	engine        *gin.Engine
	cfg           *config.Config
	pxc           *powerxclient.PowerX    // PowerX gRPC 客户端
	healthHandler *handlers.HealthHandler // 基础健康检查 handler
}

// New 创建新的路由器
func New(cfg *config.Config) *Router {
	return &Router{
		cfg: cfg,
	}
}

// SetPowerXClient 设置 PowerX gRPC 客户端
func (r *Router) SetPowerXClient(pxc *powerxclient.PowerX) {
	r.pxc = pxc
}

// Setup 设置路由
func (r *Router) Setup() *gin.Engine {
	// 设置 Gin 模式
	if r.cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 初始化 handlers
	r.initializeHandlers()

	// 创建 Gin 引擎
	r.engine = gin.New()

	// 设置全局中间件
	r.setupGlobalMiddleware()

	// 设置路由
	r.setupRoutes()

	logger.Info("Router setup completed")
	return r.engine
}

// initializeHandlers 初始化基础 handlers
func (r *Router) initializeHandlers() {
	// 只初始化基础的 handler，其他 handler 由各自的 routes 负责
	r.healthHandler = handlers.NewHealthHandler()

	logger.Info("Basic handlers initialized successfully")
}

// setupGlobalMiddleware 设置全局中间件
func (r *Router) setupGlobalMiddleware() {
	// 恢复中间件
	r.engine.Use(middleware.Recovery())

	// 请求日志中间件
	r.engine.Use(middleware.RequestLogger())

	// 安全头部中间件
	r.engine.Use(middleware.SecurityHeaders())

	// CORS 中间件
	r.engine.Use(middleware.CORS())

	// 请求 ID 中间件
	r.engine.Use(middleware.RequestID())

	// 超时中间件（30秒）
	r.engine.Use(middleware.Timeout(30 * time.Second))

	// 速率限制中间件（每分钟最多 100 个请求）
	r.engine.Use(middleware.RateLimiter(100, time.Minute))

	// 健康检查中间件（在其他中间件之前）
	r.engine.Use(middleware.HealthCheck("/healthz"))
}

// setupRoutes 设置路由
func (r *Router) setupRoutes() {
	// 设置租户认证中间件
	v1 := r.engine.Group("/v1")
	v1.Use(middleware.TenantMiddleware(r.cfg))

	// 设置业务路由的 RBAC 中间件
	gApi := r.engine.Group("/api/v1")
	gApi.Use(middleware.TenantMiddleware(r.cfg))
	rbacConfig := middleware.NewRBACConfig()
	gApi.Use(middleware.RBACMiddleware(rbacConfig))

	// 使用 API 注册器注册所有路由
	apiRegistry := api.NewRegistry(r.engine, r.pxc, db.GetGlobalDB())
	apiRegistry.RegisterRoutes(r.healthHandler)

	// 记录已注册的路由
	routes := apiRegistry.GetRegisteredRoutes()
	logger.Infof("Successfully registered %d API route groups: %v", len(routes), routes)
}

// GetEngine 获取 Gin 引擎
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}

// RegisterCustomRoutes 注册自定义路由
func (r *Router) RegisterCustomRoutes(fn func(*gin.Engine)) {
	if r.engine != nil && fn != nil {
		fn(r.engine)
	}
}

// RegisterMiddleware 注册中间件
func (r *Router) RegisterMiddleware(middleware gin.HandlerFunc) {
	if r.engine != nil {
		r.engine.Use(middleware)
	}
}
