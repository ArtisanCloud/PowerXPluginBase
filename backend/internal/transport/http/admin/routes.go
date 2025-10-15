package admin

import (
	"github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"
	adminruntime "github.com/ArtisanCloud/PowerXPlugin/internal/transport/http/admin/runtime_ops"
	"github.com/gin-gonic/gin"
)

// Register 注册 Admin 路由
func RegisterAPIRoutes(rg *gin.RouterGroup, deps *app.Deps) {
	adminHandler := NewAdminHandler(deps)
	admin := rg.Group("/admin")
	{
		// 基础管理功能
		admin.GET("/manifest", adminHandler.GetManifest) // 获取插件清单
		admin.GET("/rbac", adminHandler.GetRBACInfo)     // 获取权限信息

		runtimeOps := admin.Group("/runtime")
		adminruntime.RegisterRoutes(runtimeOps)

	}
}
