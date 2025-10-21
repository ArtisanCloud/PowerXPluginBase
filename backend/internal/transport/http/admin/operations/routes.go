package operations

import (
	"github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers Operations admin endpoints under /admin/operations.
// Concrete handlers will be implemented in subsequent phases; this scaffold
// ensures routing and RBAC wiring exist for downstream work.
func RegisterRoutes(router *gin.RouterGroup, deps *app.Deps) {
	if router == nil {
		return
	}

	operationsGroup := router.Group("/operations")

	supportHandler := NewSupportHandler(deps)
	support := operationsGroup.Group("/support")
	{
		support.GET("/playbook", supportHandler.GetPlaybook)
		support.PUT("/playbook", supportHandler.UpdatePlaybook)
		support.POST("/channels/test", supportHandler.TestChannels)
		support.GET("/metrics", supportHandler.GetMetrics)
	}

	operationsGroup.Group("/incidents")
	operationsGroup.Group("/sla")
}
