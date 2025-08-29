package task

import (
	"scrum-plugin/internal/services"
	taskhandler "scrum-plugin/internal/transport/http/admin/task"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Routes Task 后台管理路由配置
type Routes struct {
	handler *taskhandler.Handler
}

// NewRoutes 创建 Task 后台管理路由
func NewRoutes(db *gorm.DB) *Routes {
	// 创建 task service
	taskService := services.NewTaskService(db)
	
	// 创建 handler 层
	handler := taskhandler.NewHandler(taskService)
	
	return &Routes{
		handler: handler,
	}
}

// Register 注册 Task 后台管理路由
func (r *Routes) Register(rg *gin.RouterGroup) {
	tasks := rg.Group("/tasks")
	{
		// 基础 CRUD 操作（后台管理版本）
		tasks.GET("", r.handler.ListTasks)         // 获取任务列表
		tasks.GET("/:id", r.handler.GetTask)       // 获取任务详情

		// 后台管理特有功能
		tasks.PATCH("/:id/status", r.handler.UpdateTaskStatus) // 更新任务状态
		tasks.POST("/batch/status", r.handler.BatchUpdateStatus) // 批量更新状态
		
		// 统计和分析
		tasks.GET("/stats", r.handler.GetTaskStats) // 获取任务统计
	}
}

// GetPrefix 获取路由前缀
func (r *Routes) GetPrefix() string {
	return "/admin"
}

// GetName 获取路由名称
func (r *Routes) GetName() string {
	return "Admin Task API"
}