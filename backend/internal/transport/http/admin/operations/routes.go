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

	operations := router.Group("/operations")

	// Placeholder groups for upcoming phases (support, incidents, SLA).
	operations.Group("/support")
	operations.Group("/incidents")
	operations.Group("/sla")
}
