package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/powerx-plugins/scrum/internal/config"
	"github.com/powerx-plugins/scrum/internal/handlers"
	"github.com/powerx-plugins/scrum/internal/logger"
	"github.com/powerx-plugins/scrum/internal/middleware"
)

// Router 路由器结构
type Router struct {
	engine       *gin.Engine
	cfg          *config.Config
	adminHandler *handlers.AdminHandler
	taskHandler  *handlers.TaskHandler
	healthHandler *handlers.HealthHandler
}

// New 创建新的路由器
func New(cfg *config.Config, adminHandler *handlers.AdminHandler, taskHandler *handlers.TaskHandler, healthHandler *handlers.HealthHandler) *Router {
	return &Router{
		cfg:           cfg,
		adminHandler:  adminHandler,
		taskHandler:   taskHandler,
		healthHandler: healthHandler,
	}
}

// Setup 设置路由
func (r *Router) Setup() *gin.Engine {
	// 设置 Gin 模式
	if r.cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 创建 Gin 引擎
	r.engine = gin.New()

	// 设置全局中间件
	r.setupGlobalMiddleware()

	// 设置路由
	r.setupRoutes()

	logger.Info("Router setup completed")
	return r.engine
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
	// 健康检查路由（无需认证）
	r.engine.GET("/healthz", r.healthHandler.HealthCheck)
	r.engine.GET("/ping", r.healthHandler.Ping)

	// API v1 路由组
	v1 := r.engine.Group("/v1")
	{
		// 租户认证中间件
		v1.Use(middleware.TenantMiddleware(r.cfg))

		// 管理端路由（无需 RBAC）
		r.setupAdminRoutes(v1)

		// 业务路由（需要 RBAC）
		r.setupBusinessRoutes(v1)
	}
}

// setupAdminRoutes 设置管理端路由
func (r *Router) setupAdminRoutes(rg *gin.RouterGroup) {
	admin := rg.Group("/admin")
	{
		// 插件清单
		admin.GET("/manifest", r.adminHandler.GetManifest)
		
		// RBAC 信息
		admin.GET("/rbac", r.adminHandler.GetRBACInfo)
	}
}

// setupBusinessRoutes 设置业务路由
func (r *Router) setupBusinessRoutes(rg *gin.RouterGroup) {
	// RBAC 中间件
	rbacConfig := middleware.NewRBACConfig()
	rg.Use(middleware.RBACMiddleware(rbacConfig))

	// 任务路由
	r.setupTaskRoutes(rg)

	// Sprint 路由
	r.setupSprintRoutes(rg)

	// Agent 工具路由
	r.setupAgentRoutes(rg)

	// 工作流路由
	r.setupWorkflowRoutes(rg)
}

// setupTaskRoutes 设置任务路由
func (r *Router) setupTaskRoutes(rg *gin.RouterGroup) {
	tasks := rg.Group("/tasks")
	{
		// CRUD 操作
		tasks.POST("", r.taskHandler.CreateTask)
		tasks.GET("", r.taskHandler.ListTasks)
		tasks.GET("/:id", r.taskHandler.GetTask)
		tasks.PUT("/:id", r.taskHandler.UpdateTask)
		tasks.DELETE("/:id", r.taskHandler.DeleteTask)

		// 状态更新
		tasks.PATCH("/:id/status", r.taskHandler.UpdateTaskStatus)

		// TODO: 更多任务操作将在后续实现
		// 分配操作、标签操作、批量操作、统计报告等
	}
}

// setupSprintRoutes 设置 Sprint 路由
func (r *Router) setupSprintRoutes(rg *gin.RouterGroup) {
	// TODO: Sprint 路由将在后续实现
	_ = rg
}

// setupAgentRoutes 设置 Agent 工具路由
func (r *Router) setupAgentRoutes(rg *gin.RouterGroup) {
	// TODO: Agent 路由将在后续实现
	_ = rg
}

// setupWorkflowRoutes 设置工作流路由
func (r *Router) setupWorkflowRoutes(rg *gin.RouterGroup) {
	// TODO: 工作流路由将在后续实现
	_ = rg
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