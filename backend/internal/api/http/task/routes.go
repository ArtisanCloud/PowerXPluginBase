package task

import (
	"github.com/gin-gonic/gin"
	"scrum-plugin/internal/handlers"
)

// Routes Task 路由配置
type Routes struct {
	taskHandler *handlers.TaskHandler
}

// NewRoutes 创建 Task 路由
func NewRoutes(taskHandler *handlers.TaskHandler) *Routes {
	return &Routes{
		taskHandler: taskHandler,
	}
}

// Register 注册 Task 路由
func (r *Routes) Register(rg *gin.RouterGroup) {
	tasks := rg.Group("/tasks")
	{
		// CRUD 操作
		tasks.POST("", r.taskHandler.CreateTask)       // 创建任务
		tasks.GET("", r.taskHandler.ListTasks)         // 获取任务列表
		tasks.GET("/:id", r.taskHandler.GetTask)       // 获取任务详情
		tasks.PUT("/:id", r.taskHandler.UpdateTask)    // 更新任务
		tasks.DELETE("/:id", r.taskHandler.DeleteTask) // 删除任务

		// 状态更新
		tasks.PATCH("/:id/status", r.taskHandler.UpdateTaskStatus) // 更新任务状态

		// TODO: 扩展功能
		// tasks.POST("/:id/assign", r.taskHandler.AssignTask)        // 分配任务
		// tasks.DELETE("/:id/assign", r.taskHandler.UnassignTask)    // 取消分配
		// tasks.POST("/:id/labels", r.taskHandler.AddLabels)         // 添加标签
		// tasks.DELETE("/:id/labels", r.taskHandler.RemoveLabels)    // 移除标签
		// tasks.POST("/batch/status", r.taskHandler.BatchUpdateStatus) // 批量更新状态
		// tasks.POST("/batch/delete", r.taskHandler.BatchDelete)     // 批量删除
	}
}

// GetPrefix 获取路由前缀
func (r *Routes) GetPrefix() string {
	return "/api/v1"
}

// GetName 获取路由名称
func (r *Routes) GetName() string {
	return "Task API"
}
