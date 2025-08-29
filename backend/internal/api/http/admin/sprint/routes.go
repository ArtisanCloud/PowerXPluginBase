package sprint

import (
	sprinthandler "scrum-plugin/internal/transport/http/admin/sprint"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Routes Sprint 后台管理路由配置
type Routes struct {
	handler *sprinthandler.Handler
}

// NewRoutes 创建 Sprint 后台管理路由
func NewRoutes(db *gorm.DB) *Routes {
	// TODO: 创建 sprint service
	// sprintService := services.NewSprintService(db)

	// 创建 handler 层
	// handler := sprinthandler.NewHandler(sprintService)

	// 临时使用 nil，等待 SprintService 实现
	var handler *sprinthandler.Handler = nil

	return &Routes{
		handler: handler,
	}
}

// Register 注册 Sprint 后台管理路由
func (r *Routes) Register(rg *gin.RouterGroup) {
	sprints := rg.Group("/sprints")
	{
		if r.handler != nil {
			// 基础 CRUD 操作
			sprints.POST("", r.handler.CreateSprint)       // 创建 Sprint
			sprints.GET("", r.handler.ListSprints)         // 获取 Sprint 列表
			sprints.GET("/:id", r.handler.GetSprint)       // 获取 Sprint 详情
			sprints.PUT("/:id", r.handler.UpdateSprint)    // 更新 Sprint
			sprints.DELETE("/:id", r.handler.DeleteSprint) // 删除 Sprint
		} else {
			// 临时占位符路由
			sprints.GET("", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Sprint management routes will be implemented here",
					"status":  "placeholder",
					"module":  "admin.sprint",
					"endpoints": []string{
						"POST /admin/sprints - 创建 Sprint",
						"GET /admin/sprints - 获取 Sprint 列表",
						"GET /admin/sprints/:id - 获取 Sprint 详情",
						"PUT /admin/sprints/:id - 更新 Sprint",
						"DELETE /admin/sprints/:id - 删除 Sprint",
					},
				})
			})

			sprints.GET("/:id", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Sprint details will be implemented here",
					"status":  "placeholder",
					"id":      c.Param("id"),
				})
			})
		}
	}
}

// GetPrefix 获取路由前缀
func (r *Routes) GetPrefix() string {
	return "/admin"
}

// GetName 获取路由名称
func (r *Routes) GetName() string {
	return "Admin Sprint API"
}
