package sprint

import (
	"github.com/gin-gonic/gin"
)

// Routes Sprint 路由配置
type Routes struct {
	// TODO: 添加 sprintHandler 当实现了 Sprint 功能时
	// sprintHandler *handlers.SprintHandler
}

// NewRoutes 创建 Sprint 路由
func NewRoutes() *Routes {
	return &Routes{
		// sprintHandler: sprintHandler,
	}
}

// Register 注册 Sprint 路由
func (r *Routes) Register(rg *gin.RouterGroup) {
	sprints := rg.Group("/sprints")
	{
		// TODO: Sprint 路由将在后续实现
		// sprints.POST("", r.sprintHandler.CreateSprint)        // 创建 Sprint
		// sprints.GET("", r.sprintHandler.ListSprints)          // 获取 Sprint 列表
		// sprints.GET("/:id", r.sprintHandler.GetSprint)        // 获取 Sprint 详情
		// sprints.PUT("/:id", r.sprintHandler.UpdateSprint)     // 更新 Sprint
		// sprints.DELETE("/:id", r.sprintHandler.DeleteSprint)  // 删除 Sprint

		// Sprint 状态管理
		// sprints.POST("/:id/start", r.sprintHandler.StartSprint)    // 开始 Sprint
		// sprints.POST("/:id/complete", r.sprintHandler.CompleteSprint) // 完成 Sprint
		// sprints.GET("/:id/burndown", r.sprintHandler.GetBurndownChart) // 燃尽图

		_ = sprints // 避免未使用变量警告
	}
}

// GetPrefix 获取路由前缀
func (r *Routes) GetPrefix() string {
	return "/api/v1"
}

// GetName 获取路由名称
func (r *Routes) GetName() string {
	return "Sprint API"
}
