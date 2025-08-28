package admin

import (
	"github.com/gin-gonic/gin"
	"scrum-plugin/internal/handlers"
)

// Routes Admin 路由配置
type Routes struct {
	adminHandler *handlers.AdminHandler
}

// NewRoutes 创建 Admin 路由
func NewRoutes(adminHandler *handlers.AdminHandler) *Routes {
	return &Routes{
		adminHandler: adminHandler,
	}
}

// Register 注册 Admin 路由
func (r *Routes) Register(rg *gin.RouterGroup) {
	admin := rg.Group("/admin")
	{
		// 插件管理
		admin.GET("/manifest", r.adminHandler.GetManifest) // 获取插件清单
		admin.GET("/rbac", r.adminHandler.GetRBACInfo)     // 获取权限信息

		// TODO: 扩展管理功能
		// admin.GET("/status", r.adminHandler.GetStatus)       // 获取系统状态
		// admin.GET("/metrics", r.adminHandler.GetMetrics)     // 获取指标
		// admin.GET("/config", r.adminHandler.GetConfig)       // 获取配置
		// admin.PUT("/config", r.adminHandler.UpdateConfig)    // 更新配置
	}
}

// GetPrefix 获取路由前缀
func (r *Routes) GetPrefix() string {
	return "/v1"
}

// GetName 获取路由名称
func (r *Routes) GetName() string {
	return "Admin API"
}
