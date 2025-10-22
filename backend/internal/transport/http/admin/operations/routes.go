package operations

import (
	oprepo "github.com/ArtisanCloud/PowerXPlugin/internal/domain/repository/operations"
	opservice "github.com/ArtisanCloud/PowerXPlugin/internal/services/operations"
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

	if deps != nil && deps.DB != nil {
		slaRepo := oprepo.NewSLARepository(deps.DB)
		slaSvc := opservice.NewSLAService(slaRepo, deps.Config, deps.OperationsMetrics)
		slaHandler := NewSLAHandler(slaSvc)

		sla := operationsGroup.Group("/sla")
		{
			sla.GET("/profiles", slaHandler.GetProfiles)
			sla.POST("/profiles", slaHandler.UpsertProfile)
			sla.POST("/profiles/recompute", slaHandler.Recompute)
			sla.PATCH("/profiles/actuals", slaHandler.UpdateActuals)
		}
	} else {
		operationsGroup.Group("/sla")
	}
}
