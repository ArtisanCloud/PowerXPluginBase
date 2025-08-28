package health

import (
	"github.com/gin-gonic/gin"
	"scrum-plugin/internal/handlers"
)

// Routes Health 路由配置
type Routes struct {
	healthHandler *handlers.HealthHandler
}

// NewRoutes 创建 Health 路由
func NewRoutes(healthHandler *handlers.HealthHandler) *Routes {
	return &Routes{
		healthHandler: healthHandler,
	}
}

// Register 注册 Health 路由（无需认证）
func (r *Routes) Register(engine *gin.Engine) {
	// 健康检查路由（根级别，无需认证）
	engine.GET("/healthz", r.healthHandler.HealthCheck) // 健康检查
	engine.GET("/ping", r.healthHandler.Ping)           // 简单ping

	// TODO: 扩展健康检查
	// engine.GET("/readiness", r.healthHandler.ReadinessCheck) // 就绪检查
	// engine.GET("/liveness", r.healthHandler.LivenessCheck)   // 存活检查
}

// GetPrefix 获取路由前缀
func (r *Routes) GetPrefix() string {
	return ""
}

// GetName 获取路由名称
func (r *Routes) GetName() string {
	return "Health API"
}
