package router

import (
	"github.com/ArtisanCloud/PowerXPlugin/internal/config"
	"github.com/ArtisanCloud/PowerXPlugin/internal/logger"
	"github.com/ArtisanCloud/PowerXPlugin/internal/middleware"
	"github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"
	"github.com/ArtisanCloud/PowerXPlugin/internal/transport/http"
	middleware2 "github.com/ArtisanCloud/PowerXPlugin/internal/transport/http/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

// Router 路由器结构
type Router struct {
	engine *gin.Engine
	cfg    *config.Config
	deps   *app.Deps
}

// New 创建新的路由器
func NewRouter(cfg *config.Config, deps *app.Deps) *Router {
	return &Router{
		cfg:  cfg,
		deps: deps,
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
	// 设置租户认证中间件
	prefix := r.cfg.Server.APIPrefix
	if prefix == "" {
		prefix = "api/v1"
	}

	jwtCfg := r.buildJWT()
	rbacCfg := r.buildRBAC()
	abacClient := r.buildABAC()

	gApi := r.engine.Group(prefix)
	gApi.Use(middleware2.JWTAuth(jwtCfg))
	gApi.Use(middleware2.RBAC(rbacCfg, abacClient, func(m, path string) (bool, map[string]any) {
		// 标记需要 ABAC 的路由（可换成表驱动/装饰器）
		if path == "/api/v1/notes/:id" && m == "GET" {
			return true, map[string]any{"note_id": "{id}"}
		}
		return false, nil
	}))

	// 使用 API 注册器注册所有路由
	apiRegistry := http.NewRegistry(r.engine, r.deps)
	apiRegistry.RegisterAPIRoutes(gApi)

	// 记录已注册的路由表
	//apiRegistry.PrintRegisteredRoutes()
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

// —— 从配置构造 JWT 配置 —— //
func (r *Router) buildJWT() middleware.JWTAuthConfig {
	prod := r.cfg.IsProduction()
	cfg := middleware.JWTAuthConfig{
		Issuer:             "powerx",
		AcceptAudiences:    []string{"plugin:note"},
		HMACSecret:         "change-me",
		ClockSkewSeconds:   60,
		Optional:           !prod,
		AllowSignedContext: !prod,
		ContextHMACSecret:  "change-me",
		MaxCtxAgeSeconds:   300,
	}
	// TODO: 下面把你的真实配置字段映射进来（示例）：
	// if r.cfg.Auth != nil {
	//     cfg.Issuer = r.cfg.Auth.Issuer
	//     cfg.AcceptAudiences = r.cfg.Auth.Audiences
	//     cfg.HMACSecret = r.cfg.Auth.HMACSecret
	//     cfg.AllowSignedContext = r.cfg.Auth.AllowSignedContext
	//     cfg.ContextHMACSecret = r.cfg.Auth.ContextHMACSecret
	//     cfg.Optional = r.cfg.Auth.Optional && !prod
	// }
	return cfg
}

// —— 从配置构造 RBAC 配置 —— //
func (r *Router) buildRBAC() *middleware.RBACConfig {
	cfg := &middleware.RBACConfig{
		Enabled:         true,
		DefaultDeny:     true,
		SuperAdminRoles: []string{"superadmin"},
		RoutePermissions: map[string]middleware.Permission{
			"GET:/api/v1/notes/:id": {Resource: "note", Action: "read"},
		},
	}
	// TODO: 若你有 RBAC 目录/配置文件，在这里加载并覆盖 cfg.RoutePermissions
	return cfg
}

// —— 从配置构造 ABAC 客户端（可根据环境切换） —— //
func (r *Router) buildABAC() middleware.ABACClient {
	endpoint := "http://localhost:18080/check"
	if r.cfg.IsProduction() {
		endpoint = "http://pdp.powerx.svc/check"
	}
	// TODO: 如果 r.cfg 里有 PDP 地址，从配置读：
	// if r.cfg.Auth != nil && r.cfg.Auth.PDP != "" { endpoint = r.cfg.Auth.PDP }
	return middleware.NewHTTPABACClient(endpoint)
}
