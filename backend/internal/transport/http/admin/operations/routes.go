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

	var incidentHandler *IncidentHandler
	var supportHandler *SupportHandler
	var incidentDispatcher opservice.IncidentDispatcher

	if deps != nil {
		incidentDispatcher = opservice.NewLoggingIncidentDispatcher(deps.RuntimeLogger(deps.Ctx, "operations.incident", nil))
	}

	if deps != nil && deps.DB != nil {
		incidentRepo := oprepo.NewIncidentRepository(deps.DB)
		incidentSvc := opservice.NewIncidentService(incidentRepo, deps.Config, deps.OperationsMetrics, incidentDispatcher)
		incidentHandler = NewIncidentHandler(incidentSvc)

		supportHandler = NewSupportHandler(deps)
	} else {
		supportHandler = NewSupportHandler(nil)
	}

	support := operationsGroup.Group("/support")
	{
		support.GET("/playbook", supportHandler.GetPlaybook)
		support.PUT("/playbook", supportHandler.UpdatePlaybook)
		support.POST("/channels/test", supportHandler.TestChannels)
		support.GET("/metrics", supportHandler.GetMetrics)
	}

	incidents := operationsGroup.Group("/incidents")
	{
		if incidentHandler != nil {
			incidents.POST("", incidentHandler.CreateIncident)
			incidents.GET("", incidentHandler.ListIncidents)
			incidents.GET("/:incidentId", incidentHandler.GetIncident)
			incidents.PATCH("/:incidentId", incidentHandler.UpdateIncident)
			incidents.POST("/:incidentId/timeline", incidentHandler.AppendTimeline)
		}
	}

	operationsGroup.Group("/sla")
}
