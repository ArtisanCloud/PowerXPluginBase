package admin

import (
	"scrum-plugin/internal/handlers"
	"scrum-plugin/internal/api/http/admin/iam"
	"scrum-plugin/internal/api/http/admin/sprint"
	"scrum-plugin/internal/api/http/admin/task"
	"scrum-plugin/internal/services"
	iamhandler "scrum-plugin/internal/transport/http/admin/iam"
	powerxclient "scrum-plugin/internal/grpc/client"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Routes Admin 路由配置
type Routes struct {
	adminHandler *handlers.AdminHandler
	db           *gorm.DB
	powerxClient *powerxclient.PowerX
	
	// 子模块路由
	iamRoutes    *iam.Routes
	sprintRoutes *sprint.Routes
	taskRoutes   *task.Routes
}

// NewRoutes 创建 Admin 路由
func NewRoutes(db *gorm.DB, powerxClient *powerxclient.PowerX) *Routes {
	// 创建 admin handler
	adminHandler := handlers.NewAdminHandler()
	
	// 创建子模块路由
	var iamRoutes *iam.Routes
	if powerxClient != nil {
		// 创建 IAM 相关的 handler
		memberService := services.NewMemberService(db)
		memberHandler := iamhandler.NewMemberHandler(memberService, powerxClient)
		iamRoutes = iam.NewRoutes(memberHandler)
	}
	
	sprintRoutes := sprint.NewRoutes(db)
	taskRoutes := task.NewRoutes(db)

	return &Routes{
		adminHandler: adminHandler,
		db:           db,
		powerxClient: powerxClient,
		iamRoutes:    iamRoutes,
		sprintRoutes: sprintRoutes,
		taskRoutes:   taskRoutes,
	}
}

// Register 注册 Admin 路由
func (r *Routes) Register(rg *gin.RouterGroup) {
	admin := rg.Group("/admin")
	{
		// 基础管理功能
		admin.GET("/manifest", r.adminHandler.GetManifest) // 获取插件清单
		admin.GET("/rbac", r.adminHandler.GetRBACInfo)     // 获取权限信息
		
		// 注册子模块路由
		if r.iamRoutes != nil {
			r.iamRoutes.Register(admin)
		}
		r.sprintRoutes.Register(admin)
		r.taskRoutes.Register(admin)
		
		// 后台管理总览
		admin.GET("/overview", r.getOverview)
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

// getOverview 获取后台管理总览
func (r *Routes) getOverview(c *gin.Context) {
	c.JSON(200, gin.H{
		"title":   "PowerX Scrum Plugin - 后台管理",
		"version": "1.0.0",
		"modules": map[string]interface{}{
			"iam": gin.H{
				"name":        "IAM 管理",
				"description": "成员和团队管理",
				"endpoints":   []string{"/admin/iam/members", "/admin/iam/teams"},
				"available":   r.iamRoutes != nil,
			},
			"sprint": gin.H{
				"name":        "Sprint 管理",
				"description": "Sprint 生命周期管理",
				"endpoints":   []string{"/admin/sprints"},
				"available":   true,
			},
			"task": gin.H{
				"name":        "任务管理",
				"description": "任务后台管理和统计",
				"endpoints":   []string{"/admin/tasks", "/admin/tasks/stats"},
				"available":   true,
			},
		},
		"quick_links": []gin.H{
			{"name": "任务统计", "url": "/admin/tasks/stats"},
			{"name": "成员管理", "url": "/admin/iam/members"},
			{"name": "Sprint 管理", "url": "/admin/sprints"},
			{"name": "系统信息", "url": "/admin/manifest"},
		},
		"note": "使用统一的后台管理界面管理 Scrum 插件的所有功能",
	})
}
